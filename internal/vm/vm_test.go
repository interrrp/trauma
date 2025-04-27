package vm_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/interrrp/trauma/internal/bytecode"
	"github.com/interrrp/trauma/internal/vm"
)

func TestIncrement(t *testing.T) {
	vm := mustRun(t, "+++")
	assertCell(t, vm, 0, 3)

	vm = mustRun(t, "+++++")
	assertCell(t, vm, 0, 5)
}

func TestDecrement(t *testing.T) {
	vm := mustRun(t, "++-")
	assertCell(t, vm, 0, 1)

	vm = mustRun(t, "---")
	assertCell(t, vm, 0, 255-2)
}

func TestIncPtr(t *testing.T) {
	vm := mustRun(t, ">>><")
	if vm.TapePtr() != 2 {
		t.Errorf("expected tape pointer to be 2, got %d", vm.TapePtr())
	}

	vm = mustRun(t, "+>++>><+++")
	assertCell(t, vm, 0, 1)
	assertCell(t, vm, 1, 2)
	assertCell(t, vm, 2, 3)

	if _, err := runProg("<"); err == nil {
		t.Error("expected error on tape pointer underflow")
	}
}

func TestLoop(t *testing.T) {
	vm := mustRun(t, "+++[-]")
	assertCell(t, vm, 0, 0)

	if _, err := runProg("[+[+[+]+]+]"); err != nil {
		t.Errorf("error on valid syntax: %v", err)
	}

	if _, err := runProg("[[["); err == nil {
		t.Error("expected error on invalid syntax")
	}

	if _, err := runProg("[]]"); err == nil {
		t.Error("expected error on invalid syntax")
	}
}

var (
	nullReader = strings.NewReader("")
	nullWriter = io.Discard
)

func TestInput(t *testing.T) {
	vm, err := runProgWithCustomIO(",", strings.NewReader("A"), nullWriter)
	if err != nil {
		t.Errorf("error during execution: %v", err)
	}
	assertCell(t, vm, 0, 'A')
}

func TestOutput(t *testing.T) {
	var writer bytes.Buffer
	_, err := runProgWithCustomIO("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++.", nullReader, &writer)
	if err != nil {
		t.Errorf("error during execution: %v", err)
	}

	if writer.String() != "A" {
		t.Errorf("expected output to be %q, got %q", "A", writer.String())
	}
}

func mustRun(t *testing.T, program string) *vm.VM {
	r, err := runProg(program)
	if err != nil {
		t.Errorf("error during execution: %v", err)
	}
	return r
}

func runProg(program string) (*vm.VM, error) {
	bc, err := bytecode.Compile(program)
	if err != nil {
		return nil, err
	}
	return vm.Run(bc, nullReader, nullWriter)
}

func runProgWithCustomIO(program string, reader io.Reader, writer io.Writer) (*vm.VM, error) {
	bc, err := bytecode.Compile(program)
	if err != nil {
		return nil, err
	}
	return vm.Run(bc, reader, writer)
}

func assertCell(t *testing.T, vm *vm.VM, idx int, val byte) {
	actual := vm.Tape()[idx]
	if actual != val {
		t.Errorf("expected cell %d to be %d, got %d", idx, val, actual)
	}
}
