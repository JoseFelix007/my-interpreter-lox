package main

import (
	"os"
)

type Interpreter struct {
	fileContents []byte
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

func (lox *Interpreter) tokenize() error {
	lox.scanner = NewScanner(lox.fileContents)
	return lox.scanner.scanTokens()
}

func (lox *Interpreter) parse() error {
	var err error
	if lox.scanner == nil {
		err = lox.tokenize()
		lox.parser = NewParser()
		if err != nil {
			return err
		}
	}
	lox.parser.Tokens = lox.scanner.tokens
	err = lox.parser.parse()
	return err
}
