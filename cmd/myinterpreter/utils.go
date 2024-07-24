package main

import (
	"fmt"
	"os"
	"unicode"
)

func debug(message string) {
	fmt.Fprintln(os.Stderr, message)
}

func isAlpha(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || unicode.IsDigit(r)
}

func stringConditional(condition bool, firstValue string, secondValue string) string {
	if condition {
		return firstValue
	}
	return secondValue
}
