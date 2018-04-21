package main

import (
	"chip8goh/chip8"
	"os"
)

func main() {
	// Initialize the Chip8 system and load the game into the memory
	// Check for debug mode
	isDebug := os.Args[1] == "-d"
	chip8.Initialize(isDebug)

	var fnIndex int
	// Go doesn't have a ternary operator
	if isDebug {
		fnIndex = 2
	} else {
		fnIndex = 1
	}
	chip8.LoadGame(os.Args[fnIndex])

	// // Setup graphics and start game
	chip8.SdlSetupGraphics()
}
