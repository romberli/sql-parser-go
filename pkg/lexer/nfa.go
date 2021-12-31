package lexer

import (
	"github.com/romberli/go-util/constant"
	"github.com/romberli/sql-parser-go/pkg/token"
)

const (
	// keyword
	SelectString = "select"
	FromString   = "from"
	AsString     = "as"
	WhereString  = "where"
	AndString    = "and"
	// comparison operator
	GEString        = ">="
	LEString        = "<="
	NotEqual1String = "!="
	NotEqual2String = "<>"
)

var (
	MultiRuneList = map[token.Type]string{
		// keyword
		token.Select: SelectString,
		token.From:   FromString,
		token.As:     AsString,
		token.Where:  WhereString,
		token.And:    AndString,
		// comparison operator
		token.GE:        GEString,
		token.LE:        LEString,
		token.NotEqual1: NotEqual1String,
		token.NotEqual2: NotEqual2String,
	}
	SingleRuneList = map[token.Type]rune{
		// comparison operator
		token.GT:    GTRune,
		token.LT:    LTRune,
		token.Equal: EqualRune,
		// arithmetic operator
		token.Plus:     PlusRune,
		token.Minus:    MinusRune,
		token.Multiply: MultiplyRune,
		token.Divide:   DivideRune,
		token.Mod:      ModRune,
		// parenthesis
		token.LeftParenthesis:  LeftParenthesisRune,
		token.RightParenthesis: RightParenthesisRune,
		// symbol
		token.Comma:     CommaRune,
		token.Semicolon: SemicolonRune,
	}
)

type NFA struct {
	CharacterSet *CharacterSet
	Index        int
	InitState    *State
}

// NewNFA returns a new *NFA
func NewNFA(cs *CharacterSet) *NFA {
	nfa := &NFA{
		CharacterSet: cs,
		Index:        -1,
	}

	nfa.init()

	return nfa
}

// NewNFAWithDefault returns a new *NFA with default
func NewNFAWithDefault() *NFA {
	cs := NewCharacterSetWithDefault()

	return NewNFA(cs)
}

// init initialize the NFA
func (nfa *NFA) init() {
	nfa.InitState = nfa.getNewState()

	nfa.initMultiRune()
	nfa.initSingleRune()
	nfa.initIdentifier()
	nfa.initStringLiteral()
	nfa.initNumberLiteral()
}

// initMultiRune initialize the states that can recognize tokens which have multi runes
func (nfa *NFA) initMultiRune() {
	for tokenType, tokenString := range MultiRuneList {
		start := nfa.getNewState()
		// temporary state
		tempState := start
		nfa.InitState.AddNext(token.Epsilon, start)

		for _, c := range tokenString {
			s := nfa.getNewState()
			tempState.AddNext(c, s)
			tempState = s
		}

		final := nfa.getNewFinalState(tokenType)
		tempState.AddNext(token.Epsilon, final)
	}
}

// initIdentifier initialize the states that can recognize the identifier
func (nfa *NFA) initIdentifier() {
	start := nfa.getNewState()
	nfa.InitState.AddNext(token.Epsilon, start)

	for _, c := range nfa.CharacterSet.GetDigits() {
		start.AddNext(c, start)
	}

	s := nfa.getNewState()
	for _, c := range nfa.CharacterSet.GetAlphabets() {
		start.AddNext(c, s)
	}

	for _, c := range nfa.CharacterSet.GetAlphabets() {
		s.AddNext(c, s)
	}
	for _, c := range nfa.CharacterSet.GetDigits() {
		s.AddNext(c, s)
	}

	final := nfa.getNewFinalState(token.Identifier)
	s.AddNext(token.Epsilon, final)
}

// initSingleRune initialize the states that can recognize tokens which have single rune
func (nfa *NFA) initSingleRune() {
	for tokenType, c := range SingleRuneList {
		start := nfa.getNewState()
		nfa.InitState.AddNext(token.Epsilon, start)

		s := nfa.getNewState()
		start.AddNext(c, s)

		final := nfa.getNewFinalState(tokenType)
		s.AddNext(token.Epsilon, final)
	}
}

// initStringLiteral initialize the states that can recognize string literal token
func (nfa *NFA) initStringLiteral() {
	start := nfa.getNewState()
	nfa.InitState.AddNext(token.Epsilon, start)

	openQuote := nfa.getNewState()
	start.AddNext(singleQuote, openQuote)

	for _, c := range nfa.CharacterSet.GetAlphabets() {
		openQuote.AddNext(c, openQuote)
	}
	for _, c := range nfa.CharacterSet.GetDigits() {
		openQuote.AddNext(c, openQuote)
	}

	closeQuote := nfa.getNewState()
	openQuote.AddNext(singleQuote, closeQuote)

	final := nfa.getNewFinalState(token.StringLiteral)
	closeQuote.AddNext(token.Epsilon, final)
}

// initNumberLiteral initialize the states that can recognize number literal token
func (nfa *NFA) initNumberLiteral() {
	start := nfa.getNewState()
	nfa.InitState.AddNext(token.Epsilon, start)
	s := nfa.getNewState()

	for _, c := range nfa.CharacterSet.GetDigits() {
		start.AddNext(c, s)
		s.AddNext(c, s)
	}

	final := nfa.getNewFinalState(token.NumberLiteral)
	s.AddNext(token.Epsilon, final)
}

// Print prints all the states
func (nfa *NFA) Print() {
	nfa.InitState.Print()
}

// Match matches the given runes and returns proper token
func (nfa *NFA) Match(runes []rune) *token.Token {
	return nfa.match(nfa.InitState, constant.ZeroInt, runes)
}

func (nfa *NFA) match(s *State, i int, runes []rune) *token.Token {
	if i == len(runes) {
		// all input runes are matched, check the result
		if s.IsFinal {
			// final state found
			return token.NewToken(s.TokenType, string(runes))
		}
		// this state is not a final state, need to check if there is any ε-move that can transit to the final state
		nsList := s.Next[token.Epsilon]
		for _, ns := range nsList {
			if ns.IsFinal {
				// final state found
				return token.NewToken(ns.TokenType, string(runes))
			}
		}
		// all input runes are matched, and there is no ε-move that can transit to the final state, return error token
		return token.NewToken(token.Error, string(runes))
	}

	nsList := s.Next[runes[i]]
	if nsList == nil {
		nsList = s.Next[token.Epsilon]
		if nsList == nil {
			//  can't transit to any other state, return error token
			return token.NewToken(token.Error, string(runes))
		}
	} else {
		// matched an input rune, increase the matching index
		i++
	}

	for _, ns := range nsList {
		// match next rune recursively
		t := nfa.match(ns, i, runes)
		// if returning token is not an error token, it means matched some token,
		// otherwise, means this is not a good path, need to try another one
		if t.Type != token.Error {
			return t
		}
	}

	return token.NewToken(token.Error, string(runes))
}

// getNewState gets a new state
func (nfa *NFA) getNewState() *State {
	nfa.Index++
	return NewState(nfa.Index)
}

// getNewFinalState gets a new final state
func (nfa *NFA) getNewFinalState(tokenType token.Type) *State {
	final := nfa.getNewState()
	final.IsFinal = true
	final.TokenType = tokenType

	return final
}
