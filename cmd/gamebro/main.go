package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/inahym196/gamebro"
	"github.com/inahym196/gamebro/games/dog-game"
)

func init() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))
	slog.SetDefault(logger)
}

func main() {
	crt := dog.NewDogGame()

	var game ebiten.Game
	g := gamebro.NewGame(crt)

	var debugMode bool = true
	if debugMode {
		var err error
		game, err = gamebro.NewDebugGame(g)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		game = g
	}
	ebiten.SetWindowSize(gamebro.ScreenWidth*6, gamebro.ScreenHeight*6)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
