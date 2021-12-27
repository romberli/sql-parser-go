package lexer

import (
	"testing"
)

func TestNFA_Print(t *testing.T) {
	nfa := NewNFAWithDefault()
	nfa.Init()
	nfa.Print()
}
