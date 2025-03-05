package gridgen

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/mrmxf/opentsg-modules/opentsg-core/colour"
)

/*
Location is the location information of a widget.

It contains the box properties in css style for drawing shapes.
*/
type Location struct {

	// What is the alias for this location
	Alias string `json:"alias,omitempty" yaml:"alias,omitempty"`
	// CSS style fields for drawing the box
	Box Box `json:"box,omitempty" yaml:"box,omitempty"`
}

/*
Box implements the box properties of a widget.
The XY coordinates start at 0,0 in the top left of the canvas.

It tries to adhere to css style of a box https://www.w3schools.com/css/css_boxmodel.asp
Not every property is implemented however.
*/
type Box struct {
	// use a predeclared alias
	// alias must be declared before
	UseAlias    string   `json:"useAlias,omitempty" yaml:"useAlias,omitempty"`
	UseGridKeys []string `json:"useGridKeys,omitempty" yaml:"useGridKeys,omitempty"`

	// top left coordinates
	X any `json:"x,omitempty" yaml:"x,omitempty"`
	Y any `json:"y,omitempty" yaml:"y,omitempty"`
	// bottom right
	// if not used then the grid is 1 square
	X2 any `json:"x2,omitempty" yaml:"x2,omitempty"`
	Y2 any `json:"y2,omitempty" yaml:"y2,omitempty"`

	// Width and Height take second priority to
	// x2 ans y2 if they are called
	Width  any `json:"width,omitempty" yaml:"width,omitempty"`
	Height any `json:"height,omitempty" yaml:"height,omitempty"`

	// xAlignment and yAlignemnt values are not implemented yet
	XAlignment string `json:"xAlignment,omitempty" yaml:"xAlignment,omitempty"`
	YAlignment string `json:"yAlignment,omitempty" yaml:"yAlignment,omitempty"` // default top left but let them choose

	// border radius - follows the simple layout of
	// https://prykhodko.medium.com/css-border-radius-how-does-it-work-bfdf23792ac2
	// taps out at 50% of the shortest dimension.
	BorderRadius any `json:"border-radius,omitempty" yaml:"border-radius,omitempty"`
}

type TSIGProperties struct {
	Grouping string `json:"grouping"`
}

// InitAliasBox inits a map of the alias for handlers in a context
func InitAliasBox(c context.Context) context.Context {
	n := SyncMapBox{make(map[string]any), &sync.Mutex{}}

	return context.WithValue(c, aliasKeyBox, n)
}

// SyncMap  is a map with a sync.Mutex to prevent concurrent writes.
type SyncMapBox struct {
	Data map[string]any
	Mu   *sync.Mutex
}

// Get alias returns a map of the locations alias and their grid positions.
func GetAliasBox(c context.Context) SyncMapBox {
	Alias := c.Value(aliasKeyBox)
	if Alias != nil {

		return Alias.(SyncMapBox)
	}
	// else return an empty map
	var newmu sync.Mutex

	return SyncMapBox{Mu: &newmu, Data: make(map[string]any)}
}

// GeneratePatch converts the box struct into a canvas the size of the grid,
// the location generated is the upper left most corner of the grid, along with any masks that are required for non square
// grids.
func (l Location) GeneratePatch(c *context.Context) (draw.Image, image.Point, draw.Image, error) {

	alias := GetAliasBox(*c)

	// recursively use the alias
	if l.Box.UseAlias != "" {
		alias.Mu.Lock()
		item, ok := alias.Data[l.Box.UseAlias]
		alias.Mu.Unlock()

		if ok {
			// just recurse through
			if mid, ok := item.(Location); ok {
				return mid.GeneratePatch(c)
			}
		} else {
			return nil, image.Point{}, nil, fmt.Errorf("\"%s\" is not a valid grid alias", l.Box.UseAlias)
		}
	}

	// if grid keys are used these take priority over
	// any coordinates declared
	if len(l.Box.UseGridKeys) != 0 {
		// @TODO make art keys more useable
		// so they can be chained
		regArt := regexp.MustCompile(`^key:[\w]{3,10}$`)
		//	regTSIG := regexp.MustCompile(`^tsig:`)
		switch {
		case regArt.MatchString(l.Box.UseGridKeys[0]):

			return artToCanvas(l.Box.UseGridKeys[0], c)
		default:

			return l.tSIGToArea(c)

		}

	}

	// else calculate the coordinates
	return l.generatePatch(c)
}

// tSIGToArea converts a group of tsigs to a single area, with a mask matching their total tiles
func (b Location) tSIGToArea(c *context.Context) (draw.Image, image.Point, draw.Image, error) {

	matches, loc, bounds, err := b.calcTSIGToArea(c)

	if err != nil {
		return nil, image.Point{}, nil, err
	}

	base := ImageGenerator(*c, bounds)
	// fill in thr mask
	msk := image.NewAlpha16(bounds)
	for _, match := range matches {
		colour.Draw(msk, match.Shape.Sub(loc), &image.Uniform{color.Alpha16{A: 0xffff}}, image.Point{}, draw.Src)
	}

	return base, loc, msk, nil

}

// calcTSIGToArea calculates the Area a group of tsigs take up, from the minimum point of all the tiles to the maximum
// point. It may lead to some weird shapes depending on the grouping and layouts of the tiles used.
func (b Location) calcTSIGToArea(c *context.Context) ([]*Segmenter, image.Point, image.Rectangle, error) {

	//	get all the tiles associated with the grid keys
	matches, err := getTiles(c, b.Box.UseGridKeys)

	if err != nil {
		return nil, image.Point{}, image.Rectangle{}, err
	}

	frameSize := (*c).Value(sizekey).(image.Point)
	max := image.Point{}
	min := frameSize
	// calc max and min points
	// to calculate the bounds of the total imafw
	for _, match := range matches {

		if match.Shape.Min.X < min.X {
			min.X = match.Shape.Min.X
		}

		if match.Shape.Max.X > max.X {
			max.X = match.Shape.Max.X
		}

		if match.Shape.Min.Y < min.Y {
			min.Y = match.Shape.Min.Y
		}

		if match.Shape.Max.Y > max.Y {
			max.Y = match.Shape.Max.Y
		}

	}

	return matches, min, image.Rectangle{image.Point{}, max.Sub(min)}, nil
}

// generatePatch
func (b Location) generatePatch(c *context.Context) (draw.Image, image.Point, draw.Image, error) {

	// get the bounds of the area to be drawn
	dims, tsgLocation, err := b.CalcArea(c)
	if err != nil {
		return nil, image.Point{}, nil, err
	}

	// get the mask associated with the whole testcard
	mask := (*c).Value(tilemaskkey)
	var widgMask draw.Image
	if mask != nil {
		mask := mask.(draw.Image)
		// make it an alpha only to reduce memory usage
		widgMask = image.NewAlpha16(dims)
		// draw the mask based on its testcard location
		colour.Draw(widgMask, widgMask.Bounds(), mask, tsgLocation, draw.Src)
	}

	// calculate any rounded corners
	xUnit := (*c).Value(xkey).(float64)
	yUnit := (*c).Value(ykey).(float64)
	if b.Box.BorderRadius != nil {

		// find the minimum dimension
		// so that the rounding is capped
		// at 50% of the smallest
		xSize, dim := xUnit, dims.Max.X
		if xSize > yUnit {
			xSize = yUnit
		}

		if dim > dims.Max.Y {
			dim = dims.Max.Y
		}

		// get the radius the user wants
		rad, err := anyToDist(b.Box.BorderRadius, dim, xSize)
		r := int(rad)
		if err != nil {
			return nil, image.Point{}, nil, err
		}

		// truncate it down to 50%
		// if the given value is greater
		if r > dims.Max.X/2 {
			r = dims.Max.X / 2
		}

		if r > dims.Max.Y/2 {
			r = dims.Max.Y / 2
		}

		// finally get the rounded mask
		midMask := roundedMask(dims, int(r))
		if widgMask == nil {
			widgMask = midMask
		} else {
			// mask the tsig mask, with the rounded mask. Only in the bounds of the tsig mask.
			draw.DrawMask(widgMask, widgMask.Bounds(), midMask, image.Point{}, widgMask, image.Point{}, draw.Src)
		}

	}

	widgetCanvas := ImageGenerator(*c, dims)

	// log the whole location
	// if the alias is given
	if b.Alias != "" {
		aliasMap := GetAliasBox(*c)
		aliasMap.Mu.Lock() // prevent concurrent map writes
		aliasMap.Data[b.Alias] = b
		aliasMap.Mu.Unlock()
	}

	return widgetCanvas, tsgLocation, widgMask, nil
}

// CalcArea calculates the dimension of the box and the coordinate the top left is placed at.
func (l Location) CalcArea(c *context.Context) (image.Rectangle, image.Point, error) {
	if l.Box.X == nil || l.Box.Y == nil {
		// invalid coordinates received
		return image.Rectangle{}, image.Point{}, fmt.Errorf("invalid coordinates of x %v and y %v received", l.Box.X, l.Box.Y)
	}

	// get the dimensions of the test signal
	dimensions := (*c).Value(sizekey).(image.Point)
	xUnit := (*c).Value(xkey).(float64)
	yUnit := (*c).Value(ykey).(float64)

	// get the requested measurements for the start points
	y, err := anyToDist(l.Box.Y, dimensions.Y, yUnit)
	if err != nil {
		return image.Rectangle{}, image.Point{}, err
	}

	x, err := anyToDist(l.Box.X, dimensions.X, xUnit)
	if err != nil {
		return image.Rectangle{}, image.Point{}, err
	}

	var endY float64
	// switch the width in order of precedence
	switch {
	case l.Box.Y2 != nil:
		endY, err = anyToDist(l.Box.Y2, dimensions.Y, yUnit)
	case l.Box.Height != nil:
		var mid float64
		mid, err = anyToDist(l.Box.Height, dimensions.Y, yUnit)
		endY = y + mid
	default:
		// default is one y unit
		endY = y + yUnit
	}

	if err != nil {
		return image.Rectangle{}, image.Point{}, err
	}

	var endX float64
	// switch the height in order of precedence
	switch {
	case l.Box.X2 != nil:
		endX, err = anyToDist(l.Box.X2, dimensions.X, xUnit)
	case l.Box.Width != nil:
		var mid float64
		mid, err = anyToDist(l.Box.Width, dimensions.X, xUnit)
		endX = x + mid
	default:
		// default is one x unit
		endX = x + xUnit
	}

	if err != nil {
		return image.Rectangle{}, image.Point{}, err
	}

	// calculate the dimensions of the patch
	width := int(endX) - int(x)
	height := int(endY) - int(y)
	tsgLocation := image.Point{X: int(x), Y: int(y)}

	return image.Rect(0, 0, width, height), tsgLocation, nil
}

func roundedMask(rect image.Rectangle, radius int) draw.Image {

	// set up a mask
	base := image.NewAlpha16(rect)
	draw.Draw(base, base.Bounds(), &image.Uniform{color.Alpha16{A: 0xffff}}, image.Point{}, draw.Src)

	// get the four corners of the rectangle
	startPoints := []image.Point{{radius, radius}, {radius, rect.Max.Y - radius},
		{rect.Max.X - radius, radius}, {rect.Max.X - radius, rect.Max.Y - radius}}

	dir := []image.Point{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	// loop through the raidus from each corner
	for i, sp := range startPoints {
		for x := 0; x <= radius; x++ {
			for y := 0; y <= radius; y++ {
				//	r := xy
				r := xyToRadius(float64(x), float64(y))
				if r > float64(radius) {
					// truncate it back to 0
					base.Set(sp.X+(dir[i].X*x), sp.Y+(dir[i].Y*y), &color.Alpha16{})
				}
			}
		}

	}

	return base
}

// updateTSIGUnit updates a group of tsig tiles
// into a new set of tiles based on their groups.
func updateTSIGUnit(segments []Segmenter, group string) ([]Segmenter, error) {

	// if there are no groups do not update the segments
	if group == "" {
		return segments, nil
	}

	// group up the units
	out := make(map[string][]Segmenter)
	units := make(map[string]string)

	for _, s := range segments {

		dest, ok := s.Groups[group]
		if !ok {
			return nil, fmt.Errorf("no unit of %s, found for %s", group, s.ID)
		}

		out[dest] = append(out[dest], s)
		units[s.ID] = dest
	}

	outSegments := make([]Segmenter, len(out))

	// calculate the bounds of each segment, where the neighbours are
	i := 0
	for dest, outSegs := range out {

		max := image.Point{}
		min := image.Point{X: 0xfffffffffffffff, Y: 0xfffffffffffffff}
		var neighs []string
		for _, match := range outSegs {

			for _, neighbour := range match.Neighbours {

				out, ok := units[neighbour]

				if ok && dest != out {
					neigh := out
					if !slices.Contains(neighs, neigh) {
						neighs = append(neighs, neigh)
					}
				}

			}

			if match.Shape.Min.X < min.X {
				min.X = match.Shape.Min.X
			}

			if match.Shape.Max.X > max.X {
				max.X = match.Shape.Max.X
			}

			if match.Shape.Min.Y < min.Y {
				min.Y = match.Shape.Min.Y
			}

			if match.Shape.Max.Y > max.Y {
				max.Y = match.Shape.Max.Y
			}
		}

		// fmt.Println(min, max, neighs)
		outSegments[i] = Segmenter{ID: dest, ImportPosition: i, Shape: image.Rectangle{min, max}, Neighbours: neighs}

		i++
	}

	return outSegments, nil

}

func xyToRadius(x, y float64) float64 {
	return math.Sqrt(x*x + y*y)
}

// anyToDist converts the any type request to a tsg dimension
func anyToDist(a any, dimension int, unitWidth float64) (float64, error) {

	dist := fmt.Sprintf("%v", a)

	pixel := regexp.MustCompile(`^-{0,1}\d{1,}[pP][xX]$`)
	grid := regexp.MustCompile(`^\d{1,}$`)
	pcDefault := regexp.MustCompile(`^-{0,1}\d{0,2}\.{1}\d{0,}%$|^-{0,1}\d{0,2}%$|^-{0,1}(100)%$|^100\.[0]*%$`)

	switch {
	case pixel.MatchString(dist):

		// convert the pixels to an int
		pxDist, err := strconv.Atoi(dist[:len(dist)-2])
		if err != nil {
			err = fmt.Errorf("extracting %s as a integer: %v", dist, err.Error())
			return 0, err
		}
		return float64(pxDist), nil
	case pcDefault.MatchString(dist):

		// trim the %
		dist = dist[:len(dist)-1]

		perc, err := strconv.ParseFloat(dist, 64)
		if err != nil {
			return 0, fmt.Errorf("extracting %s as a percentage : %v", dist, err.Error())
		}
		// percentage * dimension
		totalWidth := (perc / 100) * float64(dimension)

		return totalWidth, nil
	case grid.MatchString(dist):
		// get the grid units
		unit, err := strconv.ParseFloat(dist, 64)
		if err != nil {
			return 0, fmt.Errorf("extracting %s as a percentage : %v", dist, err.Error())

		}
		totalWidth := unit * unitWidth

		return totalWidth, nil
	default:
		return 0, fmt.Errorf("unknown coordinate of %s", dist)
	}

}

// Get GridGeometry returns the geometry of a grid coordiante, localised to those grid coordinates.
// This is for use with widgets that utilise geometry.
// if not tsigs were called in the set up phase then a tsig that
// fills the area is given
func (l Location) GetGridGeometry(c *context.Context, unit string) ([]Segmenter, error) {

	// check for the alias first
	alias := GetAliasBox(*c)
	if l.Box.UseAlias != "" {
		alias.Mu.Lock()
		item, ok := alias.Data[l.Box.UseAlias]
		alias.Mu.Unlock()

		if ok {
			// just recurse through
			if mid, ok := item.(Location); ok {
				return mid.GetGridGeometry(c, unit)
			}
		} else {
			return nil, fmt.Errorf("\"%s\" is not a valid grid alias", l.Box.UseAlias)
		}
	}

	//
	var matches []*Segmenter
	var tsgPoint image.Point

	// extract the tsigs with gridkeys if they are used
	if len(l.Box.UseGridKeys) > 0 {
		var err error
		matches, tsgPoint, _, err = l.calcTSIGToArea(c)
		if err != nil {
			return nil, err
		}

	} else {
		// else calculate the area and check for overlaps
		var err error
		var bounds image.Rectangle
		bounds, tsgPoint, err = l.CalcArea(c)
		if err != nil {

			return nil, err
		}

		// set the bounds to cover the area
		bounds = bounds.Add(tsgPoint)

		// get the geometry of the test signal
		geometryHolder, err := getTiles(c, l.Box.UseGridKeys)
		if err != nil {
			return nil, err
		}

		// loop through every tile and check which ones
		// are contained in this patch
		for _, g := range geometryHolder {

			if bounds.Overlaps(g.Shape) {
				matches = append(matches, g)
			}
		}

	}

	offsetMatches := make([]Segmenter, len(matches))

	// offset the matches to be within the
	// bounds of the patch
	for i, match := range matches {
		mid := *match
		// go for negative offsets
		mid.Shape = mid.Shape.Add(image.Point{-tsgPoint.X, -tsgPoint.Y})
		offsetMatches[i] = mid

	}

	// order them by their import position
	sort.Slice(offsetMatches, func(i, j int) bool {
		return offsetMatches[i].ImportPosition < offsetMatches[j].ImportPosition
	})

	// update the units afterwards
	return updateTSIGUnit(offsetMatches, unit)

}

func getTiles(c *context.Context, gridKeys []string) ([]*Segmenter, error) {
	// set the bounds to cover the area
	geometry := (*c).Value(tilekey)

	var geometryHolder []*Segmenter
	if geometry != nil {
		geometryHolder = geometry.([]*Segmenter)
	} else {

		// make some tiles that fill the area
		// that match the grid coordinates
		// of the test card.

		squareX := (*c).Value(xkey).(float64)
		squareY := (*c).Value(ykey).(float64)
		frameConfig := (*c).Value(frameKey)
		if frameConfig == nil {
			frameConfig = FrameConfiguration{}
		}
		frame := frameConfig.(FrameConfiguration)
		// generate the grid tsigs.
		geometryHolder = gridToTSIG(frame.Cols, frame.Rows, squareX, squareY)

	}

	// if grid keys are not used then we stop here
	if len(gridKeys) == 0 {
		return geometryHolder, nil
	}

	var matches []*Segmenter
	regTSIG := regexp.MustCompile(`^tsig:`)

	targets := make([]string, len(gridKeys))
	// get the matches
	for i, key := range gridKeys {

		if !regTSIG.MatchString(key) {
			return nil, fmt.Errorf("invalid key of \"%s\" used", key)
		}
		targets[i] = key[5:]
	}

	for _, s := range geometryHolder {

		for _, target := range targets {

			var pass bool
			if s.ID == target {
				matches = append(matches, s)
				break
			}

			tar := strings.Split(target, ".")

			switch len(tar) {
			case 1:
				// just looking for field matches
				_, pass = s.Groups[tar[0]]
			case 2:
				// check the group value matches
				out := s.Groups[tar[0]]
				pass = (out == tar[1])
			default:
				return nil, fmt.Errorf("the key %s is invalid, please stick to a max dotpath length of 2", target)
			}

			if pass { // if we get a match then check
				matches = append(matches, s)
				break
			}
		}

	}

	// no matches
	if len(matches) == 0 {
		return nil, fmt.Errorf("no tiles found with the keys \"%s\"", gridKeys)
	}

	return matches, nil

}

// gridToTSIG generates a TSIG from the grid coordinates of the testcard
func gridToTSIG(x, y int, xscale, yscale float64) []*Segmenter {
	sections := make([]*Segmenter, x*y)
	count := 0

	for xpos := 0; xpos < x; xpos++ {
		for ypos := 0; ypos < y; ypos++ {
			// get the area of the xy coordinate
			area := image.Rect(int(float64(xpos)*xscale), int(float64(ypos)*yscale), int(float64(xpos+1)*xscale), int(float64(ypos+1)*yscale))
			// Make the excel spreadsheet name
			gridCoord := fmt.Sprintf("%v%v", gridToScale(xpos), ypos)
			neighbours := []string{}

			// generate the neighbours using simple if statements for each position
			// left
			if xpos != 0 {
				neighbours = append(neighbours, fmt.Sprintf("%v%v", gridToScale(xpos-1), ypos))
			}
			// to the right
			if xpos+1 < x {
				neighbours = append(neighbours, fmt.Sprintf("%v%v", gridToScale(xpos+1), ypos))
			}
			// above
			if ypos != 0 {
				neighbours = append(neighbours, fmt.Sprintf("%v%v", gridToScale(xpos), ypos-1))
			}
			// below
			if ypos+1 < y {
				neighbours = append(neighbours, fmt.Sprintf("%v%v", gridToScale(xpos), ypos+1))
			}

			sections[xpos*y+ypos] = &Segmenter{ID: gridCoord, Shape: area, Neighbours: neighbours, ImportPosition: count}
			count++
		}
	}

	return sections
}
