//go:build prime

package runner

import (
	"github.com/merliot/dean"
	"github.com/merliot/device/prime"
)

func Run(thinger dean.Thinger, port, portPrime, user, passwd, dialURLs, wsScheme string) {
	prime := prime.New("p1", "prime", "p1").(*prime.Prime)
	prime.SetWsScheme(wsScheme)
	server := dean.NewServer(prime, user, passwd, portPrime)
	server.AdoptThing(thinger)
	server.Run()
}
