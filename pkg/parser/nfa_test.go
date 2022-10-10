package parser

import (
	"testing"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/sql-parser-go/pkg/token"
	"github.com/stretchr/testify/assert"
)

var (
	testNFA *NFA
)

func init() {
	initTestNFA()
}

func initTestNFA() {
	tokenList := []*token.Token{
		token.NewToken(token.Select, "select"),
		token.NewToken(token.Identifier, "col1"),
		token.NewToken(token.Plus, "+"),
		token.NewToken(token.Identifier, "col2"),
		token.NewToken(token.As, "as"),
		token.NewToken(token.Identifier, "col_alias"),
		token.NewToken(token.From, "from"),
		token.NewToken(token.Identifier, "t01"),
		token.NewToken(token.Where, "where"),
		token.NewToken(token.Identifier, "id"),
		token.NewToken(token.Equal, "="),
		token.NewToken(token.NumberLiteral, "123"),
		token.NewToken(token.Plus, "+"),
		token.NewToken(token.NumberLiteral, "456"),
		token.NewToken(token.And, "and"),
		token.NewToken(token.Identifier, "col1"),
		token.NewToken(token.Equal, "="),
		token.NewToken(token.StringLiteral, "'abc'"),
		token.NewToken(token.End, constant.EmptyString),
	}

	testNFA = NewNFA(tokenList)
}

func TestNFA_All(t *testing.T) {
	TestNFA_Print(t)
	TestNFA_Match(t)
}

func TestNFA_Print(t *testing.T) {
	testNFA.Print()
}

func TestNFA_Match(t *testing.T) {
	asst := assert.New(t)

	rootNode, err := testNFA.Match()
	asst.Nil(err, "test Match() failed")
	rootNode.PrintNext()
}
