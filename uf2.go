//go:build !tinygo

package device

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func (d *Device) generateUf2(dir, target string) (err error) {

	file, err := os.CreateTemp("", "build-*.go")
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			os.Remove(file.Name())
		}
		// Otherwise leave /tmp/build-*.go file for debugging
	}()

	tmpl := d.templates.Lookup("build.tmpl")
	if tmpl == nil {
		return fmt.Errorf("template build.tmpl not found")
	}

	err = tmpl.Execute(file, d)
	if err != nil {
		return err
	}

	// Build the uf2 file
	uf2Name := d.Model + "-" + target + ".uf2"
	output := filepath.Join(dir, uf2Name)
	cmd := exec.Command("tinygo", "build", "-target", target, "-o", output, "-stack-size", "8kb", "-size", "short", file.Name())
	fmt.Println(cmd.String())
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, stdoutStderr)
	}
	fmt.Println(string(stdoutStderr))

	// For trying to repo TinyGO Issue #4206
	rmcmd := exec.Command("rm", "-f", "~/.cache/tinygo/obj-*")
	fmt.Println(rmcmd.String())
	out, err := rmcmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, out)
	}

	return nil
}

func (d *Device) GenerateUf2s(dir string) error {
	for target, _ := range d.Targets.TinyGoTargets() {
		if err := d.generateUf2(dir, target); err != nil {
			return err
		}
	}
	return nil
}
