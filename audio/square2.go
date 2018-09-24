/*H**********************************************************************
* FILENAME :        square2.go
*
* PACKAGE :			audio
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    10 Feb 2017
*
*H*/

package audio

import "github.com/moshenahmias/gopherboy/memory"

// AddrNR21 is the NR21 register address
const AddrNR21 uint16 = 0xFF16

// AddrNR22 is the NR22 register address
const AddrNR22 uint16 = 0xFF17

// AddrNR23 is the NR23 register address
const AddrNR23 uint16 = 0xFF18

// AddrNR24 is the NR24 register address
const AddrNR24 uint16 = 0xFF19

// Square2 sound channel
type Square2 struct {
	nr21 byte
	nr22 byte
	nr23 byte
	nr24 byte

	frequencyCounter int
	envalopeCounter  byte

	volume byte

	wavePos   byte
	waveState bool
}

// Read from the channel registers
func (s *Square2) Read(addr uint16) (byte, error) {

	switch addr {
	case AddrNR21:
		return s.nr21, nil
	case AddrNR22:
		return s.nr22, nil
	case AddrNR23:
		return s.nr23, nil
	case AddrNR24:
		return s.nr24, nil
	}

	return 0, memory.ReadOutOfRangeError(addr)
}

// Write to the timer registers
func (s *Square2) Write(addr uint16, data byte) error {

	switch addr {

	case AddrNR21:

		s.nr21 = data
		return nil

	case AddrNR22:

		s.nr22 = data
		return nil

	case AddrNR23:

		s.nr23 = data
		return nil

	case AddrNR24:

		s.nr24 = data

		if data&0x80 == 0x80 {
			s.initialize()
		}

		return nil
	}

	return memory.WriteOutOfRangeError(addr)
}

// ClockChanged is called after every instruction execution
func (s *Square2) ClockChanged(cycles int) error {

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

func (s *Square2) output() byte {

	if !s.enabled() || !s.waveState {
		return 0
	}

	return s.volume
}

func (s *Square2) initialize() {
	s.setLengthLoad(63)
	s.frequencyCounter = 0
	s.envalopeCounter = s.envalopePeriod()
	s.volume = s.envalopeVolume()
	s.wavePos = 0
	s.waveState = false
}

func (s *Square2) enabled() bool {
	return s.lengthLoad() > 0 && s.dac()
}

func (s *Square2) envelopeClock() {

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

func (s *Square2) lengthClock() {

	if !s.lengthEnabled() {
		return
	}

	length := s.lengthLoad()

	if length == 0 {
		return
	}

	s.setLengthLoad(length - 1)
}

func (s *Square2) waveRate() int {

	freq := int(uint(s.frequency()))

	return (2048 - freq) * 4
}

func (s *Square2) duty() byte {
	return s.nr21 >> 6
}

func (s *Square2) lengthLoad() byte {
	return s.nr21 & 0x3F
}

func (s *Square2) setLengthLoad(length byte) {
	s.nr21 = (s.nr21 & 0xC0) | (length & 0x3F)
}

func (s *Square2) lengthEnabled() bool {
	return s.nr24&0x40 == 0x40
}

func (s *Square2) envalopeVolume() byte {
	return s.nr22 >> 4
}

func (s *Square2) envalopeMode() bool {
	return s.nr22&0x08 == 0x08
}

func (s *Square2) envalopePeriod() byte {
	return s.nr22 & 0x07
}

func (s *Square2) frequency() uint16 {
	return uint16(s.nr23) | (uint16(s.nr24&0x07) << 8)
}

func (s *Square2) dac() bool {
	return s.nr22&0xF8 != 0
}
