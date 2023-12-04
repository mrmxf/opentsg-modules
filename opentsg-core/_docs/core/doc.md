# Core

## Factories and M

The input file for OpenTSG is called a factory, and can contain 1 or more references to other files.
With nesting available for the files.

When generating the widgets, the factories are processed in a depth first manner. That means every time a URI is
encountered its children and any further children are processed, before its siblings in the factory.

Each factory or widget declares which metadata keys it uses, with the "args" key
(this can be no keys).
On the generation of the widgets and factories the base metadata values
for every unique dot path are set using these keys.
This is where metadata is split from the inline update and stored in the metadata "bucket".
This base metadata "bucket" is not overwritten by later updates and is
generated on a per frame basis. It is used for applying metadata
updates to the widgets.
The workflow is the widget gets its argument keys, it searches these
keys in the metadata bucket of its parents, overwriting more generic
metadata with more specific as you proceed along the parents.
Locally declared metadata for the update will then overwrite this base metadata layer.

Wdigets can inherit any metadata that matches the declared argument keys, from their parents.
With more specific metadata overwriting previous values.

Then as the dotpath and array updates are applied, they will use these metadata values, unless
a new metadata value is called as part of that dot path.

The input factory does not have declared metadata.
