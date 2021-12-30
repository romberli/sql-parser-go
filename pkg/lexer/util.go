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
	CommaRune            = ','
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

// IsAlphabet returns if the given rune is an alphabet
func IsAlphabet(c rune) bool {
	return c >= 'a' && c <= 'z' || c == UnderBarRune
}

// IsDigit returns if the given rune is a digit
func IsDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

// IsWhiteSpace returns if the given rune is a white space
func IsWhiteSpace(c rune) bool {
	return c == SpaceRune || c == TabRune || c == ReturnRune || c == NewLineRune
}

// IsWhiteSpace returns if the given rune is either an alphabet or a digit
func IsAlphabetOrDigit(c rune) bool {
	return IsAlphabet(c) || IsDigit(c)
}
