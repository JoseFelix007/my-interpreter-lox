package main

import (
	"fmt"
	"os"
	"strings"
)

func getLexemes() map[rune]string {
	return map[rune]string{
		'(': "LEFT_PAREN",
		')': "RIGHT_PAREN",
		'{': "LEFT_BRACE",
		'}': "RIGHT_BRACE",
		'*': "STAR",
		'.': "DOT",
		',': "COMMA",
		';': "SEMICOLON",
		'+': "PLUS",
		'-': "MINUS",
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	// Uncomment this block to pass the first stage

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	errors := false
	if len(fileContents) > 0 {
		lexemes := getLexemes()
		lines := strings.Split(string(fileContents), "\n")
		for lineNumber, line := range lines {
			for _, chr := range line {
				lexeme, ok := lexemes[chr]
				if ok {
					fmt.Printf("%s %c null\n", lexeme, rune(chr))
				} else {
					errors = true
					fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", lineNumber+1, rune(chr))
				}
			}
		}
	}

	fmt.Println("EOF  null")

	if errors {
		os.Exit(65)
	}
}
