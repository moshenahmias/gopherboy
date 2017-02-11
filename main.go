/*H**********************************************************************
* FILENAME :        main.go
*
* PACKAGE :			main
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package main

import (
	"flag"
	"io/ioutil"
	"runtime"
	"sync"

	"github.com/moshenahmias/gopherboy/audio"
	"github.com/moshenahmias/gopherboy/config"
	"github.com/moshenahmias/gopherboy/cpu"
	"github.com/moshenahmias/gopherboy/display"
	"github.com/moshenahmias/gopherboy/game"
	"github.com/moshenahmias/gopherboy/joypad"
	"github.com/moshenahmias/gopherboy/memory"
	"github.com/moshenahmias/gopherboy/ui"

	"github.com/sirupsen/logrus"
)

func main() {

	// set logging level
	logrus.SetLevel(logrus.InfoLevel)

	// init command-line arguments
	argROM := flag.String("rom", "", "Path to game ROM")
	argBIOS := flag.String("bios", "", "Path to boot ROM")
	argSettings := flag.String("settings", config.DefaultSettingsFile, "Path to settings file")

	// parse command-line arguments
	flag.Parse()

	// validate rom file arg
	if len(*argROM) == 0 {
		flag.PrintDefaults()
		return
	}

	// load settings
	settings, err := config.LoadSettings(*argSettings)

	if err != nil {
		logrus.Error(err)
		return
	}

	// run
	if err := run(*argROM, *argBIOS, settings); err != nil {
		logrus.Error(err)
	}
}

func run(romFile, biosFile string, settings *config.Settings) error {

	runtime.LockOSThread()

	// init ui
	if err := ui.Initialize(); err != nil {
		return err
	}

	defer ui.Close()

	// create input listener
	input := ui.NewInput(settings.JoypadMapping)

	// create and show the window
	window, err := ui.NewWindow(
		"gopherboy",
		int(settings.Scale),
		settings.Color_0,
		settings.Color_1,
		settings.Color_2,
		settings.Color_3)

	if err != nil {
		return err
	}

	defer window.Close()

	var sound ui.Sound

	if settings.SoundDevice >= 0 {

		if err := sound.Initialize(int(settings.SoundDevice)); err != nil {
			return err
		}

		defer sound.Close()
	}

	// load bios (if available)
	var biosData []byte

	if len(biosFile) > 0 {

		var err error

		// load bios from file
		biosData, err = ioutil.ReadFile(biosFile)

		if err != nil {
			return err
		}
	}

	soundMute := false

	// restart loop
	for quit := false; !quit; {

		// set sound mute state
		sound.Mute(soundMute)

		// create mmu
		mmu := memory.NewMMU()

		// TODO: map missing IO registers to FF00-FF7F
		if err := mmu.Map(memory.NewRAM(make([]byte, 128), 0xFF00), 0xFF00, 0xFF7F); err != nil {
			return err
		}

		// create cpu core
		core, err := cpu.NewCore(mmu)

		if err != nil {
			return err
		}

		// create and map the joyp register
		joyp := joypad.NewJOYP(core, input)

		// create gpu
		gpu, err := display.NewGPU(mmu, window, core, settings.Fps)

		if err != nil {
			return err
		}

		// create apu
		_, err = audio.NewAPU(core, mmu, &sound)

		if err != nil {
			return err
		}

		// load cartridg
		cartridge, err := game.NewCartridge(romFile, core)

		if err != nil {
			return err
		}

		// assemble everything
		gameboy, err := NewGameboy(cartridge, mmu, core, biosData, joyp, gpu)

		if err != nil {
			return err
		}

		// start the game
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {

			if err := gameboy.Start(); err != nil {
				logrus.Error(err)
				input.Stop()
			}

			wg.Done()
		}()

		// wait for keyboard events
		keyEvent := input.WaitForKeyEvents()

		for keyEvent == ui.ControlEventPause || keyEvent == ui.ControlEventMute {

			// pause
			if keyEvent == ui.ControlEventPause {
				gameboy.Pause()
			}

			// mute
			if keyEvent == ui.ControlEventMute {
				soundMute = !soundMute
				sound.Mute(soundMute)
			}

			keyEvent = input.WaitForKeyEvents()
		}

		// quit
		quit = keyEvent == ui.ControlEventQuit

		// stop cpu
		gameboy.Stop()

		wg.Wait()
	}

	// bye
	return nil
}
