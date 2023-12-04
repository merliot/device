//go:build prime

package runner

import (
	"github.com/merliot/dean"
	"github.com/merliot/prime"
)

func Run(cfg Config, thinger dean.Thinger) {
	prime := prime.New("p1", "prime", "p1")
	server := dean.NewServer(prime, cfg.User, cfg.Passwd, cfg.PortPrime)
	server.AdoptThing(thinger)
	server.Run()
}
