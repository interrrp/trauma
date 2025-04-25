package interpreter

import (
	"bytes"
	"io"
	"strings"
	"testing"
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

func TestMovePtr(t *testing.T) {
	r := mustRun(t, ">>><")
	if r.tapePtr != 2 {
		t.Errorf("expected tape pointer to be 1, got %d", r.tapePtr)
	}

	r = mustRun(t, "+>++>><+++")
	assertCell(t, r, 0, 1)
	assertCell(t, r, 1, 2)
	assertCell(t, r, 2, 3)
	assertCell(t, r, 3, 0)

	if _, err := Run("<"); err == nil {
		t.Error("expected error for tape pointer underflow")
	}
}

func TestLoop(t *testing.T) {
	r := mustRun(t, "+++[-]")
	assertCell(t, r, 0, 0)

	if _, err := Run("[[[]]]"); err != nil {
		t.Errorf("error on valid syntax: %v", err)
	}

	if _, err := Run("[]]"); err == nil {
		t.Error("expected error on invalid syntax")
	}

	if _, err := Run("[[["); err == nil {
		t.Error("expected error on invalid syntax")
	}
}

func TestInput(t *testing.T) {
	r, err := RunWithCustomIO(",", strings.NewReader("A"), io.Discard)
	if err != nil {
		t.Errorf("error during execution: %v", err)
	}
	assertCell(t, r, 0, 'A')
}

func TestOutput(t *testing.T) {
	var writer bytes.Buffer
	_, err := RunWithCustomIO("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++.", strings.NewReader(""), &writer)
	if err != nil {
		t.Errorf("error during execution: %v", err)
	}

	if writer.String() != "A" {
		t.Errorf("expected output to be %q, got %q", "A", writer.String())
	}
}

func mustRun(t *testing.T, program string) *Result {
	r, err := Run(program)
	if err != nil {
		t.Errorf("error during execution: %v", err)
	}
	return r
}

func assertCell(t *testing.T, r *Result, idx int, val byte) {
	actual := r.Tape()[idx]
	if actual != val {
		t.Errorf("expected cell %d to be %d, got %d", idx, val, actual)
	}
}
