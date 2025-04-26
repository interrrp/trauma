package bytecode

type Bytecode []Instruction

func Compile(program string) (Bytecode, error) {
	if len(program) == 0 {
		return Bytecode{}, nil
	}

	c := &compiler{
		program:    program,
		programPtr: 0,
		currChar:   program[0],
	}

	if len(program) > 1 {
		c.nextChar = program[1]
	} else {
		c.nextChar = '?'
	}

	return c.compile()
}

type compiler struct {
	program    string
	programPtr int

	currChar byte
	nextChar byte
}

func (c *compiler) compile() (Bytecode, error) {
	b := Bytecode{}

	for c.programPtr < len(c.program) {
		switch c.currChar {
		case '+', '-':
			if amount := c.handleInc('+', '-'); amount != 0 {
				b = append(b, &CellInc{amount})
			}

		case '>', '<':
			if amount := c.handleInc('>', '<'); amount != 0 {
				b = append(b, &PtrInc{amount})
			}
		}

		c.advance()
	}

	return b, nil
}

func (c *compiler) handleInc(pos, neg byte) int {
	var amount int
	if c.currChar == pos {
		amount = 1
	} else if c.currChar == neg {
		amount = -1
	}

	for c.nextChar == pos || c.nextChar == neg {
		if c.nextChar == pos {
			amount++
		} else if c.nextChar == neg {
			amount--
		}
		c.advance()
	}

	return amount
}

func (c *compiler) advance() {
	c.programPtr++
	if c.programPtr >= len(c.program) {
		c.currChar = '?'
		c.nextChar = '?'
		return
	}

	c.currChar = c.program[c.programPtr]
	if c.programPtr+1 < len(c.program) {
		c.nextChar = c.program[c.programPtr+1]
	} else {
		c.nextChar = '?'
	}
}
