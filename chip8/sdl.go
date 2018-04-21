/*
	SDL.go:
	This file contains all SDL related code, including: window creation,
	rendering and keyboard handling.
*/
package chip8

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

//SdlRunning is a flag for if the use exited the game or not
var SdlRunning bool

const gameWidth = 64
const gameHeight = 32

var renderer *sdl.Renderer

//SdlSetupGraphics Sets up graphics using SDL2 library
func SdlSetupGraphics() {
	SdlRunning = true

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Chip8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	r, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	renderer = r

	renderer.SetLogicalSize(gameWidth, gameHeight)
	SdlDrawGraphics()
	for SdlRunning {
		EmulateCycle()
		if DrawFlag {
			SdlDrawGraphics()
		}

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

			if event != nil {
				switch t := event.(type) {
				case *sdl.QuitEvent:
					fmt.Println("Quit")
					SdlRunning = false
				case *sdl.KeyboardEvent:
					if t.State == sdl.PRESSED {
						keyPressed = true
					}
					if debug {
						fmt.Println("keyPressed")
					}
					switch t.Keysym.Sym {
					case sdl.K_1:
						SetKeys(0, t.State)
					case sdl.K_2:
						SetKeys(1, t.State)

					case sdl.K_3:
						SetKeys(2, t.State)

					case sdl.K_4:
						SetKeys(3, t.State)

					case sdl.K_q:
						SetKeys(4, t.State)

					case sdl.K_w:
						SetKeys(5, t.State)

					case sdl.K_e:
						SetKeys(6, t.State)

					case sdl.K_r:
						SetKeys(7, t.State)

					case sdl.K_a:
						SetKeys(8, t.State)

					case sdl.K_s:
						SetKeys(9, t.State)

					case sdl.K_d:
						SetKeys(10, t.State)

					case sdl.K_f:
						SetKeys(11, t.State)

					case sdl.K_z:
						SetKeys(12, t.State)

					case sdl.K_x:
						SetKeys(13, t.State)

					case sdl.K_c:
						SetKeys(14, t.State)

					case sdl.K_v:
						SetKeys(15, t.State)

					}
				}
			}
		}
	}

}

//SdlDrawGraphics  Draws graphics using SDL 2 library.
func SdlDrawGraphics() {
	renderer.Clear()
	for i := 0; i < len(gfx); i++ {
		x := i % gameWidth
		y := i / gameWidth
		var drawColor uint8
		if gfx[i] != 0 {
			drawColor = 255
		} else {
			drawColor = 0
		}

		renderer.SetDrawColor(drawColor, drawColor, drawColor, sdl.ALPHA_OPAQUE)
		renderer.DrawPoint(int32(x), int32(y))
	}
	renderer.Present()
}
