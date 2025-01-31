package noise

import (
	_ "embed"
)

type Config struct {
	//	Type      string       `json:"type" yaml:"type"`
	NoiseType string     `json:"noiseType" yaml:"noiseType"`
	Minimum   int        `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	Maximum   int        `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	YOffsets  Guillotine `json:"yOffset,omitempty" yaml:"yOffset,omitempty"`
}

// go for top then bottom and work from there
type Guillotine struct {
	TopLeft     int `json:"topLeft,omitempty" yaml:"topLeft,omitempty"`
	TopRight    int `json:"topRight,omitempty" yaml:"topRight,omitempty"`
	BottomRight int `json:"bottomRight,omitempty" yaml:"bottomRight,omitempty"`
	BottomLeft  int `json:"bottomLeft,omitempty" yaml:"bottomLeft,omitempty"`
}

//go:embed jsonschema/noiseschema.json
var Schema []byte
