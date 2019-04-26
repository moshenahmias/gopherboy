package memory

// MemReg is a memory register
type MemReg byte

// Read from the register
func (m *MemReg) Read(addr uint16) (byte, error) {
	return byte(*m), nil
}

// Write to the register
func (m *MemReg) Write(addr uint16, data byte) error {
	*m = MemReg(data)
	return nil
}
