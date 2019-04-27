package estimator

import (
	"testing"
	"time"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/autoscaler/cluster-autoscaler/simulator"
	. "k8s.io/autoscaler/cluster-autoscaler/utils/test"
	"k8s.io/autoscaler/cluster-autoscaler/utils/units"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
	"github.com/stretchr/testify/assert"
)

func TestBinpackingEstimate(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	estimator := NewBinpackingNodeEstimator(simulator.NewTestPredicateChecker())
	cpuPerPod := int64(350)
	memoryPerPod := int64(1000 * units.MiB)
	pod := makePod(cpuPerPod, memoryPerPod)
	pods := make([]*apiv1.Pod, 0)
	for i := 0; i < 10; i++ {
		pods = append(pods, pod)
	}
	node := &apiv1.Node{Status: apiv1.NodeStatus{Capacity: apiv1.ResourceList{apiv1.ResourceCPU: *resource.NewMilliQuantity(cpuPerPod*3-50, resource.DecimalSI), apiv1.ResourceMemory: *resource.NewQuantity(2*memoryPerPod, resource.DecimalSI), apiv1.ResourcePods: *resource.NewQuantity(10, resource.DecimalSI)}}}
	node.Status.Allocatable = node.Status.Capacity
	SetNodeReadyState(node, true, time.Time{})
	nodeInfo := schedulercache.NewNodeInfo()
	nodeInfo.SetNode(node)
	estimate := estimator.Estimate(pods, nodeInfo, []*schedulercache.NodeInfo{})
	assert.Equal(t, 5, estimate)
}
func TestBinpackingEstimateComingNodes(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	estimator := NewBinpackingNodeEstimator(simulator.NewTestPredicateChecker())
	cpuPerPod := int64(350)
	memoryPerPod := int64(1000 * units.MiB)
	pod := makePod(cpuPerPod, memoryPerPod)
	pods := make([]*apiv1.Pod, 0)
	for i := 0; i < 10; i++ {
		pods = append(pods, pod)
	}
	node := &apiv1.Node{Status: apiv1.NodeStatus{Capacity: apiv1.ResourceList{apiv1.ResourceCPU: *resource.NewMilliQuantity(cpuPerPod*3-50, resource.DecimalSI), apiv1.ResourceMemory: *resource.NewQuantity(2*memoryPerPod, resource.DecimalSI), apiv1.ResourcePods: *resource.NewQuantity(10, resource.DecimalSI)}}}
	node.Status.Allocatable = node.Status.Capacity
	SetNodeReadyState(node, true, time.Time{})
	nodeInfo := schedulercache.NewNodeInfo()
	nodeInfo.SetNode(node)
	estimate := estimator.Estimate(pods, nodeInfo, []*schedulercache.NodeInfo{nodeInfo, nodeInfo})
	assert.Equal(t, 3, estimate)
}
func TestBinpackingEstimateWithPorts(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	estimator := NewBinpackingNodeEstimator(simulator.NewTestPredicateChecker())
	cpuPerPod := int64(200)
	memoryPerPod := int64(1000 * units.MiB)
	pod := makePod(cpuPerPod, memoryPerPod)
	pod.Spec.Containers[0].Ports = []apiv1.ContainerPort{{HostPort: 5555}}
	pods := make([]*apiv1.Pod, 0)
	for i := 0; i < 8; i++ {
		pods = append(pods, pod)
	}
	node := &apiv1.Node{Status: apiv1.NodeStatus{Capacity: apiv1.ResourceList{apiv1.ResourceCPU: *resource.NewMilliQuantity(5*cpuPerPod, resource.DecimalSI), apiv1.ResourceMemory: *resource.NewQuantity(5*memoryPerPod, resource.DecimalSI), apiv1.ResourcePods: *resource.NewQuantity(10, resource.DecimalSI)}}}
	node.Status.Allocatable = node.Status.Capacity
	SetNodeReadyState(node, true, time.Time{})
	nodeInfo := schedulercache.NewNodeInfo()
	nodeInfo.SetNode(node)
	estimate := estimator.Estimate(pods, nodeInfo, []*schedulercache.NodeInfo{})
	assert.Equal(t, 8, estimate)
}
