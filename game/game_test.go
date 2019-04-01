package game

import (
	"strings"
	"testing"
)

var nilPlayer *Player

func TestUpdateValidChar(t *testing.T) {
	t.Run("Incorrect guess", func(t *testing.T) {
		maxGuesses := 6
		g := Setup("hello", maxGuesses)
		g.Update('a')
		expectedGuessedChars := map[rune]struct{}{'a': struct{}{}}
		if len(expectedGuessedChars) != len(g.guessedChars) {
			t.Errorf("number of g.guessedChars is %d", len(g.guessedChars))
		}
		for ch, _ := range g.guessedChars {
			if _, ok := expectedGuessedChars[ch]; !ok {
				t.Errorf("%c not found in g.guessedChars", ch)
			}
		}
		expectedDisplay := []string{"_", "_", "_", "_", "_"}
		if strings.Join(expectedDisplay, "") != strings.Join(g.disp, "") {
			t.Errorf("g.disp is %v", g.disp)
		}
		expectedRemainingGuesses := maxGuesses - 1
		if g.remainingGuesses != expectedRemainingGuesses {
			t.Errorf("g.remainingGuesses is %v", g.remainingGuesses)
		}
		if g.winner != nilPlayer {
			t.Errorf("winner is %v", g.winner)
		}
	})

	t.Run("Correct lowercase guess", func(t *testing.T) {
		maxGuesses := 6
		g := Setup("world", maxGuesses)
		g.Update('o')
		expectedGuessedChars := map[rune]struct{}{'o': struct{}{}}
		for ch, _ := range g.guessedChars {
			if _, ok := expectedGuessedChars[ch]; !ok {
				t.Errorf("%c not found in g.guessedChars", ch)
			}
		}
		expectedDisplay := []string{"_", "o", "_", "_", "_"}
		if strings.Join(expectedDisplay, "") != strings.Join(g.disp, "") {
			t.Errorf("`g.disp` is %v", g.disp)
		}
		expectedRemainingGuesses := maxGuesses - 1
		if g.remainingGuesses != expectedRemainingGuesses {
			t.Errorf("g.remainingGuesses is %d", g.remainingGuesses)
		}
		if g.winner != nilPlayer {
			t.Errorf("g.winner = %v", *g.winner)
		}
	})

	t.Run("Correct lowercase guess that appears multiple times", func(t *testing.T) {
		maxGuesses := 6
		g := Setup("arrrgh", maxGuesses)
		g.Update('r')
		expectedGuessedChars := map[rune]struct{}{'r': struct{}{}}
		for ch, _ := range g.guessedChars {
			if _, ok := expectedGuessedChars[ch]; !ok {
				t.Errorf("%c not found in g.guessedChars", ch)
			}
		}
		expectedDisplay := []string{"_", "r", "r", "r", "_", "_"}
		if strings.Join(expectedDisplay, "") != strings.Join(g.disp, "") {
			t.Errorf("g.disp is %v", g.disp)
		}
		expectedRemainingGuesses := maxGuesses - 1
		if g.remainingGuesses != expectedRemainingGuesses {
			t.Errorf("g.remainingGuesses is %v", g.remainingGuesses)
		}
		if g.winner != nilPlayer {
			t.Errorf("expected nil for winner, got %v", *g.winner)
		}
	})

	t.Run("Correct upper case guess", func(t *testing.T) {
		maxGuesses := 6
		g := Setup("hello", maxGuesses)
		g.Update('L')
		expectedGuessedChars := map[rune]struct{}{'l': struct{}{}}
		for ch, _ := range g.guessedChars {
			if _, ok := expectedGuessedChars[ch]; !ok {
				t.Errorf("%c not found in g.guessedChars", ch)
			}
		}
		expectedDisplay := []string{"_", "_", "l", "l", "_"}
		if strings.Join(expectedDisplay, "") != strings.Join(g.disp, "") {
			t.Errorf("g.disp is %v", g.disp)
		}
		expectedRemainingGuesses := maxGuesses - 1
		if g.remainingGuesses != expectedRemainingGuesses {
			t.Errorf("g.remainingGuesses is %v", g.remainingGuesses)
		}
		if g.winner != nilPlayer {
			t.Errorf("winner is %v", *g.winner)
		}
	})
}

func TestUpdateInvalidChar(t *testing.T) {
	maxGuesses := 6
	g := Setup("hello", maxGuesses)
	g.Update('#')
	expectedDisplay := []string{"_", "_", "_", "_", "_"}
	if strings.Join(expectedDisplay, "") != strings.Join(g.disp, "") {
		t.Errorf("g.disp is %v", g.disp)
	}
	if g.warning != "Invalid input. Please enter a letter." {
		t.Errorf("g.warning is %s", g.warning)
	}
	if len(g.guessedChars) > 0 {
		t.Errorf("g.guessedChars is %v", g.guessedChars)
	}
	expectedRemainingGuesses := maxGuesses
	if g.remainingGuesses != expectedRemainingGuesses {
		t.Errorf("g.remainingGuesses is %v", g.remainingGuesses)
	}
	if g.winner != nilPlayer {
		t.Errorf("winner is %v", *g.winner)
	}
}

func TestUpdateGameOver(t *testing.T) {
	t.Run("User wins", func(t *testing.T) {
		guesses := map[rune]struct{}{'s': struct{}{}, 'e': struct{}{}, 't': struct{}{}}
		g := Setup("test", len(guesses))
		for guess, _ := range guesses {
			g.Update(guess)
		}
		if *g.winner != User {
			t.Errorf("winner is %v", *g.winner)
		}
	})

	t.Run("Computer wins", func(t *testing.T) {
		guesses := map[rune]struct{}{'x': struct{}{}, 'e': struct{}{}, 't': struct{}{}}
		g := Setup("test", len(guesses))
		for guess, _ := range guesses {
			g.Update(guess)
		}
		if *g.winner != Computer {
			t.Errorf("winner is %v", *g.winner)
		}
	})
}
