//go:build !tinygo

package device

import (
	"fmt"
	"os"
	"os/exec"
)

func (d *Device) generateUf2(target string) error {

	file, err := os.CreateTemp("", "build-*.go")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

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
	cmd := exec.Command("tinygo", "build", "-target", target, "-o", uf2Name, "-stack-size", "8kb", "-size", "short", file.Name())
	println(cmd.String())
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, stdoutStderr)
	}
	println(string(stdoutStderr))

	return nil
}

func (d *Device) GenerateUf2s() error {
	for target, _ := range d.Targets.TinyGoTargets() {
		if err := d.generateUf2(target); err != nil {
			return err
		}
	}
	return nil
}
