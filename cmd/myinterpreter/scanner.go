package main

import (
	"fmt"
	"os"
	"unicode"
)

type ScanError struct {
	Line    int
	Message string
	Value   string
}

func (err *ScanError) print() {
	if err.Value != "" {
		fmt.Fprintf(os.Stderr, "[line %d] Error: %s: %s\n", err.Line+1, err.Message, err.Value)
	} else {
		fmt.Fprintf(os.Stderr, "[line %d] Error: %s.\n", err.Line+1, err.Message)
	}
}

type Scanner struct {
	tokens []Token
	errors []ScanError

	fileContents []byte
	start        int
	current      int
	line         int
}

func NewScanner(fileContents []byte) *Scanner {
	return &Scanner{
		fileContents: fileContents,
		start:        0,
		current:      0,
		line:         0,
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.fileContents)
}

func (s *Scanner) read() rune {
	chr := s.fileContents[s.current]
	s.current++
	return rune(chr)
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}
	chr := s.fileContents[s.current]
	return rune(chr)
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.fileContents) {
		return rune(0)
	}
	chr := s.fileContents[s.current+1]
	return rune(chr)
}

func (s *Scanner) addToken(Type string) {
	lexema := string(s.fileContents[s.start:s.current])
	s.tokens = append(s.tokens, Token{lexema, "", s.line, Type})
}

func (s *Scanner) addError(message, value string) {
	s.errors = append(s.errors, ScanError{s.line, message, value})
}

func (s *Scanner) addTokenWithLiteral(Type string, literal string) {
	lexema := string(s.fileContents[s.start:s.current])
	s.tokens = append(s.tokens, Token{lexema, literal, s.line, Type})
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if rune(s.fileContents[s.current]) != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) scanString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.read()
	}
	if s.isAtEnd() {
		s.addError("Unterminated string", "")
		return
	}

	s.read()

	literal := string(s.fileContents[s.start+1 : s.current-1])
	s.addTokenWithLiteral(STRING, literal)
}

func (s *Scanner) scanNumber() {
	for unicode.IsDigit(s.peek()) {
		s.read()
	}
	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		s.read()
		for unicode.IsDigit(s.peek()) {
			s.read()
		}
	}

	literal := string(s.fileContents[s.start:s.current])
	literalNumber, ok := parseToNumber(literal)
	if !ok {
		s.addError("Invalid number format", "")
	} else {
		s.addTokenWithLiteral(NUMBER, parseToString(literalNumber))
	}
}

func (s *Scanner) scanIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.read()
	}

	tokenType := IDENTIFIER
	keywords := s.getKeywords()
	if keyword, exists := keywords[string(s.fileContents[s.start:s.current])]; exists {
		tokenType = keyword
	}
	s.addToken(tokenType)
}

func (s *Scanner) getSingleTokens() map[rune]string {
	return map[rune]string{
		'(': LEFT_PAREN,
		')': RIGHT_PAREN,
		'{': LEFT_BRACE,
		'}': RIGHT_BRACE,
		',': COMMA,
		'.': DOT,
		'-': MINUS,
		'+': PLUS,
		';': SEMICOLON,
		'*': STAR,
	}
}

func (s *Scanner) getComparisonTokens() map[rune]string {
	return map[rune]string{
		'!': BANG,
		'=': EQUAL,
		'<': LESS,
		'>': GREATER,
	}
}

func (s *Scanner) getKeywords() map[string]string {
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

func (s *Scanner) scanToken() {
	chr := s.read()

	singleTokens := s.getSingleTokens()

	if tokenType, exists := singleTokens[chr]; exists {
		s.addToken(tokenType)
		return
	}

	comparisonTokens := s.getComparisonTokens()

	if tokenType, exists := comparisonTokens[chr]; exists {
		s.addToken(stringConditional(s.match('='), tokenType+"_EQUAL", tokenType))
		return
	}

	switch chr {
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.read()
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		s.line++
	case '"':
		s.scanString()
	default:
		if unicode.IsDigit(chr) {
			s.scanNumber()
		} else if isAlpha(chr) {
			s.scanIdentifier()
		} else {
			s.addError("Unexpected character", string(chr))
		}
	}
}

func (s *Scanner) scanTokens() ([]Token, bool) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.start = s.current
	s.addToken(EOF)

	if len(s.errors) > 0 {
		return s.tokens, false
	}
	return s.tokens, true
}

func (s *Scanner) printTokens() {
	for _, err := range s.errors {
		err.print()
	}
	for _, token := range s.tokens {
		fmt.Print(token.toString())
	}
}
