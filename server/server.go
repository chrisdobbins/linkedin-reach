package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/chrisdobbins/linkedin-reach/dictionary"
	gm "github.com/chrisdobbins/linkedin-reach/game"
	"github.com/gorilla/websocket"
)

var (
	upgrader       = websocket.Upgrader{}
	tmpl           = template.Must(template.New("").Parse(frontend))
	wordToGuess    string
	maxAttempts    int
	gameDictionary *dictionary.Dict
)

func Serve(port string, d *dictionary.Dict, maxGuesses int) {
	maxAttempts = maxGuesses
	if d == nil {
		log.Fatal("Serve: dictionary is empty")
	}
	gameDictionary = d
	http.HandleFunc("/play", playGame)
	http.HandleFunc("/", page)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func playGame(w http.ResponseWriter, r *http.Request) {
	var uiDisplay *Display
	const logErrTemplate = "playGame: %s\n"
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf(logErrTemplate, err.Error())
		http.Error(w, "Oops! Please try again.", http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	wordToGuess, err := getWord()
	if err != nil {
		log.Printf(logErrTemplate, err.Error())
		http.Error(w, "Oops! Please try again.", http.StatusInternalServerError)
		return
	}
	game, err := gm.Setup(wordToGuess, maxAttempts)
	if err != nil {
		log.Printf(logErrTemplate, err.Error())
		http.Error(w, "Oops! Please try again.", http.StatusInternalServerError)
		return
	}
	for !game.IsOver() {
		uiDisplay = transform(game.Progress())
		uiDisplay.format()
		out, err := json.Marshal(*uiDisplay)
		if err != nil {
			log.Printf(logErrTemplate, err.Error())
			http.Error(w, "Oops! Please try again.", http.StatusInternalServerError)
			return
		}
		conn.WriteMessage(websocket.TextMessage, out)
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf(logErrTemplate, err.Error())
			http.Error(w, "Oops! Please try again.", http.StatusInternalServerError)
			return
		}
		if len(msg) > 0 {
			game.Update(rune(msg[0]))
		}
	}
	uiDisplay = transform(game.Result())
	out, err := json.Marshal(*uiDisplay)
	if err != nil {
		log.Printf(logErrTemplate, err.Error())
		http.Error(w, "Oops! Please try again.", http.StatusInternalServerError)
		return
	}
	conn.WriteMessage(websocket.TextMessage, out)
}

func page(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, "wss://"+r.Host+"/play")
}

func getWord() (string, error) {
	wordToGuess, err := gameDictionary.GetOne(dictionary.WordCriteria{MaxUniqueChars: maxAttempts})
	if err != nil {
		return "", fmt.Errorf("getWord: %s", err.Error())
	}
	return wordToGuess, nil
}
