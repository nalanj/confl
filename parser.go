package confl

// peekScanner lets us peek one element ahead
type peekScanner struct {
	scanner Scanner
	peeked  []*Token
}

// Token returns the peeked token if it's there, or the next token
func (p *peekScanner) Token() *Token {
	var token *Token

	if len(p.peeked) > 0 {
		token = p.peeked[0]
		p.peeked = p.peeked[1:]
	}

	if token == nil {
		token = p.scanner.Token()
	}

	return token
}

// peek peeks from the scanner the number of tokens. If an EOF is encountered
// peeking stops.
func (p *peekScanner) peek(count int) []*Token {
	for i := 0; i < count; i++ {
		var token *Token

		if len(p.peeked) > i {
			token = p.peeked[i]
		} else {
			token = p.scanner.Token()
			p.peeked = append(p.peeked, token)
		}

		if token.Type == EOFToken {
			break
		}
	}

	return p.peeked[0:count]
}

// Parse parses a document and returns an AST
func Parse(scan Scanner) (*Map, error) {
	peek := &peekScanner{scanner: scan}

	endDelim := EOFToken
	peekedTokens := peek.peek(1)

	if len(peekedTokens) != 1 {
		return nil, &parseError{
			msg:    "Empty document",
			offset: 0,
		}
	}

	startToken := peekedTokens[0]
	if startToken.Type == EOFToken {
		return &Map{children: []Node{}}, nil
	}

	if startToken.Type == MapStartToken {
		endDelim = MapEndToken
		peek.Token()
	}

	return parseMap(peek, endDelim)
}

// parseMap parses a map
func parseMap(scan Scanner, endDelim TokenType) (*Map, error) {
	aMap := &Map{children: []Node{}}

	for {
		// scan the key
		keyNode, keyErr := parseValue(scan, true, endDelim)
		if keyErr != nil {
			return nil, keyErr
		}
		if keyNode == nil {
			return aMap, nil
		}

		// read the delimiter
		delimToken := scan.Token()
		if delimToken.Type != MapKVDelimToken {
			return nil, &parseError{
				msg:    "Illegal token, expected map delimiter `=`",
				offset: delimToken.Offset,
			}
		}

		// read and append the value
		valNode, valErr := parseValue(scan, false, endDelim)
		if valErr != nil {
			return nil, valErr
		}
		if valNode == nil {
			return nil, &parseError{
				msg:    "Illegal token, expected map value, got EOF",
				offset: -1,
			}
		}

		aMap.children = append(aMap.children, keyNode, valNode)
	}
}

// parseList parses and returns a list
func parseList(scan Scanner) (*List, error) {
	list := &List{children: []Node{}}
	for {
		// scan the next value
		node, err := parseValue(scan, true, ListEndToken)
		if err != nil {
			return nil, err
		}
		if node == nil {
			return list, nil
		}

		list.children = append(list.children, node)
	}
}

// parseValue parses and returns a node for a value type, or an error if no
// value type could be parsed. If the mapKey param is true then only those
// types that are valid for a map key are allowed
func parseValue(scan Scanner, mapKey bool, closeType TokenType) (Node, error) {

	// read the value
	token := scan.Token()

	switch {
	case token.Type == closeType:
		return nil, nil
	case token.Type == WordToken:
		return &ValueNode{nodeType: WordType, val: token.Content}, nil
	case token.Type == StringToken:
		return &ValueNode{nodeType: StringType, val: token.Content}, nil
	// case token.Type == DecoratorStartToken:
	// 	return parseDecorator(scan, mapKey)
	case token.Type == NumberToken && !mapKey:
		return &ValueNode{nodeType: NumberType, val: token.Content}, nil
	case token.Type == MapStartToken && !mapKey:
		return parseMap(scan, MapEndToken)
	case token.Type == ListStartToken && !mapKey:
		return parseList(scan)
	default:
		return nil, &parseError{
			msg:    "Illegal token",
			offset: token.Offset,
		}
	}
}
