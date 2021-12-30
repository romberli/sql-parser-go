package lexer

import (
	"fmt"
	"testing"
)

var (
	testNFA *NFA
)

func init() {
	initTestNFA()
}

func initTestNFA() {
	testNFA = NewNFAWithDefault()
	testNFA.Init()
}

func TestNFA_All(t *testing.T) {
	TestNFA_Print(t)
	TestNFA_Match(t)
}

func TestNFA_Print(t *testing.T) {
	testNFA.Print()
}

func TestNFA_Match(t *testing.T) {
	strList := []string{"select", "and", "as", "selectt", "'string'", "123", "123abc", ">=", "123."}

	for _, str := range strList {
		token := testNFA.Match([]rune(str))
		fmt.Println(token.String())
	}
}
