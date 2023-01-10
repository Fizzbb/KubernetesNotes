# Fission installation

### 1. Pods are installed in ```fission```, ```fission-function```, ```fission-builder``` three namespaces
  * ```kubectl create ns fission```
  * Instead of ```kubectl apply -f fission-all-in-one.yaml -n fission```, set ```kubectl config set-context --current --namespace=fission```
  * Then ```kubectl apply -f fission-all-v1.17.0.yaml```

### 2. Storage-svc pod stuck at pending
 * in the fission all in one yaml, add ```storageClassName: manual``` to the  ```fission-storage-pvc```, and create a local or nfs PV manually. then the PVC can pickup the PV, and the storage PVC service pod can start normally. Otherwise it will be stuck at pending.
 * PV is a cluster wide resource, i.e., no namespace is needed. Sample PV is here.
 * PV and PV's storageClassName needs to be the same (here we chose manual), in order to auto claim
 * Storage reclaim policy: retain, menas after release (PVC deletion), the PV will not be auto reused by new PVC, data is retained. choose recycled for auto reuse.
 
