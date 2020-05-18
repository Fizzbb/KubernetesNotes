# Kubernetes

## In a nutshell
- Everything in Kubernetes is about running Pods.
- How to run is has been declaratively written in the yaml files
- In the yaml file, "Kind" is like the function, "Spec" is like the inputs; You can define new "kind" as CRD to extend K8s functions

## Limitations
- Resource management: K8s understand & manage CPU and memory well, but for other types of resource, e.g. network bandwidth, disk I/O operations are not supported. 
- Scheduler: ensure workloads are fairly distributed at the current time, but if more nodes are avaliable later, no auto rescheduling.
- Size: Clustering requires communication between nodes, the number of possible communication paths, and cumulative load on the underlying database, grows exponetially with the size of cluster. Now limit to 5000 nodes.
