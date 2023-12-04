package colour

import (
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTemp(t *testing.T) {
	//	testrun()
	testrun2()
	testrun2020()
	testrun709()

}

// go test ./colour/ -bench=. -benchtime=10s

/*
func BenchmarkNRGBA64Area(b *testing.B) {
	// decode to get the colour values
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		testrun2020()
	}
}

func BenchmarkNRGBA64ACESSet(b *testing.B) {
	// decode to get the colour values

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		testrun709()
	}
} */

func BenchmarkNRGBA64Draw(b *testing.B) {
	base := NewNRGBA64(ColorSpace{ColorSpace: "rec2020"}, image.Rect(0, 0, 1000, 1000))
	for n := 0; n < b.N; n++ {
		Draw(base, base.Bounds(), &image.Uniform{&CNRGBA64{R: 0xffff, A: 0xfff0, ColorSpace: ColorSpace{ColorSpace: "rec709"}}}, image.Point{}, draw.Src)
	}
}

func BenchmarkNRGBA64ImgDraw(b *testing.B) {
	base := NewNRGBA64(ColorSpace{ColorSpace: "rec2020"}, image.Rect(0, 0, 1000, 1000))
	for n := 0; n < b.N; n++ {
		draw.Draw(base, base.Bounds(), &image.Uniform{&CNRGBA64{R: 0xffff, A: 0xfff0, ColorSpace: ColorSpace{ColorSpace: "rec709"}}}, image.Point{}, draw.Src)
	}
}

func BenchmarkNRGBA64DrawMaxAlpha(b *testing.B) {
	base := NewNRGBA64(ColorSpace{ColorSpace: "rec2020"}, image.Rect(0, 0, 1000, 1000))
	for n := 0; n < b.N; n++ {
		Draw(base, base.Bounds(), &image.Uniform{&CNRGBA64{R: 0xffff, A: 0xffff, ColorSpace: ColorSpace{ColorSpace: "rec709"}}}, image.Point{}, draw.Src)
	}
}

func BenchmarkNRGBA64ImgDrawMaxAlpha(b *testing.B) {
	base := NewNRGBA64(ColorSpace{ColorSpace: "rec2020"}, image.Rect(0, 0, 1000, 1000))
	for n := 0; n < b.N; n++ {
		draw.Draw(base, base.Bounds(), &image.Uniform{&CNRGBA64{R: 0xffff, A: 0xffff, ColorSpace: ColorSpace{ColorSpace: "rec709"}}}, image.Point{}, draw.Src)
	}
}

func TestDraw(t *testing.T) {

	/*
		tests - see draw works the same when no colour colour space is applied

		check transformations of some small squares

	*/

	// check for any deviations from go
	for i := 0; i < 5; i++ {

		baseColour := color.NRGBA64{R: uint16(rand.Int63n(65535)), G: uint16(rand.Int63n(65535)), B: uint16(rand.Int63n(65535)), A: uint16(rand.Int63n(65535))}

		colourImplementation := NewNRGBA64(ColorSpace{ColorSpace: "rec2020"}, image.Rect(0, 0, 20, 20))
		Draw(colourImplementation, colourImplementation.Bounds(), &image.Uniform{baseColour}, image.Point{}, draw.Src)

		goImplementation := NewNRGBA64(ColorSpace{ColorSpace: "rec2020"}, image.Rect(0, 0, 20, 20))
		draw.Draw(goImplementation, goImplementation.Bounds(), &image.Uniform{baseColour}, image.Point{}, draw.Src)

		hnormal := sha256.New()
		htest := sha256.New()
		hnormal.Write(goImplementation.Pix())
		htest.Write(colourImplementation.Pix())

		//td, _ := os.Create("r.png")
		//png.Encode(td, canvas)

		Convey("Checking that the go and colour implementations of draw produce the same result, when no colour space is involved", t, func() {
			Convey(fmt.Sprintf("Run using a colour of %v", baseColour), func() {
				Convey("The hashes of the image are identical", func() {
					So(htest.Sum(nil), ShouldResemble, hnormal.Sum(nil))
				})
			})
		})
	}

	colourImplementation := NewNRGBA64(ColorSpace{}, image.Rect(0, 0, 200, 200))

	colours := []color.Color{
		color.RGBA64{R: 0x7FFF, A: 0x7FFF},
		&CNRGBA64{G: 0x7FFF, A: 0x7FFF},
		&CNRGBA64{B: 0x7FFF, A: 0x7FFF},
		&CNRGBA64{R: 0x6400, G: 0x6400, A: 16384},
		&CNRGBA64{R: 0x8000, B: 0x8000, A: 49151},
		&CNRGBA64{G: 0x6100, B: 0x9900, A: 25000},
	}
	// check for any deviations from go
	for i, baseColour := range colours {

		Draw(colourImplementation, colourImplementation.Bounds(), &image.Uniform{baseColour}, image.Point{}, draw.Over)

		// get the base file
		baseFile, _ := os.Open(fmt.Sprintf("./testdata/draw/alphaDraw%v.png", i))
		baseImage, _ := png.Decode(baseFile)
		base := image.NewNRGBA64(baseImage.Bounds())
		Draw(base, base.Bounds(), baseImage, image.Point{}, draw.Src)

		hnormal := sha256.New()
		htest := sha256.New()
		hnormal.Write(base.Pix)
		htest.Write(colourImplementation.Pix())

		Convey("Checking that the draw function works with alpha", t, func() {
			Convey(fmt.Sprintf("Run adding a colour of %v to the base image", baseColour), func() {
				Convey("The hashes of the image are identical with that of the expected", func() {
					So(htest.Sum(nil), ShouldResemble, hnormal.Sum(nil))
				})
			})
		})
	}

	testColours := []CNRGBA64{{R: 35340, A: 0xffff}, {G: 30000, B: 40000, A: 0xf0f0}, {R: 0xffff, G: 0xffff, B: 0xffff}}

	target := []string{"fullalpha.png", "partialalpha.png", "noalpha.png"}

	for i, tcol := range testColours {
		colourImplementation := NewNRGBA64(ColorSpace{ColorSpace: "rec2020"}, image.Rect(0, 0, 1000, 1000))
		Draw(colourImplementation, colourImplementation.Bounds(), &image.Uniform{&tcol}, image.Point{}, draw.Src)

		baseFile, _ := os.Open("./testdata/draw/" + target[i])
	

		//PngEncode(basePng, colourImplementation.base)
			baseImage, _ := png.Decode(baseFile)

		testFormat := image.NewNRGBA64(baseImage.Bounds())
			Draw(testFormat, testFormat.Bounds(), baseImage, image.Point{}, draw.Src)
		//
		hnormal := sha256.New()
		htest := sha256.New()
			hnormal.Write(testFormat.Pix)
		htest.Write(colourImplementation.Pix())


		Convey("Checking that the transformation produces the expected results", t, func() {
			Convey(fmt.Sprintf("Run checking %v", target[i]), func() {
				Convey("The hashes of the image are identical", func() {
					So(htest.Sum(nil), ShouldResemble, hnormal.Sum(nil))
				})
			})
		})
	}

	/*
		f, _ := os.Create("./testdata/colour.png")
		png.Encode(f, base)

		basedraw := NewNRGBA64(ColorSpace{ColorSpace: "rec2020"}, image.Rect(0, 0, 1000, 1000))

		Draw(basedraw, base.Bounds(), &image.Uniform{color.NRGBA64{R: 0xffff, A: 0xfff0}}, image.Point{}, draw.Over)

		fdraw, _ := os.Create("./testdata/coloudrawr.png")
		png.Encode(fdraw, basedraw)*/

	// set some base test transformations
}

func testrun2() {
	/*

		mkae one image setting a test pattern
	*/

	base := image.NewNRGBA64(image.Rect(0, 0, 2000, 2000))
	noSpace := ColorSpace{ColorSpace: "rec709"}
	noChange := NewNRGBA64(noSpace, image.Rect(0, 0, 1000, 1000))
	b := bar{Space: noSpace}
	b.generate2(noChange)

	changeSpace := ColorSpace{ColorSpace: "rec709"}
	changeYCb := NewNRGBA64(noSpace, image.Rect(0, 0, 1000, 1000))
	cb := bar{Space: changeSpace}
	cb.generateYCbCr(changeYCb)

	change601 := ColorSpace{ColorSpace: "rec601"}
	chang601 := NewNRGBA64(noSpace, image.Rect(0, 0, 1000, 1000))
	cb601 := bar{Space: change601}
	cb601.generate2(chang601)

	change709 := ColorSpace{ColorSpace: "rec601"}
	chang709 := NewNRGBA64(noSpace, image.Rect(0, 0, 1000, 1000))
	cb709 := bar{Space: change709}
	cb709.generateYCbCr(chang709)

	Draw(base, image.Rect(0, 0, 1000, 1000), noChange, image.Point{}, draw.Over)
	Draw(base, image.Rect(1000, 0, 2000, 1000), changeYCb, image.Point{}, draw.Over)
	Draw(base, image.Rect(0, 1000, 1000, 2000), chang601, image.Point{}, draw.Over)
	Draw(base, image.Rect(1000, 1000, 2000, 2000), chang709, image.Point{}, draw.Over)

	f, _ := os.Create("./testdata/all2.png")
	PngEncode(f, base)

}

func testrun2020() {
	/*

		mkae one image setting a test pattern
	*/

	base := image.NewNRGBA64(image.Rect(0, 0, 2000, 2000))
	noSpace := ColorSpace{ColorSpace: "rec2020"}
	noChange := NewNRGBA64(noSpace, image.Rect(0, 0, 1000, 1000))
	b := bar{Space: noSpace}
	b.generate2(noChange)

	changeSpace := ColorSpace{ColorSpace: "rec709"}
	img709 := NewNRGBA64(noSpace, image.Rect(0, 0, 1000, 1000))
	cb := bar{Space: changeSpace}
	cb.generate2(img709)

	change601 := ColorSpace{ColorSpace: "rec601"}
	chang601 := NewNRGBA64(noSpace, image.Rect(0, 0, 1000, 1000))
	cb601 := bar{Space: change601}
	cb601.generate2(chang601)

	change709 := ColorSpace{ColorSpace: "p3"}
	changP3 := NewNRGBA64(noSpace, image.Rect(0, 0, 1000, 1000))
	cb709 := bar{Space: change709}
	cb709.generate2(changP3)

	Draw(base, image.Rect(0, 0, 1000, 1000), noChange, image.Point{}, draw.Over)
	Draw(base, image.Rect(1000, 0, 2000, 1000), img709, image.Point{}, draw.Over)
	Draw(base, image.Rect(0, 1000, 1000, 2000), changP3, image.Point{}, draw.Over)
	Draw(base, image.Rect(1000, 1000, 2000, 2000), chang601, image.Point{}, draw.Over)

	f, _ := os.Create("./testdata/all2020.png")
	PngEncode(f, base)

}

func testrun709() {
	/*

		mkae one image setting a test pattern
	*/

	base := image.NewNRGBA64(image.Rect(0, 0, 2000, 2000))
	noSpace := ColorSpace{ColorSpace: "rec709"}
	noChange := NewNRGBA64(noSpace, image.Rect(0, 0, 1000, 1000))
	b := bar{Space: ColorSpace{ColorSpace: "rec2020"}}
	b.generate2(noChange)

	changeSpace := ColorSpace{ColorSpace: "rec709"}
	img709 := NewNRGBA64(noSpace, image.Rect(0, 0, 1000, 1000))
	cb := bar{Space: changeSpace}
	cb.generate2(img709)

	change601 := ColorSpace{ColorSpace: "rec601"}
	chang601 := NewNRGBA64(noSpace, image.Rect(0, 0, 1000, 1000))
	cb601 := bar{Space: change601}
	cb601.generate2(chang601)

	change709 := ColorSpace{ColorSpace: "p3"}
	changP3 := NewNRGBA64(noSpace, image.Rect(0, 0, 1000, 1000))
	cb709 := bar{Space: change709}
	cb709.generate2(changP3)

	Draw(base, image.Rect(0, 0, 1000, 1000), noChange, image.Point{}, draw.Over)
	Draw(base, image.Rect(1000, 0, 2000, 1000), img709, image.Point{}, draw.Over)
	Draw(base, image.Rect(0, 1000, 1000, 2000), changP3, image.Point{}, draw.Over)
	Draw(base, image.Rect(1000, 1000, 2000, 2000), chang601, image.Point{}, draw.Over)

	f, _ := os.Create("./testdata/all709.png")
	PngEncode(f, base)

}
