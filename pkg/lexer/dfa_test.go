package lexer

import (
	"fmt"
	"testing"
)

var (
	testDFA *DFA
)

func init() {
	initTestDFA()
}

func initTestDFA() {
	testDFA = NewDFAWithDefault()
	testDFA.Init()
}

func TestDFA_All(t *testing.T) {
	TestDFA_Print(t)
	TestDFA_Match(t)
}

func TestDFA_Print(t *testing.T) {
	testDFA.Print()
}

func TestDFA_Match(t *testing.T) {
	strList := []string{"select", "and", "as", "selectt", "'string'", "123", "123abc", ">=", "123."}

	for _, str := range strList {
		token := testDFA.Match([]rune(str))
		fmt.Println(token.String())
	}
}
