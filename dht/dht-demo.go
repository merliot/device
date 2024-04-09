//go:build !tinygo && !rpi

package dht

import (
	"math"
	"math/rand"
	"time"
)

type Dht struct {
	Sensor      string
	Gpio        string
	Temperature float32 // deg C
	Humidity    float32 // %
}

func (d *Dht) Configure() {
	rand.Seed(time.Now().UnixNano())
}

func randomValue(mean, stddev float64) float32 {
	u1 := rand.Float64()
	u2 := rand.Float64()

	// Box-Muller transform to generate normally distributed random numbers
	z0 := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2.0*math.Pi*u2)
	value := mean + stddev*z0

	// Round to 1 dec place
	return float32(math.Round(value*10) / 10)
}

func (d *Dht) Read() error {
	d.Temperature = randomValue(24.1, 0.05)
	d.Humidity = randomValue(34.5, 0.05)
	return nil
}
