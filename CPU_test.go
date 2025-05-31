package main

import (
	"fmt"
	"testing"
)

func TestSetBC(t *testing.T) {
	register := Registers{}
	register.setBC(0b1)
	reg := register.getBC()
	if reg != 0b1 {
		t.Error("Expected 0b1, got ", register.getBC())
	}
}

func TestSetDE(t *testing.T) {
	register := Registers{}
	register.setDE(0b1)
	// // fmt.Println(register.getDE())
	if register.getDE() != 0b1 {
		t.Error("Expected 0b1, got ", register.getDE())
	}
}

func TestSetHL(t *testing.T) {
	register := Registers{}
	register.setHL(0b1)
	// // fmt.Println(register.getHL())
	if register.getHL() != 0b1 {
		t.Error("Expected 0b1, got ", register.getHL())
	}
}

func TestSetAF(t *testing.T) {
	register := Registers{}
	register.setAF(0b1)
	if register.getAF() != 0b1 {
		t.Error("Expected 0b1, got ", register.getAF())
	}
}

// LD A,B
func TestExecLDAB(t *testing.T) {
	//cpu := CPU{
	//	Registers: Registers{
	//		A:  0b1,
	//		B:  0b10,
	//		PC: 10,
	//	},
	//	Memory: [65536]uint8(make([]uint8, 65536)),
	//}
	cpu := NewCPU()

	cpu.Registers.A = 0b1
	cpu.Registers.B = 0b10

	cpu.Memory[cpu.Registers.PC] = 0b1111000
	cpu.execOpcodes()
	// // fmt.Println("LD A, B TEST")
	// // fmt.Println(cpu.Registers.A)
	if cpu.Registers.A != 0b10 {
		t.Error("Expected 0b10, got ", cpu.Registers.A)
	}
}

// LD B, C
func TestExecLD_B_C(t *testing.T) {
	cpu := NewCPU()

	cpu.Registers.C = 0b10
	cpu.Registers.B = 0b1
	cpu.Memory[cpu.Registers.PC] = 0b1000001
	cpu.execOpcodes()
	if cpu.Registers.B != 0b10 {
		t.Error("Expected 0b10, got ", cpu.Registers.B)
	}
}

// LD A, n
func TestExecLD_A(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b10
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Memory[cpu.Registers.PC] = 0b1111000 // LD A, B
	cpu.execOpcodes()
	if cpu.Registers.A != 0b1 {
		t.Error("Expected 0b1, got ", cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1111001 // LD A, C
	cpu.execOpcodes()
	if cpu.Registers.A != 0b10 {
		t.Error("Expected 0b10, got ", cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1111010 // LD A, D
	cpu.execOpcodes()
	if cpu.Registers.A != 0b100 {
		t.Error("Expected 0b100, got ", cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1111011 // LD A, E
	cpu.execOpcodes()
	if cpu.Registers.A != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1111100 // LD A, H
	cpu.execOpcodes()
	if cpu.Registers.A != 0b110 {
		t.Error("Expected 0b110, got ", cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1111101 // LD A, L
	cpu.execOpcodes()
	if cpu.Registers.A != 0b111 {
		t.Error("Expected 0b111, got ", cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.getBC()] = 0b100
	cpu.Memory[cpu.Registers.PC] = 0b1010 // LD A, [BC]
	cpu.execOpcodes()
	if cpu.Registers.A != 0b100 {
		t.Error("Expected 0b100, got ", cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.getDE()] = 0b111
	cpu.Memory[cpu.Registers.PC] = 0b11010 // LD A, [DE]
	cpu.execOpcodes()
	if cpu.Registers.A != 0b111 {
		t.Error("Expected 0b111, got ", cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.getHL()] = 0b11
	cpu.Memory[cpu.Registers.PC] = 0b1111110 // LD A, [HL]
	cpu.execOpcodes()
	if cpu.Registers.A != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b111110 // LD A, imm8
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.execOpcodes()
	if cpu.Registers.A != 0b1 {
		t.Error("Expected 0b1, got ", cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1111111 // LD A, A
	cpu.execOpcodes()
	if cpu.Registers.A != 0b1 {
		t.Error("Expected 0b1, got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b111110 // LD A, [imm8]
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Memory[0b001100000001] = 0b1
	cpu.execOpcodes()
	if cpu.Registers.A != 0b1 {
		t.Error("Expected 0b1, got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11111010 // LD A, [imm16]
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Memory[0b001100000001] = 0b1
	cpu.execOpcodes()
	if cpu.Registers.A != 0b1 {
		t.Error("Expected 0b1, got ", cpu.Registers.A)
	}
}

// LD B, n
func TestExecLD_B(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b10
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Memory[cpu.Registers.PC] = 0b1000111 // LD B, A
	cpu.execOpcodes()
	if cpu.Registers.B != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.B)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1000001 // LD B, C
	cpu.execOpcodes()
	if cpu.Registers.B != 0b10 {
		t.Error("Expected 0b10, got ", cpu.Registers.B)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1000010 // LD B, D
	cpu.execOpcodes()
	if cpu.Registers.B != 0b100 {
		t.Error("Expected 0b100, got ", cpu.Registers.B)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1000011 // LD B, E
	cpu.execOpcodes()
	if cpu.Registers.B != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.B)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1000100 // LD B, H
	cpu.execOpcodes()
	if cpu.Registers.B != 0b110 {
		t.Error("Expected 0b110, got ", cpu.Registers.B)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1000101 // LD B, L
	cpu.execOpcodes()
	if cpu.Registers.B != 0b111 {
		t.Error("Expected 0b111, got ", cpu.Registers.B)
	}

	cpu.Memory[cpu.Registers.getHL()] = 0b11
	cpu.Memory[cpu.Registers.PC] = 0b1000110 // LD B, [HL]
	cpu.execOpcodes()
	if cpu.Registers.B != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.B)
	}
	cpu.Memory[cpu.Registers.PC] = 0b110 // LD B, imm8
	cpu.Memory[cpu.Registers.PC+1] = 0b11
	cpu.execOpcodes()
	//// fmt.Println(cpu.getImmediate8())

	if cpu.Registers.B != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.B)
	}

	cpu.Memory[cpu.Registers.PC] = 0b1000000 // LD B, B
	cpu.execOpcodes()

	if cpu.Registers.B != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.B)
	}
}

// LD C, n
func TestExecLD_C(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b10
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Memory[cpu.Registers.PC] = 0b1001111 // LD C, A
	cpu.execOpcodes()
	if cpu.Registers.C != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.C)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1001000 // LD C, B
	cpu.execOpcodes()
	if cpu.Registers.C != 0b1 {
		t.Error("Expected 0b1, got ", cpu.Registers.C)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1001010 // LD C, D
	cpu.execOpcodes()
	if cpu.Registers.C != 0b100 {
		t.Error("Expected 0b100, got ", cpu.Registers.C)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1001011 // LD C, E
	cpu.execOpcodes()
	if cpu.Registers.C != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.C)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1001100 // LD B, H
	cpu.execOpcodes()
	if cpu.Registers.C != 0b110 {
		t.Error("Expected 0b110, got ", cpu.Registers.C)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1001101 // LD C, L
	cpu.execOpcodes()
	if cpu.Registers.C != 0b111 {
		t.Error("Expected 0b111, got ", cpu.Registers.C)
	}

	cpu.Memory[cpu.Registers.getHL()] = 0b11
	cpu.Memory[cpu.Registers.PC] = 0b1001110 // LD C, [HL]
	cpu.execOpcodes()
	if cpu.Registers.C != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.C)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1110 // LD C, imm8
	cpu.Memory[cpu.Registers.PC+1] = 0b11
	cpu.execOpcodes()

	if cpu.Registers.C != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.C)
	}

	cpu.Memory[cpu.Registers.PC] = 0b1001001 // LD C, C
	cpu.execOpcodes()
	if cpu.Registers.C != 0b11 {
		t.Error("Expected 0b111, 0b11 ", cpu.Registers.C)
	}
}

// LD D, n
func TestExecLD_D(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b10
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Memory[cpu.Registers.PC] = 0b1010111 // LD D, A
	cpu.execOpcodes()
	if cpu.Registers.D != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.D)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1010000 // LD D, B
	cpu.execOpcodes()
	if cpu.Registers.D != 0b1 {
		t.Error("Expected 0b1, got ", cpu.Registers.D)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1010001 // LD D, C
	cpu.execOpcodes()
	if cpu.Registers.D != 0b10 {
		t.Error("Expected 0b100, got ", cpu.Registers.D)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1010011 // LD D, E
	cpu.execOpcodes()
	if cpu.Registers.D != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.D)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1010100 // LD D, H
	cpu.execOpcodes()
	if cpu.Registers.D != 0b110 {
		t.Error("Expected 0b110, got ", cpu.Registers.D)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1010101 // LD D, L
	cpu.execOpcodes()
	if cpu.Registers.D != 0b111 {
		t.Error("Expected 0b111, got ", cpu.Registers.D)
	}

	cpu.Memory[cpu.Registers.getHL()] = 0b11
	cpu.Memory[cpu.Registers.PC] = 0b1010110 // LD D, [HL]
	cpu.execOpcodes()
	if cpu.Registers.D != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.D)
	}
	cpu.Memory[cpu.Registers.PC] = 0b10110 // LD D, imm8
	cpu.Memory[cpu.Registers.PC+1] = 0b11
	cpu.execOpcodes()

	if cpu.Registers.D != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.D)
	}

	cpu.Memory[cpu.Registers.PC] = 0b1010010 // LD D, D
	cpu.execOpcodes()
	if cpu.Registers.D != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.D)
	}
}

// LD E, n
func TestExecLD_E(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b10
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Memory[cpu.Registers.PC] = 0b1011111 // LD E, A
	cpu.execOpcodes()
	if cpu.Registers.E != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.E)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1011000 // LD E, B
	cpu.execOpcodes()
	if cpu.Registers.E != 0b1 {
		t.Error("Expected 0b1, got ", cpu.Registers.E)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1011001 // LD E, C
	cpu.execOpcodes()
	if cpu.Registers.E != 0b10 {
		t.Error("Expected 0b100, got ", cpu.Registers.E)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1011010 // LD E, D
	cpu.execOpcodes()
	if cpu.Registers.E != 0b100 {
		t.Error("Expected 0b100, got ", cpu.Registers.E)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1011100 // LD E, H
	cpu.execOpcodes()
	if cpu.Registers.E != 0b110 {
		t.Error("Expected 0b110, got ", cpu.Registers.E)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1011101 // LD E, L
	cpu.execOpcodes()
	if cpu.Registers.E != 0b111 {
		t.Error("Expected 0b111, got ", cpu.Registers.E)
	}

	cpu.Memory[cpu.Registers.getHL()] = 0b11
	cpu.Memory[cpu.Registers.PC] = 0b1011110 // LD E, [HL]
	cpu.execOpcodes()
	if cpu.Registers.E != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.E)
	}
	cpu.Memory[cpu.Registers.PC] = 0b11110 // LD E, imm8
	cpu.Memory[cpu.Registers.PC+1] = 0b11
	cpu.execOpcodes()

	if cpu.Registers.E != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.E)
	}

	cpu.Memory[cpu.Registers.PC] = 0b1011011 // LD E, E
	cpu.execOpcodes()
	if cpu.Registers.E != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.E)
	}
}

// LD H, n
func TestExecLD_H(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b10
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Memory[cpu.Registers.PC] = 0b1100111 // LD H, A
	cpu.execOpcodes()
	if cpu.Registers.H != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.H)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1100000 // LD H, B
	cpu.execOpcodes()
	if cpu.Registers.H != 0b1 {
		t.Error("Expected 0b1, got ", cpu.Registers.H)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1100001 // LD H, C
	cpu.execOpcodes()
	if cpu.Registers.H != 0b10 {
		t.Error("Expected 0b100, got ", cpu.Registers.H)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1100010 // LD H, D
	cpu.execOpcodes()
	if cpu.Registers.H != 0b100 {
		t.Error("Expected 0b100, got ", cpu.Registers.H)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1100011 // LD H, E
	cpu.execOpcodes()
	if cpu.Registers.H != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.H)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1100101 // LD H, L
	cpu.execOpcodes()
	if cpu.Registers.H != 0b111 {
		t.Error("Expected 0b111, got ", cpu.Registers.H)
	}

	cpu.Memory[cpu.Registers.getHL()] = 0b11
	cpu.Memory[cpu.Registers.PC] = 0b1100110 // LD H, [HL]
	cpu.execOpcodes()
	if cpu.Registers.H != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.H)
	}
	cpu.Memory[cpu.Registers.PC] = 0b100110 // LD H, imm8
	cpu.Memory[cpu.Registers.PC+1] = 0b11
	cpu.execOpcodes()

	if cpu.Registers.H != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.H)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1100100 // LD H, H
	cpu.execOpcodes()
	if cpu.Registers.H != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.H)
	}
}

// LD L, n
func TestExecLD_L(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b10
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Memory[cpu.Registers.PC] = 0b1101111 // LD L, A
	cpu.execOpcodes()
	if cpu.Registers.L != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.L)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1101000 // LD L, B
	cpu.execOpcodes()
	if cpu.Registers.L != 0b1 {
		t.Error("Expected 0b1, got ", cpu.Registers.L)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1101001 // LD L, C
	cpu.execOpcodes()
	if cpu.Registers.L != 0b10 {
		t.Error("Expected 0b100, got ", cpu.Registers.L)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1101010 // LD L, D
	cpu.execOpcodes()
	if cpu.Registers.L != 0b100 {
		t.Error("Expected 0b100, got ", cpu.Registers.L)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1101011 // LD L, E
	cpu.execOpcodes()
	if cpu.Registers.L != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.L)
	}
	cpu.Memory[cpu.Registers.PC] = 0b1101100 // LD L, H
	cpu.execOpcodes()
	if cpu.Registers.L != 0b110 {
		t.Error("Expected 0b110, got ", cpu.Registers.L)
	}

	cpu.Memory[cpu.Registers.getHL()] = 0b11
	cpu.Memory[cpu.Registers.PC] = 0b1101110 // LD L, [HL]
	cpu.execOpcodes()
	if cpu.Registers.L != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.L)
	}
	cpu.Memory[cpu.Registers.PC] = 0b101110 // LD L, imm8
	cpu.Memory[cpu.Registers.PC+1] = 0b11
	cpu.execOpcodes()

	if cpu.Registers.L != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.L)
	}

	cpu.Memory[cpu.Registers.PC] = 0b1101101 // LD L, L
	cpu.execOpcodes()
	if cpu.Registers.L != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.L)
	}
}

// LD [HL], n
func TestExecLD_HL(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b10
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111

	cpu.Memory[cpu.Registers.getHL()] = 0b1111
	cpu.Memory[cpu.Registers.PC] = 0b1110111 // LD [HL], A
	cpu.execOpcodes()
	if cpu.Memory[cpu.Registers.getHL()] != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Memory[cpu.Registers.getHL()])
	}
	cpu.Memory[cpu.Registers.PC] = 0b1110000 // LD [HL], B
	cpu.execOpcodes()
	if cpu.Memory[cpu.Registers.getHL()] != 0b1 {
		t.Error("Expected 0b1, got ", cpu.Memory[cpu.Registers.getHL()])
	}
	cpu.Memory[cpu.Registers.PC] = 0b1110001 // LD [HL], C
	cpu.execOpcodes()
	if cpu.Memory[cpu.Registers.getHL()] != 0b10 {
		t.Error("Expected 0b100, got ", cpu.Memory[cpu.Registers.getHL()])
	}
	cpu.Memory[cpu.Registers.PC] = 0b1110010 // LD [HL], D
	cpu.execOpcodes()
	if cpu.Memory[cpu.Registers.getHL()] != 0b100 {
		t.Error("Expected 0b100, got ", cpu.Memory[cpu.Registers.getHL()])
	}
	cpu.Memory[cpu.Registers.PC] = 0b1110011 // LD [HL], E
	cpu.execOpcodes()
	if cpu.Memory[cpu.Registers.getHL()] != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Memory[cpu.Registers.getHL()])
	}
	cpu.Memory[cpu.Registers.PC] = 0b1110100 // LD [HL], H
	cpu.execOpcodes()
	if cpu.Memory[cpu.Registers.getHL()] != 0b110 {
		t.Error("Expected 0b110, got ", cpu.Memory[cpu.Registers.getHL()])
	}

	cpu.Memory[cpu.Registers.getHL()] = 0b11
	cpu.Memory[cpu.Registers.PC] = 0b1110101 // LD [HL], L
	cpu.execOpcodes()
	if cpu.Memory[cpu.Registers.getHL()] != 0b111 {
		t.Error("Expected 0b111, got ", cpu.Memory[cpu.Registers.getHL()])
	}
	cpu.Memory[cpu.Registers.PC] = 0b110110 // LD L, imm8
	cpu.Memory[cpu.Registers.PC+1] = 0b11
	cpu.execOpcodes()

	if cpu.Memory[cpu.Registers.getHL()] != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Memory[cpu.Registers.getHL()])
	}
}

// LD HL, n16
func TestExecLD_HL_n16(t *testing.T) {

	cpu := NewCPU()
	cpu.Registers.B = 0b10
	cpu.Registers.H = 0b10100
	cpu.Registers.L = 0b110110
	cpu.Memory[cpu.Registers.PC] = 0b100001
	cpu.Memory[cpu.Registers.PC+1] = 0b11
	cpu.execOpcodes()
	// // fmt.Println("LD [HL], B TEST")

	if cpu.Registers.getHL() != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.getHL())
	}
}

// ADD A,n
func TestExecADDA(t *testing.T) {

	cpu := NewCPU()
	cpu.Registers.A = 0b1
	cpu.Registers.B = 0b10
	cpu.Registers.C = 0b10
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	var res uint8

	cpu.Memory[cpu.Registers.PC] = 0b10000000 // ADD A, B
	cpu.execOpcodes()
	res = uint8(0b1 + 0b10)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.PC] = 0b10000001 // ADD A, C
	cpu.execOpcodes()
	res += 0b10

	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.PC] = 0b10000010 // ADD A, D
	cpu.execOpcodes()
	res += 0b100
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.PC] = 0b10000011 // ADD A, E
	cpu.execOpcodes()
	res += 0b101
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.PC] = 0b10000100 // ADD A, H
	cpu.execOpcodes()
	res += 0b110
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.PC] = 0b10000101 // ADD A, L
	cpu.execOpcodes()
	res += 0b111
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10000110 // ADD A, [HL]
	cpu.Memory[cpu.Registers.getHL()] = 0b111
	cpu.execOpcodes()
	res += 0b111
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10000111 // ADD A, A
	cpu.execOpcodes()
	res += res
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11000110 // ADD A, imm8
	cpu.Memory[cpu.Registers.PC+1] = 0b1

	cpu.execOpcodes()
	res += 0b1
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

}

// ADC A,n
func TestExecADCA(t *testing.T) {

	cpu := NewCPU()
	cpu.Registers.A = 0b1
	cpu.Registers.B = 0b10
	cpu.Registers.C = 0b10
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.F = 0b000010000 // flagC on
	var res uint8

	cpu.Memory[cpu.Registers.PC] = 0b10001000 // ADC A, B
	cpu.execOpcodes()
	res = uint8(0b1 + 0b10 + 1)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.PC] = 0b10001001 // ADC A, C
	cpu.execOpcodes()
	res += 0b10 + 1

	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.PC] = 0b10001010 // ADC A, D
	cpu.execOpcodes()
	res += 0b100 + 1
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.PC] = 0b10001011 // ADC A, E
	cpu.execOpcodes()
	res += 0b101 + 1
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.PC] = 0b10001100 // ADC A, H
	cpu.execOpcodes()
	res += 0b110 + 1
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.PC] = 0b10001101 // ADC A, L
	cpu.execOpcodes()
	res += 0b111 + 1
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.getHL()] = 0b1
	cpu.Memory[cpu.Registers.PC] = 0b10001110 // ADC A, [HL]
	cpu.execOpcodes()
	res += 0b1 + uint8(1)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.PC] = 0b11001110 // ADC A, imm8
	cpu.Memory[cpu.Registers.PC+1] = 0b11010111
	cpu.execOpcodes()
	res += 0b11010111 + uint8(1)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.PC] = 0b10001111 // ADC A, A
	cpu.execOpcodes()
	res += res + 1

	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}
}

// SUB A, n
func TestSUB_A(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11110 //30
	cpu.Registers.B = 0b10
	cpu.Registers.C = 0b10
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111

	var res uint8
	cpu.Memory[cpu.Registers.PC] = 0b10010000 // SUB A, B
	cpu.execOpcodes()
	res = uint8(0b11110 - 0b10)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.PC] = 0b10010001 // SUB A, C
	cpu.execOpcodes()
	res -= 0b10
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.PC] = 0b10010010 // SUB A, D
	cpu.execOpcodes()
	res -= 0b100
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.PC] = 0b10010011 // SUB A, E
	cpu.execOpcodes()
	res -= 0b101
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.PC] = 0b10010100 // SUB A, H
	cpu.execOpcodes()
	res -= 0b110
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.PC] = 0b10010101 // SUB A, L
	cpu.execOpcodes()
	res -= 0b111
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.getHL()] = 0b1
	cpu.Memory[cpu.Registers.PC] = 0b10010110 // SUB A, [HL]
	cpu.execOpcodes()
	res -= 0b1
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}
	cpu.Memory[cpu.Registers.PC] = 0b11010110 // SUB A, imm8
	cpu.Memory[cpu.Registers.PC+1] = 0b11010111
	cpu.execOpcodes()
	res -= 0b11010111
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10010111 // SUB A, A
	cpu.execOpcodes()

	if cpu.Registers.A != 0 {
		t.Errorf("Expected 0, got 0b%b ", cpu.Registers.A)
	}
}

// SBC A, n
func TestSBC_A(t *testing.T) {

	cpu := NewCPU()
	cpu.Registers.A = 0b11110 //30
	cpu.Registers.B = 0b10
	cpu.Registers.C = 0b10
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.F = 0b000010000 // flagC on
	var res uint8
	cpu.Memory[cpu.Registers.PC] = 0b10011000 // SBC A, B
	//fmt.Println("FLAGC", cpu.Registers.getFlag(flagC))

	cpu.execOpcodes()
	res = uint8(0b11110 - 0b10 - 1)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.PC] = 0b10011001 // SBC A, C
	//fmt.Println("FLAGC", cpu.Registers.getFlag(flagC))

	cpu.execOpcodes()

	res -= uint8(0b10 + uint8(1))
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.PC] = 0b10011010 // SBC A, D
	//fmt.Println("FLAGC", cpu.Registers.getFlag(flagC))

	cpu.execOpcodes()
	res -= 0b100 + uint8(1)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.PC] = 0b10011011 // SBC A, E
	//fmt.Println("FLAGC", cpu.Registers.getFlag(flagC))

	cpu.execOpcodes()
	res -= 0b101 + uint8(1)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.PC] = 0b10011100 // SBC A, H
	//fmt.Println("FLAGC", cpu.Registers.getFlag(flagC))

	cpu.execOpcodes()
	res -= 0b110 + uint8(1)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.PC] = 0b10011101 // SBC A, L
	//fmt.Println("FLAGC", cpu.Registers.getFlag(flagC))

	cpu.execOpcodes()
	res -= 0b111 + uint8(1)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.getHL()] = 0b1
	cpu.Memory[cpu.Registers.PC] = 0b10011110 // SBC A, [HL]
	//fmt.Println("FLAGC", cpu.Registers.getFlag(flagC))

	cpu.execOpcodes()
	res -= 0b1 + uint8(1)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.PC] = 0b11011110 // SBC A, imm8
	cpu.Memory[cpu.Registers.PC+1] = 0b11010111
	//fmt.Println("FLAGC", cpu.Registers.getFlag(flagC))

	cpu.execOpcodes()
	res -= 0b11010111 + uint8(1)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10011111 // SBC A, A
	//fmt.Println("FLAGC", cpu.Registers.getFlag(flagC))
	cpu.execOpcodes()
	res = uint8(0b11110 - 0b11110 - 0)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}
}

// ADD HL, n
func TestExecADD_HL_BC(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.B = 0b10
	cpu.Registers.C = 0b110
	cpu.Registers.H = 0b1
	cpu.Registers.L = 0b10
	var res uint16
	cpu.Registers.setHL(0b10)
	cpu.Registers.setBC(0b100)
	cpu.Registers.setDE(0b111)
	cpu.Memory[cpu.Registers.PC] = 0b1001 // ADD HL. BC
	cpu.execOpcodes()
	res = 0b110
	if cpu.Registers.getHL() != res {
		t.Errorf("Expected 0b%b got 0b%b ", res, cpu.Registers.getHL())
	}

	cpu.Memory[cpu.Registers.PC] = 0b11001 // ADD HL. DE
	cpu.execOpcodes()
	res += 0b111
	if cpu.Registers.getHL() != res {
		t.Errorf("Expected 0b%b got 0b%b ", res, cpu.Registers.getHL())
	}

	cpu.Memory[cpu.Registers.PC] = 0b101001 // ADD HL. HL
	cpu.execOpcodes()
	res += res
	if cpu.Registers.getHL() != res {
		t.Errorf("Expected 0b%b got 0b%b ", res, cpu.Registers.getHL())
	}

	cpu.Memory[cpu.Registers.PC] = 0b101001 // ADD HL. HL
	cpu.execOpcodes()
	res += res
	if cpu.Registers.getHL() != res {
		t.Errorf("Expected 0b%b got 0b%b ", res, cpu.Registers.getHL())
	}
	cpu.Registers.SP = 0b1
	cpu.Memory[cpu.Registers.PC] = 0b111001 // ADD HL. SP
	cpu.execOpcodes()
	res += 0b1
	if cpu.Registers.getHL() != res {
		t.Errorf("Expected 0b%b got 0b%b ", res, cpu.Registers.getHL())
	}

}

// INC B
func TestINC_B(t *testing.T) {
	cpu := NewCPU()

	cpu.Registers.B = 0b1
	cpu.Memory[cpu.Registers.PC] = 0b100
	cpu.execOpcodes()
	// // fmt.Println("INC B TEST")
	// // fmt.Println(cpu.Registers.B)

	if cpu.Registers.B != 0b10 {
		t.Error("Expected 0b10, got ", cpu.Registers.B)
	}
}

// INC C
func TestINC_C(t *testing.T) {
	cpu := NewCPU()

	cpu.Registers.C = 0b1
	cpu.Memory[cpu.Registers.PC] = 0b1100
	cpu.execOpcodes()

	if cpu.Registers.C != 0b10 {
		t.Error("Expected 0b10, got ", cpu.Registers.C)
	}
}

// INC D
func TestINC_D(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.D = 0b1
	cpu.Memory[cpu.Registers.PC] = 0b10100
	cpu.execOpcodes()

	if cpu.Registers.D != 0b10 {
		t.Error("Expected 0b10, got ", cpu.Registers.D)
	}
}

// INC E
func TestINC_E(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.E = 0b1
	cpu.Memory[cpu.Registers.PC] = 0b11100
	cpu.execOpcodes()

	if cpu.Registers.E != 0b10 {
		t.Error("Expected 0b10, got ", cpu.Registers.E)
	}
}

// INC H
func TestINC_H(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.H = 0b1
	cpu.Memory[cpu.Registers.PC] = 0b100100
	cpu.execOpcodes()

	if cpu.Registers.H != 0b10 {
		t.Error("Expected 0b10, got ", cpu.Registers.H)
	}
}

// INC L
func TestINC_L(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.L = 0b1
	cpu.Memory[cpu.Registers.PC] = 0b101100
	cpu.execOpcodes()

	if cpu.Registers.L != 0b10 {
		t.Error("Expected 0b10, got ", cpu.Registers.L)
	}
}

// INC BC
func TestINC_BC(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b10
	cpu.Registers.setBC(0b100)

	cpu.Memory[cpu.Registers.PC] = 0b11
	cpu.execOpcodes()
	// // fmt.Println("INC BC TEST")
	// // fmt.Println(cpu.Registers.getBC())
	if cpu.Registers.getBC() != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.getBC())
	}
}

// INC DE
func TestINC_DE(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.D = 0b1
	cpu.Registers.E = 0b10
	cpu.Memory[cpu.Registers.PC] = 0b10011
	fmt.Println("DE", cpu.Registers.getDE())
	val := cpu.Registers.getDE()
	val++
	cpu.execOpcodes()
	if cpu.Registers.getDE() != val {
		t.Error("Expected 0x14, got ", cpu.Registers.getDE())
	}
}

// INC HL
func TestINC_HL(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.H = 0b1
	cpu.Registers.L = 0b10
	cpu.Registers.setHL(0b100)

	cpu.Memory[cpu.Registers.PC] = 0b100011
	cpu.execOpcodes()
	// // fmt.Println("INC HL TEST")
	// // fmt.Println(cpu.Registers.getHL())
	if cpu.Registers.getHL() != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.getHL())
	}
}

// INC [HL]
func TestINC_HLmem(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.L = 0b10
	cpu.Registers.H = 0b1
	cpu.Registers.setHL(0b100)
	hladdr := cpu.Registers.getHL()
	cpu.Memory[hladdr] = 0b100

	cpu.Memory[cpu.Registers.PC] = 0b110100
	cpu.execOpcodes()
	// // fmt.Println(" INC [HL] TEST")
	// // fmt.Println(cpu.Memory[cpu.Registers.getHL()])
	if cpu.Memory[cpu.Registers.getHL()] != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Memory[cpu.Registers.getHL()])
	}
}

// DEC A
func TestDEC_A(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b110
	cpu.Memory[cpu.Registers.PC] = 0b111101
	cpu.execOpcodes()

	if cpu.Registers.A != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.A)
	}
}

// DEC B
func TestDEC_B(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.B = 0b110

	cpu.Memory[cpu.Registers.PC] = 0b101
	cpu.execOpcodes()

	if cpu.Registers.B != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.B)
	}
}

// DEC C
func TestDEC_C(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.C = 0b110
	cpu.Memory[cpu.Registers.PC] = 0b1101
	cpu.execOpcodes()

	if cpu.Registers.C != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.C)
	}
}

// DEC D
func TestDEC_D(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.D = 0b110
	cpu.Memory[cpu.Registers.PC] = 0b10101
	cpu.execOpcodes()

	if cpu.Registers.D != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.D)
	}
}

// DEC E
func TestDEC_E(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.E = 0b110
	cpu.Memory[cpu.Registers.PC] = 0b11101
	cpu.execOpcodes()

	if cpu.Registers.E != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.E)
	}
}

// DEC H
func TestDEC_H(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.H = 0b110
	cpu.Memory[cpu.Registers.PC] = 0b100101
	cpu.execOpcodes()

	if cpu.Registers.H != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.H)
	}
}

// DEC L
func TestDEC_L(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.L = 0b110
	cpu.Memory[cpu.Registers.PC] = 0b101101
	cpu.execOpcodes()

	if cpu.Registers.L != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.L)
	}
}

// DEC DE
func TestDEC_DE(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.D = 0b1
	cpu.Registers.E = 0b10
	cpu.Registers.setDE(0b100)

	cpu.Memory[cpu.Registers.PC] = 0b11011
	cpu.execOpcodes()
	// // fmt.Println("DEC DE TEST")
	// // fmt.Println(cpu.Registers.getDE())
	if cpu.Registers.getDE() != 0b11 {
		t.Error("Expected 0b101, got ", cpu.Registers.getDE())
	}
}

// DEC BC
func TestDEC_BC(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b10
	cpu.Memory[cpu.Registers.PC] = 0b1011
	cpu.execOpcodes()
	//// fmt.Printf("0b%b", cpu.Registers.getBC())

	if cpu.Registers.getBC() != 0b100000001 {
		t.Error("Expected 0b100000001, got ", cpu.Registers.getBC())
	}
}

// DEC [HL]
func TestDEC_HLmem(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.H = 0b1
	cpu.Registers.L = 0b10
	cpu.Registers.setHL(0b100)
	hladdr := cpu.Registers.getHL()
	cpu.Memory[hladdr] = 0b100

	//cpu.execDEC(operands, flags)
	cpu.Memory[cpu.Registers.PC] = 0b110101
	cpu.execOpcodes()
	if cpu.Memory[cpu.Registers.getHL()] != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Memory[cpu.Registers.getHL()])
	}
}

// DEC SP
func TestDEC_SP(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.SP = 0b100

	cpu.Memory[cpu.Registers.PC] = 0b111011
	cpu.execOpcodes()
	if cpu.Registers.SP != 0b11 {
		t.Error("Expected 0b100, got ", cpu.Registers.SP)
	}
}

// SUB A, B
//func TestSUB_A_B(t *testing.T) {
//	cpu := CPU{
//		Registers: Registers{
//			A: 0b1000,
//			B: 0b110,
//		},
//		Memory: [65536]uint8(make([]uint8, 65536)),
//	}
//
//	//cpu.execSUB(operands, flags)
//	cpu.Memory[cpu.Registers.PC] = 0b10010000
//	cpu.execOpcodes()
//	// // fmt.Println("SUB A, B TEST")
//	// // fmt.Println(cpu.Registers.A)
//	if cpu.Registers.A != 0b10 {
//		t.Error("Expected 0b10, got ", cpu.Registers.A)
//	}
//
//}

// AND
func TestAND(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 10
	cpu.Memory[cpu.Registers.PC] = 0b10100000 // AND A, B
	cpu.execOpcodes()

	if cpu.Registers.A != 0b1 {
		t.Error("Expected 0b1, got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10100001 // AND A, C
	cpu.execOpcodes()
	res := 0b1 & 0b11
	if cpu.Registers.A != uint8(res) {
		t.Error("Expected 0b1, got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10100010 // AND A, D
	cpu.execOpcodes()
	res = res & 0b100
	if cpu.Registers.A != uint8(res) {
		t.Error("Expected 0b1, got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10100011 // AND A, E
	cpu.execOpcodes()
	res = res & 0b101
	if cpu.Registers.A != uint8(res) {
		t.Error("Expected 0b1, got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10100100 // AND A, H
	cpu.execOpcodes()
	res = res & 0b110
	if cpu.Registers.A != uint8(res) {
		t.Error("Expected 0b1, got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10100101 // AND A, L
	cpu.execOpcodes()
	res = res & 0b111
	if cpu.Registers.A != uint8(res) {
		t.Error("Expected 0b1, got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10100110 // AND A, [HL]
	cpu.execOpcodes()

	res = int(uint8(res) & cpu.Memory[cpu.Registers.getHL()])
	if cpu.Registers.A != uint8(res) {
		t.Error("Expected 0b1, got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10100111 // AND A, A
	cpu.execOpcodes()
	res = res & res
	if cpu.Registers.A != uint8(res) {
		t.Error("Expected 0b1, got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11100110 // AND A, imm8
	cpu.Memory[cpu.Registers.PC+1] = 0b1

	cpu.execOpcodes()
	res = res & 0b1
	if cpu.Registers.A != uint8(res) {
		t.Error("Expected 0b1, got ", cpu.Registers.A)
	}

}

// CP A, n
func TestCP(t *testing.T) {

	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Memory[cpu.Registers.PC] = 0b10111000 // CP A, B
	cpu.execOpcodes()
	if cpu.Registers.A != 0b11 && cpu.Registers.B != 0b1 {
		t.Error("Expected 0b11 and 0b1, got ", cpu.Registers.A, cpu.Registers.B)
	}
	if cpu.Registers.getFlag(flagN) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagN))
	}
	if cpu.Registers.getFlag(flagZ) != false {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagZ))
	}

	cpu.Memory[cpu.Registers.PC] = 0b10111001 // CP A, C
	cpu.execOpcodes()
	if cpu.Registers.A != 0b11 && cpu.Registers.C != 0b11 {
		t.Error("Expected 0b11 and 0b11, got ", cpu.Registers.A, cpu.Registers.C)
	}
	if cpu.Registers.getFlag(flagN) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagN))
	}
	if cpu.Registers.getFlag(flagZ) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagZ))
	}

	cpu.Memory[cpu.Registers.PC] = 0b10111010 // CP A, D
	cpu.execOpcodes()
	if cpu.Registers.A != 0b11 && cpu.Registers.D != 0b100 {
		t.Error("Expected 0b11 and 0b100, got ", cpu.Registers.A, cpu.Registers.D)
	}
	if cpu.Registers.getFlag(flagN) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagN))
	}
	if cpu.Registers.getFlag(flagZ) != false {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagZ))
	}

	cpu.Memory[cpu.Registers.PC] = 0b10111011 // CP A, E
	cpu.execOpcodes()
	if cpu.Registers.A != 0b11 && cpu.Registers.E != 0b101 {
		t.Error("Expected 0b11 and 0b101, got ", cpu.Registers.A, cpu.Registers.E)
	}
	if cpu.Registers.getFlag(flagN) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagN))
	}
	if cpu.Registers.getFlag(flagZ) != false {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagZ))
	}

	cpu.Memory[cpu.Registers.PC] = 0b10111100 // CP A, H
	cpu.execOpcodes()
	if cpu.Registers.A != 0b11 && cpu.Registers.H != 0b110 {
		t.Error("Expected 0b11 and 0b110, got ", cpu.Registers.A, cpu.Registers.H)
	}
	if cpu.Registers.getFlag(flagN) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagN))
	}
	if cpu.Registers.getFlag(flagZ) != false {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagZ))
	}

	cpu.Memory[cpu.Registers.PC] = 0b10111101 // CP A, L
	cpu.execOpcodes()
	if cpu.Registers.A != 0b11 && cpu.Registers.L != 0b111 {
		t.Error("Expected 0b11 and 0b111, got ", cpu.Registers.A, cpu.Registers.L)
	}
	if cpu.Registers.getFlag(flagN) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagN))
	}
	if cpu.Registers.getFlag(flagZ) != false {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagZ))
	}

	cpu.Memory[cpu.Registers.PC] = 0b10111110 // CP A, [HL]
	//fmt.Printf("hl 0x%04X", cpu.Registers.getHL())
	cpu.Memory[cpu.Registers.getHL()] = 0b11
	cpu.execOpcodes()

	if cpu.Registers.A != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.A)
	}
	if cpu.Registers.getFlag(flagN) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagN))
	}
	if cpu.Registers.getFlag(flagZ) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagZ))
	}

	cpu.Memory[cpu.Registers.PC] = 0b10111111 // CP A, A
	//fmt.Printf("hl 0x%04X", cpu.Registers.getHL())
	cpu.execOpcodes()

	if cpu.Registers.A != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.A)
	}
	if cpu.Registers.getFlag(flagN) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagN))
	}
	if cpu.Registers.getFlag(flagZ) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagZ))
	}

	cpu.Memory[cpu.Registers.PC] = 0b11111110 // CP A, imm8
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	//fmt.Printf("hl 0x%04X", cpu.Registers.getHL())
	cpu.execOpcodes()

	if cpu.Registers.A != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Registers.A)
	}
	if cpu.Registers.getFlag(flagN) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagN))
	}
	if cpu.Registers.getFlag(flagZ) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagZ))
	}
}

// PUSH
func TestPush(t *testing.T) {

	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 15
	cpu.Memory[cpu.Registers.PC] = 0b11000101 // PUSH BC
	cpu.execOpcodes()

	if cpu.Registers.SP != 13 {
		t.Error("Expected 13, got ", cpu.Registers.SP)
	}
	//fmt.Println((cpu.Registers.getBC() & 0xFF00) >> 8)
	//fmt.Println(cpu.Registers.getBC() & 0xFF)
	if cpu.Memory[cpu.Registers.SP+1] != 1 {
		t.Error("Expected 1, got ", cpu.Registers.SP+1)
	}
	if cpu.Memory[cpu.Registers.SP] != 3 {
		t.Error("Expected 3, got ", cpu.Registers.SP)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11010101 // PUSH DE
	cpu.execOpcodes()

	if cpu.Registers.SP != 11 {
		t.Error("Expected 11, got ", cpu.Registers.SP)
	}
	//fmt.Println((cpu.Registers.getDE() & 0xFF00) >> 8)
	//fmt.Println(cpu.Registers.getDE() & 0xFF)
	if cpu.Memory[cpu.Registers.SP+1] != 4 {
		t.Error("Expected 4, got ", cpu.Registers.SP+1)
	}
	if cpu.Memory[cpu.Registers.SP] != 5 {
		t.Error("Expected 5, got ", cpu.Registers.SP)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11100101 // PUSH HL
	cpu.execOpcodes()

	if cpu.Registers.SP != 9 {
		t.Error("Expected 11, got ", cpu.Registers.SP)
	}
	//fmt.Println((cpu.Registers.getHL() & 0xFF00) >> 8)
	//fmt.Println(cpu.Registers.getHL() & 0xFF)
	if cpu.Memory[cpu.Registers.SP+1] != 6 {
		t.Error("Expected 6, got ", cpu.Registers.SP+1)
	}
	if cpu.Memory[cpu.Registers.SP] != 7 {
		t.Error("Expected 7, got ", cpu.Registers.SP)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11110101 // PUSH AF
	cpu.execOpcodes()

	if cpu.Registers.SP != 7 {
		t.Error("Expected 11, got ", cpu.Registers.SP)
	}
	//fmt.Println("af", (cpu.Registers.getAF()&0xFF00)>>8)
	//fmt.Println(cpu.Registers.getAF() & 0xFF)
	if cpu.Memory[cpu.Registers.SP+1] != 3 {
		t.Error("Expected 3, got ", cpu.Registers.SP+1)
	}
	if cpu.Memory[cpu.Registers.SP] != 2 {
		t.Error("Expected 2, got ", cpu.Registers.SP)
	}
}

// POP
func TestPop(t *testing.T) {

	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 10
	cpu.Memory[cpu.Registers.PC] = 0b11010001 // POP DE
	cpu.Memory[cpu.Registers.SP] = 0b101
	cpu.Memory[cpu.Registers.SP+1] = 0b10
	cpu.execOpcodes()

	if cpu.Registers.SP != 12 {
		t.Error("Expected 12, got ", cpu.Registers.SP)
	}
	res := uint16(0b10<<8 | 0b101)
	//fmt.Printf("%08b", res)

	if cpu.Registers.getDE() != res {
		t.Error("Expected", res, "got ", cpu.Registers.getDE())
	}

	cpu.Memory[cpu.Registers.PC] = 0b11100001 // POP HL
	cpu.Memory[cpu.Registers.SP] = 0b10
	cpu.Memory[cpu.Registers.SP+1] = 0b1011
	cpu.execOpcodes()

	if cpu.Registers.SP != 14 {
		t.Error("Expected 14, got ", cpu.Registers.SP)
	}
	res = uint16(0b1011<<8 | 0b10)
	//fmt.Printf("%08b", res)

	if cpu.Registers.getHL() != res {
		t.Error("Expected", res, "got ", cpu.Registers.getHL())
	}

	cpu.Memory[cpu.Registers.PC] = 0b11110001 // POP BC
	cpu.Memory[cpu.Registers.SP] = 0b101
	cpu.Memory[cpu.Registers.SP+1] = 0b11
	cpu.execOpcodes()

	if cpu.Registers.SP != 16 {
		t.Error("Expected 16, got ", cpu.Registers.SP)
	}
	res = uint16(0b11<<8 | 0b101)
	//fmt.Printf("%08b", res)

	if cpu.Registers.getAF() != res {
		t.Error("Expected", res, "got ", cpu.Registers.getAF())
	}

	cpu.Memory[cpu.Registers.PC] = 0b11000001 // POP BC
	cpu.Memory[cpu.Registers.SP] = 0b111
	cpu.Memory[cpu.Registers.SP+1] = 0b1101
	cpu.execOpcodes()

	if cpu.Registers.SP != 18 {
		t.Error("Expected 18, got ", cpu.Registers.SP)
	}
	res = uint16(0b1101<<8 | 0b111)
	//fmt.Printf("%08b", res)

	if cpu.Registers.getBC() != res {
		t.Error("Expected", res, "got ", cpu.Registers.getBC())
	}
}

// RST
func TestRST(t *testing.T) {

	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0x0021
	cpu.Memory[cpu.Registers.PC] = 0b11000111 //  RST $00
	hi := (0x0022 & 0xFF00) >> 8
	lo := 0x0022 & 0xFF
	cpu.execOpcodes()

	if cpu.Registers.SP != 14 {
		t.Error("Expected 14, got ", cpu.Registers.SP)
	}
	//fmt.Printf("PC %08b", cpu.Registers.PC)
	//fmt.Printf("low is %02X", uint8(lo))
	//fmt.Printf("high is %02X", uint8(hi))
	//fmt.Printf("SP %02X SP+1 %02X", cpu.Memory[cpu.Registers.SP], cpu.Memory[cpu.Registers.SP+1])
	if cpu.Memory[cpu.Registers.SP+1] != uint8(hi) {
		t.Error("Expected ", uint8(hi), " got ", cpu.Memory[cpu.Registers.SP+1])
	}
	if cpu.Memory[cpu.Registers.SP] != uint8(lo) {
		t.Error("Expected ", uint8(lo), " got ", cpu.Memory[cpu.Registers.SP])
	}

	if cpu.Registers.PC != 0x00 {
		t.Error("Expected 0x00, got ", cpu.Registers.PC)
	}

	cpu.Registers.PC = 0x22
	cpu.Memory[cpu.Registers.PC] = 0b11001111 //  RST $08
	hi = (0x0023 & 0xFF00) >> 8
	lo = 0x0023 & 0xFF

	cpu.execOpcodes()

	if cpu.Registers.SP != 12 {
		t.Error("Expected 12, got ", cpu.Registers.SP)
	}
	if cpu.Memory[cpu.Registers.SP+1] != uint8(hi) {
		t.Error("Expected ", uint8(hi), " got ", cpu.Memory[cpu.Registers.SP+1])
	}
	if cpu.Memory[cpu.Registers.SP] != uint8(lo) {
		t.Error("Expected ", uint8(lo), " got ", cpu.Memory[cpu.Registers.SP])
	}

	if cpu.Registers.PC != 0x08 {
		t.Error("Expected 0x08, got ", cpu.Registers.PC)
	}

	cpu.Registers.PC = 0x23
	cpu.Memory[cpu.Registers.PC] = 0b11010111 //  RST $10
	hi = (0x0024 & 0xFF00) >> 8
	lo = 0x0024 & 0xFF

	cpu.execOpcodes()

	if cpu.Registers.SP != 10 {
		t.Error("Expected 04, got ", cpu.Registers.SP)
	}
	if cpu.Memory[cpu.Registers.SP+1] != uint8(hi) {
		t.Error("Expected ", uint8(hi), " got ", cpu.Memory[cpu.Registers.SP+1])
	}
	if cpu.Memory[cpu.Registers.SP] != uint8(lo) {
		t.Error("Expected ", uint8(lo), " got ", cpu.Memory[cpu.Registers.SP])
	}

	if cpu.Registers.PC != 0x10 {
		t.Error("Expected 0x10, got ", cpu.Registers.PC)
	}

	cpu.Registers.PC = 0x24
	cpu.Memory[cpu.Registers.PC] = 0b11011111 //  RST $18
	hi = (0x0025 & 0xFF00) >> 8
	lo = 0x0025 & 0xFF

	cpu.execOpcodes()

	if cpu.Registers.SP != 8 {
		t.Error("Expected 8, got ", cpu.Registers.SP)
	}
	if cpu.Memory[cpu.Registers.SP+1] != uint8(hi) {
		t.Error("Expected ", uint8(hi), " got ", cpu.Memory[cpu.Registers.SP+1])
	}
	if cpu.Memory[cpu.Registers.SP] != uint8(lo) {
		t.Error("Expected ", uint8(lo), " got ", cpu.Memory[cpu.Registers.SP])
	}

	if cpu.Registers.PC != 0x18 {
		t.Error("Expected 0x18, got ", cpu.Registers.PC)
	}

	cpu.Registers.PC = 0x25
	cpu.Memory[cpu.Registers.PC] = 0b11100111 //  RST $20
	hi = (0x0026 & 0xFF00) >> 8
	lo = 0x0026 & 0xFF

	cpu.execOpcodes()

	if cpu.Registers.SP != 6 {
		t.Error("Expected 6, got ", cpu.Registers.SP)
	}
	if cpu.Memory[cpu.Registers.SP+1] != uint8(hi) {
		t.Error("Expected ", uint8(hi), " got ", cpu.Memory[cpu.Registers.SP+1])
	}
	if cpu.Memory[cpu.Registers.SP] != uint8(lo) {
		t.Error("Expected ", uint8(lo), " got ", cpu.Memory[cpu.Registers.SP])
	}

	if cpu.Registers.PC != 0x20 {
		t.Error("Expected 0x20, got ", cpu.Registers.PC)
	}

	cpu.Registers.PC = 0x26
	cpu.Memory[cpu.Registers.PC] = 0b11101111 //  RST $28
	hi = (0x0027 & 0xFF00) >> 8
	lo = 0x0027 & 0xFF

	cpu.execOpcodes()

	if cpu.Registers.SP != 4 {
		t.Error("Expected 4, got ", cpu.Registers.SP)
	}
	if cpu.Memory[cpu.Registers.SP+1] != uint8(hi) {
		t.Error("Expected ", uint8(hi), " got ", cpu.Memory[cpu.Registers.SP+1])
	}
	if cpu.Memory[cpu.Registers.SP] != uint8(lo) {
		t.Error("Expected ", uint8(lo), " got ", cpu.Memory[cpu.Registers.SP])
	}

	if cpu.Registers.PC != 0x28 {
		t.Error("Expected 0x28, got ", cpu.Registers.PC)
	}

	cpu.Registers.PC = 0x27
	cpu.Memory[cpu.Registers.PC] = 0b11110111 //  RST $30
	hi = (0x0028 & 0xFF00) >> 8
	lo = 0x0028 & 0xFF

	cpu.execOpcodes()

	if cpu.Registers.SP != 2 {
		t.Error("Expected 2, got ", cpu.Registers.SP)
	}
	if cpu.Memory[cpu.Registers.SP+1] != uint8(hi) {
		t.Error("Expected ", uint8(hi), " got ", cpu.Memory[cpu.Registers.SP+1])
	}
	if cpu.Memory[cpu.Registers.SP] != uint8(lo) {
		t.Error("Expected ", uint8(lo), " got ", cpu.Memory[cpu.Registers.SP])
	}

	if cpu.Registers.PC != 0x30 {
		t.Error("Expected 0x30, got ", cpu.Registers.PC)
	}

	cpu.Registers.PC = 0x28
	cpu.Memory[cpu.Registers.PC] = 0b11111111 //  RST $38
	hi = (0x0029 & 0xFF00) >> 8
	lo = 0x0029 & 0xFF

	cpu.execOpcodes()

	if cpu.Registers.SP != 0 {
		t.Error("Expected 14, got ", cpu.Registers.SP)
	}
	if cpu.Memory[cpu.Registers.SP+1] != uint8(hi) {
		t.Error("Expected ", uint8(hi), " got ", cpu.Memory[cpu.Registers.SP+1])
	}
	if cpu.Memory[cpu.Registers.SP] != uint8(lo) {
		t.Error("Expected ", uint8(lo), " got ", cpu.Memory[cpu.Registers.SP])
	}

	if cpu.Registers.PC != 0x38 {
		t.Error("Expected 0x38, got ", cpu.Registers.PC)
	}
}

// BIT
func TestBIT(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b1000000 //  BIT 0, B

	cpu.execOpcodes()
	res := 0b1 & (1 << 0)
	fmt.Println(res)
	if cpu.Registers.getFlag(flagN) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagN))
	}
	if cpu.Registers.getFlag(flagH) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagH))
	}
	if cpu.Registers.getFlag(flagZ) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagZ))
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b1000001 //  BIT 0, C

	cpu.execOpcodes()
	res = 0b11 & (1 << 0)
	fmt.Println(res)
	if cpu.Registers.getFlag(flagN) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagN))
	}
	if cpu.Registers.getFlag(flagH) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagH))
	}
	if cpu.Registers.getFlag(flagZ) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagZ))
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b1000010 //  BIT 0, D

	cpu.execOpcodes()
	res = 0b100 & (1 << 0)
	fmt.Println(res)
	if cpu.Registers.getFlag(flagN) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagN))
	}
	if cpu.Registers.getFlag(flagH) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagH))
	}
	if cpu.Registers.getFlag(flagZ) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagZ))
	}

	// BIT 0 E-A, BIT 1 A-L

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b1010010 //  BIT 2, D

	cpu.execOpcodes()
	res = 0b100 & (1 << 2)
	fmt.Println(res)
	if cpu.Registers.getFlag(flagN) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagN))
	}
	if cpu.Registers.getFlag(flagH) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagH))
	}
	if cpu.Registers.getFlag(flagZ) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagZ))
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b1010101 //  BIT 2, L

	cpu.execOpcodes()
	res = 0b111 & (1 << 2)
	fmt.Println(res)
	if cpu.Registers.getFlag(flagN) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagN))
	}
	if cpu.Registers.getFlag(flagH) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagH))
	}
	if cpu.Registers.getFlag(flagZ) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagZ))
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b1011111 //  BIT 3, A

	cpu.execOpcodes()
	res = 0b11 & (1 << 3)
	fmt.Println(res)
	if cpu.Registers.getFlag(flagN) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagN))
	}
	if cpu.Registers.getFlag(flagH) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagH))
	}
	if cpu.Registers.getFlag(flagZ) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagZ))
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b1100011 //  BIT 4, E

	cpu.execOpcodes()
	res = 0b101 & (1 << 4)
	fmt.Println(res)
	if cpu.Registers.getFlag(flagN) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagN))
	}
	if cpu.Registers.getFlag(flagH) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagH))
	}
	if cpu.Registers.getFlag(flagZ) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagZ))
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b1110110 //  BIT 6, [HL]

	cpu.execOpcodes()
	res = int(cpu.Registers.getHL() & (1 << 6))
	fmt.Println(res)
	if cpu.Registers.getFlag(flagN) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagN))
	}
	if cpu.Registers.getFlag(flagH) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagH))
	}
	if cpu.Registers.getFlag(flagZ) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagZ))
	}

}

// SET
func TestSET(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b11000111 //  SET 0, A

	cpu.execOpcodes()
	res := uint8(0b11 | (1 << 0))
	fmt.Println(res)
	if cpu.Registers.A != res {
		t.Error("Expected ", res, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b11001011 // SET 1, E

	cpu.execOpcodes()
	res = uint8(0b101 | (1 << 1))
	fmt.Println(res)
	if cpu.Registers.E != res {
		t.Error("Expected ", res, " got ", cpu.Registers.E)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b11010101 // SET 2, L

	cpu.execOpcodes()
	res = uint8(0b111 | (1 << 2))
	fmt.Println(res)
	if cpu.Registers.L != res {
		t.Error("Expected ", res, " got ", cpu.Registers.L)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b11011010 // SET 3, D

	cpu.execOpcodes()
	res = uint8(0b100 | (1 << 3))
	fmt.Println(res)
	if cpu.Registers.D != res {
		t.Error("Expected ", res, " got ", cpu.Registers.D)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b11100001 // SET 4, C

	cpu.execOpcodes()
	res = uint8(0b11 | (1 << 4))
	fmt.Println(res)
	if cpu.Registers.C != res {
		t.Error("Expected ", res, " got ", cpu.Registers.C)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b11101100 // SET 5, H

	cpu.execOpcodes()
	res = uint8(0b110 | (1 << 5))
	fmt.Println(res)
	if cpu.Registers.H != res {
		t.Error("Expected ", res, " got ", cpu.Registers.H)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b11110110 // SET 6, [HL]
	cpu.Memory[cpu.Registers.getHL()] = 0b1
	cpu.execOpcodes()
	res = uint8(cpu.Memory[cpu.Registers.getHL()] | (1 << 6))
	fmt.Println(res)
	if cpu.Memory[cpu.Registers.getHL()] != res {
		t.Error("Expected ", res, " got ", cpu.Memory[cpu.Registers.getHL()])
	}

}

// RES
func TestRES(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b10000001 //  RES 0, C
	cpu.execOpcodes()
	res := 0b11 & ^(1 << 0)
	if cpu.Registers.C != uint8(res) {
		t.Error("Expected ", res, " got ", cpu.Registers.C)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b10001101 //  RES 1, L
	cpu.execOpcodes()
	res = 0b111 & ^(1 << 1)
	if cpu.Registers.L != uint8(res) {
		t.Error("Expected ", res, " got ", cpu.Registers.L)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b10010010 //  RES 2, D
	cpu.execOpcodes()
	res = 0b100 & ^(1 << 2)
	if cpu.Registers.D != uint8(res) {
		t.Error("Expected ", res, " got ", cpu.Registers.D)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b10100011 //  RES 4, E
	fmt.Printf("reg e %02X", cpu.Registers.E)
	cpu.execOpcodes()
	res = 0b101 & ^(1 << 4)
	fmt.Printf("reg e %02X", cpu.Registers.E)
	if cpu.Registers.E != uint8(res) {
		t.Error("Expected ", res, " got ", cpu.Registers.E)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b10101110 //  RES 5, [HL]
	cpu.Memory[cpu.Registers.getHL()] = 0b1
	cpu.execOpcodes()

	res = 0b1 & ^(1 << 5)
	fmt.Println("Aaaaaaaaaaaa")
	if cpu.Memory[cpu.Registers.getHL()] != uint8(res) {
		t.Error("Expected ", res, " got ", cpu.Memory[cpu.Registers.getHL()])
	}

}

// SWAP
func TestSWAP(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b110000 //  SWAP B
	cpu.execOpcodes()
	if cpu.Registers.B != 0b10000 {
		t.Error("Expected ", 0b10000, " got ", cpu.Registers.B)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b110001 //  SWAP C
	cpu.execOpcodes()
	if cpu.Registers.C != 0b110000 {
		t.Error("Expected ", 0b110000, " got ", cpu.Registers.C)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b110010 //  SWAP D
	cpu.execOpcodes()
	if cpu.Registers.D != 0b01000000 {
		t.Error("Expected ", 0b01000000, " got ", cpu.Registers.D)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b110011 //  SWAP E
	cpu.execOpcodes()
	if cpu.Registers.E != 0b01010000 {
		t.Error("Expected ", 0b01010000, " got ", cpu.Registers.E)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b110100 //  SWAP H
	cpu.execOpcodes()
	if cpu.Registers.H != 0b01100000 {
		t.Error("Expected ", 0b01100000, " got ", cpu.Registers.H)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b110101 //  SWAP L
	cpu.execOpcodes()
	if cpu.Registers.L != 0b01110000 {
		t.Error("Expected ", 0b01110000, " got ", cpu.Registers.L)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b110111 //  SWAP A
	cpu.execOpcodes()
	if cpu.Registers.A != 0b00110000 {
		t.Error("Expected ", 0b00110000, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.getHL()] = 0b11
	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b110110 //  SWAP [HL]
	cpu.execOpcodes()
	if cpu.Memory[cpu.Registers.getHL()] != 0b00110000 {
		t.Error("Expected ", 0b00110000, " got ", cpu.Registers.A)
	}

}

// SLA
func TestSLA(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b100000 //  SLA B
	//fmt.Printf("reg %04X", cpu.Registers.B)
	cpu.execOpcodes()
	//fmt.Printf("reg after %04X", cpu.Registers.B)

	reg := 0b1 << 1
	//fmt.Println("reg ", reg)
	if cpu.Registers.B != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.B)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b100001 //  SLA C
	cpu.execOpcodes()

	reg = 0b11 << 1
	if cpu.Registers.C != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.C)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b100010 //  SLA D
	cpu.execOpcodes()

	reg = 0b100 << 1
	if cpu.Registers.D != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.D)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b100011 //  SLA E
	cpu.execOpcodes()

	reg = 0b101 << 1
	if cpu.Registers.E != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.E)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b100100 //  SLA H
	cpu.execOpcodes()

	reg = 0b110 << 1
	if cpu.Registers.H != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.H)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b100101 //  SLA L
	cpu.execOpcodes()

	reg = 0b111 << 1
	if cpu.Registers.L != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.L)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b100111 //  SLA A
	cpu.execOpcodes()

	reg = 0b11 << 1
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b100110 //  SLA [HL]
	cpu.Memory[cpu.Registers.getHL()] = 0b11

	cpu.execOpcodes()

	reg = 0b11 << 1
	if cpu.Memory[cpu.Registers.getHL()] != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Memory[cpu.Registers.getHL()])
	}
}

// SRA
func TestSRA(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b101000 //  SRA B
	cpu.execOpcodes()
	reg := (0b1 >> 1) | (0b1 & 0x80)
	if cpu.Registers.B != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.B)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b101001 //  SRA C
	cpu.execOpcodes()
	reg = (0b11 >> 1) | (0b11 & 0x80)
	if cpu.Registers.C != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.C)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b101010 //  SRA D
	cpu.execOpcodes()
	reg = (0b100 >> 1) | (0b100 & 0x80)
	if cpu.Registers.D != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.D)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b101011 //  SRA E
	cpu.execOpcodes()
	reg = (0b101 >> 1) | (0b101 & 0x80)
	if cpu.Registers.E != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.E)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b101100 //  SRA H
	cpu.execOpcodes()
	reg = (0b110 >> 1) | (0b110 & 0x80)
	if cpu.Registers.H != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.H)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b101101 //  SRA L
	cpu.execOpcodes()
	reg = (0b111 >> 1) | (0b111 & 0x80)
	if cpu.Registers.L != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.L)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b101111 //  SRA A
	cpu.execOpcodes()
	reg = (0b11 >> 1) | (0b11 & 0x80)
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b101110 //  SRA [HL]
	cpu.Memory[cpu.Registers.getHL()] = 0b11
	cpu.execOpcodes()
	reg = (0b11 >> 1) | (0b11 & 0x80)
	if cpu.Memory[cpu.Registers.getHL()] != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Memory[cpu.Registers.getHL()])
	}
}

// SRL
func TestSRL(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b111000 //  SRL B
	cpu.execOpcodes()
	reg := 0b1 >> 1
	if cpu.Registers.B != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.B)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b111001 //  SRL C
	cpu.execOpcodes()
	reg = 0b11 >> 1
	if cpu.Registers.C != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.C)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b111010 //  SRL D
	cpu.execOpcodes()
	reg = 0b100 >> 1
	if cpu.Registers.D != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.D)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b111011 //  SRL E
	cpu.execOpcodes()
	reg = 0b101 >> 1
	if cpu.Registers.E != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.E)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b111100 //  SRL H
	cpu.execOpcodes()
	reg = 0b110 >> 1
	if cpu.Registers.H != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.H)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b111101 //  SRL L
	cpu.execOpcodes()
	reg = 0b111 >> 1
	if cpu.Registers.L != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.L)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b111111 //  SRL A
	cpu.execOpcodes()
	reg = 0b11 >> 1
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b111110 //  SRL [HL]
	cpu.Memory[cpu.Registers.getHL()] = 0b101
	cpu.execOpcodes()
	reg = 0b101 >> 1
	if cpu.Memory[cpu.Registers.getHL()] != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Memory[cpu.Registers.getHL()])
	}
}

// RLC
func TestRLC(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b0 //  RLC B
	cpu.execOpcodes()
	reg := 0b1 << 1
	if cpu.Registers.B != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.B)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b1 //  RLC C
	cpu.execOpcodes()
	reg = 0b11 << 1
	if cpu.Registers.C != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.C)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b10 //  RLC D
	cpu.execOpcodes()
	reg = 0b100 << 1
	if cpu.Registers.D != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.D)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b11 //  RLC E
	cpu.execOpcodes()
	reg = 0b101 << 1
	if cpu.Registers.E != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.E)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b100 //  RLC H
	cpu.execOpcodes()
	reg = 0b110 << 1
	if cpu.Registers.H != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.H)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b101 //  RLC L
	cpu.execOpcodes()
	reg = 0b111 << 1
	if cpu.Registers.L != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.L)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b111 //  RLC A
	cpu.execOpcodes()
	reg = 0b11 << 1
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b110 //  RLC [HL]
	cpu.Memory[cpu.Registers.getHL()] = 0b101
	cpu.execOpcodes()
	reg = 0b101 << 1
	if cpu.Memory[cpu.Registers.getHL()] != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Memory[cpu.Registers.getHL()])
	}
}

// RL
func TestRL(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b10000 //  RL B
	cpu.execOpcodes()
	reg := (0b1 << 1) | uint8(0)
	if cpu.Registers.B != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.B)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b10001 //  RL C
	cpu.execOpcodes()
	reg = (0b11 << 1) | uint8(0)
	if cpu.Registers.C != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.C)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b10010 //  RL D
	cpu.execOpcodes()
	reg = (0b100 << 1) | uint8(0)
	if cpu.Registers.D != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.D)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b10011 //  RL E
	cpu.execOpcodes()
	reg = (0b101 << 1) | uint8(0)
	if cpu.Registers.E != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.E)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b10100 //  RL H
	cpu.execOpcodes()
	reg = (0b110 << 1) | uint8(0)
	if cpu.Registers.H != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.H)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b10101 //  RL L
	cpu.execOpcodes()
	reg = (0b111 << 1) | uint8(0)
	if cpu.Registers.L != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.L)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b10111 //  RL A
	cpu.Registers.setFlag(flagC, true)
	cpu.execOpcodes()
	reg = (0b11 << 1) | uint8(1)
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b10110 //  RL [HL]
	cpu.Memory[cpu.Registers.getHL()] = 0b101
	cpu.execOpcodes()
	reg = (0b101 << 1) | uint8(0)
	if cpu.Memory[cpu.Registers.getHL()] != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Memory[cpu.Registers.getHL()])
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b10110 //  RL [HL] - flagC set
	cpu.Memory[cpu.Registers.getHL()] = 0b101
	cpu.Registers.setFlag(flagC, true)
	cpu.execOpcodes()
	reg = (0b101 << 1) | uint8(1)
	if cpu.Memory[cpu.Registers.getHL()] != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Memory[cpu.Registers.getHL()])
	}

}

// RRC
func TestRRC(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b1000 //  RRC B
	cpu.execOpcodes()
	reg := 0b1 >> 1
	if cpu.Registers.B != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.B)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b1001 //  RRC C
	cpu.execOpcodes()
	reg = 0b11 >> 1
	if cpu.Registers.C != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.C)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b1010 //  RRC D
	cpu.execOpcodes()
	reg = 0b100 >> 1
	if cpu.Registers.D != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.D)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b1011 //  RRC E
	cpu.execOpcodes()
	reg = 0b101 >> 1
	if cpu.Registers.E != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.E)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b1100 //  RRC H
	cpu.execOpcodes()
	reg = 0b110 >> 1
	if cpu.Registers.H != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.H)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b1111 //  RRC A
	cpu.execOpcodes()
	reg = 0b11 >> 1
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b1110 //  RRC [HL]
	cpu.Memory[cpu.Registers.getHL()] = 0b1111
	cpu.execOpcodes()
	reg = 0b1111 >> 1
	if cpu.Memory[cpu.Registers.getHL()] != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Memory[cpu.Registers.getHL()])
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b1101 //  RRC L
	cpu.execOpcodes()
	reg = 0b111 >> 1
	if cpu.Registers.L != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}
}

// RR
func TestRR(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b11000 //  RR B
	cpu.execOpcodes()
	reg := (0b1 >> 1) | 0<<7
	if cpu.Registers.B != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.B)
	}

	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b11001 //  RR C
	cpu.execOpcodes()
	//fmt.Println("FlagC", cpu.Registers.getFlag(flagC))
	reg = (0b11 >> 1) | 1<<7
	if cpu.Registers.C != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.C)
	}

	cpu.Registers.setFlag(flagC, true)
	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b11010 //  RR D
	cpu.execOpcodes()
	reg = (0b100 >> 1) | 1<<7
	if cpu.Registers.D != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.D)
	}

	cpu.Registers.setFlag(flagC, true)
	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b11011 //  RR E
	cpu.execOpcodes()
	reg = (0b101 >> 1) | 1<<7
	if cpu.Registers.E != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.E)
	}

	cpu.Registers.setFlag(flagC, true)
	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b11100 //  RR H
	cpu.execOpcodes()
	reg = (0b110 >> 1) | 1<<7
	if cpu.Registers.H != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.H)
	}

	cpu.Registers.setFlag(flagC, true)
	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b11101 //  RR L
	cpu.execOpcodes()
	reg = (0b111 >> 1) | 1<<7
	if cpu.Registers.L != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.L)
	}

	cpu.Registers.setFlag(flagC, true)
	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b11111 //  RR A
	cpu.execOpcodes()
	reg = (0b11 >> 1) | 1<<7
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Registers.setFlag(flagC, true)
	cpu.Memory[cpu.Registers.PC] = 0xCB
	cpu.Memory[cpu.Registers.PC+1] = 0b11110 //  RR [HL]
	cpu.Memory[cpu.Registers.getHL()] = 0b111
	cpu.execOpcodes()
	reg = (0b111 >> 1) | 1<<7
	if cpu.Memory[cpu.Registers.getHL()] != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Memory[cpu.Registers.getHL()])
	}
}

// LD nn, imm16
func TestLDnn(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16

	cpu.Memory[cpu.Registers.PC] = 0b1 //  LD BC, imm16
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	fmt.Println("bc", cpu.Registers.getBC())

	reg := 0b001100000001
	cpu.execOpcodes()
	//fmt.Println("reg", reg, "bc", cpu.Registers.getBC())
	if cpu.Registers.getBC() != uint16(reg) {
		t.Error("Expected ", uint16(reg), " got ", cpu.Registers.getBC())
	}

	cpu.Memory[cpu.Registers.PC] = 0b10001 //  LD DE, imm16
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	fmt.Println("bc", cpu.Registers.getBC())

	reg = 0b001100000001
	cpu.execOpcodes()
	//fmt.Println("reg", reg, "bc", cpu.Registers.getBC())
	if cpu.Registers.getDE() != uint16(reg) {
		t.Error("Expected ", uint16(reg), " got ", cpu.Registers.getDE())
	}
}

// LD [nn], r
func TestLD_mem_nn(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16

	cpu.Memory[cpu.Registers.PC] = 0b10 //  LD [BC], A
	cpu.execOpcodes()
	if cpu.Memory[cpu.Registers.getBC()] != 0b11 {
		t.Error("Expected ", 0b11, " got ", cpu.Memory[cpu.Registers.getBC()])
	}

	cpu.Memory[cpu.Registers.PC] = 0b1000 //  LD [imm16], SP
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.execOpcodes()
	if cpu.Memory[0b001100000001] != 16 {
		t.Error("Expected ", 16, " got ", cpu.Memory[0b001100000001])
	}

	cpu.Memory[cpu.Registers.PC] = 0b10010 //  LD [DE], A
	cpu.execOpcodes()
	if cpu.Memory[cpu.Registers.getDE()] != 0b11 {
		t.Error("Expected ", 0b11, " got ", cpu.Memory[cpu.Registers.getDE()])
	}

	cpu.Memory[cpu.Registers.PC] = 0b11101010 //  LD [imm16], A
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.execOpcodes()
	if cpu.Memory[0b001100000001] != 0b11 {
		t.Error("Expected ", 0b11, " got ", cpu.Memory[cpu.Registers.getDE()])
	}
}

// JR
func TestJR(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b100000 // JR NZ, imm8
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Registers.setFlag(flagZ, false)
	cpu.execOpcodes()
	if cpu.Registers.PC != uint16(int32(0b1)+int32(int8(0b11))) {
		t.Error("Expected ", uint16(int32(0b1)+int32(int8(0b11))), " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b101000 // JR Z, e8
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Registers.setFlag(flagZ, true)
	cpu.execOpcodes()
	if cpu.Registers.PC != 0b0111 {
		t.Error("Expected ", 0b0111, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b110000 // JR NC, e8
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Registers.setFlag(flagC, false)
	cpu.execOpcodes()
	if cpu.Registers.PC != 0b00001010 {
		t.Error("Expected ", 0b00001010, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b111000 // JR C, e8
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Registers.setFlag(flagC, true)
	cpu.execOpcodes()
	if cpu.Registers.PC != 0b00001101 {
		t.Error("Expected ", 0b00001101, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11000 // JR e8
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Registers.setFlag(flagC, true)
	cpu.execOpcodes()
	if cpu.Registers.PC != 0b00010000 {
		t.Error("Expected ", 0b00010000, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b100000 // JR NZ, imm8 - flagZ set
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Registers.setFlag(flagZ, true)
	cpu.execOpcodes()

	if cpu.Registers.PC != 0b00010000+2 {
		t.Error("Expected ", 0b00010000+2, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b101000 // JR Z, e8
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Registers.setFlag(flagZ, false)
	cpu.execOpcodes()
	if cpu.Registers.PC != 0b00010000+4 {
		t.Error("Expected ", 0b00010000+4, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b110000 // JR NC, e8
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Registers.setFlag(flagC, true)
	cpu.execOpcodes()
	if cpu.Registers.PC != 0b00010000+6 {
		t.Error("Expected ", 0b00010000+6, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b111000 // JR C, e8
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Registers.setFlag(flagC, false)
	cpu.execOpcodes()
	if cpu.Registers.PC != 0b00010000+8 {
		t.Error("Expected ", 0b00010000+8, " got ", cpu.Registers.PC)
	}

}

func TestHLA(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b100010 // LD [HL+], A
	hl := cpu.Registers.getHL()

	cpu.execOpcodes()
	fmt.Println("hl", cpu.Registers.getHL())
	fmt.Println("hl mem", cpu.Memory[cpu.Registers.getHL()])
	if cpu.Memory[hl] != 0b11 {
		t.Error("Expected ", 0b11, " got ", cpu.Memory[cpu.Registers.getHL()])
	}

	cpu.Memory[cpu.Registers.PC] = 0b110010 // LD [HL-], A
	hl = cpu.Registers.getHL()

	cpu.execOpcodes()
	fmt.Println("hl", cpu.Registers.getHL())
	fmt.Println("hl mem", cpu.Memory[cpu.Registers.getHL()])
	if cpu.Memory[hl] != 0b11 {
		t.Error("Expected ", 0b11, " got ", cpu.Memory[cpu.Registers.getHL()])
	}
}

// LD A, [HL+]
func TestA_HL(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b101010 // LD A, [HL+]
	cpu.Memory[cpu.Registers.getHL()] = 0b111
	cpu.execOpcodes()
	fmt.Println("hl", cpu.Registers.getHL())
	fmt.Println("hl mem", cpu.Memory[cpu.Registers.getHL()])
	if cpu.Registers.A != 0b111 {
		t.Error("Expected ", 0b111, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b111010 // LD A, [HL-]
	cpu.Memory[cpu.Registers.getHL()] = 0b111
	cpu.execOpcodes()
	fmt.Println("hl", cpu.Registers.getHL())
	fmt.Println("hl mem", cpu.Memory[cpu.Registers.getHL()])
	if cpu.Registers.A != 0b111 {
		t.Error("Expected ", 0b111, " got ", cpu.Registers.A)
	}
}

// LD SP, n16
func TestLD_SP_n16(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b110001 // LD SP, n16
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11

	reg := 0b001100000001
	cpu.execOpcodes()
	if cpu.Registers.SP != uint16(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.SP)
	}
}

// XOR
func TestXOR(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b10101000 // XOR A, B
	reg := 0b11 ^ 0b1
	cpu.execOpcodes()
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10101001 // XOR A, C
	reg = reg ^ 0b11
	cpu.execOpcodes()
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10101010 // XOR A, D
	reg = reg ^ 0b100
	cpu.execOpcodes()
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10101011 // XOR A, E
	reg = reg ^ 0b101
	cpu.execOpcodes()
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10101100 // XOR A, H
	reg = reg ^ 0b110
	cpu.execOpcodes()
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10101101 // XOR A, L
	reg = reg ^ 0b111
	cpu.execOpcodes()
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10101110 // XOR A, [HL]
	cpu.Memory[cpu.Registers.getHL()] = 0b1
	reg = reg ^ 0b1
	cpu.execOpcodes()
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10101111 // XOR A, A
	reg = reg ^ reg
	cpu.execOpcodes()
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11101110 // XOR A, imm8
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	reg = reg ^ 0b1
	cpu.execOpcodes()
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}
}

// OR
func TestOR(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b10110000 // OR A, B
	reg := 0b11 | 0b1
	cpu.execOpcodes()
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10110001 // OR A, C
	reg = reg | 0b11
	cpu.execOpcodes()
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10110010 // OR A, D
	reg = reg | 0b100
	cpu.execOpcodes()
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10110011 // OR A, E
	reg = reg | 0b101
	cpu.execOpcodes()
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10110100 // OR A, H
	reg = reg | 0b110
	cpu.execOpcodes()
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10110101 // OR A, L
	reg = reg | 0b111
	cpu.execOpcodes()
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10110110 // OR A, [HL]
	cpu.Memory[cpu.Registers.getHL()] = 0b1
	reg = reg | 0b1
	cpu.execOpcodes()
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b10110111 // OR A, A
	reg = reg | reg
	cpu.execOpcodes()
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11110110 // OR A, imm8
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	reg = reg | 0b1
	cpu.execOpcodes()
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}
}

// JP
func TestJP(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b11000010 // JP NZ, imm16
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11

	reg := 0b001100000001
	cpu.execOpcodes()
	if cpu.Registers.PC != uint16(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11001010 // JP Z, imm16
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Registers.setFlag(flagZ, true)
	reg = 0b001100000001
	cpu.execOpcodes()
	if cpu.Registers.PC != uint16(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11010010 // JP NC, imm16
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	reg = 0b001100000001
	cpu.execOpcodes()
	if cpu.Registers.PC != uint16(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11011010 // JP C, imm16
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	reg = 0b001100000001
	cpu.Registers.setFlag(flagC, true)

	cpu.execOpcodes()
	if cpu.Registers.PC != uint16(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11000011 // JP imm16
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	reg = 0b001100000001

	cpu.execOpcodes()
	if cpu.Registers.PC != uint16(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11101001 // JP HL
	cpu.execOpcodes()
	pc := cpu.Registers.PC
	if cpu.Registers.PC != cpu.Registers.getHL() {
		t.Error("Expected ", reg, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11000010 // JP NZ, imm16

	cpu.Registers.setFlag(flagZ, true)
	cpu.execOpcodes()
	if cpu.Registers.PC != pc+3 {
		t.Error("Expected ", pc+3, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11001010 // JP Z, imm16
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Registers.setFlag(flagZ, false)
	cpu.execOpcodes()
	if cpu.Registers.PC != pc+6 {
		t.Error("Expected ", pc+6, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11010010 // JP NC, imm16
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Registers.setFlag(flagC, true)
	cpu.execOpcodes()
	if cpu.Registers.PC != pc+9 {
		t.Error("Expected ", pc+9, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11011010 // JP C, imm16 - not set
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Registers.setFlag(flagC, false)

	cpu.execOpcodes()
	if cpu.Registers.PC != pc+12 {
		t.Error("Expected ", pc+12, " got ", cpu.Registers.PC)
	}

}

// CALL
func TestCALL_NZ(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b11000100 // CALL NZ, imm16
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11

	reg := 0b001100000001
	cpu.Registers.setFlag(flagZ, false)
	cpu.execOpcodes()

	if cpu.Memory[cpu.Registers.SP+1] != uint8((0b100&0xFF00)>>8) {
		t.Error("Expected ", uint8((0b100&0xFF00)>>8), " got ", cpu.Memory[cpu.Registers.SP+1])
	}
	if cpu.Memory[cpu.Registers.SP] != uint8(0b100&0xFF) {
		t.Error("Expected ", uint8(0b100&0xFF), " got ", cpu.Memory[cpu.Registers.SP])
	}

	if cpu.Registers.PC != uint16(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11000100 // CALL NZ, imm16 - set flag
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Registers.setFlag(flagZ, true)
	cpu.execOpcodes()

	if cpu.Memory[cpu.Registers.SP+1] != uint8((0b100&0xFF00)>>8) {
		t.Error("Expected ", uint8((0b100&0xFF00)>>8), " got ", cpu.Memory[cpu.Registers.SP+1])
	}
	if cpu.Memory[cpu.Registers.SP] != uint8(0b100&0xFF) {
		t.Error("Expected ", uint8(0b100&0xFF), " got ", cpu.Memory[cpu.Registers.SP])
	}

}

func TestCALL_Z(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b11001100 // CALL Z, imm16
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Registers.setFlag(flagZ, true)
	reg := 0b001100000001
	cpu.execOpcodes()

	if cpu.Memory[cpu.Registers.SP+1] != uint8((0b100&0xFF00)>>8) {
		t.Error("Expected ", uint8((0b100&0xFF00)>>8), " got ", cpu.Memory[cpu.Registers.SP+1])
	}
	if cpu.Memory[cpu.Registers.SP] != uint8(0b100&0xFF) {
		t.Error("Expected ", uint8(0b100&0xFF), " got ", cpu.Memory[cpu.Registers.SP])
	}

	if cpu.Registers.PC != uint16(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11001100 // CALL Z, imm16
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Registers.setFlag(flagZ, false)
	cpu.execOpcodes()

	if cpu.Memory[cpu.Registers.SP+1] != uint8((0b100&0xFF00)>>8) {
		t.Error("Expected ", uint8((0b100&0xFF00)>>8), " got ", cpu.Memory[cpu.Registers.SP+1])
	}
	if cpu.Memory[cpu.Registers.SP] != uint8(0b100&0xFF) {
		t.Error("Expected ", uint8(0b100&0xFF), " got ", cpu.Memory[cpu.Registers.SP])
	}

	if cpu.Registers.PC != uint16(reg)+3 {
		t.Error("Expected ", reg, " got ", cpu.Registers.PC)
	}
}

func TestCALL_NC(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b11010100 // CALL NC, imm16
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	reg := 0b001100000001
	cpu.execOpcodes()

	if cpu.Memory[cpu.Registers.SP+1] != uint8((0b100&0xFF00)>>8) {
		t.Error("Expected ", uint8((0b100&0xFF00)>>8), " got ", cpu.Memory[cpu.Registers.SP+1])
	}
	if cpu.Memory[cpu.Registers.SP] != uint8(0b100&0xFF) {
		t.Error("Expected ", uint8(0b100&0xFF), " got ", cpu.Memory[cpu.Registers.SP])
	}

	if cpu.Registers.PC != uint16(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11010100 // CALL NC, imm16 - flag set
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Registers.setFlag(flagC, true)
	cpu.execOpcodes()

	if cpu.Memory[cpu.Registers.SP+1] != uint8((0b100&0xFF00)>>8) {
		t.Error("Expected ", uint8((0b100&0xFF00)>>8), " got ", cpu.Memory[cpu.Registers.SP+1])
	}
	if cpu.Memory[cpu.Registers.SP] != uint8(0b100&0xFF) {
		t.Error("Expected ", uint8(0b100&0xFF), " got ", cpu.Memory[cpu.Registers.SP])
	}

	if cpu.Registers.PC != uint16(reg)+3 {
		t.Error("Expected ", reg+3, " got ", cpu.Registers.PC)
	}
}

func TestCALL_C(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b11011100 // CALL C, imm16
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Registers.setFlag(flagC, true)
	reg := 0b001100000001
	cpu.execOpcodes()

	if cpu.Memory[cpu.Registers.SP+1] != uint8((0b100&0xFF00)>>8) {
		t.Error("Expected ", uint8((0b100&0xFF00)>>8), " got ", cpu.Memory[cpu.Registers.SP+1])
	}
	if cpu.Memory[cpu.Registers.SP] != uint8(0b100&0xFF) {
		t.Error("Expected ", uint8(0b100&0xFF), " got ", cpu.Memory[cpu.Registers.SP])
	}

	if cpu.Registers.PC != uint16(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11011100 // CALL C, imm16 - unset flag
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	cpu.Registers.setFlag(flagC, false)
	cpu.execOpcodes()

	if cpu.Memory[cpu.Registers.SP+1] != uint8((0b100&0xFF00)>>8) {
		t.Error("Expected ", uint8((0b100&0xFF00)>>8), " got ", cpu.Memory[cpu.Registers.SP+1])
	}
	if cpu.Memory[cpu.Registers.SP] != uint8(0b100&0xFF) {
		t.Error("Expected ", uint8(0b100&0xFF), " got ", cpu.Memory[cpu.Registers.SP])
	}

	if cpu.Registers.PC != uint16(reg)+3 {
		t.Error("Expected ", reg+3, " got ", cpu.Registers.PC)
	}
}

func TestCALL_imm8(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b11001101 // CALL imm16
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	reg := 0b001100000001
	cpu.execOpcodes()

	if cpu.Memory[cpu.Registers.SP+1] != uint8((0b100&0xFF00)>>8) {
		t.Error("Expected ", uint8((0b100&0xFF00)>>8), " got ", cpu.Memory[cpu.Registers.SP+1])
	}
	if cpu.Memory[cpu.Registers.SP] != uint8(0b100&0xFF) {
		t.Error("Expected ", uint8(0b100&0xFF), " got ", cpu.Memory[cpu.Registers.SP])
	}

	if cpu.Registers.PC != uint16(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.PC)
	}
}

// LDH
func TestLDH(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b11100000 // LDH [a8], A
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	reg := 0xFF00 + uint16(0b1)
	cpu.execOpcodes()

	if cpu.Memory[reg] != cpu.Registers.A {
		t.Error("Expected ", cpu.Registers.A, " got ", cpu.Memory[reg])
	}

	cpu.Memory[cpu.Registers.PC] = 0b11100010 // LDH [a8], A
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.Memory[cpu.Registers.PC+2] = 0b11
	reg = 0xFF00 + uint16(0b11)
	cpu.execOpcodes()

	if cpu.Memory[reg] != cpu.Registers.A {
		t.Error("Expected ", cpu.Registers.A, " got ", cpu.Memory[reg])
	}

	cpu.Memory[cpu.Registers.PC] = 0b11110010 // LDH A, [C]
	reg = 0xFF00 + uint16(0b11)
	cpu.execOpcodes()

	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", uint8(reg), " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11110010 // LDH A, [C]
	reg = 0xFF00 + uint16(0b11)
	cpu.execOpcodes()

	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", uint8(reg), " got ", cpu.Registers.A)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11110000 // LDH A, [a8]
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	reg = 0xFF00 + uint16(0b1)
	cpu.execOpcodes()

	if cpu.Registers.A != cpu.Memory[reg] {
		t.Error("Expected ", cpu.Memory[reg], " got ", cpu.Registers.A)
	}
}

// ADD SP, e8
func TestADD_SP_e8(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b11101000 // ADD SP, e8
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	reg := uint16(int32(int8(16)) + int32(0b1))
	cpu.execOpcodes()

	if cpu.Registers.SP != reg {
		t.Error("Expected ", reg, " got ", cpu.Registers.SP)
	}
}

// LD HL, SP+e8
func TestLD_HL_SPe8(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			A:  0b11,
			B:  0b1,
			C:  0b11,
			D:  0b100,
			E:  0b101,
			F:  0b10,
			H:  0b110,
			L:  0b111,
			SP: 16,
			PC: 0b1,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

	cpu.Memory[cpu.Registers.PC] = 0b11111000 // LD HL, SP+e8
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	reg := uint16(int32(int8(16)) + int32(0b1))
	cpu.execOpcodes()

	if cpu.Registers.getHL() != reg {
		t.Error("Expected ", reg, " got ", cpu.Registers.getHL())
	}
}

// LD SP, HL
func TestLD_SP_HL(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b11111001 // LD SP, HL
	cpu.Memory[cpu.Registers.PC+1] = 0b1
	cpu.execOpcodes()

	if cpu.Registers.SP != cpu.Registers.getHL() {
		t.Error("Expected ", cpu.Registers.getHL(), " got ", cpu.Registers.SP)
	}
}

func TestStop(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b10000
	cpu.Memory[cpu.Registers.PC+1] = 0b0
	cpu.execOpcodes()

	if cpu.getImmediate8() != 0 {
		t.Error("Expected ", 0, " got ", cpu.getImmediate8())
	}

}

// DEC HL
func Test_DEC_HL(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b101011 //DEC HL
	val := cpu.Registers.getHL()
	val--
	cpu.execOpcodes()

	if cpu.Registers.getHL() != val {
		t.Error("Expected ", val, " got ", cpu.Registers.getHL())
	}

}

// INC SP
func Test_INC_SP(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b110011 // INC SP

	cpu.execOpcodes()

	if cpu.Registers.SP != 17 {
		t.Error("Expected ", 17, " got ", cpu.Registers.SP)
	}

}

// INC A
func Test_INC_A(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b111100 // INC A

	cpu.execOpcodes()

	if cpu.Registers.A != 0b11+1 {
		t.Error("Expected ", 0b11+1, " got ", cpu.Registers.A)
	}

}

// RET NZ
func Test_RET_NZ(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1
	cpu.Memory[cpu.Registers.PC] = 0b11000000 // RET NZ
	cpu.Memory[cpu.Registers.SP] = 0b1
	cpu.Memory[cpu.Registers.SP+1] = 0b11

	res := 0b11<<8 | 0b1
	cpu.execOpcodes()

	if cpu.Registers.PC != uint16(res) {
		t.Error("Expected ", uint16(res), " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11000000 // RET NZ
	cpu.Memory[cpu.Registers.SP] = 0b1
	cpu.Memory[cpu.Registers.SP+1] = 0b11
	cpu.Registers.setFlag(flagZ, true)
	cpu.execOpcodes()

	if cpu.Registers.PC != uint16(res)+1 {
		t.Error("Expected ", uint16(res), " got ", cpu.Registers.PC)
	}

}

// RET Z
func Test_RET_Z(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1
	cpu.Memory[cpu.Registers.PC] = 0b11001000 // RET Z
	cpu.Memory[cpu.Registers.SP] = 0b1
	cpu.Memory[cpu.Registers.SP+1] = 0b11
	cpu.Registers.setFlag(flagZ, true)

	res := 0b11<<8 | 0b1
	cpu.execOpcodes()

	if cpu.Registers.PC != uint16(res) {
		t.Error("Expected ", uint16(res), " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11001000 // RET Z - flagZ not set
	cpu.Memory[cpu.Registers.SP] = 0b1
	cpu.Memory[cpu.Registers.SP+1] = 0b11
	cpu.Registers.setFlag(flagZ, false)

	cpu.execOpcodes()

	if cpu.Registers.PC != uint16(res)+1 {
		t.Error("Expected ", uint16(res), " got ", cpu.Registers.PC)
	}
}

// RET NC
func Test_RET_NC(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1
	cpu.Memory[cpu.Registers.PC] = 0b11010000 // RET NC
	cpu.Memory[cpu.Registers.SP] = 0b1
	cpu.Memory[cpu.Registers.SP+1] = 0b11

	res := 0b11<<8 | 0b1
	cpu.execOpcodes()

	if cpu.Registers.PC != uint16(res) {
		t.Error("Expected ", uint16(res), " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11010000 // RET NC - flagC set
	cpu.Memory[cpu.Registers.SP] = 0b1
	cpu.Memory[cpu.Registers.SP+1] = 0b11
	cpu.Registers.setFlag(flagC, true)
	cpu.execOpcodes()

	if cpu.Registers.PC != uint16(res)+1 {
		t.Error("Expected ", uint16(res), " got ", cpu.Registers.PC)
	}
}

// RET C
func Test_RET_C(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b11011000 // RET C
	cpu.Memory[cpu.Registers.SP] = 0b1
	cpu.Memory[cpu.Registers.SP+1] = 0b11
	cpu.Registers.setFlag(flagC, true)
	res := 0b11<<8 | 0b1
	cpu.execOpcodes()

	if cpu.Registers.PC != uint16(res) {
		t.Error("Expected ", uint16(res), " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b0 // NOP
	cpu.execOpcodes()
	if cpu.Registers.PC != uint16(res)+1 {
		t.Error("Expected ", uint16(res)+1, " got ", cpu.Registers.PC)
	}

	cpu.Memory[cpu.Registers.PC] = 0b11011000 // RET C - flagC not set
	cpu.Memory[cpu.Registers.SP] = 0b1
	cpu.Memory[cpu.Registers.SP+1] = 0b11
	cpu.Registers.setFlag(flagC, false)
	cpu.execOpcodes()

	if cpu.Registers.PC != uint16(res)+2 {
		t.Error("Expected ", uint16(res)+2, " got ", cpu.Registers.PC)
	}
}

// RET
func TestRET(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1
	cpu.Memory[cpu.Registers.PC] = 0b11001001 // RET
	cpu.Memory[cpu.Registers.SP] = 0b1
	cpu.Memory[cpu.Registers.SP+1] = 0b11
	res := 0b11<<8 | 0b1
	cpu.execOpcodes()

	if cpu.Registers.PC != uint16(res) {
		t.Error("Expected ", res, " got ", cpu.Registers.PC)
	}

}

// RLCA
func Test_RLCA(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b111 // RLCA
	res := 0b11 << 1
	c := (0b11 & 0x80) >> 7
	fmt.Println("carry", c)
	cpu.execOpcodes()

	if cpu.Registers.A != uint8(res) {
		t.Error("Expected ", uint8(res), " got ", cpu.Registers.A)
	}
	if cpu.Registers.getFlag(flagC) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagC))
	}

}

// RRCA
func Test_RRCA(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b1111 // RRCA
	res := 0b11 >> 1
	c := 0b11 & 0x01
	fmt.Println("carry", c)
	cpu.execOpcodes()

	if cpu.Registers.A != uint8(res) {
		t.Error("Expected ", uint8(res), " got ", cpu.Registers.A)
	}
	if cpu.Registers.getFlag(flagC) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagC))
	}

}

// RLA
func Test_RLA(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b10111 // RLA
	res := (0b11 << 1) | 0
	c := 0b11 & 0x80
	fmt.Println("carry", c)
	cpu.execOpcodes()

	if cpu.Registers.A != uint8(res) {
		t.Error("Expected ", uint8(res), " got ", cpu.Registers.A)
	}
	if cpu.Registers.getFlag(flagC) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagC))
	}
	fmt.Printf("A %04b", cpu.Registers.A)

	cpu.Memory[cpu.Registers.PC] = 0b10111 // RLA
	res = (0b0110 << 1) | 1
	c = 0b0110 & 0x80
	cpu.Registers.setFlag(flagC, true)
	fmt.Println("carry", c)
	cpu.execOpcodes()

	if cpu.Registers.A != uint8(res) {
		t.Error("Expected ", uint8(res), " got ", cpu.Registers.A)
	}
	if cpu.Registers.getFlag(flagC) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagC))
	}

}

// RRA
func Test_RRA(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b11111 // RRA
	res := (0b11 >> 1) | 0<<7
	c := 0b11 & 0x01
	fmt.Println("carry", c)
	cpu.execOpcodes()

	if cpu.Registers.A != uint8(res) {
		t.Error("Expected ", uint8(res), " got ", cpu.Registers.A)
	}
	if cpu.Registers.getFlag(flagC) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagC))
	}
	fmt.Printf("A %04b", cpu.Registers.A)

	cpu.Memory[cpu.Registers.PC] = 0b11111 // RLA
	res = (0b01 >> 1) | 1<<7
	c = 0b01 & 0x01
	cpu.Registers.setFlag(flagC, true)
	fmt.Println("carry", c)
	cpu.execOpcodes()

	if cpu.Registers.A != uint8(res) {
		t.Error("Expected ", uint8(res), " got ", cpu.Registers.A)
	}
	if cpu.Registers.getFlag(flagC) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagC))
	}

}

// DAA
func Test_DAA(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b100111 // DAA
	lo := 0b11 % 10
	hi := ((0b11 - lo) % 100) / 10
	res := (hi << 4) | lo
	cpu.execOpcodes()

	if cpu.Registers.A != uint8(res) {
		t.Error("Expected ", uint8(res), " got ", cpu.Registers.A)
	}
	if cpu.Registers.getFlag(flagC) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagC))
	}
	if cpu.Registers.getFlag(flagZ) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagC))
	}
	if cpu.Registers.getFlag(flagH) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagC))
	}
}

// CPL
func TestCPL(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b101111 // CPL
	cpu.execOpcodes()
	reg := ^0b11
	if cpu.Registers.A != uint8(reg) {
		t.Error("Expected ", reg, " got ", cpu.Registers.A)
	}
}

// SCF
func TestSCF(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1

	cpu.Memory[cpu.Registers.PC] = 0b110111 // SCF
	cpu.execOpcodes()
	if cpu.Registers.getFlag(flagC) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagC))
	}
	if cpu.Registers.getFlag(flagN) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagN))
	}
	if cpu.Registers.getFlag(flagH) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagH))
	}
}

// CCF
func TestCCF(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1
	cpu.Memory[cpu.Registers.PC] = 0b111111 // CCF
	cpu.execOpcodes()
	if cpu.Registers.getFlag(flagC) != true {
		t.Error("Expected true, got ", cpu.Registers.getFlag(flagC))
	}
	if cpu.Registers.getFlag(flagN) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagC))
	}
	if cpu.Registers.getFlag(flagH) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagC))
	}

	cpu.Memory[cpu.Registers.PC] = 0b111111 // CCF
	cpu.Registers.setFlag(flagC, true)
	cpu.execOpcodes()
	if cpu.Registers.getFlag(flagC) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagC))
	}
	if cpu.Registers.getFlag(flagN) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagC))
	}
	if cpu.Registers.getFlag(flagH) != false {
		t.Error("Expected false, got ", cpu.Registers.getFlag(flagC))
	}

}

// RETI
func TestRETI(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1
	cpu.Memory[cpu.Registers.PC] = 0b11011001 // RETI
	cpu.Memory[cpu.Registers.SP] = 0b1
	cpu.Memory[cpu.Registers.SP+1] = 0b11
	res := 0b11<<8 | 0b1
	cpu.execOpcodes()

	if cpu.Registers.PC != uint16(res) {
		t.Error("Expected ", res, " got ", cpu.Registers.PC)
	}
}

// DI
func TestDI(t *testing.T) {
	cpu := NewCPU()
	cpu.Registers.A = 0b11
	cpu.Registers.B = 0b1
	cpu.Registers.C = 0b11
	cpu.Registers.D = 0b100
	cpu.Registers.E = 0b101
	cpu.Registers.F = 0b10
	cpu.Registers.H = 0b110
	cpu.Registers.L = 0b111
	cpu.Registers.SP = 16
	cpu.Registers.PC = 0b1
	cpu.Memory[cpu.Registers.PC] = 0b11110011 // DI
	cpu.execOpcodes()

	if cpu.IME != false {
		t.Error("Expected ", false, " got ", cpu.IME)
	}
}

func TestEI(t *testing.T) {
	cpu := NewCPU()

	cpu.Registers.PC = 0b0
	cpu.Memory[cpu.Registers.PC] = 0b11111011 // EI
	cpu.Memory[cpu.Registers.PC+1] = 0b0      // NOP

	cpu.Registers.PC = 0
	fmt.Println("Ime", cpu.IMEScheduled)
	cpu.execOpcodes()
	fmt.Println("Ime", cpu.IMEScheduled)

	if cpu.IMEScheduled != true {
		t.Error("Expected ", true, " got ", cpu.IMEScheduled)
	}

	if cpu.IME != false {
		t.Error("Expected false, got ", cpu.IME)
	}

	cpu.execOpcodes()
	if cpu.IME != false {
		t.Error("Expected ", false, " got ", cpu.IME)
	}
}
