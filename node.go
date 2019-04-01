package confl

// Node is an interface for interacting with an AST
type Node interface {

	// Type returns the type of the node
	Type() NodeType

	// Children returns the children of the node, if any. Maps have their
	// children organized into pairs of keys and values, so, within a map
	// Chidren()[even] is a key and Children()[odd] is a value.
	Children() []Node

	// Decorator returns the decorator for the node, or an empty string if
	// there's no decorator
	Decorator() string

	// Value returns the value of node for value type nodes
	Value() string
}
