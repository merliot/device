//go:build wioterminal

package device

import "github.com/merliot/device/target"

func (d *Device) Pins() target.GpioPins {
	return r.Targets["wioterminal"].GpioPins
}
