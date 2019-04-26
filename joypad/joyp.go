package joypad

import "github.com/moshenahmias/gopherboy/cpu"

const columnDirections byte = 0x20
const columnStartSelectAB byte = 0x10

// ButtonRight code
const ButtonRight byte = 0x01

// ButtonA code
const ButtonA byte = 0x11

// ButtonLeft code
const ButtonLeft byte = 0x02

// ButtonB code
const ButtonB byte = 0x12

// ButtonUp code
const ButtonUp byte = 0x04

// ButtonSelect code
const ButtonSelect byte = 0x14

// ButtonDown code
const ButtonDown byte = 0x08

// ButtonStart code
const ButtonStart byte = 0x18

// Keystroke represents a key press / release
type Keystroke struct {
	Button  byte
	Pressed bool
}

// Keystroker input interface
type Keystroker interface {
	GetKeystroke() *Keystroke
}

// AddrJOYP is the address for the joypad register
const AddrJOYP uint16 = 0xFF00

// JOYP register
type JOYP struct {
	data       byte
	core       *cpu.Core
	keystroker Keystroker
	state      [2]byte
}

// NewJOYP creates JOYP instance
func NewJOYP(core *cpu.Core, keystroker Keystroker) *JOYP {

	j := JOYP{core: core, keystroker: keystroker}

	j.state[0] = 0x0F
	j.state[1] = 0x0F
	j.data = 0xFF

	return &j
}

// Read joyp
func (j *JOYP) Read(addr uint16) (byte, error) {

	j.setWireState()

	return j.data, nil
}

// Write to joyp
func (j *JOYP) Write(addr uint16, data byte) error {

	j.data = (j.data & 0x0F) | (data & 0xF0)

	return nil
}

// setWireState according to the buttons state
func (j *JOYP) setWireState() error {

	if ks := j.keystroker.GetKeystroke(); ks != nil {
		btn := ks.Button

		if ks.Pressed {
			j.state[btn>>4] &= ^(btn & 0x0F)
		} else {
			j.state[btn>>4] |= (btn & 0x0F)
		}
	}

	prev := j.data & 0x0F

	if j.data&0x30 == columnStartSelectAB {

		// buttons A, B, Select or Start
		j.data = (j.data & 0xF0) | j.state[1]

	} else if j.data&0x30 == columnDirections {

		// direction buttons
		j.data = (j.data & 0xF0) | j.state[0]

	} else {

		// nothing
		j.data |= 0x0F
	}

	if (prev^(j.data&0x0F))&prev != 0x00 {

		// request joypad press interrupt
		j.core.RequestInterrupt(cpu.JoypadPressFlag)
	}

	return nil
}
