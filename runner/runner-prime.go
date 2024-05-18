//go:build prime

package runner

import (
	"github.com/merliot/dean"
	"github.com/merliot/device/prime"
)

func Run(thinger dean.Thinger, port, portPrime, user, passwd, dialURLs, wsScheme string) {
	prime := prime.NewPrime("p1", "prime", "p1", portPrime, user, passwd, thinger).(*prime.Prime)
	prime.SetWsScheme(wsScheme)
	prime.Serve()
}
