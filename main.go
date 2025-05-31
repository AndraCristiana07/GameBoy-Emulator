package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	log "github.com/mgutz/logxi/v1"
)

const screenWidth = 160
const screenHeight = 144

var logger log.Logger

func main() {
	logger = log.New("main")

	//cartridge, err := LoadCartridge("roms/tetris.gb")
	//cartridge, err := LoadCartridge("roms/tetris-recompiled.gb")
	cartridge, err := LoadCartridge("roms/Tennis.gb")
	//cartridge, err := LoadCartridge("roms/gameGB.gb")
	//cartridge, err := LoadCartridge("roms/game.gb")
	//cartridge, err := LoadCartridge("roms/hello-world.gb")

	if err != nil {
		panic(logger.Error("Error loading cartridge:", err))
	}

	cartridge.printInfo()
	cpu := NewCPU()

	cpu.loadROMFile(cartridge)

	logger.Debug("ROM banking ", cpu.Memory[0x0148])
	if cpu.Memory[0x0148] != 0x0 {
		panic(logger.Error("No banking supported"))
	}
	rl.InitWindow(1500, 500, cartridge.title)

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		cpu.frameSteps()

		lcdc := cpu.graphics.getLCDC()

		if lcdc&(1<<7) == 0 {
			continue
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkBlue)

		cpu.graphics.render()

		drawTiles(&cpu.Memory, 1000, 0)
		drawTileMap(cpu, 400, 0)

		rl.EndDrawing()
	}

}
