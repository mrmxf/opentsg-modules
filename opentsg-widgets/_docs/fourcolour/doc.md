# Four Colour

Is the four colour algorithm for TPIGS.Where four
colours are the maximum number of colours to
fill a map. This widget allows 4+ colours to be used
in the interest of computational speed for large objects.

IT has the following required field:

- `colors` - a list of TSG colors, it must be an array
of four strings or more. The algorithm will use the fewest amount of colors required.

```json
{
    "type" :  "builtin.fourcolour",
    "colors": [
        "#FF0000",
        "#00FF00",
        "#0000FF",
        "#FFFF00"
    ],
    "grid": {
      "location": "a1",
      "alias" : "A demo Alias"
    }
}
```
