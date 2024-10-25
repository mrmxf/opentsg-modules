# OpenTSG-io

## Description

opentsg-io contains the  encoding methods for custom image libraries, which are
not part of the golang/image library. This module is part of open tsg project,
but can be used independently for when you just need an image saved.

These current encoders available are:

- DPX files.
- Tiff files, which are part of the golang library, but this version does not save the
  alpha channel.
- EXR files, there are alternative libraries available that also decode exr files.
- CSV, this saves the red blue green channels as three separate files.

At the moment no decoders are contained within in the repo, but this much change
with the needs of the project.

## Limtations

These functions are made to work with opentsg, so often use `*image.NRGBA64` or `*opentsg-core/colour.NRGBA64`
instead of the image.Image interface, because we do not handle the other images types. If you would
like to contribute to improve the image options available in this package, then please do.

### DPX

The DPX encoder has the following options

- Encoding in 8, 12 and 16 bits
- All DPX files are little endian encoded

### EXR

The EXR encoder saves the files as a float32 EXR file.

### Tiff

This package is designed for images that are completely opaque, any image with a
pixel alpha value of less than a 100% will lose the alpha channel with this package.

These files are big endian encoded.

## Visuals

Depending on what test signal you are making, it can be a good idea to include screenshots
or even a video (you'll frequently see GIFs rather than actual videos). Tools
like ttygif can help, but check out Asciinema for a more sophisticated method.

## Installation

```sh
go get opentsg-io
```

## Usage

An example program can be found below,
this generates an image of 100,100 pixels,
fills it in green then saves it as a DPX file.

```go
package main

import (
  "image" 
  "image/color"
  "image/draw"
  "os"

  "github.com/mrmxf/opentsg-modules/opentsg-io/dpx"
)

func main() {

  // make an image
  canvas := image.NewNRGBA64(image.Rect(0, 0, 100, 100))
  // paint it green
  draw.Draw(canvas, canvas.Bounds(), &image.Uniform{color.RGBA{G: 0xff, A: 0xff}}, image.Point{}, draw.Src)
  //save it as a dpx

  f ,_ := os.Create("myFirst.dpx")
  dpx.Encode(f, canvas, &dpx.Options{Bitdepth: 16})
}

```
