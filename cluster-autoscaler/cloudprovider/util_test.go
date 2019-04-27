package cloudprovider

import (
	"testing"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
)

func TestBuildReadyConditions(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	conditions := BuildReadyConditions()
	foundReady := false
	for _, condition := range conditions {
		if condition.Type == apiv1.NodeReady && condition.Status == apiv1.ConditionTrue {
			foundReady = true
		}
	}
	assert.True(t, foundReady)
}
func TestBuildKubeProxy(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pod := BuildKubeProxy("kube-proxy")
	assert.NotNil(t, pod)
	assert.Equal(t, 1, len(pod.Spec.Containers))
	cpu := pod.Spec.Containers[0].Resources.Requests[apiv1.ResourceCPU]
	assert.Equal(t, int64(100), cpu.MilliValue())
}
func TestJoinStringMaps(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	map1 := map[string]string{"1": "a", "2": "b"}
	map2 := map[string]string{"3": "c", "2": "d"}
	map3 := map[string]string{"5": "e"}
	result := JoinStringMaps(map1, map2, map3)
	assert.Equal(t, map[string]string{"1": "a", "2": "d", "3": "c", "5": "e"}, result)
}
