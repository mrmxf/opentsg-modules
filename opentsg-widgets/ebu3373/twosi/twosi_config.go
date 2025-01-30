package twosi

import (
	_ "embed"
)

type Config struct {
	//	Type    string      `json:"type" yaml:"type"`
}

//go:embed jsonschema/twoschema.json
var Schema []byte
