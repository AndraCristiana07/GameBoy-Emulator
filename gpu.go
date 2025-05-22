package main

import (
	"fmt"
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"
	log "github.com/mgutz/logxi/v1"
)

var gpulogger log.Logger

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

func NewGraphics(cpu *CPU) *Graphics {
	gpulogger = log.New("gpu")
	graphic := &Graphics{cpu: cpu}
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

func (cpu *CPU) dmaTransfer(value uint8) {
	gpulogger.Debug("dmaTransfer")
	upper := uint16(value) << 8
	for i := 0; i < OAM_SIZE; i++ {
		gpulogger.Debug(fmt.Sprintf("upper is 0x%02X", upper))
		gpulogger.Debug(fmt.Sprintf("upper with i is 0x%02X", upper+uint16(i)))

		gpulogger.Debug(fmt.Sprintf("OAM at 0x%02X will be: 0x%02X", OAM_START+uint16(i), cpu.Memory[upper+uint16(i)]))
		cpu.Memory[OAM_START+uint16(i)] = cpu.Memory[upper+uint16(i)]
		gpulogger.Debug(fmt.Sprintf("DMA Transfer ->  %04X - %02X\n ", i, cpu.Memory[i]))
	}
	gpulogger.Debug("done")
}

func (cpu *CPU) readTileData(address uint16) [8][8]uint8 {
	var tile [8][8]uint8

	if address < VRAM_START || address > 0x97FF {
		// gpulogger.Debug(fmt.Sprintf("read tile data -> invalid address 0x%04X\n", address)
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
	// each tile taking 16 bytes, 2 per row
	// y row idx
	// x pixel idx
	rowAddr := address + uint16(y*2)
	bit := 7 - x

	low := graphic.cpu.memoryRead(uint16(rowAddr))
	high := graphic.cpu.memoryRead(uint16(rowAddr) + 1)
	msbLow := (low >> bit) & 1
	msbHigh := (high >> bit) & 1
	colorID := (msbHigh << 1) | msbLow
	return colorID

}

// sprites tiles
// OAM scan - mode 2
func (graphic *Graphics) spritesOAM() [height][width]uint8 {

	var spritePixels [height][width]uint8

	for y := range spritePixels {
		for x := range spritePixels[y] {
			spritePixels[y][x] = 255
		}
	}
	//obp0 := graphic.cpu.Memory[0xFF48]
	//obp1 := graphic.cpu.Memory[0xFF49]
	spriteSize := 8

	// In 8x16 sprite mode, the least significant bit of the
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

		//Byte 0 — Y Position
		//Y = Object’s vertical position on the screen + 16
		//y := int(graphic.OAM[spriteAddr]) - 16
		y := int(graphic.cpu.memoryRead(uint16(spriteAddr))) - 16

		// Byte 1 — X Position
		//X = Object’s horizontal position on the screen + 8.
		x := int(graphic.cpu.memoryRead(uint16(spriteAddr+1))) - 8

		// Byte 2 — Tile Index

		tileIndex := graphic.cpu.memoryRead(uint16(spriteAddr + 2))

		if spriteSize == 16 {
			tileIndex &= 0xFE // mask bit 0
		}

		//Byte 3 — Attributes/Flags
		attributes := graphic.cpu.memoryRead(uint16(spriteAddr + 3))

		//check if sprite is onscanline
		if int(graphic.getLY()) < y || int(graphic.getLY()) >= y+spriteSize {
			continue
		}

		visibleSprites = append(visibleSprites, Sprite{x, y, tileIndex, attributes, i})
	}

	// order by X (smaller X->bigger prority)
	// if x the same => by OAM location order
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
	ly := graphic.getLY()
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

		tileY := int(ly) - sprite.y
		if yFlip != 0 {
			tileY = spriteSize - 1 - tileY

		}
		if tileY < 0 || tileY >= 8 {
			continue
		}
		tileAddress := 0x8000 + (int(sprite.tileIndex) * 16)
		if tileAddress < VRAM_START || tileAddress > 0x97FF {
			gpulogger.Debug(fmt.Sprintf("tileAddress %04X is invalid\n", tileAddress))
			continue
		}
		tileData := graphic.cpu.readTileData(uint16(tileAddress))

		for col := 0; col < 8; col++ {
			tileX := col
			if xFlip != 0 {
				tileX = 7 - col
			}
			if tileX < 0 || tileX >= 8 {
				continue
			}
			pixelValue := tileData[tileY][tileX]

			if pixelValue == 0 {
				// gpulogger.Debug(fmt.Sprintf("skip transparent")
				continue // transparent
			}
			screenX := sprite.x + col
			screenY := ly

			// check bounds
			if screenX < 0 || screenY < 0 || screenX >= width || screenY >= height {
				// gpulogger.Debug(fmt.Sprintf("skip pixel out of bounds (x:%d, y:%d)\n", screenX, screenY)
				continue
			}

			bgPixelValue := bgPixels[screenY][screenX]
			if priority != 0 && bgPixelValue != 0 {
				// gpulogger.Debug(fmt.Sprintf("skip pixel -> priority bgPixelValue: %02X\n", bgPixelValue)
				continue
			}

			if spritePixels[ly][screenX] == 255 {
				// gpulogger.Debug(fmt.Sprintf("Drawing sprite pixel at ScreenX=%d ScreenY=%d, value=%d\n", screenX, screenY, pixelValue)
				//shade := (palette >> (pixelValue * 2)) & 0x03
				//// gpulogger.Debug("Sprite Palette", palette, "shade", shade)

				//spritePixels[graphic.getLY()][screenX] = shade
				spritePixels[ly][screenX] = pixelValue

			} else {
				gpulogger.Debug("Skip pixel. Pixel already occupied")
			}

		}

	}

	return spritePixels
}

// background tiles
func (graphic *Graphics) getBackground() [height][width]uint8 {
	var bgPixels [height][width]uint8

	scx := int(graphic.getSCX())
	scy := int(graphic.getSCY())
	//ly := int(graphic.getLY())
	lcdc := graphic.getLCDC()
	gpulogger.Debug(fmt.Sprintf("lcdc 0b%08b", lcdc))

	bgTileMapBase := uint16(0x9800)
	if (lcdc & (1 << 3)) != 0 {
		bgTileMapBase = 0x9C00
	}

	gpulogger.Debug(fmt.Sprintf("BG TileMapBase: 0x%04X\n", bgTileMapBase))
	for ly := 0; ly < height; ly++ {

		for screenX := 0; screenX < width; screenX++ {

			bgX := (screenX + scx) & 0xFF
			bgY := (ly + scy) & 0xFF

			tileMapX := bgX / 8
			tileMapY := bgY / 8
			tileIndexAddr := bgTileMapBase + uint16(tileMapY*32+tileMapX)

			tileIndex := graphic.cpu.memoryRead(tileIndexAddr)

			var tileAddress int

			if lcdc&(1<<4) != 0 {
				tileAddress = 0x8000 + int(tileIndex)*16
			} else {

				tileAddress = 0x9000 + int(int8(tileIndex))*16
			}
			gpulogger.Debug(fmt.Sprintf(" tileAddress: 0x%04X", tileAddress))

			tilePixelX := bgX % 8
			tilePixelY := bgY % 8
			colorID := graphic.getTilePixel(uint16(tileAddress), tilePixelX, tilePixelY)

			bgPixels[ly][screenX] = colorID

		}
	}
	return bgPixels
}

// window tiles
func (graphic *Graphics) getWindow() [height][width]uint8 {
	var window [height][width]uint8

	wx := graphic.getWX()
	wy := graphic.getWY()
	lcdc := graphic.getLCDC()
	// window display disabled
	if lcdc&(1<<5) == 0 {
		return window
	}

	for screenY := 0; screenY < height; screenY++ {
		if screenY < int(wy) {
			continue
		}
		for screenX := 0; screenX < width; screenX++ {

			if screenX < int(wx)-7 {
				continue
			}
			winY := screenY - int(wy)

			winX := screenX - (int(wx) - 7)
			winMapAddr := uint16(0x9800)
			if lcdc&(1<<6) != 0 {
				winMapAddr = uint16(0x9C00)
			}
			tileIdxAddr := winMapAddr + uint16((winY/8)*32+(winX/8))

			if tileIdxAddr < VRAM_START || tileIdxAddr > VRAM_END {
				continue
			}
			tileIndex := graphic.cpu.memoryRead(tileIdxAddr)
			var tileNumber int

			if graphic.getLCDC()&(1<<4) != 0 {
				tileNumber = int(tileIndex)
			} else {
				tileNumber = int(int8(tileIndex))
			}

			var tileAddr uint16
			if lcdc&(1<<4) != 0 {
				tileAddr = 0x8000 + uint16(tileNumber*16)
			} else {
				tileAddr = 0x9000 + uint16(tileNumber*16)
			}

			tileY := winY % 8
			tileX := winX % 8

			colorID := graphic.getTilePixel(uint16(tileAddr), tileX, tileY)
			//BGP Palette
			//palette := graphic.cpu.Memory[0xFF47]
			//shade := (palette >> (colorID * 2)) & 0x03
			//window[screenY][screenX] = shade
			window[screenY][screenX] = colorID
		}
	}
	return window
}

func (graphic *Graphics) renderScanline() {
	// gpulogger.Debug("In render scanline")
	ly := graphic.getLY()
	lcdc := graphic.getLCDC()

	if int(ly) >= height {
		return
	}
	//graphic.cpu.Memory[0xFF40] = 0b10010001
	bgPixels := graphic.getBackground()
	winPixels := graphic.getWindow()
	spritePixels := graphic.spritesOAM()

	for screenX := 0; screenX < width; screenX++ {
		// starting with background
		pixel := bgPixels[ly][screenX]

		// window overrides background
		if lcdc&(1<<5) != 0 && winPixels[ly][screenX] != 0 {
			pixel = winPixels[ly][screenX]

		}

		//sprites override, unless transparent
		if lcdc&(1<<1) != 0 && spritePixels[ly][screenX] != 255 {
			pixel = spritePixels[ly][screenX]
		}
		if lcdc&(1<<1) != 0 {
			spritePixel := spritePixels[ly][screenX]
			if spritePixel != 255 && spritePixel != 0 {
				pixel = spritePixel
			}
		}
		//color := colors[pixel]
		graphic.pixelBuffer[ly][screenX] = pixel
		// gpulogger.Debug(fmt.Sprintf("pixel: %d, color of pixel: %d\n", pixel, color)

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

func (graphic *Graphics) modesHandling(tCycles int) {
	// gpulogger.Debug("mode handling")
	gpulogger.Debug(fmt.Sprintf("LCDC: 0b%08b\n", graphic.getLCDC()))
	if graphic.cpu == nil {
		// gpulogger.Debug("Error: graphics.cpu is nil")
		return
	}
	ly := graphic.getLY()
	lcdc := graphic.getLCDC()

	if lcdc&(1<<7) == 0 {
		// gpulogger.Debug("LCD Disabled!")
		graphic.setMode(MODE_HBLANK)
		graphic.cycle = 0
		graphic.cpu.Memory[0xFF44] = 0 //reset
		return
	}
	graphic.cycle += tCycles

	if int(ly) >= SCANLINES_PER_FRAME {
		gpulogger.Debug("Entering VBLANK")
		graphic.setMode(MODE_VBLANK)
		if ly == SCANLINES_PER_FRAME {
			// gpulogger.Debug("Entering Vblank")
			IF := graphic.cpu.Memory[0xFF0F] | 1<<0
			graphic.cpu.memoryWrite(0xFF0F, IF)
		}
	} else {
		//if graphic.cycle >= 456-80 {
		if graphic.cycle < 80 {
			if graphic.cpu.Memory[0xFF41] != MODE_OAMSCAN {
				gpulogger.Debug("Entering OAMSCAN")
			}
			//gpulogger.Debug("Entering OAMSCAN")
			graphic.setMode(MODE_OAMSCAN)

			//} else if graphic.cycle >= 456-80-172 {
		} else if graphic.cycle < 80+172 {
			if graphic.cpu.Memory[0xFF41] != MODE_DRAWING {
				gpulogger.Debug("Entering drawing mode")
			}
			//gpulogger.Debug("Entering Drawing")
			graphic.setMode(MODE_DRAWING)

			if !graphic.drawnLine {
				graphic.renderScanline()
				//graphic.drawFrame()
				graphic.drawnLine = true
			}
		} else {
			if graphic.cpu.Memory[0xFF41] != MODE_HBLANK {
				gpulogger.Debug("Entering Hblank")
			}
			//gpulogger.Debug("Entering HBLANK")
			graphic.setMode(MODE_HBLANK)

		}

	}
	if graphic.cycle >= CYCLES_PER_LINE {
		graphic.cycle -= CYCLES_PER_LINE
		graphic.drawnLine = false

		gpulogger.Debug(fmt.Sprintf("ly inc to : %d", graphic.getLY()+1))
		graphic.cpu.Memory[0xFF44]++

		if graphic.getLY() >= TOTAL_LINES-1 {
			graphic.cpu.Memory[0xFF44] = 0

			graphic.cycle = 0

			gpulogger.Debug("reset ly, start new frame")
		}
	}
	//graphic.cpu.Memory[0xFF44] = (graphic.cpu.Memory[0xFF44] + 1) % 154
	// gpulogger.Debug(fmt.Sprintf("LY: %d Cycles: %d mode: %d  \n", graphic.getLY(), graphic.cycle, PPUmode)
}

//func (graphic *Graphics) modesHandling(tCycles int) {
//	if graphic.cpu == nil {
//		return
//	}
//
//	lcdc := graphic.getLCDC()
//
//	if lcdc&(1<<7) == 0 {
//		// LCD disabled
//		graphic.setMode(MODE_HBLANK)
//		graphic.cycle = 0
//		graphic.cpu.Memory[0xFF44] = 0 // reset LY
//		return
//	}
//
//	graphic.cycle += tCycles
//	ly := graphic.getLY()
//
//	switch graphic.getMode() {
//	case MODE_OAMSCAN:
//		if graphic.cycle >= 80 {
//			graphic.cycle -= 80
//			graphic.setMode(MODE_DRAWING)
//		}
//
//	case MODE_DRAWING:
//		if graphic.cycle >= 172 {
//			graphic.cycle -= 172
//			// drawing scanline
//			graphic.renderScanline()
//			graphic.setMode(MODE_HBLANK)
//		}
//
//	case MODE_HBLANK:
//		if graphic.cycle >= 204 {
//			graphic.cycle -= 204
//			graphic.cpu.Memory[0xFF44]++ // LY
//			ly = graphic.getLY()
//
//			if ly == 144 {
//				// Enter VBlank
//				graphic.setMode(MODE_VBLANK)
//
//				// VBlank interrupt
//				interrupt := graphic.cpu.Memory[0xFF0F]
//				interrupt |= 1 << 0
//				graphic.cpu.memoryWrite(0xFF0F, interrupt)
//
//				graphic.drawnLine = true
//			} else {
//				graphic.setMode(MODE_OAMSCAN)
//			}
//		}
//
//	case MODE_VBLANK:
//		if graphic.cycle >= 456 {
//			graphic.cycle -= 456
//			graphic.cpu.Memory[0xFF44]++
//			ly = graphic.getLY()
//
//			if ly > 153 {
//				//restart frame
//				graphic.cpu.Memory[0xFF44] = 0
//				graphic.setMode(MODE_OAMSCAN)
//			}
//		}
//	}
//}

func (graphic *Graphics) setMode(mode uint8) {
	stat := graphic.getSTAT()
	//		   clear bit         mask last 2 bits
	stat = (stat &^ 0b00000011) | (mode & 0b00000011)
	//stat = (stat & 0b11111100) | (mode & 0b00000011)

	graphic.cpu.memoryWrite(0xFF41, stat)

}

func (graphic *Graphics) getMode() uint8 {
	return graphic.cpu.memoryRead(0xFF41) & 0b00000011
}

func drawTiles(mem *[65536]uint8, startX int, startY int) {
	tileWidth := 8
	tileHeight := 8
	tileBytes := 16
	scale := 2
	tilesPerRow := 16

	totalTiles := 0x1800 / tileBytes
	for tileIdx := 0; tileIdx < totalTiles; tileIdx++ {
		tileX := tileIdx % tilesPerRow
		tileY := tileIdx / tilesPerRow

		screenBaseX := startX + tileX*tileWidth*scale
		screenBaseY := startY + tileY*tileHeight*scale

		offset := 0x8000 + tileIdx*tileBytes
		//if offset+16 > 0x1800 {
		//	gpulogger.Error("Am iesit din vram")
		//}
		tile := (*mem)[offset : offset+16]

		drawSingleTile(tile, scale, screenBaseX, screenBaseY)
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
			gpulogger.Debug(fmt.Sprintf("Draw sprites - INVALID"))
			continue
		}
		if tileIndex < 0 || tileIndex >= 384 {
			gpulogger.Debug(fmt.Sprintf("Draw sprites - out of bounds"))
			continue
		}

		tileStart := tileIndex * 16
		if tileStart+16 > len(vram) {
			gpulogger.Debug(fmt.Sprintf("Draw sprites - out of vrqm bounds"))
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
					gpulogger.Debug("transparent")
					continue
				}
				//offset := 500

				rl.DrawRectangle(int32((startX+x+tx)*scale), int32((startY+y+ty)*scale), int32(scale), int32(scale), rl.Red)

			}
		}

	}
}

func drawSingleTile(tile []uint8, scale int, tx int, ty int) {
	tileHeight := 8
	tileWidth := 8

	for y := 0; y < tileHeight; y++ {
		low := tile[y*2]
		high := tile[y*2+1]

		for x := 0; x < tileWidth; x++ {
			bit := 7 - x
			msbLow := (low >> bit) & 1
			msbHigh := (high >> bit) & 1
			colorId := (msbHigh << 1) | msbLow

			color := colors[colorId]

			// draw scaled pixel
			screenX := tx + x*scale
			screenY := ty + y*scale
			rl.DrawRectangle(int32(screenX), int32(screenY), int32(scale), int32(scale), color)
		}
	}
}

func (cpu *CPU) getTileId(memVal uint8) int {
	if cpu.graphics.getLCDC()&(1<<4) == 0 {
		return int(int8(memVal))
	}
	return int(memVal)
}

func getTileMapAddr(lcdc byte) int {
	if lcdc&(1<<3) != 0 {
		return 0x9C00
	}
	return 0x9800
}

func getTileBaseAddr(lcdc byte) int {
	if lcdc&(1<<4) != 0 {
		return 0x8000
	}
	return 0x9000
}

func drawTileMap(cpu *CPU, startX int, startY int) {
	tileMapAddr := 0x9C00
	if cpu.graphics.getLCDC()&(1<<3) == 0 {
		tileMapAddr = 0x9800
	}

	tileBase := 0x8000
	if cpu.graphics.getLCDC()&(1<<4) == 0 {
		tileBase = 0x9000
	}

	mem := &cpu.Memory

	tileWidth := 8
	tileHeight := 8
	mapWidth := 32
	mapHeight := 32
	scale := 2

	tileMap := (*mem)[tileMapAddr : tileMapAddr+1024]
	for ty := 0; ty < mapHeight; ty++ {
		for tx := 0; tx < mapWidth; tx++ {
			tileIndex := cpu.getTileId(tileMap[ty*mapWidth+tx])
			var tileAddr = tileBase + int(tileIndex)*16

			if tileAddr < 0 || tileAddr+16 > len(*mem) {
				panic(gpulogger.Error(fmt.Sprintf("Draw tile map - INVALID")))
			}

			tileData := (*mem)[tileAddr : tileAddr+16]

			baseX := startX + tx*tileWidth*scale
			baseY := startY + ty*tileHeight*scale

			drawSingleTile(tileData, scale, baseX, baseY)
		}
	}
}

func (graphic *Graphics) drawBackground() {
	scale := 2
	bgPixels := graphic.getBackground()
	for ly := 0; ly < height; ly++ {
		for screenX := 0; screenX < width; screenX++ {
			pixel := bgPixels[ly][screenX]
			rl.DrawRectangle(int32(screenX*scale), int32(ly*scale), int32(scale), int32(scale), colors[pixel])
		}
	}
}

func (graphic *Graphics) drawWindow() {
	scale := 2
	winPixels := graphic.getWindow()
	for ly := 0; ly < height; ly++ {
		for screenX := 0; screenX < width; screenX++ {
			pixel := winPixels[ly][screenX]
			rl.DrawRectangle(int32(screenX*scale), int32(ly*scale), int32(scale), int32(scale), colors[pixel])
		}
	}
}

func (graphic *Graphics) drawFrame() {
	lcdc := graphic.getLCDC()

	bgPixels := graphic.getBackground()
	winPixels := graphic.getWindow()
	spritePixels := graphic.spritesOAM()

	for ly := 0; ly < height; ly++ {
		for screenX := 0; screenX < width; screenX++ {
			var pixel uint8 = 0

			if lcdc&(1<<0) != 0 {
				pixel = bgPixels[ly][screenX]
				// window overrides background
				//if lcdc&(1<<5) != 0 && winPixels[ly][screenX] != 0 {
				if lcdc&(1<<5) != 0 {
					pixel = winPixels[ly][screenX]
				}
			}

			//sprites override
			//if lcdc&(1<<1) != 0 && spritePixels[ly][screenX] != 255 {
			if lcdc&(1<<1) != 0 {
				spriteCol := spritePixels[ly][screenX]
				if spriteCol != 255 && spriteCol != 0 {
					pixel = spriteCol
				}
				//pixel = spritePixels[ly][screenX]
			}
			graphic.pixelBuffer[ly][screenX] = pixel
		}
	}

	//graphic.drawScreen()

}

func (graphic *Graphics) render() {
	graphic.drawFrame()
	graphic.drawScreen()
}
