package cli

import (
	"fmt"
	"sort"
	"strings"

	tm "github.com/buger/goterm"
	gm "github.com/chrisdobbins/linkedin-reach/game"
)

type Display struct {
	Secret       []byte
	GuessedChars []byte
	Messages     [][]byte
}

func (d Display) Write() {
	tm.Clear()
	tm.MoveCursor(1, 1)
	box := tm.NewBox(tm.Width(), 10, 0)
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

func transform(state gm.State) (d Display) {
	remainingGuessesTemplate := "%d guess%s left"
	pluralizer := "es"
	prompt := []byte("Guess a letter: ")
	d.Messages = [][]byte{}
	d.Secret = []byte(strings.Join(state.Secret, "  "))
	d.GuessedChars = []byte{}
	for _, ch := range state.GuessedChars {
		d.GuessedChars = append(d.GuessedChars, byte(ch))
	}
	if state.RemainingGuesses == 1 {
		pluralizer = ""
	}
	d.Messages = append(d.Messages, []byte(fmt.Sprintf(remainingGuessesTemplate, state.RemainingGuesses, pluralizer)))
	d.Messages = append(d.Messages, []byte(state.Message))
	if state.EndResult == nil {
		d.Messages = append(d.Messages, append(prompt))
	}
	return
}
