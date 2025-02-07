package modules

import (
	"context"
	"embed"
	"log/slog"

	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/xeipuuv/gojsonschema"
	"gitlab.com/golang-commonmark/markdown"
	"gopkg.in/yaml.v3"

	"github.com/phsym/console-slog"
	. "github.com/smartystreets/goconvey/convey"

	slogmulti "github.com/samber/slog-multi"
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

//go:embed opentsg-core/canvaswidget/jsonschema/baseschema.json opentsg-core/config/core/jsonschema/includeschema.json opentsg-core/config/core/jsonschema/widgetEssentials.json opentsg-core/gridgen/jsonschema/tsigschema.json opentsg-widgets/addimage/jsonschema/addimageschema.json opentsg-widgets/bowtie/jsonschema/jsonschema.json opentsg-widgets/ebu3373/bars/jsonschema/barschema.json opentsg-widgets/ebu3373/luma/jsonschema/lumaschema.json opentsg-widgets/ebu3373/nearblack/jsonschema/nbschema.json opentsg-widgets/ebu3373/saturation/jsonschema/satschema.json opentsg-widgets/ebu3373/twosi/jsonschema/twoschema.json opentsg-widgets/fourcolour/jsonschema/fourColourSchema.json opentsg-widgets/framecount/jsonschema/framecounter.json opentsg-widgets/geometryText/jsonschema/geometryText.json opentsg-widgets/gradients/jsonschema/gradientSchema.json opentsg-widgets/jsonschema/framecounter.json opentsg-widgets/jsonschema/gridschema.json opentsg-widgets/noise/jsonschema/noiseschema.json opentsg-widgets/qrgen/jsonschema/qrgenschema.json opentsg-widgets/resize/jsonschema/resize.json opentsg-widgets/textbox/jsonschema/textBoxSchema.json opentsg-widgets/zoneplate/jsonschema/zoneplateschema.json
var schemas embed.FS

var schemaNames = []string{"opentsg-core/canvaswidget/jsonschema/baseschema.json", "opentsg-core/config/core/jsonschema/includeschema.json", "opentsg-core/config/core/jsonschema/widgetEssentials.json", "opentsg-core/gridgen/jsonschema/tsigschema.json", "opentsg-widgets/addimage/jsonschema/addimageschema.json", "opentsg-widgets/bowtie/jsonschema/jsonschema.json", "opentsg-widgets/ebu3373/bars/jsonschema/barschema.json", "opentsg-widgets/ebu3373/luma/jsonschema/lumaschema.json", "opentsg-widgets/ebu3373/nearblack/jsonschema/nbschema.json", "opentsg-widgets/ebu3373/saturation/jsonschema/satschema.json", "opentsg-widgets/ebu3373/twosi/jsonschema/twoschema.json", "opentsg-widgets/fourcolour/jsonschema/fourColourSchema.json", "opentsg-widgets/framecount/jsonschema/framecounter.json", "opentsg-widgets/geometryText/jsonschema/geometryText.json", "opentsg-widgets/gradients/jsonschema/gradientSchema.json", "opentsg-widgets/jsonschema/framecounter.json", "opentsg-widgets/jsonschema/gridschema.json", "opentsg-widgets/noise/jsonschema/noiseschema.json", "opentsg-widgets/qrgen/jsonschema/qrgenschema.json", "opentsg-widgets/resize/jsonschema/resize.json", "opentsg-widgets/textbox/jsonschema/textBoxSchema.json", "opentsg-widgets/zoneplate/jsonschema/zoneplateschema.json"}

type schemaBody struct {
	data     []byte
	filePath string
}

type SchemaValidator struct {
	schemas    []schemaBody
	globalPass bool

	run []run
}

// TestSlogger wraps the testing.T to write to t.Log
type TestSlogger struct {
	t *testing.T
}

// Write forwards the message to t.Log
func (t *TestSlogger) Write(p []byte) (n int, err error) {
	t.t.Log(string(p))

	return len(p), nil
}

// run is all the fields needed
// to run schemas on a directory
type run struct {
	schemaOrigin string
	schemas      []schemaBody
	target       string
	ignores      []string
}

// SchemaConfig contains an array
// of directories to get schemas from
// and the target directory to compare against the schemas
type SchemaConfig struct {
	SchemaChecks []SchemaCheck `yaml:"schemaCheck"`
}

// SchemaCheck contains the configuration options
// for checking the jsons of a directory and their schemas.
type SchemaCheck struct {
	// if nil use TSG,
	// if object or object array presume its a schema to parse (not implemented)
	// if string presume its a file or folder - if array repeat the process
	Schema any `yaml:"schema"`

	// a folder or file to check
	DirectoryToCheck string `yaml:"directory"`

	// strings that are set up to be ignored
	// string must match any part of the path to be skipped
	Ignore []string `yaml:"ignore"`
}

// NewSchemaValidatorFile validates files based off an input file
// @TODO update to reader
func NewSchemaValidatorFile(confFile string) (*SchemaValidator, error) {

	confBytes, err := os.ReadFile(confFile)
	if err != nil {
		return nil, err
	}

	var sc SchemaConfig
	err = yaml.Unmarshal(confBytes, &sc)
	if err != nil {
		return nil, err
	}

	return NewSchemaValidator(&sc)

}

// NewSchemaValidator returns a schemaValidator
// tailored to the config options.
// it can then be plugged into your testing
func NewSchemaValidator(conf *SchemaConfig) (*SchemaValidator, error) {

	// then forward it to the schemavaldiator
	schems := make([]schemaBody, len(schemaNames))
	for i, schema := range schemaNames {
		sch, _ := schemas.ReadFile(schema)
		schems[i] = schemaBody{data: sch, filePath: schema}
	}

	// open the file and unmarshal here nd see how it goes

	runners := make([]run, len(conf.SchemaChecks))
	for i, sc := range conf.SchemaChecks {
		var gotSchemas []schemaBody
		schemaOrigin := "default OTSG schemas"
		switch schMethod := sc.Schema.(type) {

		/*
			case []string:
				loop through the files
			case map[string]any:
				presume its a schema
		*/
		case []any:

			names := make([]string, len(schMethod))
			// loop through adding to the array of schemas
			for i, origin := range schMethod {

				midOrigin, err := filepath.Abs(fmt.Sprintf("%v", origin))
				if err != nil {
					return nil, err
				}

				// update the schemas with each origin
				gotSchemas, err = getSchemas(gotSchemas, midOrigin)
				if err != nil {
					return nil, err
				}

				names[i] = fmt.Sprintf("%v", origin)
			}

			schemaOrigin = strings.Join(names, ",")

		case string:
			var err error
			schemaOrigin, err = filepath.Abs(schMethod)
			if err != nil {
				return nil, err
			}
			gotSchemas, err = getSchemas(gotSchemas, schemaOrigin)

			if err != nil {
				return nil, err
			}

		default:
			gotSchemas = schems
		}

		runners[i] = run{schemaOrigin: schemaOrigin, target: sc.DirectoryToCheck,
			schemas: gotSchemas, ignores: sc.Ignore}
	}

	return &SchemaValidator{schemas: schems,
		run: runners,
	}, nil
}

// JobLogger is a no-color version of the PrettyLogger that is created
// to append to write to a given writer
func NewTestLogger(w io.Writer, t *testing.T, level slog.Level) *slog.Logger {

	Loggers := slog.New(
		slogmulti.Fanout(
			console.NewHandler(w,
				&console.HandlerOptions{Level: level, NoColor: true}),
			// create our logging error
			console.NewHandler(&TestSlogger{t: t},
				&console.HandlerOptions{Level: slog.LevelWarn, NoColor: true}),
		))

	return Loggers
}

// ValidateJsons valdiates every json, yaml and markdown snippet labelled json/yaml
// in a directory. It recursively searches every folder in the directory
func (s *SchemaValidator) ValidateJsons(t *testing.T) {

	//set the logging
	Loggers := NewTestLogger(io.Discard, t, slog.LevelInfo)

	slog.SetDefault(Loggers)

	// absolute the path before processing
	for _, runner := range s.run {
		s.globalPass = true
		if len(s.schemas) == 0 {
			panic("no schemas declared, can not validate jsons")
		}

		directory, dErr := filepath.Abs(runner.target)

		// loop through every run

		base, jErr := s.getJsons(directory, make([]string, 0), runner.ignores)

		s.globalPass = true
		valErr := s.validateJsons(base, runner.schemas)

		// run the go convey tests here

		Convey(fmt.Sprintf("Checking every json and yaml file/snippet in %s are valid", directory), t, func() {
			Convey(fmt.Sprintf("Checking files against the %s schemas", runner.schemaOrigin), func() {
				Convey("Every snippet passes a schema", func() {
					So(dErr, ShouldBeNil)
					So(jErr, ShouldBeNil)
					So(valErr, ShouldBeNil)
					So(s.globalPass, ShouldBeTrue)
				})
			})
		})
	}

}

// Get all the files in the directory with a recursive search
func (s *SchemaValidator) getJsons(directory string, files []string, ignore []string) ([]string, error) {

	// get the files in the directory
	dirs, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, dir := range dirs {
		// check if file
		if !dir.IsDir() {
			jPath := filepath.Join(directory, "/", dir.Name())
			// check file type
			if jsonFile.MatchString(dir.Name()) || yamlFile.MatchString(dir.Name()) {

				// is it a file we want to read?
				if fileFence(jPath, ignore) {
					slog.Log(context.TODO(), slog.LevelInfo, jPath, "Staus", "Skip")

					continue
				}

				files = append(files, jPath)

			} else if mdFile.MatchString(dir.Name()) {
				files = append(files, jPath)
			}

		} else {
			// continue the path search
			dpath := filepath.Join(directory, "/", dir.Name())
			files, err = s.getJsons(dpath, files, ignore)

			if err != nil {
				return nil, err
			}
		}
	}
	return files, nil
}

// validateJsons validates the json contents against the schema logging if files parr fail or succeed
func (s *SchemaValidator) validateJsons(targets []string, schemas []schemaBody) error {

	for _, target := range targets {
		// check if file

		// check file type
		if jsonFile.MatchString(target) || yamlFile.MatchString(target) {

			input, err := os.ReadFile(target)
			if err != nil {
				return err
			}

			// cleanse input to json if initially yaml
			if !json.Valid(input) { // if not json open it as yaml and save as json
				var clean any
				err := yaml.Unmarshal(input, &clean)
				if err != nil {
					slog.Log(context.TODO(), slog.LevelWarn, fmt.Sprintf("%s %s", target, err.Error()), "Staus", "Fail")

					s.globalPass = false
					continue
					// label as a fail and continue
				}

				input, err = json.Marshal(clean)
				if err != nil {
					slog.Log(context.TODO(), slog.LevelWarn, fmt.Sprintf("%s %s", target, err.Error()), "Staus", "Fail")
					s.globalPass = false
					continue
					// label as bad data and flag as a dail
				}
			}

			var schemPass bool
			// see which schemas it passes
			for _, schema := range schemas {
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
				slog.Log(context.TODO(), slog.LevelInfo, target, "Staus", "Pass")
			} else {
				s.globalPass = false
				slog.Log(context.TODO(), slog.LevelWarn, target, "Staus", "Fail")
			}

		} else if mdFile.MatchString(target) {
			input, err := os.ReadFile(target)
			if err != nil {
				return err
			}
			md := markdown.New(markdown.XHTMLOutput(true), markdown.Nofollow(true))
			tokens := md.Parse(input)

			for _, t := range tokens {

				snip, ok := t.(*markdown.Fence)
				if ok {

					var snippet []byte

					switch snip.Params {
					case "json":
						snippet = []byte(snip.Content)
					case "yaml": // no yaml snippets yet
						var clean any
						err := yaml.Unmarshal(input, &clean)
						if err != nil {
							slog.Log(context.TODO(), slog.LevelWarn, fmt.Sprintf("%s %s", target, err.Error()), "Staus", "Fail")
							s.globalPass = false
							continue
							// label as a fail and continue
						}

						snippet, err = json.Marshal(clean)
						if err != nil {
							slog.Log(context.TODO(), slog.LevelWarn, fmt.Sprintf("%s %s", target, err.Error()), "Staus", "Fail")
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
					for _, schema := range schemas {
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
						slog.Log(context.TODO(), slog.LevelInfo, target, "Staus", "Pass")
					} else {
						s.globalPass = false
						slog.Log(context.TODO(), slog.LevelWarn, target, "Staus", "Fail")
					}

				}
			}
		}

	}
	return nil
}

// fileFence fences files that may be in the errors test folders
func fileFence(filepath string, ignores []string) (skip bool) {

	for _, ignore := range ignores {
		// add more fields here as the tests grow
		if strings.Contains(filepath, ignore) {
			return true
		}
	}

	return false
}

// getSchemas recursively searches a directory
// for any schemas and adds them to the schema list.
// It presumes any files in a folder labelled jsonschema are valid and they
// are parsed with no questions asked.
func getSchemas(schemas []schemaBody, path string) ([]schemaBody, error) {

	dirs, err := os.ReadDir(path)
	if err != nil {

		return nil, err
	}

	folName := filepath.Base(path)

	if schemaFolder.MatchString(folName) {
		for _, dir := range dirs {
			if !dir.IsDir() {
				schema := filepath.Join(path + "/" + dir.Name())
				schemaBytes, err := os.ReadFile(schema)

				if err != nil {
					return nil, err
				}

				schemas = append(schemas, schemaBody{filePath: schema, data: schemaBytes})
			}
		}

		// loop through files and get the jsons
	} else {

		for _, dir := range dirs {
			if dir.IsDir() {
				dirPath := filepath.Join(path, "/", dir.Name())
				schemas, err = getSchemas(schemas, dirPath)

				if err != nil {
					return nil, err
				}
			}

		}
	}

	return schemas, nil
}
