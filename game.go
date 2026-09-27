package gamebro

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 16
	ScreenHeight = 16
)

type cartridge interface {
	Code(mmu *MMU)
}

func NewGame(crt cartridge) *Game {
	mmu := &MMU{}
	return &Game{
		crt: crt,
		mmu: mmu,
	}
}

type Game struct {
	crt    cartridge
	pixels [16 * 16 * 4]byte
	mmu    *MMU
}

func (g *Game) readTile(addr uint16) byte {
	return g.mmu.Read(0x8000 + addr)
}

func (g *Game) readTileRow(tileX, tileY, offsetY int) (lo, hi byte) {
	if tileX < 0 || tileX >= 2 {
		log.Fatalf("tileX out of range: %d", tileX)
	}
	if tileY < 0 || tileY >= 2 {
		log.Fatalf("tileY out of range: %d", tileY)
	}
	baseAddr := uint16(tileX*16 + tileY*32 + offsetY*2)
	return g.readTile(baseAddr), g.readTile(baseAddr + 1)
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
		for i, cid := range toColorIDs(g.readTileRow(tileX, tileY, offsetY)) {
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
	g.crt.Code(g.mmu)
	for ly := range 16 {
		g.renderScanline(ly)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.pixels[:])
}
func (g *Game) Layout(_, _ int) (int, int) { return ScreenWidth, ScreenHeight }
