/*H**********************************************************************
* FILENAME :        instructiontemplates.go
*
* PACKAGE :			cpu
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package cpu

// SetNr sets bit N in the register
func (c *Core) SetNr(n uint, dst *Register8) {
	bit := (byte(0x01) << n)
	dst.set(dst.get() | bit)
}

// ResNr resets bit N in the register
func (c *Core) ResNr(n uint, dst *Register8) {
	bit := ^(byte(0x01) << n)
	dst.set(dst.get() & bit)
}

// BitNr tests bit N in the register
func (c *Core) BitNr(n uint, dst *Register8) {
	c.setHalfCarryFlag(true)
	c.setSubtractFlag(false)
	bit := byte(0x01) << n
	c.setZeroFlag(dst.get()&bit == 0)
}

// SrlR shifts the register right one bit position.
// the contents of bit 0 are copied to the carry flag
// and a zero is put into bit 7.
func (c *Core) SrlR(dst *Register8) {
	c.setHalfCarryFlag(false)
	c.setSubtractFlag(false)
	dst.rotateRight()
	r := dst.get()
	c.setCarryFlag(r&0x80 == 0x80)
	dst.set(r & 0x7F)
	c.setZeroFlag(dst.get() == 0)
}

// SwapR high 4 bits with the low 4 bits
func (c *Core) SwapR(dst *Register8) {
	c.setHalfCarryFlag(false)
	c.setSubtractFlag(false)
	c.setCarryFlag(false)
	r := dst.get()
	dst.set((r << 4) | (r >> 4))
	c.setZeroFlag(dst.get() == 0)
}

// SlaR shifts the register left one bit position.
// the contents of bit 7 are copied to the carry flag
// and a zero is put into bit 0
func (c *Core) SlaR(dst *Register8) {

	c.setHalfCarryFlag(false)
	c.setSubtractFlag(false)
	dst.rotateLeft()
	r := dst.get()
	c.setCarryFlag(r&0x01 == 0x01)
	dst.set(r & 0xFE)
	c.setZeroFlag(dst.get() == 0)
}

// SraR shifts the register right one bit position.
// the contents of bit 0 are copied to the carry flag
// and the previous contents of bit 7 are unchanged.
func (c *Core) SraR(dst *Register8) {

	c.setHalfCarryFlag(false)
	c.setSubtractFlag(false)
	dst.rotateRight()
	r := dst.get()
	c.setCarryFlag(r&0x80 == 0x80)
	dst.set((r & 0x7F) | ((r << 1) & 0x80))
	c.setZeroFlag(dst.get() == 0)
}

// RrR rotates the register right one bit position.
// the contents of bit 0 are copied to the carry flag and
// the previous contents of the carry flag are copied to bit 7
func (c *Core) RrR(dst *Register8, cb bool) {

	c.setHalfCarryFlag(false)
	c.setSubtractFlag(false)
	dst.rotateRight()
	r := dst.get()
	carry := c.carryFlag()
	c.setCarryFlag(r&0x80 == 0x80)
	if carry {
		dst.set(r | 0x80)
	} else {
		dst.set(r & 0x7F)
	}

	c.setZeroFlag(cb && dst.get() == 0)
}

// RlR rotates the register left one bit position.
// the contents of bit 7 are copied to the carry flag and
// the previous contents of the carry flag are copied to bit 0.
func (c *Core) RlR(dst *Register8, cb bool) {

	c.setHalfCarryFlag(false)
	c.setSubtractFlag(false)
	dst.rotateLeft()
	r := dst.get()
	carry := c.carryFlag()
	c.setCarryFlag(r&0x01 == 0x01)
	if carry {
		dst.set(r | 0x01)
	} else {
		dst.set(r & 0xFE)
	}

	c.setZeroFlag(cb && dst.get() == 0)
}

// RrcR rotates the register right one bit position.
// the contents of bit 0 are copied to the carry flag and bit 7.
func (c *Core) RrcR(dst *Register8, cb bool) {

	c.setHalfCarryFlag(false)
	c.setSubtractFlag(false)
	dst.rotateRight()
	c.setCarryFlag(dst.get()&0x80 == 0x80)
	c.setZeroFlag(cb && dst.get() == 0)
}

// RlcR rotates the register left one bit position.
// the contents of bit 7 are copied to the carry flag and bit 0.
func (c *Core) RlcR(dst *Register8, cb bool) {

	c.setHalfCarryFlag(false)
	c.setSubtractFlag(false)
	c.setCarryFlag(dst.get()&0x80 == 0x80)
	dst.rotateLeft()
	c.setZeroFlag(cb && dst.get() == 0)
}

// Rst jumps to the given address
func (c *Core) Rst(addr uint16) error {

	c.pc.increment()
	c.sp.decrement()

	if err := c.mmu.Write(c.sp.get(), c.pc.highByte()); err != nil {
		return err
	}

	c.sp.decrement()

	if err := c.mmu.Write(c.sp.get(), c.pc.lowByte()); err != nil {
		return err
	}

	c.pc.set(addr - 1)

	return nil
}

// CallCondA16 calls the function at the immediate
// 16bit address if the contidion is met
func (c *Core) CallCondA16(cond bool) error {

	im16, err := c.loadImmediate16()

	if err != nil {
		return err
	}

	if cond {

		c.pc.increment()
		c.sp.decrement()

		if err := c.mmu.Write(c.sp.get(), c.pc.highByte()); err != nil {
			return err
		}

		c.sp.decrement()

		if err := c.mmu.Write(c.sp.get(), c.pc.lowByte()); err != nil {
			return err
		}

		c.pc.set(im16 - 1)
	}

	return nil
}

// JpCondA16 jumps to the immediate
// 16bit address if the contidion is met
func (c *Core) JpCondA16(cond bool) error {

	im16, err := c.loadImmediate16()

	if err != nil {
		return err
	}

	if cond {
		c.pc.set(im16 - 1)
	}

	return nil
}

// RetCond returns to the saved pc if the condition
// is met
func (c *Core) RetCond(cond bool) error {

	if cond {

		if v, err := c.mmu.Read(c.sp.get()); err == nil {
			c.pc.setLow(v)
		} else {
			return err
		}

		c.sp.increment()

		if v, err := c.mmu.Read(c.sp.get()); err == nil {
			c.pc.setHigh(v)
		} else {
			return err
		}

		c.sp.increment()
		c.pc.decrement()
	}

	return nil
}

// JrCondR8 jumps to the immediate 8bit offset
// if the condition is met
func (c *Core) JrCondR8(cond bool) error {

	im8, err := c.loadImmediate8()

	if err != nil {
		return err
	}

	if cond {
		c.pc.set(c.pc.get() + uint16(int8(im8)))
	}

	return nil
}

// DecR decrements the register by 1
func (c *Core) DecR(dst *Register8) {
	c.setSubtractFlag(true)
	c.setHalfCarryFlagSub8(dst.get(), 1)
	dst.decrement()
	c.setZeroFlag(dst.get() == 0)
}

// IncR increments the register by 1
func (c *Core) IncR(dst *Register8) {
	c.setSubtractFlag(false)
	c.setHalfCarryFlagAdd8(dst.get(), 1)
	dst.increment()
	c.setZeroFlag(dst.get() == 0)
}

// CpN tests is register A == n
func (c *Core) CpN(n byte) {
	c.setSubtractFlag(true)
	c.setHalfCarryFlagSub8(c.a.get(), n)
	c.setCarryFlagSub8(c.a.get(), n)
	c.setZeroFlag(c.a.get()-n == 0)
}

// OrN sets A to A or n (bitwise)
func (c *Core) OrN(n byte) {
	c.setSubtractFlag(false)
	c.setHalfCarryFlag(false)
	c.setCarryFlag(false)
	c.a.set(c.a.get() | n)
	c.setZeroFlag(c.a.get() == 0)
}

// XorN sets A to A xor n (bitwise)
func (c *Core) XorN(n byte) {
	c.setSubtractFlag(false)
	c.setHalfCarryFlag(false)
	c.setCarryFlag(false)
	c.a.set(c.a.get() ^ n)
	c.setZeroFlag(c.a.get() == 0)
}

// AndN sets A to A and n (bitwise)
func (c *Core) AndN(n byte) {
	c.setSubtractFlag(false)
	c.setHalfCarryFlag(true)
	c.setCarryFlag(false)
	c.a.set(c.a.get() & n)
	c.setZeroFlag(c.a.get() == 0)
}

// SbcAn subtracts n and the carry flag from A
func (c *Core) SbcAn(n byte) {
	var cy byte

	if c.carryFlag() {
		cy = 1
	}

	c.setSubtractFlag(true)
	c.setHalfCarryFlagTripleSub8(c.a.get(), n, cy)
	c.setCarryFlagTripleSub8(c.a.get(), n, cy)
	c.a.set(c.a.get() - n - cy)
	c.setZeroFlag(c.a.get() == 0)
}

// SubN subtracts n from A
func (c *Core) SubN(n byte) {
	c.setSubtractFlag(true)
	c.setHalfCarryFlagSub8(c.a.get(), n)
	c.setCarryFlagSub8(c.a.get(), n)
	c.a.set(c.a.get() - n)
	c.setZeroFlag(c.a.get() == 0)
}

// AdcAn adds n and the carry flag to A
func (c *Core) AdcAn(n byte) {

	var cy byte

	if c.carryFlag() {
		cy = 1
	}

	c.setSubtractFlag(false)
	c.setHalfCarryFlagTripleAdd8(c.a.get(), n, cy)
	c.setCarryFlagTripleAdd8(c.a.get(), n, cy)
	c.a.set(c.a.get() + n + cy)
	c.setZeroFlag(c.a.get() == 0)
}

// AddAn adds n to A
func (c *Core) AddAn(n byte) {
	c.setSubtractFlag(false)
	c.setHalfCarryFlagAdd8(c.a.get(), n)
	c.setCarryFlagAdd8(c.a.get(), n)
	c.a.set(c.a.get() + n)
	c.setZeroFlag(c.a.get() == 0)
}

// LdRm loads (addr) to the register
func (c *Core) LdRm(dst *Register8, addr uint16) error {

	v, err := c.mmu.Read(addr)

	if err != nil {
		return err
	}

	dst.set(v)

	return nil
}

// LdRd8 loads the immediate 8 bit value to the register
func (c *Core) LdRd8(dst *Register8) error {

	im8, err := c.loadImmediate8()

	if err != nil {
		return err
	}

	dst.set(im8)

	return nil
}

// LdRrD16 loads the immediate 16 bit value to the register
func (c *Core) LdRrD16(dst *Register16) error {

	im16, err := c.loadImmediate16()

	if err != nil {
		return err
	}

	dst.set(im16)

	return nil
}

// PopRr pops a value from the stack into the register
func (c *Core) PopRr(dst *Register16) error {

	if v, err := c.mmu.Read(c.sp.get()); err == nil {
		dst.setLow(v)
	} else {
		return err
	}

	c.sp.increment()

	if v, err := c.mmu.Read(c.sp.get()); err == nil {
		dst.setHigh(v)
	} else {
		return err
	}

	c.sp.increment()

	return nil
}

// PushRr pushes the register to the stack
func (c *Core) PushRr(src *Register16) error {

	c.sp.decrement()

	if err := c.mmu.Write(c.sp.get(), src.highByte()); err != nil {
		return err
	}

	c.sp.decrement()

	if err := c.mmu.Write(c.sp.get(), src.lowByte()); err != nil {
		return err
	}

	return nil
}

// IncRr increments the register by 1
func (c *Core) IncRr(dst *Register16) error {
	dst.increment()
	return nil
}

// DecRr decrements the register by 1
func (c *Core) DecRr(dst *Register16) error {
	dst.decrement()
	return nil
}

// AddHlRr adds the register to HL
func (c *Core) AddHlRr(src *Register16) error {

	hl := c.hl.get()
	s := src.get()
	c.setSubtractFlag(false)
	c.setHalfCarryFlagAdd16(hl, s)
	c.setCarryFlagAdd16(hl, s)
	c.hl.set(hl + s)

	return nil
}
