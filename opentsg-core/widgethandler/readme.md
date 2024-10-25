# Widgethandler

Widget handler is the multiplexer for running the widgets and
drawing them in the correct order.

It gets the z order (the order in which the widgets run) across
all the widgets types, then as the widgets are run globally, they
are placed on the test pattern in the order they were declared.

The order is the same given in the factories when the files are
all parsed, only a couple of the widgets are drawn per time to not
steal all the memory in the computer.

It contains the `Generator` interface for external widgets to be added.
This is shown below

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
