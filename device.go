// Package device is the Merliot base device.

package device

import (
	"embed"
	"fmt"
	"html"
	"html/template"
	"net/url"
	"strings"

	"github.com/merliot/dean"
	"github.com/merliot/device/target"
)

// WIfiAuth is a map of SSID:PASSPHRASE pairs
type WifiAuth map[string]string // key: ssid; value: passphrase

type Devicer interface {
	CopyWifiAuth(WifiAuth)
	SetWsScheme(string)
	SetDialURLs(string)
	SetLocked(bool)
	GetDeployParams() string
	SetDeployParams(string)
	RenderHTML(string) (template.HTML, error)
}

type params struct {
	MagicStart   string
	Ssid         string
	Passphrase   string
	Id           string
	Model        string
	Name         string
	DeployParams string
	User         string
	Passwd       string
	DialURLs     string
	MagicEnd     string
}

type Device struct {
	// Device is a Thing
	dean.Thing
	// Data passed to render templates...usually the device parent
	data any
	// Targets supported by device
	target.Targets `json:"-"`
	// WifiAuth is a map of SSID:PASSPHRASE pairs
	WifiAuth `json:"-"`
	// DialURLs is a comma seperated list of URLs the device will dial into
	DialURLs string `json:"-"`
	// DeployParams are device deploy configuration in an html param format
	DeployParams string
	deployValues url.Values
	// WsScheme is the websocket scheme to use to call back into the
	// device.  Default is ws://, which is sutable for an http:// device.
	// Set to wss:// if dialing back into an https:// device.
	WsScheme string `json:"-"`
	fs       embed.FS
	// Administratively locked
	Locked bool `json:"-"`
	// OS-specific members
	deviceOS
}

const defaultWsScheme = "ws://"

// New returns a new device identified with [id, model, name] tuple.  fs is the
// device's embedded file system.  targets is a list of targets support by the
// device.  e.g. ["rpi", "nano-rp2040"].
func New(id, model, name string, fs embed.FS, targets []string) dean.Thinger {
	fmt.Println("NEW DEVICE\r")
	d := &Device{
		Thing:    dean.NewThing(id, model, name),
		Targets:  target.MakeTargets(targets),
		WifiAuth: make(WifiAuth),
		WsScheme: defaultWsScheme,
		fs:       fs,
	}
	d.data = d
	// Do any OS-specific initialization
	d.deviceOSInit()
	return d
}

func (d *Device) SetData(data any) {
	d.data = data
}

// ParamFirstValue returns the first value os html param named by key
func (d *Device) ParamFirstValue(key string) string {
	if v, ok := d.deployValues[key]; ok {
		return v[0]
	}
	return ""
}

func (d *Device) GetDeployParams() string {
	return d.DeployParams
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

func (d *Device) SetLocked(locked bool) {
	d.Locked = locked
}
