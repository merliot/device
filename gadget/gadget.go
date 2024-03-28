package gadget

import (
	"embed"
	"net/http"

	"github.com/merliot/dean"
	"github.com/merliot/device"
)

//go:embed template
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

func (g *Gadget) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	g.API(w, req, g)
}

func (g *Gadget) Run(i *dean.Injector) {
	select {
	case <-g.quit:
	}
}
