package logic

import (
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	vpa_lister "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/listers/autoscaling.k8s.io/v1beta1"
	vpa_api_util "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/utils/vpa"
)

type ContainerResources struct{ Requests v1.ResourceList }

func newContainerResources() ContainerResources {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ContainerResources{Requests: v1.ResourceList{}}
}

type RecommendationProvider interface {
	GetContainersResourcesForPod(pod *v1.Pod) ([]ContainerResources, vpa_api_util.ContainerToAnnotationsMap, string, error)
}
type recommendationProvider struct {
	vpaLister		vpa_lister.VerticalPodAutoscalerLister
	recommendationProcessor	vpa_api_util.RecommendationProcessor
}

func NewRecommendationProvider(vpaLister vpa_lister.VerticalPodAutoscalerLister, recommendationProcessor vpa_api_util.RecommendationProcessor) *recommendationProvider {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &recommendationProvider{vpaLister: vpaLister, recommendationProcessor: recommendationProcessor}
}
func getContainersResources(pod *v1.Pod, podRecommendation vpa_types.RecommendedPodResources) []ContainerResources {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	resources := make([]ContainerResources, len(pod.Spec.Containers))
	for i, container := range pod.Spec.Containers {
		resources[i] = newContainerResources()
		recommendation := vpa_api_util.GetRecommendationForContainer(container.Name, &podRecommendation)
		if recommendation == nil {
			glog.V(2).Infof("no matching recommendation found for container %s", container.Name)
			continue
		}
		resources[i].Requests = recommendation.Target
	}
	return resources
}
func (p *recommendationProvider) getMatchingVPA(pod *v1.Pod) *vpa_types.VerticalPodAutoscaler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	configs, err := p.vpaLister.VerticalPodAutoscalers(pod.Namespace).List(labels.Everything())
	if err != nil {
		glog.Errorf("failed to get vpa configs: %v", err)
		return nil
	}
	onConfigs := make([]*vpa_types.VerticalPodAutoscaler, 0)
	for _, vpaConfig := range configs {
		if vpa_api_util.GetUpdateMode(vpaConfig) == vpa_types.UpdateModeOff {
			continue
		}
		onConfigs = append(onConfigs, vpaConfig)
	}
	glog.V(2).Infof("Let's choose from %d configs for pod %s/%s", len(onConfigs), pod.Namespace, pod.Name)
	return vpa_api_util.GetControllingVPAForPod(pod, onConfigs)
}
func (p *recommendationProvider) GetContainersResourcesForPod(pod *v1.Pod) ([]ContainerResources, vpa_api_util.ContainerToAnnotationsMap, string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	glog.V(2).Infof("updating requirements for pod %s.", pod.Name)
	vpaConfig := p.getMatchingVPA(pod)
	if vpaConfig == nil {
		glog.V(2).Infof("no matching VPA found for pod %s", pod.Name)
		return nil, nil, "", nil
	}
	var annotations vpa_api_util.ContainerToAnnotationsMap
	recommendedPodResources := &vpa_types.RecommendedPodResources{}
	if vpaConfig.Status.Recommendation != nil {
		var err error
		recommendedPodResources, annotations, err = p.recommendationProcessor.Apply(vpaConfig.Status.Recommendation, vpaConfig.Spec.ResourcePolicy, vpaConfig.Status.Conditions, pod)
		if err != nil {
			glog.V(2).Infof("cannot process recommendation for pod %s", pod.Name)
			return nil, annotations, vpaConfig.Name, err
		}
	}
	containerResources := getContainersResources(pod, *recommendedPodResources)
	return containerResources, annotations, vpaConfig.Name, nil
}
