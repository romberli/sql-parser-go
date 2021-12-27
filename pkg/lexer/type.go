package lexer

import (
	"fmt"

	"github.com/romberli/go-util/constant"
)

type TokenType int

const (
	// keyword
	Select TokenType = iota + 1
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

var (
	// epsilon
	Epsilon rune = constant.ZeroInt
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
	case GE:
		return "GreaterOrEqual"
	case GT:
		return "GreaterThan"
	case LE:
		return "LessOrEqual"
	case LT:
		return "LessThan"
	case Equal:
		return "Equal"
	case NotEqual1, NotEqual2:
		return "NotEqual"
	case Plus:
		return "Plus"
	case Minus:
		return "Minus"
	case Multiply:
		return "Multiply"
	case Divide:
		return "Divide"
	case Mod:
		return "Mod"
	case NumberLiteral:
		return "NumberLiteral"
	case StringLiteral:
		return "StringLiteral"
	case LeftParenthesis:
		return "LeftParenthesis"
	case RightParenthesis:
		return "RightParenthesis"
	case Comma, Semicolon, SingleQuote:
		return "Separator"
	case WhiteSpace:
		return "WhiteSpace"
	default:
		return "Unknown"
	}
}

type Token struct {
	TokenType TokenType
	Lexeme    string
}

func NewToken(tokenType TokenType, lexeme string) *Token {
	return &Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("{tokenType: %s, lexeme: %s", t.TokenType.String(), t.Lexeme)
}

type State struct {
	Index     int
	Value     []rune
	Next      map[rune][]*State
	IsFinal   bool
	TokenType TokenType
}

func NewState(i int) *State {
	return &State{
		Index: i,
		Next:  make(map[rune][]*State),
	}
}

func (s *State) AppendValue(c rune) {
	s.Value = append(s.Value, c)
}

func (s *State) AddNext(c rune, ns *State) {
	s.Next[c] = append(s.Next[c], ns)
}

func (s *State) Transit(c rune) []*State {
	return s.Next[c]
}

func (s *State) Print() {
	printedList := make(map[int]*State)

	if s.IsFinal {
		fmt.Println(fmt.Sprintf(
			"final state found. index: %d, tokenType: %s",
			s.Index, s.TokenType.String()))
		return
	}

	for c, nsList := range s.Next {
		for _, ns := range nsList {
			printChar := c
			if c == constant.ZeroInt {
				printChar = EpsilonRune
			}
			fmt.Println(fmt.Sprintf("state %d + intput '%c' -> state %d", s.Index, printChar, ns.Index))
			printedList[ns.Index] = ns
		}
	}
	for _, ns := range printedList {
		if s.Index != ns.Index {
			ns.Print()
		}
	}
}
