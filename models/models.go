package models

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"os"

	"github.com/merliot/dean"
	"github.com/merliot/device"
)

type Model struct {
	Module           string
	Maker            string
	device.Modeler   `json:"-"`
	Icon             string        `json:"-"`
	DescHtml         template.HTML `json:"-"`
	SupportedTargets string        `json:"-"`
}

type Models map[string]Model // keyed by model

func New(model string, maker dean.ThingMaker) Model {
	modeler := maker("proto", model, "proto").(device.Modeler)
	return Model{
		Modeler:          modeler,
		Icon:             base64.StdEncoding.EncodeToString(modeler.Icon()),
		DescHtml:         template.HTML(modeler.DescHtml()),
		SupportedTargets: modeler.SupportedTargets(),
	}
}

func Load(file string) (Models, error) {
	var models Models

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &models)
	return models, err
}
