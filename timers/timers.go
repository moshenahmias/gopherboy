/*H**********************************************************************
* FILENAME :        timers.go
*
* PACKAGE :			timers
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package timers

import (
	"github.com/moshenahmias/gopherboy/cpu"
	"github.com/moshenahmias/gopherboy/memory"
)

// AddrDIV is the DIV register address
const AddrDIV uint16 = 0xFF04

// AddrTIMA is the TIMA register address
const AddrTIMA uint16 = 0xFF05

// AddrTMA is the TMA register address
const AddrTMA uint16 = 0xFF06

// AddrTAC is the TAC register address
const AddrTAC uint16 = 0xFF07

// DividerRate in cycles
const DividerRate int = 256

// Timer emulates the gameboy's timer
type Timer struct {
	div         byte
	tima        byte
	tma         byte
	tac         byte
	divCounter  int
	timaCounter int
	core        *cpu.Core
}

// NewTimer creates GPU instance
func NewTimer(core *cpu.Core) *Timer {

	t := Timer{core: core}
	core.RegisterToClockChanges(&t)
	return &t
}

// rateOfTIMA returns the number of cycles that
// need to pass for the TIMA timer to increment
func (t *Timer) rateOfTIMA() int {

	switch t.tac & 0x03 {

	case 0x00:
		return 1024 // 4096 Hz
	case 0x01:
		return 16 // 262144 Hz
	case 0x02:
		return 64 //65536 Hz
	case 0x03:
		return 256 //16384 Hz
	}

	return 0
}

// timaEnabled returns true iff the TIMA timer is enabled
func (t *Timer) timaEnabled() bool {
	return t.tac&0x04 == 0x04
}

// Read from the timer registers
func (t *Timer) Read(addr uint16) (byte, error) {

	if addr == AddrDIV {
		return t.div, nil
	}

	if addr == AddrTIMA {
		return t.tima, nil
	}

	if addr == AddrTMA {
		return t.tma, nil
	}

	if addr == AddrTAC {
		return t.tac, nil
	}

	return 0, memory.ReadOutOfRangeError(addr)
}

// Write to the timer registers
func (t *Timer) Write(addr uint16, data byte) error {

	if addr == AddrDIV {
		t.div = 0
		return nil
	}

	if addr == AddrTIMA {
		t.tima = data
		return nil
	}

	if addr == AddrTMA {
		t.tma = data
		return nil
	}

	if addr == AddrTAC {

		if t.tac&0x07 != data&0x07 {
			t.timaCounter = 0
		}

		t.tac = data

		return nil
	}

	return memory.WriteOutOfRangeError(addr)
}

// ClockChanged is called after every instruction execution
func (t *Timer) ClockChanged(cycles int) error {

	t.divCounter += cycles

	if t.divCounter >= DividerRate {
		t.divCounter = t.divCounter - DividerRate
		t.div++
	}

	if t.timaEnabled() {

		t.timaCounter += cycles

		timaRate := t.rateOfTIMA()

		if t.timaCounter >= timaRate {

			t.timaCounter = t.timaCounter - timaRate
			t.tima++

			if t.tima == 0 {

				t.tima = t.tma

				// request timer interrupt
				t.core.RequestInterrupt(cpu.TimerOverflowFlag)
			}
		}
	}

	return nil
}
