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
}
