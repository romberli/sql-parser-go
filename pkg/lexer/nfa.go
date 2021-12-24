package lexer

const (
	SelectString = "select"
	FromString   = "from"
	AsString     = "as"
)

var (
	KeywordList = map[TokenType]string{
		Select: SelectString,
		From:   FromString,
		As:     AsString}
)

func getNFA() (int, *State) {
	var (
		i         int
		exists    bool
		lastState *State
	)

	start := NewState(i)

	for tokenType, keyword := range KeywordList {
		for _, c := range keyword {

		}

		keywordStart := NewState(i + 1)
		lastState = start
		start.AddNext(Epsilon, keywordStart)
		exists = true

		for _, c := range keyword {
			if exists {
				for nc, ns := range lastState.Next {
					if nc == Epsilon {

					}
				}
			}

			s := NewState(i + 1)
			s.AppendValue(c)
			lastState.AddNext(c, s)
			lastState = s
			isNew = false
		}
		lastState.TokenType = tokenType

	}

}
