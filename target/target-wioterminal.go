//go:build wioterminal

package target

func Pin(pin string) GpioPin {
	return supportedTargets["wioterminal"].GpioPins[pin]
}
