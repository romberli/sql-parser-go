package lexer

import (
	"github.com/romberli/go-util/constant"
	"github.com/romberli/sql-parser-go/pkg/token"
)

const maxSetChanLength = 10000

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

	setChan := make(chan *Set, maxSetChanLength)
	setChan <- dfa.InitSet

	var (
		setExists bool
	)
	for {
		if len(setChan) == constant.ZeroInt {
			close(setChan)
			break
		}
		currentSet := <-setChan
		nextRunes := map[rune]bool{}

		for _, state := range currentSet.States {
			for c := range state.Next {
				if c != token.Epsilon {
					nextRunes[c] = true
				}
			}
		}
		for c := range nextRunes {
			nextSet := dfa.getNewSet()

			for _, state := range currentSet.States {
				for _, ns := range state.Next[c] {
					for _, epsilonState := range ns.EpsilonMove() {
						nextSet.AddState(epsilonState)
					}
				}
			}

			for _, set := range allSets {
				if set.Equal(nextSet) {
					currentSet.Next[c] = set
					setExists = true
					dfa.Index--
					break
				}
			}
			if !setExists {
				allSets = append(allSets, nextSet)
				currentSet.Next[c] = nextSet
				setChan <- nextSet
			}

			setExists = false
		}
	}
}

func (dfa *DFA) Print() {
	dfa.InitSet.Print()
}

func (dfa *DFA) Match(runes []rune) *token.Token {
	tempSet := dfa.InitSet
	for _, c := range runes {
		tempSet = tempSet.Next[c]
		if tempSet == nil {
			return token.NewToken(token.Error, string(runes))
		}
	}

	if tempSet.IsFinal {
		return token.NewToken(tempSet.TokenType, string(runes))
	}

	return token.NewToken(token.Error, string(runes))
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
