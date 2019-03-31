package game

import "testing"

func TestUpdate(t *testing.T) {
	g := Setup("hello", 6)
	// testing field that's displayed to user
	// scenario: incorrect guess
	g.Update('a')
	expected := []string{"_", "_", "_", "_", "_", "_"}
	failedIndices := []int{}
	for idx, ch := range g.disp {
		if ch != expected[idx] {
			failedIndices = append(failedIndices, idx)
			t.Logf("expected %s at index %d, got %v", expected[idx], idx, ch)
		}
	}
	if len(failedIndices) > 0 {
		t.Fail()
	}

	g = Setup("world", 6)
	// correct guess
	g.Update('o')
	expected = []string{"_", "o", "_", "_", "_"}
	for idx, ch := range g.disp {
		if ch != expected[idx] {
			failedIndices = append(failedIndices, idx)
			t.Logf("expected %s at index %d, got %v", expected[idx], idx, ch)
		}
	}
	if len(failedIndices) > 0 {
		t.Fail()
	}

	// correct guess that appears multiple times in word
	g = Setup("arrrgh", 6)
	g.Update('r')
	expected = []string{"_", "r", "r", "r", "_", "_"}
	for idx, ch := range g.disp {
		if ch != expected[idx] {
			failedIndices = append(failedIndices, idx)
			t.Logf("expected %s at index %d, got %v", expected[idx], idx, ch)
		}
	}
	if len(failedIndices) > 0 {
		t.Fail()
	}

	g = Setup("hello", 6)
	// upper case but otherwise correct guess
	g.Update('L')
	expected = []string{"_", "_", "l", "l", "_"}
	for idx, ch := range g.disp {
		if ch != expected[idx] {
			failedIndices = append(failedIndices, idx)
			t.Logf("expected %s at index %d, got %v", expected[idx], idx, ch)
		}
	}
	if len(failedIndices) > 0 {
		t.Fail()
	}

	g = Setup("hello", 7)
	// invalid character
	g.Update('#')
	expected = []string{"_", "_", "_", "_", "_"}

	for idx, ch := range g.disp {
		if ch != expected[idx] {
			failedIndices = append(failedIndices, idx)
			t.Logf("expected %s at index %d, got %v", expected[idx], idx, ch)
		}
	}
	if len(failedIndices) > 0 {
		t.Fail()
	}
	if g.warning != "Invalid input. Please enter a letter." {
		t.Logf("expected g.warning to be %s, actual value is %v", "Invalid input. Please enter a letter.", g.warning)
                t.Fail()
	}
	if g.remainingGuesses != 7 {
		t.Logf("expected g.remainingGuesses to be %d, actual value is %v", 7, g.remainingGuesses)
		t.Fail()
	}
	expectedGuessedChars := map[rune]struct{}{}
	if len(g.guessedChars) > 0 {
		t.Logf("expected g.guessedChars to be %v, actual is %v", expectedGuessedChars, g.guessedChars)
		t.Fail()
	}

}
