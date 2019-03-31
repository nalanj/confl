package confl

import (
	"io"
	"io/ioutil"
)

// Parse scans and parses from a reader
func Parse(r io.Reader) (Node, error) {
	src, readErr := ioutil.ReadAll(r)
	if readErr != nil {
		return nil, readErr
	}

	scan := NewScanner(src)
	return parseScanner(scan)
}

// parseScanner parses a document and returns an AST
func parseScanner(scan *Scanner) (Node, error) {
	endDelim := EOFToken
	peekedTokens := scan.Peek(1)

	if len(peekedTokens) != 1 {
		return nil, newParseError("Empty document", scan, 0)
	}

	startToken := peekedTokens[0]
	if startToken.Type == EOFToken {
		return &mapNode{children: []Node{}}, nil
	}

	if startToken.Type == MapStartToken {
		endDelim = MapEndToken
		scan.Token()
	}

	return parseMap(scan, endDelim, "")
}

// parseMap parses a map
func parseMap(scan *Scanner, endDelim TokenType, decorator string) (*mapNode, error) {
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
		if delimToken.Type != MapKVDelimToken {
			return nil, newParseError(
				"Illegal token, expected map delimiter `=`",
				scan,
				delimToken.Offset,
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
			)
		}

		aMap.children = append(aMap.children, keyNode, valNode)
	}
}

// parseList parses and returns a list
func parseList(scan *Scanner, decorator string) (*listNode, error) {
	list := &listNode{children: []Node{}, decorator: decorator}
	for {
		// scan the next value
		node, err := parseValue(scan, false, ListEndToken, "")
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
	scan *Scanner,
	mapKey bool,
	decorator string,
) (Node, error) {

	node, err := parseValue(scan, mapKey, DecoratorEndToken, decorator)
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
func parseValue(scan *Scanner, mapKey bool, closeType TokenType, decorator string) (Node, error) {

	// read the value
	token := scan.Token()

	switch {
	case token.Type == closeType:
		return nil, nil
	case token.Type == WordToken:
		return &ValueNode{
			nodeType:  WordType,
			val:       token.Content,
			decorator: decorator,
		}, nil
	case token.Type == StringToken:
		return &ValueNode{
			nodeType:  StringType,
			val:       token.Content,
			decorator: decorator,
		}, nil
	case token.Type == DecoratorStartToken:
		return parseDecoratorContents(scan, mapKey, token.Content)
	case token.Type == NumberToken && !mapKey:
		return &ValueNode{
			nodeType:  NumberType,
			val:       token.Content,
			decorator: decorator,
		}, nil
	case token.Type == MapStartToken && !mapKey:
		return parseMap(scan, MapEndToken, decorator)
	case token.Type == ListStartToken && !mapKey:
		return parseList(scan, decorator)
	default:
		return nil, newParseError(
			"Illegal token",
			scan,
			token.Offset,
		)
	}
}
