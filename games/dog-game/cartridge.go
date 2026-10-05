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

//go:generate go run ./cmd/gen2bpp/main.go assets/src/sprite assets/sprite.rom
//go:embed assets/sprite.rom
var spriteData []byte

//go:generate go run ./cmd/gen2bpp/main.go assets/src/bg assets/bg.rom
//go:embed assets/bg.rom
var bgData []byte

//go:generate go run ./cmd/gen2bpp/main.go assets/src/window assets/window.rom
//go:embed assets/window.rom
var windowData []byte

//go:generate go run ./cmd/gentilemap/main.go assets/tilemap.rom
//go:embed assets/tilemap.rom
var tmapData []byte

type DogGame struct {
	init   bool
	cpu    *gamebro.CPU
	bg     []byte
	sprite []byte
	tmap   []byte
	window []byte
	count  int
}

func NewDogGame() *DogGame {
	return &DogGame{
		init:   true,
		bg:     bgData,
		sprite: spriteData,
		tmap:   tmapData,
		window: windowData,
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
	case addr >= 0x3000 && addr < 0x4000:
		return g.window[addr-0x3000]
	default:
		slog.Warn("not impl yet", "addr", addr)
		return 0xFF
	}
}

func readRomTile(cpu *gamebro.CPU, base uint16) (tile [BytesPerTile]byte) {
	for i := range BytesPerTile {
		tile[i] = cpu.Read(base + uint16(i))
	}
	return tile
}

func readRomBGTile(cpu *gamebro.CPU, id int) (tile [BytesPerTile]byte) {
	base := uint16(id * BytesPerTile)
	return readRomTile(cpu, base)
}

func writeTile(cpu *gamebro.CPU, addr uint16, tile [BytesPerTile]byte) {
	for i, data := range tile {
		cpu.Write(addr+uint16(i), data)
	}
}

func writeTileByID(cpu *gamebro.CPU, id int, tile [BytesPerTile]byte) {
	addr := 0x8000 + uint16(id<<4)
	writeTile(cpu, addr, tile)
}

func loadBGTile(cpu *gamebro.CPU, srcId, dstId int) {
	tile := readRomBGTile(cpu, srcId)
	writeTileByID(cpu, dstId, tile)
}

func readRomWindowTile(cpu *gamebro.CPU, id int) (tile [BytesPerTile]byte) {
	base := uint16(0x3000 + id*BytesPerTile)
	return readRomTile(cpu, base)
}

func loadWindowTile(cpu *gamebro.CPU, srcId, dstId int) {
	tile := readRomWindowTile(cpu, srcId)
	writeTileByID(cpu, dstId, tile)
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

		// Set bg tilemap
		for y := range ScreenTileHeight {
			base := 0x2000 + y*ScreenTileWidth
			for x := range ScreenTileWidth {
				tileID := g.Read(uint16(base + x))
				tmapAddr := calcTileMapAddr(0, x, y)
				cpu.Write(tmapAddr, tileID)
			}
		}
		// Set window tile
		for tileID := range len(g.window) / BytesPerTile {
			loadWindowTile(cpu, tileID, tileID|0x80)
		}

		// Set window tilemap
		tmapIndex := byte(0x80)
		for y := range ScreenTileHeight {
			for x := range ScreenTileWidth {
				tmapAddr := calcTileMapAddr(1, x, y)
				cpu.Write(tmapAddr, tmapIndex)
				tmapIndex++
			}
		}
	}

	// Set BG tile
	frameOffset := (g.count / 20) & 3
	tileID := [4][2]int{{0, 0}, {1, 2}, {0, 0}, {2, 1}}[frameOffset]
	for i, src := range tileID {
		loadBGTile(cpu, src, i|0x100)
	}
	g.count++
}
