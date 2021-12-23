package lexer

type Token struct {
	TokenType TokenType
	Lexeme    string
}

type Lexer struct {
	IsStringLiteral bool
	TokenType       TokenType
	LastChar        rune
	NextChar        rune
	Index           int
}
