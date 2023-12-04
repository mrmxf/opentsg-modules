# ACES

ACES is custom image library for aces image. It use native go float 32
for pixel values instead of float16 as go does not currently support it.
These are converted to float26 when saved in as EXR files.

ACES image are compatible with the go image library.
