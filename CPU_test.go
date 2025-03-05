package main

import (
	"fmt"
	"testing"
)

func TestSetBC(t *testing.T) {
	register := Registers{}
	register.setBC(0x0001)
	fmt.Println(register.getBC())
	if register.getBC() != 0x0001 {
		t.Error("Expected 0x0001, got ", register.getBC())
	}
}

func TestSetDE(t *testing.T) {
	register := Registers{}
	register.setDE(0x0001)
	fmt.Println(register.getDE())
	if register.getDE() != 0x0001 {
		t.Error("Expected 0x0001, got ", register.getDE())
	}
}

func TestSetHL(t *testing.T) {
	register := Registers{}
	register.setHL(0x0001)
	fmt.Println(register.getHL())
	if register.getHL() != 0x0001 {
		t.Error("Expected 0x0001, got ", register.getDE())
	}
}

//func TestAssignToA(t *testing.T) {
//	register := Registers{}
//
//}

//func TestSetFlagZ(t *testing.T) {
//	register := Registers{}
//	register.setFlag(flagZ, true)
//	fmt.Println(register.getAF())
//
//}

// LD A,B
func TestExecLDAB(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			A: 0x01,
			B: 0x02,
		},
		Memory: [8192]uint8(make([]uint8, 8192)),
	}

	operands := []map[string]string{
		{"name": "A", "immediate": "True"},
		{"name": "B", "immediate": "True"},
	}
	cpu.execLD(operands)
	fmt.Println(cpu.Registers.A)
	if cpu.Registers.A != 0x02 {
		t.Error("Expected 0x0002, got ", cpu.Registers.A)
	}
}

// LD B, [HL]
func TestExecLDHLMem(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			H: 0x12,
			L: 0x34,
			B: 0x02,
		},
		Memory: [8192]uint8(make([]uint8, 8192)),
	}

	hladdr := cpu.Registers.getHL()
	cpu.Memory[hladdr] = 0x01

	operands := []map[string]string{
		{"name": "B", "immediate": "True"},
		{"name": "HL", "immediate": "False"},
	}
	cpu.execLD(operands)
	fmt.Println(cpu.Registers.B)
	if cpu.Registers.B != 0x01 {
		t.Error("Expected 0x0002, got ", cpu.Registers.A)
	}
}

// LD [HL], B
func TestExecLDHL(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			H: 0x14,
			L: 0x36,
			B: 0x02,
		},
		Memory: [8192]uint8(make([]uint8, 8192)),
	}

	hladdr := cpu.Registers.getHL()
	cpu.Memory[hladdr] = 0x001

	operands := []map[string]string{
		{"name": "HL", "immediate": "False"},
		{"name": "B", "immediate": "True"},
	}
	cpu.execLD(operands)
	fmt.Println("aaaaaaaaa")
	fmt.Println(cpu.Registers.B)
	fmt.Println(cpu.Memory[hladdr])
	if cpu.Memory[hladdr] != 0x0002 {
		t.Error("Expected 0x0002, got ", cpu.Registers.A)
	}
}

// LD H, [HL]
func TestExecLDHHL(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			H: 0x14,
			L: 0x06,
			B: 0x02,
		},
		Memory: [8192]uint8(make([]uint8, 8192)),
	}

	hladdr := cpu.Registers.getHL()
	cpu.Memory[hladdr] = 0x001

	operands := []map[string]string{
		{"name": "H", "immediate": "True"},
		{"name": "HL", "immediate": "False"},
	}
	cpu.execLD(operands)
	fmt.Println(cpu.Registers.H)
	if cpu.Registers.H != 0x001 {
		t.Error("Expected 0x0002, got ", cpu.Registers.A)
	}
}

// ADD A,B
func TestExecADDAB(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			A: 0x01,
			B: 0x02,
		},
		Memory: [8192]uint8(make([]uint8, 8192)),
	}

	operands := []map[string]string{
		{"name": "A", "immediate": "True"},
		{"name": "B", "immediate": "True"},
	}

	flags := map[string]string{
		"Z": "Z",
		"N": "0",
		"H": "H",
		"C": "C",
	}
	cpu.execADD(operands, flags)
	fmt.Println(cpu.Registers.A)
	if cpu.Registers.A != 0x03 {
		t.Error("Expected 0x0003, got ", cpu.Registers.A)
	}

}

// ADD HL, BC
func TestExecADD_HL_BC(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			H: 0x01,
			L: 0x02,
			B: 0x02,
			C: 0x06,
		},
		Memory: [8192]uint8(make([]uint8, 8192)),
	}

	cpu.Registers.setHL(0x0002)
	cpu.Registers.setBC(0x0004)
	operands := []map[string]string{
		{"name": "HL", "immediate": "True"},
		{"name": "BC", "immediate": "True"},
	}
	flags := map[string]string{
		"Z": "Z",
		"N": "0",
		"H": "0",
		"C": "C",
	}
	cpu.execADD(operands, flags)
	fmt.Println(cpu.Registers.getHL())
	if cpu.Registers.getHL() != 0x0006 {
		t.Error("Expected 0x0006, got ", cpu.Registers.A)
	}

}

// INC B
func TestINC_B(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			B: 0x01,
		},
		Memory: [8192]uint8(make([]uint8, 8192)),
	}
	operands := []map[string]string{
		{"name": "B", "immediate": "True"},
	}
	flags := map[string]string{
		"Z": "Z",
		"N": "0",
		"H": "H",
		"C": "-",
	}
	cpu.execINC(operands, flags)
	fmt.Println(cpu.Registers.B)

	if cpu.Registers.B != 0x02 {
		t.Error("Expected 0x02, got ", cpu.Registers.A)
	}
}

// INC BC
func TestINC_BC(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			B: 0x01,
			C: 0x02,
		},
		Memory: [8192]uint8(make([]uint8, 8192)),
	}
	cpu.Registers.setBC(0x0004)
	operands := []map[string]string{
		{"name": "BC", "immediate": "True"},
	}
	flags := map[string]string{
		"Z": "-",
		"N": "-",
		"H": "-",
		"C": "-",
	}
	cpu.execINC(operands, flags)
	fmt.Println(cpu.Registers.getBC())
	if cpu.Registers.getBC() != 0x0005 {
		t.Error("Expected 0x0005, got ", cpu.Registers.A)
	}
}

// INC HL
func TestINC_HL(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			H: 0x01,
			L: 0x02,
		},
		Memory: [8192]uint8(make([]uint8, 8192)),
	}
	cpu.Registers.setHL(0x0004)
	operands := []map[string]string{
		{"name": "HL", "immediate": "True"},
	}
	flags := map[string]string{
		"Z": "-",
		"N": "-",
		"H": "-",
		"C": "-",
	}
	cpu.execINC(operands, flags)
	fmt.Println(cpu.Registers.getHL())
	if cpu.Registers.getHL() != 0x0005 {
		t.Error("Expected 0x0005, got ", cpu.Registers.A)
	}
}

// INC [HL]
func TestINC_HLmem(t *testing.T) {
	cpu := CPU{
		Registers: Registers{
			H: 0x01,
			L: 0x02,
		},
		Memory: [8192]uint8(make([]uint8, 8192)),
	}
	cpu.Registers.setHL(0x0004)
	hladdr := cpu.Registers.getHL()
	cpu.Memory[hladdr] = 0x004

	operands := []map[string]string{
		{"name": "HL", "immediate": "False"},
	}
	flags := map[string]string{
		"Z": "Z",
		"N": "0",
		"H": "H",
		"C": "-",
	}
	cpu.execINC(operands, flags)
	fmt.Println(cpu.Memory[cpu.Registers.getHL()])
	if cpu.Memory[cpu.Registers.getHL()] != 0x0005 {
		t.Error("Expected 0x0005, got ", cpu.Registers.A)
	}
}
