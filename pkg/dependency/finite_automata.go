package dependency

import (
	"github.com/romberli/sql-parser-go/pkg/token"
)

type FiniteAutomata interface {
	Init()
	Print()
	Match(runes []rune) *token.Token
}
