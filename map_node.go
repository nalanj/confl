package confl

// mapNode represents a map node, and implements Node
type mapNode struct {
	children  []Node
	decorator string
}

// Type returns the NodeType for this node
func (m *mapNode) Type() NodeType {
	return MapType
}

// Children returns the children for this node as an order list of keys and
// values as [key, value, key, value...]
func (m *mapNode) Children() []Node {
	return m.children
}

// Decorator returns the decorator for this node, or the empty string if there
// is none
func (m *mapNode) Decorator() string {
	return m.decorator
}

// Value always returns an empty string for Maps, since no value is allowed
func (m *mapNode) Value() string {
	return ""
}

// KVPair is a key value pair out of a map node
type KVPair struct {

	// Key is the key for the KVPair
	Key Node

	// Value is the value for the KVPair
	Value Node
}

// KVPairs returns a set of key/value pairs from the map
func KVPairs(n Node) []KVPair {
	if n.Type() != MapType {
		return []KVPair{}
	}

	pairs := []KVPair{}

	var key Node
	for _, node := range n.Children() {
		if key == nil {
			key = node
		} else {
			pairs = append(pairs, KVPair{key, node})
			key = nil
		}
	}

	return pairs
}
