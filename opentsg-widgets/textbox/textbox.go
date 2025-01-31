// package textbox generates textboxes.
package textbox

import (
	"image"
	"image/draw"

	"github.com/mrmxf/opentsg-modules/opentsg-core/colour"

	"github.com/mrmxf/opentsg-modules/opentsg-core/tsg"
	texter "github.com/mrmxf/opentsg-modules/opentsg-widgets/text"
)

const (
	WidgetType = "builtin.textbox"
)

func (tb TextboxJSON) Handle(resp tsg.Response, req *tsg.Request) {
	// calculate the border here

	bounds := resp.BaseImage().Bounds().Max
	var width, height float64
	if tb.BorderSize > 0 { // prevent div 0 errors
		width, height = (float64(bounds.X)*tb.BorderSize)/100, (float64(bounds.Y)*tb.BorderSize)/100
	} // else leave as 0

	// draw the borders
	borderwidth := int(height)
	if width < height {
		borderwidth = int(width)
	}

	borderColour := tb.Border.ToColour(req.PatchProperties.ColourSpace)
	colour.Draw(resp.BaseImage(), image.Rect(0, 0, borderwidth, resp.BaseImage().Bounds().Max.Y), &image.Uniform{borderColour}, image.Point{}, draw.Src)
	colour.Draw(resp.BaseImage(), image.Rect(0, 0, resp.BaseImage().Bounds().Max.X, borderwidth), &image.Uniform{borderColour}, image.Point{}, draw.Src)
	colour.Draw(resp.BaseImage(), image.Rect(resp.BaseImage().Bounds().Max.X-borderwidth, 0, resp.BaseImage().Bounds().Max.X, resp.BaseImage().Bounds().Max.Y), &image.Uniform{borderColour}, image.Point{}, draw.Src)
	colour.Draw(resp.BaseImage(), image.Rect(0, resp.BaseImage().Bounds().Max.Y-borderwidth, resp.BaseImage().Bounds().Max.X, resp.BaseImage().Bounds().Max.Y), &image.Uniform{borderColour}, image.Point{}, draw.Src)

	// get the text and background
	c := colour.NewNRGBA64(req.PatchProperties.ColourSpace, image.Rect(0, 0, resp.BaseImage().Bounds().Max.X-borderwidth*2, resp.BaseImage().Bounds().Max.Y-borderwidth*2))

	textbox := texter.NewTextboxer(req.PatchProperties.ColourSpace,

		texter.WithBackgroundColourString(tb.Back),
		texter.WithTextColourString(tb.Textc),
		texter.WithFill(tb.FillType),
		texter.WithXAlignment(tb.XAlignment),
		texter.WithYAlignment(tb.YAlignment),
		texter.WithFont(tb.Font),
	)

	err := textbox.DrawStringsHandler(c, req, tb.Text)
	if err != nil {
		resp.Write(tsg.WidgetError, err.Error())
	}

	// apply the text
	colour.Draw(resp.BaseImage(), image.Rect(borderwidth, borderwidth, resp.BaseImage().Bounds().Max.X-borderwidth, resp.BaseImage().Bounds().Max.Y-borderwidth), c, image.Point{}, draw.Src)

	resp.Write(tsg.WidgetSuccess, "success")
}
