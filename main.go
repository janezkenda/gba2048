package main

import (
	"image/color"
	"strconv"
	"unsafe"

	"machine"
	"runtime/volatile"
	"tinygo.org/x/tinyfont"

	"github.com/janezkenda/gba2048/game"
)

var display = machine.Display

type Key uint8

const (
	A Key = iota
	B
	Select
	Start
	Right
	Left
	Up
	Down
	R
	L
	None
)

var keyNames = map[Key]string{
	A:      "A",
	B:      "B",
	Select: "Select",
	Start:  "Start",
	Right:  "Right",
	Left:   "Left",
	Up:     "Up",
	Down:   "Down",
	R:      "R",
	L:      "L",
}

var keys = []Key{A, B, Select, Start, Right, Left, Up, Down, R, L}

var (
	colorFull = color.RGBA{66, 135, 245, 0}
	colorGrey = color.RGBA{100, 100, 100, 0}
)

func showRect(x int16, y int16, w int16, h int16, c color.RGBA) {
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			display.SetPixel(i, j, c)
		}
	}
}

func printField(g game.Game) {
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			el := g.Board[game.Point{X: x, Y: y}]

			if el.Value == 0 {
				showRect(42+int16(x)*40, 2+int16(y)*40, 37, 37, colorGrey)
			} else {
				clr := colorFull
				if c, ok := colors[el.Value]; ok {
					clr = *c
				}

				showRect(42+int16(x)*40, 2+int16(y)*40, 37, 37, clr)
				tinyfont.WriteLine(display, &tinyfont.Picopixel, 40+int16(x)*40+15, int16(y)*40+20, []byte(strconv.Itoa(el.Value)), color.RGBA{R: 255, G: 255, B: 255})
			}
		}
	}
}

var directionMap = map[Key]game.Direction{
	Up:    game.Up,
	Down:  game.Down,
	Left:  game.Left,
	Right: game.Right,
}

type inputReg uint16

func (i inputReg) IsPressed(k Key) bool {
	return (i & (1 << uint16(k))) == 0
}

func (i inputReg) GetPressed() Key {
	for _, k := range keys {
		if i.IsPressed(k) {
			return k
		}
	}
	return None
}

func to16bitColor(c color.RGBA) color.RGBA {
	return color.RGBA{
		R: uint8((int(c.R)*249 + 1014) >> 11),
		G: uint8((int(c.G)*253 + 505) >> 10),
		B: uint8((int(c.B)*249 + 1014) >> 11),
	}
}

var (
	colorBlueSapphire   = to16bitColor(color.RGBA{R: 5, G: 102, B: 141})
	colorMetalicSeaweed = to16bitColor(color.RGBA{R: 2, G: 128, B: 144})
	colorPersianGreen   = to16bitColor(color.RGBA{R: 0, G: 168, B: 150})
	colorMountainMeadow = to16bitColor(color.RGBA{R: 2, G: 195, B: 154})
	colorPaleSpringBud  = to16bitColor(color.RGBA{R: 240, G: 243, B: 189})
)

var colors = map[int]*color.RGBA{
	2:  &colorMountainMeadow,
	4:  &colorPersianGreen,
	8:  &colorMetalicSeaweed,
	16: &colorBlueSapphire,
}

func main() {
	display.Configure()
	showRect(0, 0, 240, 160, colorPaleSpringBud)

	keyinput := (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000130)))

	g := game.NewGame()
	printField(g)

	for {
		k := inputReg(keyinput.Get()).GetPressed()
		if k == None {
			continue
		}

		showRect(0, 0, 40, 20, colorPaleSpringBud)
		tinyfont.WriteLine(display, &tinyfont.Picopixel, 2, 5, []byte(strconv.Itoa(g.Score)), color.RGBA{R: 0, G: 0, B: 0})

		showRect(0, 140, 40, 20, colorPaleSpringBud)
		tinyfont.WriteLine(display, &tinyfont.Picopixel, 2, 150, []byte(keyNames[k]), color.RGBA{R: 0, G: 0, B: 0})

		switch k {
		case Up:
			g.Move(game.Up)
		case Down:
			g.Move(game.Down)
		case Left:
			g.Move(game.Left)
		case Right:
			g.Move(game.Right)
		default:
			continue
		}

		if g.Moved {
			g.AddNewNumber()
		}

		printField(g)
		g.Moved = false

		if g.End {
			g.GameOver()
			break
		}

		if g.Win {
			g.GameWon()
			break
		}
	}
}
