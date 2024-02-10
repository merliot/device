package skeleton

import (
	"embed"
	"net/http"

	"github.com/merliot/dean"
	"github.com/merliot/device"
)

//go:embed css images js template
var fs embed.FS

type Skeleton struct {
	*device.Device
}

var targets = []string{"demo"}

func New(id, model, name string) dean.Thinger {
	println("NEW SKELETON")
	s := &Skeleton{}
	s.Device = device.New(id, model, name, fs, targets).(*device.Device)
	return s
}

func (s *Skeleton) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.API(w, r, s)
}
