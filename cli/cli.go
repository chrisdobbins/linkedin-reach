package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chrisdobbins/linkedin-reach/dictionary"
	gm "github.com/chrisdobbins/linkedin-reach/game"
)

func PlayGame(gameDictionary *dictionary.Dict, maxAttempts int) {
	var uiDisplay Display
	var err error
	wordCriteria := dictionary.WordCriteria{
		MaxUniqueChars: maxAttempts,
	}
	wordToGuess, err := gameDictionary.GetOne(wordCriteria)
	if err != nil {
		log.Fatal(err)
	}
	game, err := gm.Setup(wordToGuess, maxAttempts)
	if err != nil {
		log.Fatalf("unable to set up game: %s", err.Error())
	}
	for !game.IsOver() {
		uiDisplay = transform(game.Progress())
		uiDisplay.Write()
		reader := bufio.NewReader(os.Stdin)
		guess, _, _ := reader.ReadRune()
		game.Update(guess)
	}

	uiDisplay = transform(game.Result())
	uiDisplay.Write()
}

func transform(state gm.State) (d Display) {
	remainingGuessesTemplate := "%d guess%s left"
	pluralizer := "es"
	prompt := []byte("Guess a letter: ")
	d.Messages = [][]byte{}
	d.Secret = []byte(strings.Join(state.Secret, ""))
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
