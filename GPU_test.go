package main

import (
	"golang.org/x/exp/slices"

	"testing"
)

func TestGetLCDC(t *testing.T) {
	cpu := NewCPU()
	graphic := NewGraphics(cpu)
	res := graphic.getLCDC()
	if res != 0x91 {
		t.Error("Expected 0x91, got ", graphic.cpu.Memory[0xFF40])

	}
}

func TestGetLY(t *testing.T) {
	cpu := NewCPU()
	graphic := NewGraphics(cpu)
	res := graphic.getLY()
	if res != 0 {
		t.Error("Expected 0x00, got ", graphic.cpu.Memory[0xFF44])

	}
}

func TestGetSTAT(t *testing.T) {
	cpu := NewCPU()
	graphic := NewGraphics(cpu)
	res := graphic.getSTAT()
	if res != 0b00000001 {
		t.Error("Expected 0b00000001, got ", graphic.cpu.Memory[0xFF41])

	}
}

func TestGetSCY(t *testing.T) {
	cpu := NewCPU()
	graphic := NewGraphics(cpu)
	res := graphic.getSCY()
	if res != 0 {
		t.Error("Expected 0, got ", graphic.cpu.Memory[0xFF42])

	}
}

func TestGetSCX(t *testing.T) {
	cpu := NewCPU()
	graphic := NewGraphics(cpu)
	res := graphic.getSCX()
	if res != 0 {
		t.Error("Expected 0, got ", graphic.cpu.Memory[0xFF43])

	}
}

func TestGetWY(t *testing.T) {
	cpu := NewCPU()
	graphic := NewGraphics(cpu)
	res := graphic.getWY()
	if res != 0 {
		t.Error("Expected 0, got ", graphic.cpu.Memory[0xFF4A])

	}
}

func TestGetWX(t *testing.T) {
	cpu := NewCPU()
	graphic := NewGraphics(cpu)
	res := graphic.getWX()
	if res != 0 {
		t.Error("Expected 0, got ", graphic.cpu.Memory[0xFF4B])

	}
}

func TestDMATransfer(t *testing.T) {
	cpu := NewCPU()
	graphic := NewGraphics(cpu)
	graphic.cpu.Memory[768] = 10

	cpu.memoryWrite(0xFF46, 3)
	//upper := uint16(3) << 8
	//// fmt.Println("upper", upper)
	res := graphic.cpu.Memory[OAM_START]
	if res != 10 {
		t.Error("Expected 10, got ", res)
	}
}

func TestReadTileData(t *testing.T) {
	cpu := NewCPU()
	graphic := NewGraphics(cpu)
	graphic.cpu.Memory[0x8003] = 0b00111100
	graphic.cpu.Memory[0x8004] = 0b01111110
	graphic.cpu.Memory[0x8005] = 0b01000010
	graphic.cpu.Memory[0x8006] = 0b01000010
	graphic.cpu.Memory[0x8007] = 0b01000010
	graphic.cpu.Memory[0x8008] = 0b01000010
	graphic.cpu.Memory[0x8009] = 0b01000010
	graphic.cpu.Memory[0x800A] = 0b01000010
	graphic.cpu.Memory[0x800B] = 0b01111110
	graphic.cpu.Memory[0x800C] = 0b01011110
	graphic.cpu.Memory[0x800D] = 0b01111110
	graphic.cpu.Memory[0x800E] = 0b00001010
	graphic.cpu.Memory[0x800F] = 0b01111100
	graphic.cpu.Memory[0x8010] = 0b01010110
	graphic.cpu.Memory[0x8011] = 0b00111000
	graphic.cpu.Memory[0x8012] = 0b01111100

	tile := cpu.readTileData(0x8003)
	//row := 0
	//col := 0
	expected := [8][8]uint8{
		{0b00, 0b10, 0b11, 0b11, 0b11, 0b11, 0b10, 0b00},
		{0b00, 0b11, 0b00, 0b00, 0b00, 0b00, 0b11, 0b00},
		{0b00, 0b11, 0b00, 0b00, 0b00, 0b00, 0b11, 0b00},
		{0b00, 0b11, 0b00, 0b00, 0b00, 0b00, 0b11, 0b00},
		{0b00, 0b11, 0b01, 0b11, 0b11, 0b11, 0b11, 0b00},
		{0b00, 0b01, 0b01, 0b01, 0b11, 0b01, 0b11, 0b00},
		{0b00, 0b11, 0b01, 0b11, 0b01, 0b11, 0b10, 0b00},
		{0b00, 0b10, 0b11, 0b11, 0b11, 0b10, 0b00, 0b00},
	}
	//// fmt.Printf("tile: %v\n", tile)
	if tile != expected {
		t.Error("Expected ", expected, "got ", tile)
	}
}

func TestGetTilePixel(t *testing.T) {
	cpu := NewCPU()
	graphic := NewGraphics(cpu)
	graphic.cpu.Memory[0x8003] = 0b00111100
	graphic.cpu.Memory[0x8004] = 0b01111110
	graphic.cpu.Memory[0x8005] = 0b01000010
	graphic.cpu.Memory[0x8006] = 0b01000010
	graphic.cpu.Memory[0x8007] = 0b01000010
	graphic.cpu.Memory[0x8008] = 0b01000010
	graphic.cpu.Memory[0x8009] = 0b01000010
	graphic.cpu.Memory[0x800A] = 0b01000010
	graphic.cpu.Memory[0x800B] = 0b01111110
	graphic.cpu.Memory[0x800C] = 0b01011110
	graphic.cpu.Memory[0x800D] = 0b01111110
	graphic.cpu.Memory[0x800E] = 0b00001010
	graphic.cpu.Memory[0x800F] = 0b01111100
	graphic.cpu.Memory[0x8010] = 0b01010110
	graphic.cpu.Memory[0x8011] = 0b00111000
	graphic.cpu.Memory[0x8012] = 0b01111100
	var row []uint8
	for x := 0; x < 8; x++ {
		colorID := graphic.getTilePixel(0x8003, x, 1)
		//// fmt.Println("colorID", colorID)
		row = append(row, colorID)
	}

	expected := []uint8{0b00, 0b11, 0b00, 0b00, 0b00, 0b00, 0b11, 0b00}
	if slices.Equal(row, expected) != true {
		t.Error("Expected ", expected, "got ", row)
	}
}
