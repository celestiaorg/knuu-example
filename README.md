# knuu-example

## Description

This repository contains an example of how to use the [knuu](https://github.com/celestiaorg/knuu) Integration Test Framework.

If you have feedback on the framework, want to report a bug or suggest an improvement, please create an issue in the [knuu](https://github.com/celestiaorg/knuu) repository.

---
## Setup

*You need to have access to a Kubernetes cluster.*

### Linux

1. Install required packages
    

    ```shell
    sudo apt-get install buildah pkg-config libgpgme-dev libdevmapper-dev btrfs-progs libbtrfs-dev
    ```

*Install **kubectl**: [install-kubectl-linux](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)*


2. Set up access to a Kubernetes cluster using your kubeconfig and create the `test` namespace.

### MacOS

1. Install required packages

    ```shell
    brew install gpgme
    ```

*Install **kubectl**: [install-kubectl-macos](https://kubernetes.io/docs/tasks/tools/install-kubectl-macos/)*

2. Set up access to a Kubernetes cluster using your kubeconfig and create the `test` namespace.

---
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

---