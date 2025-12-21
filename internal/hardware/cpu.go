package hardware

const (
	FlagC byte = 1 << 0 // carry
	FlagZ byte = 1 << 1 // zero
	FlagI byte = 1 << 2 // interrupt Disable
	FlagD byte = 1 << 3 // decimal
	FlagB byte = 1 << 4 // break
	FlagU byte = 1 << 5 // unused
	FlagV byte = 1 << 6 // overflow
	FlagN byte = 1 << 7 // negative
)

type CPU struct {
	PC uint16 // program counter
	SP byte   // stack pointer

	A, X, Y byte // registers

	P byte // status register

	cycles uint64

	Bus *Bus
}

func NewCpu(bus *Bus) *CPU {
	return &CPU{Bus: bus}
}

func (c *CPU) Reset() {
	c.A, c.X, c.Y = 0, 0, 0
	c.SP = 0xFD
	c.P = FlagU

	lo := uint16(c.Bus.Read(0xFFFC))
	hi := uint16(c.Bus.Read(0xFFFD))
	c.PC = hi<<8 | lo
}

func (c *CPU) Step() {
	opcode := c.Bus.Read(c.PC)
	c.PC++

	inst := opcodeTable[opcode]

	extra := byte(0)
	if inst.Execute != nil {
		extra = inst.Execute(c)
	}

	c.cycles = uint64(inst.Cycles + extra)
}

func (c *CPU) GetFlag(flag byte) bool {
	return (c.P & flag) != 0
}

func (c *CPU) SetFlag(flag byte, value bool) {
	if value {
		c.P |= flag
	} else {
		c.P &^= flag
	}
}
func (c *CPU) SetFlagOn(flag byte) {
	c.SetFlag(flag, true)
}

func (c *CPU) ClearFlag(flag byte) {
	c.SetFlag(flag, false)
}

// func (c *CPU) setZN(value byte) {
// 	c.SetFlag(FlagZ, value == 0)
// 	c.SetFlag(FlagN, value&0x80 != 0)
// }

// ================= Addressing Modes =================

func (c *CPU) IMP() {}

func (c *CPU) CMP(value byte) {
	result := c.A - value

	c.SetFlag(FlagC, c.A >= value)     // Carry se A >= value
	c.SetFlag(FlagZ, c.A == value)     // Zero se A == value
	c.SetFlag(FlagN, result&0x80 != 0) // Negative se bit 7 setado
}

func (c *CPU) IMM() byte {
	val := c.Bus.Read(c.PC)
	c.PC++
	return val
}

func (c *CPU) ABS() uint16 {
	lo := uint16(c.Bus.Read(c.PC))
	c.PC++
	hi := uint16(c.Bus.Read(c.PC))
	c.PC++
	return hi<<8 | lo
}

func (c *CPU) ZP() uint16 {
	addr := c.Bus.Read(c.PC)
	c.PC++
	return uint16(addr)
}

func (c *CPU) ZPX() uint16 {
	base := c.Bus.Read(c.PC)
	c.PC++
	return uint16(uint8(base + c.X))
}

func (c *CPU) ZPY() uint16 {
	base := c.Bus.Read(c.PC)
	c.PC++
	return uint16(uint8(base + c.Y))
}

func (c *CPU) ABSX() (addr uint16, pageCrossed bool) {
	base := c.ABS()
	addr = base + uint16(c.X)
	pageCrossed = (base & 0xFF00) != (addr & 0xFF00)
	return
}

func (c *CPU) ABSY() (addr uint16, pageCrossed bool) {
	base := c.ABS()
	addr = base + uint16(c.Y)
	pageCrossed = (base & 0xFF00) != (addr & 0xFF00)
	return
}

// JMP (indirect) with 6502 page-boundary bug
func (c *CPU) IND() uint16 {
	lo := uint16(c.Bus.Read(c.PC))
	c.PC++
	hi := uint16(c.Bus.Read(c.PC))
	c.PC++

	ptr := hi<<8 | lo

	var addrLo, addrHi byte
	if lo == 0xFF {
		addrLo = c.Bus.Read(ptr)
		addrHi = c.Bus.Read(ptr & 0xFF00)
	} else {
		addrLo = c.Bus.Read(ptr)
		addrHi = c.Bus.Read(ptr + 1)
	}

	return uint16(addrHi)<<8 | uint16(addrLo)
}

// (zp,X)
func (c *CPU) INDX() uint16 {
	base := c.Bus.Read(c.PC)
	c.PC++

	ptr := uint8(base + c.X)

	lo := c.Bus.Read(uint16(ptr))
	hi := c.Bus.Read(uint16(ptr + 1))

	return uint16(hi)<<8 | uint16(lo)
}

// (zp),Y
func (c *CPU) INDY() (addr uint16, pageCrossed bool) {
	base := c.Bus.Read(c.PC)
	c.PC++

	ptr := uint8(base)

	lo := c.Bus.Read(uint16(ptr))
	hi := c.Bus.Read(uint16(ptr + 1))

	baseAddr := uint16(hi)<<8 | uint16(lo)
	addr = baseAddr + uint16(c.Y)

	pageCrossed = (baseAddr & 0xFF00) != (addr & 0xFF00)
	return
}

// Relative (branches)
// func (c *CPU) REL() int8 {
// 	offset := int8(c.Bus.Read(c.PC))
// 	c.PC++
// 	return offset
// }

//
// Instructions
//

func (c *CPU) ADC(value byte) {
	carry := byte(0)

	if c.P&FlagC != 0 {
		carry = 1
	}

	result := uint16(c.A) + uint16(value) + uint16(carry)

	c.SetFlag(FlagC, result > 0xFF)
	c.SetFlag(FlagZ, c.A == 0)
	c.SetFlag(FlagN, c.A&0x80 != 0)
	c.SetFlag(FlagV, ((^(c.A ^ value))&(c.A^byte(result))&0x80) != 0)

	c.A = byte(result & 0xFF)
	c.SetFlag(FlagZ, c.A == 0)
	c.SetFlag(FlagN, c.A&0x80 != 0)
}
