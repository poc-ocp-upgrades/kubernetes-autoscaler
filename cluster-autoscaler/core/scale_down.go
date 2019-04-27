package core

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"sync"
	"time"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/clusterstate"
	"k8s.io/autoscaler/cluster-autoscaler/context"
	"k8s.io/autoscaler/cluster-autoscaler/metrics"
	"k8s.io/autoscaler/cluster-autoscaler/simulator"
	"k8s.io/autoscaler/cluster-autoscaler/utils/deletetaint"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	kube_util "k8s.io/autoscaler/cluster-autoscaler/utils/kubernetes"
	scheduler_util "k8s.io/autoscaler/cluster-autoscaler/utils/scheduler"
	apiv1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1beta1"
	kube_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kube_client "k8s.io/client-go/kubernetes"
	kube_record "k8s.io/client-go/tools/record"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/autoscaler/cluster-autoscaler/processors/status"
	"k8s.io/autoscaler/cluster-autoscaler/utils/gpu"
	"k8s.io/klog"
)

const (
	ScaleDownDisabledKey = "cluster-autoscaler.kubernetes.io/scale-down-disabled"
)
const (
	MaxKubernetesEmptyNodeDeletionTime	= 3 * time.Minute
	MaxCloudProviderNodeDeletionTime	= 5 * time.Minute
	MaxPodEvictionTime			= 2 * time.Minute
	EvictionRetryTime			= 10 * time.Second
	PodEvictionHeadroom			= 30 * time.Second
)

type NodeDeleteStatus struct {
	sync.Mutex
	deleteInProgress	bool
	nodeDeleteResults	map[string]error
}

func (n *NodeDeleteStatus) IsDeleteInProgress() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n.Lock()
	defer n.Unlock()
	return n.deleteInProgress
}
func (n *NodeDeleteStatus) SetDeleteInProgress(status bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n.Lock()
	defer n.Unlock()
	n.deleteInProgress = status
}
func (n *NodeDeleteStatus) AddNodeDeleteResult(nodeName string, result error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n.Lock()
	defer n.Unlock()
	n.nodeDeleteResults[nodeName] = result
}
func (n *NodeDeleteStatus) DrainNodeDeleteResults() map[string]error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n.Lock()
	defer n.Unlock()
	results := n.nodeDeleteResults
	n.nodeDeleteResults = make(map[string]error)
	return results
}

type scaleDownResourcesLimits map[string]int64
type scaleDownResourcesDelta map[string]int64

const scaleDownLimitUnknown = math.MinInt64

func computeScaleDownResourcesLeftLimits(nodes []*apiv1.Node, resourceLimiter *cloudprovider.ResourceLimiter, cp cloudprovider.CloudProvider, timestamp time.Time) scaleDownResourcesLimits {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	totalCores, totalMem := calculateScaleDownCoresMemoryTotal(nodes, timestamp)
	var totalGpus map[string]int64
	var totalGpusErr error
	if cloudprovider.ContainsGpuResources(resourceLimiter.GetResources()) {
		totalGpus, totalGpusErr = calculateScaleDownGpusTotal(nodes, cp, timestamp)
	}
	resultScaleDownLimits := make(scaleDownResourcesLimits)
	for _, resource := range resourceLimiter.GetResources() {
		min := resourceLimiter.GetMin(resource)
		if min > 0 {
			switch {
			case resource == cloudprovider.ResourceNameCores:
				resultScaleDownLimits[resource] = computeAboveMin(totalCores, min)
			case resource == cloudprovider.ResourceNameMemory:
				resultScaleDownLimits[resource] = computeAboveMin(totalMem, min)
			case cloudprovider.IsGpuResource(resource):
				if totalGpusErr != nil {
					resultScaleDownLimits[resource] = scaleDownLimitUnknown
				} else {
					resultScaleDownLimits[resource] = computeAboveMin(totalGpus[resource], min)
				}
			default:
				klog.Errorf("Scale down limits defined for unsupported resource '%s'", resource)
			}
		}
	}
	return resultScaleDownLimits
}
func computeAboveMin(total int64, min int64) int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if total > min {
		return total - min
	}
	return 0
}
func calculateScaleDownCoresMemoryTotal(nodes []*apiv1.Node, timestamp time.Time) (int64, int64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var coresTotal, memoryTotal int64
	for _, node := range nodes {
		if isNodeBeingDeleted(node, timestamp) {
			continue
		}
		cores, memory := getNodeCoresAndMemory(node)
		coresTotal += cores
		memoryTotal += memory
	}
	return coresTotal, memoryTotal
}
func calculateScaleDownGpusTotal(nodes []*apiv1.Node, cp cloudprovider.CloudProvider, timestamp time.Time) (map[string]int64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	type gpuInfo struct {
		name	string
		count	int64
	}
	result := make(map[string]int64)
	ngCache := make(map[string]gpuInfo)
	for _, node := range nodes {
		if isNodeBeingDeleted(node, timestamp) {
			continue
		}
		nodeGroup, err := cp.NodeGroupForNode(node)
		if err != nil {
			return nil, errors.ToAutoscalerError(errors.CloudProviderError, err).AddPrefix("can not get node group for node %v when calculating cluster gpu usage", node.Name)
		}
		if nodeGroup == nil || reflect.ValueOf(nodeGroup).IsNil() {
			nodeGroup = nil
		}
		var gpuType string
		var gpuCount int64
		var cached gpuInfo
		var cacheHit bool
		if nodeGroup != nil {
			cached, cacheHit = ngCache[nodeGroup.Id()]
			if cacheHit {
				gpuType = cached.name
				gpuCount = cached.count
			}
		}
		if !cacheHit {
			gpuType, gpuCount, err = gpu.GetNodeTargetGpus(node, nodeGroup)
			if err != nil {
				return nil, errors.ToAutoscalerError(errors.CloudProviderError, err).AddPrefix("can not get gpu count for node %v when calculating cluster gpu usage")
			}
			if nodeGroup != nil {
				ngCache[nodeGroup.Id()] = gpuInfo{name: gpuType, count: gpuCount}
			}
		}
		if gpuType == "" || gpuCount == 0 {
			continue
		}
		result[gpuType] += gpuCount
	}
	return result, nil
}
func isNodeBeingDeleted(node *apiv1.Node, timestamp time.Time) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	deleteTime, _ := deletetaint.GetToBeDeletedTime(node)
	return deleteTime != nil && (timestamp.Sub(*deleteTime) < MaxCloudProviderNodeDeletionTime || timestamp.Sub(*deleteTime) < MaxKubernetesEmptyNodeDeletionTime)
}
func noScaleDownLimitsOnResources() scaleDownResourcesLimits {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func copyScaleDownResourcesLimits(source scaleDownResourcesLimits) scaleDownResourcesLimits {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	copy := scaleDownResourcesLimits{}
	for k, v := range source {
		copy[k] = v
	}
	return copy
}
func computeScaleDownResourcesDelta(node *apiv1.Node, nodeGroup cloudprovider.NodeGroup, resourcesWithLimits []string) (scaleDownResourcesDelta, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	resultScaleDownDelta := make(scaleDownResourcesDelta)
	nodeCPU, nodeMemory := getNodeCoresAndMemory(node)
	resultScaleDownDelta[cloudprovider.ResourceNameCores] = nodeCPU
	resultScaleDownDelta[cloudprovider.ResourceNameMemory] = nodeMemory
	if cloudprovider.ContainsGpuResources(resourcesWithLimits) {
		gpuType, gpuCount, err := gpu.GetNodeTargetGpus(node, nodeGroup)
		if err != nil {
			return scaleDownResourcesDelta{}, errors.ToAutoscalerError(errors.CloudProviderError, err).AddPrefix("Failed to get node %v gpu: %v", node.Name)
		}
		resultScaleDownDelta[gpuType] = gpuCount
	}
	return resultScaleDownDelta, nil
}

type scaleDownLimitsCheckResult struct {
	exceeded		bool
	exceededResources	[]string
}

func scaleDownLimitsNotExceeded() scaleDownLimitsCheckResult {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return scaleDownLimitsCheckResult{false, []string{}}
}
func (limits *scaleDownResourcesLimits) checkScaleDownDeltaWithinLimits(delta scaleDownResourcesDelta) scaleDownLimitsCheckResult {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	exceededResources := sets.NewString()
	for resource, resourceDelta := range delta {
		resourceLeft, found := (*limits)[resource]
		if found {
			if (resourceDelta > 0) && (resourceLeft == scaleDownLimitUnknown || resourceDelta > resourceLeft) {
				exceededResources.Insert(resource)
			}
		}
	}
	if len(exceededResources) > 0 {
		return scaleDownLimitsCheckResult{true, exceededResources.List()}
	}
	return scaleDownLimitsNotExceeded()
}
func (limits *scaleDownResourcesLimits) tryDecrementLimitsByDelta(delta scaleDownResourcesDelta) scaleDownLimitsCheckResult {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := limits.checkScaleDownDeltaWithinLimits(delta)
	if result.exceeded {
		return result
	}
	for resource, resourceDelta := range delta {
		resourceLeft, found := (*limits)[resource]
		if found {
			(*limits)[resource] = resourceLeft - resourceDelta
		}
	}
	return scaleDownLimitsNotExceeded()
}

type ScaleDown struct {
	context			*context.AutoscalingContext
	clusterStateRegistry	*clusterstate.ClusterStateRegistry
	unneededNodes		map[string]time.Time
	unneededNodesList	[]*apiv1.Node
	unremovableNodes	map[string]time.Time
	podLocationHints	map[string]string
	nodeUtilizationMap	map[string]simulator.UtilizationInfo
	usageTracker		*simulator.UsageTracker
	nodeDeleteStatus	*NodeDeleteStatus
}

func NewScaleDown(context *context.AutoscalingContext, clusterStateRegistry *clusterstate.ClusterStateRegistry) *ScaleDown {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &ScaleDown{context: context, clusterStateRegistry: clusterStateRegistry, unneededNodes: make(map[string]time.Time), unremovableNodes: make(map[string]time.Time), podLocationHints: make(map[string]string), nodeUtilizationMap: make(map[string]simulator.UtilizationInfo), usageTracker: simulator.NewUsageTracker(), unneededNodesList: make([]*apiv1.Node, 0), nodeDeleteStatus: &NodeDeleteStatus{nodeDeleteResults: make(map[string]error)}}
}
func (sd *ScaleDown) CleanUp(timestamp time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	sd.usageTracker.CleanUp(timestamp.Add(-sd.context.ScaleDownUnneededTime))
}
func (sd *ScaleDown) GetCandidatesForScaleDown() []*apiv1.Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return sd.unneededNodesList
}
func (sd *ScaleDown) CleanUpUnneededNodes() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	sd.unneededNodesList = make([]*apiv1.Node, 0)
	sd.unneededNodes = make(map[string]time.Time)
}
func (sd *ScaleDown) UpdateUnneededNodes(nodes []*apiv1.Node, nodesToCheck []*apiv1.Node, pods []*apiv1.Pod, timestamp time.Time, pdbs []*policyv1.PodDisruptionBudget) errors.AutoscalerError {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	currentlyUnneededNodes := make([]*apiv1.Node, 0)
	nonExpendablePods := FilterOutExpendablePods(pods, sd.context.ExpendablePodsPriorityCutoff)
	nodeNameToNodeInfo := scheduler_util.CreateNodeNameToInfoMap(nonExpendablePods, nodes)
	utilizationMap := make(map[string]simulator.UtilizationInfo)
	sd.updateUnremovableNodes(nodes)
	filteredNodesToCheck := make([]*apiv1.Node, 0)
	for _, node := range nodesToCheck {
		if unremovableTimestamp, found := sd.unremovableNodes[node.Name]; found {
			if unremovableTimestamp.After(timestamp) {
				continue
			}
			delete(sd.unremovableNodes, node.Name)
		}
		filteredNodesToCheck = append(filteredNodesToCheck, node)
	}
	skipped := len(nodesToCheck) - len(filteredNodesToCheck)
	if skipped > 0 {
		klog.V(1).Infof("Scale-down calculation: ignoring %v nodes unremovable in the last %v", skipped, sd.context.AutoscalingOptions.UnremovableNodeRecheckTimeout)
	}
	for _, node := range filteredNodesToCheck {
		if isNodeBeingDeleted(node, timestamp) {
			klog.V(1).Infof("Skipping %s from delete considerations - the node is currently being deleted", node.Name)
			continue
		}
		if hasNoScaleDownAnnotation(node) {
			klog.V(1).Infof("Skipping %s from delete consideration - the node is marked as no scale down", node.Name)
			continue
		}
		nodeInfo, found := nodeNameToNodeInfo[node.Name]
		if !found {
			klog.Errorf("Node info for %s not found", node.Name)
			continue
		}
		utilInfo, err := simulator.CalculateUtilization(node, nodeInfo, sd.context.IgnoreDaemonSetsUtilization, sd.context.IgnoreMirrorPodsUtilization)
		if err != nil {
			klog.Warningf("Failed to calculate utilization for %s: %v", node.Name, err)
		}
		klog.V(4).Infof("Node %s - utilization %f", node.Name, utilInfo.Utilization)
		utilizationMap[node.Name] = utilInfo
		if utilInfo.Utilization >= sd.context.ScaleDownUtilizationThreshold {
			klog.V(4).Infof("Node %s is not suitable for removal - utilization too big (%f)", node.Name, utilInfo.Utilization)
			continue
		}
		currentlyUnneededNodes = append(currentlyUnneededNodes, node)
	}
	emptyNodes := make(map[string]bool)
	emptyNodesList := getEmptyNodesNoResourceLimits(currentlyUnneededNodes, pods, len(currentlyUnneededNodes), sd.context.CloudProvider)
	for _, node := range emptyNodesList {
		emptyNodes[node.Name] = true
	}
	currentlyUnneededNonEmptyNodes := make([]*apiv1.Node, 0, len(currentlyUnneededNodes))
	for _, node := range currentlyUnneededNodes {
		if !emptyNodes[node.Name] {
			currentlyUnneededNonEmptyNodes = append(currentlyUnneededNonEmptyNodes, node)
		}
	}
	currentCandidates, currentNonCandidates := sd.chooseCandidates(currentlyUnneededNonEmptyNodes)
	nodesToRemove, unremovable, newHints, simulatorErr := simulator.FindNodesToRemove(currentCandidates, nodes, nonExpendablePods, nil, sd.context.PredicateChecker, len(currentCandidates), true, sd.podLocationHints, sd.usageTracker, timestamp, pdbs)
	if simulatorErr != nil {
		return sd.markSimulationError(simulatorErr, timestamp)
	}
	additionalCandidatesCount := sd.context.ScaleDownNonEmptyCandidatesCount - len(nodesToRemove)
	if additionalCandidatesCount > len(currentNonCandidates) {
		additionalCandidatesCount = len(currentNonCandidates)
	}
	additionalCandidatesPoolSize := int(math.Ceil(float64(len(nodes)) * sd.context.ScaleDownCandidatesPoolRatio))
	if additionalCandidatesPoolSize < sd.context.ScaleDownCandidatesPoolMinCount {
		additionalCandidatesPoolSize = sd.context.ScaleDownCandidatesPoolMinCount
	}
	if additionalCandidatesPoolSize > len(currentNonCandidates) {
		additionalCandidatesPoolSize = len(currentNonCandidates)
	}
	if additionalCandidatesCount > 0 {
		klog.V(3).Infof("Finding additional %v candidates for scale down.", additionalCandidatesCount)
		additionalNodesToRemove, additionalUnremovable, additionalNewHints, simulatorErr := simulator.FindNodesToRemove(currentNonCandidates[:additionalCandidatesPoolSize], nodes, nonExpendablePods, nil, sd.context.PredicateChecker, additionalCandidatesCount, true, sd.podLocationHints, sd.usageTracker, timestamp, pdbs)
		if simulatorErr != nil {
			return sd.markSimulationError(simulatorErr, timestamp)
		}
		nodesToRemove = append(nodesToRemove, additionalNodesToRemove...)
		unremovable = append(unremovable, additionalUnremovable...)
		for key, value := range additionalNewHints {
			newHints[key] = value
		}
	}
	for _, node := range emptyNodesList {
		nodesToRemove = append(nodesToRemove, simulator.NodeToBeRemoved{Node: node, PodsToReschedule: []*apiv1.Pod{}})
	}
	result := make(map[string]time.Time)
	unneededNodesList := make([]*apiv1.Node, 0, len(nodesToRemove))
	for _, node := range nodesToRemove {
		name := node.Node.Name
		unneededNodesList = append(unneededNodesList, node.Node)
		if val, found := sd.unneededNodes[name]; !found {
			result[name] = timestamp
		} else {
			result[name] = val
		}
	}
	if len(unremovable) > 0 {
		unremovableTimeout := timestamp.Add(sd.context.AutoscalingOptions.UnremovableNodeRecheckTimeout)
		for _, node := range unremovable {
			sd.unremovableNodes[node.Name] = unremovableTimeout
		}
		klog.V(1).Infof("%v nodes found to be unremovable in simulation, will re-check them at %v", len(unremovable), unremovableTimeout)
	}
	sd.unneededNodesList = unneededNodesList
	sd.unneededNodes = result
	sd.podLocationHints = newHints
	sd.nodeUtilizationMap = utilizationMap
	sd.clusterStateRegistry.UpdateScaleDownCandidates(sd.unneededNodesList, timestamp)
	metrics.UpdateUnneededNodesCount(len(sd.unneededNodesList))
	return nil
}
func (sd *ScaleDown) updateUnremovableNodes(nodes []*apiv1.Node) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(sd.unremovableNodes) <= 0 {
		return
	}
	nodesToDelete := make(map[string]struct{}, len(sd.unremovableNodes))
	for name := range sd.unremovableNodes {
		nodesToDelete[name] = struct{}{}
	}
	for _, node := range nodes {
		if _, ok := nodesToDelete[node.Name]; ok {
			delete(nodesToDelete, node.Name)
		}
	}
	for nodeName := range nodesToDelete {
		delete(sd.unremovableNodes, nodeName)
	}
}
func (sd *ScaleDown) markSimulationError(simulatorErr errors.AutoscalerError, timestamp time.Time) errors.AutoscalerError {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	klog.Errorf("Error while simulating node drains: %v", simulatorErr)
	sd.unneededNodesList = make([]*apiv1.Node, 0)
	sd.unneededNodes = make(map[string]time.Time)
	sd.nodeUtilizationMap = make(map[string]simulator.UtilizationInfo)
	sd.clusterStateRegistry.UpdateScaleDownCandidates(sd.unneededNodesList, timestamp)
	return simulatorErr.AddPrefix("error while simulating node drains: ")
}
func (sd *ScaleDown) chooseCandidates(nodes []*apiv1.Node) ([]*apiv1.Node, []*apiv1.Node) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if sd.context.ScaleDownNonEmptyCandidatesCount <= 0 {
		return nodes, []*apiv1.Node{}
	}
	currentCandidates := make([]*apiv1.Node, 0, len(sd.unneededNodesList))
	currentNonCandidates := make([]*apiv1.Node, 0, len(nodes))
	for _, node := range nodes {
		if _, found := sd.unneededNodes[node.Name]; found {
			currentCandidates = append(currentCandidates, node)
		} else {
			currentNonCandidates = append(currentNonCandidates, node)
		}
	}
	return currentCandidates, currentNonCandidates
}
func (sd *ScaleDown) mapNodesToStatusScaleDownNodes(nodes []*apiv1.Node, nodeGroups map[string]cloudprovider.NodeGroup, evictedPodLists map[string][]*apiv1.Pod) []*status.ScaleDownNode {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var result []*status.ScaleDownNode
	for _, node := range nodes {
		result = append(result, &status.ScaleDownNode{Node: node, NodeGroup: nodeGroups[node.Name], UtilInfo: sd.nodeUtilizationMap[node.Name], EvictedPods: evictedPodLists[node.Name]})
	}
	return result
}
func (sd *ScaleDown) TryToScaleDown(allNodes []*apiv1.Node, pods []*apiv1.Pod, pdbs []*policyv1.PodDisruptionBudget, currentTime time.Time) (*status.ScaleDownStatus, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	scaleDownStatus := &status.ScaleDownStatus{NodeDeleteResults: sd.nodeDeleteStatus.DrainNodeDeleteResults()}
	nodeDeletionDuration := time.Duration(0)
	findNodesToRemoveDuration := time.Duration(0)
	defer updateScaleDownMetrics(time.Now(), &findNodesToRemoveDuration, &nodeDeletionDuration)
	nodesWithoutMaster := filterOutMasters(allNodes, pods)
	candidates := make([]*apiv1.Node, 0)
	readinessMap := make(map[string]bool)
	candidateNodeGroups := make(map[string]cloudprovider.NodeGroup)
	resourceLimiter, errCP := sd.context.CloudProvider.GetResourceLimiter()
	if errCP != nil {
		scaleDownStatus.Result = status.ScaleDownError
		return scaleDownStatus, errors.ToAutoscalerError(errors.CloudProviderError, errCP)
	}
	scaleDownResourcesLeft := computeScaleDownResourcesLeftLimits(nodesWithoutMaster, resourceLimiter, sd.context.CloudProvider, currentTime)
	nodeGroupSize := getNodeGroupSizeMap(sd.context.CloudProvider)
	resourcesWithLimits := resourceLimiter.GetResources()
	for _, node := range nodesWithoutMaster {
		if val, found := sd.unneededNodes[node.Name]; found {
			klog.V(2).Infof("%s was unneeded for %s", node.Name, currentTime.Sub(val).String())
			if hasNoScaleDownAnnotation(node) {
				klog.V(4).Infof("Skipping %s - scale down disabled annotation found", node.Name)
				continue
			}
			ready, _, _ := kube_util.GetReadinessState(node)
			readinessMap[node.Name] = ready
			if ready && !val.Add(sd.context.ScaleDownUnneededTime).Before(currentTime) {
				continue
			}
			if !ready && !val.Add(sd.context.ScaleDownUnreadyTime).Before(currentTime) {
				continue
			}
			nodeGroup, err := sd.context.CloudProvider.NodeGroupForNode(node)
			if err != nil {
				klog.Errorf("Error while checking node group for %s: %v", node.Name, err)
				continue
			}
			if nodeGroup == nil || reflect.ValueOf(nodeGroup).IsNil() {
				klog.V(4).Infof("Skipping %s - no node group config", node.Name)
				continue
			}
			size, found := nodeGroupSize[nodeGroup.Id()]
			if !found {
				klog.Errorf("Error while checking node group size %s: group size not found in cache", nodeGroup.Id())
				continue
			}
			if size <= nodeGroup.MinSize() {
				klog.V(1).Infof("Skipping %s - node group min size reached", node.Name)
				continue
			}
			scaleDownResourcesDelta, err := computeScaleDownResourcesDelta(node, nodeGroup, resourcesWithLimits)
			if err != nil {
				klog.Errorf("Error getting node resources: %v", err)
				continue
			}
			checkResult := scaleDownResourcesLeft.checkScaleDownDeltaWithinLimits(scaleDownResourcesDelta)
			if checkResult.exceeded {
				klog.V(4).Infof("Skipping %s - minimal limit exceeded for %v", node.Name, checkResult.exceededResources)
				continue
			}
			candidates = append(candidates, node)
			candidateNodeGroups[node.Name] = nodeGroup
		}
	}
	if len(candidates) == 0 {
		klog.V(1).Infof("No candidates for scale down")
		scaleDownStatus.Result = status.ScaleDownNoUnneeded
		return scaleDownStatus, nil
	}
	emptyNodes := getEmptyNodes(candidates, pods, sd.context.MaxEmptyBulkDelete, scaleDownResourcesLeft, sd.context.CloudProvider)
	if len(emptyNodes) > 0 {
		nodeDeletionStart := time.Now()
		confirmation := make(chan errors.AutoscalerError, len(emptyNodes))
		sd.scheduleDeleteEmptyNodes(emptyNodes, sd.context.ClientSet, sd.context.Recorder, readinessMap, candidateNodeGroups, confirmation)
		err := sd.waitForEmptyNodesDeleted(emptyNodes, confirmation)
		nodeDeletionDuration = time.Now().Sub(nodeDeletionStart)
		if err == nil {
			scaleDownStatus.ScaledDownNodes = sd.mapNodesToStatusScaleDownNodes(emptyNodes, candidateNodeGroups, make(map[string][]*apiv1.Pod))
			scaleDownStatus.Result = status.ScaleDownNodeDeleted
			return scaleDownStatus, nil
		}
		scaleDownStatus.Result = status.ScaleDownError
		return scaleDownStatus, err.AddPrefix("failed to delete at least one empty node: ")
	}
	findNodesToRemoveStart := time.Now()
	nonExpendablePods := FilterOutExpendablePods(pods, sd.context.ExpendablePodsPriorityCutoff)
	nodesToRemove, _, _, err := simulator.FindNodesToRemove(candidates, nodesWithoutMaster, nonExpendablePods, sd.context.ClientSet, sd.context.PredicateChecker, 1, false, sd.podLocationHints, sd.usageTracker, time.Now(), pdbs)
	findNodesToRemoveDuration = time.Now().Sub(findNodesToRemoveStart)
	if err != nil {
		scaleDownStatus.Result = status.ScaleDownError
		return scaleDownStatus, err.AddPrefix("Find node to remove failed: ")
	}
	if len(nodesToRemove) == 0 {
		klog.V(1).Infof("No node to remove")
		scaleDownStatus.Result = status.ScaleDownNoNodeDeleted
		return scaleDownStatus, nil
	}
	toRemove := nodesToRemove[0]
	utilization := sd.nodeUtilizationMap[toRemove.Node.Name]
	podNames := make([]string, 0, len(toRemove.PodsToReschedule))
	for _, pod := range toRemove.PodsToReschedule {
		podNames = append(podNames, pod.Namespace+"/"+pod.Name)
	}
	klog.V(0).Infof("Scale-down: removing node %s, utilization: %v, pods to reschedule: %s", toRemove.Node.Name, utilization, strings.Join(podNames, ","))
	sd.context.LogRecorder.Eventf(apiv1.EventTypeNormal, "ScaleDown", "Scale-down: removing node %s, utilization: %v, pods to reschedule: %s", toRemove.Node.Name, utilization, strings.Join(podNames, ","))
	simulator.RemoveNodeFromTracker(sd.usageTracker, toRemove.Node.Name, sd.unneededNodes)
	nodeDeletionStart := time.Now()
	nodeDeletionDuration = time.Now().Sub(nodeDeletionStart)
	sd.nodeDeleteStatus.SetDeleteInProgress(true)
	go func() {
		var err error
		defer func() {
			sd.nodeDeleteStatus.AddNodeDeleteResult(toRemove.Node.Name, err)
		}()
		defer sd.nodeDeleteStatus.SetDeleteInProgress(false)
		err = sd.deleteNode(toRemove.Node, toRemove.PodsToReschedule)
		if err != nil {
			klog.Errorf("Failed to delete %s: %v", toRemove.Node.Name, err)
			return
		}
		nodeGroup := candidateNodeGroups[toRemove.Node.Name]
		if readinessMap[toRemove.Node.Name] {
			metrics.RegisterScaleDown(1, gpu.GetGpuTypeForMetrics(toRemove.Node, nodeGroup), metrics.Underutilized)
		} else {
			metrics.RegisterScaleDown(1, gpu.GetGpuTypeForMetrics(toRemove.Node, nodeGroup), metrics.Unready)
		}
	}()
	scaleDownStatus.ScaledDownNodes = sd.mapNodesToStatusScaleDownNodes([]*apiv1.Node{toRemove.Node}, candidateNodeGroups, map[string][]*apiv1.Pod{toRemove.Node.Name: toRemove.PodsToReschedule})
	scaleDownStatus.Result = status.ScaleDownNodeDeleteStarted
	return scaleDownStatus, nil
}
func updateScaleDownMetrics(scaleDownStart time.Time, findNodesToRemoveDuration *time.Duration, nodeDeletionDuration *time.Duration) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	stop := time.Now()
	miscDuration := stop.Sub(scaleDownStart) - *nodeDeletionDuration - *findNodesToRemoveDuration
	metrics.UpdateDuration(metrics.ScaleDownNodeDeletion, *nodeDeletionDuration)
	metrics.UpdateDuration(metrics.ScaleDownFindNodesToRemove, *findNodesToRemoveDuration)
	metrics.UpdateDuration(metrics.ScaleDownMiscOperations, miscDuration)
}
func getEmptyNodesNoResourceLimits(candidates []*apiv1.Node, pods []*apiv1.Pod, maxEmptyBulkDelete int, cloudProvider cloudprovider.CloudProvider) []*apiv1.Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return getEmptyNodes(candidates, pods, maxEmptyBulkDelete, noScaleDownLimitsOnResources(), cloudProvider)
}
func getEmptyNodes(candidates []*apiv1.Node, pods []*apiv1.Pod, maxEmptyBulkDelete int, resourcesLimits scaleDownResourcesLimits, cloudProvider cloudprovider.CloudProvider) []*apiv1.Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	emptyNodes := simulator.FindEmptyNodesToRemove(candidates, pods)
	availabilityMap := make(map[string]int)
	result := make([]*apiv1.Node, 0)
	resourcesLimitsCopy := copyScaleDownResourcesLimits(resourcesLimits)
	resourcesNames := sets.StringKeySet(resourcesLimits).List()
	for _, node := range emptyNodes {
		nodeGroup, err := cloudProvider.NodeGroupForNode(node)
		if err != nil {
			klog.Errorf("Failed to get group for %s", node.Name)
			continue
		}
		if nodeGroup == nil || reflect.ValueOf(nodeGroup).IsNil() {
			continue
		}
		var available int
		var found bool
		if available, found = availabilityMap[nodeGroup.Id()]; !found {
			size, err := nodeGroup.TargetSize()
			if err != nil {
				klog.Errorf("Failed to get size for %s: %v ", nodeGroup.Id(), err)
				continue
			}
			available = size - nodeGroup.MinSize()
			if available < 0 {
				available = 0
			}
			availabilityMap[nodeGroup.Id()] = available
		}
		if available > 0 {
			resourcesDelta, err := computeScaleDownResourcesDelta(node, nodeGroup, resourcesNames)
			if err != nil {
				klog.Errorf("Error: %v", err)
				continue
			}
			checkResult := resourcesLimitsCopy.tryDecrementLimitsByDelta(resourcesDelta)
			if checkResult.exceeded {
				continue
			}
			available -= 1
			availabilityMap[nodeGroup.Id()] = available
			result = append(result, node)
		}
	}
	limit := maxEmptyBulkDelete
	if len(result) < limit {
		limit = len(result)
	}
	return result[:limit]
}
func (sd *ScaleDown) scheduleDeleteEmptyNodes(emptyNodes []*apiv1.Node, client kube_client.Interface, recorder kube_record.EventRecorder, readinessMap map[string]bool, candidateNodeGroups map[string]cloudprovider.NodeGroup, confirmation chan errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, node := range emptyNodes {
		klog.V(0).Infof("Scale-down: removing empty node %s", node.Name)
		sd.context.LogRecorder.Eventf(apiv1.EventTypeNormal, "ScaleDownEmpty", "Scale-down: removing empty node %s", node.Name)
		simulator.RemoveNodeFromTracker(sd.usageTracker, node.Name, sd.unneededNodes)
		go func(nodeToDelete *apiv1.Node) {
			taintErr := deletetaint.MarkToBeDeleted(nodeToDelete, client)
			if taintErr != nil {
				recorder.Eventf(nodeToDelete, apiv1.EventTypeWarning, "ScaleDownFailed", "failed to mark the node as toBeDeleted/unschedulable: %v", taintErr)
				confirmation <- errors.ToAutoscalerError(errors.ApiCallError, taintErr)
				return
			}
			var deleteErr errors.AutoscalerError
			defer func() {
				if deleteErr != nil {
					deletetaint.CleanToBeDeleted(nodeToDelete, client)
					recorder.Eventf(nodeToDelete, apiv1.EventTypeWarning, "ScaleDownFailed", "failed to delete empty node: %v", deleteErr)
				} else {
					sd.context.LogRecorder.Eventf(apiv1.EventTypeNormal, "ScaleDownEmpty", "Scale-down: empty node %s removed", nodeToDelete.Name)
				}
			}()
			deleteErr = deleteNodeFromCloudProvider(nodeToDelete, sd.context.CloudProvider, sd.context.Recorder, sd.clusterStateRegistry)
			if deleteErr == nil {
				nodeGroup := candidateNodeGroups[nodeToDelete.Name]
				if readinessMap[nodeToDelete.Name] {
					metrics.RegisterScaleDown(1, gpu.GetGpuTypeForMetrics(nodeToDelete, nodeGroup), metrics.Empty)
				} else {
					metrics.RegisterScaleDown(1, gpu.GetGpuTypeForMetrics(nodeToDelete, nodeGroup), metrics.Unready)
				}
			}
			confirmation <- deleteErr
		}(node)
	}
}
func (sd *ScaleDown) waitForEmptyNodesDeleted(emptyNodes []*apiv1.Node, confirmation chan errors.AutoscalerError) errors.AutoscalerError {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var finalError errors.AutoscalerError
	startTime := time.Now()
	for range emptyNodes {
		timeElapsed := time.Now().Sub(startTime)
		timeLeft := MaxCloudProviderNodeDeletionTime - timeElapsed
		if timeLeft < 0 {
			return errors.NewAutoscalerError(errors.TransientError, "Failed to delete nodes in time")
		}
		select {
		case err := <-confirmation:
			if err != nil {
				klog.Errorf("Problem with empty node deletion: %v", err)
				finalError = err
			}
		case <-time.After(timeLeft):
			finalError = errors.NewAutoscalerError(errors.TransientError, "Failed to delete nodes in time")
		}
	}
	return finalError
}
func (sd *ScaleDown) deleteNode(node *apiv1.Node, pods []*apiv1.Pod) errors.AutoscalerError {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	deleteSuccessful := false
	drainSuccessful := false
	if err := deletetaint.MarkToBeDeleted(node, sd.context.ClientSet); err != nil {
		sd.context.Recorder.Eventf(node, apiv1.EventTypeWarning, "ScaleDownFailed", "failed to mark the node as toBeDeleted/unschedulable: %v", err)
		return errors.ToAutoscalerError(errors.ApiCallError, err)
	}
	defer func() {
		if !deleteSuccessful {
			deletetaint.CleanToBeDeleted(node, sd.context.ClientSet)
			if !drainSuccessful {
				sd.context.Recorder.Eventf(node, apiv1.EventTypeWarning, "ScaleDownFailed", "failed to drain the node, aborting ScaleDown")
			} else {
				sd.context.Recorder.Eventf(node, apiv1.EventTypeWarning, "ScaleDownFailed", "failed to delete the node")
			}
		}
	}()
	sd.context.Recorder.Eventf(node, apiv1.EventTypeNormal, "ScaleDown", "marked the node as toBeDeleted/unschedulable")
	if err := drainNode(node, pods, sd.context.ClientSet, sd.context.Recorder, sd.context.MaxGracefulTerminationSec, MaxPodEvictionTime, EvictionRetryTime); err != nil {
		return err
	}
	drainSuccessful = true
	err := deleteNodeFromCloudProvider(node, sd.context.CloudProvider, sd.context.Recorder, sd.clusterStateRegistry)
	if err != nil {
		return err
	}
	deleteSuccessful = true
	return nil
}
func evictPod(podToEvict *apiv1.Pod, client kube_client.Interface, recorder kube_record.EventRecorder, maxGracefulTerminationSec int, retryUntil time.Time, waitBetweenRetries time.Duration) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	recorder.Eventf(podToEvict, apiv1.EventTypeNormal, "ScaleDown", "deleting pod for node scale down")
	maxTermination := int64(apiv1.DefaultTerminationGracePeriodSeconds)
	if podToEvict.Spec.TerminationGracePeriodSeconds != nil {
		if *podToEvict.Spec.TerminationGracePeriodSeconds < int64(maxGracefulTerminationSec) {
			maxTermination = *podToEvict.Spec.TerminationGracePeriodSeconds
		} else {
			maxTermination = int64(maxGracefulTerminationSec)
		}
	}
	var lastError error
	for first := true; first || time.Now().Before(retryUntil); time.Sleep(waitBetweenRetries) {
		first = false
		eviction := &policyv1.Eviction{ObjectMeta: metav1.ObjectMeta{Namespace: podToEvict.Namespace, Name: podToEvict.Name}, DeleteOptions: &metav1.DeleteOptions{GracePeriodSeconds: &maxTermination}}
		lastError = client.CoreV1().Pods(podToEvict.Namespace).Evict(eviction)
		if lastError == nil || kube_errors.IsNotFound(lastError) {
			return nil
		}
	}
	klog.Errorf("Failed to evict pod %s, error: %v", podToEvict.Name, lastError)
	recorder.Eventf(podToEvict, apiv1.EventTypeWarning, "ScaleDownFailed", "failed to delete pod for ScaleDown")
	return fmt.Errorf("Failed to evict pod %s/%s within allowed timeout (last error: %v)", podToEvict.Namespace, podToEvict.Name, lastError)
}
func drainNode(node *apiv1.Node, pods []*apiv1.Pod, client kube_client.Interface, recorder kube_record.EventRecorder, maxGracefulTerminationSec int, maxPodEvictionTime time.Duration, waitBetweenRetries time.Duration) errors.AutoscalerError {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	toEvict := len(pods)
	retryUntil := time.Now().Add(maxPodEvictionTime)
	confirmations := make(chan error, toEvict)
	for _, pod := range pods {
		go func(podToEvict *apiv1.Pod) {
			confirmations <- evictPod(podToEvict, client, recorder, maxGracefulTerminationSec, retryUntil, waitBetweenRetries)
		}(pod)
	}
	evictionErrs := make([]error, 0)
	for range pods {
		select {
		case err := <-confirmations:
			if err != nil {
				evictionErrs = append(evictionErrs, err)
			} else {
				metrics.RegisterEvictions(1)
			}
		case <-time.After(retryUntil.Sub(time.Now()) + 5*time.Second):
			return errors.NewAutoscalerError(errors.ApiCallError, "Failed to drain node %s/%s: timeout when waiting for creating evictions", node.Namespace, node.Name)
		}
	}
	if len(evictionErrs) != 0 {
		return errors.NewAutoscalerError(errors.ApiCallError, "Failed to drain node %s/%s, due to following errors: %v", node.Namespace, node.Name, evictionErrs)
	}
	allGone := true
	for start := time.Now(); time.Now().Sub(start) < time.Duration(maxGracefulTerminationSec)*time.Second+PodEvictionHeadroom; time.Sleep(5 * time.Second) {
		allGone = true
		for _, pod := range pods {
			podreturned, err := client.CoreV1().Pods(pod.Namespace).Get(pod.Name, metav1.GetOptions{})
			if err == nil && (podreturned == nil || podreturned.Spec.NodeName == node.Name) {
				klog.Errorf("Not deleted yet %v", podreturned)
				allGone = false
				break
			}
			if err != nil && !kube_errors.IsNotFound(err) {
				klog.Errorf("Failed to check pod %s/%s: %v", pod.Namespace, pod.Name, err)
				allGone = false
				break
			}
		}
		if allGone {
			klog.V(1).Infof("All pods removed from %s", node.Name)
			return nil
		}
	}
	return errors.NewAutoscalerError(errors.TransientError, "Failed to drain node %s/%s: pods remaining after timeout", node.Namespace, node.Name)
}
func cleanToBeDeleted(nodes []*apiv1.Node, client kube_client.Interface, recorder kube_record.EventRecorder) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, node := range nodes {
		cleaned, err := deletetaint.CleanToBeDeleted(node, client)
		if err != nil {
			klog.Warningf("Error while releasing taints on node %v: %v", node.Name, err)
			recorder.Eventf(node, apiv1.EventTypeWarning, "ClusterAutoscalerCleanup", "failed to clean toBeDeletedTaint: %v", err)
		} else if cleaned {
			klog.V(1).Infof("Successfully released toBeDeletedTaint on node %v", node.Name)
			recorder.Eventf(node, apiv1.EventTypeNormal, "ClusterAutoscalerCleanup", "marking the node as schedulable")
		}
	}
}
func deleteNodeFromCloudProvider(node *apiv1.Node, cloudProvider cloudprovider.CloudProvider, recorder kube_record.EventRecorder, registry *clusterstate.ClusterStateRegistry) errors.AutoscalerError {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	nodeGroup, err := cloudProvider.NodeGroupForNode(node)
	if err != nil {
		return errors.NewAutoscalerError(errors.CloudProviderError, "failed to find node group for %s: %v", node.Name, err)
	}
	if nodeGroup == nil || reflect.ValueOf(nodeGroup).IsNil() {
		return errors.NewAutoscalerError(errors.InternalError, "picked node that doesn't belong to a node group: %s", node.Name)
	}
	if err = nodeGroup.DeleteNodes([]*apiv1.Node{node}); err != nil {
		return errors.NewAutoscalerError(errors.CloudProviderError, "failed to delete %s: %v", node.Name, err)
	}
	recorder.Eventf(node, apiv1.EventTypeNormal, "ScaleDown", "node removed by cluster autoscaler")
	registry.RegisterScaleDown(&clusterstate.ScaleDownRequest{NodeGroup: nodeGroup, NodeName: node.Name, Time: time.Now(), ExpectedDeleteTime: time.Now().Add(MaxCloudProviderNodeDeletionTime)})
	return nil
}
func hasNoScaleDownAnnotation(node *apiv1.Node) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return node.Annotations[ScaleDownDisabledKey] == "true"
}

const (
	apiServerLabelKey	= "component"
	apiServerLabelValue	= "kube-apiserver"
)

func filterOutMasters(nodes []*apiv1.Node, pods []*apiv1.Pod) []*apiv1.Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	masters := make(map[string]bool)
	for _, pod := range pods {
		if pod.Namespace == metav1.NamespaceSystem && pod.Labels[apiServerLabelKey] == apiServerLabelValue {
			masters[pod.Spec.NodeName] = true
		}
	}
	others := make([]*apiv1.Node, 0, len(nodes)-len(masters))
	for _, node := range nodes {
		if !masters[node.Name] {
			others = append(others, node)
		}
	}
	return others
}
