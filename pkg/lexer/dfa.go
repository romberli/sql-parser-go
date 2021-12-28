package lexer

import (
	"github.com/romberli/sql-parser-go/pkg/token"
)

type DFA struct {
	CharacterSet *CharacterSet
	Index        int
	InitState    *State
}

func NewDFA(cs *CharacterSet) *DFA {
	return &DFA{
		CharacterSet: cs,
		Index:        -1,
	}
}

func NewDFAWithDefault() *DFA {
	cs := NewCharacterSetWithDefault()

	return NewDFA(cs)
}

func (dfa *DFA) Init() {

}

func (dfa *DFA) Print() {
	dfa.InitState.Print()
}

func (dfa *DFA) Match(runes []rune) *token.Token {
	return nil
}
