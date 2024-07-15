package main

import (
	"fmt"
	"os"
	"strings"
)

type State int
type Transition struct {
	Lexema string
	State  State
}

const (
	NORMAL State = iota
	WAITING
	WAITING_COMMENT
	IGNORE
	BREAK
)

func getTransitions() map[State]map[rune]Transition {
	return map[State]map[rune]Transition{
		NORMAL: {
			'(':  {"LEFT_PAREN", NORMAL},
			')':  {"RIGHT_PAREN", NORMAL},
			'{':  {"LEFT_BRACE", NORMAL},
			'}':  {"RIGHT_BRACE", NORMAL},
			'*':  {"STAR", NORMAL},
			'.':  {"DOT", NORMAL},
			',':  {"COMMA", NORMAL},
			';':  {"SEMICOLON", NORMAL},
			'+':  {"PLUS", NORMAL},
			'-':  {"MINUS", NORMAL},
			'=':  {"EQUAL", WAITING},
			'<':  {"LESS", WAITING},
			'>':  {"GREATER", WAITING},
			'!':  {"BANG", WAITING},
			'/':  {"SLASH", WAITING_COMMENT},
			' ':  {"SPACE", IGNORE},
			'\t': {"TAB", IGNORE},
		},
		WAITING: {
			'=': {"EQUAL", NORMAL},
		},
		WAITING_COMMENT: {
			'/': {"COMMENT", BREAK},
		},
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
		transitions := getTransitions()
		lines := strings.Split(string(fileContents), "\n")
		for lineNumber, line := range lines {
			state := NORMAL
			last_transition := Transition{"", NORMAL}
			last_chr := rune(0)
			for _, curr_chr := range line {
				transition, ok := transitions[state][curr_chr]
				if !ok && state != NORMAL {
					fmt.Printf("%s %c null\n", last_transition.Lexema, last_chr)
					last_transition = Transition{"", NORMAL}
					last_chr = rune(0)
					state = NORMAL
					transition, ok = transitions[NORMAL][curr_chr]
				}

				if !ok {
					errors = true
					fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", lineNumber+1, rune(curr_chr))
					last_transition = Transition{"", NORMAL}
					last_chr = rune(0)
					state = NORMAL
					continue
				}

				state = transition.State
				if state == BREAK {
					last_transition = Transition{"", NORMAL}
					last_chr = rune(0)
					state = NORMAL
					break
				}

				if state == IGNORE {
					last_transition = Transition{"", NORMAL}
					last_chr = rune(0)
					state = NORMAL
					continue
				}

				if state != NORMAL {
					last_chr = curr_chr
					last_transition = transition
					continue
				}

				if state == NORMAL {
					chr := string(last_chr) + string(curr_chr)
					if last_chr == rune(0) {
						chr = string(curr_chr)
					}
					lexema := ""
					if last_transition.Lexema != "" {
						lexema = last_transition.Lexema + "_"
					}
					lexema = lexema + transition.Lexema
					fmt.Printf("%s %s null\n", lexema, chr)
					last_transition = Transition{"", NORMAL}
					last_chr = rune(0)
					state = NORMAL
				}
			}
			if last_transition.Lexema != "" {
				fmt.Printf("%s %c null\n", last_transition.Lexema, last_chr)
			}
		}
	}

	fmt.Println("EOF  null")

	if errors {
		os.Exit(65)
	}
}
