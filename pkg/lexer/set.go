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

func NewSet(i int) *Set {
	return &Set{
		Index: i,
		Next:  make(map[rune]*Set),
	}
}

func (s *Set) AddState(state *State) {
	if !s.Contains(state) {
		if state.IsFinal {
			s.IsFinal = true

			i, final := s.GetFinalState()
			if final == nil {
				s.TokenType = state.TokenType
				s.States = append(s.States, state)
				return
			}

			if state.TokenType.IsKeyword() {
				s.TokenType = state.TokenType
				s.States[i] = state
				return
			}
		}

		s.States = append(s.States, state)
	}
}

func (s *Set) AddNext(c rune, ns *Set) {
	s.Next[c] = ns
}

func (s *Set) Contains(state *State) bool {
	for _, st := range s.States {
		if st.Index == state.Index {
			return true
		}
	}

	return false
}

func (s *Set) GetFinalState() (int, *State) {
	for i, state := range s.States {
		if state.IsFinal {
			return i, state
		}
	}

	return constant.ZeroInt, nil
}

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

func (s *Set) String() string {
	var states string

	for _, state := range s.States {
		states += fmt.Sprintf("%d, ", state.Index)
	}

	return fmt.Sprintf(
		"{index: %d, states: [%s]}", s.Index, strings.Trim(strings.TrimSpace(states), constant.CommaString))
}

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
			ns.Print()
		}
	}
}
