package parser

import (
	"github.com/pingcap/errors"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/sql-parser-go/pkg/ast"
	"github.com/romberli/sql-parser-go/pkg/token"
)

type LLParser struct {
	Tokens    []*token.Token
	Index     int
	InitNode  *ast.Node
	InitState *State
}

func NewLLParser(tokens []*token.Token) *LLParser {
	llp := &LLParser{
		Tokens: tokens,
		Index:  -1,
	}

	// llp.init()

	return llp
}

func (llp *LLParser) init() {
	llp.InitState = llp.getNewState()
}

func (llp *LLParser) Match() (*ast.Node, error) {
	rootNode := llp.getNewNode(ast.Root)
	err := llp.match(rootNode)
	if err != nil {
		return nil, err
	}

	return rootNode, nil
}

func (llp *LLParser) match(n *ast.Node) error {
	if llp.Index == len(llp.Tokens) {
		//  all input tokens are matched correctly
		return nil
	}

	childrenList := n.GetChildren()

	var firstSetList [][]token.Type
	for _, children := range childrenList {
		firstChild := children[constant.ZeroInt]
		firstSet := firstChild.GetFirstSet()
		firstSetList = append(firstSetList, firstSet)
	}

	// check if there is any conflict
	if token.HasIntersect(firstSetList) {
		return errors.Errorf("token conflict detected, %v", firstSetList)
	}

	var children []*ast.Node
	for i, firstSet := range firstSetList {
		if token.TypeExists(firstSet, llp.lookAhead().Type) {
			// found correct path
			children = childrenList[i]
			for _, child := range children {
				childFirst := child.GetFirstSet()
			Loop:
				if token.TypeExists(childFirst, llp.lookAhead().Type) {
					// as there is no conflict, always can add child
					n.AddChildren(child)
					if child.IsTerminal() {
						t := llp.readNext()
						if child.Type.GetTokenType() != t.Type {
							return errors.Errorf("match terminal failed. child type: %s, token type: %s", child.Type.String(), t.String())
						}

						child.SetToken(t)
						continue
					}

					err := llp.match(child)
					if err != nil {
						return err
					}

					if child.Max == -1 {
						goto Loop
					}
				}

				if child.MayEpsilon() {
					continue
				}
			}

			return nil
		}

	}

	return errors.Errorf("can't find available child. ast type: %s", n.Type.String())

}

func (llp *LLParser) lookAhead() *token.Token {
	// if nfa.Index+1 >= len(nfa.Tokens) {
	//     return nil
	// }

	return llp.Tokens[llp.Index+1]
}

func (llp *LLParser) readNext() *token.Token {
	// if nfa.Index+1 >= len(nfa.Tokens) {
	//     return nil
	// }
	llp.Index++

	return llp.Tokens[llp.Index]
}

// getNewState gets a new state
func (llp *LLParser) getNewState() *State {
	llp.Index++
	return NewState(llp.Index)
}

func (llp *LLParser) getNewNode(t ast.Type) *ast.Node {
	return ast.NewNodeWithDefault(t)
}
