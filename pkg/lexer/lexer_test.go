package lexer

import (
	"fmt"
	"testing"
)

var (
	testNFALexer *Lexer
	testDFALexer *Lexer
)

func init() {
	initTestNFA()
	testNFALexer = NewLexer(testNFA)
	initTestDFA()
	testDFALexer = NewLexer(testNFA)
}

func TestLexer_All(t *testing.T) {
	TestLexer_Lex(t)
}

func TestLexer_Lex(t *testing.T) {
	sql := `     select 123*(456+789),
            col1, col2, 'abc123_' from
            t01 where id <= 123 and col1='abc'    ;   `
	// sql := `col1, col2 from t01 where id = 123 and col1='abc';`

	tokens := testNFALexer.Lex(sql)
	fmt.Println("==========NFA==========")
	for _, token := range tokens {
		fmt.Println(token.String())
	}
	fmt.Println("==========NFA==========")

	tokens = testDFALexer.Lex(sql)
	fmt.Println("==========DFA==========")
	for _, token := range tokens {
		fmt.Println(token.String())
	}
	fmt.Println("==========DFA==========")

}
