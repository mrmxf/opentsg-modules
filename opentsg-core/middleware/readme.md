# Middleware

Middleware contains the functions for processing the
canvas after all the widgets have been generated.
This is still under development.

It is enabled in the `"builtin.canvas"`widget,
these are all disabled by default as its main use is for debugging.

It has the following fields:

- configs - the input configuration is saved.
- Average - the average image colour of the frame.
- PHash - the [phash](https://en.wikipedia.org/wiki/Perceptual_hashing) of the frame

The following json would enable the middleware to run.

```javascript
"frame analytics" : {
    "configuration": {"enabled":true},
    "average color": {"enabled":true},
    "phash": {"enabled":true}
}
```
