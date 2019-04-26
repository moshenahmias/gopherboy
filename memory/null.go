package memory

// Null is a memory unit which writes and reads nothing
type Null struct {
}

// Read nothing
func (n *Null) Read(addr uint16) (byte, error) {
	return 0, nil
}

// Write nothing
func (n *Null) Write(addr uint16, data byte) error {
	return nil
}
