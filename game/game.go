package game

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	tm "github.com/buger/goterm"
)

type Player int

const (
	User Player = iota + 1
	Computer
)

type Game struct {
	secret           string
	disp             display
	warning          string
	remainingGuesses int
	guessedChars     map[rune]struct{}
	charPositions    map[rune][]int
	winner           *Player
}

func Setup(secretWord string, guesses int) *Game {
	d := &display{}
	d.init(len(secretWord))
	g := &Game{
		secret:           secretWord,
		disp:             *d,
		remainingGuesses: guesses,
		guessedChars:     map[rune]struct{}{},
		charPositions:    map[rune][]int{},
		winner:           nil,
	}
	for idx, ch := range g.secret {
		g.charPositions[ch] = append(g.charPositions[ch], idx)
	}
	return g
}

type display []string

func (d *display) init(length int) {
	disp := []string{}
	for i := 0; i < length; i++ {
		disp = append(disp, "_")
	}
	*d = display(disp)
}

func (d *display) update(ltr rune, positions ...int) {
	dCopy := []string{}
	for i, ch := range *d {
		dCopy = append(dCopy, string(ch))
		for _, pos := range positions {
			if i == pos {
				dCopy[i] = string(ltr)
				continue
			}
		}
	}
	*d = display(dCopy)
}

func (g *Game) Update(guess rune) {
	warning := validate(guess, g.guessedChars)
	if len(warning) > 0 {
		g.warning = warning
		return
	}
	g.warning = ""
	guess = lowerCase(guess)

	g.guessedChars[guess] = struct{}{}
	for _, pos := range g.charPositions[guess] {
		g.disp.update(guess, pos)
	}
	g.remainingGuesses--
	if g.remainingGuesses == 0 {
		g.winner = new(Player)
		if strings.Join(g.disp, "") == g.secret {
			*g.winner = User
			return
		}
		*g.winner = Computer
	}
}

func lowerCase(input rune) rune {
	if input < 'a' {
		return input + 32
	}
	return input
}

func validate(guess rune, guessedLetters map[rune]struct{}) (warning string) {
	nonLtrFilter := regexp.MustCompile(`[[:^alpha:]]`)
	if guess == utf8.RuneError || len(nonLtrFilter.FindString(string(guess))) > 0 {
		warning = "Invalid input. Please enter a letter."
	}
	if _, ok := guessedLetters[guess]; ok {
		warning = "Letter already guessed. Please try again."
	}
	return
}

func (g *Game) Display() {
	tm.Clear()
	tm.MoveCursor(1, 1)
	for _, ch := range g.disp {
		tm.Print(ch)
	}
	tm.Println()
	tm.MoveCursor(1, 3)
	tm.Println(g.warning)
	tm.Println("Guessed letters: ")
	i := 0
	for ch, _ := range g.guessedChars {
		i++
		tm.Print(string(ch))
		if i < len(g.guessedChars) {
			tm.Print(",")
		}
	}
	tm.Println()
	tm.Println(fmt.Sprintf("Attempts left: %d", g.remainingGuesses))
	tm.Println("Guess a letter: ")
	tm.ResetLine("")
	tm.Flush()
}

func (g *Game) IsOver() bool {
	return g.winner != nil
}

func (g *Game) End() {
	tm.Clear()
	tm.MoveCursor(1, 1)
	if *g.winner == User {
		tm.Println(g.secret)
		tm.Println("You win!")
		tm.Flush()
		return
	}
	tm.Println("You lose :(")
	tm.Println("The word was", g.secret)
	tm.Flush()
}
