package main

import (
	"fmt"
	"os"
)

func debug(message string) {
	fmt.Fprintln(os.Stderr, message)
}
