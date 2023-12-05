package gradients

import (
	"github.com/mrmxf/opentsg-modules/opentsg-core/colour"
	"github.com/mrmxf/opentsg-modules/opentsg-core/config"
)

type Ramp struct {
	Gradients        groupContents     `json:"groupsTemplates,omitempty" yaml:"groupsTemplates,omitempty"`
	Groups           []RampProperties  `json:"groups,omitempty" yaml:"groups,omitempty"`
	WidgetProperties control           `json:"widgetProperties,omitempty" yaml:"widgetProperties,omitempty"`
	ColourSpace      colour.ColorSpace `json:"colorSpace,omitempty" yaml:"colorSpace,omitempty"`
	GridLoc          *config.Grid      `json:"grid,omitempty" yaml:"grid,omitempty"`
}

type groupContents struct {
	GroupSeparator    groupSeparator    `json:"separator,omitempty" yaml:"separator,omitempty"`
	GradientSeparator gradientSeparator `json:"gradientSeparator,omitempty" yaml:"gradientSeparator,omitempty"`
	Gradients         []Gradient        `json:"gradients,omitempty" yaml:"gradients,omitempty"`
}

type textObjectJSON struct {
	TextYPosition string  `json:"textyPosition,omitempty" yaml:"textyPosition,omitempty"`
	TextXPosition string  `json:"textxPosition,omitempty" yaml:"textxPosition,omitempty"`
	TextHeight    float64 `json:"textHeight,omitempty" yaml:"textHeight,omitempty"`
	TextColour    string  `json:"textColor,omitempty" yaml:"textColor,omitempty"`
}

type RampProperties struct {
	Colour            string `json:"color,omitempty" yaml:"color,omitempty"`
	InitialPixelValue int    `json:"initialPixelValue,omitempty" yaml:"initialPixelValue,omitempty"`
	Reverse           bool   `json:"reverse,omitempty" yaml:"reverse,omitempty"`
}
type Gradient struct {
	Height   int    `json:"height,omitempty" yaml:"height,omitempty"`
	BitDepth int    `json:"bitDepth,omitempty" yaml:"bitDepth,omitempty"`
	Label    string `json:"label,omitempty" yaml:"label,omitempty"`

	// things that are added on run throughs
	startPoint int
	reverse    bool

	// Things we generate
	base   control
	colour string
}

type groupSeparator struct {
	Height int    `json:"height" yaml:"height"`
	Colour string `json:"color" yaml:"color"`
}

type gradientSeparator struct {
	Colours []string `json:"colors" yaml:"colors"`
	Height  int      `json:"height" yaml:"height"`
	// things the user does not assign
	base control
	step int
}

type control struct {
	MaxBitDepth      int            `json:"maxBitDepth" yaml:"maxBitDepth"`
	CwRotation       string         `json:"cwRotation" yaml:"cwRotation"`
	ObjectFitFill    bool           `json:"objectFitFill" yaml:"objectFitFill"`
	PixelValueRepeat int            `json:"pixelValueRepeat" yaml:"pixelValueRepeat"`
	TextProperties   textObjectJSON `json:"textProperties" yaml:"textProperties"`
	// These are things the user does not set
	/*
		fill function - for rotation to automatically translate the fill location
		fill - get stepsize and end goal

		step size - fill or truncate. Add a multiplier


	*/

	angleType      string
	truePixelShift float64
}

var textBoxSchema = []byte(`{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"$id": "https://example.com/product.schema.json",
	"title": "Allow anything through for tests",
	"description": "An empty schema to allow custom structs to run through",
	"type": "object"
	}`)

func (r Ramp) Alias() string {
	return r.GridLoc.Alias
}

func (r Ramp) Location() string {
	return r.GridLoc.Location
}
