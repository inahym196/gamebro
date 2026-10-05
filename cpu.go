package gamebro

import (
	"log/slog"
)

type cartridge interface {
	Code(cpu *CPU)
	Read(addr uint16) byte
}

type MMU struct {
	crt     cartridge
	tiles   [16 * 8 * 3 * 8 * 2]byte
	tileMap [2 * 32 * 32]byte
	wram    [0x2000]byte // In CGB, switchable
	oam     [0xA0]byte
	ioRegs  [0x80]byte
}

func NewMMU(crt cartridge) *MMU {
	return &MMU{crt: crt}
}

func (mmu *MMU) Write(addr uint16, data byte) {
	switch {
	case addr < 0x8000:
		slog.Warn("rom cannot write", "addr", addr, "data", data)
	case addr >= 0x8000 && addr < 0x9800:
		mmu.tiles[addr-0x8000] = data
	case addr >= 0x9800 && addr < 0xA000:
		mmu.tileMap[addr-0x9800] = data
	case addr >= 0xC000 && addr < 0xE000:
		mmu.wram[addr-0xC000] = data
	case addr >= 0xFE00 && addr < 0xFEA0:
		mmu.oam[addr-0xFE00] = data
	case addr >= 0xFF00 && addr < 0xFF80:
		mmu.ioRegs[addr-0xFF00] = data
	default:
		slog.Error("not impl yet")
	}
}

func (mmu *MMU) Read(addr uint16) byte {
	switch {
	case addr < 0x8000:
		return mmu.crt.Read(addr)
	case addr >= 0x8000 && addr < 0x9800:
		return mmu.tiles[addr-0x8000]
	case addr >= 0x9800 && addr < 0xA000:
		return mmu.tileMap[addr-0x9800]
	case addr >= 0xC000 && addr < 0xE000:
		return mmu.wram[addr-0xC000]
	case addr >= 0xFE00 && addr < 0xFEA0:
		return mmu.oam[addr-0xFE00]
	case addr >= 0xFF00 && addr < 0xFF80:
		return mmu.ioRegs[addr-0xFF00]
	default:
		slog.Error("not impl yet")
		return 0xFF
	}
}

type CPU struct {
	mmu *MMU
}

func NewCPU(mmu *MMU) *CPU {
	return &CPU{mmu}
}

func (cpu *CPU) Write(addr uint16, data byte) { cpu.mmu.Write(addr, data) }
func (cpu *CPU) Read(addr uint16) byte        { return cpu.mmu.Read(addr) }

func (cpu *CPU) Step() { cpu.mmu.crt.Code(cpu) }
