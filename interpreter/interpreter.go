package interpreter

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/interrrp/trauma/bytecode"
)

// Run runs a Brainfuck program.
//
// Run will use os.Stdin and os.Stdout for the input (,) and output (.) commands,
// equivalent to RunWithCustomIO(program, os.Stdin, os.Stdout).
func Run(program string) (*Result, error) {
	return RunWithCustomIO(program, os.Stdin, os.Stdout)
}

// RunWithCustomIO runs a Brainfuck program using the given readers and writers
// for the input (,) and output (.) commands.
func RunWithCustomIO(program string, reader io.Reader, writer io.Writer) (*Result, error) {
	bc, err := bytecode.Compile(program)
	if err != nil {
		return nil, fmt.Errorf("compilation error: %v", err)
	}

	i := newBytecodeInterpreter(bc, reader, writer)
	if err := i.run(); err != nil {
		return nil, err
	}

	return newResult(i), nil
}

type bytecodeInterpreter struct {
	bytecode bytecode.Bytecode
	bcPtr    int

	tape    []byte
	tapePtr int

	loopStack []int

	reader io.Reader
	writer io.Writer
}

func newBytecodeInterpreter(bc bytecode.Bytecode, reader io.Reader, writer io.Writer) *bytecodeInterpreter {
	return &bytecodeInterpreter{
		bytecode:  bc,
		bcPtr:     0,
		tape:      []byte{0},
		tapePtr:   0,
		loopStack: []int{},
		reader:    reader,
		writer:    writer,
	}
}

type Result struct {
	tape    []byte
	tapePtr int
}

func newResult(i *bytecodeInterpreter) *Result {
	return &Result{tape: i.tape, tapePtr: i.tapePtr}
}
func (r *Result) Tape() []byte { return r.tape }
func (r *Result) TapePtr() int { return r.tapePtr }

func (i *bytecodeInterpreter) run() error {
	bufWriter := bufio.NewWriter(i.writer)

	for i.bcPtr < len(i.bytecode) {
		instr := i.bytecode[i.bcPtr]

		switch inst := instr.(type) {
		case *bytecode.CellInc:
			i.tape[i.tapePtr] = byte(int(i.tape[i.tapePtr]) + inst.Amount())

		case *bytecode.PtrInc:
			i.tapePtr += inst.Amount()

			for i.tapePtr >= len(i.tape) {
				i.tape = append(i.tape, 0)
			}

			if i.tapePtr < 0 {
				return fmt.Errorf("tape pointer underflow at instruction %d", i.bcPtr)
			}

		case *bytecode.LoopStart:
			if i.tape[i.tapePtr] == 0 {
				depth := 1
				pos := i.bcPtr

				for depth > 0 && pos < len(i.bytecode)-1 {
					pos++
					switch i.bytecode[pos].(type) {
					case *bytecode.LoopStart:
						depth++
					case *bytecode.LoopEnd:
						depth--
					}
				}

				if depth != 0 {
					return fmt.Errorf("unmatched loop start at instruction %d", i.bcPtr)
				}

				i.bcPtr = pos
			} else {
				i.loopStack = append(i.loopStack, i.bcPtr)
			}

		case *bytecode.LoopEnd:
			if i.tape[i.tapePtr] != 0 {
				if len(i.loopStack) == 0 {
					return fmt.Errorf("unmatched loop end at instruction %d", i.bcPtr)
				}

				start := i.loopStack[len(i.loopStack)-1]
				i.bcPtr = start
			} else {
				if len(i.loopStack) > 0 {
					i.loopStack = i.loopStack[:len(i.loopStack)-1]
				}
			}

		case *bytecode.Input:
			b := make([]byte, 1)
			if _, err := i.reader.Read(b); err != nil {
				return err
			}
			i.tape[i.tapePtr] = b[0]

		case *bytecode.Output:
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

		i.bcPtr++
	}

	if err := bufWriter.Flush(); err != nil {
		return err
	}

	return nil
}
