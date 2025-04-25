package interpreter

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Interpreter struct {
	prog           string
	progPtr        int
	bracketIndices map[int]int

	tape    []byte
	tapePtr int

	Reader io.Reader
	Writer io.Writer
}

func New(program string) (*Interpreter, error) {
	i := &Interpreter{
		prog:           program,
		progPtr:        0,
		bracketIndices: map[int]int{},

		tape:    []byte{0},
		tapePtr: 0,

		Reader: os.Stdin,
		Writer: os.Stdout,
	}

	if err := i.computeBracketIndices(); err != nil {
		return nil, err
	}

	return i, nil
}

func (i *Interpreter) Tape() []byte { return i.tape }

func (i *Interpreter) computeBracketIndices() error {
	var stack []int

	for idx, c := range i.prog {
		if c == '[' {
			stack = append(stack, idx)
		} else if c == ']' {
			if len(stack) == 0 {
				return fmt.Errorf("] with no matching [ at index %d", idx)
			}
			openIdx := stack[len(stack)-1]
			i.bracketIndices[idx] = openIdx
			i.bracketIndices[openIdx] = idx
			stack = stack[:len(stack)-1]
		}
	}

	if len(stack) != 0 {
		return fmt.Errorf("%d [ with no matching ]", len(stack))
	}

	return nil
}

func (i *Interpreter) Run() error {
	bufWriter := bufio.NewWriter(i.Writer)

	for i.progPtr < len(i.prog) {
		c := i.prog[i.progPtr]

		switch c {
		case '+':
			i.tape[i.tapePtr]++
		case '-':
			i.tape[i.tapePtr]--
		case '>':
			i.tapePtr++
			if i.tapePtr >= len(i.tape) {
				i.tape = append(i.tape, 0)
			}
		case '<':
			if i.tapePtr == 0 {
				return fmt.Errorf("tape pointer underflow at index %d", i.progPtr)
			}
			i.tapePtr--
		case '[':
			if i.tape[i.tapePtr] == 0 {
				i.progPtr = i.bracketIndices[i.progPtr]
				continue
			}
		case ']':
			if i.tape[i.tapePtr] != 0 {
				i.progPtr = i.bracketIndices[i.progPtr]
				continue
			}
		case ',':
			b := make([]byte, 1)
			if _, err := i.Reader.Read(b); err != nil {
				return err
			}
			i.tape[i.tapePtr] = b[0]
		case '.':
			b := i.tape[i.tapePtr]
			s := []byte{b}
			if _, err := bufWriter.Write(s); err != nil {
				return err
			}
			if b == '\n' {
				if err := bufWriter.Flush(); err != nil {
					return err
				}
			}
		}

		i.progPtr += 1
	}

	if err := bufWriter.Flush(); err != nil {
		return err
	}

	return nil
}
