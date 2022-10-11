package parser

import (
	"github.com/romberli/sql-parser-go/pkg/ast"
	"github.com/romberli/sql-parser-go/pkg/dependency"
)

type Parser struct {
	fa dependency.Parser
}

func NewParser(fa dependency.Parser) *Parser {
	return &Parser{fa}
}

func (p *Parser) Parse() (*ast.Node, error) {
	return p.fa.Match()
}
