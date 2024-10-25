package gridgen

import (
	"context"
	"fmt"
	"image"
	"image/draw"

	"github.com/mrmxf/opentsg-modules/opentsg-core/config/core"
)

type Box interface {
	generateBox(c *context.Context) (draw.Image, image.Point, draw.Image, error)
}

/*




design ideas

    "grid": {
        "alias": " NotherNameForThisLocation",
        "location": "(200,1700)-(3640,1900)"
    },

	becomes

    "box": {
        "alias": " noiseBox",
        "bounds": {"x":200, "y":1700, "w":3640, "h":1900}
    },

    "box": {
        "alias": " noiseBox2",
        "bounds": {"x":200, "y":1700, "x2":8000, "y2":4000}
    },


	"box": {
        "alias": " noiseBox2",
        "coordinates": {"x":"200px", "y":1700, "x2":8000, "y2":4000}
    },

	{"x":200, "y":1700, "w":3640, "h":1900} implicitly a top-left pinned box because it has 4 properties x, y, w, h
{"x":200, "y":1700, "x2":3640, "y2":1900} implicitly a corner pinned box because it has 4 properties x, y, x2, y2
{"cx":200, "cy":1700, "x2":3640, "y2":1900} do the

{"cx":200, "cy":1700, "radius":20px}

edge antiasliasing questions
inheritance positions questions

format for xy coordinates

it its coordinate then

x,y as pixels. Each value is the grid, no sub grid componenets yet

so 16,16 would then be used. Bin off A1, R1C1?

*/

type idea struct {

	// keep the alias from last time
	alias string
	//
	bounds bounds
}

// implement hsl(0, 100%, 50%);

// keep these ideas in mind https://www.w3schools.com/css/css_boxmodel.asp
/*
remove margin, padding and border as there will not be that mich need

thoughts - sa y the limits so people dont think this is a dierect css import
as fetures will deffo be missing

*/
type bounds struct {
	useAlias string

	// top left coordinates
	// actually any
	x, y any

	x2, y2 any

	// width height
	// can they be A or 1 etc. just mix it up
	w, h any
	// impement stuff like this
	// calc(100% - 100px)

	// centre values
	// or would they be float for half sizes etc
	cx, cy any
	// or masks like this. Leave masks out for the moment?
	//  mask-image: radial-gradient(circle, black 50%, rgba(0, 0, 0, 0.5) 50%);

	// how many x values to include for corners etc

	// circle properties
	// can then get the area around?
	radius int

	// for the mask generation
	offsets any
}

// https://www.w3schools.com/css/css_boxmodel.asp
func (b bounds) generateBox(c *context.Context) (draw.Image, image.Point, draw.Image, error) {

	// get the start point
	/*
		either xy or cx cy
		 or useALias
	*/
	aliasMap := core.GetAlias(*c)
	switch {
	case b.useAlias != "":
		// @TODO update the alias to be the
		// the image.Point, canvas size and a mask,
		// if applicable
		loc := aliasMap.Data[b.useAlias]
		if loc != "" {
			// call the function again but with the required coordinates
			mid, _ := gridSquareLocatorAndGenerator(loc, "", c)
			return mid.GImage, image.Point{mid.X, mid.Y}, mid.GMask, nil
		} else {

			return nil, image.Point{}, nil, fmt.Errorf(invalidAlias, b.useAlias)
		}
	case b.x != nil || b.y != nil:
	case b.cx != nil || b.cy != nil:
	default:
		// return no coordinate postion used
	}

	/*
		get end locatoin

		wh, r, or x2y2 for xy
		wh, r, for cxcy
	*/

	/*
		now we have the square image the mask is calculated
		which is if radius or offsets (or both)
	*/

	// get the end point

	// returns the mask, coordinate and base image and error
	return nil, image.Point{}, nil, nil
}

/*
func anyToLength(coordinate any) int {

	coord := string(fmt.Sprintf("%v", coordinate))

	regSpreadX := regexp.MustCompile(`^[a-zA-Z]{1,}$`)
	regCoord := regexp.MustCompile(`^[0-9]{1,}$`)

	regPixels := regexp.MustCompile(`^[0-9]{1,}[Pp][Xx]$`)
	regXY := regexp.MustCompile(`^\(-{0,1}[0-9]{1,5},-{0,1}[0-9]{1,5}\)-\(-{0,1}[0-9]{1,5},-{0,1}[0-9]{1,5}\)$`)
	regRC := regexp.MustCompile(`^[Rr]([\d]{2,}|[1-9]{1})[Cc]([\d]{2,}|[1-9]{1})$`)
	regRCArea := regexp.MustCompile(`^[Rr]([\d]{2,}|[1-9]{1})[Cc]([\d]{2,}|[1-9]{1}):[Rr]([\d]{2,}|[1-9]{1})[Cc]([\d]{2,}|[1-9]{1})$`)

	switch {
	case true:
	default:
	}

	return 0
}
*/
