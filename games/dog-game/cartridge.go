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

const (
	AddrPlayerPosX uint16 = 0xC000 + iota
	AddrPlayerPosY
	AddrPlayerDir
	AddrPlayerCount
)

type object struct {
	PosX, PosY byte
	TileId     byte
	Attr       byte
}

type Direction int

const (
	DirNil = iota
	DirLeft
	DirRight
)

type metaSprite2 struct {
	objs  [2]object
	dir   Direction
	count uint8
}

func NewPlayer(x, y byte, dir Direction, count uint8) metaSprite2 {
	meta2 := metaSprite2{
		objs: [2]object{
			{x, y, 0, 0x10},
			{x + 8, y, 2, 0x10},
		},
		dir:   dir,
		count: count,
	}
	if dir == DirRight {
		meta2.objs[0].TileId, meta2.objs[1].TileId = meta2.objs[1].TileId, meta2.objs[0].TileId
		meta2.objs[0].Attr |= 0x20
		meta2.objs[1].Attr |= 0x20
	}
	return meta2
}

func calcDelta(dirs byte) (dx, dy byte, dir Direction) {
	dir = DirNil
	if dirs&0x08 == 0 {
		dy++
	}
	if dirs&0x04 == 0 {
		dy--
	}
	if dirs&0x02 == 0 {
		dx--
		dir = DirLeft
	}
	if dirs&0x01 == 0 {
		dx++
		dir = DirRight
	}
	return dx, dy, dir
}

func (meta2 metaSprite2) Move(dirs byte) metaSprite2 {
	dx, dy, dir := calcDelta(dirs)
	meta2.objs[0].PosX += dx
	meta2.objs[1].PosX += dx
	meta2.objs[0].PosY += dy
	meta2.objs[1].PosY += dy
	if dx == 0 && dy == 0 {
		meta2.count = 0xFF
	} else {
		meta2.dir = dir
		meta2.count++
	}
	return meta2
}

func (meta2 metaSprite2) Jumped() bool {
	return meta2.count/16&1 == 0
}

func (meta2 metaSprite2) WriteWRAM(cpu *gamebro.CPU) {
	cpu.Write(AddrPlayerPosX, meta2.objs[0].PosX)
	cpu.Write(AddrPlayerPosY, meta2.objs[0].PosY)
	cpu.Write(AddrPlayerDir, byte(meta2.dir))
	cpu.Write(AddrPlayerCount, meta2.count)
}

func writeOAM(cpu *gamebro.CPU, id int, obj object) {
	addr := uint16(0xFE00 + id<<2)
	cpu.Write(addr, obj.PosY)
	cpu.Write(addr+1, obj.PosX)
	cpu.Write(addr+2, obj.TileId)
	cpu.Write(addr+3, obj.Attr)
}

func (meta2 metaSprite2) WriteOAM(cpu *gamebro.CPU) {
	writeOAM(cpu, 0, meta2.objs[0])
	writeOAM(cpu, 1, meta2.objs[1])
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
		dog := NewPlayer(8*11, 16*5, DirLeft, 0)
		dog.WriteWRAM(g.cpu)
	}

	// Set BG tile
	frameOffset := (g.count / 20) & 3
	tileID := [4][2]int{{0, 0}, {1, 2}, {0, 0}, {2, 1}}[frameOffset]
	for i, src := range tileID {
		g.loadBGTile(src, i|0x100)
	}
	g.count++

	dog := NewPlayer(
		g.cpu.Read(AddrPlayerPosX),
		g.cpu.Read(AddrPlayerPosY),
		Direction(g.cpu.Read(AddrPlayerDir)),
		g.cpu.Read(AddrPlayerCount),
	)

	g.cpu.Write(0xFF00, 0xEF)
	dog = dog.Move(g.cpu.Read(0xFF00))
	if dog.Jumped() {
		g.loadSpriteTile(4, 1)
		g.loadSpriteTile(5, 3)
	} else {
		g.loadSpriteTile(1, 1)
		g.loadSpriteTile(3, 3)
	}
	dog.WriteWRAM(g.cpu)
	dog.WriteOAM(g.cpu)
}
