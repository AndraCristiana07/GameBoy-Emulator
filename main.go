package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	log "github.com/mgutz/logxi/v1"
)

const screenWidth = 160
const screenHeight = 144

var logger log.Logger

func main() {
	logger = log.New("main")

	cartridge, err := LoadCartridge("roms/tetris.gb")
	//cartridge, err := LoadCartridge("blargg test/gb-test-roms/cpu_instrs/cpu_instrs.gb")
	//cartridge, err := LoadCartridge("roms/The Legend of Zelda - Links Awakening (US - EU).gb")

	if err != nil {
		logger.Error("Error loading cartridge:", err)
	}

	cartridge.printInfo()
	//cpu := CPU{}
	cpu := NewCPU()
	//// logger.Debug(fmt.Sprintf("LCDC after init: 0b%08b\n", cpu.Memory[0xFF40])

	cpu.loadROMFile(cartridge)

	//rl.InitWindow(800, 450, "raylib [core] example - basic window")
	rl.InitWindow(1500, 1000, "raylib [core] example - basic window")

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	const cyclesPerFrame = 70224
	for !rl.WindowShouldClose() {
		cyclesCurrFrame := 0
		for cyclesCurrFrame < cyclesPerFrame {
			logger.Debug(fmt.Sprintf("STAT is 0b%08b\n", cpu.Memory[0xFF41]))
			// fmt.Println("\nExecuting opcode")
			tCycles := cpu.execOpcodes()
			//cpu.Memory[0xFF40] |= 1 << 1

			cpu.timer.Update(tCycles, cpu)
			cpu.graphics.modesHandling(tCycles)

			cyclesCurrFrame += tCycles

			logger.Debug(fmt.Sprintf("tCycles: %d, total: %d, graphics.cycles: %d\n", tCycles, cyclesCurrFrame, cpu.graphics.cycle))
			//// logger.Debug(fmt.Sprintf("Memory[0x0039]: 0x%02X\n", cpu.Memory[0x0039])
			//cpu.graphics.modesHandling(tCycles)

			//// fmt.Println("after modes handle = LY: ", cpu.Memory[0xFF44], " cycles total now: ", cpu.graphics.cycle)

		}
		rl.BeginDrawing()

		rl.ClearBackground(rl.DarkBlue)

		lcdc := cpu.graphics.getLCDC()
		
		//cpu.graphics.drawFrame()
		cpu.graphics.render()
		//for addr := OAM_START; addr < OAM_END; addr++ {
		//	// logger.Debug(fmt.Sprintf("OAM[%04X]: %02X\n", addr, cpu.graphics.cpu.memoryRead(uint16(addr)))
		//}
		//Fake dma trigger
		//cpu.memoryWrite(0xFF46, 0x80)

		//TODO: erase -> for debugging
		//cpu.Memory[0xFE00+0] = 50
		//cpu.Memory[0xFE00+1] = 50
		//cpu.Memory[0xFE00+2] = 0
		//cpu.Memory[0xFE00+3] = 0
		//
		//cpu.Memory[0xFE00+4] = 58
		//cpu.Memory[0xFE00+5] = 50
		//cpu.Memory[0xFE00+6] = 1
		//cpu.Memory[0xFE00+7] = 0
		////
		//cpu.Memory[0xFE00+8] = 66
		//cpu.Memory[0xFE00+9] = 50
		//cpu.Memory[0xFE00+10] = 2
		//cpu.Memory[0xFE00+11] = 0

		for i := 0; i < 0xA0; i += 4 {
			y := cpu.Memory[0xFE00+i]
			x := cpu.Memory[0xFE00+i+1]
			tile := cpu.Memory[0xFE00+i+2]
			attr := cpu.Memory[0xFE00+i+3]

			logger.Debug(fmt.Sprintf("OAM ::: Sprite %02d - y: %02x x: %02x tile: %02x attr:%02X\n", i/4, y, x, tile, attr))
		}
		// debug
		dummyVRAM := make([]uint8, 384*16)
		for i := 0; i < len(dummyVRAM); i += 2 {
			dummyVRAM[i] = 0xFF
			dummyVRAM[i+1] = 0x00
		}

		//for i := 0; i < 160; i++ {
		//	cpu.Memory[0xFE00+uint16(i)] = cpu.Memory[0x8000+uint16(i)]
		//}

		//VRAM
		//// fmt.Println("VRAM Dump")
		for i := 0; i < 10; i++ {
			logger.Debug(fmt.Sprintf("Tile %d:\n", i))
			logger.Debug(fmt.Sprintf("Binary %08b\n", cpu.Memory[0x8000+i*16:0x8000+(i+1)*16]))
			logger.Debug(fmt.Sprintf("Hexa %02X\n", cpu.Memory[0x8000+i*16:0x8000+(i+1)*16]))
		}

		//drawTiles(cpu.Memory[0x8000:0x9800], 200, 10)
		//vram := cpu.Memory[0x8000:0x9800]
		vram := cpu.Memory[0x8000:0xA000]

		drawTiles(vram, 800, 10)

		//drawTiles(dummyVRAM, 500, 10)

		oam := cpu.Memory[0xFE00:0xFEA0]

		for i := 0; i < 5; i++ {
			tileIndex := oam[i*4+2]
			tileData := oam[tileIndex*16 : tileIndex*16+16]
			logger.Debug(fmt.Sprintf("DEBUG ::: Sprite %d tile %d data %X\n", i, tileIndex, tileData))
		}
		drawSprites(oam, vram, 50, 200)
		//drawSprites(0, 500)

		tileMapAddr := uint16(0x9800)
		//if lcdc&(1<<3) != 0 {
		//	tileMapAddr = 0x9C00
		//}
		tileBase := uint16(0x8000)
		//if lcdc&(1<<4) == 0 {
		//	tileBase = 0x8800
		//}
		ly := cpu.graphics.getLY()
		logger.Debug(fmt.Sprintf("DEBUG: LCDC: %08b tileMapAddr: 0x%04X tikeBase: 0x%04X LY %d \n", lcdc, tileMapAddr, tileBase, ly))
		drawTileMap(vram, tileMapAddr, tileBase, 500, 500)

		//scale := 4
		//for y := 0; y < screenHeight; y++ {
		//	for x := 0; x < screenWidth; x++ {
		//		color := colors[cpu.graphics.pixelBuffer[y][x]]
		//		rl.DrawRectangle(int32(x*scale), int32(y*scale), int32(scale), int32(scale), color)
		//	}
		//}

		//cpu.graphics.drawScreen()
		//cpu.graphics.drawBackground()
		logger.Debug(fmt.Sprintf("DEBUG LY: %d\n", cpu.Memory[0xFF44]))

		rl.EndDrawing()
		//cpu.run()

	}

}
