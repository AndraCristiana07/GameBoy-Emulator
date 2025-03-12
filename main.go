package main

import (
	"fmt"
)

func main() {
	//rl.InitWindow(800, 450, "raylib [core] example - basic window")
	//defer rl.CloseWindow()
	//
	//rl.SetTargetFPS(60)
	//
	//for !rl.WindowShouldClose() {
	//	rl.BeginDrawing()
	//
	//	rl.ClearBackground(rl.RayWhite)
	//	rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)
	//
	//	rl.EndDrawing()
	//}
	cartridge, err := loadCartridge("roms/The Legend of Zelda - Links Awakening (US - EU).gb")
	if err != nil {
		fmt.Printf("Error loading cartridge: %s\n", err)
		return
	}
	cartridge.printInfo()
	cpu := CPU{}

	cpu.loadROMFile(cartridge)
	cpu.run()
}
