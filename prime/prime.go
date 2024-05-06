package prime

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/merliot/dean"
	"github.com/merliot/device"
)

//go:embed css images js template
var fs embed.FS

type Child struct {
	Id     string
	Model  string
	Name   string
	Online bool
}

type Prime struct {
	*device.Device
	Child Child
	quit  chan bool
}

var targets = []string{"x86-64", "rpi"}

func New(id, model, name string) dean.Thinger {
	fmt.Println("NEW PRIME\r")
	return &Prime{
		Device: device.New(id, model, name, fs, targets).(*device.Device),
		quit:   make(chan bool),
	}
}

func (p *Prime) getState(pkt *dean.Packet) {
	p.Path = "state"
	pkt.Marshal(p).Reply()
}

func (p *Prime) online(pkt *dean.Packet, online bool) {
	p.Child.Online = online
	pkt.Broadcast()
}

func (p *Prime) connect(online bool) func(*dean.Packet) {
	return func(pkt *dean.Packet) {
		p.online(pkt, online)
	}
}

func (p *Prime) adoptedChild(pkt *dean.Packet) {
	var adopt dean.ThingMsgAdopted
	pkt.Unmarshal(&adopt)
	p.Child = Child{Id: adopt.Id, Model: adopt.Model, Name: adopt.Name}
}

func (p *Prime) Subscribers() dean.Subscribers {
	return dean.Subscribers{
		"get/state":     p.getState,
		"connected":     p.connect(true),
		"disconnected":  p.connect(false),
		"adopted/thing": p.adoptedChild,
	}
}

func (p *Prime) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.API(w, r, p)
}

func (p *Prime) Run(i *dean.Injector) {
	select {
	case <-p.quit:
	}
}

func (p *Prime) Close() {
	p.quit <- true
}
