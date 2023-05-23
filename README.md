# knuu-example

This repository contains an example of how to use the [knuu](https://github.com/celestiaorg/knuu) Integration Test Framework.

If you have feedback on the framework, want to report a bug or suggest an improvement, please create an issue in the [knuu](https://github.com/celestiaorg/knuu) repository.

## Setup

1. Install [Docker](https://docs.docker.com/get-docker/).

2. Set up access to a Kubernetes cluster using your kubeconfig and create the `test` namespace.

## Run

```shell
go test -v ./...
```

Or run only the basic examples:

```shell
go test -v ./basic
```

Or the celestia-app examples:

```shell
go test -v ./celestia_app
```
