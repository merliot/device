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
	Bottles int
}

var targets = []string{"demo"}

func New(id, model, name string) dean.Thinger {
	fmt.Println("NEW GADGET\r")
	return &Gadget{
		Device:  device.New(id, model, name, fs, targets).(*device.Device),
		Bottles: 99,
	}
}

type MsgBottles struct {
	TakeOneDown int32
}

func (g *Gadget) getState(pkt *dean.Packet) {
	pkt.SetPath("state").Marshal(g).Reply()
}

func (g *Gadget) save(pkt *dean.Packet) {
	pkt.Unmarshal(g).Broadcast()
}

func (g *Gadget) takeone(pkt *dean.Packet) {
	if g.Bottles > 0 {
		g.Bottles--
		pkt.SetPath("tookone").Broadcast()
	}
}

func (g *Gadget) Subscribers() dean.Subscribers {
	return dean.Subscribers{
		"get/state": g.getState,
		"state":     g.save,
		"takeone":   g.takeone,
	}
}

func (g *Gadget) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	g.API(w, req, g)
}
