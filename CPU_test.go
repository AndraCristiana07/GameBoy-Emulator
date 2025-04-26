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
	cpu := CPU{
		Registers: Registers{
			A: 0b1,
			B: 0b10,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

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
	cpu := CPU{
		Registers: Registers{
			B: 0b1,
			C: 0b10,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

	cpu.Memory[cpu.Registers.PC] = 0b1000001
	cpu.execOpcodes()
	if cpu.Registers.B != 0b10 {
		t.Error("Expected 0b10, got ", cpu.Registers.B)
	}
}

// LD A, n
func TestExecLD_A(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			A: 0b11,
			B: 0b1,
			C: 0b10,
			D: 0b100,
			E: 0b101,
			H: 0b110,
			L: 0b111,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}
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
}

// LD B, n
func TestExecLD_B(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			A: 0b11,
			B: 0b1,
			C: 0b10,
			D: 0b100,
			E: 0b101,
			H: 0b110,
			L: 0b111,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}
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
}

// LD C, n
func TestExecLD_C(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			A: 0b11,
			B: 0b1,
			C: 0b10,
			D: 0b100,
			E: 0b101,
			H: 0b110,
			L: 0b111,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}
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
}

// LD D, n
func TestExecLD_D(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			A: 0b11,
			B: 0b1,
			C: 0b10,
			D: 0b100,
			E: 0b101,
			H: 0b110,
			L: 0b111,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}
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
}

// LD E, n
func TestExecLD_E(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			A: 0b11,
			B: 0b1,
			C: 0b10,
			D: 0b100,
			E: 0b101,
			H: 0b110,
			L: 0b111,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}
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
}

// LD H, n
func TestExecLD_H(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			A: 0b11,
			B: 0b1,
			C: 0b10,
			D: 0b100,
			E: 0b101,
			H: 0b110,
			L: 0b111,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}
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
}

// LD L, n
func TestExecLD_L(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			A: 0b11,
			B: 0b1,
			C: 0b10,
			D: 0b100,
			E: 0b101,
			H: 0b110,
			L: 0b111,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}
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
}

// LD [HL], n
func TestExecLD_HL(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			A: 0b11,
			B: 0b1,
			C: 0b10,
			D: 0b100,
			E: 0b101,
			H: 0b110,
			L: 0b111,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}
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
	cpu := CPU{
		Registers: Registers{
			H: 0b10100,
			L: 0b110110,
			B: 0b10,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

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
	cpu := CPU{
		Registers: Registers{
			A: 0b1,
			B: 0b10,
			C: 0b10,
			D: 0b100,
			E: 0b101,
			H: 0b110,
			L: 0b111,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}
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

}

// ADC A,n
func TestExecADCA(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			A: 0b1,
			B: 0b10,
			C: 0b10,
			D: 0b100,
			E: 0b101,
			H: 0b110,
			L: 0b111,
			F: 0b000010000, // flagC on
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}
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

}

// SUB A, n
func TestSUB_A(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			A: 0b11110, //30
			B: 0b10,
			C: 0b10,
			D: 0b100,
			E: 0b101,
			H: 0b110,
			L: 0b111,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

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
}

// SBC A, n
func TestSBC_A(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			A: 0b11110, //30
			B: 0b10,
			C: 0b10,
			D: 0b100,
			E: 0b101,
			H: 0b110,
			L: 0b111,
			F: 0b000010000, // flagC on
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

	var res uint8
	cpu.Memory[cpu.Registers.PC] = 0b10011000 // SBC A, B
	cpu.execOpcodes()
	res = uint8(0b11110 - 0b10 - 1)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.PC] = 0b10011001 // SBC A, C
	cpu.execOpcodes()

	res -= uint8(0b10 + uint8(1))
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.PC] = 0b10011010 // SBC A, D
	cpu.execOpcodes()
	res -= 0b100 + uint8(1)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.PC] = 0b10011011 // SBC A, E
	cpu.execOpcodes()
	res -= 0b101 + uint8(1)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.PC] = 0b10011100 // SBC A, H
	cpu.execOpcodes()
	res -= 0b110 + uint8(1)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.PC] = 0b10011101 // SBC A, L
	cpu.execOpcodes()
	res -= 0b111 + uint8(1)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.getHL()] = 0b1
	cpu.Memory[cpu.Registers.PC] = 0b10011110 // SBC A, [HL]
	cpu.execOpcodes()
	res -= 0b1 + uint8(1)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}

	cpu.Registers.F = 0b000010000
	cpu.Memory[cpu.Registers.PC] = 0b11011110 // SBC A, imm8
	cpu.Memory[cpu.Registers.PC+1] = 0b11010111
	cpu.execOpcodes()
	res -= 0b11010111 + uint8(1)
	if cpu.Registers.A != res {
		t.Errorf("Expected 0b%b, got 0b%b ", res, cpu.Registers.A)
	}
}

// ADD HL, n
func TestExecADD_HL_BC(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			H: 0b1,
			L: 0b10,
			B: 0b10,
			C: 0b110,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

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
	cpu := CPU{
		Registers: Registers{
			B: 0b1,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

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
	cpu := CPU{
		Registers: Registers{
			C: 0b1,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

	cpu.Memory[cpu.Registers.PC] = 0b1100
	cpu.execOpcodes()

	if cpu.Registers.C != 0b10 {
		t.Error("Expected 0b10, got ", cpu.Registers.C)
	}
}

// INC D
func TestINC_D(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			D: 0b1,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

	cpu.Memory[cpu.Registers.PC] = 0b10100
	cpu.execOpcodes()

	if cpu.Registers.D != 0b10 {
		t.Error("Expected 0b10, got ", cpu.Registers.D)
	}
}

// INC E
func TestINC_E(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			E: 0b1,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

	cpu.Memory[cpu.Registers.PC] = 0b11100
	cpu.execOpcodes()

	if cpu.Registers.E != 0b10 {
		t.Error("Expected 0b10, got ", cpu.Registers.E)
	}
}

// INC H
func TestINC_H(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			H: 0b1,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

	cpu.Memory[cpu.Registers.PC] = 0b100100
	cpu.execOpcodes()

	if cpu.Registers.H != 0b10 {
		t.Error("Expected 0b10, got ", cpu.Registers.H)
	}
}

// INC L
func TestINC_L(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			L: 0b1,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

	cpu.Memory[cpu.Registers.PC] = 0b101100
	cpu.execOpcodes()

	if cpu.Registers.L != 0b10 {
		t.Error("Expected 0b10, got ", cpu.Registers.L)
	}
}

// INC BC
func TestINC_BC(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			B: 0b1,
			C: 0b10,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}
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
	cpu := CPU{
		Registers: Registers{
			D: 0b1,
			E: 0b10,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

	cpu.Memory[cpu.Registers.PC] = 0b11
	cpu.execOpcodes()
	if cpu.Registers.getDE() != 0b100000010 {
		t.Error("Expected 0b100000010, got ", cpu.Registers.getDE())
	}
}

// INC HL
func TestINC_HL(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			H: 0b1,
			L: 0b10,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}
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
	cpu := CPU{
		Registers: Registers{
			H: 0b1,
			L: 0b10,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}
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
	cpu := CPU{
		Registers: Registers{
			A: 0b110,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

	cpu.Memory[cpu.Registers.PC] = 0b111101
	cpu.execOpcodes()

	if cpu.Registers.A != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.A)
	}
}

// DEC B
func TestDEC_B(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			B: 0b110,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

	cpu.Memory[cpu.Registers.PC] = 0b101
	cpu.execOpcodes()

	if cpu.Registers.B != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.B)
	}
}

// DEC C
func TestDEC_C(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			C: 0b110,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

	cpu.Memory[cpu.Registers.PC] = 0b1101
	cpu.execOpcodes()

	if cpu.Registers.C != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.C)
	}
}

// DEC D
func TestDEC_D(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			D: 0b110,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

	cpu.Memory[cpu.Registers.PC] = 0b10101
	cpu.execOpcodes()

	if cpu.Registers.D != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.D)
	}
}

// DEC E
func TestDEC_E(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			E: 0b110,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

	cpu.Memory[cpu.Registers.PC] = 0b11101
	cpu.execOpcodes()

	if cpu.Registers.E != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.E)
	}
}

// DEC H
func TestDEC_H(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			H: 0b110,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

	cpu.Memory[cpu.Registers.PC] = 0b100101
	cpu.execOpcodes()

	if cpu.Registers.H != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.H)
	}
}

// DEC L
func TestDEC_L(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			L: 0b110,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

	cpu.Memory[cpu.Registers.PC] = 0b101101
	cpu.execOpcodes()

	if cpu.Registers.L != 0b101 {
		t.Error("Expected 0b101, got ", cpu.Registers.L)
	}
}

// DEC DE
func TestDEC_DE(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			D: 0b1,
			E: 0b10,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}
	cpu.Registers.setDE(0b100)
	//operands := []map[string]string{
	//	{"name": "DE", "immediate": "True"},
	//}
	//flags := map[string]string{
	//	"Z": "-",
	//	"N": "-",
	//	"H": "-",
	//	"C": "-",
	//}
	//cpu.execDEC(operands, flags)
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
	cpu := CPU{
		Registers: Registers{
			B: 0b1,
			C: 0b10,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

	cpu.Memory[cpu.Registers.PC] = 0b1011
	cpu.execOpcodes()
	//// fmt.Printf("0b%b", cpu.Registers.getBC())

	if cpu.Registers.getBC() != 0b100000001 {
		t.Error("Expected 0b100000001, got ", cpu.Registers.getBC())
	}
}

// DEC [HL]
func TestDEC_HLmem(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			H: 0b1,
			L: 0b10,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}
	cpu.Registers.setHL(0b100)
	hladdr := cpu.Registers.getHL()
	cpu.Memory[hladdr] = 0b100
	//
	//operands := []map[string]string{
	//	{"name": "HL", "immediate": "False"},
	//}
	//flags := map[string]string{
	//	"Z": "Z",
	//	"N": "1",
	//	"H": "H",
	//	"C": "-",
	//}

	//cpu.execDEC(operands, flags)
	cpu.Memory[cpu.Registers.PC] = 0b110101
	cpu.execOpcodes()
	if cpu.Memory[cpu.Registers.getHL()] != 0b11 {
		t.Error("Expected 0b11, got ", cpu.Memory[cpu.Registers.getHL()])
	}
}

// DEC SP
func TestDEC_SP(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			SP: 0b100,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}
	//operands := []map[string]string{
	//	{"name": "SP", "immediate": "True"},
	//}
	//flags := map[string]string{
	//	"Z": "-",
	//	"N": "-",
	//	"H": "-",
	//	"C": "-",
	//}
	//cpu.execDEC(operands, flags)
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

// AND A, B
func TestAND(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			A: 0b1,
			B: 0b1,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

	cpu.Memory[cpu.Registers.PC] = 0b10100000
	cpu.execOpcodes()

	if cpu.Registers.A != 0b1 {
		t.Error("Expected 0b1, got ", cpu.Registers.A)
	}

}

// CP A, n
func TestCP(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			A: 0b11,
			B: 0b1,
			C: 0b11,
			D: 0b100,
			E: 0b101,
			H: 0b110,
			L: 0b111,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}
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
}

// PUSH
func TestPush(t *testing.T) {
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
			SP: 15,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}
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
			SP: 10,
		},
		Memory: [65536]uint8(make([]uint8, 65536)),
	}

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
