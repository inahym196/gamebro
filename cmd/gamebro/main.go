package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/inahym196/gamebro"
)

func main() {
	ebiten.SetWindowSize(gamebro.ScreenWidth*20, gamebro.ScreenWidth*20)
	if err := ebiten.RunGame(&gamebro.Game{}); err != nil {
		log.Fatal(err)
	}
}
