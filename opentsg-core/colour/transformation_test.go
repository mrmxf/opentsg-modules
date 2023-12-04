package colour

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"math/rand"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTransfromRT(t *testing.T) {

	for i := 0; i < 10; i++ {
		// base := &CNRGBA64{R: uint16(rand.Int63n(65535)), G: uint16(rand.Int63n(65535)), B: uint16(rand.Int63n(65535)), A: 0xffff}
		R := uint16(rand.Int63n(4095))
		G := uint16(rand.Int63n(4095))
		B := uint16(rand.Int63n(4095))
		base := &CNRGBA64{R: uint16(R << 4), G: uint16(G << 4), B: uint16(B << 4), A: 0xffff, ColorSpace: ColorSpace{ColorSpace: "p3"}}
		fmt.Println(base)
		//	res := transform(ColorSpace{ColorSpace: "p3"}, ColorSpace{ColorSpace: "rec2020"}, base)

		newImage := NewNRGBA64(ColorSpace{ColorSpace: "rec2020"}, image.Rect(0, 0, 1, 1))

		Draw(newImage, newImage.Bounds(), &image.Uniform{base}, image.Point{}, draw.Src)
		f, _ := os.Create("./testdata/draw/test.png")

		PngEncode(f, newImage)
		f.Close()

		f, _ = os.Open("./testdata/draw/test.png")
		img, _ := png.Decode(f)
		testSquare := NewNRGBA64(ColorSpace{ColorSpace: "rec2020"}, image.Rect(0, 0, 1, 1))
		Draw(testSquare, testSquare.Bounds(), img, image.Point{}, draw.Src)

		finalDest := NewNRGBA64(ColorSpace{ColorSpace: "p3"}, image.Rect(0, 0, 1, 1))
		Draw(finalDest, finalDest.Bounds(), testSquare, image.Point{}, draw.Src)

		//ret := transform(ColorSpace{ColorSpace: "rec2020"}, ColorSpace{ColorSpace: "p3"}, res)
		fmt.Println(finalDest.At(0, 0))

		//	gR,gG,gB _ := finalDest.At(0, 0).RGBA

		gR, gG, gB, _ := finalDest.At(0, 0).RGBA()
		fmt.Println(finalDest.At(0, 0).RGBA())
		fmt.Println(R, uint16(math.Round(float64(gR)/16)), G, gG>>4, B, gB>>4)
		//works with a rounding method that i computationally expensve than bit sifting
		// +1 an be accounted for -1 can not as the bytes go in the other direction leading to differences
		//	fmt.Printf("R:%016b, gR:%016b, G:%016b, gG:%016b, b:%016b, gB:%016b\n", R, gR, G, gG, B, gB)
		// fmt.Printf("R:%016b, gR:%016b, G:%016b, gG:%016b, b:%016b, gB:%016b\n", R, gR>>4, G, gG>>4, B, gB>>4)
		Convey("Checking that you can roundtrip colours from different colour spaces", t, func() {
			Convey(fmt.Sprintf("Rountripping %v from rec2020 to rec709", base), func() {
				Convey("The intial 12 bit colour is preserved when translating the colour back to 12 bit from 16 bit", func() {
					So(R, ShouldResemble, uint16(math.Round(float64(gR)/16)))
					So(G, ShouldResemble, uint16(math.Round(float64(gG)/16)))
					So(B, ShouldResemble, uint16(math.Round(float64(gB)/16)))
				})
			})
		})
	}
}
