package commands

import (
	"fmt"
	"io"
)

type BasicStack struct {
	stack []string
}

func (s *BasicStack) Push(item string) {
	s.stack = append(s.stack, item)
}

func (s *BasicStack) Pop() (string, bool) {
	if len(s.stack) == 0 {
		return "", false
	}
	item := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return item, true
}

func (s *BasicStack) GetName() string {
	return "stack"
}

func (s *BasicStack) GetHelp() string {
	return "A simple stack implementation. Usage: stack push <item> | stack pop"
}

func (s *BasicStack) isValid(cmd string, args []string) bool {
	if cmd != "stack" {
		return false
	}
	if len(args) < 1 {
		return false
	}
	switch args[0] {
	case "push":
		return len(args) == 2
	case "pop":
		return len(args) == 1
	default:
		return false
	}
}

func (s *BasicStack) Execute(r io.Reader, w io.Writer, args ...string) bool {
	if !s.isValid("stack", args) {
		fmt.Fprintln(w, "Invalid command. Usage: stack push <item> | stack pop")
		return false
	}
	switch args[0] {
	case "push":
		s.Push(args[1])
		fmt.Fprintf(w, "Pushed: %s\n", args[1])
	case "pop":
		item, ok := s.Pop()
		if !ok {
			fmt.Fprintln(w, "Stack is empty")
			return false
		}
		fmt.Fprintf(w, "Popped: %s\n", item)
	}
	return true
}
