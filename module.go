//go:build !tinygo

package device

import (
	"golang.org/x/mod/modfile"
)

type Module struct {
	Path    string             // go module path
	Require []*modfile.Require // required dependencies
}

func (d Device) LoadModule() Module {
	var m Module

	data, err := d.CompositeFs.ReadFile("go.mod")
	if err != nil {
		m.Path = "include go.mod in your go:embed files"
		return m
	}

	file, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		m.Path = err.Error()
		return m
	}

	m.Path = file.Module.Mod.Path
	m.Require = file.Require

	return m
}
