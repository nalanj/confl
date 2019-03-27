package confl

// parseError represents an error while parsing
type parseError struct {

	// Message is the error message
	msg string
}

// Error returns the full error message
func (p *parseError) Error() string {
	return p.msg
}
