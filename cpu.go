// package GameBoy_Emulator
package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// TODO UNFIXED: AND, OR, XOR, PUSH, POP, SBC, CP,
// TODO: DAA, CPL, CCF, SCF, NOP, HALT, STOP, DI, EI
// TODO: JP, JR, CALL
// TODO: RST, RET, RETI

// TODO PREFIXED: SWAP
// TODO: RLCA, RLA, RRCA, RRA, RLC, RL, RRC, RR, SLA, SRA, SRL
//TODO: BIT, SET, RES

// A, B, C, D, E, F, H, and L - 8 bit
// SP, PC - 16 bit
// AF, BC, DE, and HL
// https://github.com/veandco/go-sdl2
type CPU struct {
	Registers Registers
	Memory    [8192]uint8
	//OpcodesTable map[string]map[string][]map[string]string
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

func (register *Registers) setFlag(flag uint8, on bool) {
	//var register Registers
	if on {
		register.F |= flag //set bit

	} else {
		register.F &= ^flag //clear bit
	}
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

func (cpu *CPU) getImmediate8() uint8 {
	val := cpu.Memory[cpu.Registers.PC]
	cpu.Registers.PC++
	return val
}

func (cpu *CPU) getImmediate16() uint8 {
	val := cpu.Memory[cpu.Registers.PC]
	cpu.Registers.PC += 2
	return val
}

func (cpu *CPU) execLD(operands []map[string]string) {
	dest := operands[0]["name"]
	destImmd := operands[0]["immediate"]

	src := operands[1]["name"]
	srcImmd := operands[1]["immediate"]

	fmt.Println("src:", src, "srcImmd:", srcImmd, "dest:", dest, "destImmd:", destImmd)
	var value uint8
	switch src {
	case "A":
		value = cpu.Registers.A
	case "B":
		value = cpu.Registers.B
	case "C":
		value = cpu.Registers.C
	case "D":
		value = cpu.Registers.D
	case "E":
		value = cpu.Registers.E
	case "H":
		value = cpu.Registers.H
	case "L":
		value = cpu.Registers.L
	case "HL":

		if srcImmd == "True" {
			value = uint8(cpu.Registers.getHL())
		} else {
			value = cpu.Memory[cpu.Registers.getHL()]

		}
	case "DE":
		if srcImmd == "True" {
			value = uint8(cpu.Registers.getDE())
		} else {
			value = cpu.Memory[cpu.Registers.getDE()]
		}
	case "BC":
		if srcImmd == "True" {
			value = uint8(cpu.Registers.getBC())
		} else {
			value = cpu.Memory[cpu.Registers.getBC()]
		}
	case "n8":
		value = cpu.getImmediate8()
	case "a16":
		value = cpu.getImmediate16()
	}
	switch dest {
	case "A":
		cpu.Registers.A = value
	case "B":
		cpu.Registers.B = value
	case "C":
		cpu.Registers.C = value
	case "D":
		cpu.Registers.D = value
	case "E":
		cpu.Registers.E = value
	case "H":
		cpu.Registers.H = value
	case "L":
		cpu.Registers.L = value
	case "HL":
		if destImmd == "True" {
			cpu.Registers.setHL(uint16(value))
		} else {
			cpu.Memory[cpu.Registers.getHL()] = value
		}
	case "BC":
		if destImmd == "True" {
			cpu.Registers.setBC(uint16(value))
		} else {
			cpu.Memory[cpu.Registers.getBC()] = value
		}
	case "DE":
		if destImmd == "True" {
			cpu.Registers.setDE(uint16(value))
		} else {
			cpu.Memory[cpu.Registers.getDE()] = value
		}
	case "a16":
		cpu.Memory[cpu.getImmediate16()] = value

	}
}

func (cpu *CPU) execADD(operands []map[string]string, flags map[string]string) {
	dest := operands[0]["name"]
	destImmd := operands[0]["immediate"]
	src := operands[1]["name"]
	srcImmd := operands[1]["immediate"]
	fmt.Println("src:", src, "srcImmd:", srcImmd, "dest:", dest, "destImmd:", destImmd)
	var value1 uint8
	switch src {
	case "A":
		value1 = cpu.Registers.A
	case "B":
		value1 = cpu.Registers.B
	case "C":
		value1 = cpu.Registers.C
	case "D":
		value1 = cpu.Registers.D
	case "E":
		value1 = cpu.Registers.E
	case "H":
		value1 = cpu.Registers.H
	case "L":
		value1 = cpu.Registers.L
	case "HL":
		if srcImmd == "True" {
			value1 = uint8(cpu.Registers.getHL())
		} else {
			value1 = cpu.Memory[cpu.Registers.getHL()]
		}
	case "SP":
		value1 = uint8(cpu.Registers.SP)
	case "BC":
		value1 = uint8(cpu.Registers.getBC())
	case "DE":
		value1 = uint8(cpu.Registers.getDE())

	}

	var value2 uint8
	var res uint8
	switch dest {
	case "A":
		value2 = cpu.Registers.A
		res = value1 + value2
		cpu.Registers.A = res

		if flags["Z"] == "Z" {
			cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
		}
		//Set if carry from bit 3
		if flags["H"] == "H" {
			cpu.Registers.setFlag(flagH, (cpu.Registers.A&0x0F)+(value1&0x0F) > 0x0F)
		}

	case "HL":
		if srcImmd == "True" {
			value2 = uint8(cpu.Registers.getHL())
			res = value2 + value1
			cpu.Registers.setHL(uint16(res))

		} else {
			value2 = cpu.Memory[cpu.Registers.getHL()]
			res = value2 + value1
			cpu.Memory[cpu.Registers.getHL()] = res
		}
		//Set if carry from bit 3
		if flags["H"] == "H" {
			cpu.Registers.setFlag(flagH, ((cpu.Registers.getHL()&0x0FFF)+(uint16(value1)&0x0FFF)) > 0x0FFF)
		}

	}

	//if flags["Z"] == "Z" {
	//	cpu.Registers.setFlag(flagZ, true)
	//}
	if flags["N"] == "0" {
		cpu.Registers.setFlag(flagN, false)
	}
	//Set if carry from bit 7
	if flags["C"] == "C" {
		cpu.Registers.setFlag(flagC, res > 0x0F)
	}

}

func (cpu *CPU) execSUB(operands []map[string]string, flags map[string]string) {
	dest := operands[0]["name"]
	src := operands[1]["name"]
	srcImmd := operands[1]["immediate"]
	fmt.Println("src:", src, "srcImmd:", srcImmd, "dest:", dest)
	var value1 uint8
	switch src {
	case "A":
		value1 = cpu.Registers.A
	case "B":
		value1 = cpu.Registers.B
	case "C":
		value1 = cpu.Registers.C
	case "D":
		value1 = cpu.Registers.D
	case "E":
		value1 = cpu.Registers.E
	case "H":
		value1 = cpu.Registers.H
	case "L":
		value1 = cpu.Registers.L
	case "HL":
		if srcImmd == "True" {
			value1 = uint8(cpu.Registers.getHL())

		} else {
			value1 = cpu.Memory[cpu.Registers.getHL()]
		}
	case "n8":
		value1 = cpu.getImmediate8()

	}
	var value2 uint8
	var res uint8

	value2 = cpu.Registers.A
	res = value2 - value1
	cpu.Registers.A = res

	if flags["Z"] == "Z" {
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)
	}
	if flags["N"] == "1" {
		cpu.Registers.setFlag(flagN, true)
	}
	// Set if no borrow from bit 4
	if flags["H"] == "H" {
		cpu.Registers.setFlag(flagH, (value2&0x0F) < (value1&0x0F))
	}
	//Set if no borrow ?!
	if flags["C"] == "C" {
		cpu.Registers.setFlag(flagC, value2 > value1)
	}
}
func (cpu *CPU) setINCFlags(reg uint8, flags map[string]string) {
	//Set if carry from bit 3
	halfCarry := (reg & 0x0F) == 0x0F

	if flags["H"] == "H" {
		cpu.Registers.setFlag(flagH, halfCarry)
	}
	if flags["Z"] == "Z" {
		cpu.Registers.setFlag(flagZ, reg == 0)
	}
	cpu.Registers.setFlag(flagN, false)
	//cpu.Registers.setFlag(flagC, false)
}
func (cpu *CPU) execINC(operands []map[string]string, flags map[string]string) {
	operand := operands[0]["name"]
	operandImmd := operands[0]["immediate"]
	fmt.Println("operand:", operand, "immd:", operandImmd)
	switch operand {
	case "A":
		cpu.setINCFlags(cpu.Registers.A, flags)
		cpu.Registers.A++
	case "B":
		cpu.setINCFlags(cpu.Registers.B, flags)
		cpu.Registers.B++
	case "C":
		cpu.setINCFlags(cpu.Registers.C, flags)
		cpu.Registers.C++
	case "D":
		cpu.setINCFlags(cpu.Registers.D, flags)
		cpu.Registers.D++
	case "E":
		cpu.setINCFlags(cpu.Registers.E, flags)
		cpu.Registers.E++
	case "H":
		cpu.setINCFlags(cpu.Registers.H, flags)
		cpu.Registers.H++
	case "L":
		cpu.setINCFlags(cpu.Registers.L, flags)
		cpu.Registers.L++
	case "HL":
		if operandImmd == "True" {
			value := cpu.Registers.getHL()
			value++
			cpu.Registers.setHL(value)
		} else {
			value := cpu.Memory[cpu.Registers.getHL()]
			value++
			cpu.Memory[cpu.Registers.getHL()] = value
		}
	case "BC":
		value := cpu.Registers.getBC()
		value++
		cpu.Registers.setBC(value)
	case "DE":
		value := cpu.Registers.getDE()
		value++
		cpu.Registers.setDE(value)
	case "SP":
		cpu.Registers.SP++
	}

	//if flags["H"] == "H" {
	//	cpu.Registers.setFlag(flagH, )
	//}
}

func (cpu *CPU) setDECFlags(reg uint8, flags map[string]string) {
	//Set if no borrow from bit 4
	halfCarry := (reg & 0x0F) == 0x00
	if flags["H"] == "H" {
		cpu.Registers.setFlag(flagH, halfCarry)
	}
	if flags["Z"] == "Z" {
		cpu.Registers.setFlag(flagZ, reg == 0)
	}
	cpu.Registers.setFlag(flagN, true)
}

func (cpu *CPU) execDEC(operands []map[string]string, flags map[string]string) {
	operand := operands[0]["name"]
	operandImmd := operands[0]["immediate"]
	fmt.Println("operand:", operand, "immd:", operandImmd)
	switch operand {
	case "A":
		cpu.setDECFlags(cpu.Registers.A, flags)
		cpu.Registers.A--
	case "B":
		cpu.setDECFlags(cpu.Registers.B, flags)
		cpu.Registers.B--
	case "C":
		cpu.setDECFlags(cpu.Registers.C, flags)
		cpu.Registers.C--
	case "D":
		cpu.setDECFlags(cpu.Registers.D, flags)
		cpu.Registers.D--
	case "E":
		cpu.setDECFlags(cpu.Registers.E, flags)
		cpu.Registers.E--
	case "H":
		cpu.setDECFlags(cpu.Registers.H, flags)
		cpu.Registers.H--
	case "L":
		cpu.setDECFlags(cpu.Registers.L, flags)
		cpu.Registers.L--
	case "BC":
		value := cpu.Registers.getBC()
		value--
		cpu.Registers.setBC(value)
	case "DE":
		value := cpu.Registers.getDE()
		value--
		cpu.Registers.setDE(value)
	case "HL":
		if operandImmd == "True" {
			value := cpu.Registers.getHL()
			value--
			cpu.Registers.setHL(value)
		} else {
			value := cpu.Memory[cpu.Registers.getHL()]
			value--
			cpu.Memory[cpu.Registers.getHL()] = value
		}
	case "SP":
		cpu.Registers.SP--
	}
}

func (cpu *CPU) execPUSH(operands []map[string]string) {
	// Push register pair nn onto stack.
	// Decrement Stack Pointer (SP) twice
	var value uint16
	operand := operands[0]["name"]
	operandImmd := operands[0]["immediate"]
	fmt.Println("operand:", operand, "immd:", operandImmd)
	switch operand {
	case "BC":
		value = cpu.Registers.getBC()
	case "DE":
		value = cpu.Registers.getDE()
	case "HL":
		value = cpu.Registers.getHL()
	case "AF":
		value = cpu.Registers.getAF()
	}
	cpu.Memory[cpu.Registers.SP-1] = uint8((value & 0xFF) >> 8) //hi
	cpu.Memory[cpu.Registers.SP-2] = uint8(value & 0xFF)        //lo
	cpu.Registers.SP -= 2
}

func (cpu *CPU) execPOP(operands []map[string]string) {
	//Pop two bytes off stack into register pair nn.
	// Increment Stack Pointer (SP) twice.
	cpu.Registers.SP += 2
	high := uint16(cpu.Memory[cpu.Registers.SP])
	low := uint16(cpu.Memory[cpu.Registers.SP+1])
	operand := operands[0]["name"]
	operandImmd := operands[0]["immediate"]
	fmt.Println("operand:", operand, "immd:", operandImmd)
	value := (high << 8) | low // combine low and high
	switch operand {
	case "BC":
		cpu.Registers.setBC(value)
	case "DE":
		cpu.Registers.setDE(value)
	case "HL":
		cpu.Registers.setHL(value)
	case "AF":
		cpu.Registers.setAF(value) //TODO: add flags
	}

}

func (cpu *CPU) execAND(operands []map[string]string, flags map[string]string) {
	// Logically AND n with A, result in A
	operand := operands[1]["name"]
	operandImmd := operands[1]["immediate"]
	var value uint8
	switch operand {
	case "A":
		value = cpu.Registers.A
	case "B":
		value = cpu.Registers.B
	case "C":
		value = cpu.Registers.C
	case "D":
		value = cpu.Registers.D
	case "E":
		value = cpu.Registers.E
	case "H":
		value = cpu.Registers.H
	case "L":
		value = cpu.Registers.L
	case "n8":
		value = cpu.getImmediate8()
	case "HL":
		if operandImmd == "False" {
			value = cpu.Memory[cpu.Registers.getHL()]
		}
	}
	cpu.Registers.A &= value
	if flags["Z"] == "Z" {
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)

	}
	cpu.Registers.setFlag(flagN, false) //reset
	cpu.Registers.setFlag(flagH, true)  //set
	cpu.Registers.setFlag(flagC, false) //reset

}

func (cpu *CPU) execOR(operands []map[string]string, flags map[string]string) {
	//Logical OR n with register A, result in A.
	operand := operands[1]["name"]
	operandImmd := operands[1]["immediate"]
	var value uint8
	switch operand {
	case "A":
		value = cpu.Registers.A
	case "B":
		value = cpu.Registers.B
	case "C":
		value = cpu.Registers.C
	case "D":
		value = cpu.Registers.D
	case "E":
		value = cpu.Registers.E
	case "H":
		value = cpu.Registers.H
	case "L":
		value = cpu.Registers.L
	case "n8":
		value = cpu.getImmediate8()
	case "HL":
		if operandImmd == "False" {
			value = cpu.Memory[cpu.Registers.getHL()]
		}
	}
	cpu.Registers.A |= value
	if flags["Z"] == "Z" {
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)

	}
	cpu.Registers.setFlag(flagN, false) //reset
	cpu.Registers.setFlag(flagH, false) //reset
	cpu.Registers.setFlag(flagC, false) //reset

}

func (cpu *CPU) execXOR(operands []map[string]string, flags map[string]string) {
	//Logical exclusive OR n with register A, result in A.
	operand := operands[1]["name"]
	operandImmd := operands[1]["immediate"]
	var value uint8
	switch operand {
	case "A":
		value = cpu.Registers.A
	case "B":
		value = cpu.Registers.B
	case "C":
		value = cpu.Registers.C
	case "D":
		value = cpu.Registers.D
	case "E":
		value = cpu.Registers.E
	case "H":
		value = cpu.Registers.H
	case "L":
		value = cpu.Registers.L
	case "n8":
		value = cpu.getImmediate8()
	case "HL":
		if operandImmd == "False" {
			value = cpu.Memory[cpu.Registers.getHL()]
		}
	}
	cpu.Registers.A ^= value
	if flags["Z"] == "Z" {
		cpu.Registers.setFlag(flagZ, cpu.Registers.A == 0)

	}
	cpu.Registers.setFlag(flagN, false) //reset
	cpu.Registers.setFlag(flagH, false) //reset
	cpu.Registers.setFlag(flagC, false) //reset

}

func fetchOpcodesFromJSON() (map[string]map[string][]map[string]string, error) {
	functionName := "import opcodeParser; print(opcodeParser.parse_operations())"
	fmt.Println("Function:", functionName)
	cmd := exec.Command("python3", "-c", functionName)
	instr, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return nil, err
	}

	var instructions map[string]map[string][]map[string]string
	err = json.Unmarshal(instr, &instructions)
	if err != nil {
		println(err.Error())
	}

	return instructions, nil
}

func (cpu *CPU) fetchOpcode() uint8 {
	opcode := cpu.Memory[cpu.Registers.PC]
	fmt.Println("Opcode in fetch:", opcode)
	cpu.Registers.PC++
	return opcode
}

func (cpu *CPU) execCommand() {
	opcode := cpu.fetchOpcode()
	//fmt.Printf("Executing opcode: 0x%X\n", opcode)
	fmt.Printf("Executing opcode: 0x%02X\n", opcode)

	instructions, err := fetchOpcodesFromJSON()
	if err != nil {
		fmt.Println("Error fetching opcodes:", err.Error())
		return
	}
	//fmt.Println("Instructions:", instructions)
	for instr, table := range instructions {
		if details, found := table[fmt.Sprintf("0x%02X", opcode)]; found {
			operands := details[:len(details)-1]
			flags := details[len(details)-1]

			switch instr {
			case "LD":
				cpu.execLD(operands)
			case "ADD":
				cpu.execADD(operands, flags)
			case "SUB":
				cpu.execSUB(operands, flags)
			case "INC":
				cpu.execINC(operands, flags)
			case "DEC":
				cpu.execDEC(operands, flags)
			case "PUSH":
				cpu.execPUSH(operands)
			case "POP":
				cpu.execPOP(operands)
			default:
				fmt.Printf("Unhandled operation: %s for opcode 0x%X\n", instr, opcode)
			}
			return
		}
	}

	fmt.Printf("Unknown opcode: 0x%X\n", opcode)
}

//func main() {
//	//register := Registers{}
//	//register.setBC(0x0001)
//	//fmt.Println(register.getBC())
//	//
//	//register.setFlag(flagZ, true)
//	//fmt.Println(register.getAF())
//	//cpu := CPU{}
//	//instr, _ := fetchOpcodesFromJSON("op_LD")
//	//fmt.Println(instr)
//	//opcodes, _ := fetchAllOpcodesFromJSON()
//	//fmt.Println("Opcodes:", opcodes)
//	cpu := CPU{}
//	cpu.loadOperations()
//	//cpu.execCommand()
//}
