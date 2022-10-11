package dependency

import (
	"github.com/romberli/sql-parser-go/pkg/ast"
	"github.com/romberli/sql-parser-go/pkg/token"
)

type Lexer interface {
	// Print prints all the states/sets of the finite automata
	Print()
	// Match matches the given runes and returns proper token
	Match(runes []rune) *token.Token
}

type Parser interface {
	Match() (*ast.Node, error)
}
