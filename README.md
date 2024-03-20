# knuu-example

This repository contains an example of how to use the [knuu](https://github.com/celestiaorg/knuu) Integration Test Framework.

If you have feedback on the framework, want to report a bug or suggest an improvement, please create an issue in the [knuu](https://github.com/celestiaorg/knuu) repository.

## Setup

1. Install [Docker](https://docs.docker.com/get-docker/).

2. Set up access to a Kubernetes cluster using your kubeconfig and create the `test` namespace.
> **Note:** The used namespace can be changed by setting the `KNUU_NAMESPACE` environment variable.


## Write Tests

You can find the relevant documentation in the `pkg/knuu` package at: https://pkg.go.dev/github.com/celestiaorg/knuu

## Run

```shell
make test-all
```

Or run only the basic examples:

```shell
make test-basic
```

Or run BitTwister tests:

```sh
make test-bittwister-packetloss
make test-bittwister-bandwidth
make test-bittwister-latency
make test-bittwister-jitter
```

Or the celestia-app examples:

```shell
make test-celestia-app
```

Or the celestia-node examples:

```shell
make test-celestia-node
```
