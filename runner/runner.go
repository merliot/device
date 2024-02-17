//go:build !tinygo && !prime

package runner

import (
	"github.com/merliot/dean"
)

func Run(thinger dean.Thinger, port, portPrime, user, passwd, dialURLs, wsScheme string) {
	server := dean.NewServer(thinger, user, passwd, port)
	server.Dial(dialURLs)
	server.Run()
}
