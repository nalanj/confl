package scanner

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

const (
	runeEOF rune = -1
	runeBOM rune = 0xFEFF
)

// SrcScanner is a scanner of Confl code
type SrcScanner struct {

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
func (s *SrcScanner) next() bool {
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
func Init(src []byte) *SrcScanner {
	return &SrcScanner{src: src}
}

// Token returns the next token, offset, and string content
func (s *SrcScanner) Token() Token {
	var token Token

	// advance to the first character if at the beginning of the source
	if s.offset == 0 && !s.next() {
		token.Type = Illegal
		return token
	}

	// TODO: handle BOM if at 0

	s.skipWhitespace()

	offset := s.offset
	advance := false

	switch {
	case s.ch == runeEOF:
		token.Type = EOF
	case s.isDigit():
		token.Type, token.Content = s.scanNumber()
	case s.isLetter():
		token.Type, token.Content = s.scanWord()
	case s.ch == ')':
		token.Type = DecoratorEnd
		advance = true
	case s.ch == '{':
		token.Type = MapStart
		advance = true
	case s.ch == '}':
		token.Type = MapEnd
		advance = true
	case s.ch == '=':
		token.Type = MapKVDelim
		advance = true
	case s.ch == '[':
		token.Type = ListStart
		advance = true
	case s.ch == ']':
		token.Type = ListEnd
		advance = true
	case s.ch == '#':
		token.Type, token.Content = s.scanComment()
	case s.isStringDelim():
		token.Type, token.Content = s.scanString()
	default:
		token.Type = Illegal
	}

	if advance {
		if !s.next() {
			token.Type = Illegal
			token.Content = string(s.src[offset:s.nextOffset])
		}
	}

	return token
}

// isWhitespace returns whether the current ch is whitespace
func (s *SrcScanner) isWhitespace() bool {
	if s.ch == ' ' || s.ch == '\r' || s.ch == '\n' || s.ch == '\t' {
		return true
	}
	return false
}

// isPunctuation notes if the current ch is one of the punctuation chars
func (s *SrcScanner) isPunctuation() bool {
	if s.ch == '{' || s.ch == '}' || s.ch == '[' || s.ch == ']' || s.ch == '=' ||
		s.ch == '(' || s.ch == ')' || s.ch == runeEOF {

		return true
	}
	return false
}

// isDigit returns if the current ch is a digit
func (s *SrcScanner) isDigit() bool {
	return s.ch >= '0' && s.ch <= '9' || s.ch > utf8.RuneSelf && unicode.IsDigit(s.ch)
}

// isLetter returns if the current ch is a letter
func (s *SrcScanner) isLetter() bool {
	return s.ch >= 'a' && s.ch <= 'z' || s.ch >= 'A' && s.ch <= 'Z' || s.ch > utf8.RuneSelf && unicode.IsLetter(s.ch)
}

// isStringDelim returns true if the character is a string delimiter
func (s *SrcScanner) isStringDelim() bool {
	return s.ch == '"' || s.ch == '\''
}

// skipWhitespace reads through whitespace
func (s *SrcScanner) skipWhitespace() {
	for s.isWhitespace() {
		s.next()
	}
}

// scanNumber scans numbers
func (s *SrcScanner) scanNumber() (TokenType, string) {
	startOff := s.offset
	seenDecimal := false

	for !s.isPunctuation() && !s.isWhitespace() {
		if !s.isDigit() && s.ch != '.' {
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

	return Number, string(s.src[startOff:s.offset])
}

// scanWord scans a word or a decorator
func (s *SrcScanner) scanWord() (TokenType, string) {
	startOff := s.offset

	for !s.isPunctuation() && !s.isWhitespace() {
		if !s.next() {
			return Illegal, string(s.src[startOff:s.nextOffset])
		}
	}

	content := string(s.src[startOff:s.offset])

	// if we're on a (, this is a decorator and not a word
	if s.ch == '(' {
		if !s.next() {
			return Illegal, string(s.src[startOff:s.nextOffset])
		}
		return DecoratorStart, content
	}

	return Word, content
}

// scanComment scans a comment
func (s *SrcScanner) scanComment() (TokenType, string) {
	startOff := s.offset

	for s.ch != '\n' {
		if !s.next() {
			return Illegal, string(s.src[startOff:s.nextOffset])
		}
	}

	return Comment, string(s.src[startOff:s.offset])
}

// scanString scans a string. Should be called on a string opening char
func (s *SrcScanner) scanString() (TokenType, string) {
	delim := s.ch
	startOff := s.offset
	var content []byte
	escape := false

	// skip the opening char
	if !s.next() {
		return Illegal, string(s.src[startOff:s.nextOffset])
	}
	startOff++

	for {
		if s.ch == '\\' {
			if escape {
				escape = false
			} else {
				escape = true
			}
		}

		if s.ch == delim {
			if escape {
				escape = false
			} else {
				break
			}
		}

		if !escape {
			content = append(content, s.src[s.offset:s.nextOffset]...)
		}

		if !s.next() {
			return Illegal, string(s.src[startOff:s.nextOffset])
		}
	}

	// skip the ending char
	if !s.next() {
		return Illegal, string(s.src[startOff:s.nextOffset])
	}

	return String, string(content)
}
