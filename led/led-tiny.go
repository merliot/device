//go:build tinygo

package led

import "machine"

type Led struct {
	State bool
	pin   machine.Pin
}

func (l *Led) Configure() {
	l.pin = machine.LED
	l.pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	l.pin.Low()
}

func (l *Led) Set(state bool) {
	if state {
		l.On()
	} else {
		l.Off()
	}
}

func (l *Led) On() {
	l.pin.High()
	l.State = true
}

func (l *Led) Off() {
	l.pin.Low()
	l.State = false
}
