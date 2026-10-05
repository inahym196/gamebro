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
	case addr >= 0x4000 && addr < 0x5000:
		return g.sprite[addr-0x4000]
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

func readRomSpriteTile(cpu *gamebro.CPU, id int) (tile [BytesPerTile]byte) {
	base := uint16(0x4000 + id*BytesPerTile)
	return readRomTile(cpu, base)
}

func loadSpriteTile(cpu *gamebro.CPU, srcId, dstId int) {
	tile := readRomSpriteTile(cpu, srcId)
	writeTileByID(cpu, dstId, tile)
}

type OAMField int

const (
	OAMFieldPosY OAMField = iota
	OAMFieldPosX
	OAMFieldTileId
	OAMFieldAttr
)

func WriteOAM(cpu *gamebro.CPU, id int, field OAMField, data byte) {
	addr := uint16(0xFE00 + id<<2 + int(field))
	cpu.Write(addr, data)
}

func (g *DogGame) setTileMap(mapID, x, y int, tileID byte) {
	base := [2]uint16{0x9800, 0x9C00}[mapID]
	tmapAddr := base + uint16(x|y<<5)
	g.cpu.Write(tmapAddr, tileID)
}

// TODO: ReadWriter依存にできないか？
func (g *DogGame) Code(cpu *gamebro.CPU) {
	if g.init {
		g.init = false
		g.cpu = cpu

		// Set bg tilemap
		for y := range ScreenTileHeight {
			base := 0x2000 + y*ScreenTileWidth
			for x := range ScreenTileWidth {
				tileID := g.Read(uint16(base + x))
				g.setTileMap(0, x, y, tileID)
			}
		}
		// Set window tile
		for tileID := range len(g.window) / BytesPerTile {
			loadWindowTile(cpu, tileID, tileID|0x80)
		}

		for y := range ScreenTileHeight {
			base := 0x80 + byte(y*ScreenTileWidth)
			for x := range ScreenTileWidth {
				g.setTileMap(1, x, y, base+byte(x))
			}
		}
		// Set sprite tile
		for i := range 4 {
			loadSpriteTile(cpu, i, i)
		}
	}

	// Set BG tile
	frameOffset := (g.count / 20) & 3
	tileID := [4][2]int{{0, 0}, {1, 2}, {0, 0}, {2, 1}}[frameOffset]
	for i, src := range tileID {
		loadBGTile(cpu, src, i|0x100)
	}
	g.count++

	// Set OAM
	WriteOAM(cpu, 0, OAMFieldPosY, 16*5)
	WriteOAM(cpu, 0, OAMFieldPosX, 8*11)
	WriteOAM(cpu, 0, OAMFieldTileId, 0)
	WriteOAM(cpu, 0, OAMFieldAttr, 0x10)
	WriteOAM(cpu, 1, OAMFieldPosY, 16*5)
	WriteOAM(cpu, 1, OAMFieldPosX, 8*12)
	WriteOAM(cpu, 1, OAMFieldTileId, 2)
	WriteOAM(cpu, 1, OAMFieldAttr, 0x10)
}
