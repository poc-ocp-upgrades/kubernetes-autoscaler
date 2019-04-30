package test

import (
	apiv1 "k8s.io/api/core/v1"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
)

type RecommendationBuilder interface {
	WithContainer(containerName string) RecommendationBuilder
	WithTarget(cpu, memory string) RecommendationBuilder
	WithLowerBound(cpu, memory string) RecommendationBuilder
	WithUpperBound(cpu, memory string) RecommendationBuilder
	Get() *vpa_types.RecommendedPodResources
}

func Recommendation() RecommendationBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &recommendationBuilder{}
}

type recommendationBuilder struct {
	containerName	string
	target		apiv1.ResourceList
	lowerBound	apiv1.ResourceList
	upperBound	apiv1.ResourceList
}

func (b *recommendationBuilder) WithContainer(containerName string) RecommendationBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := *b
	c.containerName = containerName
	return &c
}
func (b *recommendationBuilder) WithTarget(cpu, memory string) RecommendationBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := *b
	c.target = Resources(cpu, memory)
	return &c
}
func (b *recommendationBuilder) WithLowerBound(cpu, memory string) RecommendationBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := *b
	c.lowerBound = Resources(cpu, memory)
	return &c
}
func (b *recommendationBuilder) WithUpperBound(cpu, memory string) RecommendationBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := *b
	c.upperBound = Resources(cpu, memory)
	return &c
}
func (b *recommendationBuilder) Get() *vpa_types.RecommendedPodResources {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if b.containerName == "" {
		panic("Must call WithContainer() before Get()")
	}
	return &vpa_types.RecommendedPodResources{ContainerRecommendations: []vpa_types.RecommendedContainerResources{{ContainerName: b.containerName, Target: b.target, LowerBound: b.lowerBound, UpperBound: b.upperBound}}}
}
