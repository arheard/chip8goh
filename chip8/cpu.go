/*
	CPU.go:
	This file stores all of the opcode functions to be executed during
	the emulation cycle, as well as several convience functions.
*/
package chip8

import (
	"fmt"
	"math/rand"
)

// CPU does nothing
func cpuNULL() {}

// Opcode 00xx
func cpu0() {
	switch opcode & 0x000F {
	case 0x0000:
		// Execute opcode
		cpu00E0()
	case 0x000E:
		// Execute opcode
		cpu00EE()
	default:
		printOpcodeErr("0000")
	}
}

// 0x00E0: Clears the screen
func cpu00E0() {
	if debug {
		fmt.Println("CLS: Clearing the display.")
	}

	screenLength := 64 * 32
	for i := 0; i < screenLength; i++ {
		gfx[i] = 0
	}
	DrawFlag = true

	pc += 2
}

// 0x00EE: Returns from subroutine
func cpu00EE() {
	if debug {
		fmt.Println("RET - Return from a subroutine")
		fmt.Printf("Returning from function at %#x\n", currentSub)
	}

	sp--
	pc = stack[sp]
	pc += 2
	currentSub = pc
}

/*
	Opcode 1xxx -JP addr
	Jump to locaiton nnn
*/
func cpu1() {
	pc = opcode & 0x0FFF

	if debug {
		fmt.Printf("JP %#x\n", pc)
		fmt.Printf("Jumpiong to instruction at %#x\n", pc)
	}

}

var currentSub uint16

/*
	2nnn - CALL addr
	Call subroutine at nnn
*/
func cpu2() {
	stack[sp] = pc
	sp++
	pc = opcode & 0x0FFF
	currentSub = pc
	if debug {
		fmt.Printf("CALL  %#x\n", pc)
		fmt.Printf("Call subroutine at %#x\n", pc)
	}
}

/*
	3xkk - SE _Vx, byte
	Skip next instruction if _Vx = kk.
*/
func cpu3() {
	x, kk := xkk()

	if debug {
		fmt.Printf("SE V[%d], %d", x, kk)
		fmt.Println("skip if Vx == kk ")
		fmt.Println("x: ", x, " kk: ", kk, "_V[x]: ", _V[x])
	}
	if _V[x] == kk {
		pc += 2
	}
	pc += 2
}

/*
	4xkk - SNE _Vx, byte
	Skip next instruction if _Vx != kk.
*/
func cpu4() {
	x, kk := xkk()

	if debug {
		fmt.Printf("SNE  V[%d], %d", x, kk)
		fmt.Println("skip if Vx == kk ")
		fmt.Println("x: ", x, " kk: ", kk, "_V[x]: ", _V[x])
	}

	if _V[x] != kk {
		pc += 2
	}
	pc += 2
}

/*
	5xy0 - SE _Vx, _Vy
	Skip next instruction if _Vx = _Vy.
*/
func cpu5() {
	x, y := getxy()

	if debug {
		fmt.Printf("SE  V[%d], V[%d]", x, y)
		fmt.Println("skip if Vx == Vy ")
		fmt.Println("x: ", x, " y: ", y, "_V[x]: ", _V[x], "_V[y]: ", _V[y])
	}

	if _V[x] == _V[y] {
		pc += 2
	}
	pc += 2
}

/*
	6xkk - LD _Vx, byte
	Set _Vx = kk.

	The interpreter puts the value kk into register _Vx.
*/
func cpu6() {
	x, kk := xkk()

	_V[x] = kk
	pc += 2
}

/*
	7xkk - ADD _Vx, byte
	Set _Vx = _Vx + kk.

	Adds the value kk to the value of register _Vx, then stores the result in _Vx.
*/
func cpu7() {
	x, kk := xkk()

	_V[x] = _V[x] + kk
	pc += 2
}

// Arithmetic
func cpu8() {
	chip8op8[(opcode & 0x000F)]()
}

/*
	8xy0 - LD _Vx, _Vy
	Set _Vx = _Vy.

	Stores the value of register _Vy in register _Vx.
*/
func cpu8XY0() {
	x, y := getxy()

	_V[x] = _V[y]
	pc += 2
}

/*
	8xy1 - OR _Vx, _Vy
	Set _Vx = _Vx OR _Vy.

	Performs a bitwise OR on the values of _Vx and _Vy, then stores the result in _Vx.
	A bitwise OR compares the corrseponding bits from two values, and if either bit is 1,
	then the same bit in the result is also 1. Otherwise, it is 0.
*/
func cpu8XY1() {
	x, y := getxy()

	_V[x] = _V[x] | _V[y]
	pc += 2
}

/*
	8xy2 - AND _Vx, _Vy
	Set _Vx = _Vx AND _Vy.

	Performs a bitwise AND on the values of _Vx and _Vy, then stores the result in _Vx.
	A bitwise AND compares the corrseponding bits from two values, and if both bits are 1,
	then the same bit in the result is also 1. Otherwise, it is 0.
*/
func cpu8XY2() {
	x, y := getxy()

	_V[x] = _V[x] & _V[y]
	pc += 2
}

/*
	8xy3 - XOR _Vx, _Vy
	Set _Vx = _Vx XOR _Vy.

	Performs a bitwise exclusive OR on the values of _Vx and _Vy, then stores the result in _Vx.
	An exclusive OR compares the corrseponding bits from two values, and if the bits are not
	both the same, then the corresponding bit in the result is set to 1. Otherwise, it is 0.
*/
func cpu8XY3() {
	x, y := getxy()

	_V[x] = _V[x] ^ _V[y]
	pc += 2
}

/*
	8xy4 - ADD _Vx, _Vy
	Set _Vx = _Vx + _Vy, set _VF = carry.

	The values of _Vx and _Vy are added together. _If the result is greater than 8 bits (i.e., > 255,)
	_VF is set to 1, otherwise 0. Only the lowest 8 bits of the result are kept, and stored in _Vx.
*/
func cpu8XY4() {
	if _V[gety()] > (0xFF - _V[getx()]) {
		_V[0xF] = 1 //carry
	} else {
		_V[0xF] = 0
	}
	_V[getx()] += _V[gety()]
	pc += 2
}

/*
	8xy5 - SUB _Vx, _Vy
	Set _Vx = _Vx - _Vy, set _VF = NOT borrow.

	_If _Vx > _Vy, then _VF is set to 1, otherwise 0. Then _Vy is subtracted from _Vx,
	and the results stored in _Vx.
*/
func cpu8XY5() {
	x, y := getxy()

	if _V[x] > _V[y] {
		_V[0xF] = 1
	} else {
		_V[0xF] = 0
	}
	_V[x] = _V[x] - _V[y]
	pc += 2
}

/*
	8xy6 - SHR _Vx {, _Vy}
	Set _Vx = _Vx SHR 1.

	_If the least-significant bit of _Vx is 1, then _VF is set to 1,
	otherwise 0. Then _Vx is divided by 2.
*/
func cpu8XY6() {
	x := getx()
	if _V[x]&1 == 1 {
		_V[0xF] = 1
	} else {
		_V[0xF] = 0
	}
	_V[x] = _V[x] >> 1
	pc += 2
}

/*
	8xy7 - SUBN _Vx, _Vy
	Set _Vx = _Vy - _Vx, set _VF = NOT borrow.

	_If _Vy > _Vx, then _VF is set to 1, otherwise 0. Then _Vx is subtracted from _Vy,
	and the results stored in _Vx.
*/
func cpu8XY7() {
	x, y := getxy()

	if _V[y] > _V[x] {
		_V[0xF] = 1
	} else {
		_V[0xF] = 0
	}
	_V[x] = _V[y] - _V[x]
	pc += 2
}

/*
	8xyE - SHL _Vx {, _Vy}
	Set _Vx = _Vx SHL 1.

	_If the most-significant bit of _Vx is 1, then _VF is set to 1, otherwise to 0. Then _Vx is multiplied by 2.
*/
func cpu8XYE() {
	x := getx()
	_V[0xF] = _V[x] >> 7
	_V[x] <<= 1

	pc += 2
}

/*
	9xy0 - SNE _Vx, _Vy
	Skip next instruction if _Vx != _Vy.

	The values of _Vx and _Vy are compared, and if they are not equal,
	 the program counter is increased by 2.
*/
func cpu9() {
	x, y := getxy()

	if _V[x] != _V[y] {
		pc += 2
	}
	pc += 2
}

/*
	Annn - LD _I, addr
	Set _I = nnn.

	The value of register _I is set to nnn.
*/
func cpuA() {
	nnn := opcode & 0x0FFF
	_I = nnn
	pc += 2
}

/*
	Bnnn - JP _V0, addr
	Jump to location nnn + _V0.

	The program counter is set to nnn plus the value of _V0.
*/
func cpuB() {
	nnn := opcode & 0x0FFF
	pc = nnn + uint16(_V[0])
}

/*
	Cxkk - RND _Vx, byte
	Set _Vx = random byte AND kk.

*/
func cpuC() {
	x, kk := xkk()

	randomByte := uint8(rand.Intn(255))

	_V[x] = kk & randomByte
	pc += 2
}

func cpuD() {
	x := uint16(_V[getx()])
	y := uint16(_V[gety()])
	height := opcode & 0x000F
	var pixel, yline, xline uint16

	_V[0xF] = 0
	for yline = 0; yline < height; yline++ {
		pixel = uint16(memory[_I+yline])
		for xline = 0; xline < 8; xline++ {
			if (pixel & (0x80 >> xline)) != 0 {
				// This wraps around, stopping index oob errors
				if gfx[(x+xline+((y+yline)*64))%(gameHeight*gameWidth)] == 1 {
					_V[0xF] = 1
				}
				gfx[(x+xline+((y+yline)*64))%(gameHeight*gameWidth)] ^= 1
			}
		}
	}

	DrawFlag = true
	pc += 2

}

func cpuE() {
	switch opcode & 0x00FF {
	case 0x009E:
		cpuEX9E()
	case 0x00A1:
		cpuEXA1()
	default:
		printOpcodeErr("E000")
	}
}

/*
	Ex9E - SKP _Vx
	Skip next instruction if key with the value of _Vx is pressed.

	Checks the keyboard, and if the key corresponding to the value of _Vx is currently
	in the down position, PC is increased by 2.
*/
func cpuEX9E() {
	if key[_V[getx()]] != 0 {
		pc += 4
	} else {
		pc += 2
	}
}

/*
	ExA1 - SKNP _Vx
	Skip next instruction if key with the value of _Vx is not pressed.

	Checks the keyboard, and if the key corresponding to the value of _Vx is
	currently in the up position, PC is increased by 2.
*/
func cpuEXA1() {
	if key[_V[getx()]] == 0 {
		pc += 4
	} else {
		pc += 2
	}
}

func cpuF() {
	switch opcode & 0x00FF {
	case 0x0007:
		cpuFX007()
	case 0x000A:
		cpuFX00A()
	case 0x0015:
		cpuFX015()
	case 0x0018:
		cpuFX018()
	case 0x001E:
		cpuFX01E()
	case 0x0029:
		cpuFX029()
	case 0x0033:
		cpuFX033()
	case 0x0055:
		cpuFX055()
	case 0x0065:
		cpuFX065()
	default:
		printOpcodeErr("FX00")
	}
}

/*
	Fx07 - LD _Vx, DT
	Set _Vx = delay timer value.

	The value of DT is placed into _Vx.
*/
func cpuFX007() {
	x := getx()
	_V[x] = delayTimer
	pc += 2
}

/*
	Fx0A - LD _Vx, K
	Wait for a key press, store the value of the key in _Vx.

	All execution stops until a key is pressed, then the value of that key is stored in _Vx.
*/
func cpuFX00A() {
	if cpuWaiting {
		if keyPressed {
			x := getx()
			_V[x] = currentKeyValue

			pc += 2
			keyPressed = false
			cpuWaiting = false
		}
	} else {
		cpuWaiting = true
	}
}

/*
	Fx15 - LD DT, _Vx
	Set delay timer = _Vx.

	DT is set equal to the value of _Vx.
*/
func cpuFX015() {
	x := getx()
	delayTimer = _V[x]
	pc += 2
}

/*
	Fx18 - LD ST, _Vx
	Set sound timer = _Vx.

	ST is set equal to the value of _Vx.
*/
func cpuFX018() {
	x := getx()
	soundTimer = _V[x]
	pc += 2
}

/*
	Fx1E - ADD _I, _Vx
	Set _I = _I + _Vx.

	The values of _I and _Vx are added, and the results are stored in _I.
*/
func cpuFX01E() {
	x := getx()
	_I = _I + uint16(_V[x])
	pc += 2
}

/*
	Fx29 - LD F, _Vx
	Set _I = location of sprite for digit _Vx.

	The value of _I is set to the location for the hexadecimal sprite corresponding to
	the value of _Vx.
*/
func cpuFX029() {
	x := getx()
	_I = uint16(_V[x] * 5)
	pc += 2
}

/*
	Fx33 - LD B, _Vx
	Store BCD representation of _Vx in memory locations _I, _I+1, and _I+2.

	The interpreter takes the decimal value of _Vx, and places the hundreds
	digit in memory at location in _I, the tens digit at location _I+1, and the ones digit at location _I+2.
*/
func cpuFX033() {
	value := _V[getx()]
	ones := value % 10
	value = value / 10
	tens := value % 10
	hundreds := value / 10
	memory[_I] = hundreds
	memory[_I+1] = tens
	memory[_I+2] = ones
	pc += 2
}

/*
	Fx55 - LD [_I], _Vx
	Store registers _V0 through _Vx in memory starting at location _I.

	The interpreter copies the values of registers _V0 through _Vx into memory, starting at the address in _I.
*/
func cpuFX055() {
	x := getx()
	var i uint8

	for i = 0; i <= x; i++ {
		memory[_I+uint16(i)] = _V[i]
	}
	pc += 2
}

/*
	Fx65 - LD _Vx, [_I]
	Read registers _V0 through _Vx from memory starting at location _I.

	The interpreter reads values from memory starting at location _I into registers _V0 through _Vx.
*/
func cpuFX065() {
	x := getx()
	var i uint8
	for i = 0; i <= x; i++ {
		_V[i] = memory[_I+uint16(i)]
	}
	pc += 2
}

func xkk() (uint8, uint8) {
	x := getx()
	kk := uint8(opcode & 0x00FF)
	return x, kk
}

func getxy() (uint8, uint8) {
	return getx(), gety()
}

func getx() uint8 {
	return uint8((opcode & 0x0F00) >> 8)
}

func gety() uint8 {
	return uint8((opcode & 0x00F0) >> 4)
}

func printOpcodeErr(section string) {
	fmt.Printf("Unknown opcode [%s]: %#X\n", section, opcode)
	SdlRunning = false
	memoryDump()
}

func memoryDump() {
	fmt.Println("interpreter")
	for i := 0; i < 0x1FF; i++ {
		fmt.Printf("%#x ", memory[i])
	}
	fmt.Println()

	fmt.Println("Program RAM and Work Ram ")
	for i := 0x200; i < 0xFFF; i++ {
		fmt.Printf("%#x ", memory[i])
	}
	fmt.Println()
}
