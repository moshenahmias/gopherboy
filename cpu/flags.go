package cpu

// zeroFlag returns the zero flag state
func (c *Core) zeroFlag() bool {
	return c.f.get()&0x80 == 0x80
}

// setZeroFlag sets the zero flag state
func (c *Core) setZeroFlag(state bool) {
	if state {
		c.f.set(c.f.get() | 0x80)
	} else {
		c.f.set(c.f.get() & 0x7F)
	}
}

// subtractFlag returns the subtract flag state
func (c *Core) subtractFlag() bool {
	return c.f.get()&0x40 == 0x40
}

// setSubtractFlag sets the subtract flag state
func (c *Core) setSubtractFlag(state bool) {
	if state {
		c.f.set(c.f.get() | 0x40)
	} else {
		c.f.set(c.f.get() & 0xBF)
	}
}

// halfCarryFlag returns the half carry flag state
func (c *Core) halfCarryFlag() bool {
	return c.f.get()&0x20 == 0x20
}

// setHalfCarryFlag sets the half carry flag state
func (c *Core) setHalfCarryFlag(state bool) {
	if state {
		c.f.set(c.f.get() | 0x20)
	} else {
		c.f.set(c.f.get() & 0xDF)
	}
}

// carryFlag returns the carry flag state
func (c *Core) carryFlag() bool {
	return c.f.get()&0x10 == 0x10
}

// setCarryFlag sets the carry flag state
func (c *Core) setCarryFlag(state bool) {
	if state {
		c.f.set(c.f.get() | 0x10)
	} else {
		c.f.set(c.f.get() & 0xEF)
	}
}

// setHalfCarryFlagAdd8 sets the half carry flag when
// 'x + y' results in carrying from bit 3 (otherwise reset)
func (c *Core) setHalfCarryFlagAdd8(x, y byte) {
	c.setHalfCarryFlag(((x&0x0F)+(y&0x0F))&0x10 == 0x10)
}

// setHalfCarryFlagTripleAdd8 sets the half carry flag when
// 'x + y + z' results in carrying from bit 3 (otherwise reset)
func (c *Core) setHalfCarryFlagTripleAdd8(x, y, z byte) {
	c.setHalfCarryFlag(((x&0x0F)+(y&0x0F))&0x10 == 0x10 || (((x+y)&0x0F)+(z&0x0F))&0x10 == 0x10)
}

// setHalfCarryFlagSub8 sets the half carry flag when
// 'x - y' results in borrowing to bit 3 (otherwise reset)
func (c *Core) setHalfCarryFlagSub8(x, y byte) {
	c.setHalfCarryFlag((x & 0x0F) < (y & 0x0F))
}

// setHalfCarryFlagTripleSub8 sets the half carry flag when
// 'x - y - z' results in borrowing to bit 3 (otherwise reset)
func (c *Core) setHalfCarryFlagTripleSub8(x, y, z byte) {
	c.setHalfCarryFlag((x&0x0F) < (y&0x0F) || ((x-y)&0x0F) < (z&0x0F))
}

// setHalfCarryFlagAdd16 sets the half carry flag when
// 'x + y' results in carrying from bit 11 (otherwise reset)
func (c *Core) setHalfCarryFlagAdd16(x, y uint16) {
	c.setHalfCarryFlag(((x&0x0FFF)+(y&0x0FFF))&0x1000 == 0x1000)
}

// setHalfCarryFlagTripleAdd16 sets the half carry flag when
// 'x + y + z' results in carrying from bit 11 (otherwise reset)
func (c *Core) setHalfCarryFlagTripleAdd16(x, y, z uint16) {
	c.setHalfCarryFlag(((x&0x0FFF)+(y&0x0FFF))&0x1000 == 0x1000 || (((x+y)&0x0FFF)+(z&0x0FFF))&0x1000 == 0x1000)
}

// setHalfCarryFlagSub16 sets the half carry flag when
// 'x - y' results in borrowing to bit 11 (otherwise reset)
func (c *Core) setHalfCarryFlagSub16(x, y uint16) {
	c.setHalfCarryFlag((x & 0x0FFF) < (y & 0x0FFF))
}

// setHalfCarryFlagTripleSub16 sets the half carry flag when
// 'x - y - z' results in borrowing to bit 11 (otherwise reset)
func (c *Core) setHalfCarryFlagTripleSub16(x, y, z uint16) {
	c.setHalfCarryFlag((x&0x0FFF) < (y&0x0FFF) || ((x-y)&0x0FFF) < (z&0x0FFF))
}

// setCarryFlagAdd8 sets the carry flag when
// 'x + y' results in carrying from bit 7 (otherwise reset)
func (c *Core) setCarryFlagAdd8(x, y byte) {
	c.setCarryFlag(x > 0xFF-y)
}

// setCarryFlagTripleAdd8 sets the carry flag when
// 'x + y + z' results in carrying from bit 7 (otherwise reset)
func (c *Core) setCarryFlagTripleAdd8(x, y, z byte) {
	c.setCarryFlag(x > 0xFF-y || (x+y) > 0xFF-z)
}

// setCarryFlagSub8 sets the carry flag when
// 'x - y' results in borrowing to bit 7 (otherwise reset)
func (c *Core) setCarryFlagSub8(x, y byte) {
	c.setCarryFlag(x < y)
}

// setCarryFlagTripleSub8 sets the carry flag when
// 'x - y - z' results in borrowing to bit 7 (otherwise reset)
func (c *Core) setCarryFlagTripleSub8(x, y, z byte) {
	c.setCarryFlag(x < y || (x-y) < z)
}

// setCarryFlagAdd16 sets the carry flag when
// 'x + y' results in carrying from bit 15 (otherwise reset)
func (c *Core) setCarryFlagAdd16(x, y uint16) {
	c.setCarryFlag(x > 0xFFFF-y)
}

// setCarryFlagTripleAdd16 sets the carry flag when
// 'x + y + z' results in carrying from bit 15 (otherwise reset)
func (c *Core) setCarryFlagTripleAdd16(x, y, z uint16) {
	c.setCarryFlag(x > 0xFFFF-y || (x+y) > 0xFFFF-z)
}

// setCarryFlagSub16 sets the carry flag when
// 'x - y' results in borrowing to bit 15 (otherwise reset)
func (c *Core) setCarryFlagSub16(x, y uint16) {
	c.setCarryFlag(x < y)
}

// setCarryFlagTripleSub16 sets the carry flag when
// 'x - y' results in borrowing to bit 15 (otherwise reset)
func (c *Core) setCarryFlagTripleSub16(x, y, z uint16) {
	c.setCarryFlag(x < y || (x-y) < z)
}
