//go:build rpi

package device

import "github.com/merliot/device/target"

func (d *Device) Pins() target.GpioPins {
	return d.Targets["rpi"].GpioPins
}
