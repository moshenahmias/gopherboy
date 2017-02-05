/*H**********************************************************************
* FILENAME :        instructions.go
*
* PACKAGE :			cpu
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package cpu

// initInstructions 00 - FF
func (c *Core) initInstructions() {

	/////////////////////////////////
	// Misc / control instructions //
	/////////////////////////////////

	// NOP
	c.instructions[0x00] = func() (int, int, string, error) {
		return 1, 4, "NOP", nil
	}

	// DI
	c.instructions[0xF3] = func() (int, int, string, error) {
		c.ime = false
		return 1, 4, "DI", nil
	}

	// EI
	c.instructions[0xFB] = func() (int, int, string, error) {
		c.ime = true
		return 1, 4, "EI", nil
	}

	// PREFIX CB
	c.instructions[0xCB] = func() (int, int, string, error) {

		im8, err := c.loadImmediate8()

		if err != nil {
			return 1, 4, "PREFIX CB", err
		}

		ins := c.instructionsCB[im8]

		if ins == nil {
			return 1, 4, "PREFIX CB", noSuchInstructionError(im8)
		}

		return ins()
	}

	// STOP 0
	c.instructions[0x10] = func() (int, int, string, error) {
		c.pc.increment()
		c.stop = true
		return 2, 4, "STOP", nil
	}

	// HALT
	c.instructions[0x76] = func() (int, int, string, error) {
		c.halt = true
		return 1, 4, "HALT", nil
	}

	////////////////////////////////////////////
	// 16bit load / store / move instructions //
	////////////////////////////////////////////

	// LD BC, d16
	c.instructions[0x01] = func() (int, int, string, error) {
		return 3, 12, "LD BC, d16", c.LdRrD16(&c.bc)
	}

	// LD DE, d16
	c.instructions[0x11] = func() (int, int, string, error) {
		return 3, 12, "LD DE, d16", c.LdRrD16(&c.de)
	}

	// LD HL, d16
	c.instructions[0x21] = func() (int, int, string, error) {
		return 3, 12, "LD HL, d16", c.LdRrD16(&c.hl)
	}

	// LD SP, d16
	c.instructions[0x31] = func() (int, int, string, error) {
		return 3, 12, "LD SP, d16", c.LdRrD16(&c.sp)
	}

	// LD (a16), SP
	c.instructions[0x08] = func() (int, int, string, error) {

		im16, err := c.loadImmediate16()

		if err != nil {
			return 3, 20, "LD (a16), SP", err
		}

		if err := c.mmu.Write(im16, c.sp.lowByte()); err != nil {
			return 3, 20, "LD (a16), SP", err
		}

		return 3, 20, "LD (a16), SP", c.mmu.Write(im16+1, c.sp.highByte())
	}

	// LD HL, SP + r8
	c.instructions[0xF8] = func() (int, int, string, error) {

		im8, err := c.loadImmediate8()

		if err != nil {
			return 2, 12, "LD HL, SP + r8", err
		}

		c.setSubtractFlag(false)
		c.setZeroFlag(false)
		c.setCarryFlagAdd8(c.sp.lowByte(), im8)
		c.setHalfCarryFlagAdd8(c.sp.lowByte(), im8)
		c.hl.set(c.sp.get() + uint16(int8(im8)))

		return 2, 12, "LD HL, SP + r8", nil
	}

	// LD SP, HL
	c.instructions[0xF9] = func() (int, int, string, error) {
		c.sp.set(c.hl.get())
		return 1, 8, "LD SP, HL", nil
	}

	// POP BC
	c.instructions[0xC1] = func() (int, int, string, error) {
		return 1, 12, "POP BC", c.PopRr(&c.bc)
	}

	// POP DE
	c.instructions[0xD1] = func() (int, int, string, error) {
		return 1, 12, "POP DE", c.PopRr(&c.de)
	}

	// POP HL
	c.instructions[0xE1] = func() (int, int, string, error) {
		return 1, 12, "POP HL", c.PopRr(&c.hl)
	}

	// POP AF
	c.instructions[0xF1] = func() (int, int, string, error) {

		f := c.f.get()

		if err := c.PopRr(&c.af); err != nil {
			return 1, 12, "POP AF", err
		}

		c.f.set((c.f.get() & 0xF0) | (f & 0x0F))

		return 1, 12, "POP AF", nil
	}

	// PUSH BC
	c.instructions[0xC5] = func() (int, int, string, error) {
		return 1, 16, "PUSH BC", c.PushRr(&c.bc)
	}

	// PUSH DE
	c.instructions[0xD5] = func() (int, int, string, error) {
		return 1, 16, "PUSH DE", c.PushRr(&c.de)
	}

	// PUSH HL
	c.instructions[0xE5] = func() (int, int, string, error) {
		return 1, 16, "PUSH HL", c.PushRr(&c.hl)
	}

	// PUSH AF
	c.instructions[0xF5] = func() (int, int, string, error) {
		return 1, 16, "PUSH AF", c.PushRr(&c.af)
	}

	/////////////////////////////////////////////
	// 16bit arithmetic / logical instructions //
	/////////////////////////////////////////////

	// INC BC
	c.instructions[0x03] = func() (int, int, string, error) {
		return 1, 8, "INC BC", c.IncRr(&c.bc)
	}

	// INC DE
	c.instructions[0x13] = func() (int, int, string, error) {
		return 1, 8, "INC DE", c.IncRr(&c.de)
	}

	// INC HL
	c.instructions[0x23] = func() (int, int, string, error) {
		return 1, 8, "INC HL", c.IncRr(&c.hl)
	}

	// INC SP
	c.instructions[0x33] = func() (int, int, string, error) {
		return 1, 8, "INC SP", c.IncRr(&c.sp)
	}

	// DEC BC
	c.instructions[0x0B] = func() (int, int, string, error) {
		return 1, 8, "DEC BC", c.DecRr(&c.bc)
	}

	// DEC DE
	c.instructions[0x1B] = func() (int, int, string, error) {
		return 1, 8, "DEC DE", c.DecRr(&c.de)
	}

	// DEC HL
	c.instructions[0x2B] = func() (int, int, string, error) {
		return 1, 8, "DEC HL", c.DecRr(&c.hl)
	}

	// DEC SP
	c.instructions[0x3B] = func() (int, int, string, error) {
		return 1, 8, "DEC SP", c.DecRr(&c.sp)
	}

	// ADD HL, BC
	c.instructions[0x09] = func() (int, int, string, error) {
		return 1, 8, "ADD HL, BC", c.AddHlRr(&c.bc)
	}

	// ADD HL, DE
	c.instructions[0x19] = func() (int, int, string, error) {
		return 1, 8, "ADD HL, DE", c.AddHlRr(&c.de)
	}

	// ADD HL, HL
	c.instructions[0x29] = func() (int, int, string, error) {
		return 1, 8, "ADD HL, HL", c.AddHlRr(&c.hl)
	}

	// ADD HL, SP
	c.instructions[0x39] = func() (int, int, string, error) {
		return 1, 8, "ADD HL, SP", c.AddHlRr(&c.sp)
	}

	// ADD SP, r8
	c.instructions[0xE8] = func() (int, int, string, error) {

		im8, err := c.loadImmediate8()

		if err != nil {
			return 2, 16, "ADD SP, r8", err
		}

		c.setSubtractFlag(false)
		c.setZeroFlag(false)
		c.setHalfCarryFlagAdd8(c.sp.lowByte(), im8)
		c.setCarryFlagAdd8(c.sp.lowByte(), im8)
		c.sp.set(c.sp.get() + uint16(int8(im8)))

		return 2, 16, "ADD SP, r8", nil
	}

	///////////////////////////////////////////
	// 8bit load / store / move instructions //
	///////////////////////////////////////////

	// LD (BC), A
	c.instructions[0x02] = func() (int, int, string, error) {
		return 1, 8, "LD (BC), A", c.mmu.Write(c.bc.get(), c.a.get())
	}

	// LD (DE), A
	c.instructions[0x12] = func() (int, int, string, error) {
		return 1, 8, "LD (DE), A", c.mmu.Write(c.de.get(), c.a.get())
	}

	// LD (HL+), A
	c.instructions[0x22] = func() (int, int, string, error) {

		if err := c.mmu.Write(c.hl.get(), c.a.get()); err != nil {
			return 1, 8, "LD (HL+), A", err
		}

		c.hl.increment()

		return 1, 8, "LD (HL+), A", nil
	}

	// LD (HL-), A
	c.instructions[0x32] = func() (int, int, string, error) {

		if err := c.mmu.Write(c.hl.get(), c.a.get()); err != nil {
			return 1, 8, "LD (HL-), A", err
		}

		c.hl.decrement()

		return 1, 8, "LD (HL-), A", nil
	}

	// LD B, d8
	c.instructions[0x06] = func() (int, int, string, error) {
		return 2, 8, "LD B, d8", c.LdRd8(c.b)
	}

	// LD C, d8
	c.instructions[0x0E] = func() (int, int, string, error) {
		return 2, 8, "LD C, d8", c.LdRd8(c.c)
	}

	// LD D, d8
	c.instructions[0x16] = func() (int, int, string, error) {
		return 2, 8, "LD D, d8", c.LdRd8(c.d)
	}

	// LD E, d8
	c.instructions[0x1E] = func() (int, int, string, error) {
		return 2, 8, "LD E, d8", c.LdRd8(c.e)
	}

	// LD H, d8
	c.instructions[0x26] = func() (int, int, string, error) {
		return 2, 8, "LD H, d8", c.LdRd8(c.h)
	}

	// LD L, d8
	c.instructions[0x2E] = func() (int, int, string, error) {
		return 2, 8, "LD L, d8", c.LdRd8(c.l)
	}

	// LD A, d8
	c.instructions[0x3E] = func() (int, int, string, error) {
		return 2, 8, "LD A, d8", c.LdRd8(c.a)
	}

	// LD (HL), d8
	c.instructions[0x36] = func() (int, int, string, error) {

		im8, err := c.loadImmediate8()

		if err != nil {
			return 2, 12, "LD (HL), d8", err
		}

		return 2, 12, "LD (HL), d8", c.mmu.Write(c.hl.get(), im8)
	}

	// LD A, (BC)
	c.instructions[0x0A] = func() (int, int, string, error) {
		return 1, 8, "LD A, (BC)", c.LdRm(c.a, c.bc.get())
	}

	// LD A, (DE)
	c.instructions[0x1A] = func() (int, int, string, error) {
		return 1, 8, "LD A, (DE)", c.LdRm(c.a, c.de.get())
	}

	// LD A, (HL+)
	c.instructions[0x2A] = func() (int, int, string, error) {

		if err := c.LdRm(c.a, c.hl.get()); err != nil {
			return 1, 8, "LD A, (HL+)", err
		}

		c.hl.increment()

		return 1, 8, "LD A, (HL+)", nil
	}

	// LD A, (HL-)
	c.instructions[0x3A] = func() (int, int, string, error) {

		if err := c.LdRm(c.a, c.hl.get()); err != nil {
			return 1, 8, "LD A, (HL-)", err
		}

		c.hl.decrement()

		return 1, 8, "LD A, (HL-)", nil
	}

	// LD B, (HL)
	c.instructions[0x46] = func() (int, int, string, error) {
		return 1, 8, "LD B, (HL)", c.LdRm(c.b, c.hl.get())
	}

	// LD C, (HL)
	c.instructions[0x4E] = func() (int, int, string, error) {
		return 1, 8, "LD C, (HL)", c.LdRm(c.c, c.hl.get())
	}

	// LD D, (HL)
	c.instructions[0x56] = func() (int, int, string, error) {
		return 1, 8, "LD D, (HL)", c.LdRm(c.d, c.hl.get())
	}

	// LD E, (HL)
	c.instructions[0x5E] = func() (int, int, string, error) {
		return 1, 8, "LD E, (HL)", c.LdRm(c.e, c.hl.get())
	}

	// LD H, (HL)
	c.instructions[0x66] = func() (int, int, string, error) {
		return 1, 8, "LD H, (HL)", c.LdRm(c.h, c.hl.get())
	}

	// LD L, (HL)
	c.instructions[0x6E] = func() (int, int, string, error) {
		return 1, 8, "LD L, (HL)", c.LdRm(c.l, c.hl.get())
	}

	// LD A, (HL)
	c.instructions[0x7E] = func() (int, int, string, error) {
		return 1, 8, "LD A, (HL)", c.LdRm(c.a, c.hl.get())
	}

	// LD B, B
	c.instructions[0x40] = func() (int, int, string, error) {
		c.b.set(c.b.get())
		return 1, 4, "LD B, B", nil
	}

	// LD B, C
	c.instructions[0x41] = func() (int, int, string, error) {
		c.b.set(c.c.get())
		return 1, 4, "LD B, C", nil
	}

	// LD B, D
	c.instructions[0x42] = func() (int, int, string, error) {
		c.b.set(c.d.get())
		return 1, 4, "LD B, D", nil
	}

	// LD B, E
	c.instructions[0x43] = func() (int, int, string, error) {
		c.b.set(c.e.get())
		return 1, 4, "LD B, E", nil
	}

	// LD B, H
	c.instructions[0x44] = func() (int, int, string, error) {
		c.b.set(c.h.get())
		return 1, 4, "LD B, H", nil
	}

	// LD B, L
	c.instructions[0x45] = func() (int, int, string, error) {
		c.b.set(c.l.get())
		return 1, 4, "LD B, L", nil
	}

	// LD B, A
	c.instructions[0x47] = func() (int, int, string, error) {
		c.b.set(c.a.get())
		return 1, 4, "LD B, A", nil
	}

	// LD C, B
	c.instructions[0x48] = func() (int, int, string, error) {
		c.c.set(c.b.get())
		return 1, 4, "LD C, B", nil
	}

	// LD C, C
	c.instructions[0x49] = func() (int, int, string, error) {
		c.c.set(c.c.get())
		return 1, 4, "LD C, C", nil
	}

	// LD C, D
	c.instructions[0x4A] = func() (int, int, string, error) {
		c.c.set(c.d.get())
		return 1, 4, "LD C, D", nil
	}

	// LD C, E
	c.instructions[0x4B] = func() (int, int, string, error) {
		c.c.set(c.e.get())
		return 1, 4, "LD C, E", nil
	}

	// LD C, H
	c.instructions[0x4C] = func() (int, int, string, error) {
		c.c.set(c.h.get())
		return 1, 4, "LD C, H", nil
	}

	// LD C, L
	c.instructions[0x4D] = func() (int, int, string, error) {
		c.c.set(c.l.get())
		return 1, 4, "LD C, L", nil
	}

	// LD C, A
	c.instructions[0x4F] = func() (int, int, string, error) {
		c.c.set(c.a.get())
		return 1, 4, "LD C, A", nil
	}

	// LD D, B
	c.instructions[0x50] = func() (int, int, string, error) {
		c.d.set(c.b.get())
		return 1, 4, "LD D, B", nil
	}

	// LD D, C
	c.instructions[0x51] = func() (int, int, string, error) {
		c.d.set(c.c.get())
		return 1, 4, "LD D, C", nil
	}

	// LD D, D
	c.instructions[0x52] = func() (int, int, string, error) {
		c.d.set(c.d.get())
		return 1, 4, "LD D, D", nil
	}

	// LD D, E
	c.instructions[0x53] = func() (int, int, string, error) {
		c.d.set(c.e.get())
		return 1, 4, "LD D, E", nil
	}

	// LD D, H
	c.instructions[0x54] = func() (int, int, string, error) {
		c.d.set(c.h.get())
		return 1, 4, "LD D, H", nil
	}

	// LD D, L
	c.instructions[0x55] = func() (int, int, string, error) {
		c.d.set(c.l.get())
		return 1, 4, "LD D, L", nil
	}

	// LD D, A
	c.instructions[0x57] = func() (int, int, string, error) {
		c.d.set(c.a.get())
		return 1, 4, "LD D, A", nil
	}

	// LD E, B
	c.instructions[0x58] = func() (int, int, string, error) {
		c.e.set(c.b.get())
		return 1, 4, "LD E, B", nil
	}

	// LD E, C
	c.instructions[0x59] = func() (int, int, string, error) {
		c.e.set(c.c.get())
		return 1, 4, "LD E, C", nil
	}

	// LD E, D
	c.instructions[0x5A] = func() (int, int, string, error) {
		c.e.set(c.d.get())
		return 1, 4, "LD E, D", nil
	}

	// LD E, E
	c.instructions[0x5B] = func() (int, int, string, error) {
		c.e.set(c.e.get())
		return 1, 4, "LD E, E", nil
	}

	// LD E, H
	c.instructions[0x5C] = func() (int, int, string, error) {
		c.e.set(c.h.get())
		return 1, 4, "LD E, H", nil
	}

	// LD E, L
	c.instructions[0x5D] = func() (int, int, string, error) {
		c.e.set(c.l.get())
		return 1, 4, "LD E, L", nil
	}

	// LD E, A
	c.instructions[0x5F] = func() (int, int, string, error) {
		c.e.set(c.a.get())
		return 1, 4, "LD E, A", nil
	}

	// LD H, B
	c.instructions[0x60] = func() (int, int, string, error) {
		c.h.set(c.b.get())
		return 1, 4, "LD H, B", nil
	}

	// LD H, C
	c.instructions[0x61] = func() (int, int, string, error) {
		c.h.set(c.c.get())
		return 1, 4, "LD H, C", nil
	}

	// LD H, D
	c.instructions[0x62] = func() (int, int, string, error) {
		c.h.set(c.d.get())
		return 1, 4, "LD H, D", nil
	}

	// LD H, E
	c.instructions[0x63] = func() (int, int, string, error) {
		c.h.set(c.e.get())
		return 1, 4, "LD H, E", nil
	}

	// LD H, H
	c.instructions[0x64] = func() (int, int, string, error) {
		c.h.set(c.h.get())
		return 1, 4, "LD H, H", nil
	}

	// LD H, L
	c.instructions[0x65] = func() (int, int, string, error) {
		c.h.set(c.l.get())
		return 1, 4, "LD H, L", nil
	}

	// LD H, A
	c.instructions[0x67] = func() (int, int, string, error) {
		c.h.set(c.a.get())
		return 1, 4, "LD H, A", nil
	}

	// LD L, B
	c.instructions[0x68] = func() (int, int, string, error) {
		c.l.set(c.b.get())
		return 1, 4, "LD L, B", nil
	}

	// LD L, C
	c.instructions[0x69] = func() (int, int, string, error) {
		c.l.set(c.c.get())
		return 1, 4, "LD L, C", nil
	}

	// LD L, D
	c.instructions[0x6A] = func() (int, int, string, error) {
		c.l.set(c.d.get())
		return 1, 4, "LD L, D", nil
	}

	// LD L, E
	c.instructions[0x6B] = func() (int, int, string, error) {
		c.l.set(c.e.get())
		return 1, 4, "LD L, E", nil
	}

	// LD L, H
	c.instructions[0x6C] = func() (int, int, string, error) {
		c.l.set(c.h.get())
		return 1, 4, "LD L, H", nil
	}

	// LD L, L
	c.instructions[0x6D] = func() (int, int, string, error) {
		c.l.set(c.l.get())
		return 1, 4, "LD L, L", nil
	}

	// LD L, A
	c.instructions[0x6F] = func() (int, int, string, error) {
		c.l.set(c.a.get())
		return 1, 4, "LD L, A", nil
	}

	// LD A, B
	c.instructions[0x78] = func() (int, int, string, error) {
		c.a.set(c.b.get())
		return 1, 4, "LD A, B", nil
	}

	// LD A, C
	c.instructions[0x79] = func() (int, int, string, error) {
		c.a.set(c.c.get())
		return 1, 4, "LD A C", nil
	}

	// LD A, D
	c.instructions[0x7A] = func() (int, int, string, error) {
		c.a.set(c.d.get())
		return 1, 4, "LD A, D", nil
	}

	// LD A, E
	c.instructions[0x7B] = func() (int, int, string, error) {
		c.a.set(c.e.get())
		return 1, 4, "LD A, E", nil
	}

	// LD A, H
	c.instructions[0x7C] = func() (int, int, string, error) {
		c.a.set(c.h.get())
		return 1, 4, "LD A, H", nil
	}

	// LD A, L
	c.instructions[0x7D] = func() (int, int, string, error) {
		c.a.set(c.l.get())
		return 1, 4, "LD A, L", nil
	}

	// LD A, A
	c.instructions[0x7F] = func() (int, int, string, error) {
		c.a.set(c.a.get())
		return 1, 4, "LD A, A", nil
	}

	// LD B, (HL)
	c.instructions[0x46] = func() (int, int, string, error) {
		return 1, 8, "LD B, (HL)", c.LdRm(c.b, c.hl.get())
	}

	// LD C, (HL)
	c.instructions[0x4E] = func() (int, int, string, error) {
		return 1, 8, "LD C, (HL)", c.LdRm(c.c, c.hl.get())
	}

	// LD D, (HL)
	c.instructions[0x56] = func() (int, int, string, error) {
		return 1, 8, "LD D, (HL)", c.LdRm(c.d, c.hl.get())
	}

	// LD E, (HL)
	c.instructions[0x5E] = func() (int, int, string, error) {
		return 1, 8, "LD E, (HL)", c.LdRm(c.e, c.hl.get())
	}

	// LD H, (HL)
	c.instructions[0x66] = func() (int, int, string, error) {
		return 1, 8, "LD H, (HL)", c.LdRm(c.h, c.hl.get())
	}

	// LD L, (HL)
	c.instructions[0x6E] = func() (int, int, string, error) {
		return 1, 8, "LD L, (HL)", c.LdRm(c.l, c.hl.get())
	}

	// LD A, (HL)
	c.instructions[0x7E] = func() (int, int, string, error) {
		return 1, 8, "LD A, (HL)", c.LdRm(c.a, c.hl.get())
	}

	// LD (HL), B
	c.instructions[0x70] = func() (int, int, string, error) {
		return 1, 8, "LD (HL), B", c.mmu.Write(c.hl.get(), c.b.get())
	}

	// LD (HL), C
	c.instructions[0x71] = func() (int, int, string, error) {
		return 1, 8, "LD (HL), C", c.mmu.Write(c.hl.get(), c.c.get())
	}

	// LD (HL), D
	c.instructions[0x72] = func() (int, int, string, error) {
		return 1, 8, "LD (HL), D", c.mmu.Write(c.hl.get(), c.d.get())
	}

	// LD (HL), E
	c.instructions[0x73] = func() (int, int, string, error) {
		return 1, 8, "LD (HL), E", c.mmu.Write(c.hl.get(), c.e.get())
	}

	// LD (HL), H
	c.instructions[0x74] = func() (int, int, string, error) {
		return 1, 8, "LD (HL), H", c.mmu.Write(c.hl.get(), c.h.get())
	}

	// LD (HL), L
	c.instructions[0x75] = func() (int, int, string, error) {
		return 1, 8, "LD (HL), L", c.mmu.Write(c.hl.get(), c.l.get())
	}

	// LD (HL), A
	c.instructions[0x77] = func() (int, int, string, error) {
		return 1, 8, "LD (HL), A", c.mmu.Write(c.hl.get(), c.a.get())
	}

	// LDH (a8), A
	c.instructions[0xE0] = func() (int, int, string, error) {

		im8, err := c.loadImmediate8()

		if err != nil {
			return 2, 12, "LDH (a8), A", err
		}

		return 2, 12, "LDH (a8), A", c.mmu.Write(0xFF00+uint16(im8), c.a.get())
	}

	// LDH A, (a8)
	c.instructions[0xF0] = func() (int, int, string, error) {

		im8, err := c.loadImmediate8()

		if err != nil {
			return 2, 12, "LDH A, (a8)", err
		}

		v, err := c.mmu.Read(0xFF00 + uint16(im8))

		if err != nil {
			return 2, 12, "LDH A, (a8)", err
		}

		c.a.set(v)

		return 2, 12, "LDH A, (a8)", nil
	}

	// LD (C), A
	c.instructions[0xE2] = func() (int, int, string, error) {
		return 2, 8, "LD (C), A", c.mmu.Write(0xFF00+uint16(c.c.get()), c.a.get())
	}

	// LD A, (C)
	c.instructions[0xF2] = func() (int, int, string, error) {

		v, err := c.mmu.Read(0xFF00 + uint16(c.c.get()))

		if err != nil {
			return 1, 8, "LD A, (C)", err
		}

		c.a.set(v)

		return 2, 8, "LD A, (C)", nil
	}

	// LD (a16), A
	c.instructions[0xEA] = func() (int, int, string, error) {

		im16, err := c.loadImmediate16()

		if err != nil {
			return 3, 16, "LD (a16), A", err
		}

		return 3, 16, "LD (a16), A", c.mmu.Write(im16, c.a.get())
	}

	// LD A, (a16)
	c.instructions[0xFA] = func() (int, int, string, error) {

		im16, err := c.loadImmediate16()

		if err != nil {
			return 3, 16, "LD (a16), A", err
		}

		return 3, 16, "LD A, (a16)", c.LdRm(c.a, im16)
	}

	////////////////////////////////////////////
	// 8bit arithmetic / logical instructions //
	////////////////////////////////////////////

	// ADD A, B
	c.instructions[0x80] = func() (int, int, string, error) {
		c.AddAn(c.b.get())
		return 1, 4, "ADD A, B", nil
	}

	// ADD A, C
	c.instructions[0x81] = func() (int, int, string, error) {
		c.AddAn(c.c.get())
		return 1, 4, "ADD A, C", nil
	}

	// ADD A, D
	c.instructions[0x82] = func() (int, int, string, error) {
		c.AddAn(c.d.get())
		return 1, 4, "ADD A, D", nil
	}

	// ADD A, E
	c.instructions[0x83] = func() (int, int, string, error) {
		c.AddAn(c.e.get())
		return 1, 4, "ADD A, E", nil
	}

	// ADD A, H
	c.instructions[0x84] = func() (int, int, string, error) {
		c.AddAn(c.h.get())
		return 1, 4, "ADD A, H", nil
	}

	// ADD A, L
	c.instructions[0x85] = func() (int, int, string, error) {
		c.AddAn(c.l.get())
		return 1, 4, "ADD A, L", nil
	}

	// ADD A, (HL)
	c.instructions[0x86] = func() (int, int, string, error) {

		v, err := c.mmu.Read(c.hl.get())

		if err != nil {
			return 1, 8, "ADD A, (HL)", err
		}

		c.AddAn(v)

		return 1, 8, "ADD A, (HL)", nil
	}

	// ADD A, A
	c.instructions[0x87] = func() (int, int, string, error) {
		c.AddAn(c.a.get())
		return 1, 4, "ADD A, A", nil
	}

	// ADD A, d8
	c.instructions[0xC6] = func() (int, int, string, error) {

		im8, err := c.loadImmediate8()

		if err != nil {
			return 2, 8, "ADD A, d8", err
		}

		c.AddAn(im8)

		return 2, 8, "ADD A, d8", nil
	}

	// ADC A, B
	c.instructions[0x88] = func() (int, int, string, error) {
		c.AdcAn(c.b.get())
		return 1, 4, "ADC A, B", nil
	}

	// ADC A, C
	c.instructions[0x89] = func() (int, int, string, error) {
		c.AdcAn(c.c.get())
		return 1, 4, "ADC A, C", nil
	}

	// ADC A, D
	c.instructions[0x8A] = func() (int, int, string, error) {
		c.AdcAn(c.d.get())
		return 1, 4, "ADC A, D", nil
	}

	// ADC A, E
	c.instructions[0x8B] = func() (int, int, string, error) {
		c.AdcAn(c.e.get())
		return 1, 4, "ADC A, E", nil
	}

	// ADC A, H
	c.instructions[0x8C] = func() (int, int, string, error) {
		c.AdcAn(c.h.get())
		return 1, 4, "ADC A, H", nil
	}

	// ADC A, L
	c.instructions[0x8D] = func() (int, int, string, error) {
		c.AdcAn(c.l.get())
		return 1, 4, "ADC A, L", nil
	}

	// ADC A, (HL)
	c.instructions[0x8E] = func() (int, int, string, error) {

		v, err := c.mmu.Read(c.hl.get())

		if err != nil {
			return 1, 8, "ADC A, (HL)", err
		}

		c.AdcAn(v)

		return 1, 8, "ADC A, (HL)", nil
	}

	// ADC A, A
	c.instructions[0x8F] = func() (int, int, string, error) {
		c.AdcAn(c.a.get())
		return 1, 4, "ADC A, A", nil
	}

	// ADC A, d8
	c.instructions[0xCE] = func() (int, int, string, error) {

		im8, err := c.loadImmediate8()

		if err != nil {
			return 2, 8, "ADC A, d8", err
		}

		c.AdcAn(im8)

		return 2, 8, "ADC A, d8", nil
	}

	// SUB B
	c.instructions[0x90] = func() (int, int, string, error) {
		c.SubN(c.b.get())
		return 1, 4, "SUB B", nil
	}

	// SUB C
	c.instructions[0x91] = func() (int, int, string, error) {
		c.SubN(c.c.get())
		return 1, 4, "SUB C", nil
	}

	// SUB D
	c.instructions[0x92] = func() (int, int, string, error) {
		c.SubN(c.d.get())
		return 1, 4, "SUB D", nil
	}

	// SUB E
	c.instructions[0x93] = func() (int, int, string, error) {
		c.SubN(c.e.get())
		return 1, 4, "SUB E", nil
	}

	// SUB H
	c.instructions[0x94] = func() (int, int, string, error) {
		c.SubN(c.h.get())
		return 1, 4, "SUB H", nil
	}

	// SUB L
	c.instructions[0x95] = func() (int, int, string, error) {
		c.SubN(c.l.get())
		return 1, 4, "SUB L", nil
	}

	// SUB (HL)
	c.instructions[0x96] = func() (int, int, string, error) {

		v, err := c.mmu.Read(c.hl.get())

		if err != nil {
			return 1, 8, "SUB (HL)", err
		}

		c.SubN(v)

		return 1, 8, "SUB (HL)", nil
	}

	// SUB A
	c.instructions[0x97] = func() (int, int, string, error) {
		c.SubN(c.a.get())
		return 1, 4, "SUB A", nil
	}

	// SUB d8
	c.instructions[0xD6] = func() (int, int, string, error) {

		im8, err := c.loadImmediate8()

		if err != nil {
			return 2, 8, "SUB d8", err
		}

		c.SubN(im8)

		return 2, 8, "SUB d8", nil
	}

	// SBC A, B
	c.instructions[0x98] = func() (int, int, string, error) {
		c.SbcAn(c.b.get())
		return 1, 4, "SBC A, B", nil
	}

	// SBC A, C
	c.instructions[0x99] = func() (int, int, string, error) {
		c.SbcAn(c.c.get())
		return 1, 4, "SBC A, C", nil
	}

	// SBC A, D
	c.instructions[0x9A] = func() (int, int, string, error) {
		c.SbcAn(c.d.get())
		return 1, 4, "SBC A, D", nil
	}

	// SBC A, E
	c.instructions[0x9B] = func() (int, int, string, error) {
		c.SbcAn(c.e.get())
		return 1, 4, "SBC A, E", nil
	}

	// SBC A, H
	c.instructions[0x9C] = func() (int, int, string, error) {
		c.SbcAn(c.h.get())
		return 1, 4, "SBC A, H", nil
	}

	// SBC A, L
	c.instructions[0x9D] = func() (int, int, string, error) {
		c.SbcAn(c.l.get())
		return 1, 4, "SBC A, L", nil
	}

	// SBC A, (HL)
	c.instructions[0x9E] = func() (int, int, string, error) {

		v, err := c.mmu.Read(c.hl.get())

		if err != nil {
			return 1, 8, "SBC A, (HL)", err
		}

		c.SbcAn(v)

		return 1, 8, "SBC A, (HL)", nil
	}

	// SBC A, A
	c.instructions[0x9F] = func() (int, int, string, error) {
		c.SbcAn(c.a.get())
		return 1, 4, "SBC A, A", nil
	}

	// SBC A, d8
	c.instructions[0xDE] = func() (int, int, string, error) {

		im8, err := c.loadImmediate8()

		if err != nil {
			return 2, 8, "SBC A, d8", err
		}

		c.SbcAn(im8)

		return 2, 8, "SBC A, d8", nil
	}

	// AND B
	c.instructions[0xA0] = func() (int, int, string, error) {
		c.AndN(c.b.get())
		return 1, 4, "AND B", nil
	}

	// AND C
	c.instructions[0xA1] = func() (int, int, string, error) {
		c.AndN(c.c.get())
		return 1, 4, "AND C", nil
	}

	// AND D
	c.instructions[0xA2] = func() (int, int, string, error) {
		c.AndN(c.d.get())
		return 1, 4, "AND D", nil
	}

	// AND E
	c.instructions[0xA3] = func() (int, int, string, error) {
		c.AndN(c.e.get())
		return 1, 4, "AND E", nil
	}

	// AND H
	c.instructions[0xA4] = func() (int, int, string, error) {
		c.AndN(c.h.get())
		return 1, 4, "AND H", nil
	}

	// AND L
	c.instructions[0xA5] = func() (int, int, string, error) {
		c.AndN(c.l.get())
		return 1, 4, "AND L", nil
	}

	// AND (HL)
	c.instructions[0xA6] = func() (int, int, string, error) {

		v, err := c.mmu.Read(c.hl.get())

		if err != nil {
			return 1, 8, "AND (HL)", err
		}

		c.AndN(v)

		return 1, 8, "AND (HL)", nil
	}

	// AND A
	c.instructions[0xA7] = func() (int, int, string, error) {
		c.AndN(c.a.get())
		return 1, 4, "AND A", nil
	}

	// AND d8
	c.instructions[0xE6] = func() (int, int, string, error) {

		im8, err := c.loadImmediate8()

		if err != nil {
			return 2, 8, "AND d8", err
		}

		c.AndN(im8)

		return 2, 8, "AND d8", nil
	}

	// XOR B
	c.instructions[0xA8] = func() (int, int, string, error) {
		c.XorN(c.b.get())
		return 1, 4, "XOR B", nil
	}

	// XOR C
	c.instructions[0xA9] = func() (int, int, string, error) {
		c.XorN(c.c.get())
		return 1, 4, "XOR C", nil
	}

	// XOR D
	c.instructions[0xAA] = func() (int, int, string, error) {
		c.XorN(c.d.get())
		return 1, 4, "XOR D", nil
	}

	// XOR E
	c.instructions[0xAB] = func() (int, int, string, error) {
		c.XorN(c.e.get())
		return 1, 4, "XOR E", nil
	}

	// XOR H
	c.instructions[0xAC] = func() (int, int, string, error) {
		c.XorN(c.h.get())
		return 1, 4, "XOR H", nil
	}

	// XOR L
	c.instructions[0xAD] = func() (int, int, string, error) {
		c.XorN(c.l.get())
		return 1, 4, "XOR L", nil
	}

	// XOR (HL)
	c.instructions[0xAE] = func() (int, int, string, error) {

		v, err := c.mmu.Read(c.hl.get())

		if err != nil {
			return 1, 8, "XOR (HL)", err
		}

		c.XorN(v)

		return 1, 8, "XOR (HL)", nil
	}

	// XOR A
	c.instructions[0xAF] = func() (int, int, string, error) {
		c.XorN(c.a.get())
		return 1, 4, "XOR A", nil
	}

	// XOR d8
	c.instructions[0xEE] = func() (int, int, string, error) {

		im8, err := c.loadImmediate8()

		if err != nil {
			return 2, 8, "XOR d8", err
		}

		c.XorN(im8)

		return 2, 8, "XOR d8", nil
	}

	// OR B
	c.instructions[0xB0] = func() (int, int, string, error) {
		c.OrN(c.b.get())
		return 1, 4, "OR B", nil
	}

	// OR C
	c.instructions[0xB1] = func() (int, int, string, error) {
		c.OrN(c.c.get())
		return 1, 4, "OR C", nil
	}

	// OR D
	c.instructions[0xB2] = func() (int, int, string, error) {
		c.OrN(c.d.get())
		return 1, 4, "OR D", nil
	}

	// OR E
	c.instructions[0xB3] = func() (int, int, string, error) {
		c.OrN(c.e.get())
		return 1, 4, "OR E", nil
	}

	// OR H
	c.instructions[0xB4] = func() (int, int, string, error) {
		c.OrN(c.h.get())
		return 1, 4, "OR H", nil
	}

	// OR L
	c.instructions[0xB5] = func() (int, int, string, error) {
		c.OrN(c.l.get())
		return 1, 4, "OR L", nil
	}

	// OR (HL)
	c.instructions[0xB6] = func() (int, int, string, error) {

		v, err := c.mmu.Read(c.hl.get())

		if err != nil {
			return 1, 8, "OR (HL)", err
		}

		c.OrN(v)

		return 1, 8, "OR (HL)", nil
	}

	// OR A
	c.instructions[0xB7] = func() (int, int, string, error) {
		c.OrN(c.a.get())
		return 1, 4, "OR A", nil
	}

	// OR d8
	c.instructions[0xF6] = func() (int, int, string, error) {

		im8, err := c.loadImmediate8()

		if err != nil {
			return 2, 8, "OR d8", err
		}

		c.OrN(im8)

		return 2, 8, "OR d8", nil
	}

	// CP B
	c.instructions[0xB8] = func() (int, int, string, error) {
		c.CpN(c.b.get())
		return 1, 4, "CP B", nil
	}

	// CP C
	c.instructions[0xB9] = func() (int, int, string, error) {
		c.CpN(c.c.get())
		return 1, 4, "CP C", nil
	}

	// CP D
	c.instructions[0xBA] = func() (int, int, string, error) {
		c.CpN(c.d.get())
		return 1, 4, "CP D", nil
	}

	// CP E
	c.instructions[0xBB] = func() (int, int, string, error) {
		c.CpN(c.e.get())
		return 1, 4, "CP E", nil
	}

	// CP H
	c.instructions[0xBC] = func() (int, int, string, error) {
		c.CpN(c.h.get())
		return 1, 4, "CP H", nil
	}

	// CP L
	c.instructions[0xBD] = func() (int, int, string, error) {
		c.CpN(c.l.get())
		return 1, 4, "CP L", nil
	}

	// CP (HL)
	c.instructions[0xBE] = func() (int, int, string, error) {

		v, err := c.mmu.Read(c.hl.get())

		if err != nil {
			return 1, 8, "CP (HL)", err
		}

		c.CpN(v)

		return 1, 8, "CP (HL)", nil
	}

	// CP A
	c.instructions[0xBF] = func() (int, int, string, error) {
		c.CpN(c.a.get())
		return 1, 4, "CP A", nil
	}

	// CP d8
	c.instructions[0xFE] = func() (int, int, string, error) {

		im8, err := c.loadImmediate8()

		if err != nil {
			return 2, 8, "CP d8", err
		}

		c.CpN(im8)

		return 2, 8, "CP d8", nil
	}

	// INC B
	c.instructions[0x04] = func() (int, int, string, error) {
		c.IncR(c.b)
		return 1, 4, "INC B", nil
	}

	// INC C
	c.instructions[0x0C] = func() (int, int, string, error) {
		c.IncR(c.c)
		return 1, 4, "INC C", nil
	}

	// INC D
	c.instructions[0x14] = func() (int, int, string, error) {
		c.IncR(c.d)
		return 1, 4, "INC D", nil
	}

	// INC E
	c.instructions[0x1C] = func() (int, int, string, error) {
		c.IncR(c.e)
		return 1, 4, "INC E", nil
	}

	// INC H
	c.instructions[0x24] = func() (int, int, string, error) {
		c.IncR(c.h)
		return 1, 4, "INC H", nil
	}

	// INC L
	c.instructions[0x2C] = func() (int, int, string, error) {
		c.IncR(c.l)
		return 1, 4, "INC L", nil
	}

	// INC A
	c.instructions[0x3C] = func() (int, int, string, error) {
		c.IncR(c.a)
		return 1, 4, "INC A", nil
	}

	// INC (HL)
	c.instructions[0x34] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 1, 12, "INC (HL)", err
		}

		c.setSubtractFlag(false)
		c.setHalfCarryFlagAdd8(v, 1)
		v++
		c.setZeroFlag(v == 0)

		return 1, 12, "INC (HL)", c.mmu.Write(hl, v)
	}

	// DEC B
	c.instructions[0x05] = func() (int, int, string, error) {
		c.DecR(c.b)
		return 1, 4, "DEC B", nil
	}

	// DEC C
	c.instructions[0x0D] = func() (int, int, string, error) {
		c.DecR(c.c)
		return 1, 4, "DEC C", nil
	}

	// DEC D
	c.instructions[0x15] = func() (int, int, string, error) {
		c.DecR(c.d)
		return 1, 4, "DEC D", nil
	}

	// DEC E
	c.instructions[0x1D] = func() (int, int, string, error) {
		c.DecR(c.e)
		return 1, 4, "DEC E", nil
	}

	// DEC H
	c.instructions[0x25] = func() (int, int, string, error) {
		c.DecR(c.h)
		return 1, 4, "DEC H", nil
	}

	// DEC L
	c.instructions[0x2D] = func() (int, int, string, error) {
		c.DecR(c.l)
		return 1, 4, "DEC L", nil
	}

	// DEC A
	c.instructions[0x3D] = func() (int, int, string, error) {
		c.DecR(c.a)
		return 1, 4, "DEC A", nil
	}

	// DEC (HL)
	c.instructions[0x35] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 1, 12, "DEC (HL)", err
		}

		c.setSubtractFlag(true)
		c.setHalfCarryFlagSub8(v, 1)
		v--
		c.setZeroFlag(v == 0)

		return 1, 12, "DEC (HL)", c.mmu.Write(hl, v)
	}

	// DAA
	c.instructions[0x27] = func() (int, int, string, error) {

		a := uint16(c.a.get())

		if !c.subtractFlag() {

			if c.halfCarryFlag() || a&0x0F > 9 {

				a += 0x06
			}

			if c.carryFlag() || a > 0x9F {

				a += 0x60
			}

		} else {

			if c.halfCarryFlag() {

				a = (a - 6) & 0xFF
			}

			if c.carryFlag() {

				a -= 0x60
			}
		}

		if (a & 0x0100) == 0x0100 {
			c.setCarryFlag(true)
		}

		c.a.set(byte(a & 0x00FF))
		c.setHalfCarryFlag(false)

		c.setZeroFlag(c.a.get() == 0)

		return 1, 4, "DAA", nil
	}

	// SCF
	c.instructions[0x37] = func() (int, int, string, error) {

		c.setSubtractFlag(false)
		c.setHalfCarryFlag(false)
		c.setCarryFlag(true)

		return 1, 4, "SCF", nil
	}

	// CPL
	c.instructions[0x2F] = func() (int, int, string, error) {

		c.setHalfCarryFlag(true)
		c.setSubtractFlag(true)
		c.a.set(^c.a.get())

		return 1, 4, "CPL", nil
	}

	// CCF
	c.instructions[0x3F] = func() (int, int, string, error) {

		c.setSubtractFlag(false)
		c.setHalfCarryFlag(false)
		c.setCarryFlag(!c.carryFlag())

		return 1, 4, "CCF", nil
	}

	///////////////////
	// Jumps / calls //
	///////////////////

	// JR r8
	c.instructions[0x18] = func() (int, int, string, error) {
		return 2, 12, "JR r8", c.JrCondR8(true)
	}

	// JR NZ, r8
	c.instructions[0x20] = func() (int, int, string, error) {

		nz := !c.zeroFlag()

		if err := c.JrCondR8(nz); err != nil {
			return 2, 8, "JR NZ, r8", err
		}

		if nz {
			return 2, 12, "JR NZ, r8", nil
		}

		return 2, 8, "JR NZ, r8", nil
	}

	// JR Z, r8
	c.instructions[0x28] = func() (int, int, string, error) {

		z := c.zeroFlag()

		if err := c.JrCondR8(z); err != nil {
			return 2, 8, "JR Z, r8", err
		}

		if z {
			return 2, 12, "JR Z, r8", nil
		}

		return 2, 8, "JR Z, r8", nil
	}

	// JR NC, r8
	c.instructions[0x30] = func() (int, int, string, error) {

		nc := !c.carryFlag()

		if err := c.JrCondR8(nc); err != nil {
			return 2, 8, "JR NC, r8", err
		}

		if nc {
			return 2, 12, "JR NC, r8", nil
		}

		return 2, 8, "JR NC, r8", nil
	}

	// JR C, r8
	c.instructions[0x38] = func() (int, int, string, error) {

		carry := c.carryFlag()

		if err := c.JrCondR8(carry); err != nil {
			return 2, 8, "JR C, r8", err
		}

		if carry {
			return 2, 12, "JR C, r8", nil
		}

		return 2, 8, "JR C, r8", nil
	}

	// RET
	c.instructions[0xC9] = func() (int, int, string, error) {
		return 1, 16, "RET", c.RetCond(true)
	}

	// RETI
	c.instructions[0xD9] = func() (int, int, string, error) {

		if err := c.RetCond(true); err != nil {
			return 1, 16, "RETI", err
		}

		c.ime = true

		return 1, 16, "RETI", nil
	}

	// RET NZ
	c.instructions[0xC0] = func() (int, int, string, error) {

		nz := !c.zeroFlag()

		if err := c.RetCond(nz); err != nil {
			return 1, 8, "RET NZ", err
		}

		if nz {
			return 1, 20, "RET NZ", nil
		}

		return 1, 8, "RET NZ", nil
	}

	// RET Z
	c.instructions[0xC8] = func() (int, int, string, error) {

		z := c.zeroFlag()

		if err := c.RetCond(z); err != nil {
			return 1, 8, "RET Z", err
		}

		if z {
			return 1, 20, "RET Z", nil
		}

		return 1, 8, "RET Z", nil
	}

	// RET NC
	c.instructions[0xD0] = func() (int, int, string, error) {

		nc := !c.carryFlag()

		if err := c.RetCond(nc); err != nil {
			return 1, 8, "RET NC", err
		}

		if nc {
			return 1, 20, "RET NC", nil
		}

		return 1, 8, "RET NC", nil
	}

	// RET C
	c.instructions[0xD8] = func() (int, int, string, error) {

		carry := c.carryFlag()

		if err := c.RetCond(carry); err != nil {
			return 1, 8, "RET C", err
		}

		if carry {
			return 1, 20, "RET C", nil
		}

		return 1, 8, "RET C", nil
	}

	// JP a16
	c.instructions[0xC3] = func() (int, int, string, error) {
		return 3, 16, "JP a16", c.JpCondA16(true)
	}

	// JP NZ, a16
	c.instructions[0xC2] = func() (int, int, string, error) {

		nz := !c.zeroFlag()

		if err := c.JpCondA16(nz); err != nil {
			return 3, 12, "JP NZ, a16", err
		}

		if nz {
			return 3, 16, "JP NZ, a16", nil
		}

		return 3, 12, "JP NZ, a16", nil
	}

	// JP Z, a16
	c.instructions[0xCA] = func() (int, int, string, error) {

		z := c.zeroFlag()

		if err := c.JpCondA16(z); err != nil {
			return 3, 12, "JP Z, a16", err
		}

		if z {
			return 3, 16, "JP Z, a16", nil
		}

		return 3, 12, "JP Z, a16", nil
	}

	// JP NC, a16
	c.instructions[0xD2] = func() (int, int, string, error) {

		nc := !c.carryFlag()

		if err := c.JpCondA16(nc); err != nil {
			return 3, 12, "JP NC, a16", err
		}

		if nc {
			return 3, 16, "JP NC, a16", nil
		}

		return 3, 12, "JP NC, a16", nil
	}

	// JP C, a16
	c.instructions[0xDA] = func() (int, int, string, error) {

		carry := c.carryFlag()

		if err := c.JpCondA16(carry); err != nil {
			return 3, 12, "JP C, a16", err
		}

		if carry {
			return 3, 16, "JP C, a16", nil
		}

		return 3, 12, "JP C, a16", nil
	}

	// JP (HL)
	c.instructions[0xE9] = func() (int, int, string, error) {
		c.pc.set(c.hl.get() - 1)
		return 1, 4, "JP (HL)", nil
	}

	// CALL a16
	c.instructions[0xCD] = func() (int, int, string, error) {
		return 3, 24, "CALL a16", c.CallCondA16(true)
	}

	// CALL NZ, a16
	c.instructions[0xC4] = func() (int, int, string, error) {

		nz := !c.zeroFlag()

		if err := c.CallCondA16(nz); err != nil {
			return 3, 12, "CALL NZ, a16", err
		}

		if nz {
			return 3, 24, "CALL NZ, a16", nil
		}

		return 3, 12, "CALL NZ, a16", nil
	}

	// CALL Z, a16
	c.instructions[0xCC] = func() (int, int, string, error) {

		z := c.zeroFlag()

		if err := c.CallCondA16(z); err != nil {
			return 3, 12, "CALL Z, a16", err
		}

		if z {
			return 3, 24, "CALL Z, a16", nil
		}

		return 3, 12, "CALL Z, a16", nil
	}

	// CALL NC, a16
	c.instructions[0xD4] = func() (int, int, string, error) {

		nc := !c.carryFlag()

		if err := c.CallCondA16(nc); err != nil {
			return 3, 12, "CALL NC, a16", err
		}

		if nc {
			return 3, 24, "CALL NC, a16", nil
		}

		return 3, 12, "CALL NC, a16", nil
	}

	// CALL C, a16
	c.instructions[0xDC] = func() (int, int, string, error) {

		carry := c.carryFlag()

		if err := c.CallCondA16(carry); err != nil {
			return 3, 12, "CALL C, a16", err
		}

		if carry {
			return 3, 24, "CALL C, a16", nil
		}

		return 3, 12, "CALL C, a16", nil
	}

	// RST 00H
	c.instructions[0xC7] = func() (int, int, string, error) {
		return 1, 16, "RST 00H", c.Rst(0x0000)
	}

	// RST 08H
	c.instructions[0xCF] = func() (int, int, string, error) {
		return 1, 16, "RST 08H", c.Rst(0x0008)
	}

	// RST 10H
	c.instructions[0xD7] = func() (int, int, string, error) {
		return 1, 16, "RST 10H", c.Rst(0x0010)
	}

	// RST 18H
	c.instructions[0xDF] = func() (int, int, string, error) {
		return 1, 16, "RST 18H", c.Rst(0x0018)
	}

	// RST 20H
	c.instructions[0xE7] = func() (int, int, string, error) {
		return 1, 16, "RST 20H", c.Rst(0x0020)
	}

	// RST 28H
	c.instructions[0xEF] = func() (int, int, string, error) {
		return 1, 16, "RST 28H", c.Rst(0x0028)
	}

	// RST 30H
	c.instructions[0xF7] = func() (int, int, string, error) {
		return 1, 16, "RST 30H", c.Rst(0x0030)
	}

	// RST 38H
	c.instructions[0xFF] = func() (int, int, string, error) {
		return 1, 16, "RST 38H", c.Rst(0x0038)
	}

	/////////////////////////////
	// 8bit rotations / shifts //
	/////////////////////////////

	// RLCA
	c.instructions[0x07] = func() (int, int, string, error) {
		c.RlcR(c.a, false)
		return 1, 4, "RLCA", nil
	}

	// RRCA
	c.instructions[0x0F] = func() (int, int, string, error) {
		c.RrcR(c.a, false)
		return 1, 4, "RRCA", nil
	}

	// RLA
	c.instructions[0x17] = func() (int, int, string, error) {
		c.RlR(c.a, false)
		return 1, 4, "RLA", nil
	}

	// RRA
	c.instructions[0x1F] = func() (int, int, string, error) {
		c.RrR(c.a, false)
		return 1, 4, "RRA", nil
	}
}
