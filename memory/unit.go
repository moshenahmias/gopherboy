/*H**********************************************************************
* FILENAME :        uint.go
*
* PACKAGE :			memory
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package memory

import "fmt"

// ReadAccessViolationError creates a memory read access violation error
func ReadAccessViolationError(addr uint16) error {
	return fmt.Errorf("Read access violation at %04x", addr)
}

// WriteAccessViolationError creates a memory read access violation error
func WriteAccessViolationError(addr uint16) error {
	return fmt.Errorf("Write access violation at %04x", addr)
}

// ReadOutOfRangeError creates a memory out of range read error
func ReadOutOfRangeError(addr uint16) error {
	return fmt.Errorf("out of range read at %04x", addr)
}

// WriteOutOfRangeError creates a memory out of range write error
func WriteOutOfRangeError(addr uint16) error {
	return fmt.Errorf("out of range write at %04x", addr)
}

// Unit is an interface for all memory units
type Unit interface {
	Read(addr uint16) (byte, error)
	Write(addr uint16, data byte) error
}
