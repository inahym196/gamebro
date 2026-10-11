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

type object struct {
	PosX, PosY byte
	TileId     byte
	Attr       byte
}

// ObjectはWRAMから取得するべき
func (g *DogGame) ReadObject(id int) object {
	addr := uint16(0xFE00 + id<<2)
	return object{
		PosY:   g.cpu.Read(addr),
		PosX:   g.cpu.Read(addr + 1),
		TileId: g.cpu.Read(addr + 2),
		Attr:   g.cpu.Read(addr + 3),
	}
}

func (g *DogGame) writeOAM(id int, field OAMField, data byte) {
	addr := uint16(0xFE00 + id<<2 + int(field))
	g.cpu.Write(addr, data)
}

func (g *DogGame) setDog(x, y, tileID, attr byte) {
	g.writeOAM(0, OAMFieldPosX, x)
	g.writeOAM(0, OAMFieldPosY, y)
	g.writeOAM(0, OAMFieldTileId, tileID)
	g.writeOAM(0, OAMFieldAttr, attr)
	g.writeOAM(1, OAMFieldPosX, x+gamebro.TileSize)
	g.writeOAM(1, OAMFieldPosY, y)
	g.writeOAM(1, OAMFieldTileId, tileID+2)
	g.writeOAM(1, OAMFieldAttr, attr)
}

func (g *DogGame) moveDog(dirs byte) {

	dog1 := g.ReadObject(0)
	dog2 := g.ReadObject(1)
	if dirs&0x08 == 0 {
		dog1.PosY++
		dog2.PosY++
	}
	if dirs&0x04 == 0 {
		dog1.PosY--
		dog2.PosY--
	}
	g.writeOAM(0, OAMFieldPosY, dog1.PosY)
	g.writeOAM(1, OAMFieldPosY, dog2.PosY)

	if dirs&0x02 == 0 {
		dog1.PosX--
		dog2.PosX--
	}
	if dirs&0x01 == 0 {
		dog1.PosX++
		dog2.PosX++
	}
	g.writeOAM(0, OAMFieldPosX, dog1.PosX)
	g.writeOAM(1, OAMFieldPosX, dog2.PosX)
}

func (g *DogGame) setTileMap(mapID, x, y int, tileID byte) {
	base := [2]uint16{0x9800, 0x9C00}[mapID]
	tmapAddr := base + uint16(x|y<<5)
	g.cpu.Write(tmapAddr, tileID)
}

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
		// Set OAM
		g.setDog(8*11, 16*5, 0, 0x10)
	}

	// Set BG tile
	frameOffset := (g.count / 20) & 3
	tileID := [4][2]int{{0, 0}, {1, 2}, {0, 0}, {2, 1}}[frameOffset]
	for i, src := range tileID {
		g.loadBGTile(src, i|0x100)
	}
	g.count++

	g.cpu.Write(0xFF00, 0xEF)
	dirs := g.cpu.Read(0xFF00)
	g.moveDog(dirs)
}
