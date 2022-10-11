package parser

import (
	"github.com/pingcap/errors"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/sql-parser-go/pkg/ast"
	"github.com/romberli/sql-parser-go/pkg/token"
)

type LLOne struct {
	Tokens []*token.Token
	Index  int
}

func NewLLOne(tokens []*token.Token) *LLOne {
	return &LLOne{
		Tokens: append(tokens, token.NewToken(token.End, constant.EmptyString)),
		Index:  -1,
	}
}

func (llo *LLOne) Match() (*ast.Node, error) {
	rootNode := llo.getNewNode(ast.Root)
	err := llo.match(rootNode)
	if err != nil {
		return nil, err
	}

	if llo.Tokens[llo.Index+1].Type != token.End {
		return nil, errors.Errorf("matching token failed: matched tokens: %v, next token: %s", llo.Tokens[:llo.Index+1], llo.Tokens[llo.Index+1])
	}

	return rootNode, nil
}

func (llo *LLOne) match(n *ast.Node) error {
	if llo.Index == len(llo.Tokens)-1 {
		if llo.Tokens[llo.Index].Type == token.End {
			return nil
		}

		return errors.Errorf("matching token failed: matched tokens: %v, next token: %s", llo.Tokens[:llo.Index], llo.Tokens[llo.Index])
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
		if token.TypeExists(firstSet, llo.lookAhead().Type) {
			// found correct path
			children = childrenList[i]
			for _, child := range children {
				childFirst := child.GetFirstSet()
			Loop:
				if token.TypeExists(childFirst, llo.lookAhead().Type) {
					// as there is no conflict, always can add child
					n.AddChildren(child)
					if child.IsTerminal() {
						t := llo.readNext()
						if child.Type.GetTokenType() != t.Type {
							return errors.Errorf("match terminal failed. child type: %s, token type: %s", child.Type.String(), t.String())
						}

						child.SetToken(t)
						continue
					}

					err := llo.match(child)
					if err != nil {
						return err
					}

					if child.Max == -1 {
						goto Loop
					}

					continue
				}

				if child.MayEpsilon() {
					continue
				}

				return errors.Errorf("matching token failed: node type: %s, matched tokens: %v, next token: %s", n.Type.String(), llo.Tokens[:llo.Index], llo.Tokens[llo.Index])
			}

			return nil
		}

	}

	return errors.Errorf("matching token failed: node type: %s, matched tokens: %v, next token: %s", n.Type.String(), llo.Tokens[:llo.Index], llo.Tokens[llo.Index])
}

func (llo *LLOne) lookAhead() *token.Token {
	// if nfa.Index+1 >= len(nfa.Tokens) {
	//     return nil
	// }

	return llo.Tokens[llo.Index+1]
}

func (llo *LLOne) readNext() *token.Token {
	// if nfa.Index+1 >= len(nfa.Tokens) {
	//     return nil
	// }
	llo.Index++

	return llo.Tokens[llo.Index]
}

// getNewState gets a new state
func (llo *LLOne) getNewState() *State {
	llo.Index++
	return NewState(llo.Index)
}

func (llo *LLOne) getNewNode(t ast.Type) *ast.Node {
	return ast.NewNodeWithDefault(t)
}
