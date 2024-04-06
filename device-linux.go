// Device Linux-specific code

//go:build !tinygo

package device

import (
	"embed"
	"encoding/json"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"

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

// RenderTemplate writes the rendered template name using data to writer
func (d *Device) RenderTemplate(w http.ResponseWriter, name string, data any) {
	tmpl := d.templates.Lookup(name)
	if tmpl != nil {
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Template '"+name+"' not found", http.StatusBadRequest)
	}
}

/*
func (d *Device) showCode(w http.ResponseWriter, r *http.Request) {
	// Retrieve top-level entries
	entries, _ := fs.ReadDir(d.CompositeFs, ".")
	// Collect entry names
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	w.Header().Set("Content-Type", "text/html")
	d.RenderTemplate(w, "code.tmpl")
}
*/

func (d *Device) showState(w http.ResponseWriter, data any) {
	state, _ := json.MarshalIndent(data, "", "\t")
	d.RenderTemplate(w, "state.tmpl", string(state))
}

func (d *Device) showInstructions(w http.ResponseWriter, parts []string, data any) {
	if len(parts) >= 2 {
		target := parts[1]
		d.RenderTemplate(w, "instructions-"+target+".tmpl", data)
	} else {
		http.Error(w, "Path does not have enough elements", http.StatusBadRequest)
	}
}

// ya, i know, generics
func reverseSlice(s []string) {
	for i := 0; i < len(s)/2; i++ {
		j := len(s) - i - 1
		s[i], s[j] = s[j], s[i]
	}
}

// Set Content-Type: "text/plain" on go, css, and template files
var textFile = regexp.MustCompile("\\.(go|tmpl|css)$")

// Set Content-Type: "application/javascript" on js files
var scriptFile = regexp.MustCompile("\\.js$")

// API is the base device's API.  Derived devices can have their own API
// function to overide or extend this API.
func (d *Device) API(w http.ResponseWriter, r *http.Request, data any) {

	id, _, _ := d.Identity()

	pingPeriod := strconv.Itoa(d.PingPeriod)
	d.WebSocket = d.WsScheme + r.Host + "/ws/" + id + "/?ping-period=" + pingPeriod

	path := r.URL.Path
	parts := strings.Split(path, "/")
	reverseSlice(parts)
	switch parts[0] {
	case "", "index.html":
		d.ViewMode = ViewFull
		d.RenderTemplate(w, "index.tmpl", data)
	case "tile":
		d.ViewMode = ViewTile
		d.RenderTemplate(w, "index.tmpl", data)
	case "download-dialog":
		d.RenderTemplate(w, "download.tmpl", data)
	case "download":
		d.deploy(w, r)
	case "instructions":
		d.showInstructions(w, parts, data)
	case "info-dialog":
		d.RenderTemplate(w, "info.tmpl", data)
		/*
			case "code":
				d.showCode(w, r, data)
		*/
	case "state":
		d.showState(w, data)
	default:
		if textFile.MatchString(path) {
			w.Header().Set("Content-Type", "text/plain")
		}
		if scriptFile.MatchString(path) {
			w.Header().Set("Content-Type", "application/javascript")
		}
		http.FileServer(http.FS(d.CompositeFs)).ServeHTTP(w, r)
	}
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
