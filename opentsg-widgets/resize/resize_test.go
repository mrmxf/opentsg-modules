package resize

import (
	"crypto/sha256"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"testing"

	"github.com/mrmxf/opentsg-modules/opentsg-core/colour"
	"github.com/mrmxf/opentsg-modules/opentsg-core/tsg"
	examplejson "github.com/mrmxf/opentsg-modules/opentsg-widgets/exampleJson"
	"github.com/mrmxf/opentsg-modules/opentsg-widgets/text"
	"github.com/mrmxf/opentsg-modules/opentsg-widgets/utils/parameters"
	"github.com/nfnt/resize"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExample(t *testing.T) {

	examples := []Config{
		{XDetections: []*parameters.DistanceField{{Dist: "3840px"}, {Dist: "1920px"}},
			YDetections: []*parameters.DistanceField{{Dist: "1080px"}, {Dist: "540px"}}},

		{XDetections: []*parameters.DistanceField{{Dist: "75%"}, {Dist: "50%"}},
			YDetections: []*parameters.DistanceField{{Dist: "50%"}, {Dist: "25%"}},
			Graticule:   graticule{Position: text.AlignmentMiddle, TextColor: "#C2A649"}},

		{XStep: &parameters.DistanceField{Dist: "500px"}, YStep: &parameters.DistanceField{Dist: "250px"}},

		{XStep: &parameters.DistanceField{Dist: "10%"}, YStep: &parameters.DistanceField{Dist: "10%"},
			XStepEnd: &parameters.DistanceField{Dist: "50%"}, YStepEnd: &parameters.DistanceField{Dist: "50%"}},

		{XStep: &parameters.DistanceField{Dist: "10%"}, XStepEnd: &parameters.DistanceField{Dist: "50%"},
			Graticule: graticule{Position: text.AlignmentLeft, TextColor: "#C2A649", GraticuleColour: "#9A3A73"}},
	}

	desc := []string{"pixel", "percentage", "stepper", "limitedStepper", "graticule"}

	for i, example := range examples {

		examplejson.SaveExampleJsonRequest(example, &tsg.Request{FrameProperties: tsg.FrameProperties{FrameDimensions: image.Point{3840, 2160}}},
			WidgetType, desc[i], true)
	}

}
func TestResizeDirections(t *testing.T) {

	canvasSize := []image.Point{{X: 4096, Y: 2160}, {X: 1920, Y: 1080}, {X: 1280, Y: 720}}
	destX := []string{"3840px", "1280px", "960px"}
	destXInt := []uint{3840, 1280, 960}

	for i, p := range canvasSize {
		base := image.NewNRGBA64(image.Rectangle{image.Point{}, p})

		out := tsg.TestResponder{BaseImg: base}
		Config{XDetections: []*parameters.DistanceField{{Dist: destX[i]}}}.Handle(
			&out, &tsg.Request{FrameProperties: tsg.FrameProperties{FrameDimensions: p}})

		// f, _ := os.Create(fmt.Sprintf("./testdata/resize%v.png", p.X))
		// png.Encode(f, base)

		file, _ := os.Open(fmt.Sprintf("./testdata/resize%v.png", p.X))
		baseVals, _ := png.Decode(file)

		// Assign the colour to the correct type of image NGRBA64 and replace the colour values
		readImage := image.NewNRGBA64(baseVals.Bounds())
		colour.Draw(readImage, readImage.Bounds(), baseVals, image.Point{0, 0}, draw.Over)

		// Make a hash of the pixels of each image
		hnormal := sha256.New()
		htest := sha256.New()
		hnormal.Write(readImage.Pix)
		htest.Write(base.Pix)

		Convey("Checking the detectors are generated", t, func() {
			Convey(fmt.Sprintf("generating an x detection of %v for a canvas of %v", destX, canvasSize[i]), func() {
				Convey("No error is returned and the file matches exactly", func() {
					So(out.Message, ShouldResemble, "success")
					So(htest.Sum(nil), ShouldResemble, hnormal.Sum(nil))
				})
			})
		})

		down := resize.Resize(destXInt[i], uint(p.Y), base, resize.Lanczos3)
		up := resize.Resize(uint(p.X), uint(p.Y), down, resize.Lanczos3)

		fl, _ := os.Create(fmt.Sprintf("./testdata/resizeLanczos%v.png", p.X))
		png.Encode(fl, up)

		downbi := resize.Resize(destXInt[i], uint(p.Y), base, resize.Bicubic)
		upbi := resize.Resize(uint(p.X), uint(p.Y), downbi, resize.Bicubic)

		fb, _ := os.Create(fmt.Sprintf("./testdata/resizeBicubic%v.png", p.X))
		png.Encode(fb, upbi)

	}

	destY := []string{"1080px"}
	destYint := []uint{1080}

	for i := range destY {
		base := image.NewNRGBA64(image.Rectangle{image.Point{}, canvasSize[i]})

		out := tsg.TestResponder{BaseImg: base}
		Config{YDetections: []*parameters.DistanceField{{Dist: destY[i]}}}.Handle(
			&out, &tsg.Request{FrameProperties: tsg.FrameProperties{FrameDimensions: canvasSize[i]}})

		// f, _ := os.Create(fmt.Sprintf("./testdata/resize%v.png", canvasSize[i].Y))
		// png.Encode(f, base)

		file, _ := os.Open(fmt.Sprintf("./testdata/resize%v.png", canvasSize[i].Y))

		baseVals, _ := png.Decode(file)

		// Assign the colour to the correct type of image NGRBA64 and replace the colour values
		readImage := image.NewNRGBA64(baseVals.Bounds())
		colour.Draw(readImage, readImage.Bounds(), baseVals, image.Point{0, 0}, draw.Over)

		// Make a hash of the pixels of each image
		hnormal := sha256.New()
		htest := sha256.New()
		hnormal.Write(readImage.Pix)
		htest.Write(base.Pix)

		Convey("Checking the detectors are generated", t, func() {
			Convey(fmt.Sprintf("generating an y detection of %v for a canvas of %v", destY, canvasSize[i]), func() {
				Convey("No error is returned and the file matches exactly", func() {
					So(out.Message, ShouldResemble, "success")
					So(htest.Sum(nil), ShouldResemble, hnormal.Sum(nil))
				})
			})
		})

		down := resize.Resize(uint(canvasSize[i].X), destYint[i], base, resize.Lanczos3)
		up := resize.Resize(uint(canvasSize[i].X), uint(canvasSize[i].Y), down, resize.Lanczos3)

		fl, _ := os.Create(fmt.Sprintf("./testdata/resizeLanczos%v.png", canvasSize[i].Y))
		png.Encode(fl, up)

		downbi := resize.Resize(uint(canvasSize[i].X), destYint[i], base, resize.Bicubic)
		upbi := resize.Resize(uint(canvasSize[i].X), uint(canvasSize[i].Y), downbi, resize.Bicubic)

		fb, _ := os.Create(fmt.Sprintf("./testdata/resizeBicubic%v.png", canvasSize[i].Y))
		png.Encode(fb, upbi)

	}
}

func TestBoxes(t *testing.T) {

	boxSizes := []int{3, 6, 27, 11, 7}
	imageSize := []image.Point{
		{1920, 1080},
		{1080, 1920},
		{1920, 1080},
		{1080, 1920},
		{100, 1000},
	}

	baseImageSize := image.Point{X: 1920, Y: 1080}

	for i, size := range boxSizes {

		fill := make([]*parameters.DistanceField, size)

		for i := range fill {
			fill[i] = &parameters.DistanceField{Dist: "1280px"}
		}
		myImage := image.NewNRGBA64(image.Rectangle{image.Point{}, imageSize[i]})
		out := tsg.TestResponder{BaseImg: myImage}
		Config{XDetections: fill}.Handle(
			&out, &tsg.Request{FrameProperties: tsg.FrameProperties{FrameDimensions: baseImageSize}})

		// f, _ := os.Create(fmt.Sprintf("./testdata/%vbox.png", size))
		// png.Encode(f, myImage)
		// Open the image to compare to
		file, _ := os.Open(fmt.Sprintf("./testdata/%vbox.png", size))
		// Decode to get the colour values
		baseVals, _ := png.Decode(file)

		// Assign the colour to the correct type of image NGRBA64 and replace the colour values
		readImage := image.NewNRGBA64(baseVals.Bounds())
		colour.Draw(readImage, readImage.Bounds(), baseVals, image.Point{0, 0}, draw.Over)

		// Make a hash of the pixels of each image
		hnormal := sha256.New()
		htest := sha256.New()
		hnormal.Write(readImage.Pix)
		htest.Write(myImage.Pix)

		// Save the file
		Convey("Checking the boxes fill the space easily", t, func() {
			Convey(fmt.Sprintf("seeing how %v boxes fill a space of %v", size, imageSize[i]), func() {
				Convey("No error is returned and the file matches exactly", func() {
					So(out.Message, ShouldResemble, "success")
					So(htest.Sum(nil), ShouldResemble, hnormal.Sum(nil))
				})
			})
		})
	}
}

func TestSteps(t *testing.T) {

	steps := []parameters.DistanceField{{Dist: 10}, {Dist: 15}, {Dist: "100px"}, {Dist: "256px"}}
	ends := []parameters.DistanceField{{Dist: 10}, {Dist: 0}, {Dist: "800px"}, {Dist: "512px"}}
	graticules := []string{text.AlignmentLeft, text.AlignmentRight, text.AlignmentTop, text.AlignmentBottom}
	imageSize := image.Point{1920, 1080}

	for i, step := range steps {

		myImage := image.NewNRGBA64(image.Rectangle{image.Point{}, imageSize})

		out := tsg.TestResponder{BaseImg: myImage}
		Config{XStep: &step, XStepEnd: &ends[i],
			Graticule: graticule{Position: graticules[i]}}.Handle(
			&out, &tsg.Request{FrameProperties: tsg.FrameProperties{FrameDimensions: imageSize}})

		// f, _ := os.Create(fmt.Sprintf("./testdata/steps/step%v%s.png", step.Dist, graticules[i]))
		// png.Encode(f, myImage)
		// Open the image to compare to
		file, _ := os.Open(fmt.Sprintf("./testdata/steps/step%v%s.png", step.Dist, graticules[i]))
		// Decode to get the colour values
		baseVals, _ := png.Decode(file)

		// Assign the colour to the correct type of image NGRBA64 and replace the colour values
		readImage := image.NewNRGBA64(baseVals.Bounds())
		colour.Draw(readImage, readImage.Bounds(), baseVals, image.Point{0, 0}, draw.Over)

		// Make a hash of the pixels of each image
		hnormal := sha256.New()
		htest := sha256.New()
		hnormal.Write(readImage.Pix)
		htest.Write(myImage.Pix)

		// Save the file
		Convey("Checking the boxes fill the space easily", t, func() {
			Convey(fmt.Sprintf("seeing how %v boxes fill a space of %v", step, imageSize), func() {
				Convey("No error is returned and the file matches exactly", func() {
					So(out.Message, ShouldResemble, "success")
					So(htest.Sum(nil), ShouldResemble, hnormal.Sum(nil))
				})
			})
		})
	}
}
