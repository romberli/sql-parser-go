package lexer

import (
	"github.com/romberli/sql-parser-go/pkg/dependency"
	"github.com/romberli/sql-parser-go/pkg/token"
)

type Lexer struct {
	fa dependency.Lexer
}

// NewLexer returns a new *Lexer
func NewLexer(fa dependency.Lexer) *Lexer {
	return &Lexer{
		fa: fa,
	}
}

// GetFiniteAutomata returns the finite automata of the lexer
func (l *Lexer) GetFiniteAutomata() dependency.Lexer {
	return l.fa
}

// Lex scans the input string and returns a token list
func (l *Lexer) Lex(sql string) []*token.Token {
	var (
		isStringLiteral bool
		runes           []rune
		tokens          []*token.Token
	)

	sqlRunes := []rune(sql)
	length := len(sqlRunes)

	for i, c := range sqlRunes {
		if c == singleQuote && !isStringLiteral {
			// string literal starts
			isStringLiteral = true
			runes = append(runes, c)
			continue
		}

		if c == singleQuote && isStringLiteral {
			// string literal ends
			isStringLiteral = false
			runes = append(runes, c)
			// match string literal token
			tokens = append(tokens, l.GetFiniteAutomata().Match(runes))
			runes = nil
			continue
		}

		if isStringLiteral {
			runes = append(runes, c)
			if i == length-1 {
				// string literal does not end with single quote
				tokens = append(tokens, token.NewToken(token.Error, string(runes)))
			}
			continue
		}

		if IsWhiteSpace(c) {
			continue
		}

		switch c {
		case GTRune:
			runes = append(runes, c)
			if i >= length-1 || sqlRunes[i+1] != EqualRune {
				tokens = append(tokens, l.GetFiniteAutomata().Match(runes))
				runes = nil
			}
		case LTRune:
			runes = append(runes, c)
			if i >= length-1 || (sqlRunes[i+1] != EqualRune && sqlRunes[i+1] != GTRune) {
				tokens = append(tokens, l.GetFiniteAutomata().Match(runes))
				runes = nil
			}
		case ExclamationRune:
			runes = append(runes, c)
			if i >= length-1 || sqlRunes[i+1] != EqualRune {
				tokens = append(tokens, l.GetFiniteAutomata().Match(runes))
				runes = nil
			}
		case EqualRune, PlusRune, MinusRune, MultiplyRune, DivideRune, LeftParenthesisRune, RightParenthesisRune,
			SemicolonRune, CommaRune:
			runes = append(runes, c)
			tokens = append(tokens, l.GetFiniteAutomata().Match(runes))
			runes = nil
		default:
			// match tokens that contains multi runes
			runes = append(runes, c)
			if i >= length-1 || !IsAlphabetOrDigit(sqlRunes[i+1]) {
				tokens = append(tokens, l.GetFiniteAutomata().Match(runes))
				runes = nil
			}
		}
	}

	return tokens
}
