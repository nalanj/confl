package confl

// TokenType represents the various types of tokens
type TokenType uint

const (
	// IllegalToken represents an illegal token
	IllegalToken TokenType = iota

	// EOFToken represents a token for the end of the file
	EOFToken

	// NumberToken represents a number token
	NumberToken

	// WordToken represents a word token
	WordToken

	// MapStartToken represents the start of a map token
	MapStartToken

	// MapEndToken represents the end of a map token
	MapEndToken

	// MapKVDelimToken represents a key value delimiter token
	MapKVDelimToken

	// ListStartToken represents the start of a list token
	ListStartToken

	// ListEndToken represents the end of a list token
	ListEndToken

	// CommentToken represents a comment token
	CommentToken

	// StringToken represents a string token
	StringToken

	// DecoratorStartToken represents the start of a decorator
	DecoratorStartToken

	// DecoratorEndToken represents the end of a decorator
	DecoratorEndToken
)

// Token is a token from the scanner
type Token struct {

	// Type of the token
	Type TokenType

	// Offset of the token in the source
	Offset int

	// Content of the token
	Content string
}
