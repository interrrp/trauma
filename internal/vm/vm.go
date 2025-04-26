package vm

import (
	"bufio"
	"fmt"
	"io"

	"github.com/interrrp/trauma/internal/bytecode"
)

func Run(bc bytecode.Bytecode, reader io.Reader, writer io.Writer) (*Result, error) {
	v := newVm(bc, reader, writer)
	return v.run()
}

type vm struct {
	bytecode  bytecode.Bytecode
	bcPtr     int
	loopStack []int

	tape    []byte
	tapePtr int

	reader io.Reader
	writer io.Writer
}

func newVm(bc bytecode.Bytecode, reader io.Reader, writer io.Writer) *vm {
	return &vm{
		bytecode:  bc,
		bcPtr:     0,
		loopStack: []int{},

		tape:    []byte{0},
		tapePtr: 0,

		reader: reader,
		writer: writer,
	}
}

func (v *vm) run() (*Result, error) {
	bufWriter := bufio.NewWriter(v.writer)

	for v.bcPtr < len(v.bytecode) {
		instr := v.bytecode[v.bcPtr]

		switch inst := instr.(type) {
		case *bytecode.CellInc:
			v.tape[v.tapePtr] = byte(int(v.tape[v.tapePtr]) + inst.Amount())

		case *bytecode.CellEmpty:
			v.tape[v.tapePtr] = 0

		case *bytecode.PtrInc:
			v.tapePtr += inst.Amount()

			for v.tapePtr >= len(v.tape) {
				v.tape = append(v.tape, 0)
			}

			if v.tapePtr < 0 {
				return nil, fmt.Errorf("tape pointer underflow at instruction %d", v.bcPtr)
			}

		case *bytecode.LoopStart:
			if v.tape[v.tapePtr] == 0 {
				depth := 1
				pos := v.bcPtr

				for depth > 0 && pos < len(v.bytecode)-1 {
					pos++
					switch v.bytecode[pos].(type) {
					case *bytecode.LoopStart:
						depth++
					case *bytecode.LoopEnd:
						depth--
					}
				}

				if depth != 0 {
					return nil, fmt.Errorf("unmatched loop start at instruction %d", v.bcPtr)
				}

				v.bcPtr = pos
			} else {
				v.loopStack = append(v.loopStack, v.bcPtr)
			}

		case *bytecode.LoopEnd:
			if v.tape[v.tapePtr] != 0 {
				if len(v.loopStack) == 0 {
					return nil, fmt.Errorf("unmatched loop end at instruction %d", v.bcPtr)
				}

				start := v.loopStack[len(v.loopStack)-1]
				v.bcPtr = start
			} else {
				if len(v.loopStack) > 0 {
					v.loopStack = v.loopStack[:len(v.loopStack)-1]
				}
			}

		case *bytecode.Input:
			b := make([]byte, 1)
			if _, err := v.reader.Read(b); err != nil {
				return nil, err
			}
			v.tape[v.tapePtr] = b[0]

		case *bytecode.Output:
			b := v.tape[v.tapePtr]
			s := []byte{b}
			if _, err := bufWriter.Write(s); err != nil {
				return nil, err
			}
			if b == '\n' {
				if err := bufWriter.Flush(); err != nil {
					return nil, err
				}
			}
		}

		v.bcPtr++
	}

	if err := bufWriter.Flush(); err != nil {
		return nil, err
	}

	return newResult(v), nil
}
