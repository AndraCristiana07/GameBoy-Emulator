package main

import (
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
	log "github.com/mgutz/logxi/v1"
)

// const screenWidth = 160
// const screenHeight = 144

var logger log.Logger

func main() {
	logger = log.New("main")
	romPath := os.Args[1]
	if romPath == "" {
		fmt.Print("No ROM given")
	}
	cartridge, err := LoadCartridge(romPath)

	if err != nil {
		panic(logger.Error("Error loading cartridge:", err))
	}

	cartridge.printInfo()
	cpu := NewCPU()

	cpu.loadROMFile(cartridge)
	cpulogger.Debug(fmt.Sprintf("at 0x2000: 0x%x", cpu.Memory[0x2000]))

	logger.Debug("ROM banking ", cpu.Memory[0x0148])
	if cpu.Memory[0x0147] != 0x0 {
		panic(logger.Error("No banking supported"))
	}
	rl.InitWindow(1500, 500, cartridge.title)

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkBlue)

		cpu.frameSteps()

		lcdc := cpu.graphics.getLCDC()
		for i := 0; i < 0xA0; i += 4 {
			y := cpu.Memory[0xFE00+i]
			x := cpu.Memory[0xFE00+i+1]
			tile := cpu.Memory[0xFE00+i+2]
			attr := cpu.Memory[0xFE00+i+3]
			logger.Debug(fmt.Sprintf("OAM ::: Sprite %02d - y: %02x x: %02x tile: %02x attr:%02X\n", i/4, y, x, tile, attr))
		}
		if lcdc&(1<<7) == 0 {
			continue
		}

		drawTiles(&cpu.Memory, 1000, 0)
		drawTileMap(cpu, 400, 0)

		rl.EndDrawing()
	}

}
