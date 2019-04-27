package estimator

import (
	"testing"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/autoscaler/cluster-autoscaler/utils/units"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
	"github.com/stretchr/testify/assert"
)

func makePod(cpuPerPod, memoryPerPod int64) *apiv1.Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &apiv1.Pod{Spec: apiv1.PodSpec{Containers: []apiv1.Container{{Resources: apiv1.ResourceRequirements{Requests: apiv1.ResourceList{apiv1.ResourceCPU: *resource.NewMilliQuantity(cpuPerPod, resource.DecimalSI), apiv1.ResourceMemory: *resource.NewQuantity(memoryPerPod, resource.DecimalSI)}}}}}}
}
func TestEstimate(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cpuPerPod := int64(500)
	memoryPerPod := int64(1000 * units.MiB)
	pod := makePod(cpuPerPod, memoryPerPod)
	pods := []*apiv1.Pod{}
	for i := 0; i < 5; i++ {
		podCopy := *pod
		pods = append(pods, &podCopy)
	}
	node := &apiv1.Node{Status: apiv1.NodeStatus{Capacity: apiv1.ResourceList{apiv1.ResourceCPU: *resource.NewMilliQuantity(3*cpuPerPod, resource.DecimalSI), apiv1.ResourceMemory: *resource.NewQuantity(2*memoryPerPod, resource.DecimalSI), apiv1.ResourcePods: *resource.NewQuantity(10, resource.DecimalSI)}}}
	nodeInfo := schedulercache.NewNodeInfo()
	nodeInfo.SetNode(node)
	estimator := NewBasicNodeEstimator()
	estimate := estimator.Estimate(pods, nodeInfo, []*schedulercache.NodeInfo{})
	assert.Equal(t, 3, estimate)
	assert.Equal(t, int64(500*5), estimator.cpuSum.MilliValue())
	assert.Equal(t, int64(5*memoryPerPod), estimator.memorySum.Value())
	assert.Equal(t, 5, estimator.GetCount())
	assert.Contains(t, estimator.GetDebug(), "CPU")
}
func TestEstimateWithComing(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cpuPerPod := int64(500)
	memoryPerPod := int64(1000 * units.MiB)
	pod := makePod(cpuPerPod, memoryPerPod)
	pods := []*apiv1.Pod{}
	for i := 0; i < 5; i++ {
		podCopy := *pod
		pods = append(pods, &podCopy)
	}
	node := &apiv1.Node{Status: apiv1.NodeStatus{Capacity: apiv1.ResourceList{apiv1.ResourceCPU: *resource.NewMilliQuantity(3*cpuPerPod, resource.DecimalSI), apiv1.ResourceMemory: *resource.NewQuantity(2*memoryPerPod, resource.DecimalSI), apiv1.ResourcePods: *resource.NewQuantity(10, resource.DecimalSI)}}}
	node.Status.Allocatable = node.Status.Capacity
	nodeInfo := schedulercache.NewNodeInfo()
	nodeInfo.SetNode(node)
	estimator := NewBasicNodeEstimator()
	estimate := estimator.Estimate(pods, nodeInfo, []*schedulercache.NodeInfo{nodeInfo, nodeInfo})
	assert.Equal(t, 1, estimate)
	assert.Contains(t, estimator.GetDebug(), "CPU")
	assert.Equal(t, int64(500*5), estimator.cpuSum.MilliValue())
	assert.Equal(t, int64(5*memoryPerPod), estimator.memorySum.Value())
	assert.Equal(t, 5, estimator.GetCount())
}
func TestEstimateWithPorts(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cpuPerPod := int64(500)
	memoryPerPod := int64(1000 * units.MiB)
	pod := makePod(cpuPerPod, memoryPerPod)
	pod.Spec.Containers[0].Ports = []apiv1.ContainerPort{{HostPort: 5555}}
	pods := []*apiv1.Pod{}
	for i := 0; i < 5; i++ {
		pods = append(pods, pod)
	}
	node := &apiv1.Node{Status: apiv1.NodeStatus{Capacity: apiv1.ResourceList{apiv1.ResourceCPU: *resource.NewMilliQuantity(3*cpuPerPod, resource.DecimalSI), apiv1.ResourceMemory: *resource.NewQuantity(2*memoryPerPod, resource.DecimalSI), apiv1.ResourcePods: *resource.NewQuantity(10, resource.DecimalSI)}}}
	nodeInfo := schedulercache.NewNodeInfo()
	nodeInfo.SetNode(node)
	estimator := NewBasicNodeEstimator()
	estimate := estimator.Estimate(pods, nodeInfo, []*schedulercache.NodeInfo{})
	assert.Contains(t, estimator.GetDebug(), "CPU")
	assert.Equal(t, 5, estimate)
}
