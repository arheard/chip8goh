package chip8

import (
	"fmt"
	"log"
	"os"
	"time"
)

const cycleTime = 16666670      // Time between each cycle
var opcode uint16               // Current opcode being processed
var memory [4096]uint8          // Chip8 RAM
var _V [16]uint8                // Registers
var _I uint16                   // Memeory Index
var pc uint16                   // Program counter
var chip8OpcodeTable [16]func() // Array of function pointers for opcode handling
var chip8op8 [15]func()         // Sub array of funciton handlers for 8xxx specific handling
var gfx [64 * 32]uint8          // Graphics buffer
var stack [16]uint16
var sp uint16              	// Stack Pointer
var key [16]uint8          	// Rey mapping key[index] = pressed (0 | 1)
var keyPressed bool        	// Flag for if a key was pressed during the last frame
var currentKeyValue uint8  	// Value for the key pressed within the last frame
var chip8Fontset [80]uint8 	// Stors internal representation of fontset
var delayTimer uint8
var soundTimer uint8
var lastCycleTime uint64 	// Epoch time of last cycle, used for assuring framerate
var cpuWaiting bool      	// Flag to determine if the program should perform instructions or not

//DrawFlag flag indicating if the scren should be drawn
var DrawFlag bool

var debug bool

//Initialize initalizes emulator
func Initialize(isDebug bool) {
	debug = isDebug

	chip8Fontset = [80]uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80} // F

	chip8OpcodeTable = [16]func(){
		cpu0, cpu1, cpu2, cpu3, cpu4, cpu5, cpu6, cpu7,
		cpu8, cpu9, cpuA, cpuB, cpuC, cpuD, cpuE, cpuF}

	chip8op8 = [15]func(){
		cpu8XY0, cpu8XY1, cpu8XY2, cpu8XY3, cpu8XY4, cpu8XY5, cpu8XY6, cpu8XY7,
		cpuNULL, cpuNULL, cpuNULL, cpuNULL, cpuNULL, cpuNULL, cpu8XYE}

	pc = 0x200 // Program counter starts at 0x200

	// Load fontset
	for i := 0; i < 80; i++ {
		memory[i] = chip8Fontset[i]
	}

}

//LoadGame loads game file into memory.
func LoadGame(filename string) {
	// Open program file
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Read program  file into buffer
	buffer := make([]byte, 4096-512) // total program / data space size
	bufferSize, err := f.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	// Copy program file into memory
	for i := 0; i < bufferSize; i++ {
		memory[i+512] = buffer[i]
	}
}

//EmulateCycle emulates a cycle of the emulator
func EmulateCycle() {
	currentTime := uint64(time.Now().UnixNano())
	if currentTime-lastCycleTime >= cycleTime {
		if !cpuWaiting {
			// Fetch opcode
			opcode = (uint16(memory[pc]) << 8) | uint16(memory[pc+1])
			if debug {
				fmt.Printf("%#x\n", opcode)
			}

			// Execute opcode
			chip8OpcodeTable[((opcode & 0xF000) >> 12)]()

			// Update timers
			if delayTimer > 0 {
				delayTimer--
			}

			if soundTimer > 0 {
				if soundTimer == 1 {
					fmt.Println("BEEP!")
				}
				soundTimer--
			}
		} else {
			// Check if input has been recieved from keyboard
			cpuFX00A()
		}
	}
}

//SetKeys sets keys pressed in emulator
func SetKeys(keyCode uint8, pressed uint8) {
	key[keyCode] = pressed
	currentKeyValue = keyCode
	if debug {
		fmt.Println("currentKeyValue: ", currentKeyValue)
	}
}
