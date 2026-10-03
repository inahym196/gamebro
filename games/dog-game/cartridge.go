package dog

import (
	_ "embed"
	"log/slog"

	"github.com/inahym196/gamebro"
)

const (
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
func readRomTile(cpu *gamebro.CPU, id int) (tile [gamebro.BytesPerTile]byte) {
	base := id * gamebro.BytesPerTile
	for i := range gamebro.BytesPerTile {
		addr := uint16(base + i)
		tile[i] = cpu.Read(addr)
	}
	return tile
}

func loadTile(cpu *gamebro.CPU, srcId int, dstId int) {
	tile := readRomTile(cpu, srcId)
	cpu.WriteTile(dstId, tile)
}

func calcTileMapAddr(mapID, x, y int) uint16 {
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
	return base + uint16(x|y<<5)
}

func (g *DogGame) Code(cpu *gamebro.CPU) {
	if g.init {
		g.init = false
		g.cpu = cpu
		for tileID := range len(g.bg) / gamebro.BytesPerTile {
			loadTile(cpu, tileID, tileID)
		}
		for y := range ScreenTileHeight {
			base := 0x2000 + y*ScreenTileWidth
			for x := range ScreenTileWidth {
				tileID := g.Read(uint16(base + x))
				tmapAddr := calcTileMapAddr(0, x, y)
				cpu.Write(tmapAddr, tileID)
			}
		}
	}
}
