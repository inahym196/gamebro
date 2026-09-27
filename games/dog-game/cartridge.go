package dog

import (
	_ "embed"
	"log/slog"

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
		tile[i] = g.bg[base+i]
	}
	return tile
}

func (g *DogGame) writeTile(id int, tile [16]byte) {
	if id < 0 || id >= 384 {
		slog.Error("write tile: out of range", "tileID", id)
	}
	base := 0x8000 + uint16(id<<4)
	for i, data := range tile {
		g.cpu.Write(base+uint16(i), data)
	}
}

func (g *DogGame) writeTileMap(mapID, x, y int, data byte) {
	if mapID < 0 || mapID >= 2 {
		slog.Error("write tileMap: out of range", "mapID", mapID)
	}
	if x < 0 || x >= 32 {
		slog.Error("write tileMap: out of range", "x", x)
	}
	if y < 0 || x >= 32 {
		slog.Error("write tileMap: out of range", "y", y)
	}
	base := [2]uint16{0x9800, 0x9C00}[mapID]
	addr := base + uint16(x|y<<5)
	g.cpu.Write(addr, data)
}

func (g *DogGame) Code(cpu *gamebro.CPU) {
	if g.init {
		g.init = false
		g.cpu = cpu
		for tileID := range len(g.bg) / 16 {
			tile := g.fetchRomTile(tileID)
			g.writeTile(tileID, tile)
		}
		for y := range 18 {
			for x := range 20 {
				tileID := g.tmap[x+y*20]
				g.writeTileMap(0, x, y, tileID)
			}
		}
	}
}
