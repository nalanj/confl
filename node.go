package confl

// Node is an interface for interacting with an AST
type Node interface {

	// Type returns the type of the node
	Type() NodeType

	// Children returns the children of the node, if any
	Children() []Node

	// Decorator returns the decorator for the node, or an empty string if
	// there's no decorator
	Decorator() string

	// Value returns the value of node for value type nodes
	Value() string
}
