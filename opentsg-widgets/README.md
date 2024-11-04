# tsg-widgets

`opentsg-widgets` is the library of builtin widgets,
for the [Open Test Signal Generator](https://opentsg.studio/).
Feel free to use any as a blueprint to make your own widgets,
for use with openTSG.

## Examples

Examples of the json for each type can be found at the `exampleJson` folder.
These examples contain all the fields unique to the widget. The name of the folder
is the same as the `"type"` field that would be declared for that widget. examples that have nested
folders such as builtin.ebu3373/bars, have their widget names including the `/`, so
`builtin.ebu3373/bars` would be the widget type.

The positional fields of `"grid"` are expected to be in **every** widget input file. And have the
layout as shown below. They may not be found in every/any example json.
The widget type is also required. Each widget has a unique `"type"`,
so OpenTSG can identify and use the widget.

```json
"type" : "builtin.example",
"grid": {
    "location": "a1:b2",
    "alias" : "A demo Alias"
}

```

## Widget Properties

This section contains the properties of the widgets.
This contains the design behind the widget, the fields
and contents it uses. And an example JSON

- [AddImage](_docs/addimage/doc.md)
- [Ebu3373](_docs/ebu3373/doc.md)
- [Fourcolour](_docs/fourcolour/doc.md)
- [FrameCount](_docs/framecount/doc.md)
- [Gradients](_docs/gradients/doc.md)
- [Noise](_docs/noise/doc.md)
- [QrGen](_docs/qrgen/doc.md)
- [TextBox](_docs/textbox/doc.md)
- [ZonePlate](_docs/zoneplate/doc.md)

## Notes for developers

There are several stages for developing a widget to use in OpenTSG.

The first is the configuration of the widget. It
requires the following items:

- The struct with the `*config.Grid` and `colour.ColorSpace` fields.
- The json schema of the object, preferably embedded in the code.
(one less file to track)
- The functions to match the widget handler Generator interface.
- Most importantly an idea of what the widget is there to test for and
achieve.

The widget handler Generator interface

```go
// Generator contains the method for running widgets to generate segments of the test chart.
type Generator interface {
    // Generate the widget in the bounds of the image
    Generate(draw.Image, ...any) error
    // Loc returns the location of the grid for  gridgen.ParamToCanvas
    Location() string
    // Alias returns the alias of the grid for  gridgen.ParamToCanvas
    Alias() string
}
```

An example widget would look like this.
This example is a simple widget that fills in a canvas as a solid green.

```go

import (
    _ "embed"

    "github.com/mrmxf/opentsg-modules/opentsg-core/colour"
    "github.com/mrmxf/opentsg-modules/opentsg-core/config"
)


type exampleJSON struct {
  // Required fields
  GridLoc     *config.Grid      `json:"grid,omitempty" yaml:"grid,omitempty"`
  ColourSpace colour.ColorSpace `json:"colorSpace,omitempty" yaml:"colorSpace,omitempty"`

  // Any other fields relating to the test pattern

}

// the example schema (this is not a real file)
//go:embed jsonschema/exampleSchema.json
var exampleSchema []byte

// return the Alias of the location
func (ej exampleJSON) Alias() string {
    return ej.GridLoc.Alias
}

// return the exact location
func (ej exampleJSON) Location() string {
    return ej.GridLoc.Location
}

// It just fills in the canvas as green (for this demo)
func (ej exampleJSON) Generate(canvas draw.Image, opts ...any) error {
    draw.Draw(canvas, canvas.Bounds(),  &image.Uniform{color.RGBA{G: 0xff, A: 0xff}}, image.Point{}, draw.Over)
    return nil
}
```

Now we have the groundwork for the widget we need to tie it
altogether in an exported function, so that it can be used by OpenTSG.
Like this function given below.

```go

func ExampleGenerate(canvasChan chan draw.Image, debug bool, c *context.Context, wg, wgc *sync.WaitGroup, logs *errhandle.Logger) {
    defer wg.Done()

    // set up the configuration so ExampleJson is identified
    conf := widgethandler.GenConf[exampleJSON]{Debug: debug, Schema: exampleSchema, WidgetType: "example"}
    widgethandler.WidgetRunner(canvasChan, conf, c, logs, wgc) 
}
```

which can be added to tsg by calling the following code,
before the TSG object is run. Or if you want to add it to the standard library of
openTSG, then add it to line 222 of `/opentsg-modules/opentsg-core/tsg/framedraw.go`

The below is a demo of adding an external widget
without changing the opentsg library.

```golang
    opentsg, configErr := tsg.FileImport(commandInputs, *profile, *debug, myFlags...)
    //handle configERR

    // Add the customWidget here!
    opentsg.AddCustomWidgets(example.ExampleGenerate)

    // run opentsg
    opentsg.Draw(*debug, *outputmnt, *outputLog)
```
