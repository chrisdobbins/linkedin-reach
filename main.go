package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/chrisdobbins/linkedin-reach/cli"
	"github.com/chrisdobbins/linkedin-reach/dictionary"
	"github.com/chrisdobbins/linkedin-reach/server"
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
-guesses, --guesses: Configures the maximum allowed number of guesses. Default is %d
-serve, --serve: Starts web version of this game` + "\n"

var (
	wordToGuess    string
	helpFlag       bool
	maxAttempts    int
	port           string
	shouldServe    bool
	gameDictionary *dictionary.Dict
)

func init() {
	rand.Seed(time.Now().UnixNano())

	flag.Usage = func() { fmt.Fprintf(os.Stderr, fmt.Sprintf(help, defaultMaxAttempts)) }
	flag.IntVar(&maxAttempts, "guesses", defaultMaxAttempts, guessesUsage)
	flag.BoolVar(&shouldServe, "serve", false, "whether to start web version of game")
	port = "8080"
	if os.Getenv("PORT") != "" {
		shouldServe = true
		port = os.Getenv("PORT")
	}

	newDict, err := dictionary.New()
	if err != nil {
		log.Fatal(err)
	}
	gameDictionary = &newDict
}

func main() {
	flag.Parse()
	if &maxAttempts == nil {
		maxAttempts = defaultMaxAttempts
	}
	if shouldServe {
		server.Serve(port, gameDictionary, maxAttempts)
		return
	}
	cli.PlayGame(gameDictionary, maxAttempts)
}
