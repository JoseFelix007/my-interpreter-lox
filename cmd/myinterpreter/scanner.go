package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

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
	fileContents   []byte
	tokens         []Token
	errors         []ScanError
	transitions    map[State]map[rune]Transition
	currentState   State
	lastTransition Transition
	literal        string
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
		fileContents:   fileContents,
		currentState:   NORMAL,
		lastTransition: Transition{rune(0), "", NORMAL},
		literal:        "",
		transitions:    getTransitions(),
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
		fmt.Fprintf(os.Stderr, "[line %d] Error: %s: %s\n", err.Line, err.Message, err.Value)
	} else {
		fmt.Fprintf(os.Stderr, "[line %d] Error: %s.\n", err.Line, err.Message)
	}
}

func (s *Scanner) clearScanValues() {
	s.currentState = NORMAL
	s.lastTransition = Transition{rune(0), "", NORMAL}
	s.literal = ""
}

func (s *Scanner) scanTokens() error {
	if len(s.fileContents) == 0 {
		return nil
	}

	lines := strings.Split(string(s.fileContents), "\n")
	for lineNumber, line := range lines {
		s.clearScanValues()
		for i, currentChar := range line {
			transition, ok := s.transitions[s.currentState][currentChar]
			if !ok && s.currentState == WAITING_STRING {
				s.literal += string(currentChar)
				continue
			}

			if s.currentState == WAITING_NUMBER {
				if unicode.IsDigit(currentChar) || currentChar == '.' {
					s.literal += string(currentChar)
					continue
				}
				if s.literal != "" {
					if strings.Count(s.literal, ".") > 1 {
						s.addError(lineNumber+1, "Invalid number format", s.literal)
					} else {
						s.addToken(s.literal, s.literal, i, lineNumber+1, NUMBER)
					}
					s.clearScanValues()
				}
				transition, ok = s.transitions[NORMAL][currentChar]
			}

			if !ok && s.currentState != NORMAL {
				s.addToken(string(s.lastTransition.Char), "", i, lineNumber+1, s.lastTransition.Type)
				s.clearScanValues()
				transition, ok = s.transitions[NORMAL][currentChar]
			}

			if !ok {
				s.addError(lineNumber+1, "Unexpected character", string(currentChar))
				s.clearScanValues()
				continue
			}

			s.currentState = transition.State
			if s.currentState == BREAK {
				s.clearScanValues()
				break
			}

			if s.currentState == IGNORE {
				s.clearScanValues()
				continue
			}

			if s.currentState == WAITING_NUMBER {
				if unicode.IsDigit(currentChar) || currentChar == '.' {
					s.literal += string(currentChar)
					continue
				}
			}

			if s.currentState != NORMAL {
				s.lastTransition = transition
				continue
			}

			if s.currentState == NORMAL {
				lexema := string(s.lastTransition.Char) + string(currentChar)
				if s.lastTransition.Char == rune(0) {
					lexema = string(currentChar)
				}
				tokenType := ""
				if s.lastTransition.Type != "" {
					tokenType = s.lastTransition.Type + "_"
				}
				if s.lastTransition.State == WAITING_STRING {
					tokenType = ""
					lexema = "\"" + s.literal + "\""
				}
				tokenType = tokenType + transition.Type
				s.addToken(lexema, s.literal, i, lineNumber+1, tokenType)
				s.clearScanValues()
			}
		}
		if s.literal != "" {
			if s.currentState == WAITING_NUMBER {
				if strings.Count(s.literal, ".") > 1 {
					s.addError(lineNumber+1, "Invalid number format", s.literal)
				} else {
					s.addToken(s.literal, s.literal, len(line), lineNumber+1, NUMBER)
				}
				continue
			}
			s.addError(lineNumber+1, "Unterminated string", "")
			continue
		}
		if s.lastTransition.Type != "" {
			s.addToken(string(s.lastTransition.Char), "", len(line), lineNumber+1, s.lastTransition.Type)
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
