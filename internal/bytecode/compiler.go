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
				b = append(b, &Inc{amount})
			}

		case '>', '<':
			if amount := c.handleInc('>', '<'); amount != 0 {
				b = append(b, &IncPtr{amount})
			}

		case '[':
			if dist, ok := c.patternMove(); ok {
				b = append(b, &Move{dist})
			} else if c.pattern("[-]") || c.pattern("[+]") {
				b = append(b, &Clear{})
			} else {
				b = append(b, &LoopStart{})
			}

		case ']':
			b = append(b, &LoopEnd{})

		case ',':
			b = append(b, &Input{})
		case '.':
			b = append(b, &Output{})
		}

		c.advance()
	}

	return b, nil
}

func (c *compiler) patternMove() (int, bool) {
	remaining := len(c.program) - c.programPtr
	if remaining < 6 {
		return 0, false
	}
	if c.currChar != '[' || c.nextChar != '-' {
		return 0, false
	}
	i := c.programPtr + 2
	countR := 0
	for ; i < len(c.program) && c.program[i] == '>'; i++ {
		countR++
	}
	if countR < 1 {
		return 0, false
	}
	if i >= len(c.program) || c.program[i] != '+' {
		return 0, false
	}
	i++
	countL := 0
	for ; i < len(c.program) && c.program[i] == '<'; i++ {
		countL++
	}
	if countL != countR {
		return 0, false
	}
	if i >= len(c.program) || c.program[i] != ']' {
		return 0, false
	}
	steps := i - c.programPtr
	for j := 0; j < steps; j++ {
		c.advance()
	}
	return countR, true
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

func (c *compiler) pattern(s string) bool {
	if len(c.program)-1-c.programPtr < len(s)-1 {
		return false
	}
	if c.program[c.programPtr:c.programPtr+len(s)] == s {
		for range len(s) - 1 {
			c.advance()
		}
		return true
	}
	return false
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
