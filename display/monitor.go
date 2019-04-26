package display

// Pixel represents a single pixel color in the monitor
type Pixel byte

// PixelWhite is color 0
const PixelWhite Pixel = 0

// PixelLightGray is color 1
const PixelLightGray Pixel = 1

// PixelDarkGray is color 2
const PixelDarkGray Pixel = 2

// PixelBlack is color 3 (darkest)
const PixelBlack Pixel = 3

// ScreenWidth is the physical screen width
const ScreenWidth int = 160

// ScreenHeight is the physical screen height
const ScreenHeight int = 144

// Frame is a 160x144 matrix that represents a single frame
type Frame [ScreenWidth][ScreenHeight]Pixel

// Monitor of the gameboy
type Monitor interface {
	Close() error
	DrawFrame(f *Frame) error
}
