# Prepare local environment for K8s development on Ubuntu

Refer to K8s community [development guide.](https://github.com/kubernetes/community/blob/master/contributors/devel/development.md)

- Hardware requirement: 8 GB RAM, 50 GB disk
- Install dependency 
  - Install GNU Development Tools
    ```sh
    sudo apt update
    sudo apt install build-essential
    ```
  - Install [docker](https://docs.docker.com/engine/install/ubuntu/)
  - Install jq, a command line json processor
    ```sh
    sudo apt-get install jq
    ```
  - Install Go
    ```sh
    # download
    curl -O https://storage.googleapis.com/golang/go1.15.5.linux-amd64.tar.gz
    # install (unzip and place to /usr/local)
    tar -C /usr/local -xzf go1.15.5.linux-amd64.tar.gz
    # export path, add the following to ~/.bashrc to make it permanent 
    export PATH=$PATH:/usr/local/go/bin
    # set up go env variable
    go env -w GOBIN=/usr/local/go/bin
    go env -w GOPATH=/usr/local/go
    ```
  - Install etcd, to test Kubernetes, a key-value store etcd is needed, after download Kubernetes source code, run the following
    ```sh
    ./hack/install-etcd.sh
    export PATH="$GOPATH/src/k8s.io/kubernetes/third_party/etcd:${PATH}"
    ```
 - Fork Kubernetes
