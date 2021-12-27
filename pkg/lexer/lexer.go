package lexer

type Lexer struct {
	IsStringLiteral bool
	TokenType       TokenType
	LastChar        rune
	NextChar        rune
	Index           int
}
