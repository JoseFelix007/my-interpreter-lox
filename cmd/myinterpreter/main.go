package main

import (
	"fmt"
	"os"

	"errors"
	"strconv"
	"strings"
	"unicode"
)

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

	scanner := NewScanner(fileContents)
	err = scanner.scanTokens()
	scanner.printTokens()
	if err != nil {
		os.Exit(65)
	}
}

type Token struct {
	Lexema  string
	Literal string
	Num     int
	Line    int
	Type    string
}
type ScanError struct {
	Line    int
	Message string
	Value   string
}
type Scanner struct {
	tokens       []Token
	errors       []ScanError
	transitions  map[State]map[rune]Transition
	currentState State

	fileContents   []byte
	lines          []string
	prevCursorChar int
	cursorChar     int
	cursorLine     int
}
type State int
type Transition struct {
	Char  rune
	Type  string
	State State
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
const (
	LEFT_PAREN     = "LEFT_PAREN"
	RIGHT_PAREN    = "RIGHT_PAREN"
	LEFT_BRACE     = "LEFT_BRACE"
	RIGHT_BRACE    = "RIGHT_BRACE"
	STAR           = "STAR"
	DOT            = "DOT"
	COMMA          = "COMMA"
	SEMICOLON      = "SEMICOLON"
	PLUS           = "PLUS"
	MINUS          = "MINUS"
	EQUAL          = "EQUAL"
	LESS           = "LESS"
	GREATER        = "GREATER"
	BANG           = "BANG"
	SLASH          = "SLASH"
	SPACE          = "SPACE"
	TAB            = "TAB"
	QUOTATION_MARK = "QUOTATION_MARK"
	COMMENT        = "COMMENT"
	STRING         = "STRING"
	NUMBER         = "NUMBER"
)

func getTransitions() map[State]map[rune]Transition {
	transitions := map[State]map[rune]Transition{
		NORMAL: {
			'(':  {'(', LEFT_PAREN, NORMAL},
			')':  {')', RIGHT_PAREN, NORMAL},
			'{':  {'{', LEFT_BRACE, NORMAL},
			'}':  {'}', RIGHT_BRACE, NORMAL},
			'*':  {'*', STAR, NORMAL},
			'.':  {'.', DOT, NORMAL},
			',':  {',', COMMA, NORMAL},
			';':  {';', SEMICOLON, NORMAL},
			'+':  {'+', PLUS, NORMAL},
			'-':  {'-', MINUS, NORMAL},
			'=':  {'=', EQUAL, WAITING_EQUAL},
			'<':  {'<', LESS, WAITING_EQUAL},
			'>':  {'>', GREATER, WAITING_EQUAL},
			'!':  {'!', BANG, WAITING_EQUAL},
			'/':  {'/', SLASH, WAITING_COMMENT},
			' ':  {' ', SPACE, IGNORE},
			'\t': {'\t', TAB, IGNORE},
			'"':  {'"', QUOTATION_MARK, WAITING_STRING},
		},
		WAITING_EQUAL: {
			'=': {'=', EQUAL, NORMAL},
		},
		WAITING_COMMENT: {
			'/': {'/', COMMENT, BREAK},
		},
		WAITING_STRING: {
			'"': {'"', STRING, NORMAL},
		},
		WAITING_NUMBER: {
			'.': {'.', DOT, WAITING_NUMBER},
		},
	}

	for i := '0'; i <= '9'; i++ {
		transitions[NORMAL][i] = Transition{rune(i), NUMBER, WAITING_NUMBER}
		transitions[WAITING_NUMBER][i] = Transition{rune(i), NUMBER, WAITING_NUMBER}
	}

	return transitions
}

func NewScanner(fileContents []byte) *Scanner {
	return &Scanner{
		currentState: NORMAL,
		transitions:  getTransitions(),

		fileContents:   fileContents,
		lines:          strings.Split(string(fileContents), "\n"),
		prevCursorChar: 0,
		cursorChar:     0,
		cursorLine:     0,
	}
}

func (s *Scanner) addToken(lexema, literal string, num, line int, tokenType string) {
	s.tokens = append(s.tokens, Token{lexema, literal, num, line, tokenType})
}

func (s *Scanner) printToken(token Token) {
	if token.Literal == "" {
		token.Literal = "null"
	}
	fmt.Printf("%s %s %s\n", token.Type, token.Lexema, token.Literal)
}

func (s *Scanner) addError(line int, message, value string) {
	s.errors = append(s.errors, ScanError{line, message, value})
}

func (s *Scanner) printError(err ScanError) {
	if err.Value != "" {
		fmt.Fprintf(os.Stderr, "[line %d] Error: %s: %s\n", err.Line+1, err.Message, err.Value)
	} else {
		fmt.Fprintf(os.Stderr, "[line %d] Error: %s.\n", err.Line+1, err.Message)
	}
}

func (s *Scanner) read() rune {
	if s.isEOL() {
		return rune(0)
	}
	chr := s.lines[s.cursorLine][s.cursorChar]
	s.cursorChar++
	return rune(chr)
}

func (s *Scanner) peek() rune {
	if s.isEOL() {
		return rune(0)
	}
	chr := s.lines[s.cursorLine][s.cursorChar]
	return rune(chr)
}

func (s *Scanner) isEOL() bool {
	return s.cursorChar >= len(s.lines[s.cursorLine])
}

func (s *Scanner) isWaiting() bool {
	return s.currentState == WAITING_STRING ||
		s.currentState == WAITING_EQUAL ||
		s.currentState == WAITING_COMMENT ||
		s.currentState == WAITING_NUMBER
}

func (s *Scanner) scanNumber() {
	for unicode.IsDigit(s.peek()) {
		s.read()
	}
	if s.peek() == '.' {
		s.read()
		if !unicode.IsDigit(s.peek()) {
			s.addToken(string('.'), string('.'), s.prevCursorChar, s.cursorLine, DOT)
			return
		}
		for unicode.IsDigit(s.peek()) {
			s.read()
		}
	}

	line := s.lines[s.cursorLine]
	literalFloat, err := strconv.ParseFloat(line[s.prevCursorChar:s.cursorChar], 64)
	if err != nil {
		s.addError(s.cursorLine, "Invalid number format", "")
	} else {
		literal := strconv.FormatFloat(literalFloat, 'f', -1, 64)
		s.addToken(literal, literal, s.prevCursorChar, s.cursorLine, NUMBER)
	}
}

func (s *Scanner) scanString() {
	transition, ok := s.transitions[s.currentState][s.peek()]
	for !s.isEOL() && !ok {
		s.read()
		s.currentState = transition.State
		transition, ok = s.transitions[s.currentState][s.peek()]
	}
	if s.isEOL() && !ok {
		s.addError(s.cursorLine, "Unterminated string", "")
	}
	if ok {
		s.read()
		line := s.lines[s.cursorLine]
		literal := line[s.prevCursorChar:s.cursorChar]
		s.addToken(literal, literal, s.prevCursorChar, s.cursorLine, STRING)
	}
}

func (s *Scanner) scanToken(initialTransition Transition) {
	tokenType := initialTransition.Type
	s.prevCursorChar = s.cursorChar - 1
	for s.isWaiting() {
		nextTransition, ok := s.transitions[s.currentState][s.peek()]
		if !ok {
			break
		}

		s.read()
		s.currentState = nextTransition.State
		tokenType += "_" + nextTransition.Type
	}

	line := s.lines[s.cursorLine]
	literal := line[s.prevCursorChar:s.cursorChar]
	s.addToken(literal, "", s.prevCursorChar, s.cursorLine, tokenType)
	s.currentState = NORMAL
}

func (s *Scanner) scanTokens() error {
	if len(s.fileContents) == 0 {
		return nil
	}

	for lineIdx, line := range s.lines {
		s.cursorChar = 0
		s.cursorLine = lineIdx

		s.currentState = NORMAL
		s.prevCursorChar = s.cursorChar
		for s.cursorChar < len(line) {
			if s.currentState == WAITING_STRING {
				s.prevCursorChar = s.cursorChar - 1
				s.scanString()
				continue
			}

			if unicode.IsDigit(s.peek()) || s.peek() == '.' {
				s.prevCursorChar = s.cursorChar
				s.scanNumber()
				continue
			}

			chr := s.read()
			transition, ok := s.transitions[s.currentState][chr]

			if !ok {
				s.addError(s.cursorLine, "Unexpected character", string(chr))
				s.currentState = NORMAL
				continue
			}

			s.currentState = transition.State
			if s.currentState == BREAK {
				break
			}

			if s.currentState == IGNORE {
				s.currentState = NORMAL
				continue
			}

			if s.currentState == WAITING_STRING {
				continue
			}

			if s.currentState == WAITING_COMMENT {
				_, ok = s.transitions[s.currentState][s.peek()]
				if ok {
					continue
				}
			}

			s.scanToken(transition)
		}
	}

	if len(s.errors) > 0 {
		return errors.New("")
	}
	return nil
}

func (s *Scanner) printTokens() {
	for _, err := range s.errors {
		s.printError(err)
	}
	for _, token := range s.tokens {
		s.printToken(token)
	}

	fmt.Println("EOF  null")
}
