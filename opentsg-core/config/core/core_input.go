package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/mrmxf/opentsg-modules/opentsg-core/config/validator"
	"github.com/mrmxf/opentsg-modules/opentsg-core/credentials"
	"gopkg.in/yaml.v3"
)

// File import reads a factory json file and extracts all the included files and factories recursively.
// It returns a context holding them, the number of frames to be run and any errors encountered.
func FileImport(inputFile, profile string, debug bool, httpKeys ...string) (context.Context, int, error) {
	cont := context.Background()
	authDecoder, err := credentials.AuthInit(profile, httpKeys...)
	if err != nil {
		return cont, 0, fmt.Errorf("0000 %v", err)
	}
	inputFile, _ = filepath.Abs(inputFile)
	inputBytes, err := os.ReadFile(inputFile)
	if err != nil {
		return cont, 0, fmt.Errorf("0001 %v", err)
	}

	data := make(validator.JSONLines)
	err = validator.Liner(inputBytes, inputFile, "factory", data)
	if err != nil { // return just the first error? figure out error handling
		return cont, 0, err
	}

	errs := validator.SchemaValidator(incschema, inputBytes, inputFile, data)
	if err != nil { // return just the first error? TODO figure out error handling
		return cont, 0, errs[0]
	}

	// Take the input
	var inputFactory factory
	err = yaml.Unmarshal(inputBytes, &inputFactory)
	if err != nil {
		return cont, 0, fmt.Errorf("0002 %v when opening %s", err, inputFile)
	}
	if len(inputFactory.Create) == 0 {
		return cont, 0, fmt.Errorf("0003 No frames declared in %s", inputFile)
	}

	holder := base{importedFactories: make(map[string]factory), importedWidgets: make(map[string]json.RawMessage),
		jsonFileLines: data, authBody: authDecoder, metadataParams: map[string][]string{}}

	err = holder.factoryInit(inputFactory, filepath.Dir(inputFile), "", []int{})
	if err != nil {
		return cont, 0, err
	}

	// init a global alias map for use in grid gen
	aliasMap := SyncMap{make(map[string]string), &sync.Mutex{}}
	cont = context.WithValue(cont, aliasKey, aliasMap)
	// add the base factory and the imported schemas
	cont = context.WithValue(cont, updates, inputFactory)
	cont = context.WithValue(cont, frameHolders, holder)
	// add the working directory everything is relative to
	cont = context.WithValue(cont, factoryDir, filepath.Dir(inputFile))

	return cont, len(inputFactory.Create), nil
}

// Make a map of all the factories and their information
// so that are all open at once and there is no reopening of files and changing of data
func (b *base) factoryInit(jsonFactory factory, path, parent string, positions []int) error {

	if len(positions) > 30 {
		return fmt.Errorf("0004 recursive set initialisation file detected, the maximum dotpath depth of 30 has been reached")
	}

	for i, f := range jsonFactory.Include {
		fileBytes, errHTTP := b.authBody.Decode(f.URI)
		// run the inputPath without a the end of ../ to trim the file

		var err error
		// generate the input path per run to stop overwriting errors
		inputPath := path
		if errHTTP != nil {
			// check again as an extension of the url
			oldpath := inputPath
			inputPath, _ = url.JoinPath(inputPath, f.URI)
			// inputPath = filepath.Clean(filepath.Join(inputPath, f.URI))

			fileBytes, errHTTP = b.authBody.Decode(inputPath)

			// then check the local files
			if errHTTP != nil {
				//retry the file as a path
				inputPath = filepath.Join(oldpath, f.URI)
				inputPath, _ = filepath.Abs(inputPath)
				// inputPath = filepath.Clean(filepath.Join(inputPath, f.URI))
				fileBytes, err = os.ReadFile(inputPath)
			}
		}

		if err == nil {

			// check if the bytes have children by being a json factory
			var newF factory
			err := yaml.Unmarshal(fileBytes, &newF)
			if err != nil {
				return fmt.Errorf("0005 error parsing %s: %v", inputPath, err)
			}

			if _, ok := b.importedWidgets[parent+f.Name]; ok {
				return fmt.Errorf("0006 the alias %s is repeated, every alias is required to be unique", parent+f.Name)
			} else if _, ok := b.importedFactories[parent+f.Name]; ok {
				return fmt.Errorf("0006 the alias %s is repeated, every alias is required to be unique", parent+f.Name)
			}

			// schema validation to sort between widgets and factories
			factLines := make(validator.JSONLines)
			err = validator.Liner(fileBytes, inputPath, "factory", factLines) // treat it as a factory update
			if err != nil {
				return err
			}

			var validatorError error
			if err := validator.SchemaValidator(incschema, fileBytes, parent, factLines); err != nil {
				// @TODO include a better error handling method
				b.importedWidgets[parent+f.Name] = fileBytes
				validatorError = validator.Liner(fileBytes, inputPath, "widget", b.jsonFileLines)
			} else {
				// schema check here?
				b.importedFactories[parent+f.Name] = newF
				validatorError = validator.Liner(fileBytes, inputPath, "factory", b.jsonFileLines)
			}

			if validatorError != nil {

				return validatorError
			}

			// pass thorugh the factory as it won't run for length 0
			err = b.factoryInit(newF, filepath.Dir(inputPath), parent+f.Name+".", append(positions, i))
			if err != nil { // return the error up the chain
				return err
			}
		} else {

			return err
			// fmt.Println(fmt.Errorf("Error opening %v:, %v\n", p, err))
		}

		b.metadataParams[parent+f.Name] = f.Args
	}

	return nil
}

/*
// The resource search algorthim
func (b *base) FileSearch(URI, path string) ([]byte, error) {
	fileBytes, errHTTP := b.authBody.Decode(URI)


	/* If they find it they use it.

We're searching for, path/name.ext

	1 . look relative to _main.json look for path(main)/path/name.ext
	2. relative to path(parent)/path/name.ext
	3. while parent(path)^n/path/name.ext. 

	look relative to wd of the executable
    
    look relative to the prefix specified in env OPENTSG_HOME
	*/

	/*
	file.Abs(URI + path)
	
	
	// compare to any os.Getwd()
	

	var err error
	// generate the input path per run to stop overwriting errors
	inputPath := path
	if errHTTP != nil {
		// check again as an extension of the url
		oldpath := inputPath
		inputPath, _ = url.JoinPath(inputPath, URI)
		// inputPath = filepath.Clean(filepath.Join(inputPath, f.URI))

		fileBytes, errHTTP = b.authBody.Decode(inputPath)

		// then check the local files
		if errHTTP != nil {
			//retry the file as a path
			inputPath = filepath.Join(oldpath, f.URI)
			inputPath, _ = filepath.Abs(inputPath)
			// inputPath = filepath.Clean(filepath.Join(inputPath, f.URI))
			fileBytes, err = os.ReadFile(inputPath)
		}
	}
	fmt.Println(err)
	return fileBytes, errHTTP
}
*/