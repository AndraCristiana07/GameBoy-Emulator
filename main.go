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
	cpu := CPU{}
	graphics := Graphics{}
	//test pixel
	graphics.OAM[0] = 50 + 16 //y
	graphics.OAM[1] = 50 + 8  //x
	graphics.OAM[2] = 2       //idx
	graphics.OAM[3] = 0       //attrib

	cpu.loadROMFile(cartridge)

	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		graphics.spritesOAM()
		//graphics.testPixelDrawing()
		rl.DrawText("Gameboy emulator...!", 190, 200, 20, rl.LightGray)

		rl.EndDrawing()
		//cpu.run()

	}

}
