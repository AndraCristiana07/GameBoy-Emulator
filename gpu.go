package main

import (
	"fmt"
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const width = 160
const height = 144

const bgWidth = 256
const bgHeight = 256

const VRAM_START = 0x8000
const VRAM_END = 0x9FFF
const VRAM_SIZE = VRAM_END - VRAM_START

const OAM_START = 0xFE00
const OAM_END = 0xFE9F

const OAM_SIZE = OAM_END - OAM_START + 1
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

var colors = [4]rl.Color{rl.White, rl.LightGray, rl.DarkGray, rl.Black}

type TilePixelID uint8

const (
	zero TilePixelID = iota
	one
	two
	three
)

// tile = array of 8 rows where a row is an array of 8 TileValues
type tile = [8][8]uint8

type Sprite struct {
	x, y       int
	tileIndex  byte
	attributes byte
	OAMOrder   int
}
type Graphics struct {
	cpu *CPU
	//pixelBuffer [height][width]uint8
	pixelBuffer [height][width]uint8

	drawnLine bool
	cycle     int
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

func NewGraphics(cpu *CPU) *Graphics {
	graphic := &Graphics{cpu: cpu}
	// fmt.Println("New graphics called. address", &graphic)
	graphic.cycle = 0
	return graphic
}

func (graphic *Graphics) getLCDC() uint8 {
	return graphic.cpu.Memory[0xFF40]
}

func (graphic *Graphics) getLY() uint8 {
	return graphic.cpu.Memory[0xFF44]
}

func (graphic *Graphics) getSTAT() uint8 {
	return graphic.cpu.Memory[0xFF41]
}

func (graphic *Graphics) getSCY() uint8 {
	return graphic.cpu.Memory[0xFF42]
}
func (graphic *Graphics) getSCX() uint8 {
	return graphic.cpu.Memory[0xFF43]
}

func (graphic *Graphics) getWY() uint8 {
	return graphic.cpu.Memory[0xFF4A]
}

func (graphic *Graphics) getWX() uint8 {
	return graphic.cpu.Memory[0xFF4B]
}

func (graphic *Graphics) readVRAM(address uint16) byte {
	if address >= VRAM_START && address <= VRAM_END {
		return graphic.cpu.Memory[address-VRAM_START]
	}
	panic("nu ar trebui sa se intample asa ceva")
	//return 0
}

func (cpu *CPU) dmaTransfer(value uint8) {
	println("dmaTransfer")
	//upper := uint16(value) / 100
	upper := uint16(value) << 8
	//dest := uint16(OAM_START)
	for i := 0; i < OAM_SIZE; i++ {
		fmt.Printf("upper is 0x%02X\n", upper)
		fmt.Printf("upper with i is 0x%02X\n", upper+uint16(i))

		fmt.Printf("OAM at 0x%02X will be: 0x%02X\n", OAM_START+uint16(i), cpu.Memory[upper+uint16(i)])
		cpu.Memory[OAM_START+uint16(i)] = cpu.Memory[upper+uint16(i)]
		// fmt.Printf("DMA Transfer ->  %04X - %02X\n ", i, graphic.cpu.Memory[i])
	}
	println("done")
}

func (cpu *CPU) readTileData(address uint16) [8][8]uint8 {
	var tile [8][8]uint8

	if address < VRAM_START || address > 0x97FF {
		// fmt.Printf("read tile data -> invalid address 0x%04X\n", address)
		return tile
	}
	for row := 0; row < 8; row++ {
		low := cpu.memoryRead(address + uint16(row*2))
		high := cpu.memoryRead(address + uint16(row*2) + 1)

		for col := 0; col < 8; col++ {
			bit := 7 - col
			msbLow := (low >> bit) & 1
			msbHigh := (high >> bit) & 1
			colorID := (msbHigh << 1) | msbLow
			tile[row][col] = uint8(colorID)
		}
	}
	return tile
}

func (graphic *Graphics) getTilePixel(address uint16, x int, y int) uint8 {
	//each tile taking 16 bytes, 2 per row
	// y row idx
	// x pixel idx
	rowAddr := address + uint16(y*2)
	bit := 7 - x
	low := graphic.cpu.memoryRead(uint16(rowAddr))
	high := graphic.cpu.memoryRead(uint16(rowAddr) + 1)
	//msbLow := low & (1 << bit)
	//msbHigh := high & (1 << bit)
	msbLow := (low >> bit) & 1
	msbHigh := (high >> bit) & 1
	colorID := (msbHigh << 1) | msbLow
	return colorID

}

// sprites tiles
// OAM scan - mode 2
func (graphic *Graphics) spritesOAM() [height][width]uint8 {

	// fmt.Println("In sprites render")
	var spritePixels [height][width]uint8

	for y := range spritePixels {
		for x := range spritePixels[y] {
			spritePixels[y][x] = 255
		}
	}
	//obp0 := graphic.cpu.Memory[0xFF48]
	//obp1 := graphic.cpu.Memory[0xFF49]
	spriteSize := 8

	//In 8x16 sprite mode, the least significant bit of the
	// sprite pattern number is ignored and treated as 0.

	if graphic.getLCDC()&(1<<2) != 0 {
		spriteSize = 16
	}

	// 10 sprites visible at a time
	var visibleSprites []Sprite

	// display up to 40 movable objects (or sprites)
	for i := 0; i < 40; i++ {
		// each sprite consists of 4 bytes
		spriteAddr := OAM_START + i*4
		// fmt.Printf("spriteAddr is 0x%04X\n", spriteAddr)

		//Byte 0 — Y Position
		//Y = Object’s vertical position on the screen + 16
		//y := int(graphic.OAM[spriteAddr]) - 16
		y := int(graphic.cpu.memoryRead(uint16(spriteAddr))) - 16

		//Byte 1 — X Position
		//X = Object’s horizontal position on the screen + 8.
		x := int(graphic.cpu.memoryRead(uint16(spriteAddr+1))) - 8

		//Byte 2 — Tile Index
		//// fmt.Printf("Value at address 0x%04X is %d\n", spriteAddr+2, graphic.cpu.memoryRead(uint16(spriteAddr+2)))

		tileIndex := graphic.cpu.memoryRead(uint16(spriteAddr + 2))

		if spriteSize == 16 {
			tileIndex &= 0xFE // mask bit 0
		}

		//Byte 3 — Attributes/Flags
		//attributes := graphic.OAM[spriteAddr+3]
		attributes := graphic.cpu.memoryRead(uint16(spriteAddr + 3))
		//// fmt.Printf("value for spriteaddr: %d\n", graphic.cpu.memoryRead(uint16(spriteAddr)))
		//// fmt.Printf("Sprite %d -> OAM Addr: 0x%04X | X: %d Y: %d | Tile Index: %d | Attributes: 0b%08b\n", i, spriteAddr, x, y, tileIndex, attributes)

		//check if sprite is onscanline
		if int(graphic.getLY()) < y || int(graphic.getLY()) >= y+spriteSize {
			continue
		}

		visibleSprites = append(visibleSprites, Sprite{x, y, tileIndex, attributes, i})
	}

	// order by X (smaller X->bigger prority)
	//if x the same => by OAM location order
	sort.Slice(visibleSprites, func(i, j int) bool {
		if visibleSprites[i].x == visibleSprites[j].x {
			return visibleSprites[i].OAMOrder < visibleSprites[j].OAMOrder
		}
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

		priority := sprite.attributes & (1 << 7)

		//DMGPallete := sprite.attributes & (1 << 4)

		//palette := obp0
		//if DMGPallete != 0 {
		//	palette = obp1
		//}

		tileY := int(graphic.getLY()) - sprite.y
		if yFlip != 0 {
			tileY = spriteSize - 1 - tileY

		}
		if tileY < 0 || tileY >= 8 {
			//// fmt.Printf("tileY %d is out of bounds (LY=%d, spriteY=%dm spriteSize=%d) \n", tileY, graphic.getLY(), sprite.y, spriteSize)
			continue
		}
		tileAddress := 0x8000 + (int(sprite.tileIndex) * 16)
		if tileAddress < VRAM_START || tileAddress > 0x97FF {
			// fmt.Printf("tileAddress %04X is invalid\n", tileAddress)
			continue
		}
		tileData := graphic.cpu.readTileData(uint16(tileAddress))
		// fmt.Printf("Tile data: %d\n", tileData)

		for col := 0; col < 8; col++ {
			tileX := col
			if xFlip != 0 {
				tileX = 7 - col
			}
			// fmt.Printf("tiledata: %d, tileX: %d, tileY: %d\n", tileData, tileX, tileY)
			if tileX < 0 || tileX >= 8 || tileY < 0 || tileY >= 8 {
				// fmt.Printf("Out of bounds! tileX=%d, tileY=%d  (oamOrder=%d, tileIdx=%d, sprtieX=%d, spriteY=%d)\n", tileX, tileY, sprite.OAMOrder, sprite.tileIndex, sprite.x, sprite.y)
			}
			// fmt.Printf("About to fetch pixel for sprite: tileX=%d, tileY=%d, tileIndex=%d\n", tileX, tileY, sprite.tileIndex)
			pixelValue := tileData[tileY][tileX]

			// fmt.Printf("Pixel value: %02X\n", pixelValue)
			if pixelValue == 0 {
				// fmt.Printf("skip transparent")
				continue // transparent
			}
			screenX := sprite.x + col
			screenY := graphic.getLY()

			// check bounds
			if screenX < 0 || screenY < 0 || screenX >= width || screenY >= height {
				// fmt.Printf("skip pixel out of bounds (x:%d, y:%d)\n", screenX, screenY)
				continue
			}

			bgPixelValue := bgPixels[screenY][screenX]
			if priority != 0 && bgPixelValue != 0 {
				// fmt.Printf("skip pixel -> priority bgPixelValue: %02X\n", bgPixelValue)
				continue
			}

			if spritePixels[graphic.getLY()][screenX] == 255 {
				// fmt.Printf("Drawing sprite pixel at ScreenX=%d ScreenY=%d, value=%d\n", screenX, screenY, pixelValue)
				//shade := (palette >> (pixelValue * 2)) & 0x03
				//// fmt.Println("Sprite Palette", palette, "shade", shade)

				//spritePixels[graphic.getLY()][screenX] = shade
				spritePixels[graphic.getLY()][screenX] = pixelValue

			} else {
				// fmt.Println("Skip pixel. Pixel already occupied")
			}

		}

	}

	return spritePixels
}

// background tiles
func (graphic *Graphics) getBackground() [height][width]uint8 {
	// $9800-$9BFF and $9C00-$9FFF
	var bgPixels [height][width]uint8
	//palette := graphic.cpu.Memory[0xFF47]

	for screenY := 0; screenY < height; screenY++ {
		//bgY := (screenY + int(graphic.getLY()) + int(graphic.getSCY())) & 0xFF
		bgY := (screenY + int(graphic.getSCY())) & 0xFF
		for screenX := 0; screenX < width; screenX++ {

			//bgX := (screenX + int(graphic.getSCX())) & 0x1F
			bgX := (screenX + int(graphic.getSCX())) & 0xFF

			bgMapAddr := uint16(0x9800) //default
			if graphic.getLCDC()&(1<<3) != 0 {
				bgMapAddr = uint16(0x9C00)
			}
			tileIdxAddr := bgMapAddr + uint16((bgY/8)*32+(bgX/8))

			if tileIdxAddr < VRAM_START || tileIdxAddr > VRAM_END {
				continue
			}
			tileIndex := graphic.readVRAM(tileIdxAddr)
			var tileNumber int
			//$8000-$97FFbgY
			if graphic.getLCDC()&(1<<4) != 0 {
				tileNumber = int(tileIndex)
			} else {
				tileNumber = int(int8(tileIndex)) // make it signed from unsigned
			}

			if tileNumber < 0 || tileNumber >= 384 {
				// fmt.Printf("Invalid tile index: %d tileIdxAddr: 0x%04X\n", tileNumber, tileIdxAddr)
				continue
			}

			//var tileAddr uint16
			//if graphic.getLCDC()&(1<<4) != 0 {
			//	tileAddr = 0x8000 + uint16(tileNumber*16)
			//} else {
			//	tileAddr = 0x9000 + uint16(tileNumber*16)
			//}
			//
			//// fmt.Println("Tile address background", tileAddr)
			tileY := bgY % 8
			tileX := bgX % 8
			//colorID := graphic.getTilePixel(uint16(tileNumber), tileX, tileY)
			colorID := graphic.getTilePixel(uint16(tileNumber*16), tileX, tileY)

			//colorID := graphic.getTilePixel(tileAddr, tileX, tileY)

			//BGP Palette
			//shade := (palette >> (colorID * 2)) & 0x03
			//// fmt.Println("Background palette:", palette, "shade:", shade)
			//bgPixels[screenY][screenX] = shade
			bgPixels[screenY][screenX] = colorID
		}
	}
	return bgPixels

}

// window tiles
func (graphic *Graphics) getWindow() [height][width]uint8 {
	var window [height][width]uint8

	// window display disabled
	if graphic.getLCDC()&(1<<5) == 0 {
		return window
	}

	for screenY := 0; screenY < height; screenY++ {
		if screenY < int(graphic.getWY()) {
			continue
		}
		for screenX := 0; screenX < width; screenX++ {

			if screenX < int(graphic.getWX())-7 {
				continue
			}
			winY := screenY - int(graphic.getWY())

			winX := screenX - (int(graphic.getWX()) - 7)
			winMapAddr := uint16(0x9800)
			if graphic.getLCDC()&(1<<6) != 0 {
				winMapAddr = uint16(0x9C00)
			}
			tileIdxAddr := winMapAddr + uint16((winY/8)*32+(winX/8))

			if tileIdxAddr < VRAM_START || tileIdxAddr > VRAM_END {
				continue
			}
			tileIndex := graphic.readVRAM(tileIdxAddr)

			var tileNumber int

			if graphic.getLCDC()&(1<<4) != 0 {
				tileNumber = int(tileIndex)
			} else {
				tileNumber = int(int8(tileIndex))
			}

			//var tileAddr uint16
			//if graphic.getLCDC()&(1<<4) != 0 {
			//	tileAddr = 0x8000 + uint16(tileNumber*16)
			//} else {
			//	tileAddr = 0x9000 + uint16(tileNumber*16)
			//}
			//
			//// fmt.Println("Tile address window", tileAddr)

			tileY := winY % 8
			tileX := winX % 8
			//colorID := graphic.getTilePixel(tileNumber, tileX, tileY)
			//colorID := graphic.getTilePixel(uint16(tileNumber), tileX, tileY)

			colorID := graphic.getTilePixel(uint16(tileNumber*16), tileX, tileY)
			//BGP Palette
			//palette := graphic.cpu.Memory[0xFF47]
			//shade := (palette >> (colorID * 2)) & 0x03
			//// fmt.Println("Window pallete", palette, "shade", shade)

			//window[screenY][screenX] = shade
			window[screenY][screenX] = colorID

		}
	}
	return window
}

func (graphic *Graphics) renderScanline() {
	// fmt.Println("In render scanline")
	if int(graphic.getLY()) >= height {
		return
	}
	ly := graphic.getLY()

	bgPixels := graphic.getBackground()
	winPixels := graphic.getWindow()
	spritePixels := graphic.spritesOAM()

	//// fmt.Printf("bgPixels: %d, winPixels: %d, spritePixels: %d\n", bgPixels, winPixels, spritePixels)
	// fmt.Printf("LCDC: 0b%08b and sprite flag %t\n", graphic.getLCDC(), (graphic.getLCDC()&(1<<1)) != 0)
	for screenX := 0; screenX < width; screenX++ {
		// fmt.Printf("bgPixels: %d, winPixels: %d, spritePixels: %d at LY:%d and screenX:%d\n", bgPixels[graphic.getLY()][screenX], winPixels[graphic.getLY()][screenX], spritePixels[graphic.getLY()][screenX], graphic.getLY(), screenX)
		var pixel uint8 = 0 // default white

		// starting with background
		pixel = bgPixels[ly][screenX]

		// window overrides background
		if graphic.getLCDC()&(1<<5) != 0 && winPixels[ly][screenX] != 0 {
			pixel = winPixels[ly][screenX]

		}

		//sprites override, unless transparent
		//if graphic.getLCDC()&(1<<1) != 0 && spritePixels[graphic.getLY()][screenX] != 255 {
		//	// fmt.Printf("Sprite pixel override at X: %d, value:%d\n", screenX, spritePixels[screenX])
		//	pixel = spritePixels[graphic.getLY()][screenX]
		//}
		if graphic.getLCDC()&(1<<1) != 0 {
			spritePixel := spritePixels[ly][screenX]
			if spritePixel != 255 && spritePixel != 0 {
				pixel = spritePixel
				// fmt.Printf("Sprite pixel override at X: %d, value:%d\n", screenX, spritePixels[screenX])

			}
		}
		//color := colors[pixel]
		graphic.pixelBuffer[ly][screenX] = pixel
		// fmt.Printf("pixel: %d, color of pixel: %d\n", pixel, color)

		//rl.DrawPixel(int32(screenX), int32(graphic.getLY()), color)
		//scale := 2
		//rl.DrawRectangle(int32(screenX*scale), int32(int(ly)*scale), int32(scale), int32(scale), color)

	}
}

func (graphic *Graphics) drawScreen() {
	scale := 2
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			colorIdx := graphic.pixelBuffer[y][x]
			color := colors[colorIdx]
			rl.DrawRectangle(int32(x*scale), int32(y*scale), int32(scale), int32(scale), color)
		}
	}
}

func (graphic *Graphics) modesHandeling(tCycles int) {
	// fmt.Println("mode handling")
	// fmt.Printf("LCDC: 0b%08b\n", graphic.getLCDC())
	if graphic.cpu == nil {
		// fmt.Println("Error: graphics.cpu is nil")
		return
	}
	if graphic.getLCDC()&(1<<7) == 0 {
		// fmt.Println("LCD Disabled!")
		graphic.setMode(MODE_HBLANK)
		graphic.cycle = 0
		return
	}
	// fmt.Println("Before : Entering with tCycles: ", tCycles, " current cycle: ", graphic.cycle, " mode: ", PPUmode)
	graphic.cycle += tCycles
	//// fmt.Println("after adding cycle -> total: ", graphic.cycle, " added cycles: ", tCycles)

	if int(graphic.getLY()) >= SCANLINES_PER_FRAME {
		fmt.Println("Entering VBLANK")
		graphic.setMode(MODE_VBLANK)
		if graphic.getLY() == SCANLINES_PER_FRAME {
			// fmt.Println("Entering Vblank")
			IF := graphic.cpu.Memory[0xFF0F] | 1<<0
			graphic.cpu.memoryWrite(0xFF0F, IF)
		}
	} else {
		//if graphic.cycle >= 456-80 {
		if graphic.cycle < 80 {
			if graphic.cpu.Memory[0xFF41] != MODE_OAMSCAN {
				fmt.Println("Entering OAMSCAN")
			}
			//fmt.Println("Entering OAMSCAN")
			graphic.setMode(MODE_OAMSCAN)

			//} else if graphic.cycle >= 456-80-172 {
		} else if graphic.cycle < 80+172 {
			if graphic.cpu.Memory[0xFF41] != MODE_DRAWING {
				fmt.Println("Entering drawing mode")
			}
			//fmt.Println("Entering Drawing")
			graphic.setMode(MODE_DRAWING)

			if !graphic.drawnLine {
				graphic.renderScanline()
				graphic.drawnLine = true
			}
		} else {
			if graphic.cpu.Memory[0xFF41] != MODE_HBLANK {
				fmt.Println("Entering Hblank")
			}
			//fmt.Println("Entering HBLANK")
			graphic.setMode(MODE_HBLANK)

		}

	}
	if graphic.cycle >= CYCLES_PER_LINE {
		graphic.cycle -= CYCLES_PER_LINE
		graphic.drawnLine = false

		fmt.Printf("ly inc to : %d\n", graphic.getLY()+1)
		graphic.cpu.Memory[0xFF44]++

		if graphic.getLY() >= TOTAL_LINES-1 {
			graphic.cpu.Memory[0xFF44] = 0

			graphic.cycle = 0

			fmt.Println("reset ly, start new frame")
		}
	}
	//graphic.cpu.Memory[0xFF44] = (graphic.cpu.Memory[0xFF44] + 1) % 154
	// fmt.Printf("LY: %d Cycles: %d mode: %d  \n", graphic.getLY(), graphic.cycle, PPUmode)
}

//	7				  6					5					4
//
// LCD & PPU enable	Window tile map		Window enable	BG & Window tiles
//
//	3			2			1				0
//
// BG tile map	OBJ size	OBJ enable	BG & Window enable / priority
func (graphic *Graphics) lcdControlBits() (byte, byte, byte, byte, byte, byte, byte, byte) {
	LCDEnable := graphic.getLCDC() & (1 << 7)
	windowTileMapArea := graphic.getLCDC() & (1 << 6) // 0 = 9800–9BFF; 1 = 9C00–9FFF
	windowEnable := graphic.getLCDC() & (1 << 5)
	bgWinTileDataArea := graphic.getLCDC() & (1 << 4) //0 = 8800–97FF; 1 = 8000–8FFF
	bgTileDataArea := graphic.getLCDC() & (1 << 3)    //0 = 9800–9BFF; 1 = 9C00–9FFF
	objSize := graphic.getLCDC() & (1 << 2)           //0 = 8×8; 1 = 8×16
	objEnable := graphic.getLCDC() & (1 << 1)
	bgWinEnable := graphic.getLCDC() & (1 << 0)
	return LCDEnable, windowTileMapArea, windowEnable, bgWinTileDataArea, bgTileDataArea, objSize, objEnable, bgWinEnable
}

//	7			6					5				4
//			LYC int select	Mode 2 int select	Mode 1 int select
//	3						2		  1	 0
//
// Mode 0 int select		LYC == LY	PPU mode
func (graphic *Graphics) lcdStatusBits() (byte, byte, byte, byte, byte, byte) {
	LYCSelect := graphic.getLY() & (1 << 6) //If set, selects the LYC == LY condition for the STAT interrupt
	mode2 := graphic.getLY() & (1 << 5)     //f set, selects the Mode 2 condition for the STAT interrupt.
	mode1 := graphic.getLY() & (1 << 4)     //If set, selects the Mode 1 condition for the STAT interrupt.
	mode0 := graphic.getLY() & (1 << 3)     //If set, selects the Mode 0 condition for the STAT interrupt.
	LYCeqLY := graphic.getLY() & (1 << 2)   //Set when LY contains the same value as LYC; it is constantly updated.
	PPUMode := graphic.getLY()&(1<<1) | graphic.getLY()&(1<<0)
	return LYCSelect, mode2, mode1, mode0, LYCeqLY, PPUMode
}

func (graphic *Graphics) setMode(mode uint8) {
	stat := graphic.getSTAT()
	//		   clear bit         mask last 2 bits
	stat = (stat &^ 0b00000011) | (mode & 0b00000011)
	graphic.cpu.memoryWrite(0xFF41, stat)

}

func drawTiles(vram []uint8, startX int, startY int) {
	tileWidth := 8
	tileHeight := 8
	tileBytes := 16
	scale := 2
	tilesPerRow := 16

	totalTiles := len(vram) / tileBytes
	for tileIdx := 0; tileIdx < totalTiles; tileIdx++ {
		tileX := tileIdx % tilesPerRow
		tileY := tileIdx / tilesPerRow
		//tileX := (tileIdx % tilesPerRow) * tileWidth * scale
		//tileY := (tileIdx / tilesPerRow) * tileHeight * scale

		screenBaseX := startX + tileX*tileWidth*scale
		screenBaseY := startY + tileY*tileHeight*scale

		offset := tileIdx * tileBytes
		if offset+16 > len(vram) {
			continue
		}
		tile := vram[offset : offset+16]

		for y := 0; y < tileHeight; y++ {
			low := tile[y*2]
			high := tile[y*2+1]
			//low := vram[offset+y*2]
			//high := vram[offset+y*2+1]

			for x := 0; x < tileWidth; x++ {
				bit := 7 - x
				msbLow := (low >> bit) & 1
				msbHigh := (high >> bit) & 1
				colorId := (msbHigh << 1) | msbLow

				color := colors[colorId]

				// draw scaled pixel
				screenX := screenBaseX + x*scale
				screenY := screenBaseY + y*scale
				//screenX := startX + tileX + (x * scale)
				//screenY := startY + tileY + (y * scale)
				rl.DrawRectangle(int32(screenX), int32(screenY), int32(scale), int32(scale), color)
			}
		}
	}
}

func drawSprites(oam []uint8, vram []uint8, startX int, startY int) {
	spriteCount := 40
	spriteWidth := 8
	spriteHeight := 8
	scale := 2

	for i := 0; i < spriteCount; i++ {
		spriteAddr := i * 4
		y := int(oam[spriteAddr]) - 16
		x := int(oam[spriteAddr+1]) - 8
		tileIndex := int(oam[spriteAddr+2])
		//attributes := int(oam[spriteAddr+3])

		if tileIndex*16+16 > len(vram) {
			fmt.Printf("Draw sprites - INVALID")
			continue
		}
		if tileIndex < 0 || tileIndex >= 384 {
			fmt.Printf("Draw sprites - out of bounds")
			continue
		}

		tileStart := tileIndex * 16
		if tileStart+16 > len(vram) {
			fmt.Printf("Draw sprites - out of vrqm bounds")
			continue
		}

		tileData := vram[tileStart : tileStart+16]

		for ty := 0; ty < spriteHeight; ty++ {
			low := tileData[ty*2]
			high := tileData[ty*2+1]
			for tx := 0; tx < spriteWidth; tx++ {
				bit := 7 - tx
				msbLow := (low >> bit) & 1
				msbHigh := (high >> bit) & 1
				colorId := (msbHigh << 1) | msbLow
				if colorId == 0 {
					fmt.Printf("transparent")
					continue
				}
				//offset := 500
				//rl.DrawRectangle(int32(startX+x+tx*scale), int32(startY+y+ty*scale), int32(scale), int32(scale), colors[colorId])

				rl.DrawRectangle(int32((startX+x+tx)*scale), int32((startY+y+ty)*scale), int32(scale), int32(scale), rl.Red)

				//rl.DrawRectangle(int32((x+tx)*scale), int32((y+ty)*scale+offset), int32(scale), int32(scale), rl.Red)
				//rl.DrawPixel(int32(startX+x+tx), int32(startY+y+ty), colors[colorId])
				//rl.DrawText(fmt.Sprintf("%02d", i), int32((startX+x+tx)*scale), int32((startY+y+ty)*scale), 10, rl.Red)
			}
		}

	}
}

func drawTileMap(vram []uint8, tileMapAddr uint16, tileBaseAddr uint16, startX int, startY int) {
	tileWidth := 8
	tileHeight := 8
	mapWidth := 32
	mapHeight := 32
	scale := 2

	//tileMap := vram[tileMapAddr-VRAM_START : tileMapAddr-VRAM_END]
	tileMapOffset := int(tileMapAddr - 0x8000)
	tileMap := vram[tileMapOffset : tileMapOffset+1024]
	for ty := 0; ty < mapHeight; ty++ {
		for tx := 0; tx < mapWidth; tx++ {
			tileIndex := tileMap[ty*mapWidth+tx]
			var tileAddr int
			if tileBaseAddr == 0x8000 {
				tileAddr = int(tileIndex) * 16
			} else {
				tileAddr = int(int8(tileIndex)) * 16
				tileAddr += 0x1000 //0x8800
			}

			if tileAddr < 0 || tileAddr+16 > len(vram) {
				fmt.Printf("Draw tile map - INVALID")
				continue
			}
			tileData := vram[tileAddr : tileAddr+16]
			for py := 0; py < tileHeight; py++ {
				low := tileData[py*2]
				high := tileData[py*2+1]
				for px := 0; px < tileWidth; px++ {
					bit := 7 - px
					l := (low >> bit) & 1
					h := (high >> bit) & 1
					colorID := (h << 1) | l
					color := colors[colorID]
					baseX := startX + tx*tileWidth*scale
					baseY := startY + ty*tileHeight*scale
					//rl.DrawRectangle(int32((startX+tx*tileWidth+px)*scale), int32((startY+ty*tileHeight+py)*scale), int32(scale), int32(scale), color)
					rl.DrawRectangle(int32(baseX+px*scale), int32(baseY+py*scale), int32(scale), int32(scale), color)

				}
			}
		}
	}
}
