package scheduler

import (
	"testing"
	. "k8s.io/autoscaler/cluster-autoscaler/utils/test"
	apiv1 "k8s.io/api/core/v1"
	"github.com/stretchr/testify/assert"
)

func TestCreateNodeNameToInfoMap(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	p1 := BuildTestPod("p1", 1500, 200000)
	p1.Spec.NodeName = "node1"
	p2 := BuildTestPod("p2", 3000, 200000)
	p2.Spec.NodeName = "node2"
	p3 := BuildTestPod("p3", 3000, 200000)
	p3.Spec.NodeName = "node3"
	var priority int32 = 100
	podWaitingForPreemption := BuildTestPod("w1", 1500, 200000)
	podWaitingForPreemption.Spec.Priority = &priority
	podWaitingForPreemption.Annotations = map[string]string{NominatedNodeAnnotationKey: "node1"}
	n1 := BuildTestNode("node1", 2000, 2000000)
	n2 := BuildTestNode("node2", 2000, 2000000)
	res := CreateNodeNameToInfoMap([]*apiv1.Pod{p1, p2, p3, podWaitingForPreemption}, []*apiv1.Node{n1, n2})
	assert.Equal(t, 2, len(res))
	assert.Equal(t, p1, res["node1"].Pods()[0])
	assert.Equal(t, podWaitingForPreemption, res["node1"].Pods()[1])
	assert.Equal(t, n1, res["node1"].Node())
	assert.Equal(t, p2, res["node2"].Pods()[0])
	assert.Equal(t, n2, res["node2"].Node())
}
