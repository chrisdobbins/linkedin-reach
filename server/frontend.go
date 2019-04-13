package server

const frontend = `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
body {
font-family: helvetica, arial, sans-serif;
    font-size: 30px;
    color: aqua;
    background-color: #01013F;
}
</style>
</head>
<body>
<div style="height: 50vh; width: 50vw; margin: 0 auto;" id="board">
  <div id="state">
    <div id="secret" style="height: 40%; width: 100%; font-size: 20px;"></div>
    <div id="remainingGuesses" style="height: 30%; width: 30%"></div>
    <div id="guessedChars"></div>
    <form id="inputform">
      <input id="input" type="text" placeholder="Guess a letter">
      <button id="send">Submit Guess</button>
    </form>
  </div>
  <div id="message"></div>
</div>
<script>  
window.addEventListener("load", function(evt) {
    var input = document.getElementById("input");
    var ws;
    var print = function(parsedMsg) {
       document.querySelector("#secret").innerHTML = parsedMsg.secret.split("").join(" ");
        document.querySelector("#guessedChars").innerHTML = "Guessed letters: "+ parsedMsg.guessedChars;
       document.querySelector("#remainingGuesses").innerHTML = "Remaining guesses: " + parsedMsg.remainingGuesses;
       document.querySelector("#message").innerHTML = parsedMsg.message;
    };
        if (!ws) {
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
        }
        ws.onclose = function(evt) {
            ws = null;
        }
        ws.onmessage = function(evt) {
           document.querySelector("#input").value = "";
           var data = JSON.parse(evt.data);
           if (!data) {
              return
           }
           print(data);
           if (data.gameOver) {
             var state = document.querySelector("#state"); 
             var form = document.querySelector("#inputform");
             var remainingGuesses = document.querySelector("#remainingGuesses");
             var secret = document.querySelector("#secret");
             var guessedChars = document.querySelector("#guessedChars");
             var nodes = [form, remainingGuesses, secret, guessedChars];
             nodes.forEach(function(node){
               state.removeChild(node);
             });
           }
        }
        ws.onerror = function(evt) {
        }
        }
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        if (input.value && input.value.length > 0) {
        ws.send(input.value);
        }
        return false;
    };
});
</script>
</body>
</html>
`
