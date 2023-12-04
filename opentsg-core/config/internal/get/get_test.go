package get

/*check the keys and any results
parse some things as json then compare the results to the expected

run for the true/false with the different types of maps each time

go for should bein due to the way maps can be ordered
*/

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestArrayWithMap(t *testing.T) {

	body, _ := os.ReadFile("./testdata/arraymaps.json")
	base := make(map[string]any)

	err := json.Unmarshal(body, &base)

	res, key := orderer(Get(base, []string{}, true))

	expecKey := [][]string{{"arrayWithMap", "0"}, {"arrayWithMap", "1"}, {"arrayWithMap", "2", "map", "with"}, {"base"}, {"map", "with", "nested"}}
	expecRes := []interface{}{float64(3), float64(3), "layers", "base", "values"}

	Convey("Checking that all keys and values are extracted", t, func() {
		Convey("using ./testdata/arraymaps.json as the input file then running get with array dotpath", func() {
			Convey(fmt.Sprintf("The results match %v and the keys match %v", expecRes, expecKey), func() {
				So(err, ShouldBeNil)
				So(key, ShouldResemble, expecKey)
				So(res, ShouldResemble, expecRes)
			})
		})
	})

	res, key = orderer(Get(base, []string{}, false))

	expecKey = [][]string{{"arrayWithMap"}, {"base"}, {"map", "with", "nested"}}
	expecRes = []interface{}{[]interface{}{float64(3), float64(3), map[string]interface{}{"map": map[string]interface{}{"with": "layers"}}}, "base", "values"}

	Convey("Checking that all keys and values are extracted", t, func() {
		Convey("using ./testdata/arraymaps.json as the input file then running get without the array dotpath", func() {
			Convey(fmt.Sprintf("The results match %v and the keys match %v", expecRes, expecKey), func() {
				So(key, ShouldResemble, expecKey)
				So(res, ShouldResemble, expecRes)
			})
		})
	})
}

func TestArray(t *testing.T) {

	body, _ := os.ReadFile("./testdata/array.json")
	base := make(map[string]any)

	err := json.Unmarshal(body, &base)

	res, key := orderer(Get(base, []string{}, true))

	expecKey := [][]string{{"arrayWithMap", "0"}, {"arrayWithMap", "1"}, {"arrayWithMap", "2"}, {"base"}, {"map", "with", "nested"}}
	expecRes := []interface{}{float64(3), float64(3), "no map here", "base", "values"}

	Convey("Checking that all keys and values are extracted", t, func() {
		Convey("using ./testdata/array.json as the input file then running get with array dotpath", func() {
			Convey(fmt.Sprintf("The results match %v and the keys match %v", expecRes, expecKey), func() {
				So(err, ShouldBeNil)
				So(key, ShouldResemble, expecKey)
				So(res, ShouldResemble, expecRes)
			})
		})
	})

	res, key = orderer(Get(base, []string{}, false))

	expecKey = [][]string{{"arrayWithMap"}, {"base"}, {"map", "with", "nested"}}
	expecRes = []interface{}{[]interface{}{float64(3), float64(3), "no map here"}, "base", "values"}

	Convey("Checking that all keys and values are extracted", t, func() {
		Convey("using ./testdata/array.json as the input file then running get without the array dotpath", func() {
			Convey(fmt.Sprintf("The results match %v and the keys match %v", expecRes, expecKey), func() {
				So(key, ShouldResemble, expecKey)
				So(res, ShouldResemble, expecRes)
			})
		})
	})
}

// orderer organises the results and values in order of the first string of every 2d value
// used for keeping the tests identical
func orderer(v []any, s [][]string) ([]any, [][]string) {
	// sort everything in order

	for i := range s {

		for j := 0; j < len(s)-1; j++ {
			// check alphabetical
			if s[i][0][0] < s[j][0][0] {
				s[j], s[i] = s[i], s[j]
				v[j], v[i] = v[i], v[j]
				// else for ordering the array maps with more than one item
				// and not checking themselves
			} else if s[i][0][0] == s[j][0][0] && i != j {
				if s[i][1][0] < s[j][1][0] {
					s[j], s[i] = s[i], s[j]
					v[j], v[i] = v[i], v[j]

				}
			}
		}
	}

	return v, s
}
