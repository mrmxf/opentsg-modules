package textbox

import (
	"github.com/mrmxf/opentsg-modules/opentsg-core/colour"
	"github.com/mrmxf/opentsg-modules/opentsg-core/config"
)

type TextboxJSON struct {
	// Type       string       `json:"type" yaml:"type"`
	Text []string `json:"text,omitempty" yaml:"text,omitempty"`

	GridLoc     *config.Grid      `json:"grid,omitempty" yaml:"grid,omitempty"`
	ColourSpace colour.ColorSpace `json:"colorSpace,omitempty" yaml:"colorSpace,omitempty"`
	Border      string            `json:"borderColor,omitempty" yaml:"borderColor,omitempty"`
	BorderSize  float64           `json:"borderSize,omitempty" yaml:"borderSize,omitempty"`
	Font        string            `json:"font,omitempty" yaml:"font,omitempty"`

	Back       string `json:"backgroundColor,omitempty" yaml:"backgroundColor,omitempty"`
	Textc      string `json:"textColor,omitempty" yaml:"textColor,omitempty"`
	FillType   string `json:"fillType,omitempty" yaml:"fillType,omitempty"`
	XAlignment string `json:"xAlignment,omitempty" yaml:"xAlignment,omitempty"`
	YAlignment string `json:"yAlignment,omitempty" yaml:"yAlignment,omitempty"`
}

var textBoxSchema = []byte(`{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"$id": "https://example.com/product.schema.json",
	"title": "Allow anything through for tests",
	"description": "An empty schema to allow custom structs to run through",
	"type": "object"
	}`)

func (tb TextboxJSON) Alias() string {
	return tb.GridLoc.Alias
}

func (tb TextboxJSON) Location() string {
	return tb.GridLoc.Location
}
