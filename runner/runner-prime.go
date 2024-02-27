//go:build prime

package runner

import (
	"flag"
	"fmt"

	"github.com/merliot/dean"
	"github.com/merliot/device"
	"github.com/merliot/device/prime"
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

	prime := prime.New("p1", "prime", "p1").(*prime.Prime)
	prime.SetWsScheme(wsScheme)
	server := dean.NewServer(prime, user, passwd, portPrime)
	server.AdoptThing(device)
	server.Run()
}
