package device

import (
	"html"
	"net/url"

	"github.com/merliot/dean"
)

// key: ssid; value: passphrase
type WifiAuth map[string]string

type Devicer interface {
	Load()
	SetWifiAuth(WifiAuth)
}

type Device struct {
	dean.Thing
	Targets      `json:"-"`
	WifiAuth     `json:"-"`
	DeployParams string
	deviceOS
}

func New(id, model, name string, targets []string) dean.Thinger {
	println("NEW COMMON")
	d := &Device{}
	d.Thing = dean.NewThing(id, model, name)
	d.Targets = makeTargets(targets)
	d.WifiAuth = make(WifiAuth)
	d.deviceOSInit()
	return c
}

func (d *Device) ParseDeployParams() url.Values {
	values, _ := url.ParseQuery(d.DeployParams)
	return values
}

func (d *Device) SetDeployParams(params string) {
       d.DeployParams = html.UnescapeString(params)
}

func (d *Device) SetWifiAuth(auth WifiAuth) {
	d.WifiAuth = auth
}
