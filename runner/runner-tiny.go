//go:build tinygo

package runner

import (
	"github.com/merliot/dean"
)

func (d *Device) Serve(thinger dean.Thinger) {
	tinynet.NetConnect(ssid, pass)
	runner := dean.NewRunner(thinger)
	runner.Dial()
	runner.Run()
}
