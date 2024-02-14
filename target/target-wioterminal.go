//go:build wioterminal

package target

func Pin(pin string) GpioPin {
	gpio, ok := supportedTargets["wioterminal"].GpioPins[pin]
	return gpio, ok
}
