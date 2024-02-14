//go:build wioterminal

package target

func Pin(pin string) GpioPin {
	return Targets["wioterminal"].GpioPins[pin]
}
