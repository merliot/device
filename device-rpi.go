//go:build rpi

package device

import (
	"strconv"

	"github.com/merliot/device/target"
	"gobot.io/x/gobot/v2/drivers/gpio"
)

// FailSafe by turning off all gpios
func (d *Device) FailSafe() {
	for _, pin := range target.Targets["rpi"].GpioPins {
		rpin := strconv.Itoa(pin)
		driver = gpio.NewDirectPinDriver(target.Adaptor, rpin)
		driver.start()
		driver.off()
	}

}

func (d *Device) Setup() {
	target.Adaptor.Connect()
}
