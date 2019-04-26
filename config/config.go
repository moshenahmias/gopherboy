package config

import (
	"bufio"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/jsonpb"
)

// DefaultSettingsFile is the default settings file name
const DefaultSettingsFile = "settings.json"

var DefaultSettings = Settings{

	JoypadMapping: map[int32]EJoypad{
		1073741906: EJoypad_JoypadUp,
		1073741905: EJoypad_JoypadDown,
		1073741904: EJoypad_JoypadLeft,
		1073741903: EJoypad_JoypadRight,
		122:        EJoypad_JoypadA,
		120:        EJoypad_JoypadB,
		32:         EJoypad_JoypadSelect,
		13:         EJoypad_JoypadStart},

	SoundDevice: 0,
	Fps:         60,
	Scale:       2,
	Color_0:     0x00C3D6AA,
	Color_1:     0x008EA86C,
	Color_2:     0x004D642D,
	Color_3:     0x00283A10}

// createDefaultSettingsFile in the given path
func createDefaultSettingsFile(path string) error {

	// create file
	f, err := os.Create(path)

	if err != nil {
		return err
	}

	defer f.Close()

	// init writer
	w := bufio.NewWriter(f)
	defer w.Flush()

	// marshal json to file
	var m jsonpb.Marshaler

	m.EmitDefaults = true
	m.Indent = "    "

	if err := m.Marshal(w, &DefaultSettings); err != nil {
		return err
	}

	return nil
}

// LoadSettings loads the settings file from the disk
func LoadSettings(path string) (*Settings, error) {

	// check if file exists
	if _, err := os.Stat(path); err != nil {

		// create default config file
		if err := createDefaultSettingsFile(path); err != nil {
			return nil, err
		}
	}

	// read file from disk
	b, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	// unmarshal json to settings object
	settings := &Settings{}

	if err := jsonpb.UnmarshalString(string(b), settings); err != nil {
		return nil, err
	}

	return settings, nil
}
