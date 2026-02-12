//Danny Radosevich

//all the color definitions and functions for coloring output

package commands

import (
	"fmt"
	"io"
)

// List of colors
type color string

const (
	Reset   color = "\x1b[0m"
	Red     color = "\x1b[31m"
	Green   color = "\x1b[32m"
	Yellow  color = "\x1b[33m"
	Blue    color = "\x1b[34m"
	Magenta color = "\x1b[35m"
	Cyan    color = "\x1b[36m"
	White   color = "\x1b[37m"
	Gold    color = "\x1b[38;2;255;196;47m" // #FFC42F
)

func (c color) ColorStart(w io.Writer) {
	fmt.Fprint(w, c)
}

func (c color) ColorEnd(w io.Writer) {
	fmt.Fprint(w, Reset)
}
