package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"sort"
)

const width = 160
const height = 144

const bgWidth = 256
const bgHeight = 256

const VRAM_START = 0x8000
const VRAM_END = 0x97FF
const VRAM_SIZE = VRAM_END - VRAM_START

const MODE_HBLANK = 0
const MODE_VBLANK = 1
const MODE_OAMSCAN = 2
const MODE_DRAWING = 3

const CYCLES_PER_LINE = 456

//// cycles per mode
//const CYCLES_MODE_2 = 80  //searchin objects
//const CYCLES_MODE_3 = 172 // drawing 172-289
//const CYCLES_MODE_0 = 87  // hblank 87-204

const SCANLINES_PER_FRAME = 144
const VBLANK_LINES = 10
const TOTAL_LINES = 154

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

type Sprite struct {
	x, y       int
	tileIndex  byte
	attributes byte
}
type Graphics struct {
	CPU     *CPU
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
	OBP1 uint8 // FF49\

	mode int

	cycle int
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

func (graphic *Graphics) getFromMemory(address uint16) uint8 {
	switch {
	case address >= VRAM_START && address <= VRAM_END:
		return graphic.readVRAM(address)
	case address >= 0xFE00 && address <= 0xFE9F:
		return graphic.readOAM(address)
	case address == 0xFF44:
		return graphic.LY
	case address == 0xFF45:
		return graphic.LYC
	case address == 0xFF40:
		// TODO: get for LCDC
		graphic.getLCDC()
	case address == 0xFF41:
		//TODO: dct for STAT
		graphic.getSTAT()
	case address == 0xFF42:
		return graphic.SCY
	case address == 0xFF43:
		return graphic.SCX
	case address == 0xFF4A:
		return graphic.WY
	case address == 0xFF4B:
		return graphic.WX
	case address == 0xFF47:
		return graphic.BGP
	case address == 0xFF48:
		return graphic.OBP0
	case address == 0xFF49:
		return graphic.OBP1

	default:
		fmt.Printf("GPU read unknown address")
	}
	return 0

}

func (graphic *Graphics) getSTAT() uint8 {
	var bit1, bit2, bit3, bit4, bit5 uint8
	LYCSelect, mode2, mode1, mode0, LYCeqLY, _ := graphic.lcdStatusBits()
	if LYCSelect != 0 {
		bit1 = 0b01000000
	} else {
		bit1 = 0
	}
	if mode2 != 0 {
		bit2 = 0b00100000
	} else {
		bit2 = 0
	}
	if mode1 != 0 {
		bit3 = 0b00010000
	} else {
		bit3 = 0
	}
	if mode0 != 0 {
		bit4 = 0b00001000
	} else {
		bit4 = 0
	}
	if LYCeqLY != 0 {
		bit5 = 0b00000100
	} else {
		bit5 = 0
	}
	return bit1 | bit2 | bit3 | bit4 | bit5 | uint8(graphic.mode)
}

func (graphic *Graphics) setSTAT(value uint8) {

	LYCSelect, mode2, mode1, mode0, LYCeqLY, _ := graphic.lcdStatusBits()
	LYCSelect |= value & 0b01000000
	mode2 |= value & 0b00100000
	mode1 |= value & 0b00010000
	mode0 |= value & 0b00001000
	LYCeqLY |= value & 0b00000100
}

func (graphic *Graphics) getLCDC() uint8 {
	LCDEnable, windowTileMapArea, windowEnable, bgWinTileDataArea, bgTileDataArea, objSize, objEnable, bgWinEnable := graphic.lcdControlBits()
	var bit1, bit2, bit3, bit4, bit5, bit6, bit7, bit8 uint8
	if LCDEnable != 0 {
		bit1 = 0b10000000

	} else {
		bit1 = 0
	}
	if windowTileMapArea != 0 {
		bit2 = 0b01000000
	} else {
		bit2 = 0
	}
	if windowEnable != 0 {
		bit3 = 0b00100000
	} else {
		bit3 = 0
	}
	if bgWinTileDataArea != 0 {
		bit4 = 0b00010000
	} else {
		bit4 = 0
	}
	if bgTileDataArea != 0 {
		bit5 = 0b00001000
	} else {
		bit5 = 0
	}
	if objSize != 0 {
		bit6 = 0b00000100
	} else {
		bit6 = 0
	}
	if objEnable != 0 {
		bit7 = 0b00000010
	} else {
		bit7 = 0
	}
	if bgWinEnable != 0 {
		bit8 = 0b00000001
	} else {
		bit8 = 0
	}
	return bit1 | bit2 | bit3 | bit4 | bit5 | bit6 | bit7 | bit8
}

func (graphic *Graphics) setLCDC(value uint8) {
	LCDEnable, windowTileMapArea, windowEnable, bgWinTileDataArea, bgTileDataArea, objSize, objEnable, bgWinEnable := graphic.lcdControlBits()
	LCDEnable |= value & 0b10000000
	windowTileMapArea |= value & 0b01000000
	windowEnable |= value & 0b00100000
	bgWinTileDataArea |= value & 0b00010000
	bgTileDataArea |= value & 0b00001000
	objSize |= value & 0b00000100
	objEnable |= value & 0b00000010
	bgWinEnable |= value & 0b00000001

}

func (graphic *Graphics) set(address uint16, value uint8) {
	switch {
	case address >= VRAM_START && address <= VRAM_END:
		graphic.writeVRAM(address, value)
	case address >= 0xFE00 && address <= 0xFE9F:
		graphic.writeOAM(address, value)
	case address == 0xFF44:
		graphic.LY = value
	case address == 0xFF45:
		graphic.LYC = value
	case address == 0xFF40:
		// TODO: set for LCDC
	case address == 0xFF41:
		//TODO: set stat
	case address == 0xFF42:
		graphic.SCY = value
	case address == 0xFF43:
		graphic.SCX = value
	case address == 0xFF4A:
		graphic.WY = value
	case address == 0xFF4B:
		graphic.WX = value
	case address == 0xFF47:
		graphic.BGP = value
	case address == 0xFF48:
		graphic.OBP0 = value
	case address == 0xFF49:
		graphic.OBP1 = value

	}

}

// sprites tiles
// OAM scan - mode 2
func (graphic *Graphics) spritesOAM() [width]uint8 {

	var spritePixels [width]uint8

	for i := range spritePixels {
		spritePixels[i] = 255 //default transparent
	}

	//// sprite disable
	//if graphic.LCDC&(1<<1) == 0 {
	//	return nil
	//}
	spriteSize := 8

	//In 8x16 sprite mode, the least significant bit of the
	// sprite pattern number is ignored and treated as 0.

	if graphic.LCDC&(1<<2) != 0 {
		spriteSize = 16
	}

	// 10 sprites visible at a time
	var visibleSprites []Sprite

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

		//check if sprite is onscanline
		if int(graphic.LY)+16 < y || int(graphic.LY)+16 > y+spriteSize-1 {
			continue
		}

		visibleSprites = append(visibleSprites, Sprite{x, y, tileIndex, attributes})

		// TODO: sort sprites by priority
		// order by X (smaller X->bigger prority)
		//if x the same => by OAM location order
		sort.Slice(visibleSprites, func(i, j int) bool {
			return visibleSprites[i].x < visibleSprites[j].x
		})

		// render 10 sprites
		if len(visibleSprites) > 10 {
			visibleSprites = visibleSprites[:10]
		}

		bgPixels := graphic.getBackground()

		for _, sprite := range visibleSprites {

			//					7			6	  5			 4		     3		 2	1	0
			//Attributes	Priority	Y flip	X flip	 DMG palette 	Bank	CGB palette
			yFlip := sprite.attributes & (1 << 6)
			xFlip := sprite.attributes & (1 << 5)

			tileData := graphic.tileSet[sprite.tileIndex]
			fmt.Printf("Tile data: %d\n", tileData)

			priority := attributes & (1 << 7)

			//DMGPallete := attributes & (1 << 4)

			tileY := graphic.LY - uint8(sprite.y)
			if yFlip != 0 {
				tileY = uint8(spriteSize) - 1 - graphic.LY - uint8(sprite.y)
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
				screenX := sprite.x + col
				screenY := uint8(sprite.y) + graphic.LY

				// check bounds
				if screenX < 0 || screenY < 0 || screenX >= width || screenY >= height {
					continue
				}

				//TODO: priority
				bgPixelValue := bgPixels[screenY][screenX]
				if priority != 0 && bgPixelValue != 0 {
					continue
				}

				if spritePixels[screenX] == 255 {
					spritePixels[screenX] = pixelValue
				}
				//rl.DrawPixel(int32(screenX), int32(screenY), colors[pixelValue])

			}

		}

	}
	return spritePixels
}

// background tiles
func (graphic *Graphics) getBackground() [height][width]uint8 {
	// $9800-$9BFF and $9C00-$9FFF
	var bgPixels [height][width]uint8
	for screenY := 0; screenY < height; screenY++ {
		for screenX := 0; screenX < width; screenX++ {
			bgY := (screenY + int(graphic.LY) + int(graphic.SCY)) & 0xFF

			bgX := (screenX + int(graphic.SCX)) & 0x1F
			bgMapAddr := uint16(0x9800) //default
			if graphic.LCDC&(1<<3) != 0 {
				bgMapAddr = uint16(0x9C00)
			}
			tileIdxAddr := bgMapAddr + uint16((bgY/8)*32+(bgX/8))

			if tileIdxAddr < VRAM_START || tileIdxAddr > VRAM_END {
				continue
			}
			//tileIndex := graphic.VRAM[tileIdxAddr]
			tileIndex := graphic.readVRAM(tileIdxAddr)
			var tileNumber int
			//$8000-$97FFbgY
			if graphic.LCDC&(1<<4) != 0 {
				tileNumber = int(tileIndex)
			} else {
				tileNumber = int(int8(tileIndex)) // make it signed from unsigned
			}

			tileData := graphic.tileSet[tileNumber]
			tileY := bgY % 8
			tileX := bgX % 8
			bgPixels[screenY][screenX] = tileData[tileY][tileX]

		}
	}
	return bgPixels

}

// window tiles
func (graphic *Graphics) getWindow() [height][width]uint8 {
	var window [height][width]uint8

	// window display disabled
	if graphic.LCDC&(1<<5) == 0 {
		return window
	}

	for screenY := 0; screenY < height; screenY++ {
		if screenY < int(graphic.WY) {
			continue
		}
		for screenX := 0; screenX < width; screenX++ {
			//if graphic.WX < 7 {
			//	continue
			//}
			if screenX < int(graphic.WX)-7 {
				continue
			}
			winY := screenY - int(graphic.WY)

			winX := screenX - (int(graphic.WX) - 7)
			winMapAddr := uint16(0x9800)
			if graphic.LCDC&(1<<6) != 0 {
				winMapAddr = uint16(0x9C00)
			}
			tileIdxAddr := winMapAddr + uint16((winY/8)*32+(winX/8))
			//tileIndex := graphic.VRAM[tileIdxAddr]
			tileIndex := graphic.readVRAM(tileIdxAddr)

			var tileNumber int

			if graphic.LCDC&(1<<4) != 0 {
				tileNumber = int(tileIndex)
			} else {
				tileNumber = int(int8(tileIndex))
			}

			tileData := graphic.tileSet[tileNumber]
			tileY := winY % 8
			tileX := winX % 8
			window[screenY][screenX] = tileData[tileY][tileX]

		}
	}
	return window
}

func (graphic *Graphics) renderScanline() {
	if int(graphic.LY) >= height {
		return
	}

	bgPixels := graphic.getBackground()
	winPixels := graphic.getWindow()
	spritePixels := graphic.spritesOAM()

	for screenX := 0; screenX < width; screenX++ {
		pixel := bgPixels[graphic.LY][screenX]
		if graphic.LCDC&(1<<5) != 0 && winPixels[graphic.LY][screenX] != 0 {
			pixel = winPixels[graphic.LY][screenX]

		}
		if graphic.LCDC&(1<<1) != 0 && spritePixels[screenX] != 255 {
			pixel = spritePixels[screenX]
		}
		rl.DrawPixel(int32(screenX), int32(graphic.LY), colors[pixel])
	}
}

func (graphic *Graphics) modesHandeling(tCycles int) {
	cpu := CPU{}
	if graphic.LCDC&(1<<7) == 0 {
		return //LCD disabled
	}
	graphic.cycle += tCycles
	LYCSelect, mode2, mode1, mode0, _, _ := graphic.lcdStatusBits()
	// at a time total amount of 80 T-Cycles
	for graphic.cycle >= 80 {
		switch graphic.mode {
		case MODE_OAMSCAN:
			//mode 2 80 cycles
			if graphic.cycle >= 80 {
				graphic.mode = MODE_DRAWING
				graphic.cycle -= 80

			}
			if mode2 != 0 {
				//TODO: handle m2 interrupt
				cpu.IF |= 1 << 1

			}

		case MODE_DRAWING:
			//mode 3 172-289 cycles
			if graphic.cycle >= 172 {
				graphic.mode = MODE_HBLANK
				graphic.cycle -= 172
				graphic.renderScanline()

			}

		case MODE_HBLANK:
			// mide 0 87-204cycles
			if graphic.cycle >= 204 {
				graphic.cycle -= 204
				graphic.LY++

				if graphic.LY == SCANLINES_PER_FRAME {
					//enter VBlank
					graphic.mode = MODE_VBLANK
					// TODO: hanfle vblank interrupt
					cpu.IF |= 1 << 0
				} else {
					graphic.mode = MODE_OAMSCAN // start next scanline
				}
			}
			if mode0 != 0 {
				//TODO: handle m0 interrupt
				cpu.IF |= 1 << 1
			}

		case MODE_VBLANK:
			//mode 1 4560 cycles
			if graphic.cycle >= CYCLES_PER_LINE {
				graphic.LY++
				graphic.cycle -= CYCLES_PER_LINE
				if graphic.LY > TOTAL_LINES-1 { //end of vblank
					graphic.LY = 0
					graphic.mode = MODE_OAMSCAN // start new scanline

				}
			}
			if mode1 != 0 {
				//TODO: handle m1 interrupt
				cpu.IF |= 1 << 1
			}
		}
		if graphic.LY == graphic.LYC {
			//set stat bit 2
			graphic.STAT |= 1 << 2
			if LYCSelect != 0 {
				cpu.IF |= 1 << 1
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
