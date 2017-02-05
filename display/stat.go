/*H**********************************************************************
* FILENAME :        stat.go
*
* PACKAGE :			display
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package display

// ModeDuringHBlank - The LCD controller is in the H-Blank period and
// the CPU can access both the display RAM (8000h-9FFFh) and OAM (FE00h-FE9Fh)
const ModeDuringHBlank byte = 0x00

// ModeDuringVBlank - The LCD contoller is in the V-Blank period (or the
// display is disabled) and the CPU can access both the display RAM (8000h-9FFFh)
// and OAM (FE00h-FE9Fh)
const ModeDuringVBlank byte = 0x01

// ModeSearchingOAM - The LCD controller is reading from OAM memory.
// The CPU <cannot> access OAM memory (FE00h-FE9Fh) during this period.
const ModeSearchingOAM byte = 0x02

// ModeTransferingDataToLCD - The LCD controller is reading from both OAM and VRAM,
// The CPU <cannot> access OAM and VRAM during this period. CGB Mode: Cannot access
// Palette Data (FF69,FF6B) either.
const ModeTransferingDataToLCD byte = 0x03

// STAT is the LCD Status register
type STAT byte

// coincidenceInterruptEnabled returns true iff bit 6 is set
func (s *STAT) coincidenceInterruptEnabled() bool {
	return *s&0x40 == 0x40
}

// oamInterruptEnabled returns true iff bit 5 is set
func (s *STAT) oamInterruptEnabled() bool {
	return *s&0x20 == 0x20
}

// vBlankInterruptEnabled returns true iff bit 4 is set
func (s *STAT) vBlankInterruptEnabled() bool {
	return *s&0x10 == 0x10
}

// hBlankInterruptEnabled returns true iff bit 3 is set
func (s *STAT) hBlankInterruptEnabled() bool {
	return *s&0x08 == 0x08
}

// modeFlag return the current mode
func (s *STAT) modeFlag() byte {
	return byte(*s) & 0x03
}

// setModeFlag to 0, 1, 2, or 3
func (s *STAT) setModeFlag(mode byte) {
	if mode < 4 {
		*s = STAT((byte(*s) & 0xFC) | mode)
	}
}

// coincidenceFlag returns true iff bit 2 is set
func (s *STAT) coincidenceFlag() bool {
	return *s&0x04 == 0x04
}

// setCoincidenceFlag by setting or resetting bit 2
func (s *STAT) setCoincidenceFlag(state bool) {
	if state {
		*s = *s | 0x04
	} else {
		*s = *s & 0xFB
	}
}

// Read from the register
func (s *STAT) Read(addr uint16) (byte, error) {
	return byte(*s), nil
}

// Write to the register
func (s *STAT) Write(addr uint16, data byte) error {
	*s = STAT((byte(*s) & 0x07) | (data & 0xF8))
	return nil
}
