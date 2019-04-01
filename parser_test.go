package confl

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleParse() {
	src := []byte(`
		# a simple list of hosts and details
		mail.confl.org={
			purpose=mail
			os=os([linux 4])
			connects_to=[dc.confl.org]
		}

		dc.confl.org={
			purpose=domain_controller
			os=os([windows 10])
		}

		web.confl.org={
			purpose=web_server
			os=os([freebsd 12])
			connects_to=[dc.confl.org mail.confl.org]
		}
		`)

	doc, err := Parse(bytes.NewReader(src))
	if err != nil {
		panic(err)
	}

	fmt.Println(len(doc.Children()))

	// Output: 6
}

func TestParseScanner(t *testing.T) {
	tests := []struct {
		name string
		src  string
		doc  *mapNode
		err  bool
	}{

		{
			"implicit document map",
			`test=23 "also"=this`,
			&mapNode{
				children: []Node{
					&valueNode{nodeType: WordType, val: "test"},
					&valueNode{nodeType: NumberType, val: "23"},
					&valueNode{nodeType: StringType, val: "also"},
					&valueNode{nodeType: WordType, val: "this"},
				},
			},
			false,
		},

		{
			"implicit document map, illegal end token",
			`test=23 "also"=this}`,
			nil,
			true,
		},

		{
			"explicit document map",
			`{test=23 "also"=this}`,
			&mapNode{
				children: []Node{
					&valueNode{nodeType: WordType, val: "test"},
					&valueNode{nodeType: NumberType, val: "23"},
					&valueNode{nodeType: StringType, val: "also"},
					&valueNode{nodeType: WordType, val: "this"},
				},
			},
			false,
		},

		{
			"explicit document map, illegal end token",
			`{test=23 "also"=this`,
			nil,
			true,
		},

		{
			"nested map",
			`map={key=value}`,
			&mapNode{
				children: []Node{
					&valueNode{nodeType: WordType, val: "map"},
					&mapNode{
						children: []Node{
							&valueNode{nodeType: WordType, val: "key"},
							&valueNode{nodeType: WordType, val: "value"},
						},
					},
				},
			},
			false,
		},

		{
			"nested list",
			`list=[item1 item2]`,
			&mapNode{
				children: []Node{
					&valueNode{nodeType: WordType, val: "list"},
					&listNode{
						children: []Node{
							&valueNode{nodeType: WordType, val: "item1"},
							&valueNode{nodeType: WordType, val: "item2"},
						},
					},
				},
			},
			false,
		},

		{
			"simple decorator",
			`dec(test)=23 "also"=this`,
			&mapNode{
				children: []Node{
					&valueNode{nodeType: WordType, val: "test", decorator: "dec"},
					&valueNode{nodeType: NumberType, val: "23"},
					&valueNode{nodeType: StringType, val: "also"},
					&valueNode{nodeType: WordType, val: "this"},
				},
			},
			false,
		},

		{
			"decorator on list as map key errors",
			`dec({key=val})=val`,
			nil,
			true,
		},

		{
			"decorator with map",
			`key=dec({decKey=val})`,
			&mapNode{
				children: []Node{
					&valueNode{nodeType: WordType, val: "key"},
					&mapNode{
						children: []Node{
							&valueNode{nodeType: WordType, val: "decKey"},
							&valueNode{nodeType: WordType, val: "val"},
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
			scan := newScanner([]byte(test.src))
			doc, err := parseScanner(scan)
			assert.Equal(t, test.err, err != nil)
			assert.Equal(t, test.doc, doc)
		})
	}
}
