package vm

import (
	"bufio"
	"fmt"
	"io"

	"github.com/interrrp/trauma/internal/bytecode"
)

func Run(code bytecode.Bytecode, reader io.Reader, writer io.Writer) (*Result, error) {
	return newVm(code, reader, writer).run()
}

const tapePreallocSize = 5_000

type vm struct {
	code         bytecode.Bytecode
	codePtr      int
	loopIndexMap []int

	tape    []byte
	tapePtr int

	reader io.Reader
	writer io.Writer
}

func newVm(code bytecode.Bytecode, reader io.Reader, writer io.Writer) *vm {
	return &vm{
		code:    code,
		codePtr: 0,

		tape:    make([]byte, tapePreallocSize),
		tapePtr: 0,

		reader: reader,
		writer: writer,
	}
}

func (v *vm) buildLoopIndexMap() error {
	var stack []int

	v.loopIndexMap = make([]int, len(v.code))
	for i := range v.loopIndexMap {
		v.loopIndexMap[i] = -1
	}

	for currentIndex, inst := range v.code {
		switch inst.Kind() {
		case bytecode.LoopStart:
			stack = append(stack, currentIndex)

		case bytecode.LoopEnd:
			if len(stack) == 0 {
				return fmt.Errorf("unmatched ending bracket at instruction %d", currentIndex)
			}

			startIndex := stack[len(stack)-1]
			closeIndex := currentIndex
			stack = stack[:len(stack)-1]

			v.loopIndexMap[startIndex] = closeIndex
			v.loopIndexMap[closeIndex] = startIndex
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

	for v.codePtr < len(v.code) {
		currentCell := &v.tape[v.tapePtr]

		inst := v.code[v.codePtr]
		switch inst.Kind() {
		case bytecode.Inc:
			*currentCell = byte(int(*currentCell) + inst.Amount())

		case bytecode.Clear:
			*currentCell = 0

		case bytecode.Move:
			v.tape[v.tapePtr+inst.Amount()] += *currentCell
			*currentCell = 0

		case bytecode.IncPtr:
			v.tapePtr += inst.Amount()

			for v.tapePtr >= len(v.tape) {
				v.tape = append(v.tape, 0)
			}

			if v.tapePtr < 0 {
				return nil, fmt.Errorf("tape pointer underflow at instruction %d", v.codePtr)
			}

		case bytecode.LoopStart:
			if *currentCell == 0 {
				v.codePtr = v.loopIndexMap[v.codePtr]
			}

		case bytecode.LoopEnd:
			if *currentCell != 0 {
				v.codePtr = v.loopIndexMap[v.codePtr]
			}

		case bytecode.Input:
			b := make([]byte, 1)
			if _, err := v.reader.Read(b); err != nil {
				return nil, err
			}
			*currentCell = b[0]

		case bytecode.Output:
			s := []byte{*currentCell}
			if _, err := bufWriter.Write(s); err != nil {
				return nil, err
			}
			if *currentCell == '\n' {
				if err := bufWriter.Flush(); err != nil {
					return nil, err
				}
			}
		}

		v.codePtr++
	}

	if err := bufWriter.Flush(); err != nil {
		return nil, err
	}

	return newResult(v), nil
}
