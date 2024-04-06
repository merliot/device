//go:build tinygo

package dht

import (
	"errors"
	"machine"

	"github.com/merliot/device/target"
	"tinygo.org/x/drivers/dht"
)

type Dht struct {
	Sensor      string
	Gpio        string
	Temperature float32
	Humidity    float32
	pin         machine.Pin
	dht         dht.Device
}

func (d *Dht) Configure() {
	sensor := dht.DHT11
	switch d.Sensor {
	case "DHT22":
		sensor = dht.DHT22
	}
	d.pin = machine.NoPin
	if pin, ok := target.Pin(d.Gpio); ok {
		d.pin = machine.Pin(pin)
		d.dht = dht.New(d.pin, sensor)
	}
}

func (d *Dht) Read() error {
	if d.pin == machine.NoPin {
		return errors.New("Gpio pin not configured")
	}
	temp, hum, err := d.dht.Measurements()
	if err != nil {
		return err
	}
	d.Temperature = float32(temp) / 10.0
	d.Humidity = float32(hum) / 10.0
	return nil
}
