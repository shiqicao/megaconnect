# MegaConnect
Reference implementation of the Megaspace InterConnect Protocol [TODO - link].

# Prerequisites
- Install [Go].
  For convenience, also add `$(go env GOPATH)/bin` to your `PATH` env.
- Install [dep].
- [Optional] Install [gRPC go support][grpc-go].
  Required only if you need to update protobuf specs.

# Initial Setup
1. Clone this repo under `$(go env GOPATH)/src/github.com/megaspacelab/megaconnect`.
  (`go get` doens't work very well with private repos.)
1. From repo root, run
   ```
   ./init-dev.sh
   ```
   This downloads all dependencies into `vendor/` and sets up git hooks.

# Build
To compile only, run
```
make
```

To compile and install binaries into Go path, run
```
make install
```

After making changes to protobuf specs (`grpc/*.proto`), compile the new specs by running
```
make protos
```

# Run
Several binaries are produced out of this repo.

## flow-manager
Manages and distributes monitors to chain managers (connectors), and aggregates reports from them.
[TODO - add cli args and usage]

## example-connect
Example connector implementation that periodically reports a fake new block.
```
example-connect --debug
```

## wfc
Workflow compiler.
[TODO - add usage]


[go]: https://golang.org/dl/
[dep]: https://golang.github.io/dep/docs/installation.html
[grpc-go]: https://grpc.io/docs/quickstart/go.html#prerequisites
