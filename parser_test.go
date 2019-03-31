package confl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockScanner outputs tokens directly
type mockScanner struct {
	tokens  []*Token
	current int
}

// Token returns the next token from the mock scanner
func (s *mockScanner) Token() *Token {
	cur := s.current
	s.current++
	return s.tokens[cur]
}

// Peek returns the next count tokens, up to an eof
func (s *mockScanner) Peek(count int) []*Token {
	end := s.current + count
	if end > len(s.tokens) {
		end = len(s.tokens)
	}

	return s.tokens[s.current:end]
}

func TestParser(t *testing.T) {
	tests := []struct {
		name   string
		tokens []*Token
		doc    *Map
		err    bool
	}{

		{
			"implicit document map",
			[]*Token{
				&Token{Type: WordToken, Content: "test"},
				&Token{Type: MapKVDelimToken},
				&Token{Type: NumberToken, Content: "23"},
				&Token{Type: StringToken, Content: "also"},
				&Token{Type: MapKVDelimToken},
				&Token{Type: WordToken, Content: "this"},
				&Token{Type: EOFToken},
			},
			&Map{
				children: []Node{
					&ValueNode{nodeType: WordType, val: "test"},
					&ValueNode{nodeType: NumberType, val: "23"},
					&ValueNode{nodeType: StringType, val: "also"},
					&ValueNode{nodeType: WordType, val: "this"},
				},
			},
			false,
		},

		{
			"implicit document map, illegal end token",
			[]*Token{
				&Token{Type: WordToken, Content: "test"},
				&Token{Type: MapKVDelimToken},
				&Token{Type: NumberToken, Content: "23"},
				&Token{Type: StringToken, Content: "also"},
				&Token{Type: MapKVDelimToken},
				&Token{Type: WordToken, Content: "this"},
				&Token{Type: MapEndToken},
				&Token{Type: EOFToken},
			},
			nil,
			true,
		},

		{
			"explicit document map",
			[]*Token{
				&Token{Type: MapStartToken},
				&Token{Type: WordToken, Content: "test"},
				&Token{Type: MapKVDelimToken},
				&Token{Type: NumberToken, Content: "23"},
				&Token{Type: StringToken, Content: "also"},
				&Token{Type: MapKVDelimToken},
				&Token{Type: WordToken, Content: "this"},
				&Token{Type: MapEndToken},
				&Token{Type: EOFToken},
			},
			&Map{
				children: []Node{
					&ValueNode{nodeType: WordType, val: "test"},
					&ValueNode{nodeType: NumberType, val: "23"},
					&ValueNode{nodeType: StringType, val: "also"},
					&ValueNode{nodeType: WordType, val: "this"},
				},
			},
			false,
		},

		{
			"explicit document map, illegal end token",
			[]*Token{
				&Token{Type: MapStartToken},
				&Token{Type: WordToken, Content: "test"},
				&Token{Type: MapKVDelimToken},
				&Token{Type: NumberToken, Content: "23"},
				&Token{Type: StringToken, Content: "also"},
				&Token{Type: MapKVDelimToken},
				&Token{Type: WordToken, Content: "this"},
				&Token{Type: EOFToken},
			},
			nil,
			true,
		},

		{
			"nested map",
			[]*Token{
				&Token{Type: WordToken, Content: "map"},
				&Token{Type: MapKVDelimToken},
				&Token{Type: MapStartToken},
				&Token{Type: WordToken, Content: "key"},
				&Token{Type: MapKVDelimToken},
				&Token{Type: WordToken, Content: "value"},
				&Token{Type: MapEndToken},
				&Token{Type: EOFToken},
			},
			&Map{
				children: []Node{
					&ValueNode{nodeType: WordType, val: "map"},
					&Map{
						children: []Node{
							&ValueNode{nodeType: WordType, val: "key"},
							&ValueNode{nodeType: WordType, val: "value"},
						},
					},
				},
			},
			false,
		},

		{
			"nested list",
			[]*Token{
				&Token{Type: WordToken, Content: "list"},
				&Token{Type: MapKVDelimToken},
				&Token{Type: ListStartToken},
				&Token{Type: WordToken, Content: "item1"},
				&Token{Type: WordToken, Content: "item2"},
				&Token{Type: ListEndToken},
				&Token{Type: EOFToken},
			},
			&Map{
				children: []Node{
					&ValueNode{nodeType: WordType, val: "list"},
					&List{
						children: []Node{
							&ValueNode{nodeType: WordType, val: "item1"},
							&ValueNode{nodeType: WordType, val: "item2"},
						},
					},
				},
			},
			false,
		},

		{
			"simple decorator",
			[]*Token{
				&Token{Type: DecoratorStartToken, Content: "dec"},
				&Token{Type: WordToken, Content: "test"},
				&Token{Type: DecoratorEndToken},
				&Token{Type: MapKVDelimToken},
				&Token{Type: NumberToken, Content: "23"},
				&Token{Type: StringToken, Content: "also"},
				&Token{Type: MapKVDelimToken},
				&Token{Type: WordToken, Content: "this"},
				&Token{Type: EOFToken},
			},
			&Map{
				children: []Node{
					&ValueNode{nodeType: WordType, val: "test", decorator: "dec"},
					&ValueNode{nodeType: NumberType, val: "23"},
					&ValueNode{nodeType: StringType, val: "also"},
					&ValueNode{nodeType: WordType, val: "this"},
				},
			},
			false,
		},

		{
			"decorator on list as map key errors",
			[]*Token{
				&Token{Type: DecoratorStartToken, Content: "dec"},
				&Token{Type: MapStartToken},
				&Token{Type: WordToken, Content: "key"},
				&Token{Type: MapKVDelimToken},
				&Token{Type: WordToken, Content: "val"},
				&Token{Type: MapEndToken},
				&Token{Type: MapKVDelimToken},
				&Token{Type: WordToken, Content: "val"},
				&Token{Type: EOFToken},
			},
			nil,
			true,
		},

		{
			"decorator with map",
			[]*Token{
				&Token{Type: WordToken, Content: "key"},
				&Token{Type: MapKVDelimToken},
				&Token{Type: DecoratorStartToken, Content: "dec"},
				&Token{Type: MapStartToken},
				&Token{Type: WordToken, Content: "decKey"},
				&Token{Type: MapKVDelimToken},
				&Token{Type: WordToken, Content: "val"},
				&Token{Type: MapEndToken},
				&Token{Type: DecoratorEndToken},
				&Token{Type: EOFToken},
			},
			&Map{
				children: []Node{
					&ValueNode{nodeType: WordType, val: "key"},
					&Map{
						children: []Node{
							&ValueNode{nodeType: WordType, val: "decKey"},
							&ValueNode{nodeType: WordType, val: "val"},
						},
						decorator: "dec",
					},
				},
			},
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scan := &mockScanner{tokens: test.tokens}
			doc, err := Parse(scan)
			assert.Equal(t, test.err, err != nil)
			assert.Equal(t, test.doc, doc)
		})
	}
}
