//go:build rpi

package target

import (
	"gobot.io/x/gobot/v2/platforms/raspi"
)

var Adaptor *raspi.Adaptor = raspi.NewAdaptor()

func Pin(pin string) (GpioPin, bool) {
	gpio, ok := supportedTargets["rpi"].GpioPins[pin]
	return gpio, ok
}
