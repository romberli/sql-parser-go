package lexer

import (
	"fmt"
	"strings"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/sql-parser-go/pkg/token"
)

type Set struct {
	States    []*State
	Index     int
	Next      map[rune]*Set
	IsFinal   bool
	TokenType token.Type
}

// NewSet returns a new *Set
func NewSet(i int) *Set {
	return &Set{
		Index: i,
		Next:  make(map[rune]*Set),
	}
}

// AddState add the given state into the set,
// - if a same state already exists in the set, it will be ignored
// - if the new state is a final state and there is already a final state exists in the set,
// 	 only one state of which token type is keyword stays in the set
func (s *Set) AddState(state *State) {
	if !s.Contains(state) {
		if state.IsFinal {
			s.IsFinal = true

			i, final := s.GetFinalState()
			if final == nil {
				// there is no final state in the set
				s.TokenType = state.TokenType
				s.States = append(s.States, state)
				return
			}

			if state.TokenType.IsKeyword() {
				// keyword token type has the top priority, it will replace the old final state
				s.TokenType = state.TokenType
				s.States[i] = state
				return
			}
		}

		s.States = append(s.States, state)
	}
}

// AddNext adds the next set of given rune
func (s *Set) AddNext(c rune, ns *Set) {
	s.Next[c] = ns
}

// Contains returns if the given state is in the set
func (s *Set) Contains(state *State) bool {
	for _, st := range s.States {
		if st.Index == state.Index {
			return true
		}
	}

	return false
}

// GetFinalState returns the final state in the set
func (s *Set) GetFinalState() (int, *State) {
	for i, state := range s.States {
		if state.IsFinal {
			return i, state
		}
	}

	return constant.ZeroInt, nil
}

// Equal returns if the twe sets contains the same states, the order of the states does not matter
func (s *Set) Equal(other *Set) bool {
	if len(s.States) != len(other.States) {
		return false
	}

	for _, state := range s.States {
		for !other.Contains(state) {
			return false
		}
	}

	return true
}

// String returns the string representation of the set
func (s *Set) String() string {
	var states string

	for _, state := range s.States {
		states += fmt.Sprintf("%d, ", state.Index)
	}

	return fmt.Sprintf(
		"{index: %d, states: [%s]}", s.Index, strings.Trim(strings.TrimSpace(states), constant.CommaString))
}

// Print prints the set and all the next sets recursively
func (s *Set) Print() {
	printedList := make(map[int]*Set)

	if s.IsFinal {
		fmt.Println(fmt.Sprintf(
			"final set found. index: %d, tokenType: %s",
			s.Index, s.TokenType.String()))
		// return
	}

	for c, ns := range s.Next {
		printChar := c
		if c == token.Epsilon {
			printChar = EpsilonRune
		}
		fmt.Println(fmt.Sprintf(
			"set %s + intput '%c' -> set %s", s.String(), printChar, ns.String()))
		printedList[ns.Index] = ns
	}

	for _, ns := range printedList {
		if s.Index != ns.Index {
			// print recursively
			ns.Print()
		}
	}
}
