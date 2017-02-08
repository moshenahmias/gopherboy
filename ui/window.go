/*H**********************************************************************
* FILENAME :        window.go
*
* PACKAGE :			ui
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package ui

import (
	"unsafe"

	"github.com/moshenahmias/gopherboy/display"

	"github.com/veandco/go-sdl2/sdl"
)

// Window represents the device's lcd screen
type Window struct {
	scale    int
	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture
	colors   [4]uint32
}

// NewWindow creates Window instance
func NewWindow(title string, scale int, color0, color1, color2, color3 uint32) (*Window, error) {

	if scale < 1 {
		scale = 1
	}

	// calc window width and height
	width := display.ScreenWidth * scale
	height := display.ScreenHeight * scale

	// create the windows
	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width, height, sdl.WINDOW_SHOWN)

	if err != nil {
		return nil, err
	}

	// create the renderer
	renderer, err := sdl.CreateRenderer(window, -1, 0)

	if err != nil {
		window.Destroy()
		return nil, err
	}

	// create the texture
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, width, height)

	if err != nil {
		renderer.Destroy()
		window.Destroy()
		return nil, err
	}

	lcd := Window{scale: scale, window: window, renderer: renderer, texture: texture}

	lcd.colors[0] = color0
	lcd.colors[1] = color1
	lcd.colors[2] = color2
	lcd.colors[3] = color3

	return &lcd, nil
}

// Close destroys the texture, renderer and window
func (l *Window) Close() error {

	l.texture.Destroy()
	l.renderer.Destroy()
	l.window.Destroy()

	return nil
}

// DrawFrame to window
func (l *Window) DrawFrame(f *display.Frame) error {

	var pixels unsafe.Pointer
	pitch := 0

	if err := l.texture.Lock(nil, &pixels, &pitch); err != nil {
		return err
	}

	for j := 0; j < display.ScreenHeight; j++ {
		for i := 0; i < display.ScreenWidth; i++ {
			for x := i * l.scale; x < i*l.scale+l.scale; x++ {
				for y := j * l.scale; y < j*l.scale+l.scale; y++ {
					p := unsafe.Pointer(uintptr(pixels) + unsafe.Sizeof(uint32(0))*uintptr(y*(pitch/4)+x))
					*(*uint32)(p) = l.color((*f)[i][j])
				}
			}
		}
	}

	l.texture.Update(nil, pixels, pitch)
	l.texture.Unlock()
	l.renderer.Clear()
	l.renderer.Copy(l.texture, nil, nil)
	l.renderer.Present()

	return nil
}

// color converts color code to the actual color
func (l *Window) color(code display.Pixel) uint32 {

	if code > 3 {
		return 0xFFFFFF
	}

	return l.colors[code]
}
