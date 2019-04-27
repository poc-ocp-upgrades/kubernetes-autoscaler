package nodegroupset

import (
	"testing"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/autoscaler/cluster-autoscaler/utils/gpu"
	. "k8s.io/autoscaler/cluster-autoscaler/utils/test"
	kubeletapis "k8s.io/kubernetes/pkg/kubelet/apis"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
	"github.com/stretchr/testify/assert"
)

func checkNodesSimilar(t *testing.T, n1, n2 *apiv1.Node, comparator NodeInfoComparator, shouldEqual bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	checkNodesSimilarWithPods(t, n1, n2, []*apiv1.Pod{}, []*apiv1.Pod{}, comparator, shouldEqual)
}
func checkNodesSimilarWithPods(t *testing.T, n1, n2 *apiv1.Node, pods1, pods2 []*apiv1.Pod, comparator NodeInfoComparator, shouldEqual bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ni1 := schedulercache.NewNodeInfo(pods1...)
	ni1.SetNode(n1)
	ni2 := schedulercache.NewNodeInfo(pods2...)
	ni2.SetNode(n2)
	assert.Equal(t, shouldEqual, comparator(ni1, ni2))
}
func TestIdenticalNodesSimilar(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	n1 := BuildTestNode("node1", 1000, 2000)
	n2 := BuildTestNode("node2", 1000, 2000)
	checkNodesSimilar(t, n1, n2, IsNodeInfoSimilar, true)
}
func TestNodesSimilarVariousRequirements(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	n1 := BuildTestNode("node1", 1000, 2000)
	n2 := BuildTestNode("node2", 1000, 2000)
	n2.Status.Capacity[apiv1.ResourceCPU] = *resource.NewMilliQuantity(1001, resource.DecimalSI)
	checkNodesSimilar(t, n1, n2, IsNodeInfoSimilar, false)
	n3 := BuildTestNode("node3", 1000, 2000)
	n3.Status.Allocatable[apiv1.ResourceCPU] = *resource.NewMilliQuantity(999, resource.DecimalSI)
	checkNodesSimilar(t, n1, n3, IsNodeInfoSimilar, true)
	n4 := BuildTestNode("node4", 1000, 2000)
	n4.Status.Allocatable[apiv1.ResourceCPU] = *resource.NewMilliQuantity(500, resource.DecimalSI)
	checkNodesSimilar(t, n1, n4, IsNodeInfoSimilar, false)
	n5 := BuildTestNode("node5", 1000, 2000)
	n5.Status.Capacity[gpu.ResourceNvidiaGPU] = *resource.NewQuantity(1, resource.DecimalSI)
	n5.Status.Allocatable[gpu.ResourceNvidiaGPU] = n5.Status.Capacity[gpu.ResourceNvidiaGPU]
	checkNodesSimilar(t, n1, n5, IsNodeInfoSimilar, false)
}
func TestNodesSimilarVariousRequirementsAndPods(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	n1 := BuildTestNode("node1", 1000, 2000)
	p1 := BuildTestPod("pod1", 500, 1000)
	p1.Spec.NodeName = "node1"
	n2 := BuildTestNode("node2", 1000, 2000)
	n2.Status.Allocatable[apiv1.ResourceCPU] = *resource.NewMilliQuantity(500, resource.DecimalSI)
	n2.Status.Allocatable[apiv1.ResourceMemory] = *resource.NewQuantity(1000, resource.DecimalSI)
	checkNodesSimilarWithPods(t, n1, n2, []*apiv1.Pod{p1}, []*apiv1.Pod{}, IsNodeInfoSimilar, false)
	n3 := BuildTestNode("node3", 1000, 2000)
	p3 := BuildTestPod("pod3", 500, 1000)
	p3.Spec.NodeName = "node3"
	checkNodesSimilarWithPods(t, n1, n3, []*apiv1.Pod{p1}, []*apiv1.Pod{p3}, IsNodeInfoSimilar, true)
	n4 := BuildTestNode("node4", 1000, 2000)
	n4.Status.Allocatable[apiv1.ResourceCPU] = *resource.NewMilliQuantity(999, resource.DecimalSI)
	p4 := BuildTestPod("pod4", 501, 1001)
	p4.Spec.NodeName = "node4"
	checkNodesSimilarWithPods(t, n1, n4, []*apiv1.Pod{p1}, []*apiv1.Pod{p4}, IsNodeInfoSimilar, true)
}
func TestNodesSimilarVariousLabels(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	n1 := BuildTestNode("node1", 1000, 2000)
	n1.ObjectMeta.Labels["test-label"] = "test-value"
	n1.ObjectMeta.Labels["character"] = "winnie the pooh"
	n2 := BuildTestNode("node2", 1000, 2000)
	n2.ObjectMeta.Labels["test-label"] = "test-value"
	checkNodesSimilar(t, n1, n2, IsNodeInfoSimilar, false)
	n2.ObjectMeta.Labels["character"] = "winnie the pooh"
	checkNodesSimilar(t, n1, n2, IsNodeInfoSimilar, true)
	n1.ObjectMeta.Labels[kubeletapis.LabelHostname] = "node1"
	n2.ObjectMeta.Labels[kubeletapis.LabelHostname] = "node2"
	checkNodesSimilar(t, n1, n2, IsNodeInfoSimilar, true)
	n1.ObjectMeta.Labels[kubeletapis.LabelZoneFailureDomain] = "mars-olympus-mons1-b"
	n2.ObjectMeta.Labels[kubeletapis.LabelZoneFailureDomain] = "us-houston1-a"
	checkNodesSimilar(t, n1, n2, IsNodeInfoSimilar, true)
	n1.ObjectMeta.Labels["beta.kubernetes.io/fluentd-ds-ready"] = "true"
	n2.ObjectMeta.Labels["beta.kubernetes.io/fluentd-ds-ready"] = "false"
	checkNodesSimilar(t, n1, n2, IsNodeInfoSimilar, true)
	n1.ObjectMeta.Labels["beta.kubernetes.io/fluentd-ds-ready"] = "true"
	delete(n2.ObjectMeta.Labels, "beta.kubernetes.io/fluentd-ds-ready")
	checkNodesSimilar(t, n1, n2, IsNodeInfoSimilar, true)
}
