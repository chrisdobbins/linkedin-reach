package cli

import (
	"bufio"
	"log"
	"os"

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
