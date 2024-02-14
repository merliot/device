//go:build !rpi && !tinygo

package device

func (d *Device) Setup()    {}
func (d *Device) FailSafe() {}
