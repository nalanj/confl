package confl

// Scanner is an interface for tokenizing Confl source
type Scanner interface {
	Token() *Token
}
