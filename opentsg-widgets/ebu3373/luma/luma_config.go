package luma

import (
	_ "embed"
)

type LumaJSON struct {
}

//go:embed jsonschema/lumaschema.json
var Schema []byte

/*
func (l lumaJSON) Alias() string {
	return l.GridLoc.Alias
}

func (l lumaJSON) Location() string {
	return l.GridLoc.Location
}*/
