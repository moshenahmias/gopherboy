package cpu

import (
	"github.com/moshenahmias/gopherboy/memory"

	"fmt"
	"time"
)

// Frequency of the cpu cycles per seconds
const Frequency int = 4194304

func noSuchInstructionError(opcode byte) error {
	return fmt.Errorf("No such instruction %02x", opcode)
}

// TimedUnit is a cpu cycles observer
type TimedUnit interface {
	ClockChanged(cycles int) error
}

type Instruction func() (int, int, string, error)

// Core represents the CPU's core, it includes
// its registers, access to the MMU and instructions.
type Core struct {
	sp Register16 // stack pointer
	pc Register16 // program counter
	af Register16 // AF register
	bc Register16 // BC register
	de Register16 // DE register
	hl Register16 // HL register

	a *Register8 // reference to register A
	f *Register8 // reference to register F
	b *Register8 // reference to register B
	c *Register8 // reference to register C
	d *Register8 // reference to register D
	e *Register8 // reference to register E
	h *Register8 // reference to register H
	l *Register8 // reference to register L

	instructions   [256]Instruction // instruction set
	instructionsCB [256]Instruction // cb instruction set

	mmu        memory.Unit   // MMU
	ime        bool          // interrupt master enable
	ier        memory.MemReg // interrupt enable register
	ifr        memory.MemReg // interrupt flags register
	quit       bool          // stop flag
	halt       bool          // halt flag
	stop       bool          // bool flag
	timedUnits []TimedUnit   // clocked units

	throttle int

	pause bool
}

// NewCore creates Core instance
func NewCore(mmu *memory.MMU) (*Core, error) {

	c := Core{mmu: mmu}

	c.throttle = 5500

	c.a = c.af.high()
	c.f = c.af.low()
	c.b = c.bc.high()
	c.c = c.bc.low()
	c.d = c.de.high()
	c.e = c.de.low()
	c.h = c.hl.high()
	c.l = c.hl.low()

	if err := mmu.Map(&c.ifr, 0xFF0F, 0xFF0F); err != nil {
		return nil, err
	}

	if err := mmu.Map(&c.ier, 0xFFFF, 0xFFFF); err != nil {
		return nil, err
	}

	c.initInstructions()
	c.initInstructionsCB()

	return &c, nil
}

// RegisterToClockChanges that take place after every instruction execution
func (c *Core) RegisterToClockChanges(unit TimedUnit) {
	c.timedUnits = append(c.timedUnits, unit)
}

// Throttle the cpu speed
func (c *Core) Throttle(tooFast bool) {

	if tooFast {
		c.throttle += 10
	} else {
		c.throttle -= 10
	}
}

// Start the cpu activity at address 'pc'
func (c *Core) Start(pc uint16) error {

	c.pc.set(pc)

	for !c.quit {

		for c.pause && !c.quit {
			time.Sleep(time.Millisecond * 100)
		}

		opcode, err := c.mmu.Read(c.pc.get())

		if err != nil {
			return c.wrapError(err, "pc read failed")
		}

		ins := c.instructions[opcode]

		if ins == nil {
			return c.wrapError(noSuchInstructionError(opcode), "instruction fetch failed")
		}

		_, cycles, name, err := ins()

		if err != nil {
			return c.wrapErrorf(err, "%s %02x failed", name, opcode)
		}

		for i := 0; i < c.throttle; i++ {
			// do absolutely nothing
		}

		c.pc.increment()

		for do := true; do; do = (c.halt || c.stop) && !c.quit {

			for c.pause && !c.quit {
				time.Sleep(time.Millisecond * 100)
			}

			if c.halt || c.stop {
				if _, err := c.mmu.Read(0xFF00); err != nil {
					return c.wrapError(err, "joyp read (during stop) failed")
				}
			}

			for _, u := range c.timedUnits {
				if err := u.ClockChanged(cycles); err != nil {
					return c.wrapError(err, "ClockChanged() failed")
				}
			}

			if err := c.handleInterrupts(); err != nil {
				return c.wrapError(err, "HandleInterrupts() failed")
			}
		}
	}

	return nil
}

// Pause the cpu (or resume if already paused)
func (c *Core) Pause() {
	c.pause = !c.pause
}

// Stop the execition loop
func (c *Core) Stop() {
	c.quit = true
}

// wrapErrorf is the formatted version of wrapError
func (c *Core) wrapErrorf(err error, format string, v ...interface{}) error {
	return c.wrapError(err, fmt.Sprintf(format, v...))
}

// wrapError with a custom message and cpu info
func (c *Core) wrapError(err error, message string) error {

	return fmt.Errorf(
		"%s - %s ([A: %02x] [BC: %04x] [DE: %04x] [HL: %04x] [SP: %04x] [PC: %04x] [ZNHC: %04b] [IME: %t] [IE: %02x] [IF: %02x])",
		message,
		err,
		c.a.get(),
		c.bc.get(),
		c.de.get(),
		c.hl.get(),
		c.sp.get(),
		c.pc.get(),
		c.f.get()>>4,
		c.ime,
		c.ier,
		c.ifr)
}

// loadImmediate8 bit immediate value
func (c *Core) loadImmediate8() (byte, error) {

	c.pc.increment()
	val, err := c.mmu.Read(c.pc.get())

	if err != nil {
		return 0, err
	}

	return val, nil
}

// loadImmediate16 bit immediate value
func (c *Core) loadImmediate16() (uint16, error) {

	c.pc.increment()
	lval, err := c.mmu.Read(c.pc.get())

	if err != nil {
		return 0, err
	}

	c.pc.increment()
	hval, err := c.mmu.Read(c.pc.get())

	if err != nil {
		return 0, err
	}

	val := uint16(hval)
	val = (val << 8) | uint16(lval)

	return val, nil
}
