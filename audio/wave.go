/*H**********************************************************************
* FILENAME :        wave.go
*
* PACKAGE :			audio
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    11 Feb 2017
*
*H*/

package audio

import "github.com/moshenahmias/gopherboy/memory"

// AddrNR30 is the NR30 register address
const AddrNR30 uint16 = 0xFF1A

// AddrNR31 is the NR31 register address
const AddrNR31 uint16 = 0xFF1B

// AddrNR32 is the NR32 register address
const AddrNR32 uint16 = 0xFF1C

// AddrNR33 is the NR33 register address
const AddrNR33 uint16 = 0xFF1D

// AddrNR34 is the NR34 register address
const AddrNR34 uint16 = 0xFF1E

// AddrWaveTableStart is the starting address of the wave table
const AddrWaveTableStart uint16 = 0xFF30

// AddrWaveTableEnd is the ending address of the wave table
const AddrWaveTableEnd uint16 = 0xFF3F

// Wave sound channel
type Wave struct {
	nr30 byte
	nr31 byte
	nr32 byte
	nr33 byte
	nr34 byte

	waveTable [16]byte

	frequencyCounter int

	wavePos byte
}

// Read from the channel registers
func (w *Wave) Read(addr uint16) (byte, error) {

	switch addr {
	case AddrNR30:
		return w.nr30, nil
	case AddrNR31:
		return w.nr31, nil
	case AddrNR32:
		return w.nr32, nil
	case AddrNR33:
		return w.nr33, nil
	case AddrNR24:
		return w.nr34, nil
	}

	if AddrWaveTableStart <= addr && addr <= AddrWaveTableEnd {
		return w.waveTable[addr-AddrWaveTableStart], nil
	}

	return 0, memory.ReadOutOfRangeError(addr)
}

// Write to the timer registers
func (w *Wave) Write(addr uint16, data byte) error {

	switch addr {

	case AddrNR30:

		w.nr30 = data
		return nil

	case AddrNR31:

		w.nr31 = data
		return nil

	case AddrNR32:

		w.nr32 = data

		return nil

	case AddrNR33:

		w.nr33 = data
		return nil

	case AddrNR34:

		w.nr34 = data

		if data&0x80 == 0x80 {
			w.initialize()
		}

		return nil
	}

	if AddrWaveTableStart <= addr && addr <= AddrWaveTableEnd {
		w.waveTable[addr-AddrWaveTableStart] = data
		return nil
	}

	return memory.WriteOutOfRangeError(addr)
}

// ClockChanged is called after every instruction execution
func (w *Wave) ClockChanged(cycles int) error {

	w.frequencyCounter += cycles

	waveRate := w.waveRate()

	if w.frequencyCounter >= waveRate {

		w.frequencyCounter = w.frequencyCounter - waveRate
		w.wavePos = (w.wavePos + 1) % 32
	}

	return nil
}

func (w *Wave) output() byte {

	if !w.enabled() {
		return 0
	}

	return w.currentSample() >> w.volumeShifts()
}

func (w *Wave) initialize() {
	w.setLengthLoad(255)
	w.frequencyCounter = 0
	w.wavePos = 0
}

func (w *Wave) enabled() bool {
	return w.lengthLoad() > 0 && w.dac()
}

func (w *Wave) lengthClock() {

	if !w.lengthEnabled() {
		return
	}

	length := w.lengthLoad()

	if length == 0 {
		return
	}

	w.setLengthLoad(length - 1)
}

func (w *Wave) waveRate() int {

	freq := int(uint(w.frequency()))

	return (2048 - freq) * 4
}

func (w *Wave) lengthLoad() byte {
	return w.nr31
}

func (w *Wave) setLengthLoad(length byte) {
	w.nr31 = length
}

func (w *Wave) lengthEnabled() bool {
	return w.nr34&0x40 == 0x40
}

func (w *Wave) frequency() uint16 {
	return uint16(w.nr33) | (uint16(w.nr34&0x07) << 8)
}

func (w *Wave) volumeShifts() byte {

	switch (w.nr32 << 1) >> 6 {
	case 0:
		return 4
	case 1:
		return 0
	case 2:
		return 1
	}

	return 2
}

func (w *Wave) currentSample() byte {

	if w.wavePos%2 == 0 {

		return w.waveTable[w.wavePos/2] & 0x0F
	}

	return w.waveTable[w.wavePos/2] >> 4
}

func (w *Wave) dac() bool {
	return w.nr30&0x80 == 0x80
}
