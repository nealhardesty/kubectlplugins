package util

import (
	"io"
)

var (
	esc       = "\033["
	clearLine = []byte(esc + "2K\r")
	moveUp    = []byte(esc + "1A")
	moveDown  = []byte(esc + "1B")
)

// ClearLine erases the current terminal line
func ClearLine(out io.Writer) {
	out.Write(clearLine)
}
