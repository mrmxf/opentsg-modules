package geometrytext

import (
	_ "embed"
)

type Config struct {
	TextColour string `json:"textColor" yaml:"textColor"`
}

//go:embed jsonschema/geometryText.json
var Schema []byte
