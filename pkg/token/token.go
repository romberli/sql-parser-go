package token

import (
	"fmt"

	"github.com/romberli/go-util/constant"
)

type Type int

const (
	// keyword
	Select Type = iota + 1
	From
	As
	Where
	And
	Or
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
	// Epsilon
	Epsilon
	// End
	End
	// error
	Error
)

var (
	// epsilon
	EpsilonRune rune = constant.ZeroInt
	KeywordList      = []Type{Select, From, As, Where, And, Or}
)

// String returns the string representation of the token type
func (t Type) String() string {
	switch t {
	case Select:
		return "SelectKeyword"
	case From:
		return "FromKeyword"
	case As:
		return "AsKeyword"
	case And:
		return "AndKeyword"
	case Or:
		return "OrKeyword"
	case Where:
		return "WhereKeyword"
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
	case Comma:
		return "Comma"
	case Semicolon:
		return "Semicolon"
	case SingleQuote:
		return "SingleQuote"
	case WhiteSpace:
		return "WhiteSpace"
	case Epsilon:
		return "ε"
	case End:
		return "End"
	case Error:
		return "Error"
	default:
		return "Unknown"
	}
}

// IsKeyword returns if the token type is a keyword
func (t Type) IsKeyword() bool {
	for _, keyword := range KeywordList {
		if t == keyword {
			return true
		}
	}

	return false
}

type Token struct {
	Type   Type
	Lexeme string
}

// NewToken returns a new *Token
func NewToken(tokenType Type, lexeme string) *Token {
	return &Token{
		Type:   tokenType,
		Lexeme: lexeme,
	}
}

// String returns the string representation of the token
func (t *Token) String() string {
	return fmt.Sprintf("{tokenType: %s, lexeme: %s}", t.Type.String(), t.Lexeme)
}
