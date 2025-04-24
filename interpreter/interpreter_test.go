package interpreter

import (
	"bytes"
	"strings"
	"testing"
)

func TestIncrement(t *testing.T) {
	i := run(t, "+++")
	if i.Tape[0] != 3 {
		t.Errorf("expected first cell to be 3, got %d", i.Tape[0])
	}

	i = run(t, "+++++")
	if i.Tape[0] != 5 {
		t.Errorf("expected first cell to be 5, got %d", i.Tape[0])
	}
}

func TestDecrement(t *testing.T) {
	i := run(t, "++-")
	if i.Tape[0] != 1 {
		t.Errorf("expected first cell to be 1, got %d", i.Tape[0])
	}

	i = run(t, "---")
	if i.Tape[0] != 255-2 {
		t.Errorf("expected first cell to overflow to %d, got %d", 255-2, i.Tape[0])
	}
}

func TestMovePtr(t *testing.T) {
	i := run(t, ">>><")
	if i.tapePtr != 2 {
		t.Errorf("expected tape pointer to be 1, got %d", i.tapePtr)
	}

	i = run(t, "+>++>><+++")
	if i.Tape[0] != 1 {
		t.Errorf("expected first cell to be 1, got %d", i.Tape[0])
	}
	if i.Tape[1] != 2 {
		t.Errorf("expected second cell to be 2, got %d", i.Tape[1])
	}
	if i.Tape[2] != 3 {
		t.Errorf("expected third cell to be 3, got %d", i.Tape[2])
	}
	if i.Tape[3] != 0 {
		t.Errorf("expected fourth cell to be 0, got %d", i.Tape[3])
	}
}

func TestLoop(t *testing.T) {
	i := run(t, "+++[-]")
	if i.Tape[0] != 0 {
		t.Errorf("expected first cell to be 0, got %d", i.Tape[0])
	}

	if _, err := New("[[[]]]"); err != nil {
		t.Errorf("error on valid syntax: %v", err)
	}

	if _, err := New("[]]"); err == nil {
		t.Error("expected an error on invalid syntax")
	}

	if _, err := New("[[["); err == nil {
		t.Error("expected an error on invalid syntax")
	}
}

func TestInput(t *testing.T) {
	i, err := New(",")
	if err != nil {
		t.Errorf("error during interpreter creation: %v", err)
	}
	i.Reader = strings.NewReader("A")

	if err := i.Run(); err != nil {
		t.Errorf("error during execution: %v", err)
	}

	if i.Tape[0] != 'A' {
		t.Errorf("expected first cell to be ASCII value of %q, got %d", 'a', i.Tape[0])
	}
}

func TestOutput(t *testing.T) {
	i, err := New("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++.")
	if err != nil {
		t.Errorf("error during interpreter creation: %v", err)
	}

	var writer bytes.Buffer
	i.Writer = &writer

	if err := i.Run(); err != nil {
		t.Errorf("error during execution: %v", err)
	}

	if writer.String() != "A" {
		t.Errorf("expected output to be %q, got %q", "A", writer.String())
	}
}

func run(t *testing.T, prog string) *Interpreter {
	i, err := New(prog)
	if err != nil {
		t.Errorf("error during interpreter creation: %v", err)
	}
	if err := i.Run(); err != nil {
		t.Errorf("error during execution: %v", err)
	}
	return i
}
