# How to customize Kubernetes scheduler 

Using kubernetes schedulering framework, build an out-of-tree scheduler to replace the default kuber-scheduler, based on [Kubernetes SIG scheduler-plugins](https://github.com/kubernetes-sigs/scheduler-plugins)

## Steps


## Catches
1. The ```go.mod``` file, scheduler-plugins depends on kubernetes, but ```go mod tidy``` will include kubernetes's staging dependencys, which have version 0.0.0. Needs to replace manually. or just use the go.mod file from "sigs.k8s.io/scheduler-plugins/"
2. The kube-scheduler deploy yaml file is located at ```/etc/kubernetes/manifests/kube-scheduler.yaml```, Kubernetes check it on the background, any changes in that file will cause scheduler pod re-deployed. If syntax error or unsupported configs are added in the yaml file, the kube-scheduler will disappear, util a valid yaml is updated. So always backup the default yaml before changes.
3. For the ```/etc/kubernetes/manifests/kube-scheduler.yaml```, don't load the scheduler config file (```--config=XXX.yaml```) with configMap. Use hostPath file mount. Although an all-in-one yaml can be created using configMap, but the scheduler cannot be launched in this way. The guess is that only file changes are tracked by Kuberenets not the configMap object.
4. Use ```kube-scheduler``` command arguments ```--v=5``` or some other values to see the debug logs. Corresponding in the code write ```klog.V(5).InfoS("")```
