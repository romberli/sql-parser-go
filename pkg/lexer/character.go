package lexer

const (
	// ascii boundary
	digitStart    = 48
	digitEnd      = 57
	alphabetStart = 97
	alphabetEnd   = 122

	underBarRune = '_'
	singleQuote  = '\''
	EpsilonRune  = 'ε'
)

type CharacterSet struct {
	Alphabets []rune
	Digits    []rune
}

// NewCharacterSet returns a new *CharacterSet
func NewCharacterSet(alphabets, digits []rune) *CharacterSet {
	return &CharacterSet{
		Alphabets: alphabets,
		Digits:    digits,
	}
}

// NewCharacterSetWithDefault returns a new *CharacterSet with default
func NewCharacterSetWithDefault() *CharacterSet {
	alphabets := []rune{underBarRune}

	for i := alphabetStart; i <= alphabetEnd; i++ {
		alphabets = append(alphabets, rune(i))
	}

	var digits []rune
	for i := digitStart; i <= digitEnd; i++ {
		digits = append(digits, rune(i))
	}

	return NewCharacterSet(alphabets, digits)
}

// GetAlphabets returns the alphabet runes
func (cs *CharacterSet) GetAlphabets() []rune {
	return cs.Alphabets
}

// GetAlphabets returns the digit runes
func (cs *CharacterSet) GetDigits() []rune {
	return cs.Digits
}
