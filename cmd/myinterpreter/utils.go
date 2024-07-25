package main

import (
	"fmt"
	"os"
	"strconv"
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

func parseToNumber(value string) (float64, bool) {
	valueFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}
	return valueFloat, true
}

func parseToString(value float64) string {
	prec := -1
	if value == float64(int(value)) {
		prec = 1
	}
	return strconv.FormatFloat(value, 'f', prec, 64)
}

func parseBoolToString(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func isTruthy(value string) string {
	if value == "nil" || value == "false" {
		return "false"
	}
	return "true"
}

func negate(value string) string {
	if value == "false" {
		return "true"
	}
	return "false"
}
