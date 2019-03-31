package confl

// listNode represents a list node, and implements Node
type listNode struct {
	children  []Node
	decorator string
}

// Type returns the NodeType for this node
func (l *listNode) Type() NodeType {
	return ListType
}

// Children returns the children for this node as an order list of keys and
// values as [key, value, key, value...]
func (l *listNode) Children() []Node {
	return l.children
}

// Decorator returns the decorator for this node, or the empty string if there
// is none
func (l *listNode) Decorator() string {
	return l.decorator
}

// Value always returns an empty string for Lists, since no value is allowed
func (l *listNode) Value() string {
	return ""
}
