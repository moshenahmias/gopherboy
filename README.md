# gopherboy
A Nintendo Game Boy emulator written in golang.

1. Passes blargg's cpu tests 01 - 11.
2. Currently supports only none-MBC and MBC1 roms.
3. Sound not implemented.

### TODO list

1. MBC2.
2. Save / Load game state.
3. Sound.

### Dependencies

1. go-sdl2: https://github.com/veandco/go-sdl2
2. logrus: https://github.com/sirupsen/logrus
3. protobuf: https://github.com/golang/protobuf

### How to build

1. Follow the instructions for go-sdl2, logrus and protobuf in the links above.
2. go get github.com/moshenahmias/gopherboy
3. go install github.com/moshenahmias/gopherboy

### How to run

```
Usage of gopherboy:

  -bios string
        Path to boot ROM
        
  -rom string
        Path to game ROM
        
  -settings string
        Path to settings file (default "settings.json")    
```

### Default Keyboard Mapping

| Operation     | Key           |
| ------------- |:-------------:| 
| Up            | Up Arrow      |
| Down          | Down Arrow    |
| Left          | Left Arrow    | 
| Right         | Right Arrow   | 
| A             | Z             | 
| B             | X             | 
| Select        | Space         | 
| Start         | Enter         | 
| Reset         | F1            | 
| Exit          | ESC           | 

### Settings

You can change the following settings via the *settings.json* file:

##### Joypad Mapping:

```
"joypadMapping": {
        "<sdl keycode>": <joypad keycode>
        ...}       
```

Joypad keycodes: Up = 0, Down = 1, Left = 2, Right = 3, A = 4, B = 5, Select = 6, Start = 7

SDL keycodes: https://wiki.libsdl.org/SDLKeycodeLookup

##### Other Settings:

FPS rate, screen size and color palette.

### Screenshots

![Super Mario Land](images/gopherboy1.png)&nbsp;
![Contra - The Alien Wars](images/gopherboy2.png)&nbsp;
![BattleCity](images/gopherboy3.png)&nbsp;
![Tetris](images/gopherboy4.png)&nbsp;
