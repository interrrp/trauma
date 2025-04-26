package bytecode

import "testing"

func TestCellInc(t *testing.T) {
	assertBytecode(t, "+++--", "CellInc 1")
	assertBytecode(t, "++---", "CellInc -1")
	assertBytecode(t, "+-")
}

func TestPtrInc(t *testing.T) {
	assertBytecode(t, ">>><<", "PtrInc 1")
	assertBytecode(t, ">><<<", "PtrInc -1")
	assertBytecode(t, "><")
}

func TestLoops(t *testing.T) {
	assertBytecode(t, "[+[>>]--]",
		"LoopStart",
		"CellInc 1",
		"LoopStart",
		"PtrInc 2",
		"LoopEnd",
		"CellInc -2",
		"LoopEnd")
}

func TestIO(t *testing.T) {
	assertBytecode(t, "++[,].>.",
		"CellInc 2",
		"LoopStart",
		"Input",
		"LoopEnd",
		"Output",
		"PtrInc 1",
		"Output")

	assertBytecode(t, ",,,...",
		"Input",
		"Input",
		"Input",
		"Output",
		"Output",
		"Output")
}

func assertBytecode(t *testing.T, program string, expected ...string) {
	b, err := Compile(program)
	if err != nil {
		t.Errorf("error during compilation: %v", err)
	}

	if len(expected) != len(b) {
		t.Logf("expected %d instruction(s), got %d:", len(expected), len(b))
		for i, inst := range b {
			t.Logf("%d. %s", i+1, inst.String())
		}
		t.Fail()
	}

	for i, inst := range b {
		exp := expected[i]
		real := inst.String()
		if exp != real {
			t.Errorf("expected %q at index %d, got %q", exp, i, real)
		}
	}
}
