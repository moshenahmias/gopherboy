/*H**********************************************************************
* FILENAME :        memreg.go
*
* PACKAGE :			memory
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

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
