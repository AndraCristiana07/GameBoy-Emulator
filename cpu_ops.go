package main

import "fmt"

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
	cpu.xorA(cpu.Registers.A)
	return 4
}

// LD r16,n16 - r16=HL
func (cpu *CPU) op21() int {
	cpu.ldr16n16(cpu.Registers.setHL)
	return 12
}

// LD r8, n8 - r8=C
func (cpu *CPU) op0e() int {
	cpu.ldr8n8(&cpu.Registers.C)
	return 8
}

// LD r8, n8 - r8=B
func (cpu *CPU) op06() int {
	cpu.ldr8n8(&cpu.Registers.B)
	return 8
}

// LD [HLD],A
func (cpu *CPU) op32() int {
	cpu.ldmemhlda()
	return 8
}

// DEC r8 - r8=B
func (cpu *CPU) op05() int {
	cpu.decr8(&cpu.Registers.B)
	return 4
}

// JR cc, e8 - cc=NZ
func (cpu *CPU) op20() int {
	return cpu.jrCCe8(!cpu.Registers.getFlag(flagZ))
}

///////// review 1 end

// DEC r8 - r8=C
func (cpu *CPU) op0d() int {
	cpu.decr8(&cpu.Registers.C)
	return 4
}

// LD r8, n8 - r8=A
func (cpu *CPU) op3e() int {
	cpu.ldr8n8(&cpu.Registers.A)
	return 8
}

// DI
func (cpu *CPU) opf3() int {
	cpu.IME = false
	return 4
}

// LDH [a8], r8 - r8=A
func (cpu *CPU) ope0() int {
	cpu.ldhmema8r8(cpu.Registers.A)
	return 12
}

// LDH r8, [a8] - r8=A
func (cpu *CPU) opf0() int {
	cpu.ldhr8mema8(&cpu.Registers.A)
	return 12
}

// CP A, n8
func (cpu *CPU) opfe() int {
	n8 := cpu.getImmediate8()
	cpu.cpa(n8)
	return 8
}

// skip review pentru ca nu e folosit
// SUB A, r8 - r8=H
func (cpu *CPU) op94() int {
	cpu.subar8(cpu.Registers.H)
	return 4
}

// LD [HL], n8
func (cpu *CPU) op36() int {
	n8 := cpu.getImmediate8()
	cpu.ldmemhl(n8)
	return 12
}

// LD [a16], A
func (cpu *CPU) opea() int {
	cpu.ldmema16A()
	return 16
}

// LD SP, n16
func (cpu *CPU) op31() int {
	cpu.ldspn16()
	return 12
}

// skip review pentru ca nu e folosit
// RST $38
func (cpu *CPU) opff() int {
	cpu.call(0x38)
	return 16
}

// LD A, [HLI]
func (cpu *CPU) op2a() int {
	cpu.ldamemhli()
	return 8
}

// LDH [C], A
func (cpu *CPU) ope2() int {
	cpu.ldhmemca()
	return 8
}

// INC r8 - r8=C
func (cpu *CPU) op0c() int {
	cpu.incr8(&cpu.Registers.C)
	return 4
}

// CALL a16
func (cpu *CPU) opcd() int {
	cpu.calla16()
	return 24
}

// LD r16, n16 - r16=BC
func (cpu *CPU) op01() int {
	cpu.ldr16n16(cpu.Registers.setBC)
	return 12
}

// DEC R16 - r16=BC
func (cpu *CPU) op0b() int {
	cpu.decr16(cpu.Registers.getBC, cpu.Registers.setBC)
	return 8
}

// LD r8, r8 - r8=A,B
func (cpu *CPU) op78() int {
	cpu.ldr8r8(&cpu.Registers.A, cpu.Registers.B)
	return 4
}

// OR A, r8 - r8=C
func (cpu *CPU) opb1() int {
	cpu.orA(cpu.Registers.C)
	return 4
}

// RET
func (cpu *CPU) opc9() int {
	cpu.ret()
	return 16
}

// EI
func (cpu *CPU) opfb() int {
	cpu.ei()
	return 4
}

// PUSH r16 - r16=AF
func (cpu *CPU) opf5() int {
	cpu.push(cpu.Registers.getAF())
	return 16
}

// PUSH r16 - r16=BC
func (cpu *CPU) opc5() int {
	cpu.push(cpu.Registers.getBC())
	return 16
}

// PUSH r16 - r16=DE
func (cpu *CPU) opd5() int {
	cpu.push(cpu.Registers.getDE())
	return 16
}

// PUSH r16 - r16=HL
func (cpu *CPU) ope5() int {
	cpu.push(cpu.Registers.getHL())
	return 16
}

// AND A, r8 - r8=A
func (cpu *CPU) opa7() int {
	cpu.andA(cpu.Registers.A)
	return 4
}

// JR CC, e8 - CC=Z
func (cpu *CPU) op28() int {
	return cpu.jrCCe8(cpu.Registers.getFlag(flagZ))
}

// RET CC - cc=NZ
func (cpu *CPU) opc0() int {
	return cpu.retcc(!cpu.Registers.getFlag(flagZ))
}

// LD r8, [a16] - r8=A
func (cpu *CPU) opfa() int {
	cpu.ldr8mema16(&cpu.Registers.A)
	return 16
}

// RET CC - cc=Z
func (cpu *CPU) opc8() int {
	//fmt.Println("RET CC")
	return cpu.retcc(cpu.Registers.getFlag(flagZ))
}

// skip review pentru ca nu e folosit
// SCF
func (cpu *CPU) op37() int {
	cpu.scf()
	return 4
}

// skip review pentru ca nu e folosit
// INC r8 - r8=E
func (cpu *CPU) op1c() int {
	cpu.incr8(&cpu.Registers.E)
	return 4
}

// skip review pentru ca nu e folosit
// INC r8 - r8=H
func (cpu *CPU) op24() int {
	cpu.incr8(&cpu.Registers.H)
	return 4
}

// skip review pentru ca nu e folosit
// ADD HL, r16 - r16=BC
func (cpu *CPU) op09() int {
	cpu.addhlr16(cpu.Registers.getBC)
	return 8
}

///////// review 2 end

// DEC r8 - r8=A
func (cpu *CPU) op3d() int {
	cpu.decr8(&cpu.Registers.A)
	return 4
}

// INC [HL]
func (cpu *CPU) op34() int {
	cpu.incmemhl()
	return 12
}

// INC r8 - r8=A
func (cpu *CPU) op3c() int {
	cpu.incr8(&cpu.Registers.A)
	return 4
}

// POP HL
func (cpu *CPU) ope1() int {
	cpu.Registers.setHL(cpu.pop())
	return 12
}

// POP r16 - r16=DE
func (cpu *CPU) opd1() int {
	cpu.Registers.setDE(cpu.pop())
	return 12
}

// POP r16 - r16=BC
func (cpu *CPU) opc1() int {
	cpu.Registers.setBC(cpu.pop())
	return 12
}

// POP r16 - r16=AF
func (cpu *CPU) opf1() int {
	cpu.Registers.setAF(cpu.pop())
	return 12
}

// RETI
func (cpu *CPU) opd9() int {
	cpu.reti()
	return 16
}

// CPL
func (cpu *CPU) op2f() int {
	cpu.cpl()
	return 4
}

// AND A, e8
func (cpu *CPU) ope6() int {
	e8 := cpu.getImmediate8()
	cpulogger.Debug(fmt.Sprintf("Imm in and A e8 %04X", e8))
	cpu.andA(e8)
	return 8
}

// LD r8, r8 - r8=B, A
func (cpu *CPU) op47() int {
	cpu.ldr8r8(&cpu.Registers.B, cpu.Registers.A)
	return 4
}

// OR A r8 - r8=B
func (cpu *CPU) opb0() int {
	cpu.orA(cpu.Registers.B)
	return 4
}

// LD r8, r8 - r8=C, A
func (cpu *CPU) op4f() int {
	cpu.ldr8r8(&cpu.Registers.C, cpu.Registers.A)
	return 4
}

// XOR A, r8 - r8=C
func (cpu *CPU) opa9() int {
	cpu.xorA(cpu.Registers.C)
	return 4
}

// AND A, r8 - r8=C
func (cpu *CPU) opa1() int {
	cpu.andA(cpu.Registers.C)
	return 4
}

// LD r8, r8 - r8=A, C
func (cpu *CPU) op79() int {
	cpu.ldr8r8(&cpu.Registers.A, cpu.Registers.C)
	return 4
}

// RST $28
func (cpu *CPU) opef() int {
	cpu.call(0x28)
	return 16
}

// ADD A, r8 - r8=A
func (cpu *CPU) op87() int {
	cpu.addar8(cpu.Registers.A)
	return 4
}

// LD r8,r8 - r8=E,A
func (cpu *CPU) op5f() int {
	cpu.ldr8r8(&cpu.Registers.E, cpu.Registers.A)
	return 4
}

// LD r8, n8 - r8=D
func (cpu *CPU) op16() int {
	cpu.ldr8n8(&cpu.Registers.D)
	return 8
}

// ADD HL, r16 - r16=DE
func (cpu *CPU) op19() int {
	cpu.addhlr16(cpu.Registers.getDE)
	return 8
}

// LD r8, [HL] - r8=E
func (cpu *CPU) op5e() int {
	cpu.ldr8memhl(&cpu.Registers.E)
	return 8
}

// INC HL
func (cpu *CPU) op23() int {
	cpu.inchl()
	return 8
}

// LD r8. [HL] r8=D
func (cpu *CPU) op56() int {
	cpu.ldr8memhl(&cpu.Registers.D)
	return 8
}

// JP HL
func (cpu *CPU) ope9() int {
	hl := cpu.Registers.getHL()
	cpu.jump(hl)
	return 4
}

// LD r16, n16 - r16=DE
func (cpu *CPU) op11() int {
	cpu.ldr16n16(cpu.Registers.setDE)
	return 12
}

// LD [r16], A - a16=DE
func (cpu *CPU) op12() int {
	cpu.ldmemr16A(cpu.Registers.getDE)
	return 8
}

// INC r16 - r16=DE
func (cpu *CPU) op13() int {
	cpulogger.Debug(fmt.Sprintf("Inc DE %04X", cpu.Registers.getDE()))
	cpu.incr16(cpu.Registers.getDE, cpu.Registers.setDE)
	cpulogger.Debug(fmt.Sprintf("Inc DE after %04X", cpu.Registers.getDE()))
	return 8
}

// LD A,[r16] - r16=DE
func (cpu *CPU) op1a() int {
	cpu.ldamemr16(cpu.Registers.getDE)
	return 8
}

// LD [HLI], A
func (cpu *CPU) op22() int {
	cpu.ldmemhlia()
	return 8
}

// LD r8,r8 - r8=A,H
func (cpu *CPU) op7c() int {
	cpu.ldr8r8(&cpu.Registers.A, cpu.Registers.H)
	return 4
}

// JP cc, a16 - cc=Z
func (cpu *CPU) opca() int {
	return cpu.jpCCa16(cpu.Registers.getFlag(flagZ))
}

// LD A, [r16] - r16=HL
func (cpu *CPU) op7e() int {
	cpu.ldamemr16(cpu.Registers.getHL)
	return 8
}

// JR e8
func (cpu *CPU) op18() int {
	cpu.jre8()
	return 12
}

// DEC r8 - r8=L
func (cpu *CPU) op2d() int {
	cpu.decr8(&cpu.Registers.L)
	return 4
}

// LD  A, [HLD]
func (cpu *CPU) op3a() int {
	cpu.ldamemhld()
	return 8
}

// LD r8, r8 - r8=D,A
func (cpu *CPU) op57() int {
	cpu.ldr8r8(&cpu.Registers.D, cpu.Registers.A)
	return 4
}

// LD r8, r8 - r8=A,E
func (cpu *CPU) op7b() int {
	cpu.ldr8r8(&cpu.Registers.A, cpu.Registers.E)
	return 4
}

// LD r8, r8 - r8=A,D
func (cpu *CPU) op7a() int {
	cpu.ldr8r8(&cpu.Registers.A, cpu.Registers.D)
	return 4
}

// LD A, r16 - r16=[BC]
func (cpu *CPU) op0a() int {
	cpu.ldamemr16(cpu.Registers.getBC)
	return 8
}

// LD r8, r8 - r8=A,L
func (cpu *CPU) op7d() int {
	cpu.ldr8r8(&cpu.Registers.A, cpu.Registers.L)
	return 4
}

// ADD A, n8
func (cpu *CPU) opc6() int {
	cpu.addan8()
	return 8
}

// LD r8, r8 - r8=L,A
func (cpu *CPU) op6f() int {
	cpu.ldr8r8(&cpu.Registers.L, cpu.Registers.A)
	return 4
}

// LD r8, r8 - r8=E,L
func (cpu *CPU) op5d() int {
	cpu.ldr8r8(&cpu.Registers.E, cpu.Registers.L)
	return 4
}

// LD r8, r8 - r8=D,H
func (cpu *CPU) op54() int {
	cpu.ldr8r8(&cpu.Registers.D, cpu.Registers.H)
	return 4
}

// INC r8 - r8=L
func (cpu *CPU) op2c() int {
	cpu.incr8(&cpu.Registers.L)
	return 4
}

// OR A, n8
func (cpu *CPU) opf6() int {
	n8 := cpu.getImmediate8()
	cpu.orA(n8)
	return 8
}

// DEC [HL]
func (cpu *CPU) op35() int {
	cpu.decmemhl()
	return 12
}

// JR cc, e8 - cc=NC
func (cpu *CPU) op30() int {
	return cpu.jrCCe8(!cpu.Registers.getFlag(flagC))
}

// LD r8, r8 - r8=L,E
func (cpu *CPU) op6b() int {
	cpu.ldr8r8(&cpu.Registers.L, cpu.Registers.E)
	return 4
}

// LD [r16], A r16=BC
func (cpu *CPU) op02() int {
	cpu.ldmemr16A(cpu.Registers.getBC)
	return 8
}

// LD [r16], A r16=HL
func (cpu *CPU) op77() int {
	cpu.ldmemr16A(cpu.Registers.getHL)
	return 8
}

// INC r16 - r16=BC
func (cpu *CPU) op03() int {
	cpu.incr16(cpu.Registers.getBC, cpu.Registers.setBC)
	return 8
}

// SBC A, r8 - r8=E
func (cpu *CPU) op9b() int {
	cpu.sbcar8(cpu.Registers.E)
	return 4
}

// JP cc, a16 cc=C
func (cpu *CPU) opda() int {
	return cpu.jpCCa16(cpu.Registers.getFlag(flagC))
}

// RLCA
func (cpu *CPU) op07() int {
	cpu.rlca()
	return 4
}

// LD r8, r8 - r8=H,A
func (cpu *CPU) op67() int {
	cpu.ldr8r8(&cpu.Registers.H, cpu.Registers.A)
	return 4
}

// RST $18
func (cpu *CPU) opdf() int {
	cpu.call(0x18)
	return 16
}

// ADD HL, SP
func (cpu *CPU) op39() int {
	cpu.addhlsp()
	return 8
}

// LD L, H
func (cpu *CPU) op6c() int {
	cpu.ldr8r8(&cpu.Registers.L, cpu.Registers.H)
	return 4
}

// LD SP, HL
func (cpu *CPU) opf9() int {
	cpu.ldsphl()
	return 8
}

// RST $08
func (cpu *CPU) opcf() int {
	cpu.call(0x08)
	return 16
}

// LD r8, [HL] - r8=C
func (cpu *CPU) op4e() int {
	cpu.ldr8memhl(&cpu.Registers.C)
	return 8
}

// LD r8, [HL] - r8=B
func (cpu *CPU) op46() int {
	cpu.ldr8memhl(&cpu.Registers.B)
	return 8
}

// LD r8, r8 - r8=L,C
func (cpu *CPU) op69() int {
	cpu.ldr8r8(&cpu.Registers.L, cpu.Registers.C)
	return 4
}

// LD r8, r8 - r8=H,B
func (cpu *CPU) op60() int {
	cpu.ldr8r8(&cpu.Registers.H, cpu.Registers.B)
	return 4
}

// ADD A, r8 - r8=L
func (cpu *CPU) op85() int {
	cpu.addar8(cpu.Registers.L)
	return 4
}

// JP NZ, a16
func (cpu *CPU) opc2() int {
	return cpu.jpCCa16(!cpu.Registers.getFlag(flagZ))
}

// LD [HL], r8 - r8=E
func (cpu *CPU) op73() int {
	cpu.ldmemhl(cpu.Registers.E)
	return 8
}

// LD [HL], r8 - r8=D
func (cpu *CPU) op72() int {
	cpu.ldmemhl(cpu.Registers.D)
	return 8
}

// LD [HL], r8 - r8=C
func (cpu *CPU) op71() int {
	cpu.ldmemhl(cpu.Registers.C)
	return 8
}

// LD r8, n8 - r8=E
func (cpu *CPU) op1e() int {
	cpu.ldr8n8(&cpu.Registers.E)
	return 8
}

// LD r8, r8 - r8=H,D
func (cpu *CPU) op62() int {
	cpu.ldr8r8(&cpu.Registers.H, cpu.Registers.D)
	return 4
}

// LD r8, r8 - r8=B,B
func (cpu *CPU) op40() int {
	cpu.ldr8r8(&cpu.Registers.B, cpu.Registers.B)
	return 4
}

// after tetris

// INC r8 - r8=B
func (cpu *CPU) op04() int {
	cpu.incr8(&cpu.Registers.B)
	return 4
}

// INC r8 - r8=D
func (cpu *CPU) op14() int {
	cpu.incr8(&cpu.Registers.D)
	return 4
}

// INC SP
func (cpu *CPU) op33() int {
	cpu.incsp()
	return 8
}

// DEC r8 - r8=D
func (cpu *CPU) op15() int {
	cpu.decr8(&cpu.Registers.D)
	return 4
}

// DEC r16 - r16=DE
func (cpu *CPU) op1b() int {
	cpu.decr16(cpu.Registers.getDE, cpu.Registers.setDE)
	return 8
}

// DEC r8 - r8=E
func (cpu *CPU) op1d() int {
	cpu.decr8(&cpu.Registers.E)
	return 4
}

// DEC r8 - r8=H
func (cpu *CPU) op25() int {
	cpu.decr8(&cpu.Registers.H)
	return 4
}

// DEC r16 - r16=HL
func (cpu *CPU) op2b() int {
	cpu.decr16(cpu.Registers.getHL, cpu.Registers.setHL)
	return 8
}

// DEC SP
func (cpu *CPU) op3b() int {
	cpu.decsp()
	return 8
}

// ADD HL, r16 - r16=HL
func (cpu *CPU) op29() int {
	cpu.addhlr16(cpu.Registers.getHL)
	return 8
}

// ADD A, r8 - r8=B
func (cpu *CPU) op80() int {
	cpu.addar8(cpu.Registers.B)
	return 4
}

// ADD A, r8 - r8=C
func (cpu *CPU) op81() int {
	cpu.addar8(cpu.Registers.C)
	return 4
}

// ADD A, r8 - r8=D
func (cpu *CPU) op82() int {
	cpu.addar8(cpu.Registers.D)
	return 4
}

// ADD A, r8 - r8=E
func (cpu *CPU) op83() int {
	cpu.addar8(cpu.Registers.E)
	return 4
}

// ADD A, r8 - r8=H
func (cpu *CPU) op84() int {
	cpu.addar8(cpu.Registers.H)
	return 4
}

// ADD SP, e8
func (cpu *CPU) ope8() int {
	cpu.addspe8()
	return 16
}

// LD r8, n8 - r8=H
func (cpu *CPU) op26() int {
	cpu.ldr8n8(&cpu.Registers.H)
	return 8
}

// SUB A, r8 - r8=L
func (cpu *CPU) op95() int {
	cpu.subar8(cpu.Registers.L)
	return 4
}

// LD r8, [HL] - r8=H
func (cpu *CPU) op66() int {
	cpu.ldr8memhl(&cpu.Registers.H)
	return 8
}

// HALT
func (cpu *CPU) op76() int {
	cpu.halt()
	return 4
}

// RRCA
func (cpu *CPU) op0f() int {
	cpu.rrca()
	return 4
}

// OR A, r8 - r8=D
func (cpu *CPU) opb2() int {
	cpu.orA(cpu.Registers.D)
	return 4
}

// JP CC, a16 - cc=NC
func (cpu *CPU) opd2() int {
	return cpu.jpCCa16(!cpu.Registers.getFlag(flagC))
}

// OR A, r8 - r8=E
func (cpu *CPU) opb3() int {
	cpu.orA(cpu.Registers.E)
	return 4
}

// AND A, r8 - r8=B
func (cpu *CPU) opa0() int {
	cpu.andA(cpu.Registers.B)
	return 4
}

// LD r8, n8 - r8=L
func (cpu *CPU) op2e() int {
	cpu.ldr8n8(&cpu.Registers.L)
	return 8
}

// ADC A, n8
func (cpu *CPU) opce() int {
	cpu.adcan8()
	return 8
}

// CP A r8 - r8=B
func (cpu *CPU) opb8() int {
	cpu.cpa(cpu.Registers.B)
	return 4
}

// JR CC, e8 - cc=C
func (cpu *CPU) op38() int {
	return cpu.jrCCe8(cpu.Registers.getFlag(flagC))
}

// CALL CC, a16 - cc=Z
func (cpu *CPU) opcc() int {
	return cpu.callcca16(cpu.Registers.getFlag(flagZ))
}

// XOR A, n8
func (cpu *CPU) opee() int {
	n8 := cpu.getImmediate8()
	cpu.xorA(n8)
	return 8
}

// SUB A, n8
func (cpu *CPU) opd6() int {
	cpu.suban8()
	return 8
}

// RET CC - cc=C
func (cpu *CPU) opd8() int {
	return cpu.retcc(cpu.Registers.getFlag(flagC))
}

// CP A, [HL]
func (cpu *CPU) opbe() int {
	cpu.cpamemhl()
	return 8
}

// SUB A, r8 - r8=C
func (cpu *CPU) op91() int {
	cpu.subar8(cpu.Registers.C)
	return 4
}

// SUB A, r8 - r8=B
func (cpu *CPU) op90() int {
	cpu.subar8(cpu.Registers.B)
	return 4
}

// RET CC - cc=NC
func (cpu *CPU) opd0() int {
	return cpu.retcc(!cpu.Registers.getFlag(flagC))
}

// ADC A, r8 - r8=C
func (cpu *CPU) op89() int {
	cpu.adcar8(cpu.Registers.C)
	return 4
}

// LD r8, r8 - r8=H,C
func (cpu *CPU) op61() int {
	cpu.ldr8r8(&cpu.Registers.H, cpu.Registers.C)
	return 4
}

// RLA
func (cpu *CPU) op17() int {
	cpu.rla()
	return 4
}

// CP A, r8 - r8=C
func (cpu *CPU) opb9() int {
	cpu.cpa(cpu.Registers.C)
	return 4
}

// SUB A, r8 - r8=E
func (cpu *CPU) op93() int {
	cpu.subar8(cpu.Registers.E)
	return 4
}

// XOR A, [HL]
func (cpu *CPU) opae() int {
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	cpu.xorA(value)
	return 8
}

// SBC A, r8 - r8=D
func (cpu *CPU) op9a() int {
	cpu.sbcar8(cpu.Registers.D)
	return 4
}

// STOP n8
func (cpu *CPU) op10() int {
	cpu.stop()
	return 4
}

// SUB A, [HL]
func (cpu *CPU) op96() int {
	cpu.subamemhl()
	return 8
}

// SBC A, [HL]
func (cpu *CPU) op9e() int {
	cpu.sbcamemhl()
	return 8
}

// ADC A, r8 - r8=H
func (cpu *CPU) op8c() int {
	cpu.adcar8(cpu.Registers.H)
	return 4
}

// SBC A, n8
func (cpu *CPU) opde() int {
	cpu.sbcan8()
	return 8
}

// SBC A, r8 - r8=C
func (cpu *CPU) op99() int {
	cpu.sbcar8(cpu.Registers.C)
	return 4
}

// SBC A, r8 - r8=B
func (cpu *CPU) op98() int {
	cpu.sbcar8(cpu.Registers.B)
	return 4
}

// ADD A, [HL]
func (cpu *CPU) op86() int {
	cpu.addamemhl()
	return 8
}

// ADC A [HL]
func (cpu *CPU) op8e() int {
	cpu.adcamemhl()
	return 8
}

// SBC A, r8 - r8=H
func (cpu *CPU) op9c() int {
	cpu.sbcar8(cpu.Registers.H)
	return 4
}

// CP A, r8 - r8=L
func (cpu *CPU) opbd() int {
	cpu.cpa(cpu.Registers.L)
	return 4
}

// CCF
func (cpu *CPU) op3f() int {
	cpu.ccf()
	return 4
}

// LD r8, r8 - r8=D,C
func (cpu *CPU) op51() int {
	cpu.ldr8r8(&cpu.Registers.D, cpu.Registers.C)
	return 4
}

// SUB A, r8 - r8=D
func (cpu *CPU) op92() int {
	cpu.subar8(cpu.Registers.D)
	return 4
}

// SUB A, r8 - r8=A
func (cpu *CPU) op97() int {
	cpu.subar8(cpu.Registers.A)
	return 4
}

// AND A, r8 - r8=D
func (cpu *CPU) opa2() int {
	cpu.andA(cpu.Registers.D)
	return 4
}

// AND A, r8 - r8=E
func (cpu *CPU) opa3() int {
	cpu.andA(cpu.Registers.E)
	return 4
}

// AND A, r8 - r8=H
func (cpu *CPU) opa4() int {
	cpu.andA(cpu.Registers.H)
	return 4
}

// AND A, r8 - r8=L
func (cpu *CPU) opa5() int {
	cpu.andA(cpu.Registers.L)
	return 4
}

// AND A, [HL]
func (cpu *CPU) opa6() int {
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	cpu.andA(value)
	return 8
}

// OR A, r8 - r8=H
func (cpu *CPU) opb4() int {
	cpu.orA(cpu.Registers.H)
	return 4
}

// OR A, r8 - r8=L
func (cpu *CPU) opb5() int {
	cpu.orA(cpu.Registers.L)
	return 4
}

// OR A, [HL]
func (cpu *CPU) opb6() int {
	hl := cpu.Registers.getHL()
	value := cpu.memoryRead(hl)
	cpu.orA(value)
	return 8
}

// OR A, r8 - r8=A
func (cpu *CPU) opb7() int {
	cpu.orA(cpu.Registers.A)
	return 4
}

// XOR A, r8 - r8=B
func (cpu *CPU) opa8() int {
	cpu.xorA(cpu.Registers.B)
	return 4
}

// XOR A, r8 - r8=D
func (cpu *CPU) opaa() int {
	cpu.xorA(cpu.Registers.D)
	return 4
}

// XOR A, r8 - r8=E
func (cpu *CPU) opab() int {
	cpu.xorA(cpu.Registers.E)
	return 4
}

// XOR A, r8 - r8=H
func (cpu *CPU) opac() int {
	cpu.xorA(cpu.Registers.H)
	return 4
}

// XOR A, r8 - r8=L
func (cpu *CPU) opad() int {
	cpu.xorA(cpu.Registers.L)
	return 4
}

// ADC A, r8 - r8=B
func (cpu *CPU) op88() int {
	cpu.adcar8(cpu.Registers.B)
	return 4
}

// ADC A, r8 - r8=D
func (cpu *CPU) op8a() int {
	cpu.adcar8(cpu.Registers.D)
	return 4
}

// ADC A, r8 - r8=E
func (cpu *CPU) op8b() int {
	cpu.adcar8(cpu.Registers.E)
	return 4
}

// ADC A, r8 - r8=L
func (cpu *CPU) op8d() int {
	cpu.adcar8(cpu.Registers.L)
	return 4
}

// ADC A, r8 - r8=A
func (cpu *CPU) op8f() int {
	cpu.adcar8(cpu.Registers.A)
	return 4
}

// SBC A, r8 - r8=L
func (cpu *CPU) op9d() int {
	cpu.sbcar8(cpu.Registers.L)
	return 4
}

// SBC A, r8 - r8=A
func (cpu *CPU) op9f() int {
	cpu.sbcar8(cpu.Registers.A)
	return 4
}

// CP A, r8 - r8=D
func (cpu *CPU) opba() int {
	cpu.cpa(cpu.Registers.D)
	return 4
}

// CP A, r8 - r8=E
func (cpu *CPU) opbb() int {
	cpu.cpa(cpu.Registers.E)
	return 4
}

// CP A, r8 - r8=H
func (cpu *CPU) opbc() int {
	cpu.cpa(cpu.Registers.H)
	return 4
}

// CP A, r8 - r8=A
func (cpu *CPU) opbf() int {
	cpu.cpa(cpu.Registers.A)
	return 4
}

// CALL cc - cc=NZ
func (cpu *CPU) opc4() int {
	return cpu.callcca16(!cpu.Registers.getFlag(flagZ))
}

// CALL cc - cc=C
func (cpu *CPU) opdc() int {
	return cpu.callcca16(cpu.Registers.getFlag(flagC))
}

// ILLEGAL
func (cpu *CPU) opfc() int {
	return 4
}

// CALL cc - cc=NC
func (cpu *CPU) opd4() int {
	return cpu.callcca16(!cpu.Registers.getFlag(flagC))
}

// RRA
func (cpu *CPU) op1f() int {
	cpu.rra()
	return 4
}

// LD [n16], SP
func (cpu *CPU) op08() int {
	cpu.ldmemn16sp()
	return 20
}

// LD HL, SP+e8
func (cpu *CPU) opf8() int {
	cpu.ldhlspe8()
	return 12
}

// LD [HL], r8 - r8=B
func (cpu *CPU) op70() int {
	cpu.ldmemhl(cpu.Registers.B)
	return 8
}

// LD [HL], r8 - r8=H
func (cpu *CPU) op74() int {
	cpu.ldmemhl(cpu.Registers.H)
	return 8
}

// LD [HL], r8 - r8=L
func (cpu *CPU) op75() int {
	cpu.ldmemhl(cpu.Registers.L)
	return 8
}

// LD r8, [HL] - r8=L
func (cpu *CPU) op6e() int {
	cpu.ldr8memhl(&cpu.Registers.L)
	return 8
}

// LD r8,r8 - r8=A,A
func (cpu *CPU) op7f() int {
	cpu.ldr8r8(&cpu.Registers.A, cpu.Registers.A)
	return 4
}

// LD r8,r8 - r8=B,C
func (cpu *CPU) op41() int {
	cpu.ldr8r8(&cpu.Registers.B, cpu.Registers.C)
	return 4
}

// LD r8,r8 - r8=B,D
func (cpu *CPU) op42() int {
	cpu.ldr8r8(&cpu.Registers.B, cpu.Registers.D)
	return 4
}

// LD r8,r8 - r8=B,E
func (cpu *CPU) op43() int {
	cpu.ldr8r8(&cpu.Registers.B, cpu.Registers.E)
	return 4
}

// LD r8,r8 - r8=B,H
func (cpu *CPU) op44() int {
	cpu.ldr8r8(&cpu.Registers.B, cpu.Registers.H)
	return 4
}

// LD r8,r8 - r8=B,L
func (cpu *CPU) op45() int {
	cpu.ldr8r8(&cpu.Registers.B, cpu.Registers.L)
	return 4
}

// LD r8,r8 - r8=C,B
func (cpu *CPU) op48() int {
	cpu.ldr8r8(&cpu.Registers.C, cpu.Registers.B)
	return 4
}

// LD r8,r8 - r8=C,C
func (cpu *CPU) op49() int {
	cpu.ldr8r8(&cpu.Registers.C, cpu.Registers.C)
	return 4
}

// LD r8,r8 - r8=C,D
func (cpu *CPU) op4a() int {
	cpu.ldr8r8(&cpu.Registers.C, cpu.Registers.D)
	return 4
}

// LD r8,r8 - r8=C,E
func (cpu *CPU) op4b() int {
	cpu.ldr8r8(&cpu.Registers.C, cpu.Registers.E)
	return 4
}

// LD r8,r8 - r8=C,H
func (cpu *CPU) op4c() int {
	cpu.ldr8r8(&cpu.Registers.C, cpu.Registers.H)
	return 4
}

// LD r8,r8 - r8=C,L
func (cpu *CPU) op4d() int {
	cpu.ldr8r8(&cpu.Registers.C, cpu.Registers.L)
	return 4
}

// ///// generated
// RLC r8 - r8=B
func (cpu *CPU) cbop00() int {
	cpu.rlcr8(&cpu.Registers.B)
	return 8
}

// RLC r8 - r8=C
func (cpu *CPU) cbop01() int {
	cpu.rlcr8(&cpu.Registers.C)
	return 8
}

// RLC r8 - r8=D
func (cpu *CPU) cbop02() int {
	cpu.rlcr8(&cpu.Registers.D)
	return 8
}

// RLC r8 - r8=E
func (cpu *CPU) cbop03() int {
	cpu.rlcr8(&cpu.Registers.E)
	return 8
}

// RLC r8 - r8=H
func (cpu *CPU) cbop04() int {
	cpu.rlcr8(&cpu.Registers.H)
	return 8
}

// RLC r8 - r8=L
func (cpu *CPU) cbop05() int {
	cpu.rlcr8(&cpu.Registers.L)
	return 8
}

// RLC [HL]
func (cpu *CPU) cbop06() int {
	cpu.rlcmemhl()
	return 16
}

// RLC r8 - r8=A
func (cpu *CPU) cbop07() int {
	cpu.rlcr8(&cpu.Registers.A)
	return 8
}

// RRC r8 - r8=B
func (cpu *CPU) cbop08() int {
	cpu.rrcr8(&cpu.Registers.B)
	return 8
}

// RRC r8 - r8=C
func (cpu *CPU) cbop09() int {
	cpu.rrcr8(&cpu.Registers.C)
	return 8
}

// RRC r8 - r8=D
func (cpu *CPU) cbop0a() int {
	cpu.rrcr8(&cpu.Registers.D)
	return 8
}

// RRC r8 - r8=E
func (cpu *CPU) cbop0b() int {
	cpu.rrcr8(&cpu.Registers.E)
	return 8
}

// RRC r8 - r8=H
func (cpu *CPU) cbop0c() int {
	cpu.rrcr8(&cpu.Registers.H)
	return 8
}

// RRC r8 - r8=L
func (cpu *CPU) cbop0d() int {
	cpu.rrcr8(&cpu.Registers.L)
	return 8
}

// RRC [HL]
func (cpu *CPU) cbop0e() int {
	cpu.rrcmemhl()
	return 16
}

// RRC r8 - r8=A
func (cpu *CPU) cbop0f() int {
	cpu.rrcr8(&cpu.Registers.A)
	return 8
}

// RL r8 - r8=B
func (cpu *CPU) cbop10() int {
	cpu.rlr8(&cpu.Registers.B)
	return 8
}

// RL r8 - r8=C
func (cpu *CPU) cbop11() int {
	cpu.rlr8(&cpu.Registers.C)
	return 8
}

// RL r8 - r8=D
func (cpu *CPU) cbop12() int {
	cpu.rlr8(&cpu.Registers.D)
	return 8
}

// RL r8 - r8=E
func (cpu *CPU) cbop13() int {
	cpu.rlr8(&cpu.Registers.E)
	return 8
}

// RL r8 - r8=H
func (cpu *CPU) cbop14() int {
	cpu.rlr8(&cpu.Registers.H)
	return 8
}

// RL r8 - r8=L
func (cpu *CPU) cbop15() int {
	cpu.rlr8(&cpu.Registers.L)
	return 8
}

// RL [HL]
func (cpu *CPU) cbop16() int {
	cpu.rlmemhl()
	return 16
}

// RL r8 - r8=A
func (cpu *CPU) cbop17() int {
	cpu.rlr8(&cpu.Registers.A)
	return 8
}

// RR r8 - r8=B
func (cpu *CPU) cbop18() int {
	cpu.rrr8(&cpu.Registers.B)
	return 8
}

// RR r8 - r8=C
func (cpu *CPU) cbop19() int {
	cpu.rrr8(&cpu.Registers.C)
	return 8
}

// RR r8 - r8=D
func (cpu *CPU) cbop1a() int {
	cpu.rrr8(&cpu.Registers.D)
	return 8
}

// RR r8 - r8=E
func (cpu *CPU) cbop1b() int {
	cpu.rrr8(&cpu.Registers.E)
	return 8
}

// RR r8 - r8=H
func (cpu *CPU) cbop1c() int {
	cpu.rrr8(&cpu.Registers.H)
	return 8
}

// RR r8 - r8=L
func (cpu *CPU) cbop1d() int {
	cpu.rrr8(&cpu.Registers.L)
	return 8
}

// RR [HL]
func (cpu *CPU) cbop1e() int {
	cpu.rrmemhl()
	return 16
}

// RR r8 - r8=A
func (cpu *CPU) cbop1f() int {
	cpu.rrr8(&cpu.Registers.A)
	return 8
}

// SLA r8 - r8=B
func (cpu *CPU) cbop20() int {
	cpu.slar8(&cpu.Registers.B)
	return 8
}

// SLA r8 - r8=C
func (cpu *CPU) cbop21() int {
	cpu.slar8(&cpu.Registers.C)
	return 8
}

// SLA r8 - r8=D
func (cpu *CPU) cbop22() int {
	cpu.slar8(&cpu.Registers.D)
	return 8
}

// SLA r8 - r8=E
func (cpu *CPU) cbop23() int {
	cpu.slar8(&cpu.Registers.E)
	return 8
}

// SLA r8 - r8=H
func (cpu *CPU) cbop24() int {
	cpu.slar8(&cpu.Registers.H)
	return 8
}

// SLA r8 - r8=L
func (cpu *CPU) cbop25() int {
	cpu.slar8(&cpu.Registers.L)
	return 8
}

// SLA [HL]
func (cpu *CPU) cbop26() int {
	cpu.slamemhl()
	return 16
}

// SLA r8 - r8=A
func (cpu *CPU) cbop27() int {
	cpu.slar8(&cpu.Registers.A)
	return 8
}

// SRA r8 - r8=B
func (cpu *CPU) cbop28() int {
	cpu.srar8(&cpu.Registers.B)
	return 8
}

// SRA r8 - r8=C
func (cpu *CPU) cbop29() int {
	cpu.srar8(&cpu.Registers.C)
	return 8
}

// SRA r8 - r8=D
func (cpu *CPU) cbop2a() int {
	cpu.srar8(&cpu.Registers.D)
	return 8
}

// SRA r8 - r8=E
func (cpu *CPU) cbop2b() int {
	cpu.srar8(&cpu.Registers.E)
	return 8
}

// SRA r8 - r8=H
func (cpu *CPU) cbop2c() int {
	cpu.srar8(&cpu.Registers.H)
	return 8
}

// SRA r8 - r8=L
func (cpu *CPU) cbop2d() int {
	cpu.srar8(&cpu.Registers.L)
	return 8
}

// SRA [HL]
func (cpu *CPU) cbop2e() int {
	cpu.sramemhl()
	return 16
}

// SRA r8 - r8=A
func (cpu *CPU) cbop2f() int {
	cpu.srar8(&cpu.Registers.A)
	return 8
}

// SWAP r8 - r8=B
func (cpu *CPU) cbop30() int {
	cpu.swapr8(&cpu.Registers.B)
	return 8
}

// SWAP r8 - r8=C
func (cpu *CPU) cbop31() int {
	cpu.swapr8(&cpu.Registers.C)
	return 8
}

// SWAP r8 - r8=D
func (cpu *CPU) cbop32() int {
	cpu.swapr8(&cpu.Registers.D)
	return 8
}

// SWAP r8 - r8=E
func (cpu *CPU) cbop33() int {
	cpu.swapr8(&cpu.Registers.E)
	return 8
}

// SWAP r8 - r8=H
func (cpu *CPU) cbop34() int {
	cpu.swapr8(&cpu.Registers.H)
	return 8
}

// SWAP r8 - r8=L
func (cpu *CPU) cbop35() int {
	cpu.swapr8(&cpu.Registers.L)
	return 8
}

// SWAP [HL]
func (cpu *CPU) cbop36() int {
	cpu.swapmemhl()
	return 16
}

// SWAP r8 - r8=A
func (cpu *CPU) cbop37() int {
	cpu.swapr8(&cpu.Registers.A)
	return 8
}

// SRL r8 - r8=B
func (cpu *CPU) cbop38() int {
	cpu.srlr8(&cpu.Registers.B)
	return 8
}

// SRL r8 - r8=C
func (cpu *CPU) cbop39() int {
	cpu.srlr8(&cpu.Registers.C)
	return 8
}

// SRL r8 - r8=D
func (cpu *CPU) cbop3a() int {
	cpu.srlr8(&cpu.Registers.D)
	return 8
}

// SRL r8 - r8=E
func (cpu *CPU) cbop3b() int {
	cpu.srlr8(&cpu.Registers.E)
	return 8
}

// SRL r8 - r8=H
func (cpu *CPU) cbop3c() int {
	cpu.srlr8(&cpu.Registers.H)
	return 8
}

// SRL r8 - r8=L
func (cpu *CPU) cbop3d() int {
	cpu.srlr8(&cpu.Registers.L)
	return 8
}

// SRL [HL]
func (cpu *CPU) cbop3e() int {
	cpu.srlmemhl()
	return 16
}

// SRL r8 - r8=A
func (cpu *CPU) cbop3f() int {
	cpu.srlr8(&cpu.Registers.A)
	return 8
}

// BIT u3, r8 - u3=0,r8=B
func (cpu *CPU) cbop40() int {
	cpu.bitu3r8(0, &cpu.Registers.B)
	return 8
}

// BIT u3, r8 - u3=0,r8=C
func (cpu *CPU) cbop41() int {
	cpu.bitu3r8(0, &cpu.Registers.C)
	return 8
}

// BIT u3, r8 - u3=0,r8=D
func (cpu *CPU) cbop42() int {
	cpu.bitu3r8(0, &cpu.Registers.D)
	return 8
}

// BIT u3, r8 - u3=0,r8=E
func (cpu *CPU) cbop43() int {
	cpu.bitu3r8(0, &cpu.Registers.E)
	return 8
}

// BIT u3, r8 - u3=0,r8=H
func (cpu *CPU) cbop44() int {
	cpu.bitu3r8(0, &cpu.Registers.H)
	return 8
}

// BIT u3, r8 - u3=0,r8=L
func (cpu *CPU) cbop45() int {
	cpu.bitu3r8(0, &cpu.Registers.L)
	return 8
}

// BIT u3, [HL] - u3=0
func (cpu *CPU) cbop46() int {
	cpu.bitu3memhl(0)
	return 12
}

// BIT u3, r8 - u3=0,r8=A
func (cpu *CPU) cbop47() int {
	cpu.bitu3r8(0, &cpu.Registers.A)
	return 8
}

// BIT u3, r8 - u3=1,r8=B
func (cpu *CPU) cbop48() int {
	cpu.bitu3r8(1, &cpu.Registers.B)
	return 8
}

// BIT u3, r8 - u3=1,r8=C
func (cpu *CPU) cbop49() int {
	cpu.bitu3r8(1, &cpu.Registers.C)
	return 8
}

// BIT u3, r8 - u3=1,r8=D
func (cpu *CPU) cbop4a() int {
	cpu.bitu3r8(1, &cpu.Registers.D)
	return 8
}

// BIT u3, r8 - u3=1,r8=E
func (cpu *CPU) cbop4b() int {
	cpu.bitu3r8(1, &cpu.Registers.E)
	return 8
}

// BIT u3, r8 - u3=1,r8=H
func (cpu *CPU) cbop4c() int {
	cpu.bitu3r8(1, &cpu.Registers.H)
	return 8
}

// BIT u3, r8 - u3=1,r8=L
func (cpu *CPU) cbop4d() int {
	cpu.bitu3r8(1, &cpu.Registers.L)
	return 8
}

// BIT u3, [HL] - u3=1
func (cpu *CPU) cbop4e() int {
	cpu.bitu3memhl(1)
	return 12
}

// BIT u3, r8 - u3=1,r8=A
func (cpu *CPU) cbop4f() int {
	cpu.bitu3r8(1, &cpu.Registers.A)
	return 8
}

// BIT u3, r8 - u3=2,r8=B
func (cpu *CPU) cbop50() int {
	cpu.bitu3r8(2, &cpu.Registers.B)
	return 8
}

// BIT u3, r8 - u3=2,r8=C
func (cpu *CPU) cbop51() int {
	cpu.bitu3r8(2, &cpu.Registers.C)
	return 8
}

// BIT u3, r8 - u3=2,r8=D
func (cpu *CPU) cbop52() int {
	cpu.bitu3r8(2, &cpu.Registers.D)
	return 8
}

// BIT u3, r8 - u3=2,r8=E
func (cpu *CPU) cbop53() int {
	cpu.bitu3r8(2, &cpu.Registers.E)
	return 8
}

// BIT u3, r8 - u3=2,r8=H
func (cpu *CPU) cbop54() int {
	cpu.bitu3r8(2, &cpu.Registers.H)
	return 8
}

// BIT u3, r8 - u3=2,r8=L
func (cpu *CPU) cbop55() int {
	cpu.bitu3r8(2, &cpu.Registers.L)
	return 8
}

// BIT u3, [HL] - u3=2
func (cpu *CPU) cbop56() int {
	cpu.bitu3memhl(2)
	return 12
}

// BIT u3, r8 - u3=2,r8=A
func (cpu *CPU) cbop57() int {
	cpu.bitu3r8(2, &cpu.Registers.A)
	return 8
}

// BIT u3, r8 - u3=3,r8=B
func (cpu *CPU) cbop58() int {
	cpu.bitu3r8(3, &cpu.Registers.B)
	return 8
}

// BIT u3, r8 - u3=3,r8=C
func (cpu *CPU) cbop59() int {
	cpu.bitu3r8(3, &cpu.Registers.C)
	return 8
}

// BIT u3, r8 - u3=3,r8=D
func (cpu *CPU) cbop5a() int {
	cpu.bitu3r8(3, &cpu.Registers.D)
	return 8
}

// BIT u3, r8 - u3=3,r8=E
func (cpu *CPU) cbop5b() int {
	cpu.bitu3r8(3, &cpu.Registers.E)
	return 8
}

// BIT u3, r8 - u3=3,r8=H
func (cpu *CPU) cbop5c() int {
	cpu.bitu3r8(3, &cpu.Registers.H)
	return 8
}

// BIT u3, r8 - u3=3,r8=L
func (cpu *CPU) cbop5d() int {
	cpu.bitu3r8(3, &cpu.Registers.L)
	return 8
}

// BIT u3, [HL] - u3=3
func (cpu *CPU) cbop5e() int {
	cpu.bitu3memhl(3)
	return 12
}

// BIT u3, r8 - u3=3,r8=A
func (cpu *CPU) cbop5f() int {
	cpu.bitu3r8(3, &cpu.Registers.A)
	return 8
}

// BIT u3, r8 - u3=4,r8=B
func (cpu *CPU) cbop60() int {
	cpu.bitu3r8(4, &cpu.Registers.B)
	return 8
}

// BIT u3, r8 - u3=4,r8=C
func (cpu *CPU) cbop61() int {
	cpu.bitu3r8(4, &cpu.Registers.C)
	return 8
}

// BIT u3, r8 - u3=4,r8=D
func (cpu *CPU) cbop62() int {
	cpu.bitu3r8(4, &cpu.Registers.D)
	return 8
}

// BIT u3, r8 - u3=4,r8=E
func (cpu *CPU) cbop63() int {
	cpu.bitu3r8(4, &cpu.Registers.E)
	return 8
}

// BIT u3, r8 - u3=4,r8=H
func (cpu *CPU) cbop64() int {
	cpu.bitu3r8(4, &cpu.Registers.H)
	return 8
}

// BIT u3, r8 - u3=4,r8=L
func (cpu *CPU) cbop65() int {
	cpu.bitu3r8(4, &cpu.Registers.L)
	return 8
}

// BIT u3, [HL] - u3=4
func (cpu *CPU) cbop66() int {
	cpu.bitu3memhl(4)
	return 12
}

// BIT u3, r8 - u3=4,r8=A
func (cpu *CPU) cbop67() int {
	cpu.bitu3r8(4, &cpu.Registers.A)
	return 8
}

// BIT u3, r8 - u3=5,r8=B
func (cpu *CPU) cbop68() int {
	cpu.bitu3r8(5, &cpu.Registers.B)
	return 8
}

// BIT u3, r8 - u3=5,r8=C
func (cpu *CPU) cbop69() int {
	cpu.bitu3r8(5, &cpu.Registers.C)
	return 8
}

// BIT u3, r8 - u3=5,r8=D
func (cpu *CPU) cbop6a() int {
	cpu.bitu3r8(5, &cpu.Registers.D)
	return 8
}

// BIT u3, r8 - u3=5,r8=E
func (cpu *CPU) cbop6b() int {
	cpu.bitu3r8(5, &cpu.Registers.E)
	return 8
}

// BIT u3, r8 - u3=5,r8=H
func (cpu *CPU) cbop6c() int {
	cpu.bitu3r8(5, &cpu.Registers.H)
	return 8
}

// BIT u3, r8 - u3=5,r8=L
func (cpu *CPU) cbop6d() int {
	cpu.bitu3r8(5, &cpu.Registers.L)
	return 8
}

// BIT u3, [HL] - u3=5
func (cpu *CPU) cbop6e() int {
	cpu.bitu3memhl(5)
	return 12
}

// BIT u3, r8 - u3=5,r8=A
func (cpu *CPU) cbop6f() int {
	cpu.bitu3r8(5, &cpu.Registers.A)
	return 8
}

// BIT u3, r8 - u3=6,r8=B
func (cpu *CPU) cbop70() int {
	cpu.bitu3r8(6, &cpu.Registers.B)
	return 8
}

// BIT u3, r8 - u3=6,r8=C
func (cpu *CPU) cbop71() int {
	cpu.bitu3r8(6, &cpu.Registers.C)
	return 8
}

// BIT u3, r8 - u3=6,r8=D
func (cpu *CPU) cbop72() int {
	cpu.bitu3r8(6, &cpu.Registers.D)
	return 8
}

// BIT u3, r8 - u3=6,r8=E
func (cpu *CPU) cbop73() int {
	cpu.bitu3r8(6, &cpu.Registers.E)
	return 8
}

// BIT u3, r8 - u3=6,r8=H
func (cpu *CPU) cbop74() int {
	cpu.bitu3r8(6, &cpu.Registers.H)
	return 8
}

// BIT u3, r8 - u3=6,r8=L
func (cpu *CPU) cbop75() int {
	cpu.bitu3r8(6, &cpu.Registers.L)
	return 8
}

// BIT u3, [HL] - u3=6
func (cpu *CPU) cbop76() int {
	cpu.bitu3memhl(6)
	return 12
}

// BIT u3, r8 - u3=6,r8=A
func (cpu *CPU) cbop77() int {
	cpu.bitu3r8(6, &cpu.Registers.A)
	return 8
}

// BIT u3, r8 - u3=7,r8=B
func (cpu *CPU) cbop78() int {
	cpu.bitu3r8(7, &cpu.Registers.B)
	return 8
}

// BIT u3, r8 - u3=7,r8=C
func (cpu *CPU) cbop79() int {
	cpu.bitu3r8(7, &cpu.Registers.C)
	return 8
}

// BIT u3, r8 - u3=7,r8=D
func (cpu *CPU) cbop7a() int {
	cpu.bitu3r8(7, &cpu.Registers.D)
	return 8
}

// BIT u3, r8 - u3=7,r8=E
func (cpu *CPU) cbop7b() int {
	cpu.bitu3r8(7, &cpu.Registers.E)
	return 8
}

// BIT u3, r8 - u3=7,r8=H
func (cpu *CPU) cbop7c() int {
	cpu.bitu3r8(7, &cpu.Registers.H)
	return 8
}

// BIT u3, r8 - u3=7,r8=L
func (cpu *CPU) cbop7d() int {
	cpu.bitu3r8(7, &cpu.Registers.L)
	return 8
}

// BIT u3, [HL] - u3=7
func (cpu *CPU) cbop7e() int {
	cpu.bitu3memhl(7)
	return 12
}

// BIT u3, r8 - u3=7,r8=A
func (cpu *CPU) cbop7f() int {
	cpu.bitu3r8(7, &cpu.Registers.A)
	return 8
}

// RES u3, r8 - u3=0,r8=B
func (cpu *CPU) cbop80() int {
	cpu.resu3r8(0, &cpu.Registers.B)
	return 8
}

// RES u3, r8 - u3=0,r8=C
func (cpu *CPU) cbop81() int {
	cpu.resu3r8(0, &cpu.Registers.C)
	return 8
}

// RES u3, r8 - u3=0,r8=D
func (cpu *CPU) cbop82() int {
	cpu.resu3r8(0, &cpu.Registers.D)
	return 8
}

// RES u3, r8 - u3=0,r8=E
func (cpu *CPU) cbop83() int {
	cpu.resu3r8(0, &cpu.Registers.E)
	return 8
}

// RES u3, r8 - u3=0,r8=H
func (cpu *CPU) cbop84() int {
	cpu.resu3r8(0, &cpu.Registers.H)
	return 8
}

// RES u3, r8 - u3=0,r8=L
func (cpu *CPU) cbop85() int {
	cpu.resu3r8(0, &cpu.Registers.L)
	return 8
}

// RES u3, [HL] - u3=0
func (cpu *CPU) cbop86() int {
	cpu.resu3memhl(0)
	return 16
}

// RES u3, r8 - u3=0,r8=A
func (cpu *CPU) cbop87() int {
	cpu.resu3r8(0, &cpu.Registers.A)
	return 8
}

// RES u3, r8 - u3=1,r8=B
func (cpu *CPU) cbop88() int {
	cpu.resu3r8(1, &cpu.Registers.B)
	return 8
}

// RES u3, r8 - u3=1,r8=C
func (cpu *CPU) cbop89() int {
	cpu.resu3r8(1, &cpu.Registers.C)
	return 8
}

// RES u3, r8 - u3=1,r8=D
func (cpu *CPU) cbop8a() int {
	cpu.resu3r8(1, &cpu.Registers.D)
	return 8
}

// RES u3, r8 - u3=1,r8=E
func (cpu *CPU) cbop8b() int {
	cpu.resu3r8(1, &cpu.Registers.E)
	return 8
}

// RES u3, r8 - u3=1,r8=H
func (cpu *CPU) cbop8c() int {
	cpu.resu3r8(1, &cpu.Registers.H)
	return 8
}

// RES u3, r8 - u3=1,r8=L
func (cpu *CPU) cbop8d() int {
	cpu.resu3r8(1, &cpu.Registers.L)
	return 8
}

// RES u3, [HL] - u3=1
func (cpu *CPU) cbop8e() int {
	cpu.resu3memhl(1)
	return 16
}

// RES u3, r8 - u3=1,r8=A
func (cpu *CPU) cbop8f() int {
	cpu.resu3r8(1, &cpu.Registers.A)
	return 8
}

// RES u3, r8 - u3=2,r8=B
func (cpu *CPU) cbop90() int {
	cpu.resu3r8(2, &cpu.Registers.B)
	return 8
}

// RES u3, r8 - u3=2,r8=C
func (cpu *CPU) cbop91() int {
	cpu.resu3r8(2, &cpu.Registers.C)
	return 8
}

// RES u3, r8 - u3=2,r8=D
func (cpu *CPU) cbop92() int {
	cpu.resu3r8(2, &cpu.Registers.D)
	return 8
}

// RES u3, r8 - u3=2,r8=E
func (cpu *CPU) cbop93() int {
	cpu.resu3r8(2, &cpu.Registers.E)
	return 8
}

// RES u3, r8 - u3=2,r8=H
func (cpu *CPU) cbop94() int {
	cpu.resu3r8(2, &cpu.Registers.H)
	return 8
}

// RES u3, r8 - u3=2,r8=L
func (cpu *CPU) cbop95() int {
	cpu.resu3r8(2, &cpu.Registers.L)
	return 8
}

// RES u3, [HL] - u3=2
func (cpu *CPU) cbop96() int {
	cpu.resu3memhl(2)
	return 16
}

// RES u3, r8 - u3=2,r8=A
func (cpu *CPU) cbop97() int {
	cpu.resu3r8(2, &cpu.Registers.A)
	return 8
}

// RES u3, r8 - u3=3,r8=B
func (cpu *CPU) cbop98() int {
	cpu.resu3r8(3, &cpu.Registers.B)
	return 8
}

// RES u3, r8 - u3=3,r8=C
func (cpu *CPU) cbop99() int {
	cpu.resu3r8(3, &cpu.Registers.C)
	return 8
}

// RES u3, r8 - u3=3,r8=D
func (cpu *CPU) cbop9a() int {
	cpu.resu3r8(3, &cpu.Registers.D)
	return 8
}

// RES u3, r8 - u3=3,r8=E
func (cpu *CPU) cbop9b() int {
	cpu.resu3r8(3, &cpu.Registers.E)
	return 8
}

// RES u3, r8 - u3=3,r8=H
func (cpu *CPU) cbop9c() int {
	cpu.resu3r8(3, &cpu.Registers.H)
	return 8
}

// RES u3, r8 - u3=3,r8=L
func (cpu *CPU) cbop9d() int {
	cpu.resu3r8(3, &cpu.Registers.L)
	return 8
}

// RES u3, [HL] - u3=3
func (cpu *CPU) cbop9e() int {
	cpu.resu3memhl(3)
	return 16
}

// RES u3, r8 - u3=3,r8=A
func (cpu *CPU) cbop9f() int {
	cpu.resu3r8(3, &cpu.Registers.A)
	return 8
}

// RES u3, r8 - u3=4,r8=B
func (cpu *CPU) cbopa0() int {
	cpu.resu3r8(4, &cpu.Registers.B)
	return 8
}

// RES u3, r8 - u3=4,r8=C
func (cpu *CPU) cbopa1() int {
	cpu.resu3r8(4, &cpu.Registers.C)
	return 8
}

// RES u3, r8 - u3=4,r8=D
func (cpu *CPU) cbopa2() int {
	cpu.resu3r8(4, &cpu.Registers.D)
	return 8
}

// RES u3, r8 - u3=4,r8=E
func (cpu *CPU) cbopa3() int {
	cpu.resu3r8(4, &cpu.Registers.E)
	return 8
}

// RES u3, r8 - u3=4,r8=H
func (cpu *CPU) cbopa4() int {
	cpu.resu3r8(4, &cpu.Registers.H)
	return 8
}

// RES u3, r8 - u3=4,r8=L
func (cpu *CPU) cbopa5() int {
	cpu.resu3r8(4, &cpu.Registers.L)
	return 8
}

// RES u3, [HL] - u3=4
func (cpu *CPU) cbopa6() int {
	cpu.resu3memhl(4)
	return 16
}

// RES u3, r8 - u3=4,r8=A
func (cpu *CPU) cbopa7() int {
	cpu.resu3r8(4, &cpu.Registers.A)
	return 8
}

// RES u3, r8 - u3=5,r8=B
func (cpu *CPU) cbopa8() int {
	cpu.resu3r8(5, &cpu.Registers.B)
	return 8
}

// RES u3, r8 - u3=5,r8=C
func (cpu *CPU) cbopa9() int {
	cpu.resu3r8(5, &cpu.Registers.C)
	return 8
}

// RES u3, r8 - u3=5,r8=D
func (cpu *CPU) cbopaa() int {
	cpu.resu3r8(5, &cpu.Registers.D)
	return 8
}

// RES u3, r8 - u3=5,r8=E
func (cpu *CPU) cbopab() int {
	cpu.resu3r8(5, &cpu.Registers.E)
	return 8
}

// RES u3, r8 - u3=5,r8=H
func (cpu *CPU) cbopac() int {
	cpu.resu3r8(5, &cpu.Registers.H)
	return 8
}

// RES u3, r8 - u3=5,r8=L
func (cpu *CPU) cbopad() int {
	cpu.resu3r8(5, &cpu.Registers.L)
	return 8
}

// RES u3, [HL] - u3=5
func (cpu *CPU) cbopae() int {
	cpu.resu3memhl(5)
	return 16
}

// RES u3, r8 - u3=5,r8=A
func (cpu *CPU) cbopaf() int {
	cpu.resu3r8(5, &cpu.Registers.A)
	return 8
}

// RES u3, r8 - u3=6,r8=B
func (cpu *CPU) cbopb0() int {
	cpu.resu3r8(6, &cpu.Registers.B)
	return 8
}

// RES u3, r8 - u3=6,r8=C
func (cpu *CPU) cbopb1() int {
	cpu.resu3r8(6, &cpu.Registers.C)
	return 8
}

// RES u3, r8 - u3=6,r8=D
func (cpu *CPU) cbopb2() int {
	cpu.resu3r8(6, &cpu.Registers.D)
	return 8
}

// RES u3, r8 - u3=6,r8=E
func (cpu *CPU) cbopb3() int {
	cpu.resu3r8(6, &cpu.Registers.E)
	return 8
}

// RES u3, r8 - u3=6,r8=H
func (cpu *CPU) cbopb4() int {
	cpu.resu3r8(6, &cpu.Registers.H)
	return 8
}

// RES u3, r8 - u3=6,r8=L
func (cpu *CPU) cbopb5() int {
	cpu.resu3r8(6, &cpu.Registers.L)
	return 8
}

// RES u3, [HL] - u3=6
func (cpu *CPU) cbopb6() int {
	cpu.resu3memhl(6)
	return 16
}

// RES u3, r8 - u3=6,r8=A
func (cpu *CPU) cbopb7() int {
	cpu.resu3r8(6, &cpu.Registers.A)
	return 8
}

// RES u3, r8 - u3=7,r8=B
func (cpu *CPU) cbopb8() int {
	cpu.resu3r8(7, &cpu.Registers.B)
	return 8
}

// RES u3, r8 - u3=7,r8=C
func (cpu *CPU) cbopb9() int {
	cpu.resu3r8(7, &cpu.Registers.C)
	return 8
}

// RES u3, r8 - u3=7,r8=D
func (cpu *CPU) cbopba() int {
	cpu.resu3r8(7, &cpu.Registers.D)
	return 8
}

// RES u3, r8 - u3=7,r8=E
func (cpu *CPU) cbopbb() int {
	cpu.resu3r8(7, &cpu.Registers.E)
	return 8
}

// RES u3, r8 - u3=7,r8=H
func (cpu *CPU) cbopbc() int {
	cpu.resu3r8(7, &cpu.Registers.H)
	return 8
}

// RES u3, r8 - u3=7,r8=L
func (cpu *CPU) cbopbd() int {
	cpu.resu3r8(7, &cpu.Registers.L)
	return 8
}

// RES u3, [HL] - u3=7
func (cpu *CPU) cbopbe() int {
	cpu.resu3memhl(7)
	return 16
}

// RES u3, r8 - u3=7,r8=A
func (cpu *CPU) cbopbf() int {
	cpu.resu3r8(7, &cpu.Registers.A)
	return 8
}

// SET u3, r8 - u3=0,r8=B
func (cpu *CPU) cbopc0() int {
	cpu.setu3r8(0, &cpu.Registers.B)
	return 8
}

// SET u3, r8 - u3=0,r8=C
func (cpu *CPU) cbopc1() int {
	cpu.setu3r8(0, &cpu.Registers.C)
	return 8
}

// SET u3, r8 - u3=0,r8=D
func (cpu *CPU) cbopc2() int {
	cpu.setu3r8(0, &cpu.Registers.D)
	return 8
}

// SET u3, r8 - u3=0,r8=E
func (cpu *CPU) cbopc3() int {
	cpu.setu3r8(0, &cpu.Registers.E)
	return 8
}

// SET u3, r8 - u3=0,r8=H
func (cpu *CPU) cbopc4() int {
	cpu.setu3r8(0, &cpu.Registers.H)
	return 8
}

// SET u3, r8 - u3=0,r8=L
func (cpu *CPU) cbopc5() int {
	cpu.setu3r8(0, &cpu.Registers.L)
	return 8
}

// SET u3, [HL] - u3=0
func (cpu *CPU) cbopc6() int {
	cpu.setu3memhl(0)
	return 16
}

// SET u3, r8 - u3=0,r8=A
func (cpu *CPU) cbopc7() int {
	cpu.setu3r8(0, &cpu.Registers.A)
	return 8
}

// SET u3, r8 - u3=1,r8=B
func (cpu *CPU) cbopc8() int {
	cpu.setu3r8(1, &cpu.Registers.B)
	return 8
}

// SET u3, r8 - u3=1,r8=C
func (cpu *CPU) cbopc9() int {
	cpu.setu3r8(1, &cpu.Registers.C)
	return 8
}

// SET u3, r8 - u3=1,r8=D
func (cpu *CPU) cbopca() int {
	cpu.setu3r8(1, &cpu.Registers.D)
	return 8
}

// SET u3, r8 - u3=1,r8=E
func (cpu *CPU) cbopcb() int {
	cpu.setu3r8(1, &cpu.Registers.E)
	return 8
}

// SET u3, r8 - u3=1,r8=H
func (cpu *CPU) cbopcc() int {
	cpu.setu3r8(1, &cpu.Registers.H)
	return 8
}

// SET u3, r8 - u3=1,r8=L
func (cpu *CPU) cbopcd() int {
	cpu.setu3r8(1, &cpu.Registers.L)
	return 8
}

// SET u3, [HL] - u3=1
func (cpu *CPU) cbopce() int {
	cpu.setu3memhl(1)
	return 16
}

// SET u3, r8 - u3=1,r8=A
func (cpu *CPU) cbopcf() int {
	cpu.setu3r8(1, &cpu.Registers.A)
	return 8
}

// SET u3, r8 - u3=2,r8=B
func (cpu *CPU) cbopd0() int {
	cpu.setu3r8(2, &cpu.Registers.B)
	return 8
}

// SET u3, r8 - u3=2,r8=C
func (cpu *CPU) cbopd1() int {
	cpu.setu3r8(2, &cpu.Registers.C)
	return 8
}

// SET u3, r8 - u3=2,r8=D
func (cpu *CPU) cbopd2() int {
	cpu.setu3r8(2, &cpu.Registers.D)
	return 8
}

// SET u3, r8 - u3=2,r8=E
func (cpu *CPU) cbopd3() int {
	cpu.setu3r8(2, &cpu.Registers.E)
	return 8
}

// SET u3, r8 - u3=2,r8=H
func (cpu *CPU) cbopd4() int {
	cpu.setu3r8(2, &cpu.Registers.H)
	return 8
}

// SET u3, r8 - u3=2,r8=L
func (cpu *CPU) cbopd5() int {
	cpu.setu3r8(2, &cpu.Registers.L)
	return 8
}

// SET u3, [HL] - u3=2
func (cpu *CPU) cbopd6() int {
	cpu.setu3memhl(2)
	return 16
}

// SET u3, r8 - u3=2,r8=A
func (cpu *CPU) cbopd7() int {
	cpu.setu3r8(2, &cpu.Registers.A)
	return 8
}

// SET u3, r8 - u3=3,r8=B
func (cpu *CPU) cbopd8() int {
	cpu.setu3r8(3, &cpu.Registers.B)
	return 8
}

// SET u3, r8 - u3=3,r8=C
func (cpu *CPU) cbopd9() int {
	cpu.setu3r8(3, &cpu.Registers.C)
	return 8
}

// SET u3, r8 - u3=3,r8=D
func (cpu *CPU) cbopda() int {
	cpu.setu3r8(3, &cpu.Registers.D)
	return 8
}

// SET u3, r8 - u3=3,r8=E
func (cpu *CPU) cbopdb() int {
	cpu.setu3r8(3, &cpu.Registers.E)
	return 8
}

// SET u3, r8 - u3=3,r8=H
func (cpu *CPU) cbopdc() int {
	cpu.setu3r8(3, &cpu.Registers.H)
	return 8
}

// SET u3, r8 - u3=3,r8=L
func (cpu *CPU) cbopdd() int {
	cpu.setu3r8(3, &cpu.Registers.L)
	return 8
}

// SET u3, [HL] - u3=3
func (cpu *CPU) cbopde() int {
	cpu.setu3memhl(3)
	return 16
}

// SET u3, r8 - u3=3,r8=A
func (cpu *CPU) cbopdf() int {
	cpu.setu3r8(3, &cpu.Registers.A)
	return 8
}

// SET u3, r8 - u3=4,r8=B
func (cpu *CPU) cbope0() int {
	cpu.setu3r8(4, &cpu.Registers.B)
	return 8
}

// SET u3, r8 - u3=4,r8=C
func (cpu *CPU) cbope1() int {
	cpu.setu3r8(4, &cpu.Registers.C)
	return 8
}

// SET u3, r8 - u3=4,r8=D
func (cpu *CPU) cbope2() int {
	cpu.setu3r8(4, &cpu.Registers.D)
	return 8
}

// SET u3, r8 - u3=4,r8=E
func (cpu *CPU) cbope3() int {
	cpu.setu3r8(4, &cpu.Registers.E)
	return 8
}

// SET u3, r8 - u3=4,r8=H
func (cpu *CPU) cbope4() int {
	cpu.setu3r8(4, &cpu.Registers.H)
	return 8
}

// SET u3, r8 - u3=4,r8=L
func (cpu *CPU) cbope5() int {
	cpu.setu3r8(4, &cpu.Registers.L)
	return 8
}

// SET u3, [HL] - u3=4
func (cpu *CPU) cbope6() int {
	cpu.setu3memhl(4)
	return 16
}

// SET u3, r8 - u3=4,r8=A
func (cpu *CPU) cbope7() int {
	cpu.setu3r8(4, &cpu.Registers.A)
	return 8
}

// SET u3, r8 - u3=5,r8=B
func (cpu *CPU) cbope8() int {
	cpu.setu3r8(5, &cpu.Registers.B)
	return 8
}

// SET u3, r8 - u3=5,r8=C
func (cpu *CPU) cbope9() int {
	cpu.setu3r8(5, &cpu.Registers.C)
	return 8
}

// SET u3, r8 - u3=5,r8=D
func (cpu *CPU) cbopea() int {
	cpu.setu3r8(5, &cpu.Registers.D)
	return 8
}

// SET u3, r8 - u3=5,r8=E
func (cpu *CPU) cbopeb() int {
	cpu.setu3r8(5, &cpu.Registers.E)
	return 8
}

// SET u3, r8 - u3=5,r8=H
func (cpu *CPU) cbopec() int {
	cpu.setu3r8(5, &cpu.Registers.H)
	return 8
}

// SET u3, r8 - u3=5,r8=L
func (cpu *CPU) cboped() int {
	cpu.setu3r8(5, &cpu.Registers.L)
	return 8
}

// SET u3, [HL] - u3=5
func (cpu *CPU) cbopee() int {
	cpu.setu3memhl(5)
	return 16
}

// SET u3, r8 - u3=5,r8=A
func (cpu *CPU) cbopef() int {
	cpu.setu3r8(5, &cpu.Registers.A)
	return 8
}

// SET u3, r8 - u3=6,r8=B
func (cpu *CPU) cbopf0() int {
	cpu.setu3r8(6, &cpu.Registers.B)
	return 8
}

// SET u3, r8 - u3=6,r8=C
func (cpu *CPU) cbopf1() int {
	cpu.setu3r8(6, &cpu.Registers.C)
	return 8
}

// SET u3, r8 - u3=6,r8=D
func (cpu *CPU) cbopf2() int {
	cpu.setu3r8(6, &cpu.Registers.D)
	return 8
}

// SET u3, r8 - u3=6,r8=E
func (cpu *CPU) cbopf3() int {
	cpu.setu3r8(6, &cpu.Registers.E)
	return 8
}

// SET u3, r8 - u3=6,r8=H
func (cpu *CPU) cbopf4() int {
	cpu.setu3r8(6, &cpu.Registers.H)
	return 8
}

// SET u3, r8 - u3=6,r8=L
func (cpu *CPU) cbopf5() int {
	cpu.setu3r8(6, &cpu.Registers.L)
	return 8
}

// SET u3, [HL] - u3=6
func (cpu *CPU) cbopf6() int {
	cpu.setu3memhl(6)
	return 16
}

// SET u3, r8 - u3=6,r8=A
func (cpu *CPU) cbopf7() int {
	cpu.setu3r8(6, &cpu.Registers.A)
	return 8
}

// SET u3, r8 - u3=7,r8=B
func (cpu *CPU) cbopf8() int {
	cpu.setu3r8(7, &cpu.Registers.B)
	return 8
}

// SET u3, r8 - u3=7,r8=C
func (cpu *CPU) cbopf9() int {
	cpu.setu3r8(7, &cpu.Registers.C)
	return 8
}

// SET u3, r8 - u3=7,r8=D
func (cpu *CPU) cbopfa() int {
	cpu.setu3r8(7, &cpu.Registers.D)
	return 8
}

// SET u3, r8 - u3=7,r8=E
func (cpu *CPU) cbopfb() int {
	cpu.setu3r8(7, &cpu.Registers.E)
	return 8
}

// SET u3, r8 - u3=7,r8=H
func (cpu *CPU) cbopfc() int {
	cpu.setu3r8(7, &cpu.Registers.H)
	return 8
}

// SET u3, r8 - u3=7,r8=L
func (cpu *CPU) cbopfd() int {
	cpu.setu3r8(7, &cpu.Registers.L)
	return 8
}

// SET u3, [HL] - u3=7
func (cpu *CPU) cbopfe() int {
	cpu.setu3memhl(7)
	return 16
}

// SET u3, r8 - u3=7,r8=A
func (cpu *CPU) cbopff() int {
	cpu.setu3r8(7, &cpu.Registers.A)
	return 8
}
