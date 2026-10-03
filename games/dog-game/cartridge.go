package dog

import (
	_ "embed"
	"log/slog"

	"github.com/inahym196/gamebro"
)

const (
	BytesPerTile     = gamebro.TileSize * 2
	ScreenTileWidth  = gamebro.ScreenWidth / gamebro.TileSize
	ScreenTileHeight = gamebro.ScreenHeight / gamebro.TileSize
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

func (g *DogGame) Read(addr uint16) byte {
	switch {
	case addr < 0x1000:
		return g.bg[addr]
	case addr >= 0x1000 && addr < 0x2000:
		return g.sprite[addr-0x1000]
	case addr >= 0x2000 && addr < 0x3000:
		return g.tmap[addr-0x2000]
	default:
		slog.Warn("not impl yet", "addr", addr)
		return 0xFF
	}
}

func (g *DogGame) readTile(id int) (tile [BytesPerTile]byte) {
	base := id * BytesPerTile
	for i := range BytesPerTile {
		addr := uint16(base + i)
		tile[i] = g.cpu.Read(addr)
	}
	return tile
}

func (g *DogGame) writeTile(id int, tile [BytesPerTile]byte) {
	if id < 0 || id >= 384 {
		slog.Error("write tile: out of range", "tileID", id)
	}
	base := 0x8000 + uint16(id<<4)
	for i, data := range tile {
		addr := base + uint16(i)
		g.cpu.Write(addr, data)
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
		for tileID := range len(g.bg) / gamebro.BytesPerTile {
			g.cpu.LoadTile(tileID, tileID)
		}
		for y := range ScreenTileHeight {
			base := y * ScreenTileWidth
			for x := range ScreenTileWidth {
				addr := uint16(0x2000 + base + x)
				tileID := g.Read(addr)
				g.writeTileMap(0, x, y, tileID)
			}
		}
	}
}
