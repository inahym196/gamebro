package gamebro

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/inahym196/gamebro/games/dog-game"
)

const (
	ScreenWidth  = 16
	ScreenHeight = 16
)

var dogTile = [16 * 16]byte{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 3, 3, 0, 0, 0, 0, 0, 0, 3, 3, 0, 0, 0, 0, 0,
	0, 3, 2, 3, 0, 0, 0, 0, 3, 2, 3, 0, 0, 0, 0, 0,
	0, 3, 1, 2, 3, 0, 0, 3, 2, 1, 3, 0, 0, 0, 0, 0,
	0, 3, 1, 2, 2, 3, 3, 2, 2, 1, 3, 0, 0, 3, 3, 0,
	0, 3, 2, 2, 2, 1, 2, 2, 2, 2, 3, 0, 3, 2, 2, 3,
	0, 3, 2, 2, 2, 1, 2, 2, 2, 2, 2, 3, 3, 3, 2, 3,
	0, 3, 2, 2, 3, 1, 1, 3, 2, 2, 2, 2, 2, 2, 2, 3,
	3, 2, 1, 1, 1, 3, 1, 1, 1, 2, 2, 2, 2, 2, 3, 3,
	0, 3, 2, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 3, 0,
	0, 3, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 3, 0,
	0, 3, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 1, 3, 0,
	0, 0, 3, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 1, 3, 0,
	0, 0, 0, 3, 1, 3, 3, 3, 1, 3, 3, 3, 3, 1, 3, 0,
	0, 0, 0, 3, 3, 0, 0, 0, 3, 3, 0, 0, 3, 3, 3, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

func NewGame(crt *dog.DogGame) *Game {
	return &Game{
		crt: crt,
	}
}

type Game struct {
	crt    *dog.DogGame
	pixels [16 * 16 * 4]byte
}

func (g *Game) Update() error {
	g.crt.Code()
	for i, color := range dogTile {
		g.pixels[i*4+0] = colorMap(color)
		g.pixels[i*4+1] = colorMap(color)
		g.pixels[i*4+2] = colorMap(color)
		g.pixels[i*4+3] = 0xff
	}
	return nil
}

func colorMap(color byte) byte {
	switch color {
	case 3:
		return 0x00
	case 2:
		return 0x55
	case 1:
		return 0xAA
	case 0:
		return 0xFF
	default:
		panic("unexpected color")
	}
}
func (g *Game) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.pixels[:])
}
func (g *Game) Layout(_, _ int) (int, int) { return ScreenWidth, ScreenHeight }
