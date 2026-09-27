package gamebro

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 16
	ScreenHeight = 16
)

type cartridge interface {
	Code(cpu *CPU)
}

func NewGame(crt cartridge) *Game {
	return &Game{
		crt: crt,
		cpu: NewCPU(crt),
	}
}

type Game struct {
	crt    cartridge
	pixels [16 * 16 * 4]byte
	cpu    *CPU
}

func (g *Game) readTileRow(id byte, offsetY int) (lo, hi byte) {
	base := 0x8000 + uint16(id<<4)
	return g.cpu.Read(base + uint16(offsetY<<1)), g.cpu.Read(base + uint16(offsetY<<1) + 1)
}

func (g *Game) fetchBGWindowTileRow(tileX, tileY, offsetY, mapID int) (lo, hi byte) {
	addr := [2]uint16{0x9800, 0x9C00}[mapID] + uint16(tileY<<5|tileX)
	tileID := g.cpu.Read(addr)
	return g.readTileRow(tileID, offsetY)
}

func (g *Game) fetchBGTileRow(tileX, tileY, offsetY int) (lo, hi byte) {
	return g.fetchBGWindowTileRow(tileX, tileY, offsetY, 0)
}

func colorMap(colorID int, palette uint8) byte {
	cmap := [4]uint8{0xFF, 0xAA, 0x55, 0x00}
	paletteID := int((palette >> (colorID << 1)) & 0b11)
	return cmap[paletteID]
}

func toColorIDs(lo, hi byte) [8]int {
	var cNums [8]int
	for x := range 8 {
		shift := 7 - x
		_lo := (lo >> shift) & 1
		_hi := (hi >> shift) & 1
		cNums[x] = int(_hi<<1 | _lo)
	}
	return cNums
}

func (g *Game) renderScanline(ly int) {
	tileY := ly / 8
	offsetY := ly % 8
	for tileX := range 16 / 8 {
		baseX := tileX * 8
		lo, hi := g.fetchBGTileRow(tileX, tileY, offsetY)
		for i, cid := range toColorIDs(lo, hi) {
			c := colorMap(cid, 0b11100100)
			idx := (ly*16 + baseX + i) * 4
			g.pixels[idx] = c
			g.pixels[idx+1] = c
			g.pixels[idx+2] = c
			g.pixels[idx+3] = 255
		}
	}
}

func (g *Game) Update() error {
	g.cpu.Step()
	for ly := range 16 {
		g.renderScanline(ly)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.pixels[:])
}
func (g *Game) Layout(_, _ int) (int, int) { return ScreenWidth, ScreenHeight }
