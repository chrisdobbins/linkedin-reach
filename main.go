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
	"strings"
	"time"

	gm "github.com/chrisdobbins/linkedin-reach/game"
)

var (
	wordToGuess string
	dictionary  *dict
)

const maxAttempts = 6

var game *gm.Game

func init() {
	rand.Seed(time.Now().UnixNano())
	dictionary = &dict{}
	var err error
	if err = dictionary.populate(); err != nil {
		log.Fatal(err)
	}
	wordToGuess, err = getWord(*dictionary, maxAttempts)
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

func getWord(d dict, maxAttempts int) (string, error) {
	if len(d.words) == 0 {
		return "", errors.New("no guessable words")
	}
	positionMap := map[rune][]int{}
	wordIndex := rand.Intn(len(d.words))
	word := d.words[wordIndex]

	for idx, ch := range word {
		positionMap[ch] = append(positionMap[ch], idx)
	}
	// ensure that word can be guessed within
	// configured # of attempts
	if (len(positionMap) > maxAttempts) || (len(word) == 0) {
		newDict := d
		newDict.words = append(append([]string{}, d.words[:wordIndex]...), d.words[wordIndex+1:]...)
		return getWord(newDict, maxAttempts)
	}
	return word, nil
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
