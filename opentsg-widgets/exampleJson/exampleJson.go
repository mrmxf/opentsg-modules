package examplejson

import (
	"context"
	"encoding/json"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/mrmxf/opentsg-modules/opentsg-core/widgethandler"
)

/*



 */

func SaveExampleJson(example widgethandler.Generator, folder, name string, saveImage bool) {

	jsonExample, _ := json.MarshalIndent(example, "", "    ")

	// get the demo folder so it can be found in new repos
	pwd, _ := os.Getwd()

	// how many folders up do we need to go to save in exampleJson
	ups := strings.Count(folder, "/")
	basePath := "../"
	for i := 0; i < ups; i++ {
		basePath += "../"
	}

	base := filepath.Join(pwd, basePath+"exampleJson/")

	// check a folder exists
	if _, err := os.Stat(base + string(os.PathSeparator) + folder); os.IsNotExist(err) {
		os.MkdirAll(base+string(os.PathSeparator)+folder, 0777)
	}

	if saveImage {
		baseImage := image.NewNRGBA64(image.Rect(0, 0, 500, 500))
		mockCont := context.Background()
		example.Generate(baseImage, &mockCont)
		fImg, _ := os.Create(base + string(os.PathSeparator) + folder + string(os.PathSeparator) + name + "-example.png")
		png.Encode(fImg, baseImage)

		// Add the type and location fields
		var updater map[string]any
		json.Unmarshal(jsonExample, &updater)
		updater["props"] = map[string]any{"type": folder,
			"location": map[string]any{
				"alias": "A demo Alias",
				"box":   map[string]any{"x": 1, "y": 1}}}
		jsonExample, _ = json.MarshalIndent(updater, "", "    ")

	}

	f, _ := os.Create(base + string(os.PathSeparator) + folder + string(os.PathSeparator) + name + "-example.json")
	f.Write(jsonExample)

}
