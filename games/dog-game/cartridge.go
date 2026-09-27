package dog

import (
	_ "embed"

	"github.com/inahym196/gamebro"
)

//go:embed assets/sprite.rom
var spriteData []byte

type DogGame struct {
	init   bool
	sprite []byte
}

func NewDogGame() *DogGame {
	return &DogGame{init: true, sprite: spriteData}
}

func (g *DogGame) fetchRomTile(id int) (tile [16]byte) {
	base := id * 16
	for i := range 16 {
		tile[i] = g.sprite[base+i]
	}
	return tile
}

func (g *DogGame) Code(mmu *gamebro.MMU) {
	if g.init {
		g.init = false
		for tileID := range 4 {
			tile := g.fetchRomTile(tileID)
			for i, b := range tile {
				mmu.WriteTile(tileID*16+i, b)
			}
		}
	}
}
