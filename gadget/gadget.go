package gadget

import (
	"embed"
	"fmt"
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
	fmt.Println("NEW GADGET\r")
	return &Gadget{
		Device: device.New(id, model, name, fs, targets).(*device.Device),
		quit:   make(chan bool),
	}
}

type MsgBottles struct {
	TakeOneDown int32
}

func (g *Gadget) save(pkt *dean.Packet) {
	println("GADGET SAVE")
	pkt.Unmarshal(g).Broadcast()
	pkt.SetPath("bottles").Marshal(&MsgBottles{99}).Reply()
}

func (g *Gadget) getState(pkt *dean.Packet) {
	pkt.SetPath("state").Marshal(g).Reply()
}

func (g *Gadget) bottles(pkt *dean.Packet) {
	var bottles MsgBottles
	pkt.Unmarshal(&bottles)
	bottles.TakeOneDown--
	if bottles.TakeOneDown > 0 {
		pkt.Marshal(&bottles).Reply()
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
