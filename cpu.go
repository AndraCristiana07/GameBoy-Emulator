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

//	type Memory struct {
//		ROM  []byte // 0x0000 - 0x3FFF <-> 0x4000-0x7FFF
//		VRAM []byte // 0x8000 - 0x9FFF
//		ERAM []byte // 0xA000 - 0xBFFF
//		WRAM []byte // 0xC000 - 0xDFFF
//		OAM  []byte // 0xFE00 - 0xFE9F
//		IO   []byte // 0xFF00 - 0xFF7F
//		HRAM []byte // 0xFF80 - 0xFFFE
//		IE   byte   // 0xFFFF
//	}
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
	//// cpulogger.Debug(fmt.Sprintf(("cpu.grahics", cpu.graphics))
	//// cpulogger.Debug(fmt.Sprintf(("cpu.grahics.cpu", cpu.graphics.cpu))
	cpu.Registers.PC = 0x100
	cpu.Registers.setAF(0x01B0)
	cpu.Registers.setBC(0x0013)
	cpu.Registers.setDE(0x00D8)
	cpu.Registers.setHL(0x014D)
	cpu.Registers.SP = 0xFFFE

	cpu.graphics = NewGraphics(cpu)

	//cpu.Memory[0xFF40] = 0x91
	cpu.Memory[0xFF40] = 0b10010001 // LCDC
	cpu.Memory[0xFF40] |= 1 << 7    //LCDC enable

	cpu.Memory[0xFF41] = 0b00000001 //STAT

	cpu.Memory[0xFF42] = 0x00 //SCY
	cpu.Memory[0xFF43] = 0x00 // SCX

	cpu.Memory[0xFF44] = 0 //first scanline, LY

	cpu.Memory[0xFF45] = 0x00 //LYC
	cpu.Memory[0xFF47] = 0xFC //BGP
	cpu.Memory[0xFF48] = 0xFF //OBP0
	cpu.Memory[0xFF49] = 0xFF // OBP1

	cpu.Memory[0xFF4A] = 0x00 // WY
	cpu.Memory[0xFF4B] = 0x00 // WX

	cpu.Memory[0xFFFF] = 0x00 //IE

	// Timer and divider
	cpu.Memory[0xFF05] = 0x00 // TIMA: Timer counter
	cpu.Memory[0xFF06] = 0x00 // TMA: Timer modulo
	cpu.Memory[0xFF07] = 0x00 // TAC: Timer control

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

func (cpu *CPU) getImmediate8() uint8 {
	//val := cpu.Memory[cpu.Registers.PC]
	val := cpu.memoryRead(cpu.Registers.PC)
	cpu.Registers.PC++
	cpulogger.Debug(fmt.Sprintf("immediate 8 val: 0x%04X, in memory: 0x%02X  pc now: 0x%04X ", val, cpu.Memory[val], cpu.Registers.PC))

	return val
}

func (cpu *CPU) getImmediate16() uint16 {
	//val := cpu.Memory[cpu.Registers.PC]
	//val1 := cpu.Memory[cpu.Registers.PC]
	//val2 := cpu.Memory[cpu.Registers.PC+1]
	val1 := cpu.memoryRead(cpu.Registers.PC)
	val2 := cpu.memoryRead(cpu.Registers.PC + 1)
	res := uint16(val1) | uint16(val2)<<8
	cpu.Registers.PC += 2
	cpulogger.Debug(fmt.Sprintf("immediate 16 val: 0x%04X, in memory: 0x%02X val1: 0x%02X and val2: 0x%02X pc now: 0x%04X", res, cpu.Memory[res], val1, val2, cpu.Registers.PC))
	return res
}

func (cpu *CPU) inc(reg uint8) uint8 {
	cpulogger.Debug("INC")

	//Set if carry from bit 3
	halfCarry := (reg & 0x0F) == 0x0F
	cpulogger.Debug(fmt.Sprintf("inc reg val: 0x%04X", reg))
	res := reg + 1
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, halfCarry)
	cpulogger.Debug(fmt.Sprintf("inc reg val+1: 0x%04X", res))
	return uint8(res)

}

func (cpu *CPU) dec(reg uint8) uint8 {
	cpulogger.Debug("DEC")
	res := reg - 1

	//Set if no borrow from bit 4
	halfCarry := (reg & 0x0F) == 0x00
	cpu.Registers.setFlag(flagZ, res == 0)
	cpu.Registers.setFlag(flagN, true)
	cpu.Registers.setFlag(flagH, halfCarry)
	return res

}

func (cpu *CPU) addA(reg uint8) uint8 {
	cpulogger.Debug("ADD")

	a := cpu.Registers.A
	res := uint16(a) + uint16(reg)
	halfCarry := ((a & 0x0F) + (reg & 0x0F)) > 0x0F //if lower nibble overflows (carry from bit 3 -> 4)
	carry := res > 0xFF                             // if res overflows 8 bits
	cpu.Registers.setFlag(flagZ, uint8(res) == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, halfCarry)
	cpu.Registers.setFlag(flagC, carry)
	return uint8(res & 0xFF)
}

func (cpu *CPU) adcA(reg uint8) uint8 {
	cpulogger.Debug("ADC")

	a := cpu.Registers.A
	carry := uint8(0)
	if cpu.Registers.getFlag(flagC) {
		carry = uint8(1)
	}

	res := uint16(a) + uint16(reg) + uint16(carry)
	halfCarry := ((a & 0x0F) + (reg & 0x0F) + carry) > 0x0F //if lower nibble overflows (carry from bit 3 -> 4)
	carryFlag := res > 0xFF                                 // if res overflows 8 bits
	cpu.Registers.setFlag(flagZ, uint8(res) == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, halfCarry)
	cpu.Registers.setFlag(flagC, carryFlag)
	return uint8(res & 0xFF)
}

func (cpu *CPU) subA(reg uint8) uint8 {
	cpulogger.Debug("SUB")

	a := cpu.Registers.A
	res := uint16(a) - uint16(reg)
	halfCarry := (a & 0x0F) < (reg & 0x0F) // if no borrow from bit 4
	carryFlag := res < 0                   // if no borrow
	cpu.Registers.setFlag(flagZ, uint8(res) == 0)
	cpu.Registers.setFlag(flagN, true)
	cpu.Registers.setFlag(flagH, halfCarry)
	cpu.Registers.setFlag(flagC, carryFlag)
	return uint8(res & 0xFF)
}

func (cpu *CPU) sbcA(reg uint8) uint8 {
	cpulogger.Debug("SBC")

	carry := uint8(0)
	if cpu.Registers.getFlag(flagC) {
		carry = uint8(1)
	}
	a := cpu.Registers.A
	res := uint16(a) - uint16(reg) - uint16(carry)
	halfCarry := (a&0x0F)-(reg&0x0F)-carry < 0 // if no borrow from bit 4
	carryFlag := res < 0                       // if no borrow
	cpu.Registers.setFlag(flagZ, uint8(res) == 0)
	cpu.Registers.setFlag(flagN, true)
	cpu.Registers.setFlag(flagH, halfCarry)
	cpu.Registers.setFlag(flagC, carryFlag)
	return uint8(res & 0xFF)
}

func (cpu *CPU) cpA(reg uint8) {
	cpulogger.Debug("CP")

	a := cpu.Registers.A
	res := uint16(a) - uint16(reg)
	halfCarry := (a & 0x0F) < (reg & 0x0F) // if no borrow from bit 4
	carryFlag := res < 0                   // if no borrow
	cpu.Registers.setFlag(flagZ, uint8(res) == 0)
	cpu.Registers.setFlag(flagN, true)
	cpu.Registers.setFlag(flagH, halfCarry)
	cpu.Registers.setFlag(flagC, carryFlag)

}

func (cpu *CPU) fetchOpcode() uint8 {
	//opcode := cpu.Memory[cpu.Registers.PC]
	opcode := cpu.memoryRead(cpu.Registers.PC)
	cpu.Registers.PC++
	return opcode
}

func (cpu *CPU) fetchCBOpcode() uint8 {
	//cbOpcode := cpu.Memory[cpu.Registers.PC]
	cbOpcode := cpu.memoryRead(cpu.Registers.PC)
	cpulogger.Debug(fmt.Sprintf("CB Opcode in fetch: 0x%02X", cbOpcode))
	cpu.Registers.PC++
	return cbOpcode
}

func (cpu *CPU) push(n uint16) {
	cpulogger.Debug("PUSH")
	hi := (n & 0xFF00) >> 8
	lo := n & 0xFF
	cpu.Registers.SP -= 2

	//cpu.Memory[cpu.Registers.SP+1] = uint8(hi)
	//cpu.Memory[cpu.Registers.SP] = uint8(lo)
	cpu.memoryWrite(cpu.Registers.SP+1, uint8(hi))
	cpu.memoryWrite(cpu.Registers.SP, uint8(lo))
}

func (cpu *CPU) pop() uint16 {
	cpulogger.Debug("POP")
	//lo := uint16(cpu.Memory[cpu.Registers.SP])
	//hi := uint16(cpu.Memory[cpu.Registers.SP+1])
	lo := uint16(cpu.memoryRead(cpu.Registers.SP))
	hi := uint16(cpu.memoryRead(cpu.Registers.SP + 1))
	//// cpulogger.Debug(fmt.Sprintf("POP -> lower: 0x%02X and highrr: 0x%02X\n", lo, hi)
	cpu.Registers.SP += 2

	//res := uint16(hi)<<8 | uint16(lo)
	//n = res
	return hi<<8 | lo

}

func (cpu *CPU) execRST(address uint16) {
	cpulogger.Debug("RST")
	// Push present address onto stack.
	// Jump to address $0000 + n.
	cpu.Registers.SP -= 2
	hi := (cpu.Registers.PC & 0xFF00) >> 8
	lo := cpu.Registers.PC & 0xFF
	//cpulogger.Debug(fmt.Sprintf("RST- SP before: 0x%04X\n", cpu.Registers.SP)
	//cpu.Memory[cpu.Registers.SP] = uint8(cpu.Registers.PC & 0xFF) //lower byte
	//cpu.memoryWrite(cpu.Registers.SP, uint8(cpu.Registers.PC&0xFF)) //lower byte
	cpu.memoryWrite(cpu.Registers.SP, uint8(lo)) //lower byte

	cpu.memoryWrite(cpu.Registers.SP+1, uint8(hi)) //upper byte

	cpu.Registers.PC = address
	cpulogger.Debug(fmt.Sprintf("PC now: 0x%04X and in memory 0x%02X\n", cpu.Registers.PC, cpu.Memory[cpu.Registers.PC]))
	//cpulogger.Debug(fmt.Sprintf("Pc after  0x%04X sp = 0x%04X\n", cpu.Registers.PC, cpu.Registers.SP)

}

func (cpu *CPU) execBIT(reg uint8, bit uint8) {
	//Test bit b in register r
	cpulogger.Debug("BIT")

	res := reg & (1 << bit)
	cpu.Registers.setFlag(flagZ, res == 0x00)
	cpu.Registers.setFlag(flagN, false) //reset
	cpu.Registers.setFlag(flagH, true)  //set
}

func (cpu *CPU) execBITHL(bit uint8) {
	//Test bit b in register r
	cpulogger.Debug("BIT HL")

	//res := cpu.Memory[cpu.Registers.getHL()] & (1 << bit)
	res := cpu.memoryRead(cpu.Registers.getHL()) & (1 << bit)

	cpu.Registers.setFlag(flagZ, res == 0x00)
	cpu.Registers.setFlag(flagN, false) //reset
	cpu.Registers.setFlag(flagH, true)  //set
}

func (cpu *CPU) execSET(reg uint8, bit uint8) uint8 {
	//Set bit b in register r
	cpulogger.Debug("SET")

	return reg | (1 << bit)
}

func (cpu *CPU) execSETHL(bit uint8) {
	//Set bit b in register r
	cpulogger.Debug("SET HL")

	//cpu.Memory[cpu.Registers.getHL()] = cpu.Memory[cpu.Registers.getHL()] | (1 << bit)
	cpu.memoryWrite(cpu.Registers.getHL(), cpu.Memory[cpu.Registers.getHL()]|(1<<bit))
}

func (cpu *CPU) execRES(reg uint8, bit uint8) uint8 {
	//Reset bit b in register r
	cpulogger.Debug("RES")

	return reg & ^(1 << bit)
}

func (cpu *CPU) execRESHL(bit uint8) {
	//Reset bit b in register r
	cpulogger.Debug("RES HL")

	addr := cpu.Registers.getHL()
	val := cpu.memoryRead(addr)
	//cpu.Memory[cpu.Registers.getHL()] = cpu.Memory[cpu.Registers.getHL()] & ^(1 << bit)
	cpu.memoryWrite(addr, val & ^(1<<bit))
}

func (cpu *CPU) execSWAP(reg uint8) uint8 {
	cpulogger.Debug("SWAP")

	//Swap upper & lower nibles of n
	reg = (reg >> 4) | ((reg & 0x0F) << 4)
	cpu.Registers.setFlag(flagZ, reg == 0x00)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, false)
	return reg
}

func (cpu *CPU) execSWAPHL() uint8 {
	cpulogger.Debug("SWAP HL")

	//Swap upper & lower nibles of n
	//cpu.Memory[cpu.Registers.getHL()] = (cpu.Memory[cpu.Registers.getHL()] >> 4) | (cpu.Memory[cpu.Registers.getHL()] << 4)
	val := cpu.memoryRead(cpu.Registers.getHL())
	res := val>>4 | (val&0x0F)<<4
	cpu.memoryWrite(cpu.Registers.getHL(), res) //??
	cpu.Registers.setFlag(flagZ, res == 0x00)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	cpu.Registers.setFlag(flagC, false)
	return res
}

func (cpu *CPU) execSLA(reg uint8) uint8 {
	//Shift n left into Carry. LSB of n set to 0.
	cpulogger.Debug("SLA")

	carry := (reg & 0x80) >> 7
	res := reg << 1
	cpu.Registers.setFlag(flagC, carry == 0x01)
	cpu.Registers.setFlag(flagZ, res == 0x00)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	return res
}

func (cpu *CPU) execSLAHL() {
	//Shift n left into Carry. LSB of n set to 0.
	cpulogger.Debug("SLA HL")

	//carry := (cpu.Memory[cpu.Registers.getHL()] & 0x80) >> 7
	carry := (cpu.memoryRead(cpu.Registers.getHL()) & 0x80) >> 7
	//cpu.Memory[cpu.Registers.getHL()] = cpu.Memory[cpu.Registers.getHL()] << 1
	//addr := cpu.memoryRead(cpu.Registers.getHL())
	addr := cpu.Registers.getHL()
	val := cpu.memoryRead(addr) << 1
	cpu.memoryWrite(addr, val)
	cpu.Registers.setFlag(flagC, carry == 0x01)
	//cpu.Registers.setFlag(flagZ, cpu.Memory[cpu.Registers.getHL()] == 0x00)
	cpu.Registers.setFlag(flagZ, val == 0x00)

	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
}

func (cpu *CPU) execSRA(reg uint8) uint8 {
	//Shift n right into Carry. MSB doesn't change.
	cpulogger.Debug("SRA")

	carry := reg & 0x01
	reg = (reg >> 1) | (reg & 0x80)
	cpu.Registers.setFlag(flagC, carry == 0x01)
	cpu.Registers.setFlag(flagZ, reg == 0x00)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	return reg
}

func (cpu *CPU) execSRAHL() {
	//Shift n right into Carry. MSB doesn't change.
	cpulogger.Debug("SRA HL")
	hl := cpu.memoryRead(cpu.Registers.getHL())
	carry := hl & 0x01
	//cpu.Memory[cpu.Registers.getHL()] = (cpu.Memory[cpu.Registers.getHL()] >> 1) | (cpu.Memory[cpu.Registers.getHL()] & 0x80)
	reg := (hl >> 1) | (hl & 0x80)
	cpu.memoryWrite(cpu.Registers.getHL(), reg)
	cpu.Registers.setFlag(flagC, carry == 0x01)
	cpu.Registers.setFlag(flagZ, cpu.Memory[cpu.Registers.getHL()] == 0x00)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
}
func (cpu *CPU) execSRL(reg uint8) uint8 {
	//Shift n right into Carry. MSB set to 0.
	cpulogger.Debug("SRL")

	carry := reg & 0x01
	reg = reg >> 1
	cpu.Registers.setFlag(flagC, carry == 0x01)
	cpu.Registers.setFlag(flagZ, reg == 0x00)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	return reg
}

func (cpu *CPU) execSRLHL() {
	//Shift n right into Carry. MSB set to 0.
	cpulogger.Debug("SRL HL")
	hl := cpu.memoryRead(cpu.Registers.getHL())
	carry := hl & 0x01
	//cpu.Memory[cpu.Registers.getHL()] = cpu.Memory[cpu.Registers.getHL()] >> 1
	reg := hl >> 1
	cpu.memoryWrite(cpu.Registers.getHL(), reg)
	cpu.Registers.setFlag(flagC, carry == 0x01)
	cpu.Registers.setFlag(flagZ, cpu.Memory[cpu.Registers.getHL()] == 0x00)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
}

func (cpu *CPU) execRLC(reg uint8) uint8 {
	//Rotate n left. Old bit 7 to Carry flag
	cpulogger.Debug("RLC")
	carry := (reg & 0x80) >> 7
	reg = reg << 1
	cpu.Registers.setFlag(flagC, carry == 0x00)
	cpu.Registers.setFlag(flagZ, reg == 0x00)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	return reg

}

func (cpu *CPU) execRLCHL() {
	cpulogger.Debug("RLC HL")

	//Rotate n left. Old bit 7 to Carry flag
	carry := (cpu.Memory[cpu.Registers.getHL()] & 0x80) >> 7
	//cpu.Memory[cpu.Registers.getHL()] = cpu.Memory[cpu.Registers.getHL()] << 1
	addr := cpu.memoryRead(cpu.Registers.getHL())
	val := addr << 1
	cpu.memoryWrite(cpu.Registers.getHL(), val)
	cpu.Registers.setFlag(flagC, carry == 0x00)
	cpu.Registers.setFlag(flagZ, cpu.Memory[cpu.Registers.getHL()] == 0x00)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)

}

func (cpu *CPU) execRL(reg uint8) uint8 {
	cpulogger.Debug("RL")
	// Rotate n left through Carry flag.

	oldCarry := uint8(0)
	if cpu.Registers.getFlag(flagC) {
		oldCarry = 1
	}
	newCarry := reg & 0x80 //store bit 7
	reg = (reg << 1) | oldCarry
	cpu.Registers.setFlag(flagC, newCarry == 0x01)
	cpu.Registers.setFlag(flagZ, reg == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
	return reg
}

func (cpu *CPU) execRLHL() {
	cpulogger.Debug("RL HL")
	// Rotate n left through Carry flag.

	oldCarry := uint8(0)
	if cpu.Registers.getFlag(flagC) {
		oldCarry = 1
	}
	hl := cpu.memoryRead(cpu.Registers.getHL())
	newCarry := cpu.memoryRead(cpu.Registers.getHL()) & 0x80 //store bit 7
	reg := (hl << 1) | oldCarry
	cpu.memoryWrite(cpu.Registers.getHL(), reg)
	//cpu.Registers.setHL(uint16((cpu.Memory[cpu.Registers.getHL()] << 1) | oldCarry))
	cpu.Registers.setFlag(flagC, newCarry == 0x01)
	cpu.Registers.setFlag(flagZ, cpu.Memory[cpu.Registers.getHL()] == 0)
	cpu.Registers.setFlag(flagN, false)
	cpu.Registers.setFlag(flagH, false)
}

func (cpu *CPU) execRRC(reg uint8) uint8 {
	cpulogger.Debug("RRC")
	// Rotate n right. Old bit 0 to Carry flag
	c := reg&0x01 == 0x01
	reg >>= 1
	cpu.Registers.setFlag(flagZ, reg == 0)
	cpu.Registers.setFlag(flagC, c)
	cpu.Registers.setFlag(flagN, false) //reset
	cpu.Registers.setFlag(flagH, false) //reset
	return reg
}

func (cpu *CPU) execRRCHL() {
	cpulogger.Debug("RRC HL")
	// Rotate n right. Old bit 0 to Carry flag
	//c := cpu.Memory[cpu.Registers.getHL()]&0x01 == 0x01
	hl := cpu.memoryRead(cpu.Registers.getHL())
	carry := hl&0x01 == 0x01
	reg := hl >> 1
	cpu.memoryWrite(cpu.Registers.getHL(), reg)
	//cpu.Memory[cpu.Registers.getHL()] >>= 1
	cpu.Registers.setFlag(flagZ, hl == 0)
	cpu.Registers.setFlag(flagC, carry)
	cpu.Registers.setFlag(flagN, false) //reset
	cpu.Registers.setFlag(flagH, false) //reset
}

func (cpu *CPU) execRR(reg uint8) uint8 {
	cpulogger.Debug("RR")
	//Rotate n right through Carry flag.
	oldCarry := uint8(0)
	if cpu.Registers.getFlag(flagC) {
		oldCarry = 1
	}
	newCarry := reg & 0x01 //store bit 0
	reg = (reg >> 1) | oldCarry<<7
	cpu.Registers.setFlag(flagZ, reg == 0)
	cpu.Registers.setFlag(flagC, newCarry == 0x01)
	cpu.Registers.setFlag(flagN, false) //reset
	cpu.Registers.setFlag(flagH, false) //reset
	return reg
}

func (cpu *CPU) execRRHL() {
	cpulogger.Debug("RR HL")
	//Rotate n right through Carry flag.
	oldCarry := uint8(0)
	if cpu.Registers.getFlag(flagC) {
		oldCarry = 1
	}
	hl := cpu.memoryRead(cpu.Registers.getHL())

	newCarry := hl & 0x01 //store bit 0

	reg := (hl >> 1) | oldCarry<<7
	cpu.memoryWrite(cpu.Registers.getHL(), reg)
	//cpu.Registers.setHL(uint16((hl >> 1) | oldCarry<<7))
	cpu.Registers.setFlag(flagZ, hl == 0)
	cpu.Registers.setFlag(flagC, newCarry == 0x01)
	cpu.Registers.setFlag(flagN, false) //reset
	cpu.Registers.setFlag(flagH, false) //reset
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
	//cpu.IE := cpu.Memory[0xFFF] //Interrupt enable
	//cpu.IF := cpu.Memory[0xF0F] //Interrupt flag
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
	//TODO: add more if they exist
	if address >= 0xC000 && address <= 0xCFFF {
		if address == 0xC010 {
			cpulogger.Debug(fmt.Sprintf("writing to 0xC010 value: 0x%02X", value))
		}
		cpu.Memory[address] = value
		cpu.Memory[address+0x2000] = value //?
	} else if address >= 0xE000 && address <= 0xDDFF {
		cpu.Memory[address] = value
		cpu.Memory[address-0x2000] = value
	} else if address >= 0xFF04 && address <= 0xFF07 {
		cpu.timer.Write(address, value)
	} else if address >= VRAM_START && address <= VRAM_END {
		//cpu.graphics.writeVRAM(address, value)
		cpu.Memory[address] = value
		cpulogger.Debug(fmt.Sprintf("VRAM write ->  address: 0x%02X, value: 0b%08b", address, value))
	} else if address >= OAM_START && address <= OAM_END {
		cpulogger.Debug(fmt.Sprintf("OAM write ->  address: 0x%02X, value: 0b%08b", address, value))
		cpu.Memory[address] = value
		//DEBUG
		//} else if address == 0xFF01 {
		//	cpulogger.Debug(fmt.Sprintf("writing to 0xFF01 value: 0x%02X\n", value)
		//} else if address == 0xFF02 {
		//	cpulogger.Debug(fmt.Sprintf("writing to 0xFF02 value: 0x%02X\n", value)
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
	//if address >= 0xFF04 && address <= 0xFF07 {
	//	return cpu.timer.Read(address)
	//} else if address >= VRAM_START && address <= VRAM_END {
	//	return cpu.graphics.readVRAM(address)
	//} else if address >= OAM_START && address <= OAM_END {
	//	return cpu.graphics.readOAM(address - OAM_START)
	//} else if address >= 0xFF40 && address <= 0xFF4B {
	//	return cpu.graphics.getFromMemory(address)
	//} else {
	//if address == 0xFF40 {
	//	return cpu.graphics.getLCDC()
	//} else if address == 0xFF41 {
	//	return cpu.graphics.getSTAT()
	//} else {
	return cpu.Memory[address]
	//}

	//}
}

func (cpu *CPU) execOpcodes() int {
	if cpu.halted {
		return 0
	}
	if cpu.stopped {
		return 0
	}
	//cpu.Memory[0xFF40] |= 1 << 7
	var tCycles int
	// cpulogger.Debug(fmt.Sprintf("Before instructiins => PC: 0x%04X | Memory[0x0039]: 0x%02X\n", cpu.Registers.PC, cpu.Memory[0x0039])

	//cpulogger.Debug(fmt.Sprintf("Executing opcode: 0x%02X\n", opcode)
	opcode := cpu.fetchOpcode()
	cpulogger.Debug("", opcode)
	cpulogger.Debug("", cpu.Registers.PC)
	cpulogger.Debug("", cpu.Registers.SP)
	cpulogger.Debug("Opcode in fetch: 0x%X ; PC now:  0x%02X ; SP now: 0x%02X", opcode, cpu.Registers.PC, cpu.Registers.SP)

	//// cpulogger.Debug(fmt.Sprintf("pc: 0x%04X and opcode: 0x%02X\n", cpu.Registers.PC, opcode)
	switch opcode {

	case 0b1: // 0x01 -> LD BC, imm16
		cpu.Registers.setBC(cpu.getImmediate16())
		tCycles = 12

	case 0b10: // 0x02 -> LD [BC], A
		cpu.memoryWrite(cpu.Registers.getBC(), cpu.Registers.A)
		tCycles = 8

	case 0b110: // 0x06 -> LD B, imm8
		cpu.Registers.B = cpu.getImmediate8()
		tCycles = 8

	case 0b1000: // 0x08 -> LD [imm16], SP
		cpu.memoryWrite(cpu.getImmediate16(), uint8(cpu.Registers.SP))
		tCycles = 20

	case 0b1001: // 0x09 -> ADD HL, BC
		cpu.Registers.setHL(cpu.Registers.getHL() + cpu.Registers.getBC())
		tCycles = 8

	case 0b1010: // 0x0A -> LD A, [BC]
		cpu.Registers.A = cpu.memoryRead(cpu.Registers.getBC())
		tCycles = 8

	case 0b1110: // 0x0E -> LD C, imm8
		cpu.Registers.C = cpu.getImmediate8()
		tCycles = 8

	case 0b10001: // 0x11 -> LD DE, n16
		cpu.Registers.setDE(uint16(cpu.getImmediate16()))
		tCycles = 12

	case 0b10010: // 0x12 -> LD [DE], A
		cpu.memoryWrite(cpu.Registers.getDE(), cpu.Registers.A)
		tCycles = 8

	case 0b10110: // 0x16 -> LD D, imm8
		cpu.Registers.D = cpu.getImmediate8()
		tCycles = 8

	case 0b11001: // 0x19 -> ADD HL, DE
		cpu.Registers.setHL(cpu.Registers.getHL() + cpu.Registers.getDE())
		tCycles = 8

	case 0b11010: // 0x1A -> LD A, [DE]
		cpu.Registers.A = cpu.memoryRead(cpu.Registers.getDE())
		tCycles = 8

	case 0b11110: // 0x1E -> LD E, imm8
		cpu.Registers.E = cpu.getImmediate8()
		tCycles = 8

	case 0b100000: // 0x20 -> JR NZ, imm8
		// If flagZ is not set then add n to current
		// address and jump to it
		n := int8(cpu.getImmediate8())

		//// cpulogger.Debug(fmt.Sprintf("Immediate 8 in 0x20: 0x%04X", n)
		if cpu.Registers.getFlag(flagZ) == false {

			//cpu.Registers.PC += uint16(n)
			cpulogger.Debug("JR NZ taken: ")
			cpulogger.Debug(fmt.Sprintf("PC before jr : 0x%04X", cpu.Registers.PC))
			//cpulogger.Debug(fmt.Sprintf("PC before jr : 0x%04X\n", cpu.Registers.PC)
			cpulogger.Debug(fmt.Sprintf("PC will now be %02X ", uint16(int32(cpu.Registers.PC)+int32(n))))

			cpu.Registers.PC = uint16(int32(cpu.Registers.PC) + int32(n))
			cpulogger.Debug(fmt.Sprintf("PC after jr nz imm : 0x%04X", cpu.Registers.PC))
			tCycles = 12

		} else {
			cpulogger.Debug("JR NZ not taken")
			tCycles = 20

		}
		//// cpulogger.Debug(fmt.Sprintf(string(flagZ))
		//// cpulogger.Debug(fmt.Sprintf(string(n))

	case 0b100001: // 0x21 -> LD HL, n16  // pc 0x247 ? imm 0xff26
		n := cpu.getImmediate16()
		cpu.Registers.setHL(uint16(n))
		// cpucpulogger.Debug("LD HL : immed: ", n)
		tCycles = 12

	case 0b100010: // 0x22 -> LD [HL+], A
		cpu.memoryWrite(cpu.Registers.getHL(), cpu.Registers.A)
		cpu.Registers.setHL(cpu.Registers.getHL() + 1)
		tCycles = 8

	case 0b100110: // 0x26 -> LD H, imm8
		cpu.Registers.H = cpu.getImmediate8()
		tCycles = 8

	case 0b101000: // 0x28 -> JR Z, e8
		// If flagZ is set then add n to current
		// address and jump to it
		//cpulogger.Debug(fmt.Sprintf("pc before jr %04X", cpu.Registers.PC)
		n := int8(cpu.getImmediate8())
		if cpu.Registers.getFlag(flagZ) {
			//cpu.Registers.PC += uint16(n)
			cpu.Registers.PC = uint16(int32(cpu.Registers.PC) + int32(n))
			//cpulogger.Debug(fmt.Sprintf("pc after jr %04X", cpu.Registers.PC)

			tCycles = 12
		} else {
			tCycles = 8
		}
		// cpulogger.Debug(fmt.Sprintf(string(flagZ))
		// cpulogger.Debug(fmt.Sprintf(string(n))

	case 0b101001: // 0x29 -> ADD HL, HL
		cpu.Registers.setHL(cpu.Registers.getHL() + cpu.Registers.getHL())
		tCycles = 8

	case 0b101010: // 0x2A -> LD A, [HL+]
		cpu.Registers.A = cpu.memoryRead(cpu.Registers.getHL())
		cpu.Registers.setHL(cpu.Registers.getHL() + 1)
		tCycles = 8

	case 0b101110: // 0x2E -> LD L, imm8
		cpu.Registers.L = cpu.getImmediate8()
		tCycles = 8

	case 0b110000: // 0x30 -> JR NC, e8
		// If flagC is not set then add n to current
		// address and jump to it
		n := int8(cpu.getImmediate8())
		cpulogger.Debug(fmt.Sprintf("pc before %04X", cpu.Registers.PC))
		if !cpu.Registers.getFlag(flagC) {
			//cpu.Registers.PC += uint16(n)
			cpu.Registers.PC = uint16(int32(cpu.Registers.PC) + int32(n))
			cpulogger.Debug(fmt.Sprintf("pc after %04X", cpu.Registers.PC))

			tCycles = 12
		} else {
			tCycles = 8

		}

	case 0b110001: // 0x31 -> LD SP, n16
		cpu.Registers.SP = uint16(cpu.getImmediate16())
		tCycles = 12

	case 0b110010: // 0x32 -> LD [HL-], A  //??
		//cpulogger.Debug(fmt.Sprintf()
		cpu.memoryWrite(cpu.Registers.getHL(), cpu.Registers.A)

		cpu.Registers.setHL(cpu.Registers.getHL() - 1)
		tCycles = 8

	case 0b110110: // 0x36 -> LD [HL], imm8
		cpu.memoryWrite(cpu.Registers.getHL(), cpu.getImmediate8())
		tCycles = 12

	case 0b111000: // 0x38 -> JR C, e8
		// If flagC is set then add n to current
		//address and jump to it
		n := int8(cpu.getImmediate8())
		if cpu.Registers.getFlag(flagC) {
			//cpu.Registers.PC += uint16(n)
			cpu.Registers.PC = uint16(int32(cpu.Registers.PC) + int32(n))

			tCycles = 12
		} else {
			tCycles = 8
		}

	case 0b111001: // 0x39 -> ADD HL, SP
		cpu.Registers.setHL(cpu.Registers.getHL() + (cpu.Registers.SP))
		tCycles = 8

	case 0b111010: // 0x3A -> LD A, [HL-]
		cpu.Registers.A = cpu.memoryRead(cpu.Registers.getHL())
		cpu.Registers.setHL(cpu.Registers.getHL() - 1)
		tCycles = 8

	case 0b111110: // 0x3E -> LD A, imm8
		cpu.Registers.A = cpu.getImmediate8()
		cpulogger.Debug(fmt.Sprintf("REgisters B %d", cpu.Registers.B))
		cpulogger.Debug(fmt.Sprintf("Register C %d", cpu.Registers.C))
		cpulogger.Debug(fmt.Sprintf("Register BC %d", cpu.Registers.getBC()))
		tCycles = 8

	case 0b1000000: // 0x40 -> LD B, B
		//cpu.Registers.B = cpu.Registers.B
		tCycles = 4

	case 0b1000001: // 0x41 -> LD B, C
		cpu.Registers.B = cpu.Registers.C
		tCycles = 4

	case 0b1000010: // 0x42 -> LD B, D
		cpu.Registers.B = cpu.Registers.D
		tCycles = 4

	case 0b1000011: // 0x43 -> LD B, E
		cpu.Registers.B = cpu.Registers.E
		tCycles = 4

	case 0b1000100: // 0x44 -> LD B, H
		cpu.Registers.B = cpu.Registers.H
		tCycles = 4

	case 0b1000101: // 0x45 -> LD B, L
		cpu.Registers.B = cpu.Registers.L
		tCycles = 4

	case 0b1000110: // 0x46 -> LD B, [HL]
		cpu.Registers.B = cpu.memoryRead(cpu.Registers.getHL())
		tCycles = 8

	case 0b1000111: // 0x47 -> LD B, A
		cpu.Registers.B = cpu.Registers.A
		tCycles = 4

	case 0b1001000: // 0x48 -> LD C, B
		cpu.Registers.C = cpu.Registers.B
		tCycles = 4

	case 0b1001001: // 0x49 -> LD C, C
		//cpu.Registers.C = cpu.Registers.C
		tCycles = 4

	case 0b1001010: // 0x4A -> LD C, D
		cpu.Registers.C = cpu.Registers.D
		tCycles = 4

	case 0b1001011: // 0x4B -> LD C, E
		cpu.Registers.C = cpu.Registers.E
		tCycles = 4

	case 0b1001100: // 0x4C -> LD C, H
		cpu.Registers.C = cpu.Registers.H
		tCycles = 4

	case 0b1001101: // 0x4D -> LD C, L
		cpu.Registers.C = cpu.Registers.L
		tCycles = 4

	case 0b1001110: // 0x4E -> LD C, [HL]
		cpu.Registers.C = cpu.memoryRead(cpu.Registers.getHL())
		tCycles = 8

	case 0b1001111: // 0x4F -> LD C, A
		cpu.Registers.C = cpu.Registers.A
		tCycles = 4

	case 0b1010000: // 0x50 -> LD D, B
		cpu.Registers.D = cpu.Registers.B
		tCycles = 4

	case 0b1010001: // 0x51 -> LD D, C
		cpu.Registers.D = cpu.Registers.C
		tCycles = 4

	case 0b1010010: // 0x52 -> LD D, D
		//cpu.Registers.D = cpu.Registers.D
		tCycles = 4

	case 0b1010011: // 0x53 -> LD D, E
		cpu.Registers.D = cpu.Registers.E
		tCycles = 4

	case 0b1010100: // 0x54 -> LD D, H
		cpu.Registers.D = cpu.Registers.H
		tCycles = 4

	case 0b1010101: // 0x55 -> LD D, L
		cpu.Registers.D = cpu.Registers.L
		tCycles = 4

	case 0b1010110: // 0x56 -> LD D, [HL]
		cpu.Registers.D = cpu.memoryRead(cpu.Registers.getHL())
		tCycles = 8

	case 0b1010111: // 0x57 -> LD D, A
		cpu.Registers.D = cpu.Registers.A
		tCycles = 4

	case 0b1011000: // 0x58 -> LD E, B
		cpu.Registers.E = cpu.Registers.B
		tCycles = 4

	case 0b1011001: // 0x59 -> LD E, C
		cpu.Registers.E = cpu.Registers.C
		tCycles = 4

	case 0b1011010: // 0x5A -> LD E, D
		cpu.Registers.E = cpu.Registers.D
		tCycles = 4

	case 0b1011011: // 0x5B -> LD E, E
		//cpu.Registers.E = cpu.Registers.E
		tCycles = 4

	case 0b1011100: // 0x5C -> LD E, H
		cpu.Registers.E = cpu.Registers.H
		tCycles = 4

	case 0b1011101: // 0x5D -> LD E, L
		cpu.Registers.E = cpu.Registers.L
		tCycles = 4

	case 0b1011110: // 0x5E -> LD E, [HL]
		cpu.Registers.E = cpu.memoryRead(cpu.Registers.getHL())
		tCycles = 8

	case 0b1011111: // 0x5F -> LD E, A
		cpu.Registers.E = cpu.Registers.A
		tCycles = 4

	case 0b1100000: // 0x60 -> LD H, B
		cpu.Registers.H = cpu.Registers.B
		tCycles = 4

	case 0b1100001: // 0x61 -> LD H, C
		cpu.Registers.H = cpu.Registers.C
		tCycles = 4

	case 0b1100010: // 0x62 -> LD H, D
		cpu.Registers.H = cpu.Registers.D
		tCycles = 4

	case 0b1100011: // 0x63 -> LD H, E
		cpu.Registers.H = cpu.Registers.E
		tCycles = 4

	case 0b1100100: // 0x64 -> LD H, H
		tCycles = 4
		break

	case 0b1100101: // 0x65 -> LD H, L
		cpu.Registers.H = cpu.Registers.L
		tCycles = 4

	case 0b1100110: // 0x66 -> LD H, [HL]
		cpu.Registers.H = cpu.memoryRead(cpu.Registers.getHL())
		tCycles = 8

	case 0b1100111: // 0x67 -> LD H, A
		cpu.Registers.H = cpu.Registers.A
		tCycles = 4

	case 0b1101000: // 0x68 -> LD L, B
		cpu.Registers.L = cpu.Registers.B
		tCycles = 4

	case 0b1101001: // 0x69 -> LD L, C
		cpu.Registers.L = cpu.Registers.C
		tCycles = 4

	case 0b1101010: // 0x6A -> LD L, D
		cpu.Registers.L = cpu.Registers.D
		tCycles = 4

	case 0b1101011: // 0x6B -> LD L, E
		cpu.Registers.L = cpu.Registers.E
		tCycles = 4

	case 0b1101100: // 0x6C -> LD L, H
		cpu.Registers.L = cpu.Registers.H
		tCycles = 4

	case 0b1101101: // 0x6D -> LD L, L
		//cpu.Registers.L = cpu.Registers.L
		tCycles = 4

	case 0b1101110: // 0x6E -> LD L, [HL]
		cpu.Registers.L = cpu.memoryRead(cpu.Registers.getHL())
		tCycles = 8

	case 0b1101111: // 0x6F -> LD L, A
		cpu.Registers.L = cpu.Registers.A
		tCycles = 4

	case 0b1110000: // 0x70 -> LD [HL], B
		cpu.memoryWrite(cpu.Registers.getHL(), cpu.Registers.B)
		tCycles = 8

	case 0b1110001: // 0x71 -> LD [HL], C
		cpu.memoryWrite(cpu.Registers.getHL(), cpu.Registers.C)
		tCycles = 8

	case 0b1110010: // 0x72 -> LD [HL], D
		cpu.memoryWrite(cpu.Registers.getHL(), cpu.Registers.D)
		tCycles = 8

	case 0b1110011: // 0x73 -> LD [HL], E
		cpu.memoryWrite(cpu.Registers.getHL(), cpu.Registers.E)
		tCycles = 8

	case 0b1110100: // 0x74 -> LD [HL], H
		cpu.memoryWrite(cpu.Registers.getHL(), cpu.Registers.H)
		tCycles = 8

	case 0b1110101: // 0x75 -> LD [HL], L
		cpu.memoryWrite(cpu.Registers.getHL(), cpu.Registers.L)
		tCycles = 8

	case 0b1110111: // 0x77 -> LD [HL], A
		cpu.memoryWrite(cpu.Registers.getHL(), cpu.Registers.A)
		tCycles = 8

	case 0b1111000: // 0x78 -> LD A, B
		cpu.Registers.A = cpu.Registers.B
		tCycles = 4

	case 0b1111001: // 0x79 -> LD A, C
		cpu.Registers.A = cpu.Registers.C
		tCycles = 4

	case 0b1111010: // 0x7A -> LD A, D
		cpu.Registers.A = cpu.Registers.D
		tCycles = 4

	case 0b1111011: // 0x7B -> LD A, E
		cpu.Registers.A = cpu.Registers.E
		tCycles = 4

	case 0b1111100: // 0x7C -> LD A, H
		cpu.Registers.A = cpu.Registers.H
		tCycles = 4

	case 0b1111101: // 0x7D -> LD A, L
		cpu.Registers.A = cpu.Registers.L
		tCycles = 4

	case 0b1111110: // 0x7E -> LD A, [HL]
		cpu.Registers.A = cpu.memoryRead(cpu.Registers.getHL())
		tCycles = 8

	case 0b1111111: // 0x7F -> LD A, A
		//cpu.Registers.A = cpu.Registers.A
		tCycles = 4

	case 0b10000000: // 0x80 -> ADD A, B
		//cpu.Registers.A += cpu.Registers.B
		//cpu.Registers.setFlag()
		cpu.Registers.A = cpu.addA(cpu.Registers.B)
		tCycles = 4

	case 0b10000001: // 0x81 -> ADD A, C
		cpu.Registers.A = cpu.addA(cpu.Registers.C)

		tCycles = 4

	case 0b10000010: // 0x82 -> ADD A, D
		cpu.Registers.A = cpu.addA(cpu.Registers.D)

		tCycles = 4

	case 0b10000011: // 0x83 -> ADD A, E
		cpu.Registers.A = cpu.addA(cpu.Registers.E)

		tCycles = 4

	case 0b10000100: // 0x84 -> ADD A, H
		cpu.Registers.A = cpu.addA(cpu.Registers.H)

		tCycles = 4

	case 0b10000101: // 0x85 -> ADD A, L
		cpu.Registers.A = cpu.addA(cpu.Registers.L)
		tCycles = 4

	case 0b10000110: // 0x86 -> ADD A, [HL]
		val := cpu.Memory[cpu.Registers.getHL()]
		cpu.Registers.A = cpu.addA(val)
		tCycles = 8

	case 0b10000111: // 0x87 -> ADD A, A
		cpu.Registers.A = cpu.addA(cpu.Registers.A)
		tCycles = 4

	case 0b10001000: // 0x88 -> ADC A, B
		cpu.Registers.A = cpu.adcA(cpu.Registers.B)
		tCycles = 4

	case 0b10001001: // 0x89 -> ADC A, C
		cpu.Registers.A = cpu.adcA(cpu.Registers.C)
		tCycles = 4

	case 0b10001010: // 0x8A -> ADC A, D
		cpu.Registers.A = cpu.adcA(cpu.Registers.D)
		tCycles = 4

	case 0b10001011: // 0x8B -> ADC A, E
		cpu.Registers.A = cpu.adcA(cpu.Registers.E)
		tCycles = 4

	case 0b10001100: // 0x8C -> ADC A, H
		cpu.Registers.A = cpu.adcA(cpu.Registers.H)
		tCycles = 4

	case 0b10001101: // 0x8D -> ADC A, L
		cpu.Registers.A = cpu.adcA(cpu.Registers.L)
		tCycles = 4

	case 0b10001110: // 0x8E -> ADC A, [HL]
		//carry := uint8(1)
		//res := cpu.Registers.A + uint8(cpu.Registers.getHL()) + carry
		val := cpu.Memory[cpu.Registers.getHL()]
		cpu.Registers.A = cpu.adcA(val)
		tCycles = 8

	case 0b10001111: // 0x8F -> ADC A, A
		cpu.Registers.A = cpu.adcA(cpu.Registers.A)
		tCycles = 4

	case 0b10010000: // 0x90 -> SUB A, B
		cpu.Registers.A = cpu.subA(cpu.Registers.B)
		tCycles = 4

	case 0b10010001: // 0x91 -> SUB A, C
		cpu.Registers.A = cpu.subA(cpu.Registers.C)
		tCycles = 4

	case 0b10010010: // 0x92 -> SUB A, D
		cpu.Registers.A = cpu.subA(cpu.Registers.D)
		tCycles = 4

	case 0b10010011: // 0x93 -> SUB A, E
		cpu.Registers.A = cpu.subA(cpu.Registers.E)
		tCycles = 4

	case 0b10010100: // 0x94 -> SUB A, H
		cpu.Registers.A = cpu.subA(cpu.Registers.H)
		tCycles = 4

	case 0b10010101: // 0x95 -> SUB A, L
		cpu.Registers.A = cpu.subA(cpu.Registers.L)
		tCycles = 4

	case 0b10010110: // 0x96 -> SUB A, [HL]
		val := cpu.Memory[cpu.Registers.getHL()]
		cpu.Registers.A = cpu.subA(val)
		tCycles = 8

	case 0b10010111: // 0x97 -> SUB A, A
		cpu.Registers.A = cpu.subA(cpu.Registers.A)
		tCycles = 4

	case 0b10011000: // 0x98 -> SBC A, B
		cpu.Registers.A = cpu.sbcA(cpu.Registers.B)
		tCycles = 4

	case 0b10011001: // 0x99 -> SBC A, C
		cpu.Registers.A = cpu.sbcA(cpu.Registers.C)
		tCycles = 4

	case 0b10011010: // 0x9A -> SBC A, D
		cpu.Registers.A = cpu.sbcA(cpu.Registers.D)
		tCycles = 4

	case 0b10011011: // 0x9B -> SBC A, E
		cpu.Registers.A = cpu.sbcA(cpu.Registers.E)
		tCycles = 4

	case 0b10011100: // 0x9C -> SBC A, H
		cpu.Registers.A = cpu.sbcA(cpu.Registers.H)
		tCycles = 4

	case 0b10011101: // 0x9D -> SBC A, L
		cpu.Registers.A = cpu.sbcA(cpu.Registers.L)
		tCycles = 4

	case 0b10011110: // 0x9E -> SBC A, [HL]
		val := cpu.Memory[cpu.Registers.getHL()]
		cpu.Registers.A = cpu.sbcA(val)
		tCycles = 8

	case 0b10011111: // 0x9F -> SBC A, A
		cpu.Registers.A = cpu.sbcA(cpu.Registers.A)
		tCycles = 4

	case 0b10100000: // 0xA0 -> AND A, B
		cpu.Registers.A &= cpu.Registers.B
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, true)  //set
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10100001: // 0xA1 -> AND A, C
		cpu.Registers.A &= cpu.Registers.C
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, true)  //set
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10100010: // 0xA2 -> AND A, D
		cpu.Registers.A &= cpu.Registers.D
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, true)  //set
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10100011: // 0xA3 -> AND A, E
		cpu.Registers.A &= cpu.Registers.E
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, true)  //set
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10100100: // 0xA4 -> AND A, H
		cpu.Registers.A &= cpu.Registers.H
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, true)  //set
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10100101: // 0xA5 -> AND A, L
		cpu.Registers.A &= cpu.Registers.L
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, true)  //set
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10100110: // 0xA6 -> AND A, [HL]
		cpu.Registers.A &= cpu.Memory[cpu.Registers.getHL()]
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, true)  //set
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 8

	case 0b10100111: // 0xA7 -> AND A, A
		cpu.Registers.A &= cpu.Registers.A
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, true)  //set
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10101000: // 0xA8 -> XOR A, B
		cpu.Registers.A ^= cpu.Registers.B
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10101001: // 0xA9 -> XOR A, C
		cpu.Registers.A ^= cpu.Registers.C
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10101010: // 0xAA -> XOR A, D
		cpu.Registers.A ^= cpu.Registers.D
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10101011: // 0xAB -> XOR A, E
		cpu.Registers.A ^= cpu.Registers.E
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10101100: // 0xAC -> XOR A, H
		cpu.Registers.A ^= cpu.Registers.H
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10101101: // 0xAD -> XOR A, L
		cpu.Registers.A ^= cpu.Registers.L
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10101110: // 0xAE -> XOR A, [HL]
		cpu.Registers.A ^= cpu.Memory[cpu.Registers.getHL()]
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 8

	case 0b10101111: // 0xAF -> XOR A, A
		cpu.Registers.A ^= cpu.Registers.A
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10110000: // 0xB0 -> OR A, B
		cpu.Registers.A |= cpu.Registers.B
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10110001: // 0xB1 -> OR A, C
		cpu.Registers.A |= cpu.Registers.C
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10110010: // 0xB2 -> OR A, D
		cpu.Registers.A |= cpu.Registers.D
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10110011: // 0xB3 -> OR A, E
		cpu.Registers.A |= cpu.Registers.E
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10110100: // 0xB4 -> OR A, H
		cpu.Registers.A |= cpu.Registers.H
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10110101: // 0xB5 -> OR A, L
		cpu.Registers.A |= cpu.Registers.L
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10110110: // 0xB6 -> OR A, [HL]
		cpu.Registers.A |= cpu.Memory[cpu.Registers.getHL()]
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 8

	case 0b10110111: // 0xB7 -> OR A, A
		cpu.Registers.A |= cpu.Registers.A
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 4

	case 0b10111000: // 0xB8 -> CP A, B
		cpu.cpA(cpu.Registers.B)
		tCycles = 4

	case 0b10111001: // 0xB9 -> CP A, C
		cpu.cpA(cpu.Registers.C)
		tCycles = 4

	case 0b10111010: // 0xBA -> CP A, D
		cpu.cpA(cpu.Registers.D)
		tCycles = 4

	case 0b10111011: // 0xBB -> CP A, E
		cpu.cpA(cpu.Registers.E)
		tCycles = 4

	case 0b10111100: // 0xBC -> CP A, H
		cpu.cpA(cpu.Registers.H)
		tCycles = 4

	case 0b10111101: // 0xBD -> CP A, L
		cpu.cpA(cpu.Registers.L)
		tCycles = 4

	case 0b10111110: // 0xBE -> CP A, [HL]
		val := cpu.Memory[cpu.Registers.getHL()]
		cpu.cpA(val)
		tCycles = 8

	case 0b10111111: // 0xBF -> CP A, A
		cpu.cpA(cpu.Registers.A)
		tCycles = 4

	case 0b11000010: // 0xC2 -> JP NZ, imm16
		//Jump to address n if flagZ is not set
		n := cpu.getImmediate16()
		if !cpu.Registers.getFlag(flagZ) {
			cpu.Registers.PC = uint16(n)
			tCycles = 16
		} else {
			tCycles = 12
			cpulogger.Debug("JP not taken")
		}

	case 0b11000100: // 0xC4 -> CALL NZ, imm16
		//Push address of next instruction onto stack and then
		// jump to address nn if flagZ is not set
		n := cpu.getImmediate16()
		if !cpu.Registers.getFlag(flagZ) {
			//cpulogger.Debug(fmt.Sprintf("Pc %04b", cpu.Registers.PC)
			cpu.push(cpu.Registers.PC)
			cpu.Registers.PC = uint16(n)
			tCycles = 24
		} else {
			tCycles = 12
		}

	case 0b11000110: // 0xC6 -> ADD A, imm8
		val := cpu.getImmediate8()
		cpu.Registers.A = cpu.addA(val)
		tCycles = 8

	case 0b11001010: // 0xCA -> JP Z, imm16
		//Jump to address n if flagZ is set
		n := cpu.getImmediate16()
		if cpu.Registers.getFlag(flagZ) {
			cpu.Registers.PC = uint16(n)
			tCycles = 16

		} else {
			tCycles = 12

		}

	case 0b11001100: // 0xCC -> CALL Z, imm16
		//Push address of next instruction onto stack and then
		// jump to address nn if flagZ is  set
		n := cpu.getImmediate16()
		if cpu.Registers.getFlag(flagZ) {
			cpulogger.Debug(fmt.Sprintf("PC %d", &cpu.Registers.PC))
			cpu.push(cpu.Registers.PC)
			cpu.Registers.PC = uint16(n)
			tCycles = 24
		} else {
			tCycles = 12

		}

	case 0b11001110: // 0xCE -> ADC A, imm8
		//carry := uint8(1)
		//res := cpu.Registers.A + cpu.getImmediate8() + carry
		val := cpu.getImmediate8()
		cpu.Registers.A = cpu.adcA(val)
		tCycles = 8

	case 0b11010010: // 0xD2 -> JP NC, imm16
		//Jump to address n if flagC is not set
		n := cpu.getImmediate16()
		if !cpu.Registers.getFlag(flagC) {
			cpu.Registers.PC = uint16(n)
			tCycles = 16
		} else {
			tCycles = 12

		}

	case 0b11010100: // 0xD4 -> CALL NC, imm16
		//Push address of next instruction onto stack and then
		// jump to address nn if flagC is not set
		n := cpu.getImmediate16()
		if !cpu.Registers.getFlag(flagC) {
			cpu.push(cpu.Registers.PC)
			cpu.Registers.PC = uint16(n)
			tCycles = 24
		} else {
			tCycles = 12

		}

	case 0b11010110: // 0xD6 -> SUB A, imm8
		val := cpu.getImmediate8()
		cpu.Registers.A = cpu.subA(val)
		tCycles = 8

	case 0b11011010: // 0xDA -> JP C, imm16
		//Jump to address n if flagC is  set
		n := cpu.getImmediate16()
		if cpu.Registers.getFlag(flagC) {
			cpu.Registers.PC = uint16(n)
			tCycles = 20

		} else {
			tCycles = 8
		}

	case 0b11011100: // 0xDC -> CALL C, imm16
		//Push address of next instruction onto stack and then
		// jump to address nn if flagC is  set
		n := cpu.getImmediate16()
		if cpu.Registers.getFlag(flagC) {
			cpu.push(cpu.Registers.PC)
			cpu.Registers.PC = uint16(n)
			tCycles = 24

		} else {
			tCycles = 12
		}

	case 0b11011110: // 0xDE -> SBC A, imm8
		val := cpu.getImmediate8()
		cpu.Registers.A = cpu.sbcA(val)
		tCycles = 8

	case 0b11100000: // 0xE0 -> LDH [a8], A
		//cpu.memoryWrite(uint16(cpu.getImmediate8())|0xFF00, cpu.Registers.A)
		n := cpu.getImmediate8()
		addr := 0xFF00 + uint16(n)
		cpu.memoryWrite(addr, cpu.Registers.A)
		tCycles = 12

	case 0b11100010: // 0xE2 -> LDH [C], A / LD [$FF00+C], A
		//cpu.memoryWrite(uint16(cpu.Registers.C)|0xFF00, cpu.Registers.A)
		addr := 0xFF00 + uint16(cpu.Registers.C)
		cpu.memoryWrite(addr, cpu.Registers.A)
		tCycles = 8

	case 0b11100110: // 0xE6 -> AND A, imm8
		cpu.Registers.A &= cpu.getImmediate8()
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, true)  //set
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 8

	case 0b11101000: // 0xE8 -> ADD SP, e8
		sp := cpu.Registers.SP
		imm := int8(cpu.getImmediate8())
		res := uint16(int32(sp) + int32(imm))
		//halfCarry := ((int32(sp) & 0x0F) + (int32(imm) & 0x0F)) > 0x0F //if lower nibble overflows (carry from bit 3 -> 4)
		//carry := res > 0xFF                                            // if res overflows 8 bits
		val1 := sp & 0xFF
		val2 := uint16(imm) & 0xFF
		halfCarry := ((val1 & 0x0F) + (val2 & 0x0F)) > 0x0F //if lower nibble overflows (carry from bit 3 -> 4)
		carry := val1+val2 > 0xFF                           // if res overflows 8 bits
		cpu.Registers.SP = res
		cpu.Registers.setFlag(flagZ, false)
		cpu.Registers.setFlag(flagN, false)
		cpu.Registers.setFlag(flagH, halfCarry)
		cpu.Registers.setFlag(flagC, carry)
		tCycles = 16

	case 0b11111000: // 0xF8 -> LD HL, SP+e8
		n := int8(cpu.getImmediate8())
		sp := cpu.Registers.SP
		res := uint16(int32(sp) + int32(n))
		val1 := sp & 0xFF
		val2 := uint16(n) & 0xFF
		halfCarry := ((val1 & 0x0F) + (val2 & 0x0F)) > 0x0F //if lower nibble overflows (carry from bit 3 -> 4)
		carry := val1+val2 > 0xFF                           // if res overflows 8 bits
		cpu.Registers.setHL(res)
		cpu.Registers.setFlag(flagZ, false)
		cpu.Registers.setFlag(flagN, false)
		cpu.Registers.setFlag(flagH, halfCarry)
		cpu.Registers.setFlag(flagC, carry)
		tCycles = 12

	case 0b11101010: // 0xEA -> LD [imm16], A
		cpu.memoryWrite(cpu.getImmediate16(), cpu.Registers.A)
		tCycles = 16

	case 0b11101110: // 0xEE -> XOR A, imm8
		cpu.Registers.A ^= cpu.getImmediate8()
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 8

	case 0b11110000: // 0xF0 -> LDH A, [a8]
		cpulogger.Debug(fmt.Sprintf("REgister A before: %d", cpu.Registers.A))
		n := cpu.getImmediate8()
		addr := 0xFF00 + uint16(n)
		cpu.Registers.A = cpu.memoryRead(addr)
		//cpu.Registers.A = cpu.memoryRead(uint16(cpu.getImmediate8()) | 0xFF00)
		cpulogger.Debug(fmt.Sprintf("REgister A after: %d", cpu.Registers.A))

		tCycles = 12

	case 0b11110010: // 0xF2 -> LDH A, [C]
		addr := 0xFF00 + uint16(cpu.Registers.C)
		//cpu.Registers.A = cpu.memoryRead(uint16(cpu.Registers.C) | 0xFF00)
		cpu.Registers.A = cpu.memoryRead(addr)
		tCycles = 8

	case 0b11110110: // 0xF6 -> OR A, imm8
		cpu.Registers.A |= cpu.getImmediate8()
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		cpu.Registers.setFlag(flagC, false) //reset
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		tCycles = 8

	case 0b11111001: // 0xF9 -> LD SP, HL
		cpu.Registers.SP = cpu.Registers.getHL()
		tCycles = 8

	case 0b11111010: // 0xFA -> LD A, [imm16]
		cpu.Registers.A = cpu.memoryRead(cpu.getImmediate16())
		tCycles = 16

	case 0b11111110: // 0xFE -> CP A, imm8
		val := cpu.getImmediate8()
		cpu.cpA(val)
		// one operand cases
		tCycles = 8

	case 0b11: // 0x03 -> INC BC
		value := cpu.Registers.getBC()
		value++
		cpu.Registers.setBC(value)
		tCycles = 8

	case 0b100: // 0x04 -> INC B
		//flags := map[string]string{"Z": "Z", "N": "0", "H": "H", "C": "-"}
		//cpu.Registers.B++
		cpu.Registers.B = cpu.inc(cpu.Registers.B)

		tCycles = 4

	case 0b101: // 0x05 -> DEC B
		//flags := map[string]string{"Z": "Z", "N": "1", "H": "H", "C": "-"}
		//cpu.Registers.B--
		cpulogger.Debug(fmt.Sprintf("Register B before %d", cpu.Registers.B))
		//cpulogger.Debug(fmt.Sprintf("carry flag in  dec before: %v", cpu.Registers.getFlag(flagC))
		cpu.Registers.B = cpu.dec(cpu.Registers.B)
		cpulogger.Debug(fmt.Sprintf("Register B after %0d", cpu.Registers.B))
		//cpulogger.Debug(fmt.Sprintf("carry flag in  dec after: %v", cpu.Registers.getFlag(flagC))

		tCycles = 4

	case 0b1011: // 0x0B -> DEC BC
		value := cpu.Registers.getBC()
		value--
		cpu.Registers.setBC(value)
		tCycles = 8

	case 0b1100: // 0x0C -> INC C
		//flags := map[string]string{"Z": "Z", "N": "0", "H": "H", "C": "-"}
		//cpu.Registers.C++
		cpu.Registers.C = cpu.inc(cpu.Registers.C)
		tCycles = 4

	case 0b1101: // 0x0D -> DEC C
		//flags := map[string]string{"Z": "Z", "N": "1", "H": "H", "C": "-"}
		//cpu.Registers.C--
		cpulogger.Debug(fmt.Sprintf("REgister C before dec is %d", cpu.Registers.C))
		cpulogger.Debug(fmt.Sprintf("Flag Z before is: %t", cpu.Registers.getFlag(flagZ)))
		cpu.Registers.C = cpu.dec(cpu.Registers.C)
		cpulogger.Debug(fmt.Sprintf("REgister C after dec is %d", cpu.Registers.C))
		cpulogger.Debug(fmt.Sprintf("Flag Z after is: %t", cpu.Registers.getFlag(flagZ)))

		tCycles = 4

	case 0b10000: // 0x10 -> STOP imm8
		//TODO
		//Halt CPU & LCD display until button pressed
		imm := cpu.getImmediate8()
		if imm != 0x00 {
			cpulogger.Debug(fmt.Sprintf("STOP -> imm is 0x%02X", imm))
		}
		cpu.Registers.PC++
		//opcode = cpu.fetchOpcode()
		//cpu.stopped = true
		tCycles = 4

	case 0b10011: // 0x13 -> INC DE
		value := cpu.Registers.getDE()
		value++
		cpu.Registers.setDE(value)
		tCycles = 8

	case 0b10100: // 0x14 -> INC D  //? 255 -> 0 ?
		//flags := map[string]string{"Z": "Z", "N": "0", "H": "H", "C": "-"}
		//cpu.Registers.D++
		cpu.Registers.D = cpu.inc(cpu.Registers.D)
		cpulogger.Debug("Flag Z", cpu.Registers.getFlag(flagZ))
		cpulogger.Debug("Flag N", cpu.Registers.getFlag(flagN))
		cpulogger.Debug("Flag H", cpu.Registers.getFlag(flagH))
		tCycles = 4

	case 0b10101: // 0x15 -> DEC D
		//flags := map[string]string{"Z": "Z", "N": "1", "H": "H", "C": "-"}
		//cpu.Registers.D--
		cpu.Registers.D = cpu.dec(cpu.Registers.D)
		tCycles = 4

	case 0b11000: // 0x18 -> JR e8
		n := int8(cpu.getImmediate8())
		// cpulogger.Debug(fmt.Sprintf("Immediate 8 at 0x18: 0x%04X", n)
		//cpu.Registers.PC += uint16(n)
		cpu.Registers.PC = uint16(int32(cpu.Registers.PC) + int32(n))

		tCycles = 12

	case 0b11011: // 0x1B -> DEC DE
		value := cpu.Registers.getDE()
		value--
		cpu.Registers.setDE(value)
		tCycles = 8

	case 0b11100: // 0x1C -> INC E
		//flags := map[string]string{"Z": "Z", "N": "0", "H": "H", "C": "-"}
		//cpu.Registers.E++
		cpu.Registers.E = cpu.inc(cpu.Registers.E)
		tCycles = 4

	case 0b11101: // 0x1D -> DEC E
		//flags := map[string]string{"Z": "Z", "N": "1", "H": "H", "C": "-"}
		//cpu.Registers.E--
		cpu.Registers.E = cpu.dec(cpu.Registers.E)
		tCycles = 4

	case 0b100011: // 0x23 -> INC HL
		value := cpu.Registers.getHL()
		value++
		cpu.Registers.setHL(value)
		tCycles = 8

	case 0b100100: // 0x24 -> INC H
		//flags := map[string]string{"Z": "Z", "N": "0", "H": "H", "C": "-"}
		//cpu.Registers.H++
		cpu.Registers.H = cpu.inc(cpu.Registers.H)
		tCycles = 4

	case 0b100101: // 0x25 -> DEC H
		//flags := map[string]string{"Z": "Z", "N": "1", "H": "H", "C": "-"}
		//cpu.Registers.H--
		cpu.Registers.H = cpu.dec(cpu.Registers.H)
		tCycles = 4

	case 0b101011: // 0x2B -> DEC HL
		value := cpu.Registers.getHL()
		value--
		cpu.Registers.setHL(value)
		tCycles = 8

	case 0b101100: // 0x2C -> INC L
		//flags := map[string]string{"Z": "Z", "N": "0", "H": "H", "C": "-"}
		//cpu.Registers.L++
		cpu.Registers.L = cpu.inc(cpu.Registers.L)
		tCycles = 4

	case 0b101101: // 0x2D -> DEC L
		//flags := map[string]string{"Z": "Z", "N": "1", "H": "H", "C": "-"}
		//cpu.Registers.L--
		cpu.Registers.L = cpu.dec(cpu.Registers.L)
		tCycles = 4

	case 0b110011: // 0x33 -> INC SP
		cpu.Registers.SP++
		tCycles = 8

	case 0b110100: // 0x34 -> INC [HL]
		cpu.Memory[cpu.Registers.getHL()]++
		tCycles = 12

	case 0b110101: // 0x35 -> DEC [HL]
		cpu.Memory[cpu.Registers.getHL()]--
		tCycles = 12

	case 0b111011: // 0x3B -> DEC SP
		cpu.Registers.SP--
		tCycles = 8

	case 0b111100: // 0x3C -> INC A
		//flags := map[string]string{"Z": "Z", "N": "0", "H": "H", "C": "-"}
		//cpu.Registers.A++
		cpu.Registers.A = cpu.inc(cpu.Registers.A)
		tCycles = 4

	case 0b111101: // 0x3D -> DEC A
		//flags := map[string]string{"Z": "Z", "N": "1", "H": "H", "C": "-"}
		//cpu.Registers.A--
		cpu.Registers.A = cpu.dec(cpu.Registers.A)
		tCycles = 4

	case 0b11000000: // 0xC0 -> RET NZ

		//Pop two bytes from stack & jump to that address if flagZ not set
		if !cpu.Registers.getFlag(flagZ) {
			cpu.Registers.PC = cpu.pop()
			tCycles = 20

		} else {
			tCycles = 8
		}

	case 0b11000001: // 0xC1 -> POP BC
		//cpu.pop(cpu.Registers.getBC())
		cpu.Registers.setBC(cpu.pop())
		tCycles = 12

	case 0b11000011: // 0xC3 -> JP imm16
		n := cpu.getImmediate16()
		// cpulogger.Debug(fmt.Sprintf("Immediate value 0x%02X at PC: 0x%02X\n", n, cpu.Registers.PC)
		cpu.Registers.PC = uint16(n)
		tCycles = 16

	case 0b11000101: // 0xC5 -> PUSH BC
		cpu.push(cpu.Registers.getBC())
		tCycles = 16

	case 0b11000111: // 0xC7 -> RST $00
		cpu.execRST(0x00)
		tCycles = 16

	case 0b11001000: // 0xC8 -> RET Z
		//Pop two bytes from stack & jump to that address if flagZ set
		if cpu.Registers.getFlag(flagZ) {
			cpu.Registers.PC = cpu.pop()
			tCycles = 20
		} else {
			tCycles = 8
		}

	case 0b11001101: // 0xCD -> CALL imm16
		n := cpu.getImmediate16()
		cpu.push(cpu.Registers.PC)
		cpu.Registers.PC = uint16(n)
		tCycles = 24

	case 0b11001111: // 0xCF -> RST $08
		cpu.execRST(0x08)
		tCycles = 16

	case 0b11010000: // 0xD0 -> RET NC
		//Pop two bytes from stack & jump to that address if flagC not set
		if !cpu.Registers.getFlag(flagC) {
			cpu.Registers.PC = cpu.pop()
			tCycles = 20
		} else {
			tCycles = 8
		}

	case 0b11010001: // 0xD1 -> POP DE
		cpu.Registers.setDE(cpu.pop())
		tCycles = 12

	case 0b11010101: // 0xD5 -> PUSH DE
		cpu.push(cpu.Registers.getDE())
		tCycles = 16

	case 0b11010111: // 0xD7 -> RST $10
		cpu.execRST(0x10)
		tCycles = 16

	case 0b11011000: // 0xD8 -> RET C
		//Pop two bytes from stack & jump to that address if flagC set
		if cpu.Registers.getFlag(flagC) {
			cpu.Registers.PC = cpu.pop()
			tCycles = 20

		} else {
			tCycles = 8

		}

	case 0b11011111: // 0xDF -> RST $18
		cpu.execRST(0x18)
		tCycles = 16

	case 0b11100001: // 0xE1 -> POP HL
		cpu.Registers.setHL(cpu.pop())
		tCycles = 12

	case 0b11100101: // 0xE5 -> PUSH HL
		cpu.push(cpu.Registers.getHL())
		tCycles = 16

	case 0b11100111: // 0xE7 -> RST $20
		cpu.execRST(0x20)
		tCycles = 16

	case 0b11101001: // 0xE9 -> JP HL
		n := cpu.Registers.getHL()
		cpu.Registers.PC = n
		tCycles = 4

	case 0b11101111: // 0xEF -> RST $28
		cpu.execRST(0x28)
		tCycles = 16

	case 0b11110001: // 0xF1 -> POP AF
		cpu.Registers.setAF(cpu.pop())
		tCycles = 12

	case 0b11110101: // 0xF5 -> PUSH AF
		cpu.push(cpu.Registers.getAF())
		tCycles = 16

	case 0b11110111: // 0xF7 -> RST $30
		cpu.execRST(0x30)
		tCycles = 16

	case 0b11111111: // 0xFF -> RST $38
		cpu.execRST(0x38)
		// cases with 0 operands
		tCycles = 16

	case 0b0: // 0x00 -> NOP
		tCycles = 4

	case 0b111: // 0x07 -> RLCA

		// Rotate A left. Old bit 7 to Carry flag
		c := (cpu.Registers.A&0x80)>>7 == 0x01
		cpu.Registers.setFlag(flagC, c)
		cpu.Registers.A <<= 1
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		cpu.Registers.setFlag(flagN, false)
		cpu.Registers.setFlag(flagH, false)
		tCycles = 4

	case 0b1111: // 0x0F -> RRCA

		// Rotate A right. Old bit 0 to Carry flag.
		c := cpu.Registers.A&0x01 == 0x01
		cpu.Registers.A >>= 1
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		cpu.Registers.setFlag(flagC, c)
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		tCycles = 4

	case 0b10111: // 0x17 -> RLA
		// Rotate A left through Carry flag.

		oldCarry := uint8(0)
		if cpu.Registers.getFlag(flagC) {
			oldCarry = 1
		}
		newCarry := cpu.Registers.A & 0x80 //store bit 7
		cpu.Registers.A = (cpu.Registers.A << 1) | oldCarry
		cpu.Registers.setFlag(flagC, newCarry == 0x01)
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		cpu.Registers.setFlag(flagN, false)
		cpu.Registers.setFlag(flagH, false)
		tCycles = 4

	case 0b11111: // 0x1F -> RRA
		// Rotate A right through Carry flag.

		oldCarry := uint8(0)
		if cpu.Registers.getFlag(flagC) {
			oldCarry = 1
		}
		newCarry := cpu.Registers.A & 0x01 //store bit 0
		cpu.Registers.A = (cpu.Registers.A >> 1) | oldCarry<<7
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		cpu.Registers.setFlag(flagC, newCarry == 0x01)
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		tCycles = 4

	case 0b100111: // 0x27 -> DAA
		// Decimal adjust register A.
		// This instruction adjusts register A so that the
		// correct representation of Binary Coded Decimal (BCD)
		// is obtained.
		lo := cpu.Registers.A % 10
		hi := ((cpu.Registers.A - lo) % 100) / 10
		cpu.Registers.A = (hi << 4) | lo
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		cpu.Registers.setFlag(flagH, false) //reset
		cpu.Registers.setFlag(flagC, true)  //set because flags["C"] == "C" ?!?
		tCycles = 4

	case 0b101111: // 0x2F -> CPL
		//Complement A register. (Flip all bits.)
		cpu.Registers.A = ^cpu.Registers.A
		cpu.Registers.setFlag(flagN, true) // set
		cpu.Registers.setFlag(flagH, true) //set
		tCycles = 4

	case 0b110111: // 0x37 -> SCF
		// Set Carry flag.
		cpu.Registers.setFlag(flagC, true)  //set
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		tCycles = 4

	case 0b111111: // 0x3F -> CCF
		//Complement carry flag.
		// If C flag is set, then reset it.
		// If C flag is reset, then set it.
		if !cpu.Registers.getFlag(flagC) {
			cpu.Registers.setFlag(flagC, true)
		} else {
			cpu.Registers.setFlag(flagC, false)
		}
		cpu.Registers.setFlag(flagN, false) //reset
		cpu.Registers.setFlag(flagH, false) //reset
		tCycles = 4

	case 0b1110110: // 0x76 -> HALT

		// Power down CPU until an interrupt occurs. Use this
		// when ever possible to reduce energy consumption.
		cpu.halted = true
		tCycles = 4

	case 0b11001001: // 0xC9 -> RET
		cpu.Registers.PC = cpu.pop()
		// cpulogger.Debug(fmt.Sprintf("RET-> PC=0x%02X", cpu.Registers.PC)
		tCycles = 16

	case 0b11001011: // 0xCB -> PREFIX
		// go to prefixed
		cbOpcode := cpu.fetchCBOpcode()
		// // cpulogger.Debug(fmt.Sprintf("Executing CB prefixed opcode: 0x%02X ", cbOpcode)
		switch cbOpcode {
		case 0b0: // 0x00 -> RLC B
			cpu.Registers.B = cpu.execRLC(cpu.Registers.B)
			tCycles = 8

		case 0b1: // 0x01 -> RLC C
			cpu.Registers.C = cpu.execRLC(cpu.Registers.C)
			tCycles = 8

		case 0b10: // 0x02 -> RLC D
			cpu.Registers.D = cpu.execRLC(cpu.Registers.D)
			tCycles = 8

		case 0b11: // 0x03 -> RLC E
			cpu.Registers.E = cpu.execRLC(cpu.Registers.E)
			tCycles = 8

		case 0b100: // 0x04 -> RLC H
			cpu.Registers.H = cpu.execRLC(cpu.Registers.H)
			tCycles = 8

		case 0b101: // 0x05 -> RLC L
			cpu.Registers.L = cpu.execRLC(cpu.Registers.L)
			tCycles = 8

		case 0b110: // 0x06 -> RLC [HL]
			cpu.execRLCHL()
			tCycles = 16

		case 0b111: // 0x07 -> RLC A
			cpu.Registers.A = cpu.execRLC(cpu.Registers.A)
			tCycles = 8

		case 0b1000: // 0x08 -> RRC B
			cpu.Registers.B = cpu.execRRC(cpu.Registers.B)
			tCycles = 8

		case 0b1001: // 0x09 -> RRC C
			cpu.Registers.C = cpu.execRRC(cpu.Registers.C)
			tCycles = 8

		case 0b1010: // 0x0A -> RRC D
			cpu.Registers.D = cpu.execRRC(cpu.Registers.D)
			tCycles = 8

		case 0b1011: // 0x0B -> RRC E
			cpu.Registers.E = cpu.execRRC(cpu.Registers.E)
			tCycles = 8

		case 0b1100: // 0x0C -> RRC H
			cpu.Registers.H = cpu.execRRC(cpu.Registers.H)
			tCycles = 8

		case 0b1101: // 0x0D -> RRC L
			cpu.Registers.L = cpu.execRRC(cpu.Registers.L)
			tCycles = 8

		case 0b1110: // 0x0E -> RRC [HL]
			cpu.execRRCHL()
			tCycles = 16

		case 0b1111: // 0x0F -> RRC A
			cpu.Registers.A = cpu.execRRC(cpu.Registers.A)
			tCycles = 8

		case 0b10000: // 0x10 -> RL B
			cpu.Registers.B = cpu.execRL(cpu.Registers.B)
			tCycles = 8

		case 0b10001: // 0x11 -> RL C
			cpu.Registers.C = cpu.execRL(cpu.Registers.C)
			tCycles = 8

		case 0b10010: // 0x12 -> RL D
			cpu.Registers.D = cpu.execRL(cpu.Registers.D)
			tCycles = 8

		case 0b10011: // 0x13 -> RL E
			cpu.Registers.E = cpu.execRL(cpu.Registers.E)
			tCycles = 8

		case 0b10100: // 0x14 -> RL H
			cpu.Registers.H = cpu.execRL(cpu.Registers.H)
			tCycles = 8

		case 0b10101: // 0x15 -> RL L
			cpu.Registers.L = cpu.execRL(cpu.Registers.L)
			tCycles = 8

		case 0b10110: // 0x16 -> RL [HL]
			cpu.execRLHL()
			tCycles = 16

		case 0b10111: // 0x17 -> RL A
			cpu.Registers.A = cpu.execRL(cpu.Registers.A)
			tCycles = 8

		case 0b11000: // 0x18 -> RR B
			cpu.Registers.B = cpu.execRR(cpu.Registers.B)
			tCycles = 8

		case 0b11001: // 0x19 -> RR C
			cpu.Registers.C = cpu.execRR(cpu.Registers.C)
			tCycles = 8

		case 0b11010: // 0x1A -> RR D
			cpu.Registers.D = cpu.execRR(cpu.Registers.D)
			tCycles = 8

		case 0b11011: // 0x1B -> RR E
			cpu.Registers.E = cpu.execRR(cpu.Registers.E)
			tCycles = 8

		case 0b11100: // 0x1C -> RR H
			cpu.Registers.H = cpu.execRR(cpu.Registers.H)
			tCycles = 8

		case 0b11101: // 0x1D -> RR L
			cpu.Registers.L = cpu.execRR(cpu.Registers.L)
			tCycles = 8

		case 0b11110: // 0x1E -> RR [HL]
			cpu.execRRHL()
			tCycles = 16

		case 0b11111: // 0x1F -> RR A
			cpu.Registers.A = cpu.execRR(cpu.Registers.A)
			tCycles = 8

		case 0b100000: // 0x20 -> SLA B
			cpu.Registers.B = cpu.execSLA(cpu.Registers.B)
			tCycles = 8

		case 0b100001: // 0x21 -> SLA C
			cpu.Registers.C = cpu.execSLA(cpu.Registers.C)
			tCycles = 8

		case 0b100010: // 0x22 -> SLA D
			cpu.Registers.D = cpu.execSLA(cpu.Registers.D)
			tCycles = 8

		case 0b100011: // 0x23 -> SLA E
			cpu.Registers.E = cpu.execSLA(cpu.Registers.E)
			tCycles = 8

		case 0b100100: // 0x24 -> SLA H
			cpu.Registers.H = cpu.execSLA(cpu.Registers.H)
			tCycles = 8

		case 0b100101: // 0x25 -> SLA L
			cpu.Registers.L = cpu.execSLA(cpu.Registers.L)
			tCycles = 8

		case 0b100110: // 0x26 -> SLA [HL]
			cpu.execSLAHL()
			tCycles = 16

		case 0b100111: // 0x27 -> SLA A
			cpu.Registers.A = cpu.execSLA(cpu.Registers.A)
			tCycles = 8

		case 0b101000: // 0x28 -> SRA B
			cpu.Registers.B = cpu.execSRA(cpu.Registers.B)
			tCycles = 8

		case 0b101001: // 0x29 -> SRA C
			cpu.Registers.C = cpu.execSRA(cpu.Registers.C)
			tCycles = 8

		case 0b101010: // 0x2A -> SRA D
			cpu.Registers.D = cpu.execSRA(cpu.Registers.D)
			tCycles = 8

		case 0b101011: // 0x2B -> SRA E
			cpu.Registers.E = cpu.execSRA(cpu.Registers.E)
			tCycles = 8

		case 0b101100: // 0x2C -> SRA H
			cpu.Registers.H = cpu.execSRA(cpu.Registers.H)
			tCycles = 8

		case 0b101101: // 0x2D -> SRA L
			cpu.Registers.L = cpu.execSRA(cpu.Registers.L)
			tCycles = 8

		case 0b101110: // 0x2E -> SRA [HL]
			cpu.execSRAHL()
			tCycles = 16

		case 0b101111: // 0x2F -> SRA A
			cpu.Registers.A = cpu.execSRA(cpu.Registers.A)
			tCycles = 8

		case 0b110000: // 0x30 -> SWAP B
			cpu.Registers.B = cpu.execSWAP(cpu.Registers.B)
			tCycles = 8

		case 0b110001: // 0x31 -> SWAP C
			cpu.Registers.C = cpu.execSWAP(cpu.Registers.C)
			tCycles = 8

		case 0b110010: // 0x32 -> SWAP D
			cpu.Registers.D = cpu.execSWAP(cpu.Registers.D)
			tCycles = 8

		case 0b110011: // 0x33 -> SWAP E
			cpu.Registers.E = cpu.execSWAP(cpu.Registers.E)
			tCycles = 8

		case 0b110100: // 0x34 -> SWAP H
			cpu.Registers.H = cpu.execSWAP(cpu.Registers.H)
			tCycles = 8

		case 0b110101: // 0x35 -> SWAP L
			cpu.Registers.L = cpu.execSWAP(cpu.Registers.L)
			tCycles = 8

		case 0b110110: // 0x36 -> SWAP [HL]
			cpu.execSWAPHL()
			tCycles = 16

		case 0b110111: // 0x37 -> SWAP A
			cpu.Registers.A = cpu.execSWAP(cpu.Registers.A)
			tCycles = 8

		case 0b111000: // 0x38 -> SRL B
			cpu.Registers.B = cpu.execSRL(cpu.Registers.B)
			tCycles = 8

		case 0b111001: // 0x39 -> SRL C
			cpu.Registers.C = cpu.execSRL(cpu.Registers.C)
			tCycles = 8

		case 0b111010: // 0x3A -> SRL D
			cpu.Registers.D = cpu.execSRL(cpu.Registers.D)
			tCycles = 8

		case 0b111011: // 0x3B -> SRL E
			cpu.Registers.E = cpu.execSRL(cpu.Registers.E)
			tCycles = 8

		case 0b111100: // 0x3C -> SRL H
			cpu.Registers.H = cpu.execSRL(cpu.Registers.H)
			tCycles = 8

		case 0b111101: // 0x3D -> SRL L
			cpu.Registers.L = cpu.execSRL(cpu.Registers.L)
			tCycles = 8

		case 0b111110: // 0x3E -> SRL [HL]
			cpu.execSRLHL()
			tCycles = 16

		case 0b111111: // 0x3F -> SRL A
			cpu.Registers.A = cpu.execSRL(cpu.Registers.A)
			tCycles = 8

		case 0b1000000: // 0x40 -> BIT 0, B
			cpu.execBIT(cpu.Registers.B, 0)
			tCycles = 8

		case 0b1000001: // 0x41 -> BIT 0, C
			cpu.execBIT(cpu.Registers.C, 0)
			tCycles = 8

		case 0b1000010: // 0x42 -> BIT 0, D
			cpu.execBIT(cpu.Registers.D, 0)
			tCycles = 8

		case 0b1000011: // 0x43 -> BIT 0, E
			cpu.execBIT(cpu.Registers.E, 0)
			tCycles = 8

		case 0b1000100: // 0x44 -> BIT 0, H
			cpu.execBIT(cpu.Registers.H, 0)
			tCycles = 8

		case 0b1000101: // 0x45 -> BIT 0, L
			cpu.execBIT(cpu.Registers.L, 0)
			tCycles = 8

		case 0b1000110: // 0x46 -> BIT 0, [HL]
			cpu.execBITHL(0)
			tCycles = 12

		case 0b1000111: // 0x47 -> BIT 0, A
			cpu.execBIT(cpu.Registers.A, 0)
			tCycles = 8

		case 0b1001000: // 0x48 -> BIT 1, B
			cpu.execBIT(cpu.Registers.B, 1)
			tCycles = 8

		case 0b1001001: // 0x49 -> BIT 1, C
			cpu.execBIT(cpu.Registers.C, 1)
			tCycles = 8

		case 0b1001010: // 0x4A -> BIT 1, D
			cpu.execBIT(cpu.Registers.D, 1)
			tCycles = 8

		case 0b1001011: // 0x4B -> BIT 1, E
			cpu.execBIT(cpu.Registers.E, 1)
			tCycles = 8

		case 0b1001100: // 0x4C -> BIT 1, H
			cpu.execBIT(cpu.Registers.H, 1)
			tCycles = 8

		case 0b1001101: // 0x4D -> BIT 1, L
			cpu.execBIT(cpu.Registers.L, 1)
			tCycles = 8

		case 0b1001110: // 0x4E -> BIT 1, [HL]
			cpu.execBITHL(1)
			tCycles = 12

		case 0b1001111: // 0x4F -> BIT 1, A
			cpu.execBIT(cpu.Registers.A, 1)
			tCycles = 8

		case 0b1010000: // 0x50 -> BIT 2, B
			cpu.execBIT(cpu.Registers.B, 2)
			tCycles = 8

		case 0b1010001: // 0x51 -> BIT 2, C
			cpu.execBIT(cpu.Registers.C, 2)
			tCycles = 8

		case 0b1010010: // 0x52 -> BIT 2, D
			cpu.execBIT(cpu.Registers.D, 2)
			tCycles = 8

		case 0b1010011: // 0x53 -> BIT 2, E
			cpu.execBIT(cpu.Registers.E, 2)
			tCycles = 8

		case 0b1010100: // 0x54 -> BIT 2, H
			cpu.execBIT(cpu.Registers.H, 2)
			tCycles = 8

		case 0b1010101: // 0x55 -> BIT 2, L
			cpu.execBIT(cpu.Registers.L, 2)
			tCycles = 8

		case 0b1010110: // 0x56 -> BIT 2, [HL]
			cpu.execBITHL(2)
			tCycles = 12

		case 0b1010111: // 0x57 -> BIT 2, A
			cpu.execBIT(cpu.Registers.A, 2)
			tCycles = 8

		case 0b1011000: // 0x58 -> BIT 3, B
			cpu.execBIT(cpu.Registers.B, 3)
			tCycles = 8

		case 0b1011001: // 0x59 -> BIT 3, C
			cpu.execBIT(cpu.Registers.C, 3)
			tCycles = 8

		case 0b1011010: // 0x5A -> BIT 3, D
			cpu.execBIT(cpu.Registers.D, 3)
			tCycles = 8

		case 0b1011011: // 0x5B -> BIT 3, E
			cpu.execBIT(cpu.Registers.E, 3)
			tCycles = 8

		case 0b1011100: // 0x5C -> BIT 3, H
			cpu.execBIT(cpu.Registers.H, 3)
			tCycles = 8

		case 0b1011101: // 0x5D -> BIT 3, L
			cpu.execBIT(cpu.Registers.L, 3)
			tCycles = 8

		case 0b1011110: // 0x5E -> BIT 3, [HL]
			cpu.execBITHL(3)
			tCycles = 12

		case 0b1011111: // 0x5F -> BIT 3, A
			cpu.execBIT(cpu.Registers.A, 3)
			tCycles = 8

		case 0b1100000: // 0x60 -> BIT 4, B
			cpu.execBIT(cpu.Registers.B, 4)
			tCycles = 8

		case 0b1100001: // 0x61 -> BIT 4, C
			cpu.execBIT(cpu.Registers.C, 4)
			tCycles = 8

		case 0b1100010: // 0x62 -> BIT 4, D
			cpu.execBIT(cpu.Registers.D, 4)
			tCycles = 8

		case 0b1100011: // 0x63 -> BIT 4, E
			cpu.execBIT(cpu.Registers.E, 4)
			tCycles = 8

		case 0b1100100: // 0x64 -> BIT 4, H
			cpu.execBIT(cpu.Registers.H, 4)
			tCycles = 8

		case 0b1100101: // 0x65 -> BIT 4, L
			cpu.execBIT(cpu.Registers.L, 4)
			tCycles = 8

		case 0b1100110: // 0x66 -> BIT 4, [HL]
			cpu.execBITHL(4)
			tCycles = 12

		case 0b1100111: // 0x67 -> BIT 4, A
			cpu.execBIT(cpu.Registers.A, 4)
			tCycles = 8

		case 0b1101000: // 0x68 -> BIT 5, B
			cpu.execBIT(cpu.Registers.B, 5)
			tCycles = 8

		case 0b1101001: // 0x69 -> BIT 5, C
			cpu.execBIT(cpu.Registers.C, 5)
			tCycles = 8

		case 0b1101010: // 0x6A -> BIT 5, D
			cpu.execBIT(cpu.Registers.D, 5)
			tCycles = 8

		case 0b1101011: // 0x6B -> BIT 5, E
			cpu.execBIT(cpu.Registers.E, 5)
			tCycles = 8

		case 0b1101100: // 0x6C -> BIT 5, H
			cpu.execBIT(cpu.Registers.H, 5)
			tCycles = 8

		case 0b1101101: // 0x6D -> BIT 5, L
			cpu.execBIT(cpu.Registers.L, 5)
			tCycles = 8

		case 0b1101110: // 0x6E -> BIT 5, [HL]
			cpu.execBITHL(5)
			tCycles = 12

		case 0b1101111: // 0x6F -> BIT 5, A
			cpu.execBIT(cpu.Registers.A, 5)
			tCycles = 8

		case 0b1110000: // 0x70 -> BIT 6, B
			cpu.execBIT(cpu.Registers.B, 6)
			tCycles = 8

		case 0b1110001: // 0x71 -> BIT 6, C
			cpu.execBIT(cpu.Registers.C, 6)
			tCycles = 8

		case 0b1110010: // 0x72 -> BIT 6, D
			cpu.execBIT(cpu.Registers.D, 6)
			tCycles = 8

		case 0b1110011: // 0x73 -> BIT 6, E
			cpu.execBIT(cpu.Registers.E, 6)
			tCycles = 8

		case 0b1110100: // 0x74 -> BIT 6, H
			cpu.execBIT(cpu.Registers.H, 6)
			tCycles = 8

		case 0b1110101: // 0x75 -> BIT 6, L
			cpu.execBIT(cpu.Registers.L, 6)
			tCycles = 8

		case 0b1110110: // 0x76 -> BIT 6, [HL]
			cpu.execBITHL(6)
			tCycles = 12

		case 0b1110111: // 0x77 -> BIT 6, A
			cpu.execBIT(cpu.Registers.A, 6)
			tCycles = 8

		case 0b1111000: // 0x78 -> BIT 7, B
			cpu.execBIT(cpu.Registers.B, 7)
			tCycles = 8

		case 0b1111001: // 0x79 -> BIT 7, C
			cpu.execBIT(cpu.Registers.C, 7)
			tCycles = 8

		case 0b1111010: // 0x7A -> BIT 7, D
			cpu.execBIT(cpu.Registers.D, 7)
			tCycles = 8

		case 0b1111011: // 0x7B -> BIT 7, E
			cpu.execBIT(cpu.Registers.E, 7)
			tCycles = 8

		case 0b1111100: // 0x7C -> BIT 7, H
			cpu.execBIT(cpu.Registers.H, 7)
			tCycles = 8

		case 0b1111101: // 0x7D -> BIT 7, L
			cpu.execBIT(cpu.Registers.L, 7)
			tCycles = 8

		case 0b1111110: // 0x7E -> BIT 7, [HL]
			cpu.execBITHL(7)
			tCycles = 12

		case 0b1111111: // 0x7F -> BIT 7, A
			cpu.execBIT(cpu.Registers.A, 7)
			tCycles = 8

		case 0b10000000: // 0x80 -> RES 0, B
			cpu.Registers.B = cpu.execRES(cpu.Registers.B, 0)
			tCycles = 8

		case 0b10000001: // 0x81 -> RES 0, C
			cpu.Registers.C = cpu.execRES(cpu.Registers.C, 0)
			tCycles = 8

		case 0b10000010: // 0x82 -> RES 0, D
			cpu.Registers.D = cpu.execRES(cpu.Registers.D, 0)
			tCycles = 8

		case 0b10000011: // 0x83 -> RES 0, E
			cpu.Registers.E = cpu.execRES(cpu.Registers.E, 0)
			tCycles = 8

		case 0b10000100: // 0x84 -> RES 0, H
			cpu.Registers.H = cpu.execRES(cpu.Registers.H, 0)
			tCycles = 8

		case 0b10000101: // 0x85 -> RES 0, L
			cpu.Registers.L = cpu.execRES(cpu.Registers.L, 0)
			tCycles = 8

		case 0b10000110: // 0x86 -> RES 0, [HL]
			cpu.execRESHL(0)
			tCycles = 16

		case 0b10000111: // 0x87 -> RES 0, A
			cpu.Registers.A = cpu.execRES(cpu.Registers.A, 0)
			// cpulogger.Debug(fmt.Sprintf("Now executing RES at 0x87")
			tCycles = 8

		case 0b10001000: // 0x88 -> RES 1, B
			cpu.Registers.B = cpu.execRES(cpu.Registers.B, 1)
			tCycles = 8

		case 0b10001001: // 0x89 -> RES 1, C
			cpu.Registers.C = cpu.execRES(cpu.Registers.C, 1)
			tCycles = 8

		case 0b10001010: // 0x8A -> RES 1, D
			cpu.Registers.D = cpu.execRES(cpu.Registers.D, 1)
			tCycles = 8

		case 0b10001011: // 0x8B -> RES 1, E
			cpu.Registers.E = cpu.execRES(cpu.Registers.E, 1)
			tCycles = 8

		case 0b10001100: // 0x8C -> RES 1, H
			cpu.Registers.H = cpu.execRES(cpu.Registers.H, 1)
			tCycles = 8

		case 0b10001101: // 0x8D -> RES 1, L
			cpu.Registers.L = cpu.execRES(cpu.Registers.L, 1)
			tCycles = 8

		case 0b10001110: // 0x8E -> RES 1, [HL]
			cpu.execRESHL(1)
			tCycles = 16

		case 0b10001111: // 0x8F -> RES 1, A
			cpu.Registers.A = cpu.execRES(cpu.Registers.A, 1)
			tCycles = 8

		case 0b10010000: // 0x90 -> RES 2, B
			cpu.Registers.B = cpu.execRES(cpu.Registers.B, 2)
			tCycles = 8

		case 0b10010001: // 0x91 -> RES 2, C
			cpu.Registers.C = cpu.execRES(cpu.Registers.C, 2)
			tCycles = 8

		case 0b10010010: // 0x92 -> RES 2, D
			cpu.Registers.D = cpu.execRES(cpu.Registers.D, 2)
			tCycles = 8

		case 0b10010011: // 0x93 -> RES 2, E
			cpu.Registers.E = cpu.execRES(cpu.Registers.E, 2)
			tCycles = 8

		case 0b10010100: // 0x94 -> RES 2, H
			cpu.Registers.H = cpu.execRES(cpu.Registers.H, 2)
			tCycles = 8

		case 0b10010101: // 0x95 -> RES 2, L
			cpu.Registers.L = cpu.execRES(cpu.Registers.L, 2)
			tCycles = 8

		case 0b10010110: // 0x96 -> RES 2, [HL]
			cpu.execRESHL(2)
			tCycles = 16

		case 0b10010111: // 0x97 -> RES 2, A
			cpu.Registers.A = cpu.execRES(cpu.Registers.A, 2)
			tCycles = 8

		case 0b10011000: // 0x98 -> RES 3, B
			cpu.Registers.B = cpu.execRES(cpu.Registers.B, 3)
			tCycles = 8

		case 0b10011001: // 0x99 -> RES 3, C
			cpu.Registers.C = cpu.execRES(cpu.Registers.C, 3)
			tCycles = 8

		case 0b10011010: // 0x9A -> RES 3, D
			cpu.Registers.D = cpu.execRES(cpu.Registers.D, 3)
			tCycles = 8

		case 0b10011011: // 0x9B -> RES 3, E
			cpu.Registers.E = cpu.execRES(cpu.Registers.E, 3)
			tCycles = 8

		case 0b10011100: // 0x9C -> RES 3, H
			cpu.Registers.H = cpu.execRES(cpu.Registers.H, 3)
			tCycles = 8

		case 0b10011101: // 0x9D -> RES 3, L
			cpu.Registers.L = cpu.execRES(cpu.Registers.L, 3)
			tCycles = 8

		case 0b10011110: // 0x9E -> RES 3, [HL]
			cpu.execRESHL(3)
			tCycles = 16

		case 0b10011111: // 0x9F -> RES 3, A
			cpu.Registers.A = cpu.execRES(cpu.Registers.A, 3)
			tCycles = 8

		case 0b10100000: // 0xA0 -> RES 4, B
			cpu.Registers.B = cpu.execRES(cpu.Registers.B, 4)
			tCycles = 8

		case 0b10100001: // 0xA1 -> RES 4, C
			cpu.Registers.C = cpu.execRES(cpu.Registers.C, 4)
			tCycles = 8

		case 0b10100010: // 0xA2 -> RES 4, D
			cpu.Registers.D = cpu.execRES(cpu.Registers.D, 4)
			tCycles = 8

		case 0b10100011: // 0xA3 -> RES 4, E
			cpu.Registers.E = cpu.execRES(cpu.Registers.E, 4)
			tCycles = 8

		case 0b10100100: // 0xA4 -> RES 4, H
			cpu.Registers.H = cpu.execRES(cpu.Registers.H, 4)
			tCycles = 8

		case 0b10100101: // 0xA5 -> RES 4, L
			cpu.Registers.L = cpu.execRES(cpu.Registers.L, 4)
			tCycles = 8

		case 0b10100110: // 0xA6 -> RES 4, [HL]
			cpu.execRESHL(4)
			tCycles = 16

		case 0b10100111: // 0xA7 -> RES 4, A
			cpu.Registers.A = cpu.execRES(cpu.Registers.A, 4)
			tCycles = 8

		case 0b10101000: // 0xA8 -> RES 5, B
			cpu.Registers.B = cpu.execRES(cpu.Registers.B, 5)
			tCycles = 8

		case 0b10101001: // 0xA9 -> RES 5, C
			cpu.Registers.C = cpu.execRES(cpu.Registers.C, 5)
			tCycles = 8

		case 0b10101010: // 0xAA -> RES 5, D
			cpu.Registers.D = cpu.execRES(cpu.Registers.D, 5)
			tCycles = 8

		case 0b10101011: // 0xAB -> RES 5, E
			cpu.Registers.E = cpu.execRES(cpu.Registers.E, 5)
			tCycles = 8

		case 0b10101100: // 0xAC -> RES 5, H
			cpu.Registers.H = cpu.execRES(cpu.Registers.H, 5)
			tCycles = 8

		case 0b10101101: // 0xAD -> RES 5, L
			cpu.Registers.L = cpu.execRES(cpu.Registers.L, 5)
			tCycles = 8

		case 0b10101110: // 0xAE -> RES 5, [HL]
			cpu.execRESHL(5)
			tCycles = 16

		case 0b10101111: // 0xAF -> RES 5, A
			cpu.Registers.A = cpu.execRES(cpu.Registers.A, 5)
			tCycles = 8

		case 0b10110000: // 0xB0 -> RES 6, B
			cpu.Registers.B = cpu.execRES(cpu.Registers.B, 6)
			tCycles = 8

		case 0b10110001: // 0xB1 -> RES 6, C
			cpu.Registers.C = cpu.execRES(cpu.Registers.C, 6)
			tCycles = 8

		case 0b10110010: // 0xB2 -> RES 6, D
			cpu.Registers.D = cpu.execRES(cpu.Registers.D, 6)
			tCycles = 8

		case 0b10110011: // 0xB3 -> RES 6, E
			cpu.Registers.E = cpu.execRES(cpu.Registers.E, 6)
			tCycles = 8

		case 0b10110100: // 0xB4 -> RES 6, H
			cpu.Registers.H = cpu.execRES(cpu.Registers.H, 6)
			tCycles = 8

		case 0b10110101: // 0xB5 -> RES 6, L
			cpu.Registers.L = cpu.execRES(cpu.Registers.L, 6)
			tCycles = 8

		case 0b10110110: // 0xB6 -> RES 6, [HL]
			cpu.execRESHL(6)
			tCycles = 16

		case 0b10110111: // 0xB7 -> RES 6, A
			cpu.Registers.A = cpu.execRES(cpu.Registers.A, 6)
			tCycles = 8

		case 0b10111000: // 0xB8 -> RES 7, B
			cpu.Registers.B = cpu.execRES(cpu.Registers.B, 7)
			tCycles = 8

		case 0b10111001: // 0xB9 -> RES 7, C
			cpu.Registers.C = cpu.execRES(cpu.Registers.C, 7)
			tCycles = 8

		case 0b10111010: // 0xBA -> RES 7, D
			cpu.Registers.D = cpu.execRES(cpu.Registers.D, 7)
			tCycles = 8

		case 0b10111011: // 0xBB -> RES 7, E
			cpu.Registers.E = cpu.execRES(cpu.Registers.E, 7)
			tCycles = 8

		case 0b10111100: // 0xBC -> RES 7, H
			cpu.Registers.H = cpu.execRES(cpu.Registers.H, 7)
			tCycles = 8

		case 0b10111101: // 0xBD -> RES 7, L
			cpu.Registers.L = cpu.execRES(cpu.Registers.L, 7)
			tCycles = 8

		case 0b10111110: // 0xBE -> RES 7, [HL]
			cpu.execRESHL(7)
			tCycles = 16

		case 0b10111111: // 0xBF -> RES 7, A
			cpu.Registers.A = cpu.execRES(cpu.Registers.A, 7)
			tCycles = 8

		case 0b11000000: // 0xC0 -> SET 0, B
			cpu.Registers.B = cpu.execSET(cpu.Registers.B, 0)
			tCycles = 8

		case 0b11000001: // 0xC1 -> SET 0, C
			cpu.Registers.C = cpu.execSET(cpu.Registers.C, 0)
			tCycles = 8

		case 0b11000010: // 0xC2 -> SET 0, D
			cpu.Registers.D = cpu.execSET(cpu.Registers.D, 0)
			tCycles = 8

		case 0b11000011: // 0xC3 -> SET 0, E
			cpu.Registers.E = cpu.execSET(cpu.Registers.E, 0)
			tCycles = 8

		case 0b11000100: // 0xC4 -> SET 0, H
			cpu.Registers.H = cpu.execSET(cpu.Registers.H, 0)
			tCycles = 8

		case 0b11000101: // 0xC5 -> SET 0, L
			cpu.Registers.L = cpu.execSET(cpu.Registers.L, 0)
			tCycles = 8

		case 0b11000110: // 0xC6 -> SET 0, [HL]
			cpu.execSETHL(0)
			tCycles = 16

		case 0b11000111: // 0xC7 -> SET 0, A
			cpu.Registers.A = cpu.execSET(cpu.Registers.A, 0)
			tCycles = 8

		case 0b11001000: // 0xC8 -> SET 1, B
			cpu.Registers.B = cpu.execSET(cpu.Registers.B, 1)
			tCycles = 8

		case 0b11001001: // 0xC9 -> SET 1, C
			cpu.Registers.C = cpu.execSET(cpu.Registers.C, 1)
			tCycles = 8

		case 0b11001010: // 0xCA -> SET 1, D
			cpu.Registers.D = cpu.execSET(cpu.Registers.D, 1)
			tCycles = 8

		case 0b11001011: // 0xCB -> SET 1, E
			cpu.Registers.E = cpu.execSET(cpu.Registers.E, 1)
			tCycles = 8

		case 0b11001100: // 0xCC -> SET 1, H
			cpu.Registers.H = cpu.execSET(cpu.Registers.H, 1)
			tCycles = 8

		case 0b11001101: // 0xCD -> SET 1, L
			cpu.Registers.L = cpu.execSET(cpu.Registers.L, 1)
			tCycles = 8

		case 0b11001110: // 0xCE -> SET 1, [HL]
			cpu.execSETHL(1)
			tCycles = 16

		case 0b11001111: // 0xCF -> SET 1, A
			cpu.Registers.A = cpu.execSET(cpu.Registers.A, 1)
			tCycles = 8

		case 0b11010000: // 0xD0 -> SET 2, B
			cpu.Registers.B = cpu.execSET(cpu.Registers.B, 2)
			tCycles = 8

		case 0b11010001: // 0xD1 -> SET 2, C
			cpu.Registers.C = cpu.execSET(cpu.Registers.C, 2)
			tCycles = 8

		case 0b11010010: // 0xD2 -> SET 2, D
			cpu.Registers.D = cpu.execSET(cpu.Registers.D, 2)
			tCycles = 8

		case 0b11010011: // 0xD3 -> SET 2, E
			cpu.Registers.E = cpu.execSET(cpu.Registers.E, 2)
			tCycles = 8

		case 0b11010100: // 0xD4 -> SET 2, H
			cpu.Registers.H = cpu.execSET(cpu.Registers.H, 2)
			tCycles = 8

		case 0b11010101: // 0xD5 -> SET 2, L
			cpu.Registers.L = cpu.execSET(cpu.Registers.L, 2)
			tCycles = 8

		case 0b11010110: // 0xD6 -> SET 2, [HL]
			cpu.execSETHL(2)
			tCycles = 16

		case 0b11010111: // 0xD7 -> SET 2, A
			cpu.Registers.A = cpu.execSET(cpu.Registers.A, 2)
			tCycles = 8

		case 0b11011000: // 0xD8 -> SET 3, B
			cpu.Registers.B = cpu.execSET(cpu.Registers.B, 3)
			tCycles = 8

		case 0b11011001: // 0xD9 -> SET 3, C
			cpu.Registers.C = cpu.execSET(cpu.Registers.C, 3)
			tCycles = 8

		case 0b11011010: // 0xDA -> SET 3, D
			cpu.Registers.D = cpu.execSET(cpu.Registers.D, 3)
			tCycles = 8

		case 0b11011011: // 0xDB -> SET 3, E
			cpu.Registers.E = cpu.execSET(cpu.Registers.E, 3)
			tCycles = 8

		case 0b11011100: // 0xDC -> SET 3, H
			cpu.Registers.H = cpu.execSET(cpu.Registers.H, 3)
			tCycles = 8

		case 0b11011101: // 0xDD -> SET 3, L
			cpu.Registers.L = cpu.execSET(cpu.Registers.L, 3)
			tCycles = 8

		case 0b11011110: // 0xDE -> SET 3, [HL]
			cpu.execSETHL(3)
			tCycles = 16

		case 0b11011111: // 0xDF -> SET 3, A
			cpu.Registers.A = cpu.execSET(cpu.Registers.A, 3)
			tCycles = 8

		case 0b11100000: // 0xE0 -> SET 4, B
			cpu.Registers.B = cpu.execSET(cpu.Registers.B, 4)
			tCycles = 8

		case 0b11100001: // 0xE1 -> SET 4, C
			cpu.Registers.C = cpu.execSET(cpu.Registers.C, 4)
			tCycles = 8

		case 0b11100010: // 0xE2 -> SET 4, D
			cpu.Registers.D = cpu.execSET(cpu.Registers.D, 4)
			tCycles = 8

		case 0b11100011: // 0xE3 -> SET 4, E
			cpu.Registers.E = cpu.execSET(cpu.Registers.E, 4)
			tCycles = 8

		case 0b11100100: // 0xE4 -> SET 4, H
			cpu.Registers.H = cpu.execSET(cpu.Registers.H, 4)
			tCycles = 8

		case 0b11100101: // 0xE5 -> SET 4, L
			cpu.Registers.L = cpu.execSET(cpu.Registers.L, 4)
			tCycles = 8

		case 0b11100110: // 0xE6 -> SET 4, [HL]
			cpu.execSETHL(4)
			tCycles = 16

		case 0b11100111: // 0xE7 -> SET 4, A
			cpu.Registers.A = cpu.execSET(cpu.Registers.A, 4)
			tCycles = 8

		case 0b11101000: // 0xE8 -> SET 5, B
			cpu.Registers.B = cpu.execSET(cpu.Registers.B, 5)
			tCycles = 8

		case 0b11101001: // 0xE9 -> SET 5, C
			cpu.Registers.C = cpu.execSET(cpu.Registers.C, 5)
			tCycles = 8

		case 0b11101010: // 0xEA -> SET 5, D
			cpu.Registers.D = cpu.execSET(cpu.Registers.D, 5)
			tCycles = 8

		case 0b11101011: // 0xEB -> SET 5, E
			cpu.Registers.E = cpu.execSET(cpu.Registers.E, 5)
			tCycles = 8

		case 0b11101100: // 0xEC -> SET 5, H
			cpu.Registers.H = cpu.execSET(cpu.Registers.H, 5)
			tCycles = 8

		case 0b11101101: // 0xED -> SET 5, L
			cpu.Registers.L = cpu.execSET(cpu.Registers.L, 5)
			tCycles = 8

		case 0b11101110: // 0xEE -> SET 5, [HL]
			cpu.execSETHL(5)
			tCycles = 16

		case 0b11101111: // 0xEF -> SET 5, A
			cpu.Registers.A = cpu.execSET(cpu.Registers.A, 5)
			tCycles = 8

		case 0b11110000: // 0xF0 -> SET 6, B
			cpu.Registers.B = cpu.execSET(cpu.Registers.B, 6)
			tCycles = 8

		case 0b11110001: // 0xF1 -> SET 6, C
			cpu.Registers.C = cpu.execSET(cpu.Registers.C, 6)
			tCycles = 8

		case 0b11110010: // 0xF2 -> SET 6, D
			cpu.Registers.D = cpu.execSET(cpu.Registers.D, 6)
			tCycles = 8

		case 0b11110011: // 0xF3 -> SET 6, E
			cpu.Registers.E = cpu.execSET(cpu.Registers.E, 6)
			tCycles = 8

		case 0b11110100: // 0xF4 -> SET 6, H
			cpu.Registers.H = cpu.execSET(cpu.Registers.H, 6)
			tCycles = 8

		case 0b11110101: // 0xF5 -> SET 6, L
			cpu.Registers.L = cpu.execSET(cpu.Registers.L, 6)
			tCycles = 8

		case 0b11110110: // 0xF6 -> SET 6, [HL]
			cpu.execSETHL(6)
			tCycles = 16

		case 0b11110111: // 0xF7 -> SET 6, A
			cpu.Registers.A = cpu.execSET(cpu.Registers.A, 6)
			tCycles = 8

		case 0b11111000: // 0xF8 -> SET 7, B
			cpu.Registers.B = cpu.execSET(cpu.Registers.B, 7)
			tCycles = 8

		case 0b11111001: // 0xF9 -> SET 7, C
			cpu.Registers.C = cpu.execSET(cpu.Registers.C, 7)
			tCycles = 8

		case 0b11111010: // 0xFA -> SET 7, D
			cpu.Registers.D = cpu.execSET(cpu.Registers.D, 7)
			tCycles = 8

		case 0b11111011: // 0xFB -> SET 7, E
			cpu.Registers.E = cpu.execSET(cpu.Registers.E, 7)
			tCycles = 8

		case 0b11111100: // 0xFC -> SET 7, H
			cpu.Registers.H = cpu.execSET(cpu.Registers.H, 7)
			tCycles = 8

		case 0b11111101: // 0xFD -> SET 7, L
			cpu.Registers.L = cpu.execSET(cpu.Registers.L, 7)
			tCycles = 8

		case 0b11111110: // 0xFE -> SET 7, [HL]
			cpu.execSETHL(7)
			tCycles = 16

		case 0b11111111: // 0xFF -> SET 7, A
			cpu.Registers.A = cpu.execSET(cpu.Registers.A, 7)
			tCycles = 8

		default:
			panic("Not an instruction in cb")
		}
	case 0b11010011: // 0xD3 -> ILLEGAL_D3
		// cpucpulogger.Debug("ILLEGAL_D3")
		tCycles = 4
		break
	case 0b11011001: // 0xD9 -> RETI

		cpulogger.Debug("RETI instruction")
		//Pop two bytes from stack & jump to that address then
		// enable interrupts.
		for i := 0; i < 8; i += 2 {
			// cpulogger.Debug(fmt.Sprintf("Stack[SP+%d]: 0x%02X%02X", i, cpu.Memory[cpu.Registers.SP+uint16(i)+1], cpu.Memory[cpu.Registers.SP+uint16(i)])
		}
		// cpulogger.Debug(fmt.Sprintf("SP before RETI : 0x%04X", cpu.Registers.SP)
		// cpulogger.Debug(fmt.Sprintf("Stack[SP]: 0x%02X, Stack[SP+1] : 0x%02X", cpu.Memory[cpu.Registers.SP], cpu.Memory[cpu.Registers.SP+1])

		cpu.Registers.PC = cpu.pop()
		cpu.IME = true
		// cpulogger.Debug(fmt.Sprintf("RETI -> PC: 0x%04X", cpu.Registers.PC)
		tCycles = 16

	case 0b11011011: // 0xDB -> ILLEGAL_DB
		// cpucpulogger.Debug("ILLEGAL_DB")
		tCycles = 4
		break
	case 0b11011101: // 0xDD -> ILLEGAL_DD
		// cpucpulogger.Debug("ILLEGAL_DD")
		tCycles = 4
		break
	case 0b11100011: // 0xE3 -> ILLEGAL_E3
		// cpucpulogger.Debug("ILLEGAL_E3")
		tCycles = 4
		break
	case 0b11100100: // 0xE4 -> ILLEGAL_E4
		// cpucpulogger.Debug("ILLEGAL_E4")
		tCycles = 4
		break
	case 0b11101011: // 0xEB -> ILLEGAL_EB
		// cpucpulogger.Debug("ILLEGAL_EB")
		tCycles = 4
		break
	case 0b11101100: // 0xEC -> ILLEGAL_EC

		// cpucpulogger.Debug("ILLEGAL_EC")
		tCycles = 4
		break
	case 0b11101101: // 0xED -> ILLEGAL_ED
		// cpucpulogger.Debug("ILLEGAL_ED")
		tCycles = 4
		break
	case 0b11110011: // 0xF3 -> DI

		// This instruction disables interrupts but not
		// immediately. Interrupts are disabled after
		// instruction after DI is executed.
		cpulogger.Debug("DI instruction")
		cpu.IME = false
		tCycles = 4
	case 0b11110100: // 0xF4 -> ILLEGAL_F4
		cpulogger.Debug("ILLEGAL_F4")
		tCycles = 4
		break
	case 0b11111011: // 0xFB -> EI
		//The effect of ei is delayed by one instruction !
		// Enable interrupts. This intruction enables interrupts
		// but not immediately. Interrupts are enabled after
		// instruction after EI is executed.
		cpulogger.Debug("EI instruction")
		cpulogger.Debug(fmt.Sprintf(""))
		//cpu.IMEScheduled = true
		cpu.IME = true
		tCycles = 4
	case 0b11111100: // 0xFC -> ILLEGAL_FC
		cpulogger.Debug("ILLEGAL_FC")
		tCycles = 4
		break
	case 0b11111101: // 0xFD -> ILLEGAL_FD
		cpulogger.Debug("ILLEGAL_FD")
		tCycles = 4
		break

	default:
		panic("Not an operation")
	}
	//if cpu.Registers.PC > 0xFFFF {
	//	// cpucpulogger.Debug("PC out of bounds")
	//	break
	//}
	//cpu.Memory[0xFF40] |= 1 << 7

	if cpu.handleInterruptions() {
		return 0
	}
	cpu.checkSchedule()

	return tCycles
}

func (cpu *CPU) loadROMFile(cartridge *Cartridge) {
	cpu.Cartridge = cartridge
	copy(cpu.Memory[:], cartridge.ROMdata)
	//copy(cpu.Memory[:0xFE01], cartridge.ROMdata[:0xFE01])
	//copy(cpu.Memory[0x0000:0x8000], cartridge.ROMdata[:0x8000])

	//copy(graphic.VRAM[:], cartridge.ROMdata[:VRAM_SIZE])
	//copy(graphic.OAM[:], cartridge.ROMdata[OAM_START:OAM_END])
	cpu.Registers.PC = 0x100
}
