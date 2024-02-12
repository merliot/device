//go:build tinygo

package led

import "machine"

type Led struct {
	State bool
	pin   machine.Pin
}

func New() Led {
	led := Led{pin: machine.LED}
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	return led
}
