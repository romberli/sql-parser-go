package parser

import (
	"fmt"

	"github.com/romberli/sql-parser-go/pkg/node"
	"github.com/romberli/sql-parser-go/pkg/token"
)

type State struct {
	Index  int
	Node   *node.Node
	Parent *node.Node
	Next   map[token.Type][]*State
}

func NewState(i int) *State {
	return &State{
		Index: i,
		Next:  make(map[token.Type][]*State),
	}
}

func (s *State) SetNode(node *node.Node) {
	s.Node = node
}

func (s *State) SetParent(parent *node.Node) {
	s.Parent = parent
}

func (s *State) AddNext(t token.Type, ns *State) {
	s.Next[t] = append(s.Next[t], ns)
}

func (s *State) EpsilonMove() []*State {
	states := []*State{s}

	s.epsilonMove(&states)

	return states
}

func (s *State) epsilonMove(states *[]*State) {
	for _, state := range s.Next[token.Epsilon] {
		*states = append(*states, state)
		state.epsilonMove(states)
	}
}

func (s *State) Print() {
	printedList := make(map[int]*State)

	s.print(printedList)
}

func (s *State) print(printedList map[int]*State) {
	if s.Node != nil {
		if s.Node.Type == node.End {
			fmt.Println(fmt.Sprintf(
				"final state found. index: %d, nodeType: %s",
				s.Index, s.Node.Type.String()))
			return
		}

		fmt.Println(fmt.Sprintf(
			"Node is not null. index: %d, node type: %s",
			s.Index, s.Node.Type.String()))
	}

	for t, nsList := range s.Next {
		for _, ns := range nsList {
			fmt.Println(fmt.Sprintf(
				"state found. index: %d, next state: %d, token type: %s",
				s.Index, ns.Index, t.String()))
			_, ok := printedList[ns.Index]
			if !ok {
				printedList[ns.Index] = ns
				ns.print(printedList)
			}
		}
	}
}
