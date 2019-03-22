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

	// MapStart represents the start of a map token
	MapStart

	// MapEnd represents the end of a map token
	MapEnd
)
