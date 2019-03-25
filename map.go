package confl

// Map represents a map node, and implements Node
type Map struct {
	children  []Node
	decorator string
}

// Type returns the NodeType for this node
func (m *Map) Type() NodeType {
	return MapType
}

// Children returns the children for this node as an order list of keys and
// values as [key, value, key, value...]
func (m *Map) Children() []Node {
	return m.children
}

// Decorator returns the decorator for this node, or the empty string if there
// is none
func (m *Map) Decorator() string {
	return m.decorator
}

// Value always returns an empty string for Maps, since no value is allowed
func (m *Map) Value() string {
	return ""
}
