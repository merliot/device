//go:build !tinygo

package device

import (
	"embed"
	"encoding/json"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/merliot/dean"
)

//go:embed css images js template
var deviceFs embed.FS

type deviceOS struct {
	WebSocket   string            `json:"-"`
	PingPeriod  int               `json:"-"`
	CompositeFs *dean.CompositeFS `json:"-"`
	templates   *template.Template
}

func (d *Device) deviceOSInit() {
	d.PingPeriod = 4
	//d.PingPeriod = 60
	d.CompositeFs = dean.NewCompositeFS()
	d.CompositeFs.AddFS(deviceFs)
	d.templates = d.CompositeFs.ParseFS("template/*")
}

func RenderTemplate(templates *template.Template, w http.ResponseWriter, name string, data any) {
	tmpl := templates.Lookup(name)
	if tmpl != nil {
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Template '"+name+"' not found", http.StatusBadRequest)
	}
}

func (d *Device) showCode(templates *template.Template, w http.ResponseWriter, r *http.Request) {
	// Retrieve top-level entries
	entries, _ := fs.ReadDir(d.CompositeFs, ".")
	// Collect entry names
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	w.Header().Set("Content-Type", "text/html")
	RenderTemplate(templates, w, "code.tmpl", names)
}

func ShowState(templates *template.Template, w http.ResponseWriter, data any) {
	state, _ := json.MarshalIndent(data, "", "\t")
	RenderTemplate(templates, w, "state.tmpl", string(state))
}

// Set Content-Type: "text/plain" on go, css, and template files
var textFile = regexp.MustCompile("\\.(go|tmpl|js|css)$")

// Set Content-Type: "application/javascript" on js files
var scriptFile = regexp.MustCompile("\\.(go|tmpl|js|css)$")

func (d *Device) API(templates *template.Template, w http.ResponseWriter, r *http.Request) {

	id, _, _ := d.Identity()

	pingPeriod := strconv.Itoa(d.PingPeriod)
	d.WebSocket = wsScheme + r.Host + "/ws/" + id + "/?ping-period=" + pingPeriod

	path := r.URL.Path
	switch strings.TrimPrefix(path, "/") {
	case "", "index.html":
		RenderTemplate(templates, w, "index.tmpl", d)
	case "download":
		RenderTemplate(templates, w, "download.tmpl", d)
	case "info":
		RenderTemplate(templates, w, "info.tmpl", d)
	case "deploy":
		d.deploy(templates, w, r)
	case "code":
		d.showCode(templates, w, r)
	case "state":
		ShowState(templates, w, d)
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

func (d *Device) Load() {
	bytes, err := os.ReadFile("devs/" + d.Id + ".json")
	if err == nil {
		json.Unmarshal(bytes, &d.DeployParams)
	}
}

func (d *Device) Save() {
	bytes, err := json.MarshalIndent(d.DeployParams, "", "\t")
	if _, err := os.Stat("devs/"); os.IsNotExist(err) {
		// If the directory doesn't exist, create it
		os.Mkdir("devs/", os.ModePerm)
	}
	if err == nil {
		os.WriteFile("devs/"+d.Id+".json", bytes, 0600)
	}
}
