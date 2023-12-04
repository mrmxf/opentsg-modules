# tpg core

tpg-core contains all the engine and core functionality for running openTSG.

see the documentation for every section here:

- [Aces](_docs/aces/doc.md)
- [Canvaswidget](_docs/canvaswidget/doc.md)
- [Colour](_docs/colour/doc.md)
- [Colourgen](_docs/colourgen/doc.md)
- [Core](_docs/core/doc.md)
- [Credentials](_docs/credentials/doc.md)
- [ErrHandle](_docs/errHandle/doc.md)
- [Gridgen](_docs/gridgen/doc.md)
- [Middleware](_docs/middleware/doc.md)
- [tpg](_docs/tpg/doc.md)
- [widgethandler](_docs/widgethandler/doc.md)

## ColourSpace Documentation

If a colour space is declared than the images use the colour space functionality. This is an *image.NRGBA64 wrapped with a
colour space, so that if a colour from a different space is added to the canvas
it can be transformed using one of the builtin transform functions.
These combine with the CNRGBA64 colours, which are the same as color.NRGBA64 but with a colourspace.
When setting a colourspace aware image with a colour space aware colour, the colours are transformed,
to match the destination (image) colourspace, on a per pixel basis.

This transformations are included in the NRGBA64 Set method and when using the Draw and DrawMask functions.

Different transformation methods will be implemented. Currently only matrix transformations are used.
to go from RGB to XYZ space and then to XYZ to RGB.


To add colour space to opentsg add the following json to the base image json.
Then add the same json code to any widget you would like to use a colourspace.

```json
"ColorSpace": {"ColorSpace" : "rec709"}
```

Current available colour spaces are :

- "rec709"
- "rec2020"
- "p3"
- "rec601"

## Factories and metadata

The input file for OpenTSG is called a factory, and can contain 1 or more references to other files.
With nesting available for the files.

When generating the widgets, the factories are processed in a depth first manner. That means every time a URI is
encountered its children and any further children are processed, before its siblings in the factory.

Each factory or widget declares which metadata keys it uses, with the "args" key
(this can be no keys).
On the generation of the widgets and factories the base metadata values
for every unique dot path are set using these keys.
This is where metadata is split from the inline update and stored in the metadata "bucket".
This base metadata "bucket" is not overwritten by later updates and is
generated on a per frame basis. It is used for applying metadata
updates to the widgets.
The workflow is the widget gets its argument keys, it searches these
keys in the metadata bucket of its parents, overwriting more generic
metadata with more specific as you proceed along the parents.
Locally declared metadata for the update will then overwirte this base metadata layer.

Wdigets can inherit any metadata that matches the declared argument keys, from their parents.
With more specific metadata overwriting previous values.

Then as the dotpath and array updates are applied, they will use these metadata values, unless
a new metadata value is called as part of that dot path.

The input factory does not have declared metadata.

The create array function in the first file. Each object in the array is a frame.
Within sub factories the create, sets the order the widgets are run. E,g if you want a big
picture to run first and then have smaller widgets placed on top it would be first.
