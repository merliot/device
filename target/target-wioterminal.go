//go:build wioterminal

package target

func Pin(pin string) (GpioPin, bool) {
	gpio, ok := supportedTargets["wioterminal"].GpioPins[pin]
	return gpio, ok
}
