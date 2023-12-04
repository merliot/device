//go:build tinygo

package runner

import (
	"github.com/merliot/dean"
	_ "github.com/merliot/dean/tinynet"
)

func Run(cfg Config, thinger dean.Thinger) {
	runner := dean.NewRunner(thinger, cfg.User, cfg.Passwd)
	runner.Dial(cfg.DialURLs)
	runner.Run()
}