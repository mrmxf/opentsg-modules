package bench

import (
	"encoding/json"
	"image/color"
	"image/draw"
	"os"
	"reflect"
)

func tester(box draw.Image, colors map[int]color.Color) {

	for y := 0; y < box.Bounds().Max.Y; y++ {
		c := colors[y%4]
		for x := 0; x < box.Bounds().Max.X; x++ {
			box.Set(x, y, c)
		}
	}
}

func mapTest() bool {
	b, err := os.ReadFile("my.json")
	if err != nil {
		return false
	}

	body := make(map[string]interface{})
	err = json.Unmarshal(b, &body)
	if err != nil {
		return false
	}

	names := make([]string, len(body))
	count := 0
	for k := range body {
		names[count] = k
		count++
	}

	base := []string{"picture", "eyeColor", "friends", "favoriteFruit", "_id", "guid", "isActive", "email", "address", "about", "latitude", "tags", "age", "name", "gender", "phone", "longitude", "greeting", "index", "balance", "company", "registered"}

	// check if they match
	return reflect.DeepEqual(names, base)
}

/*

first run

BenchmarkNRGBA64-16         1692          27731194 ns/op
BenchmarkACES-16            1899          24619949 ns/op


second run

BenchmarkNRGBA64-16                         1885          25632873 ns/op
BenchmarkNRGBA64ACESColour-16                600          79332951 ns/op
BenchmarkACES-16                            2102          22791580 ns/op
BenchmarkACESNRGBAColour-16                  556          85354098 ns/op

third run

BenchmarkNRGBA64-16                         2048          23647763 ns/op
BenchmarkNRGBA64ACESColour-16                643          75871906 ns/op
BenchmarkACES-16                            2142          22323483 ns/op
BenchmarkACESNRGBAColour-16                  507          85208727 ns/op

fourth run where nrgba64 speed up has been removed
BenchmarkNRGBA64-16                         1814          24401593 ns/op
BenchmarkNRGBA64ACESColour-16                618          86061670 ns/op
BenchmarkACES-16                            1920          27455221 ns/op
BenchmarkACESNRGBAColour-16                  506          90274354 ns/op
*/
