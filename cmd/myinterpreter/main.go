package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

type State int
type Transition struct {
	Lexema string
	State  State
}

const (
	NORMAL State = iota
	WAITING_EQUAL
	WAITING_COMMENT
	WAITING_STRING
	WAITING_NUMBER
	IGNORE
	BREAK
)

func getTransitions() map[State]map[rune]Transition {
	transitions := map[State]map[rune]Transition{
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
			'=':  {"EQUAL", WAITING_EQUAL},
			'<':  {"LESS", WAITING_EQUAL},
			'>':  {"GREATER", WAITING_EQUAL},
			'!':  {"BANG", WAITING_EQUAL},
			'/':  {"SLASH", WAITING_COMMENT},
			' ':  {"SPACE", IGNORE},
			'\t': {"TAB", IGNORE},
			'"':  {"QUOTATION_MARK", WAITING_STRING},
		},
		WAITING_EQUAL: {
			'=': {"EQUAL", NORMAL},
		},
		WAITING_COMMENT: {
			'/': {"COMMENT", BREAK},
		},
		WAITING_STRING: {
			'"': {"STRING", NORMAL},
		},
		WAITING_NUMBER: {
			'.': {"DOT", WAITING_NUMBER},
		},
	}

	for i := '0'; i <= '9'; i++ {
		transitions[NORMAL][i] = Transition{"NUMBER", WAITING_NUMBER}
		transitions[WAITING_NUMBER][i] = Transition{"NUMBER", WAITING_NUMBER}
	}

	return transitions
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
			literal := ""
			args := "null"
			for _, curr_chr := range line {
				transition, ok := transitions[state][curr_chr]
				if !ok && state == WAITING_STRING {
					literal += string(curr_chr)
					continue
				}

				if state == WAITING_NUMBER {
					if unicode.IsDigit(curr_chr) || curr_chr == '.' {
						literal += string(curr_chr)
						continue
					}
					if literal != "" {
						if strings.Count(literal, ".") > 1 {
							errors = true
							fmt.Fprintf(os.Stderr, "[line %d] Error: Invalid number format: %s\n", lineNumber+1, literal)
						} else {
							// number, err := strconv.ParseFloat(literal, 64)
							// if err != nil {
							// 	fmt.Fprintf(os.Stderr, "[line %d] Error: Invalid number\n", lineNumber+1)
							// }
							fmt.Printf("NUMBER %s %s\n", literal, literal)
						}
						last_transition = Transition{"", NORMAL}
						last_chr = rune(0)
						state = NORMAL
						literal = ""
					}
					transition, ok = transitions[NORMAL][curr_chr]
				}

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

				if state == WAITING_NUMBER {
					if unicode.IsDigit(curr_chr) || curr_chr == '.' {
						literal += string(curr_chr)
						continue
					}
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
					if literal != "" {
						lexema = ""
						chr = "\"" + literal + "\""
						args = literal
					}
					lexema = lexema + transition.Lexema
					fmt.Printf("%s %s %s\n", lexema, chr, args)
					last_transition = Transition{"", NORMAL}
					last_chr = rune(0)
					state = NORMAL
					literal = ""
					args = "null"
				}
			}
			if literal != "" {
				if state == WAITING_NUMBER {
					if strings.Count(literal, ".") > 1 {
						errors = true
						fmt.Fprintf(os.Stderr, "[line %d] Error: Invalid number format: %s\n", lineNumber+1, literal)
					} else {
						// number, err := strconv.ParseFloat(literal, 64)
						// if err != nil {
						// 	fmt.Fprintf(os.Stderr, "[line %d] Error: Invalid number\n", lineNumber+1)
						// }
						fmt.Printf("NUMBER %s %s\n", literal, literal)
					}
					continue
				}
				errors = true
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", lineNumber+1)
				continue
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
