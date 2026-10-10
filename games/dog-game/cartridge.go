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

const (
	BGBaseAddr      uint16 = 0x0000
	SpriteBaseAddr  uint16 = 0x1000
	TileMapBaseAddr uint16 = 0x2000
	WinBaseAddr     uint16 = 0x3000
)

func (g *DogGame) readRomTile(base uint16) (tile [BytesPerTile]byte) {
	for i := range BytesPerTile {
		tile[i] = g.cpu.Read(base + uint16(i))
	}
	return tile
}

func (g *DogGame) readRomBGTile(id int) (tile [BytesPerTile]byte) {
	base := BGBaseAddr + uint16(id*BytesPerTile)
	return g.readRomTile(base)
}

func (g *DogGame) writeTile(addr uint16, tile [BytesPerTile]byte) {
	for i, data := range tile {
		g.cpu.Write(addr+uint16(i), data)
	}
}

func (g *DogGame) writeTileByID(id int, tile [BytesPerTile]byte) {
	addr := 0x8000 + uint16(id<<4)
	g.writeTile(addr, tile)
}

func (g *DogGame) loadBGTile(srcId, dstId int) {
	tile := g.readRomBGTile(srcId)
	g.writeTileByID(dstId, tile)
}

func (g *DogGame) readRomWindowTile(id int) (tile [BytesPerTile]byte) {
	base := WinBaseAddr + uint16(id*BytesPerTile)
	return g.readRomTile(base)
}

func (g *DogGame) loadWindowTile(srcId, dstId int) {
	tile := g.readRomWindowTile(srcId)
	g.writeTileByID(dstId, tile)
}

func (g *DogGame) readRomSpriteTile(id int) (tile [BytesPerTile]byte) {
	base := SpriteBaseAddr + uint16(id*BytesPerTile)
	return g.readRomTile(base)
}

func (g *DogGame) loadSpriteTile(srcId, dstId int) {
	tile := g.readRomSpriteTile(srcId)
	g.writeTileByID(dstId, tile)
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
			base := TileMapBaseAddr + uint16(y*ScreenTileWidth)
			for x := range ScreenTileWidth {
				tileID := g.Read(base + uint16(x))
				g.setTileMap(0, x, y, tileID)
			}
		}
		// Set window tile
		for tileID := range len(g.window) / BytesPerTile {
			g.loadWindowTile(tileID, tileID|0x80)
		}

		for y := range ScreenTileHeight {
			base := 0x80 + byte(y*ScreenTileWidth)
			for x := range ScreenTileWidth {
				g.setTileMap(1, x, y, base+byte(x))
			}
		}
		// Set sprite tile
		for i := range 4 {
			g.loadSpriteTile(i, i)
		}
	}

	// Set BG tile
	frameOffset := (g.count / 20) & 3
	tileID := [4][2]int{{0, 0}, {1, 2}, {0, 0}, {2, 1}}[frameOffset]
	for i, src := range tileID {
		g.loadBGTile(src, i|0x100)
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
