package noderesources

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	schedulerconfig "k8s.io/kube-scheduler/config/v1"
	"k8s.io/kubernetes/pkg/scheduler/framework"

	"sigs.k8s.io/scheduler-plugins/pkg/apis/config"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Allocatable is a score plugin that favors nodes based on their allocatable
// resources.
type Allocatable struct {
	handle    framework.Handle
	clientset *kubernetes.Clientset
	resourceAllocationScorer
}

var _ = framework.ScorePlugin(&Allocatable{})

// AllocatableName is the name of the plugin used in the Registry and configurations.
const AllocatableName = "NodeResourcesAllocatable"

// Name returns name of the plugin. It is used in logs, etc.
func (alloc *Allocatable) Name() string {
	return AllocatableName
}

func validateResources(resources []schedulerconfig.ResourceSpec) error {
	for _, resource := range resources {
		if resource.Weight <= 0 {
			return fmt.Errorf("resource Weight of %v should be a positive value, got %v", resource.Name, resource.Weight)
		}
		// No upper bound on weight.
	}
	return nil
}

// Score invoked at the score extension point.
func (alloc *Allocatable) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	klog.V(5).InfoS("Alnair scheduler is working on score plugin, pod name:", pod.Name)
	//patch timestamp annotation on pod
	err := UpdatePodAnnotations(alloc.clientset, pod)
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("cannot patch timestamp to pod %s, err: %v", pod.Name, err))
	}
	klog.V(5).InfoS("Alnair add annotation to pod ", pod.Name)
	nodeInfo, err := alloc.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("getting node %q from Snapshot: %v", nodeName, err))
	}

	// alloc.score favors nodes with least allocatable or most allocatable resources.
	// It calculates the sum of the node's weighted allocatable resources.
	//
	// Note: the returned "score" is negative for least allocatable, and positive for most allocatable.
	return alloc.score(pod, nodeInfo)
}

// ScoreExtensions of the Score plugin.
func (alloc *Allocatable) ScoreExtensions() framework.ScoreExtensions {
	return alloc
}

func clientsetInit() (*kubernetes.Clientset, error) {
	// creates the in-cluster config

	var (
		config *rest.Config
		err    error
	)
	config, err = rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", "/etc/kubernetes/scheduler.conf") //use the default path for now, pass through arg later
	}
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset, err

}

func UpdatePodAnnotations(clientset *kubernetes.Clientset, pod *v1.Pod) error {
	//dont use deep copy to newpod and update, will copy object version, casue the following error
	//err: Operation cannot be fulfilled on pods "XXX": the object has been modified; please apply your changes to the latest version and try again
	patchData := map[string]interface{}{"metadata": map[string]map[string]string{"annotations": {
		"Scheduler-Timestamp": fmt.Sprintf("%d", time.Now().UnixNano())}}}
	//patchData := {"metadata": {"annotations": {"Scheduler-TimeStamp": fmt.Sprintf("%d", time.Now().UnixNano())}}

	namespace := pod.Namespace
	podName := pod.Name

	playLoadBytes, _ := json.Marshal(patchData)
	_, err := clientset.CoreV1().Pods(namespace).Patch(context.TODO(), podName, types.StrategicMergePatchType, playLoadBytes, metav1.PatchOptions{})

	if err != nil {
		klog.V(5).ErrorS(err, "Alnair Pod Patch fail")
		return fmt.Errorf("Alnair %v pod Patch fail %v", podName, err)
	}

	return nil
}

// NewAllocatable initializes a new plugin and returns it.
func NewAllocatable(allocArgs runtime.Object, h framework.Handle) (framework.Plugin, error) {
	// Start with default values.
	klog.V(5).InfoS("alnair scheduler plugin is enabled")
	mode := config.Least
	resToWeightMap := defaultResourcesToWeightMap

	// Update values from args, if specified.
	if allocArgs != nil {
		args, ok := allocArgs.(*config.NodeResourcesAllocatableArgs)
		if !ok {
			return nil, fmt.Errorf("want args to be of type NodeResourcesAllocatableArgs, got %T", allocArgs)
		}
		if args.Mode != "" {
			mode = args.Mode
			if mode != config.Least && mode != config.Most {
				return nil, fmt.Errorf("invalid mode, got %s", mode)
			}
		}
		if len(args.Resources) > 0 {
			if err := validateResources(args.Resources); err != nil {
				return nil, err
			}

			resToWeightMap = make(resourceToWeightMap)
			for _, resource := range args.Resources {
				resToWeightMap[v1.ResourceName(resource.Name)] = resource.Weight
			}
		}
	}
	cs, err := clientsetInit()
	if err != nil {
		return nil, fmt.Errorf("Alnair Cannot initialize in-cluster kubernetes config")
	}
	return &Allocatable{
		handle: h,
		resourceAllocationScorer: resourceAllocationScorer{
			Name:                AllocatableName,
			scorer:              resourceScorer(resToWeightMap, mode),
			resourceToWeightMap: resToWeightMap,
		},
		clientset: cs,
	}, nil
}

func resourceScorer(resToWeightMap resourceToWeightMap, mode config.ModeType) func(resourceToValueMap, resourceToValueMap) int64 {
	return func(requested, allocable resourceToValueMap) int64 {
		// TODO: consider volumes in scoring.
		var nodeScore, weightSum int64
		for resource, weight := range resToWeightMap {
			resourceScore := score(allocable[resource], mode)
			nodeScore += resourceScore * weight
			weightSum += weight
		}
		return nodeScore / weightSum
	}
}

func score(capacity int64, mode config.ModeType) int64 {
	switch config.ModeType(mode) {
	case config.Least:
		return -1 * capacity
	case config.Most:
		return capacity
	}

	klog.V(10).InfoS("No match for mode", "mode", mode)
	return 0
}

// NormalizeScore invoked after scoring all nodes.
func (alloc *Allocatable) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	// Find highest and lowest scores.
	var highest int64 = -math.MaxInt64
	var lowest int64 = math.MaxInt64
	for _, nodeScore := range scores {
		if nodeScore.Score > highest {
			highest = nodeScore.Score
		}
		if nodeScore.Score < lowest {
			lowest = nodeScore.Score
		}
	}

	// Transform the highest to lowest score range to fit the framework's min to max node score range.
	oldRange := highest - lowest
	newRange := framework.MaxNodeScore - framework.MinNodeScore
	for i, nodeScore := range scores {
		if oldRange == 0 {
			scores[i].Score = framework.MinNodeScore
		} else {
			scores[i].Score = ((nodeScore.Score - lowest) * newRange / oldRange) + framework.MinNodeScore
		}
	}

	return nil
}
