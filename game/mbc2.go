/*H**********************************************************************
* FILENAME :        mbc2.go
*
* PACKAGE :			game
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    06 Jan 2017
*
*H*/

package game

import "github.com/moshenahmias/gopherboy/memory"

// MBC2 (max 256KByte ROM and 512x4 bits RAM)
type MBC2 struct {
	rom        *memory.ROM
	otherBanks *memory.ROM
	ram        *memory.RAM
	bankROM    uint32
	enableRAM  bool
}

// NewMBC2 creates mbc2 instance
func NewMBC2(rom []byte) *MBC2 {

	m := MBC2{
		rom:        memory.NewROM(rom, 0x0000),
		otherBanks: memory.NewROM(rom, 0x4000),
		ram:        memory.NewRAM(make([]byte, 512), 0xA000),
		bankROM:    1}

	return &m
}

// Read from address 'addr' at the target bank
func (m *MBC2) Read(addr uint16) (byte, error) {

	// rom bank 0
	if 0x0000 <= addr && addr <= 0x3FFF {

		if err := m.rom.SetWindow(0); err != nil {
			return 0, err
		}

		return m.rom.Read(addr)
	}

	// other rom banks
	if 0x4000 <= addr && addr <= 0x7FFF {

		if err := m.otherBanks.SetWindow(m.bankROM * 16384); err != nil {
			return 0, err
		}

		return m.otherBanks.Read(addr)
	}

	// ram
	if 0xA000 <= addr && addr <= 0xA1FF {
		return m.ram.Read(addr)
	}

	return 0, memory.ReadOutOfRangeError(addr)
}

// Write 'data' to address 'addr' at the target bank
// or change the MBC control registers
func (m *MBC2) Write(addr uint16, data byte) error {

	// enable ram
	if 0x0000 <= addr && addr <= 0x1FFF {

		if addr&0x0100 == 0 {
			m.enableRAM = data&0x0F == 0x0A
		}

		return nil
	}

	// rom bank
	if 0x2000 <= addr && addr <= 0x3FFF {

		if addr&0x0100 == 0x0100 {

			if data == 0 {
				data = 1
			}

			m.bankROM = uint32(data & 0x0F)
		}

		return nil
	}

	// ram
	if 0xA000 <= addr && addr <= 0xA1FF {

		if !m.enableRAM {
			return nil
		}

		return m.ram.Write(addr, data&0x0F)
	}

	return memory.WriteOutOfRangeError(addr)
}
