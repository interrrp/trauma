package bytecode_test

import (
	"testing"

	"github.com/interrrp/trauma/internal/bytecode"
)

func TestInc(t *testing.T) {
	assertBytecode(t, "+++--", "Inc 1")
	assertBytecode(t, "++---", "Inc -1")
	assertBytecode(t, "+-")
}

func TestIncPtr(t *testing.T) {
	assertBytecode(t, ">>><<", "IncPtr 1")
	assertBytecode(t, ">><<<", "IncPtr -1")
	assertBytecode(t, "><")
}

func TestLoops(t *testing.T) {
	assertBytecode(t, "[+[>>]--]",
		"LoopStart",
		"Inc 1",
		"LoopStart",
		"IncPtr 2",
		"LoopEnd",
		"Inc -2",
		"LoopEnd")
}

func TestIO(t *testing.T) {
	assertBytecode(t, "++[,].>.",
		"Inc 2",
		"LoopStart",
		"Input",
		"LoopEnd",
		"Output",
		"IncPtr 1",
		"Output")

	assertBytecode(t, ",,,...",
		"Input",
		"Input",
		"Input",
		"Output",
		"Output",
		"Output")
}

func TestClear(t *testing.T) {
	assertBytecode(t, "+++[-]",
		"Inc 3",
		"Clear")

	assertBytecode(t, "+[++[-]]",
		"Inc 1",
		"LoopStart",
		"Inc 2",
		"Clear",
		"LoopEnd")
}

func assertBytecode(t *testing.T, program string, expected ...string) {
	actual, err := bytecode.Compile(program)
	if err != nil {
		t.Errorf("error during compilation: %v", err)
	}

	if len(expected) != len(actual) {
		t.Logf("expected %d instruction(s), got %d:", len(expected), len(actual))
		for i, inst := range actual {
			t.Logf("%d. %s", i+1, inst.String())
		}
		t.Fail()
	}

	for i, inst := range actual {
		exp := expected[i]
		real := inst.String()
		if exp != real {
			t.Errorf("expected %q at index %d, got %q", exp, i, real)
		}
	}
}
