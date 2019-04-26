package main

import (
	"github.com/moshenahmias/gopherboy/cpu"
	"github.com/moshenahmias/gopherboy/display"
	"github.com/moshenahmias/gopherboy/game"
	"github.com/moshenahmias/gopherboy/joypad"
	"github.com/moshenahmias/gopherboy/memory"
	"github.com/moshenahmias/gopherboy/timers"
)

// Gameboy console
type Gameboy struct {
	core *cpu.Core
	mmu  *memory.MMU
	bios bool
}

// NewGameboy creates Gameboy instance
func NewGameboy(
	cartridge *game.Cartridge,
	mmu *memory.MMU,
	core *cpu.Core,
	biosData []byte,
	joyp *joypad.JOYP,
	gpu *display.GPU) (*Gameboy, error) {

	// map FEA0-FEFF (unused)
	if err := mmu.Map(&memory.Null{}, 0xFEA0, 0xFEFF); err != nil {
		return nil, err
	}

	// map joypad
	if err := mmu.Map(joyp, joypad.AddrJOYP, joypad.AddrJOYP); err != nil {
		return nil, err
	}

	// create and map the timer
	timer := timers.NewTimer(core)

	// map timer
	if err := mmu.Map(timer, 0xFF04, 0xFF07); err != nil {
		return nil, err
	}

	bios := len(biosData) > 0

	if bios {

		// map bios
		if err := mmu.Map(memory.NewROM(biosData, 0), 0x0000, 0x00FF); err != nil {
			return nil, err
		}

		// map cartridge
		if err := mmu.Map(cartridge, 0x0100, 0x7FFF); err != nil {
			return nil, err
		}

	} else {

		// map cartridge
		if err := mmu.Map(cartridge, 0x0000, 0x7FFF); err != nil {
			return nil, err
		}
	}

	// map external ram
	if err := mmu.Map(cartridge, 0xA000, 0xBFFF); err != nil {
		return nil, err
	}

	// map bios unmapper
	if err := mmu.Map(memory.NewBiosUnmapper(mmu, cartridge), 0xFF50, 0xFF50); err != nil {
		return nil, err
	}

	// map working ram
	wram := memory.NewRAM(make([]byte, 8192), 0xC000)
	if err := mmu.Map(wram, 0xC000, 0xDFFF); err != nil {
		return nil, err
	}

	// map shadow
	for src := uint16(0xE000); src <= uint16(0xFDFF); src++ {
		if err := mmu.Map(memory.NewEchoer(mmu, src-0x2000), src, src); err != nil {
			return nil, err
		}
	}

	// map zero page information ram
	zpram := memory.NewRAM(make([]byte, 127), 0xFF80)
	if err := mmu.Map(zpram, 0xFF80, 0xFFFE); err != nil {
		return nil, err
	}

	core.RegisterToClockChanges(gpu)

	return &Gameboy{core: core, mmu: mmu, bios: bios}, nil
}

// Start the cpu
func (g *Gameboy) Start() error {

	if g.bios {
		return g.core.Start(0x0000)
	}

	return g.core.Start(0x0100)
}

// Pause the cpu
func (g *Gameboy) Pause() {
	g.core.Pause()
}

// Stop the cpu
func (g *Gameboy) Stop() {
	g.core.Stop()
}
