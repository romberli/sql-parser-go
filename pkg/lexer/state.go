package lexer

import (
	"fmt"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/sql-parser-go/pkg/token"
)

type State struct {
	Index     int
	Value     []rune
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

func (s *State) AppendValue(c rune) {
	s.Value = append(s.Value, c)
}

func (s *State) AddNext(c rune, ns *State) {
	s.Next[c] = append(s.Next[c], ns)
}

func (s *State) Transit(c rune) []*State {
	return s.Next[c]
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
