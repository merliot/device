//go:build nano_rp2040

package target

func Pin(pin string) (GpioPin, bool) {
	return Targets["nano-rp2040"].GpioPins[pin]
}
