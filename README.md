# OpenTSG-Modules

Welcome to opentsg-modules
This contains all the components for [open tsg](https://opentsg.studio/) to run.

Please see the individual documentation for each module:

- [tsg-core](opentsg-core/README.md)
- [tsg-widgets](opentsg-widgets/README.md)
- [tsg-io](opentsg-io/README.md)

## Dev Build Notes

Build the program with the following command.
Make sure you have the latest version of go installed

```cmd
go build
```

```cmd
./opentsg-modules --c node-example/loader.json
```

The configuration file of `node-example/loader.json` imports
`node-example/order.json`, which sets the order the widgets are generated in.
