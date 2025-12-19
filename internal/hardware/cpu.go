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
	PC uint16 //program counter
	SP byte   //stack pointer

	A, X, Y byte //registers

	P byte

	Bus *Bus
}

func NewCpu(bus *Bus) *CPU {
	return &CPU{
		Bus: bus,
	}
}

func (c *CPU) Reset() {
	c.A = 0
	c.X = 0
	c.Y = 0
	c.SP = 0xFD
	c.P = FlagU
	lo := uint16(c.Bus.Read(0xFFFC))
	hi := uint16(c.Bus.Read(0xFFFD))
	c.PC = hi<<8 | lo
}
