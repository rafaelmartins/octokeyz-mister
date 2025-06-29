package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/rafaelmartins/octokeyz-mister/internal/cleanup"
	"github.com/rafaelmartins/octokeyz-mister/internal/inotify"
	"github.com/rafaelmartins/octokeyz-mister/internal/ipaddr"
	"github.com/rafaelmartins/octokeyz-mister/internal/uinput"
	"github.com/rafaelmartins/octokeyz-mister/internal/vkbd"
	"rafaelmartins.com/p/octokeyz"
)

func updateCore(dev *octokeyz.Device) error {
	v, err := os.ReadFile("/tmp/CORENAME")
	if err != nil {
		return err
	}
	return dev.DisplayLine(octokeyz.DisplayLine4, fmt.Sprintf("Core: %s", strings.TrimSpace(string(v))), octokeyz.DisplayLineAlignLeft)
}

func main() {
	defer cleanup.Cleanup()

	dev, err := octokeyz.GetDevice("")
	cleanup.Check(err)

	cleanup.Check(dev.Open())
	cleanup.Register(dev)

	kbd, err := vkbd.New(dev, octokeyz.BUTTON_5, map[octokeyz.ButtonID]vkbd.ButtonMapping{
		octokeyz.BUTTON_1: {
			Normal: []uinput.Key{uinput.KEY_F12},
			Mod:    []uinput.Key{uinput.KEY_LEFTALT, uinput.KEY_F12},
		},
		octokeyz.BUTTON_2: {
			Normal: []uinput.Key{uinput.KEY_ESC},
		},
		octokeyz.BUTTON_3: {
			Normal: []uinput.Key{uinput.KEY_UP},
		},
		octokeyz.BUTTON_4: {
			Normal: []uinput.Key{uinput.KEY_ENTER},
		},
		octokeyz.BUTTON_6: {
			Normal: []uinput.Key{uinput.KEY_LEFT},
			Mod:    []uinput.Key{uinput.KEY_LEFTCTRL, uinput.KEY_LEFTALT, uinput.KEY_RIGHTALT},
		},
		octokeyz.BUTTON_7: {
			Normal: []uinput.Key{uinput.KEY_DOWN},
			Mod:    []uinput.Key{uinput.KEY_LEFTSHIFT, uinput.KEY_LEFTCTRL, uinput.KEY_LEFTALT, uinput.KEY_RIGHTALT},
		},
		octokeyz.BUTTON_8: {
			Normal: []uinput.Key{uinput.KEY_RIGHT},
		},
	})
	cleanup.Check(err)
	cleanup.Register(kbd)

	dev.AddHandler(octokeyz.BUTTON_5, func(b *octokeyz.Button) error {
		return dev.Led(octokeyz.LedFlash)
	})

	cleanup.Check(dev.DisplayLine(octokeyz.DisplayLine1, "MiSTer FPGA", octokeyz.DisplayLineAlignCenter))

	in, err := inotify.New()
	cleanup.Check(err)
	cleanup.Register(in)
	cleanup.Check(in.AddWatch("/tmp/CORENAME", inotify.IN_CLOSE_WRITE))
	cleanup.Check(updateCore(dev))

	go func() {
		cleanup.Check(in.Listen(func(f string, ev inotify.InotifyEvent) error {
			return updateCore(dev)
		}))
	}()

	ip, err := ipaddr.NewMonitor("eth0", "wlan0")
	cleanup.Check(err)
	cleanup.Register(ip)

	go func() {
		cleanup.Check(ip.Run(func(itf string, ip net.IP) error {
			line := octokeyz.DisplayLine6
			if itf == "eth0" {
				line = octokeyz.DisplayLine7
			}
			if ip == nil {
				return dev.DisplayClearLine(line)
			}
			return dev.DisplayLine(line, fmt.Sprintf("%s: %s", strings.ToUpper(itf[:len(itf)-1]), ip), octokeyz.DisplayLineAlignLeft)
		}))
	}()

	cleanup.Check(dev.Listen(nil))
}
