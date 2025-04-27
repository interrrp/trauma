package bytecode

import "fmt"

type Kind string

const (
	Inc       = "Inc"
	Clear     = "Clear"
	Move      = "Move"
	IncPtr    = "IncPtr"
	LoopStart = "LoopStart"
	LoopEnd   = "LoopEnd"
	Input     = "Input"
	Output    = "Output"
)

type Instruction struct {
	kind   Kind
	amount int
}

func (i *Instruction) Kind() Kind  { return i.kind }
func (i *Instruction) Amount() int { return i.amount }

func (i *Instruction) String() string {
	if i.Amount() == 0 {
		return string(i.Kind())
	}
	return fmt.Sprintf("%s %d", i.Kind(), i.Amount())
}
