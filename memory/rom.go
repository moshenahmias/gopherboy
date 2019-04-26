package memory

import "fmt"

// ROM is a read only memory unit
type ROM struct {
	data       []byte
	window     []byte
	addrOffset uint16
}

// NewROM creates ROM instance
func NewROM(data []byte, addrOffset uint16) *ROM {

	r := ROM{data: data, addrOffset: addrOffset}
	r.SetWindow(0)

	return &r
}

// SetWindow sets the virtual data view
// read access is done against that view
func (r *ROM) SetWindow(offset uint32) error {

	if uint32(len(r.data)) <= offset {
		return fmt.Errorf("window offset out of range (%d)", offset)
	}

	r.window = r.data[offset:]

	return nil
}

// Read from address 'addr'
func (r *ROM) Read(addr uint16) (byte, error) {

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
func (r *ROM) Write(addr uint16, data byte) error {
	return nil
}
