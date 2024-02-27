//go:build tinygo

package runner

import (
	"github.com/merliot/dean"
	_ "github.com/merliot/dean/tinynet"
	_ "github.com/merliot/dean/tinynet/connect"
	"github.com/merliot/device"
)

func Run(device *device.Device, port, portPrime, user, passwd, dialURLs, wsScheme string) {
	runner := dean.NewRunner(device, user, passwd)
	runner.Dial(dialURLs)
	runner.Run()
}
