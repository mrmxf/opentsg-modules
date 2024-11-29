package framecount

import (
	_ "embed"

	"github.com/mrmxf/opentsg-modules/opentsg-core/colour"
	"github.com/mrmxf/opentsg-modules/opentsg-core/config"
)

type Config struct {
	//	Type         string            `json:"type" yaml:"type"`
	FrameCounter      bool        `json:"frameCounter,omitempty" yaml:"frameCounter,omitempty"`
	Imgpos            interface{} `json:"gridPosition,omitempty" yaml:"gridPosition,omitempty"`
	TextColour        string      `json:"textColor,omitempty" yaml:"textColor,omitempty"`
	BackColour        string      `json:"backgroundColor,omitempty" yaml:"backgroundColor,omitempty"`
	Font              string      `json:"font,omitempty" yaml:"font,omitempty"`
	FontSize          float64     `json:"fontSize,omitempty" yaml:"fontSize,omitempty"`
	config.WidgetGrid `yaml:",inline"`
	ColourSpace       colour.ColorSpace `json:"colorSpace,omitempty" yaml:"colorSpace,omitempty"`

	//	DesignScale  string       `json:"designScale" yaml:"designScale"`
	// This is added in for metadata purposes
	frameNumber int `json:"frameNumber"`
}

// start the count at -1 as it is incremented before being returned
var framecount = -1

//go:embed jsonschema/framecounter.json
var Schema []byte

func (f *Config) getFrames() bool {
	if f.FrameCounter {
		framecount++
	}

	return f.FrameCounter
}

func framePos() int {
	return framecount
}

/*
func (f frameJSON) Alias() string {
	return f.GridLoc.Alias
}

func (f frameJSON) Location() string {
	return f.GridLoc.Location
}*/
