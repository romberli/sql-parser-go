package ast

import (
	"fmt"
	"strings"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/sql-parser-go/pkg/token"
)

type Node struct {
	Type  Type
	Token *token.Token
	Min   int
	Max   int

	Children []*Node
}

func NewNode(nodeType Type, min, max int) *Node {
	return newNode(nodeType, min, max)
}

func NewNodeWithDefault(nodeType Type) *Node {
	return newNode(nodeType, 1, 1)
}

func newNode(nodeType Type, min, max int) *Node {
	node := &Node{
		Type: nodeType,
		Min:  min,
		Max:  max,
	}

	// ast.init()

	return node
}

func (n *Node) IsTerminal() bool {
	return n.Type > Epsilon
}

func (n *Node) MayEpsilon() bool {
	return n.Min == constant.ZeroInt
}

func (n *Node) GetChildren() [][]*Node {
	if n.IsTerminal() {
		return nil
	}

	switch n.Type {
	case Root:
		return [][]*Node{{
			NewNode(SelectStatement, 1, 1),
			NewNode(StatementTerminator, 0, 1),
		}}
	case SelectStatement:
		return [][]*Node{{
			NewNode(SelectKeyword, 1, 1),
			NewNode(ColumnList, 1, 1),
			NewNode(FromKeyword, 1, 1),
			NewNode(TableName, 1, 1),
			NewNode(WhereClause, 0, 1),
		}}
	case ColumnList:
		return [][]*Node{{
			NewNode(ColumnIdentifier, 1, 1),
			NewNode(OtherColumns, 0, -1),
		}}
	case ColumnIdentifier:
		return [][]*Node{{
			NewNode(ColumnWithAlias, 1, 1),
		}}
	case OtherColumns:
		return [][]*Node{{
			NewNode(CommaOperator, 1, 1),
			NewNode(ColumnWithAlias, 1, 1),
		}}
	case ColumnWithAlias:
		return [][]*Node{{
			NewNode(ColumnExpression, 1, 1),
			NewNode(AliasName, 0, 1),
		}}
	case ColumnExpression:
		return [][]*Node{{
			NewNode(ColumnName, 1, 1),
			NewNode(OtherExpression, 0, -1),
		}}
	case ColumnName:
		return [][]*Node{
			{NewNode(Identifier, 1, 1)},
			{NewNode(LiteralExpression, 1, 1)},
		}
	case OtherExpression:
		return [][]*Node{{
			NewNode(ExpressionOperator, 1, 1),
			NewNode(ColumnName, 1, 1),
		}}
	case ExpressionOperator:
		return [][]*Node{
			{NewNode(PlusOperator, 1, 1)},
			{NewNode(MinusOperator, 1, 1)},
		}
	case LiteralExpression:
		return [][]*Node{{
			NewNode(Literal, 1, 1),
			NewNode(OtherLiteral, 0, -1),
		}}
	case Literal:
		return [][]*Node{
			{NewNode(NumberLiteral, 1, 1)},
			{NewNode(StringLiteral, 1, 1)},
		}
	case OtherLiteral:
		return [][]*Node{{
			NewNode(ExpressionOperator, 1, 1),
			NewNode(Literal, 1, 1),
		}}
	case AliasName:
		return [][]*Node{
			{
				NewNode(AsKeyword, 1, 1),
				NewNode(Identifier, 1, 1),
			},
			{
				NewNode(Identifier, 1, 1),
			},
		}
	case TableName:
		return [][]*Node{{
			NewNode(Identifier, 1, 1),
			NewNode(AliasName, 0, 1),
		}}
	case WhereClause:
		return [][]*Node{{
			NewNode(WhereKeyword, 1, 1),
			NewNode(ColumnComparison, 1, 1),
			NewNode(OtherColumnComparison, 0, -1),
		}}
	case ColumnComparison:
		return [][]*Node{{
			NewNode(ColumnName, 1, 1),
			NewNode(OtherColumnName, 0, 1),
		}}
	case OtherColumnName:
		return [][]*Node{{
			NewNode(ComparisonOperator, 1, 1),
			NewNode(ColumnName, 1, 1),
		}}
	case OtherColumnComparison:
		return [][]*Node{{
			NewNode(WhereOperator, 1, 1),
			NewNode(ColumnComparison, 1, 1),
		}}
	case ComparisonOperator:
		return [][]*Node{
			{NewNode(GreaterOrEqualOperator, 1, 1)},
			{NewNode(GreaterThanOperator, 1, 1)},
			{NewNode(LessOrEqualOperator, 1, 1)},
			{NewNode(LessThanOperator, 1, 1)},
			{NewNode(EqualOperator, 1, 1)},
			{NewNode(NotEqual1Operator, 1, 1)},
			{NewNode(NotEqual2Operator, 1, 1)},
		}
	case WhereOperator:
		return [][]*Node{
			{NewNode(AndKeyword, 1, 1)},
			{NewNode(OrKeyword, 1, 1)},
		}
	case StatementTerminator:
		return [][]*Node{
			{NewNode(SemicolonOperator, 1, 1)},
		}
	}

	return nil
}

func (n *Node) GetFirstSet() []token.Type {
	var tokenTypeList []token.Type

	if n.IsTerminal() {
		if !token.TypeExists(tokenTypeList, n.Type.GetTokenType()) {
			tokenTypeList = append(tokenTypeList, n.Type.GetTokenType())
		}
	} else {
		for _, children := range n.GetChildren() {
			for _, t := range children[constant.ZeroInt].GetFirstSet() {
				if !token.TypeExists(tokenTypeList, t) {
					tokenTypeList = append(tokenTypeList, t)
				}
			}
		}

		if n.MayEpsilon() {
			followSet := n.GetFollowSet()
			// if followSet == nil {
			//     if !token.TypeExists(tokenTypeList, token.Epsilon) {
			//         tokenTypeList = append(tokenTypeList, token.Epsilon)
			//     }
			// } else {
			//     for _, t := range followSet {
			//         if !token.TypeExists(tokenTypeList, t) {
			//             tokenTypeList = append(tokenTypeList, t)
			//         }
			//     }
			// }
			for _, t := range followSet {
				if !token.TypeExists(tokenTypeList, t) {
					tokenTypeList = append(tokenTypeList, t)
				}
			}
		}
	}

	return tokenTypeList
}

func (n *Node) GetFollowSet() []token.Type {
	var tokenTypeList []token.Type

	switch n.Type {
	case SelectStatement:
		tokenTypeList = append(tokenTypeList, NewNode(StatementTerminator, 0, 1).GetFirstSet()...)
		tokenTypeList = append(tokenTypeList, token.End)
	case ColumnList:
		tokenTypeList = append(tokenTypeList, token.From)
	case ColumnIdentifier:
		tokenTypeList = append(tokenTypeList, NewNode(OtherColumns, 0, -1).GetFirstSet()...)
	case ColumnExpression:
		tokenTypeList = append(tokenTypeList, NewNode(AliasName, 0, 1).GetFirstSet()...)
	case ColumnName:
		tokenTypeList = append(tokenTypeList, NewNode(OtherExpression, 0, -1).GetFirstSet()...)
	case ExpressionOperator:
		tokenTypeList = append(tokenTypeList, NewNode(ColumnName, 1, 1).GetFirstSet()...)
		tokenTypeList = append(tokenTypeList, NewNode(LiteralExpression, 1, 1).GetFirstSet()...)
	case Literal:
		tokenTypeList = append(tokenTypeList, NewNode(OtherLiteral, 0, -1).GetFirstSet()...)
	case TableName:
		tokenTypeList = append(tokenTypeList, NewNode(WhereClause, 0, 1).GetFirstSet()...)
	case ColumnComparison:
		tokenTypeList = append(tokenTypeList, NewNode(OtherColumnComparison, 0, -1).GetFirstSet()...)
	case ComparisonOperator:
		tokenTypeList = append(tokenTypeList, NewNode(ColumnName, 1, 1).GetFirstSet()...)
	case WhereOperator:
		tokenTypeList = append(tokenTypeList, NewNode(ColumnComparison, 1, 1).GetFirstSet()...)
	default:
		// do nothing here because this ast is terminal
		// or there is no following tokens in the production rule that this ast is used
		//
		// OtherColumns
		// ColumnWithAlias
		// OtherExpression
		// LiteralExpression
		// OtherLiteral
		// AliasName
		// WhereClause
		// OtherColumnName
		// OtherColumnComparison
		// StatementTerminator
	}

	return tokenTypeList
}

func (n *Node) SetToken(t *token.Token) {
	n.Token = t
}

func (n *Node) SetRepeatTime(min, max int) {
	n.Min = min
	n.Max = max
}

func (n *Node) AddChildren(next *Node) {
	n.Children = append(n.Children, next)
}

func (n *Node) RemoveLastChild() {
	if len(n.Children) > constant.ZeroInt {
		n.Children = n.Children[:len(n.Children)-1]
	}
}

func (n *Node) PrintChildren() {
	n.printChildren(0, 2)
}

func (n *Node) printChildren(i, j int) {
	prefix := fmt.Sprintf("%s|%s", strings.Repeat(constant.SpaceString, i), strings.Repeat(constant.DashString, j))

	if i == constant.ZeroInt {
		prefix = constant.EmptyString
	}

	nodeString := fmt.Sprintf("%s%s", prefix, n.Type.String())
	if n.Token != nil {
		nodeString += fmt.Sprintf("(%s)", n.Token.String())
	}
	fmt.Println(nodeString)
	i += 3
	for _, child := range n.Children {
		// j += 2
		child.printChildren(i, j)
	}
}
