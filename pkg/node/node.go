package node

import (
	"fmt"
	"strings"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/sql-parser-go/pkg/token"
)

type Node struct {
	Type         Type
	MayBeEpsilon bool
	Token        *token.Token

	Children  [][]*Node
	FirstSet  []token.Type
	FollowSet []token.Type

	Next []*Node
}

func NewNode(nodeType Type, mayBeEpsilon bool) *Node {
	return newNode(nodeType, mayBeEpsilon)
}

func NewNodeWithDefault(nodeType Type) *Node {
	return newNode(nodeType, false)
}

func newNode(nodeType Type, mayBeEpsilon bool) *Node {
	node := &Node{
		Type:         nodeType,
		MayBeEpsilon: mayBeEpsilon,
	}

	// node.init()

	return node
}

func (n *Node) IsTerminal() bool {
	return n.Type > Epsilon
}

func (n *Node) init() {
	n.Children = n.GetChildren()
	// n.FirstSet = n.GetFirstSet()
	// n.FollowSet = n.GetFollowSet()
}

func (n *Node) GetChildren() [][]*Node {
	switch n.Type {
	case Root:
		return [][]*Node{{
			NewNode(SelectStatement, false),
			NewNode(SemicolonOperator, true),
		}}
	case SelectStatement:
		return [][]*Node{{
			NewNode(SelectKeyword, false),
			NewNode(ColumnList, false),
			NewNode(FromKeyword, false),
			NewNode(TableName, false),
			NewNode(WhereClause, true),
		}}
	case ColumnList:
		return [][]*Node{{
			NewNode(ColumnIdentifier, false),
			NewNode(OtherColumns, true),
		}}
	case ColumnIdentifier:
		return [][]*Node{{
			NewNode(ColumnWithAlias, false),
		}}
	case OtherColumns:
		return [][]*Node{{
			NewNode(CommaOperator, false),
			NewNode(ColumnWithAlias, false),
		}}
	case ColumnWithAlias:
		return [][]*Node{{
			NewNode(ColumnExpression, false),
			NewNode(AliasName, true),
		}}
	case ColumnExpression:
		return [][]*Node{{
			NewNode(ColumnName, false),
			NewNode(OtherExpression, true),
		}}
	case ColumnName:
		return [][]*Node{{
			NewNode(ExpressionOperator, false),
			NewNode(ColumnName, false),
		}}
	case OtherExpression:
		return [][]*Node{{
			NewNode(ExpressionOperator, false),
			NewNode(ColumnName, false),
		}}
	case ExpressionOperator:
		return [][]*Node{{
			NewNode(PlusOperator, false),
			NewNode(MinusOperator, false),
		}}
	case LiteralExpression:
		return [][]*Node{{
			NewNode(Literal, false),
			NewNode(OtherLiteral, true),
		}}
	case Literal:
		return [][]*Node{{
			NewNode(NumberLiteral, false),
			NewNode(StringLiteral, false),
		}}
	case OtherLiteral:
		return [][]*Node{{
			NewNode(ExpressionOperator, false),
			NewNode(Literal, false),
		}}
	case AliasName:
		return [][]*Node{
			{
				NewNode(AsKeyword, false),
				NewNode(Identifier, false),
			},
			{
				NewNode(Identifier, false),
			},
		}
	case TableName:
		return [][]*Node{{
			NewNode(Identifier, false),
			NewNode(AliasName, true),
		}}
	case WhereClause:
		return [][]*Node{{
			NewNode(WhereKeyword, false),
			NewNode(ColumnComparison, false),
			NewNode(OtherColumnComparison, true),
		}}
	case ColumnComparison:
		return [][]*Node{{
			NewNode(ColumnName, false),
			NewNode(OtherColumnName, true),
		}}
	case OtherColumnName:
		return [][]*Node{{
			NewNode(ComparisonOperator, false),
			NewNode(ColumnName, false),
		}}
	case OtherColumnComparison:
		return [][]*Node{{
			NewNode(WhereOperator, false),
			NewNode(ColumnComparison, false),
		}}
	case ComparisonOperator:
		return [][]*Node{
			{NewNode(GreaterOrEqualOperator, false)},
			{NewNode(GreaterThanOperator, false)},
			{NewNode(LessOrEqualOperator, false)},
			{NewNode(LessThanOperator, false)},
			{NewNode(EqualOperator, false)},
			{NewNode(NotEqual1Operator, false)},
			{NewNode(NotEqual2Operator, false)},
		}
	case WhereOperator:
		return [][]*Node{
			{NewNode(AndKeyword, false)},
			{NewNode(OrKeyword, false)},
		}
	case StatementTerminator:
		return [][]*Node{
			{NewNode(SemicolonOperator, false)},
		}
	}

	return nil
}

func (n *Node) GetFirstSet() []token.Type {
	var tokenTypeList []token.Type

	if n.IsTerminal() {
		tokenTypeList = append(tokenTypeList, n.Type.GetTokenType())
	} else {
		for _, children := range n.Children {
			tokenTypeList = append(tokenTypeList, children[constant.ZeroInt].GetFirstSet()...)
		}

		if n.MayBeEpsilon {
			tokenTypeList = append(tokenTypeList, n.GetFollowSet()...)
		}
	}

	return tokenTypeList
}

func (n *Node) GetFollowSet() []token.Type {
	var tokenTypeList []token.Type

	switch n.Type {
	case SelectStatement:
		tokenTypeList = append(tokenTypeList, NewNode(StatementTerminator, true).GetFirstSet()...)
		tokenTypeList = append(tokenTypeList, token.End)
	case ColumnList:
		tokenTypeList = append(tokenTypeList, token.From)
	case ColumnIdentifier:
		tokenTypeList = append(tokenTypeList, NewNode(OtherColumns, true).GetFirstSet()...)
	case ColumnExpression:
		tokenTypeList = append(tokenTypeList, NewNode(AliasName, true).GetFirstSet()...)
	case ColumnName:
		tokenTypeList = append(tokenTypeList, NewNode(OtherExpression, true).GetFirstSet()...)
	case ExpressionOperator:
		tokenTypeList = append(tokenTypeList, NewNode(ColumnName, false).GetFirstSet()...)
		tokenTypeList = append(tokenTypeList, NewNode(LiteralExpression, false).GetFirstSet()...)
	case Literal:
		tokenTypeList = append(tokenTypeList, NewNode(OtherLiteral, true).GetFirstSet()...)
	case TableName:
		tokenTypeList = append(tokenTypeList, NewNode(WhereClause, true).GetFirstSet()...)
	case ColumnComparison:
		tokenTypeList = append(tokenTypeList, NewNode(OtherColumnComparison, true).GetFirstSet()...)
	case ComparisonOperator:
		tokenTypeList = append(tokenTypeList, NewNode(ColumnName, false).GetFirstSet()...)
	case WhereOperator:
		tokenTypeList = append(tokenTypeList, NewNode(ColumnComparison, false).GetFirstSet()...)
	default:
		// do nothing here because this node is terminal
		// or there is no following tokens in the production rule that this node is used
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

func (n *Node) AddNext(next *Node) {
	n.Next = append(n.Next, next)
}

func (n *Node) RemoveNext(i int) {
	if len(n.Next) > i {
		n.Next = append(n.Next[:i], n.Next[i+1:]...)
	}
}

func (n *Node) RemoveLast() {
	if len(n.Next) > constant.ZeroInt {
		n.Next = n.Next[:len(n.Next)-1]
	}
}

func (n *Node) PrintNext() {
	n.printNext(0, 2)
}

func (n *Node) printNext(i, j int) {
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
	for _, child := range n.Next {
		// j += 2
		child.printNext(i, j)
	}
}
