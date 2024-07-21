package main

import (
	"errors"
	"fmt"
	"os"
)

// import (
// 	"errors"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"unicode"
// )

type ParseError struct {
	Token   Token
	message string
}

type Parser struct {
	Tokens []Token
	Cursor int
	Exprs  []Expr
	errors []ParseError
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
	print() string
}

type ExprLiteral struct {
	Value string
}

func (expr *ExprLiteral) print() string {
	return expr.Value
}

type ExprGroup struct {
	Expr Expr
}

func (expr *ExprGroup) print() string {
	return fmt.Sprintf("(group %s)", expr.Expr.print())
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

func (p *Parser) matchSome(types []string) bool {
	for _, Type := range types {
		if p.match(Type) {
			return true
		}
	}
	return false
}

func (p *Parser) error(token Token, message string) {
	p.errors = append(p.errors, ParseError{Token: token, message: message})
}

func (p *Parser) consume(Type string, message string) error {
	if p.check(Type) {
		p.read()
		return nil
	}
	p.error(p.previous(), message)
	return errors.New("")
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
	if p.matchSome([]string{NUMBER, STRING}) {
		return &ExprLiteral{
			Value: p.previous().Literal,
		}
	}
	if p.match(LEFT_PAREN) {
		expr := p.expression()
		err := p.consume(RIGHT_PAREN, "Unmatched parentheses.")
		if err != nil {
			return nil
		}
		return &ExprGroup{
			Expr: expr,
		}
	}
	return nil
}

func (p *Parser) unary() Expr {
	return p.primary()
}

func (p *Parser) factor() Expr {
	return p.unary()
}

func (p *Parser) term() Expr {
	return p.factor()
}

func (p *Parser) comparison() Expr {
	return p.term()
}

func (p *Parser) equality() Expr {
	return p.comparison()
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) parse() error {
	for !p.isAtEnd() {
		expr := p.expression()
		if expr != nil {
			p.Exprs = append(p.Exprs, expr)
		}
	}

	if len(p.errors) > 0 {
		return errors.New("")
	}
	return nil
}

func (err *ParseError) print() {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err.message)
}

func (p *Parser) print() {
	for _, err := range p.errors {
		err.print()
	}

	for _, expr := range p.Exprs {
		fmt.Printf("%s\n", expr.print())
	}
}
