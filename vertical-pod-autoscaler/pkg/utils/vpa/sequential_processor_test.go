package api

import (
	"testing"
	"k8s.io/api/core/v1"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	"github.com/stretchr/testify/assert"
)

type fakeProcessor struct{ message string }

func (p *fakeProcessor) Apply(podRecommendation *vpa_types.RecommendedPodResources, policy *vpa_types.PodResourcePolicy, conditions []vpa_types.VerticalPodAutoscalerCondition, pod *v1.Pod) (*vpa_types.RecommendedPodResources, ContainerToAnnotationsMap, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := podRecommendation.DeepCopy()
	result.ContainerRecommendations[0].ContainerName += p.message
	containerToAnnotationsMap := ContainerToAnnotationsMap{"trace": []string{p.message}}
	return result, containerToAnnotationsMap, nil
}
func TestSequentialProcessor(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	name1 := "processor1"
	name2 := "processor2"
	tested := NewSequentialProcessor([]RecommendationProcessor{&fakeProcessor{name1}, &fakeProcessor{name2}})
	rec1 := &vpa_types.RecommendedPodResources{ContainerRecommendations: []vpa_types.RecommendedContainerResources{{ContainerName: ""}}}
	result, annotations, _ := tested.Apply(rec1, nil, nil, nil)
	assert.Equal(t, name1+name2, result.ContainerRecommendations[0].ContainerName)
	assert.Contains(t, annotations, "trace")
	assert.Contains(t, annotations["trace"], name1)
	assert.Contains(t, annotations["trace"], name2)
}
