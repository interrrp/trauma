package bytecode

import "fmt"

type Instruction interface {
	String() string
}

type CellInc struct{ amount int }

func (ci *CellInc) Amount() int    { return ci.amount }
func (ci *CellInc) String() string { return fmt.Sprintf("CellInc %d", ci.Amount()) }

type PtrInc struct{ amount int }

func (pi *PtrInc) Amount() int    { return pi.amount }
func (pi *PtrInc) String() string { return fmt.Sprintf("PtrInc %d", pi.Amount()) }
