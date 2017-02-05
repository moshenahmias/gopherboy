/*H**********************************************************************
* FILENAME :        palette.go
*
* PACKAGE :			display
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package display

// Color in palette
type Color byte

// ColorWhite represents color 0
const ColorWhite Color = 0

// ColorLightGray represents color 1
const ColorLightGray Color = 1

// ColorDarkGray represents color 2
const ColorDarkGray Color = 2

// ColorBlack represents color 3
const ColorBlack Color = 3

// ColorTransparent represents no-color
const ColorTransparent Color = 4

// Palette register
type Palette byte

// toColor from color code
func (p *Palette) toColor(code byte) Color {
	return Color((byte(*p) << (8 - ((code + 1) * 2))) >> 6) // 0 - 3
}

// Read from the register
func (p *Palette) Read(addr uint16) (byte, error) {
	return byte(*p), nil
}

// Write to the register
func (p *Palette) Write(addr uint16, data byte) error {
	*p = Palette(data)
	return nil
}
