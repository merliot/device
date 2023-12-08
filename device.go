package device

import (
	"html"
	"net/url"
	"strings"

	"github.com/merliot/dean"
	"github.com/merliot/target"
)

// key: ssid; value: passphrase
type WifiAuth map[string]string

type Devicer interface {
	Load()
	SetWifiAuth(WifiAuth)
}

type Device struct {
	dean.Thing
	target.Targets `json:"-"`
	WifiAuth       `json:"-"`
	DeployParams   string
	deviceOS
}

func New(id, model, name string, targets []string) dean.Thinger {
	println("NEW COMMON")
	d := &Device{}
	d.Thing = dean.NewThing(id, model, name)
	d.Targets = target.MakeTargets(targets)
	d.WifiAuth = make(WifiAuth)
	d.deviceOSInit()
	return d
}

func (d *Device) ParseDeployParams() url.Values {
	values, _ := url.ParseQuery(d.DeployParams)
	return values
}

func (d *Device) SetDeployParams(params string) {
	d.DeployParams = html.UnescapeString(params)
}

func (d *Device) SetWifiAuth(ssids, passphrases string) {
	d.WifiAuth = make(WifiAuth)
	keys := strings.Split(ssids, ",")
	values := strings.Split(passphrases, ",")
	for i, key := range keys {
		if i < len(values) {
			d.WifiAuth[key] = values[i]
		}
	}
}
