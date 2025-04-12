package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const screenWidth = 160
const screenHeight = 144

func main() {
	cartridge, err := LoadCartridge("roms/tetris.gb")
	//cartridge, err := LoadCartridge("roms/The Legend of Zelda - Links Awakening (US - EU).gb")

	if err != nil {
		// fmt.Printf("Error loading cartridge: %s\n", err)
		return
	}

	cartridge.printInfo()
	//cpu := CPU{}
	cpu := NewCPU()
	//// fmt.Printf("LCDC after init: 0b%08b\n", cpu.Memory[0xFF40])

	cpu.loadROMFile(cartridge)

	//cpu.Memory[0xFF40] = 0x91
	//cpu.Memory[0xFF40] |= 1 << 1

	//rl.InitWindow(800, 450, "raylib [core] example - basic window")
	rl.InitWindow(1500, 1000, "raylib [core] example - basic window")

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	const cyclesPerFrame = 70224
	for !rl.WindowShouldClose() {

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)
		//cpu.Memory[0xFF40] = 0x91
		//cpu.Memory[0xFF40] |= 1 << 1
		//// fmt.Println("Starting new frame")
		//// fmt.Printf("LCDC at start: 0b%08b\n", cpu.Memory[0xFF40])

		//LCDC := cpu.Memory[0xFF40]
		//if LCDC&(1<<7) == 0 {
		//cpu.Memory[0xFF40] |= 1 << 1
		//}

		cyclesCurrFrame := 0

		//for addr := OAM_START; addr < OAM_END; addr++ {
		//	// fmt.Printf("OAM[%04X]: %02X\n", addr, cpu.graphics.cpu.memoryRead(uint16(addr)))
		//}
		//Fake dma trigger
		//cpu.memoryWrite(0xFF46, 0x80)

		for i := 0; i < 0xA0; i += 4 {
			y := cpu.Memory[0xFE00+i]
			x := cpu.Memory[0xFE00+i+1]
			tile := cpu.Memory[0xFE00+i+2]
			attr := cpu.Memory[0xFE00+i+3]
			fmt.Printf("OAM ::: Sprite %02d - y: %02x x: %02x tile: %02x attr:%02X\n", i/4, y, x, tile, attr)
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
			fmt.Printf("Tile %d:\n", i)
			fmt.Printf("Binary %08b\n", cpu.Memory[0x8000+i*16:0x8000+(i+1)*16])
			fmt.Printf("Hexa %02X\n", cpu.Memory[0x8000+i*16:0x8000+(i+1)*16])
		}

		//drawTiles(cpu.Memory[0x8000:0x9800], 200, 10)
		vram := cpu.Memory[0x8000:0x9800]
		drawTiles(vram, 800, 10)

		//drawTiles(dummyVRAM, 500, 10)

		for cyclesCurrFrame < cyclesPerFrame {
			fmt.Printf("STAT is 0b%08b\n", cpu.Memory[0xFF41])
			// fmt.Println("\nExecuting opcode")
			tCycles := cpu.execOpcodes()
			//cpu.Memory[0xFF40] |= 1 << 1

			cpu.timer.Update(tCycles, cpu)
			cyclesCurrFrame += tCycles

			//// fmt.Printf("tCycles: %d, total: %d, graphics.cycles: %d\n", tCycles, cyclesCurrFrame, cpu.graphics.cycle)
			//// fmt.Printf("Memory[0x0039]: 0x%02X\n", cpu.Memory[0x0039])

			cpu.graphics.modesHandeling(tCycles)
			//// fmt.Println("after modes handle = LY: ", cpu.Memory[0xFF44], " cycles total now: ", cpu.graphics.cycle)

		}
		//// fmt.Printf("Ending frame\n PC now %v ", cpu.Registers.PC)
		//tCycles := cpu.execOpcodes()
		//// fmt.Printf("vycles %v", tCycles)
		//graphics.modesHandeling(tCycles)

		//graphics.modesHandeling(70224)
		//graphics.testPixelDrawing()
		//rl.DrawText("Gameboy emulator...!", 190, 200, 20, rl.LightGray)

		rl.EndDrawing()
		//cpu.run()

	}

}
