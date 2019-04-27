package core

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/clusterstate"
	"k8s.io/autoscaler/cluster-autoscaler/clusterstate/utils"
	"k8s.io/autoscaler/cluster-autoscaler/context"
	"k8s.io/autoscaler/cluster-autoscaler/metrics"
	"k8s.io/autoscaler/cluster-autoscaler/simulator"
	"k8s.io/autoscaler/cluster-autoscaler/utils/daemonset"
	"k8s.io/autoscaler/cluster-autoscaler/utils/deletetaint"
	"k8s.io/autoscaler/cluster-autoscaler/utils/drain"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	"k8s.io/autoscaler/cluster-autoscaler/utils/glogx"
	"k8s.io/autoscaler/cluster-autoscaler/utils/gpu"
	kube_util "k8s.io/autoscaler/cluster-autoscaler/utils/kubernetes"
	scheduler_util "k8s.io/autoscaler/cluster-autoscaler/utils/scheduler"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1 "k8s.io/api/extensions/v1beta1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	kube_client "k8s.io/client-go/kubernetes"
	kubeletapis "k8s.io/kubernetes/pkg/kubelet/apis"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
	"k8s.io/klog"
)

const (
	ReschedulerTaintKey = "CriticalAddonsOnly"
)

type podSchedulableInfo struct {
	spec		apiv1.PodSpec
	labels		map[string]string
	schedulingError	*simulator.PredicateError
}
type podSchedulableMap map[string][]podSchedulableInfo

func (psi *podSchedulableInfo) match(pod *apiv1.Pod) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return reflect.DeepEqual(pod.Labels, psi.labels) && apiequality.Semantic.DeepEqual(pod.Spec, psi.spec)
}
func (podMap podSchedulableMap) get(pod *apiv1.Pod) (*simulator.PredicateError, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ref := drain.ControllerRef(pod)
	if ref == nil {
		return nil, false
	}
	uid := string(ref.UID)
	if infos, found := podMap[uid]; found {
		for _, info := range infos {
			if info.match(pod) {
				return info.schedulingError, true
			}
		}
	}
	return nil, false
}
func (podMap podSchedulableMap) set(pod *apiv1.Pod, err *simulator.PredicateError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ref := drain.ControllerRef(pod)
	if ref == nil {
		return
	}
	uid := string(ref.UID)
	podMap[uid] = append(podMap[uid], podSchedulableInfo{spec: pod.Spec, labels: pod.Labels, schedulingError: err})
}
func FilterOutSchedulable(unschedulableCandidates []*apiv1.Pod, nodes []*apiv1.Node, allScheduled []*apiv1.Pod, podsWaitingForLowerPriorityPreemption []*apiv1.Pod, predicateChecker *simulator.PredicateChecker, expendablePodsPriorityCutoff int) []*apiv1.Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	unschedulablePods := []*apiv1.Pod{}
	nonExpendableScheduled := FilterOutExpendablePods(allScheduled, expendablePodsPriorityCutoff)
	nodeNameToNodeInfo := scheduler_util.CreateNodeNameToInfoMap(append(nonExpendableScheduled, podsWaitingForLowerPriorityPreemption...), nodes)
	podSchedulable := make(podSchedulableMap)
	loggingQuota := glogx.PodsLoggingQuota()
	for _, pod := range unschedulableCandidates {
		cachedError, found := podSchedulable.get(pod)
		if found {
			if cachedError != nil {
				unschedulablePods = append(unschedulablePods, pod)
			} else {
				glogx.V(4).UpTo(loggingQuota).Infof("Pod %s marked as unschedulable can be scheduled (based on simulation run for other pod owned by the same controller). Ignoring in scale up.", pod.Name)
			}
			continue
		}
		nodeName, err := predicateChecker.FitsAny(pod, nodeNameToNodeInfo)
		var predicateError *simulator.PredicateError
		if err != nil {
			predicateError = simulator.NewPredicateError("FitsAny", err, nil, nil)
			unschedulablePods = append(unschedulablePods, pod)
		} else {
			glogx.V(4).UpTo(loggingQuota).Infof("Pod %s marked as unschedulable can be scheduled on %s. Ignoring in scale up.", pod.Name, nodeName)
		}
		podSchedulable.set(pod, predicateError)
	}
	glogx.V(4).Over(loggingQuota).Infof("%v other pods marked as unschedulable can be scheduled.", -loggingQuota.Left())
	return unschedulablePods
}
func FilterOutExpendableAndSplit(unschedulableCandidates []*apiv1.Pod, expendablePodsPriorityCutoff int) ([]*apiv1.Pod, []*apiv1.Pod) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	unschedulableNonExpendable := []*apiv1.Pod{}
	waitingForLowerPriorityPreemption := []*apiv1.Pod{}
	for _, pod := range unschedulableCandidates {
		if pod.Spec.Priority != nil && int(*pod.Spec.Priority) < expendablePodsPriorityCutoff {
			klog.V(4).Infof("Pod %s has priority below %d (%d) and will scheduled when enough resources is free. Ignoring in scale up.", pod.Name, expendablePodsPriorityCutoff, *pod.Spec.Priority)
		} else if annot, found := pod.Annotations[scheduler_util.NominatedNodeAnnotationKey]; found && len(annot) > 0 {
			waitingForLowerPriorityPreemption = append(waitingForLowerPriorityPreemption, pod)
			klog.V(4).Infof("Pod %s will be scheduled after low prioity pods are preempted on %s. Ignoring in scale up.", pod.Name, annot)
		} else {
			unschedulableNonExpendable = append(unschedulableNonExpendable, pod)
		}
	}
	return unschedulableNonExpendable, waitingForLowerPriorityPreemption
}
func FilterOutExpendablePods(pods []*apiv1.Pod, expendablePodsPriorityCutoff int) []*apiv1.Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := []*apiv1.Pod{}
	for _, pod := range pods {
		if pod.Spec.Priority == nil || int(*pod.Spec.Priority) >= expendablePodsPriorityCutoff {
			result = append(result, pod)
		}
	}
	return result
}
func CheckPodsSchedulableOnNode(context *context.AutoscalingContext, pods []*apiv1.Pod, nodeGroupId string, nodeInfo *schedulercache.NodeInfo) map[*apiv1.Pod]*simulator.PredicateError {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	schedulingErrors := map[*apiv1.Pod]*simulator.PredicateError{}
	loggingQuota := glogx.PodsLoggingQuota()
	podSchedulable := make(podSchedulableMap)
	for _, pod := range pods {
		if _, repeated := schedulingErrors[pod]; repeated {
			klog.Warningf("Pod %v appears multiple time on pods list, will only count it once in scale-up simulation", pod)
		}
		err, found := podSchedulable.get(pod)
		if found {
			schedulingErrors[pod] = err
			if err != nil {
				glogx.V(2).UpTo(loggingQuota).Infof("Pod %s can't be scheduled on %s. Used cached predicate check results", pod.Name, nodeGroupId)
			}
		}
		if !found {
			err = context.PredicateChecker.CheckPredicates(pod, nil, nodeInfo)
			podSchedulable.set(pod, err)
			schedulingErrors[pod] = err
			if err != nil {
				klog.V(2).Infof("Pod %s can't be scheduled on %s, predicate failed: %v", pod.Name, nodeGroupId, err.VerboseError())
			}
		}
	}
	glogx.V(2).Over(loggingQuota).Infof("%v other pods can't be scheduled on %s.", -loggingQuota.Left(), nodeGroupId)
	return schedulingErrors
}
func GetNodeInfosForGroups(nodes []*apiv1.Node, nodeInfoCache map[string]*schedulercache.NodeInfo, cloudProvider cloudprovider.CloudProvider, kubeClient kube_client.Interface, daemonsets []*extensionsv1.DaemonSet, predicateChecker *simulator.PredicateChecker) (map[string]*schedulercache.NodeInfo, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(map[string]*schedulercache.NodeInfo)
	seenGroups := make(map[string]bool)
	processNode := func(node *apiv1.Node) (bool, string, errors.AutoscalerError) {
		nodeGroup, err := cloudProvider.NodeGroupForNode(node)
		if err != nil {
			return false, "", errors.ToAutoscalerError(errors.CloudProviderError, err)
		}
		if nodeGroup == nil || reflect.ValueOf(nodeGroup).IsNil() {
			return false, "", nil
		}
		id := nodeGroup.Id()
		if _, found := result[id]; !found {
			nodeInfo, err := simulator.BuildNodeInfoForNode(node, kubeClient)
			if err != nil {
				return false, "", err
			}
			sanitizedNodeInfo, err := sanitizeNodeInfo(nodeInfo, id)
			if err != nil {
				return false, "", err
			}
			result[id] = sanitizedNodeInfo
			return true, id, nil
		}
		return false, "", nil
	}
	for _, node := range nodes {
		if !kube_util.IsNodeReadyAndSchedulable(node) {
			continue
		}
		added, id, typedErr := processNode(node)
		if typedErr != nil {
			return map[string]*schedulercache.NodeInfo{}, typedErr
		}
		if added && nodeInfoCache != nil {
			if nodeInfoCopy, err := deepCopyNodeInfo(result[id]); err == nil {
				nodeInfoCache[id] = nodeInfoCopy
			}
		}
	}
	for _, nodeGroup := range cloudProvider.NodeGroups() {
		id := nodeGroup.Id()
		seenGroups[id] = true
		if _, found := result[id]; found {
			continue
		}
		if nodeInfoCache != nil {
			if nodeInfo, found := nodeInfoCache[id]; found {
				if nodeInfoCopy, err := deepCopyNodeInfo(nodeInfo); err == nil {
					result[id] = nodeInfoCopy
					continue
				}
			}
		}
		nodeInfo, err := GetNodeInfoFromTemplate(nodeGroup, daemonsets, predicateChecker)
		if err != nil {
			if err == cloudprovider.ErrNotImplemented {
				continue
			} else {
				klog.Errorf("Unable to build proper template node for %s: %v", id, err)
				return map[string]*schedulercache.NodeInfo{}, errors.ToAutoscalerError(errors.CloudProviderError, err)
			}
		}
		result[id] = nodeInfo
	}
	for id := range nodeInfoCache {
		if _, ok := seenGroups[id]; !ok {
			delete(nodeInfoCache, id)
		}
	}
	for _, node := range nodes {
		if !kube_util.IsNodeReadyAndSchedulable(node) {
			added, _, typedErr := processNode(node)
			if typedErr != nil {
				return map[string]*schedulercache.NodeInfo{}, typedErr
			}
			nodeGroup, err := cloudProvider.NodeGroupForNode(node)
			if err != nil {
				return map[string]*schedulercache.NodeInfo{}, errors.ToAutoscalerError(errors.CloudProviderError, err)
			}
			if added {
				klog.Warningf("Built template for %s based on unready/unschedulable node %s", nodeGroup.Id(), node.Name)
			}
		}
	}
	return result, nil
}
func GetNodeInfoFromTemplate(nodeGroup cloudprovider.NodeGroup, daemonsets []*extensionsv1.DaemonSet, predicateChecker *simulator.PredicateChecker) (*schedulercache.NodeInfo, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	id := nodeGroup.Id()
	baseNodeInfo, err := nodeGroup.TemplateNodeInfo()
	if err != nil {
		return nil, errors.ToAutoscalerError(errors.CloudProviderError, err)
	}
	pods := daemonset.GetDaemonSetPodsForNode(baseNodeInfo, daemonsets, predicateChecker)
	pods = append(pods, baseNodeInfo.Pods()...)
	fullNodeInfo := schedulercache.NewNodeInfo(pods...)
	fullNodeInfo.SetNode(baseNodeInfo.Node())
	sanitizedNodeInfo, typedErr := sanitizeNodeInfo(fullNodeInfo, id)
	if typedErr != nil {
		return nil, typedErr
	}
	return sanitizedNodeInfo, nil
}
func FilterOutNodesFromNotAutoscaledGroups(nodes []*apiv1.Node, cloudProvider cloudprovider.CloudProvider) ([]*apiv1.Node, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make([]*apiv1.Node, 0)
	for _, node := range nodes {
		nodeGroup, err := cloudProvider.NodeGroupForNode(node)
		if err != nil {
			return []*apiv1.Node{}, errors.ToAutoscalerError(errors.CloudProviderError, err)
		}
		if nodeGroup == nil || reflect.ValueOf(nodeGroup).IsNil() {
			result = append(result, node)
		}
	}
	return result, nil
}
func deepCopyNodeInfo(nodeInfo *schedulercache.NodeInfo) (*schedulercache.NodeInfo, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	newPods := make([]*apiv1.Pod, 0)
	for _, pod := range nodeInfo.Pods() {
		newPods = append(newPods, pod.DeepCopy())
	}
	newNodeInfo := schedulercache.NewNodeInfo(newPods...)
	if err := newNodeInfo.SetNode(nodeInfo.Node().DeepCopy()); err != nil {
		return nil, errors.ToAutoscalerError(errors.InternalError, err)
	}
	return newNodeInfo, nil
}
func sanitizeNodeInfo(nodeInfo *schedulercache.NodeInfo, nodeGroupName string) (*schedulercache.NodeInfo, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	sanitizedNode, err := sanitizeTemplateNode(nodeInfo.Node(), nodeGroupName)
	if err != nil {
		return nil, err
	}
	sanitizedPods := make([]*apiv1.Pod, 0)
	for _, pod := range nodeInfo.Pods() {
		sanitizedPod := pod.DeepCopy()
		sanitizedPod.Spec.NodeName = sanitizedNode.Name
		sanitizedPods = append(sanitizedPods, sanitizedPod)
	}
	sanitizedNodeInfo := schedulercache.NewNodeInfo(sanitizedPods...)
	if err := sanitizedNodeInfo.SetNode(sanitizedNode); err != nil {
		return nil, errors.ToAutoscalerError(errors.InternalError, err)
	}
	return sanitizedNodeInfo, nil
}
func sanitizeTemplateNode(node *apiv1.Node, nodeGroup string) (*apiv1.Node, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	newNode := node.DeepCopy()
	nodeName := fmt.Sprintf("template-node-for-%s-%d", nodeGroup, rand.Int63())
	newNode.Labels = make(map[string]string, len(node.Labels))
	for k, v := range node.Labels {
		if k != kubeletapis.LabelHostname {
			newNode.Labels[k] = v
		} else {
			newNode.Labels[k] = nodeName
		}
	}
	newNode.Name = nodeName
	newTaints := make([]apiv1.Taint, 0)
	for _, taint := range node.Spec.Taints {
		switch taint.Key {
		case ReschedulerTaintKey:
			klog.V(4).Infof("Removing rescheduler taint when creating template from node %s", node.Name)
		case deletetaint.ToBeDeletedTaint:
			klog.V(4).Infof("Removing autoscaler taint when creating template from node %s", node.Name)
		default:
			newTaints = append(newTaints, taint)
		}
	}
	newNode.Spec.Taints = newTaints
	return newNode, nil
}
func removeOldUnregisteredNodes(unregisteredNodes []clusterstate.UnregisteredNode, context *context.AutoscalingContext, currentTime time.Time, logRecorder *utils.LogEventRecorder) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	removedAny := false
	for _, unregisteredNode := range unregisteredNodes {
		if unregisteredNode.UnregisteredSince.Add(context.MaxNodeProvisionTime).Before(currentTime) {
			klog.V(0).Infof("Removing unregistered node %v", unregisteredNode.Node.Name)
			nodeGroup, err := context.CloudProvider.NodeGroupForNode(unregisteredNode.Node)
			if err != nil {
				klog.Warningf("Failed to get node group for %s: %v", unregisteredNode.Node.Name, err)
				return removedAny, err
			}
			if nodeGroup == nil || reflect.ValueOf(nodeGroup).IsNil() {
				klog.Warningf("No node group for node %s, skipping", unregisteredNode.Node.Name)
				continue
			}
			size, err := nodeGroup.TargetSize()
			if err != nil {
				klog.Warningf("Failed to get node group size, err: %v", err)
				continue
			}
			if nodeGroup.MinSize() >= size {
				klog.Warningf("Failed to remove node %s: node group min size reached, skipping unregistered node removal", unregisteredNode.Node.Name)
				continue
			}
			err = nodeGroup.DeleteNodes([]*apiv1.Node{unregisteredNode.Node})
			if err != nil {
				klog.Warningf("Failed to remove node %s: %v", unregisteredNode.Node.Name, err)
				logRecorder.Eventf(apiv1.EventTypeWarning, "DeleteUnregisteredFailed", "Failed to remove node %s: %v", unregisteredNode.Node.Name, err)
				return removedAny, err
			}
			logRecorder.Eventf(apiv1.EventTypeNormal, "DeleteUnregistered", "Removed unregistered node %v", unregisteredNode.Node.Name)
			removedAny = true
		}
	}
	return removedAny, nil
}
func fixNodeGroupSize(context *context.AutoscalingContext, clusterStateRegistry *clusterstate.ClusterStateRegistry, currentTime time.Time) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	fixed := false
	for _, nodeGroup := range context.CloudProvider.NodeGroups() {
		incorrectSize := clusterStateRegistry.GetIncorrectNodeGroupSize(nodeGroup.Id())
		if incorrectSize == nil {
			continue
		}
		if incorrectSize.FirstObserved.Add(context.MaxNodeProvisionTime).Before(currentTime) {
			delta := incorrectSize.CurrentSize - incorrectSize.ExpectedSize
			if delta < 0 {
				klog.V(0).Infof("Decreasing size of %s, expected=%d current=%d delta=%d", nodeGroup.Id(), incorrectSize.ExpectedSize, incorrectSize.CurrentSize, delta)
				if err := nodeGroup.DecreaseTargetSize(delta); err != nil {
					return fixed, fmt.Errorf("Failed to decrease %s: %v", nodeGroup.Id(), err)
				}
				fixed = true
			}
		}
	}
	return fixed, nil
}
func getPotentiallyUnneededNodes(context *context.AutoscalingContext, nodes []*apiv1.Node) []*apiv1.Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make([]*apiv1.Node, 0, len(nodes))
	nodeGroupSize := getNodeGroupSizeMap(context.CloudProvider)
	for _, node := range nodes {
		nodeGroup, err := context.CloudProvider.NodeGroupForNode(node)
		if err != nil {
			klog.Warningf("Error while checking node group for %s: %v", node.Name, err)
			continue
		}
		if nodeGroup == nil || reflect.ValueOf(nodeGroup).IsNil() {
			klog.V(4).Infof("Skipping %s - no node group config", node.Name)
			continue
		}
		size, found := nodeGroupSize[nodeGroup.Id()]
		if !found {
			klog.Errorf("Error while checking node group size %s: group size not found", nodeGroup.Id())
			continue
		}
		if size <= nodeGroup.MinSize() {
			klog.V(1).Infof("Skipping %s - node group min size reached", node.Name)
			continue
		}
		result = append(result, node)
	}
	return result
}
func hasHardInterPodAffinity(affinity *apiv1.Affinity) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if affinity == nil {
		return false
	}
	if affinity.PodAffinity != nil {
		if len(affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution) > 0 {
			return true
		}
	}
	if affinity.PodAntiAffinity != nil {
		if len(affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution) > 0 {
			return true
		}
	}
	return false
}
func anyPodHasHardInterPodAffinity(pods []*apiv1.Pod) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, pod := range pods {
		if hasHardInterPodAffinity(pod.Spec.Affinity) {
			return true
		}
	}
	return false
}
func ConfigurePredicateCheckerForLoop(unschedulablePods []*apiv1.Pod, schedulablePods []*apiv1.Pod, predicateChecker *simulator.PredicateChecker) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	podsWithAffinityFound := anyPodHasHardInterPodAffinity(unschedulablePods)
	if !podsWithAffinityFound {
		podsWithAffinityFound = anyPodHasHardInterPodAffinity(schedulablePods)
	}
	predicateChecker.SetAffinityPredicateEnabled(podsWithAffinityFound)
	if !podsWithAffinityFound {
		klog.V(1).Info("No pod using affinity / antiaffinity found in cluster, disabling affinity predicate for this loop")
	}
}
func getNodeCoresAndMemory(node *apiv1.Node) (int64, int64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cores := getNodeResource(node, apiv1.ResourceCPU)
	memory := getNodeResource(node, apiv1.ResourceMemory)
	return cores, memory
}
func getNodeResource(node *apiv1.Node, resource apiv1.ResourceName) int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	nodeCapacity, found := node.Status.Capacity[resource]
	if !found {
		return 0
	}
	nodeCapacityValue := nodeCapacity.Value()
	if nodeCapacityValue < 0 {
		nodeCapacityValue = 0
	}
	return nodeCapacityValue
}
func getNodeGroupSizeMap(cloudProvider cloudprovider.CloudProvider) map[string]int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	nodeGroupSize := make(map[string]int)
	for _, nodeGroup := range cloudProvider.NodeGroups() {
		size, err := nodeGroup.TargetSize()
		if err != nil {
			klog.Errorf("Error while checking node group size %s: %v", nodeGroup.Id(), err)
			continue
		}
		nodeGroupSize[nodeGroup.Id()] = size
	}
	return nodeGroupSize
}
func UpdateClusterStateMetrics(csr *clusterstate.ClusterStateRegistry) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if csr == nil || reflect.ValueOf(csr).IsNil() {
		return
	}
	metrics.UpdateClusterSafeToAutoscale(csr.IsClusterHealthy())
	readiness := csr.GetClusterReadiness()
	metrics.UpdateNodesCount(readiness.Ready, readiness.Unready+readiness.LongNotStarted, readiness.NotStarted, readiness.LongUnregistered, readiness.Unregistered)
}
func getOldestCreateTime(pods []*apiv1.Pod) time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	oldest := time.Now()
	for _, pod := range pods {
		if oldest.After(pod.CreationTimestamp.Time) {
			oldest = pod.CreationTimestamp.Time
		}
	}
	return oldest
}
func getOldestCreateTimeWithGpu(pods []*apiv1.Pod) (bool, time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	oldest := time.Now()
	gpuFound := false
	for _, pod := range pods {
		if gpu.PodRequestsGpu(pod) {
			gpuFound = true
			if oldest.After(pod.CreationTimestamp.Time) {
				oldest = pod.CreationTimestamp.Time
			}
		}
	}
	return gpuFound, oldest
}
func UpdateEmptyClusterStateMetrics() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	metrics.UpdateClusterSafeToAutoscale(false)
	metrics.UpdateNodesCount(0, 0, 0, 0, 0)
}
