package confl

// TokenType represents the various types of tokens
type tokenType uint

const (
	// illegalToken represents an illegal token
	illegalToken tokenType = iota

	// eofToken represents a token for the end of the file
	eofToken

	// numberToken represents a number token
	numberToken

	// wordToken represents a word token
	wordToken

	// mapStartToken represents the start of a map token
	mapStartToken

	// mapEndToken represents the end of a map token
	mapEndToken

	// mapKVDelimToken represents a key value delimiter token
	mapKVDelimToken

	// listStartToken represents the start of a list token
	listStartToken

	// listEndToken represents the end of a list token
	listEndToken

	// stringToken represents a string token
	stringToken

	// decoratorStartToken represents the start of a decorator
	decoratorStartToken

	// decoratorEndToken represents the end of a decorator
	decoratorEndToken
)

// typeString converts a token type to a string
func (t tokenType) String() string {
	switch t {
	case decoratorEndToken:
		return ")"
	case listEndToken:
		return "]"
	case mapEndToken:
		return "}"
	case eofToken:
		return "EOF"
	default:
		panic("Cannot convert token type")
	}
}

// token is a token from the scanner
type token struct {

	// Type of the token
	Type tokenType

	// Offset of the token in the source
	Offset int

	// Content of the token
	Content string
}
