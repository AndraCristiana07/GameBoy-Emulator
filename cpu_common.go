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

// STOP n8
func (cpu *CPU) stop() {
	n8 := cpu.getImmediate8()
	if n8 != 0x00 {
		cpulogger.Debug(fmt.Sprintf("STOP -> imm is 0x%02X", n8))
	}
	cpu.Registers.PC++
	cpu.stopped = true
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

// LD [HL], r8/n8
func (cpu *CPU) ldmemhl(other uint8) {
	cpu.memoryWrite(cpu.Registers.getHL(), other)
}

// DEC r8
func (cpu *CPU) decr8(register *uint8) {
	value := *register
	*register--
	cpu.Registers.setFlag(flagZ, *register == 0)
	cpu.Registers.setFlag(flagN, true)
	// Set if borrow from bit 4
	//halfCarry := (reg & 0b10000) != (*register & 0b10000)
	//halfCarry := (value&0b1111)-1 > 0b1111
	halfCarry := value & 0b1111
	cpu.Registers.setFlag(flagH, halfCarry == 0x00)
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

// /////////////////////////
// LDH r8, [a8]
func (cpu *CPU) ldhr8mema8(register *uint8) {
	a8 := cpu.getImmediate8()
	addr := 0xFF00 + uint16(a8)
	cpulogger.Debug(fmt.Sprintf("-ldh r8, [a8] %04X %04X - a8: %04X", addr, cpu.memoryRead(addr), a8))
	*register = cpu.memoryRead(addr)

}

// compares with A
func (cpu *CPU) cpa(other uint8) {
	a := cpu.Registers.A
	//res := uint16(a) - uint16(other)
	res := a - other

	// if borrow from bit 4
	//halfCarry := (cpu.Registers.A & 0b10000) != uint8(res&0b10000)
	halfCarry := (a & 0b1111) < (other & 0b1111)
	// if borrow
	carryFlag := uint16(other) > uint16(a)
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, true)
	cpu.Registers.setFlag(flagH, halfCarry)
	cpu.Registers.setFlag(flagC, carryFlag)
}

// SUB A, r8
func (cpu *CPU) subar8(other uint8) {
	a := cpu.Registers.A
	res := a - other
	// if borrow from bit 4
	//halfCarry := (cpu.Registers.A & 0b01111) != uint8(res&0b10000)
	halfCarry := (a & 0b1111) < (other & 0b1111)
	//halfCarry := (cpu.Registers.A & 0b10000) < (other & 0b10000)
	// if borrow
	carryFlag := uint16(other) > uint16(a)
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
	value := *register
	*register++
	cpu.Registers.setFlag(flagZ, *register == 0)
	cpu.Registers.setFlag(flagN, false)
	// Set if overflow from bit 3
	//halfCarry := value > 0b1111
	halfCarry := value & 0b1111
	//halfCarry := (value&0b1111)+1 > 0b1111
	cpu.Registers.setFlag(flagH, halfCarry == 0)
}

// /////////////////////////////////////
// CALL a16
func (cpu *CPU) calla16() {
	a16 := cpu.getImmediate16()
	cpu.push(cpu.Registers.PC)
	cpu.jump(a16)
}

// CALL
func (cpu *CPU) call(adress uint16) {
	cpu.push(cpu.Registers.PC)
	cpu.jump(adress)
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

// //////////////////////////
// ADD HL, r16
func (cpu *CPU) addhlr16(r16getter func() uint16) {
	hl := cpu.Registers.getHL()
	r16 := r16getter()
	res := hl + r16
	cpu.Registers.setHL(res)
	cpu.Registers.setFlag(flagN, false)
	//halfCarry := res > 0b111111111111
	halfCarry := ((hl & 0b111111111111) + (r16 & 0b111111111111)) > 0b111111111111
	cpu.Registers.setFlag(flagH, halfCarry)
	//carry := res > 0b1111111111111111
	//carry := uint32(hl)+uint32(r16) > 0b1111111111111111
	carry := uint32(hl)+uint32(r16) > 0b1111111111111111

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
	//swap := (val>>4)&0b00001111 | (val<<4)&0b11110000
	swap := (val >> 4) | (val << 4)

	*register = swap
	cpu.Registers.setFlag(flagZ, val == 0x00)
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
	r16 := r16getter()
	cpulogger.Debug(fmt.Sprintf("addr %x in memory %x register A %x", r16, cpu.Memory[r16], cpu.Registers.A))
	cpu.memoryWrite(r16, cpu.Registers.A)

}

// ///////////////////////////////
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
	cpu.hldec()
}

// LD [HLI], A
func (cpu *CPU) ldmemhlia() {
	hl := cpu.Registers.getHL()
	cpu.memoryWrite(hl, cpu.Registers.A)
	cpu.hlinc()
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
	cpu.hlinc()
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

// //////////////////////////////////
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
	halfCarry := val & 0b1111

	cpu.Registers.setFlag(flagH, halfCarry == 0x00)

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
	//halfCarry := (a & 0b10000) != (res & 0b10000)
	halfCarry := (a & 0b1111) < (n & 0b1111)

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
	cpu.hldec()
}

// ADD HL, SP
func (cpu *CPU) addhlsp() {
	hl := cpu.Registers.getHL()
	sp := cpu.Registers.SP
	res := hl + sp

	cpu.Registers.setHL(res)

	//cpu.Registers.setFlag(flagZ, false)
	cpu.Registers.setFlag(flagN, false)
	// Set if overflow from bit 11
	//halfCarry := ((hl & 0b111111111111) + (sp & 0b111111111111)) > 0b111111111111
	halfCarry := ((hl & 0b111111111111) + (sp & 0b111111111111)) > 0b111111111111
	cpu.Registers.setFlag(flagH, halfCarry)
	// Set if overflow from bit 15
	carry := uint32(hl)+uint32(sp) > 0b1111111111111111
	//carry := (hl)+(sp) > 0b1111111111111111

	cpu.Registers.setFlag(flagC, carry)

}

// LD SP, HL
func (cpu *CPU) ldsphl() {

	hl := cpu.Registers.getHL()
	cpu.Registers.SP = hl

}

// SLA r8
func (cpu *CPU) slar8(register *uint8) {
	val := *register
	res := val << 1
	// bit 0 of register is reset to 0
	//res &= 0b11111110
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	// set when a rotate/shift operation shifts out a “1” bit
	carry := val & 0b10000000
	//carry := (val & 0b10000000) >> 7

	cpu.Registers.setFlag(flagC, carry != 0x00)
	*register = res
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

// /////////////////////////////
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
	//halfCarry := res > 0b1111
	halfCarry := ((sp & 0b1111) + (uint16(e8) & 0b1111)) > 0b1111
	cpu.Registers.setFlag(flagH, halfCarry)
	//Set if overflow from bit 7
	//carry := res > 0b11111111
	carry := ((sp & 0b11111111) + (uint16(e8) & 0b11111111)) > 0b11111111
	cpu.Registers.setFlag(flagC, carry)
	cpu.Registers.SP = res

}

// SRL r8
func (cpu *CPU) srlr8(register *uint8) {
	val := *register
	res := val >> 1
	// bit 7 of register is reset to 0
	//res &= 0b01111111
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	carry := val & 0b00000001
	cpu.Registers.setFlag(flagC, carry != 0x00)
	*register = res
}

// SET u3, r8
func (cpu *CPU) setu3r8(u3 uint8, register *uint8) {
	*register |= 1 << u3
}

// HALT
func (cpu *CPU) halt() {
	//cpu.halted = true
	if cpu.IME {
		cpu.halted = true
		cpu.haltBug = false
	} else {
		IE := cpu.memoryRead(0xFFFF)
		IF := cpu.memoryRead(0xFF0F)

		interruptions := IE & IF
		if interruptions == 0 {
			cpu.halted = true
			cpu.haltBug = false
		} else {
			cpu.halted = false
			cpu.haltBug = true
		}
	}
}

// RRCA
func (cpu *CPU) rrca() {
	a := cpu.Registers.A
	carry := a & 0b00000001
	res := (cpu.Registers.A >> 1) | (carry << 7)
	//res := a >> 1
	//if carry == 0x01 {
	//	res = a>>1 | (0b10000000)
	//}
	cpu.Registers.A = res
	cpu.Registers.setFlag(flagZ, false)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, carry == 0x01)
}

// ADC A, n8
func (cpu *CPU) adcan8() {
	carry := uint8(0)
	if cpu.Registers.getFlag(flagC) {
		carry = uint8(1)
	}
	n8 := cpu.getImmediate8()
	a := cpu.Registers.A
	res := a + carry + n8
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	// Set if overflow from bit 3
	halfCarry := ((a & 0b1111) + (n8 & 0b1111) + (carry & 0b1111)) > 0b1111
	cpu.Registers.setFlag(flagH, halfCarry)
	// Set if overflow from bit 7
	newCarry := (uint16(a) + uint16(n8) + uint16(carry)) > 0b11111111
	cpu.Registers.setFlag(flagC, newCarry)
	cpu.Registers.A = res

}

// CALL CC a16
func (cpu *CPU) callcca16(flag bool) int {
	a16 := cpu.getImmediate16()
	if flag {
		cpu.push(cpu.Registers.PC)
		cpu.jump(a16)
		return 24
	} else {
		return 12
	}

}

// //////////////////////////////
// SUB A, n8
func (cpu *CPU) suban8() {
	n8 := cpu.getImmediate8()
	a := cpu.Registers.A
	res := a - n8
	// if borrow from bit 4
	//halfCarry := (cpu.Registers.A & 0b01111) != uint8(res&0b10000)
	halfCarry := (a & 0b1111) < (n8 & 0b1111)
	//halfCarry := (cpu.Registers.A & 0b10000) < (n8 & 0b10000)
	// if borrow
	carryFlag := uint16(n8) > uint16(a)
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, true)
	cpu.Registers.setFlag(flagH, halfCarry)
	cpu.Registers.setFlag(flagC, carryFlag)
	cpu.Registers.A = res
}

// CP A, [HL]
func (cpu *CPU) cpamemhl() {
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	a := cpu.Registers.A
	res := a - value
	// if borrow from bit 4
	//halfCarry := cpu.Registers.A&0b10000 != res&0b10000
	halfCarry := (a & 0b1111) < (value & 0b1111)
	// if borrow
	carryFlag := value > a
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, true)
	cpu.Registers.setFlag(flagH, halfCarry)
	cpu.Registers.setFlag(flagC, carryFlag)
}

// SET u3, [HL]
func (cpu *CPU) setu3memhl(u3 uint8) {
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	res := value | 1<<u3
	cpu.memoryWrite(hl, res)
}

// RL r8
func (cpu *CPU) rlr8(register *uint8) {
	carry := uint8(0)
	if cpu.Registers.getFlag(flagC) {
		carry = uint8(1)
	}
	value := *register
	newCarry := (value & 0b10000000) >> 7
	res := (value << 1) | carry
	cpu.Registers.setFlag(flagZ, res == 0x00)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, newCarry == 0x01)
	*register = res

}

// ADC A, r8
func (cpu *CPU) adcar8(register uint8) {
	carry := uint8(0)
	if cpu.Registers.getFlag(flagC) {
		carry = uint8(1)
	}
	a := cpu.Registers.A
	res := a + carry + register
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	// Set if overflow from bit 3
	halfCarry := ((a & 0b1111) + (register & 0b1111) + (carry & 0b1111)) > 0b1111
	cpu.Registers.setFlag(flagH, halfCarry)
	// Set if overflow from bit 7
	newCarry := (uint16(a) + uint16(register) + uint16(carry)) > 0b11111111
	cpu.Registers.setFlag(flagC, newCarry)
	cpu.Registers.A = res
}

// RLA
func (cpu *CPU) rla() {
	carry := uint8(0)
	if cpu.Registers.getFlag(flagC) {
		carry = uint8(1)
	}
	a := cpu.Registers.A
	newCarry := (a & 0b10000000) >> 7
	res := (a << 1) | carry
	cpu.Registers.setFlag(flagZ, false)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, newCarry == 0x01)
	cpu.Registers.A = res
}

// SUB A, [HL]
func (cpu *CPU) subamemhl() {
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	a := cpu.Registers.A
	res := a - value
	// if borrow from bit 4
	//halfCarry := (cpu.Registers.A & 0b10000) != uint8(res&0b10000)
	halfCarry := (a & 0b1111) < (value & 0b1111)
	// if borrow
	carryFlag := uint16(value) > uint16(a)
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, true)
	cpu.Registers.setFlag(flagH, halfCarry)
	cpu.Registers.setFlag(flagC, carryFlag)
	cpu.Registers.A = res
}

// SBC A, [HL]
func (cpu *CPU) sbcamemhl() {
	carry := uint8(0)
	if cpu.Registers.getFlag(flagC) {
		carry = uint8(1)
	}
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	n := value + carry
	a := cpu.Registers.A
	res := a - n
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, true)
	// Set if borrow from bit 4
	//halfCarry := (a & 0b10000) != (res & 0b10000)
	halfCarry := (a & 0b1111) < (value&0b1111)+carry

	cpu.Registers.setFlag(flagH, halfCarry)
	// Set if borrow (i.e. if (r8 + carry) > A)
	newCarry := uint16(value)+uint16(carry) > uint16(a)
	cpu.Registers.setFlag(flagC, newCarry)
	cpu.Registers.A = res
}

// ////////////////////////
// RR r8
func (cpu *CPU) rrr8(register *uint8) {
	value := *register
	//res := value >> 1
	carry := uint8(0)
	if cpu.Registers.getFlag(flagC) {
		carry = uint8(0b10000000)
	}
	//if cpu.Registers.getFlag(flagC) {
	//	res = uint8(0b10000000) | (value >> 1)
	//}
	newCarry := value & 0b00000001
	res := (value >> 1) | carry
	*register = res

	cpu.Registers.setFlag(flagZ, res == 0x00)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, newCarry != 0x00)
}

// SBC A, n8
func (cpu *CPU) sbcan8() {
	carry := uint8(0)
	if cpu.Registers.getFlag(flagC) {
		carry = uint8(1)
	}
	n8 := cpu.getImmediate8()
	n := n8 + carry
	a := cpu.Registers.A
	res := a - n
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, true)
	// Set if borrow from bit 4
	//halfCarry := (a & 0b10000) != (res & 0b10000)
	halfCarry := (a & 0b1111) < (n8&0b1111)+carry

	cpu.Registers.setFlag(flagH, halfCarry)
	// Set if borrow (i.e. if (r8 + carry) > A)
	newCarry := uint16(n8)+uint16(carry) > uint16(a)

	//halfCarry := (a & 0b1111) < (n & 0b1111)
	//cpu.Registers.setFlag(flagH, halfCarry)
	//// Set if borrow (i.e. if (r8 + carry) > A)
	//newCarry := uint16(n) > uint16(a)
	cpu.Registers.setFlag(flagC, newCarry)
	cpu.Registers.A = res
}

// SRA r8
func (cpu *CPU) srar8(register *uint8) {
	val := *register
	bit7 := val & 0b10000000
	res := val >> 1
	res |= bit7
	// bit 7 of register is unchanged
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	// set when a rotate/shift operation shifts out a “1” bit
	carry := val & 0b00000001
	cpu.Registers.setFlag(flagC, carry != 0x00)
	*register = res
}

// ADD A, [HL]
func (cpu *CPU) addamemhl() {
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	a := cpu.Registers.A
	res := a + value
	cpu.Registers.A = res

	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	// Set if overflow from bit 3
	halfCarry := ((a & 0b1111) + (value & 0b1111)) > 0b1111
	//halfCarry := res > 0b1111
	cpu.Registers.setFlag(flagH, halfCarry)
	// Set if overflow from bit 7
	carry := (uint16(a) + uint16(value)) > 0b11111111
	cpu.Registers.setFlag(flagC, carry)
}

// ADC A, [HL]
func (cpu *CPU) adcamemhl() {
	carry := uint8(0)
	if cpu.Registers.getFlag(flagC) {
		carry = uint8(1)
	}
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	a := cpu.Registers.A
	res := a + carry + value

	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	// Set if overflow from bit 3
	halfCarry := ((a & 0b1111) + (value & 0b1111) + (carry & 0b1111)) > 0b1111
	cpu.Registers.setFlag(flagH, halfCarry)
	// Set if overflow from bit 7
	newCarry := (uint16(a) + uint16(value) + uint16(carry)) > 0b11111111
	cpu.Registers.setFlag(flagC, newCarry)
	cpu.Registers.A = res
}

// CCF
func (cpu *CPU) ccf() {
	carry := true
	if cpu.Registers.getFlag(flagC) {
		carry = false
	}
	cpu.Registers.setFlag(flagC, carry)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
}

// RLC r8
func (cpu *CPU) rlcr8(register *uint8) {
	value := *register
	carry := (value & 0b10000000) >> 7
	res := (value << 1) | carry
	*register = res
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, carry == 0x01)
}

// RLC [HL]
func (cpu *CPU) rlcmemhl() {
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	carry := (value & 0b10000000) >> 7
	res := (value << 1) | carry
	cpu.memoryWrite(hl, res)
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, carry == 0x01)
}

// ////////////////////////////
// RRC r8
func (cpu *CPU) rrcr8(register *uint8) {
	value := *register
	carry := value & 0b00000001
	res := (value >> 1) | (carry << 7)
	//res := value >> 1
	//if carry == 0x01 {
	//	res = 0b10000000 | (value >> 1)
	//}
	*register = res
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, carry != 0x00)
}

// RRC [HL]
func (cpu *CPU) rrcmemhl() {
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	carry := value & 0b00000001
	//res := (value >> 1) | (carry << 7)
	res := value >> 1
	if carry == 0x01 {
		res = 0b10000000 | (value >> 1)
	}
	cpu.memoryWrite(hl, res)
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, carry != 0x00)
}

// RR [HL]
func (cpu *CPU) rrmemhl() {
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	carry := uint8(0)
	if cpu.Registers.getFlag(flagC) {
		carry = uint8(0b10000000)
	}
	//res := value >> 1
	//if cpu.Registers.getFlag(flagC) {
	//	res = uint8(0b10000000) | (value >> 1)
	//}
	res := (value >> 1) | carry
	newCarry := value & 0b00000001

	//res := carry | (value >> 1)
	cpu.Registers.setFlag(flagZ, res == 0x00)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, newCarry != 0x00)
	cpu.memoryWrite(hl, res)
}

// SRL [HL]
func (cpu *CPU) srlmemhl() {
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	carry := value & 0b00000001
	res := value >> 1
	// bit 7 of register is reset to 0 ?
	//res &= 0b01111111
	cpu.memoryWrite(hl, res)
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, carry != 0x00)
}

// SRA [HL]
func (cpu *CPU) sramemhl() {
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	bit7 := value & 0b10000000
	res := (value >> 1) | bit7
	carry := value & 0b00000001

	// bit 7 of register is unchanged
	cpu.memoryWrite(hl, res)

	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	// set when a rotate/shift operation shifts out a “1” bit
	cpu.Registers.setFlag(flagC, carry != 0x00)
}

// SLA [HL]
func (cpu *CPU) slamemhl() {
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	// set when a rotate/shift operation shifts out a “1” bit
	carry := value & 0b10000000
	//carry := (value & 0b10000000) >> 7

	res := value << 1
	//// bit 0 of register is reset to 0
	//res &= 0b11111110

	cpu.memoryWrite(hl, res)

	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)

	cpu.Registers.setFlag(flagC, carry != 0x00)
}

// SWAP [HL]
func (cpu *CPU) swapmemhl() {
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	//swap := (value>>4)&0b00001111 | (value<<4)&0b11110000
	swap := (value >> 4) | (value << 4)

	cpu.memoryWrite(hl, swap)
	cpu.Registers.setFlag(flagZ, value == 0x00)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, false)
}

// RL [HL]
func (cpu *CPU) rlmemhl() {
	carry := uint8(0)
	if cpu.Registers.getFlag(flagC) {
		carry = uint8(1)
	}
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	//newCarry := value & 0b10000000
	newCarry := (value & 0b10000000) >> 7
	res := (value << 1) | carry

	cpu.memoryWrite(hl, res)
	cpu.Registers.setFlag(flagZ, res == 0x00)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, newCarry == 0x01)

}

// RRA
func (cpu *CPU) rra() {
	a := cpu.Registers.A
	carry := uint8(0)
	if cpu.Registers.getFlag(flagC) {
		carry = uint8(0b10000000)
	}
	res := (a >> 1) | carry
	//res := a >> 1
	//if cpu.Registers.getFlag(flagC) {
	//	res = (a >> 1) | (0b10000000)
	//}
	newCarry := a & 0b00000001
	//res := carry | (value >> 1)
	cpu.Registers.A = res

	cpu.Registers.setFlag(flagZ, false)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, newCarry != 0x00)
}

/////////////////

// LD [n16], SP TODO
func (cpu *CPU) ldmemn16sp() {
	sp := cpu.Registers.SP
	n16 := cpu.getImmediate16()
	addr := cpu.memoryRead(n16)
	//addr2 := cpu.memoryRead(n16 + 1)
	value := sp & 0xFF
	value2 := sp >> 8
	cpu.memoryWrite(uint16(addr), uint8(value))
	cpu.memoryWrite(uint16(addr+1), uint8(value2))
}

// LD HL, SP+e8
func (cpu *CPU) ldhlspe8() {
	e8 := int(cpu.getImmediate8())
	sp := cpu.Registers.SP
	value := uint16(int32(e8) + int32(sp))
	cpu.Registers.setHL(value)
	cpu.Registers.setFlag(flagZ, false)
	cpu.Registers.setFlag(flagN, false)
	halfCarry := ((sp & 0b1111) + (uint16(e8) & 0b1111)) > 0b1111
	cpu.Registers.setFlag(flagH, halfCarry)
	carry := ((sp & 0b11111111) + (uint16(e8) & 0b11111111)) > 0b11111111
	cpu.Registers.setFlag(flagC, carry)
}

func (cpu *CPU) hldec() {
	hl := cpu.Registers.getHL()
	cpu.Registers.setHL(hl - 1)
}

func (cpu *CPU) hlinc() {
	hl := cpu.Registers.getHL()
	cpu.Registers.setHL(hl + 1)
}
