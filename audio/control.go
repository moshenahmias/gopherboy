package audio

import "github.com/moshenahmias/gopherboy/memory"

const (
	channel1 = 0x01
	channel2 = 0x02
	channel3 = 0x04
	channel4 = 0x08
)

// AddrNR50 is the NR50 register address
const AddrNR50 uint16 = 0xFF24

// AddrNR51 is the NR51 register address
const AddrNR51 uint16 = 0xFF25

// AddrNR52 is the NR52 register address
const AddrNR52 uint16 = 0xFF26

// Control and status manager
type Control struct {
	nr50 byte
	nr51 byte
	nr52 byte
}

// Read from the control unit registers
func (c *Control) Read(addr uint16) (byte, error) {

	switch addr {
	case AddrNR50:
		return c.nr50, nil
	case AddrNR51:
		return c.nr51, nil
	case AddrNR52:
		return c.nr52, nil
	}

	return 0, memory.ReadOutOfRangeError(addr)
}

// Write to the control unit registers
func (c *Control) Write(addr uint16, data byte) error {

	switch addr {

	case AddrNR50:

		c.nr50 = data
		return nil

	case AddrNR51:

		c.nr51 = data
		return nil

	case AddrNR52:

		c.nr52 = (c.nr52 & 0x7F) | (data & 0x80)

		return nil
	}

	return memory.WriteOutOfRangeError(addr)
}

func (c *Control) soundOn() bool {
	return c.nr52&0x80 == 0x80
}

func (c *Control) vinL() bool {
	return c.nr50&0x80 == 0x80
}

func (c *Control) vinR() bool {
	return c.nr50&0x08 == 0x08
}

func (c *Control) leftVolume() byte {
	return (c.nr50 << 1) >> 5
}

func (c *Control) rightVolume() byte {
	return c.nr50 & 0x07
}

func (c *Control) leftChannels() byte {
	return c.nr51 >> 4
}

func (c *Control) rightChannels() byte {
	return c.nr51 & 0x0F
}
