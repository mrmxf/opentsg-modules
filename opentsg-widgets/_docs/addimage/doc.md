# Add Image

Add image adds an image onto the test pattern.

It has the following required fields:

- `image` - the local or online location of image to be added.

And the following optional fields:

- `imageFill` - the fill type of the image, which are:
  - `fill` - Stretch the image to fill the X and Y axis
  - `x scale` - Scale the image to fit the X axis
  - `y scale` - Scale the image to fit the Y axis
  - `xy scale` -  Scale the image to fit the X or Y axis, which ever one is reached first.

```json
{
    "type" :  "builtin.addimage",
    "image": "./exmaple/exampleImage.png",
    "imageFill": "fill",
    "grid": {
      "location": "a1",
      "alias" : "A demo Alias"
    }
}
```
