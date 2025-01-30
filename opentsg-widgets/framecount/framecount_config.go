package framecount

import (
	_ "embed"
)

type Config struct {
	//	Type         string            `json:"type" yaml:"type"`
	FrameCounter bool        `json:"frameCounter,omitempty" yaml:"frameCounter,omitempty"`
	Imgpos       interface{} `json:"gridPosition,omitempty" yaml:"gridPosition,omitempty"`
	TextColour   string      `json:"textColor,omitempty" yaml:"textColor,omitempty"`
	BackColour   string      `json:"backgroundColor,omitempty" yaml:"backgroundColor,omitempty"`
	Font         string      `json:"font,omitempty" yaml:"font,omitempty"`
	FontSize     float64     `json:"fontSize,omitempty" yaml:"fontSize,omitempty"`

	//	DesignScale  string       `json:"designScale" yaml:"designScale"`

}

//go:embed jsonschema/framecounter.json
var Schema []byte
