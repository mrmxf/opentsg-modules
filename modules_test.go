package modules

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/xeipuuv/gojsonschema"
	"gitlab.com/golang-commonmark/markdown"
	"gopkg.in/yaml.v3"
)

func init() {
	schemaFolder = regexp.MustCompile(`^[jJ]son[sS]chema$`)
	jsonFile = regexp.MustCompile(`[jJ][sS][oO][nN]$`)
	yamlFile = regexp.MustCompile(`[yY][aA][mM][lL]$`)
	mdFile = regexp.MustCompile(`[mM][dD]$`)
}

var schemaFolder *regexp.Regexp
var jsonFile *regexp.Regexp
var yamlFile *regexp.Regexp
var mdFile *regexp.Regexp

func TestXxx(t *testing.T) {

	our, _ := os.ReadDir(".")
	fmt.Println(our)
	sg := SchmeaGetter{schemas: make([]schemaBody, 0),
		skips: make([]string, 0),
		fails: make([]string, 0), pass: make([]string, 0)}
	fmt.Println(sg.GetRecursiveSchemas("."))
	fmt.Println(len(sg.schemas))

	directory, _ := filepath.Abs(".")
	fmt.Println(directory)
	sg.ValidateJsons(directory)

	for _, p := range sg.pass {
		fmt.Println(p)
	}

	for _, p := range sg.skips {
		fmt.Println(p)
	}
	for _, f := range sg.fails {
		fmt.Println(f)
	}
	fmt.Println("test", sg.globalPass)
}

func (s *SchmeaGetter) ValidateJsons(directory string) error {

	dirs, err := os.ReadDir(directory)

	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			jPath := filepath.Join(directory, "/", dir.Name())
			if jsonFile.MatchString(dir.Name()) || yamlFile.MatchString(dir.Name()) {

				// is it a file we want to read?
				if fileFence(jPath) {
					s.skips = append(s.skips, fmt.Sprintf("Skip %s", jPath))
					continue
				}

				input, _ := os.ReadFile(jPath)

				// cleanse input to json if initially yaml
				if !json.Valid(input) { // if not json open it as yaml and save as json
					var clean any
					err := yaml.Unmarshal(input, &clean)
					if err != nil {
						s.fails = append(s.fails, fmt.Sprintf("Fail %s %s", jPath, err.Error()))
						s.globalPass = false
						continue
						// label as a fail and continue
					}

					input, err = json.Marshal(clean)
					if err != nil {
						s.fails = append(s.fails, fmt.Sprintf("Fail %s %s", jPath, err.Error()))
						s.globalPass = false
						continue
						// label as bad data and flag as a dail
					}
				}

				var schemPass bool
				// see which schemas it passes
				for _, schema := range s.schemas {
					// Loop here
					schemaLoader := gojsonschema.NewBytesLoader(schema.data)
					documentLoader := gojsonschema.NewBytesLoader(input)

					result, err := gojsonschema.Validate(schemaLoader, documentLoader)

					if err != nil {

						continue
					} else if !result.Valid() {
						continue
					}

					schemPass = true
					// mark as a pass and break
					break
				}

				if schemPass {
					s.pass = append(s.pass, fmt.Sprintf("PASS %s", jPath))
				} else {
					s.globalPass = false
					s.fails = append(s.fails, fmt.Sprintf("Fail %s \n", jPath))
				}

			} else if mdFile.MatchString(dir.Name()) {
				input, _ := os.ReadFile(jPath)
				md := markdown.New(markdown.XHTMLOutput(true), markdown.Nofollow(true))
				tokens := md.Parse(input)

				for _, t := range tokens {

					snip, ok := t.(*markdown.Fence)
					if ok {

						var snippet []byte
						if strings.Contains(snip.Params, ".valid") {
							fmt.Println(snip.Params)
						}

						switch snip.Params {
						case "json":
							snippet = []byte(snip.Content)
						case "yaml":
							var clean any
							err := yaml.Unmarshal(input, &clean)
							if err != nil {
								s.fails = append(s.fails, fmt.Sprintf("Fail %s %s", jPath, err.Error()))
								s.globalPass = false
								continue
								// label as a fail and continue
							}

							snippet, err = json.Marshal(clean)
							if err != nil {
								s.fails = append(s.fails, fmt.Sprintf("Fail %s %s", jPath, err.Error()))
								s.globalPass = false
								continue
								// label as bad data and flag as a dail
							}
						default:
							// not a token we want to deal with
							continue

						}

						var schemPass bool
						// see which schemas it passes
						for _, schema := range s.schemas {
							// Loop here
							schemaLoader := gojsonschema.NewBytesLoader(schema.data)
							documentLoader := gojsonschema.NewBytesLoader(snippet)

							result, err := gojsonschema.Validate(schemaLoader, documentLoader)

							if err != nil {

								continue
							} else if !result.Valid() {
								continue
							}

							schemPass = true
							// mark as a pass and break
							break
						}

						if schemPass {
							s.pass = append(s.pass, fmt.Sprintf("PASS %s", jPath))
						} else {
							s.globalPass = false
							s.fails = append(s.fails, fmt.Sprintf("Fail %s \n", jPath))
						}

					}
				}
			}

		} else {
			dpath := filepath.Join(directory, "/", dir.Name())
			s.ValidateJsons(dpath)

		}

	}

	return nil
}

// fileFence fences files that may be in the errors test folders
func fileFence(filepath string) (skip bool) {
	if strings.Contains(filepath, "errors") {
		return true
	}

	return false
}

type SchmeaGetter struct {
	schemas    []schemaBody
	globalPass bool
	fails      []string
	skips      []string
	pass       []string
}

type schemaBody struct {
	data []byte
	file string
}

func (s *SchmeaGetter) GetRecursiveSchemas(directory string) error {

	/*
		loop through
	*/

	directory, err := filepath.Abs(directory)
	if err != nil {
		return err
	}

	dirs, err := os.ReadDir(directory)

	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			err := s.getRecursiveSchemas(dir, directory)

			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (s *SchmeaGetter) getRecursiveSchemas(dir fs.DirEntry, parent string) error {

	/*
		loop through
	*/
	path, _ := filepath.Abs(filepath.Join(parent, "/", dir.Name()))

	dirs, err := os.ReadDir(path)
	if err != nil {

		return err
	}

	if schemaFolder.MatchString(dir.Name()) {
		for _, dir := range dirs {
			if !dir.IsDir() {
				schema := filepath.Join(path + "/" + dir.Name())
				schemaBytes, err := os.ReadFile(schema)

				if err != nil {
					return err
				}

				s.schemas = append(s.schemas, schemaBody{file: schema, data: schemaBytes})
			}
		}

		// loop through files and get the jsons
	} else {

		for _, dir := range dirs {
			if dir.IsDir() {
				err := s.getRecursiveSchemas(dir, path)

				if err != nil {
					return err
				}
			}

		}
	}

	return nil
}
