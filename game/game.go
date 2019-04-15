package game

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Player int

const (
	User Player = iota + 1
	Computer
	errInvalidGuess          = "Invalid input.\nPlease enter a letter."
	errInvalidMaxGuessesTmpl = "unable to set up game with %d guesses"
	errLtrAlreadyGuessed     = "Letter already guessed.\nPlease try again."
)

type Game struct {
	secret           string
	disp             display
	message          string
	remainingGuesses int
	guessedChars     map[rune]struct{}
	charPositions    map[rune][]int
	winner           *Player
}

type State struct {
	Secret           []string
	RemainingGuesses int
	Message          string
	GuessedChars     []rune
	EndResult        *Player
}

func (g *Game) Progress() State {
	guessedChars := []rune{}
	for ch, _ := range g.guessedChars {
		guessedChars = append(guessedChars, ch)
	}
	return State{
		Secret:           []string(g.disp),
		RemainingGuesses: g.remainingGuesses,
		Message:          g.message,
		GuessedChars:     guessedChars,
	}
}

func Setup(secretWord string, guesses int) (*Game, error) {
	if guesses <= 0 {
		return &Game{}, fmt.Errorf(errInvalidMaxGuessesTmpl, guesses)
	}
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
	return g, nil
}

func (g *Game) Update(guess rune) {
	warning := validate(guess, g.guessedChars)
	if len(warning) > 0 {
		g.message = warning
		return
	}
	g.message = ""
	guess = unicode.ToLower(guess)

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

func (g *Game) IsOver() bool {
	return g.winner != nil
}

func (g *Game) Result() (s State) {
	s.EndResult = g.winner
	if *g.winner == User {
		s.Message = "You win!"
		return
	}
	s.Message = fmt.Sprintf("You lose :(\n The word was %s", g.secret)
	return
}

func validate(guess rune, guessedLetters map[rune]struct{}) (warning string) {
	nonLtrFilter := regexp.MustCompile(`[[:^alpha:]]`)
	if guess == utf8.RuneError || len(nonLtrFilter.FindString(string(guess))) > 0 {
		warning = errInvalidGuess
	}
	if _, ok := guessedLetters[guess]; ok {
		warning = errLtrAlreadyGuessed
	}
	return
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


