/*H**********************************************************************
* FILENAME :        gpu.go
*
* PACKAGE :			display
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package display

import (
	"fmt"
	"sort"
	"time"

	"github.com/moshenahmias/gopherboy/cpu"
	"github.com/moshenahmias/gopherboy/memory"
)

// AddrLCDC is the LCDC register address
const AddrLCDC uint16 = 0xFF40

// AddrSTAT is the STAT register address
const AddrSTAT uint16 = 0xFF41

// AddrSCY is the Scroll Y register address
const AddrSCY uint16 = 0xFF42

// AddrSCX is the Scroll X register address
const AddrSCX uint16 = 0xFF43

// AddrLY is the Y-Coordinate register address
const AddrLY uint16 = 0xFF44

// AddrLYC is the LY Compare register address
const AddrLYC uint16 = 0xFF45

// AddrDMA is the DMA transfer unit address
const AddrDMA uint16 = 0xFF46

// AddrBGP is the BG & Window Palette Data register address
const AddrBGP uint16 = 0xFF47

// AddrOBP0 is the Object Palette 0 Data register address
const AddrOBP0 uint16 = 0xFF48

// AddrOBP1 is the Object Palette 0 Data register address
const AddrOBP1 uint16 = 0xFF49

// AddrWinY is the Window Y Position register address
const AddrWinY uint16 = 0xFF4A

// AddrWinX is the Window Y Position register address
const AddrWinX uint16 = 0xFF4B

// Layer is a part of a single frame
type Layer [ScreenWidth][ScreenHeight]Color

// GPU renders the background, window and sprites
type GPU struct {
	monitor Monitor
	core    *cpu.Core
	mmu     *memory.MMU

	lcdc LCDC
	scy  memory.MemReg
	scx  memory.MemReg
	wy   memory.MemReg
	wx   memory.MemReg
	lyc  memory.MemReg
	ly   byte
	lx   byte

	bgp Palette
	obp [2]Palette

	stat STAT

	bgLayer     Layer
	winLayer    Layer
	spriteLayer [2]Layer

	vram *memory.RAM
	oam  *memory.RAM

	sprites SpriteAttrs

	cyclesCounter int

	displayEnabled    bool
	spritesEnabled    bool
	backgroundEnabled bool

	ignoreVBlankInt bool
	ignoreHBlankInt bool
	ignoreLYCInt    bool
	ignoreOAMInt    bool

	frameDuration time.Time
	fps           int64
}

// NewGPU creates GPU instance
func NewGPU(mmu *memory.MMU, monitor Monitor, core *cpu.Core, fps uint32) (*GPU, error) {

	if fps == 0 {
		fps = 1
	}

	g := GPU{mmu: mmu, monitor: monitor, core: core, fps: int64(fps)}

	g.lcdc = 0x91

	// lcdc
	if err := mmu.Map(&g.lcdc, AddrLCDC, AddrLCDC); err != nil {
		return nil, err
	}

	// stat
	if err := mmu.Map(&g.stat, AddrSTAT, AddrSTAT); err != nil {
		return nil, err
	}

	// scy
	if err := mmu.Map(&g.scy, AddrSCY, AddrSCY); err != nil {
		return nil, err
	}

	// scx
	if err := mmu.Map(&g.scx, AddrSCX, AddrSCX); err != nil {
		return nil, err
	}

	// wx
	if err := mmu.Map(&g.wx, AddrWinX, AddrWinX); err != nil {
		return nil, err
	}

	// wy
	if err := mmu.Map(&g.wy, AddrWinY, AddrWinY); err != nil {
		return nil, err
	}

	// ly
	if err := mmu.Map(&g, AddrLY, AddrLY); err != nil {
		return nil, err
	}

	// lyc
	if err := mmu.Map(&g.lyc, AddrLYC, AddrLYC); err != nil {
		return nil, err
	}

	// bgp
	if err := mmu.Map(&g.bgp, AddrBGP, AddrBGP); err != nil {
		return nil, err
	}

	// obp0
	if err := mmu.Map(&g.obp[0], AddrOBP0, AddrOBP0); err != nil {
		return nil, err
	}

	// obp1
	if err := mmu.Map(&g.obp[1], AddrOBP1, AddrOBP1); err != nil {
		return nil, err
	}

	// dma
	if err := mmu.Map(&g, AddrDMA, AddrDMA); err != nil {
		return nil, err
	}

	// vram
	g.vram = memory.NewRAM(make([]byte, 8192), 0x8000)
	if err := mmu.Map(&g, 0x8000, 0x9FFF); err != nil {
		panic(err)
	}

	// oam
	g.oam = memory.NewRAM(make([]byte, 160), 0xFE00)
	if err := mmu.Map(&g, 0xFE00, 0xFE9F); err != nil {
		panic(err)
	}

	g.displayEnabled = false
	g.initialize()

	return &g, nil
}

// Read from ly, dma, oam or vram
func (g *GPU) Read(addr uint16) (byte, error) {

	if addr == AddrLY {
		return g.ly, nil
	}

	if addr == AddrDMA {
		return 0, nil
	}

	if 0x8000 <= addr && addr <= 0x9FFF {
		return g.vram.Read(addr)
	}

	if 0xFE00 <= addr && addr <= 0xFE9F {
		return g.oam.Read(addr)
	}

	return 0, memory.ReadOutOfRangeError(addr)
}

// Write to ly, dma, oam or vram
func (g *GPU) Write(addr uint16, data byte) error {

	if addr == AddrLY {
		g.initialize()
		return nil
	}

	if addr == AddrDMA {

		//if g.stat.ModeFlag() == ModeSearchingOAM || g.stat.ModeFlag() == ModeTransferingDataToLCD {
		//	return nil
		//}

		return g.transferDataToOAM(data)
	}

	if 0x8000 <= addr && addr <= 0x9FFF {

		//if g.stat.ModeFlag() == ModeTransferingDataToLCD {
		//	return nil
		//}

		return g.vram.Write(addr, data)
	}

	if 0xFE00 <= addr && addr <= 0xFE9F {

		//if g.stat.ModeFlag() == ModeSearchingOAM || g.stat.ModeFlag() == ModeTransferingDataToLCD {
		//	return nil
		//}

		return g.oam.Write(addr, data)
	}

	return memory.WriteOutOfRangeError(addr)
}

// transferDataToOAM from the given address * 100
func (g *GPU) transferDataToOAM(from byte) error {

	src := uint16(from) << 8

	for i := uint16(0); i < 0xA0; i++ {

		v, err := g.mmu.Read(src + i)

		if err != nil {
			return err
		}

		if err := g.mmu.Write(0xFE00+i, v); err != nil {
			return err
		}
	}

	return nil
}

// scanPixels scans next n (or less) pixels
// returns true at the end of the scanline
func (g *GPU) scanPixels(n byte) (bool, error) {

	for ; g.lx < 160 && n > 0; g.lx++ {
		if err := g.renderPixel(g.lx, g.ly); err != nil {
			return false, err
		}

		if g.lx%16 == 0 {
			if err := g.renderSprite(); err != nil {
				return false, err
			}
		}

		n--
	}

	if g.lx == 160 {
		g.lx = 0
		return true, nil
	}

	return false, nil
}

// renderPixel (x,y) for the background and window
func (g *GPU) renderPixel(x, y byte) error {

	if err := g.renderBackgroundPixel(x, y); err != nil {
		return err
	}

	if err := g.renderWindowPixel(x, y); err != nil {
		return err
	}

	return nil
}

// initialize the GPU
func (g *GPU) initialize() {

	g.frameDuration = time.Now()
	g.ly = 0
	g.lx = 0
	g.cyclesCounter = 0
	g.stat.setCoincidenceFlag(false)
	g.stat.setModeFlag(ModeSearchingOAM)
	g.ignoreHBlankInt = false
	g.ignoreLYCInt = false
	g.ignoreOAMInt = false
	g.ignoreVBlankInt = false
	g.sprites = nil

	g.spritesEnabled = g.lcdc.spritesEnabled()
	g.backgroundEnabled = g.lcdc.backgroundEnabled()

	for x := 0; x < ScreenWidth; x++ {
		for y := 0; y < ScreenHeight; y++ {
			g.winLayer[x][y] = ColorTransparent
			g.bgLayer[x][y] = ColorWhite
			g.spriteLayer[SpriteAboveBackground][x][y] = ColorTransparent
			g.spriteLayer[SpriteBelowBackground][x][y] = ColorTransparent
		}
	}
}

// ClockChanged is called after every instruction execution
func (g *GPU) ClockChanged(cycles int) error {

	if !g.lcdc.displayEnabled() {

		if g.displayEnabled {

			// shut down the monitor

			g.initialize()

			if err := g.updateMonitor(); err != nil {
				return err
			}

			g.displayEnabled = false

			return nil
		}

		return nil
	}

	g.displayEnabled = true

	//////////////////////
	// Coincidence flag //
	//////////////////////

	if byte(g.ly) == byte(g.lyc) && !g.stat.coincidenceFlag() {

		g.stat.setCoincidenceFlag(true)

		if g.stat.coincidenceInterruptEnabled() && !g.ignoreLYCInt {

			// request lcd status interrupt
			g.core.RequestInterrupt(cpu.LCDStatusTriggersFlag)

			if g.lyc <= 143 {

				if g.lyc > 0 {
					g.ignoreOAMInt = true
				}

				if g.lyc == 143 {
					g.ignoreVBlankInt = true
				}

				g.ignoreHBlankInt = true
			}

			return nil
		}

	} else if byte(g.ly) != byte(g.lyc) {

		g.stat.setCoincidenceFlag(false)
	}

	switch g.stat.modeFlag() {

	//////////////////////////////
	// During Searching OAM-RAM //
	//////////////////////////////

	case ModeSearchingOAM:

		g.cyclesCounter += cycles

		if g.cyclesCounter < 80 {
			return nil
		}

		g.cyclesCounter = g.cyclesCounter - 80

		for x := 0; x < ScreenWidth; x++ {
			g.spriteLayer[SpriteAboveBackground][x][g.ly] = ColorTransparent
			g.spriteLayer[SpriteBelowBackground][x][g.ly] = ColorTransparent
		}

		if g.spritesEnabled {

			attrs, err := g.searchOAM()

			if err != nil {
				return err
			}

			if len(attrs) > 10 {
				attrs = attrs[0:10]
			}

			g.sprites = attrs
		}

		g.stat.setModeFlag(ModeTransferingDataToLCD)

	///////////////////////////////////////////
	// During Transfering Data to LCD Driver //
	///////////////////////////////////////////

	case ModeTransferingDataToLCD:

		hb, err := g.scanPixels(byte(cycles))

		if err != nil {
			return err
		}

		if hb {

			g.stat.setModeFlag(ModeDuringHBlank)

			if g.stat.hBlankInterruptEnabled() && !g.ignoreHBlankInt {
				// request lcd status interrupt
				g.core.RequestInterrupt(cpu.LCDStatusTriggersFlag)

				if g.ly == 143 {
					g.ignoreVBlankInt = true
				}

				g.ignoreOAMInt = true
			}

			g.ignoreHBlankInt = false

			g.ly++
		}

	////////////////////
	// During H-Blank //
	////////////////////

	case ModeDuringHBlank:

		g.cyclesCounter += cycles

		if g.cyclesCounter < 204 {
			return nil
		}

		g.cyclesCounter = g.cyclesCounter - 204

		if g.ly == 144 {

			// draw rendered frame
			if err := g.updateMonitor(); err != nil {
				return err
			}

			g.stat.setModeFlag(ModeDuringVBlank)

			if g.stat.vBlankInterruptEnabled() && !g.ignoreVBlankInt {

				// request lcd status interrupt
				g.core.RequestInterrupt(cpu.LCDStatusTriggersFlag)

				g.ignoreOAMInt = true
				g.ignoreLYCInt = true
			}

			g.ignoreVBlankInt = false

			// request vertical blank interrupt
			g.core.RequestInterrupt(cpu.VerticalBlankFlag)

		} else if g.ly < 144 {

			g.stat.setModeFlag(ModeSearchingOAM)

			if g.stat.oamInterruptEnabled() && !g.ignoreOAMInt {
				// request lcd status interrupt
				g.core.RequestInterrupt(cpu.LCDStatusTriggersFlag)
			}

			g.ignoreOAMInt = false

		} else {

			return fmt.Errorf("ly > 144 during h-blank (ly = %d)", g.ly)
		}

	////////////////////
	// During V-Blank //
	////////////////////

	case ModeDuringVBlank:

		g.cyclesCounter += cycles

		if g.cyclesCounter < 456 {
			return nil
		}

		g.cyclesCounter = g.cyclesCounter - 456

		g.ly++

		if g.ly > 153 {

			g.ignoreLYCInt = false

			g.spritesEnabled = g.lcdc.spritesEnabled()
			g.backgroundEnabled = g.lcdc.backgroundEnabled()

			g.ly = 0
			g.stat.setModeFlag(ModeSearchingOAM)

			if g.stat.oamInterruptEnabled() && !g.ignoreOAMInt {

				// request lcd status interrupt
				g.core.RequestInterrupt(cpu.LCDStatusTriggersFlag)
			}

			g.ignoreOAMInt = false
		}
	}

	return nil
}

// renderBackgroundPixel renders (x, y) pixel of the background
func (g *GPU) renderBackgroundPixel(x, y byte) error {

	if !g.backgroundEnabled {
		g.bgLayer[x][y] = ColorWhite
		return nil
	}

	scx := byte(g.scx)
	scy := byte(g.scy)

	// find the tile in map and pixel within the tile for (x, y)
	tileX := (x + scx) / 8  // 0 - 31
	tileY := (y + scy) / 8  // 0 - 31
	tilePX := (x + scx) % 8 // 0 - 7
	tilePY := (y + scy) % 8 // 0 - 7

	// find the tile offset in map
	tileOff := uint16(tileX) + (uint16(tileY) * 32) // 0 - 1023

	// set map address
	var mapAddr uint16

	mapN := g.lcdc.backgroundMap()

	if mapN == 0 {
		mapAddr = 0x9800
	} else {
		mapAddr = 0x9C00
	}

	// get tile id
	tileID, err := g.mmu.Read(mapAddr + tileOff)

	if err != nil {
		return err
	}

	// set tile address
	tileAddr := g.getTileAddr(tileID)

	// get row bytes
	rowByte1N := tilePY * 2    // 0 - 15
	rowByte2N := rowByte1N + 1 // 0 - 15

	rowByte1, err := g.mmu.Read(tileAddr + uint16(rowByte1N))

	if err != nil {
		return err
	}

	rowByte2, err := g.mmu.Read(tileAddr + uint16(rowByte2N))

	if err != nil {
		return err
	}

	// extract color code
	rowByte1 = (rowByte1 << tilePX) >> 7
	rowByte2 = ((rowByte2 << tilePX) >> 7) << 1

	colorCode := rowByte1 | rowByte2

	// set color for pixel (x, y)
	g.bgLayer[x][y] = g.bgp.toColor(colorCode)

	return nil
}

// renderWindowPixel renders (x, y) pixel of the window
func (g *GPU) renderWindowPixel(x, y byte) error {

	if !g.lcdc.windowEnabled() {
		g.winLayer[x][y] = ColorTransparent
		return nil
	}

	wx := byte(g.wx)

	if wx < 7 {
		return nil
	}

	wy := byte(g.wy)

	if x < (wx-7) || y < wy {
		g.winLayer[x][y] = ColorTransparent
		return nil
	}

	// find the tile in map and pixel within the tile for (x, y)
	tileX := (x - (wx - 7)) / 8 // 0 - 19
	tileY := (y - wy) / 8       // 0 - 17

	tilePX := (x - (wx - 7)) % 8 // 0 - 7
	tilePY := (y - wy) % 8       // 0 - 7

	// find the tile offset in map
	tileOff := uint16(tileX) + (uint16(tileY) * 32) // 0 - 563

	// set map address
	var mapAddr uint16

	mapN := g.lcdc.windowMap()

	if mapN == 0 {
		mapAddr = 0x9800
	} else {
		mapAddr = 0x9C00
	}

	// get tile id
	tileID, err := g.mmu.Read(mapAddr + tileOff)

	if err != nil {
		return err
	}

	// set tile address
	tileAddr := g.getTileAddr(tileID)

	// get row bytes
	rowByte1N := tilePY * 2    // 0 - 15
	rowByte2N := rowByte1N + 1 // 0 - 15

	rowByte1, err := g.mmu.Read(tileAddr + uint16(rowByte1N))

	if err != nil {
		return err
	}

	rowByte2, err := g.mmu.Read(tileAddr + uint16(rowByte2N))

	if err != nil {
		return err
	}

	// extract color code
	rowByte1 = (rowByte1 << tilePX) >> 7
	rowByte2 = ((rowByte2 << tilePX) >> 7) << 1

	colorCode := rowByte1 | rowByte2

	// set color for pixel (x, y)
	g.winLayer[x][y] = g.bgp.toColor(colorCode)

	return nil
}

// getTileAddr from tile id and tileset
func (g *GPU) getTileAddr(id byte) uint16 {

	if g.lcdc.tileset() == 0 {

		if id&0x80 == 0 {
			return 0x9000 + (16 * uint16(id)) // id: 0 - 127
		}

		return 0x8800 + (16 * uint16(id+128)) // id: (-128) - (-1)
	}

	return 0x8000 + (16 * uint16(id)) // id: 0 - 255
}

// createFrame from background, window and sprites
func (g *GPU) createFrame() *Frame {

	var f Frame

	for x := 0; x < ScreenWidth; x++ {
		for y := 0; y < ScreenHeight; y++ {

			// merge bg and window
			if g.winLayer[x][y] == ColorTransparent {
				f[x][y] = Pixel(g.bgLayer[x][y])
			} else {
				f[x][y] = Pixel(g.winLayer[x][y])
			}

			// merge sprites
			if g.spriteLayer[SpriteAboveBackground][x][y] != ColorTransparent {

				f[x][y] = Pixel(g.spriteLayer[SpriteAboveBackground][x][y])

			} else {

				if f[x][y] == PixelWhite {

					if g.spriteLayer[SpriteBelowBackground][x][y] != ColorTransparent {

						f[x][y] = Pixel(g.spriteLayer[SpriteBelowBackground][x][y])
					}
				}
			}
		}
	}

	return &f
}

// updateMonitor with the bg, window and sprites rendered layers
func (g *GPU) updateMonitor() error {

	timeSinceLastFrame := time.Since(g.frameDuration).Nanoseconds()
	realTime := time.Second.Nanoseconds() / g.fps

	if timeSinceLastFrame < realTime {
		time.Sleep(time.Duration(realTime - timeSinceLastFrame))
	}

	err := g.monitor.DrawFrame(g.createFrame())

	g.frameDuration = time.Now()

	return err
}

/////////////
// Sprites //
/////////////

// getSpriteAttribute from OAM
func (g *GPU) getSpriteAttribute(addr uint16) (*SpriteAttr, error) {

	var attr SpriteAttr

	for i := 0; i < 4; i++ {

		if b, err := g.mmu.Read(addr + uint16(i)); err == nil {
			attr[i] = b
		} else {
			return nil, err
		}
	}

	return &attr, nil
}

// searchOAM for current ly sprites
func (g *GPU) searchOAM() (SpriteAttrs, error) {

	var attrs SpriteAttrs

	for id := 0; id < 40; id++ {

		attrAddr := 0xFE00 + uint16(id*4)

		// get the sprite attribute
		attr, err := g.getSpriteAttribute(attrAddr)

		if err != nil {
			return nil, err
		}

		w := g.lcdc.spriteWidth()

		if offScreen(attr.coordinateX(), 8, 8, 159) || offScreen(attr.coordinateY(), 16, w, 143) {
			// sprite is off screen
			continue
		}

		ys, ye, _ := calcCoords(attr.coordinateY(), 16, w, 143)

		if ys <= g.ly && g.ly <= ye {
			attrs = append(attrs, attr)
		}
	}

	sort.Sort(attrs)

	return attrs, nil
}

// renderSprite from sprites list
func (g *GPU) renderSprite() error {

	if len(g.sprites) == 0 {
		return nil
	}

	attr := g.sprites[0]

	y := g.ly

	w := g.lcdc.spriteWidth()

	xs, xe, spx := calcCoords(attr.coordinateX(), 8, 8, 159)
	ys, _, spy := calcCoords(attr.coordinateY(), 16, w, 143)
	flipX := attr.flipX()
	flipY := attr.flipY()
	priority := attr.priority()
	palette := g.obp[attr.palette()]
	tileID := attr.tileID()
	spy = spy + (y - ys)

	for x := xs; x <= xe; x++ {

		colorCode, err := g.spriteColorCode(tileID, w, spx, spy, flipX, flipY)

		if err != nil {
			return err
		}

		if g.spriteLayer[priority][x][y] == ColorTransparent && colorCode != 0 {

			g.spriteLayer[priority][x][y] = palette.toColor(colorCode)
		}

		spx++
	}

	g.sprites = g.sprites[1:]

	return nil
}

// spriteColorCode returns the color code for a given pixel in a given sprite
func (g *GPU) spriteColorCode(id byte, width byte, x, y byte, flipX, flipY bool) (byte, error) {

	// horizontal flip
	if flipX {
		x = 7 - x
	}

	// vertical flip
	if flipY {
		y = width - y - 1
	}

	var addr uint16

	if width == 8 {
		// 8x8
		addr = 0x8000 + (uint16(id)*16 + uint16(y)*2)
	} else {
		// 8x16
		if y < 8 {
			addr = 0x8000 + (uint16(id&0xFE)*16 + uint16(y)*2)
		} else {
			addr = 0x8000 + (uint16(id|0x01)*16 + uint16(y-8)*2)
		}
	}

	// get row bytes
	rowByte0, err := g.mmu.Read(addr)

	if err != nil {
		return 0, err
	}

	rowByte1, err := g.mmu.Read(addr + 1)

	if err != nil {
		return 0, err
	}

	// extract color code
	rowByte0 = (rowByte0 << x) >> 7
	rowByte1 = ((rowByte1 << x) >> 7) << 1

	return rowByte0 | rowByte1, nil
}

func offScreen(i, off, len, max byte) bool {
	return max+off < i || i <= off-len
}

func calcCoords(i, off, len, max byte) (byte, byte, byte) {

	var start byte
	var end byte
	var spr byte

	if off-len < i && i < off {

		start = 0
		end = len - (off - i) - 1
		spr = off - i

	} else if off <= i && i <= max+off-len+1 {

		start = i - off
		end = i - off + len - 1

	} else if max+off-len+2 <= i && i <= max+off {

		start = i - off
		end = max

	} else {

		// off screen
		panic("sprite off screen")
	}

	return start, end, spr
}
