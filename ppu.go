package gamebro

type PPU struct {
	mmu        *MMU
	linePixels [16]uint8
}

func NewPPU(mmu *MMU) *PPU {
	return &PPU{mmu: mmu}
}

func (ppu *PPU) LinePixels() []uint8 { return ppu.linePixels[:] }

func (ppu *PPU) readTileRow(id byte, offsetY int) (lo, hi byte) {
	base := 0x8000 + uint16(id<<4)
	return ppu.mmu.Read(base + uint16(offsetY<<1)), ppu.mmu.Read(base + uint16(offsetY<<1) + 1)
}

func (ppu *PPU) fetchBGWindowTileRow(tileX, tileY, offsetY, mapID int) (lo, hi byte) {
	addr := [2]uint16{0x9800, 0x9C00}[mapID] + uint16(tileY<<5|tileX)
	tileID := ppu.mmu.Read(addr)
	return ppu.readTileRow(tileID, offsetY)
}

func (ppu *PPU) fetchBGTileRow(tileX, tileY, offsetY int) (lo, hi byte) {
	return ppu.fetchBGWindowTileRow(tileX, tileY, offsetY, 0)
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

func (ppu *PPU) renderScanline(ly int) {
	tileY := ly / 8
	offsetY := ly % 8
	for tileX := range 16 / 8 {
		lo, hi := ppu.fetchBGTileRow(tileX, tileY, offsetY)
		baseX := tileX * 8
		for i, cid := range toColorIDs(lo, hi) {
			bgp := uint8(0b11100100)
			ppu.linePixels[baseX+i] = colorMap(cid, bgp)
		}
	}
}
