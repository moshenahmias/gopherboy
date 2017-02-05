/*H**********************************************************************
* FILENAME :        ui.go
*
* PACKAGE :			ui
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package ui

import "github.com/veandco/go-sdl2/sdl"

// Initialize things
func Initialize() error {

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return err
	}

	return nil
}

// Close things
func Close() {
	sdl.Quit()
}
