package gamebro

import "log/slog"

type mmu struct {
	tiles   [16 * 8 * 3 * 8 * 2]byte
	tileMap [2 * 32 * 32]byte
	crt     cartridge
}

type CPU struct {
	*mmu
}

func NewCPU(crt cartridge) *CPU {
	return &CPU{&mmu{crt: crt}}
}

func (cpu *CPU) Write(addr uint16, data byte) {
	switch {
	case addr >= 0x8000 && addr < 0xA000:
		cpu.tiles[addr-0x8000] = data
	default:
		slog.Error("not impl yet")
	}
}

func (cpu *CPU) Read(addr uint16) byte {
	switch {
	case addr >= 0x8000 && addr < 0xA000:
		return cpu.tiles[addr-0x8000]
	default:
		slog.Error("not impl yet")
		return 0xFF
	}
}

func (cpu *CPU) Step() { cpu.crt.Code(cpu) }
