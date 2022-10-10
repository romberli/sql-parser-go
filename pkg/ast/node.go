package ast

import (
	"fmt"
	"strings"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/sql-parser-go/pkg/node"
	"github.com/romberli/sql-parser-go/pkg/token"
)

type Node struct {
	Type     node.Type
	Token    *token.Token
	Children []*Node
}

func NewNode(t node.Type, token *token.Token) *Node {
	return newNode(t, token)
}

func NewNodeWithDefault(t node.Type) *Node {
	return newNode(t, nil)
}

func newNode(t node.Type, token *token.Token) *Node {
	return &Node{
		Type:  t,
		Token: token,
	}
}

func (n *Node) IsTerminal() bool {
	return n.Type.IsTerminal()
}

func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

func (n *Node) RemoveChild(index int) {
	if len(n.Children) > index {
		n.Children = append(n.Children[:index], n.Children[index+1:]...)
	}
}

func (n *Node) RemoveLastChild() {
	if len(n.Children) > constant.ZeroInt {
		n.Children = n.Children[:len(n.Children)-1]
	}
}

func (n *Node) Print() {
	n.print(constant.ZeroInt, constant.ZeroInt)
}

func (n *Node) print(i, j int) {
	prefix := fmt.Sprintf("%s|%s", strings.Repeat(constant.SpaceString, i), strings.Repeat(constant.UnderBarString, j))
	fmt.Printf("%s%s\n", prefix, n.Type.String())
	for _, child := range n.Children {
		i += 2
		j += 2
		child.print(i, j)
	}
}
