package confl

// Token represents a token from a Confl file
type Token uint

const (
	// Illegal represents an illegal token
	Illegal Token = iota

	// EOF represents a token for the end of the file
	EOF

	// Number represents a number token
	Number

	// Word represents a word token
	Word

	// MapStart represents the start of a map token
	MapStart

	// MapEnd represents the end of a map token
	MapEnd

	// MapKVDelim represents a key value delimiter token
	MapKVDelim

	// ListStart represents the start of a list token
	ListStart

	// ListEnd represents the end of a list token
	ListEnd

	// Comment represents a comment token
	Comment

	// String represents a string token
	String

	// DecoratorStart represents the start of a decorator
	DecoratorStart

	// DecoratorEnd represents the end of a decorator
	DecoratorEnd
)
