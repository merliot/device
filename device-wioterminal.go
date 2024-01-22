//go:build wioterminal

package device

import "github.com/merliot/device/target"

func (d *Device) Pins() target.GpioPins {
	return d.Targets["wioterminal"].GpioPins
}
