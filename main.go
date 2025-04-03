package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const screenWidth = 160
const screenHeight = 144

func main() {
	cartridge, err := LoadCartridge("roms/The Legend of Zelda - Links Awakening (US - EU).gb")
	if err != nil {
		fmt.Printf("Error loading cartridge: %s\n", err)
		return
	}

	cartridge.printInfo()
	//cpu := CPU{}
	cpu := NewCPU()
	

	cpu.loadROMFile(cartridge, cpu.graphics)

	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	const cyclesPerFrame = 70224
	for !rl.WindowShouldClose() {

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		fmt.Println("Startinf new frame")
		cyclesCurrFrame := 0

		//for addr := OAM_START; addr < OAM_END; addr++ {
		//	fmt.Printf("OAM[%04X]: %02X\n", addr, cpu.graphics.cpu.memoryRead(uint16(addr)))
		//}

		for cyclesCurrFrame < cyclesPerFrame {
			fmt.Println("")
			fmt.Println("Executing opcode")
			tCycles := cpu.execOpcodes()
			cpu.timer.Update(tCycles, cpu)
			cyclesCurrFrame += tCycles

			fmt.Printf("tCycles: %d, total: %d, graphics.cycles: %d\n", tCycles, cyclesCurrFrame, cpu.graphics.cycle)
			fmt.Printf("MEmory[0x0039]: 0x%02X\n", cpu.Memory[0x0039])

			cpu.graphics.modesHandeling(tCycles)
			fmt.Println("after modes handle = LY: ", cpu.graphics.LY, " mode: ", cpu.graphics.mode, " cycles total now: ", cpu.graphics.cycle)

		}
		fmt.Println("Ending frame")
		fmt.Printf("PC now %v", cpu.Registers.PC)
		//tCycles := cpu.execOpcodes()
		//fmt.Printf("vycles %v", tCycles)
		//graphics.modesHandeling(tCycles)

		//graphics.modesHandeling(70224)
		//graphics.testPixelDrawing()
		rl.DrawText("Gameboy emulator...!", 190, 200, 20, rl.LightGray)

		rl.EndDrawing()
		//cpu.run()

	}

}
