# knuu-example

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
