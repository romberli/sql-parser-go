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
		return "selectKeyword"
	case From:
		return "fromKeyword"
	case As:
		return "asKeyword"
	case And:
		return "andKeyword"
	case Or:
		return "orKeyword"
	case Where:
		return "whereKeyword"
	case Identifier:
		return "identifier"
	case GE:
		return "greaterOrEqual"
	case GT:
		return "greaterThan"
	case LE:
		return "lessOrEqual"
	case LT:
		return "lessThan"
	case Equal:
		return "equal"
	case NotEqual1, NotEqual2:
		return "notEqual"
	case Plus:
		return "plus"
	case Minus:
		return "minus"
	case Multiply:
		return "multiply"
	case Divide:
		return "divide"
	case Mod:
		return "mod"
	case NumberLiteral:
		return "numberLiteral"
	case StringLiteral:
		return "stringLiteral"
	case LeftParenthesis:
		return "leftParenthesis"
	case RightParenthesis:
		return "rightParenthesis"
	case Comma:
		return "comma"
	case Semicolon:
		return "semicolon"
	case SingleQuote:
		return "singleQuote"
	case WhiteSpace:
		return "whiteSpace"
	case Epsilon:
		return "ε"
	case End:
		return "end"
	case Error:
		return "error"
	default:
		return "unknown"
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
	return fmt.Sprintf(`{tokenType: %s, lexeme: %s}`, t.Type.String(), t.Lexeme)
}
