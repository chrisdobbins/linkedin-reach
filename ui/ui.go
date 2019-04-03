package ui

import (
	tm "github.com/buger/goterm"
)

type Display struct {
	Secret       []byte
	GuessedChars []byte
	Messages     [][]byte
}

func (d Display) Write() {
	tm.Clear()
	tm.MoveCursor(1, 1)
	box := tm.NewBox(60|tm.PCT, 10, 0)
	box.Write(d.format())
	tm.Print(box.String())
	tm.MoveCursor(3, 9)
	tm.Flush()
}

func (d Display) format() (ret []byte) {
	ret = []byte{}
	newLine := []byte("\n")
	ret = append(ret, append(d.Secret, newLine...)...)
	ret = append(ret, append(d.GuessedChars, newLine...)...)
	for _, msgs := range d.Messages {
		ret = append(ret, append(msgs, newLine...)...)
	}
	return
}
