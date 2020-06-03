package game

import "math/rand"

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type (
	Entry struct {
		Value   int
		Blocked bool
	}

	Game struct {
		Board [4][4]*Entry
		Score int
		Moved bool
		End   bool
		Win   bool
	}
)

func New() Game {
	g := Game{
		Board: [4][4]*Entry{},
	}

	for x, l := range g.Board {
		for y := range l {
			g.Board[x][y] = &Entry{}
		}
	}

	return g
}

func (g *Game) AddNewNumber() {
	type point struct {
		x, y int
	}

	var zeroes []point

	for x, l := range g.Board {
		for y := range l {
			if g.Board[x][y].Value == 0 {
				zeroes = append(zeroes, point{x, y})
			}
		}
	}

	if len(zeroes) == 0 {
		return
	}

	num := 2
	if rand.Intn(10) == 0 {
		num = 4
	}

	pt := zeroes[rand.Intn(len(zeroes))]

	g.Board[pt.x][pt.y].Value = num

	if g.canMove() {
		return
	}

	g.End = true
}

func (g *Game) canMove() bool {
	for x, l := range g.Board {
		for y := range l {
			if g.Board[x][y].Value == 0 {
				return true
			}
		}
	}

	for x, l := range g.Board {
		for y := range l {
			val := g.Board[x][y].Value
			if g.testAdd(x+1, y, val) ||
				g.testAdd(x-1, y, val) ||
				g.testAdd(x, y+1, val) ||
				g.testAdd(x, y-1, val) {
				return true
			}
		}
	}

	return false
}

func (g *Game) testAdd(x, y, value int) bool {
	if x < 0 || x > 3 || y < 0 || y > 3 {
		return false
	}

	return g.Board[x][y].Value == value
}

func (g *Game) mv(pt, other *Entry) {
	if pt.Value != 0 && other.Value == pt.Value &&
		!pt.Blocked && !other.Blocked {
		pt.Value = 0

		other.Value *= 2
		other.Blocked = true

		if other.Value == 2048 {
			g.Win = true
		}

		g.Score += other.Value
		g.Moved = true
	} else if other.Value == 0 && pt.Value != 0 {
		other.Value = pt.Value
		pt.Value = 0

		g.Moved = true
	}
}

func (g *Game) Move(d Direction) {
	switch d {
	case Up:
		for x := 0; x < 4; x++ {
			for y := 1; y < 4; y++ {
				if g.Board[x][y].Value == 0 {
					continue
				}

				for i := y; i > 0; i-- {
					pt := g.Board[x][i]
					other := g.Board[x][i-1]

					g.mv(pt, other)
				}
			}
		}
	case Down:
		for x := 0; x < 4; x++ {
			for y := 2; y >= 0; y-- {
				if g.Board[x][y].Value == 0 {
					continue
				}

				for i := y; i < 3; i++ {
					pt := g.Board[x][i]
					other := g.Board[x][i+1]

					g.mv(pt, other)
				}
			}
		}
	case Left:
		for y := 0; y < 4; y++ {
			for x := 1; x < 4; x++ {
				if g.Board[x][y].Value == 0 {
					continue
				}

				for i := x; i > 0; i-- {
					pt := g.Board[i][y]
					other := g.Board[i-1][y]

					g.mv(pt, other)
				}
			}
		}
	case Right:
		for y := 0; y < 4; y++ {
			for x := 2; x >= 0; x-- {
				if g.Board[x][y].Value == 0 {
					continue
				}

				for i := x; i < 3; i++ {
					pt := g.Board[i][y]
					other := g.Board[i+1][y]

					g.mv(pt, other)
				}
			}
		}
	}

	for x, l := range g.Board {
		for y := range l {
			g.Board[x][y].Blocked = false
		}
	}
}
