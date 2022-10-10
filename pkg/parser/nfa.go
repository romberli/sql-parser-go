package parser

import (
	"fmt"

	"github.com/pingcap/errors"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/sql-parser-go/pkg/ast"
	"github.com/romberli/sql-parser-go/pkg/node"
	"github.com/romberli/sql-parser-go/pkg/token"
)

type NFA struct {
	Tokens    []*token.Token
	Index     int
	InitState *State
	tree      *ast.Tree
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
	final.SetNode(node.NewNodeWithDefault(node.End))

	nfa.InitState = start
	rootEnd.AddNext(token.End, final)
}

func (nfa *NFA) Match() (*node.Node, error) {
	err := nfa.match(nfa.InitState, constant.ZeroInt)
	if err != nil {
		return nil, err
	}

	return nfa.InitState.Next[token.Epsilon][constant.ZeroInt].Node, nil
}

func (nfa *NFA) match(s *State, i int) error {
	if s.Node != nil {
		fmt.Println(fmt.Sprintf("matching %s started", s.Node.Type.String()))
	}
	if i == len(nfa.Tokens) {
		//  all input tokens are matched correctly
		return nil
	}

	// newNode := node
	// if s.Node != nil {
	//     // start to parse a new non-terminal
	//     newNode = ast.NewNodeWithDefault(s.Node.Type)
	//     node.(newNode)
	// }
	if s.Node != nil && s.Parent != nil {
		s.Parent.AddNext(s.Node)
	}

	t := nfa.Tokens[i]
	nsList := s.Next[t.Type]
	if nsList == nil {
		nsList = s.Next[token.Epsilon]
		if nsList == nil {
			// can't transit to any other state, return error
			if s.Node != nil {
				s.Parent.RemoveLast()
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
		if err != nil {
			// if s.Node != nil {
			//     s.Parent.RemoveLast()
			// }
		} else {
			return nil
		}
	}

	matchedTokens := nfa.Tokens[:i]
	nextToken := nfa.Tokens[i]

	if s.Node != nil {
		s.Parent.RemoveLast()
	}
	fmt.Println(fmt.Sprintf("matching %s failed", s.Node.Type.String()))
	return errors.Errorf("error when matching token. matched tokens: %s, next token: %s",
		matchedTokens, nextToken)
}

func (nfa *NFA) Print() {
	nfa.InitState.Print()
}

func (nfa *NFA) lookAhead() *token.Token {
	// if nfa.Index+1 >= len(nfa.Tokens) {
	//     return nil
	// }

	return nfa.Tokens[nfa.Index+1]
}

func (nfa *NFA) readNext() *token.Token {
	// if nfa.Index+1 >= len(nfa.Tokens) {
	//     return nil
	// }
	nfa.Index++

	return nfa.Tokens[nfa.Index]
}

// getNewState gets a new state
func (nfa *NFA) getNewState() *State {
	nfa.Index++
	return NewState(nfa.Index)
}

func (nfa *NFA) parseRoot() (*State, *State) {
	start := nfa.getNewState()
	rootNode := node.NewNodeWithDefault(node.Root)
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

func (nfa *NFA) parseSelectStatement(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	selectNode := node.NewNodeWithDefault(node.SelectStatement)
	start.SetNode(selectNode)
	start.SetParent(parent)
	selectKeyword := nfa.getNewState()
	selectKeyword.SetNode(node.NewNodeWithDefault(node.SelectKeyword))
	selectKeyword.SetParent(selectNode)
	columnListStart, columnListEnd := nfa.parseColumnList(selectNode)
	fromKeyword := nfa.getNewState()
	fromKeyword.SetNode(node.NewNodeWithDefault(node.FromKeyword))
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

func (nfa *NFA) parseColumnList(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	columnListNode := node.NewNodeWithDefault(node.ColumnList)
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

func (nfa *NFA) parseColumnIdentifier(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	columnIdentifierNode := node.NewNodeWithDefault(node.ColumnIdentifier)
	start.SetNode(columnIdentifierNode)
	start.SetParent(parent)
	columnWithAliasStart, columnWithAliasEnd := nfa.parseColumnWithAlias(columnIdentifierNode)
	end := nfa.getNewState()

	start.AddNext(token.Epsilon, columnWithAliasStart)
	columnWithAliasEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseOtherColumns(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	otherColumnsNode := node.NewNodeWithDefault(node.OtherColumns)
	start.SetNode(otherColumnsNode)
	start.SetParent(parent)
	commaOperator := nfa.getNewState()
	commaOperator.SetNode(node.NewNodeWithDefault(node.CommaOperator))
	commaOperator.SetParent(otherColumnsNode)
	columnWithAliasStart, columnWithAliasEnd := nfa.parseColumnWithAlias(otherColumnsNode)
	end := nfa.getNewState()

	start.AddNext(token.Comma, commaOperator)
	commaOperator.AddNext(token.Epsilon, columnWithAliasStart)
	columnWithAliasEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseColumnWithAlias(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	columnWithAliasNode := node.NewNodeWithDefault(node.ColumnWithAlias)
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

func (nfa *NFA) parseColumnExpression(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	columnExpression := node.NewNodeWithDefault(node.ColumnExpression)
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

func (nfa *NFA) parseColumnName(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	columnNameNode := node.NewNodeWithDefault(node.ColumnName)
	start.SetNode(columnNameNode)
	start.SetParent(parent)
	identifier := nfa.getNewState()
	identifier.SetNode(node.NewNodeWithDefault(node.Identifier))
	identifier.SetParent(columnNameNode)
	literalExpressionStart, literalExpressionEnd := nfa.parseLiteralExpression(columnNameNode)
	end := nfa.getNewState()

	start.AddNext(token.Identifier, identifier)
	start.AddNext(token.Epsilon, literalExpressionStart)
	identifier.AddNext(token.Epsilon, end)
	literalExpressionEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseOtherExpression(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	otherExpressionNode := node.NewNodeWithDefault(node.OtherExpression)
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

func (nfa *NFA) parseExpressionOperator(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	expressionOperatorNode := node.NewNodeWithDefault(node.ExpressionOperator)
	start.SetNode(expressionOperatorNode)
	start.SetParent(parent)
	plusOperator := nfa.getNewState()
	plusOperator.SetNode(node.NewNodeWithDefault(node.PlusOperator))
	plusOperator.SetParent(expressionOperatorNode)
	minusOperator := nfa.getNewState()
	minusOperator.SetNode(node.NewNodeWithDefault(node.MinusOperator))
	minusOperator.SetParent(expressionOperatorNode)
	end := nfa.getNewState()

	start.AddNext(token.Plus, plusOperator)
	plusOperator.AddNext(token.Epsilon, end)
	start.AddNext(token.Minus, minusOperator)
	minusOperator.AddNext(token.Epsilon, end)
	// start.AddNext(token.Multiply, end)
	// start.AddNext(token.Divide, end)

	return start, end
}

func (nfa *NFA) parseLiteralExpression(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	literalExpressionNode := node.NewNodeWithDefault(node.LiteralExpression)
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

func (nfa *NFA) parseLiteral(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	literalNode := node.NewNodeWithDefault(node.Literal)
	start.SetNode(literalNode)
	start.SetParent(parent)
	numberLiteral := nfa.getNewState()
	numberLiteral.SetNode(node.NewNodeWithDefault(node.NumberLiteral))
	numberLiteral.SetParent(literalNode)
	stringLiteral := nfa.getNewState()
	stringLiteral.SetNode(node.NewNodeWithDefault(node.StringLiteral))
	stringLiteral.SetParent(literalNode)
	end := nfa.getNewState()

	start.AddNext(token.NumberLiteral, numberLiteral)
	numberLiteral.AddNext(token.Epsilon, end)
	start.AddNext(token.StringLiteral, stringLiteral)
	stringLiteral.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseOtherLiteral(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	otherLiteralNode := node.NewNodeWithDefault(node.OtherLiteral)
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

func (nfa *NFA) parseAliasName(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	aliasNameNode := node.NewNodeWithDefault(node.AliasName)
	start.SetNode(aliasNameNode)
	start.SetParent(parent)
	asKeyword := nfa.getNewState()
	asKeyword.SetNode(node.NewNodeWithDefault(node.AsKeyword))
	asKeyword.SetParent(aliasNameNode)
	identifier := nfa.getNewState()
	identifier.SetNode(node.NewNodeWithDefault(node.Identifier))
	identifier.SetParent(aliasNameNode)
	end := nfa.getNewState()

	start.AddNext(token.As, asKeyword)
	asKeyword.AddNext(token.Identifier, identifier)
	identifier.AddNext(token.Epsilon, end)
	start.AddNext(token.Identifier, identifier)

	return start, end
}

func (nfa *NFA) parseTableName(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	tableNameNode := node.NewNodeWithDefault(node.TableName)
	start.SetNode(tableNameNode)
	start.SetParent(parent)
	identifier := nfa.getNewState()
	identifier.SetNode(node.NewNodeWithDefault(node.Identifier))
	identifier.SetParent(tableNameNode)
	aliasNameStart, aliasNameEnd := nfa.parseAliasName(tableNameNode)
	end := nfa.getNewState()

	start.AddNext(token.Identifier, identifier)
	identifier.AddNext(token.Epsilon, aliasNameStart)
	identifier.AddNext(token.Epsilon, end)
	aliasNameEnd.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseWhereClause(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	whereClauseNode := node.NewNodeWithDefault(node.WhereClause)
	start.SetNode(whereClauseNode)
	start.SetParent(parent)
	whereKeyword := nfa.getNewState()
	whereKeyword.SetNode(node.NewNodeWithDefault(node.WhereKeyword))
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

func (nfa *NFA) parseColumnComparison(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	columnComparisonNode := node.NewNodeWithDefault(node.ColumnComparison)
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

func (nfa *NFA) parseOtherColumnName(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	otherColumnNameNode := node.NewNodeWithDefault(node.OtherColumnName)
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

func (nfa *NFA) parseOtherColumnComparison(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	otherColumnComparisonNode := node.NewNodeWithDefault(node.OtherColumnComparison)
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

func (nfa *NFA) parseComparisonOperator(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	comparisonOperatorNode := node.NewNodeWithDefault(node.ComparisonOperator)
	start.SetNode(comparisonOperatorNode)
	start.SetParent(parent)
	greaterOrEqualOperator := nfa.getNewState()
	greaterOrEqualOperator.SetNode(node.NewNodeWithDefault(node.GreaterOrEqualOperator))
	greaterOrEqualOperator.SetParent(comparisonOperatorNode)
	greaterThanOperator := nfa.getNewState()
	greaterThanOperator.SetNode(node.NewNodeWithDefault(node.GreaterThanOperator))
	greaterThanOperator.SetParent(comparisonOperatorNode)
	lessOrEqualOperator := nfa.getNewState()
	lessOrEqualOperator.SetNode(node.NewNodeWithDefault(node.LessOrEqualOperator))
	lessOrEqualOperator.SetParent(comparisonOperatorNode)
	lessThanOperator := nfa.getNewState()
	lessThanOperator.SetNode(node.NewNodeWithDefault(node.LessThanOperator))
	lessThanOperator.SetParent(comparisonOperatorNode)
	equalOperator := nfa.getNewState()
	equalOperator.SetNode(node.NewNodeWithDefault(node.EqualOperator))
	equalOperator.SetParent(comparisonOperatorNode)
	notEqual1Operator := nfa.getNewState()
	notEqual1Operator.SetNode(node.NewNodeWithDefault(node.NotEqual1Operator))
	notEqual1Operator.SetParent(comparisonOperatorNode)
	notEqual2Operator := nfa.getNewState()
	notEqual2Operator.SetNode(node.NewNodeWithDefault(node.NotEqual2Operator))
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

func (nfa *NFA) parseWhereOperator(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	whereOperatorNode := node.NewNodeWithDefault(node.WhereOperator)
	start.SetNode(whereOperatorNode)
	start.SetParent(parent)
	andKeyword := nfa.getNewState()
	andKeyword.SetNode(node.NewNodeWithDefault(node.AndKeyword))
	andKeyword.SetParent(whereOperatorNode)
	orKeyword := nfa.getNewState()
	orKeyword.SetNode(node.NewNodeWithDefault(node.OrKeyword))
	orKeyword.SetParent(whereOperatorNode)
	end := nfa.getNewState()

	start.AddNext(token.And, andKeyword)
	andKeyword.AddNext(token.Epsilon, end)
	start.AddNext(token.Or, orKeyword)
	orKeyword.AddNext(token.Epsilon, end)

	return start, end
}

func (nfa *NFA) parseStatementTerminator(parent *node.Node) (*State, *State) {
	start := nfa.getNewState()
	statementTerminatorNode := node.NewNodeWithDefault(node.StatementTerminator)
	start.SetNode(statementTerminatorNode)
	start.SetParent(parent)
	semicolonOperator := nfa.getNewState()
	semicolonOperator.SetNode(node.NewNodeWithDefault(node.SemicolonOperator))
	semicolonOperator.SetParent(statementTerminatorNode)
	end := nfa.getNewState()

	start.AddNext(token.Semicolon, semicolonOperator)
	semicolonOperator.AddNext(token.Epsilon, end)

	return start, end
}
