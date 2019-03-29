package main

import (
	"math/rand"
	"testing"
	"time"
)

func TestGetWord(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	d := &dict{}
        maxAttempts := 6

	// "guessable" words are those that have <= 6 unique letters
	guessable := []string{"apple", "foil", "foie", "zygote", "delivered"}
	guessableIdx := rand.Intn(len(guessable))
	unguessable := []string{"regulator", "deliverance", "fragility"}

	// scenario: guessable word exists in `d.words`
	expected := guessable[guessableIdx]
	d.words = append([]string{guessable[guessableIdx]}, unguessable...)
	word, positionMap, err := getWord(*d, maxAttempts)
	if err != nil {
		t.Fail()
	}
	if positionMap == nil {
		t.Fail()
	}
	for _, ch := range expected {
		if _, ok := positionMap[ch]; !ok {
			t.Logf("%s not found in positionMap", string(ch))
			t.Fail()
		}
	}
	if word != expected {
		t.Logf("expected word: %s, actual word: %s", expected, word)
		t.Fail()
	}

	// scenario: all words are unguessable within 6 attempts
	d.words = unguessable
	word, positionMap, err = getWord(*d, maxAttempts)
	expected = ""
	if err == nil {
		t.Logf("expected error to be %v, actual error was: %s", nil, err.Error())
		t.Fail()
	}
	if word != expected {
		t.Logf("expected word: %s, actual word: %s", expected, word)
		t.Fail()
	}
	if positionMap != nil {
		t.Logf("expected: %v, received: %v", nil, positionMap)
		t.Fail()
	}

	// scenario: empty slice
	d.words = []string{}
	_, _, err = getWord(*d,maxAttempts) 
	if err == nil {
		t.Logf("expected error, got nil")
		t.Fail()
	}
}
