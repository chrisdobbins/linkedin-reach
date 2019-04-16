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

This game is currently deployed at [http://limitless-bastion-18713.herokuapp.com/](http://limitless-bastion-18713.herokuapp.com/)*.

## Installation ##

The easiest way is to download the most recent release for your OS. You can also build from source by cloning this repo and running `go build`, assuming Go is already installed.

### Raspberry Pi ###

For the Raspberry Pi, there is a setup script because this is meant to be played on a 3.5" display. This script installs the ARM binary in `/usr/local/bin`, then configures the Pi to use the display, in that order. The display driver installation is the last part of the script because the OS will restart after installation is complete.

 `wpa_supplicant` is a sample config file meant to go into Raspbian's `/boot/` directory after it has been burned onto the SD card and while the SD card is still mounted. The Pi will automatically move this file into `/etc/wpa_supplicant/` on booting, thereby causing it to connect to the specified wifi network automatically upon booting. To enable headless access via ssh, simply create a file named `ssh` in `/boot/`.

***

*http is not a typo. I realize that this is not a good practice and I would never do this for a production-ready application.
