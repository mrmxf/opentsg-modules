# tsg

tsg is the opentsg engine that runs all the image generation.

it is create by running

```go
openTSG, err := tsg.BuildOpenTSG(inputFile string, profile string, debug bool, httpKeys ...string)
```

## Customisation

OpenTSG is designed to be customisable, with the ability to include
custom widgets and encode functions, without having to touch this repo to
make those changes.

## Adding widgets

When designing the widget code, make sure it follows the layout in
the notes for developers section in the opentsg-widgets [README](./../../../opentsg-widgets/README.md)

```go
    opentsg, configErr := tsg.FileImport(commandInputs, *profile, *debug, myFlags...)
    //handle configErr

    // Add the customWidget here!
    opentsg.HandleFunc("example.example" ,example.ExampleGenerate)
```

## Adding save functions

You can add external save functions with the following lines

First make sure the function matches the Encoder type in
below.

```go
// Encoder is a function for encoding the openTSG output into a
// specified format.
type Encoder func(io.Writer, image.Image, EncodeOptions) error

// EncodeOptions contains an extra options for encoding a file
type EncodeOptions struct {
    // the target bitdepth an image is saved to
    // only relevant for DPX files
    BitDepth int
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
    opentsg.EncoderFunc("jpg",jpegEncode)

    // run opentsg
    opentsg.Draw(*debug, *outputmnt, *outputLog)
}

// wrap the standard library jpeg encoder
func jpegEncode(w io.Writer, img draw.Image, _ tsg.EncodeOptions) error {
    return jpeg.Encode(w, img, &jpeg.Options{Quality:100})
}
```
