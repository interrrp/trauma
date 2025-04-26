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
