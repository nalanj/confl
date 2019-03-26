package confl

// Parse parses a document and returns an AST
func Parse(scan Scanner) (*Map, error) {
	docMap := &Map{children: []Node{}}

	if err := parseMap(scan, docMap); err != nil {
		return nil, err
	}

	return docMap, nil
}

// parseMap parses a map
func parseMap(scan Scanner, aMap *Map) error {
	for {

		// scan the key
		keyNode, keyErr := parseValue(scan, true)
		if keyErr != nil {
			return keyErr
		}
		if keyNode == nil {
			return nil
		}

		// read the delimiter
		delimToken := scan.Token()
		if delimToken.Type != MapKVDelimToken {
			return &parseError{
				msg:    "Illegal token, expected map delimiter `=`",
				offset: delimToken.Offset,
			}
		}

		// read and append the value
		valNode, valErr := parseValue(scan, false)
		if valErr != nil {
			return valErr
		}
		if valNode == nil {
			return &parseError{
				msg:    "Illegal token, expected map value, got EOF",
				offset: -1,
			}
		}

		aMap.children = append(aMap.children, keyNode, valNode)
	}
}

// parseValue parses and returns a node for a value type, or an error if no
// value type could be parsed. If the mapKey param is true then only those
// types that are valid for a map key are allowed
func parseValue(scan Scanner, mapKey bool) (Node, error) {

	// read the value
	token := scan.Token()

	switch {
	case token.Type == EOFToken:
		return nil, nil
	case token.Type == WordToken:
		return &ValueNode{nodeType: WordType, val: token.Content}, nil
	case token.Type == StringToken:
		return &ValueNode{nodeType: StringType, val: token.Content}, nil
	case token.Type == NumberToken && !mapKey:
		return &ValueNode{nodeType: NumberType, val: token.Content}, nil
	default:
		return nil, &parseError{
			msg:    "Illegal token",
			offset: token.Offset,
		}
	}
}
