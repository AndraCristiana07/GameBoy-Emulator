package main

import (
	"fmt"
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

func TestDrawScreen(t *testing.T) {
	cpu := NewCPU()
	graphic := NewGraphics(cpu)

	graphic.pixelBuffer[0][0] = 1
	graphic.pixelBuffer[1][0] = 0
	graphic.pixelBuffer[2][0] = 2
	graphic.pixelBuffer[3][0] = 1
	graphic.pixelBuffer[0][1] = 2
	graphic.pixelBuffer[1][1] = 1
	graphic.pixelBuffer[2][1] = 0
	graphic.pixelBuffer[3][1] = 1

}

func TestSpritesOAM1(t *testing.T) {
	graphic := NewGraphics(NewCPU())
	graphic.cpu.Memory[0xFF40] = 0b00000010 // OBJ enable
	graphic.cpu.Memory[0xFF44] = 10         // LY = 10

	spriteIndex := 0
	graphic.cpu.Memory[OAM_START+spriteIndex*4] = 10 + 16  // y=10
	graphic.cpu.Memory[OAM_START+spriteIndex*4+1] = 16 + 8 // x=16
	graphic.cpu.Memory[OAM_START+spriteIndex*4+2] = 0      // tile
	graphic.cpu.Memory[OAM_START+spriteIndex*4+3] = 0      // attributes

	//tile 0 in VRAM (8x8)
	tileBase := uint16(0x8000)
	for i := 0; i < 16; i += 2 {
		graphic.cpu.Memory[tileBase+uint16(i)] = 0b11111111
		graphic.cpu.Memory[tileBase+uint16(i+1)] = 0b00000000
	}

	spritePixels := graphic.spritesOAM()

	if spritePixels[10][16] != 1 {
		t.Errorf("Expected to be 1, got %d", spritePixels[10][16])
	}

}

func TestSpritesOAM2(t *testing.T) {
	graphic := NewGraphics(NewCPU())
	graphic.cpu.Memory[0xFF40] = 0b00000010 // OBJ enable
	graphic.cpu.Memory[0xFF44] = 12         // LY = 12

	spriteIndex := 0
	graphic.cpu.Memory[OAM_START+spriteIndex*4] = 12 + 16  // y=12
	graphic.cpu.Memory[OAM_START+spriteIndex*4+1] = 16 + 8 // x=16
	graphic.cpu.Memory[OAM_START+spriteIndex*4+2] = 0      // tile
	graphic.cpu.Memory[OAM_START+spriteIndex*4+3] = 0      // attributes

	//tile 0 in VRAM (8x8)
	tileBase := uint16(0x8000)
	for i := 0; i < 16; i += 2 {
		graphic.cpu.Memory[tileBase+uint16(i)] = 0b00000000
		graphic.cpu.Memory[tileBase+uint16(i+1)] = 0b11111111
	}

	spritePixels := graphic.spritesOAM()

	if spritePixels[12][16] != 2 {
		t.Errorf("Expected to be 3, got %d", spritePixels[12][16])
	}
}

func TestSpritesOAM3(t *testing.T) {
	graphic := NewGraphics(NewCPU())
	graphic.cpu.Memory[0xFF40] = 0b00000010 // OBJ enable
	graphic.cpu.Memory[0xFF44] = 14         // LY = 14

	spriteIndex := 0
	graphic.cpu.Memory[OAM_START+spriteIndex*4] = 14 + 16  // y=14
	graphic.cpu.Memory[OAM_START+spriteIndex*4+1] = 16 + 8 // x=16
	graphic.cpu.Memory[OAM_START+spriteIndex*4+2] = 0      // tile
	graphic.cpu.Memory[OAM_START+spriteIndex*4+3] = 0      // attributes

	//tile 0 in VRAM (8x8)
	tileBase := uint16(0x8000)
	for i := 0; i < 16; i += 2 {
		graphic.cpu.Memory[tileBase+uint16(i)] = 0b11111111
		graphic.cpu.Memory[tileBase+uint16(i+1)] = 0b11111111
	}

	spritePixels := graphic.spritesOAM()

	if spritePixels[14][16] != 3 {
		t.Errorf("Expected to be 3, got %d", spritePixels[14][16])
	}
}

func TestSpritesOAM4(t *testing.T) {
	graphic := NewGraphics(NewCPU())
	graphic.cpu.Memory[0xFF40] = 0b00000110 // OBJ enable
	graphic.cpu.Memory[0xFF44] = 14         // LY = 14

	spriteIndex := 0
	graphic.cpu.Memory[OAM_START+spriteIndex*4] = 14 + 16  // y=14
	graphic.cpu.Memory[OAM_START+spriteIndex*4+1] = 16 + 8 // x=16
	graphic.cpu.Memory[OAM_START+spriteIndex*4+2] = 0      // tile
	graphic.cpu.Memory[OAM_START+spriteIndex*4+3] = 0      // attributes

	//tile 0 in VRAM (8x8)
	tileBase := uint16(0x8000)
	for i := 0; i < 16; i += 2 {
		graphic.cpu.Memory[tileBase+uint16(i)] = 0b11111111
		graphic.cpu.Memory[tileBase+uint16(i+1)] = 0b11111111
	}

	spritePixels := graphic.spritesOAM()

	if spritePixels[14][16] != 3 {
		t.Errorf("Expected to be 3, got %d", spritePixels[14][16])
	}
}

func TestGetBackground(t *testing.T) {
	graphic := &Graphics{
		cpu: &CPU{
			Memory: [65536]uint8{},
		},
	}

	graphic.cpu.Memory[0xFF40] = 0b10010001
	graphic.cpu.Memory[0xFF44] = 0 // LY = 0
	graphic.cpu.Memory[0xFF42] = 0 // SCY
	graphic.cpu.Memory[0xFF43] = 0 // SCX

	// tile index 0 (0x9800)
	graphic.cpu.Memory[0x9800] = 0

	tileAddr := uint16(0x8000)
	for i := 0; i < 16; i += 2 {
		graphic.cpu.Memory[tileAddr+uint16(i)] = 0b11111111
		graphic.cpu.Memory[tileAddr+uint16(i+1)] = 0b00000000
	}
	bg := graphic.getBackground()

	if bg[0][0] != 1 {
		t.Errorf("Expected pixel (0, 0) to be 1, got %d", bg[0][0])
	}
}

func TestGetWindow(t *testing.T) {
	graphic := &Graphics{
		cpu: &CPU{
			Memory: [65536]uint8{},
		},
	}

	//graphic.cpu.Memory[0xFF40] = 0b11100000
	graphic.cpu.Memory[0xFF40] = 0b10110001
	graphic.cpu.Memory[0xFF44] = 0 // LY = 0
	graphic.cpu.Memory[0xFF4A] = 0 // WY
	graphic.cpu.Memory[0xFF4B] = 7 // WX

	graphic.cpu.Memory[0x9800] = 0

	if graphic.getLCDC()&(1<<5) == 0 {
		fmt.Println("AAAAAAAAAAAAAAAANOOO")
	}
	tileAddr := uint16(0x8000)
	for i := 0; i < 16; i += 2 {
		graphic.cpu.Memory[tileAddr+uint16(i)] = 0b11111111
		graphic.cpu.Memory[tileAddr+uint16(i+1)] = 0b00000000

	}

	win := graphic.getWindow()

	if win[0][7] != 1 {
		t.Errorf("Expected window pixel at (0, 7) to be 1, got %d", win[0][7])
	}
}
func TestSetMode(t *testing.T) {
	graphic := NewGraphics(NewCPU())
	graphic.cpu.Memory[0xFF41] = 0b00011001
	graphic.setMode(MODE_HBLANK)
	if graphic.cpu.Memory[0xFF41] != 0b00011000 {
		t.Errorf("Expected to be %d, got %d", 0b00011000, graphic.cpu.Memory[0xFF41])
	}
}

func TestModesHandling(t *testing.T) {
	graphic := NewGraphics(NewCPU())
	graphic.cpu.Memory[0xFF40] = 0b10000110
	graphic.cpu.Memory[0xFF44] = 10 // LY = 10
	graphic.modesHandling(20)
	if graphic.getSTAT() != 0b00000010 {
		t.Errorf("Expected to be %08b, got %08b", 0b00000010, graphic.getSTAT())
	}

	graphic.cpu.Memory[0xFF40] = 0b00000110
	graphic.cpu.Memory[0xFF44] = 10 // LY = 10
	graphic.modesHandling(20)
	if graphic.getLCDC()&(1<<7) == 0 {
		fmt.Println("AAAA")
	}
	if graphic.getSTAT() != 0b00000000 {
		t.Errorf("Expected to be %08b, got %08b", 0b00000000, graphic.getSTAT())
	}

	graphic.cpu.Memory[0xFF40] = 0b10110110
	graphic.cpu.Memory[0xFF44] = 145 // LY = 145
	graphic.modesHandling(20)
	if graphic.getLCDC()&(1<<7) == 0 {
		fmt.Println("AAAA")
	}
	if graphic.getSTAT() != 0b00000001 {
		t.Errorf("Expected to be %08b, got %08b", 0b00000001, graphic.getSTAT())
	}

	graphic.cpu.Memory[0xFF40] = 0b10110110
	graphic.cpu.Memory[0xFF44] = 12 // LY = 12
	graphic.modesHandling(170)
	if graphic.getLCDC()&(1<<7) == 0 {
		fmt.Println("AAAA")
	}
	if graphic.getSTAT() != 0b00000011 {
		t.Errorf("Expected to be %08b, got %08b", 0b00000011, graphic.getSTAT())
	}

	graphic.cpu.Memory[0xFF40] = 0b10110110
	graphic.cpu.Memory[0xFF44] = 12 // LY = 12
	graphic.modesHandling(270)
	if graphic.getLCDC()&(1<<7) == 0 {
		fmt.Println("AAAA")
	}
	if graphic.getSTAT() != 0b00000000 {
		t.Errorf("Expected to be %08b, got %08b", 0b00000000, graphic.getSTAT())
	}
}

func TestMapAddress(t *testing.T) {
	//graphic := NewGraphics(NewCPU())
	cpu := NewCPU()
	cpu.Memory[0xFF40] = 0b10110110
	//drawTileMap(cpu, 10, 15)
	mapAddr := getTileMapAddr(cpu.Memory[0xFF40])
	baseAddr := getTileBaseAddr(cpu.Memory[0xFF40])
	if mapAddr != 0x9800 {
		t.Errorf("Expected to be 0x9C00, got %04x", mapAddr)
	}
	if baseAddr != 0x8000 {
		t.Errorf("Expected to be 0x8000, got %04x", baseAddr)
	}

	cpu.Memory[0xFF40] = 0b10101110
	mapAddr = getTileMapAddr(cpu.Memory[0xFF40])
	baseAddr = getTileBaseAddr(cpu.Memory[0xFF40])
	if mapAddr != 0x9C00 {
		t.Errorf("Expected to be 0x9C00, got %04x", mapAddr)
	}
	if baseAddr != 0x9000 {
		t.Errorf("Expected to be 0x8000, got %04x", baseAddr)
	}

	cpu.Memory[0xFF40] = 0b10111110
	mapAddr = getTileMapAddr(cpu.Memory[0xFF40])
	baseAddr = getTileBaseAddr(cpu.Memory[0xFF40])
	if mapAddr != 0x9C00 {
		t.Errorf("Expected to be 0x9C00, got %04x", mapAddr)
	}
	if baseAddr != 0x8000 {
		t.Errorf("Expected to be 0x8000, got %04x", baseAddr)
	}

	cpu.Memory[0xFF40] = 0b10100110
	mapAddr = getTileMapAddr(cpu.Memory[0xFF40])
	baseAddr = getTileBaseAddr(cpu.Memory[0xFF40])
	if mapAddr != 0x9800 {
		t.Errorf("Expected to be 0x9C00, got %04x", mapAddr)
	}
	if baseAddr != 0x9000 {
		t.Errorf("Expected to be 0x8000, got %04x", baseAddr)
	}
}
