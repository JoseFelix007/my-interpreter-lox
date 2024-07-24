package main

import (
	"errors"
	"fmt"
	"os"
)

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
	if expr.Expr == nil {
		return ""
	} else {
		return fmt.Sprintf("(group %s)", expr.Expr.print())
	}
}

type ExprUnary struct {
	Token *Token
	Right Expr
}

func (expr *ExprUnary) print() string {
	if expr.Right == nil {
		return ""
	} else {
		token := ""
		if expr.Token != nil {
			token = expr.Token.Lexema
		}
		return fmt.Sprintf("(%s %s)", token, expr.Right.print())
	}
}

type ExprBinary struct {
	Left  Expr
	Token *Token
	Right Expr
}

func (expr *ExprBinary) print() string {
	if expr.Left == nil || expr.Right == nil || expr.Token == nil {
		return ""
	} else {
		return fmt.Sprintf("(%s %s %s)", expr.Token.Lexema, expr.Left.print(), expr.Right.print())
	}
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
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

func (p *Parser) primary() (Expr, bool) {
	if p.match(TRUE) {
		return &ExprLiteral{
			Value: "true",
		}, true
	}
	if p.match(FALSE) {
		return &ExprLiteral{
			Value: "false",
		}, true
	}
	if p.match(NIL) {
		return &ExprLiteral{
			Value: "nil",
		}, true
	}
	if p.matchSome([]string{NUMBER, STRING}) {
		return &ExprLiteral{
			Value: p.previous().Literal,
		}, true
	}
	if p.match(LEFT_PAREN) {
		expr, ok := p.expression()
		err := p.consume(RIGHT_PAREN, "Unmatched parentheses.")
		if err != nil {
			return nil, false
		}
		return &ExprGroup{
			Expr: expr,
		}, ok
	}
	p.error(p.read(), "Expect expression.")
	return nil, false
}

func (p *Parser) unary() (Expr, bool) {
	if p.matchSome([]string{BANG, MINUS}) {
		operator := p.previous()
		right, ok := p.unary()
		return &ExprUnary{
			Token: &operator,
			Right: right,
		}, ok
	}
	return p.primary()
}

func (p *Parser) factor() (Expr, bool) {
	left, ok := p.unary()

	for p.matchSome([]string{SLASH, STAR}) && ok {
		var right Expr
		operator := p.previous()
		right, ok = p.unary()
		left = &ExprBinary{
			Left:  left,
			Token: &operator,
			Right: right,
		}
	}

	return left, ok
}

func (p *Parser) term() (Expr, bool) {
	left, ok := p.factor()

	for p.matchSome([]string{MINUS, PLUS}) && ok {
		var right Expr
		operator := p.previous()
		right, ok = p.factor()
		left = &ExprBinary{
			Left:  left,
			Token: &operator,
			Right: right,
		}
	}

	return left, ok
}

func (p *Parser) comparison() (Expr, bool) {
	left, ok := p.term()

	for p.matchSome([]string{GREATER, GREATER_EQUAL, LESS, LESS_EQUAL}) && ok {
		var right Expr
		operator := p.previous()
		right, ok = p.term()
		left = &ExprBinary{
			Left:  left,
			Token: &operator,
			Right: right,
		}
	}

	return left, ok
}

func (p *Parser) equality() (Expr, bool) {
	left, ok := p.comparison()

	for p.matchSome([]string{BANG_EQUAL, EQUAL_EQUAL}) && ok {
		var right Expr
		operator := p.previous()
		right, ok = p.comparison()
		left = &ExprBinary{
			Left:  left,
			Token: &operator,
			Right: right,
		}
	}

	return left, ok
}

func (p *Parser) expression() (Expr, bool) {
	return p.equality()
}

func (p *Parser) parse() bool {
	for !p.isAtEnd() {
		expr, ok := p.expression()
		if ok {
			p.Exprs = append(p.Exprs, expr)
		}
	}

	return len(p.errors) <= 0
}

func (err *ParseError) print() {
	if err.Token.Type == EOF {
		fmt.Fprintf(os.Stderr, "[line %d] Error at end: %s\n", err.Token.Line, err.message)
	} else {
		fmt.Fprintf(os.Stderr, "[line %d] Error at '%s': %s\n", err.Token.Line, err.Token.Lexema, err.message)
	}
}

func (p *Parser) print() {
	for _, err := range p.errors {
		err.print()
	}

	for _, expr := range p.Exprs {
		fmt.Printf("%s\n", expr.print())
	}
}
