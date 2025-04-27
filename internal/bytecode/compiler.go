package bytecode

type Bytecode []Instruction

func Compile(program string) (Bytecode, error) {
	if len(program) == 0 {
		return Bytecode{}, nil
	}

	c := &compiler{
		program:     program,
		programPtr:  0,
		currentChar: program[0],
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

	currentChar byte
	nextChar    byte
}

func (c *compiler) compile() (Bytecode, error) {
	b := Bytecode{}

	for c.programPtr < len(c.program) {
		switch c.currentChar {
		case '+', '-':
			if amount := c.sumRepeatableCommands('+', '-'); amount != 0 {
				b = append(b, Instruction{Inc, amount})
			}

		case '>', '<':
			if amount := c.sumRepeatableCommands('>', '<'); amount != 0 {
				b = append(b, Instruction{IncPtr, amount})
			}

		case '[':
			if distance, ok := c.findMovePattern(); ok {
				b = append(b, Instruction{Move, distance})
			} else if c.findPattern("[-]") || c.findPattern("[+]") {
				b = append(b, Instruction{Clear, 0})
			} else {
				b = append(b, Instruction{LoopStart, 0})
			}

		case ']':
			b = append(b, Instruction{LoopEnd, 0})

		case ',':
			b = append(b, Instruction{Input, 0})
		case '.':
			b = append(b, Instruction{Output, 0})
		}

		c.advance()
	}

	return b, nil
}

func (c *compiler) findMovePattern() (int, bool) {
	remaining := len(c.program) - c.programPtr
	if remaining < 6 {
		return 0, false
	}
	if c.currentChar != '[' || c.nextChar != '-' {
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

func (c *compiler) sumRepeatableCommands(pos, neg byte) int {
	var amount int
	if c.currentChar == pos {
		amount = 1
	} else if c.currentChar == neg {
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

func (c *compiler) findPattern(s string) bool {
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
		c.currentChar = '?'
		c.nextChar = '?'
		return
	}

	c.currentChar = c.program[c.programPtr]
	if c.programPtr+1 < len(c.program) {
		c.nextChar = c.program[c.programPtr+1]
	} else {
		c.nextChar = '?'
	}
}
