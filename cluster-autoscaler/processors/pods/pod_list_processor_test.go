package pods

import (
	"testing"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/autoscaler/cluster-autoscaler/context"
	. "k8s.io/autoscaler/cluster-autoscaler/utils/test"
)

func TestPodListProcessor(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	context := &context.AutoscalingContext{}
	p1 := BuildTestPod("p1", 40, 0)
	p2 := BuildTestPod("p2", 400, 0)
	n1 := BuildTestNode("n1", 100, 1000)
	n2 := BuildTestNode("n1", 100, 1000)
	unschedulablePods := []*apiv1.Pod{p1}
	allScheduled := []*apiv1.Pod{p2}
	nodes := []*apiv1.Node{n1, n2}
	podListProcessor := NewDefaultPodListProcessor()
	gotUnschedulablePods, gotAllScheduled, err := podListProcessor.Process(context, unschedulablePods, allScheduled, nodes)
	if len(gotUnschedulablePods) != 1 || len(gotAllScheduled) != 1 || err != nil {
		t.Errorf("Error podListProcessor.Process() = %v, %v, %v want %v, %v, nil ", gotUnschedulablePods, gotAllScheduled, err, unschedulablePods, allScheduled)
	}
}
