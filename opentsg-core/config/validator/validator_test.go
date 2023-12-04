package validator

import (
	"fmt"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSeveralErrrors(t *testing.T) {
	// make a basic json open it wth liner
	// run a schema validator on it to fail

	fakeSon, _ := os.ReadFile("./testdata/jsonlines/badstripe.json")
	fakeSchema, _ := os.ReadFile("./testdata/jsonlines/stripeschema.json")

	// generate a basic map with the fake json
	d := make(map[uint64]fileAndLocation)
	err := Liner(fakeSon, "./testdata/jsonlines/badstripe.json", "test", d)

	results := SchemaValidator(fakeSchema, fakeSon, "testsuite", d)
	expec := []error{fmt.Errorf("0026 Additional property bad is not allowed at line 74 in ./testdata/jsonlines/badstripe.json"),
		fmt.Errorf("0026 stripes.groupHeader.color.0 must be one of the following: \"red\", \"green\", \"blue\", \"black\", \"white\", \"gray\", \"grey\" at line 10 in ./testdata/jsonlines/badstripe.json"), fmt.Errorf("0026 Must be less than or equal to 4000 at line 12 in ./testdata/jsonlines/badstripe.json")}
	Convey("Checking that errors are caught and their line position is returned", t, func() {
		Convey("using three different errors", func() {
			Convey(fmt.Sprintf("Errors of %v are returned", expec), func() {
				So(err, ShouldBeNil)
				So(results, shouldBeInAnyOrder, expec)
			})
		})
	})
}

var madeFile = "error.json"

func TestTypes(t *testing.T) {
	fakeSchema, _ := os.ReadFile("./testdata/jsonlines/stripeschema.json")

	extras := []string{`, "text":"not real"`, `,"text":["fake"]`, `, "text":["fake", 6, "whatever" ]`,
		`,"text":1034`, `,"text":103.4`, `,"text":[{"surprise":"map"}, 1034, "3"]`}
	expected := []string{", given: string at line 68", ", given: array at line 68", ", given: array at line 68",
		", given: integer at line 68", ", given: number at line 68", ", given: array at line 68", "given: array at line 68"}
	for i, e := range extras {

		lines := make(JSONLines)
		badBytes := addWidget(madeFile, e, lines) // make a file with the error to be caught

		results := SchemaValidator(fakeSchema, badBytes, "testsuite", lines)

		Convey("Checking that errors are caught and their line position is returned for different field types that don't match the expected", t, func() {
			Convey("Using an update of"+e, func() {
				Convey(fmt.Sprintf("An error of %v is returned", expected[i]), func() {
					So(results, ShouldResemble, []error{fmt.Errorf("0026 Invalid type. Expected: object%s in error.json", expected[i])})
				})
			})
		})
	}

	extraProps := []string{`,"text":{"surprise":5634}`, `,"text":{"surprise":{"lowermap":5634}}`,
		`,"text":{"surprise":[3,"hello"]}`, `,"text":{"surprise":"big surprise"}`, `,"text": {"surprise": [3, {"surprise": {"map": 12}}]}`}
	mapErr := []error{fmt.Errorf("0026 Additional property surprise is not allowed at line 68 in error.json")}

	for _, e := range extraProps {

		lines := make(JSONLines)
		badBytes := addWidget(madeFile, e, lines) // make a file with the error to be caught

		results := SchemaValidator(fakeSchema, badBytes, "testsuite", lines)

		Convey("Checking that errors are caught and their line position is returned for different field types within an object with a map", t, func() {
			Convey("Using an update of "+e, func() {
				Convey(fmt.Sprintf("An error of %v is returned", mapErr), func() {
					So(results, ShouldResemble, mapErr)
				})
			})
		})
	}

	// check double updates work

	doubleWidget := []string{`,"text":"bad property"`}
	doubleLoader := []string{`,"loader": {"text":"bad property"}`}
	doubleResult := []error{fmt.Errorf("0026 Invalid type. Expected: object, given: string at line 68,11 in error.json,loader.json")}
	for i, e := range doubleWidget {
		lines := make(JSONLines)
		badBytes := addWidget(madeFile, e, lines)         // make a file with the error to be caught
		addFactory("loader.json", doubleLoader[i], lines) // mock the loader with the same mistake

		results := SchemaValidator(fakeSchema, badBytes, "testsuite", lines)

		Convey("Checking that errors are caught and their line positions are returned for when two files have the same error within them", t, func() {
			Convey(fmt.Sprintf("Using updates of %v and %v ", e, doubleLoader[i]), func() {
				Convey(fmt.Sprintf("An error of %v is returned highlighting both files involved", doubleResult), func() {
					So(results, ShouldResemble, doubleResult)
				})
			})
		})
	}

}

func TestYaml(t *testing.T) {
	// testing yaml sources are readinto go json schema
	fakeSon, _ := os.ReadFile("./testdata/jsonlines/badstripe.yaml")
	fakeSchema, _ := os.ReadFile("./testdata/jsonlines/stripeschema.json")

	// generate a basic map with the fake json
	d := make(map[uint64]fileAndLocation)
	err := Liner(fakeSon, "./testdata/jsonlines/badstripe.json", "test", d)

	results := SchemaValidator(fakeSchema, fakeSon, "testsuite", d)
	expec := []error{fmt.Errorf("0026 stripes.groupHeader.color.0 must be one of the following: \"red\", \"green\", \"blue\", \"black\", \"white\", \"gray\", \"grey\" at line 9 in ./testdata/jsonlines/badstripe.json")}

	Convey("Checking that errors are caught and their line position is returned", t, func() {
		Convey("using three different errors", func() {
			Convey(fmt.Sprintf("Errors of %v are returned", expec), func() {
				So(err, ShouldBeNil)
				So(results, ShouldResemble, expec)
			})
		})
	})
}
func TestDoubleSource(t *testing.T) {
	fakeSchema, _ := os.ReadFile("./testdata/jsonlines/stripeschema.json")

	extraAdd := `,"bad":{
		"surprise":5634,
		 "some":"value"}`
	factroyAdd := `,"loader": {"bad":{"additional":"problem"}}`
	complete := `,"bad":{
		"surprise":5634,
		 "some":"value",
		 "additional":"problem"}`

	errorLocations := [][]error{{fmt.Errorf("0026  Additional property bad is not allowed at line 70 in error.json")},
		{fmt.Errorf("0026  Additional property bad is not allowed at line 69 in error.json")},
		{fmt.Errorf("0026  Additional property bad is not allowed at line 11 in loader.json")}}

	lines := make(JSONLines)
	addWidget(madeFile, extraAdd, lines)         // make a file with the error to be caught
	addFactory("loader.json", factroyAdd, lines) // mock the loader with the same mistake

	genJSON := addWidget(madeFile, complete, make(JSONLines)) // mock the generated json

	results := SchemaValidator(fakeSchema, genJSON, "testsuite", lines)

	Convey("Checking that errors are caught and that the unknown file is returned for when a bad property is made up of multiple fields", t, func() {
		Convey(fmt.Sprintf("Using updates of %v and %v in the json and loader respectively", extraAdd, factroyAdd), func() {
			Convey(fmt.Sprintf("An error is found to be one of %v ", errorLocations), func() {
				So(results, ShouldBeIn, errorLocations)
			})
		})
	})

	/*repeat the addwidget several tiems over*/
	repeatLines := make(JSONLines)
	addWidget(madeFile, extraAdd, repeatLines)
	addWidget(madeFile, extraAdd, repeatLines)
	addWidget(madeFile, extraAdd, repeatLines)
	badBytes := addWidget(madeFile, extraAdd, repeatLines) // make a file with the error to be caught

	results = SchemaValidator(fakeSchema, badBytes, "testsuite", repeatLines)
	errorMult := []error{fmt.Errorf("0026 Additional property bad is not allowed at line 68 in error.json")}

	Convey("Checking that multiple additions of the same file and line don't lead to repeated files", t, func() {
		Convey(fmt.Sprintf("Using updates of %v in the json", extraAdd), func() {
			Convey(fmt.Sprintf("An error is found to be %v where the file has not been repeated", errorMult), func() {
				So(results, ShouldResemble, errorMult)
			})
		})
	})

}

// cont mocker generates the context from the input file that is generated by the user
// the file follows the template of base
// this mocks the draw set up and allows us to run the widget
func addWidget(file, extra string, positions JSONLines) []byte {
	// make the stripes a passable schema
	complete := fmt.Sprintf(mockPlain, extra)
	_ = Liner([]byte(complete), file, "widget", positions)

	return []byte(complete)
}

func addFactory(file, extra string, positions JSONLines) {
	// make the stripes a passable schema
	complete := fmt.Sprintf(mockBase, extra)
	_ = Liner([]byte(complete), file, "factory", positions)
}

func shouldBeInAnyOrder(actual interface{}, expected ...interface{}) string {
	if len(expected) != 1 {
		return fmt.Sprintf("got %v expected only 1 object is allowed", len(expected))
	}
	actualArray, ok := actual.([]error)
	expectedArray, oke := expected[0].([]error)

	if !ok || !oke {
		return "both items are required to be an arrays"
	} else if len(actualArray) != len(expectedArray) {
		return fmt.Sprintf("mismatch received an array of length %v  and expected one of length %v", len(actualArray), len(expected))
	}

	for _, a := range actualArray {

		var match bool
		for j := 0; j < len(actualArray); j++ {
			// just compare the error string at the moment
			if a.Error() == expectedArray[j].Error() {
				match = true

				break
			}
		}
		if !match {
			return fmt.Sprintf("%v was not found in the expected array", a)
		}
	}

	return ""
}

var mockBase = `{
    "include": [
      {
        "uri": "notparsedfile.json",
        "name": "load"
      }
    ],
    "create": [
      {
        "load": {}
		%s
      }
    ]
  }`

var mockPlain = `{
    "type": "builtin.ramps",
    "minimum": 0,
    "maximum": 4095,
    "depth": 12,
    "fillType": "fill",
    "stripes": {
        "groupHeader": {
            "color": [
                "white"
            ],
            "height": 45
        },
        "interstripes": {
            "color": [
                "black"
            ],
            "height": 1
        },
        "ramps": {
            "fill": "gradient",
            "bitdepth": [
                10,8
            ],
            "labels": [
                "10b",
                "8b"
            ],
            "height": 10,
            "rampGroups": {
                "ared": {
                    "color": "red",
                    "rampstart": 0,
                    "direction": 1
                },
                "bgreen": {
                    "color": "green",
                    "rampstart": 0,
                    "direction": 1
                },
                "cblue": {
                    "color": "blue",
                    "rampstart": 0,
                    "direction": 1
                },
                "dred": {
                    "color": "red",
                    "rampstart": 0,
                    "direction": -1
                },
                "egreen": {
                    "color": "green",
                    "rampstart": 0,
                    "direction": -1
                },
                "fblue": {
                    "color": "blue",
                    "rampstart": 0,
                    "direction": -1
                }
            }
        }
    },
    "grid": {
        "location": "a5:p8",
        "alias": "bottom"
    }
	%s
}`
