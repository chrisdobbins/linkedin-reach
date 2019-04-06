package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/chrisdobbins/linkedin-reach/dictionary"
	gm "github.com/chrisdobbins/linkedin-reach/game"
	"github.com/gorilla/websocket"
)

var (
	wordToGuess    string
	maxAttempts    int
	uiDisplay      Display
	gameDictionary *dictionary.Dict
)

func getWord() (string, error) {
	wordToGuess, err := gameDictionary.GetOne(dictionary.WordCriteria{MaxUniqueChars: maxAttempts})
	if err != nil {
		return "", fmt.Errorf("getWord: %s", err.Error())
	}
	return wordToGuess, nil
}

func Serve(port string, d *dictionary.Dict, maxGuesses int) {
	maxAttempts = maxGuesses
	if d == (&dictionary.Dict{}) {
		log.Fatal("Serve: dictionary is nil")
	}
	gameDictionary = d
	http.HandleFunc("/play", playGame)
	http.HandleFunc("/", page)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s",port), nil))
}

func playGame(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("unable to upgrade connection: %s", err.Error())
		http.Error(w, "Oops! Please try again.", http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	wordToGuess, err := getWord()
	if err != nil {
		log.Printf("playGame: %s\n", err.Error())
		http.Error(w, "Oops! Please try again.", http.StatusInternalServerError)
		return
	}
	game, err := gm.Setup(wordToGuess, maxAttempts)
	if err != nil {
		log.Println(fmt.Sprintf("gm.Setup: %s", err.Error()))
		http.Error(w, "Oops! Please try again.", http.StatusInternalServerError)
		return
	}
	for !game.IsOver() {
		uiDisplay = transform(game.Progress())
		out, err := json.Marshal(uiDisplay)
		if err != nil {
			log.Printf("unable to marshal JSON: %s", err.Error())
			http.Error(w, "Oops! Please try again.", http.StatusInternalServerError)
			return
		}
		conn.WriteMessage(websocket.TextMessage, out)
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("err conn.ReadMessage: %s\n", err.Error())
			http.Error(w, "Oops! Please try again.", http.StatusInternalServerError)
			return
		}
		if len(msg) > 0 {
			game.Update(rune(msg[0]))
		}
	}
	uiDisplay = transform(game.Result())
	out, err := json.Marshal(uiDisplay)
	if err != nil {
		log.Printf("unable to unmarshal %v: %s", out, err.Error())
		http.Error(w, "Oops! Please try again.", http.StatusInternalServerError)
		return
	}
	conn.WriteMessage(websocket.TextMessage, out)
}

type Display struct {
	Secret           string `json:"secret"`
	GuessedChars     string `json:"guessedChars"`
	Message          string `json:"message"`
	RemainingGuesses int    `json:"remainingGuesses"`
        GameOver bool `json:"gameOver"`
}

func transform(state gm.State) (d Display) {
	d.Secret = strings.Join(state.Secret, "")
	gc := []string{}
	for _, ch := range state.GuessedChars {
		gc = append(gc, string(ch))
	}
	d.GuessedChars = strings.Join(gc, ",")
	d.RemainingGuesses = state.RemainingGuesses
	d.Message = state.Message
        if state.EndResult != nil {
             d.GameOver = true
        }
	return
}

func page(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, "ws://"+r.Host+"/play")
}

var (
	upgrader = websocket.Upgrader{}
	tmpl     = template.Must(template.New("").Parse(frontend))
)
