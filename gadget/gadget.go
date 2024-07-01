package gadget

import (
	"embed"
	"fmt"

	"github.com/merliot/dean"
	"github.com/merliot/device"
)

//go:embed css js template
var fs embed.FS

type Gadget struct {
	*device.Device
	Bottles int
}

type Update struct {
	Bottles int
}

var targets = []string{"demo"}

func New(id, model, name string) dean.Thinger {
	fmt.Println("NEW GADGET\r")
	var g = &Gadget{
		Device:  device.New(id, model, name, fs, targets).(*device.Device),
		Bottles: 4,
	}
	g.SetData(g)
	return g
}

func (g *Gadget) getState(pkt *dean.Packet) {
	pkt.SetPath("state").Marshal(g).Reply()
}

func (g *Gadget) save(pkt *dean.Packet) {
	println("gadget save")
	pkt.Unmarshal(g).Broadcast()
}

func (g *Gadget) takeone(pkt *dean.Packet) {
	if g.Bottles > 0 {
		g.Bottles--
		pkt.SetPath("tookone").Marshal(&Update{g.Bottles}).Reply().Broadcast()
	}
}

func (g *Gadget) Subscribers() dean.Subscribers {
	return dean.Subscribers{
		"get/state": g.getState,
		"state":     g.save,
		"takeone":   g.takeone,
	}
}

/*
func (g *Gadget) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.API(w, r, g)
}
*/
