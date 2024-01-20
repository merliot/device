package device

import (
	"html"
	"net/url"
	"strings"

	"github.com/merliot/dean"
	"github.com/merliot/device/target"
)

// key: ssid; value: passphrase
type WifiAuth map[string]string

type Devicer interface {
	Load()
	SetWifiAuth(WifiAuth)
}

type params struct {
	Ssid         string
	Passphrase   string
	Id           string
	Model        string
	Name         string
	DeployParams string
	User         string
	Passwd       string
	DialURLs     string
}

type Device struct {
	dean.Thing
	target.Targets `json:"-"`
	WifiAuth       `json:"-"`
	DeployParams   string
	DialURLs       string `json:"-"`
	deviceOS
}

func New(id, model, name string, targets []string) dean.Thinger {
	println("NEW DEVICE")
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

func (d *Device) SetDialURLs(urls string) {
	d.DialURLs = urls
}
