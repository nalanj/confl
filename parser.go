package confl

// peekScanner lets us peek one element ahead
type peekScanner struct {
	scanner Scanner
	peeked  *Token
}

// Token returns the peeked token if it's there, or the next token
func (p *peekScanner) Token() *Token {
	token := p.peeked
	p.peeked = nil

	if token == nil {
		token = p.scanner.Token()
	}
	return token
}

// peek looks ahead one token space. Subsequent calls will return the same
// peeked token until Token is called
func (p *peekScanner) peek() *Token {
	token := p.peeked
	if token == nil {
		p.peeked = p.scanner.Token()
		token = p.peeked
	}

	return token
}

// Parse parses a document and returns an AST
func Parse(scan Scanner) (*Map, error) {
	peek := &peekScanner{scanner: scan}

	startToken := peek.peek()
	if startToken.Type == MapStartToken {
		peek.Token()
	}

	return parseMap(peek)
}

// parseMap parses a map
func parseMap(scan Scanner) (*Map, error) {
	aMap := &Map{children: []Node{}}

	for {
		// scan the key
		keyNode, keyErr := parseValue(scan, true, MapEndToken)
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
		valNode, valErr := parseValue(scan, false, MapEndToken)
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
	case token.Type == EOFToken || token.Type == closeType:
		return nil, nil
	case token.Type == WordToken:
		return &ValueNode{nodeType: WordType, val: token.Content}, nil
	case token.Type == StringToken:
		return &ValueNode{nodeType: StringType, val: token.Content}, nil
	case token.Type == NumberToken && !mapKey:
		return &ValueNode{nodeType: NumberType, val: token.Content}, nil
	case token.Type == MapStartToken && !mapKey:
		return parseMap(scan)
	case token.Type == ListStartToken && !mapKey:
		return parseList(scan)
	default:
		return nil, &parseError{
			msg:    "Illegal token",
			offset: token.Offset,
		}
	}
}
