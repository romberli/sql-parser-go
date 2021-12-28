package lexer

import (
	"fmt"
	"strings"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/sql-parser-go/pkg/token"
)

type State struct {
	Index     int
	Next      map[rune][]*State
	IsFinal   bool
	TokenType token.Type
}

func NewState(i int) *State {
	return &State{
		Index: i,
		Next:  make(map[rune][]*State),
	}
}

func (s *State) AddNext(c rune, ns *State) {
	s.Next[c] = append(s.Next[c], ns)
}

func (s *State) Print() {
	printedList := make(map[int]*State)

	if s.IsFinal {
		fmt.Println(fmt.Sprintf(
			"final state found. index: %d, tokenType: %s",
			s.Index, s.TokenType.String()))
		return
	}

	for c, nsList := range s.Next {
		for _, ns := range nsList {
			printChar := c
			if c == constant.ZeroInt {
				printChar = EpsilonRune
			}
			fmt.Println(fmt.Sprintf("state %d + intput '%c' -> state %d", s.Index, printChar, ns.Index))
			printedList[ns.Index] = ns
		}
	}

	for _, ns := range printedList {
		if s.Index != ns.Index {
			ns.Print()
		}
	}
}

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
	s.States = append(s.States, state)
}

func (s *Set) AddNext(c rune, ns *Set) {
	s.Next[c] = ns
}

func (s *Set) String() string {
	var states string

	for _, state := range s.States {
		states += fmt.Sprintf("%d, ", state.Index)
	}

	return fmt.Sprintf("{index: %d, states: [%s]}", s.Index, strings.Trim(strings.TrimSpace(states), constant.CommaString))
}

func (s *Set) Print() {
	printedList := make(map[int]*Set)

	if s.IsFinal {
		fmt.Println(fmt.Sprintf(
			"final state found. index: %d, tokenType: %s",
			s.Index, s.TokenType.String()))
		return
	}

	for c, ns := range s.Next {
		printChar := c
		if c == constant.ZeroInt {
			printChar = EpsilonRune
		}
		fmt.Println(fmt.Sprintf("state {%d %s + intput '%c' -> state %d", s.Index, s.String(), printChar, ns.Index))
		printedList[ns.Index] = ns
	}

	for _, ns := range printedList {
		if s.Index != ns.Index {
			ns.Print()
		}
	}
}
