package modules

import (
	"testing"
)

func TestJsonSnippets(t *testing.T) {

	// design some tests to always pass etc
	/*
		// here we generate the schemas for the modules to use as imported files
		files := ""
		filesString := ""
		for _, s := range sg.schemas {
			sc := strings.ReplaceAll(s.filePath, directory+"/", "")
			files += sc + " "
			filesString += "\"" + sc + "\"" + ","
		}
		embedder := fmt.Sprintf("//go:embed %s\nvar schemas embed.FS", files)
		array := fmt.Sprintf("var schemaNames = []string{%s}", filesString)
		fmt.Println(embedder)
		fmt.Println(array)
	*/

	sv, _ := NewSchemaValidator(&SchemaConfig{[]SchemaCheck{
		{DirectoryToCheck: "./opentsg-core"},
		{DirectoryToCheck: "./opentsg-io", Schema: "."},
		{DirectoryToCheck: "./opentsg-widgets"},
	}})

	sv.ValidateJsons(t)

	svf, _ := NewSchemaValidatorFile("testdata/schema.yaml")
	svf.ValidateJsons(t)

}
