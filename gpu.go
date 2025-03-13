package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const width = 160
const height = 144

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
	spriteSize := 8 //TODO based on lcdc to be 8 or 16

	//if graphic.LCDC&(1<<2) != 0 {
	//	spriteSize = 16
	//}

	//test sprite
	//for i := 0; i < 12; i++ {
	//	fmt.Printf("sprite %d: y: %d  x: %d  attrib: %d\n", i, graphic.OAM[i*4], graphic.OAM[i*4+1], graphic.OAM[i*4+2])
	//}

	// display up to 40 movable objects (or sprites)
	for i := 0; i < 12; i++ {
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
				tileY = spriteSize - row
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
