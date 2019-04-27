package labels

import (
	"testing"
	. "k8s.io/autoscaler/cluster-autoscaler/utils/test"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
)

func TestCalculateNodeSelectorStats(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	p1 := BuildTestPod("p1", 500, 0)
	p1.Spec.NodeSelector = map[string]string{"A": "X", "B": "Y"}
	p2 := BuildTestPod("p2", 800, 0)
	p2.Spec.NodeSelector = map[string]string{"A": "Z12345"}
	p3 := BuildTestPod("p3", 100, 0)
	p3.Spec.NodeSelector = map[string]string{"A": "X", "B": "Y"}
	pods := []*apiv1.Pod{p1, p2, p3}
	stats := calculateNodeSelectorStats(pods)
	sortNodeSelectorStats(stats)
	assert.Equal(t, 2, len(stats))
	assert.Equal(t, p2.Spec.NodeSelector, stats[0].nodeSelector)
	assert.Equal(t, int64(800), stats[0].totalCpu.MilliValue())
	assert.Equal(t, p1.Spec.NodeSelector, stats[1].nodeSelector)
	assert.Equal(t, int64(600), stats[1].totalCpu.MilliValue())
}
func TestBestLabelSet(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	p1 := BuildTestPod("p1", 500, 0)
	p1.Spec.NodeSelector = map[string]string{"A": "X", "C": "Y"}
	p2 := BuildTestPod("p2", 200, 0)
	p2.Spec.NodeSelector = map[string]string{"A": "Z12345"}
	p3 := BuildTestPod("p3", 100, 0)
	p3.Spec.NodeSelector = map[string]string{"A": "X", "B": "Y"}
	p4 := BuildTestPod("p3", 100, 0)
	p4.Spec.NodeSelector = map[string]string{"A": "X", "cloud.google.com/gke": "true"}
	expectedResult := map[string]string{"A": "X", "C": "Y", "B": "Y"}
	assert.Equal(t, expectedResult, BestLabelSet([]*apiv1.Pod{p1, p2, p3, p4}))
}
