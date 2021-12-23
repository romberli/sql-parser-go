package lexer

const (
	// char
	UnderBarRune = '_'
	// comparison operator
	GTRune          = '>'
	LTRune          = '<'
	EqualRune       = '='
	ExclamationRune = '!'
	// arithmetic operator
	PlusRune     = '+'
	MinusRune    = '-'
	MultiplyRune = '*'
	DivideRune   = '/'
	ModRune      = '%'
	// separator
	ColonRune            = ','
	SemicolonRune        = ';'
	LeftParenthesisRune  = '('
	RightParenthesisRune = ')'
	SingleQuoteRune      = '\''
	// white space
	SpaceRune   = ' '
	TabRune     = '\t'
	ReturnRune  = '\r'
	NewLineRune = '\n'
)

func IsAlphabet(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == UnderBarRune
}

func IsNumber(c rune) bool {
	return c >= '0' && c <= '9'
}

func IsComparisonOperator(c rune) bool {
	return c == GTRune || c == LTRune || c == EqualRune || c == ExclamationRune
}

func IsArithmeticOperator(c rune) bool {
	return c == PlusRune || c == MinusRune || c == MultiplyRune || c == DivideRune || c == ModRune
}

func IsSeparator(c rune) bool {
	return c == ColonRune || c == SemicolonRune || c == LeftParenthesisRune || c == RightParenthesisRune || c == SingleQuoteRune
}

func IsWhiteSpace(c rune) bool {
	return c == SpaceRune || c == TabRune || c == ReturnRune || c == NewLineRune
}

func ShouldMatch(c rune) bool {
	return IsWhiteSpace(c) || IsComparisonOperator(c) || IsArithmeticOperator(c) || IsSeparator(c)
}
