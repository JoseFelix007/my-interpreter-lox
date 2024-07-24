package main

import "fmt"

type Token struct {
	Lexema  string
	Literal string
	Line    int
	Type    string
}

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

func (token *Token) toString() string {
	if token.Literal == "" {
		token.Literal = "null"
	}
	return fmt.Sprintf("%s %s %s\n", token.Type, token.Lexema, token.Literal)
}
