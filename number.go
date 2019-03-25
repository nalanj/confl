package confl

// Number represents a number node, and implements Node
type Number struct {

	// val is the value of the number as a string
	val string

	// decorator is the decorator for the number, if any
	decorator string
}

// Type returns the node type for the node, TypeNumber
func (n *Number) Type() NodeType {
	return NumberType
}

// Children always returns an empty set of nodes since children aren't allowed
// on Number
func (n *Number) Children() []Node {
	return []Node{}
}

// Decorator returns the decorator on the number if there is one
func (n *Number) Decorator() string {
	return n.decorator
}

// Value returns the value of the number as a string
func (n *Number) Value() string {
	return n.val
}
