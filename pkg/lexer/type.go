package lexer

type TokenType int

const (
	// keyword
	Select TokenType = iota
	From
	As
	Where
	And
	// identifier
	Identifier
	// comparison operator
	GE
	GT
	LE
	LT
	Equal
	NotEqual1
	NotEqual2
	// arithmetic operator
	Plus
	Minus
	Multiply
	Divide
	Mod
	// number literal
	NumberLiteral
	// string literal
	StringLiteral
	// separator
	Comma
	Semicolon
	LeftParenthesis
	RightParenthesis
	SingleQuote
	// white space
	WhiteSpace
)

func (tt TokenType) String() string {
	switch tt {
	case Select:
		return "SelectKeyword"
	case From:
		return "FromKeyword"
	case As:
		return "AsKeyword"
	case Where:
		return "WhereKeyword"
	case And:
		return "AndKeyword"
	case Identifier:
		return "Identifier"
	case GE, GT, LE, LT, Equal, NotEqual1, NotEqual2:
		return "ComparisonOperator"
	case Plus, Minus, Multiply, Divide, Mod:
		return "ArithmeticOperator"
	case NumberLiteral:
		return "NumberLiteral"
	case StringLiteral:
		return "StringLiteral"
	case Comma, Semicolon, LeftParenthesis, RightParenthesis, SingleQuote:
		return "Separator"
	case WhiteSpace:
		return "WhiteSpace"
	default:
		return "Unknown"
	}
}

type State struct {
	Index     int
	Value     string
	Next      map[rune][]*State
	IsFinal   bool
	TokenType TokenType
}

func NewState(i int) *State {
	return &State{
		Index: i,
	}
}

func (s *State) AddNext(c rune, ns *State) {
	s.Next[c] = append(s.Next[c], ns)
}

func (s *State) Transit(c rune) []*State {
	return s.Next[c]
}
