package server

import (
	"sort"
	gm "github.com/chrisdobbins/linkedin-reach/game"
	"strings"
)

type Display struct {
	Secret           string `json:"secret"`
	GuessedChars     string `json:"guessedChars"`
	Message          string `json:"message"`
	RemainingGuesses int    `json:"remainingGuesses"`
	GameOver         bool   `json:"gameOver"`
}

func (d *Display) format() {
	if len(d.GuessedChars) == 0 {
		return
	}
	gc := strings.Split(d.GuessedChars, "")
	var sortChars = func(i, j int) bool { return gc[i] < gc[j] }
	sort.Slice(gc, sortChars)
	d.GuessedChars = strings.Join(gc, ",")
}

func transform(state gm.State) (*Display) {
	d := &Display{}
	d.Secret = strings.Join(state.Secret, "")
	gc := []string{}
	for _, ch := range state.GuessedChars {
		gc = append(gc, string(ch))
	}
	d.GuessedChars = strings.Join(gc, "")
	d.RemainingGuesses = state.RemainingGuesses
	d.Message = state.Message
	if state.EndResult != nil {
		d.GameOver = true
	}
	return d
}
