package noderesources

import (
	v1 "k8s.io/api/core/v1"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/features"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	schedutil "k8s.io/kubernetes/pkg/scheduler/util"
)

// resourceToWeightMap contains resource name and weight.
type resourceToWeightMap map[v1.ResourceName]int64

// defaultResourcesToWeightMap is used to set default resourceToWeight map for CPU and memory.
// The base unit for CPU is millicore, while the base using for memory is a byte.
// The default CPU weight is 1<<20 and default memory weight is 1. That means a millicore
// has a weighted score equivalent to 1 MiB.
var defaultResourcesToWeightMap = resourceToWeightMap{v1.ResourceMemory: 1, v1.ResourceCPU: 1 << 20}

// resourceAllocationScorer contains information to calculate resource allocation score.
type resourceAllocationScorer struct {
	Name                string
	scorer              func(requested, allocatable resourceToValueMap) int64
	resourceToWeightMap resourceToWeightMap
}

// resourceToValueMap contains resource name and score.
type resourceToValueMap map[v1.ResourceName]int64

// score will use `scorer` function to calculate the score.
func (r *resourceAllocationScorer) score(
	pod *v1.Pod,
	nodeInfo *framework.NodeInfo) (int64, *framework.Status) {
	node := nodeInfo.Node()
	if node == nil {
		return 0, framework.NewStatus(framework.Error, "node not found")
	}
	if r.resourceToWeightMap == nil {
		return 0, framework.NewStatus(framework.Error, "resources not found")
	}
	requested := make(resourceToValueMap, len(r.resourceToWeightMap))
	allocatable := make(resourceToValueMap, len(r.resourceToWeightMap))
	for resource := range r.resourceToWeightMap {
		allocatable[resource], requested[resource] = calculateResourceAllocatableRequest(nodeInfo, pod, resource)
	}

	score := r.scorer(requested, allocatable)

	if klog.V(10).Enabled() {
		klog.InfoS("Resources and score",
			"podName", pod.Name, "nodeName", node.Name, "scorer", r.Name,
			"allocatableResources", allocatable, "requestedResources", requested,
			"score", score)
	}

	return score, nil
}

// calculateResourceAllocatableRequest returns resources Allocatable and Requested values
func calculateResourceAllocatableRequest(nodeInfo *framework.NodeInfo, pod *v1.Pod, resource v1.ResourceName) (int64, int64) {
	podRequest := calculatePodResourceRequest(pod, resource)
	switch resource {
	case v1.ResourceCPU:
		return nodeInfo.Allocatable.MilliCPU, (nodeInfo.NonZeroRequested.MilliCPU + podRequest)
	case v1.ResourceMemory:
		return nodeInfo.Allocatable.Memory, (nodeInfo.NonZeroRequested.Memory + podRequest)

	case v1.ResourceEphemeralStorage:
		return nodeInfo.Allocatable.EphemeralStorage, (nodeInfo.Requested.EphemeralStorage + podRequest)
	default:
		if schedutil.IsScalarResourceName(resource) {
			return nodeInfo.Allocatable.ScalarResources[resource], (nodeInfo.Requested.ScalarResources[resource] + podRequest)
		}
	}
	if klog.V(10).Enabled() {
		klog.InfoS("Requested resource not considered for node score calculation",
			"resource", resource,
		)
	}
	return 0, 0
}

// calculatePodResourceRequest returns the total non-zero requests. If Overhead is defined for the pod and the
// PodOverhead feature is enabled, the Overhead is added to the result.
// podResourceRequest = max(sum(podSpec.Containers), podSpec.InitContainers) + overHead
func calculatePodResourceRequest(pod *v1.Pod, resource v1.ResourceName) int64 {
	var podRequest int64
	for i := range pod.Spec.Containers {
		container := &pod.Spec.Containers[i]
		value := schedutil.GetRequestForResource(resource, &container.Resources.Requests, true)
		podRequest += value
	}

	for i := range pod.Spec.InitContainers {
		initContainer := &pod.Spec.InitContainers[i]
		value := schedutil.GetRequestForResource(resource, &initContainer.Resources.Requests, true)
		if podRequest < value {
			podRequest = value
		}
	}

	// If Overhead is being utilized, add to the total requests for the pod
	if pod.Spec.Overhead != nil && utilfeature.DefaultFeatureGate.Enabled(features.PodOverhead) {
		if quantity, found := pod.Spec.Overhead[resource]; found {
			podRequest += quantity.Value()
		}
	}

	return podRequest
}
