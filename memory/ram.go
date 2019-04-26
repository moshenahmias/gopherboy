package memory

import (
	"fmt"
)

// RAM is a read/write memory unit
type RAM struct {
	data       []byte
	window     []byte
	addrOffset uint16
}

// NewRAM creates RAM instance
func NewRAM(data []byte, addrOffset uint16) *RAM {

	r := RAM{data: data, addrOffset: addrOffset}
	r.SetWindow(0)

	return &r
}

// SetWindow sets the virtual data view
// read/write access is done against that view
func (r *RAM) SetWindow(offset uint32) error {

	if uint32(len(r.data)) <= offset {
		return fmt.Errorf("window offset out of range (%d)", offset)
	}

	r.window = r.data[offset:]

	return nil
}

// Read from address 'addr'
func (r *RAM) Read(addr uint16) (byte, error) {

	if r.addrOffset > addr {
		return 0, ReadOutOfRangeError(addr)
	}

	physicalAddr := addr - r.addrOffset

	if uint(len(r.window)) <= uint(physicalAddr) {
		return 0, ReadOutOfRangeError(addr)
	}

	return r.window[physicalAddr], nil
}

// Write 'data' to address 'addr'
func (r *RAM) Write(addr uint16, data byte) error {

	if r.addrOffset > addr {
		return WriteOutOfRangeError(addr)
	}

	physicalAddr := addr - r.addrOffset

	if uint(len(r.window)) <= uint(physicalAddr) {
		return WriteOutOfRangeError(addr)
	}

	r.window[physicalAddr] = data

	return nil
}
