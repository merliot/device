//go:build !tinygo && !prime

package runner

import (
	"github.com/merliot/dean"
)

func Run(thinger dean.Thinger) {
	server := dean.NewServer(thinger)
	server.Dial()
	server.Run()
}
