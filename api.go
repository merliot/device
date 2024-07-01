//go:build !tinygo

package device

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"strconv"
)

func (d *Device) setupAPI() {
	d.ServeMux = http.NewServeMux()
	d.ServeMux.Handle("/", http.FileServer(http.FS(d.CompositeFs)))
	d.ServeMux.Handle("/{$}", d.showIndex())
	d.ServeMux.Handle("/state", d.showState())
	d.ServeMux.Handle("/code", d.showCode())
	d.ServeMux.Handle("/tile", d.renderTemplate("index-tile.tmpl"))
	d.ServeMux.Handle("/info-dialog", d.renderTemplate("info.tmpl"))
	d.ServeMux.Handle("/download-dialog", d.renderTemplate("download.tmpl"))
	d.ServeMux.HandleFunc("/download", d.deploy)
	d.ServeMux.Handle("/instructions/{target}", d.showInstructions())
}

func (d *Device) render(w io.Writer, name string, data any) error {
	tmpl := d.templates.Lookup(name)
	if tmpl == nil {
		return fmt.Errorf("Template '%s' not found", name)
	}
	return tmpl.Execute(w, data)
}

// Render template name to HTML
func (d *Device) RenderHTML(name string) (template.HTML, error) {
	var buf bytes.Buffer
	if err := d.render(&buf, name, d); err != nil {
		return template.HTML(""), err
	}
	return template.HTML(buf.String()), nil
}

// RenderTemplate writes the rendered template name using data to writer
func (d *Device) RenderTemplate(w http.ResponseWriter, name string, data any) {
	if err := d.render(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (d *Device) renderTemplateData(w http.ResponseWriter, name string, data any) {
	if err := d.render(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (d *Device) renderTemplate(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d.renderTemplateData(w, name, d.data)
	})
}

func (d *Device) showIndex() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pingPeriod := strconv.Itoa(d.PingPeriod)
		d.WebSocket = d.WsScheme + r.Host + "/ws/?ping-period=" + pingPeriod + "&trunk"
		d.renderTemplateData(w, "index.tmpl", d.data)
	})
}

func (d *Device) showCode() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve top-level entries
		entries, _ := fs.ReadDir(d.CompositeFs, ".")
		// Collect entry names
		names := make([]string, 0, len(entries))
		for _, entry := range entries {
			names = append(names, entry.Name())
		}
		d.renderTemplateData(w, "code.tmpl", names)
	})
}

func (d *Device) showState() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, _ := json.MarshalIndent(d.data, "", "\t")
		d.renderTemplateData(w, "state.tmpl", string(state))
	})
}

func (d *Device) showInstructions() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		target := r.PathValue("target")
		d.renderTemplateData(w, "instructions-"+target+".tmpl", d.data)
	})
}
