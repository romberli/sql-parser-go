package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testLLParser *LLOne
)

func init() {
	initTestLLParser()
}

func initTestLLParser() {
	testLLParser = NewLLOne(initTokenList())
}

func TestLLParser_All(t *testing.T) {
	TestLLParser_Match(t)
}

func TestLLParser_Match(t *testing.T) {
	asst := assert.New(t)

	rootNode, err := testLLParser.Match()
	asst.Nil(err, "test Match() failed")
	if err == nil {
		rootNode.PrintChildren()
	}
}
