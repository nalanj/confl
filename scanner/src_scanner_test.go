package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScan(t *testing.T) {
	tests := []struct {
		name     string
		document []byte
		tokens   []TokenType
		contents []string
	}{
		{
			"simple integer number",
			[]byte("12"),
			[]TokenType{Number, EOF},
			[]string{"12", ""},
		},
		{
			"simple decimal number",
			[]byte("12.3"),
			[]TokenType{Number, EOF},
			[]string{"12.3", ""},
		},
		{
			"illegal: two decimal number",
			[]byte("1.2.3"),
			[]TokenType{Illegal},
			[]string{"1.2."},
		},
		{
			"word",
			[]byte(" testing "),
			[]TokenType{Word, EOF},
			[]string{"testing", ""},
		},
		{
			"empty map",
			[]byte("{}"),
			[]TokenType{MapStart, MapEnd, EOF},
			[]string{"", "", ""},
		},
		{
			"map",
			[]byte("{word=12}"),
			[]TokenType{MapStart, Word, MapKVDelim, Number, MapEnd, EOF},
			[]string{"", "word", "", "12", "", ""},
		},
		{
			"spacey map",
			[]byte("  {  word = 12\t}"),
			[]TokenType{MapStart, Word, MapKVDelim, Number, MapEnd, EOF},
			[]string{"", "word", "", "12", "", ""},
		},
		{
			"empty list",
			[]byte("[]"),
			[]TokenType{ListStart, ListEnd, EOF},
			[]string{"", "", ""},
		},
		{
			"list with a couple of items",
			[]byte("[word 1.2]"),
			[]TokenType{ListStart, Word, Number, ListEnd, EOF},
			[]string{"", "word", "1.2", "", ""},
		},
		{
			"comment",
			[]byte("word # comment\nword"),
			[]TokenType{Word, Comment, Word, EOF},
			[]string{"word", "# comment", "word", ""},
		},
		{
			"simple string with double quote",
			[]byte("\"a string\""),
			[]TokenType{String, EOF},
			[]string{"a string", ""},
		},
		{
			"simple string with single quote",
			[]byte("'a string'"),
			[]TokenType{String, EOF},
			[]string{"a string", ""},
		},
		{
			"string with escaped inner double quote",
			[]byte("\"a \\\" string\""),
			[]TokenType{String, EOF},
			[]string{"a \" string", ""},
		},
		{
			"string with escaped inner single quote",
			[]byte("'a \\' string'"),
			[]TokenType{String, EOF},
			[]string{"a ' string", ""},
		},
		{
			"string with line breaks",
			[]byte("'a \nstring'"),
			[]TokenType{String, EOF},
			[]string{"a \nstring", ""},
		},
		{
			"simple decorator",
			[]byte("decorator(12)"),
			[]TokenType{DecoratorStart, Number, DecoratorEnd, EOF},
			[]string{"decorator", "12", "", ""},
		},
		{
			"complex example",
			[]byte(`
			# Simple wifi configuration
			device(wifi0)={
				network="Pretty fly for a wifi"
				key="Some long wpa key"
				dhcp=true

				dns=["10.0.0.1" "10.0.0.2"]
				gateway="10.0.0.1"

				vpn={host="12.12.12.12" user=frank pass=secret key=path("/etc/vpn.key")}
			}	
			`),
			[]TokenType{
				Comment,
				DecoratorStart, Word, DecoratorEnd, MapKVDelim, MapStart,
				Word, MapKVDelim, String,
				Word, MapKVDelim, String,
				Word, MapKVDelim, Word,
				Word, MapKVDelim, ListStart, String, String, ListEnd,
				Word, MapKVDelim, String,
				Word, MapKVDelim,
				MapStart,
				Word, MapKVDelim, String,
				Word, MapKVDelim, Word,
				Word, MapKVDelim, Word,
				Word, MapKVDelim, DecoratorStart, String, DecoratorEnd,
				MapEnd,
				MapEnd, EOF,
			},
			[]string{
				"# Simple wifi configuration",
				"device", "wifi0", "", "", "",
				"network", "", "Pretty fly for a wifi",
				"key", "", "Some long wpa key",
				"dhcp", "", "true",
				"dns", "", "", "10.0.0.1", "10.0.0.2", "",
				"gateway", "", "10.0.0.1",
				"vpn", "",
				"",
				"host", "", "12.12.12.12",
				"user", "", "frank",
				"pass", "", "secret",
				"key", "", "path", "/etc/vpn.key", "",
				"",
				"", "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens := []TokenType{}
			contents := []string{}

			s := Init(test.document)

			for {
				token := s.Token()

				tokens = append(tokens, token.Type)
				contents = append(contents, token.Content)

				if token.Type == Illegal || token.Type == EOF {
					break
				}
			}

			assert.Equal(t, test.tokens, tokens)
			assert.Equal(t, test.contents, contents)
		})
	}
}
