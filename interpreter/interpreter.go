package interpreter

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Run runs a Brainfuck program.
//
// Run will use os.Stdin and os.Stdout for the input (,) and output (.) commands,
// equivalent to RunWithCustomIO(program, os.Stdin, os.Stdout).
func Run(program string) (*Result, error) {
	i := interpreter{
		program:         program,
		programPtr:      0,
		bracketIndexMap: map[int]int{},

		tape:    []byte{0},
		tapePtr: 0,

		reader: os.Stdin,
		writer: os.Stdout,
	}

	if err := i.buildBracketIndexMap(); err != nil {
		return nil, err
	}

	if err := i.run(); err != nil {
		return nil, err
	}

	return newResult(&i), nil
}

// RunWithCustomIO runs a Brainfuck program using the given readers and writers
// for the input (,) and output (.) commands.
func RunWithCustomIO(program string, reader io.Reader, writer io.Writer) (*Result, error) {
	i := interpreter{
		program:         program,
		programPtr:      0,
		bracketIndexMap: map[int]int{},

		tape:    []byte{0},
		tapePtr: 0,

		reader: reader,
		writer: writer,
	}

	if err := i.buildBracketIndexMap(); err != nil {
		return nil, err
	}

	if err := i.run(); err != nil {
		return nil, err
	}

	return newResult(&i), nil
}

type Result struct {
	tape    []byte
	tapePtr int
}

func newResult(i *interpreter) *Result {
	return &Result{tape: i.tape, tapePtr: i.tapePtr}
}
func (r *Result) Tape() []byte { return r.tape }
func (r *Result) TapePtr() int { return r.tapePtr }

type interpreter struct {
	program         string
	programPtr      int
	bracketIndexMap map[int]int

	tape    []byte
	tapePtr int

	reader io.Reader
	writer io.Writer
}

func (i *interpreter) buildBracketIndexMap() error {
	var stack []int

	for idx, c := range i.program {
		if c == '[' {
			stack = append(stack, idx)
		} else if c == ']' {
			if len(stack) == 0 {
				return fmt.Errorf("] with no matching [ at index %d", idx)
			}
			openIdx := stack[len(stack)-1]
			i.bracketIndexMap[idx] = openIdx
			i.bracketIndexMap[openIdx] = idx
			stack = stack[:len(stack)-1]
		}
	}

	if len(stack) != 0 {
		return fmt.Errorf("%d [ with no matching ]", len(stack))
	}

	return nil
}

func (i *interpreter) run() error {
	bufWriter := bufio.NewWriter(i.writer)

	for i.programPtr < len(i.program) {
		switch i.program[i.programPtr] {
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
				return fmt.Errorf("tape pointer underflow at index %d", i.programPtr)
			}
			i.tapePtr--

		case '[':
			if i.tape[i.tapePtr] == 0 {
				i.programPtr = i.bracketIndexMap[i.programPtr]
				continue
			}
		case ']':
			if i.tape[i.tapePtr] != 0 {
				i.programPtr = i.bracketIndexMap[i.programPtr]
				continue
			}

		case ',':
			b := make([]byte, 1)
			if _, err := i.reader.Read(b); err != nil {
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

		i.programPtr += 1
	}

	if err := bufWriter.Flush(); err != nil {
		return err
	}

	return nil
}
