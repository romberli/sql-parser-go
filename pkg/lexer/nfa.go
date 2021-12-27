package lexer

const (
	// keyword
	SelectString = "select"
	FromString   = "from"
	AsString     = "as"
	WhereString  = "where"
	AndString    = "and"
	// comparison operator
	GEString        = ">="
	GTString        = ">"
	LEString        = "<="
	LTString        = "<"
	EqualString     = "=="
	NotEqual1String = "!="
	NotEqual2String = "<>"
)

var (
	MultiRuneList = map[TokenType]string{
		// keyword
		Select: SelectString,
		From:   FromString,
		As:     AsString,
		Where:  WhereString,
		And:    AndString,
		// comparison operator
		GE:        GEString,
		GT:        GTString,
		LE:        LEString,
		LT:        LTString,
		Equal:     EqualString,
		NotEqual1: NotEqual1String,
		NotEqual2: NotEqual2String,
	}
	SingleRuneList = map[TokenType]rune{
		// arithmetic operator
		Plus:     PlusRune,
		Minus:    MinusRune,
		Multiply: MultiplyRune,
		Divide:   DivideRune,
		Mod:      ModRune,
		// parenthesis
		LeftParenthesis:  LeftParenthesisRune,
		RightParenthesis: RightParenthesisRune,
	}
)

type NFA struct {
	CharacterSet *CharacterSet
	Index        int
	InitState    *State
}

func NewNFA(cs *CharacterSet) *NFA {
	return &NFA{
		CharacterSet: cs,
		Index:        -1,
	}
}

func NewNFAWithDefault() *NFA {
	cs := NewCharacterSetWithDefault()

	return NewNFA(cs)
}

func (nfa *NFA) Init() {
	nfa.InitState = nfa.getNewState()

	nfa.initMultiRune()
	nfa.initSingleRune()
	nfa.initIdentifier()
	nfa.initStringLiteral()
	nfa.initNumberLiteral()
}

func (nfa *NFA) initMultiRune() {
	for tokenType, tokenString := range MultiRuneList {
		start := nfa.getNewState()
		// temporary state
		tempState := start
		nfa.InitState.AddNext(Epsilon, start)

		for _, c := range tokenString {
			s := nfa.getNewState()
			tempState.AddNext(c, s)
			tempState = s
		}

		final := nfa.getFinalState(tokenType)
		tempState.AddNext(Epsilon, final)
	}
}

func (nfa *NFA) initIdentifier() {
	start := nfa.getNewState()
	nfa.InitState.AddNext(Epsilon, start)

	for _, c := range nfa.CharacterSet.GetDigits() {
		start.AddNext(c, start)
	}

	s := nfa.getNewState()
	for _, c := range nfa.CharacterSet.GetAlphabets() {
		start.AddNext(c, s)
	}
	for _, c := range nfa.CharacterSet.GetAlphabets() {
		s.AddNext(c, s)
	}

	final := nfa.getFinalState(Identifier)
	s.AddNext(Epsilon, final)
}

func (nfa *NFA) initSingleRune() {
	for tokenType, c := range SingleRuneList {
		start := nfa.getNewState()
		nfa.InitState.AddNext(Epsilon, start)

		s := nfa.getNewState()
		start.AddNext(c, s)

		final := nfa.getFinalState(tokenType)
		s.AddNext(Epsilon, final)
	}
}

func (nfa *NFA) initStringLiteral() {
	start := nfa.getNewState()
	nfa.InitState.AddNext(Epsilon, start)

	for _, c := range nfa.CharacterSet.GetAlphabets() {
		start.AddNext(c, start)
	}
	for _, c := range nfa.CharacterSet.GetDigits() {
		start.AddNext(c, start)
	}

	final := nfa.getFinalState(StringLiteral)
	start.AddNext(Epsilon, final)
}

func (nfa *NFA) initNumberLiteral() {
	start := nfa.getNewState()
	nfa.InitState.AddNext(Epsilon, start)
	s := nfa.getNewState()

	for _, c := range nfa.CharacterSet.GetDigits() {
		start.AddNext(c, s)
		s.AddNext(c, s)
	}

	final := nfa.getFinalState(NumberLiteral)
	s.AddNext(Epsilon, final)
}

func (nfa *NFA) Print() {
	nfa.InitState.Print()
}

func (nfa *NFA) Match() *Token {
	return nil
}

func (nfa *NFA) getNewState() *State {
	nfa.Index++
	return NewState(nfa.Index)
}

func (nfa *NFA) getFinalState(tokenType TokenType) *State {
	final := nfa.getNewState()
	final.IsFinal = true
	final.TokenType = tokenType

	return final
}
