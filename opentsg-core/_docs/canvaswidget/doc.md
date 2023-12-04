# Canvas Widget

Canvas widget contains all of the properties for the base
test pattern. All of the widgets are based off of this widget.
Canvas widget is a mock widget, it always runs first as it initialises
the rest of the image properties for the frame, without it nothing would configure.

```go
// ConfigVals is the go struct of all the configuration values that may be called by an input.
type ConfigVals struct {
	Type        string            `json:"type" yaml:"type"`
	Name        []string          `json:"name,omitempty" yaml:"name,omitempty"`
	ColourSpace colour.ColorSpace `json:"ColorSpace,omitempty" yaml:"ColorSpace,omitempty"`
	Framesize   config.Framesize  `json:"frameSize,omitempty" yaml:"frameSize,omitempty"`
	LineWidth   float64           `json:"linewidth,omitempty" yaml:"linewidth,omitempty"`
	FileDepth   int               `json:"filedepth,omitempty" yaml:"filedepth,omitempty"`
	GridRows    int               `json:"gridRows,omitempty" yaml:"gridRows,omitempty"`
	GridColumns int               `json:"gridColumns,omitempty" yaml:"gridColumns,omitempty"`
	BaseImage   string            `json:"baseImage,omitempty" yaml:"baseImage,omitempty"`
	Geometry    string            `json:"geometry,omitempty" yaml:"geometry,omitempty"`
	LineColor   string            `json:"lineColor,omitempty" yaml:"lineColor,omitempty"`
	Background  string            `json:"backgroundFillColor,omitempty" yaml:"backgroundFillColor,omitempty"`
	ImageType   string            `json:"imageType,omitempty" yaml:"imageType,omitempty"`
	Analytics   analytics         `json:"frame analytics" yaml:"frame analytics"`
}

type analytics struct {
	Configs enabled `json:"configuration" yaml:"configuration"`
	Average enabled `json:"average color" yaml:"average color"`
	PHash   enabled `json:"phash" yaml:"phash"`
}

```
