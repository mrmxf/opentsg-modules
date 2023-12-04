# tpg-widgets

`opentsg-widgets` is the library of builtin widgets,
for the [Open Test Signal Generator](https://opentsg.io/).
Feel free to use any as a demo to make your own widgets.

## Examples

Examples of the json for each type can be found at the `exampleJson` folder.
These examples contain all the fields unique to the widget. The name of the folder
is the same as the type that would be declared for that widget. examples that have nested
folders such as builtin.ebu3373/bars, have their widget names including the `/`, so
`builtin.ebu3373/bars` would be the widget type.

The positional fields of grid are expected to be in **every** widget. And have the
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
