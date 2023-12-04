package zoneplate

// This package is for generating lines of different interpolation, it is not currently used in this test card
/*
func horizontal(canvas *image.NRGBA64, s scaleint) [][]color.NRGBA64 {

	canvasArray := make([][]color.NRGBA64, 0)
	// Go y then x directions
	for i := 0; i < 500; i++ {
		fillLine := make([]color.NRGBA64, 0)
		fillrgb := 0

		barcount := 0
		barfill := s.scale(0)
		for j := 0; j < 1000; j++ {
			fill := color.NRGBA64{uint16(fillrgb) << 4, uint16(fillrgb) << 4, uint16(fillrgb) << 4, uint16(0xffff)}

			fillLine = append(fillLine, fill)

			// When it has looped so that the widths match the one stipulate by the scale
			barcount++
			// Alternate between 12 bit black and white
			if barcount == barfill {
				if fillrgb == 0 {
					fillrgb = fillrgb + 4095
				} else {
					fillrgb = fillrgb - 4095
				}
				// Reset bar values to the new scale
				barcount = 0
				barfill = s.scale(j)
			}
		}
		canvasArray = append(canvasArray, [][]color.NRGBA64{fillLine}...)

	}
	return canvasArray
}

func threeBox(canvas *image.NRGBA64) {
	// Generate 3 boxes in arrays and append them every width along by calling the function
	// Then draw and save
	// Can make swap an interface thst takes the x,y and returns a value, make it easier for constants and functions together

	y := constOne(1)
	x := linear(5) // Put an if function here that can choose the type
	z := constTwenty(20)
	swap := []scaleint{&y, &x, &z}
	for i := 0; i < 3; i++ {
		// Move the x along 1000 pieces

		canvasArray := horizontal(canvas, swap[i])
		canvas = canvasmaker.CanvasFill(canvas, canvasArray, i*1000, 0)
	}
}

// Different types of gradient

type constOne int
type linear int
type log int
type exp int
type constTwenty int

func (l *constOne) scale(x int) int {
	return 1
}

func (l *log) scale(x int) int {
	// E.0 is 1
	// Ln(k) = 20 k is 485165195.4
	k := 485165.1954
	// Stops log 0 error
	if x == 0 {
		return 1
	}

	pos := math.Floor(math.Log(k * float64(x)))
	return int(pos)
}

func (l *linear) scale(x int) int {
	pos := math.Floor((float64(x) / 1000) * 19)
	return int(1 + pos)
}

func (e *exp) scale(x int) int {
	// E.0 is 1
	// E.1000 is 20
	// Log20 is 2.99573
	k := 0.00299573
	pos := math.Pow(math.E, (k * float64(x)))

	return int(pos)
}

func (l *constTwenty) scale(x int) int {
	return 20
}

type scaleint interface {
	scale(x int) int
}

*/
