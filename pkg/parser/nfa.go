package parser

import (
	"fmt"

	"github.com/pingcap/errors"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/sql-parser-go/pkg/ast"
	"github.com/romberli/sql-parser-go/pkg/token"
)

type NFA struct {
	Tokens    []*token.Token
	Index     int
	InitState *State
}

// NewNFA returns a new *NFA
func NewNFA(tokens []*token.Token) *NFA {
	nfa := &NFA{
		Tokens: tokens,
		Index:  -2,
	}

	nfa.init()

	return nfa
}

func (nfa *NFA) init() {
	start := nfa.getNewState()
	rootStart, rootEnd := nfa.parseRoot()
	start.AddNext(token.Epsilon, rootStart)
	final := nfa.getNewState()
	final.SetNode(ast.NewNodeWithDefault(ast.End))

	nfa.InitState = start
	rootEnd.AddNext(token.End, final)
}

func (nfa *NFA) Match() (*ast.Node, error) {
	err := nfa.match(nfa.InitState, constant.ZeroInt)
	if err != nil {
		return nil, err
	}

	return nfa.InitState.Next[token.Epsilon][constant.ZeroInt].Node, nil
}

func (nfa *NFA) match(s *State, i int) error {
	// if s.Node != nil {
	//     fmt.Println(fmt.Sprintf("matching %s started", s.Node.Type.String()))
	// }
	if i == len(nfa.Tokens) {
		//  all input tokens are matched correctly
		return nil
	}

	// newNode := ast
	// if s.Node != nil {
	//     // start to parse a new non-terminal
	//     newNode = ast.NewNodeWithDefault(s.Node.Type)
	//     ast.(newNode)
	// }
	if s.Node != nil && s.Parent != nil {
		s.Parent.AddChildren(s.Node)
	}

	t := nfa.Tokens[i]
	nsList := s.Next[t.Type]
	if nsList == nil {
		nsList = s.Next[token.Epsilon]
		if nsList == nil {
			// can't transit to any other state, return error
			if s.Node != nil {
				s.Parent.RemoveLastChild()
			}
			fmt.Println(fmt.Sprintf("matching %s failed", s.Node.Type.String()))
			return errors.Errorf("error when matching token. matched tokens: %s, next token: %s",
				nfa.Tokens[:i], nfa.Tokens[i])
		}
	} else {
		// matched a token, set the token
		nsList[constant.ZeroInt].Node.SetToken(t)
		// increasing the index
		i++
	}

	for _, ns := range nsList {
		err := nfa.match(ns, i)
		if err == nil {
			return nil
		}
	}

	matchedTokens := nfa.Tokens[:i]
	nextToken := nfa.Tokens[i]

	if s.Node != nil {
		s.Parent.RemoveLastChild()
	}
	fmt.Println(fmt.Sprintf("matching %s failed", s.Node.Type.String()))
	return errors.Errorf("error when matching token. matched tokens: %s, next token: %s",
		matchedTokens, nextToken)
}

func (nfa *NFA) Print() {
	nfa.InitState.Print()
}

// getNewState gets a new state
func (nfa *NFA) getNewState() *State {
	nfa.Index++
	return NewState(nfa.Index)
}

func (nfa *NFA) parseRoot() (*State, *State) {
	start := nfa.getNewState()
	rootNode := ast.NewNodeWithDefault(ast.Root)
	start.SetNode(rootNode)
	selectStatementStart, selectStatementEnd := nfa.parseSelectStatement(rootNode)
	statementTerminatorStart, statementTerminatorEnd := nfa.parseStatementTerminator(rootNode)
	end := nfa.getNewState()

	start.AddNext(token.Epsilon, selectStatementStart)
	selectStatementEnd.AddNext(token.Epsilon, statementTerminatorStart)
	selectStatementEnd.AddNext(token.Epsilon, end)
	statementTerminatorEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseSelectStatement(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	selectNode := ast.NewNodeWithDefault(ast.SelectStatement)
	start.SetNode(selectNode)
	start.SetParent(parent)
	selectKeyword := nfa.getNewState()
	selectKeyword.SetNode(ast.NewNodeWithDefault(ast.SelectKeyword))
	selectKeyword.SetParent(selectNode)
	columnListStart, columnListEnd := nfa.parseColumnList(selectNode)
	fromKeyword := nfa.getNewState()
	fromKeyword.SetNode(ast.NewNodeWithDefault(ast.FromKeyword))
	fromKeyword.SetParent(selectNode)
	tableNameStart, tableNameEnd := nfa.parseTableName(selectNode)
	whereClauseStart, whereClauseEnd := nfa.parseWhereClause(selectNode)
	end := nfa.getNewState()

	start.AddNext(token.Select, selectKeyword)
	selectKeyword.AddNext(token.Epsilon, columnListStart)
	columnListEnd.AddNext(token.From, fromKeyword)
	fromKeyword.AddNext(token.Epsilon, tableNameStart)
	tableNameEnd.AddNext(token.Epsilon, whereClauseStart)
	tableNameEnd.AddNext(token.Epsilon, end)
	whereClauseEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseColumnList(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	columnListNode := ast.NewNodeWithDefault(ast.ColumnList)
	start.SetNode(columnListNode)
	start.SetParent(parent)
	columnIdentifierStart, columnIdentifierEnd := nfa.parseColumnIdentifier(columnListNode)
	otherColumnsStart, otherColumnsEnd := nfa.parseOtherColumns(columnListNode)
	end := nfa.getNewState()

	start.AddNext(token.Epsilon, columnIdentifierStart)
	columnIdentifierEnd.AddNext(token.Epsilon, otherColumnsStart)
	columnIdentifierEnd.AddNext(token.Epsilon, end)
	otherColumnsEnd.AddNext(token.Epsilon, otherColumnsStart)
	otherColumnsEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseColumnIdentifier(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	columnIdentifierNode := ast.NewNodeWithDefault(ast.ColumnIdentifier)
	start.SetNode(columnIdentifierNode)
	start.SetParent(parent)
	columnWithAliasStart, columnWithAliasEnd := nfa.parseColumnWithAlias(columnIdentifierNode)
	end := nfa.getNewState()

	start.AddNext(token.Epsilon, columnWithAliasStart)
	columnWithAliasEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseOtherColumns(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	otherColumnsNode := ast.NewNodeWithDefault(ast.OtherColumns)
	start.SetNode(otherColumnsNode)
	start.SetParent(parent)
	commaOperator := nfa.getNewState()
	commaOperator.SetNode(ast.NewNodeWithDefault(ast.CommaOperator))
	commaOperator.SetParent(otherColumnsNode)
	columnWithAliasStart, columnWithAliasEnd := nfa.parseColumnWithAlias(otherColumnsNode)
	end := nfa.getNewState()

	start.AddNext(token.Comma, commaOperator)
	commaOperator.AddNext(token.Epsilon, columnWithAliasStart)
	columnWithAliasEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseColumnWithAlias(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	columnWithAliasNode := ast.NewNodeWithDefault(ast.ColumnWithAlias)
	start.SetNode(columnWithAliasNode)
	start.SetParent(parent)
	columnExpressionStart, columnExpressionEnd := nfa.parseColumnExpression(columnWithAliasNode)
	aliasNameStart, aliasNameEnd := nfa.parseAliasName(columnWithAliasNode)
	end := nfa.getNewState()

	start.AddNext(token.Epsilon, columnExpressionStart)
	columnExpressionEnd.AddNext(token.Epsilon, aliasNameStart)
	columnExpressionEnd.AddNext(token.Epsilon, end)
	aliasNameEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseColumnExpression(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	columnExpression := ast.NewNodeWithDefault(ast.ColumnExpression)
	start.SetNode(columnExpression)
	start.SetParent(parent)
	columnNameStart, columnNameEnd := nfa.parseColumnName(columnExpression)
	otherExpressionStart, otherExpressionEnd := nfa.parseOtherExpression(columnExpression)
	end := nfa.getNewState()

	start.AddNext(token.Epsilon, columnNameStart)
	columnNameEnd.AddNext(token.Epsilon, otherExpressionStart)
	columnNameEnd.AddNext(token.Epsilon, end)
	otherExpressionEnd.AddNext(token.Epsilon, otherExpressionStart)
	otherExpressionEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseColumnName(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	columnNameNode := ast.NewNodeWithDefault(ast.ColumnName)
	start.SetNode(columnNameNode)
	start.SetParent(parent)
	identifier := nfa.getNewState()
	identifier.SetNode(ast.NewNodeWithDefault(ast.Identifier))
	identifier.SetParent(columnNameNode)
	literalExpressionStart, literalExpressionEnd := nfa.parseLiteralExpression(columnNameNode)
	end := nfa.getNewState()

	start.AddNext(token.Identifier, identifier)
	start.AddNext(token.Epsilon, literalExpressionStart)
	identifier.AddNext(token.Epsilon, end)
	literalExpressionEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseOtherExpression(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	otherExpressionNode := ast.NewNodeWithDefault(ast.OtherExpression)
	start.SetNode(otherExpressionNode)
	start.SetParent(parent)
	expressionOperatorStart, expressionOperatorEnd := nfa.parseExpressionOperator(otherExpressionNode)
	columnNameStart, columnNameEnd := nfa.parseColumnName(otherExpressionNode)
	end := nfa.getNewState()

	start.AddNext(token.Epsilon, expressionOperatorStart)
	expressionOperatorEnd.AddNext(token.Epsilon, columnNameStart)
	columnNameEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseExpressionOperator(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	expressionOperatorNode := ast.NewNodeWithDefault(ast.ExpressionOperator)
	start.SetNode(expressionOperatorNode)
	start.SetParent(parent)
	plusOperator := nfa.getNewState()
	plusOperator.SetNode(ast.NewNodeWithDefault(ast.PlusOperator))
	plusOperator.SetParent(expressionOperatorNode)
	minusOperator := nfa.getNewState()
	minusOperator.SetNode(ast.NewNodeWithDefault(ast.MinusOperator))
	minusOperator.SetParent(expressionOperatorNode)
	end := nfa.getNewState()

	start.AddNext(token.Plus, plusOperator)
	plusOperator.AddNext(token.Epsilon, end)
	start.AddNext(token.Minus, minusOperator)
	minusOperator.AddNext(token.Epsilon, end)
	// start.AddChildren(token.Multiply, end)
	// start.AddChildren(token.Divide, end)

	return start, end
}

func (nfa *NFA) parseLiteralExpression(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	literalExpressionNode := ast.NewNodeWithDefault(ast.LiteralExpression)
	start.SetNode(literalExpressionNode)
	start.SetParent(parent)
	literalStart, literalEnd := nfa.parseLiteral(literalExpressionNode)
	otherLiteralStart, otherLiteralEnd := nfa.parseOtherLiteral(literalExpressionNode)
	end := nfa.getNewState()

	start.AddNext(token.Epsilon, literalStart)
	literalEnd.AddNext(token.Epsilon, otherLiteralStart)
	literalEnd.AddNext(token.Epsilon, end)
	otherLiteralEnd.AddNext(token.Epsilon, otherLiteralStart)
	otherLiteralEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseLiteral(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	literalNode := ast.NewNodeWithDefault(ast.Literal)
	start.SetNode(literalNode)
	start.SetParent(parent)
	numberLiteral := nfa.getNewState()
	numberLiteral.SetNode(ast.NewNodeWithDefault(ast.NumberLiteral))
	numberLiteral.SetParent(literalNode)
	stringLiteral := nfa.getNewState()
	stringLiteral.SetNode(ast.NewNodeWithDefault(ast.StringLiteral))
	stringLiteral.SetParent(literalNode)
	end := nfa.getNewState()

	start.AddNext(token.NumberLiteral, numberLiteral)
	numberLiteral.AddNext(token.Epsilon, end)
	start.AddNext(token.StringLiteral, stringLiteral)
	stringLiteral.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseOtherLiteral(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	otherLiteralNode := ast.NewNodeWithDefault(ast.OtherLiteral)
	start.SetNode(otherLiteralNode)
	start.SetParent(parent)
	expressionOperatorStart, expressionOperatorEnd := nfa.parseExpressionOperator(otherLiteralNode)
	literalStart, literalEnd := nfa.parseLiteral(otherLiteralNode)
	end := nfa.getNewState()

	start.AddNext(token.Epsilon, expressionOperatorStart)
	expressionOperatorEnd.AddNext(token.Epsilon, literalStart)
	literalEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseAliasName(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	aliasNameNode := ast.NewNodeWithDefault(ast.AliasName)
	start.SetNode(aliasNameNode)
	start.SetParent(parent)
	asKeyword := nfa.getNewState()
	asKeyword.SetNode(ast.NewNodeWithDefault(ast.AsKeyword))
	asKeyword.SetParent(aliasNameNode)
	identifier := nfa.getNewState()
	identifier.SetNode(ast.NewNodeWithDefault(ast.Identifier))
	identifier.SetParent(aliasNameNode)
	end := nfa.getNewState()

	start.AddNext(token.As, asKeyword)
	asKeyword.AddNext(token.Identifier, identifier)
	identifier.AddNext(token.Epsilon, end)
	start.AddNext(token.Identifier, identifier)

	return start, end
}

func (nfa *NFA) parseTableName(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	tableNameNode := ast.NewNodeWithDefault(ast.TableName)
	start.SetNode(tableNameNode)
	start.SetParent(parent)
	identifier := nfa.getNewState()
	identifier.SetNode(ast.NewNodeWithDefault(ast.Identifier))
	identifier.SetParent(tableNameNode)
	aliasNameStart, aliasNameEnd := nfa.parseAliasName(tableNameNode)
	end := nfa.getNewState()

	start.AddNext(token.Identifier, identifier)
	identifier.AddNext(token.Epsilon, aliasNameStart)
	identifier.AddNext(token.Epsilon, end)
	aliasNameEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseWhereClause(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	whereClauseNode := ast.NewNodeWithDefault(ast.WhereClause)
	start.SetNode(whereClauseNode)
	start.SetParent(parent)
	whereKeyword := nfa.getNewState()
	whereKeyword.SetNode(ast.NewNodeWithDefault(ast.WhereKeyword))
	whereKeyword.SetParent(whereClauseNode)
	columnComparisonStart, columnComparisonEnd := nfa.parseColumnComparison(whereClauseNode)
	otherColumnComparisonStart, otherColumnComparisonEnd := nfa.parseOtherColumnComparison(whereClauseNode)
	end := nfa.getNewState()

	start.AddNext(token.Where, whereKeyword)
	whereKeyword.AddNext(token.Epsilon, columnComparisonStart)
	columnComparisonEnd.AddNext(token.Epsilon, otherColumnComparisonStart)
	columnComparisonEnd.AddNext(token.Epsilon, end)
	otherColumnComparisonEnd.AddNext(token.Epsilon, otherColumnComparisonStart)
	otherColumnComparisonEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseColumnComparison(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	columnComparisonNode := ast.NewNodeWithDefault(ast.ColumnComparison)
	start.SetNode(columnComparisonNode)
	start.SetParent(parent)
	columnNameStart, columnNameEnd := nfa.parseColumnName(columnComparisonNode)
	otherColumnNameStart, otherColumnNameEnd := nfa.parseOtherColumnName(columnComparisonNode)
	end := nfa.getNewState()

	start.AddNext(token.Epsilon, columnNameStart)
	columnNameEnd.AddNext(token.Epsilon, otherColumnNameStart)
	columnNameEnd.AddNext(token.Epsilon, end)
	otherColumnNameEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseOtherColumnName(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	otherColumnNameNode := ast.NewNodeWithDefault(ast.OtherColumnName)
	start.SetNode(otherColumnNameNode)
	start.SetParent(parent)
	comparisonOperatorStart, comparisonOperatorEnd := nfa.parseComparisonOperator(otherColumnNameNode)
	columnNameStart, columnNameEnd := nfa.parseColumnName(otherColumnNameNode)
	end := nfa.getNewState()

	start.AddNext(token.Epsilon, comparisonOperatorStart)
	comparisonOperatorEnd.AddNext(token.Epsilon, columnNameStart)
	columnNameEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseOtherColumnComparison(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	otherColumnComparisonNode := ast.NewNodeWithDefault(ast.OtherColumnComparison)
	start.SetNode(otherColumnComparisonNode)
	start.SetParent(parent)
	whereOperatorStart, whereOperatorEnd := nfa.parseWhereOperator(otherColumnComparisonNode)
	columnComparisonStart, columnComparisonEnd := nfa.parseColumnComparison(otherColumnComparisonNode)
	end := nfa.getNewState()

	start.AddNext(token.Epsilon, whereOperatorStart)
	whereOperatorEnd.AddNext(token.Epsilon, columnComparisonStart)
	columnComparisonEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseComparisonOperator(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	comparisonOperatorNode := ast.NewNodeWithDefault(ast.ComparisonOperator)
	start.SetNode(comparisonOperatorNode)
	start.SetParent(parent)
	greaterOrEqualOperator := nfa.getNewState()
	greaterOrEqualOperator.SetNode(ast.NewNodeWithDefault(ast.GreaterOrEqualOperator))
	greaterOrEqualOperator.SetParent(comparisonOperatorNode)
	greaterThanOperator := nfa.getNewState()
	greaterThanOperator.SetNode(ast.NewNodeWithDefault(ast.GreaterThanOperator))
	greaterThanOperator.SetParent(comparisonOperatorNode)
	lessOrEqualOperator := nfa.getNewState()
	lessOrEqualOperator.SetNode(ast.NewNodeWithDefault(ast.LessOrEqualOperator))
	lessOrEqualOperator.SetParent(comparisonOperatorNode)
	lessThanOperator := nfa.getNewState()
	lessThanOperator.SetNode(ast.NewNodeWithDefault(ast.LessThanOperator))
	lessThanOperator.SetParent(comparisonOperatorNode)
	equalOperator := nfa.getNewState()
	equalOperator.SetNode(ast.NewNodeWithDefault(ast.EqualOperator))
	equalOperator.SetParent(comparisonOperatorNode)
	notEqual1Operator := nfa.getNewState()
	notEqual1Operator.SetNode(ast.NewNodeWithDefault(ast.NotEqual1Operator))
	notEqual1Operator.SetParent(comparisonOperatorNode)
	notEqual2Operator := nfa.getNewState()
	notEqual2Operator.SetNode(ast.NewNodeWithDefault(ast.NotEqual2Operator))
	notEqual2Operator.SetParent(comparisonOperatorNode)
	end := nfa.getNewState()

	start.AddNext(token.GE, greaterOrEqualOperator)
	greaterOrEqualOperator.AddNext(token.Epsilon, end)
	start.AddNext(token.GT, greaterThanOperator)
	greaterThanOperator.AddNext(token.Epsilon, end)
	start.AddNext(token.LE, lessOrEqualOperator)
	lessOrEqualOperator.AddNext(token.Epsilon, end)
	start.AddNext(token.LT, lessThanOperator)
	lessThanOperator.AddNext(token.Epsilon, end)
	start.AddNext(token.Equal, equalOperator)
	equalOperator.AddNext(token.Epsilon, end)
	start.AddNext(token.NotEqual1, notEqual1Operator)
	notEqual1Operator.AddNext(token.Epsilon, end)
	start.AddNext(token.NotEqual2, notEqual2Operator)
	notEqual2Operator.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseWhereOperator(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	whereOperatorNode := ast.NewNodeWithDefault(ast.WhereOperator)
	start.SetNode(whereOperatorNode)
	start.SetParent(parent)
	andKeyword := nfa.getNewState()
	andKeyword.SetNode(ast.NewNodeWithDefault(ast.AndKeyword))
	andKeyword.SetParent(whereOperatorNode)
	orKeyword := nfa.getNewState()
	orKeyword.SetNode(ast.NewNodeWithDefault(ast.OrKeyword))
	orKeyword.SetParent(whereOperatorNode)
	end := nfa.getNewState()

	start.AddNext(token.And, andKeyword)
	andKeyword.AddNext(token.Epsilon, end)
	start.AddNext(token.Or, orKeyword)
	orKeyword.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseStatementTerminator(parent *ast.Node) (*State, *State) {
	start := nfa.getNewState()
	statementTerminatorNode := ast.NewNodeWithDefault(ast.StatementTerminator)
	start.SetNode(statementTerminatorNode)
	start.SetParent(parent)
	semicolonOperator := nfa.getNewState()
	semicolonOperator.SetNode(ast.NewNodeWithDefault(ast.SemicolonOperator))
	semicolonOperator.SetParent(statementTerminatorNode)
	end := nfa.getNewState()

	start.AddNext(token.Semicolon, semicolonOperator)
	semicolonOperator.AddNext(token.Epsilon, end)

	return start, end
}
