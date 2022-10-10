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

// NewDFA returns a new *DFA
func NewDFA(cs *CharacterSet) *DFA {
	dfa := &DFA{
		CharacterSet: cs,
		Index:        -1,
	}

	dfa.init()

	return dfa
}

// NewDFAWithDefault returns a new *DFA with default
func NewDFAWithDefault() *DFA {
	cs := NewCharacterSetWithDefault()

	return NewDFA(cs)
}

// init initialize the DFA
func (dfa *DFA) init() {
	nfa := NewNFA(dfa.CharacterSet)
	dfa.NFA = nfa

	dfa.InitSet = dfa.getNewSet()

	allSets := []*Set{dfa.InitSet}

	// initialize the first set
	for _, state := range nfa.InitState.EpsilonMove() {
		dfa.InitSet.AddState(state)
	}

	// initialize a channel to save the sets that need be processed
	setChan := make(chan *Set, maxSetChanLength)
	// put the init set to the channel
	setChan <- dfa.InitSet

	var (
		setExists bool
	)
	for {
		if len(setChan) == constant.ZeroInt {
			// all sets are processed
			close(setChan)
			break
		}
		// get a set from the channel
		currentSet := <-setChan
		nextRunes := map[rune]bool{}
		// get all the next runes except the MayBeEpsilon of the states in the set
		for _, state := range currentSet.States {
			for c := range state.Next {
				if c != token.EpsilonRune {
					nextRunes[c] = true
				}
			}
		}

		for c := range nextRunes {
			// create a new set
			nextSet := dfa.getNewSet()

			for _, state := range currentSet.States {
				for _, ns := range state.Next[c] {
					// get the next state of c and also all the epsilon move states
					for _, epsilonState := range ns.EpsilonMove() {
						nextSet.AddState(epsilonState)
					}
				}
			}

			for _, set := range allSets {
				if set.Equal(nextSet) {
					// this set already exists, use the old one as the next set of the current set
					currentSet.Next[c] = set
					setExists = true
					dfa.Index--
					break
				}
			}
			if !setExists {
				// this is a brand-new set, add this to the all set
				allSets = append(allSets, nextSet)
				// use the new one as the next set of the current set
				currentSet.Next[c] = nextSet
				// send the new set to the channel and wait to be processed
				setChan <- nextSet
			}

			setExists = false
		}
	}
}

// Print prints all the sets including the state in it
func (dfa *DFA) Print() {
	dfa.InitSet.Print()
}

// Match matches the given runes and returns proper token
func (dfa *DFA) Match(runes []rune) *token.Token {
	tempSet := dfa.InitSet

	for _, c := range runes {
		// transit to the next set
		tempSet = tempSet.Next[c]
		if tempSet == nil {
			return token.NewToken(token.Error, string(runes))
		}
	}

	if tempSet.IsFinal {
		// final set found
		return token.NewToken(tempSet.TokenType, string(runes))
	}

	return token.NewToken(token.Error, string(runes))
}

// getNewSet gets a new set
func (dfa *DFA) getNewSet() *Set {
	dfa.Index++
	return NewSet(dfa.Index)
}
