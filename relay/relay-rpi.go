//go:build rpi

package relay

import (
	"strconv"

	"github.com/merliot/device/target"
	"gobot.io/x/gobot/v2/drivers/gpio"
)

type Relay struct {
	Name   string
	Gpio   string
	State  bool
	driver *gpio.RelayDriver
}

func (r Relay) Configure() {
	if pin, ok := target.Pin(r.Gpio); ok {
		spin := strconv.Itoa(int(pin))
		r.driver = gpio.NewRelayDriver(target.Adaptor, spin)
		r.driver.Start()
		r.driver.Off()
	}
}

func (r Relay) On() {
	if r.driver != nil {
		r.driver.On()
	}
}

func (r Relay) Off() {
	if r.driver != nil {
		r.driver.Off()
	}
}
