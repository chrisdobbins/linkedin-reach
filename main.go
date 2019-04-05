package main

import (
	"bufio"
	"flag"
	"fmt"
	_ "go/build"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/chrisdobbins/linkedin-reach/dictionary"
	gm "github.com/chrisdobbins/linkedin-reach/game"
	"github.com/chrisdobbins/linkedin-reach/server"
	"github.com/chrisdobbins/linkedin-reach/ui"
)

const defaultMaxAttempts = 6

const guessesUsage = "Configures th maximum allowed number of guesses."
const help = `
This is a word-guessing game similar to hangman. 
Rules:
You are allowed a certain number of guesses for a word. Each word is guaranteed to be guessable within the allowed number of guesses.
Each guess must be an ASCII letter; all other inputs will be rejected, though they will not count against your remaining guesses. Guesses are case-insensitive.
Good luck and have fun!

Basic options:
-h, --help: Brings up this message
-guesses, --guesses: Configures the maximum allowed number of guesses. Default is %d
-serve, --serve: Starts web version of this game` + "\n"

var (
	wordToGuess  string
	helpFlag     bool
	maxAttempts  int
	serveAddress string
	shouldServe  bool
)

func init() {
	rand.Seed(time.Now().UnixNano())

	flag.Usage = func() { fmt.Fprintf(os.Stderr, fmt.Sprintf(help, defaultMaxAttempts)) }
	flag.IntVar(&maxAttempts, "guesses", defaultMaxAttempts, guessesUsage)
	flag.BoolVar(&shouldServe, "serve", false, "whether to start web version of game")
	serveAddress = "localhost:8080"

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
	var uiDisplay ui.Display
	if shouldServe {
		server.Serve(serveAddress, wordToGuess, defaultMaxAttempts)
	} else { // debug else statement
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
}

func transform(state gm.State) (d ui.Display) {
	remainingGuessesTemplate := "%d guesses left"
	prompt := []byte("Guess a letter: ")
	d.Messages = [][]byte{}
	d.Secret = []byte(strings.Join(state.Secret, ""))
	d.GuessedChars = []byte{}
	for _, ch := range state.GuessedChars {
		d.GuessedChars = append(d.GuessedChars, byte(ch))
	}
	d.Messages = append(d.Messages, []byte(fmt.Sprintf(remainingGuessesTemplate, state.RemainingGuesses)))
	d.Messages = append(d.Messages, []byte(state.Message))
	if state.EndResult == nil {
		d.Messages = append(d.Messages, append(prompt))
	}
	return
}
