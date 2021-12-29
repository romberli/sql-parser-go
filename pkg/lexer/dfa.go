package lexer

import (
	"github.com/romberli/sql-parser-go/pkg/token"
)

type DFA struct {
	CharacterSet *CharacterSet
	NFA          *NFA
	Index        int
	InitSet      *Set
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
	nfa := NewNFA(dfa.CharacterSet)
	nfa.Init()
	dfa.NFA = nfa

	dfa.InitSet = dfa.getNewSet()

	allSets := []*Set{dfa.InitSet}

	for _, state := range nfa.InitState.EpsilonMove() {
		dfa.InitSet.AddState(state)
	}

	setChan := make(chan *Set)
	setChan <- dfa.InitSet

	var (
		ok        bool
		setExists bool
		nextSet   *Set
	)
	for {
		nextSet, ok = <-setChan
		if !ok {
			break
		}

		nextRunes := map[rune]bool{}

		for _, state := range nextSet.States {
			for c := range state.Next {
				if c != token.Epsilon {
					nextRunes[c] = true
				}
			}
		}
		for c := range nextRunes {
			newSet := dfa.getNewSet()

			for _, state := range dfa.InitSet.States {
				for _, st := range state.Next[c] {
					for _, epsilonState := range st.EpsilonMove() {
						if !newSet.Contains(epsilonState) {
							newSet.AddState(epsilonState)
						}
					}
				}
			}

			for _, as := range allSets {
				if as.Equal(newSet) {
					setExists = true
					dfa.Index--
					break
				}
			}
			if !setExists {
				allSets = append(allSets, newSet)
				setChan <- newSet
				setExists = false
			}

			setExists = false
		}
	}
}

func (dfa *DFA) Print() {
	dfa.InitSet.Print()
}

func (dfa *DFA) Match(runes []rune) *token.Token {
	return nil
}

func (dfa *DFA) getNewSet() *Set {
	dfa.Index++
	return NewSet(dfa.Index)
}

func (dfa *DFA) getFinalState(tokenType token.Type) *Set {
	final := dfa.getNewSet()
	final.IsFinal = true
	final.TokenType = tokenType

	return final
}
