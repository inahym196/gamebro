package dog

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed assets/sprite.rom
var spriteData []byte

type DogGame struct {
	sprite []byte
}

func NewDogGame() *DogGame {
	return &DogGame{spriteData}
}

func (g *DogGame) fetchRomTile(id int) (tile [16]byte) {
	base := id * 16
	for i := range 16 {
		tile[i] = g.sprite[base+i]
	}
	return tile
}

func (g *DogGame) Code() {
	for i := range 4 {
		fmt.Printf("tileID=%d: ", i)
		tile := g.fetchRomTile(i)
		for _, b := range tile {
			fmt.Printf("%02x ", b)
		}
		fmt.Println()
	}
	os.Exit(0)
	//return mmu.SetTile(0, tile)
}
