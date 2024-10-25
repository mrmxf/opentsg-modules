package gradients

import (
	"context"
	"image"
	"image/draw"

	"github.com/mrmxf/opentsg-modules/opentsg-core/colour"
	"github.com/mrmxf/opentsg-modules/opentsg-widgets/text"
)

// labels places the label on the stripe based on the angle of the stripe, the text does not change angle
func (txt textObjectJSON) labels(target draw.Image, colourSpace colour.ColorSpace, label, angle string) {

	var canvas draw.Image

	// get the flat row
	switch angle {
	case rotate180, noRotation:
		bounds := target.Bounds()
		bounds.Max.Y = int((txt.TextHeight) * float64(bounds.Max.Y) / 100)
		canvas = colour.NewNRGBA64(colourSpace, bounds)
		// canvas = image.NewNRGBA64(bounds)
	case rotate270, rotate90:
		// canvas = image.NewNRGBA64(image.Rect(0, 0, target.Bounds().Dy(), (con.TextProperties.TextHeight*target.Bounds().Dx())/100))
		canvas = colour.NewNRGBA64(colourSpace, image.Rect(0, 0, target.Bounds().Dy(), int((txt.TextHeight)*float64(target.Bounds().Max.X)/100)))
	}

	mc := context.Background()

	txtBox := text.NewTextboxer(colourSpace,
		text.WithTextColourString(txt.TextColour),
		text.WithXAlignment(txt.TextXPosition),
		text.WithYAlignment(txt.TextYPosition),
		text.WithFont(text.FontPixel),
		text.WithFill(text.FillTypeFull),
	)

	txtBox.DrawString(canvas, &mc, label)

	// rotate the text and transpose it on
	// @TODO figure out how to make this more efficent
	b := canvas.Bounds().Max
	var intermediate draw.Image
	intermediate = image.NewNRGBA64(target.Bounds())
	switch angle {
	case rotate90:
		for x := 0; x <= b.X; x++ {
			for y := 0; y <= b.Y; y++ {
				c := canvas.At(x, y)
				intermediate.Set(b.Y-y, x, c)
			}
		}
	case rotate270:
		for x := 0; x <= b.X; x++ {
			for y := 0; y <= b.Y; y++ {
				c := canvas.At(x, y)
				intermediate.Set(y, b.X-x, c)
			}
		}

	case rotate180:

		for x := 0; x <= b.X; x++ {
			for y := 0; y <= b.Y; y++ {
				c := canvas.At(x, y)
				intermediate.Set(b.X-x, b.Y-y, c)
			}
		}

	default:
		intermediate = canvas

	}

	// add the label
	colour.Draw(target, target.Bounds(), intermediate, image.Point{}, draw.Over)
}
