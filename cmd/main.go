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

	cpu := hardware.NewCpu(bus)
	cpu.Reset()

	fmt.Println("Starting CPU execution trace...")

	for i := 0; i < 2; i++ {
		pc := cpu.PC
		opcode := bus.Mem[pc]

		var value byte
		switch opcode {
		case 0x69: // Only for debug
			value = bus.Mem[pc+1]
		default:
			value = 0
		}

		fmt.Printf("Before Step: PC: $%04X | Opcode: 0x%02X | A: %d | Adding: %d | C:%v Z:%v N:%v V:%v\n",
			pc, opcode, cpu.A, value,
			cpu.GetFlag(hardware.FlagC),
			cpu.GetFlag(hardware.FlagZ),
			cpu.GetFlag(hardware.FlagN),
			cpu.GetFlag(hardware.FlagV))

		cpu.Step()

		fmt.Printf("After Step:  PC: $%04X | A: %d | C:%v Z:%v N:%v V:%v\n\n",
			cpu.PC, cpu.A,
			cpu.GetFlag(hardware.FlagC),
			cpu.GetFlag(hardware.FlagZ),
			cpu.GetFlag(hardware.FlagN),
			cpu.GetFlag(hardware.FlagV))
	}

	// Print final state
	fmt.Println("------")
	fmt.Println("\nFinal CPU state:")
	fmt.Printf("A: %d\n", cpu.A)
	fmt.Printf("Carry flag: %v\n", cpu.GetFlag(hardware.FlagC))
	fmt.Printf("Zero flag: %v\n", cpu.GetFlag(hardware.FlagZ))
	fmt.Printf("Negative flag: %v\n", cpu.GetFlag(hardware.FlagN))
	fmt.Printf("Overflow flag: %v\n", cpu.GetFlag(hardware.FlagV))
}
