package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/chrisdobbins/linkedin-reach/dictionary"
	gm "github.com/chrisdobbins/linkedin-reach/game"
)

const defaultMaxAttempts = 6

const guessesUsage = "Configures the maximum allowed number of guesses."
const help = `
This is a word-guessing game similar to hangman. 
Rules:
You are allowed a certain number of guesses for a word. Each word is guaranteed to be guessable within the allowed number of guesses.
Each guess must be an ASCII letter; all other inputs will be rejected, though they will not count against your remaining guesses. Guesses are case-insensitive.
Good luck and have fun!

Basic options:
-h, --help: Brings up this message
-guesses, --guesses: Configures the maximum allowed number of guesses. Default is %d` + "\n"

var (
	wordToGuess string
	game        *gm.Game
	helpFlag    bool
	maxAttempts int
)

func init() {
	rand.Seed(time.Now().UnixNano())

	flag.Usage = func() { fmt.Fprintf(os.Stderr, fmt.Sprintf(help, defaultMaxAttempts)) }
	flag.IntVar(&maxAttempts, "guesses", defaultMaxAttempts, guessesUsage)

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
	flag.Parse()

	game = gm.Setup(wordToGuess, maxAttempts)
	for !game.IsOver() {
		game.Display()
		reader := bufio.NewReader(os.Stdin)
		guess, _, _ := reader.ReadRune()
		game.Update(guess)
	}
	game.End()
}
