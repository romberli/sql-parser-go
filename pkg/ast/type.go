package ast

import (
	"github.com/romberli/sql-parser-go/pkg/token"
)

type Type int

const (
	// non-terminal
	Root Type = iota
	SelectStatement
	SimpleSelectStatement
	ColumnList
	TableName
	WhereClause
	ColumnIdentifier
	OtherColumns
	ColumnWithAlias
	ColumnExpression
	AliasName
	ColumnName
	OtherExpression
	ExpressionOperator
	LiteralExpression
	Literal
	OtherLiteral
	ColumnComparison
	OtherColumnName
	OtherColumnComparison
	WhereOperator
	ComparisonOperator
	StatementTerminator
	Epsilon
	// terminal
	SelectKeyword
	FromKeyword
	AsKeyword
	WhereKeyword
	AndKeyword
	OrKeyword
	Identifier
	StringLiteral
	NumberLiteral
	SemicolonOperator
	CommaOperator
	PlusOperator
	MinusOperator
	GreaterOrEqualOperator
	GreaterThanOperator
	LessOrEqualOperator
	LessThanOperator
	EqualOperator
	NotEqual1Operator
	NotEqual2Operator
	End
)

func (t Type) String() string {
	switch t {
	case Root:
		return "Root"
	case Epsilon:
		return "MayBeEpsilon"
	case SelectStatement:
		return "SelectStatement"
	case SimpleSelectStatement:
		return "SimpleSelectStatement"
	case ColumnList:
		return "ColumnList"
	case TableName:
		return "TableName"
	case WhereClause:
		return "WhereClause"
	case ColumnIdentifier:
		return "ColumnIdentifier"
	case OtherColumns:
		return "OtherColumns"
	case ColumnWithAlias:
		return "ColumnWithAlias"
	case ColumnExpression:
		return "ColumnExpression"
	case AliasName:
		return "AliasName"
	case ColumnName:
		return "ColumnName"
	case OtherExpression:
		return "OtherExpression"
	case ExpressionOperator:
		return "ExpressionOperator"
	case LiteralExpression:
		return "LiteralExpression"
	case Literal:
		return "Literal"
	case OtherLiteral:
		return "OtherLiteral"
	case ColumnComparison:
		return "ColumnComparison"
	case OtherColumnName:
		return "OtherColumnName"
	case OtherColumnComparison:
		return "OtherColumnComparison"
	case WhereOperator:
		return "WhereOperator"
	case ComparisonOperator:
		return "ComparisonOperator"
	case StatementTerminator:
		return "StatementTerminator"
	case SelectKeyword:
		return "selectKeyword"
	case FromKeyword:
		return "fromKeyword"
	case AsKeyword:
		return "asKeyword"
	case WhereKeyword:
		return "whereKeyword"
	case AndKeyword:
		return "andKeyword"
	case OrKeyword:
		return "orKeyword"
	case Identifier:
		return "identifier"
	case StringLiteral:
		return "stringLiteral"
	case NumberLiteral:
		return "numberLiteral"
	case SemicolonOperator:
		return "semicolonOperator"
	case CommaOperator:
		return "commaOperator"
	case PlusOperator:
		return "plusOperator"
	case MinusOperator:
		return "minusOperator"
	case GreaterOrEqualOperator:
		return "greaterOrEqualOperator"
	case GreaterThanOperator:
		return "greaterThanOperator"
	case LessOrEqualOperator:
		return "lessOrEqualOperator"
	case LessThanOperator:
		return "lessThanOperator"
	case EqualOperator:
		return "equalOperator"
	case NotEqual1Operator:
		return "notEqual1Operator"
	case NotEqual2Operator:
		return "notEqual2Operator"
	case End:
		return "end"
	default:
		return "Unknown"
	}
}

func (t Type) IsTerminal() bool {
	return t > Epsilon
}

func (t Type) GetTokenType() token.Type {
	if t.IsTerminal() {
		switch t {
		case SelectKeyword:
			return token.Select
		case FromKeyword:
			return token.From
		case AsKeyword:
			return token.As
		case WhereKeyword:
			return token.Where
		case AndKeyword:
			return token.And
		case OrKeyword:
			return token.Or
		case Identifier:
			return token.Identifier
		case StringLiteral:
			return token.StringLiteral
		case NumberLiteral:
			return token.NumberLiteral
		case SemicolonOperator:
			return token.Semicolon
		case CommaOperator:
			return token.Comma
		case PlusOperator:
			return token.Plus
		case MinusOperator:
			return token.Minus
		case GreaterOrEqualOperator:
			return token.GE
		case GreaterThanOperator:
			return token.GT
		case LessOrEqualOperator:
			return token.LE
		case LessThanOperator:
			return token.LT
		case EqualOperator:
			return token.Equal
		case NotEqual1Operator:
			return token.NotEqual1
		case NotEqual2Operator:
			return token.NotEqual2
		case End:
			return token.End
		}
	}

	return token.Error
}
