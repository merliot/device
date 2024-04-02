//go:build tinygo

package vl53l1x

import (
	"machine"

	dev "tinygo.org/x/drivers/vl53l1x"
)

type Vl53l1x struct {
	device dev.Device
}

// Configure VL53L1x time-of-flight distance sensor
func (v *Vl53l1x) Configure() {
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 400000,
	})
	v.device = dev.New(machine.I2C0)
	if v.device.Connected() {
		v.device.Configure(true)
		v.device.SetMeasurementTimingBudget(50000)
		v.device.StartContinuous(50)
	}
}

func (v *Vl53l1x) Distance() (dist int32, ok bool) {
	if v.device.Connected() {
		v.device.Read(true)
		dist = v.device.Distance() // mm
		ok = true
	}
	return
}
