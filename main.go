package main

import (
	"bufio"
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
	dict        []string
	wordToGuess string
)

func updateDisplay(displayWord string, displayMap []string, warnings ...string) {
	tm.Clear()
	tm.MoveCursor(1, 1)
	for idx, _ := range displayWord {
		tm.Print(displayMap[idx])
	}
	tm.MoveCursor(1, 2)
	for _, warn := range warnings {
		tm.Println(warn)
	}
	tm.Println("Guess a letter: ")
	tm.Flush()
}

func createInitialDisplay(word string) (display []string) {
	for _ = range word {
		display = append(display, "_")
	}
	return display
}

func findChar(ch rune, positionMap map[rune][]int) (bool, []int) {
	if len(positionMap[ch]) == 0 {
		return false, []int{}
	}
	return true, positionMap[ch]
}

func getWord(maxAttempts int) (string, map[rune][]int) {
	positionMap := map[rune][]int{}
	dict, err := getDictionary()
	if err != nil {
		log.Fatal(err)
	}
	word := dict[rand.Intn(len(dict))]
	for idx, ch := range word {
		positionMap[ch] = append(positionMap[ch], idx)
	}
	// ensure that word can be guessed within
	// configured # of attempts
	if len(positionMap) > maxAttempts && len(word) > maxAttempts {
		return getWord(maxAttempts)
	}
	return word, positionMap
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	maxAttempts := 6
	wordToGuess, positionMap := getWord(maxAttempts)
	display := createInitialDisplay(wordToGuess)
	var warnings []string
	for attempts := 0; attempts < maxAttempts; {
		tm.Println(attempts)
		tm.Println(wordToGuess)
		updateDisplay(wordToGuess, display, warnings...)
		reader := bufio.NewReader(os.Stdin)
		re := regexp.MustCompile(`[[:^alpha:]]`)
		guess, _, _ := reader.ReadRune()

		if guess == utf8.RuneError || len(re.Find([]byte(string(guess)))) > 0 {
			warnings = []string{"Invalid input. Please enter a letter."}
			continue
		} else {
			warnings = []string{}
		}
		for _, g := range positionMap[guess] {
			display[g] = string(wordToGuess[g])
		}
		attempts++
	}
	if strings.Join(display, "") != wordToGuess {
		tm.Clear()
		tm.Println("Game over.")
		tm.Flush()
	} else {
		tm.Clear()
		tm.Println("you win")
		tm.Flush()
	}
}

func getDictionary() ([]string, error) {
	errorTemplate := "unable to retrieve word list: %s"
	resp, err := http.Get("http://app.linkedin-reach.io/words")
	if err != nil {
		return []string{}, fmt.Errorf(errorTemplate, err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, fmt.Errorf(errorTemplate, err.Error())
	}

	return strings.Split(string(body), "\n"), nil
}
