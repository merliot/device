//go:build !tinygo

package device

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/merliot/uf2"
)

func genFile(templates *template.Template, template string, name string,
	values map[string]string) error {

	tmpl := templates.Lookup(template)
	if tmpl == nil {
		return fmt.Errorf("Template '%s' not found", template)
	}

	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, values)
}

func (d *Device) serveFile(dir, filename string, w http.ResponseWriter, r *http.Request) error {

	// Calculate MD5 checksum of installer
	cmd := exec.Command("md5sum", dir+"/"+filename)
	println(cmd.String())
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, stdoutStderr)
	}
	md5sum := bytes.Fields(stdoutStderr)[0]
	md5sumBase64 := base64.StdEncoding.EncodeToString(md5sum)

	// Set the Content-Disposition header to suggest the original filename for download
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	// Set the MD5 checksum header
	w.Header().Set("Content-MD5", md5sumBase64)

	http.ServeFile(w, r, dir+"/"+filename)

	return nil
}

func (d *Device) deployGo(dir string, values map[string]string, envs []string,
	templates *template.Template, w http.ResponseWriter, r *http.Request) error {

	// Generate build.go from server.tmpl
	if err := genFile(templates, "server.tmpl", dir+"/build.go", values); err != nil {
		return err
	}

	// Generate installer.go from installer.tmpl
	if err := genFile(templates, "installer.tmpl", dir+"/installer.go", values); err != nil {
		return err
	}

	// Generate model.service from service.tmpl
	if err := genFile(templates, "service.tmpl", dir+"/"+d.Model+".service", values); err != nil {
		return err
	}

	// Generate model.conf from log.tmpl
	if err := genFile(templates, "log.tmpl", dir+"/"+d.Model+".conf", values); err != nil {
		return err
	}

	// Build build.go -> model (binary)

	// substitute "-" for "_" in target, ala TinyGo, when using as tag
	target := strings.Replace(values["target"], "-", "_", -1)

	cmd := exec.Command("go", "build", "-o", dir+"/"+d.Model, "-tags", target, dir+"/build.go")
	println(cmd.String())
	cmd.Env = append(cmd.Environ(), envs...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, stdoutStderr)
	}

	// Build installer and serve as download-able file

	installer := d.Id + "-installer"
	cmd = exec.Command("go", "build", "-o", dir+"/"+installer, dir+"/installer.go")
	println(cmd.String())
	cmd.Env = append(cmd.Environ(), envs...)
	stdoutStderr, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, stdoutStderr)
	}

	return d.serveFile(dir, installer, w, r)
}

/*
func (d *Device) deployTinyGo(dir string, values map[string]string, envs []string,
	templates *template.Template, w http.ResponseWriter, r *http.Request) error {

	// Generate build.go from runner.tmpl
	if err := genFile(templates, "runner.tmpl", dir+"/build.go", values); err != nil {
		return err
	}

	// Build build.go -> uf2 binary

	installer := d.Id + "-installer.uf2"
	target := values["target"]

	cmd := exec.Command("tinygo", "build", "-target", target, "-stack-size", "8kb",
		"-o", dir+"/"+installer, dir+"/build.go")
	println(cmd.String())
	cmd.Env = append(cmd.Environ(), envs...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, stdoutStderr)
	}

	return d.serveFile(dir, installer, w, r)
}
*/

func (d *Device) deployTinyGoUF2(dir string, values map[string]string, envs []string,
	templates *template.Template, w http.ResponseWriter, r *http.Request) error {

	var p = params{
		Ssid:         values["ssid"],
		Passphrase:   values["passphrase"],
		Id:           values["id"],
		Model:        values["model"],
		Name:         values["name"],
		DeployParams: values["deployParams"],
		User:         values["user"],
		Passwd:       values["passwd"],
		DialURLs:     values["dialUrls"],
	}

	target := values["target"]

	base := target + ".uf2"
	installer := d.Id + "-installer.uf2"

	// Re-write the base uf2 file and save as the installer uf2 file.
	// The paramsMem area is replace by json-encoded params.

	uf2, err := uf2.Read(base)
	if err != nil {
		return err
	}

	oldBytes := bytes.Repeat([]byte{byte('x')}, 2048)
	newBytes := make([]byte, 2048)

	newParams, err := json.Marshal(p)
	if err != nil {
		return err
	}
	copy(newBytes, newParams)

	uf2.ReplaceBytes(oldBytes, newBytes)

	err = uf2.Write(dir + "/" + installer)
	if err != nil {
		return err
	}

	return d.serveFile(dir, installer, w, r)
}

func (d *Device) buildValues(r *http.Request) (map[string]string, error) {

	var values = make(map[string]string)

	// Squash request params down to first value for each key.  The resulting
	// map[string]string is much nicer to pass to html/template as data value.

	for k, v := range r.URL.Query() {
		if len(v) > 0 {
			values[k] = v[0]
		}
	}

	values["deployParams"] = d.DeployParams
	values["id"] = d.Id
	values["model"] = d.Model
	values["modelStruct"] = strings.Title(d.Model)
	values["name"] = d.Name

	if ssid, ok := values["ssid"]; ok {
		values["passphrase"] = d.WifiAuth[ssid]
	}

	values["hub"] = wsScheme + r.Host + "/ws/?ping-period=4"
	println("\n\n\nHUB", values["hub"], "\n\n\n")

	if values["backuphub"] != "" {
		u, err := url.Parse(values["backuphub"])
		if err != nil {
			return nil, err
		}
		scheme := "ws://"
		if u.Scheme == "https" {
			scheme = "wss://"
		}
		values["backuphub"] = scheme + u.Host + "/ws/?ping-period=4"
	}

	values["dialUrls"] = values["hub"] + "," + values["backuphub"]

	if user, passwd, ok := r.BasicAuth(); ok {
		values["user"] = user
		values["passwd"] = passwd
	}

	return values, nil
}

func (d *Device) buildEnvs(values map[string]string) []string {
	envs := []string{}
	switch values["target"] {
	case "demo", "x86-64":
		envs = []string{"CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64"}
	case "rpi":
		// TODO: do we want more targets for GOARM=7|8?
		envs = []string{"CGO_ENABLED=0", "GOOS=linux", "GOARCH=arm", "GOARM=5"}
	}
	return envs
}

func (d *Device) _deploy(templates *template.Template, w http.ResponseWriter, r *http.Request) error {

	values, err := d.buildValues(r)
	if err != nil {
		return err
	}

	envs := d.buildEnvs(values)

	// Create temp build directory
	dir, err := os.MkdirTemp("./", d.Id+"-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	//println(dir)

	switch values["target"] {
	case "demo", "x86-64", "rpi":
		return d.deployGo(dir, values, envs, templates, w, r)
	case "nano-rp2040", "wioterminal", "pyportal":
		//return d.deployTinyGo(dir, values, envs, templates, w, r)
		return d.deployTinyGoUF2(dir, values, envs, templates, w, r)
	default:
		return errors.New("Target not supported")
	}

	return nil
}

func (d *Device) deploy(templates *template.Template, w http.ResponseWriter, r *http.Request) {
	d.DeployParams = r.URL.RawQuery
	if err := d._deploy(templates, w, r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	d.Save()
}
