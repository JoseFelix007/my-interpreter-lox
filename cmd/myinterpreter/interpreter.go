package main

import (
	"fmt"
	"os"
)

type Interpreter struct {
	fileContents []byte
	tokens       []Token
	exprs        []Expr
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
	exprs, ok := lox.parser.parse()
	lox.exprs = exprs
	return ok
}

func (lox *Interpreter) evaluate() bool {
	if lox.parser == nil {
		if ok := lox.parse(); !ok {
			return false
		}
	}
	for _, expr := range lox.exprs {
		result := expr.evaluate()
		fmt.Printf("%s\n", result)
	}
	return true
}
