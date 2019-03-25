package confl

// parseError represents an error while parsing
type parseError struct {

	// Message is the error message
	msg string

	// offset is the location where the error occurred
	offset int
}

// Error returns the full error message
func (p *parseError) Error() string {
	return p.msg
}
