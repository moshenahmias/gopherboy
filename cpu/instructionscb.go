/*H**********************************************************************
* FILENAME :        instructionscb.go
*
* PACKAGE :			cpu
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package cpu

// InitInstructionsCB 00 - FF
func (c *Core) initInstructionsCB() {

	//////////////////////////////////////////////////
	// 8bit rotations / shifts and bit instructions //
	//////////////////////////////////////////////////

	// RLC B
	c.instructionsCB[0x00] = func() (int, int, string, error) {
		c.insRlcR(c.b, true)
		return 2, 8, "RLC B", nil
	}

	// RLC C
	c.instructionsCB[0x01] = func() (int, int, string, error) {
		c.insRlcR(c.c, true)
		return 2, 8, "RLC C", nil
	}

	// RLC D
	c.instructionsCB[0x02] = func() (int, int, string, error) {
		c.insRlcR(c.d, true)
		return 2, 8, "RLC D", nil
	}

	// RLC E
	c.instructionsCB[0x03] = func() (int, int, string, error) {
		c.insRlcR(c.e, true)
		return 2, 8, "RLC E", nil
	}

	// RLC H
	c.instructionsCB[0x04] = func() (int, int, string, error) {
		c.insRlcR(c.h, true)
		return 2, 8, "RLC H", nil
	}

	// RLC L
	c.instructionsCB[0x05] = func() (int, int, string, error) {
		c.insRlcR(c.l, true)
		return 2, 8, "RLC L", nil
	}

	// RLC (HL)
	c.instructionsCB[0x06] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "RLC (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insRlcR(&r, true)

		return 2, 16, "RLC (HL)", c.mmu.Write(hl, r.get())
	}

	// RLC A
	c.instructionsCB[0x07] = func() (int, int, string, error) {
		c.insRlcR(c.a, true)
		return 2, 8, "RLC A", nil
	}

	// RRC B
	c.instructionsCB[0x08] = func() (int, int, string, error) {
		c.insRrcR(c.b, true)
		return 2, 8, "RRC B", nil
	}

	// RRC C
	c.instructionsCB[0x09] = func() (int, int, string, error) {
		c.insRrcR(c.c, true)
		return 2, 8, "RRC C", nil
	}

	// RRC D
	c.instructionsCB[0x0A] = func() (int, int, string, error) {
		c.insRrcR(c.d, true)
		return 2, 8, "RRC D", nil
	}

	// RRC E
	c.instructionsCB[0x0B] = func() (int, int, string, error) {
		c.insRrcR(c.e, true)
		return 2, 8, "RRC E", nil
	}

	// RRC H
	c.instructionsCB[0x0C] = func() (int, int, string, error) {
		c.insRrcR(c.h, true)
		return 2, 8, "RRC H", nil
	}

	// RRC L
	c.instructionsCB[0x0D] = func() (int, int, string, error) {
		c.insRrcR(c.l, true)
		return 2, 8, "RRC L", nil
	}

	// RRC (HL)
	c.instructionsCB[0x0E] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "RRC (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insRrcR(&r, true)

		return 2, 16, "RRC (HL)", c.mmu.Write(hl, r.get())
	}

	// RRC A
	c.instructionsCB[0x0F] = func() (int, int, string, error) {
		c.insRrcR(c.a, true)
		return 2, 8, "RRC A", nil
	}

	// RL B
	c.instructionsCB[0x10] = func() (int, int, string, error) {
		c.insRlR(c.b, true)
		return 2, 8, "RL B", nil
	}

	// RL C
	c.instructionsCB[0x11] = func() (int, int, string, error) {
		c.insRlR(c.c, true)
		return 2, 8, "RL C", nil
	}

	// RL D
	c.instructionsCB[0x12] = func() (int, int, string, error) {
		c.insRlR(c.d, true)
		return 2, 8, "RL D", nil
	}

	// RL E
	c.instructionsCB[0x13] = func() (int, int, string, error) {
		c.insRlR(c.e, true)
		return 2, 8, "RL E", nil
	}

	// RL H
	c.instructionsCB[0x14] = func() (int, int, string, error) {
		c.insRlR(c.h, true)
		return 2, 8, "RL H", nil
	}

	// RL L
	c.instructionsCB[0x15] = func() (int, int, string, error) {
		c.insRlR(c.l, true)
		return 2, 8, "RL L", nil
	}

	// RL (HL)
	c.instructionsCB[0x16] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "RL (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insRlR(&r, true)

		return 2, 16, "RL (HL)", c.mmu.Write(hl, r.get())
	}

	// RL A
	c.instructionsCB[0x17] = func() (int, int, string, error) {
		c.insRlR(c.a, true)
		return 2, 8, "RL A", nil
	}

	// RR B
	c.instructionsCB[0x18] = func() (int, int, string, error) {
		c.insRrR(c.b, true)
		return 2, 8, "RR B", nil
	}

	// RR C
	c.instructionsCB[0x19] = func() (int, int, string, error) {
		c.insRrR(c.c, true)
		return 2, 8, "RR C", nil
	}

	// RR D
	c.instructionsCB[0x1A] = func() (int, int, string, error) {
		c.insRrR(c.d, true)
		return 2, 8, "RR D", nil
	}

	// RR E
	c.instructionsCB[0x1B] = func() (int, int, string, error) {
		c.insRrR(c.e, true)
		return 2, 8, "RR E", nil
	}

	// RR H
	c.instructionsCB[0x1C] = func() (int, int, string, error) {
		c.insRrR(c.h, true)
		return 2, 8, "RR H", nil
	}

	// RR L
	c.instructionsCB[0x1D] = func() (int, int, string, error) {
		c.insRrR(c.l, true)
		return 2, 8, "RR L", nil
	}

	// RR (HL)
	c.instructionsCB[0x1E] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "RR (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insRrR(&r, true)

		return 2, 16, "RR (HL)", c.mmu.Write(hl, r.get())
	}

	// RR A
	c.instructionsCB[0x1F] = func() (int, int, string, error) {
		c.insRrR(c.a, true)
		return 2, 8, "RR A", nil
	}

	// SLA B
	c.instructionsCB[0x20] = func() (int, int, string, error) {
		c.insSlaR(c.b)
		return 2, 8, "SLA B", nil
	}

	// SLA C
	c.instructionsCB[0x21] = func() (int, int, string, error) {
		c.insSlaR(c.c)
		return 2, 8, "SLA C", nil
	}

	// SLA D
	c.instructionsCB[0x22] = func() (int, int, string, error) {
		c.insSlaR(c.d)
		return 2, 8, "SLA D", nil
	}

	// SLA E
	c.instructionsCB[0x23] = func() (int, int, string, error) {
		c.insSlaR(c.e)
		return 2, 8, "SLA E", nil
	}

	// SLA H
	c.instructionsCB[0x24] = func() (int, int, string, error) {
		c.insSlaR(c.h)
		return 2, 8, "SLA H", nil
	}

	// SLA L
	c.instructionsCB[0x25] = func() (int, int, string, error) {
		c.insSlaR(c.l)
		return 2, 8, "SLA L", nil
	}

	// SLA (HL)
	c.instructionsCB[0x26] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "SLA (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insSlaR(&r)

		return 2, 16, "SLA (HL)", c.mmu.Write(hl, r.get())
	}

	// SLA A
	c.instructionsCB[0x27] = func() (int, int, string, error) {
		c.insSlaR(c.a)
		return 2, 8, "SLA A", nil
	}

	// SRA B
	c.instructionsCB[0x28] = func() (int, int, string, error) {
		c.insSraR(c.b)
		return 2, 8, "SRA B", nil
	}

	// SRA C
	c.instructionsCB[0x29] = func() (int, int, string, error) {
		c.insSraR(c.c)
		return 2, 8, "SRA C", nil
	}

	// SRA D
	c.instructionsCB[0x2A] = func() (int, int, string, error) {
		c.insSraR(c.d)
		return 2, 8, "SRA D", nil
	}

	// SRA E
	c.instructionsCB[0x2B] = func() (int, int, string, error) {
		c.insSraR(c.e)
		return 2, 8, "SRA E", nil
	}

	// SRA H
	c.instructionsCB[0x2C] = func() (int, int, string, error) {
		c.insSraR(c.h)
		return 2, 8, "SRA H", nil
	}

	// SRA L
	c.instructionsCB[0x2D] = func() (int, int, string, error) {
		c.insSraR(c.l)
		return 2, 8, "SRA L", nil
	}

	// SRA (HL)
	c.instructionsCB[0x2E] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "SRA (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insSraR(&r)

		return 2, 16, "SRA (HL)", c.mmu.Write(hl, r.get())
	}

	// SRA A
	c.instructionsCB[0x2F] = func() (int, int, string, error) {
		c.insSraR(c.a)
		return 2, 8, "SRA A", nil
	}

	// SWAP B
	c.instructionsCB[0x30] = func() (int, int, string, error) {
		c.insSwapR(c.b)
		return 2, 8, "SWAP B", nil
	}

	// SWAP C
	c.instructionsCB[0x31] = func() (int, int, string, error) {
		c.insSwapR(c.c)
		return 2, 8, "SWAP C", nil
	}

	// SWAP D
	c.instructionsCB[0x32] = func() (int, int, string, error) {
		c.insSwapR(c.d)
		return 2, 8, "SWAP D", nil
	}

	// SWAP E
	c.instructionsCB[0x33] = func() (int, int, string, error) {
		c.insSwapR(c.e)
		return 2, 8, "SWAP E", nil
	}

	// SWAP H
	c.instructionsCB[0x34] = func() (int, int, string, error) {
		c.insSwapR(c.h)
		return 2, 8, "SWAP H", nil
	}

	// SWAP L
	c.instructionsCB[0x35] = func() (int, int, string, error) {
		c.insSwapR(c.l)
		return 2, 8, "SWAP L", nil
	}

	// SWAP (HL)
	c.instructionsCB[0x36] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "SWAP (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insSwapR(&r)

		return 2, 16, "SWAP (HL)", c.mmu.Write(hl, r.get())
	}

	// SWAP A
	c.instructionsCB[0x37] = func() (int, int, string, error) {
		c.insSwapR(c.a)
		return 2, 8, "SWAP A", nil
	}

	// SRL B
	c.instructionsCB[0x38] = func() (int, int, string, error) {
		c.insSinsRlR(c.b)
		return 2, 8, "SRL B", nil
	}

	// SRL C
	c.instructionsCB[0x39] = func() (int, int, string, error) {
		c.insSinsRlR(c.c)
		return 2, 8, "SRL C", nil
	}

	// SRL D
	c.instructionsCB[0x3A] = func() (int, int, string, error) {
		c.insSinsRlR(c.d)
		return 2, 8, "SRL D", nil
	}

	// SRL E
	c.instructionsCB[0x3B] = func() (int, int, string, error) {
		c.insSinsRlR(c.e)
		return 2, 8, "SRL E", nil
	}

	// SRL H
	c.instructionsCB[0x3C] = func() (int, int, string, error) {
		c.insSinsRlR(c.h)
		return 2, 8, "SRL H", nil
	}

	// SRL L
	c.instructionsCB[0x3D] = func() (int, int, string, error) {
		c.insSinsRlR(c.l)
		return 2, 8, "SRL L", nil
	}

	// SRL (HL)
	c.instructionsCB[0x3E] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "SRL (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insSinsRlR(&r)

		return 2, 16, "SRL (HL)", c.mmu.Write(hl, r.get())
	}

	// SRL A
	c.instructionsCB[0x3F] = func() (int, int, string, error) {
		c.insSinsRlR(c.a)
		return 2, 8, "SRL A", nil
	}

	// BIT 0, B
	c.instructionsCB[0x40] = func() (int, int, string, error) {
		c.insBitNr(0, c.b)
		return 2, 8, "BIT 0, B", nil
	}

	// BIT 0, C
	c.instructionsCB[0x41] = func() (int, int, string, error) {
		c.insBitNr(0, c.c)
		return 2, 8, "BIT 0, C", nil
	}

	// BIT 0, D
	c.instructionsCB[0x42] = func() (int, int, string, error) {
		c.insBitNr(0, c.d)
		return 2, 8, "BIT 0, D", nil
	}

	// BIT 0, E
	c.instructionsCB[0x43] = func() (int, int, string, error) {
		c.insBitNr(0, c.e)
		return 2, 8, "BIT 0, E", nil
	}

	// BIT 0, H
	c.instructionsCB[0x44] = func() (int, int, string, error) {
		c.insBitNr(0, c.h)
		return 2, 8, "BIT 0, H", nil
	}

	// BIT 0, L
	c.instructionsCB[0x45] = func() (int, int, string, error) {
		c.insBitNr(0, c.l)
		return 2, 8, "BIT 0, L", nil
	}

	// BIT 0, (HL)
	c.instructionsCB[0x46] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "BIT 0, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insBitNr(0, &r)

		return 2, 16, "BIT 0, (HL)", nil
	}

	// BIT 0, A
	c.instructionsCB[0x47] = func() (int, int, string, error) {
		c.insBitNr(0, c.a)
		return 2, 8, "BIT 0, A", nil
	}

	// BIT 1, B
	c.instructionsCB[0x48] = func() (int, int, string, error) {
		c.insBitNr(1, c.b)
		return 2, 8, "BIT 1, B", nil
	}

	// BIT 1, C
	c.instructionsCB[0x49] = func() (int, int, string, error) {
		c.insBitNr(1, c.c)
		return 2, 8, "BIT 1, C", nil
	}

	// BIT 1, D
	c.instructionsCB[0x4A] = func() (int, int, string, error) {
		c.insBitNr(1, c.d)
		return 2, 8, "BIT 1, D", nil
	}

	// BIT 1, E
	c.instructionsCB[0x4B] = func() (int, int, string, error) {
		c.insBitNr(1, c.e)
		return 2, 8, "BIT 1, E", nil
	}

	// BIT 1, H
	c.instructionsCB[0x4C] = func() (int, int, string, error) {
		c.insBitNr(1, c.h)
		return 2, 8, "BIT 1, H", nil
	}

	// BIT 1, L
	c.instructionsCB[0x4D] = func() (int, int, string, error) {
		c.insBitNr(1, c.l)
		return 2, 8, "BIT 1, L", nil
	}

	// BIT 1, (HL)
	c.instructionsCB[0x4E] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "BIT 1, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insBitNr(1, &r)

		return 2, 16, "BIT 1, (HL)", nil
	}

	// BIT 1, A
	c.instructionsCB[0x4F] = func() (int, int, string, error) {
		c.insBitNr(1, c.a)
		return 2, 8, "BIT 1, A", nil
	}

	// BIT 2, B
	c.instructionsCB[0x50] = func() (int, int, string, error) {
		c.insBitNr(2, c.b)
		return 2, 8, "BIT 2, B", nil
	}

	// BIT 2, C
	c.instructionsCB[0x51] = func() (int, int, string, error) {
		c.insBitNr(2, c.c)
		return 2, 8, "BIT 2, C", nil
	}

	// BIT 2, D
	c.instructionsCB[0x52] = func() (int, int, string, error) {
		c.insBitNr(2, c.d)
		return 2, 8, "BIT 2, D", nil
	}

	// BIT 2, E
	c.instructionsCB[0x53] = func() (int, int, string, error) {
		c.insBitNr(2, c.e)
		return 2, 8, "BIT 2, E", nil
	}

	// BIT 2, H
	c.instructionsCB[0x54] = func() (int, int, string, error) {
		c.insBitNr(2, c.h)
		return 2, 8, "BIT 2, H", nil
	}

	// BIT 2, L
	c.instructionsCB[0x55] = func() (int, int, string, error) {
		c.insBitNr(2, c.l)
		return 2, 8, "BIT 2, L", nil
	}

	// BIT 2, (HL)
	c.instructionsCB[0x56] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "BIT 2, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insBitNr(2, &r)

		return 2, 16, "BIT 2, (HL)", nil
	}

	// BIT 2, A
	c.instructionsCB[0x57] = func() (int, int, string, error) {
		c.insBitNr(2, c.a)
		return 2, 8, "BIT 2, A", nil
	}

	// BIT 3, B
	c.instructionsCB[0x58] = func() (int, int, string, error) {
		c.insBitNr(3, c.b)
		return 2, 8, "BIT 3, B", nil
	}

	// BIT 3, C
	c.instructionsCB[0x59] = func() (int, int, string, error) {
		c.insBitNr(3, c.c)
		return 2, 8, "BIT 3, C", nil
	}

	// BIT 3, D
	c.instructionsCB[0x5A] = func() (int, int, string, error) {
		c.insBitNr(3, c.d)
		return 2, 8, "BIT 3, D", nil
	}

	// BIT 3, E
	c.instructionsCB[0x5B] = func() (int, int, string, error) {
		c.insBitNr(3, c.e)
		return 2, 8, "BIT 3, E", nil
	}

	// BIT 3, H
	c.instructionsCB[0x5C] = func() (int, int, string, error) {
		c.insBitNr(3, c.h)
		return 2, 8, "BIT 3, H", nil
	}

	// BIT 3, L
	c.instructionsCB[0x5D] = func() (int, int, string, error) {
		c.insBitNr(3, c.l)
		return 2, 8, "BIT 3, L", nil
	}

	// BIT 3, (HL)
	c.instructionsCB[0x5E] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "BIT 3, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insBitNr(3, &r)

		return 2, 16, "BIT 3, (HL)", nil
	}

	// BIT 3, A
	c.instructionsCB[0x5F] = func() (int, int, string, error) {
		c.insBitNr(3, c.a)
		return 2, 8, "BIT 3, A", nil
	}

	// BIT 4, B
	c.instructionsCB[0x60] = func() (int, int, string, error) {
		c.insBitNr(4, c.b)
		return 2, 8, "BIT 4, B", nil
	}

	// BIT 4, C
	c.instructionsCB[0x61] = func() (int, int, string, error) {
		c.insBitNr(4, c.c)
		return 2, 8, "BIT 4, C", nil
	}

	// BIT 4, D
	c.instructionsCB[0x62] = func() (int, int, string, error) {
		c.insBitNr(4, c.d)
		return 2, 8, "BIT 4, D", nil
	}

	// BIT 4, E
	c.instructionsCB[0x63] = func() (int, int, string, error) {
		c.insBitNr(4, c.e)
		return 2, 8, "BIT 4, E", nil
	}

	// BIT 4, H
	c.instructionsCB[0x64] = func() (int, int, string, error) {
		c.insBitNr(4, c.h)
		return 2, 8, "BIT 4, H", nil
	}

	// BIT 4, L
	c.instructionsCB[0x65] = func() (int, int, string, error) {
		c.insBitNr(4, c.l)
		return 2, 8, "BIT 4, L", nil
	}

	// BIT 4, (HL)
	c.instructionsCB[0x66] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "BIT 4, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insBitNr(4, &r)

		return 2, 16, "BIT 4, (HL)", nil
	}

	// BIT 4, A
	c.instructionsCB[0x67] = func() (int, int, string, error) {
		c.insBitNr(4, c.a)
		return 2, 8, "BIT 4, A", nil
	}

	// BIT 5, B
	c.instructionsCB[0x68] = func() (int, int, string, error) {
		c.insBitNr(5, c.b)
		return 2, 8, "BIT 5, B", nil
	}

	// BIT 5, C
	c.instructionsCB[0x69] = func() (int, int, string, error) {
		c.insBitNr(5, c.c)
		return 2, 8, "BIT 5, C", nil
	}

	// BIT 5, D
	c.instructionsCB[0x6A] = func() (int, int, string, error) {
		c.insBitNr(5, c.d)
		return 2, 8, "BIT 5, D", nil
	}

	// BIT 5, E
	c.instructionsCB[0x6B] = func() (int, int, string, error) {
		c.insBitNr(5, c.e)
		return 2, 8, "BIT 5, E", nil
	}

	// BIT 5, H
	c.instructionsCB[0x6C] = func() (int, int, string, error) {
		c.insBitNr(5, c.h)
		return 2, 8, "BIT 5, H", nil
	}

	// BIT 5, L
	c.instructionsCB[0x6D] = func() (int, int, string, error) {
		c.insBitNr(5, c.l)
		return 2, 8, "BIT 5, L", nil
	}

	// BIT 5, (HL)
	c.instructionsCB[0x6E] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "BIT 5, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insBitNr(5, &r)

		return 2, 16, "BIT 5, (HL)", nil
	}

	// BIT 5, A
	c.instructionsCB[0x6F] = func() (int, int, string, error) {
		c.insBitNr(5, c.a)
		return 2, 8, "BIT 5, A", nil
	}

	// BIT 6, B
	c.instructionsCB[0x70] = func() (int, int, string, error) {
		c.insBitNr(6, c.b)
		return 2, 8, "BIT 6, B", nil
	}

	// BIT 6, C
	c.instructionsCB[0x71] = func() (int, int, string, error) {
		c.insBitNr(6, c.c)
		return 2, 8, "BIT 6, C", nil
	}

	// BIT 6, D
	c.instructionsCB[0x72] = func() (int, int, string, error) {
		c.insBitNr(6, c.d)
		return 2, 8, "BIT 6, D", nil
	}

	// BIT 6, E
	c.instructionsCB[0x73] = func() (int, int, string, error) {
		c.insBitNr(6, c.e)
		return 2, 8, "BIT 6, E", nil
	}

	// BIT 6, H
	c.instructionsCB[0x74] = func() (int, int, string, error) {
		c.insBitNr(6, c.h)
		return 2, 8, "BIT 6, H", nil
	}

	// BIT 6, L
	c.instructionsCB[0x75] = func() (int, int, string, error) {
		c.insBitNr(6, c.l)
		return 2, 8, "BIT 6, L", nil
	}

	// BIT 6, (HL)
	c.instructionsCB[0x76] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "BIT 6, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insBitNr(6, &r)

		return 2, 16, "BIT 6, (HL)", nil
	}

	// BIT 6, A
	c.instructionsCB[0x77] = func() (int, int, string, error) {
		c.insBitNr(6, c.a)
		return 2, 8, "BIT 6, A", nil
	}

	// BIT 7, B
	c.instructionsCB[0x78] = func() (int, int, string, error) {
		c.insBitNr(7, c.b)
		return 2, 8, "BIT 7, B", nil
	}

	// BIT 7, C
	c.instructionsCB[0x79] = func() (int, int, string, error) {
		c.insBitNr(7, c.c)
		return 2, 8, "BIT 7, C", nil
	}

	// BIT 7, D
	c.instructionsCB[0x7A] = func() (int, int, string, error) {
		c.insBitNr(7, c.d)
		return 2, 8, "BIT 7, D", nil
	}

	// BIT 7, E
	c.instructionsCB[0x7B] = func() (int, int, string, error) {
		c.insBitNr(7, c.e)
		return 2, 8, "BIT 7, E", nil
	}

	// BIT 7, H
	c.instructionsCB[0x7C] = func() (int, int, string, error) {
		c.insBitNr(7, c.h)
		return 2, 8, "BIT 7, H", nil
	}

	// BIT 7, L
	c.instructionsCB[0x7D] = func() (int, int, string, error) {
		c.insBitNr(7, c.l)
		return 2, 8, "BIT 7, L", nil
	}

	// BIT 7, (HL)
	c.instructionsCB[0x7E] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "BIT 7, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insBitNr(7, &r)

		return 2, 16, "BIT 7, (HL)", nil
	}

	// BIT 7, A
	c.instructionsCB[0x7F] = func() (int, int, string, error) {
		c.insBitNr(7, c.a)
		return 2, 8, "BIT 7, A", nil
	}

	// RES 0, B
	c.instructionsCB[0x80] = func() (int, int, string, error) {
		c.insResNr(0, c.b)
		return 2, 8, "RES 0, B", nil
	}

	// RES 0, C
	c.instructionsCB[0x81] = func() (int, int, string, error) {
		c.insResNr(0, c.c)
		return 2, 8, "RES 0, C", nil
	}

	// RES 0, D
	c.instructionsCB[0x82] = func() (int, int, string, error) {
		c.insResNr(0, c.d)
		return 2, 8, "RES 0, D", nil
	}

	// RES 0, E
	c.instructionsCB[0x83] = func() (int, int, string, error) {
		c.insResNr(0, c.e)
		return 2, 8, "RES 0, E", nil
	}

	// RES 0, H
	c.instructionsCB[0x84] = func() (int, int, string, error) {
		c.insResNr(0, c.h)
		return 2, 8, "RES 0, H", nil
	}

	// RES 0, L
	c.instructionsCB[0x85] = func() (int, int, string, error) {
		c.insResNr(0, c.l)
		return 2, 8, "RES 0, L", nil
	}

	// RES 0, (HL)
	c.instructionsCB[0x86] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "RES 0, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insResNr(0, &r)

		return 2, 16, "RES 0, (HL)", c.mmu.Write(hl, r.get())
	}

	// RES 0, A
	c.instructionsCB[0x87] = func() (int, int, string, error) {
		c.insResNr(0, c.a)
		return 2, 8, "RES 0, A", nil
	}

	// RES 1, B
	c.instructionsCB[0x88] = func() (int, int, string, error) {
		c.insResNr(1, c.b)
		return 2, 8, "RES 1, B", nil
	}

	// RES 1, C
	c.instructionsCB[0x89] = func() (int, int, string, error) {
		c.insResNr(1, c.c)
		return 2, 8, "RES 1, C", nil
	}

	// RES 1, D
	c.instructionsCB[0x8A] = func() (int, int, string, error) {
		c.insResNr(1, c.d)
		return 2, 8, "RES 1, D", nil
	}

	// RES 1, E
	c.instructionsCB[0x8B] = func() (int, int, string, error) {
		c.insResNr(1, c.e)
		return 2, 8, "RES 1, E", nil
	}

	// RES 1, H
	c.instructionsCB[0x8C] = func() (int, int, string, error) {
		c.insResNr(1, c.h)
		return 2, 8, "RES 1, H", nil
	}

	// RES 1, L
	c.instructionsCB[0x8D] = func() (int, int, string, error) {
		c.insResNr(1, c.l)
		return 2, 8, "RES 1, L", nil
	}

	// RES 1, (HL)
	c.instructionsCB[0x8E] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "RES 1, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insResNr(1, &r)

		return 2, 16, "RES 1, (HL)", c.mmu.Write(hl, r.get())
	}

	// RES 1, A
	c.instructionsCB[0x8F] = func() (int, int, string, error) {
		c.insResNr(1, c.a)
		return 2, 8, "RES 1, A", nil
	}

	// RES 2, B
	c.instructionsCB[0x90] = func() (int, int, string, error) {
		c.insResNr(2, c.b)
		return 2, 8, "RES 2, B", nil
	}

	// RES 2, C
	c.instructionsCB[0x91] = func() (int, int, string, error) {
		c.insResNr(2, c.c)
		return 2, 8, "RES 2, C", nil
	}

	// RES 2, D
	c.instructionsCB[0x92] = func() (int, int, string, error) {
		c.insResNr(2, c.d)
		return 2, 8, "RES 2, D", nil
	}

	// RES 2, E
	c.instructionsCB[0x93] = func() (int, int, string, error) {
		c.insResNr(2, c.e)
		return 2, 8, "RES 2, E", nil
	}

	// RES 2, H
	c.instructionsCB[0x94] = func() (int, int, string, error) {
		c.insResNr(2, c.h)
		return 2, 8, "RES 2, H", nil
	}

	// RES 2, L
	c.instructionsCB[0x95] = func() (int, int, string, error) {
		c.insResNr(2, c.l)
		return 2, 8, "RES 2, L", nil
	}

	// RES 2, (HL)
	c.instructionsCB[0x96] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "RES 2, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insResNr(2, &r)

		return 2, 16, "RES 2, (HL)", c.mmu.Write(hl, r.get())
	}

	// RES 2, A
	c.instructionsCB[0x97] = func() (int, int, string, error) {
		c.insResNr(2, c.a)
		return 2, 8, "RES 2, A", nil
	}

	// RES 3, B
	c.instructionsCB[0x98] = func() (int, int, string, error) {
		c.insResNr(3, c.b)
		return 2, 8, "RES 3, B", nil
	}

	// RES 3, C
	c.instructionsCB[0x99] = func() (int, int, string, error) {
		c.insResNr(3, c.c)
		return 2, 8, "RES 3, C", nil
	}

	// RES 3, D
	c.instructionsCB[0x9A] = func() (int, int, string, error) {
		c.insResNr(3, c.d)
		return 2, 8, "RES 3, D", nil
	}

	// RES 3, E
	c.instructionsCB[0x9B] = func() (int, int, string, error) {
		c.insResNr(3, c.e)
		return 2, 8, "RES 3, E", nil
	}

	// RES 3, H
	c.instructionsCB[0x9C] = func() (int, int, string, error) {
		c.insResNr(3, c.h)
		return 2, 8, "RES 3, H", nil
	}

	// RES 3, L
	c.instructionsCB[0x9D] = func() (int, int, string, error) {
		c.insResNr(3, c.l)
		return 2, 8, "RES 3, L", nil
	}

	// RES 3, (HL)
	c.instructionsCB[0x9E] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "RES 3, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insResNr(3, &r)

		return 2, 16, "RES 3, (HL)", c.mmu.Write(hl, r.get())
	}

	// RES 3, A
	c.instructionsCB[0x9F] = func() (int, int, string, error) {
		c.insResNr(3, c.a)
		return 2, 8, "RES 3, A", nil
	}

	// RES 4, B
	c.instructionsCB[0xA0] = func() (int, int, string, error) {
		c.insResNr(4, c.b)
		return 2, 8, "RES 4, B", nil
	}

	// RES 4, C
	c.instructionsCB[0xA1] = func() (int, int, string, error) {
		c.insResNr(4, c.c)
		return 2, 8, "RES 4, C", nil
	}

	// RES 4, D
	c.instructionsCB[0xA2] = func() (int, int, string, error) {
		c.insResNr(4, c.d)
		return 2, 8, "RES 4, D", nil
	}

	// RES 4, E
	c.instructionsCB[0xA3] = func() (int, int, string, error) {
		c.insResNr(4, c.e)
		return 2, 8, "RES 4, E", nil
	}

	// RES 4, H
	c.instructionsCB[0xA4] = func() (int, int, string, error) {
		c.insResNr(4, c.h)
		return 2, 8, "RES 4, H", nil
	}

	// RES 4, L
	c.instructionsCB[0xA5] = func() (int, int, string, error) {
		c.insResNr(4, c.l)
		return 2, 8, "RES 4, L", nil
	}

	// RES 4, (HL)
	c.instructionsCB[0xA6] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "RES 4, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insResNr(4, &r)

		return 2, 16, "RES 4, (HL)", c.mmu.Write(hl, r.get())
	}

	// RES 4, A
	c.instructionsCB[0xA7] = func() (int, int, string, error) {
		c.insResNr(4, c.a)
		return 2, 8, "RES 4, A", nil
	}

	// RES 5, B
	c.instructionsCB[0xA8] = func() (int, int, string, error) {
		c.insResNr(5, c.b)
		return 2, 8, "RES 5, B", nil
	}

	// RES 5, C
	c.instructionsCB[0xA9] = func() (int, int, string, error) {
		c.insResNr(5, c.c)
		return 2, 8, "RES 5, C", nil
	}

	// RES 5, D
	c.instructionsCB[0xAA] = func() (int, int, string, error) {
		c.insResNr(5, c.d)
		return 2, 8, "RES 5, D", nil
	}

	// RES 5, E
	c.instructionsCB[0xAB] = func() (int, int, string, error) {
		c.insResNr(5, c.e)
		return 2, 8, "RES 5, E", nil
	}

	// RES 5, H
	c.instructionsCB[0xAC] = func() (int, int, string, error) {
		c.insResNr(5, c.h)
		return 2, 8, "RES 5, H", nil
	}

	// RES 5, L
	c.instructionsCB[0xAD] = func() (int, int, string, error) {
		c.insResNr(5, c.l)
		return 2, 8, "RES 5, L", nil
	}

	// RES 5, (HL)
	c.instructionsCB[0xAE] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "RES 5, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insResNr(5, &r)

		return 2, 16, "RES 5, (HL)", c.mmu.Write(hl, r.get())
	}

	// RES 5, A
	c.instructionsCB[0xAF] = func() (int, int, string, error) {
		c.insResNr(5, c.a)
		return 2, 8, "RES 5, A", nil
	}

	// RES 6, B
	c.instructionsCB[0xB0] = func() (int, int, string, error) {
		c.insResNr(6, c.b)
		return 2, 8, "RES 6, B", nil
	}

	// RES 6, C
	c.instructionsCB[0xB1] = func() (int, int, string, error) {
		c.insResNr(6, c.c)
		return 2, 8, "RES 6, C", nil
	}

	// RES 6, D
	c.instructionsCB[0xB2] = func() (int, int, string, error) {
		c.insResNr(6, c.d)
		return 2, 8, "RES 6, D", nil
	}

	// RES 6, E
	c.instructionsCB[0xB3] = func() (int, int, string, error) {
		c.insResNr(6, c.e)
		return 2, 8, "RES 6, E", nil
	}

	// RES 6, H
	c.instructionsCB[0xB4] = func() (int, int, string, error) {
		c.insResNr(6, c.h)
		return 2, 8, "RES 6, H", nil
	}

	// RES 6, L
	c.instructionsCB[0xB5] = func() (int, int, string, error) {
		c.insResNr(6, c.l)
		return 2, 8, "RES 6, L", nil
	}

	// RES 6, (HL)
	c.instructionsCB[0xB6] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "RES 6, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insResNr(6, &r)

		return 2, 16, "RES 6, (HL)", c.mmu.Write(hl, r.get())
	}

	// RES 6, A
	c.instructionsCB[0xB7] = func() (int, int, string, error) {
		c.insResNr(6, c.a)
		return 2, 8, "RES 6, A", nil
	}

	// RES 7, B
	c.instructionsCB[0xB8] = func() (int, int, string, error) {
		c.insResNr(7, c.b)
		return 2, 8, "RES 7, B", nil
	}

	// RES 7, C
	c.instructionsCB[0xB9] = func() (int, int, string, error) {
		c.insResNr(7, c.c)
		return 2, 8, "RES 7, C", nil
	}

	// RES 7, D
	c.instructionsCB[0xBA] = func() (int, int, string, error) {
		c.insResNr(7, c.d)
		return 2, 8, "RES 7, D", nil
	}

	// RES 7, E
	c.instructionsCB[0xBB] = func() (int, int, string, error) {
		c.insResNr(7, c.e)
		return 2, 8, "RES 7, E", nil
	}

	// RES 7, H
	c.instructionsCB[0xBC] = func() (int, int, string, error) {
		c.insResNr(7, c.h)
		return 2, 8, "RES 7, H", nil
	}

	// RES 7, L
	c.instructionsCB[0xBD] = func() (int, int, string, error) {
		c.insResNr(7, c.l)
		return 2, 8, "RES 7, L", nil
	}

	// RES 7, (HL)
	c.instructionsCB[0xBE] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "RES 7, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insResNr(7, &r)

		return 2, 16, "RES 7, (HL)", c.mmu.Write(hl, r.get())
	}

	// RES 7, A
	c.instructionsCB[0xBF] = func() (int, int, string, error) {
		c.insResNr(7, c.a)
		return 2, 8, "RES 7, A", nil
	}

	// SET 0, B
	c.instructionsCB[0xC0] = func() (int, int, string, error) {
		c.insSetNr(0, c.b)
		return 2, 8, "SET 0, B", nil
	}

	// SET 0, C
	c.instructionsCB[0xC1] = func() (int, int, string, error) {
		c.insSetNr(0, c.c)
		return 2, 8, "SET 0, C", nil
	}

	// SET 0, D
	c.instructionsCB[0xC2] = func() (int, int, string, error) {
		c.insSetNr(0, c.d)
		return 2, 8, "SET 0, D", nil
	}

	// SET 0, E
	c.instructionsCB[0xC3] = func() (int, int, string, error) {
		c.insSetNr(0, c.e)
		return 2, 8, "SET 0, E", nil
	}

	// SET 0, H
	c.instructionsCB[0xC4] = func() (int, int, string, error) {
		c.insSetNr(0, c.h)
		return 2, 8, "SET 0, H", nil
	}

	// SET 0, L
	c.instructionsCB[0xC5] = func() (int, int, string, error) {
		c.insSetNr(0, c.l)
		return 2, 8, "SET 0, L", nil
	}

	// SET 0, (HL)
	c.instructionsCB[0xC6] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "SET 0, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insSetNr(0, &r)

		return 2, 16, "SET 0, (HL)", c.mmu.Write(hl, r.get())
	}

	// SET 0, A
	c.instructionsCB[0xC7] = func() (int, int, string, error) {
		c.insSetNr(0, c.a)
		return 2, 8, "SET 0, A", nil
	}

	// SET 1, B
	c.instructionsCB[0xC8] = func() (int, int, string, error) {
		c.insSetNr(1, c.b)
		return 2, 8, "SET 1, B", nil
	}

	// SET 1, C
	c.instructionsCB[0xC9] = func() (int, int, string, error) {
		c.insSetNr(1, c.c)
		return 2, 8, "SET 1, C", nil
	}

	// SET 1, D
	c.instructionsCB[0xCA] = func() (int, int, string, error) {
		c.insSetNr(1, c.d)
		return 2, 8, "SET 1, D", nil
	}

	// SET 1, E
	c.instructionsCB[0xCB] = func() (int, int, string, error) {
		c.insSetNr(1, c.e)
		return 2, 8, "SET 1, E", nil
	}

	// SET 1, H
	c.instructionsCB[0xCC] = func() (int, int, string, error) {
		c.insSetNr(1, c.h)
		return 2, 8, "SET 1, H", nil
	}

	// SET 1, L
	c.instructionsCB[0xCD] = func() (int, int, string, error) {
		c.insSetNr(1, c.l)
		return 2, 8, "SET 1, L", nil
	}

	// SET 1, (HL)
	c.instructionsCB[0xCE] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "SET 1, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insSetNr(1, &r)

		return 2, 16, "SET 1, (HL)", c.mmu.Write(hl, r.get())
	}

	// SET 1, A
	c.instructionsCB[0xCF] = func() (int, int, string, error) {
		c.insSetNr(1, c.a)
		return 2, 8, "SET 1, A", nil
	}

	// SET 2, B
	c.instructionsCB[0xD0] = func() (int, int, string, error) {
		c.insSetNr(2, c.b)
		return 2, 8, "SET 2, B", nil
	}

	// SET 2, C
	c.instructionsCB[0xD1] = func() (int, int, string, error) {
		c.insSetNr(2, c.c)
		return 2, 8, "SET 2, C", nil
	}

	// SET 2, D
	c.instructionsCB[0xD2] = func() (int, int, string, error) {
		c.insSetNr(2, c.d)
		return 2, 8, "SET 2, D", nil
	}

	// SET 2, E
	c.instructionsCB[0xD3] = func() (int, int, string, error) {
		c.insSetNr(2, c.e)
		return 2, 8, "SET 2, E", nil
	}

	// SET 2, H
	c.instructionsCB[0xD4] = func() (int, int, string, error) {
		c.insSetNr(2, c.h)
		return 2, 8, "SET 2, H", nil
	}

	// SET 2, L
	c.instructionsCB[0xD5] = func() (int, int, string, error) {
		c.insSetNr(2, c.l)
		return 2, 8, "SET 2, L", nil
	}

	// SET 2, (HL)
	c.instructionsCB[0xD6] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "SET 2, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insSetNr(2, &r)

		return 2, 16, "SET 2, (HL)", c.mmu.Write(hl, r.get())
	}

	// SET 2, A
	c.instructionsCB[0xD7] = func() (int, int, string, error) {
		c.insSetNr(2, c.a)
		return 2, 8, "SET 2, A", nil
	}

	// SET 3, B
	c.instructionsCB[0xD8] = func() (int, int, string, error) {
		c.insSetNr(3, c.b)
		return 2, 8, "SET 3, B", nil
	}

	// SET 3, C
	c.instructionsCB[0xD9] = func() (int, int, string, error) {
		c.insSetNr(3, c.c)
		return 2, 8, "SET 3, C", nil
	}

	// SET 3, D
	c.instructionsCB[0xDA] = func() (int, int, string, error) {
		c.insSetNr(3, c.d)
		return 2, 8, "SET 3, D", nil
	}

	// SET 3, E
	c.instructionsCB[0xDB] = func() (int, int, string, error) {
		c.insSetNr(3, c.e)
		return 2, 8, "SET 3, E", nil
	}

	// SET 3, H
	c.instructionsCB[0xDC] = func() (int, int, string, error) {
		c.insSetNr(3, c.h)
		return 2, 8, "SET 3, H", nil
	}

	// SET 3, L
	c.instructionsCB[0xDD] = func() (int, int, string, error) {
		c.insSetNr(3, c.l)
		return 2, 8, "SET 3, L", nil
	}

	// SET 3, (HL)
	c.instructionsCB[0xDE] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "SET 3, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insSetNr(3, &r)

		return 2, 16, "SET 3, (HL)", c.mmu.Write(hl, r.get())
	}

	// SET 3, A
	c.instructionsCB[0xDF] = func() (int, int, string, error) {
		c.insSetNr(3, c.a)
		return 2, 8, "SET 3, A", nil
	}

	// SET 4, B
	c.instructionsCB[0xE0] = func() (int, int, string, error) {
		c.insSetNr(4, c.b)
		return 2, 8, "SET 4, B", nil
	}

	// SET 4, C
	c.instructionsCB[0xE1] = func() (int, int, string, error) {
		c.insSetNr(4, c.c)
		return 2, 8, "SET 4, C", nil
	}

	// SET 4, D
	c.instructionsCB[0xE2] = func() (int, int, string, error) {
		c.insSetNr(4, c.d)
		return 2, 8, "SET 4, D", nil
	}

	// SET 4, E
	c.instructionsCB[0xE3] = func() (int, int, string, error) {
		c.insSetNr(4, c.e)
		return 2, 8, "SET 4, E", nil
	}

	// SET 4, H
	c.instructionsCB[0xE4] = func() (int, int, string, error) {
		c.insSetNr(4, c.h)
		return 2, 8, "SET 4, H", nil
	}

	// SET 4, L
	c.instructionsCB[0xE5] = func() (int, int, string, error) {
		c.insSetNr(4, c.l)
		return 2, 8, "SET 4, L", nil
	}

	// SET 4, (HL)
	c.instructionsCB[0xE6] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "SET 4, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insSetNr(4, &r)

		return 2, 16, "SET 4, (HL)", c.mmu.Write(hl, r.get())
	}

	// SET 4, A
	c.instructionsCB[0xE7] = func() (int, int, string, error) {
		c.insSetNr(4, c.a)
		return 2, 8, "SET 4, A", nil
	}

	// SET 5, B
	c.instructionsCB[0xE8] = func() (int, int, string, error) {
		c.insSetNr(5, c.b)
		return 2, 8, "SET 5, B", nil
	}

	// SET 5, C
	c.instructionsCB[0xE9] = func() (int, int, string, error) {
		c.insSetNr(5, c.c)
		return 2, 8, "SET 5, C", nil
	}

	// SET 5, D
	c.instructionsCB[0xEA] = func() (int, int, string, error) {
		c.insSetNr(5, c.d)
		return 2, 8, "SET 5, D", nil
	}

	// SET 5, E
	c.instructionsCB[0xEB] = func() (int, int, string, error) {
		c.insSetNr(5, c.e)
		return 2, 8, "SET 5, E", nil
	}

	// SET 5, H
	c.instructionsCB[0xEC] = func() (int, int, string, error) {
		c.insSetNr(5, c.h)
		return 2, 8, "SET 5, H", nil
	}

	// SET 5, L
	c.instructionsCB[0xED] = func() (int, int, string, error) {
		c.insSetNr(5, c.l)
		return 2, 8, "SET 5, L", nil
	}

	// SET 5, (HL)
	c.instructionsCB[0xEE] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "SET 5, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insSetNr(5, &r)

		return 2, 16, "SET 5, (HL)", c.mmu.Write(hl, r.get())
	}

	// SET 5, A
	c.instructionsCB[0xEF] = func() (int, int, string, error) {
		c.insSetNr(5, c.a)
		return 2, 8, "SET 5, A", nil
	}

	// SET 6, B
	c.instructionsCB[0xF0] = func() (int, int, string, error) {
		c.insSetNr(6, c.b)
		return 2, 8, "SET 6, B", nil
	}

	// SET 6, C
	c.instructionsCB[0xF1] = func() (int, int, string, error) {
		c.insSetNr(6, c.c)
		return 2, 8, "SET 6, C", nil
	}

	// SET 6, D
	c.instructionsCB[0xF2] = func() (int, int, string, error) {
		c.insSetNr(6, c.d)
		return 2, 8, "SET 6, D", nil
	}

	// SET 6, E
	c.instructionsCB[0xF3] = func() (int, int, string, error) {
		c.insSetNr(6, c.e)
		return 2, 8, "SET 6, E", nil
	}

	// SET 6, H
	c.instructionsCB[0xF4] = func() (int, int, string, error) {
		c.insSetNr(6, c.h)
		return 2, 8, "SET 6, H", nil
	}

	// SET 6, L
	c.instructionsCB[0xF5] = func() (int, int, string, error) {
		c.insSetNr(6, c.l)
		return 2, 8, "SET 6, L", nil
	}

	// SET 6, (HL)
	c.instructionsCB[0xF6] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "SET 6, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insSetNr(6, &r)

		return 2, 16, "SET 6, (HL)", c.mmu.Write(hl, r.get())
	}

	// SET 6, A
	c.instructionsCB[0xF7] = func() (int, int, string, error) {
		c.insSetNr(6, c.a)
		return 2, 8, "SET 6, A", nil
	}

	// SET 7, B
	c.instructionsCB[0xF8] = func() (int, int, string, error) {
		c.insSetNr(7, c.b)
		return 2, 8, "SET 7, B", nil
	}

	// SET 7, C
	c.instructionsCB[0xF9] = func() (int, int, string, error) {
		c.insSetNr(7, c.c)
		return 2, 8, "SET 7, C", nil
	}

	// SET 7, D
	c.instructionsCB[0xFA] = func() (int, int, string, error) {
		c.insSetNr(7, c.d)
		return 2, 8, "SET 7, D", nil
	}

	// SET 7, E
	c.instructionsCB[0xFB] = func() (int, int, string, error) {
		c.insSetNr(7, c.e)
		return 2, 8, "SET 7, E", nil
	}

	// SET 7, H
	c.instructionsCB[0xFC] = func() (int, int, string, error) {
		c.insSetNr(7, c.h)
		return 2, 8, "SET 7, H", nil
	}

	// SET 7, L
	c.instructionsCB[0xFD] = func() (int, int, string, error) {
		c.insSetNr(7, c.l)
		return 2, 8, "SET 7, L", nil
	}

	// SET 7, (HL)
	c.instructionsCB[0xFE] = func() (int, int, string, error) {

		hl := c.hl.get()
		v, err := c.mmu.Read(hl)

		if err != nil {
			return 2, 16, "SET 7, (HL)", err
		}

		r := Register8{}
		r.set(v)
		c.insSetNr(7, &r)

		return 2, 16, "SET 7, (HL)", c.mmu.Write(hl, r.get())
	}

	// SET 7, A
	c.instructionsCB[0xFF] = func() (int, int, string, error) {
		c.insSetNr(7, c.a)
		return 2, 8, "SET 7, A", nil
	}
}
