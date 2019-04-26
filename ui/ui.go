package ui

import "github.com/veandco/go-sdl2/sdl"

// Initialize things
func Initialize() error {

	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO); err != nil {
		return err
	}

	return nil
}

// Close things
func Close() {
	sdl.Quit()
}
