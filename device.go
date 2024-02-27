package device

import (
	"embed"
	"html"
	"net/url"
	"strings"

	"github.com/merliot/dean"
	"github.com/merliot/device/target"
)

// key: ssid; value: passphrase
type WifiAuth map[string]string

type Devicer interface {
	Load(string) error
	CopyWifiAuth(WifiAuth)
	SetWsScheme(string)
}

type Modeler interface {
	Icon() []byte
	DescHtml() []byte
	SupportedTargets() string
	GenerateUf2s() error
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

type ViewMode int

const (
	ViewFull ViewMode = iota
	ViewTile
	ViewTileVert
	ViewTileHorz
)

type Device struct {
	dean.Thing
	target.Targets `json:"-"`
	ModelStruct    string `json:"-"`
	WifiAuth       `json:"-"`
	DialURLs       string `json:"-"`
	DeployParams   string
	deployValues   url.Values
	ViewMode       `json:"-"`
	WsScheme       string `json:"-"`
	fs             embed.FS
	deviceOS
}

func New(id, model, name string, fs embed.FS, targets []string) dean.Thinger {
	println("NEW DEVICE")
	d := &Device{
		Thing:       dean.NewThing(id, model, name),
		Targets:     target.MakeTargets(targets),
		ModelStruct: strings.Title(model),
		WifiAuth:    make(WifiAuth),
		WsScheme:    "ws://",
		fs:          fs,
	}
	d.deviceOSInit()
	return d
}

func (d *Device) ParamFirstValue(key string) string {
	if v, ok := d.deployValues[key]; ok {
		return v[0]
	}
	return ""
}

func (d *Device) SetDeployParams(params string) {
	d.DeployParams = html.UnescapeString(params)
	d.deployValues, _ = url.ParseQuery(d.DeployParams)
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

func (d *Device) CopyWifiAuth(auth WifiAuth) {
	d.WifiAuth = auth
}

func (d *Device) SetDialURLs(urls string) {
	d.DialURLs = urls
}

func (d *Device) SetWsScheme(scheme string) {
	d.WsScheme = scheme
}
