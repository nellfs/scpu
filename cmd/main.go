package main

import "github.com/nellfs/scpu/scpu/internal/hardware"

func main() {
	cpu := hardware.NewCpu()
	cpu.Reset()
}
