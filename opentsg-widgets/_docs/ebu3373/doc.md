# EBU3733

EBU 3373 is a test pattern for testing HD Hybrid Log Gamma workflows,
for OpenTSG the individual segments of the test pattern have been separated
into widgets , so that you can use just the components you require.
More technical details of each segment
can be found [here.](https://tech.ebu.ch/publications/tech3373)

## Bars

The ebu3373 bars are generated to fill the space.
The bars are in 10 bit BT2020 colour space, and contain
colour used for testing conversions to
Rec 709.

It has no fields, only the location is required.

```json
{
    "type":"builtin.ebu3373/bars",
    "grid": {
      "location": "a1",
      "alias" : "A demo Alias"
    }
}
```

## Luma

A horizontal luma ramp from the 10 bit value 4 to 1019,
is generated to fill the widget area. It is used for
checking linear mapping functions and for clipping
of lights and darks.

It has no fields, only the location is required.

```json
{
    "type":"builtin.ebu3373/luma",
    "grid": {
      "location": "a1",
      "alias" : "A demo Alias"
    }
}
```

## Near black

Alternating segments of 10 bit 0% black (RGB10(64,64,64)) and
-4%, -2%, -1%, 1%, 2%, 4% black. This is for checking
sub blacks are not removed in the production chain.
This widget is generated to fill the widget area.

It has no fields, only the location is required.

```json
{
    "type":"builtin.ebu3373/nearblack",
    "grid": {
      "location": "a1",
      "alias" : "A demo Alias"
    }
}
```

## Saturation

The saturation generates the steps of red, green and blue
from minimum to maximum saturation. These saturations
are used to establish the colour space of the viewing device.

It has the following optional fields

- `colors` - An array of the "red", "green" and "blue". The saturation appears in the order
they were declared. If only one colour is chosen then only that colour is used. If `colors` 
is not declared then the default order of, red,green and blue is used.

```json
{
    "type":"builtin.ebu3373/saturation"
    "colors": [
        "red",
        "green",
        "blue"
    ],
    "grid": {
      "location": "a1",
      "alias" : "A demo Alias"
    }
}
```

## Two SI

Two Si is the two sample interleave pattern, for when SMPTE ST 425-5:2019
is used. The order of the four 3G-SDI video cables can be checked.
It consists of 8 combos of the four ABCD cables and lines. Each letter and
set of lines are linked to a cable. IF they don't match the expected
layout then the cables are not correctly ordered.

It has no fields, only the location is required.

```json
{
    "type": "builtin.ebu3373/twosi",
    "grid": {
      "location": "a1",
      "alias" : "A demo Alias"
    }
}
```
