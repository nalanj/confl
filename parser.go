package confl

import (
	"fmt"
	"io"
	"io/ioutil"
)

// Parse scans and parses from a reader
func Parse(r io.Reader) (Node, error) {
	src, readErr := ioutil.ReadAll(r)
	if readErr != nil {
		return nil, readErr
	}

	scan := newScanner(src)
	return parseScanner(scan)
}

// parseScanner parses a document and returns an AST
func parseScanner(scan *scanner) (Node, error) {
	endDelim := eofToken
	peekedTokens := scan.Peek(1)

	if len(peekedTokens) != 1 {
		return nil, newParseError("Empty document", scan, &token{})
	}

	startToken := peekedTokens[0]
	if startToken.Type == eofToken {
		return &mapNode{children: []Node{}}, nil
	}

	if startToken.Type == mapStartToken {
		endDelim = mapEndToken
		scan.Token()
	}

	return parseMap(scan, endDelim, "")
}

// parseMap parses a map
func parseMap(scan *scanner, endDelim tokenType, decorator string) (*mapNode, error) {
	aMap := &mapNode{children: []Node{}, decorator: decorator}

	for {
		// scan the key
		keyNode, keyErr := parseValue(scan, true, endDelim, "")
		if keyErr != nil {
			return nil, keyErr
		}
		if keyNode == nil {
			return aMap, nil
		}

		// read the delimiter
		delimToken := scan.Token()
		if delimToken.Type != mapKVDelimToken {
			return nil, newParseError(
				"Illegal token, expected map delimiter `=`",
				scan,
				delimToken,
			)
		}

		// read and append the value
		valNode, valErr := parseValue(scan, false, endDelim, "")
		if valErr != nil {
			return nil, valErr
		}
		if valNode == nil {
			return nil, newParseError(
				"Illegal token, expected map value, got EOF",
				scan,
				&token{Content: "", Offset: len(scan.src)},
			)
		}

		aMap.children = append(aMap.children, keyNode, valNode)
	}
}

// parseList parses and returns a list
func parseList(scan *scanner, decorator string) (*listNode, error) {
	list := &listNode{children: []Node{}, decorator: decorator}
	for {
		// scan the next value
		node, err := parseValue(scan, false, listEndToken, "")
		if err != nil {
			return nil, err
		}
		if node == nil {
			return list, nil
		}

		list.children = append(list.children, node)
	}
}

// parseDecoratorContents parses the node in a decorator
func parseDecoratorContents(
	scan *scanner,
	mapKey bool,
	decorator string,
) (Node, error) {

	node, err := parseValue(scan, mapKey, decoratorEndToken, decorator)
	if err != nil {
		return nil, err
	}

	// eat the closing decorator delimiter
	scan.Token()

	return node, nil
}

// parseValue parses and returns a node for a value type, or an error if no
// value type could be parsed. If the mapKey param is true then only those
// types that are valid for a map key are allowed
func parseValue(scan *scanner, mapKey bool, closeType tokenType, decorator string) (Node, error) {

	// read the value
	token := scan.Token()

	switch {
	case token.Type == closeType:
		return nil, nil
	case token.Type == mapEndToken || token.Type == listEndToken || token.Type == decoratorEndToken || token.Type == eofToken:
		return nil, newParseError(
			fmt.Sprintf(
				"Illegal closing token: got %s, expected %s",
				token.Type,
				closeType,
			),
			scan,
			token,
		)
	case token.Type == wordToken:
		return &valueNode{
			nodeType:  WordType,
			val:       token.Content,
			decorator: decorator,
		}, nil
	case token.Type == stringToken:
		return &valueNode{
			nodeType:  StringType,
			val:       token.Content,
			decorator: decorator,
		}, nil
	case token.Type == decoratorStartToken:
		return parseDecoratorContents(scan, mapKey, token.Content)
	case token.Type == numberToken:
		if mapKey {
			return nil, newParseError(
				"Numbers aren't allowed as map keys",
				scan,
				token,
			)
		}

		return &valueNode{
			nodeType:  NumberType,
			val:       token.Content,
			decorator: decorator,
		}, nil
	case token.Type == mapStartToken:
		if mapKey {
			return nil, newParseError(
				"Maps aren't allowed as map keys",
				scan,
				token,
			)
		}
		return parseMap(scan, mapEndToken, decorator)
	case token.Type == listStartToken && !mapKey:
		if mapKey {
			return nil, newParseError(
				"Lists aren't allowed as map keys",
				scan,
				token,
			)
		}

		return parseList(scan, decorator)
	default:
		return nil, newParseError(
			"Illegal token",
			scan,
			token,
		)
	}
}
