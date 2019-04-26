package audio

import (
	"github.com/moshenahmias/gopherboy/cpu"
	"github.com/moshenahmias/gopherboy/memory"
)

// DutyTable of square wave 1 and 2
var DutyTable = [4]byte{
	0x01,
	0x81,
	0x87,
	0x7E}

// Audioer outputs sound via a samples buffer
type Audioer interface {
	Queue(samples []byte) error
	Frequency() int
	BufferSize() uint16
	SamplesCount() uint32
}

// APU emulates the device sound chip
type APU struct {
	fs             *FrameSequencer
	ch1            Square1
	ch2            Square2
	ch3            Wave
	ch4            Noise
	control        Control
	audioer        Audioer
	samplesCounter int
	throttle       int
}

// NewAPU creates APU instance
func NewAPU(core *cpu.Core, mmu *memory.MMU, audioer Audioer) (*APU, error) {

	a := APU{audioer: audioer}

	// init all channels
	a.ch1.initialize()
	a.ch2.initialize()
	a.ch3.initialize()
	a.ch4.initialize()

	fs := FrameSequencer{ch1: &a.ch1, ch2: &a.ch2, ch3: &a.ch3, ch4: &a.ch4}

	a.fs = &fs

	// map sound channel 1
	if err := mmu.Map(&a.ch1, AddrNR10, AddrNR14); err != nil {
		return nil, err
	}

	// map sound channel 2
	if err := mmu.Map(&a.ch2, AddrNR21, AddrNR24); err != nil {
		return nil, err
	}

	// map sound channel 3
	if err := mmu.Map(&a.ch3, AddrNR30, AddrNR34); err != nil {
		return nil, err
	}

	// map sound control
	if err := mmu.Map(&a.control, AddrNR50, AddrNR52); err != nil {
		return nil, err
	}

	if err := mmu.Map(&a.ch3, AddrWaveTableStart, AddrWaveTableEnd); err != nil {
		return nil, err
	}

	// map sound channel 4
	if err := mmu.Map(&a.ch4, AddrNR41, AddrNR44); err != nil {
		return nil, err
	}

	// register to the cpu's clock
	core.RegisterToClockChanges(&a)

	return &a, nil
}

// ClockChanged is called after every instruction execution
func (a *APU) ClockChanged(cycles int) error {

	if err := a.ch1.ClockChanged(cycles); err != nil {
		return err
	}

	if err := a.ch2.ClockChanged(cycles); err != nil {
		return err
	}

	if err := a.ch3.ClockChanged(cycles); err != nil {
		return err
	}

	if err := a.ch4.ClockChanged(cycles); err != nil {
		return err
	}

	if err := a.fs.ClockChanged(cycles); err != nil {
		return err
	}

	a.samplesCounter += cycles

	sc := a.audioer.SamplesCount()
	bs := a.audioer.BufferSize()

	if sc < uint32((3*bs)/4) && a.throttle > -10 {
		a.throttle--
	} else if sc == uint32(bs) && a.throttle < -5 {
		a.throttle++
	}

	freq := (cpu.Frequency / a.audioer.Frequency()) + a.throttle

	for a.samplesCounter >= freq {

		a.samplesCounter = a.samplesCounter - freq

		var sampleLeft byte
		var sampleRight byte

		if a.control.soundOn() {

			leftChannels := a.control.leftChannels()
			rightChannels := a.control.rightChannels()

			var sampleLeftCh1 byte
			var sampleLeftCh2 byte
			var sampleLeftCh3 byte
			var sampleLeftCh4 byte
			var sampleRightCh1 byte
			var sampleRightCh2 byte
			var sampleRightCh3 byte
			var sampleRightCh4 byte

			if leftChannels&channel1 != 0 {
				sampleLeftCh1 = a.ch1.output()
			}

			if leftChannels&channel2 != 0 {
				sampleLeftCh2 = a.ch2.output()
			}

			if leftChannels&channel3 != 0 {
				sampleLeftCh3 = a.ch3.output()
			}

			if leftChannels&channel4 != 0 {
				sampleLeftCh4 = a.ch4.output()
			}

			if rightChannels&channel1 != 0 {
				sampleRightCh1 = a.ch1.output()
			}

			if rightChannels&channel2 != 0 {
				sampleRightCh2 = a.ch2.output()
			}

			if rightChannels&channel3 != 0 {
				sampleRightCh3 = a.ch3.output()
			}

			if rightChannels&channel4 != 0 {
				sampleRightCh4 = a.ch4.output()
			}

			sampleLeft = sampleLeftCh1 + sampleLeftCh2 + sampleLeftCh3 + sampleLeftCh4
			sampleRight = sampleRightCh1 + sampleRightCh2 + sampleRightCh3 + sampleRightCh4

			sampleLeft += a.control.leftVolume()
			sampleRight += a.control.rightVolume()
		}

		if err := a.audioer.Queue([]byte{sampleLeft, sampleRight}); err != nil {
			return err
		}
	}

	return nil
}
