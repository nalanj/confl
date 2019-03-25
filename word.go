package confl

// Word represents a word node, and implements Node
type Word struct {

	// val is the value of the word as a string
	val string

	// decorator is the decorator for the word, if any
	decorator string
}

// Type returns the node type for the node, TypeWord
func (w *Word) Type() NodeType {
	return WordType
}

// Children always returns an empty set of nodes since children aren't allowed
// on Words
func (w *Word) Children() []Node {
	return []Node{}
}

// Decorator returns the decorator on the word if there is one
func (w *Word) Decorator() string {
	return w.decorator
}

// Value returns the value of the word as a string
func (w *Word) Value() string {
	return w.val
}
