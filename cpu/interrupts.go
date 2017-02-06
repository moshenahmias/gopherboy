/*H**********************************************************************
* FILENAME :        interrupts.go
*
* PACKAGE :			cpu
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package cpu

import "github.com/moshenahmias/gopherboy/memory"

// AddrIE is the address of the Interrupt Enable register
const AddrIE uint16 = 0xFFFF

// AddrIF is the address of the Interrupt Flags register
const AddrIF uint16 = 0xFF0F

// AddrVerticalBlank is the address of the Vertical Blank ISR
const AddrVerticalBlank uint16 = 0x0040

// AddrLCDStatusTriggers is the address of the LCD Status Triggers ISR
const AddrLCDStatusTriggers uint16 = 0x0048

// AddrTimerOverflow is the address of the Timer Overflow ISR
const AddrTimerOverflow uint16 = 0x0050

// AddrSerialLink is the address of the Serial Link ISR
const AddrSerialLink uint16 = 0x0058

// AddrJoypadPress is the address of the Joypad Press ISR
const AddrJoypadPress uint16 = 0x0060

// VerticalBlankFlag in the IF register
const VerticalBlankFlag byte = 0x01

// LCDStatusTriggersFlag in the IF register
const LCDStatusTriggersFlag byte = 0x02

// TimerOverflowFlag in the IF register
const TimerOverflowFlag byte = 0x04

// SerialLinkFlag in the IF register
const SerialLinkFlag byte = 0x08

// JoypadPressFlag in the IF register
const JoypadPressFlag byte = 0x10

// RequestInterrupt sets the requested interrupt bit in the IF register
func (c *Core) RequestInterrupt(flag byte) {
	c.ifr = memory.MemReg(byte(c.ifr) | flag)
}

// handleInterrupts handles the highest enabled and requested interrupt
func (c *Core) handleInterrupts() error {

	ifr := byte(c.ifr)

	if ifr != 0 {

		c.halt = false
		c.stop = false
	}

	if !c.ime {
		return nil
	}

	ier := byte(c.ier)
	mask := ier & ifr

	// vertical blank
	if mask&VerticalBlankFlag != 0 {

		c.ifr = memory.MemReg(byte(c.ifr) & (^VerticalBlankFlag))
		return c.jumpToISR(AddrVerticalBlank)
	}

	// lcd status triggers
	if mask&LCDStatusTriggersFlag != 0 {

		c.ifr = memory.MemReg(byte(c.ifr) & (^LCDStatusTriggersFlag))
		return c.jumpToISR(AddrLCDStatusTriggers)
	}

	// timer overflow
	if mask&TimerOverflowFlag != 0 {

		c.ifr = memory.MemReg(byte(c.ifr) & (^TimerOverflowFlag))
		return c.jumpToISR(AddrTimerOverflow)
	}

	// serial link
	if mask&SerialLinkFlag != 0 {

		c.ifr = memory.MemReg(byte(c.ifr) & (^SerialLinkFlag))
		return c.jumpToISR(AddrSerialLink)
	}

	// joypad press
	if mask&JoypadPressFlag != 0 {

		c.ifr = memory.MemReg(byte(c.ifr) & (^JoypadPressFlag))
		return c.jumpToISR(AddrJoypadPress)
	}

	return nil
}

// jumpToISR disables the IME, saves the current PC and
// jumps to the given ISR address
func (c *Core) jumpToISR(addr uint16) error {

	c.ime = false

	c.sp.decrement()

	if err := c.mmu.Write(c.sp.get(), c.pc.highByte()); err != nil {
		return err
	}

	c.sp.decrement()

	if err := c.mmu.Write(c.sp.get(), c.pc.lowByte()); err != nil {
		return err
	}

	c.pc.set(addr)

	return nil
}
