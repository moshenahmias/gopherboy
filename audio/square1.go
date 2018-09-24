/*H**********************************************************************
* FILENAME :        square1.go
*
* PACKAGE :			audio
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    10 Feb 2017
*
*H*/

package audio

import "github.com/moshenahmias/gopherboy/memory"

// AddrNR10 is the NR10 register address
const AddrNR10 uint16 = 0xFF10

// AddrNR11 is the NR11 register address
const AddrNR11 uint16 = 0xFF11

// AddrNR12 is the NR12 register address
const AddrNR12 uint16 = 0xFF12

// AddrNR13 is the NR13 register address
const AddrNR13 uint16 = 0xFF13

// AddrNR14 is the NR14 register address
const AddrNR14 uint16 = 0xFF14

// Square1 sound channel
type Square1 struct {
	nr10 byte
	nr11 byte
	nr12 byte
	nr13 byte
	nr14 byte

	frequencyCounter int
	envalopeCounter  byte
	sweepCounter     byte

	volume byte

	wavePos   byte
	waveState bool

	freqShadow uint16
}

// Read from the channel registers
func (s *Square1) Read(addr uint16) (byte, error) {

	switch addr {
	case AddrNR10:
		return s.nr10, nil
	case AddrNR11:
		return s.nr11, nil
	case AddrNR12:
		return s.nr12, nil
	case AddrNR13:
		return s.nr13, nil
	case AddrNR14:
		return s.nr14, nil
	}

	return 0, memory.ReadOutOfRangeError(addr)
}

// Write to the timer registers
func (s *Square1) Write(addr uint16, data byte) error {

	switch addr {

	case AddrNR10:

		s.nr10 = data
		return nil

	case AddrNR11:

		s.nr11 = data
		return nil

	case AddrNR12:

		s.nr12 = data
		return nil

	case AddrNR13:

		s.nr13 = data
		return nil

	case AddrNR14:

		s.nr14 = data

		if data&0x80 == 0x80 {
			s.initialize()
		}

		return nil
	}

	return memory.WriteOutOfRangeError(addr)
}

// ClockChanged is called after every instruction execution
func (s *Square1) ClockChanged(cycles int) error {

	s.frequencyCounter += cycles

	waveRate := s.waveRate()
	if s.frequencyCounter >= waveRate {

		s.frequencyCounter = s.frequencyCounter - waveRate

		waveform := DutyTable[s.duty()]

		s.waveState = (waveform>>(7-s.wavePos))&0x01 == 0x01

		s.wavePos = (s.wavePos + 1) % 8
	}

	return nil
}

func (s *Square1) output() byte {

	if !s.enabled() || !s.waveState {
		return 0
	}

	return s.volume
}

func (s *Square1) initialize() {
	s.setLengthLoad(63)
	s.frequencyCounter = 0
	s.envalopeCounter = s.envalopePeriod()
	s.sweepCounter = s.sweepPeriod()
	s.volume = s.envalopeVolume()
	s.wavePos = 0
	s.waveState = false
	s.freqShadow = s.frequency()
	s.doSweepShift()
}

func (s *Square1) enabled() bool {
	return s.lengthLoad() > 0 && s.freqShadow < 2048 && s.dac()
}

func (s *Square1) sweepClock() {

	if s.sweepCounter == 0 {
		return
	}

	s.sweepCounter--

	if s.sweepCounter > 0 {
		return
	}

	s.sweepCounter = s.sweepPeriod()

	s.doSweepShift()
}

func (s *Square1) doSweepShift() {

	shift := s.sweepShift()

	if shift == 0 {
		return
	}

	freq := s.freqShadow
	shiftedFreq := freq >> shift

	if s.sweepNegate() {
		s.freqShadow = freq - shiftedFreq
	} else {
		s.freqShadow = freq + shiftedFreq
	}

	if s.freqShadow < 2048 {
		s.setFrequency(s.freqShadow)
	}
}

func (s *Square1) envelopeClock() {

	if s.envalopeCounter == 0 {
		return
	}

	s.envalopeCounter--

	if s.envalopeMode() {

		// increase
		if s.volume < 15 {
			s.volume++
		}

	} else {

		// decrease
		if s.volume > 0 {
			s.volume--
		}
	}
}

func (s *Square1) lengthClock() {

	if !s.lengthEnabled() {
		return
	}

	length := s.lengthLoad()

	if length == 0 {
		return
	}

	s.setLengthLoad(length - 1)
}

func (s *Square1) waveRate() int {

	freq := int(uint(s.frequency()))

	return (2048 - freq) * 4
}

func (s *Square1) sweepPeriod() byte {
	return (s.nr10 << 1) >> 5
}

func (s *Square1) sweepNegate() bool {
	return s.nr10&0x08 == 0x08
}

func (s *Square1) sweepShift() byte {
	return s.nr10 & 0x07
}

func (s *Square1) duty() byte {
	return s.nr11 >> 6
}

func (s *Square1) lengthLoad() byte {
	return s.nr11 & 0x3F
}

func (s *Square1) setLengthLoad(length byte) {
	s.nr11 = (s.nr11 & 0xC0) | (length & 0x3F)
}

func (s *Square1) lengthEnabled() bool {
	return s.nr14&0x40 == 0x40
}

func (s *Square1) envalopeVolume() byte {
	return s.nr12 >> 4
}

func (s *Square1) envalopeMode() bool {
	return s.nr12&0x08 == 0x08
}

func (s *Square1) envalopePeriod() byte {
	return s.nr12 & 0x07
}

func (s *Square1) frequency() uint16 {
	return uint16(s.nr13) | (uint16(s.nr14&0x07) << 8)
}

func (s *Square1) setFrequency(frequency uint16) {
	s.nr13 = byte(frequency & 0x000F)
	s.nr14 = (s.nr14 & 0xF8) | byte((frequency>>8)&0x0007)
}

func (s *Square1) dac() bool {
	return s.nr12&0xF8 != 0
}
