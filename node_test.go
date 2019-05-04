package confl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeIsText(t *testing.T) {
	tests := []struct {
		name     string
		nodeType NodeType
		text     bool
	}{
		{"true for word", WordType, true},
		{"true for string", StringType, true},
		{"false for number", NumberType, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.text, IsText(&valueNode{nodeType: test.nodeType}))
		})
	}
}
