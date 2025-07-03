// package GameBoy_Emulator
package main

import (
	"fmt"

	log "github.com/mgutz/logxi/v1"
)

var cpulogger log.Logger

// A, B, C, D, E, F, H, and L - 8 bit
// SP, PC - 16 bit
// AF, BC, DE, and HL
type CPU struct {
	Registers    Registers
	Memory       [65536]uint8
	timer        Timer
	Cartridge    *Cartridge
	graphics     *Graphics
	joypad       *Joypad
	IME          bool // interrupt master enable
	IMEScheduled bool //enable IME after one instr
	halted       bool
	haltBug      bool
	stopped      bool
	IE           uint8 // FFFF — IE: Interrupt enable
	IF           uint8 //FF0F — IF: Interrupt flag
}

type Registers struct {
	A, B, C, D, E, F, H, L uint8
	SP, PC                 uint16
}

// flags
const flagZ uint8 = 1 << 7 // zero flag
const flagN uint8 = 1 << 6 // sub flag
const flagH uint8 = 1 << 5 // half carry flag
const flagC uint8 = 1 << 4 // carry flag

func NewCPU() *CPU {
	cpulogger = log.New("cpu")

	cpu := &CPU{
		joypad: &Joypad{},
	}
	cpu.IME = false

	return cpu

}

func (register *Registers) setFlag(flag uint8, on bool) {
	if on {
		register.F |= flag //set bit
	} else {
		register.F &= ^flag //clear bit
	}
}

func (register *Registers) getFlag(flag uint8) bool {
	return register.F&flag != 0
}

func (register *Registers) getAF() uint16 {
	return uint16(register.A)<<8 | uint16(register.F)
}

func (register *Registers) setAF(value uint16) {
	register.A = uint8(value >> 8)
	register.F = uint8(value & 0xFF)
}

// B-hi C-lo

func (register *Registers) getBC() uint16 {
	//			most significant	least significant
	return uint16(register.B)<<8 | uint16(register.C)
}

func (register *Registers) setBC(value uint16) {
	register.B = uint8(value >> 8)   //take upper bits
	register.C = uint8(value & 0xFF) // take lower bits
}

// D-hi E-lo

func (register *Registers) getDE() uint16 {
	return uint16(register.D)<<8 | uint16(register.E)
}

func (register *Registers) setDE(value uint16) {
	register.D = uint8(value >> 8)   //take upper bits
	register.E = uint8(value & 0xFF) // take lower bits
}

// H-hi L-lo

func (register *Registers) getHL() uint16 {
	return uint16(register.H)<<8 | uint16(register.L)
}

func (register *Registers) setHL(value uint16) {
	register.H = uint8(value >> 8)   //take upper bits
	register.L = uint8(value & 0xFF) // take lower bits
}

func (cpu *CPU) getIE() uint8 {
	return cpu.memoryRead(0xFFFF)
}

func (cpu *CPU) getIF() uint8 {
	return cpu.memoryRead(0xFF0F)
}
func (cpu *CPU) checkSchedule() {
	if cpu.IMEScheduled {
		cpu.IME = true
		cpu.IMEScheduled = false
	}
}

func (cpu *CPU) fetchOpcode() uint8 {
	opcode := cpu.memoryRead(cpu.Registers.PC)
	if cpu.haltBug {
		cpu.haltBug = false // PC not incremented
	} else {
		cpu.Registers.PC++
	}
	return opcode
}

func (cpu *CPU) push(n uint16) {
	hi := (n & 0xFF00) >> 8
	lo := n & 0xFF
	// hi := n >> 8
	// lo := n
	cpu.Registers.SP -= 1

	cpu.memoryWrite(cpu.Registers.SP, uint8(hi))
	//cpu.memoryWrite(cpu.Registers.SP, byte(n>>8))

	cpu.Registers.SP -= 1
	if cpu.Registers.SP == 0xFF80 {
		panic(cpulogger.Error("Stack smash"))
	}
	cpu.memoryWrite(cpu.Registers.SP, uint8(lo))
}

func (cpu *CPU) pop() uint16 {
	lo := uint16(cpu.memoryRead(cpu.Registers.SP))
	hi := uint16(cpu.memoryRead(cpu.Registers.SP + 1))
	value := hi<<8 | lo
	cpu.Registers.SP += 2

	if cpu.Registers.SP == 0xFFFE+1 {
		panic(cpulogger.Error("Stack smash"))
	}
	return value
}

func (cpu *CPU) handleInterruptions() bool {
	// The IME (interrupt master enable) flag is reset by DI
	// and prohibits all interrupts. It is set by EI and
	// acknowledges the interrupt setting by the IE register.
	// 1. When an interrupt is generated, the IF flag will be
	// set.
	// 2. If the IME flag is set & the corresponding IE flag
	// is set, the following 3 steps are performed.
	// 3. Reset the IME flag and prevent all interrupts.
	// 4. The PC (program counter) is pushed onto the stack.
	// 5. Jump to the starting address of the interrupt.
	if !cpu.IME {
		return false
	}
	IE := cpu.getIE()
	IF := cpu.getIF()

	interruptions := IE & IF
	if interruptions == 0 {
		return false
	}
	//The priorities follow the order of the bits in the IE and IF registers:
	//Bit 0 (VBlank) has the highest priority, and
	//Bit 4 (Joypad) has the lowest priority.
	//	7 6	5	  4		  3		 2	     1	  0
	//IF		Joypad	Serial	Timer	LCD	VBlank
	var addr uint16
	var bit byte
	if interruptions&0x01 != 0 { //VBlank
		addr = 0x40
		bit = 0x01
	} else if interruptions&0x02 != 0 { //LCD
		addr = 0x48
		bit = 0x02
	} else if interruptions&0x04 != 0 { //Timer
		addr = 0x50
		bit = 0x04
	} else if interruptions&0x08 != 0 { //Serial
		addr = 0x58
		bit = 0x08
	} else if interruptions&0x10 != 0 { // Joypad
		addr = 0x60
		bit = 0x10
	}
	if addr != 0 {
		cpu.IME = false                  //reset the IME flag and prevent all interrupts
		cpu.memoryWrite(0xFF0F, IF&^bit) //clear bit
		cpu.push(cpu.Registers.PC)
		cpu.Registers.PC = addr
		cpu.graphics.cycle += 20
		return true
	}
	return false
}

func (cpu *CPU) memoryWrite(address uint16, value byte) {
	if address == 0xFF00 {
		cpu.joypad.write(value)
		cpulogger.Debug(fmt.Sprintf("write in 0xFF00 %08b", cpu.Memory[address]))
	} else if address >= 0xFF04 && address <= 0xFF07 {
		cpu.timer.Write(address, value)
	} else if address == 0xFF46 {
		cpu.Memory[address] = value
		cpu.dmaTransfer(value)
	} else if address == 0xFF89 {
		cpulogger.Debug(fmt.Sprintf("write in 0xFF89 %08b", value))
		cpu.Memory[address] = value
	} else if address >= 0xC000 && address <= 0xCFFF {
		cpu.Memory[address] = value
		// if address+0x2000 <= 0xFFFF {
		cpu.Memory[address+0x2000] = value
		// }
	} else if address >= 0xE000 && address <= 0xFDFF {
		cpu.Memory[address] = value
		cpu.Memory[address-0x2000] = value
	} else if address >= 0xFEA0 && address <= 0xFEFF {
		return
	} else if address >= VRAM_START && address <= VRAM_END {
		//cpu.graphics.writeVRAM(address, value)
		cpu.Memory[address] = value
	} else if address >= OAM_START && address <= OAM_END {
		if address == 0xFE00 {
			logger.Debug(fmt.Sprintf("write in 0xFE00 %08b", value))
		}
		cpu.Memory[address] = value
	} else if address == 0xFF40 {
		cpulogger.Debug(fmt.Sprintf("!!LCDC WRITE: 0x%02X\n", value))
		cpu.Memory[address] = value
	} else if address < 0x8000 {
		return
	} else {
		cpu.Memory[address] = value
	}
}
func (cpu *CPU) memoryRead(address uint16) byte {
	if address == 0xFF00 {
		return cpu.joypad.read()
	} else if address >= 0xFEA0 && address <= 0xFEFF {
		return 0xFF //TODO : idk
	} else {
		return cpu.Memory[address]
	}

}

func (cpu *CPU) execOpcodes() int {
	IE := cpu.getIE()
	IF := cpu.getIF()
	if cpu.halted {
		cpulogger.Debug(fmt.Sprintf("halted with IF: %x and IE: %x", IF, IE))
		if (IE & IF & 0x1F) != 0 {
			cpu.halted = false
		}
	}
	if cpu.stopped {

		cpulogger.Debug("stopped")
		return 0
	}
	if cpu.handleInterruptions() {
		cpulogger.Debug("interruptions !!1")
		return 5
	}
	var tCycles = -1
	prefix := cpu.fetchOpcode()

	isPrefixed := (prefix == 0xcb)
	opcode := prefix
	if cpu.Registers.getDE() == 0x9a05 {
		fmt.Println(cpu.Registers.PC - 1)
		fmt.Println("STOPP")
	}
	cpulogger.Debug(fmt.Sprintf("Executing opcode 0x%x @PC=0x%x A=0x%x F=0x%x DE=0x%x HL=0x%x BC=0x%x IE&IF= 0x%0x", opcode, cpu.Registers.PC-1, cpu.Registers.A, cpu.Registers.F, cpu.Registers.getDE(), cpu.Registers.getHL(), cpu.Registers.getBC(), IE&IF))
	if prefix == 0xcb {
		// prefixed
		opcode = cpu.fetchOpcode()
		cpulogger.Debug(fmt.Sprintf("Executing in cbprefixed opcode 0x%x @PC=0x%x A=0x%x F=0x%x ", opcode, cpu.Registers.PC-1, cpu.Registers.A, cpu.Registers.F))
		switch opcode {
		case 0x37:
			tCycles = cpu.cbop37()
		case 0x87:
			tCycles = cpu.cbop87()
		case 0x7f:
			tCycles = cpu.cbop7f()
		case 0x86:
			tCycles = cpu.cbop86()
		case 0x27:
			tCycles = cpu.cbop27()
		case 0x50:
			tCycles = cpu.cbop50()
		case 0x60:
			tCycles = cpu.cbop60()
		case 0x68:
			tCycles = cpu.cbop68()
		case 0x58:
			tCycles = cpu.cbop58()
		case 0x7e:
			tCycles = cpu.cbop7e()
		case 0x6f:
			tCycles = cpu.cbop6f()
		case 0x28:
			tCycles = cpu.cbop28()
		case 0x29:
			tCycles = cpu.cbop29()

		// ----done tetris----
		// ----start tennis---
		case 0xbf:
			tCycles = cpu.cbopbf()
		case 0x3f:
			tCycles = cpu.cbop3f()
		case 0xbe:
			tCycles = cpu.cbopbe()
		case 0xc7:
			tCycles = cpu.cbopc7()
		case 0xff:
			tCycles = cpu.cbopff()
		case 0x12:
			tCycles = cpu.cbop12()
		//////////////////
		case 0x00:
			tCycles = cpu.cbop00()
			////// joypad
		case 0x40:
			tCycles = cpu.cbop40()
		case 0xa9:
			tCycles = cpu.cbopa9()
		case 0x13:
			tCycles = cpu.cbop13()
		case 0x33:
			tCycles = cpu.cbop33()
		case 0xf7:
			tCycles = cpu.cbopf7()
		case 0xfd:
			tCycles = cpu.cbopfd()
		case 0x49:
			tCycles = cpu.cbop49()
		case 0xb8:
			tCycles = cpu.cbopb8()
		case 0xf8:
			tCycles = cpu.cbopf8()
		////
		case 0x57:
			tCycles = cpu.cbop57()
		case 0x5f:
			tCycles = cpu.cbop5f()
		case 0x70:
			tCycles = cpu.cbop70()
		case 0x78:
			tCycles = cpu.cbop78()
		case 0x47:
			tCycles = cpu.cbop47()
		case 0x4f:
			tCycles = cpu.cbop4f()
		case 0x77:
			tCycles = cpu.cbop77()
		case 0x48:
			tCycles = cpu.cbop48()
		case 0x41:
			tCycles = cpu.cbop41()
		case 0xb0:
			tCycles = cpu.cbopb0()
		case 0xfe:
			tCycles = cpu.cbopfe()
		case 0x23:
			tCycles = cpu.cbop23()
		case 0x11:
			tCycles = cpu.cbop11()
		case 0x3d:
			tCycles = cpu.cbop3d()
		case 0xc9:
			tCycles = cpu.cbopc9()
		case 0xd1:
			tCycles = cpu.cbopd1()
		case 0x69:
			tCycles = cpu.cbop69()
		case 0x61:
			tCycles = cpu.cbop61()
		case 0x71:
			tCycles = cpu.cbop71()
		case 0x79:
			tCycles = cpu.cbop79()
		case 0x10:
			tCycles = cpu.cbop10()
		case 0x21:
			tCycles = cpu.cbop21()
		case 0x17:
			tCycles = cpu.cbop17()
		case 0x39:
			tCycles = cpu.cbop39()
		case 0x1c:
			tCycles = cpu.cbop1c()
		case 0x1d:
			tCycles = cpu.cbop1d()
		case 0xb6:
			tCycles = cpu.cbopb6()
		case 0x2f:
			tCycles = cpu.cbop2f()
		case 0xd9:
			tCycles = cpu.cbopd9()
		case 0xc1:
			tCycles = cpu.cbopc1()
		case 0x38:
			tCycles = cpu.cbop38()
		case 0x67:
			tCycles = cpu.cbop67()
		case 0xee:
			tCycles = cpu.cbopee()
		case 0x81:
			tCycles = cpu.cbop81()
		case 0x89:
			tCycles = cpu.cbop89()
		case 0x5e:
			tCycles = cpu.cbop5e()
		case 0xde:
			tCycles = cpu.cbopde()
		case 0xef:
			tCycles = cpu.cbopef()
		case 0x46:
			tCycles = cpu.cbop46()
		case 0x01:
			tCycles = cpu.cbop01()
		case 0x02:
			tCycles = cpu.cbop02()
		case 0x03:
			tCycles = cpu.cbop03()
		case 0x04:
			tCycles = cpu.cbop04()
		case 0x05:
			tCycles = cpu.cbop05()
		case 0x06:
			tCycles = cpu.cbop06()
		case 0x07:
			tCycles = cpu.cbop07()
		case 0x08:
			tCycles = cpu.cbop08()
		case 0x09:
			tCycles = cpu.cbop09()
		case 0x0a:
			tCycles = cpu.cbop0a()
		case 0x0b:
			tCycles = cpu.cbop0b()
		case 0x0c:
			tCycles = cpu.cbop0c()
		case 0x0d:
			tCycles = cpu.cbop0d()
		case 0x0e:
			tCycles = cpu.cbop0e()
		case 0x0f:
			tCycles = cpu.cbop0f()
		case 0x18:
			tCycles = cpu.cbop18()
		case 0x19:
			tCycles = cpu.cbop19()
		case 0x1a:
			tCycles = cpu.cbop1a()
		case 0x1b:
			tCycles = cpu.cbop1b()
		case 0x1e:
			tCycles = cpu.cbop1e()
		case 0x1f:
			tCycles = cpu.cbop1f()
		case 0x3a:
			tCycles = cpu.cbop3a()
		case 0x3b:
			tCycles = cpu.cbop3b()
		case 0x3c:
			tCycles = cpu.cbop3c()

		case 0x3e:
			tCycles = cpu.cbop3e()
		case 0x2a:
			tCycles = cpu.cbop2a()
		case 0x2b:
			tCycles = cpu.cbop2b()
		case 0x2c:
			tCycles = cpu.cbop2c()
		case 0x2d:
			tCycles = cpu.cbop2d()
		case 0x2e:
			tCycles = cpu.cbop2e()
		case 0x20:
			tCycles = cpu.cbop20()
		case 0x22:
			tCycles = cpu.cbop22()
		case 0x24:
			tCycles = cpu.cbop24()
		case 0x25:
			tCycles = cpu.cbop25()
		case 0x26:
			tCycles = cpu.cbop26()
		case 0x30:
			tCycles = cpu.cbop30()
		case 0x31:
			tCycles = cpu.cbop31()
		case 0x32:
			tCycles = cpu.cbop32()
		case 0x34:
			tCycles = cpu.cbop34()
		case 0x35:
			tCycles = cpu.cbop35()
		case 0x36:
			tCycles = cpu.cbop36()
		case 0x8e:
			tCycles = cpu.cbop8e()
		case 0x8f:
			tCycles = cpu.cbop8f()
		case 0x90:
			tCycles = cpu.cbop90()
		case 0x91:
			tCycles = cpu.cbop91()
		case 0x92:
			tCycles = cpu.cbop92()
		case 0x93:
			tCycles = cpu.cbop93()
		case 0x94:
			tCycles = cpu.cbop94()
		case 0x95:
			tCycles = cpu.cbop95()
		case 0x96:
			tCycles = cpu.cbop96()
		case 0x97:
			tCycles = cpu.cbop97()
		case 0x98:
			tCycles = cpu.cbop98()
		case 0x99:
			tCycles = cpu.cbop99()
		case 0x9a:
			tCycles = cpu.cbop9a()
		case 0x9b:
			tCycles = cpu.cbop9b()
		case 0x9c:
			tCycles = cpu.cbop9c()
		case 0x9d:
			tCycles = cpu.cbop9d()
		case 0x9e:
			tCycles = cpu.cbop9e()
		case 0x9f:
			tCycles = cpu.cbop9f()
		case 0xa0:
			tCycles = cpu.cbopa0()
		case 0xa1:
			tCycles = cpu.cbopa1()
		case 0xa2:
			tCycles = cpu.cbopa2()
		case 0xa3:
			tCycles = cpu.cbopa3()
		case 0xa4:
			tCycles = cpu.cbopa4()
		case 0xa5:
			tCycles = cpu.cbopa5()
		case 0xa6:
			tCycles = cpu.cbopa6()
		case 0xa7:
			tCycles = cpu.cbopa7()
		case 0xa8:
			tCycles = cpu.cbopa8()
		case 0xaa:
			tCycles = cpu.cbopaa()
		case 0xab:
			tCycles = cpu.cbopab()
		case 0xac:
			tCycles = cpu.cbopac()
		case 0xad:
			tCycles = cpu.cbopad()
		case 0xae:
			tCycles = cpu.cbopae()
		case 0xaf:
			tCycles = cpu.cbopaf()
		case 0x14:
			tCycles = cpu.cbop14()
		case 0x15:
			tCycles = cpu.cbop15()
		case 0x16:
			tCycles = cpu.cbop16()
		case 0xf0:
			tCycles = cpu.cbopf0()
		case 0xb1:
			tCycles = cpu.cbopb1()
		case 0xb2:
			tCycles = cpu.cbopb2()
		case 0xb3:
			tCycles = cpu.cbopb3()
		case 0xb4:
			tCycles = cpu.cbopb4()
		case 0xb5:
			tCycles = cpu.cbopb5()
		case 0xb7:
			tCycles = cpu.cbopb7()
		case 0xb9:
			tCycles = cpu.cbopb9()
		case 0xba:
			tCycles = cpu.cbopba()
		case 0xbb:
			tCycles = cpu.cbopbb()
		case 0xbc:
			tCycles = cpu.cbopbc()
		case 0xbd:
			tCycles = cpu.cbopbd()
		case 0x4d:
			tCycles = cpu.cbop4d()
		case 0x4e:
			tCycles = cpu.cbop4e()
		case 0x51:
			tCycles = cpu.cbop51()
		case 0x52:
			tCycles = cpu.cbop52()
		case 0x53:
			tCycles = cpu.cbop53()
		case 0x54:
			tCycles = cpu.cbop54()
		case 0x55:
			tCycles = cpu.cbop55()
		case 0x56:
			tCycles = cpu.cbop56()
		case 0x59:
			tCycles = cpu.cbop59()
		case 0x5a:
			tCycles = cpu.cbop5a()
		case 0x5b:
			tCycles = cpu.cbop5b()
		case 0x5c:
			tCycles = cpu.cbop5c()
		case 0x5d:
			tCycles = cpu.cbop5d()
		case 0x62:
			tCycles = cpu.cbop62()
		case 0x63:
			tCycles = cpu.cbop63()
		case 0x64:
			tCycles = cpu.cbop64()
		case 0x65:
			tCycles = cpu.cbop65()
		case 0x66:
			tCycles = cpu.cbop66()
		case 0x6a:
			tCycles = cpu.cbop6a()
		case 0x6b:
			tCycles = cpu.cbop6b()
		case 0x6c:
			tCycles = cpu.cbop6c()
		case 0x6d:
			tCycles = cpu.cbop6d()
		case 0x6e:
			tCycles = cpu.cbop6e()
		case 0x72:
			tCycles = cpu.cbop72()
		case 0x73:
			tCycles = cpu.cbop73()
		case 0x74:
			tCycles = cpu.cbop74()
		case 0x75:
			tCycles = cpu.cbop75()
		case 0x76:
			tCycles = cpu.cbop76()
		case 0x7a:
			tCycles = cpu.cbop7a()
		case 0x7b:
			tCycles = cpu.cbop7b()
		case 0x7c:
			tCycles = cpu.cbop7c()
		case 0x7d:
			tCycles = cpu.cbop7d()
		case 0xc0:
			tCycles = cpu.cbopc0()
		case 0xc2:
			tCycles = cpu.cbopc2()
		case 0xc3:
			tCycles = cpu.cbopc3()
		case 0xc4:
			tCycles = cpu.cbopc4()
		case 0xc5:
			tCycles = cpu.cbopc5()
		case 0xc6:
			tCycles = cpu.cbopc6()
		case 0xc8:
			tCycles = cpu.cbopc8()
		case 0xca:
			tCycles = cpu.cbopca()
		case 0xcb:
			tCycles = cpu.cbopcb()
		case 0xcc:
			tCycles = cpu.cbopcc()
		case 0xcd:
			tCycles = cpu.cbopcd()
		case 0xce:
			tCycles = cpu.cbopce()
		case 0xcf:
			tCycles = cpu.cbopcf()
		case 0xd0:
			tCycles = cpu.cbopd0()
		case 0xd2:
			tCycles = cpu.cbopd2()
		case 0xd3:
			tCycles = cpu.cbopd3()
		case 0xd4:
			tCycles = cpu.cbopd4()
		case 0xd5:
			tCycles = cpu.cbopd5()
		case 0xd6:
			tCycles = cpu.cbopd6()
		case 0xd7:
			tCycles = cpu.cbopd7()
		case 0xd8:
			tCycles = cpu.cbopd8()
		case 0xda:
			tCycles = cpu.cbopda()
		case 0xdb:
			tCycles = cpu.cbopdb()
		case 0xdc:
			tCycles = cpu.cbopdc()
		case 0xdd:
			tCycles = cpu.cbopdd()
		case 0xdf:
			tCycles = cpu.cbopdf()
		case 0xe0:
			tCycles = cpu.cbope0()
		case 0xe1:
			tCycles = cpu.cbope1()
		case 0xe2:
			tCycles = cpu.cbope2()
		case 0xe3:
			tCycles = cpu.cbope3()
		case 0xe4:
			tCycles = cpu.cbope4()
		case 0xe5:
			tCycles = cpu.cbope5()
		case 0xe6:
			tCycles = cpu.cbope6()
		case 0xe7:
			tCycles = cpu.cbope7()
		case 0xe8:
			tCycles = cpu.cbope8()
		case 0xe9:
			tCycles = cpu.cbope9()
		case 0xea:
			tCycles = cpu.cbopea()
		case 0xeb:
			tCycles = cpu.cbopeb()
		case 0xec:
			tCycles = cpu.cbopec()
		case 0xed:
			tCycles = cpu.cboped()
		case 0xf1:
			tCycles = cpu.cbopf1()
		case 0xf2:
			tCycles = cpu.cbopf2()
		case 0xf3:
			tCycles = cpu.cbopf3()
		case 0xf4:
			tCycles = cpu.cbopf4()
		case 0xf5:
			tCycles = cpu.cbopf5()
		case 0xf6:
			tCycles = cpu.cbopf6()
		case 0xf9:
			tCycles = cpu.cbopf9()
		case 0xfa:
			tCycles = cpu.cbopfa()
		case 0xfb:
			tCycles = cpu.cbopfb()
		case 0xfc:
			tCycles = cpu.cbopfc()
		case 0x42:
			tCycles = cpu.cbop42()
		case 0x43:
			tCycles = cpu.cbop43()
		case 0x44:
			tCycles = cpu.cbop44()
		case 0x45:
			tCycles = cpu.cbop45()
		case 0x4a:
			tCycles = cpu.cbop4a()
		case 0x4b:
			tCycles = cpu.cbop4b()
		case 0x4c:
			tCycles = cpu.cbop4c()
		case 0x80:
			tCycles = cpu.cbop80()
		case 0x82:
			tCycles = cpu.cbop82()
		case 0x83:
			tCycles = cpu.cbop83()
		case 0x84:
			tCycles = cpu.cbop84()
		case 0x85:
			tCycles = cpu.cbop85()
		case 0x88:
			tCycles = cpu.cbop88()
		case 0x8a:
			tCycles = cpu.cbop8a()
		case 0x8b:
			tCycles = cpu.cbop8b()
		case 0x8c:
			tCycles = cpu.cbop8c()
		case 0x8d:
			tCycles = cpu.cbop8d()

		default:
			panic(cpulogger.Error(fmt.Sprintf("[CB] Opcode 0x%x is not implemented. PC=0x%x", opcode, cpu.Registers.PC-1)))
		}
	} else {
		// unprefixed
		switch opcode {
		case 0x00:
			tCycles = cpu.op00()
		case 0xc3:
			tCycles = cpu.opc3()
		case 0xaf:
			tCycles = cpu.opaf()
		case 0x21:
			tCycles = cpu.op21()
		case 0x0e:
			tCycles = cpu.op0e()
		case 0x06:
			tCycles = cpu.op06()
		case 0x32:
			tCycles = cpu.op32()
		case 0x05:
			tCycles = cpu.op05()
		case 0x20:
			tCycles = cpu.op20()
		case 0x0d:
			tCycles = cpu.op0d()
		case 0x3e:
			tCycles = cpu.op3e()
		case 0xf3:
			tCycles = cpu.opf3()
		case 0xe0:
			tCycles = cpu.ope0()
		case 0xf0:
			tCycles = cpu.opf0()
		case 0xfe:
			tCycles = cpu.opfe()
		case 0x36:
			tCycles = cpu.op36()
		case 0xea:
			tCycles = cpu.opea()
		case 0x31:
			tCycles = cpu.op31()
		case 0x2a:
			tCycles = cpu.op2a()
		case 0xe2:
			tCycles = cpu.ope2()
		case 0x0c:
			tCycles = cpu.op0c()
		case 0xcd:
			tCycles = cpu.opcd()
		case 0x01:
			tCycles = cpu.op01()
		case 0x0b:
			tCycles = cpu.op0b()
		case 0x78:
			tCycles = cpu.op78()
		case 0xb1:
			tCycles = cpu.opb1()
		case 0xc9:
			tCycles = cpu.opc9()
		case 0xfb:
			tCycles = cpu.opfb()
		case 0xf5:
			tCycles = cpu.opf5()
		case 0xc5:
			tCycles = cpu.opc5()
		case 0xd5:
			tCycles = cpu.opd5()
		case 0xe5:
			tCycles = cpu.ope5()
		case 0xa7:
			tCycles = cpu.opa7()
		case 0x28:
			tCycles = cpu.op28()
		case 0xc0:
			tCycles = cpu.opc0()
		case 0xfa:
			tCycles = cpu.opfa()
		case 0xc8:
			tCycles = cpu.opc8()
		case 0x3d:
			tCycles = cpu.op3d()
		case 0x34:
			tCycles = cpu.op34()
		case 0x3c:
			tCycles = cpu.op3c()
		case 0xe1:
			tCycles = cpu.ope1()
		case 0xd1:
			tCycles = cpu.opd1()
		case 0xc1:
			tCycles = cpu.opc1()
		case 0xf1:
			tCycles = cpu.opf1()
		case 0xd9:
			tCycles = cpu.opd9()
		case 0x2f:
			tCycles = cpu.op2f()
		case 0xe6:
			tCycles = cpu.ope6()
		case 0x47:
			tCycles = cpu.op47()
		case 0xb0:
			tCycles = cpu.opb0()
		case 0x4f:
			tCycles = cpu.op4f()
		case 0xa9:
			tCycles = cpu.opa9()
		case 0xa1:
			tCycles = cpu.opa1()
		case 0x79:
			tCycles = cpu.op79()
		case 0xef:
			tCycles = cpu.opef()
		case 0x87:
			tCycles = cpu.op87()
		case 0x5f:
			tCycles = cpu.op5f()
		case 0x16:
			tCycles = cpu.op16()
		case 0x19:
			tCycles = cpu.op19()
		case 0x5e:
			tCycles = cpu.op5e()
		case 0x23:
			tCycles = cpu.op23()
		case 0x56:
			tCycles = cpu.op56()
		case 0xe9:
			tCycles = cpu.ope9()
		case 0x11:
			tCycles = cpu.op11()
		case 0x12:
			tCycles = cpu.op12()
		case 0x13:
			tCycles = cpu.op13()
		case 0x1a:
			tCycles = cpu.op1a()
		case 0x22:
			tCycles = cpu.op22()
		case 0x7c:
			tCycles = cpu.op7c()
		case 0x1c:
			tCycles = cpu.op1c()
		case 0xca:
			tCycles = cpu.opca()
		case 0x7e:
			tCycles = cpu.op7e()
		case 0x18:
			tCycles = cpu.op18()
		case 0x2d:
			tCycles = cpu.op2d()
		case 0x3a:
			tCycles = cpu.op3a()
		case 0x57:
			tCycles = cpu.op57()
		case 0x7b:
			tCycles = cpu.op7b()
		case 0x7a:
			tCycles = cpu.op7a()
		case 0x0a:
			tCycles = cpu.op0a()
		case 0x7d:
			tCycles = cpu.op7d()
		case 0xc6:
			tCycles = cpu.opc6()
		case 0x6f:
			tCycles = cpu.op6f()
		case 0x5d:
			tCycles = cpu.op5d()
		case 0x54:
			tCycles = cpu.op54()
		case 0x2c:
			tCycles = cpu.op2c()
		case 0x09:
			tCycles = cpu.op09()
		case 0xf6:
			tCycles = cpu.opf6()
		case 0x35:
			tCycles = cpu.op35()
		case 0x30:
			tCycles = cpu.op30()
		case 0x6b:
			tCycles = cpu.op6b()
		case 0x02:
			tCycles = cpu.op02()
		case 0x77:
			tCycles = cpu.op77()
		case 0x03:
			tCycles = cpu.op03()
		case 0x9b:
			tCycles = cpu.op9b()
		case 0xda:
			tCycles = cpu.opda()
		case 0x07:
			tCycles = cpu.op07()
		case 0x67:
			tCycles = cpu.op67()
		case 0x4e:
			tCycles = cpu.op4e()
		case 0x46:
			tCycles = cpu.op46()
		case 0x69:
			tCycles = cpu.op69()
		case 0x60:
			tCycles = cpu.op60()
		case 0x85:
			tCycles = cpu.op85()
		case 0xc2:
			tCycles = cpu.opc2()
		case 0x73:
			tCycles = cpu.op73()
		case 0x72:
			tCycles = cpu.op72()
		case 0x71:
			tCycles = cpu.op71()
		case 0x1e:
			tCycles = cpu.op1e()
		case 0x62:
			tCycles = cpu.op62()
		case 0x40:
			tCycles = cpu.op40()
		// ----done tetris----
		// ----start tennis---
		case 0x26:
			tCycles = cpu.op26()
		case 0x95:
			tCycles = cpu.op95()
		case 0xcf:
			tCycles = cpu.opcf()
		case 0x66:
			tCycles = cpu.op66()
		case 0x81:
			tCycles = cpu.op81()
		case 0x29:
			tCycles = cpu.op29()
		case 0x76:
			tCycles = cpu.op76()
		case 0x80:
			tCycles = cpu.op80()
		/////
		case 0x04:
			tCycles = cpu.op04()
		case 0x0f:
			tCycles = cpu.op0f()
		case 0xb2:
			tCycles = cpu.opb2()
		case 0xd2:
			tCycles = cpu.opd2()
		case 0xb3:
			tCycles = cpu.opb3()
		case 0xa0:
			tCycles = cpu.opa0()
		case 0x2e:
			tCycles = cpu.op2e()
		case 0xce:
			tCycles = cpu.opce()
		case 0xb8:
			tCycles = cpu.opb8()
		case 0x38:
			tCycles = cpu.op38()
		case 0xcc:
			tCycles = cpu.opcc()
		case 0xee:
			tCycles = cpu.opee()
		case 0xdf:
			tCycles = cpu.opdf()
		case 0xd6:
			tCycles = cpu.opd6()
		case 0x14:
			tCycles = cpu.op14()
		case 0x83:
			tCycles = cpu.op83()
		case 0xd8:
			tCycles = cpu.opd8()
		case 0xbe:
			tCycles = cpu.opbe()
		case 0x91:
			tCycles = cpu.op91()
		case 0x90:
			tCycles = cpu.op90()
		case 0xd0:
			tCycles = cpu.opd0()
		case 0x89:
			tCycles = cpu.op89()
		case 0x6c:
			tCycles = cpu.op6c()
		case 0x61:
			tCycles = cpu.op61()
		case 0x17:
			tCycles = cpu.op17()
		case 0xb9:
			tCycles = cpu.opb9()
		case 0x93:
			tCycles = cpu.op93()
		case 0xae:
			tCycles = cpu.opae()
		case 0x9a:
			tCycles = cpu.op9a()
		case 0x25:
			tCycles = cpu.op25()
		case 0x96:
			tCycles = cpu.op96()
		case 0x9e:
			tCycles = cpu.op9e()
		case 0x8c:
			tCycles = cpu.op8c()
		case 0xde:
			tCycles = cpu.opde()
		case 0x99:
			tCycles = cpu.op99()
		case 0x98:
			tCycles = cpu.op98()
		case 0x86:
			tCycles = cpu.op86()
		case 0x8e:
			tCycles = cpu.op8e()
		case 0x9c:
			tCycles = cpu.op9c()
		case 0x24:
			tCycles = cpu.op24()
		case 0xbd:
			tCycles = cpu.opbd()
		case 0x3f:
			tCycles = cpu.op3f()
		case 0x51:
			tCycles = cpu.op51()
		case 0x1b:
			tCycles = cpu.op1b()
			//////// joypad
		case 0xc4:
			tCycles = cpu.opc4()
		case 0xdc:
			tCycles = cpu.opdc()
		case 0x27:
			tCycles = cpu.op27()

		case 0xff:
			tCycles = cpu.opff()
		case 0x1f:
			tCycles = cpu.op1f()
		//////////////
		case 0x33:
			tCycles = cpu.op33()
		case 0x15:
			tCycles = cpu.op15()
		case 0x1d:
			tCycles = cpu.op1d()
		case 0x2b:
			tCycles = cpu.op2b()
		case 0x3b:
			tCycles = cpu.op3b()
		case 0x39:
			tCycles = cpu.op39()
		case 0x82:
			tCycles = cpu.op82()
		case 0x84:
			tCycles = cpu.op84()
		case 0xe8:
			tCycles = cpu.ope8()
		case 0x92:
			tCycles = cpu.op92()
		case 0x94:
			tCycles = cpu.op94()
		case 0x97:
			tCycles = cpu.op97()
		case 0xa2:
			tCycles = cpu.opa2()
		case 0xa3:
			tCycles = cpu.opa3()
		case 0xa4:
			tCycles = cpu.opa4()
		case 0xa5:
			tCycles = cpu.opa5()
		case 0xa6:
			tCycles = cpu.opa6()
		case 0xb4:
			tCycles = cpu.opb4()
		case 0xb5:
			tCycles = cpu.opb5()
		case 0xb6:
			tCycles = cpu.opb6()
		case 0xb7:
			tCycles = cpu.opb7()
		case 0xa8:
			tCycles = cpu.opa8()
		case 0xaa:
			tCycles = cpu.opaa()
		case 0xab:
			tCycles = cpu.opab()
		case 0xac:
			tCycles = cpu.opac()
		case 0xad:
			tCycles = cpu.opad()
		case 0x88:
			tCycles = cpu.op88()
		case 0x8a:
			tCycles = cpu.op8a()
		case 0x8b:
			tCycles = cpu.op8b()
		case 0x8d:
			tCycles = cpu.op8d()
		case 0x8f:
			tCycles = cpu.op8f()
		case 0x9d:
			tCycles = cpu.op9d()
		case 0x9f:
			tCycles = cpu.op9f()
		case 0xba:
			tCycles = cpu.opba()
		case 0xbb:
			tCycles = cpu.opbb()
		case 0xbc:
			tCycles = cpu.opbc()
		case 0xbf:
			tCycles = cpu.opbf()
		case 0x10:
			tCycles = cpu.op10()
		case 0xd4:
			tCycles = cpu.opd4()
		case 0x37:
			tCycles = cpu.op37()
		case 0xf9:
			tCycles = cpu.opf9()
		case 0x08:
			tCycles = cpu.op08()
		case 0xf8:
			tCycles = cpu.opf8()
		case 0x70:
			tCycles = cpu.op70()
		case 0x74:
			tCycles = cpu.op74()
		case 0x75:
			tCycles = cpu.op75()
		case 0x6e:
			tCycles = cpu.op6e()
		case 0x7f:
			tCycles = cpu.op7f()
		case 0x41:
			tCycles = cpu.op41()
		case 0x42:
			tCycles = cpu.op42()
		case 0x43:
			tCycles = cpu.op43()
		case 0x44:
			tCycles = cpu.op44()
		case 0x45:
			tCycles = cpu.op45()
		case 0x48:
			tCycles = cpu.op48()
		case 0x49:
			tCycles = cpu.op49()
		case 0x4a:
			tCycles = cpu.op4a()
		case 0x4b:
			tCycles = cpu.op4b()
		case 0x4c:
			tCycles = cpu.op4c()
		case 0x4d:
			tCycles = cpu.op4d()
		case 0x5b:
			tCycles = cpu.op5b()
		///////////////////////
		case 0x50:
			tCycles = cpu.op50()
		case 0x52:
			tCycles = cpu.op52()
		case 0x53:
			tCycles = cpu.op53()
		case 0x55:
			tCycles = cpu.op55()
		case 0x58:
			tCycles = cpu.op58()
		case 0x59:
			tCycles = cpu.op59()
		case 0x5a:
			tCycles = cpu.op5a()
		case 0x5c:
			tCycles = cpu.op5c()
		case 0x63:
			tCycles = cpu.op63()
		case 0x64:
			tCycles = cpu.op64()
		case 0x65:
			tCycles = cpu.op65()
		case 0x68:
			tCycles = cpu.op68()
		case 0x6a:
			tCycles = cpu.op6a()
		case 0x6d:
			tCycles = cpu.op6d()
		case 0xf2:
			tCycles = cpu.opf2()
		case 0xc7:
			tCycles = cpu.opc7()
		case 0xd7:
			tCycles = cpu.opd7()
		case 0xe7:
			tCycles = cpu.ope7()
		case 0xf7:
			tCycles = cpu.opf7()
		// ILLEGALS
		case 0xd3:
			tCycles = cpu.opd3()
		case 0xdb:
			tCycles = cpu.opdb()
		case 0xdd:
			tCycles = cpu.opdd()
		case 0xe3:
			tCycles = cpu.ope3()
		case 0xe4:
			tCycles = cpu.ope4()
		case 0xeb:
			tCycles = cpu.opeb()
		case 0xec:
			tCycles = cpu.opec()
		case 0xed:
			tCycles = cpu.oped()
		case 0xf4:
			tCycles = cpu.opf4()
		case 0xfc:
			tCycles = cpu.opfc()
		case 0xfd:
			tCycles = cpu.opfd()

		default:
			panic(cpulogger.Error(fmt.Sprintf("Opcode 0x%x is not implemented. PC=0x%x", opcode, cpu.Registers.PC-1)))
		}

	}

	if tCycles < 0 {
		if isPrefixed {
			panic(cpulogger.Error(fmt.Sprintf("[CB] Opcode 0x%x did not change tCycles", opcode)))
		} else {
			panic(cpulogger.Error(fmt.Sprintf("Opcode 0x%x did not change tCycles", opcode)))
		}
	}

	return tCycles
}

func (cpu *CPU) loadROMFile(cartridge *Cartridge) {
	cpu.Cartridge = cartridge
	//copy(cpu.Memory[:], cartridge.ROMdata)
	copy(cpu.Memory[0x0000:0x8000], cartridge.ROMdata[:0x8000])

	//cartridge.bootROM

	cpu.Registers.PC = 0x100
	cpu.Registers.setAF(0x01B0)
	cpu.Registers.setBC(0x0013)
	cpu.Registers.setDE(0x00D8)
	cpu.Registers.setHL(0x014D)
	cpu.Registers.SP = 0xFFFE

	cpu.graphics = NewGraphics(cpu)

	cpu.Memory[0xFF00] = 0xCF // Joypad
	cpu.Memory[0xFF01] = 0x00 // Serial transfer data
	cpu.Memory[0xFF02] = 0x7E // Serial transfer control

	// Timer and divider
	cpu.Memory[0xFF04] = 0x18 // DIV - Divider register
	cpu.Memory[0xFF05] = 0x00 // TIMA: Timer counter
	cpu.Memory[0xFF06] = 0x00 // TMA: Timer modulo
	cpu.Memory[0xFF07] = 0x00 // TAC: Timer control

	cpu.Memory[0xFF0F] = 0xE1 // IF

	//Audio
	cpu.Memory[0xFF10] = 0x80 //NR10: Channel 1 sweep
	cpu.Memory[0xFF11] = 0xBF //NR11: Channel 1 length timer & duty cycle
	cpu.Memory[0xFF12] = 0xF3 // NR12: Channel 1 volume & envelope
	cpu.Memory[0xFF14] = 0xBF // NR14: Channel 1 period high & control
	cpu.Memory[0xFF16] = 0x3F //NR21 ($FF16) → NR11
	cpu.Memory[0xFF17] = 0x00 //NR22 ($FF17) → NR12
	cpu.Memory[0xFF19] = 0xBF // NR24 ($FF19) → NR14
	cpu.Memory[0xFF1A] = 0x7F //NR30: Channel 3 DAC enable
	cpu.Memory[0xFF1B] = 0xFF // NR31: Channel 3 length timer [write-only]
	cpu.Memory[0xFF1C] = 0x9F // NR32: Channel 3 output level
	cpu.Memory[0xFF1E] = 0xBF // NR34: Channel 3 period high & control
	cpu.Memory[0xFF20] = 0xFF //NR41: Channel 4 length timer [write-only]
	cpu.Memory[0xFF21] = 0x00 // NR42: Channel 4 volume & envelope
	cpu.Memory[0xFF22] = 0x00 // NR43: Channel 4 frequency & randomness
	cpu.Memory[0xFF23] = 0xBF // NR44: Channel 4 control
	cpu.Memory[0xFF24] = 0x77 // NR50: Master volume & VIN panning
	cpu.Memory[0xFF25] = 0xF3 // NR51: Sound panning
	cpu.Memory[0xFF26] = 0xF1 //NR52: Audio master control

	cpu.Memory[0xFF40] = 0b10010001 // LCDC

	cpu.Memory[0xFF41] = 0b00000001 //STAT

	cpu.Memory[0xFF42] = 0x00 //SCY
	cpu.Memory[0xFF43] = 0x00 // SCX

	cpu.Memory[0xFF44] = 0 //first scanline, LY

	cpu.Memory[0xFF45] = 0x00 //LYC
	cpu.Memory[0xFF46] = 0xFF //DMA
	cpu.Memory[0xFF47] = 0xFC //BGP
	cpu.Memory[0xFF48] = 0xFF //OBP0
	cpu.Memory[0xFF49] = 0xFF // OBP1

	cpu.Memory[0xFF50] = 0x00 // BOOT ROM OFF

	cpu.Memory[0xFF4A] = 0x00 // WY
	cpu.Memory[0xFF4B] = 0x00 // WX

	cpu.Memory[0xFFFF] = 0x00 //IE
}
func (cpu *CPU) frameSteps() {
	const cyclesPerFrame = 70224
	cyclesCurrFrame := 0
	for cyclesCurrFrame < cyclesPerFrame {
		tCycles := cpu.execOpcodes()
		cpu.graphics.modesHandling(tCycles)
		cpu.checkSchedule()

		cpu.handleInterruptions()
		cpu.timer.Update(tCycles, cpu)
		cpu.joypad.UpdateJoypad()
		cyclesCurrFrame += tCycles
	}

}
