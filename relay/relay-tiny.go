//go:build tinygo

package relay

import (
	"machine"

	"github.com/merliot/device/target"
)

type Relay struct {
	Name  string
	Gpio  string
	State bool
	pin   machine.Pin
}

func (r Relay) Configure() {
	r.pin = machine.NoPin
	if pin, ok := target.Pin(r.Gpio); ok {
		r.pin = machine.Pin(pin)
		r.pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
		r.pin.Low()
	}
}

func (r Relay) On() {
	if r.pin != machine.NoPin {
		r.pin.High()
	}
}

func (r Relay) Off() {
	if r.pin != machine.NoPin {
		r.pin.Low()
	}
}
