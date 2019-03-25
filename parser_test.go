package confl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockScanner outputs tokens directly
type mockScanner struct {
	tokens       []Token
	currentToken int
}

// Token returns the next token from the mock scanner
func (s *mockScanner) Token() Token {
	cur := s.currentToken
	s.currentToken++
	return s.tokens[cur]
}

func TestParser(t *testing.T) {
	tests := []struct {
		name   string
		tokens []Token
		doc    *Map
	}{

		{
			"implicit document map",
			[]Token{
				Token{Type: WordToken, Content: "test"},
				Token{Type: MapKVDelimToken},
				Token{Type: NumberToken, Content: "23"},
				Token{Type: EOFToken},
			},
			&Map{
				children: []Node{
					&Word{val: "test"},
					&Number{val: "23"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scan := &mockScanner{tokens: test.tokens}
			doc, err := Parse(scan)
			assert.Nil(t, err)
			assert.Equal(t, test.doc, doc)
		})
	}
}
