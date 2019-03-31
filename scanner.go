package confl

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

const (
	runeEOF rune = -1
	runeBOM rune = 0xFEFF
)

// scanner is a scanner of Confl code
type scanner struct {

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

	// peeked is the set of peeked tokens
	peeked []*Token
}

// next returns the next character from the scanner
func (s *scanner) next() bool {
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

// newScanner returns a new scanner based on the given source
func newScanner(src []byte) *scanner {
	return &scanner{src: src, peeked: []*Token{}}
}

// Token returns the peeked token if it's there, or the next token
func (s *scanner) Token() *Token {
	var token *Token

	if len(s.peeked) > 0 {
		token = s.peeked[0]
		s.peeked = s.peeked[1:]
	}

	if token == nil {
		token = s.nextToken()
	}

	return token
}

// Peek returns up to count tokens, stopping on EOF if one is hit
func (s *scanner) Peek(count int) []*Token {
	for i := 0; i < count; i++ {
		var token *Token

		if len(s.peeked) > i {
			token = s.peeked[i]
		} else {
			token = s.nextToken()
			s.peeked = append(s.peeked, token)
		}

		if token.Type == EOFToken {
			break
		}
	}

	return s.peeked[0:count]
}

// nextToken returns the next token, offset, and string content
func (s *scanner) nextToken() *Token {
	var token Token

	// advance to the first character if at the beginning of the source
	if s.offset == 0 && !s.next() {
		token.Type = IllegalToken
		return &token
	}

	// TODO: handle BOM if at 0

	for s.skipWhitespace() || s.skipComment() {
	}

	offset := s.offset
	advance := false

	switch {
	case s.ch == runeEOF:
		token.Type = EOFToken
	case s.isDigit():
		token.Type, token.Content = s.scanNumber()
	case s.isLetter():
		token.Type, token.Content = s.scanWord()
	case s.ch == ')':
		token.Type = DecoratorEndToken
		advance = true
	case s.ch == '{':
		token.Type = MapStartToken
		advance = true
	case s.ch == '}':
		token.Type = MapEndToken
		advance = true
	case s.ch == '=':
		token.Type = MapKVDelimToken
		advance = true
	case s.ch == '[':
		token.Type = ListStartToken
		advance = true
	case s.ch == ']':
		token.Type = ListEndToken
		advance = true
	case s.isStringDelim():
		token.Type, token.Content = s.scanString()
	default:
		token.Type = IllegalToken
	}

	if advance {
		if !s.next() {
			token.Type = IllegalToken
			token.Content = string(s.src[offset:s.nextOffset])
		}
	}

	return &token
}

// isWhitespace returns whether the current ch is whitespace
func (s *scanner) isWhitespace() bool {
	if s.ch == ' ' || s.ch == '\r' || s.ch == '\n' || s.ch == '\t' {
		return true
	}
	return false
}

// isPunctuation notes if the current ch is one of the punctuation chars
func (s *scanner) isPunctuation() bool {
	if s.ch == '{' || s.ch == '}' || s.ch == '[' || s.ch == ']' || s.ch == '=' ||
		s.ch == '(' || s.ch == ')' || s.ch == '#' || s.ch == runeEOF {

		return true
	}
	return false
}

// isDigit returns if the current ch is a digit
func (s *scanner) isDigit() bool {
	return s.ch >= '0' && s.ch <= '9' || s.ch > utf8.RuneSelf && unicode.IsDigit(s.ch)
}

// isLetter returns if the current ch is a letter
func (s *scanner) isLetter() bool {
	return s.ch >= 'a' && s.ch <= 'z' || s.ch >= 'A' && s.ch <= 'Z' || s.ch > utf8.RuneSelf && unicode.IsLetter(s.ch)
}

// isStringDelim returns true if the character is a string delimiter
func (s *scanner) isStringDelim() bool {
	return s.ch == '"' || s.ch == '\''
}

// skipWhitespace reads through whitespace
func (s *scanner) skipWhitespace() bool {
	skipped := false

	for s.isWhitespace() {
		skipped = true
		s.next()
	}

	return skipped
}

// skipComment ignores a comment
func (s *scanner) skipComment() bool {
	skipped := false

	if s.ch == '#' {
		skipped = true

		for s.ch != '\n' {
			s.next()
		}
	}

	return skipped
}

// scanNumber scans numbers
func (s *scanner) scanNumber() (TokenType, string) {
	startOff := s.offset
	seenDecimal := false

	for !s.isPunctuation() && !s.isWhitespace() {
		if !s.isDigit() && s.ch != '.' {
			return IllegalToken, string(s.src[startOff:s.nextOffset])
		}

		if s.ch == '.' {
			if seenDecimal {
				return IllegalToken, string(s.src[startOff:s.nextOffset])
			}

			seenDecimal = true
		}

		if !s.next() {
			return IllegalToken, string(s.src[startOff:s.nextOffset])
		}
	}

	return NumberToken, string(s.src[startOff:s.offset])
}

// scanWord scans a word or a decorator
func (s *scanner) scanWord() (TokenType, string) {
	startOff := s.offset

	for !s.isPunctuation() && !s.isWhitespace() {
		if !s.next() {
			return IllegalToken, string(s.src[startOff:s.nextOffset])
		}
	}

	content := string(s.src[startOff:s.offset])

	// if we're on a (, this is a decorator and not a word
	if s.ch == '(' {
		if !s.next() {
			return IllegalToken, string(s.src[startOff:s.nextOffset])
		}
		return DecoratorStartToken, content
	}

	return WordToken, content
}

// scanString scans a string. Should be called on a string opening char
func (s *scanner) scanString() (TokenType, string) {
	delim := s.ch
	startOff := s.offset
	var content []byte
	escape := false

	// skip the opening char
	if !s.next() {
		return IllegalToken, string(s.src[startOff:s.nextOffset])
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
			return IllegalToken, string(s.src[startOff:s.nextOffset])
		}
	}

	// skip the ending char
	if !s.next() {
		return IllegalToken, string(s.src[startOff:s.nextOffset])
	}

	return StringToken, string(content)
}
