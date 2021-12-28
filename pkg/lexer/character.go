package lexer

const (
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

func NewCharacterSet(alphabets, digits []rune) *CharacterSet {
	return &CharacterSet{
		Alphabets: alphabets,
		Digits:    digits,
	}
}

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

func (cs *CharacterSet) Init() {
	cs.Alphabets = append(cs.Alphabets, underBarRune)
	for i := alphabetStart; i <= alphabetEnd; i++ {
		cs.Alphabets = append(cs.Alphabets, rune(i))
	}

	for i := digitStart; i <= digitEnd; i++ {
		cs.Digits = append(cs.Digits, rune(i))
	}
}

func (cs *CharacterSet) GetAlphabets() []rune {
	return cs.Alphabets
}

func (cs *CharacterSet) GetDigits() []rune {
	return cs.Digits
}
