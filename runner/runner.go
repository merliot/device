//go:build !tinygo && !prime

package runner

import (
	"github.com/merliot/dean"
)

func Run(cfg Config, thinger dean.Thinger) {
	server := dean.NewServer(thinger, cfg.User, cfg.Passwd, cfg.Port)
	server.Dial(cfg.DialURLs)
	server.Run()
}
