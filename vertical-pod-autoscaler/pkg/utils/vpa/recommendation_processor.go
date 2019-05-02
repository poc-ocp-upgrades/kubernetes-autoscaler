package api

import (
 "k8s.io/api/core/v1"
 vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
)

type ContainerToAnnotationsMap = map[string][]string
type RecommendationProcessor interface {
 Apply(podRecommendation *vpa_types.RecommendedPodResources, policy *vpa_types.PodResourcePolicy, conditions []vpa_types.VerticalPodAutoscalerCondition, pod *v1.Pod) (*vpa_types.RecommendedPodResources, ContainerToAnnotationsMap, error)
}
