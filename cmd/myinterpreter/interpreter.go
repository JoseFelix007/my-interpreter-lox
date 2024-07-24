package main

import (
	"os"
)

type Interpreter struct {
	fileContents []byte
	tokens       []Token
	scanner      *Scanner
	parser       *Parser
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (lox *Interpreter) readFile(filename string) error {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	lox.fileContents = fileContents
	return nil
}

func (lox *Interpreter) tokenize() bool {
	lox.scanner = NewScanner(lox.fileContents)
	tokens, ok := lox.scanner.scanTokens()
	lox.tokens = tokens
	return ok
}

func (lox *Interpreter) parse() bool {
	if lox.scanner == nil {
		ok := lox.tokenize()
		lox.parser = NewParser()
		if !ok {
			return false
		}
	}
	lox.parser.Tokens = lox.tokens
	return lox.parser.parse()
}
