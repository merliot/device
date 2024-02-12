package skeleton

import (
	"embed"
	"net/http"

	"github.com/merliot/dean"
	"github.com/merliot/device"
	"github.com/merliot/device/led"
)

//go:embed css images js template
var fs embed.FS

var targets = []string{"demo"}

type Skeleton struct {
	*device.Device
	led.Led
}

type MsgClick struct {
	dean.ThingMsg
	State bool
}

func New(id, model, name string) dean.Thinger {
	println("NEW SKELETON")
	return &Skeleton{
		Device: device.New(id, model, name, fs, targets).(*device.Device),
		Led:    led.New(),
	}
}

func (s *Skeleton) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.API(w, r, s)
}

func (s *Skeleton) save(msg *dean.Msg) {
	msg.Unmarshal(s).Broadcast()
}

func (s *Skeleton) getState(msg *dean.Msg) {
	s.Path = "state"
	msg.Marshal(s).Reply()
}

func (s *Skeleton) click(msg *dean.Msg) {
	msg.Unmarshal(&s.Led)
	if s.IsMetal() {
		s.Led.Set(s.State)
	}
	msg.Broadcast()
}

func (s *Skeleton) Subscribers() dean.Subscribers {
	return dean.Subscribers{
		"state":     s.save,
		"get/state": s.getState,
		"click":     s.click,
	}
}

func (s *Skeleton) Run(i *dean.Injector) {
	select {}
}
