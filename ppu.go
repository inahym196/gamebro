package gamebro

const (
	TileSize = 8
)

type PPU struct {
	mmu        *MMU
	linePixels [ScreenWidth]uint8
}

func NewPPU(mmu *MMU) *PPU {
	return &PPU{mmu: mmu}
}

func (ppu *PPU) LinePixels() []uint8 { return ppu.linePixels[:] }

func (ppu *PPU) readTileRow(id byte, offsetY int, method80 bool) (lo, hi byte) {
	tileAddr := 0x8000 + uint16(id)<<4
	if !method80 {
		tileAddr = uint16(0x9000 + int(int8(id))<<4)
	}
	rowAddr := tileAddr + uint16(offsetY<<1)
	return ppu.mmu.Read(rowAddr), ppu.mmu.Read(rowAddr + 1)
}

func (ppu *PPU) fetchBGWindowTileRow(tileX, tileY, offsetY int, base uint16) (lo, hi byte) {
	addr := base | uint16(tileY<<5|tileX)
	tileID := ppu.mmu.Read(addr)
	return ppu.readTileRow(tileID, offsetY, false)
}

func (ppu *PPU) fetchWindowTileRow(tileX, tileY, offsetY int) (lo, hi byte) {
	// TODO: LCDC.6
	base := [2]uint16{0x9800, 0x9C00}[1]
	return ppu.fetchBGWindowTileRow(tileX, tileY, offsetY, base)
}

func (ppu *PPU) fetchBGTileRow(tileX, tileY, offsetY int) (lo, hi byte) {
	// TODO: LCDC.3
	base := [2]uint16{0x9800, 0x9C00}[0]
	return ppu.fetchBGWindowTileRow(tileX, tileY, offsetY, base)
}

func colorMap(colorID int, palette uint8) byte {
	cmap := [4]uint8{0xFF, 0xAA, 0x55, 0x00}
	paletteID := int((palette >> (colorID << 1)) & 0b11)
	return cmap[paletteID]
}

func toColorIDs(lo, hi byte) [TileSize]int {
	var cNums [TileSize]int
	for x := range TileSize {
		shift := 7 - x
		_lo := (lo >> shift) & 1
		_hi := (hi >> shift) & 1
		cNums[x] = int(_hi<<1 | _lo)
	}
	return cNums
}

func (ppu *PPU) renderBGLine(ly int) {
	tileY := ly / TileSize
	offsetY := ly % TileSize

	wy := 128
	windowLine := false
	if ly >= wy {
		windowLine = true
		winY := ly - wy
		tileY = winY / TileSize
		offsetY = winY % TileSize
	}
	for tileX := range ScreenWidth / TileSize {
		var lo, hi byte
		if windowLine {
			lo, hi = ppu.fetchWindowTileRow(tileX, tileY, offsetY)
		} else {
			lo, hi = ppu.fetchBGTileRow(tileX, tileY, offsetY)
		}
		baseX := tileX * TileSize
		for i, cid := range toColorIDs(lo, hi) {
			bgp := uint8(0b11100100)
			ppu.linePixels[baseX+i] = colorMap(cid, bgp)
		}
	}
}

func (ppu *PPU) RenderScanline(ly int) {
	ppu.renderBGLine(ly)
}
