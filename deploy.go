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
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/merliot/device/uf2"
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
	cmd := exec.Command("md5sum", filepath.Join(dir, filename))
	fmt.Println(cmd.String())
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

	http.ServeFile(w, r, filepath.Join(dir, filename))

	return nil
}

func (d *Device) deployGo(dir string, values map[string]string, envs []string,
	templates *template.Template, w http.ResponseWriter, r *http.Request) error {

	// Generate build.go from server.tmpl
	if err := genFile(templates, "server.tmpl", filepath.Join(dir, "build.go"), values); err != nil {
		return err
	}

	// Generate installer.go from installer.tmpl
	if err := genFile(templates, "installer.tmpl", filepath.Join(dir, "installer.go"), values); err != nil {
		return err
	}

	// Generate model.service from service.tmpl
	if err := genFile(templates, "service.tmpl", filepath.Join(dir, d.Model+".service"), values); err != nil {
		return err
	}

	// Generate model.conf from log.tmpl
	if err := genFile(templates, "log.tmpl", filepath.Join(dir, d.Model+".conf"), values); err != nil {
		return err
	}

	// Build build.go -> model (binary)

	// substitute "-" for "_" in target, ala TinyGo, when using as tag
	target := strings.Replace(values["target"], "-", "_", -1)

	cmd := exec.Command("go", "build", "-ldflags", "-s -w", "-o", filepath.Join(dir, d.Model),
		"-tags", target, filepath.Join(dir, "build.go"))
	fmt.Println(cmd.String())
	cmd.Env = append(cmd.Environ(), envs...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, stdoutStderr)
	}

	// Build installer and serve as download-able file

	installer := d.Id + "-installer"
	cmd = exec.Command("go", "build", "-ldflags", "-s -w", "-o", filepath.Join(dir, installer), filepath.Join(dir, "installer.go"))
	fmt.Println(cmd.String())
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
	fmt.Println(cmd.String())
	cmd.Env = append(cmd.Environ(), envs...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, stdoutStderr)
	}

	return d.serveFile(dir, installer, w, r)
}
*/

// Random string to embed in UF2 we can search for later to locate params
const UF2Magic = "Call the Doctor!  Miss you Dan."

func (d *Device) deployTinyGoUF2(dir string, values map[string]string, envs []string,
	templates *template.Template, w http.ResponseWriter, r *http.Request) error {

	var p = params{
		MagicStart:   UF2Magic,
		Ssid:         values["ssid"],
		Passphrase:   values["passphrase"],
		Id:           values["id"],
		Model:        values["model"],
		Name:         values["name"],
		DeployParams: values["deployParams"],
		User:         values["user"],
		Passwd:       values["passwd"],
		DialURLs:     values["dialUrls"],
		MagicEnd:     UF2Magic,
	}

	target := values["target"]

	base := filepath.Join("uf2s", d.Model+"-"+target+".uf2")
	installer := d.Id + "-installer.uf2"

	// Re-write the base uf2 file and save as the installer uf2 file.
	// The paramsMem area is replaced by json-encoded params.

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

	err = uf2.Write(filepath.Join(dir, installer))
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
	values["modelStruct"] = d.ModelStruct
	values["name"] = d.Name
	values["modulePath"] = d.Module.Path

	if ssid, ok := values["ssid"]; ok {
		values["passphrase"] = d.WifiAuth[ssid]
	}

	values["hub"] = d.WsScheme + r.Host + "/ws/?ping-period=4"
	values["dialUrls"] = values["hub"] + "," + d.DialURLs

	if user, passwd, ok := r.BasicAuth(); ok {
		values["user"] = user
		values["passwd"] = passwd
	}

	//fmt.Println(values)
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
	//fmt.Println(dir)

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

func (d *Device) deploy(w http.ResponseWriter, r *http.Request) {
	if d.Locked {
		http.Error(w, "Refusing to download, device is locked", http.StatusLocked)
		return
	}
	d.SetDeployParams(r.URL.RawQuery)
	if err := d._deploy(d.templates, w, r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
