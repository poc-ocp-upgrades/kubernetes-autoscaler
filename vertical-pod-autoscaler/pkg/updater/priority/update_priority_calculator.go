package priority

import (
	"flag"
	"math"
	"sort"
	"time"
	"github.com/golang/glog"
	apiv1 "k8s.io/api/core/v1"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	vpa_api_util "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/utils/vpa"
)

const (
	defaultUpdateThreshold		= 0.10
	podLifetimeUpdateThreshold	= time.Hour * 12
)

var (
	evictAfterOOMThreshold = flag.Duration("evict-after-oom-treshold", 10*time.Minute, `Evict pod that has only one container and it OOMed in less than
		evict-after-oom-treshold since start.`)
)

type UpdatePriorityCalculator struct {
	resourcesPolicy			*vpa_types.PodResourcePolicy
	conditions				[]vpa_types.VerticalPodAutoscalerCondition
	pods					[]podPriority
	config					*UpdateConfig
	recommendationProcessor	vpa_api_util.RecommendationProcessor
}
type UpdateConfig struct{ MinChangePriority float64 }

func NewUpdatePriorityCalculator(policy *vpa_types.PodResourcePolicy, conditions []vpa_types.VerticalPodAutoscalerCondition, config *UpdateConfig, processor vpa_api_util.RecommendationProcessor) UpdatePriorityCalculator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if config == nil {
		config = &UpdateConfig{MinChangePriority: defaultUpdateThreshold}
	}
	return UpdatePriorityCalculator{resourcesPolicy: policy, conditions: conditions, config: config, recommendationProcessor: processor}
}
func (calc *UpdatePriorityCalculator) AddPod(pod *apiv1.Pod, recommendation *vpa_types.RecommendedPodResources, now time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	processedRecommendation, _, err := calc.recommendationProcessor.Apply(recommendation, calc.resourcesPolicy, calc.conditions, pod)
	if err != nil {
		glog.V(2).Infof("cannot process recommendation for pod %s: %v", pod.Name, err)
		return
	}
	updatePriority := calc.getUpdatePriority(pod, processedRecommendation)
	quickOOM := false
	if len(pod.Status.ContainerStatuses) == 1 {
		terminationState := pod.Status.ContainerStatuses[0].LastTerminationState
		if terminationState.Terminated != nil && terminationState.Terminated.Reason == "OOMKilled" && terminationState.Terminated.FinishedAt.Time.Sub(terminationState.Terminated.StartedAt.Time) < *evictAfterOOMThreshold {
			quickOOM = true
			glog.V(2).Infof("quick OOM detected in pod %v", pod.Name)
		}
	}
	if !updatePriority.outsideRecommendedRange && !quickOOM {
		if pod.Status.StartTime == nil {
			glog.V(2).Infof("not updating pod %v, missing field pod.Status.StartTime", pod.Name)
			return
		}
		if now.Before(pod.Status.StartTime.Add(podLifetimeUpdateThreshold)) {
			glog.V(2).Infof("not updating a short-lived pod %v, request within recommended range", pod.Name)
			return
		}
		if updatePriority.resourceDiff < calc.config.MinChangePriority {
			glog.V(2).Infof("not updating pod %v, resource diff too low: %v", pod.Name, updatePriority)
			return
		}
	}
	glog.V(2).Infof("pod accepted for update %v with priority %v", pod.Name, updatePriority.resourceDiff)
	calc.pods = append(calc.pods, updatePriority)
}
func (calc *UpdatePriorityCalculator) GetSortedPods(admission PodEvictionAdmission) []*apiv1.Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sort.Sort(byPriority(calc.pods))
	result := []*apiv1.Pod{}
	for _, podPrio := range calc.pods {
		if admission == nil || admission.Admit(podPrio.pod, podPrio.recommendation) {
			result = append(result, podPrio.pod)
		} else {
			glog.V(2).Infof("pod removed from update queue by PodEvictionAdmission: %v", podPrio.pod.Name)
		}
	}
	return result
}
func (calc *UpdatePriorityCalculator) getUpdatePriority(pod *apiv1.Pod, recommendation *vpa_types.RecommendedPodResources) podPriority {
	_logClusterCodePath()
	defer _logClusterCodePath()
	outsideRecommendedRange := false
	scaleUp := false
	totalRequestPerResource := make(map[apiv1.ResourceName]int64)
	totalRecommendedPerResource := make(map[apiv1.ResourceName]int64)
	for _, podContainer := range pod.Spec.Containers {
		recommendedRequest := vpa_api_util.GetRecommendationForContainer(podContainer.Name, recommendation)
		if recommendedRequest == nil {
			continue
		}
		for resourceName, recommended := range recommendedRequest.Target {
			totalRecommendedPerResource[resourceName] += recommended.MilliValue()
			lowerBound, hasLowerBound := recommendedRequest.LowerBound[resourceName]
			upperBound, hasUpperBound := recommendedRequest.UpperBound[resourceName]
			if request, hasRequest := podContainer.Resources.Requests[resourceName]; hasRequest {
				totalRequestPerResource[resourceName] += request.MilliValue()
				if recommended.MilliValue() > request.MilliValue() {
					scaleUp = true
				}
				if (hasLowerBound && request.Cmp(lowerBound) < 0) || (hasUpperBound && request.Cmp(upperBound) > 0) {
					outsideRecommendedRange = true
				}
			} else {
				scaleUp = true
				outsideRecommendedRange = true
			}
		}
	}
	resourceDiff := 0.0
	for resource, totalRecommended := range totalRecommendedPerResource {
		totalRequest := math.Max(float64(totalRequestPerResource[resource]), 1.0)
		resourceDiff += math.Abs(totalRequest-float64(totalRecommended)) / totalRequest
	}
	return podPriority{pod: pod, outsideRecommendedRange: outsideRecommendedRange, scaleUp: scaleUp, resourceDiff: resourceDiff, recommendation: recommendation}
}

type podPriority struct {
	pod						*apiv1.Pod
	outsideRecommendedRange	bool
	scaleUp					bool
	resourceDiff			float64
	recommendation			*vpa_types.RecommendedPodResources
}
type byPriority []podPriority

func (list byPriority) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(list)
}
func (list byPriority) Swap(i, j int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	list[i], list[j] = list[j], list[i]
}
func (list byPriority) Less(i, j int) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if list[i].scaleUp != list[j].scaleUp {
		return list[i].scaleUp
	}
	return list[i].resourceDiff > list[j].resourceDiff
}
