package gamebro

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 160
	ScreenHeight = 144
)

func NewGame(crt cartridge) *Game {
	mmu := NewMMU(crt)
	return &Game{
		pbuf: NewPixelBuffer(ScreenWidth, ScreenHeight),
		cpu:  NewCPU(mmu),
		ppu:  NewPPU(mmu),
	}
}

type Game struct {
	pbuf *PixelBuffer
	cpu  *CPU
	ppu  *PPU
}

func (g *Game) Update() error {
	g.cpu.Step()
	for ly := range ScreenHeight {
		g.ppu.RenderScanline(ly)
		g.pbuf.SetScanlineGray(ly, g.ppu.LinePixels())
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.pbuf.pixels)
}
func (g *Game) Layout(_, _ int) (int, int) { return ScreenWidth, ScreenHeight }
