package prime

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/merliot/dean"
	"github.com/merliot/device"
)

//go:embed css images js template
var fs embed.FS

type Child struct {
	Id             string
	Model          string
	Name           string
	Online         bool
	device.Devicer `json:"-"`
}

func (c *Child) ModelTitle() template.JS {
	return template.JS(strings.Title(c.Model))
}

type Prime struct {
	*device.Device
	server *dean.Server
	Child  Child
	quit   chan bool
}

var targets = []string{"x86-64", "rpi"}

func New(id, model, name string) dean.Thinger {
	fmt.Println("NEW PRIME\r")
	return &Prime{
		Device: device.New(id, model, name, fs, targets).(*device.Device),
		quit:   make(chan bool),
	}
}

func NewPrime(id, model, name, portPrime, user, passwd string, thinger dean.Thinger) dean.Thinger {
	p := New(id, model, name).(*Prime)
	p.server = dean.NewServer(p, user, passwd, portPrime)
	p.Child.Devicer = thinger.(device.Devicer)
	p.server.AdoptThing(thinger)
	return p
}

func (p *Prime) Serve() {
	p.server.Run()
}

func (p *Prime) getState(pkt *dean.Packet) {
	pkt.SetPath("state").Marshal(p).Reply()
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
	p.Child.Id = adopt.Id
	p.Child.Model = adopt.Model
	p.Child.Name = adopt.Name
}

func (p *Prime) foo(pkt *dean.Packet) {
	println("FOO")
}

func (p *Prime) Subscribers() dean.Subscribers {
	return dean.Subscribers{
		"get/state":     p.getState,
		"connected":     p.connect(true),
		"disconnected":  p.connect(false),
		"adopted/thing": p.adoptedChild,
		"state":         p.foo,
	}
}

// ya, i know, generics
func reverseSlice(s []string) {
	for i := 0; i < len(s)/2; i++ {
		j := len(s) - i - 1
		s[i], s[j] = s[j], s[i]
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
	p.server.Close()
	p.quit <- true
}
