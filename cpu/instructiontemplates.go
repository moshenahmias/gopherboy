package cpu

// insSetNr sets bit N in the register
func (c *Core) insSetNr(n uint, dst *Register8) {
	bit := (byte(0x01) << n)
	dst.set(dst.get() | bit)
}

// insResNr resets bit N in the register
func (c *Core) insResNr(n uint, dst *Register8) {
	bit := ^(byte(0x01) << n)
	dst.set(dst.get() & bit)
}

// insBitNr tests bit N in the register
func (c *Core) insBitNr(n uint, dst *Register8) {
	c.setHalfCarryFlag(true)
	c.setSubtractFlag(false)
	bit := byte(0x01) << n
	c.setZeroFlag(dst.get()&bit == 0)
}

// insSinsRlR shifts the register right one bit position.
// the contents of bit 0 are copied to the carry flag
// and a zero is put into bit 7.
func (c *Core) insSinsRlR(dst *Register8) {
	c.setHalfCarryFlag(false)
	c.setSubtractFlag(false)
	dst.rotateRight()
	r := dst.get()
	c.setCarryFlag(r&0x80 == 0x80)
	dst.set(r & 0x7F)
	c.setZeroFlag(dst.get() == 0)
}

// insSwapR high 4 bits with the low 4 bits
func (c *Core) insSwapR(dst *Register8) {
	c.setHalfCarryFlag(false)
	c.setSubtractFlag(false)
	c.setCarryFlag(false)
	r := dst.get()
	dst.set((r << 4) | (r >> 4))
	c.setZeroFlag(dst.get() == 0)
}

// insSlaR shifts the register left one bit position.
// the contents of bit 7 are copied to the carry flag
// and a zero is put into bit 0
func (c *Core) insSlaR(dst *Register8) {

	c.setHalfCarryFlag(false)
	c.setSubtractFlag(false)
	dst.rotateLeft()
	r := dst.get()
	c.setCarryFlag(r&0x01 == 0x01)
	dst.set(r & 0xFE)
	c.setZeroFlag(dst.get() == 0)
}

// insSraR shifts the register right one bit position.
// the contents of bit 0 are copied to the carry flag
// and the previous contents of bit 7 are unchanged.
func (c *Core) insSraR(dst *Register8) {

	c.setHalfCarryFlag(false)
	c.setSubtractFlag(false)
	dst.rotateRight()
	r := dst.get()
	c.setCarryFlag(r&0x80 == 0x80)
	dst.set((r & 0x7F) | ((r << 1) & 0x80))
	c.setZeroFlag(dst.get() == 0)
}

// insRrR rotates the register right one bit position.
// the contents of bit 0 are copied to the carry flag and
// the previous contents of the carry flag are copied to bit 7
func (c *Core) insRrR(dst *Register8, cb bool) {

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

// insRlR rotates the register left one bit position.
// the contents of bit 7 are copied to the carry flag and
// the previous contents of the carry flag are copied to bit 0.
func (c *Core) insRlR(dst *Register8, cb bool) {

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

// insRrcR rotates the register right one bit position.
// the contents of bit 0 are copied to the carry flag and bit 7.
func (c *Core) insRrcR(dst *Register8, cb bool) {

	c.setHalfCarryFlag(false)
	c.setSubtractFlag(false)
	dst.rotateRight()
	c.setCarryFlag(dst.get()&0x80 == 0x80)
	c.setZeroFlag(cb && dst.get() == 0)
}

// insRlcR rotates the register left one bit position.
// the contents of bit 7 are copied to the carry flag and bit 0.
func (c *Core) insRlcR(dst *Register8, cb bool) {

	c.setHalfCarryFlag(false)
	c.setSubtractFlag(false)
	c.setCarryFlag(dst.get()&0x80 == 0x80)
	dst.rotateLeft()
	c.setZeroFlag(cb && dst.get() == 0)
}

// insRst jumps to the given address
func (c *Core) insRst(addr uint16) error {

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

// insCallCondA16 calls the function at the immediate
// 16bit address if the contidion is met
func (c *Core) insCallCondA16(cond bool) error {

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

// insJpCondA16 jumps to the immediate
// 16bit address if the contidion is met
func (c *Core) insJpCondA16(cond bool) error {

	im16, err := c.loadImmediate16()

	if err != nil {
		return err
	}

	if cond {
		c.pc.set(im16 - 1)
	}

	return nil
}

// insRetCond returns to the saved pc if the condition
// is met
func (c *Core) insRetCond(cond bool) error {

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

// insJrCondR8 jumps to the immediate 8bit offset
// if the condition is met
func (c *Core) insJrCondR8(cond bool) error {

	im8, err := c.loadImmediate8()

	if err != nil {
		return err
	}

	if cond {
		c.pc.set(c.pc.get() + uint16(int8(im8)))
	}

	return nil
}

// insDecR decrements the register by 1
func (c *Core) insDecR(dst *Register8) {
	c.setSubtractFlag(true)
	c.setHalfCarryFlagSub8(dst.get(), 1)
	dst.decrement()
	c.setZeroFlag(dst.get() == 0)
}

// insIncR increments the register by 1
func (c *Core) insIncR(dst *Register8) {
	c.setSubtractFlag(false)
	c.setHalfCarryFlagAdd8(dst.get(), 1)
	dst.increment()
	c.setZeroFlag(dst.get() == 0)
}

// insCpN tests is register A == n
func (c *Core) insCpN(n byte) {
	c.setSubtractFlag(true)
	c.setHalfCarryFlagSub8(c.a.get(), n)
	c.setCarryFlagSub8(c.a.get(), n)
	c.setZeroFlag(c.a.get()-n == 0)
}

// insOrN sets A to A or n (bitwise)
func (c *Core) insOrN(n byte) {
	c.setSubtractFlag(false)
	c.setHalfCarryFlag(false)
	c.setCarryFlag(false)
	c.a.set(c.a.get() | n)
	c.setZeroFlag(c.a.get() == 0)
}

// insXorN sets A to A xor n (bitwise)
func (c *Core) insXorN(n byte) {
	c.setSubtractFlag(false)
	c.setHalfCarryFlag(false)
	c.setCarryFlag(false)
	c.a.set(c.a.get() ^ n)
	c.setZeroFlag(c.a.get() == 0)
}

// insAndN sets A to A and n (bitwise)
func (c *Core) insAndN(n byte) {
	c.setSubtractFlag(false)
	c.setHalfCarryFlag(true)
	c.setCarryFlag(false)
	c.a.set(c.a.get() & n)
	c.setZeroFlag(c.a.get() == 0)
}

// insSbcAn subtracts n and the carry flag from A
func (c *Core) insSbcAn(n byte) {
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

// insSubN subtracts n from A
func (c *Core) insSubN(n byte) {
	c.setSubtractFlag(true)
	c.setHalfCarryFlagSub8(c.a.get(), n)
	c.setCarryFlagSub8(c.a.get(), n)
	c.a.set(c.a.get() - n)
	c.setZeroFlag(c.a.get() == 0)
}

// insAdcAn adds n and the carry flag to A
func (c *Core) insAdcAn(n byte) {

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

// insAddAn adds n to A
func (c *Core) insAddAn(n byte) {
	c.setSubtractFlag(false)
	c.setHalfCarryFlagAdd8(c.a.get(), n)
	c.setCarryFlagAdd8(c.a.get(), n)
	c.a.set(c.a.get() + n)
	c.setZeroFlag(c.a.get() == 0)
}

// insLdRm loads (addr) to the register
func (c *Core) insLdRm(dst *Register8, addr uint16) error {

	v, err := c.mmu.Read(addr)

	if err != nil {
		return err
	}

	dst.set(v)

	return nil
}

// insLdRd8 loads the immediate 8 bit value to the register
func (c *Core) insLdRd8(dst *Register8) error {

	im8, err := c.loadImmediate8()

	if err != nil {
		return err
	}

	dst.set(im8)

	return nil
}

// insLdRrD16 loads the immediate 16 bit value to the register
func (c *Core) insLdRrD16(dst *Register16) error {

	im16, err := c.loadImmediate16()

	if err != nil {
		return err
	}

	dst.set(im16)

	return nil
}

// insPopRr pops a value from the stack into the register
func (c *Core) insPopRr(dst *Register16) error {

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

// insPushRr pushes the register to the stack
func (c *Core) insPushRr(src *Register16) error {

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

// insIncRr increments the register by 1
func (c *Core) insIncRr(dst *Register16) error {
	dst.increment()
	return nil
}

// insDecRr decrements the register by 1
func (c *Core) insDecRr(dst *Register16) error {
	dst.decrement()
	return nil
}

// insAddHlRr adds the register to HL
func (c *Core) insAddHlRr(src *Register16) error {

	hl := c.hl.get()
	s := src.get()
	c.setSubtractFlag(false)
	c.setHalfCarryFlagAdd16(hl, s)
	c.setCarryFlagAdd16(hl, s)
	c.hl.set(hl + s)

	return nil
}
