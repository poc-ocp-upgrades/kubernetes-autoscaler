package api

import (
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	"github.com/golang/glog"
)

func NewCappingRecommendationProcessor() RecommendationProcessor {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &cappingRecommendationProcessor{}
}

type cappingAction string

var (
	cappedToMinAllowed	cappingAction	= "capped to minAllowed"
	cappedToMaxAllowed	cappingAction	= "capped to maxAllowed"
	cappedToLimit		cappingAction	= "capped to container limit"
)

func toCappingAnnotation(resourceName apiv1.ResourceName, action cappingAction) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%s %s", resourceName, action)
}

type cappingRecommendationProcessor struct{}

func (c *cappingRecommendationProcessor) Apply(podRecommendation *vpa_types.RecommendedPodResources, policy *vpa_types.PodResourcePolicy, conditions []vpa_types.VerticalPodAutoscalerCondition, pod *apiv1.Pod) (*vpa_types.RecommendedPodResources, ContainerToAnnotationsMap, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if podRecommendation == nil && policy == nil {
		return nil, nil, nil
	}
	if podRecommendation == nil {
		podRecommendation = new(vpa_types.RecommendedPodResources)
	}
	updatedRecommendations := []vpa_types.RecommendedContainerResources{}
	containerToAnnotationsMap := ContainerToAnnotationsMap{}
	for _, containerRecommendation := range podRecommendation.ContainerRecommendations {
		container := getContainer(containerRecommendation.ContainerName, pod)
		if container == nil {
			glog.V(2).Infof("no matching Container found for recommendation %s", containerRecommendation.ContainerName)
			continue
		}
		updatedContainerResources, containerAnnotations, err := getCappedRecommendationForContainer(*container, &containerRecommendation, policy)
		if len(containerAnnotations) != 0 {
			containerToAnnotationsMap[containerRecommendation.ContainerName] = containerAnnotations
		}
		if err != nil {
			return nil, nil, fmt.Errorf("cannot update recommendation for container name %v", container.Name)
		}
		updatedRecommendations = append(updatedRecommendations, *updatedContainerResources)
	}
	return &vpa_types.RecommendedPodResources{ContainerRecommendations: updatedRecommendations}, containerToAnnotationsMap, nil
}
func getCappedRecommendationForContainer(container apiv1.Container, containerRecommendation *vpa_types.RecommendedContainerResources, policy *vpa_types.PodResourcePolicy) (*vpa_types.RecommendedContainerResources, []string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if containerRecommendation == nil {
		return nil, nil, fmt.Errorf("no recommendation available for container name %v", container.Name)
	}
	containerPolicy := GetContainerResourcePolicy(container.Name, policy)
	cappedRecommendations := containerRecommendation.DeepCopy()
	cappingAnnotations := make([]string, 0)
	process := func(recommendation apiv1.ResourceList, genAnnotations bool) {
		annotations := applyVPAPolicy(recommendation, containerPolicy)
		if genAnnotations {
			cappingAnnotations = append(cappingAnnotations, annotations...)
		}
		annotations = capRecommendationToContainerLimit(recommendation, container)
		if genAnnotations {
			cappingAnnotations = append(cappingAnnotations, annotations...)
		}
	}
	process(cappedRecommendations.Target, true)
	process(cappedRecommendations.LowerBound, false)
	process(cappedRecommendations.UpperBound, false)
	return cappedRecommendations, cappingAnnotations, nil
}
func capRecommendationToContainerLimit(recommendation apiv1.ResourceList, container apiv1.Container) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	annotations := make([]string, 0)
	for resourceName, limit := range container.Resources.Limits {
		recommendedValue, found := recommendation[resourceName]
		if found && recommendedValue.MilliValue() > limit.MilliValue() {
			recommendation[resourceName] = limit
			annotations = append(annotations, toCappingAnnotation(resourceName, cappedToLimit))
		}
	}
	return annotations
}
func applyVPAPolicy(recommendation apiv1.ResourceList, policy *vpa_types.ContainerResourcePolicy) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if policy == nil {
		return nil
	}
	annotations := make([]string, 0)
	for resourceName, recommended := range recommendation {
		cappedToMin, isCapped := maybeCapToMin(recommended, resourceName, policy)
		recommendation[resourceName] = cappedToMin
		if isCapped {
			annotations = append(annotations, toCappingAnnotation(resourceName, cappedToMinAllowed))
		}
		cappedToMax, isCapped := maybeCapToMax(cappedToMin, resourceName, policy)
		recommendation[resourceName] = cappedToMax
		if isCapped {
			annotations = append(annotations, toCappingAnnotation(resourceName, cappedToMaxAllowed))
		}
	}
	return annotations
}
func applyVPAPolicyForContainer(containerName string, containerRecommendation *vpa_types.RecommendedContainerResources, policy *vpa_types.PodResourcePolicy) (*vpa_types.RecommendedContainerResources, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if containerRecommendation == nil {
		return nil, fmt.Errorf("no recommendation available for container name %v", containerName)
	}
	cappedRecommendations := containerRecommendation.DeepCopy()
	containerPolicy := GetContainerResourcePolicy(containerName, policy)
	if containerPolicy == nil {
		return cappedRecommendations, nil
	}
	process := func(recommendation apiv1.ResourceList) {
		for resourceName, recommended := range recommendation {
			cappedToMin, _ := maybeCapToMin(recommended, resourceName, containerPolicy)
			recommendation[resourceName] = cappedToMin
			cappedToMax, _ := maybeCapToMax(cappedToMin, resourceName, containerPolicy)
			recommendation[resourceName] = cappedToMax
		}
	}
	process(cappedRecommendations.Target)
	process(cappedRecommendations.LowerBound)
	process(cappedRecommendations.UpperBound)
	return cappedRecommendations, nil
}
func maybeCapToMin(recommended resource.Quantity, resourceName apiv1.ResourceName, containerPolicy *vpa_types.ContainerResourcePolicy) (resource.Quantity, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	min, found := containerPolicy.MinAllowed[resourceName]
	if found && !min.IsZero() && recommended.Cmp(min) < 0 {
		return min, true
	}
	return recommended, false
}
func maybeCapToMax(recommended resource.Quantity, resourceName apiv1.ResourceName, containerPolicy *vpa_types.ContainerResourcePolicy) (resource.Quantity, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	max, found := containerPolicy.MaxAllowed[resourceName]
	if found && !max.IsZero() && recommended.Cmp(max) > 0 {
		return max, true
	}
	return recommended, false
}
func ApplyVPAPolicy(podRecommendation *vpa_types.RecommendedPodResources, policy *vpa_types.PodResourcePolicy) (*vpa_types.RecommendedPodResources, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if podRecommendation == nil {
		return nil, nil
	}
	if policy == nil {
		return podRecommendation, nil
	}
	updatedRecommendations := []vpa_types.RecommendedContainerResources{}
	for _, containerRecommendation := range podRecommendation.ContainerRecommendations {
		containerName := containerRecommendation.ContainerName
		updatedContainerResources, err := applyVPAPolicyForContainer(containerName, &containerRecommendation, policy)
		if err != nil {
			return nil, fmt.Errorf("cannot apply policy on recommendation for container name %v", containerName)
		}
		updatedRecommendations = append(updatedRecommendations, *updatedContainerResources)
	}
	return &vpa_types.RecommendedPodResources{ContainerRecommendations: updatedRecommendations}, nil
}
func GetRecommendationForContainer(containerName string, recommendation *vpa_types.RecommendedPodResources) *vpa_types.RecommendedContainerResources {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if recommendation != nil {
		for i, containerRec := range recommendation.ContainerRecommendations {
			if containerRec.ContainerName == containerName {
				recommendationCopy := recommendation.ContainerRecommendations[i]
				return &recommendationCopy
			}
		}
	}
	return nil
}
func getContainer(containerName string, pod *apiv1.Pod) *apiv1.Container {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i, container := range pod.Spec.Containers {
		if container.Name == containerName {
			return &pod.Spec.Containers[i]
		}
	}
	return nil
}
