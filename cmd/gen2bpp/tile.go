package main

import (
	"fmt"
	"image"
	"image/color"
	"iter"
	"os"
)

type Tile [8][8]uint8

// タイルの一行分の色コードを受け取り、2bpp(lo, hi byte)に変換して返す
func Encode2bpp(colors [8]uint8) (lo, hi byte) {
	for i, c := range colors {
		shift := 7 - i
		lo |= c & 1 << shift
		hi |= (c >> 1) & 1 << shift
	}
	return lo, hi
}

func EncodeTile(tile Tile) (bpps [8][2]byte) {
	for iy := range 8 {
		lo, hi := Encode2bpp(tile[iy])
		bpps[iy][0] = lo
		bpps[iy][1] = hi
	}
	return bpps
}

type TilesImage struct {
	img           image.Image
	width, height int
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func NewTilesImage(path string) (*TilesImage, error) {

	img, err := loadImage(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load image: %w", err)
	}
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	tileSize := 8

	if width%tileSize != 0 || height%tileSize != 0 {
		return nil, fmt.Errorf("no evenly divisible by the tile size %d, (%d,%d)", tileSize, width, height)
	}
	return &TilesImage{img, width, height}, nil
}

func colorToGBValue(c color.Color) uint8 {
	r, g, b, _ := c.RGBA()
	r8 := uint8(r >> 8)
	g8 := uint8(g >> 8)
	b8 := uint8(b >> 8)

	if r8 == 255 && g8 == 255 && b8 == 255 {
		return 0 // #FFF(00)
	}
	if r8 == 170 && g8 == 170 && b8 == 170 {
		return 1 // #AAA(01)
	}
	if r8 == 85 && g8 == 85 && b8 == 85 {
		return 2 // #555(10)
	}
	if r8 == 0 && g8 == 0 && b8 == 0 {
		return 3 // #000(11)
	}
	panic(fmt.Sprintf("undefined color: R: %d G: %d B: %d", r8, g8, b8))
}

func (ti *TilesImage) tileRowAt(startX, y int) (row [8]uint8) {
	for ix := range 8 {
		c := ti.img.At(startX+ix, y)
		row[ix] = colorToGBValue(c)
	}
	return row
}

func (ti *TilesImage) tileAt(tileX, tileY int) (tile Tile) {
	for iy := range 8 {
		tile[iy] = ti.tileRowAt(tileX*8, tileY+iy)
	}
	return tile
}

func (ti *TilesImage) Tiles() iter.Seq2[int, Tile] {
	return func(yield func(int, Tile) bool) {
		for y := range ti.height / 8 {
			for x := range ti.width / 8 {
				if !yield(x+y*8, ti.tileAt(x, y)) {
					return
				}
			}
		}
	}
}
