# Gridgen

Gridgen handles the grid generation of the test card.
It handles the grid coordinate system. Scales the rows and columns
to generate grids that fit the size of the canvas.

It generates the initial grid pattern, that is used as the base
test pattern.
It generates the coordinates and
that the widgets use, as well as applying
any geometry from TPIGs.

It also has the art key functionality, to allow for different
backgrounds to be used, where the background is to be preserved.

## TPIG (Test Pattern Input Geometry)

TPIGs are used for generating test patterns to fit 3d shapes.
They can be included within the canvas options json, with
the following code.

```json
{
    "type": "builtin.canvasoptions",
    ...
    "geometry":"./path/to/TPIG.json"
}

```

The `"geometry"` field specifies the TPIG file, no other input is required
to include TPIGs.

### How do TPIGS work

The geometry is the first thing openTSG checks when generating the base test pattern.
If a TPIG is present, the TPIG is unwrapped and sets the layout of the test pattern
that can be written to. The OpenTSG grid system still overlays this base image.

OpenTSG then runs normally, except that any pixels that aren't within the TPIG boundaries
are not filled in.

After OpenTSG has run the image is saved as normal, unless carving is required.
If carving is required then the image is then carved up into the smaller images
and each carved image is saved individually,
as {filename}{carve}.ext. The flat image is also saved
Find out more about carving with TPIGS [here](#carving)

### What is in a TPIG

A TPIG is a json file, documenting the flat layout, any
carving required and individual tiles of the object.

It has the following layout:

```json
{
    "Tile layout": [
        {
            "Name": "A000",
            "Tags": [],
            "Layout": {
                "Carve": {
                    "Destination": "C1",
                    "X": 0,
                    "Y": 0
                },
                "Flat": {
                    "X": 0,
                    "Y": 0
                },
                "XY": {
                    "X": 10,
                    "Y": 10
                }
            }
        }
    ],
    "Dimensions": {
        "Flat": {
            "X0": 0,
            "Y0": 0,
            "X1": 30,
            "Y1": 30
        }
    },
    "Carve": {
        "C1": {
            "X0": 0,
            "Y0": 0,
            "X1": 30,
            "Y1": 30
        }
    }
}
```

The `"Tile layout"` field contains a data array for each tile face. It contains

- `"Name"` - the name of that polygon, often used for labelling. - OPTIONAL
- `"Tags"` - any tags associated with that polygon, see [here](#tpig-tags) for more information about tags. - OPTIONAL
- `"Layout"` - the XY coordinate layout of the face. It has the following sub fields:
  - `"Carve"` - the name of the carve destination, and the XY coordinates - OPTIONAL
  - `""Flat"` - The initial XY coordinates of the object in its flat layout. - REQUIRED
  - `"XY"` - The height and width of the face. - REQUIRED

For carve and flat only the initial XY coordinates are required,
as the XY width and height gives the size of the tile.

An example `"Tile layout"` is given below.

```json
{
            "Name": "A000",
            "Tags": [],
            "Layout": {
                "Carve": {
                    "Destination": "C1",
                    "X": 0,
                    "Y": 0
                },
                "Flat": {
                    "X": 0,
                    "Y": 0
                },
                "XY": {
                    "X": 10,
                    "Y": 10
                }
            }
        }
```

The `"Dimensions"` contains the fields for the flat object dimensions.
The example below has a size of 30 by 30 pixels, X0,Y0 are the minimum coordinates
and X1,Y1 represent the maximum coordinates.

```json
"Dimensions": {
        "Flat": {
            "X0": 0,
            "Y0": 0,
            "X1": 30,
            "Y1": 30
        }
    }
```

The `"Carve"` field is a map of carve name and its size, to allow
different carve destinations to have different dimensions.

e.g.

```json
 "Carve": {
        "Carve1": {
            "X0": 0,
            "Y0": 0,
            "X1": 30,
            "Y1": 30
        },
        "Carve2": {
            "X0": 0,
            "Y0": 0,
            "X1": 550,
            "Y1": 2000
        }
    }
```

The example above has two carve destinations, carve1 and carve2 of varying sizes.
Where X0,Y0 are the minimum coordinates
and X1,Y1 represent the maximum coordinates.

Please note coordinates, start at (0,0) in the top left of the canvas
and go up to (max_X,max_Y) in the bottom right.

### TPIG tags

Each face has a tag field, for adding additional information about the face.
Some of these tags are used by the widgets, the ones currently in use are:

- `"neighbour:example"` - this is used to signify what the polygon's neighbours are. It is currently used by the four colour widget

You can add your own custom tags to describe extra info about the tile, as long
as it does not conflict with existing fields.

### Carving

Carving is required when several UV maps are required to put a picture on an object.
It is used for very large or complex objects,
where several inputs are required to light up different
sections of the display.

To implement this in a TPIG the carve field needs to be declared
with the name of the carved image and its dimensions like so.

In the example below the carve has one file called C1,
it has dimensions of a width of 300 and a height of 300.

```json
"Carve": {
        "C1": {
            "X0": 0,
            "Y0": 0,
            "X1": 300,
            "Y1": 300
        }
    }
```

Then each tile to be carved requires a `"Carve"` field with the `"Destination"`
of the carve matching the name of a carve location,
and the XY coordinate within the carved image e.g.

In the example below, the tile has a flat layout at (0,0) and
when it is carved it is assigned C1 at (10,10)

```json
 {
            "Name": "A000",
            "Tags": [],
            "Layout": {
                "Carve": {
                    "Destination": "C1",
                    "X": 10,
                    "Y": 10
                },
                "Flat": {
                    "X": 0,
                    "Y": 0
                },
                "XY": {
                    "X": 10,
                    "Y": 10
                }
            }
        }
```

### TPIGS and widgets

Some widgets are designed for use with TPIG as part of the input, for
example the four colour widget.
When these widgets are run without a TPIG, the grid coordinate system
is used as the input geometry instead.

The complete list of TPIG widgets is given below.

- [Four colour](../../../opentsg-widgets/_docs/fourcolour/doc.md)

### Building a TPIG

These are the steps we undertook to build a simple TPIG
for a house in the [TPIG demo](https://github.com/mrmxf/opentsg-node/blob/main/READMETPIG.md).

We started with making an obj of a cube.
We ignored the normals and only focused on geometric vertices and texture
coordinates of the tiles.

The first step was making the coordinates of the cube,
so that it unwrapped as shown in the readme. This did not work first time
as it turns out blender has to be told what axis is up, when importing objs,
as there is no standard.
Blenders default is y up, but the obj was designed for z to be the up axis.

This leads to many minutes pondering if the axis lied to me in the obj format,
but it was only blender making everything wonky.

Next was getting the uv map correct, so when it was unfolded the textures lined up.
This meant that tiles textures aren't reversed or upside down, due to the
u,v coordinates not lining up.

This meant going through the UV map ensuring coordinates
for neighbouring vertices are the same. And ensuring the
sizes of the map scales correctly to the physical sie of
the display, to prevent squeezing or stretching of the image.
e.g. left panel top right matches the top left of the next panel.

Once the obj then had a correct UV map and the tiles
all lined up as intended we moved onto the TPIG design.

The TPIG transferred the UV map to json,
gave a scale to set the base image, we took squares
to be 1080x1080 pixels. For example a UV coordinate
of 0.25,0.3333 then transferred to a tpig coordinate
of 1080,1080, for the
unwrapped cube having a length of 4320 and a height
of pixels.

Each tiles uv map was then scaled to give pixel precise
values. No carving was added as the unwrapped image
was the uv map. If the uv map differs from the flat then carving
would have been added.
