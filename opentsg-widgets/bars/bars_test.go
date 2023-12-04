package bars

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
	// . "github.com/smartystreets/goconvey/convey"
)

func TestTemp(t *testing.T) {
	heights := []float64{7 * b, 1 * b, 1 * b, 3 * b}
	fills := [][]bars{{{width: d, color: gray40}, {width: f, color: white75}, {width: c, color: yellow75}, {width: c, color: cyan75}, {width: c, color: green75}, {width: c, color: mag75}, {width: c, color: red75}, {width: f, color: blue75}, {width: d, color: gray40}},
		{{width: d, color: cyan100}, {width: f, color: white100}, {width: 5*c + f, color: white75}, {width: d, color: blue100}},
		{{width: d, color: yellow100}, {width: f, color: black0}, {width: 5 * c, fill: yRamp}, {width: f, color: white100}, {width: d, color: red100}},
		{{width: d, color: gray15}, {width: k, fill: superBlack}, {width: g, fill: superWhite}, {width: h, color: black0}, {width: (69 / 1920.0), color: black2Neg}, {width: (68 / 1920.0), color: black0}, {width: (69 / 1920.0), color: black2Pos}, {width: (68 / 1920.0), color: black0}, {width: (69 / 1920.0), color: black4Pos}, {width: c, color: black0}, {width: d, color: gray15}}}

	HD := image.NewNRGBA64(image.Rect(0, 0, 1920, 1080))
	bar(HD, fills, heights)

	fe, _ := os.Create("test.png")
	png.Encode(fe, HD)

	fmt.Println(color.YCbCrToRGB(219, 16, 138))

	fillsSD := [][]bars{{{width: sdwidth, color: white75SD}, {width: sdwidth, color: yellow75SD}, {width: sdwidth, color: cyan75SD}, {width: sdwidth, color: green75SD}, {width: sdwidth, color: mag75SD}, {width: sdwidth, color: red75SD}, {width: sdwidth, color: blue75SD}},
		{{width: sdwidth, color: blue75SD}, {width: sdwidth, color: black0SD}, {width: sdwidth, color: mag75SD}, {width: sdwidth, color: black0SD}, {width: sdwidth, color: cyan75SD}, {width: sdwidth, color: black0}, {width: sdwidth, color: white75SD}},
		{{width: qwidth, color: ISD}, {width: qwidth, color: white100SD}, {width: qwidth, color: QSD}, {width: qwidth, color: black0SD}, {width: (sdwidth / 3.0), color: black4NegSD}, {width: (sdwidth / 3.0), color: black0SD}, {width: (sdwidth / 3.0), color: black2PosSD}, {width: sdwidth, color: black0SD}}}

	heightsSD := []float64{0.67, 0.08, 0.25}

	SD := image.NewNRGBA(image.Rect(0, 0, 720, 480))
	bar(SD, fillsSD, heightsSD)

	fmt.Println(SD.At(234, 300))

	fsd, _ := os.Create("testSD.png")
	png.Encode(fsd, SD)

}
