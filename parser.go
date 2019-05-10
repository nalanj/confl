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
	return parseMap(scan, eofToken, "")
}

// parseMap parses a map
func parseMap(
	scan *scanner,
	endDelim tokenType,
	decorator string,
) (*mapNode, error) {
	aMap := &mapNode{children: []Node{}, decorator: decorator}
	keys := make(map[string]struct{})

	for {
		// scan the key
		keyStart := scan.nextOffset
		keyNode, keyErr := parseValue(scan, true, endDelim, "")
		keyEnd := scan.nextOffset
		if keyErr != nil {
			return nil, keyErr
		}
		if keyNode == nil {
			return aMap, nil
		}
		if _, ok := keys[keyNode.Value()]; ok {
			return nil, newParseError(
				fmt.Sprintf("Duplicate key %s", keyNode.Value()),
				scan,
				keyStart,
				keyEnd-keyStart,
			)
		}

		// read the delimiter
		delimToken := scan.Token()
		if delimToken.Type != mapKVDelimToken {
			return nil, newParseError(
				"Illegal token, expected map delimiter `=`",
				scan,
				delimToken.Offset,
				len(delimToken.Content),
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
				len(scan.src),
				0,
			)
		}

		aMap.children = append(aMap.children, keyNode, valNode)
		keys[keyNode.Value()] = struct{}{}
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
	case token.Type == mapEndToken ||
		token.Type == listEndToken ||
		token.Type == decoratorEndToken ||
		token.Type == eofToken:
		return nil, newParseError(
			fmt.Sprintf(
				"Illegal closing token: got %s, expected %s",
				token.Type,
				closeType,
			),
			scan,
			token.Offset,
			len(token.Content),
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
				token.Offset,
				len(token.Content),
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
				token.Offset,
				len(token.Content),
			)
		}
		return parseMap(scan, mapEndToken, decorator)
	case token.Type == listStartToken && !mapKey:
		if mapKey {
			return nil, newParseError(
				"Lists aren't allowed as map keys",
				scan,
				token.Offset,
				len(token.Content),
			)
		}

		return parseList(scan, decorator)
	default:
		return nil, newParseError(
			"Illegal token",
			scan,
			token.Offset,
			len(token.Content),
		)
	}
}
