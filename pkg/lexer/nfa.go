package lexer

const (
	SelectString = "select"
	FromString   = "from"
	AsString     = "as"
)

func getNFA() {

}

func getSelectNFA(i int) (int, *State) {
	start := NewState(i)

	for _, c := range SelectString {
		s := NewState(i + 1)
	}
}
