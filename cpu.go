package gamebro

import "log/slog"

type mmu struct {
	crt     cartridge
	tiles   [16 * 8 * 3 * 8 * 2]byte
	tileMap [2 * 32 * 32]byte
	wram    [0x2000]byte // In CGB, switchable
	oam     [0xA0]byte
	ioRegs  [0x80]byte
}

type CPU struct {
	*mmu
}

func NewCPU(crt cartridge) *CPU {
	return &CPU{&mmu{crt: crt}}
}

func (cpu *CPU) Write(addr uint16, data byte) {
	switch {
	case addr >= 0x8000 && addr < 0x9800:
		cpu.tiles[addr-0x8000] = data
	case addr >= 0x9800 && addr < 0xA000:
		cpu.tileMap[addr-0x9800] = data
	case addr >= 0xC000 && addr < 0xE000:
		cpu.wram[addr-0xC000] = data
	case addr >= 0xFF00 && addr < 0xFF80:
		cpu.ioRegs[addr-0xFF00] = data
	default:
		slog.Error("not impl yet")
	}
}

func (cpu *CPU) Read(addr uint16) byte {
	switch {
	case addr >= 0x8000 && addr < 0x9800:
		return cpu.tiles[addr-0x8000]
	case addr >= 0x9800 && addr < 0xA000:
		return cpu.tileMap[addr-0x9800]
	case addr >= 0xC000 && addr < 0xE000:
		return cpu.wram[addr-0xC000]
	case addr >= 0xFF00 && addr < 0xFF80:
		return cpu.ioRegs[addr-0xFF00]
	default:
		slog.Error("not impl yet")
		return 0xFF
	}
}

func (cpu *CPU) Step() { cpu.crt.Code(cpu) }
