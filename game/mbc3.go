package game

import "github.com/moshenahmias/gopherboy/memory"
import "github.com/moshenahmias/gopherboy/cpu"

import "time"

// MBC3 (max 2MByte ROM and/or 32KByte RAM and Timer)
type MBC3 struct {
	rom               *memory.ROM
	otherBanks        *memory.ROM
	ram               *memory.RAM
	bankRAM           uint32
	bankROM           uint32
	enableTimerAndRAM bool
	rtcCode           byte
	rtc               [5]byte
	rtcSnapshot       []byte
	latch             bool
	cyclesCounter     int
}

// NewMBC3 creates mbc1 instance
func NewMBC3(rom []byte, ram []byte) *MBC3 {

	m := MBC3{
		rom:           memory.NewROM(rom, 0x0000),
		otherBanks:    memory.NewROM(rom, 0x4000),
		ram:           memory.NewRAM(ram, 0xA000),
		bankROM:       1,
		rtcCode:       0,
		latch:         false,
		cyclesCounter: 0}

	now := time.Now()

	m.setDays(uint16(now.Day()))
	m.setHours(byte(now.Hour()))
	m.setMinutes(byte(now.Minute()))
	m.setSeconds(byte(now.Second()))

	return &m
}

func (m *MBC3) seconds() byte {
	return m.rtc[0]
}

func (m *MBC3) setSeconds(seconds byte) {
	m.rtc[0] = seconds
}

func (m *MBC3) minutes() byte {
	return m.rtc[1]
}

func (m *MBC3) setMinutes(minutes byte) {
	m.rtc[1] = minutes
}

func (m *MBC3) hours() byte {
	return m.rtc[2]
}

func (m *MBC3) setHours(hours byte) {
	m.rtc[2] = hours
}

func (m *MBC3) days() uint16 {
	return uint16(m.rtc[3]) | (uint16(m.rtc[4]&0x01) << 8)
}

func (m *MBC3) setDays(days uint16) {
	m.rtc[3] = byte(days & 0x00FF)
	m.rtc[4] = (m.rtc[4] & 0xFE) | byte((days>>8)&0x0001)
}

func (m *MBC3) halt() bool {
	return m.rtc[4]&0x40 == 0x40
}

func (m *MBC3) setDayCarry() {
	m.rtc[4] |= 0x80
}

// ClockChanged is called after every instruction execution
func (m *MBC3) ClockChanged(cycles int) error {

	m.cyclesCounter += cycles

	if m.cyclesCounter >= cpu.Frequency {

		m.cyclesCounter = m.cyclesCounter - cpu.Frequency

		if !m.halt() {

			// seconds
			m.setSeconds(m.seconds() + 1)

			if m.seconds() == 60 {

				m.setSeconds(0)

				// minutes
				m.setMinutes(m.minutes() + 1)

				if m.minutes() == 60 {

					m.setMinutes(0)

					// hours
					m.setHours(m.hours() + 1)

					if m.hours() == 24 {

						m.setHours(0)

						// days
						m.setDays(m.days() + 1)

						if m.days() == 0 {
							m.setDayCarry()
						}
					}
				}
			}
		}
	}

	return nil
}

// Read from address 'addr' at the target bank or rtc registers
func (m *MBC3) Read(addr uint16) (byte, error) {

	// rom bank 0
	if 0x0000 <= addr && addr <= 0x3FFF {

		if err := m.rom.SetWindow(0); err != nil {
			return 0, err
		}

		return m.rom.Read(addr)
	}

	// other rom banks
	if 0x4000 <= addr && addr <= 0x7FFF {

		if err := m.otherBanks.SetWindow(m.bankROM * 16384); err != nil {
			return 0, err
		}

		return m.otherBanks.Read(addr)
	}

	// ram / rtc
	if 0xA000 <= addr && addr <= 0xBFFF {

		// ram
		if 0x00 <= m.rtcCode && m.rtcCode <= 0x03 {

			if err := m.ram.SetWindow(m.bankRAM * 8192); err != nil {
				return 0, err
			}

			return m.ram.Read(addr)
		}

		// rtc
		if 0x08 <= m.rtcCode && m.rtcCode <= 0x0C {

			if m.rtcSnapshot != nil {
				return m.rtcSnapshot[m.rtcCode-0x08], nil
			}

			return m.rtc[m.rtcCode-0x08], nil
		}

		return 0, nil
	}

	return 0, memory.ReadOutOfRangeError(addr)
}

// Write 'data' to address 'addr' at the target bank
// or change the MBC control registers
func (m *MBC3) Write(addr uint16, data byte) error {

	// enable ram
	if 0x0000 <= addr && addr <= 0x1FFF {

		if data == 0x00 {
			m.enableTimerAndRAM = false
		} else if data == 0x0A {
			m.enableTimerAndRAM = true
		}

		return nil
	}

	// rom bank
	if 0x2000 <= addr && addr <= 0x3FFF {

		if data == 0 {
			data = 1
		}

		m.bankROM = uint32(data & 0x7F)

		return nil
	}

	// ram
	if 0x4000 <= addr && addr <= 0x5FFF {

		m.rtcCode = data

		if 0x00 <= data && data <= 0x03 {

			m.bankRAM = uint32(data)
		}

		return nil
	}

	// latch clock data
	if 0x6000 <= addr && addr <= 0x7FFF {

		if data == 0x00 {

			m.latch = true
			return nil
		}

		if data == 0x01 {

			if m.latch {

				if m.rtcSnapshot == nil {

					m.rtcSnapshot = make([]byte, 5)
					copy(m.rtcSnapshot, m.rtc[:])

				} else {

					m.rtcSnapshot = nil
				}

			}
		}

		m.latch = false

		return nil
	}

	// ram / rtc
	if 0xA000 <= addr && addr <= 0xBFFF {

		if !m.enableTimerAndRAM {
			return nil
		}

		// ram
		if 0x00 <= m.rtcCode && m.rtcCode <= 0x03 {

			if err := m.ram.SetWindow(m.bankRAM * 8192); err != nil {
				return err
			}

			return m.ram.Write(addr, data)
		}

		// rtc
		if 0x08 <= m.rtcCode && m.rtcCode <= 0x0C {

			m.rtc[m.rtcCode-0x08] = data
		}

		return nil
	}

	return memory.WriteOutOfRangeError(addr)
}
