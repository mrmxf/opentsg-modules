# tsg

tsg is the opentsg engine that runs all the image generation.

it is create by running

```go
openTSG, err := tsg.BuildOpenTSG(inputFile string, profile string, debug bool, httpKeys ...string)
```

## Customisation

OpenTSG is designed to be customisable, with the ability to include
custom widgets and save functions, without having to touch this repo to
make those changes.

## Adding widgets

When designing the widget code, make sure it follows the layout in
the notes for developers section in the opentsg-widgets [README](./../../../opentsg-widgets/README.md)

```go
    opentsg, configErr := tsg.FileImport(commandInputs, *profile, *debug, myFlags...)
    //handle configErr

    // Add the customWidget here!
    opentsg.AddCustomWidgets(example.ExampleGenerate)
```

## Adding save functions

You can add external save functions with the following lines

First make sure the function matches the SaveFunction format in
NameSave below.

```go
type NameSave struct {
 Extension    string
 SaveFunction func(io.Writer, draw.Image, int) error
}
```

This demo below wraps the go standard library jpeg encoder as
a custom function. This is then added to the openTSG engine,
where it can now save output files with the extension `"jpg"` as a JPEG file.

```go
package example

import (
    "image/jpeg"
    "github.com/mrmxf/opentsg-modules/opentsg-core/tsg"
)

func main () {
    opentsg, configErr := tsg.FileImport(commandInputs, *profile, *debug, myFlags...)
    //handle configErr

    // Add the custom jpeg saver here!
    opentsg.AddCustomSave([]tsg.NameSaves{{Extension: "jpg",SaveFunction: jpegEncode}})

    // run opentsg
    opentsg.Draw(*debug, *outputmnt, *outputLog)
}

// wrap the standard library jpeg encoder
func jpegEncode(w io.Writer, img draw.Image, _ int) error {
    return jpeg.Encode(w, img, &jpeg.Options{Quality:100})
}
```
