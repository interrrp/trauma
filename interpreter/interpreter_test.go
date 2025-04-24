package interpreter

import (
	"bytes"
	"strings"
	"testing"
)

func TestIncrement(t *testing.T) {
	i := run(t, "+++")
	expectCell(t, i, 0, 3)

	i = run(t, "+++++")
	expectCell(t, i, 0, 5)
}

func TestDecrement(t *testing.T) {
	i := run(t, "++-")
	expectCell(t, i, 0, 1)

	i = run(t, "---")
	expectCell(t, i, 0, 255-2)
}

func TestMovePtr(t *testing.T) {
	i := run(t, ">>><")
	if i.tapePtr != 2 {
		t.Errorf("expected tape pointer to be 1, got %d", i.tapePtr)
	}

	i = run(t, "+>++>><+++")
	expectCell(t, i, 0, 1)
	expectCell(t, i, 1, 2)
	expectCell(t, i, 2, 3)
	expectCell(t, i, 3, 0)

	i = mustNew(t, "<")
	if err := i.Run(); err == nil {
		t.Error("expected error for tape pointer underflow")
	}
}

func TestLoop(t *testing.T) {
	i := run(t, "+++[-]")
	expectCell(t, i, 0, 0)

	if _, err := New("[[[]]]"); err != nil {
		t.Errorf("error on valid syntax: %v", err)
	}

	if _, err := New("[]]"); err == nil {
		t.Error("expected error on invalid syntax")
	}

	if _, err := New("[[["); err == nil {
		t.Error("expected error on invalid syntax")
	}
}

func TestInput(t *testing.T) {
	i := mustNew(t, ",")
	i.Reader = strings.NewReader("A")
	mustRun(t, i)
	expectCell(t, i, 0, 'A')
}

func TestOutput(t *testing.T) {
	i := mustNew(t, "+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++.")
	var writer bytes.Buffer
	i.Writer = &writer

	mustRun(t, i)
	if writer.String() != "A" {
		t.Errorf("expected output to be %q, got %q", "A", writer.String())
	}
}

func expectCell(t *testing.T, i *Interpreter, idx int, val uint8) {
	if i.Tape[idx] != val {
		t.Errorf("expected cell %d to be %d, got %d", idx, val, i.Tape[idx])
	}
}

func mustNew(t *testing.T, prog string) *Interpreter {
	i, err := New(prog)
	if err != nil {
		t.Errorf("unexpected error during interpreter creation: %v", err)
	}
	return i
}

func mustRun(t *testing.T, i *Interpreter) {
	if err := i.Run(); err != nil {
		t.Errorf("unexpected error during execution: %v", err)
	}
}

func run(t *testing.T, prog string) *Interpreter {
	i := mustNew(t, prog)
	mustRun(t, i)
	return i
}
