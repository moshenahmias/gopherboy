/*H**********************************************************************
* FILENAME :        framesequencer.go
*
* PACKAGE :			audio
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    10 Fab 2017
*
*H*/

package audio

// LengthRate is the length period in cpu cycles
const LengthRate int = 16384

// EnvelopeRate is the envelope period in cpu cycles
const EnvelopeRate int = 65536

// SweepeRate is the sweep period in cpu cycles
const SweepeRate int = 32768

// FrameSequencer generates low frequency clocks for the
// modulation units. It is clocked by a 512 Hz timer.
type FrameSequencer struct {
	sweepCounter    int
	envelopeCounter int
	lengthCounter   int
	ch1             *Square1
	ch2             *Square2
	ch3             *Wave
	ch4             *Noise
}

// ClockChanged is called after every instruction execution
func (f *FrameSequencer) ClockChanged(cycles int) error {

	f.sweepCounter += cycles
	f.envelopeCounter += cycles
	f.lengthCounter += cycles

	// 128hz
	if f.sweepCounter >= SweepeRate {

		f.sweepCounter = f.sweepCounter - SweepeRate
		f.ch1.sweepClock()
	}

	// 256hz
	if f.lengthCounter >= LengthRate {

		f.lengthCounter = f.lengthCounter - LengthRate

		f.ch1.lengthClock()
		f.ch2.lengthClock()
		f.ch3.lengthClock()
		f.ch4.lengthClock()
	}

	// 64hz
	if f.envelopeCounter >= EnvelopeRate {

		f.envelopeCounter = f.envelopeCounter - EnvelopeRate

		f.ch1.envelopeClock()
		f.ch2.envelopeClock()
		f.ch4.envelopeClock()
	}

	return nil
}
