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
	Registers Registers
	Memory    [65536]uint8
	timer     Timer
	Cartridge *Cartridge
	graphics  *Graphics
	//OpcodesTable map[string]map[string][]map[string]string
	IME          bool // interrupt master enable
	IMEScheduled bool //enable IME after one instr
	halted       bool
	stopped      bool
	IE           uint8 // FFFF — IE: Interrupt enable
	IF           uint8 //FF0F — IF: Interrupt flag
}

type Registers struct {
	A, B, C, D, E, F, H, L uint8
	SP, PC                 uint16

	// PC- 0x100 init
	//AF, BC, DE, HL uint16
}

// flags
const flagZ uint8 = 1 << 7 // zero flag
const flagN uint8 = 1 << 6 // sub flag
const flagH uint8 = 1 << 5 // half carry flag
const flagC uint8 = 1 << 4 // carry flag

func NewCPU() *CPU {
	cpulogger = log.New("cpu")

	cpu := &CPU{}
	cpu.IME = false

	//cpu.graphics.cpu = cpu

	return cpu

}

func (register *Registers) setFlag(flag uint8, on bool) {
	//var register Registers
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
	//var register Registers
	return uint16(register.A)<<8 | uint16(register.F)
}

func (register *Registers) setAF(value uint16) {
	//var register Registers
	register.A = uint8(value >> 8)
	register.F = uint8(value & 0xFF)
}

// B-hi C-lo

func (register *Registers) getBC() uint16 {
	//var register Registers
	//			most significant	least significant
	return uint16(register.B)<<8 | uint16(register.C)
}

func (register *Registers) setBC(value uint16) {
	//var register Registers
	register.B = uint8(value >> 8)   //take upper bits
	register.C = uint8(value & 0xFF) // take lower bits
}

// D-hi E-lo

func (register *Registers) getDE() uint16 {
	//var register Registers
	return uint16(register.D)<<8 | uint16(register.E)
}

func (register *Registers) setDE(value uint16) {
	//var register Registers
	register.D = uint8(value >> 8)   //take upper bits
	register.E = uint8(value & 0xFF) // take lower bits
}

// H-hi L-lo

func (register *Registers) getHL() uint16 {
	//var register Registers
	return uint16(register.H)<<8 | uint16(register.L)
}

func (register *Registers) setHL(value uint16) {
	//var register Registers
	register.H = uint8(value >> 8)   //take upper bits
	register.L = uint8(value & 0xFF) // take lower bits
}
func (cpu *CPU) checkSchedule() {
	if cpu.IMEScheduled {
		cpu.IME = true
		cpu.IMEScheduled = false
	}
}

func (cpu *CPU) fetchOpcode() uint8 {
	//opcode := cpu.Memory[cpu.Registers.PC]
	opcode := cpu.memoryRead(cpu.Registers.PC)
	cpu.Registers.PC++
	return opcode
}

func (cpu *CPU) push(n uint16) {
	hi := (n & 0xFF00) >> 8
	lo := n & 0xFF
	cpu.Registers.SP -= 2

	if cpu.Registers.SP == 0xFF80 {
		panic(cpulogger.Error("Stack smash"))
	}

	cpu.memoryWrite(cpu.Registers.SP+1, uint8(hi))
	cpu.memoryWrite(cpu.Registers.SP, uint8(lo))
}

func (cpu *CPU) pop() uint16 {
	lo := uint16(cpu.memoryRead(cpu.Registers.SP))
	hi := uint16(cpu.memoryRead(cpu.Registers.SP + 1))
	cpu.Registers.SP += 2

	if cpu.Registers.SP == 0xFFFE+1 {
		panic(cpulogger.Error("Stack smash"))
	}
	return hi<<8 | lo
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
	IE := cpu.memoryRead(0xFFFF)
	IF := cpu.memoryRead(0xFF0F)

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
		//	refs
		//  https://gbdev.gg8.se/forums/viewtopic.php?id=845
		//  https://gbdev.io/pandocs/Joypad_Input.html#ff00--p1joyp-joypad
		//TODO: write just to bits 5,4 for joypad - p1
		current := cpu.Memory[address]
		cpu.Memory[address] = current&0b11001111 | value&0b00110000
		cpulogger.Debug(fmt.Sprintf("write in 0xFF00 %08b", cpu.Memory[address]))
	} else if address == 0xFF89 {
		cpulogger.Debug(fmt.Sprintf("write in 0xFF89 %08b", value))
		cpu.Memory[address] = value
	} else if address == 0x6C60 {
		cpulogger.Debug(fmt.Sprintf("write in 0x6C60 %x", value))
		cpu.Memory[address] = value
	} else if address >= 0xC000 && address <= 0xCFFF {
		cpu.Memory[address] = value
		cpu.Memory[address+0x2000] = value // ?
	} else if address >= 0xE000 && address <= 0xDDFF {
		cpu.Memory[address] = value
		cpu.Memory[address-0x2000] = value
	} else if address >= 0xFF04 && address <= 0xFF07 {
		cpu.timer.Write(address, value)
	} else if address >= VRAM_START && address <= VRAM_END {
		//cpu.graphics.writeVRAM(address, value)
		cpu.Memory[address] = value
	} else if address >= OAM_START && address <= OAM_END {
		cpu.Memory[address] = value
	} else if address == 0xFF46 {
		cpu.Memory[address] = value
		cpu.dmaTransfer(value)
	} else if address == 0xFF40 {
		cpulogger.Debug(fmt.Sprintf("!!LCDC WRITE: 0x%02X\n", value))
		cpu.Memory[address] = value
	} else {
		cpu.Memory[address] = value
	}
}
func (cpu *CPU) memoryRead(address uint16) byte {
	return cpu.Memory[address]
}
func (cpu *CPU) execOpcodes() int {
	if cpu.halted {
		return 0
	}
	if cpu.stopped {
		return 0
	}
	var tCycles int = -1
	prefix := cpu.fetchOpcode()

	isPrefixed := (prefix == 0xcb)
	opcode := prefix

	cpulogger.Debug(fmt.Sprintf("Executing opcode 0x%x @PC=0x%x A=0x%x F=0x%x DE=0x%x HL=0x%x", opcode, cpu.Registers.PC-1, cpu.Registers.A, cpu.Registers.F, cpu.Registers.getDE(), cpu.Registers.getHL()))
	if prefix == 0xcb {
		// prefixed
		opcode = cpu.fetchOpcode()
		cpulogger.Debug(fmt.Sprintf("Executing in cbprefixed opcode 0x%x @PC=0x%x A=0x%x F=0x%x ", opcode, cpu.Registers.PC-1, cpu.Registers.A, cpu.Registers.F))
		// TODO: ROM01:6C61	cp $00 (fe)-> should be cb before but it's not
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
		////case 0xff:
		////	tCycles = cpu.opff()
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
	copy(cpu.Memory[:], cartridge.ROMdata)
	//cartridge.bootROM // TODO urgent

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

	cpu.Memory[0xFF4A] = 0x00 // WY
	cpu.Memory[0xFF4B] = 0x00 // WX

	cpu.Memory[0xFFFF] = 0x00 //IE
}
func (cpu *CPU) frameSteps() {
	const cyclesPerFrame = 70224
	cyclesCurrFrame := 0
	for cyclesCurrFrame < cyclesPerFrame {
		tCycles := cpu.execOpcodes()
		cpu.handleInterruptions()
		cpu.checkSchedule()
		cpu.timer.Update(tCycles, cpu)
		cpu.graphics.modesHandling(tCycles)
		cyclesCurrFrame += tCycles
	}
}
