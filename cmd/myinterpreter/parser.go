package main

import "fmt"

// import (
// 	"errors"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"unicode"
// )

type Parser struct {
	Tokens []Token
	Cursor int
	Exprs  []Expr
}

func NewParser() *Parser {
	return &Parser{}
}

const (
	EXPRESSION = "EXPRESSION"
	EQUALITY   = "EQUALITY"
	COMPARISON = "COMPARISON"
	TERM       = "TERM"
	FACTOR     = "FACTOR"
	UNARY      = "UNARY"
	PRIMARY    = "PRIMARY"
)

type Expr interface {
	print()
}

type ExprLiteral struct {
	Value string
}

func (expr *ExprLiteral) print() {
	fmt.Printf("%s\n", expr.Value)
}

type ExprGroup struct {
	Expr Expr
}

type ExprUnary struct {
	Token Token
	Right Expr
}

type ExprBinary struct {
	Left  Expr
	Token Token
	Right Expr
}

func (p *Parser) isAtEnd() bool {
	return p.Cursor >= len(p.Tokens)
}

func (p *Parser) read() Token {
	token := p.Tokens[p.Cursor]
	p.Cursor++
	return token
}

func (p *Parser) peek() Token {
	return p.Tokens[p.Cursor]
}

func (p *Parser) previous() Token {
	return p.Tokens[p.Cursor-1]
}

func (p *Parser) check(Type string) bool {
	if p.isAtEnd() {
		return false
	}
	token := p.peek()
	return token.Type == Type
}

func (p *Parser) match(Type string) bool {
	if p.check(Type) {
		p.read()
		return true
	}
	return false
}

func (p *Parser) matchMany(types []string) bool {
	for _, Type := range types {
		if p.match(Type) {
			return true
		}
	}
	return false
}

func (p *Parser) primary() Expr {
	if p.match(TRUE) {
		return &ExprLiteral{
			Value: "true",
		}
	}
	if p.match(FALSE) {
		return &ExprLiteral{
			Value: "false",
		}
	}
	if p.match(NIL) {
		return &ExprLiteral{
			Value: "nil",
		}
	}
	return nil
}

func (p *Parser) parse() {
	for !p.isAtEnd() {
		p.Exprs = append(p.Exprs, p.primary())
	}
}

func (p *Parser) print() {
	for _, expr := range p.Exprs {
		expr.print()
	}
}
