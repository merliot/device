//go:build tinygo

package runner

import (
	"github.com/merliot/dean"
	_ "github.com/merliot/dean/tinynet"
	_ "github.com/merliot/dean/tinynet/connect"
)

func Run(thinger dean.Thinger, port, portPrime, user, passwd, dialURLs, wsScheme string) {
	runner := dean.NewRunner(thinger, user, passwd)
	runner.Dials(dialURLs)
	runner.Run()
}
