package main

// TODO toate cazurile din tabele https://gbdev.io/gb-opcodes/optables/

// nop
func (cpu *CPU) op00() int {
	return 4
}

// JP a16
func (cpu *CPU) opc3() int {
	a16 := cpu.getImmediate16()
	cpu.jump(a16)
	return 16
}

// XOR A,r8 - r8=A
func (cpu *CPU) opaf() int {
	cpu.xor(cpu.Registers.A)
	return 4
}

func (cpu *CPU) op21() int {

	return
}
