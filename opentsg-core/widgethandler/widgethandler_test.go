package widgethandler

import (
	"context"
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/mrmxf/opentsg-modules/opentsg-core/canvaswidget"
	"github.com/mrmxf/opentsg-modules/opentsg-core/colour"
	"github.com/mrmxf/opentsg-modules/opentsg-core/colourgen"
	"github.com/mrmxf/opentsg-modules/opentsg-core/config/core"
	errhandle "github.com/mrmxf/opentsg-modules/opentsg-core/errHandle"
	"github.com/mrmxf/opentsg-modules/opentsg-core/gridgen"
	. "github.com/smartystreets/goconvey/convey"
)

var mockSchema = []byte(`{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"$id": "https://example.com/product.schema.json",
	"title": "Allow anything through for tests",
	"description": "An empty schema to allow custom structs to run through",
	"type": "object"
	}`)

/*
var tightSchema = []byte(`{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"$id": "https://example.com/product.schema.json",
	"title": "Allow anything through for tests",
	"description": "An empty schema to allow custom structs to run through",
	"type": "object",
	"properties": {
		"type": {
			"type": "string",
			"enum": ["test"]}
	},
	"required": ["type"]
}`)*/

var base = `{
    "include": [
      {
        "uri": "%s",
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

// cont mocker generates the context from the input file that is generated by the user
// the file follows the template of base
// this mocks the draw set up and allows us to run the widget
func contMocker(file, target, extra string) (*context.Context, draw.Image) {
	body := fmt.Sprintf(base, target, extra)
	mockFile, _ := os.Create(file)
	_, _ = mockFile.Write([]byte(body))
	mockFile.Close()

	c, _, _ := core.FileImport(file, "", false)

	cFrame, _ := core.FrameWidgetsGenerator(c, 0, false)

	mockC := MetaDataInit(cFrame)
	canvaswidget.LoopInit(mockC)

	canvas, _ := gridgen.GridGen(mockC)

	return mockC, canvas
}

func TestWidgetRun(t *testing.T) {

	// generate all the context items
	fileName := "./testdata/mockcont/base.json"

	mockC, canvas := contMocker(fileName, "./collater.json", "")

	mLog := errhandle.LogInit("stdout", "")

	mockCanvas := GenConf[canvaswidget.ConfigVals]{true, mockSchema, "builtin.canvasoptions", nil}
	mockConfig := GenConf[test]{true, mockSchema, "mocktest", nil}

	canvasChan := make(chan draw.Image, 1)
	canvasChan <- canvas

	// run the functions
	var wg sync.WaitGroup
	wg.Add(2)
	// Replicate the context that would normally be used
	go WidgetRunner(canvasChan, mockCanvas, mockC, mLog, &wg)
	go WidgetRunner(canvasChan, mockConfig, mockC, mLog, &wg)
	time.Sleep(time.Second) // mock the external wait function

	// open the test image and extract its pixels
	file, _ := os.Open("./testdata/mockcont/base.png")
	// decode to get the colour values
	baseVals, _ := png.Decode(file)

	// assign the colour to the correct type of image NGRBA64 and replace the colour values
	readImage := image.NewNRGBA64(baseVals.Bounds())
	colour.Draw(readImage, readImage.Bounds(), baseVals, image.Point{}, draw.Src)

	hnormal := sha256.New()
	htest := sha256.New()
	hnormal.Write(readImage.Pix)
	htest.Write(canvas.(*colour.NRGBA64).Pix())

	// td, _ := os.Create("r.png")
	// png.Encode(td, canvas)

	Convey("Checking that generator runs for a single function with a map of 4 colours", t, func() {
		Convey("Run using ./collater.json for the image widgets", func() {
			Convey("A image with 4 coloured squares is generated and sha256 is identical to ./testdata/gen_simple_test.png", func() {
				So(htest.Sum(nil), ShouldResemble, hnormal.Sum(nil))
			})
		})
	})
	os.Remove(fileName)
}

func TestZposRun(t *testing.T) {

	target := []string{"./testdata/mockcont/base_hook.png", "./testdata/mockcont/base.png"} // , "./testdata/gen_simple_test.png"}
	files := []string{"./testdata/mockcont/base_hook.json", "./testdata/mockcont/base_hook_first.json"}
	additions := []string{"./collater_hook.json", "./collater_hook_reverse.json"}
	extras := []string{"", ""}
	results := []string{"A image with 4 coloured squares is generated, with a black square in the top left from the z order of collater_hook",
		"A image with 4 coloured squares is generated, the black square is overdrawn as it is declared first of the 5 widgets"}

	for i, ftarget := range target {

		mLog := errhandle.LogInit("stdout", "")
		// generate a mocked context for use with gen
		mockC, canvas := contMocker(files[i], additions[i], extras[i])

		mockConfig := GenConf[test]{true, mockSchema, "mocktest", nil}
		mockCanvas := GenConf[canvaswidget.ConfigVals]{true, mockSchema, "builtin.canvasoptions", nil}
		mockhook := GenConf[testhook]{true, mockSchema, "mockhook", nil}

		// add canvas
		canvasChan := make(chan draw.Image, 1)
		canvasChan <- canvas

		var wg sync.WaitGroup
		wg.Add(3)
		// Replicate the context that would normally be used
		go WidgetRunner(canvasChan, mockCanvas, mockC, mLog, &wg)
		go WidgetRunner(canvasChan, mockConfig, mockC, mLog, &wg)
		go WidgetRunner(canvasChan, mockhook, mockC, mLog, &wg)
		time.Sleep(time.Second) // mock the external wait function

		// mock an extra waitgroup
		time.Sleep(1 * time.Second)
		file, _ := os.Open(ftarget)
		// decode to get the colour values
		baseVals, _ := png.Decode(file)

		// assign the colour to the correct type of image NGRBA64 and replace the colour values
		readImage := image.NewNRGBA64(baseVals.Bounds())
		colour.Draw(readImage, readImage.Bounds(), baseVals, image.Point{}, draw.Src)

		hnormal := sha256.New()
		htest := sha256.New()
		hnormal.Write(readImage.Pix)
		htest.Write(canvas.(*colour.NRGBA64).Pix())

		// td, _ := os.Create(ftarget + "r.png")
		// png.Encode(td, canvas)

		Convey("Checking that generator runs for the zOrder across two functions", t, func() {
			Convey(fmt.Sprintf("using the file of %s to generate the widgets", additions[i]), func() {
				Convey(results[i], func() {
					So(htest.Sum(nil), ShouldResemble, hnormal.Sum(nil))
				})
			})
		})
		os.Remove(files[i])
	}
}

func TestErrorZpos(t *testing.T) {

	target := []string{"./testdata/mockcont/base_missed.png", "./testdata/mockcont/base.png", "./testdata/mockcont/base_hook.png"}
	files := []string{"./testdata/mockcont/base_location.json", "./testdata/mockcont/base_error.json", "./testdata/mockcont/base_missed.json"}
	additions := []string{"./collater_hook.json", "./collater_hook.json", "./collater_hook_missed.json"}
	extras := []string{`, "load.red": {"position":"not real"} , "load.green": {"position":"not real"}`, `, "load.black": {"err":"the wheels have come off"}`, ""}
	results := []string{"The zpos is still updated and the black square is written when half the widgets have failed",
		"That the failed widget zpos due to poor bounds is updated allowing the next zpos to run and the black square is still written",
		"That the missed widget zpos is updated allowing the next zpos to run and the black square is still written"}

	for i, ftarget := range target {

		mLog := errhandle.LogInit("stdout", "")
		// generate a mocked context for use with gen
		mockC, canvas := contMocker(files[i], additions[i], extras[i])

		mockConfig := GenConf[test]{true, mockSchema, "mocktest", nil}
		mockCanvas := GenConf[canvaswidget.ConfigVals]{true, mockSchema, "builtin.canvasoptions", nil}
		mockhook := GenConf[testhook]{true, mockSchema, "mockhook", nil}

		// add canvas
		canvasChan := make(chan draw.Image, 1)
		canvasChan <- canvas

		var wg sync.WaitGroup
		wg.Add(4)
		var wg2 sync.WaitGroup
		wg2.Add(1)
		// Replicate the context that would normally be used
		go WidgetRunner(canvasChan, mockCanvas, mockC, mLog, &wg)
		go WidgetRunner(canvasChan, mockConfig, mockC, mLog, &wg)
		go WidgetRunner(canvasChan, mockhook, mockC, mLog, &wg)
		go MockMissedGen(canvasChan, true, mockC, &wg2, &wg, mLog)
		time.Sleep(time.Second) // mock the external wait function

		// mock an extra waitgroup
		time.Sleep(1 * time.Second)
		file, _ := os.Open(ftarget)
		// decode to get the colour values
		baseVals, _ := png.Decode(file)

		// assign the colour to the correct type of image NGRBA64 and replace the colour values
		readImage := image.NewNRGBA64(baseVals.Bounds())
		colour.Draw(readImage, readImage.Bounds(), baseVals, image.Point{}, draw.Src)

		hnormal := sha256.New()
		htest := sha256.New()
		hnormal.Write(readImage.Pix)
		htest.Write(canvas.(*colour.NRGBA64).Pix())

		// td, _ := os.Create(fmt.Sprintf("%vr.png", ftarget))
		// png.Encode(td, canvas)

		Convey("Checking that generator runs in the zOrder when errors are emitted", t, func() {
			Convey(fmt.Sprintf("using the file of %s with the additional call of %s to generate the widgets", additions[i], extras[i]), func() {
				Convey(results[i], func() {
					So(htest.Sum(nil), ShouldResemble, hnormal.Sum(nil))
				})
			})
		})
		os.Remove(files[i])
	}
}

// Structs for testing the generator function
type test struct {
	Colour   string `json:"colour,omitempty"` // ensure the json tag is used
	Position string `json:"position,omitempty"`
}

// Mock generator functions
func (tt test) Generate(i draw.Image, t ...any) error {
	c := colourgen.HexToColour(tt.Colour, colour.ColorSpace{})
	// fmt.Println(tt.Fill)
	colour.Draw(i.(*colour.NRGBA64), i.Bounds(), &image.Uniform{c}, image.Point{}, draw.Src)

	return nil
}

func (tt test) Alias() string {

	return ""
}

func (tt test) Location() string {

	return tt.Position
}

type testhook struct {
	Loc string `json:"location"`
	Err string `json:"err"`
}

func (tt testhook) Generate(i draw.Image, t ...any) error {
	if tt.Err != "" {

		return fmt.Errorf("%v", tt.Err)
	}
	colour.Draw(i, i.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	return nil
}

func (tt testhook) Alias() string {

	return ""
}

func (tt testhook) Location() string {
	if tt.Loc != "" {

		return tt.Loc
	}

	return "a0"
}

func TestPut(t *testing.T) {

	type testString struct {
		Content string `json:"content"`
	}
	type testInt struct {
		Content int `json:"content"`
	}

	tests := []any{testString{"test1"}, testInt{6}}
	names := []string{"testStringAlias", "testIntAlias"}
	added := []any{"test1", 6} // gets shifted to float64 from json back to reality
	// gen and expected should grow at the same rate
	expected := make(map[string]map[any]interface{})
	dummy := context.TODO()
	newC := MetaDataInit(dummy)

	for i, ts := range tests {
		// map to add
		gen := make(map[core.AliasIdentity]interface{})
		gen[core.AliasIdentity{Alias: names[i], ZPos: 0}] = ts
		err := put(gen, newC)

		// add to the expected total
		expected[names[i]] = make(map[any]interface{})
		expected[names[i]]["content"] = added[i]

		imageGeneration := (*newC).Value(metakey).(metadata)

		//imageGeneration := Extract(newC, "testkey", "")
		Convey("Checking that put assigns values to image generation", t, func() {
			Convey(fmt.Sprintf("using %v as the input map", gen), func() {
				Convey(fmt.Sprintf("the expected map of %v is returned as the saved data", expected), func() {
					So(err, ShouldBeNil)
					So(imageGeneration.data, ShouldResemble, expected)
				})
			})
		})
	}
}

func TestExtract(t *testing.T) {

	type threeLayer struct {
		CheckString string `json:"checkString,omitempty" yaml:"checkString,omitempty"`
	}

	type twoLayer struct {
		Nest     threeLayer `json:"threelayer,omitempty" yaml:"threelayer,omitempty"`
		CheckInt int        `json:"checkInt" yaml:"checkInt"`
	}

	type nested struct {
		Nest twoLayer `json:"twolayer,omitempty" yaml:"twolayer,omitempty"`
	}

	var top nested
	var middle twoLayer
	var bottom threeLayer
	bottom.CheckString = "bottom"
	middle.CheckInt = 0
	middle.Nest = bottom
	top.Nest = middle

	gen := make(map[core.AliasIdentity]nested)
	gen[core.AliasIdentity{Alias: "testKey", ZPos: 0}] = top
	dummy := context.TODO()
	newC := MetaDataInit(dummy)
	_ = put(gen, newC)

	keys := [][]string{{"twolayer", "threelayer", "checkString"},
		{"twolayer", "checkInt"},
		{},
	}
	expec := make(map[any]interface{})
	mid := make(map[any]interface{})
	bot := make(map[any]interface{})
	bot["checkString"] = "bottom"
	mid["threelayer"] = bot
	mid["checkInt"] = 0

	expec["twolayer"] = mid
	expected := []interface{}{
		"bottom", 0, expec,
	}

	for i, key := range keys {
		result := Extract(newC, "testKey", key...)

		Convey("Checking that extract can iterativley search through maps", t, func() {
			Convey(fmt.Sprintf("Searching with these keys %v", key), func() {
				Convey(fmt.Sprintf("A value of %v is returned", expected[i]), func() {
					So(result, ShouldResemble, expected[i])
				})
			})
		})
	}
}

// go test ./widgethandler/ -bench=. -benchtime=1s
/*
func BenchmarkSpill(b *testing.B) {

	mLog := log.Default()
	//generate a mocked context for use with gen
	mockC, mockBase := contMocker()

	mockConfig := GenConf[test]{true, mockSchema, "mocktest", nil}

	//add canvas
	canvasChan := make(chan draw.Image, 1)
	canvas := image.NewNRGBA64(image.Rect(0, 0, 1000, 1000))
	canvasChan <- canvas

	//Replicate the context that would normally be used
	mockC = context.WithValue(mockC, "widget bases", mockBase)
	//mockCMHook := GetContext(mockC)
	input := mockBase.Data["mocktest"]
	in := make(map[config.AliasIdentity]test)
	for k, v := range input {
		in[k] = v.(test)
	}
	mockCMHook := GetContext(mockC)

	for n := 0; n < b.N; n++ {
		spill(canvasChan, mockConfig, mockCMHook, mLog, "test", in)

	}
}

func BenchmarkParam(b *testing.B) {

	mLog := log.Default()
	//generate a mocked context for use with gen
	mockC, mockBase := contMocker()

	mockConfig := GenConf[test]{true, mockSchema, "mocktest", nil}

	//add canvas
	canvasChan := make(chan draw.Image, 1)
	canvas := image.NewNRGBA64(image.Rect(0, 0, 1000, 1000))
	canvasChan <- canvas

	//Replicate the context that would normally be used
	mockC = context.WithValue(mockC, "widget bases", mockBase)
	//mockCMHook := GetContext(mockC)
	input := mockBase.Data["mocktest"]
	in := make(map[config.AliasIdentity]test)
	for k, v := range input {
		in[k] = v.(test)
	}
	for n := 0; n < b.N; n++ {
		mockCMHook := GetContext(mockC)
		pram(canvasChan, mockConfig, mockCMHook, mLog, "test", in)

	}
}
*/