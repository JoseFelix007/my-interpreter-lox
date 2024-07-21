package main

import (
	"fmt"
	"os"
)

type CommandFunc func(*Interpreter) error

var commands = map[string]CommandFunc{
	"tokenize": func(lox *Interpreter) error {
		err := lox.tokenize()
		lox.scanner.printTokens()
		return err
	},
	"parse": func(lox *Interpreter) error {
		err := lox.parse()
		lox.parser.print()
		return err
	},
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	debug("Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" && command != "parse" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	lox := NewInterpreter()
	commandFunc, ok := commands[os.Args[1]]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	err := lox.readFile(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	err = commandFunc(lox)
	if err != nil {
		os.Exit(65)
	}
}
