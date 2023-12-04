// package anglegen converts strings to floats for the opentsg angles
package anglegen

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

// AngleCalc returns the float equivalent of a pi string e.g. π*1/4 or degree values as a string, converted to radians.
// patterns to be used are :
//
// π\*(\d){1,4}/{1}(\d){1,4}$
//
// ^π\*(\d){1,4}$
//
// ^(\d){0,}(\.){0,1}(\d){0,}$
func AngleCalc(ang string) (float64, error) {
	// piUnit := []rune{960} //[]rune{960} this is how the pi image is changed

	/*regex here the types of radian and degree*/
	raidanFrac := regexp.MustCompile(`π\*(\d){1,4}/{1}(\d){1,4}$`)
	radian := regexp.MustCompile(`^π\*(\d){1,4}$`)
	degree := regexp.MustCompile(`^(\d){0,}(\.){0,1}(\d){0,}$`)

	// check angle type
	switch {
	case radian.MatchString(ang) || raidanFrac.MatchString(ang):
		// convert the string to a float64
		rad := stringTofraction(ang[3:])

		return (math.Pi * rad), nil
	case degree.MatchString(ang):
		degree, err := strconv.ParseFloat(ang, 64)

		return (math.Pi * (degree / 180)), err
	default:

		return 0, fmt.Errorf("%s is not a valid angle", ang)
	}
}

func stringTofraction(form string) float64 {
	var angle float64
	var pos int
	div := false

	for i, c := range form {
		if c == 47 { // 47 is "/"
			div = true
			pos = i
		}
	}

	if div {
		num, _ := strconv.ParseFloat(form[:pos], 64)
		dom, _ := strconv.ParseFloat(form[pos+1:], 64)
		angle = num / dom

	} else {
		angle, _ = strconv.ParseFloat(form, 64)
	}

	return angle
}
