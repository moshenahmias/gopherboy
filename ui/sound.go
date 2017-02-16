/*H**********************************************************************
* FILENAME :        sound.go
*
* PACKAGE :			ui
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    16 Fab 2017
*
*H*/

package ui

import "github.com/veandco/go-sdl2/sdl"

// Sound plays samples from the apu
type Sound struct {
	dev         sdl.AudioDeviceID
	mute        bool
	initialized bool
}

// Queue samples to buffer
func (s *Sound) Queue(samples []byte) error {

	if s.mute || !s.initialized {
		return nil
	}

	if err := sdl.QueueAudio(s.dev, samples); err != nil {
		return err
	}

	return nil
}

// Mute or unmute the sound
func (s *Sound) Mute(mute bool) {

	s.mute = mute

	if mute {

		sdl.ClearQueuedAudio(s.dev)
		sdl.PauseAudioDevice(s.dev, true)

	} else {

		sdl.PauseAudioDevice(s.dev, false)
	}
}

// Frequency of sound (samples/sec)
func (s *Sound) Frequency() int {
	return 44100
}

// SamplesCount return the current number of samples in buffer
func (s *Sound) SamplesCount() uint32 {
	return sdl.GetQueuedAudioSize(s.dev)
}

// BufferSize is the size of the samples
// buffer
func (s *Sound) BufferSize() uint16 {
	return 512
}

// Close the sound
func (s *Sound) Close() error {

	sdl.CloseAudioDevice(s.dev)

	return nil
}

// Initialize the sound
func (s *Sound) Initialize(dev int) error {

	if dev < 0 {
		return nil
	}

	device := sdl.GetAudioDeviceName(dev, false)

	spec := sdl.AudioSpec{
		Freq:     int32(s.Frequency()),
		Format:   sdl.AUDIO_U8,
		Channels: 2,
		Samples:  s.BufferSize(),
		Callback: nil}

	devID, err := sdl.OpenAudioDevice(device, false, &spec, nil, 0)

	if err != nil {
		return err
	}

	s.dev = devID

	sdl.PauseAudioDevice(s.dev, false)

	s.initialized = true

	return nil
}
