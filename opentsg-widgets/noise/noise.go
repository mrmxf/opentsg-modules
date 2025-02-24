// Package noise generates images of noise
package noise

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"math/rand"
	"time"

	"github.com/mrmxf/opentsg-modules/opentsg-core/colour"
	"github.com/mrmxf/opentsg-modules/opentsg-core/tsg"
)

const (
	WidgetType = "builtin.noise"
)

const (
	whiteNoise = "white noise"
)

var randnum = randSeed

func randSeed() int64 {
	return time.Now().Unix()
}

func (c Config) Handle(resp tsg.Response, req *tsg.Request) {

	// Have a seed variable tht is taken out for testing purposes
	random := rand.New(rand.NewSource(randnum()))

	var max int
	if c.Maximum != 0 {
		max = c.Maximum
	} else {
		// Revert to the default
		max = 4095
	}
	min := c.Minimum

	if max < min {
		resp.Write(tsg.WidgetError, fmt.Sprintf("0141 The minimum noise value %v is greater than the maximum noise value %v", min, max))
		return
	}

	if c.NoiseType == whiteNoise { // upgrade to switch statement when more types come in
		err := c.whitenoise(random, resp.BaseImage(), req.PatchProperties.ColourSpace, min, max)
		if err != nil {
			resp.Write(tsg.WidgetError, err.Error())
			return
		}
	}

	resp.Write(tsg.WidgetSuccess, "success")

}

func (n Config) whitenoise(random *rand.Rand, canvas draw.Image, cspace colour.ColorSpace, min, max int) error {
	b := canvas.Bounds().Max

	yStart := 0
	TopOffset := 0

	if n.YOffsets.TopLeft != 0 || n.YOffsets.TopRight != 0 {
		if n.YOffsets.TopLeft > n.YOffsets.TopRight {
			yStart = n.YOffsets.TopLeft
			TopOffset = n.YOffsets.TopLeft - n.YOffsets.TopRight
		} else {
			yStart = n.YOffsets.TopRight
			TopOffset = -(n.YOffsets.TopRight - n.YOffsets.TopLeft)
		}
	}

	yMax := b.Y
	BottomOffset := 0

	if n.YOffsets.BottomLeft != 0 || n.YOffsets.BottomRight != 0 {
		if n.YOffsets.BottomLeft > n.YOffsets.BottomRight {
			yMax = b.Y - n.YOffsets.BottomLeft
			BottomOffset = n.YOffsets.BottomLeft - n.YOffsets.BottomRight
		} else {
			yMax = b.Y - n.YOffsets.BottomRight
			BottomOffset = -(n.YOffsets.BottomRight - n.YOffsets.BottomLeft)
		}
	}

	if yMax < yStart {
		return fmt.Errorf("0DEV vertical offset overlap, the offsets go past the middle in both directions. Box height : %v, top offset %v, bottom offset %v", b.Y, TopOffset, BottomOffset)
	}

	triangle(random, canvas, b, cspace, true, yStart-int(math.Abs(float64(TopOffset))), TopOffset, max, min)

	for y := yStart; y < yMax; y++ {
		for x := 0; x < b.X; x++ {
			colourPos := uint16(random.Intn(max-min)+min) << 4

			canvas.Set(x, y, &colour.CNRGBA64{R: colourPos, G: colourPos, B: colourPos, A: 0xffff, ColorSpace: cspace})
		}
	}

	// dp bottom half
	triangle(random, canvas, b, cspace, false, yMax, BottomOffset, max, min)
	/*
		Get the block height

		get the x shift per y increase

		go along each x shift doing that increase based on x or y chnage


	*/
	return nil
}

func triangle(random *rand.Rand, canvas draw.Image, b image.Point, colourSpace colour.ColorSpace, top bool, yMax, offset, max, min int) {
	if offset != 0 {
		yOffset := int(math.Abs(float64(offset)))

		xShift := b.X / yOffset

		xCount := 0
		xPos := 0

		// set it up to walk backwards
		if offset < 0 {
			xPos = b.X - xShift
		}

		off := 0

		for xCount <= yOffset {
			if top {
				off = yOffset - xCount
			}

			for y := yMax + off; y < yMax+xCount+1+off; y++ {
				for x := xPos; x < xPos+xShift; x++ {
					colourPos := uint16(random.Intn(max-min)+min) << 4

					canvas.Set(x, y, &colour.CNRGBA64{R: colourPos, G: colourPos, B: colourPos, A: 0xffff, ColorSpace: colourSpace})
				}
			}

			xCount++
			if offset < 0 {
				xPos -= xShift
			} else {
				xPos += xShift
			}
		}
	}
}
