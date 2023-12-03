//go:build tinygo

package runner

import (
	"github.com/merliot/dean"
	_ "github.com/merliot/dean/tinynet"
)

func Run(thinger dean.Thinger) {
	runner := dean.NewRunner(thinger)
	runner.Dial()
	runner.Run()
}
