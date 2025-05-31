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
	cpu.cpAn8(n8)
	return 8
}

// skip review pentru ca nu e folosit
// SUB A, r8 - r8=H
func (cpu *CPU) op94() int {
	cpu.subAr8(cpu.Registers.H)
	return 4
}

// LD [HL], n8
func (cpu *CPU) op36() int {
	n8 := cpu.getImmediate8()
	cpu.ldmemhlr8(n8)
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
	cpu.rst(0x38)
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

// SWAP r8 - r8=A
func (cpu *CPU) cbop37() int {
	cpu.swapr8(&cpu.Registers.A)
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
	cpu.rst(0x28)
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

// RES u3, r8 - u3=0, r8=A
func (cpu *CPU) cbop87() int {
	cpu.resu3r8(0, &cpu.Registers.A)
	return 8
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

// BIT u3, r8 - u3=7, r8=A
func (cpu *CPU) cbop7f() int {
	cpu.bitu3r8(7, &cpu.Registers.A)
	return 8
}

// OR A, n8
func (cpu *CPU) opf6() int {
	n8 := cpu.getImmediate8()
	cpu.orA(n8)
	return 8
}

// RES u3, [HL] - u3=0
func (cpu *CPU) cbop86() int {
	cpu.resu3memhl(0)
	return 16
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
	cpu.rst(0x18)
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
	cpu.ldspr16(cpu.Registers.getHL)
	return 8
}

// RST $08
func (cpu *CPU) opcf() int {
	cpu.rst(0x08)
	return 16
}

// SLA r8 - r8=A
func (cpu *CPU) cbop27() int {
	cpu.slar8(&cpu.Registers.A)
	return 8
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
	cpu.ldmemhlr8(cpu.Registers.E)
	return 8
}

// LD [HL], r8 - r8=D
func (cpu *CPU) op72() int {
	cpu.ldmemhlr8(cpu.Registers.D)
	return 8
}

// LD [HL], r8 - r8=C
func (cpu *CPU) op71() int {
	cpu.ldmemhlr8(cpu.Registers.C)
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

// BIT u3, r8 - u3=2,r8=B
func (cpu *CPU) cbop50() int {
	cpu.bitu3r8(2, &cpu.Registers.B)
	return 8
}

// BIT u3, r8 - u3=4,r8=B
func (cpu *CPU) cbop60() int {
	cpu.bitu3r8(4, &cpu.Registers.B)
	return 8
}

// BIT u3, r8 - u3=5,r8=B
func (cpu *CPU) cbop68() int {
	cpu.bitu3r8(5, &cpu.Registers.B)
	return 8
}

// BIT u3, r8 - u3=3,r8=B
func (cpu *CPU) cbop58() int {
	cpu.bitu3r8(3, &cpu.Registers.B)
	return 8
}

// BIT u3, [HL] - u3=7
func (cpu *CPU) cbop7e() int {
	cpu.bitu3memhl(7)
	return 12
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

// RES u3, r8 - u3=7,r8=A
func (cpu *CPU) cbopbf() int {
	cpu.resu3r8(7, &cpu.Registers.A)
	return 8
}

// LD r8, n8 - r8=H
func (cpu *CPU) op26() int {
	cpu.ldr8n8(&cpu.Registers.H)
	return 8
}

// SUB A, r8 - r8=L
func (cpu *CPU) op95() int {
	cpu.subar8(&cpu.Registers.L)
	return 4
}
