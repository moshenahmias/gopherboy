/*H**********************************************************************
* FILENAME :        romemreg.go
*
* PACKAGE :			memory
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package memory

// ROMemReg is a read only memory register
type ROMemReg byte

// Read from the register
func (r *ROMemReg) Read(addr uint16) (byte, error) {
	return byte(*r), nil
}

// Write to the register
func (r *ROMemReg) Write(addr uint16, data byte) error {
	return nil
}
