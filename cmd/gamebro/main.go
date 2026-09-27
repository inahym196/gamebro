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
	game := gamebro.NewGame(crt)
	ebiten.SetWindowSize(gamebro.ScreenWidth*20, gamebro.ScreenWidth*20)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
