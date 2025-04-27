package vm

import (
	"bufio"
	"fmt"
	"io"

	"github.com/interrrp/trauma/internal/bytecode"
)

func Run(code bytecode.Bytecode, reader io.Reader, writer io.Writer) (*VM, error) {
	vm := newVM(code, reader, writer)
	if err := vm.run(); err != nil {
		return vm, err
	}
	return vm, nil
}

const tapePreallocSize = 5_000

type VM struct {
	code         bytecode.Bytecode
	codePtr      int
	loopIndexMap []int

	tape    []byte
	tapePtr int

	reader io.Reader
	writer io.Writer
}

func newVM(code bytecode.Bytecode, reader io.Reader, writer io.Writer) *VM {
	return &VM{
		code:    code,
		codePtr: 0,

		tape:    make([]byte, tapePreallocSize),
		tapePtr: 0,

		reader: reader,
		writer: writer,
	}
}

func (vm *VM) Tape() []byte { return vm.tape }
func (vm *VM) TapePtr() int { return vm.tapePtr }

func (vm *VM) buildLoopIndexMap() error {
	var stack []int

	vm.loopIndexMap = make([]int, len(vm.code))
	for i := range vm.loopIndexMap {
		vm.loopIndexMap[i] = -1
	}

	for currentIndex, inst := range vm.code {
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

			vm.loopIndexMap[startIndex] = closeIndex
			vm.loopIndexMap[closeIndex] = startIndex
		}
	}

	if len(stack) != 0 {
		return fmt.Errorf("unmatched starting bracket(s) at instruction(s) %v", stack)
	}

	return nil
}

func (vm *VM) run() error {
	if err := vm.buildLoopIndexMap(); err != nil {
		return err
	}

	bufWriter := bufio.NewWriter(vm.writer)

	for vm.codePtr < len(vm.code) {
		currentCell := &vm.tape[vm.tapePtr]
		inst := vm.code[vm.codePtr]

		switch inst.Kind() {
		case bytecode.Inc:
			*currentCell = byte(int(*currentCell) + inst.Amount())

		case bytecode.Clear:
			*currentCell = 0

		case bytecode.Move:
			vm.tape[vm.tapePtr+inst.Amount()] += *currentCell
			*currentCell = 0

		case bytecode.IncPtr:
			vm.tapePtr += inst.Amount()

			for vm.tapePtr >= len(vm.tape) {
				vm.tape = append(vm.tape, 0)
			}

			if vm.tapePtr < 0 {
				return fmt.Errorf("tape pointer underflow at instruction %d", vm.codePtr)
			}

		case bytecode.LoopStart:
			if *currentCell == 0 {
				vm.codePtr = vm.loopIndexMap[vm.codePtr]
			}

		case bytecode.LoopEnd:
			if *currentCell != 0 {
				vm.codePtr = vm.loopIndexMap[vm.codePtr]
			}

		case bytecode.Input:
			b := make([]byte, 1)
			if _, err := vm.reader.Read(b); err != nil {
				return err
			}
			*currentCell = b[0]

		case bytecode.Output:
			s := []byte{*currentCell}
			if _, err := bufWriter.Write(s); err != nil {
				return err
			}
			if *currentCell == '\n' {
				if err := bufWriter.Flush(); err != nil {
					return err
				}
			}
		}

		vm.codePtr++
	}

	if err := bufWriter.Flush(); err != nil {
		return err
	}

	return nil
}
