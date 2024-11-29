package geometrytext

import (
	"github.com/mrmxf/opentsg-modules/opentsg-core/colour"
	"github.com/mrmxf/opentsg-modules/opentsg-core/config"

	_ "embed"
)

type Config struct {
	TextColour  string            `json:"textColor" yaml:"textColor"`
	GridLoc     *config.Grid      `json:"grid,omitempty" yaml:"grid,omitempty"`
	ColourSpace colour.ColorSpace `json:"colorSpace,omitempty" yaml:"colorSpace,omitempty"`
}

//go:embed jsonschema/geometryText.json
var Schema []byte

func (f Config) Alias() string {
	return f.GridLoc.Alias
}

func (f Config) Location() string {
	return f.GridLoc.Location
}
