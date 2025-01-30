package saturation

import (
	_ "embed"
)

type Config struct {
	// Type    string       `json:"type" yaml:"type"`
	Colours []string `json:"colors,omitempty" yaml:"colors,omitempty"`
}

//go:embed jsonschema/satschema.json
var Schema []byte

/*
func (s saturationJSON) Alias() string {
	return s.GridLoc.Alias
}

func (s saturationJSON) Location() string {
	return s.GridLoc.Location
}*/
