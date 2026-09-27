package gamebro

type MMU struct {
	tiles [16 * 8 * 3 * 8 * 2]byte
	crt   cartridge
}

func (m *MMU) WriteTile(addr int, data byte) {
	if addr >= len(m.tiles) {
		// TODO: log
		panic("")
	}
	m.tiles[addr] = data
}

func (m *MMU) ReadTile(addr int) byte {
	if addr >= len(m.tiles) {
		// TODO: log
		panic("")
	}
	return m.tiles[addr]
}

func (m *MMU) ReadCode() func(mmu *MMU) {
	return m.crt.Code
}
