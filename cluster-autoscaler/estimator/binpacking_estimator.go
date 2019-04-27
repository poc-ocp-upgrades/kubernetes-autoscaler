package estimator

import (
	"sort"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/autoscaler/cluster-autoscaler/simulator"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

type podInfo struct {
	score	float64
	pod	*apiv1.Pod
}
type BinpackingNodeEstimator struct{ predicateChecker *simulator.PredicateChecker }

func NewBinpackingNodeEstimator(predicateChecker *simulator.PredicateChecker) *BinpackingNodeEstimator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &BinpackingNodeEstimator{predicateChecker: predicateChecker}
}
func (estimator *BinpackingNodeEstimator) Estimate(pods []*apiv1.Pod, nodeTemplate *schedulercache.NodeInfo, upcomingNodes []*schedulercache.NodeInfo) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	podInfos := calculatePodScore(pods, nodeTemplate)
	sort.Slice(podInfos, func(i, j int) bool {
		return podInfos[i].score > podInfos[j].score
	})
	nodeWithPod := func(nodeInfo *schedulercache.NodeInfo, pod *apiv1.Pod) *schedulercache.NodeInfo {
		podsOnNode := nodeInfo.Pods()
		podsOnNode = append(podsOnNode, pod)
		newNodeInfo := schedulercache.NewNodeInfo(podsOnNode...)
		newNodeInfo.SetNode(nodeInfo.Node())
		return newNodeInfo
	}
	newNodes := make([]*schedulercache.NodeInfo, 0)
	newNodes = append(newNodes, upcomingNodes...)
	for _, podInfo := range podInfos {
		found := false
		for i, nodeInfo := range newNodes {
			if err := estimator.predicateChecker.CheckPredicates(podInfo.pod, nil, nodeInfo); err == nil {
				found = true
				newNodes[i] = nodeWithPod(nodeInfo, podInfo.pod)
				break
			}
		}
		if !found {
			newNodes = append(newNodes, nodeWithPod(nodeTemplate, podInfo.pod))
		}
	}
	return len(newNodes) - len(upcomingNodes)
}
func calculatePodScore(pods []*apiv1.Pod, nodeTemplate *schedulercache.NodeInfo) []*podInfo {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	podInfos := make([]*podInfo, 0, len(pods))
	for _, pod := range pods {
		cpuSum := resource.Quantity{}
		memorySum := resource.Quantity{}
		for _, container := range pod.Spec.Containers {
			if request, ok := container.Resources.Requests[apiv1.ResourceCPU]; ok {
				cpuSum.Add(request)
			}
			if request, ok := container.Resources.Requests[apiv1.ResourceMemory]; ok {
				memorySum.Add(request)
			}
		}
		score := float64(0)
		if cpuAllocatable, ok := nodeTemplate.Node().Status.Allocatable[apiv1.ResourceCPU]; ok && cpuAllocatable.MilliValue() > 0 {
			score += float64(cpuSum.MilliValue()) / float64(cpuAllocatable.MilliValue())
		}
		if memAllocatable, ok := nodeTemplate.Node().Status.Allocatable[apiv1.ResourceMemory]; ok && memAllocatable.Value() > 0 {
			score += float64(memorySum.Value()) / float64(memAllocatable.Value())
		}
		podInfos = append(podInfos, &podInfo{score: score, pod: pod})
	}
	return podInfos
}
