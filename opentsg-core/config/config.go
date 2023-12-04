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

// Grid gives the grid system with the coordinates and an alias
type Grid struct {
	Location string `json:"location,omitempty" yaml:"location,omitempty"`
	Alias    string `json:"alias,omitempty" yaml:"alias,omitempty"`
}
