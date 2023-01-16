## Fission installation

### 1. Control plane pods are installed in ```fission``` namespace, function pods are launched under ```fission-function``` namespaces
  * ```kubectl create ns fission```
  * Instead of ```kubectl apply -f fission-all-in-one.yaml -n fission```, set ```kubectl config set-context --current --namespace=fission```
  * Then ```kubectl apply -f fission-all-v1.17.0.yaml```
  * There a 3rd namespace ```fission-builder``` too. 
### 2. Storage-svc pod stuck at pending
 * in the fission all in one yaml, add [storageClassName: manual](fission-all-v1.17.0.yaml#L60) to the  ```fission-storage-pvc```, and create a local or nfs PV manually. then the PVC can pickup the PV, and the storage PVC service pod can start normally. Otherwise it will be stuck at pending.
 * PV is a cluster wide resource, i.e., no namespace is needed. Sample PV is [here](nfs_pv.yaml).
 * PV and PV's storageClassName needs to be the same (here we chose manual), in order to auto claim
 * Storage reclaim policy: retain, menas after release (PVC deletion), the PV will not be auto reused by new PVC, data is retained. choose recycled for auto reuse.

### 3. install fission cli for easy debug
```
curl -Lo fission https://github.com/fission/fission/releases/download/v1.17.0/fission-v1.17.0-linux-amd64 \
    && chmod +x fission && sudo mv fission /usr/local/bin/
```
### 4. create new fission env, different from fission default [python-env](https://github.com/fission/environments/blob/master/python/Dockerfile)

* Generate a specs folder
```
fission spec init
```
* Copy the [env-gpu.yaml](env-gpu.yaml) to the specs folder, this yaml file has some pariticula syntax, current is a working version with fission 1.17. More examples can be found in [here](https://fission.io/docs/usage/spec/)
* Apply the spec to create a new env
```
fission spec apply
```
* check a new env is created
```
fission env list
```
## Fission Priciples

### 1. A function Pod consists with two containers: Fetcher and Runtimer. Fetcher fetches user function into function pod during specialization stage, Rutime is a container contains necessary lanuage environment to run user function. 

### 2. Function pods are managed by pool manager. a warm pool of pods are pre-launched (default=3).
