package main

// TODO toate cazurile din lista https://rgbds.gbdev.io/docs/v0.9.2/gbz80.7

func (cpu *CPU) getImmediate16() uint16 {
	val1 := cpu.memoryRead(cpu.Registers.PC)
	val2 := cpu.memoryRead(cpu.Registers.PC + 1)
	res := uint16(val1) | uint16(val2)<<8
	cpu.Registers.PC += 2
	return res
}

func (cpu *CPU) jump(addr uint16) {
	cpu.Registers.PC = addr
}

// writes to A
func (cpu *CPU) xor(other uint8) {
	res := cpu.Registers.A ^ other

	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, false)

	cpu.Registers.A = res
}
