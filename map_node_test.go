package confl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapKVPairs(t *testing.T) {
	m := &mapNode{
		children: []Node{
			&valueNode{nodeType: WordType, val: "key1"},
			&valueNode{nodeType: WordType, val: "val1"},
			&valueNode{nodeType: WordType, val: "key2"},
			&valueNode{nodeType: WordType, val: "val2"},
		},
	}

	assert.Equal(
		t,
		KVPairs(m),
		[][2]Node{
			[2]Node{
				&valueNode{nodeType: WordType, val: "key1"},
				&valueNode{nodeType: WordType, val: "val1"},
			},
			[2]Node{
				&valueNode{nodeType: WordType, val: "key2"},
				&valueNode{nodeType: WordType, val: "val2"},
			},
		},
	)
}
