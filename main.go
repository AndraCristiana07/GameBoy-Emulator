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
	graphics.OAM[0] = 50 + 16 //y
	graphics.OAM[1] = 50 + 8  //x
	graphics.OAM[2] = 2       //idx
	graphics.OAM[3] = 0       //attrib

	graphics.OAM[4] = 51 + 16 //y
	graphics.OAM[5] = 51 + 8  //x
	graphics.OAM[6] = 3       //idx
	graphics.OAM[7] = 0       //attrib

	graphics.OAM[8] = 52 + 16 //y
	graphics.OAM[9] = 52 + 8  //x
	graphics.OAM[10] = 4      //idx
	graphics.OAM[11] = 0      //attrib

	graphics.OAM[12] = 53 + 16 //y
	graphics.OAM[13] = 53 + 8  //x
	graphics.OAM[14] = 5       //idx
	graphics.OAM[15] = 0       //attrib

	graphics.OAM[12] = 54 + 16 //y
	graphics.OAM[13] = 54 + 8  //x
	graphics.OAM[14] = 6       //idx
	graphics.OAM[15] = 0       //attrib

	graphics.OAM[16] = 55 + 16 //y
	graphics.OAM[17] = 55 + 8  //x
	graphics.OAM[18] = 7       //idx
	graphics.OAM[19] = 0       //attrib

	graphics.OAM[20] = 56 + 16 //y
	graphics.OAM[21] = 56 + 8  //x
	graphics.OAM[22] = 8       //idx
	graphics.OAM[23] = 0       //attrib

	graphics.OAM[24] = 57 + 16 //y
	graphics.OAM[25] = 57 + 8  //x
	graphics.OAM[26] = 9       //idx
	graphics.OAM[27] = 0       //attrib

	graphics.OAM[28] = 58 + 16 //y
	graphics.OAM[29] = 58 + 8  //x
	graphics.OAM[30] = 10      //idx
	graphics.OAM[31] = 0       //attrib

	graphics.OAM[32] = 59 + 16 //y
	graphics.OAM[33] = 59 + 8  //x
	graphics.OAM[34] = 11      //idx
	graphics.OAM[35] = 0       //attrib

	graphics.OAM[36] = 60 + 16 //y
	graphics.OAM[37] = 60 + 8  //x
	graphics.OAM[38] = 12      //idx
	graphics.OAM[39] = 0       //attrib

	graphics.OAM[40] = 61 + 16 //y
	graphics.OAM[41] = 61 + 8  //x
	graphics.OAM[42] = 13      //idx
	graphics.OAM[43] = 0       //attrib

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
