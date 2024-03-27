package gadget

import (
	"embed"

	"github.com/merliot/dean"
	"github.com/merliot/device"
)

var fs embed.FS

type Gadget struct {
	*device.Device
	quit chan (bool)
}

var targets = []string{"demo"}

func New(id, model, name string) dean.Thinger {
	return &Gadget{
		Device: device.New(id, model, name, fs, targets).(*device.Device),
		quit:   make(chan (bool)),
	}
}

func (g *Gadget) Run(i *dean.Injector) {
	select {
	case <-g.quit:
	}
}
