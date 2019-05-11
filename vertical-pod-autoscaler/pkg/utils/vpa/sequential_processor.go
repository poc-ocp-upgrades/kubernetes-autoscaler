package api

import (
	"k8s.io/api/core/v1"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
)

func NewSequentialProcessor(processors []RecommendationProcessor) RecommendationProcessor {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &sequentialRecommendationProcessor{processors: processors}
}

type sequentialRecommendationProcessor struct{ processors []RecommendationProcessor }

func (p *sequentialRecommendationProcessor) Apply(podRecommendation *vpa_types.RecommendedPodResources, policy *vpa_types.PodResourcePolicy, conditions []vpa_types.VerticalPodAutoscalerCondition, pod *v1.Pod) (*vpa_types.RecommendedPodResources, ContainerToAnnotationsMap, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	recommendation := podRecommendation
	accumulatedContainerToAnnotationsMap := ContainerToAnnotationsMap{}
	for _, processor := range p.processors {
		var (
			err							error
			containerToAnnotationsMap	ContainerToAnnotationsMap
		)
		recommendation, containerToAnnotationsMap, err = processor.Apply(recommendation, policy, conditions, pod)
		for container, newAnnotations := range containerToAnnotationsMap {
			annotations, found := accumulatedContainerToAnnotationsMap[container]
			if found {
				accumulatedContainerToAnnotationsMap[container] = append(annotations, newAnnotations...)
			} else {
				accumulatedContainerToAnnotationsMap[container] = newAnnotations
			}
		}
		if err != nil {
			return nil, accumulatedContainerToAnnotationsMap, err
		}
	}
	return recommendation, accumulatedContainerToAnnotationsMap, nil
}
