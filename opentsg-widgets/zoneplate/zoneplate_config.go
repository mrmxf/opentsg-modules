package zoneplate

import (
	_ "embed"

	"github.com/mrmxf/opentsg-modules/opentsg-widgets/utils/parameters"
)

type ZConfig struct {
	PlateType   string                 `json:"plateType,omitempty" yaml:"plateType,omitempty"`
	WaveType    string                 `json:"waveType,omitempty" yaml:"waveType,omitempty"`
	Startcolour string                 `json:"startColor,omitempty" yaml:"startColor,omitempty"`
	Colors      []parameters.HexString `json:"colors,omitempty" yaml:"colors,omitempty"`

	Frequency parameters.AngleField `json:"frequency,omitempty" yaml:"frequency,omitempty"`
	// embed the angle
	parameters.RotationAngle `yaml:",inline"`
	parameters.Offset        `yaml:",inline"`
}

//go:embed jsonschema/zoneplateschema.json
var Schema []byte
