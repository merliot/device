//go:build nano_rp2040

package device

import "github.com/merliot/device/target"

func (d *Device) Pins() target.GpioPins {
	return d.Targets["nano-rp2040"].GpioPins
}
