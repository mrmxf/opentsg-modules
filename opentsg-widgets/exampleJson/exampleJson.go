package examplejson

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/mrmxf/opentsg-modules/opentsg-core/widgethandler"
)

/*



 */

// be able to change the base some how
const (
	base = "/workspace/opentsg-modules/opentsg-widgets/exampleJson/"
)

func SaveExampleJson(example widgethandler.Generator, folder, name string, saveImage bool) {

	jsonExample, _ := json.MarshalIndent(example, "", "    ")

	// check a folder exists
	if _, err := os.Stat(base + string(os.PathSeparator) + folder); os.IsNotExist(err) {
		os.MkdirAll(base+string(os.PathSeparator)+folder, 0777)
	}

	if saveImage {
		baseImage := image.NewNRGBA64(image.Rect(0, 0, 500, 500))
		mockCont := context.Background()
		example.Generate(baseImage, &mockCont)
		fImg, e := os.Create(base + string(os.PathSeparator) + folder + string(os.PathSeparator) + name + "-example.png")
		fmt.Println(png.Encode(fImg, baseImage))
		fmt.Println(e, "HERE", base+string(os.PathSeparator)+folder+string(os.PathSeparator)+name+"-example.png")
	}

	f, _ := os.Create(base + string(os.PathSeparator) + folder + string(os.PathSeparator) + name + "-example.json")
	f.Write(jsonExample)

}
