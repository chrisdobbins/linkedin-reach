package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	tm "github.com/buger/goterm"
)

var (
	dictionary *dict
)

const maxAttempts = 6

func init() {
	rand.Seed(time.Now().UnixNano())
	dictionary = &dict{}
	if err := dictionary.populate(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	wordToGuess, positionMap, err := getWord(*dictionary, maxAttempts)
	if err != nil {
		log.Fatal(err)
	}
	display := createInitialDisplay(wordToGuess)
	tm.Println(wordToGuess)
	var warnings []string
	guessedChars := []rune{}

	for attempts := 0; attempts < maxAttempts && strings.Join(display, "") != wordToGuess; {
		attemptsLeft := maxAttempts - attempts
		updateTerminal(wordToGuess, attemptsLeft, guessedChars, display, warnings...)

		reader := bufio.NewReader(os.Stdin)
		guess, _, _ := reader.ReadRune()
		if isInvalid(guess) {
			warnings = []string{"Invalid input. Please enter a letter."}
			continue
		}
		warnings = []string{}
		guessedChars = append(guessedChars, guess)
		attempts++
		for _, g := range positionMap[guess] {
			display[g] = string(wordToGuess[g])
		}
	}
	endGame(wordToGuess, display)
}

func endGame(word string, display []string) {
	tm.Clear()
	tm.MoveCursor(1, 1)
	if strings.Join(display, "") == word {
		tm.Println(word)
		tm.Println("You win!")
		tm.Flush()
		return
	}

	tm.Println("You lose :(")
	tm.Println("The word was", word)
	tm.Flush()
}

func updateTerminal(displayWord string, attemptsLeft int, guessedChars []rune, displayMap []string, warnings ...string) {
	tm.Clear()
	tm.MoveCursor(1, 1)
	for idx, _ := range displayWord {
		tm.Print(displayMap[idx])
	}
	for i, warn := range warnings {
		tm.MoveCursor(1, i+2)
		tm.Println(warn)
	}
	tm.Println()
	tm.Println("Guessed letters: ")
	for i, ch := range guessedChars {
		tm.Print(string(ch))
		if i < len(guessedChars)-1 {
			tm.Print(",")
		}
	}
	tm.Println()
	tm.Println(fmt.Sprintf("%d attempts left", attemptsLeft))
	tm.Println("Guess a letter: ")
	tm.ResetLine("")
	tm.Flush()
}

func createInitialDisplay(word string) (display []string) {
	for _ = range word {
		display = append(display, "_")
	}
	return display
}

type dict struct {
	words []string
}

func (d *dict) populate() error {
	words, err := getDictionary()
	if err != nil {
		return err
	}
	d.words = words
	return nil
}

// testing this first
func getWord(d dict, maxAttempts int) (string, map[rune][]int, error) {
	if len(d.words) == 0 {
		return "", nil, errors.New("no guessable words")
	}
	positionMap := map[rune][]int{}
	wordIndex := rand.Intn(len(d.words))
	word := d.words[wordIndex]

	for idx, ch := range word {
		positionMap[ch] = append(positionMap[ch], idx)
	}
	// ensure that word can be guessed within
	// configured # of attempts
	if (len(positionMap) > maxAttempts && len(word) > maxAttempts) || (len(word) == 0) {
		newDict := d
		newDict.words = append(append([]string{}, d.words[:wordIndex]...), d.words[wordIndex+1:]...)
		return getWord(newDict, maxAttempts)
	}
	return word, positionMap, nil
}

func isInvalid(guess rune) bool {
	nonLtrFilter := regexp.MustCompile(`[[:^alpha:]]`)
	return (guess == utf8.RuneError || len(nonLtrFilter.FindString(string(guess))) > 0)
}

func getDictionary() ([]string, error) {
	errorTemplate := "unable to retrieve word list: %s"
	url := "http://app.linkedin-reach.io/words"
	resp, err := http.Get(url)
	if err != nil {
		return []string{}, fmt.Errorf(errorTemplate, err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []string{}, fmt.Errorf(errorTemplate, fmt.Sprintf("error getting words from %s", url))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, fmt.Errorf(errorTemplate, err.Error())
	}

	return strings.Split(string(body), "\n"), nil
}
