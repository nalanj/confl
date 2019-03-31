package confl

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
			[]TokenType{NumberToken, EOFToken},
			[]string{"12", ""},
		},
		{
			"simple decimal number",
			[]byte("12.3"),
			[]TokenType{NumberToken, EOFToken},
			[]string{"12.3", ""},
		},
		{
			"illegal: two decimal number",
			[]byte("1.2.3"),
			[]TokenType{IllegalToken},
			[]string{"1.2."},
		},
		{
			"word",
			[]byte(" testing "),
			[]TokenType{WordToken, EOFToken},
			[]string{"testing", ""},
		},
		{
			"empty map",
			[]byte("{}"),
			[]TokenType{MapStartToken, MapEndToken, EOFToken},
			[]string{"", "", ""},
		},
		{
			"map",
			[]byte("{word=12}"),
			[]TokenType{
				MapStartToken,
				WordToken, MapKVDelimToken, NumberToken,
				MapEndToken, EOFToken,
			},
			[]string{"", "word", "", "12", "", ""},
		},
		{
			"spacey map",
			[]byte("  {  word = 12\t}"),
			[]TokenType{
				MapStartToken,
				WordToken, MapKVDelimToken, NumberToken,
				MapEndToken, EOFToken,
			},
			[]string{"", "word", "", "12", "", ""},
		},
		{
			"empty list",
			[]byte("[]"),
			[]TokenType{ListStartToken, ListEndToken, EOFToken},
			[]string{"", "", ""},
		},
		{
			"list with a couple of items",
			[]byte("[word 1.2]"),
			[]TokenType{
				ListStartToken,
				WordToken, NumberToken,
				ListEndToken, EOFToken,
			},
			[]string{"", "word", "1.2", "", ""},
		},
		{
			"comment",
			[]byte("word # comment\nword"),
			[]TokenType{
				WordToken, WordToken, EOFToken,
			},
			[]string{"word", "word", ""},
		},
		{
			"simple string with double quote",
			[]byte("\"a string\""),
			[]TokenType{StringToken, EOFToken},
			[]string{"a string", ""},
		},
		{
			"simple string with single quote",
			[]byte("'a string'"),
			[]TokenType{StringToken, EOFToken},
			[]string{"a string", ""},
		},
		{
			"string with escaped inner double quote",
			[]byte("\"a \\\" string\""),
			[]TokenType{StringToken, EOFToken},
			[]string{"a \" string", ""},
		},
		{
			"string with escaped inner single quote",
			[]byte("'a \\' string'"),
			[]TokenType{StringToken, EOFToken},
			[]string{"a ' string", ""},
		},
		{
			"string with line breaks",
			[]byte("'a \nstring'"),
			[]TokenType{StringToken, EOFToken},
			[]string{"a \nstring", ""},
		},
		{
			"simple decorator",
			[]byte("decorator(12)"),
			[]TokenType{DecoratorStartToken, NumberToken, DecoratorEndToken, EOFToken},
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
				DecoratorStartToken, WordToken, DecoratorEndToken, MapKVDelimToken, MapStartToken,
				WordToken, MapKVDelimToken, StringToken,
				WordToken, MapKVDelimToken, StringToken,
				WordToken, MapKVDelimToken, WordToken,
				WordToken, MapKVDelimToken, ListStartToken, StringToken, StringToken, ListEndToken,
				WordToken, MapKVDelimToken, StringToken,
				WordToken, MapKVDelimToken,
				MapStartToken,
				WordToken, MapKVDelimToken, StringToken,
				WordToken, MapKVDelimToken, WordToken,
				WordToken, MapKVDelimToken, WordToken,
				WordToken, MapKVDelimToken, DecoratorStartToken, StringToken, DecoratorEndToken,
				MapEndToken,
				MapEndToken, EOFToken,
			},
			[]string{
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

			s := NewScanner(test.document)

			for {
				token := s.Token()

				tokens = append(tokens, token.Type)
				contents = append(contents, token.Content)

				if token.Type == IllegalToken || token.Type == EOFToken {
					break
				}
			}

			assert.Equal(t, test.tokens, tokens)
			assert.Equal(t, test.contents, contents)
		})
	}
}
