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
    # set up go env variable, add the following to ~/.bashrc to make it permanent
    export GOROOT=/usr/local/go
    export GOPATH=$HOME/go   # location for future go get libararies
    export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
    ```
 - Fork Kubernetes from https://github.com/kubernetes/kubernetes
 - Clone to $HOME/go/src/k8s.io, 1.1 GB
   ```sh
   git clone git@github.com:Fizzbb/kubernetes.git
   ```
 - Add and fetch upstream in $HOME/go/src/k8s.io/kubernetes
   track kubernetes master's changes, not your local forked version
   ```sh
   git remote add upstream https://github.com/kubernetes/kubernetes.git
   git fetch upstream
   git branch --set-upstream-to=upstream/master master
   ```
   Branch 'master' set up to track remote branch 'master' from 'upstream'.
 - Download other needed go utilities
   ```sh
   go get -u github.com/jteeuwen/go-bindata/go-bindata
   go get -u github.com/cloudflare/cfssl/cmd/...
   ```
 - Install etcd, to test Kubernetes, a key-value store etcd is needed
   ```sh
   $GOPATH/src/k8s.io/kubernetes/hack/install-etcd.sh
   export PATH=$PATH:$GOPATH/src/k8s.io/kubernetes/third_party/etcd
   ```
- Finally, build and start your local k8s cluster 
  ```sh
  $GOPATH/src/k8s.io/kubernetes/hack/local-up-cluster.sh
  ```
  Got cannot touch file or pemission denied error for binary in $GOPATH/src/k8s.io/kubernetes/_output/bin
