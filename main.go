package main

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/chrisdobbins/linkedin-reach/dictionary"
	gm "github.com/chrisdobbins/linkedin-reach/game"
)

var (
	wordToGuess string
)

const maxAttempts = 6

var game *gm.Game

func init() {
	rand.Seed(time.Now().UnixNano())
	wordCriteria := dictionary.WordCriteria{
		MaxUniqueChars: maxAttempts,
	}
	gameDictionary, err := dictionary.New()
	if err != nil {
		log.Fatal(err)
	}
	wordToGuess, err = gameDictionary.GetOne(wordCriteria)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	game = gm.Setup(wordToGuess, maxAttempts)
	for !game.IsOver() {
		game.Display()
		reader := bufio.NewReader(os.Stdin)
		guess, _, _ := reader.ReadRune()
		game.Update(guess)
	}
	game.End()
}
