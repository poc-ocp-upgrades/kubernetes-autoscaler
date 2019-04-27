package logic

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/recommender/model"
)

func TestMinResourcesApplied(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	constEstimator := NewConstEstimator(model.Resources{model.ResourceCPU: model.CPUAmountFromCores(0.001), model.ResourceMemory: model.MemoryAmountFromBytes(1e6)})
	recommender := podResourceRecommender{constEstimator, constEstimator, constEstimator}
	containerNameToAggregateStateMap := model.ContainerNameToAggregateStateMap{"container-1": &model.AggregateContainerState{}}
	recommendedResources := recommender.GetRecommendedPodResources(containerNameToAggregateStateMap)
	assert.Equal(t, model.CPUAmountFromCores(*podMinCPUMillicores/1000), recommendedResources["container-1"].Target[model.ResourceCPU])
	assert.Equal(t, model.MemoryAmountFromBytes(*podMinMemoryMb*1024*1024), recommendedResources["container-1"].Target[model.ResourceMemory])
}
func TestMinResourcesSplitAcrossContainers(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	constEstimator := NewConstEstimator(model.Resources{model.ResourceCPU: model.CPUAmountFromCores(0.001), model.ResourceMemory: model.MemoryAmountFromBytes(1e6)})
	recommender := podResourceRecommender{constEstimator, constEstimator, constEstimator}
	containerNameToAggregateStateMap := model.ContainerNameToAggregateStateMap{"container-1": &model.AggregateContainerState{}, "container-2": &model.AggregateContainerState{}}
	recommendedResources := recommender.GetRecommendedPodResources(containerNameToAggregateStateMap)
	assert.Equal(t, model.CPUAmountFromCores((*podMinCPUMillicores/1000)/2), recommendedResources["container-1"].Target[model.ResourceCPU])
	assert.Equal(t, model.CPUAmountFromCores((*podMinCPUMillicores/1000)/2), recommendedResources["container-1"].Target[model.ResourceCPU])
	assert.Equal(t, model.MemoryAmountFromBytes((*podMinMemoryMb*1024*1024)/2), recommendedResources["container-2"].Target[model.ResourceMemory])
	assert.Equal(t, model.MemoryAmountFromBytes((*podMinMemoryMb*1024*1024)/2), recommendedResources["container-2"].Target[model.ResourceMemory])
}
