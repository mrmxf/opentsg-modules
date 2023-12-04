package anglegen

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGoodInputs(t *testing.T) {
	goodStrings := []string{"π*2/3", "π*1"}
	ExpectedRadianResult := []float64{2.0943951023931953, 3.141592653589793}

	for i, angle := range goodStrings {

		res, err := AngleCalc(angle)

		Convey("Checking radians are calculated correctly", t, func() {
			Convey(fmt.Sprintf("using an angle of %v", angle), func() {
				Convey(fmt.Sprintf("No error is returned and the calculated angle is %v", ExpectedRadianResult[i]), func() {
					So(err, ShouldBeNil)
					So(res, ShouldResemble, ExpectedRadianResult[i])
				})
			})
		})
	}

	goodDegree := []string{"45", "180"}
	ExpectedDegreeResult := []float64{0.7853981633974483, 3.141592653589793}

	for i, angle := range goodDegree {

		res, err := AngleCalc(angle)

		Convey("Checking radians are calculated correctly", t, func() {
			Convey(fmt.Sprintf("using an angle of %v", angle), func() {
				Convey(fmt.Sprintf("No error is returned and the calculated angle is %v", ExpectedDegreeResult[i]), func() {
					So(err, ShouldBeNil)
					So(res, ShouldResemble, ExpectedDegreeResult[i])
				})
			})
		})
	}

}

func TestBadInputs(t *testing.T) {
	goodDegree := []string{"4.5.0", "thirtyfive"}

	for _, angle := range goodDegree {

		res, err := AngleCalc(angle)

		Convey("Checking incorrect strings are caught", t, func() {
			Convey(fmt.Sprintf("using an angle of %v", angle), func() {
				Convey("An error is returned highlighting the angle in question", func() {
					So(err, ShouldResemble, fmt.Errorf("%s is not a valid angle", angle))
					So(res, ShouldEqual, 0)
				})
			})
		})
	}

}
