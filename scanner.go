package confl

import (
	"errors"
	"unicode/utf8"
)

const (
	runeEOF rune = -1
	runeBOM rune = 0xFEFF
)

// Scanner is a scanner of Confl code
type Scanner struct {

	// offset within the document
	offset int

	// next offset within the document
	nextOffset int

	// src is the source of the document
	src []byte

	// ch is the current rune
	ch rune

	// err is the current error
	err error
}

// next returns the next character from the scanner
func (s *Scanner) next() bool {
	if s.nextOffset < len(s.src) {
		s.offset = s.nextOffset

		r, w := rune(s.src[s.offset]), 1

		if r == 0 {
			s.err = errors.New("Illegal character \\0")
			return false
		} else if r >= utf8.RuneSelf {
			r, w = utf8.DecodeRune(s.src[s.offset:])
			if r == utf8.RuneError && w == 1 {
				s.err = errors.New("illegal utf-8 encoding")
				return false
			} else if r == runeBOM && s.offset > 0 {
				s.err = errors.New("illegal byte order mark")
				return false
			}
		}

		s.nextOffset += w
		s.ch = r
	} else {
		s.offset = len(s.src)
		s.ch = runeEOF
	}

	return true
}

// Init returns a new scanner based on the given source
func Init(src []byte) *Scanner {
	return &Scanner{src: src}
}

// Token returns the next token, offset, and string content
func (s *Scanner) Token() (Token, int, string) {
	if !s.next() {
		return Illegal, s.offset, ""
	}

	// TODO: handle BOM if at 0

	s.skipWhitespace()

	switch {
	case s.ch == runeEOF:
		return EOF, s.offset, ""
	case s.ch >= '0' && s.ch <= '9':
		offset := s.offset
		token, numStr := s.scanNumber()
		return token, offset, numStr
	case s.ch == '{':
		return MapStart, s.offset, ""
	case s.ch == '}':
		return MapEnd, s.offset, ""
	default:
		return Illegal, s.offset, ""
	}
}

// isWhitespace returns whether the current ch is whitespace
func (s *Scanner) isWhitespace() bool {
	if s.ch == ' ' || s.ch == '\r' || s.ch == '\n' || s.ch == '\t' {
		return true
	}
	return false
}

// isPunctuation notes if the current ch is one of the punctuation chars
func (s *Scanner) isPunctuation() bool {
	if s.ch == '{' || s.ch == '}' || s.ch == '[' || s.ch == ']' || s.ch == '=' || s.ch == runeEOF {
		return true
	}
	return false
}

// skipWhitespace reads through whitespace
func (s *Scanner) skipWhitespace() {
	for s.isWhitespace() {
		s.next()
	}
}

// scanNumber scans numbers
func (s *Scanner) scanNumber() (Token, string) {
	startOff := s.offset
	seenDecimal := false

	for !s.isPunctuation() && !s.isWhitespace() {
		if (s.ch <= '0' || s.ch >= '9') && s.ch != '.' {
			return Illegal, string(s.src[startOff:s.nextOffset])
		}

		if s.ch == '.' {
			if seenDecimal {
				return Illegal, string(s.src[startOff:s.nextOffset])
			}

			seenDecimal = true
		}

		if !s.next() {
			return Illegal, string(s.src[startOff:s.nextOffset])
		}
	}

	return Number, string(s.src[startOff:s.nextOffset])
}
