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

// TODO: これより上、ppuレシーバではなくpkg関数にするべき

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

type IORegsField uint16

const (
	IORegsOBP0 IORegsField = 0xFF48
	IORegsOBP1 IORegsField = 0xFF49
)

func ReadIORegs(mmu *MMU, field IORegsField) byte {
	return mmu.Read(uint16(field))
}

type object struct {
	PosX, PosY byte
	TileId     byte
	Attr       byte
}

func ReadObject(mmu *MMU, id int) object {
	addr := uint16(0xFE00 + id<<2)
	return object{
		PosY:   mmu.Read(addr),
		PosX:   mmu.Read(addr + 1),
		TileId: mmu.Read(addr + 2),
		Attr:   mmu.Read(addr + 3),
	}
}

func (ppu *PPU) renderSpriteLine(ly int) error {
	nSprites := 0
	// TODO: LCDC.2 OBJ Size
	is8x16 := true
	height := 8
	if is8x16 {
		height = 16
	}

	for i := range 40 {
		object := ReadObject(ppu.mmu, i)
		palette := byte(0b11100100) //ReadIORegs(ppu.mmu, IORegsOBP0)
		if object.Attr&0x10 > 0 {
			palette = 0b11010000 //ReadIORegs(ppu.mmu, IORegsOBP1)
		}

		spriteY := int(object.PosY) - 16
		if ly < spriteY || ly >= spriteY+height {
			continue
		}
		nSprites++
		if nSprites > 10 {
			break
		}

		offsetY := ly - spriteY

		// TODO: LCDC.2 OBJ Size
		tileId := object.TileId
		if is8x16 {
			if offsetY < 8 {
				tileId = tileId & 0xFE
			} else {
				tileId = tileId | 0x01
				offsetY -= 8
			}
		}
		lo, hi := ppu.readTileRow(tileId, offsetY&0x7, true)
		cids := toColorIDs(lo, hi)
		baseX := int(object.PosX) - 8
		for i, cid := range cids {
			screenX := baseX + i
			if screenX < 0 || screenX > ScreenWidth {
				continue
			}
			if cid == 0 {
				continue
			}
			ppu.linePixels[screenX] = colorMap(cid, palette)
		}
	}
	return nil
}

func (ppu *PPU) RenderScanline(ly int) {
	ppu.renderBGLine(ly)
	ppu.renderSpriteLine(ly)
}
