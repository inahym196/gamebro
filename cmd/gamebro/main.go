package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/inahym196/gamebro"
	"github.com/inahym196/gamebro/games/dog-game"
)

func main() {
	crt := dog.NewDogGame()
	game := gamebro.NewGame(crt)
	ebiten.SetWindowSize(gamebro.ScreenWidth*20, gamebro.ScreenWidth*20)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
