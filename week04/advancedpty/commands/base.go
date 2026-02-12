package commands

import "io"

type Base struct {
	Name, Help string
	Action     func(input io.Reader, output io.Writer, args ...string) bool
}

func (b Base) String() string { return b.Name }

// String returns the name of the command,
// so that it can be printed in a list of commands

func (b Base) GetName() string { return b.Name }

// GetName returns the name of the command,
// so that it can be used to identify the command when executing

func (b Base) GetHelp() string { return b.Help }

// GetHelp returns the help string for the command,
// so that it can be printed in a list of commands or when the user asks for help

func (b Base) Execute(r io.Reader, w io.Writer, args ...string) bool {
	return b.Action(r, w, args...)
}
