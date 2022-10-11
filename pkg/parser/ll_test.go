package parser

import (
	"testing"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/sql-parser-go/pkg/token"
	"github.com/stretchr/testify/assert"
)

var (
	testLLParser *LLParser
)

func init() {
	initTestNFA()
	initTestLLParser()
}

func initTokenList() []*token.Token {
	return []*token.Token{
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
}

func initTestLLParser() {
	testLLParser = NewLLParser(initTokenList())
}

func TestLLParser_All(t *testing.T) {
	TestLLParser_Print(t)
	TestLLParser_Match(t)
}

func TestLLParser_Print(t *testing.T) {
	// testLLParser.Print()
}

func TestLLParser_Match(t *testing.T) {
	asst := assert.New(t)

	rootNode, err := testLLParser.Match()
	asst.Nil(err, "test Match() failed")
	rootNode.PrintChildren()
}
