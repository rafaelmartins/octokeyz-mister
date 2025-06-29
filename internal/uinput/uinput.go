package uinput

import (
	"encoding/binary"
	"os"
	"slices"
	"syscall"
	"time"
)

const (
	_UI_DEV_CREATE  = uintptr(0x5501)
	_UI_DEV_DESTROY = uintptr(0x5502)
	_UI_SET_EVBIT   = uintptr(0x40045564)
	_UI_SET_KEYBIT  = uintptr(0x40045565)

	_EV_SYN = 0x00
	_EV_KEY = 0x01
)

type Key uint16

const (
	KEY_ESC              Key = 1
	KEY_1                Key = 2
	KEY_2                Key = 3
	KEY_3                Key = 4
	KEY_4                Key = 5
	KEY_5                Key = 6
	KEY_6                Key = 7
	KEY_7                Key = 8
	KEY_8                Key = 9
	KEY_9                Key = 10
	KEY_0                Key = 11
	KEY_MINUS            Key = 12
	KEY_EQUAL            Key = 13
	KEY_BACKSPACE        Key = 14
	KEY_TAB              Key = 15
	KEY_Q                Key = 16
	KEY_W                Key = 17
	KEY_E                Key = 18
	KEY_R                Key = 19
	KEY_T                Key = 20
	KEY_Y                Key = 21
	KEY_U                Key = 22
	KEY_I                Key = 23
	KEY_O                Key = 24
	KEY_P                Key = 25
	KEY_LEFTBRACE        Key = 26
	KEY_RIGHTBRACE       Key = 27
	KEY_ENTER            Key = 28
	KEY_LEFTCTRL         Key = 29
	KEY_A                Key = 30
	KEY_S                Key = 31
	KEY_D                Key = 32
	KEY_F                Key = 33
	KEY_G                Key = 34
	KEY_H                Key = 35
	KEY_J                Key = 36
	KEY_K                Key = 37
	KEY_L                Key = 38
	KEY_SEMICOLON        Key = 39
	KEY_APOSTROPHE       Key = 40
	KEY_GRAVE            Key = 41
	KEY_LEFTSHIFT        Key = 42
	KEY_BACKSLASH        Key = 43
	KEY_Z                Key = 44
	KEY_X                Key = 45
	KEY_C                Key = 46
	KEY_V                Key = 47
	KEY_B                Key = 48
	KEY_N                Key = 49
	KEY_M                Key = 50
	KEY_COMMA            Key = 51
	KEY_DOT              Key = 52
	KEY_SLASH            Key = 53
	KEY_RIGHTSHIFT       Key = 54
	KEY_KPASTERISK       Key = 55
	KEY_LEFTALT          Key = 56
	KEY_SPACE            Key = 57
	KEY_CAPSLOCK         Key = 58
	KEY_F1               Key = 59
	KEY_F2               Key = 60
	KEY_F3               Key = 61
	KEY_F4               Key = 62
	KEY_F5               Key = 63
	KEY_F6               Key = 64
	KEY_F7               Key = 65
	KEY_F8               Key = 66
	KEY_F9               Key = 67
	KEY_F10              Key = 68
	KEY_NUMLOCK          Key = 69
	KEY_SCROLLLOCK       Key = 70
	KEY_KP7              Key = 71
	KEY_KP8              Key = 72
	KEY_KP9              Key = 73
	KEY_KPMINUS          Key = 74
	KEY_KP4              Key = 75
	KEY_KP5              Key = 76
	KEY_KP6              Key = 77
	KEY_KPPLUS           Key = 78
	KEY_KP1              Key = 79
	KEY_KP2              Key = 80
	KEY_KP3              Key = 81
	KEY_KP0              Key = 82
	KEY_KPDOT            Key = 83
	KEY_102ND            Key = 86
	KEY_F11              Key = 87
	KEY_F12              Key = 88
	KEY_RO               Key = 89
	KEY_KATAKANA         Key = 90
	KEY_HENKAN           Key = 92
	KEY_KATAKANAHIRAGANA Key = 93
	KEY_MUHENKAN         Key = 94
	KEY_KPENTER          Key = 96
	KEY_RIGHTCTRL        Key = 97
	KEY_KPSLASH          Key = 98
	KEY_SYSRQ            Key = 99
	KEY_RIGHTALT         Key = 100
	KEY_HOME             Key = 102
	KEY_UP               Key = 103
	KEY_PAGEUP           Key = 104
	KEY_LEFT             Key = 105
	KEY_RIGHT            Key = 106
	KEY_END              Key = 107
	KEY_DOWN             Key = 108
	KEY_PAGEDOWN         Key = 109
	KEY_INSERT           Key = 110
	KEY_DELETE           Key = 111
	KEY_POWER            Key = 116
	KEY_KPEQUAL          Key = 117
	KEY_PAUSE            Key = 119
	KEY_HANGUEL          Key = 122
	KEY_HANJA            Key = 123
	KEY_YEN              Key = 124
	KEY_LEFTMETA         Key = 125
	KEY_RIGHTMETA        Key = 126
	KEY_COMPOSE          Key = 127
	KEY_F13              Key = 183
	KEY_F14              Key = 184
	KEY_F15              Key = 185
	KEY_F16              Key = 186
	KEY_F17              Key = 187
	KEY_F18              Key = 188
	KEY_F19              Key = 189
	KEY_F20              Key = 190
	KEY_F21              Key = 191
	KEY_F22              Key = 192
	KEY_F23              Key = 193
	KEY_F24              Key = 194
)

type inputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

type uinputUserDev struct {
	Name [80]byte
	ID   struct {
		Bustype uint16
		Vendor  uint16
		Product uint16
		Version uint16
	}
	FFEffectsMax uint32
	Absmax       [64]int32
	Absmin       [64]int32
	Absfuzz      [64]int32
	Absflat      [64]int32
}

type Device struct {
	fp *os.File
}

func NewDevice(keys []Key) (*Device, error) {
	fp, err := os.OpenFile("/dev/uinput", syscall.O_WRONLY|syscall.O_NONBLOCK, 0660)
	if err != nil {
		return nil, err
	}

	rv := &Device{
		fp: fp,
	}

	uiud := uinputUserDev{}
	copy(uiud.Name[:], "octokeyz")
	uiud.ID.Bustype = 0x03 // USB
	uiud.ID.Version = 4

	if err := rv.ioctl(_UI_SET_EVBIT, _EV_KEY); err != nil {
		return nil, err
	}

	for _, key := range keys {
		if err := rv.ioctl(_UI_SET_KEYBIT, uintptr(key)); err != nil {
			return nil, err
		}
	}

	if err := binary.Write(rv.fp, binary.LittleEndian, uiud); err != nil {
		return nil, err
	}

	if err := rv.ioctl(_UI_DEV_CREATE, 0); err != nil {
		return nil, err
	}

	time.Sleep(1 * time.Second)
	return rv, nil
}

func (d *Device) ioctl(request uintptr, arg uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, d.fp.Fd(), request, arg)
	if errno != 0 {
		return errno
	}
	return nil
}

func (d *Device) Close() error {
	if err := d.ioctl(_UI_DEV_DESTROY, 0); err != nil {
		return err
	}
	return d.fp.Close()
}

func (d *Device) sendEvents(kc []Key, value int32) error {
	sl := slices.Clone(kc)
	if value == 0 {
		slices.Reverse(sl)
	}

	for _, k := range sl {
		ts := getNowTimeval()
		if err := binary.Write(d.fp, binary.LittleEndian, inputEvent{
			Time:  ts,
			Type:  _EV_KEY,
			Code:  uint16(k),
			Value: value,
		}); err != nil {
			return err
		}

		if err := binary.Write(d.fp, binary.LittleEndian, inputEvent{
			Time:  ts,
			Type:  _EV_SYN,
			Code:  0,
			Value: 0,
		}); err != nil {
			return err
		}

		time.Sleep(10 * time.Millisecond)
	}
	return nil
}

func (d *Device) Press(kc ...Key) error {
	return d.sendEvents(kc, 1)
}

func (d *Device) Release(kc ...Key) error {
	return d.sendEvents(kc, 0)
}
