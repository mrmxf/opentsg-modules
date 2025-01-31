package geometrytext

import (
	"image"
	"image/draw"
	"math"

	"github.com/mrmxf/opentsg-modules/opentsg-core/colour"
	"github.com/mrmxf/opentsg-modules/opentsg-core/tsg"
	"github.com/mrmxf/opentsg-modules/opentsg-widgets/text"
)

const WidgetType = "builtin.geometrytext"

func (gt Config) Handle(resp tsg.Response, req *tsg.Request) {

	flats := req.PatchProperties.Geometry
	// fmt.Println(len(flats), gt.GridLoc)
	// This is too intensive as text box does way more than this widget needs

	geomBox := text.NewTextboxer(req.PatchProperties.ColourSpace,
		text.WithFont(text.FontPixel),
		text.WithTextColourString(gt.TextColour),
	)

	// extract colours here and text

	///	cont := context.Background()
	for _, f := range flats {

		segment := req.GenerateSubImage(resp.BaseImage(), image.Rect(0, 0, f.Shape.Dx(), f.Shape.Dy()))

		// go for 2:1 height to width
		// find which length of letters fits this ratio the best
		tagLength := float64(len(f.ID))
		bestLength := tagLength
		bestRatio := float64(f.Shape.Dy()) / (float64(f.Shape.Dx()) / tagLength)
		for i := 1.0; i < tagLength; i *= 2 {

			newRatio := (float64(f.Shape.Dy()) / i) / (float64(f.Shape.Dx()) / (tagLength / i))
			if math.Abs(2-newRatio) < math.Abs(2-bestRatio) {
				bestLength = math.Round(tagLength / i)
				bestRatio = newRatio
			}
		}

		// split the lines into the best length
		lines := int(math.Round(tagLength / bestLength))
		tagLines := make([]string, lines)
		for i := 0; i < len(tagLines); i++ {
			end := i*int(bestLength) + int(bestLength)
			if end > len(f.ID) {
				end = len(f.ID)
			}
			tagLines[i] = f.ID[i*int(bestLength) : end]
		}
		// lines := strings.Split(f.Name, " ")
		geomBox.DrawStringsHandler(segment, req, tagLines)
		colour.Draw(resp.BaseImage(), f.Shape, segment, image.Point{}, draw.Over)
		// geomBox.DrawStrings(f.Shape, cont, lines)

	}

	resp.Write(tsg.WidgetSuccess, "success")
}
