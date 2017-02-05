/*H**********************************************************************
* FILENAME :        biosunmapper.go
*
* PACKAGE :			memory
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package memory

// BiosUnmapper memory unit is used to unmap the 256B bios from
// 0x0000 and map the game cartridge when the bios execution ends
type BiosUnmapper struct {
	mmu       *MMU
	cartridge Unit
}

// NewBiosUnmapper creates BiosUnmapper instance
func NewBiosUnmapper(mmu *MMU, cartridge Unit) *BiosUnmapper {
	bu := BiosUnmapper{mmu: mmu, cartridge: cartridge}
	return &bu
}

// Read access is not allowed for this memory unit
func (b *BiosUnmapper) Read(addr uint16) (byte, error) {
	return 0, ReadAccessViolationError(addr)
}

// Write will unmap the 256B bios from 0x0000 and map the game cartridge instead
func (b *BiosUnmapper) Write(addr uint16, data byte) error {

	if addr != 0xFF50 {
		return WriteAccessViolationError(addr)
	}

	if data == 0x01 {
		b.mmu.Map(b.cartridge, 0x0000, 0x00FF)
	}

	return nil
}
