//go:build !tinygo && !prime

package runner

import (
	"flag"
	"fmt"

	"github.com/merliot/dean"
	"github.com/merliot/device"
)

func Run(device *device.Device, port, portPrime, user, passwd, dialURLs, wsScheme string) {
	uf2 := flag.Bool("uf2", false, "Generate uf2 files")
	flag.Parse()

	if *uf2 {
		if err := device.GenerateUf2s(); err != nil {
			fmt.Println(err.Error())
		}
		return
	}

	server := dean.NewServer(device, user, passwd, port)
	server.Dial(dialURLs)
	server.Run()
}
