package cli

import (
	"sort"
	"strings"

	tm "github.com/buger/goterm"
)

type Display struct {
	Secret       []byte
	GuessedChars []byte
	Messages     [][]byte
}

func (d Display) Write() {
	tm.Clear()
	tm.MoveCursor(1, 1)
	box := tm.NewBox(60|tm.PCT, 10, 0)
	box.Write(d.format())
	tm.Print(box.String())
	tm.MoveCursor(3, 9)
	tm.Flush()
}

func (d Display) format() []byte {
	rawFormattedDisplay := []byte{}
	guessedChars := d.GuessedChars
	var sortChars = func(i, j int) bool { return guessedChars[i] < guessedChars[j] }
	newLine := []byte("\n")
	rawFormattedDisplay = append(rawFormattedDisplay, append(d.Secret, newLine...)...)
	sort.Slice(guessedChars, sortChars)
	formattedGuessedChars := strings.Join(strings.Split(string(guessedChars), ""), ",")
	rawFormattedDisplay = append(rawFormattedDisplay, append([]byte(formattedGuessedChars), newLine...)...)
	for _, msgs := range d.Messages {
		rawFormattedDisplay = append(rawFormattedDisplay, append(msgs, newLine...)...)
	}
	return rawFormattedDisplay
}
