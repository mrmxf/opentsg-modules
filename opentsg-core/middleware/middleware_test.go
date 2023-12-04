package middleware

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

/*todo wait for middleware to be readded in future updaes*/
func TestMinMaxFrame(t *testing.T) { 
	mockFunctions := make([]Functions, 1)
	MinMax := [][2]int{{0, 9}, {3, 10}, {4, 9999}}

	for _, mm := range MinMax {

		mockFunctions[0].AnyOf = []string{fmt.Sprintf("framecount(%v,%v)", mm[0], mm[1])}
		err, _, _ := Check(mockFunctions, 8)

		Convey("Checking an exisiting file is read", t, func() {
			Convey("using a ./testdata/apitest.json as the input file", func() {
				Convey("No error is returned", func() {
					So(err, ShouldBeTrue)
				})
			})
		})
	}

	MinMaxBad := [][2]int{{5, 3}, {13, 10}, {43533, 9999}}

	for _, mm := range MinMaxBad {

		mockFunctions[0].AnyOf = []string{fmt.Sprintf("framecount(%v,%v)", mm[0], mm[1])}
		err, _, _ := Check(mockFunctions, 8)

		Convey("Checking an exisiting file is read", t, func() {
			Convey("using a ./testdata/apitest.json as the input file", func() {
				Convey("No error is returned", func() {
					So(err, ShouldBeFalse)
				})
			})
		})
	}
}
