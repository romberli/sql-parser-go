package dependency

import (
	"github.com/romberli/sql-parser-go/pkg/token"
)

type FiniteAutomata interface {
	Print()
	Match(runes []rune) *token.Token
}
