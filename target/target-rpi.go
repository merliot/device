//go:build rpi

package target

import (
	"gobot.io/x/gobot/v2/platforms/raspi"
)

var Adapter *raspi.Adapter = raspi.NewAdapter()

func Pin(pin string) (GpioPin, bool) {
	gpio, ok := supportedTargets["rpi"].GpioPins[pin]
	return gpio, ok
}
