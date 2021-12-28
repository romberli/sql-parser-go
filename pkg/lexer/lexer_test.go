package lexer

import (
	"fmt"
	"testing"
)

var (
	testLexer *Lexer
)

func init() {
	initTestNFA()
	testLexer = NewLexer(testNFA)
}

func TestLexer_All(t *testing.T) {
	TestLexer_Lex(t)
}

func TestLexer_Lex(t *testing.T) {
	sql := `     select 123*(456+789),
            col1, col2, 'abc123_' from
            t01 where id <= 123 and col1='abc'    ;   `
	// sql := `col1, col2 from t01 where id = 123 and col1='abc';`
	tokens := testLexer.Lex(sql)

	for _, token := range tokens {
		fmt.Println(token.String())
	}
}
