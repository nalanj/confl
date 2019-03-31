package confl

import "fmt"

// parseError represents an error while parsing
type parseError struct {

	// msg is the error message
	msg string

	// src is the source around where the error happened
	src []byte

	// offset is the offset in src where the error happened
	offset int
}

// Error returns the full error message
func (p *parseError) Error() string {
	return fmt.Sprintf(
		"%s\n%s\n",
		p.msg,
		p.src,
	)
}

// newParseError returns a new parse error based on the given msg, scanner, and
// offset
func newParseError(msg string, scan *Scanner, offset int) *parseError {
	start := offset - 10
	if start < 0 {
		start = 0
	}

	end := offset + 20
	if end > len(scan.src) {
		end = len(scan.src)
	}

	return &parseError{
		msg:    msg,
		src:    scan.src[start:end],
		offset: offset - start,
	}
}
