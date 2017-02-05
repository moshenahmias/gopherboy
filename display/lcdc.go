/*H**********************************************************************
* FILENAME :        lcdc.go
*
* PACKAGE :			display
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package display

// LCDC is the LCD Control register
type LCDC byte

// Read from the register
func (l *LCDC) Read(addr uint16) (byte, error) {
	return byte(*l), nil
}

// Write to the register
func (l *LCDC) Write(addr uint16, data byte) error {
	*l = LCDC(data)
	return nil
}

// displayEnabled returns true iff bit 7 is set
func (l *LCDC) displayEnabled() bool {
	return byte(*l)&0x80 == 0x80
}

// windowEnabled returns true iff bit 5 is set
func (l *LCDC) windowEnabled() bool {
	return byte(*l)&0x20 == 0x20
}

// spritesEnabled returns true iff bit 1 is set
func (l *LCDC) spritesEnabled() bool {
	return byte(*l)&0x02 == 0x02
}

// backgroundEnabled returns true iff bit 0 is set
func (l *LCDC) backgroundEnabled() bool {
	return byte(*l)&0x01 == 0x01
}

// spriteWidth returns the sprite width (8 or 16)
func (l *LCDC) spriteWidth() byte {

	if byte(*l)&0x04 == 0x04 {
		return 16
	}

	return 8
}

// tileset returns the tileset id (0 or 1)
func (l *LCDC) tileset() byte {
	return (byte(*l) << 3) >> 7
}

// backgroundMap returns the map id for the background (0 or 1)
func (l *LCDC) backgroundMap() byte {
	return (byte(*l) << 4) >> 7
}

// windowMap returns the map id for the window (0 or 1)
func (l *LCDC) windowMap() byte {
	return (byte(*l) << 1) >> 7
}
