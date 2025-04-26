package vm

import (
	"bufio"
	"fmt"
	"io"

	"github.com/interrrp/trauma/internal/bytecode"
)

func Run(bc bytecode.Bytecode, reader io.Reader, writer io.Writer) (*Result, error) {
	return newVm(bc, reader, writer).run()
}

const tapePreallocSize = 5_000

type vm struct {
	bytecode     bytecode.Bytecode
	bytecodePtr  int
	loopIndexMap []int

	tape    []byte
	tapePtr int

	reader io.Reader
	writer io.Writer
}

func newVm(bc bytecode.Bytecode, reader io.Reader, writer io.Writer) *vm {
	return &vm{
		bytecode:    bc,
		bytecodePtr: 0,

		tape:    make([]byte, tapePreallocSize),
		tapePtr: 0,

		reader: reader,
		writer: writer,
	}
}

func (v *vm) buildLoopIndexMap() error {
	var stack []int

	v.loopIndexMap = make([]int, len(v.bytecode))
	for i := range v.loopIndexMap {
		v.loopIndexMap[i] = -1
	}

	for i, inst := range v.bytecode {
		switch inst.(type) {
		case *bytecode.LoopStart:
			stack = append(stack, i)

		case *bytecode.LoopEnd:
			if len(stack) == 0 {
				return fmt.Errorf("unmatched ending bracket at instruction %d", i)
			}

			startIdx := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			v.loopIndexMap[startIdx] = i
			v.loopIndexMap[i] = startIdx
		}
	}

	if len(stack) != 0 {
		return fmt.Errorf("unmatched starting bracket(s) at instruction(s) %v", stack)
	}

	return nil
}

func (v *vm) run() (*Result, error) {
	if err := v.buildLoopIndexMap(); err != nil {
		return nil, err
	}

	bufWriter := bufio.NewWriter(v.writer)

	for v.bytecodePtr < len(v.bytecode) {
		i := v.bytecode[v.bytecodePtr]
		switch inst := i.(type) {
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
				return nil, fmt.Errorf("tape pointer underflow at instruction %d", v.bytecodePtr)
			}

		case *bytecode.LoopStart:
			if v.tape[v.tapePtr] == 0 {
				v.bytecodePtr = v.loopIndexMap[v.bytecodePtr]
			}

		case *bytecode.LoopEnd:
			if v.tape[v.tapePtr] != 0 {
				v.bytecodePtr = v.loopIndexMap[v.bytecodePtr]
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

		v.bytecodePtr++
	}

	if err := bufWriter.Flush(); err != nil {
		return nil, err
	}

	return newResult(v), nil
}
