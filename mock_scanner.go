package confl

// mockScanner outputs tokens directly
type mockScanner struct {
	tokens  []*Token
	current int
}

// Token returns the next token from the mock scanner
func (s *mockScanner) Token() *Token {
	cur := s.current
	s.current++
	return s.tokens[cur]
}

// Peek returns the next count tokens, up to an eof
func (s *mockScanner) Peek(count int) []*Token {
	end := s.current + count
	if end > len(s.tokens) {
		end = len(s.tokens)
	}

	return s.tokens[s.current:end]
}
