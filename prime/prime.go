package prime

import (
	"embed"
	"html/template"
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
	Child     Child
	templates *template.Template
}

var targets = []string{"x86-64", "rpi"}

func New(id, model, name string) dean.Thinger {
	println("NEW PRIME")
	p := &Prime{}
	p.Device = device.New(id, model, name, targets).(*device.Device)
	p.CompositeFs.AddFS(fs)
	p.templates = p.CompositeFs.ParseFS("template/*")
	return p
}

func (p *Prime) getState(msg *dean.Msg) {
	p.Path = "state"
	msg.Marshal(p).Reply()
}

func (p *Prime) online(msg *dean.Msg, online bool) {
	p.Child.Online = online
	msg.Broadcast()
}

func (p *Prime) connect(online bool) func(*dean.Msg) {
	return func(msg *dean.Msg) {
		p.online(msg, online)
	}
}

func (p *Prime) adoptedChild(msg *dean.Msg) {
	var adopt dean.ThingMsgAdopted
	msg.Unmarshal(&adopt)
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
	p.API(p.templates, w, r)
}
