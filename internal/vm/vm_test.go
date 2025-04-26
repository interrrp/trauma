package vm

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/interrrp/trauma/internal/bytecode"
)

func TestIncrement(t *testing.T) {
	r := mustRun(t, "+++")
	assertCell(t, r, 0, 3)

	r = mustRun(t, "+++++")
	assertCell(t, r, 0, 5)
}

func TestDecrement(t *testing.T) {
	r := mustRun(t, "++-")
	assertCell(t, r, 0, 1)

	r = mustRun(t, "---")
	assertCell(t, r, 0, 255-2)
}

func TestIncPtr(t *testing.T) {
	r := mustRun(t, ">>><")
	if r.TapePtr() != 2 {
		t.Errorf("expected tape pointer to be 2, got %d", r.TapePtr())
	}

	r = mustRun(t, "+>++>><+++")
	assertCell(t, r, 0, 1)
	assertCell(t, r, 1, 2)
	assertCell(t, r, 2, 3)

	if _, err := runProg("<"); err == nil {
		t.Error("expected error on tape pointer underflow")
	}
}

func TestMove(t *testing.T) {
	r := mustRun(t, "+++[->+<]")
	assertCell(t, r, 0, 0)
	assertCell(t, r, 1, 3)

	r = mustRun(t, "+>+++[-<+>]")
	assertCell(t, r, 0, 4)
	assertCell(t, r, 1, 0)
}

func TestLoop(t *testing.T) {
	r := mustRun(t, "+++[-]")
	assertCell(t, r, 0, 0)

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
	r, err := runProgWithCustomIO(",", strings.NewReader("A"), nullWriter)
	if err != nil {
		t.Errorf("error during execution: %v", err)
	}
	assertCell(t, r, 0, 'A')
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

func mustRun(t *testing.T, program string) *Result {
	r, err := runProg(program)
	if err != nil {
		t.Errorf("error during execution: %v", err)
	}
	return r
}

func runProg(program string) (*Result, error) {
	bc, err := bytecode.Compile(program)
	if err != nil {
		return nil, err
	}
	return Run(bc, nullReader, nullWriter)
}

func runProgWithCustomIO(program string, reader io.Reader, writer io.Writer) (*Result, error) {
	bc, err := bytecode.Compile(program)
	if err != nil {
		return nil, err
	}
	return Run(bc, reader, writer)
}

func assertCell(t *testing.T, r *Result, idx int, val byte) {
	actual := r.Tape()[idx]
	if actual != val {
		t.Errorf("expected cell %d to be %d, got %d", idx, val, actual)
	}
}
