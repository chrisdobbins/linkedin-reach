package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	gm "github.com/chrisdobbins/linkedin-reach/game"
	"github.com/gorilla/websocket"
)

var (
	wordToGuess string
	maxAttempts int
	uiDisplay   Display
)

func Serve(addr, word string, maxGuesses int) {
	maxAttempts = maxGuesses
	fmt.Println("serve max attemtps: ", maxAttempts)
	wordToGuess = word
	http.HandleFunc("/play", playGame)
	http.HandleFunc("/", page)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func playGame(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer conn.Close()

	game, err := gm.Setup(wordToGuess, maxAttempts)
	if err != nil {
		log.Println(fmt.Sprintf("gm.Setup: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	for !game.IsOver() {
		uiDisplay = transform(game.Progress())
		out, err := json.Marshal(uiDisplay)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		conn.WriteMessage(websocket.TextMessage, out)
		_, msg, err := conn.ReadMessage()

		if err != nil {
			fmt.Printf("err conn.ReadMessage: %s\n", err.Error())
			return
		}
		if len(msg) == 0 {
			return
		}
		game.Update(rune(msg[0]))
	}
	uiDisplay = transform(game.Result())
	out, err := json.Marshal(uiDisplay)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	conn.WriteMessage(websocket.TextMessage, out)
}

type Display struct {
	Secret           string `json:"secret"`
	GuessedChars     string `json:"guessedChars"`
	Message          string `json:"messages"`
	RemainingGuesses int    `json:"remainingGuesses`
}

func transform(state gm.State) (d Display) {
	d.Secret = strings.Join(state.Secret, "")
	gc := []string{}
	for _, ch := range state.GuessedChars {
		gc = append(gc, string(ch))
	}
	d.GuessedChars = strings.Join(gc, ",")
	fmt.Println(state.Message)
	d.RemainingGuesses = state.RemainingGuesses
	d.Message = state.Message
	return
}

func page(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, "ws://"+r.Host+"/play")
}

var (
	upgrader = websocket.Upgrader{}
	tmpl     = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var debug = document.getElementById("debug");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        debug.appendChild(d);
    };
        if (!ws) {
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        }
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="secret"></div>
</td></tr>
<tr>
<td>
<div id="remainingGuesses"></div>
</td>
<tr>
<td>
<div id="debug"></div>
</td>
</tr>
</tr>
<tr>
<td>
<div id="guessedChars"></div>
</td>
</tr>
</table>
</body>
</html>
`))
)
