package confl

// ValueNode is a value node
type ValueNode struct {

	// nodeType is the type of this node
	nodeType NodeType

	// val is the value of the string
	val string

	// decorator is the decorator for the string, if any
	decorator string
}

// Type returns the node type for the node
func (n *ValueNode) Type() NodeType {
	return n.nodeType
}

// Children always returns an empty set of nodes since children aren't allowed
// on value nodes
func (n *ValueNode) Children() []Node {
	return []Node{}
}

// Decorator returns the decorator on the node if there is one
func (n *ValueNode) Decorator() string {
	return n.decorator
}

// Value returns the value of the node as a string
func (n *ValueNode) Value() string {
	return n.val
}
