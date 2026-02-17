//Danny Radosevich

//extended command example

package commands

import (
	"fmt"
	"io"
	"strings"
)

// create a Command interface that all commands will implement
// It is exported so it can be used by other packages (like the main package)
type Command interface {
	GetName() string
	GetHelp() string
	Execute(r io.Reader, w io.Writer, args ...string) bool //variadic args for command parameters
	//variadic means it can take any number of string arguments, which will be passed as a slice
}

//now need to specify the commands that will implement this interface, such as echo, color, etc.

// first an error for un registered commands
var ErrDupeCommand = fmt.Errorf("duplicate  command")

// a command slice
var commands []Command

//now we need to register the commands we want to use, such as echo and color, or any future

func init() {
	//register the echo command
	RegisterCommand(Base{
		Name: "echo",
		Help: "Echoes the arguments back to the user",
		Action: func(r io.Reader, w io.Writer, args ...string) bool {
			Echo(args, w)
			return true
		},
	})
	//register the color command
	RegisterCommand(Base{
		Name: "color",
		Help: "Changes the color of the text. Usage: color <start|stop> [color name]",
		Action: func(r io.Reader, w io.Writer, args ...string) bool {
			fmt.Printf("args: %v\n", args)
			if len(args) < 1 {
				fmt.Fprintln(w, "Usage: color <start|stop>[color name]")
				return false
			}
			if strings.ToLower(args[0]) == "stop" {
				Reset.ColorEnd(w)
				return true
			} else {

				switch strings.ToLower(args[1]) {
				case "red":
					Red.ColorStart(w)
				case "green":
					Green.ColorStart(w)
				case "yellow":
					Yellow.ColorStart(w)
				case "blue":
					Blue.ColorStart(w)
				case "magenta":
					Magenta.ColorStart(w)
				case "cyan":
					Cyan.ColorStart(w)
				case "white":
					White.ColorStart(w)
				case "gold":
					Gold.ColorStart(w)
				default:
					fmt.Fprintf(w, "Unknown color: %s\n", args[0])
				}
			}
			return true
		},
	})
}

func RegisterCommand(c Command) error {
	//get command name
	name := c.GetName()
	//ensure not already a command
	for index, com := range commands {
		//unique commands in order
		switch strings.Compare(name, com.GetName()) {
		case 0:
			//command already exists, return our new error
			return ErrDupeCommand
		case 1:
			//can append the command to it's place in the list, which is sorted by name
			commands = append(commands[:index], append([]Command{c}, commands[index:]...)...)
		case -1:
			//keep looking for the right place to insert the command
			continue

		}

	}
	commands = append(commands, c)
	//if we get here, the command is unique and can be added to the end of the list

	return nil
}

func GetCommand(name string) Command {
	for _, c := range commands {
		if c.GetName() == name {
			return c
		}
	}
	return nil
}
