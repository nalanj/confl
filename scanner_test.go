package confl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScan(t *testing.T) {
	tests := []struct {
		name     string
		document []byte
		tokens   []tokenType
		contents []string
	}{
		{
			"simple integer number",
			[]byte("12"),
			[]tokenType{numberToken, eofToken},
			[]string{"12", ""},
		},
		{
			"simple decimal number",
			[]byte("12.3"),
			[]tokenType{numberToken, eofToken},
			[]string{"12.3", ""},
		},
		{
			"decimal starting with 0",
			[]byte("0.3"),
			[]tokenType{numberToken, eofToken},
			[]string{"0.3", ""},
		},
		{
			"hex number lower case",
			[]byte("0x3"),
			[]tokenType{numberToken, eofToken},
			[]string{"0x3", ""},
		},
		{
			"hex number upper case",
			[]byte("0X3"),
			[]tokenType{numberToken, eofToken},
			[]string{"0X3", ""},
		},

		{
			"illegal: two decimal number",
			[]byte("1.2.3"),
			[]tokenType{illegalToken},
			[]string{"1.2."},
		},
		{
			"word",
			[]byte(" testing "),
			[]tokenType{wordToken, eofToken},
			[]string{"testing", ""},
		},
		{
			"empty map",
			[]byte("{}"),
			[]tokenType{mapStartToken, mapEndToken, eofToken},
			[]string{"", "", ""},
		},
		{
			"map",
			[]byte("{word=12}"),
			[]tokenType{
				mapStartToken,
				wordToken, mapKVDelimToken, numberToken,
				mapEndToken, eofToken,
			},
			[]string{"", "word", "", "12", "", ""},
		},
		{
			"spacey map",
			[]byte("  {  word = 12\t}"),
			[]tokenType{
				mapStartToken,
				wordToken, mapKVDelimToken, numberToken,
				mapEndToken, eofToken,
			},
			[]string{"", "word", "", "12", "", ""},
		},
		{
			"empty list",
			[]byte("[]"),
			[]tokenType{listStartToken, listEndToken, eofToken},
			[]string{"", "", ""},
		},
		{
			"list with a couple of items",
			[]byte("[word 1.2]"),
			[]tokenType{
				listStartToken,
				wordToken, numberToken,
				listEndToken, eofToken,
			},
			[]string{"", "word", "1.2", "", ""},
		},
		{
			"comment",
			[]byte("# comment\nword"),
			[]tokenType{
				wordToken, eofToken,
			},
			[]string{"word", ""},
		},
		{
			"simple string with double quote",
			[]byte("\"a string\""),
			[]tokenType{stringToken, eofToken},
			[]string{"a string", ""},
		},
		{
			"simple string with single quote",
			[]byte("'a string'"),
			[]tokenType{stringToken, eofToken},
			[]string{"a string", ""},
		},
		{
			"string with escaped inner double quote",
			[]byte("\"a \\\" string\""),
			[]tokenType{stringToken, eofToken},
			[]string{"a \" string", ""},
		},
		{
			"string with escaped inner single quote",
			[]byte("'a \\' string'"),
			[]tokenType{stringToken, eofToken},
			[]string{"a ' string", ""},
		},
		{
			"string with line breaks",
			[]byte("'a \nstring'"),
			[]tokenType{stringToken, eofToken},
			[]string{"a \nstring", ""},
		},
		{
			"simple decorator",
			[]byte("decorator(12)"),
			[]tokenType{decoratorStartToken, numberToken, decoratorEndToken, eofToken},
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
			[]tokenType{
				decoratorStartToken, wordToken, decoratorEndToken, mapKVDelimToken, mapStartToken,
				wordToken, mapKVDelimToken, stringToken,
				wordToken, mapKVDelimToken, stringToken,
				wordToken, mapKVDelimToken, wordToken,
				wordToken, mapKVDelimToken, listStartToken, stringToken, stringToken, listEndToken,
				wordToken, mapKVDelimToken, stringToken,
				wordToken, mapKVDelimToken,
				mapStartToken,
				wordToken, mapKVDelimToken, stringToken,
				wordToken, mapKVDelimToken, wordToken,
				wordToken, mapKVDelimToken, wordToken,
				wordToken, mapKVDelimToken, decoratorStartToken, stringToken, decoratorEndToken,
				mapEndToken,
				mapEndToken, eofToken,
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
			tokens := []tokenType{}
			contents := []string{}

			s := newScanner(test.document)

			for {
				token := s.Token()

				tokens = append(tokens, token.Type)
				contents = append(contents, token.Content)

				if token.Type == illegalToken || token.Type == eofToken {
					break
				}
			}

			assert.Equal(t, test.tokens, tokens)
			assert.Equal(t, test.contents, contents)
		})
	}
}
