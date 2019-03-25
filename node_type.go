package confl

// NodeType represents the available node types
type NodeType int

const (
	// NumberType is the NodeType for numbers
	NumberType NodeType = iota

	// WordType is the NodeType for words
	WordType

	// StringType is the NodeType for strings
	StringType

	// MapType is the NodeType for maps
	MapType

	// ListType is the NodeType for lists
	ListType
)
