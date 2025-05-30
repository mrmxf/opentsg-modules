// package nearblack generates the ebu3373 nearblack bar
package nearblack

import (
	"image"
	"image/draw"

	"github.com/mrmxf/opentsg-modules/opentsg-core/colour"
	"github.com/mrmxf/opentsg-modules/opentsg-core/tsg"
)

const (
	WidgetType = "builtin.ebu3373/nearblack"
)

var (
	neg4  = colour.CNRGBA64{R: 2048, G: 2048, B: 2048, A: 0xffff}
	neg2  = colour.CNRGBA64{R: 3072, G: 3072, B: 3072, A: 0xffff}
	neg1  = colour.CNRGBA64{R: 3584, G: 3584, B: 3584, A: 0xffff}
	black = colour.CNRGBA64{R: 4096, G: 4096, B: 4096, A: 0xffff}
	pos1  = colour.CNRGBA64{R: 4608, G: 4608, B: 4608, A: 0xffff}
	pos2  = colour.CNRGBA64{R: 5120, G: 5120, B: 5120, A: 0xffff}
	pos4  = colour.CNRGBA64{R: 6144, G: 6144, B: 6144, A: 0xffff}

	grey = colour.CNRGBA64{R: 26496, G: 26496, B: 26496, A: 0xffff}
)

func (nb Config) Handle(resp tsg.Response, req *tsg.Request) {

	b := resp.BaseImage().Bounds().Max
	greyRun := grey
	greyRun.UpdateColorSpace(req.PatchProperties.ColourSpace)
	colour.Draw(resp.BaseImage(), resp.BaseImage().Bounds(), &image.Uniform{&greyRun}, image.Point{}, draw.Src)
	// Scale everything so it fits the shape of the canvas
	wScale := (float64(b.X) / 3840.0)
	startPoint := wScale * 480
	off := wScale * 206

	order := []colour.CNRGBA64{neg4, neg2, neg1, pos1, pos2, pos4}
	area := image.Rect(int(startPoint), 0, int(startPoint+off*2), b.Y)
	colour.Draw(resp.BaseImage(), area, &image.Uniform{&black}, image.Point{}, draw.Src)
	startPoint += off * 2
	for _, c := range order {
		// alternate through the colours
		fill := c
		fill.UpdateColorSpace(req.PatchProperties.ColourSpace)
		colour.Draw(resp.BaseImage(), image.Rect(int(startPoint), 0, int(startPoint+off), b.Y), &image.Uniform{&c}, image.Point{}, draw.Src)
		startPoint += off
		// append with the 0% black
		blackRun := black
		blackRun.UpdateColorSpace(req.PatchProperties.ColourSpace)
		colour.Draw(resp.BaseImage(), image.Rect(int(startPoint), 0, int(startPoint+off), b.Y), &image.Uniform{&blackRun}, image.Point{}, draw.Src)
		startPoint += off
	}

	resp.Write(tsg.WidgetSuccess, "success")
}
