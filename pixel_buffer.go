package gamebro

import "fmt"

type PixelBuffer struct {
	pixels []byte
	width  int
	height int
}

func NewPixelBuffer(w, h int) *PixelBuffer {
	return &PixelBuffer{
		pixels: make([]byte, w*h*4),
		width:  w,
		height: h,
	}
}

func (buf *PixelBuffer) SetScanlineGray(y int, lineBytes []uint8) error {
	if y < 0 || y >= buf.height {
		return fmt.Errorf("out of range: %d", y)
	}
	if len(lineBytes) > buf.width {
		return fmt.Errorf("out of range: %d", len(lineBytes))
	}
	rowOffset := y * buf.width * 4
	for x, c := range lineBytes {
		idx := rowOffset + x*4
		buf.pixels[idx] = c
		buf.pixels[idx+1] = c
		buf.pixels[idx+2] = c
		buf.pixels[idx+3] = 0xFF
	}
	return nil
}

func (buf *PixelBuffer) Pixels() []byte { return buf.pixels }
