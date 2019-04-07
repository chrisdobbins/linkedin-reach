# Hangman (or something like it) #

[![Build Status](https://travis-ci.com/chrisdobbins/linkedin-reach.svg?branch=master)](https://travis-ci.com/chrisdobbins/linkedin-reach)

This is a word-guessing game similar to hangman. 

Rules:
* You are allowed a certain number of guesses for a word. Each word is guaranteed to be guessable within the allowed number of guesses.
* Each guess must be an ASCII letter; all other inputs will be rejected, though they will not count against your remaining guesses. Guesses are case-insensitive.

Good luck and have fun!

Basic options:
* `-h`, `--help`: Brings up this message
* `-guesses`, `--guesses`: Configures the maximum allowed number of guesses. Default is 6.
* `-serve`, `--serve`: Starts web app instead of terminal app. Default port is 8080, but can be changed by setting the `PORT` environment variable.
