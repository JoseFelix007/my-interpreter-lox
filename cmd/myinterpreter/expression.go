package main

import (
	"fmt"
)

type Expr interface {
	toString() string
	evaluate() string
}

type ExprLiteral struct {
	Value string
}

func (expr *ExprLiteral) toString() string {
	return expr.Value
}

func (expr *ExprLiteral) evaluate() string {
	return expr.Value
}

type ExprGroup struct {
	Expr Expr
}

func (expr *ExprGroup) toString() string {
	if expr.Expr == nil {
		return ""
	} else {
		return fmt.Sprintf("(group %s)", expr.Expr.toString())
	}
}

func (expr *ExprGroup) evaluate() string {
	return expr.Expr.evaluate()
}

type ExprUnary struct {
	Token *Token
	Right Expr
}

func (expr *ExprUnary) toString() string {
	if expr.Right == nil {
		return ""
	} else {
		token := ""
		if expr.Token != nil {
			token = expr.Token.Lexema
		}
		return fmt.Sprintf("(%s %s)", token, expr.Right.toString())
	}
}

func (expr *ExprUnary) evaluate() string {
	right := expr.Right.evaluate()

	switch expr.Token.Type {
	case MINUS:
		return "-" + right
	case BANG:
		return negate(isTruthy(right))
	}
	return "nil"
}

type ExprBinary struct {
	Left  Expr
	Token *Token
	Right Expr
}

func (expr *ExprBinary) toString() string {
	if expr.Left == nil || expr.Right == nil || expr.Token == nil {
		return ""
	} else {
		return fmt.Sprintf("(%s %s %s)", expr.Token.Lexema, expr.Left.toString(), expr.Right.toString())
	}
}

func (expr *ExprBinary) evaluate() string {
	left := expr.Left.evaluate()
	right := expr.Right.evaluate()

	leftNumber, leftOk := parseToNumber(left)
	rightNumber, rightOk := parseToNumber(right)
	switch expr.Token.Type {
	case GREATER:
		if leftOk && rightOk {
			return parseBoolToString(leftNumber > rightNumber)
		}
	case GREATER_EQUAL:
		if leftOk && rightOk {
			return parseBoolToString(leftNumber >= rightNumber)
		}
	case LESS:
		if leftOk && rightOk {
			return parseBoolToString(leftNumber < rightNumber)
		}
	case LESS_EQUAL:
		if leftOk && rightOk {
			return parseBoolToString(leftNumber <= rightNumber)
		}
	case BANG_EQUAL:
		return parseBoolToString(left != right)
	case EQUAL_EQUAL:
		return parseBoolToString(left == right)
	case MINUS:
		if leftOk && rightOk {
			return parseToString(leftNumber - rightNumber)
		}
	case PLUS:
		if leftOk && rightOk {
			return parseToString(leftNumber + rightNumber)
		}
		return left + right
	case SLASH:
		if leftOk && rightOk {
			return parseToString(leftNumber / rightNumber)
		}
	case STAR:
		if leftOk && rightOk {
			return parseToString(leftNumber * rightNumber)
		}
	}
	return "nil"
}
