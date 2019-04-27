package core

import (
	"time"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/clusterstate"
	"k8s.io/autoscaler/cluster-autoscaler/clusterstate/utils"
	"k8s.io/autoscaler/cluster-autoscaler/config"
	"k8s.io/autoscaler/cluster-autoscaler/context"
	"k8s.io/autoscaler/cluster-autoscaler/estimator"
	"k8s.io/autoscaler/cluster-autoscaler/expander"
	"k8s.io/autoscaler/cluster-autoscaler/metrics"
	ca_processors "k8s.io/autoscaler/cluster-autoscaler/processors"
	"k8s.io/autoscaler/cluster-autoscaler/simulator"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	"k8s.io/autoscaler/cluster-autoscaler/utils/gpu"
	"k8s.io/autoscaler/cluster-autoscaler/utils/tpu"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/autoscaler/cluster-autoscaler/processors/status"
	"k8s.io/autoscaler/cluster-autoscaler/utils/backoff"
	"k8s.io/klog"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

const (
	unschedulablePodTimeBuffer		= 2 * time.Second
	unschedulablePodWithGpuTimeBuffer	= 30 * time.Second
	nodesNotReadyAfterStartTimeout		= 10 * time.Minute
)

type StaticAutoscaler struct {
	*context.AutoscalingContext
	clusterStateRegistry	*clusterstate.ClusterStateRegistry
	startTime		time.Time
	lastScaleUpTime		time.Time
	lastScaleDownDeleteTime	time.Time
	lastScaleDownFailTime	time.Time
	scaleDown		*ScaleDown
	processors		*ca_processors.AutoscalingProcessors
	initialized		bool
	nodeInfoCache		map[string]*schedulercache.NodeInfo
}

func NewStaticAutoscaler(opts config.AutoscalingOptions, predicateChecker *simulator.PredicateChecker, autoscalingKubeClients *context.AutoscalingKubeClients, processors *ca_processors.AutoscalingProcessors, cloudProvider cloudprovider.CloudProvider, expanderStrategy expander.Strategy, estimatorBuilder estimator.EstimatorBuilder, backoff backoff.Backoff) *StaticAutoscaler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	autoscalingContext := context.NewAutoscalingContext(opts, predicateChecker, autoscalingKubeClients, cloudProvider, expanderStrategy, estimatorBuilder)
	clusterStateConfig := clusterstate.ClusterStateRegistryConfig{MaxTotalUnreadyPercentage: opts.MaxTotalUnreadyPercentage, OkTotalUnreadyCount: opts.OkTotalUnreadyCount, MaxNodeProvisionTime: opts.MaxNodeProvisionTime}
	clusterStateRegistry := clusterstate.NewClusterStateRegistry(autoscalingContext.CloudProvider, clusterStateConfig, autoscalingContext.LogRecorder, backoff)
	scaleDown := NewScaleDown(autoscalingContext, clusterStateRegistry)
	return &StaticAutoscaler{AutoscalingContext: autoscalingContext, startTime: time.Now(), lastScaleUpTime: time.Now(), lastScaleDownDeleteTime: time.Now(), lastScaleDownFailTime: time.Now(), scaleDown: scaleDown, processors: processors, clusterStateRegistry: clusterStateRegistry, nodeInfoCache: make(map[string]*schedulercache.NodeInfo)}
}
func (a *StaticAutoscaler) cleanUpIfRequired() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if a.initialized {
		return
	}
	if readyNodes, err := a.ReadyNodeLister().List(); err != nil {
		klog.Errorf("Failed to list ready nodes, not cleaning up taints: %v", err)
	} else {
		cleanToBeDeleted(readyNodes, a.AutoscalingContext.ClientSet, a.Recorder)
	}
	a.initialized = true
}
func (a *StaticAutoscaler) RunOnce(currentTime time.Time) errors.AutoscalerError {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	a.cleanUpIfRequired()
	unschedulablePodLister := a.UnschedulablePodLister()
	scheduledPodLister := a.ScheduledPodLister()
	pdbLister := a.PodDisruptionBudgetLister()
	scaleDown := a.scaleDown
	autoscalingContext := a.AutoscalingContext
	klog.V(4).Info("Starting main loop")
	stateUpdateStart := time.Now()
	allNodes, readyNodes, typedErr := a.obtainNodeLists()
	if typedErr != nil {
		return typedErr
	}
	if a.actOnEmptyCluster(allNodes, readyNodes, currentTime) {
		return nil
	}
	daemonsets, err := a.ListerRegistry.DaemonSetLister().List()
	if err != nil {
		klog.Errorf("Failed to get daemonset list")
		return errors.ToAutoscalerError(errors.ApiCallError, err)
	}
	nodeInfosForGroups, autoscalerError := GetNodeInfosForGroups(readyNodes, a.nodeInfoCache, autoscalingContext.CloudProvider, autoscalingContext.ClientSet, daemonsets, autoscalingContext.PredicateChecker)
	if err != nil {
		return autoscalerError.AddPrefix("failed to build node infos for node groups: ")
	}
	typedErr = a.updateClusterState(allNodes, nodeInfosForGroups, currentTime)
	if typedErr != nil {
		return typedErr
	}
	metrics.UpdateDurationFromStart(metrics.UpdateState, stateUpdateStart)
	scaleUpStatus := &status.ScaleUpStatus{Result: status.ScaleUpNotTried}
	scaleUpStatusProcessorAlreadyCalled := false
	scaleDownStatus := &status.ScaleDownStatus{Result: status.ScaleDownNotTried}
	scaleDownStatusProcessorAlreadyCalled := false
	defer func() {
		if autoscalingContext.WriteStatusConfigMap {
			status := a.clusterStateRegistry.GetStatus(currentTime)
			utils.WriteStatusConfigMap(autoscalingContext.ClientSet, autoscalingContext.ConfigNamespace, status.GetReadableString(), a.AutoscalingContext.LogRecorder)
		}
		if !scaleUpStatusProcessorAlreadyCalled && a.processors != nil && a.processors.ScaleUpStatusProcessor != nil {
			a.processors.ScaleUpStatusProcessor.Process(a.AutoscalingContext, scaleUpStatus)
		}
		if !scaleDownStatusProcessorAlreadyCalled && a.processors != nil && a.processors.ScaleDownStatusProcessor != nil {
			a.processors.ScaleDownStatusProcessor.Process(a.AutoscalingContext, scaleDownStatus)
		}
		err := a.processors.AutoscalingStatusProcessor.Process(a.AutoscalingContext, a.clusterStateRegistry, currentTime)
		if err != nil {
			klog.Errorf("AutoscalingStatusProcessor error: %v.", err)
		}
	}()
	unregisteredNodes := a.clusterStateRegistry.GetUnregisteredNodes()
	if len(unregisteredNodes) > 0 {
		klog.V(1).Infof("%d unregistered nodes present", len(unregisteredNodes))
		removedAny, err := removeOldUnregisteredNodes(unregisteredNodes, autoscalingContext, currentTime, autoscalingContext.LogRecorder)
		if err != nil {
			if removedAny {
				klog.Warningf("Some unregistered nodes were removed, but got error: %v", err)
			} else {
				klog.Errorf("Failed to remove unregistered nodes: %v", err)
			}
			return errors.ToAutoscalerError(errors.CloudProviderError, err)
		}
		if removedAny {
			klog.V(0).Infof("Some unregistered nodes were removed, skipping iteration")
			return nil
		}
	}
	if !a.clusterStateRegistry.IsClusterHealthy() {
		klog.Warning("Cluster is not ready for autoscaling")
		scaleDown.CleanUpUnneededNodes()
		autoscalingContext.LogRecorder.Eventf(apiv1.EventTypeWarning, "ClusterUnhealthy", "Cluster is unhealthy")
		return nil
	}
	fixedSomething, err := fixNodeGroupSize(autoscalingContext, a.clusterStateRegistry, currentTime)
	if err != nil {
		klog.Errorf("Failed to fix node group sizes: %v", err)
		return errors.ToAutoscalerError(errors.CloudProviderError, err)
	}
	if fixedSomething {
		klog.V(0).Infof("Some node group target size was fixed, skipping the iteration")
		return nil
	}
	metrics.UpdateLastTime(metrics.Autoscaling, time.Now())
	allUnschedulablePods, err := unschedulablePodLister.List()
	if err != nil {
		klog.Errorf("Failed to list unscheduled pods: %v", err)
		return errors.ToAutoscalerError(errors.ApiCallError, err)
	}
	metrics.UpdateUnschedulablePodsCount(len(allUnschedulablePods))
	allScheduled, err := scheduledPodLister.List()
	if err != nil {
		klog.Errorf("Failed to list scheduled pods: %v", err)
		return errors.ToAutoscalerError(errors.ApiCallError, err)
	}
	allUnschedulablePods, allScheduled, err = a.processors.PodListProcessor.Process(a.AutoscalingContext, allUnschedulablePods, allScheduled, allNodes)
	if err != nil {
		klog.Errorf("Failed to process pod list: %v", err)
		return errors.ToAutoscalerError(errors.InternalError, err)
	}
	ConfigurePredicateCheckerForLoop(allUnschedulablePods, allScheduled, a.PredicateChecker)
	scaleDownForbidden := false
	unschedulablePodsWithoutTPUs := tpu.ClearTPURequests(allUnschedulablePods)
	unschedulablePods, unschedulableWaitingForLowerPriorityPreemption := FilterOutExpendableAndSplit(unschedulablePodsWithoutTPUs, a.ExpendablePodsPriorityCutoff)
	klog.V(4).Infof("Filtering out schedulables")
	filterOutSchedulableStart := time.Now()
	unschedulablePodsToHelp := FilterOutSchedulable(unschedulablePods, readyNodes, allScheduled, unschedulableWaitingForLowerPriorityPreemption, a.PredicateChecker, a.ExpendablePodsPriorityCutoff)
	metrics.UpdateDurationFromStart(metrics.FilterOutSchedulable, filterOutSchedulableStart)
	if len(unschedulablePodsToHelp) != len(unschedulablePods) {
		klog.V(2).Info("Schedulable pods present")
		scaleDownForbidden = true
	} else {
		klog.V(4).Info("No schedulable pods")
	}
	unschedulablePodsToHelp = a.filterOutYoungPods(unschedulablePodsToHelp, currentTime)
	if len(unschedulablePodsToHelp) == 0 {
		scaleUpStatus.Result = status.ScaleUpNotNeeded
		klog.V(1).Info("No unschedulable pods")
	} else if a.MaxNodesTotal > 0 && len(readyNodes) >= a.MaxNodesTotal {
		scaleUpStatus.Result = status.ScaleUpNoOptionsAvailable
		klog.V(1).Info("Max total nodes in cluster reached")
	} else if allPodsAreNew(unschedulablePodsToHelp, currentTime) {
		scaleDownForbidden = true
		scaleUpStatus.Result = status.ScaleUpInCooldown
		klog.V(1).Info("Unschedulable pods are very new, waiting one iteration for more")
	} else {
		scaleUpStart := time.Now()
		metrics.UpdateLastTime(metrics.ScaleUp, scaleUpStart)
		scaleUpStatus, typedErr := ScaleUp(autoscalingContext, a.processors, a.clusterStateRegistry, unschedulablePodsToHelp, readyNodes, daemonsets, nodeInfosForGroups)
		metrics.UpdateDurationFromStart(metrics.ScaleUp, scaleUpStart)
		if a.processors != nil && a.processors.ScaleUpStatusProcessor != nil {
			a.processors.ScaleUpStatusProcessor.Process(autoscalingContext, scaleUpStatus)
			scaleUpStatusProcessorAlreadyCalled = true
		}
		if typedErr != nil {
			klog.Errorf("Failed to scale up: %v", typedErr)
			return typedErr
		}
		if scaleUpStatus.Result == status.ScaleUpSuccessful {
			a.lastScaleUpTime = currentTime
			scaleDownStatus.Result = status.ScaleDownInCooldown
			return nil
		}
	}
	if a.ScaleDownEnabled {
		pdbs, err := pdbLister.List()
		if err != nil {
			scaleDownStatus.Result = status.ScaleDownError
			klog.Errorf("Failed to list pod disruption budgets: %v", err)
			return errors.ToAutoscalerError(errors.ApiCallError, err)
		}
		unneededStart := time.Now()
		klog.V(4).Infof("Calculating unneeded nodes")
		scaleDown.CleanUp(currentTime)
		potentiallyUnneeded := getPotentiallyUnneededNodes(autoscalingContext, allNodes)
		typedErr := scaleDown.UpdateUnneededNodes(allNodes, potentiallyUnneeded, append(allScheduled, unschedulableWaitingForLowerPriorityPreemption...), currentTime, pdbs)
		if typedErr != nil {
			scaleDownStatus.Result = status.ScaleDownError
			klog.Errorf("Failed to scale down: %v", typedErr)
			return typedErr
		}
		metrics.UpdateDurationFromStart(metrics.FindUnneeded, unneededStart)
		if klog.V(4) {
			for key, val := range scaleDown.unneededNodes {
				klog.Infof("%s is unneeded since %s duration %s", key, val.String(), currentTime.Sub(val).String())
			}
		}
		scaleDownInCooldown := scaleDownForbidden || a.lastScaleUpTime.Add(a.ScaleDownDelayAfterAdd).After(currentTime) || a.lastScaleDownFailTime.Add(a.ScaleDownDelayAfterFailure).After(currentTime) || a.lastScaleDownDeleteTime.Add(a.ScaleDownDelayAfterDelete).After(currentTime)
		calculateUnneededOnly := scaleDownInCooldown || scaleDown.nodeDeleteStatus.IsDeleteInProgress()
		klog.V(4).Infof("Scale down status: unneededOnly=%v lastScaleUpTime=%s "+"lastScaleDownDeleteTime=%v lastScaleDownFailTime=%s scaleDownForbidden=%v isDeleteInProgress=%v", calculateUnneededOnly, a.lastScaleUpTime, a.lastScaleDownDeleteTime, a.lastScaleDownFailTime, scaleDownForbidden, scaleDown.nodeDeleteStatus.IsDeleteInProgress())
		if scaleDownInCooldown {
			scaleDownStatus.Result = status.ScaleDownInCooldown
		} else if scaleDown.nodeDeleteStatus.IsDeleteInProgress() {
			scaleDownStatus.Result = status.ScaleDownInProgress
		} else {
			klog.V(4).Infof("Starting scale down")
			a.processors.NodeGroupManager.RemoveUnneededNodeGroups(autoscalingContext)
			scaleDownStart := time.Now()
			metrics.UpdateLastTime(metrics.ScaleDown, scaleDownStart)
			scaleDownStatus, typedErr := scaleDown.TryToScaleDown(allNodes, allScheduled, pdbs, currentTime)
			metrics.UpdateDurationFromStart(metrics.ScaleDown, scaleDownStart)
			if scaleDownStatus.Result == status.ScaleDownNodeDeleted {
				a.lastScaleDownDeleteTime = currentTime
				a.clusterStateRegistry.Recalculate()
			}
			if a.processors != nil && a.processors.ScaleDownStatusProcessor != nil {
				a.processors.ScaleDownStatusProcessor.Process(autoscalingContext, scaleDownStatus)
				scaleDownStatusProcessorAlreadyCalled = true
			}
			if typedErr != nil {
				klog.Errorf("Failed to scale down: %v", err)
				a.lastScaleDownFailTime = currentTime
				return typedErr
			}
		}
	}
	return nil
}
func (a *StaticAutoscaler) filterOutYoungPods(allUnschedulablePods []*apiv1.Pod, currentTime time.Time) []*apiv1.Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var oldUnschedulablePods []*apiv1.Pod
	newPodScaleUpDelay := a.AutoscalingOptions.NewPodScaleUpDelay
	for _, pod := range allUnschedulablePods {
		podAge := currentTime.Sub(pod.CreationTimestamp.Time)
		if podAge > newPodScaleUpDelay {
			oldUnschedulablePods = append(oldUnschedulablePods, pod)
		} else {
			klog.V(3).Infof("Pod %s is %.3f seconds old, too new to consider unschedulable", pod.Name, podAge.Seconds())
		}
	}
	return oldUnschedulablePods
}
func (a *StaticAutoscaler) ExitCleanUp() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	a.processors.CleanUp()
	if !a.AutoscalingContext.WriteStatusConfigMap {
		return
	}
	utils.DeleteStatusConfigMap(a.AutoscalingContext.ClientSet, a.AutoscalingContext.ConfigNamespace)
}
func (a *StaticAutoscaler) obtainNodeLists() ([]*apiv1.Node, []*apiv1.Node, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	allNodes, err := a.AllNodeLister().List()
	if err != nil {
		klog.Errorf("Failed to list all nodes: %v", err)
		return nil, nil, errors.ToAutoscalerError(errors.ApiCallError, err)
	}
	readyNodes, err := a.ReadyNodeLister().List()
	if err != nil {
		klog.Errorf("Failed to list ready nodes: %v", err)
		return nil, nil, errors.ToAutoscalerError(errors.ApiCallError, err)
	}
	allNodes, readyNodes = gpu.FilterOutNodesWithUnreadyGpus(allNodes, readyNodes)
	return allNodes, readyNodes, nil
}
func (a *StaticAutoscaler) actOnEmptyCluster(allNodes, readyNodes []*apiv1.Node, currentTime time.Time) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(allNodes) == 0 {
		a.onEmptyCluster("Cluster has no nodes.", true)
		return true
	}
	if len(readyNodes) == 0 {
		a.onEmptyCluster("Cluster has no ready nodes.", currentTime.After(a.startTime.Add(nodesNotReadyAfterStartTimeout)))
		return true
	}
	return false
}
func (a *StaticAutoscaler) updateClusterState(allNodes []*apiv1.Node, nodeInfosForGroups map[string]*schedulercache.NodeInfo, currentTime time.Time) errors.AutoscalerError {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	err := a.AutoscalingContext.CloudProvider.Refresh()
	if err != nil {
		klog.Errorf("Failed to refresh cloud provider config: %v", err)
		return errors.ToAutoscalerError(errors.CloudProviderError, err)
	}
	err = a.clusterStateRegistry.UpdateNodes(allNodes, nodeInfosForGroups, currentTime)
	if err != nil {
		klog.Errorf("Failed to update node registry: %v", err)
		a.scaleDown.CleanUpUnneededNodes()
		return errors.ToAutoscalerError(errors.CloudProviderError, err)
	}
	UpdateClusterStateMetrics(a.clusterStateRegistry)
	return nil
}
func (a *StaticAutoscaler) onEmptyCluster(status string, emitEvent bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	klog.Warningf(status)
	a.scaleDown.CleanUpUnneededNodes()
	UpdateEmptyClusterStateMetrics()
	if a.AutoscalingContext.WriteStatusConfigMap {
		utils.WriteStatusConfigMap(a.AutoscalingContext.ClientSet, a.AutoscalingContext.ConfigNamespace, status, a.AutoscalingContext.LogRecorder)
	}
	if emitEvent {
		a.AutoscalingContext.LogRecorder.Eventf(apiv1.EventTypeWarning, "ClusterUnhealthy", status)
	}
}
func allPodsAreNew(pods []*apiv1.Pod, currentTime time.Time) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if getOldestCreateTime(pods).Add(unschedulablePodTimeBuffer).After(currentTime) {
		return true
	}
	found, oldest := getOldestCreateTimeWithGpu(pods)
	return found && oldest.Add(unschedulablePodWithGpuTimeBuffer).After(currentTime)
}
