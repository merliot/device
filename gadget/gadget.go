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
	quit chan bool
}

var targets = []string{"demo"}

func New(id, model, name string) dean.Thinger {
	return &Gadget{
		Device: device.New(id, model, name, fs, targets).(*device.Device),
		quit:   make(chan bool),
	}
}

type Bottles struct {
	Path        string
	TakeOneDown int32
}

func (g *Gadget) save(msg *dean.Msg) {
	msg.Unmarshal(g).Broadcast()
	msg.Marshal(&Bottles{"bottles", 99}).Reply()
}

func (g *Gadget) getState(msg *dean.Msg) {
	g.Path = "state"
	msg.Marshal(g).Reply()
}

func (g *Gadget) bottles(msg *dean.Msg) {
	var bottles Bottles
	msg.Unmarshal(&bottles)
	bottles.TakeOneDown--
	if bottles.TakeOneDown > 0 {
		msg.Marshal(&bottles).Reply()
	} else {
		g.quit <- true
	}
}

func (g *Gadget) Subscribers() dean.Subscribers {
	return dean.Subscribers{
		"state":     g.save,
		"get/state": g.getState,
		"bottles":   g.bottles,
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
