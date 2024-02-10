//go:build !tinygo

package device

import (
	"embed"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/merliot/dean"
)

//go:embed css images js template favicon.ico
var deviceFs embed.FS

const defaultPingPeriod int = 4

type deviceOS struct {
	WebSocket   string            `json:"-"`
	PingPeriod  int               `json:"-"`
	CompositeFs *dean.CompositeFS `json:"-"`
	filePath    string
	templates   *template.Template
}

func (d *Device) deviceOSInit() {
	d.PingPeriod = defaultPingPeriod
	d.CompositeFs = dean.NewCompositeFS()
	d.CompositeFs.AddFS(deviceFs)
	d.CompositeFs.AddFS(d.fs)
	d.templates = d.CompositeFs.ParseFS("template/*")
}

func (d *Device) renderTemplate(w http.ResponseWriter, name string, data any) {
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
	d.renderTemplate(w, "code.tmpl")
}
*/

func (d *Device) showState(w http.ResponseWriter) {
	state, _ := json.MarshalIndent(d, "", "\t")
	d.renderTemplate(w, "state.tmpl", string(state))
}

// Set Content-Type: "text/plain" on go, css, and template files
var textFile = regexp.MustCompile("\\.(go|tmpl|css)$")

// Set Content-Type: "application/javascript" on js files
var scriptFile = regexp.MustCompile("\\.js$")

func (d *Device) API(w http.ResponseWriter, r *http.Request, data any) {

	id, _, _ := d.Identity()

	pingPeriod := strconv.Itoa(d.PingPeriod)
	d.WebSocket = wsScheme + r.Host + "/ws/" + id + "/?ping-period=" + pingPeriod

	path := r.URL.Path
	switch strings.TrimPrefix(path, "/") {
	case "", "index.html":
		d.ViewMode = ViewFull
		d.renderTemplate(w, "index.tmpl", data)
	case "tile":
		d.ViewMode = ViewTile
		d.renderTemplate(w, "index.tmpl", data)
	case "download-dialog":
		d.renderTemplate(w, "download.tmpl", data)
	case "download":
		d.deploy(w, r)
	case "info-dialog":
		d.renderTemplate(w, "info.tmpl", data)
		/*
			case "code":
				d.showCode(w, r, data)
		*/
	case "state":
		d.showState(w)
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

func (d *Device) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.API(w, r, d)
}

func (d *Device) Load(filePath string) error {
	d.filePath = filePath
	bytes, err := os.ReadFile(d.filePath)
	if err == nil {
		json.Unmarshal(bytes, &d.DeployParams)
	}
	return err
}

func (d *Device) Save() error {
	dir := filepath.Dir(d.filePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err == nil {
		var bytes []byte
		bytes, err = json.MarshalIndent(d.DeployParams, "", "\t")
		if err == nil {
			err = os.WriteFile(d.filePath, bytes, 0600)
		}
	}
	return err
}

func (d *Device) Icon() []byte {
	icon, _ := deviceFs.ReadFile("images/icon.png")
	return icon
}

func (d *Device) DescHtml() []byte {
	return []byte("Missing DescHtml()")
}

func (d *Device) SupportedTargets() string {
	return "Missing SupportedTargets()"
}
