package gamebro

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/inahym196/minifont"
	"github.com/tinne26/ptxt"
)

type DebugGame struct {
	isDebug     bool
	game        *Game
	tilesScreen *ebiten.Image
	gameScreen  *ebiten.Image
	mapScreen   [2]*ebiten.Image
	text        *ptxt.Renderer
}

func NewDebugGame(game *Game) (*DebugGame, error) {
	strand, err := ptxt.NewStrand(minifont.Font())
	if err != nil {
		return nil, err
	}
	renderer := ptxt.NewRenderer()
	renderer.SetStrand(strand)
	renderer.SetColor(color.RGBA{255, 255, 255, 255})

	const blockW, blockH, blockNum = 16, 8, 3

	return &DebugGame{
		isDebug:     true,
		game:        game,
		gameScreen:  ebiten.NewImage(ScreenWidth, ScreenHeight),
		tilesScreen: ebiten.NewImage(TileSize*blockW, TileSize*blockH*blockNum),
		mapScreen: [2]*ebiten.Image{
			ebiten.NewImage(256, 192),
			ebiten.NewImage(256, 192),
		},
		text: renderer,
	}, nil
}

func drawTileRow(buf *PixelBuffer, x, y int, lo, hi byte) {
	for offsetX, cid := range toColorIDs(lo, hi) {
		c := colorMap(cid, 0b1110_0100)
		buf.SetPixelGray(x+offsetX, y, c)
	}
}

func drawTile(buf *PixelBuffer, startX, startY int, row [16]byte) {
	for i := range len(row) / 2 {
		drawTileRow(buf, startX, startY+i, row[i*2], row[i*2+1])
	}
}

func updateTilesScreen(tilesImage *ebiten.Image, tiles [6144]byte) {
	tilesImage.Clear()
	bounds := tilesImage.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	buf := NewPixelBuffer(w, h)
	const TilesPerRow = 16
	for tileY := range len(tiles) / (TilesPerRow * 16) {
		for tileX := range TilesPerRow {
			tileIdx := (tileY*TilesPerRow + tileX) * 16
			tile := [16]byte(tiles[tileIdx : tileIdx+16])
			drawTile(buf, tileX*TileSize, tileY*TileSize, tile)
		}
	}
	tilesImage.WritePixels(buf.Pixels())
}

func checkeredColor(x, y int) color.RGBA {
	if (x+y)%2 == 0 {
		return color.RGBA{0xAA, 0xAA, 0xAA, 0xFF}
	}
	return color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
}

func updateMapScreen(mapImage *ebiten.Image, tileMap []byte, text *ptxt.Renderer) {
	if len(tileMap) != 1024 {
		panic("tileMap len is not 1024")
	}
	mapImage.Clear()
	const (
		glyphWidth  = 4
		lineSpacing = 6
	)

	for y := range 32 {
		for x := range 32 {
			text.SetColor(checkeredColor(x, y))
			s := fmt.Sprintf("%02X", tileMap[y*32+x])
			text.Draw(mapImage, s, x*glyphWidth*2, (y+1)*lineSpacing)
		}
	}
}

func (g *DebugGame) Update() error {
	if err := g.game.Update(); err != nil {
		return err
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.isDebug = !g.isDebug
	}
	if g.isDebug {
		updateTilesScreen(g.tilesScreen, g.game.cpu.mmu.tiles)
		updateMapScreen(g.mapScreen[0], g.game.cpu.mmu.tileMap[:1024], g.text)
		updateMapScreen(g.mapScreen[1], g.game.cpu.mmu.tileMap[1024:], g.text)
	}
	return nil
}

func drawViewPort(screen *ebiten.Image, scx, scy uint8) {
	vector.StrokeRect(screen, float32(scx)-1, float32(scy)*0.75-1, 8*20, 6*18+1, 1, color.RGBA{0xAA, 0, 0, 0xFF}, false)
	if scx > 0x60 {
		vector.StrokeRect(screen, -1, float32(scy)*0.75, float32(scx-0x60), 6*18, 1, color.RGBA{0xAA, 0, 0, 0xFF}, false)
	}
	if scy > 0x70 {
		vector.StrokeRect(screen, float32(scx)-1, -1, 8*20, float32(scy-0x70)*0.75+1, 1, color.RGBA{0xAA, 0, 0, 0xFF}, false)
	}
	if scx > 0x60 && scy > 0x70 {
		vector.StrokeRect(screen, -1, -2, float32(scx-0x60), float32(scy-0x70)*0.75, 1, color.RGBA{0xAA, 0, 0, 0xFF}, false)
	}
}

func (g *DebugGame) Draw(screen *ebiten.Image) {
	g.game.Draw(g.gameScreen)
	screen.DrawImage(g.gameScreen, nil)

	if g.isDebug {
		bounds := g.gameScreen.Bounds()
		gameW, gameH := bounds.Dx(), bounds.Dy()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, float64(gameH)+TileSize)
		screen.DrawImage(g.tilesScreen, op)

		msg := fmt.Sprintf("TPS: %.2f", ebiten.ActualTPS())
		_, gh := g.Layout(0, 0)
		ebitenutil.DebugPrintAt(screen, msg, 0, gh-16)

		drawViewPort(g.mapScreen[0], 0, 0)

		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(gameW)+TileSize, 0)
		screen.DrawImage(g.mapScreen[0], op)

		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(gameW)+TileSize, 192+6)
		screen.DrawImage(g.mapScreen[1], op)
	}
}

func (g *DebugGame) Layout(w, h int) (int, int) {
	gw, gh := g.game.Layout(w, h)
	if g.isDebug {
		return gw * 3, gh * 3
	}
	return gw, gh
}
