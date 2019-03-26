package confl

// List represents a list node, and implements Node
type List struct {
	children  []Node
	decorator string
}

// Type returns the NodeType for this node
func (l *List) Type() NodeType {
	return ListType
}

// Children returns the children for this node as an order list of keys and
// values as [key, value, key, value...]
func (l *List) Children() []Node {
	return l.children
}

// Decorator returns the decorator for this node, or the empty string if there
// is none
func (l *List) Decorator() string {
	return l.decorator
}

// Value always returns an empty string for Lists, since no value is allowed
func (l *List) Value() string {
	return ""
}
