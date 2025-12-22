package hardware

var opcodeTable [256]Instruction

type Instruction struct {
	Name    string
	Execute func(c *CPU) (extracycle byte)
	Bytes   byte
	Cycles  byte
}

func init() {
	// --- ADC Immediate (#oper) ---
	opcodeTable[0x69] = Instruction{
		Name: "ADC",
		Execute: func(c *CPU) byte {
			c.ADC(c.IMM())
			return 0
		},
		Bytes:  2,
		Cycles: 2,
	}

	// --- ADC Zero Page (oper) ---
	opcodeTable[0x65] = Instruction{
		Name: "ADC",
		Execute: func(c *CPU) byte {
			addr := c.ZP()
			c.ADC(c.Bus.Read(addr))
			return 0
		},
		Bytes:  2,
		Cycles: 3,
	}

	// --- ADC Zero Page,X (oper,X) ---
	opcodeTable[0x75] = Instruction{
		Name: "ADC",
		Execute: func(c *CPU) byte {
			addr := c.ZPX()
			c.ADC(c.Bus.Read(addr))
			return 0
		},
		Bytes:  2,
		Cycles: 4,
	}

	// --- ADC Absolute (oper) ---
	opcodeTable[0x6D] = Instruction{
		Name: "ADC",
		Execute: func(c *CPU) byte {
			addr := c.ABS()
			c.ADC(c.Bus.Read(addr))
			return 0
		},
		Bytes:  3,
		Cycles: 4,
	}

	// --- ADC Absolute,X (oper,X) ---
	opcodeTable[0x7D] = Instruction{
		Name: "ADC",
		Execute: func(c *CPU) byte {
			addr, pageCrossed := c.ABSX() // ignore page crossing for now
			c.ADC(c.Bus.Read(addr))
			if pageCrossed {
				return 1
			}
			return 0
		},
		Bytes:  3,
		Cycles: 4}

	// --- ADC Absolute,Y (oper,Y) ---
	opcodeTable[0x79] = Instruction{
		Name: "ADC",
		Execute: func(c *CPU) byte {
			addr, pageCrossed := c.ABSY() // ignore page crossing for now
			c.ADC(c.Bus.Read(addr))
			if pageCrossed {
				return 1
			}
			return 0
		},
		Bytes:  3,
		Cycles: 4}

	// --- ADC (Indirect,X) ---
	opcodeTable[0x61] = Instruction{
		Name: "ADC",
		Execute: func(c *CPU) byte {
			addr := c.INDX()
			c.ADC(c.Bus.Read(addr))
			return 0
		},
		Bytes:  2,
		Cycles: 6,
	}

	// --- ADC (Indirect),Y ---
	opcodeTable[0x71] = Instruction{
		Name: "ADC",
		Execute: func(c *CPU) byte {
			addr, pageCrossed := c.INDY()
			c.ADC(c.Bus.Read(addr))
			if pageCrossed {
				return 1
			}
			return 0
		},
		Bytes:  2,
		Cycles: 5,
	}

	// --- AND - Bitwise And ---
	opcodeTable[0x29] = Instruction{
		Name: "AND",
		Execute: func(c *CPU) byte {
			addr := c.IMM()
			c.AND(c.Bus.Read(uint16(addr)))
			return 0
		}, Bytes: 2, Cycles: 2,
	}
	opcodeTable[0x25] = Instruction{
		Name: "AND",
		Execute: func(c *CPU) byte {
			addr := c.ZP()
			c.AND(c.Bus.Read(uint16(addr)))
			return 0
		}, Bytes: 2, Cycles: 3,
	}

	opcodeTable[0x35] = Instruction{
		Name: "AND",
		Execute: func(c *CPU) byte {
			addr := c.ZPX()
			c.AND(c.Bus.Read(uint16(addr)))
			return 0
		}, Bytes: 2, Cycles: 4,
	}

	opcodeTable[0x2D] = Instruction{
		Name: "AND",
		Execute: func(c *CPU) byte {
			addr := c.ABS()
			c.AND(c.Bus.Read(uint16(addr)))
			return 0
		}, Bytes: 3, Cycles: 4,
	}

	opcodeTable[0x3D] = Instruction{
		Name: "AND",
		Execute: func(c *CPU) byte {
			addr, pageCrossed := c.ABSX()
			c.AND(c.Bus.Read(uint16(addr)))
			if pageCrossed {
				return 1
			}
			return 0
		}, Bytes: 3, Cycles: 4,
	}

	opcodeTable[0x39] = Instruction{
		Name: "AND",
		Execute: func(c *CPU) byte {
			addr, pageCrossed := c.ABSY()
			c.AND(c.Bus.Read(uint16(addr)))
			if pageCrossed {
				return 1
			}
			return 0
		}, Bytes: 3, Cycles: 4,
	}

	opcodeTable[0x21] = Instruction{
		Name: "AND",
		Execute: func(c *CPU) byte {
			addr := c.INDX()
			c.AND(c.Bus.Read(uint16(addr)))
			return 0
		}, Bytes: 2, Cycles: 6,
	}
	opcodeTable[0x31] = Instruction{
		Name: "AND",
		Execute: func(c *CPU) byte {
			addr, pageCrossed := c.INDY()
			c.AND(c.Bus.Read(uint16(addr)))
			if pageCrossed {
				return 1
			}
			return 0
		}, Bytes: 2, Cycles: 6,
	}
}

func (c *CPU) ADC(value byte) {
	carry := byte(0)
	if c.GetFlag(FlagC) {
		carry = 1
	}
	result := uint16(c.A) + uint16(value) + uint16(carry)

	c.SetFlag(FlagC, result > 0xFF)
	c.SetFlag(FlagV, ((^(c.A ^ value))&(c.A^byte(result))&0x80) != 0)

	c.A = byte(result & 0xFF)
	c.SetFlag(FlagZ, c.A == 0)
	c.SetFlag(FlagN, c.A&0x80 != 0)
}

func (c *CPU) AND(value byte) {
	c.A = c.A & value
	c.SetFlag(FlagZ, c.A == 0)
	c.SetFlag(FlagN, c.A&0x80 != 0)
}

// func (c *CPU) ASL(addr uint16) {
// 	value := c.Bus.Read(addr)

// 	c.SetFlag(FlagC, value&0x80 != 0)

// 	value <<= 1
// 	c.Bus.Write(addr, value)

// 	c.SetFlag(FlagZ, value == 0)
// 	c.SetFlag(FlagN, value&0x80 != 0)
// }
