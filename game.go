package gamebro

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 160
	ScreenHeight = 144
)

func NewGame(crt cartridge) *Game {
	joypad := NewJoyPad()
	mmu := NewMMU(crt, joypad)
	return &Game{
		pbuf:   NewPixelBuffer(ScreenWidth, ScreenHeight),
		cpu:    NewCPU(mmu),
		ppu:    NewPPU(mmu),
		joypad: joypad,
	}
}

type Game struct {
	pbuf   *PixelBuffer
	cpu    *CPU
	ppu    *PPU
	joypad *JoyPad
}

func dirInput() (input byte) {
	if ebiten.IsKeyPressed(ebiten.KeyJ) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		input |= 0x08
	}
	if ebiten.IsKeyPressed(ebiten.KeyK) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		input |= 0x04
	}
	if ebiten.IsKeyPressed(ebiten.KeyH) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		input |= 0x02
	}
	if ebiten.IsKeyPressed(ebiten.KeyL) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		input |= 0x01
	}
	return input
}

func buttonInput() (input byte) {
	if ebiten.IsKeyPressed(ebiten.KeyM) || ebiten.IsKeyPressed(ebiten.KeyEnter) {
		input |= 0x08
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeySpace) {
		input |= 0x04
	}
	return input
}

func (g *Game) Update() error {
	// joypadはlow-activeのため反転
	g.joypad.SetStates(^buttonInput(), ^dirInput())
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
