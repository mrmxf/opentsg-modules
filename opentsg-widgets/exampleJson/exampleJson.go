package examplejson

import (
	"encoding/json"
	"os"
)

/*



 */

// be able to change the base some how
const (
	base = "/workspace/opentsg-widgets/exampleJson/"
)

func SaveExampleJson(example any, folder, name string) {

	jsonExample, _ := json.MarshalIndent(example, "", "    ")

	// check a folder exists
	if _, err := os.Stat(base + string(os.PathSeparator) + folder); os.IsNotExist(err) {
		os.MkdirAll(base+string(os.PathSeparator)+folder, 0777)
	}

	f, _ := os.Create(base + string(os.PathSeparator) + folder + string(os.PathSeparator) + name + "-example.json")

	f.Write(jsonExample)

}
