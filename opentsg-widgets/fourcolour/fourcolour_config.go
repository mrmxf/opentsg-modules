package fourcolour

import (
	"github.com/mrmxf/opentsg-modules/opentsg-widgets/utils/parameters"
)

type Config struct {
	Colourpallette []parameters.HexString `json:"colors" yaml:"colors"`
}

var Schema = []byte(`{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"$id": "https://example.com/product.schema.json",
	"title": "Allow anything through for tests",
	"description": "An empty schema to allow custom structs to run through",
	"type": "object"
	}`)
