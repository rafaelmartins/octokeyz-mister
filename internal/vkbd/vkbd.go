package vkbd

import (
	"fmt"
	"slices"

	"github.com/rafaelmartins/octokeyz-mister/internal/uinput"
	"rafaelmartins.com/p/octokeyz"
)

type ButtonMapping struct {
	Normal []uinput.Key
	Mod    []uinput.Key
}

type VKBD struct {
	ui      *uinput.Device
	okz     *octokeyz.Device
	mod     octokeyz.Modifier
	mapping map[octokeyz.ButtonID]ButtonMapping
}

func New(dev *octokeyz.Device, modifier octokeyz.ButtonID, mapping map[octokeyz.ButtonID]ButtonMapping) (*VKBD, error) {
	keys := []uinput.Key{}
	for btn, mapp := range mapping {
		if btn == modifier {
			return nil, fmt.Errorf("vkbd: %s used as modifier, can't be used a action button", btn)
		}
		for _, m := range mapp.Normal {
			if !slices.Contains(keys, m) {
				keys = append(keys, m)
			}
		}
		for _, m := range mapp.Mod {
			if !slices.Contains(keys, m) {
				keys = append(keys, m)
			}
		}
	}
	slices.Sort(keys)

	ui, err := uinput.NewDevice(keys)
	if err != nil {
		return nil, err
	}

	rv := &VKBD{
		ui:      ui,
		okz:     dev,
		mapping: mapping,
	}

	if err := dev.AddHandler(modifier, rv.mod.Handler); err != nil {
		return nil, err
	}
	for btn := range mapping {
		if err := dev.AddHandler(btn, rv.handler); err != nil {
			return nil, err
		}
	}
	return rv, nil
}

func (v *VKBD) Close() error {
	if v.ui != nil {
		return v.ui.Close()
	}
	return nil
}

func (v *VKBD) handler(b *octokeyz.Button) error {
	bm, ok := v.mapping[b.GetID()]
	if !ok {
		return fmt.Errorf("vkbd: handler called for button %s, but no mapping found", b)
	}

	c := bm.Normal
	if v.mod.Pressed() {
		c = bm.Mod
	}

	if err := v.ui.Press(c...); err != nil {
		return err
	}
	b.WaitForRelease()
	return v.ui.Release(c...)
}
