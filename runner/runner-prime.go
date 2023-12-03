//go:build prime

package runner

import (
	"github.com/merliot/dean"
	"github.com/merliot/prime"
)

func Run(thinger dean.Thinger) {
	prime := prime.New("p1", "prime", "p1")
	server := dean.NewServer(prime)
	server.AdoptThing(thinger)
	server.Run()
}
