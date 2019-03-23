package confl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScan(t *testing.T) {
	tests := []struct {
		name     string
		document []byte
		tokens   []Token
		values   []string
	}{
		{
			"simple integer number",
			[]byte("12"),
			[]Token{Number, EOF},
			[]string{"12", ""},
		},
		{
			"simple decimal number",
			[]byte("12.3"),
			[]Token{Number, EOF},
			[]string{"12.3", ""},
		},
		{
			"illegal: two decimal number",
			[]byte("1.2.3"),
			[]Token{Illegal},
			[]string{"1.2."},
		},
		{
			"word",
			[]byte(" testing "),
			[]Token{Word, EOF},
			[]string{"testing", ""},
		},
		{
			"empty map",
			[]byte("{}"),
			[]Token{MapStart, MapEnd, EOF},
			[]string{"", "", ""},
		},
		{
			"map",
			[]byte("{word=12}"),
			[]Token{MapStart, Word, MapKVDelim, Number, MapEnd, EOF},
			[]string{"", "word", "", "12", "", ""},
		},
		{
			"spacey map",
			[]byte("  {  word = 12\t}"),
			[]Token{MapStart, Word, MapKVDelim, Number, MapEnd, EOF},
			[]string{"", "word", "", "12", "", ""},
		},
		{
			"empty list",
			[]byte("[]"),
			[]Token{ListStart, ListEnd, EOF},
			[]string{"", "", ""},
		},
		{
			"list with a couple of items",
			[]byte("[word 1.2]"),
			[]Token{ListStart, Word, Number, ListEnd, EOF},
			[]string{"", "word", "1.2", "", ""},
		},
		{
			"comment",
			[]byte("word # comment\nword"),
			[]Token{Word, Comment, Word, EOF},
			[]string{"word", "# comment", "word", ""},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens := []Token{}
			values := []string{}

			s := Init(test.document)

			for {
				token, _, value := s.Token()

				tokens = append(tokens, token)
				values = append(values, value)

				if token == Illegal || token == EOF {
					break
				}
			}

			assert.Equal(t, test.tokens, tokens)
			assert.Equal(t, test.values, values)
		})
	}
}
