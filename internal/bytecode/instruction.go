package bytecode

import "fmt"

type Instruction interface {
	Name() string
	String() string
}

type Inc struct{ amount int }

func (i *Inc) Name() string   { return "Inc" }
func (i *Inc) String() string { return fmt.Sprintf("%s %d", i.Name(), i.Amount()) }
func (i *Inc) Amount() int    { return i.amount }

type Clear struct{}

func (c *Clear) Name() string   { return "Clear" }
func (c *Clear) String() string { return c.Name() }

type Move struct{ distance int }

func (m *Move) Name() string   { return "Move" }
func (m *Move) String() string { return fmt.Sprintf("%s %d", m.Name(), m.Distance()) }
func (m *Move) Distance() int  { return m.distance }

type IncPtr struct{ amount int }

func (pi *IncPtr) Name() string   { return "IncPtr" }
func (pi *IncPtr) String() string { return fmt.Sprintf("%s %d", pi.Name(), pi.Amount()) }
func (pi *IncPtr) Amount() int    { return pi.amount }

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
