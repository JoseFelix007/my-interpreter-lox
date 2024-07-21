package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
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
	tokens       []Token
	errors       []ScanError
	transitions  map[State]map[rune]Transition
	keywords     map[string]string
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
	IDENTIFIER     = "IDENTIFIER"
	//Multiple
	LESS_EQUAL    = "LESS_EQUAL"
	GREATER_EQUAL = "GREATER_EQUAL"
	BANG_EQUAL    = "BANG_EQUAL"
	EQUAL_EQUAL   = "EQUAL_EQUAL"
	//KeyWords
	AND    = "AND"
	CLASS  = "CLASS"
	ELSE   = "ELSE"
	FALSE  = "FALSE"
	FOR    = "FOR"
	FUN    = "FUN"
	IF     = "IF"
	NIL    = "NIL"
	OR     = "OR"
	PRINT  = "PRINT"
	RETURN = "RETURN"
	SUPER  = "SUPER"
	THIS   = "THIS"
	TRUE   = "TRUE"
	VAR    = "VAR"
	WHILE  = "WHILE"
	EOF    = "EOF"
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

func getKeywords() map[string]string {
	return map[string]string{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"for":    FOR,
		"fun":    FUN,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
	}
}

func NewScanner(fileContents []byte) *Scanner {
	return &Scanner{
		currentState: NORMAL,
		transitions:  getTransitions(),
		keywords:     getKeywords(),

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
	hasStartDigits := false
	// hasEndDigits := false
	for unicode.IsDigit(s.peek()) {
		hasStartDigits = true
		s.read()
	}
	if s.peek() == '.' {
		s.read()
		if !hasStartDigits {
			s.addToken(string('.'), "", s.prevCursorChar, s.cursorLine, DOT)
			return
		}
		if hasStartDigits && !unicode.IsDigit(s.peek()) {
			s.cursorChar--
		}
		for unicode.IsDigit(s.peek()) {
			// hasEndDigits = true
			s.read()
		}
	}

	line := s.lines[s.cursorLine]
	lexema := line[s.prevCursorChar:s.cursorChar]
	literalFloat, err := strconv.ParseFloat(lexema, 64)
	if err != nil {
		s.addError(s.cursorLine, "Invalid number format", "")
	} else {
		prec := -1
		if literalFloat == float64(int(literalFloat)) {
			prec = 1
		}
		literal := strconv.FormatFloat(literalFloat, 'f', prec, 64)
		s.addToken(lexema, literal, s.prevCursorChar, s.cursorLine, NUMBER)
	}
}

func (s *Scanner) scanString() {
	transition, ok := s.transitions[s.currentState][s.peek()]
	for !s.isEOL() && !ok {
		s.read()
		transition, ok = s.transitions[s.currentState][s.peek()]
	}
	if s.isEOL() && !ok {
		s.addError(s.cursorLine, "Unterminated string", "")
	}
	if ok {
		s.currentState = transition.State
		s.read()
		line := s.lines[s.cursorLine]
		lexema := line[s.prevCursorChar:s.cursorChar]
		literal := line[s.prevCursorChar+1 : s.cursorChar-1]
		s.addToken(lexema, literal, s.prevCursorChar, s.cursorLine, STRING)
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

func isAlpha(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || unicode.IsDigit(r)
}

func (s *Scanner) scanIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.read()
	}

	line := s.lines[s.cursorLine]
	lexema := line[s.prevCursorChar:s.cursorChar]

	tokenType := IDENTIFIER
	keyword, exists := s.keywords[lexema]
	if exists {
		tokenType = keyword
	}
	s.addToken(lexema, "", s.prevCursorChar, s.cursorLine, tokenType)
}

func (s *Scanner) scanTokens() error {
	if len(s.fileContents) == 0 {
		s.addToken("", "", 0, 0, EOF)
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

			if isAlpha(s.peek()) {
				s.prevCursorChar = s.cursorChar
				s.scanIdentifier()
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

	s.addToken("", "", 0, s.cursorLine+1, EOF)

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
}
