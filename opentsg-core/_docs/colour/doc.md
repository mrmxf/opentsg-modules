# ColourSpace Documentation

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
