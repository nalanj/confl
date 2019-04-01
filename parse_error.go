package confl

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// parseError represents an error while parsing
type parseError struct {

	// msg is the error message
	msg string

	// src is the source around where the error happened
	src []byte

	// offset is the offset in src where the error happened
	offset int

	// length is the length of the token where the error happened
	length int
}

// Error returns the full error message
func (p *parseError) Error() string {
	focusEnd := p.offset + p.length

	pre := string(p.src[0:p.offset])
	focus := string(p.src[p.offset:focusEnd])
	post := string(p.src[focusEnd:len(p.src)])

	return fmt.Sprintf(
		"%s\n%s%s%s\n%s%s%s\n",
		p.msg,
		pre, focus, post,
		strings.Repeat(" ", utf8.RuneCountInString(pre)),
		strings.Repeat("^", utf8.RuneCountInString(focus)),
		strings.Repeat(" ", utf8.RuneCountInString(post)),
	)
}

// newParseError returns a new parse error based on the given msg, scanner, and
// offset
func newParseError(msg string, scan *scanner, offset int, length int) *parseError {
	start := offset - 10
	if start < 0 {
		start = 0
	}

	end := offset + length + 20
	if end > len(scan.src) {
		end = len(scan.src)
	}

	// length has to at least be 1 and less than end
	if length == 0 {
		length = 1
	}
	if offset+length > end {
		length = end - offset
	}

	src := scan.src[start:end]

	// if offset is EOF, show it
	if offset == end {
		src = append(src, []byte("(EOF)")...)
		length = 5
	}

	return &parseError{
		msg:    msg,
		src:    src,
		offset: offset - start,
		length: length,
	}
}
