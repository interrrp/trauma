package main

import (
	"fmt"
	"os"

	"github.com/interrrp/trauma/interpreter"
)

func main() {
	if len(os.Args) != 2 {
		fail("usage: %s <path to brainfuck file>", os.Args[0])
	}
	path := os.Args[1]

	b, err := os.ReadFile(path)
	if err != nil {
		fail("failed to read file: %v", err)
	}
	prog := string(b)

	if _, err := interpreter.Run(prog); err != nil {
		fail("%v", err)
	}
}

func fail(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
