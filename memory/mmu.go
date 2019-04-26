package memory

import (
	"fmt"
	"os"
)

// MMU is the gateway for all other memory units
type MMU struct {
	mapping []Unit
}

// NewMMU creates MMU instance
func NewMMU() *MMU {
	return &MMU{mapping: make([]Unit, 65536)}
}

// Dump the memory to file
func (m *MMU) Dump(filename string) error {

	f, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer f.Close()

	for addr := uint(0); addr <= 0xFFFF; addr++ {

		if m.mapping[addr] == nil {
			f.WriteString(fmt.Sprintf("$%04x: not mapped\n", uint16(addr)))

		} else {
			data, err := m.mapping[addr].Read(uint16(addr))

			if err != nil {
				f.WriteString(fmt.Sprintf("$%04x: %s\n", uint16(addr), err))
			} else {
				f.WriteString(fmt.Sprintf("$%04x: %02x\n", uint16(addr), data))
			}
		}

	}

	return nil
}

// Map 'unit' to 'from -> to'
func (m *MMU) Map(unit Unit, from, to uint16) error {

	if from > to || uint(from) >= uint(len(m.mapping)) || uint(to) >= uint(len(m.mapping)) {
		return fmt.Errorf("invalid mapping from: %04x to %04x)", from, to)
	}

	i := uint(from)
	j := uint(to)

	for ; i <= j; i++ {
		m.mapping[i] = unit
	}

	return nil
}

// Read from address 'addr'
func (m *MMU) Read(addr uint16) (byte, error) {

	if uint(len(m.mapping)) <= uint(addr) {
		return 0, ReadOutOfRangeError(addr)
	}

	if m.mapping[addr] == nil {
		return 0, ReadAccessViolationError(addr)
	}

	return m.mapping[addr].Read(addr)
}

// Write 'data' to address 'addr'
func (m *MMU) Write(addr uint16, data byte) error {

	if uint(len(m.mapping)) <= uint(addr) {
		return WriteOutOfRangeError(addr)
	}

	if m.mapping[addr] == nil {
		return WriteAccessViolationError(addr)
	}

	return m.mapping[addr].Write(addr, data)
}
