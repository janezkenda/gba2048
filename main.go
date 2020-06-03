package main

import (
	"image/color"
	"strconv"

	"github.com/janezkenda/tinygba/display"
	"github.com/janezkenda/tinygba/keypad"
	"tinygo.org/x/tinyfont"

	"github.com/janezkenda/gba2048/game"
)

func log2(x int) int {
	var res int

	for x > 1 {
		x >>= 1
		res++
	}
	return res
}

var gradient = append(
	display.Gradient(display.ClrYellow, 0x6007, 5),
	append(
		display.Gradient(0x6007, display.ClrRed, 5),
		display.Gradient(display.ClrRed, display.ClrBlack, 10)...,
	)...,
)

func printField(g game.Game, m3 *display.Mode3) {
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			el := g.Board[x][y]

			size := 37
			xbase, ybase := 43, 3

			x0 := xbase + x*size + x*2
			y0 := ybase + y*size + y*2

			if el.Value == 0 {
				m3.Rect(x0, y0, x0+size+2, y0+size+2, display.ClrWhite)
				m3.Rect(x0+16, y0+16, x0+size-15, y0+size-15, display.ClrBlack)
			} else {
				clr := gradient[log2(el.Value)-1]
				m3.Rect(x0, y0, x0+size, y0+size, clr)
				m3.Frame(x0, y0, x0+size, y0+size, display.ClrBlack)
				m3.Line(x0+size, y0+1, x0+size, y0+size, display.ClrBlack)
				m3.Line(x0+1, y0+size, x0+size, y0+size, display.ClrBlack)
				tinyfont.WriteLine(m3, &tinyfont.Picopixel, 40+int16(x)*40+15, int16(y)*40+20, []byte(strconv.Itoa(el.Value)), color.RGBA{R: 255, G: 255, B: 255})
			}
		}
	}
}

var m3 = display.NewMode3()

func main() {
	// m3 := display.NewMode3()
	m3.Fill(display.ClrWhite)
	m3.Frame(41, 1, 200, 160, display.ClrBlack)

	g := game.New()
	g.AddNewNumber()
	g.AddNewNumber()
	printField(g, m3)

	for {
		keypad.Poll()
		if keypad.Up.Hit() {
			g.Move(game.Up)
		} else if keypad.Down.Hit() {
			g.Move(game.Down)
		} else if keypad.Left.Hit() {
			g.Move(game.Left)
		} else if keypad.Right.Hit() {
			g.Move(game.Right)
		} else {
			continue
		}

		if g.Moved {
			g.AddNewNumber()
		}

		m3.Rect(0, 0, 40, 20, display.ClrWhite)
		tinyfont.WriteLine(m3, &tinyfont.Picopixel, 2, 5, []byte(strconv.Itoa(g.Score)), color.RGBA{R: 0, G: 0, B: 0})
		printField(g, m3)

		g.Moved = false

		if g.End || g.Win {
			break
		}
	}
}
