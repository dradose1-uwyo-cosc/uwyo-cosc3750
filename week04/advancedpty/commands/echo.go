//Danny Radosevich
//basic echo

package commands

import "io"

func Echo(args []string, w io.Writer) {
	for _, arg := range args {
		w.Write([]byte(arg + " "))
	}
	w.Write([]byte("\n"))
}
