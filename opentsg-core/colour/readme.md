# ColourSpace Documentation

The colour library allows the images to use [colour spaces](https://en.wikipedia.org/wiki/Color_space),
and more importantly transform between the colour spaces for different images.

If a colour space is declared  at any point by a widget or the canvas
than the images associated with it use the colour space functionality.
This is an *image.NRGBA64 wrapped with a
colour space, so that if a colour from a different space is added to the canvas
it can be transformed to match the destination colour space.

colour space aware colours are also available as colour.CNRGBA64,
which are the same as color.NRGBA64 but with a colour space field.
When setting a colourspace aware image with a colour space aware colour, the colours are transformed,
to match the destination (image) colourspace, on a per pixel basis. This means each pixel is transformed
individually to remain as accurate as possible.

These colourspace transformations are included in the NRGBA64 `Set` method and when using the
 `colour.Draw` and `colour.DrawMask` functions. These work the same as the image.Draw library,
 with the additional checks for colour space.

## Using colour spaces

To add colour space to a widget the following json is required.

```javascript
"props":{
    "ColorSpace": {"ColorSpace" : "rec709"}
}
```

Current available colour spaces are :

- `"rec709"`
- `"rec2020"`
- `"p3"`
- `"rec601"`

Different transformation methods will be implemented in the future,
such as look up tables. Currently only matrix transformations are used.
to go from RGB to XYZ space and then to XYZ to RGB.
