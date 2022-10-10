package ast

type Tree struct {
	Root *Node
}

func NewTree(root *Node) *Tree {
	return &Tree{
		Root: root,
	}
}
