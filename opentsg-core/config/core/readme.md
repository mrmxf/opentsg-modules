# Core

The core module contains all the core components of OpenTSG.
It organises how the configuration factories are passed
and how the global context that supports OpenTSG is created
and handled.

## Factories

The input file for OpenTSG is called a factory, and can contain
1 or more references to other input files.
These input files can be nested several layers deep, (the default max is 30 to
prevent infinite recursion).

When generating the widgets, the factories are processed in a depth first manner.
That means every time a URI is encountered its children and any further children
are processed, before processing its siblings in the factory. This is
the order in which widgets will be run when generating the test pattern.

This method also generates a unique dotpath id for every widget.
As the names have to be unique for each include statement in a factory,
their children inherit this name as part of the name dotpath, building a
dot path that is unique to each widget.

### Factory arguments

Arguments are available to allow substitutions in the factories and
reduce the amount of input files and repeated code.

Each factory or widget declares which metadata keys it uses, with the "args" key
(this can be no keys). This would be implemented in the json like so.

```javascript
"include": [
    {"uri":"canvas.json", "name":"canvas",  "args":["title","update"]}
]
```

When the widgets and factories are first parsed by openTSG,
each widget is assigned the metadata keys from its parents.
Next the json updates in the `"create"` field are handled.

This metadata keys and fields are split from the update and stored in the metadata "bucket".
This base metadata "bucket" is not overwritten by later updates and is
generated on a per frame basis. It is used for applying metadata
updates to the widgets.

The widget gets its argument keys, it searches these
keys in the metadata bucket of its parents, overwriting more generic
metadata with more specific as you proceed along the parents to the child.
Locally declared metadata for the update will then overwrite this base metadata layer.

Wdigets can inherit any metadata that matches the declared argument keys, from their parents.
With more specific metadata overwriting previous values. e.g. The first factory decalres a
title with `{"title":"exmaple"}`, then any child that has the argument `title` will use
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

## Widget run order

The order widgets run in can be set with the create function of nested factories.
The order in which they are declared in the array is the order in which they run.
If several widgets are declared in the same position then their order will be
random, so make sure they have no overlap when doing this.
This is because the top level create array declares what factories
 are used in each frame.
