package vm

type Result struct {
	tape    []byte
	tapePtr int
}

func newResult(v *vm) *Result {
	return &Result{tape: v.tape, tapePtr: v.tapePtr}
}

func (r *Result) Tape() []byte { return r.tape }
func (r *Result) TapePtr() int { return r.tapePtr }
