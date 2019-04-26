package game

import "github.com/moshenahmias/gopherboy/memory"

// NullMBC the ROM is directly mapped to memory at 0000-7FFFh
type NullMBC struct {
	rom *memory.ROM
}

// NewNullMBC creates NullMBC instance
func NewNullMBC(rom []byte) *NullMBC {
	m := NullMBC{rom: memory.NewROM(rom, 0x0000)}
	return &m
}

// Read from address 'addr'
func (n *NullMBC) Read(addr uint16) (byte, error) {

	if 0x0000 <= addr && addr <= 0x7FFF {
		return n.rom.Read(addr)
	}

	return 0, memory.ReadOutOfRangeError(addr)
}

// Write 'data' to address 'addr'
func (n *NullMBC) Write(addr uint16, data byte) error {

	if 0x0000 <= addr && addr <= 0x7FFF {
		return n.rom.Write(addr, data)
	}

	return memory.WriteOutOfRangeError(addr)
}
