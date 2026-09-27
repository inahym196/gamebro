package gamebro

import "log/slog"

type MMU struct {
	tiles   [16 * 8 * 3 * 8 * 2]byte
	tileMap [2 * 32 * 32]byte
	crt     cartridge
}

func (m *MMU) Write(addr uint16, data byte) {
	switch {
	case addr >= 0x8000 && addr < 0xA000:
		m.tiles[addr-0x8000] = data
	default:
		slog.Error("not impl yet")
	}
}

func (m *MMU) Read(addr uint16) byte {
	switch {
	case addr >= 0x8000 && addr < 0xA000:
		return m.tiles[addr-0x8000]
	default:
		slog.Error("not impl yet")
		return 0xFF
	}
}

func (m *MMU) ReadCode() func(mmu *MMU) {
	return m.crt.Code
}
