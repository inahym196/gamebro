package dog

import (
	_ "embed"

	"github.com/inahym196/gamebro"
)

//go:embed assets/sprite.rom
var spriteData []byte

//go:embed assets/bg.rom
var bgData []byte

//go:embed assets/tilemap.rom
var tmapData []byte

type DogGame struct {
	init   bool
	cpu    *gamebro.CPU
	bg     []byte
	sprite []byte
	tmap   []byte
}

func NewDogGame() *DogGame {
	return &DogGame{
		init:   true,
		bg:     bgData,
		sprite: spriteData,
		tmap:   tmapData,
	}
}

func (g *DogGame) fetchRomTile(id int) (tile [16]byte) {
	base := id * 16
	for i := range 16 {
		tile[i] = g.sprite[base+i]
	}
	return tile
}

func (g *DogGame) writeTile(addr uint16, data byte) {
	g.cpu.Write(addr+0x8000, data)
}

func (g *DogGame) Code(cpu *gamebro.CPU) {
	if g.init {
		g.init = false
		g.cpu = cpu
		for tileID := range 4 {
			tile := g.fetchRomTile(tileID)
			base := uint16(tileID << 4)
			for i, b := range tile {
				g.writeTile(base+uint16(i), b)
			}
		}
	}
}
