/*H**********************************************************************
* FILENAME :        noise.go
*
* PACKAGE :			audio
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    10 Feb 2017
*
*H*/

package audio

import (
	"github.com/moshenahmias/gopherboy/cpu"
	"github.com/moshenahmias/gopherboy/memory"
)

// AddrNR41 is the NR41 register address
const AddrNR41 uint16 = 0xFF20

// AddrNR42 is the NR42 register address
const AddrNR42 uint16 = 0xFF21

// AddrNR43 is the NR43 register address
const AddrNR43 uint16 = 0xFF22

// AddrNR44 is the NR44 register address
const AddrNR44 uint16 = 0xFF23

// Noise sound channel
type Noise struct {
	nr41 byte
	nr42 byte
	nr43 byte
	nr44 byte

	frequencyCounter int
	envalopeCounter  byte

	volume byte

	waveState bool

	lfsr uint16
}

// Read from the channel registers
func (n *Noise) Read(addr uint16) (byte, error) {

	switch addr {
	case AddrNR41:
		return n.nr41, nil
	case AddrNR42:
		return n.nr42, nil
	case AddrNR43:
		return n.nr43, nil
	case AddrNR44:
		return n.nr44, nil
	}

	return 0, memory.ReadOutOfRangeError(addr)
}

// Write to the timer registers
func (n *Noise) Write(addr uint16, data byte) error {

	switch addr {

	case AddrNR41:

		n.nr41 = data
		return nil

	case AddrNR42:

		n.nr42 = data
		return nil

	case AddrNR43:

		n.nr43 = data
		return nil

	case AddrNR44:

		n.nr44 = data

		if data&0x80 == 0x80 {
			n.trigger()
		}

		return nil
	}

	return memory.WriteOutOfRangeError(addr)
}

// ClockChanged is called after every instruction execution
func (n *Noise) ClockChanged(cycles int) error {

	n.frequencyCounter += cycles

	waveRate := n.waveRate()

	if n.frequencyCounter >= waveRate {

		n.frequencyCounter = n.frequencyCounter - waveRate

		xor := (n.lfsr ^ (n.lfsr >> 1)) & 0x0001
		n.lfsr = n.lfsr >> 1

		if n.narrowLSFR() {

			n.lfsr = (n.lfsr & 0x003F) | (xor << 6)

		} else {

			n.lfsr = (n.lfsr & 0x3FFF) | (xor << 14)
		}

		n.waveState = (n.lfsr & 0x0001) != 0x0001
	}

	return nil
}

func (n *Noise) output() byte {

	if !n.enabled() || !n.waveState {
		return 0
	}

	return n.volume
}

func (n *Noise) initialize() {

	n.lfsr = 0xFFFF
	n.trigger()
}

func (n *Noise) trigger() {
	n.setLengthLoad(63)
	n.frequencyCounter = 0
	n.envalopeCounter = n.envalopePeriod()
	n.volume = n.envalopeVolume()
	n.waveState = false
}

func (n *Noise) enabled() bool {
	return n.lengthLoad() > 0 && n.dac()
}

func (n *Noise) envelopeClock() {

	if n.envalopeCounter == 0 {
		return
	}

	n.envalopeCounter--

	if n.envalopeMode() {

		// increase
		if n.volume < 15 {
			n.volume++
		}

	} else {

		// decrease
		if n.volume > 0 {
			n.volume--
		}
	}
}

func (n *Noise) lengthClock() {

	if !n.lengthEnabled() {
		return
	}

	length := n.lengthLoad()

	if length == 0 {
		return
	}

	n.setLengthLoad(length - 1)
}

func (n *Noise) waveRate() int {

	f := cpu.Frequency / 8
	r := n.divisor()

	var q int

	if r == 0 {
		q = f * 2
	} else {
		q = f / int(uint(r))
	}

	return q >> (n.clockShift() + 1)
}

func (n *Noise) lengthLoad() byte {
	return n.nr41 & 0x3F
}

func (n *Noise) setLengthLoad(length byte) {
	n.nr41 = (n.nr41 & 0xC0) | (length & 0x3F)
}

func (n *Noise) lengthEnabled() bool {
	return n.nr44&0x40 == 0x40
}

func (n *Noise) envalopeVolume() byte {
	return n.nr42 >> 4
}

func (n *Noise) envalopeMode() bool {
	return n.nr42&0x08 == 0x08
}

func (n *Noise) envalopePeriod() byte {
	return n.nr42 & 0x07
}

func (n *Noise) divisor() byte {
	return n.nr43 & 0x07
}

func (n *Noise) narrowLSFR() bool {
	return n.nr43&0x08 == 0x08
}

func (n *Noise) clockShift() byte {
	return n.nr43 >> 4
}

func (n *Noise) dac() bool {
	return n.nr42&0xF8 != 0
}
