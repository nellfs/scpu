package main

import (
	"fmt"

	"github.com/nellfs/scpu/internal/hardware"
)

func main() {
	bus := &hardware.Bus{}

	// Example program: ADC #10; ADC #20
	bus.Mem[0x8000] = 0x69 // ADC immediate
	bus.Mem[0x8001] = 10
	bus.Mem[0x8002] = 0x69 // ADC immediate
	bus.Mem[0x8003] = 20

	// Set the reset vector to 0x8000
	bus.Mem[0xFFFC] = 0x00
	bus.Mem[0xFFFD] = 0x80

	// Print the memory program
	fmt.Printf("Program bytes:\n")
	for i := uint16(0x8000); i <= 0x8003; i++ {
		fmt.Printf("$%04X: 0x%02X\n", i, bus.Mem[i])
	}
}
