// Device Linux-specific code

//go:build !tinygo

package device

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/merliot/dean"
)

//go:embed css images html js template favicon.ico
var deviceFs embed.FS

const defaultPingPeriod int = 4

// Linux device structure
type deviceOS struct {
	// WebSocket is the websocket address for a client to dial back into
	// the device
	WebSocket string `json:"-"`
	// PingPeriod is the ping period, measured in seconds
	PingPeriod int `json:"-"`
	// CompositeFs is the device's fs.  Derived devices will overlay
	// their fs onto this fs.
	CompositeFs *dean.CompositeFS `json:"-"`
	// Module is extracted go.mod info
	Module    `json:"-"`
	templates *template.Template
}

// Linux device struct init
func (d *Device) deviceOSInit() {
	d.PingPeriod = defaultPingPeriod
	d.CompositeFs = dean.NewCompositeFS()
	d.CompositeFs.AddFS(deviceFs)
	d.CompositeFs.AddFS(d.fs)
	d.Module = d.LoadModule()
	d.templates = d.CompositeFs.ParseFS("template/*")
}

// ServeHTTP is the base device's web server handler for /.  Derived devices
// can override or extend with their own ServeHTTP function.
func (d *Device) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.API(w, r, d)
}

// Icon is the base device's icon image
func (d *Device) Icon() []byte {
	icon, _ := d.fs.ReadFile("images/icon.png")
	return icon
}

// DescHtml is the base device's description
func (d *Device) DescHtml() []byte {
	desc, _ := d.fs.ReadFile("html/desc.html")
	return desc
}

// Supported Targets is the base device's targets in full format
func (d *Device) SupportedTargets() string {
	return d.Targets.FullNames()
}
