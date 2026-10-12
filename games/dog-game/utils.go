package dog

func (g *DogGame) readRomTile(base uint16) (tile [BytesPerTile]byte) {
	for i := range BytesPerTile {
		tile[i] = g.cpu.Read(base + uint16(i))
	}
	return tile
}

func (g *DogGame) readRomBGTile(id int) (tile [BytesPerTile]byte) {
	base := BGBaseAddr + uint16(id*BytesPerTile)
	return g.readRomTile(base)
}

func (g *DogGame) writeTile(addr uint16, tile [BytesPerTile]byte) {
	for i, data := range tile {
		g.cpu.Write(addr+uint16(i), data)
	}
}

func (g *DogGame) writeTileByID(id int, tile [BytesPerTile]byte) {
	addr := 0x8000 + uint16(id<<4)
	g.writeTile(addr, tile)
}

func (g *DogGame) loadBGTile(srcId, dstId int) {
	tile := g.readRomBGTile(srcId)
	g.writeTileByID(dstId, tile)
}

func (g *DogGame) readRomWindowTile(id int) (tile [BytesPerTile]byte) {
	base := WinBaseAddr + uint16(id*BytesPerTile)
	return g.readRomTile(base)
}

func (g *DogGame) loadWindowTile(srcId, dstId int) {
	tile := g.readRomWindowTile(srcId)
	g.writeTileByID(dstId, tile)
}

func (g *DogGame) readRomSpriteTile(id int) (tile [BytesPerTile]byte) {
	base := SpriteBaseAddr + uint16(id*BytesPerTile)
	return g.readRomTile(base)
}

func (g *DogGame) loadSpriteTile(srcId, dstId int) {
	tile := g.readRomSpriteTile(srcId)
	g.writeTileByID(dstId, tile)
}
