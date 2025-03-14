package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const width = 160
const height = 144

const bgWidth = 256
const bgHeight = 256

const VRAM_START = 0x8000
const VRAM_END = 0x97FF
const VRAM_SIZE = VRAM_END - VRAM_START

var colors = []rl.Color{rl.White, rl.LightGray, rl.DarkGray, rl.Black}

type TilePixleID uint8

const (
	zero TilePixleID = iota
	one
	two
	three
)

//var tilePixleId uint8

// tile = array of 8 rows where a row is an array of 8 TileValues
type tile = [8][8]uint8

type Graphics struct {
	VRAM    [VRAM_SIZE]byte
	OAM     [160]byte //Object Attribute Memory
	LCDC    byte      // LCD control
	tileSet [384]tile

	LY   uint8 // FF44, value from 0-153
	LYC  uint8 // FF45
	STAT uint8 //FF41

	//SCY, SCX: Background viewport Y position, X position
	SCY uint8 // FF42
	SCX uint8 // FF43

	//Window Y position, X position plus 7
	WY uint8 // FF4A
	WX uint8 // FF4B

	BGP uint8 //FF47

	OBP0 uint8 //FF48
	OBP1 uint8 // FF49
}

// Priority: 0 = No, 1 = BG and Window colors 1–3 are drawn over this OBJ
// Y flip: 0 = Normal, 1 = Entire OBJ is vertically mirrored
// X flip: 0 = Normal, 1 = Entire OBJ is horizontally mirrored
// DMG palette [Non CGB Mode only]: 0 = OBP0, 1 = OBP1
// Bank [CGB Mode Only]: 0 = Fetch tile from VRAM bank 0, 1 = Fetch tile from VRAM bank 1
// CGB palette [CGB Mode Only]: Which of OBP0–7 to use
//type Attributes struct {
//	priority   bool
//	yFlip      bool
//	xFlip      bool
//	DMGPalette bool
//}

func (graphic *Graphics) readVRAM(address uint16) byte {
	if address >= VRAM_START && address <= VRAM_END {
		return graphic.VRAM[address-VRAM_START]
	}
	return 0
}

func (graphic *Graphics) writeVRAM(address uint16, value byte) {
	index := address - VRAM_START
	graphic.VRAM[index] = value
	if index >= 0x1800 {
		return
	}
	// tiles rows encoded in 2 bytes
	// first byte on even address
	normalizedIndex := index & 0xFFFE

	// get bytes for tile row
	byte1 := graphic.VRAM[normalizedIndex]
	byte2 := graphic.VRAM[normalizedIndex+1]

	tileIndex := index / 16      // each tile 16 bytes,  8 rows tall, 2 bytes => 16 bytes
	rowIndex := (index % 16) / 2 // 2 bytes per row

	for pixel := 0; pixel < 8; pixel++ {
		//pixel 0 - left most bit (bit 7)
		mask := 1 << (7 - pixel)
		lsb := byte1 & byte(mask) // first byte -> least significant bit
		msb := byte2 & byte(mask) // second byte -> most significant bit
		//var value byte
		var pixelValue TilePixleID
		if lsb != 0 && msb != 0 {
			//var tilepixelid = three
			pixelValue = three
		} else if lsb != 0 {
			pixelValue = one
		} else if msb != 0 {
			pixelValue = two
		} else {
			pixelValue = zero
		}
		graphic.tileSet[tileIndex][rowIndex][pixel] = uint8(pixelValue)
	}

}

func (graphic *Graphics) readOAM(address uint16) byte {
	if address >= 0xFE00 && address <= 0xFE9F {
		return graphic.OAM[address]
	}
	return 0
}

func (graphic *Graphics) writeOAM(address uint16, value byte) {
	//if address >= 0xFE00 && address <= 0xFE9F {
	graphic.OAM[address] = value
	//}
}

func (graphic *Graphics) spritesOAM() {
	spriteSize := 8

	//In 8x16 sprite mode, the least significant bit of the
	// sprite pattern number is ignored and treated as 0.

	if graphic.LCDC&(1<<2) != 0 {
		spriteSize = 16
	}

	//test sprite
	//for i := 0; i < 12; i++ {
	//	fmt.Printf("sprite %d: y: %d  x: %d  attrib: %d\n", i, graphic.OAM[i*4], graphic.OAM[i*4+1], graphic.OAM[i*4+2])
	//}

	// display up to 40 movable objects (or sprites)
	for i := 0; i < 40; i++ {
		// each sprite consists of 4 bytes
		spriteAddr := i * 4

		//Byte 0 — Y Position
		//Y = Object’s vertical position on the screen + 16
		y := int(graphic.OAM[spriteAddr]) - 16

		//Byte 1 — X Position
		//X = Object’s horizontal position on the screen + 8.
		x := int(graphic.OAM[spriteAddr+1]) - 8

		//Byte 2 — Tile Index
		tileIndex := graphic.OAM[spriteAddr+2]
		fmt.Printf("Tile index: %d\n", tileIndex)

		if spriteSize == 16 {
			tileIndex &= 0xFF // mask bit 0
		}
		//if tileIndex >= 384 {
		//	continue
		//}
		//Byte 3 — Attributes/Flags
		attributes := graphic.OAM[spriteAddr+3]

		//					7			6	  5			 4		     3		 2	1	0
		//Attributes	Priority	Y flip	X flip	 DMG palette 	Bank	CGB palette
		//priority := attributes & (1 << 7)
		yFlip := attributes & (1 << 6)
		xFlip := attributes & (1 << 5)
		//DMGPallete := attributes & (1 << 4)

		tileData := graphic.tileSet[tileIndex]
		fmt.Printf("Tile data: %d\n", tileData)

		for row := 0; row < 8; row++ {
			tileY := row
			if yFlip != 0 {
				tileY = spriteSize - 1 - row
			}
			for col := 0; col < 8; col++ {
				tileX := col
				if xFlip != 0 {
					tileX = 7 - col
				}
				pixelValue := tileData[tileY][tileX]
				fmt.Printf("Pixel value: %02X\n", pixelValue)
				if pixelValue == 0 {
					continue // transparent
				}

				screenX := x + col
				screenY := y + row
				// TODO: priority check

				// check bounds
				if screenX < 0 || screenY < 0 || screenX >= width || screenY >= height {
					continue
				}
				fmt.Printf("ScreenX: %X , ScreenY: %X\n", screenX, screenY)
				// draw pixels
				rl.DrawPixel(int32(screenX), int32(screenY), colors[pixelValue])
				fmt.Printf("Screen X: %d, Y: %d\n", screenX, screenY)
				//rl.DrawPixel(int32(screenX), int32(screenY), rl.Blue)

			}
		}

	}
}

func (graphic *Graphics) testPixelDrawing() {
	for y := 0; y < width; y++ {
		for x := 0; x < height; x++ {
			rl.DrawPixel(int32(x), int32(y), rl.Black)
		}
	}
}

//		7				  6					5					4
//LCD & PPU enable	Window tile map		Window enable	BG & Window tiles

//	3			2			1				0
//
// BG tile map	OBJ size	OBJ enable	BG & Window enable / priority
func (graphic *Graphics) lcdControlBits() (byte, byte, byte, byte, byte, byte, byte, byte) {
	LCDEnable := graphic.LCDC & (1 << 7)
	windowTileMapArea := graphic.LCDC & (1 << 6) // 0 = 9800–9BFF; 1 = 9C00–9FFF
	windowEnable := graphic.LCDC & (1 << 5)
	bgWinTileDataArea := graphic.LCDC & (1 << 4) //0 = 8800–97FF; 1 = 8000–8FFF
	bgTileDataArea := graphic.LCDC & (1 << 3)    //0 = 9800–9BFF; 1 = 9C00–9FFF
	objSize := graphic.LCDC & (1 << 2)           //0 = 8×8; 1 = 8×16
	objEnable := graphic.LCDC & (1 << 1)
	bgWinEnable := graphic.LCDC & (1 << 0)
	return LCDEnable, windowTileMapArea, windowEnable, bgWinTileDataArea, bgTileDataArea, objSize, objEnable, bgWinEnable
}

//	7			6					5				4
//			LYC int select	Mode 2 int select	Mode 1 int select
//	3						2		  1	 0
//
// Mode 0 int select		LYC == LY	PPU mode
func (graphic *Graphics) lcdStatusBits() (byte, byte, byte, byte, byte, byte) {
	LYCSelect := graphic.STAT & (1 << 6) //If set, selects the LYC == LY condition for the STAT interrupt
	mode2 := graphic.STAT & (1 << 5)     //f set, selects the Mode 2 condition for the STAT interrupt.
	mode1 := graphic.STAT & (1 << 4)     //If set, selects the Mode 1 condition for the STAT interrupt.
	mode0 := graphic.STAT & (1 << 3)     //If set, selects the Mode 0 condition for the STAT interrupt.
	LYCeqLY := graphic.STAT & (1 << 2)   //Set when LY contains the same value as LYC; it is constantly updated.
	PPUMode := graphic.STAT&(1<<1) | graphic.STAT&(1<<0)
	return LYCSelect, mode2, mode1, mode0, LYCeqLY, PPUMode
}

//
//func (graphics *Graphics) PPUcalcCoord() {
//	bottom := (graphics.SCY + 143) % 255
//	right := (graphics.SCX + 159) % 255
//}
