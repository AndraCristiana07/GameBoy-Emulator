package main

import "fmt"

// TODO toate cazurile din lista https://rgbds.gbdev.io/docs/v0.9.2/gbz80.7
// conventii
// - pentru paranteze drepte foloseste keyword-ul "mem" in numele functiei
//   e.g. LD A,[n16] -> ldamemn16

func (cpu *CPU) getImmediate8() uint8 {
	res := cpu.memoryRead(cpu.Registers.PC)
	cpu.Registers.PC++
	return res
}

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

// XOR - writes to A
func (cpu *CPU) xorA(other uint8) {
	res := cpu.Registers.A ^ other

	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, false)

	cpu.Registers.A = res
}

// LD r16,n16
func (cpu *CPU) ldr16n16(r16setter func(value uint16)) {
	n16 := cpu.getImmediate16()
	r16setter(n16)
}

// LD r8,n8
func (cpu *CPU) ldr8n8(register *uint8) {
	n8 := cpu.getImmediate8()
	*register = n8
}

// LD [HL],r8
func (cpu *CPU) ldmemhlr8(register uint8) {
	cpu.memoryWrite(cpu.Registers.getHL(), register)
}

// DEC r8
func (cpu *CPU) decr8(register *uint8) {
	reg := *register
	*register--
	cpu.Registers.setFlag(flagZ, *register == 0)
	cpu.Registers.setFlag(flagN, true)
	// Set if borrow from bit 4
	halfCarry := (reg & 0b10000) != (*register & 0b10000)
	cpu.Registers.setFlag(flagH, halfCarry)
}

// JR cc, e8
func (cpu *CPU) jrCCe8(flag bool) int {
	// If flag then add n to current
	// address and jump to it
	e8 := int8(cpu.getImmediate8())
	cpulogger.Debug(fmt.Sprintf("jrCCe8 e8 is: %X", e8))
	if flag {
		cpu.jump(uint16(int32(cpu.Registers.PC) + int32(e8)))
		return 12
	} else {
		return 8
	}
}

// LDH [a8], r8
func (cpu *CPU) ldhmema8r8(register uint8) {
	a8 := cpu.getImmediate8()
	addr := 0xFF00 + uint16(a8)

	if a8 == 0 {
		cpulogger.Info(fmt.Sprintf("%x %x %x", cpu.Registers.PC, cpu.Registers.A, addr))
		cpulogger.Info(fmt.Sprintf("%x %x %x %x", cpu.memoryRead(0xFF00), cpu.memoryRead(0xFF00+1), cpu.memoryRead(0xFF00+2), cpu.memoryRead(0xFF00+3)))
	}

	// cpu.Memory[0xFF00+a8] = reg
	cpu.memoryWrite(addr, register)
	cpulogger.Info(fmt.Sprintf("after operation %x ", cpu.memoryRead(addr)))
}

// LDH r8, [a8]
func (cpu *CPU) ldhr8mema8(register *uint8) {
	a8 := cpu.getImmediate8()
	addr := 0xFF00 + uint16(a8)
	cpulogger.Debug(fmt.Sprintf("-ldh r8, [a8] %04X %04X - a8: %04X", addr, cpu.memoryRead(addr), a8))
	*register = cpu.memoryRead(addr)

}

// compares with A
func (cpu *CPU) cpAn8(other uint8) {
	res := uint16(cpu.Registers.A) - uint16(other)
	// if borrow from bit 4
	halfCarry := (cpu.Registers.A & 0b10000) != uint8(res&0b10000)
	// if borrow
	carryFlag := other > cpu.Registers.A
	cpu.Registers.setFlag(flagZ, uint8(res) == 0)
	cpu.Registers.setFlag(flagN, true)
	cpu.Registers.setFlag(flagH, halfCarry)
	cpu.Registers.setFlag(flagC, carryFlag)
}

// SUB A, r8
func (cpu *CPU) subAr8(other uint8) {
	res := cpu.Registers.A - other
	// if borrow from bit 4
	halfCarry := (cpu.Registers.A & 0b01111) != uint8(res&0b10000)
	// if borrow
	carryFlag := other > cpu.Registers.A
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, true)
	cpu.Registers.setFlag(flagH, halfCarry)
	cpu.Registers.setFlag(flagC, carryFlag)
	cpu.Registers.A = res
}

// LD [a16], A
func (cpu *CPU) ldmema16A() {
	a16 := cpu.getImmediate16()
	cpu.memoryWrite(a16, cpu.Registers.A)
}

// LD SP, n16
func (cpu *CPU) ldspn16() {
	n16 := cpu.getImmediate16()
	cpu.Registers.SP = n16
}

// RST
func (cpu *CPU) rst(address uint16) {
	cpu.push(cpu.Registers.PC)
	cpu.jump(address)
}

// LD r8, [HL]
func (cpu *CPU) ldr8memhl(register *uint8) {
	*register = cpu.memoryRead(cpu.Registers.getHL())
	cpulogger.Debug(fmt.Sprintf("hl in ld %x and in memory %x", cpu.Registers.getHL(), cpu.Memory[cpu.Registers.getHL()]))

}

// LDH [C], A
func (cpu *CPU) ldhmemca() {
	address := 0xFF00 + uint16(cpu.Registers.C)
	cpu.memoryWrite(address, cpu.Registers.A)
}

// INC r8
func (cpu *CPU) incr8(register *uint8) {
	*register++
	cpu.Registers.setFlag(flagZ, *register == 0)
	cpu.Registers.setFlag(flagN, false)
	// Set if overflow from bit 3
	halfCarry := *register > 0b1111
	cpu.Registers.setFlag(flagH, halfCarry)
}

// CALL a16
func (cpu *CPU) calla16() {
	a16 := cpu.getImmediate16()
	cpu.push(cpu.Registers.PC)
	cpu.jump(a16)
}

// DEC r16
func (cpu *CPU) decr16(r16getter func() uint16, r16setter func(value uint16)) {
	value := r16getter()
	value--
	r16setter(value)
}

// LD r8, r8
func (cpu *CPU) ldr8r8(dest *uint8, source uint8) {
	*dest = source
}

// OR - writes to A
func (cpu *CPU) orA(other uint8) {
	res := cpu.Registers.A | other

	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, false)
	cpu.Registers.A = res
}

// RET
func (cpu *CPU) ret() {
	cpu.Registers.PC = cpu.pop()
}

// EI
func (cpu *CPU) ei() {
	cpu.IMEScheduled = true
}

// AND - writes to A
func (cpu *CPU) andA(other uint8) {
	res := cpu.Registers.A & other

	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, true)
	cpu.Registers.setFlag(flagC, false)
	cpu.Registers.A = res
}

// RET CC
func (cpu *CPU) retcc(flag bool) int {
	if flag {
		cpu.Registers.PC = cpu.pop()
		return 20
	} else {
		return 8
	}
}

// LD r8, [a16]
func (cpu *CPU) ldr8mema16(register *uint8) {
	a16 := cpu.getImmediate16()
	addr := cpu.memoryRead(a16)
	*register = addr
}

// SCF
func (cpu *CPU) scf() {
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, true)
}

// ADD HL, r16
func (cpu *CPU) addhlr16(r16getter func() uint16) {
	hl := cpu.Registers.getHL()
	res := hl + r16getter()
	cpu.Registers.setHL(res)
	cpu.Registers.setFlag(flagN, false)
	//halfCarry := res > 0b111111111111
	halfCarry := ((hl & 0b111111111111) + (r16getter() & 0b111111111111)) > 0b111111111111
	cpu.Registers.setFlag(flagH, halfCarry)
	//carry := res > 0b1111111111111111
	carry := uint32(hl)+uint32(r16getter()) > 0b1111111111111111
	cpu.Registers.setFlag(flagC, carry)
}

// INC [HL]
func (cpu *CPU) incmemhl() {
	hl := cpu.Registers.getHL()
	val := cpu.memoryRead(hl)
	newVal := val + 1
	cpu.memoryWrite(hl, newVal)
	cpu.Registers.setFlag(flagZ, newVal == 0)
	cpu.Registers.setFlag(flagN, false)
	// Set if overflow from bit 3
	//halfCarry := newVal > 0b1111
	halfCarry := ((val & 0b1111) + 1) > 0b1111
	cpu.Registers.setFlag(flagH, halfCarry)
}

// RETI
func (cpu *CPU) reti() {
	cpu.ei()
	cpu.ret()
}

// CPL
func (cpu *CPU) cpl() {
	cpu.Registers.A = ^cpu.Registers.A
	cpu.Registers.setFlag(flagN, true)
	cpu.Registers.setFlag(flagH, true)
}

// SWAP r8
func (cpu *CPU) swapr8(register *uint8) {
	val := *register
	swap := (val>>4)&0b00001111 | (val<<4)&0b11110000
	*register = swap
	cpu.Registers.setFlag(flagZ, swap == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, false)
}

// ADD A, r8
func (cpu *CPU) addar8(register uint8) {
	a := cpu.Registers.A
	res := a + register
	cpu.Registers.A = res
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	// Set if overflow from bit 3
	halfCarry := ((a & 0b1111) + (register & 0b1111)) > 0b1111
	//halfCarry := res > 0b1111
	cpu.Registers.setFlag(flagH, halfCarry)
	// Set if overflow from bit 7
	carry := (uint16(a) + uint16(register)) > 0b11111111
	cpu.Registers.setFlag(flagC, carry)
}

// INC HL
func (cpu *CPU) inchl() {
	hl := cpu.Registers.getHL()
	cpu.Registers.setHL(hl + 1)
}

// RES u3, r8
func (cpu *CPU) resu3r8(u3 uint8, register *uint8) {
	value := *register
	*register = value &^ (1 << u3)
}

// LD [r16], A
func (cpu *CPU) ldmemr16A(r16getter func() uint16) {
	cpulogger.Debug(fmt.Sprintf("addr %x in memory %x register A %x", r16getter(), cpu.Memory[r16getter()], cpu.Registers.A))
	cpu.memoryWrite(r16getter(), cpu.Registers.A)

}

// INC r16
func (cpu *CPU) incr16(r16getter func() uint16, r16setter func(value uint16)) {
	value := r16getter()
	value++
	r16setter(value)
}

// LD A,[r16]
func (cpu *CPU) ldamemr16(r16getter func() uint16) {
	r16 := cpu.memoryRead(r16getter())
	cpu.Registers.A = r16
}

// LD [HLD], A
func (cpu *CPU) ldmemhlda() {
	hl := cpu.Registers.getHL()
	cpu.memoryWrite(hl, cpu.Registers.A)
	cpu.Registers.setHL(hl - 1)
}

// LD [HLI], A
func (cpu *CPU) ldmemhlia() {
	hl := cpu.Registers.getHL()
	cpu.memoryWrite(hl, cpu.Registers.A)
	cpu.Registers.setHL(hl + 1)
}

// JP cc, a16
func (cpu *CPU) jpCCa16(flag bool) int {
	a16 := cpu.getImmediate16()
	if flag {
		cpu.jump(a16)
		return 16
	} else {
		return 12
	}
}

// JR e8
func (cpu *CPU) jre8() {
	e8 := int8(cpu.getImmediate8())
	cpu.jump(uint16(int32(cpu.Registers.PC) + int32(e8)))
}

// LD A, [HLI]
func (cpu *CPU) ldamemhli() {
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	cpu.Registers.A = value
	cpu.Registers.setHL(hl + 1)
}

// ADD A, n8
func (cpu *CPU) addan8() {
	n8 := cpu.getImmediate8()
	a := cpu.Registers.A
	res := a + n8

	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	halfCarry := ((a & 0b1111) + (n8 & 0b1111)) > 0b1111
	cpu.Registers.setFlag(flagH, halfCarry)
	carry := (uint16(a) + uint16(n8)) > 0b11111111
	cpu.Registers.setFlag(flagC, carry)
	cpu.Registers.A = res

}

// BIT u3, r8
func (cpu *CPU) bitu3r8(u3 uint8, register *uint8) {
	value := *register
	res := value & (1 << u3)
	cpu.Registers.setFlag(flagZ, res == 0x00)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, true)
}

// RES u3, [HL]
func (cpu *CPU) resu3memhl(u3 uint8) {
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	res := value &^ (1 << u3)
	cpu.memoryWrite(hl, res)
}

// DEC [HL]
func (cpu *CPU) decmemhl() {
	hl := cpu.Registers.getHL()
	val := cpu.memoryRead(hl)
	newVal := val - 1
	cpu.memoryWrite(hl, newVal)
	cpu.Registers.setFlag(flagZ, newVal == 0)
	cpu.Registers.setFlag(flagN, true)
	// Set if borrow from bit 4
	//halfCarry := (uint8(newVal) & 0b10000) != (uint8(val) & 0b10000) //TODO?
	halfCarry := val&0b1111 == 0

	cpu.Registers.setFlag(flagH, halfCarry)

}

// SBC A, r8
func (cpu *CPU) sbcar8(register uint8) {
	carry := uint8(0)
	if cpu.Registers.getFlag(flagC) {
		carry = uint8(1)
	}
	n := register + carry
	a := cpu.Registers.A
	res := a - n
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, true)
	// Set if borrow from bit 4
	halfCarry := (a & 0b10000) != (res & 0b10000)
	cpu.Registers.setFlag(flagH, halfCarry)
	// Set if borrow (i.e. if (r8 + carry) > A)
	newCarry := n > a
	cpu.Registers.setFlag(flagC, newCarry)
	cpu.Registers.A = res
}

// RLCA
func (cpu *CPU) rlca() {
	carry := (cpu.Registers.A & 0b10000000) >> 7
	res := (cpu.Registers.A << 1) | carry
	cpu.Registers.A = res
	cpu.Registers.setFlag(flagZ, false)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, carry == 0x01)

}

// LD A, [HLD]
func (cpu *CPU) ldamemhld() {
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	cpu.Registers.A = value
	cpu.Registers.setHL(hl - 1)
}

// ADD HL, SP
func (cpu *CPU) addhlsp() {
	hl := cpu.Registers.getHL()
	res := hl + cpu.Registers.SP
	cpu.Registers.setHL(res)
}

// LD SP, HL
func (cpu *CPU) ldspr16(r16getter func() uint16) {
	reg := r16getter()
	sp := cpu.Registers.SP
	res := sp + reg
	cpu.Registers.SP = res
	cpu.Registers.setFlag(flagN, false)
	//halfCarry := res > 0b111111111111
	halfCarry := ((sp & 0b111111111111) + (reg & 0b111111111111)) > 0b111111111111
	cpu.Registers.setFlag(flagH, halfCarry)
	//carry := res > 0b1111111111111111
	carry := ((sp & 0b1111111111111111) + (reg & 0b1111111111111111)) > 0b1111111111111111
	cpu.Registers.setFlag(flagC, carry)
}

// SLA r8
func (cpu *CPU) slar8(register *uint8) {
	val := *register
	res := val << 1
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	// set when a rotate/shift operation shifts out a “1” bit
	carry := (val&0b10000000)>>7 == 0x01
	cpu.Registers.setFlag(flagH, carry)
}

// BIT u3, [HL]
func (cpu *CPU) bitu3memhl(u3 uint8) {
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	res := value & (1 << u3)
	cpu.Registers.setFlag(flagZ, res == 0x00)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, true)
}

// INC SP
func (cpu *CPU) incsp() {
	sp := cpu.Registers.SP
	cpu.Registers.SP = sp + 1
}

// DEC SP
func (cpu *CPU) decsp() {
	sp := cpu.Registers.SP
	cpu.Registers.SP = sp - 1
}

// ADD SP, e8
func (cpu *CPU) addspe8() {
	sp := cpu.Registers.SP
	e8 := int8(cpu.getImmediate8())
	res := uint16(int32(sp) + int32(e8))
	cpu.Registers.setFlag(flagZ, false)
	cpu.Registers.setFlag(flagN, false)
	// Set if overflow from bit 3
	halfCarry := res > 0b1111
	cpu.Registers.setFlag(flagH, halfCarry)
	//Set if overflow from bit 7
	carry := res > 0b11111111
	cpu.Registers.setFlag(flagC, carry)
	cpu.Registers.SP = res

}

// SUB A, r8
func (cpu *CPU) subar8(register *uint8) {
	a := cpu.Registers.A
	value := *register
	res := a - value
	// Set if borrow from bit 4
	halfCarry := (a & 0b10000) < (value & 0b10000)
	// Set if borrow (i.e. if r8 > A)
	carryFlag := value > a
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, true)
	cpu.Registers.setFlag(flagH, halfCarry)
	cpu.Registers.setFlag(flagC, carryFlag)
	cpu.Registers.A = res & 0b11111111
}
