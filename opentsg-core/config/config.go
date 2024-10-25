// package config contains the configuration information for setting up the
// initial state of opentsg. It is used for generating each frames properties
// and contains the global variables for the grid information
package config

/////////////////////////////////////////
// Standard functions and structs to use//
////////////////////////////////////////

// framesize contains the width and height of the image to be generated
type Framesize struct {
	W int `json:"w,omitempty"`
	H int `json:"h,omitempty"`
}

// Position gives the x,y coordinates of a widget
type Position struct {
	X float64 `json:"x,omitempty" yaml:"x,omitempty"`
	Y float64 `json:"y,omitempty" yaml:"y,omitempty"`
}

// WidgetGrid is the gridgen layout to be used by all widgets.
// simply import it by embedding it in your struct
/*
e.g.

type mydemo struct {

gridgen.Grid
}
*/
type WidgetGrid struct {
	GridLoc *Grid `json:"grid,omitempty" yaml:"grid,omitempty"`
}

// Grid gives the grid system with the coordinates and an alias
type Grid struct {
	Location string `json:"location,omitempty" yaml:"location,omitempty"`
	Alias    string `json:"alias,omitempty" yaml:"alias,omitempty"`
}

// default return values
func (w WidgetGrid) Alias() string {
	return w.GridLoc.Alias
}

func (w WidgetGrid) Location() string {
	return w.GridLoc.Location
}

/*

check x,y?
do not mix and match for the moment

each one has a field that can be handled


*/
