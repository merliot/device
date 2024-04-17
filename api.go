//go:build !tinygo

package device

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

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

func (d *Device) showCode(w http.ResponseWriter, r *http.Request) {
	// Retrieve top-level entries
	entries, _ := fs.ReadDir(d.CompositeFs, ".")
	// Collect entry names
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	w.Header().Set("Content-Type", "text/html")
	d.RenderTemplate(w, "code.tmpl", names)
}

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

// Set Content-Type: "text/plain" on go, mod, css, and template files
var textFile = regexp.MustCompile("\\.(go|mod|tmpl|css)$")

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
	case "":
		if len(parts) > 1 && parts[1] != "" {
			// Must be a directory
			http.FileServer(http.FS(d.CompositeFs)).ServeHTTP(w, r)
			return
		}
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
	case "code":
		d.showCode(w, r)
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
