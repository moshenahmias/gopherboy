/*H**********************************************************************
* FILENAME :        mbc1.go
*
* PACKAGE :			game
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package game

import "github.com/moshenahmias/gopherboy/memory"

// MBC1 (max 2MByte ROM and/or 32KByte RAM)
type MBC1 struct {
	rom        *memory.ROM
	otherBanks *memory.ROM
	ram        *memory.RAM
	mode       byte
	bankRAM    uint32
	bankROM0   uint32
	bankROM1   uint32
	enableRAM  bool
}

// NewMBC1 creates mbc1 instance
func NewMBC1(rom []byte, ram []byte) *MBC1 {

	m := MBC1{
		rom:        memory.NewROM(rom, 0x0000),
		otherBanks: memory.NewROM(rom, 0x4000),
		ram:        memory.NewRAM(ram, 0xA000),
		bankROM0:   1,
		bankROM1:   1}

	return &m
}

// Read from address 'addr' at the target bank
func (m *MBC1) Read(addr uint16) (byte, error) {

	// rom bank 0
	if 0x0000 <= addr && addr <= 0x3FFF {

		if err := m.rom.SetWindow(0); err != nil {
			return 0, err
		}

		return m.rom.Read(addr)
	}

	// other rom banks
	if 0x4000 <= addr && addr <= 0x7FFF {

		var bank uint32

		if m.mode == 0 {

			// mode 0
			bank = m.bankROM0

		} else {

			// mode 1
			bank = m.bankROM1
		}

		if err := m.otherBanks.SetWindow(bank * 16384); err != nil {
			return 0, err
		}

		return m.otherBanks.Read(addr)
	}

	// ram
	if 0xA000 <= addr && addr <= 0xBFFF {

		var bank uint32

		if m.mode == 1 {
			bank = m.bankRAM
		}

		if err := m.ram.SetWindow(bank * 8192); err != nil {
			return 0, err
		}

		return m.ram.Read(addr)
	}

	return 0, memory.ReadOutOfRangeError(addr)
}

// Write 'data' to address 'addr' at the target bank
// or change the MBC control registers
func (m *MBC1) Write(addr uint16, data byte) error {

	// enable ram
	if 0x0000 <= addr && addr <= 0x1FFF {
		m.enableRAM = data&0x0F == 0x0A
		return nil
	}

	// rom bank
	if 0x2000 <= addr && addr <= 0x3FFF {

		if data == 0 {
			data = 1
		}

		m.bankROM1 = uint32(data & 0x1F)
		m.bankROM0 = (m.bankROM0 & 0xFFFFFFE0) | m.bankROM1

		return nil
	}

	// ram/rom bank
	if 0x4000 <= addr && addr <= 0x5FFF {

		m.bankRAM = uint32(data & 0x03)
		m.bankROM0 = (m.bankROM0 & 0x0000001F) | (m.bankRAM << 5)

		return nil
	}

	// mode
	if 0x6000 <= addr && addr <= 0x7FFF {
		m.mode = data & 0x01
		return nil
	}

	// ram
	if 0xA000 <= addr && addr <= 0xBFFF {

		if !m.enableRAM {
			return nil
		}

		var bank uint32

		if m.mode == 1 {
			bank = m.bankRAM
		}

		if err := m.ram.SetWindow(bank * 8192); err != nil {
			return err
		}

		return m.ram.Write(addr, data)
	}

	return memory.WriteOutOfRangeError(addr)
}
