/*H**********************************************************************
* FILENAME :        cartridge.go
*
* PACKAGE :			game
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package game

import (
	"errors"
	"io/ioutil"

	"github.com/moshenahmias/gopherboy/memory"

	"fmt"
)

// ErrCorrupted is returned when the loaded ROM is corrupted
var ErrCorrupted = errors.New("ErrCorrupted")

// Cartridge represents the GB Classic game cartridge
type Cartridge struct {
	mbc memory.Unit
}

// NewCartridge creates Cartridge instance
func NewCartridge(fileROM string) (*Cartridge, error) {

	// load rom from file
	romData, err := ioutil.ReadFile(fileROM)

	if err != nil {
		return nil, err
	}

	rom := memory.NewROM(romData, 0)

	var ramData []byte

	// read ram type
	if ramType, err := rom.Read(0x0149); err == nil {

		// alocate ram
		switch ramType {

		case 0x00: // No RAM or MBC2
			break
		case 0x01: // 2 KBytes
			ramData = make([]byte, 2048)
		case 0x02: // 8 Kbytes
			ramData = make([]byte, 8192)
		case 0x03: //32 KBytes
			ramData = make([]byte, 32768)
		default:
			return nil, ErrCorrupted
		}

	} else {
		// failed to read ram type
		return nil, err
	}

	// read MBC type
	mbcType, err := rom.Read(0x0147)

	if err != nil {
		return nil, err
	}

	c := Cartridge{}

	// create MBC
	switch mbcType {

	case 0x00: // rom only

		c.mbc = NewNullMBC(romData)

	case 0x01, 0x02, 0x03: // MBC1

		c.mbc = NewMBC1(romData, ramData)

	case 0x05, 0x06: // MBC2

		c.mbc = NewMBC2(romData)

	default:

		return nil, fmt.Errorf("cartridge type not supported (%x)", mbcType)
	}

	return &c, nil
}

// Read from address 'addr'
func (c *Cartridge) Read(addr uint16) (byte, error) {
	return c.mbc.Read(addr)
}

// Write 'data' to address 'addr'
func (c *Cartridge) Write(addr uint16, data byte) error {
	return c.mbc.Write(addr, data)
}
