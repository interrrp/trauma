package bytecode

import "fmt"

type Instruction interface {
	Name() string
	String() string
}

type CellInc struct{ amount int }

func (ci *CellInc) Name() string   { return "CellInc" }
func (ci *CellInc) Amount() int    { return ci.amount }
func (ci *CellInc) String() string { return fmt.Sprintf("%s %d", ci.Name(), ci.Amount()) }

type CellEmpty struct{}

func (ce *CellEmpty) Name() string   { return "EmptyCell" }
func (ce *CellEmpty) String() string { return ce.Name() }

type PtrInc struct{ amount int }

func (pi *PtrInc) Name() string   { return "PtrInc" }
func (pi *PtrInc) Amount() int    { return pi.amount }
func (pi *PtrInc) String() string { return fmt.Sprintf("%s %d", pi.Name(), pi.Amount()) }

type LoopStart struct{}

func (ls *LoopStart) Name() string   { return "LoopStart" }
func (ls *LoopStart) String() string { return ls.Name() }

type LoopEnd struct{}

func (le *LoopEnd) Name() string   { return "LoopEnd" }
func (le *LoopEnd) String() string { return le.Name() }

type Input struct{}

func (i *Input) Name() string   { return "Input" }
func (i *Input) String() string { return i.Name() }

type Output struct{}

func (o *Output) Name() string   { return "Output" }
func (o *Output) String() string { return o.Name() }
