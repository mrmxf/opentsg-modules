package nearblack

import (
	_ "embed"
)

type Config struct {
}

//go:embed jsonschema/nbschema.json
var Schema []byte

/*
func (nb nearblackJSON) Alias() string {
	return nb.GridLoc.Alias
}

func (nb nearblackJSON) Location() string {
	return nb.GridLoc.Location
}
*/
