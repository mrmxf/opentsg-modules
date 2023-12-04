# qr gen

Generates a qr code from the user input.
It has the following required fields

- `code` - the text to be made into a qr code

It has the following optional fields:

- `gridPosition` - the relative x,y positions as percentages
of the grid the inhabit.

```json
{
    "type" :  "builtin.qrcode",
    "code": "https://opentsg.io/",
    "gridPosition": {
        "x":0,
        "y": 0
    },
    "grid": {
      "location": "a1",
      "alias" : "A demo Alias"
    }
}
```
