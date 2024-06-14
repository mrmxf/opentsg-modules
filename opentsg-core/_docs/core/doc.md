# Core

## Factories

The input file for OpenTSG is called a factory, and can contain 1 or more references to other files.
With nesting available for the files.

When generating the widgets, the factories are processed in a depth first manner. That means every time a URI is
encountered its children and any further children are processed, before processing its siblings in the factory.

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
With more specific metadata overwriting previous values. e.g. The first factory decalres a
title with `{"title":"exmaple"}`, then any child that has the argument `child` will use
value of `"example"`

Then as the dotpath and array updates are applied, they will use these metadata values, unless
a new metadata value is called as part of that dot path.

### Input File Search Order

When referencing importing files using the `"uri":"example.json"` method the file is searched for
in several locations. It takes the following steps:

1. It searches for the path as a url.
2. It searches relative to the main json that was called.
3. It searches in age order (oldest first) the folders that other factory files
were called from.
4. it searches relative to the folder the executable was called from
5. it searches relative to the environment variable `OPENTSG_HOME`

Then if no file is found then an error is returned.

For step 3 the factories are searched depth first, so only folder locations
that are parents of the factory are included. No width based searches occur.

### Metadata Arguments update order

When passing metadata through the input files, the child metadata overwrites
the parent metadata, if applicable. As part of this the metadata is mustached
as you go along. So a parent will declare `"title": "TestTitle"` then the child
will use `"title": "{{title}}-update"`, resulting in a final argument value of
`"title": "TestTitle-update"` reaching the widget. This means arguments can
be built up throughout the initialisation process.

## TPIG


