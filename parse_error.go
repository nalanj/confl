package confl

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// ParseError represents an error while parsing
type ParseError struct {

	// msg is the error message
	msg string

	// src is the source around where the error happened
	src []byte

	// offset is the offset in src where the error happened
	offset int

	// length is the length of the token where the error happened
	length int

	// line is the line for the error
	line int
}

// Error returns the error message
func (p *ParseError) Error() string {
	return p.msg
}

// ErrorWithCode returns a multi-line formatted version of the error including
// the code where the error occurred.
func (p *ParseError) ErrorWithCode() string {
	focusEnd := p.offset + p.length

	line := fmt.Sprintf("Line %d: ", p.line)
	pre := string(p.src[0:p.offset])
	focus := string(p.src[p.offset:focusEnd])
	post := string(p.src[focusEnd:len(p.src)])

	return fmt.Sprintf(
		"%s\n%s%s%s%s\n%s%s%s%s\n",
		p.msg,
		line, pre, focus, post,
		strings.Repeat(" ", utf8.RuneCountInString(line)),
		strings.Repeat(" ", utf8.RuneCountInString(pre)),
		strings.Repeat("^", utf8.RuneCountInString(focus)),
		strings.Repeat(" ", utf8.RuneCountInString(post)),
	)
}

// newParseError returns a new parse error based on the given msg, scanner, and
// offset
func newParseError(msg string, scan *scanner, offset, length int) *ParseError {
	start := scan.lineStart
	if start < 0 {
		start = 0
	}

	// if the offset wraps lines, just go with the last line
	if offset < start {
		offset = start
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

	return &ParseError{
		msg:    msg,
		src:    src,
		offset: offset - start,
		length: length,
		line:   scan.line,
	}
}
