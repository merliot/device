//go:build !rpi && !tinygo

package relay

type Relay struct {
	Name  string
	Gpio  string
	State bool
}

func (r Relay) Configure() {
}

func (r Relay) On() {
}

func (r Relay) Off() {
}
