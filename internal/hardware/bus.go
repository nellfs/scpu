package hardware

type Bus struct {
	Mem [64 * 1024]byte
}

func NewBus() *Bus {
	return &Bus{}
}

// Returns a byte from memory
func (b *Bus) Read(addr uint16) byte {
	return b.Mem[addr]
}

func (b *Bus) Write(addr uint16, data byte) {
	b.Mem[addr] = data
}
