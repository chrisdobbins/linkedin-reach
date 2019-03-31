package game

import (
	"strings"
	"testing"
)

func TestUpdate(t *testing.T) {
	g := Setup("hello", 6)
	// testing field that's displayed to user
	// scenario: valid but incorrect guess
	g.Update('a')
	// test `g.guessedChars` field
	expectedGuessedChars := map[rune]struct{}{'a': struct{}{}}
	if len(expectedGuessedChars) != len(g.guessedChars) {
		t.Logf("expected %d guessedChars in g.guessedChars, actual number was %d", len(expectedGuessedChars), len(g.guessedChars))
		t.Fail()
	}
	missedChars := []rune{}
	for ch, _ := range g.guessedChars {
		if _, ok := expectedGuessedChars[ch]; !ok {
			t.Logf("expected %c not found in g.guessedChars", ch)
			missedChars = append(missedChars, ch)
		}
	}
	if len(missedChars) > 0 {
		t.Fail()
	}
	// test g.disp
	expectedDisplay := []string{"_", "_", "_", "_", "_"}
	if strings.Join(expectedDisplay, "") != strings.Join(g.disp, "") {
		t.Logf("expected `g.disp` to be: %v, got %v", expectedDisplay, g.disp)
		t.Fail()
	}
	// END incorrect, but valid test

	// scenario: correct lowercase guess
	g = Setup("world", 6)
	g.Update('o')
	// test g.guessedChars
	missedChars = []rune{}
	expectedGuessedChars = map[rune]struct{}{'o': struct{}{}}
	for ch, _ := range g.guessedChars {
		if _, ok := expectedGuessedChars[ch]; !ok {
			t.Logf("expected %c not found in g.guessedChars", ch)
			missedChars = append(missedChars, ch)
		}
	}
	if len(missedChars) > 0 {
		t.Fail()
	}
	// test g.disp
	expectedDisplay = []string{"_", "o", "_", "_", "_"}
	if strings.Join(expectedDisplay, "") != strings.Join(g.disp, "") {
		t.Logf("expected `g.disp` to be: %v, got %v", expectedDisplay, g.disp)
		t.Fail()
	}
	// end correct lowercase guess

	// scenario: correct lowercase guess that appears multiple times in word
	g = Setup("arrrgh", 6)
	g.Update('r')
	// test `g.guessedChars` field
	missedChars = []rune{}
	expectedGuessedChars = map[rune]struct{}{'r': struct{}{}}
	for ch, _ := range g.guessedChars {
		if _, ok := expectedGuessedChars[ch]; !ok {
			t.Logf("expected %c not found in g.guessedChars", ch)
			missedChars = append(missedChars, ch)
		}
	}
	if len(missedChars) > 0 {
		t.Fail()
	}
	// test `g.expectedDisplay`
	expectedDisplay = []string{"_", "r", "r", "r", "_", "_"}
	if strings.Join(expectedDisplay, "") != strings.Join(g.disp, "") {
		t.Logf("expected `g.disp` to be: %v, got %v", expectedDisplay, g.disp)
		t.Fail()
	}
	// end correct lowercase guess

	// scenario: upper case but otherwise correct guess
	g = Setup("hello", 6)
	g.Update('L')
	// test `g.guessedChars`
	expectedGuessedChars = map[rune]struct{}{'l': struct{}{}}
	missedChars = []rune{}
	for ch, _ := range g.guessedChars {
		if _, ok := expectedGuessedChars[ch]; !ok {
			t.Logf("expected %c not found in g.guessedChars", ch)
			missedChars = append(missedChars, ch)
		}
	}
	if len(missedChars) > 0 {
		t.Fail()
	}
	// test `g.disp`
	expectedDisplay = []string{"_", "_", "l", "l", "_"}
	if strings.Join(expectedDisplay, "") != strings.Join(g.disp, "") {
		t.Logf("expected `g.disp` to be: %v, got %v", expectedDisplay, g.disp)
		t.Fail()
	}
	// end test uppercase correct guess

	// scenario: invalid character guess
	g = Setup("hello", 7)
	g.Update('#')
	expectedDisplay = []string{"_", "_", "_", "_", "_"}
	if strings.Join(expectedDisplay, "") != strings.Join(g.disp, "") {
		t.Logf("expected `g.disp` to be: %v, got %v", expectedDisplay, g.disp)
		t.Fail()
	}
	// test g.warning
	if g.warning != "Invalid input. Please enter a letter." {
		t.Logf("expected g.warning to be %s, actual value is %v", "Invalid input. Please enter a letter.", g.warning)
		t.Fail()
	}
	// test g.remainingGuesses
	if g.remainingGuesses != 7 {
		t.Logf("expected g.remainingGuesses to be %d, actual value is %v", 7, g.remainingGuesses)
		t.Fail()
	}
	// test g.guessedChars
	expectedGuessedChars = map[rune]struct{}{}
	if len(g.guessedChars) > 0 {
		t.Logf("expected g.guessedChars to be %v, actual is %v", expectedGuessedChars, g.guessedChars)
		t.Fail()
	}

}
