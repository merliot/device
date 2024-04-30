package models

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/merliot/dean"
)

type Modeler interface {
	Icon() []byte
	DescHtml() []byte
	SupportedTargets() string
	GenerateUf2s(string) error
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type Model struct {
	Module           string
	Maker            string
	Modeler          `json:"-"`
	Icon             string           `json:"-"`
	DescHtml         template.HTML    `json:"-"`
	SupportedTargets string           `json:"-"`
	ServeHTTP        http.HandlerFunc `json:"-"`
}

type Models map[string]Model // keyed by model

func New(model string, maker dean.ThingMaker) Model {
	modeler := maker("proto", model, "proto").(Modeler)
	return Model{
		Modeler:          modeler,
		Icon:             base64.StdEncoding.EncodeToString(modeler.Icon()),
		DescHtml:         template.HTML(modeler.DescHtml()),
		SupportedTargets: modeler.SupportedTargets(),
		ServeHTTP:        modeler.ServeHTTP,
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
