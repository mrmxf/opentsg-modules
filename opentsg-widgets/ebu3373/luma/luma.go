// package luma generates the ebu3373 luma bar
package luma

import (
	"math"

	"github.com/mrmxf/opentsg-modules/opentsg-core/colour"
	"github.com/mrmxf/opentsg-modules/opentsg-core/tsg"
)

const (
	WidgetType = "builtin.ebu3373/luma"
)

func (l LumaJSON) Handle(resp tsg.Response, req *tsg.Request) {

	b := resp.BaseImage().Bounds().Max
	// Set so the image scales at any size
	// Will need a caveat for widths < 1015 that the values will increase by more than one 10 bit pixel at a time
	wScale3 := 3.0 * (float64(b.X) / 3840)
	// Ceil the block width for floats and the extra pixel then goes to the start
	// This is following the design of the charts
	blockWidth := int(math.Ceil((float64(b.X) - wScale3*1015.0) / 2.0))

	for x := 0; x < b.X; x++ {
		var setColour colour.CNRGBA64
		// Check the x position and set the relevant colour
		switch {
		case x < blockWidth:
			setColour = colour.CNRGBA64{R: 4096, G: 4096, B: 4096, A: 0xffff, ColorSpace: req.PatchProperties.ColourSpace}
		case x >= (blockWidth + int(math.Ceil(wScale3*1015.0))):
			setColour = colour.CNRGBA64{R: 46144, G: 46144, B: 46144, A: 0xffff, ColorSpace: req.PatchProperties.ColourSpace}
		case x >= blockWidth && x < (blockWidth+int(math.Ceil(wScale3*1015.0))):
			// Calculate the changer per pixel and add to the base off 4
			pixVal := (float32(x-blockWidth) / float32(wScale3)) + 4.0
			// Floor the value and assign it as a 16 bit value
			// Aces.RGBA128{uint16(pixVal) << 6, uint16(pixVal) << 6, uint16(pixVal) << 6, 0xffff}
			setColour = colour.CNRGBA64{R: uint16(pixVal) << 6, G: uint16(pixVal) << 6, B: uint16(pixVal) << 6, A: 0xffff, ColorSpace: req.PatchProperties.ColourSpace}
		}

		// Set for the same colour for the depth of the ramp
		for y := 0; y < b.Y; y++ {
			resp.BaseImage().Set(x, y, &setColour)
		}
	}

	resp.Write(tsg.WidgetSuccess, "success")
}
