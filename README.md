# knuu-example

This repository contains an example of how to use the [knuu](https://github.com/celestiaorg/knuu) Integration Test Framework.

If you have feedback on the framework, want to report a bug or suggest an improvement, please create an issue in the [knuu](https://github.com/celestiaorg/knuu) repository.

## Setup

1. Install required packages

    ```shell
    sudo apt-get install buildah pkg-config libgpgme-dev libdevmapper-dev btrfs-progs libbtrfs-dev
    ```

2. Set up access to a Kubernetes cluster using your kubeconfig and create the `test` namespace.

## Run

```shell
go run .
```
