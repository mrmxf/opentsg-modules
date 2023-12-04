package colour

import (
	"image"
	"image/color"
	"image/draw"
)

type bar struct {
	Space ColorSpace
}

type bars struct {
	width float64
	color Color
}

var ( /*
		gray40   = color.NRGBA64{R: 414 << 6, G: 414 << 6, B: 414 << 6, A: 0xffff}
		white75  = color.NRGBA64{R: 721 << 6, G: 721 << 6, B: 721 << 6, A: 0xffff}
		yellow75 = color.NRGBA64{R: 721 << 6, G: 721 << 6, B: 64 << 6, A: 0xffff}
		cyan75   = color.NRGBA64{R: 64 << 6, G: 721 << 6, B: 721 << 6, A: 0xffff}
		green75  = color.NRGBA64{R: 64 << 6, G: 721 << 6, B: 64 << 6, A: 0xffff}
		mag75    = color.NRGBA64{R: 721 << 6, G: 64 << 6, B: 721 << 6, A: 0xffff}
		red75    = color.NRGBA64{R: 721 << 6, G: 64 << 6, B: 64 << 6, A: 0xffff}
		blue75   = color.NRGBA64{R: 64 << 6, G: 64 << 6, B: 721 << 6, A: 0xffff}*/

	gray40   = &CNRGBA64{R: 414 << 6, G: 414 << 6, B: 414 << 6, A: 0xffff}
	white75  = &CNRGBA64{R: 721 << 6, G: 721 << 6, B: 721 << 6, A: 0xffff}
	yellow75 = &CNRGBA64{R: 721 << 6, G: 721 << 6, B: 64 << 6, A: 0xffff}
	cyan75   = &CNRGBA64{R: 64 << 6, G: 721 << 6, B: 721 << 6, A: 0xffff}
	green75  = &CNRGBA64{R: 64 << 6, G: 721 << 6, B: 64 << 6, A: 0xffff}
	mag75    = &CNRGBA64{R: 721 << 6, G: 64 << 6, B: 721 << 6, A: 0xffff}
	red75    = &CNRGBA64{R: 721 << 6, G: 64 << 6, B: 64 << 6, A: 0xffff}
	blue75   = &CNRGBA64{R: 64 << 6, G: 64 << 6, B: 721 << 6, A: 0xffff}

	gray40YCbCr   = &CyCbCr{Y: 104, Cb: 128, Cr: 128}
	white75YCbCr  = &CyCbCr{Y: 180, Cb: 128, Cr: 128}
	yellow75YCbCr = &CyCbCr{Y: 168, Cb: 44, Cr: 136}
	cyan75YCbCr   = &CyCbCr{Y: 145, Cb: 147, Cr: 44}
	green75YCbCr  = &CyCbCr{Y: 133, Cb: 63, Cr: 52}
	mag75YCbCr    = &CyCbCr{Y: 63, Cb: 193, Cr: 204}
	red75YCbCr    = &CyCbCr{Y: 51, Cb: 109, Cr: 212}
	blue75YCbCr   = &CyCbCr{Y: 28, Cb: 212, Cr: 120}
)

const (
	//widths
	d = 240 / 1920.0
	f = 205 / 1920.0
	c = 206 / 1920.0
	b = 1 / 12.0
	k = 309 / 1920.0
	g = 411 / 1920.0
	h = 171 / 1920.0
)

/*
func (br bar) generate(canvas Image) {
	b := canvas.Bounds().Max
	w := 0.0
	twidth := 0.0

	fills := []bars{{width: d, color: gray40}, {width: f, color: white75}, {width: c, color: yellow75}, {width: c, color: cyan75}, {width: c, color: green75}, {width: c, color: mag75}, {width: c, color: red75}, {width: f, color: blue75}, {width: d, color: gray40}}

	for _, f := range fills {
		twidth += f.width * float64(b.X)
		area := image.Rect(int(w), int(0), int(w+f.width*float64(b.X)), b.Y)

		canvas.Draw(area, f.color, draw.Src, br.Space)
		//draw.Draw(canvas, area, fill, image.Point{}, draw.Src)

		w += f.width * float64(b.X)
	}

}*/

func (br bar) generate2(canvas Image) {
	b := canvas.Bounds().Max
	w := 0.0
	twidth := 0.0

	fills := []bars{{width: d, color: gray40}, {width: f, color: white75}, {width: c, color: yellow75}, {width: c, color: cyan75}, {width: c, color: green75}, {width: c, color: mag75}, {width: c, color: red75}, {width: f, color: blue75}, {width: d, color: gray40}}

	for _, f := range fills {
		twidth += f.width * float64(b.X)
		area := image.Rect(int(w), int(0), int(w+f.width*float64(b.X)), b.Y)
		fill := f.color
		fill.UpdateColorSpace(br.Space)
		/*	fmt.Println(fill)
			fmt.Println(fill.RGBA())
			R, G, B, _ := fill.RGBA()
			fmt.Println(uint8(R>>8), uint8(G>>8), uint8(B>>8), R>>8, G, B)
			Y, cb, cr := RGBToYCbCr(uint8(R>>8), uint8(G>>8), uint8(B>>8))
			fmt.Println(RGBToYCbCr(uint8(R>>8), uint8(G>>8), uint8(B>>8)))
			fmt.Println(Y, cb, cr)
			fmt.Println(color.YCbCrToRGB(uint8(math.Round(Y)), uint8(math.Round(cb)), uint8(math.Round(cr))))
			fmt.Println(YCbCrToRGB(Y, cb, cr))*/
		// canvas.Draw(area, f.color, draw.Src, br.Space)
		Draw(canvas, area, &image.Uniform{fill}, image.Point{}, draw.Src)

		w += f.width * float64(b.X)
	}

}

func (br bar) generateYCbCr(canvas Image) {
	b := canvas.Bounds().Max
	w := 0.0
	twidth := 0.0

	// fills := []bars{{width: d, color: gray40YCbCr}, {width: f, color: white75YCbCr}, {width: c, color: yellow75YCbCr}, {width: c, color: cyan75YCbCr}, {width: c, color: green75YCbCr}, {width: c, color: mag75YCbCr}, {width: c, color: red75YCbCr}, {width: f, color: blue75YCbCr}, {width: d, color: gray40YCbCr}}
	fills := []bars{{width: d, color: gray40}, {width: f, color: white75}, {width: c, color: yellow75}, {width: c, color: cyan75}, {width: c, color: green75}, {width: c, color: mag75}, {width: c, color: red75}, {width: f, color: blue75}, {width: d, color: gray40}}

	for _, f := range fills {
		twidth += f.width * float64(b.X)
		area := image.Rect(int(w), int(0), int(w+f.width*float64(b.X)), b.Y)

		R, G, B, _ := f.color.RGBA()
		//	fmt.Println(uint8(R>>8), uint8(G>>8), uint8(B>>8))
		Y, cb, cr := color.RGBToYCbCr(uint8(R>>8), uint8(G>>8), uint8(B>>8))
		//	fmt.Println(Y, cb, cr)
		//	fmt.Println(fill.RGBA())

		fill := &CyCbCr{Y: Y, Cb: cb, Cr: cr, Space: br.Space}

		// canvas.Draw(area, f.color, draw.Src, br.Space)
		Draw(canvas, area, &image.Uniform{fill}, image.Point{}, draw.Src)

		w += f.width * float64(b.X)
	}

}
