/*H**********************************************************************
* FILENAME :        input.go
*
* PACKAGE :			ui
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package ui

import (
	"errors"

	"github.com/moshenahmias/gopherboy/config"
	"github.com/moshenahmias/gopherboy/joypad"

	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

// ErrInvalidJoypadCode returned for an invalid joypad code
var ErrInvalidJoypadCode = errors.New("Invalid joypad code")

const maxQueueSize int = 20

// ControlEvent is a key event that is not one
// of the 8 gameboy key events
type ControlEvent byte

// ControlEventQuit signals a quit request
const ControlEventQuit ControlEvent = 0

// ControlEventReset signals a reset request
const ControlEventReset ControlEvent = 1

// ControlEventPause signals a pause request
const ControlEventPause ControlEvent = 2

// Input is a Keystroker implementer
type Input struct {
	keystrokes []joypad.Keystroke
	m          sync.Mutex
	mapping    map[int32]config.EJoypad
	stop       bool
}

// NewInput creates Input instance
func NewInput(mapping map[int32]config.EJoypad) *Input {

	return &Input{mapping: mapping, stop: false}
}

// convertJoypadCode to internal system code
func (i *Input) convertJoypadCode(code config.EJoypad) (byte, error) {

	switch code {

	case config.EJoypad_JoypadRight: // ->

		return joypad.ButtonRight, nil

	case config.EJoypad_JoypadA: // A

		return joypad.ButtonA, nil

	case config.EJoypad_JoypadLeft: // <-

		return joypad.ButtonLeft, nil

	case config.EJoypad_JoypadB: // B

		return joypad.ButtonB, nil

	case config.EJoypad_JoypadUp: // ^

		return joypad.ButtonUp, nil

	case config.EJoypad_JoypadSelect: // select

		return joypad.ButtonSelect, nil

	case config.EJoypad_JoypadDown: // V

		return joypad.ButtonDown, nil

	case config.EJoypad_JoypadStart: // start

		return joypad.ButtonStart, nil
	}

	return 0, ErrInvalidJoypadCode
}

// GetKeystroke returns a gameboy key code
func (i *Input) GetKeystroke() *joypad.Keystroke {

	i.m.Lock()
	defer i.m.Unlock()

	if len(i.keystrokes) == 0 {
		return nil
	}

	ks := i.keystrokes[0]

	i.keystrokes = i.keystrokes[1:]

	return &ks
}

// Stop waiting for key events
func (i *Input) Stop() {
	i.stop = true
}

// AddKeyEvent to queue
func (i *Input) AddKeyEvent(code sdl.Keycode, pressed bool) {

	if v, found := i.mapping[int32(code)]; found {

		if btn, err := i.convertJoypadCode(v); err == nil {

			i.m.Lock()

			i.keystrokes = append(i.keystrokes, joypad.Keystroke{Button: btn, Pressed: pressed})

			if len(i.keystrokes) > maxQueueSize {
				i.keystrokes = i.keystrokes[1:]
			}

			i.m.Unlock()
		}
	}
}

// WaitForKeyEvents blocks until a key is pressed or unpressed
func (i *Input) WaitForKeyEvents() ControlEvent {

	i.keystrokes = nil
	i.stop = false

	for !i.stop {

		ev := sdl.WaitEvent()

		if ev == nil {
			continue
		}

		switch t := ev.(type) {

		case *sdl.QuitEvent:

			return ControlEventQuit

		case *sdl.KeyDownEvent:

			if t.Keysym.Sym == sdl.K_ESCAPE {

				return ControlEventQuit
			}

			if t.Keysym.Sym == sdl.K_F1 {

				return ControlEventReset
			}

			if t.Keysym.Sym == sdl.K_F2 {

				return ControlEventPause
			}

			i.AddKeyEvent(t.Keysym.Sym, true)

		case *sdl.KeyUpEvent:

			i.AddKeyEvent(t.Keysym.Sym, false)
		}
	}

	return ControlEventQuit
}
