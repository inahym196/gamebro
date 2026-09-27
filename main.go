package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 160
	screenHeight = 144
)

type Game struct{}

func (g *Game) Update() error { return nil }
func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "hello world")
}
func (g *Game) Layout(_, _ int) (int, int) { return screenWidth, screenHeight }

func main() {
	ebiten.SetWindowSize(screenWidth*5, screenHeight*5)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
