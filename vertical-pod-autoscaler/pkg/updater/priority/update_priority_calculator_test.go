package priority

import (
	"testing"
	"time"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/utils/test"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/stretchr/testify/assert"
)

const (
	containerName = "container1"
)

func TestSortPriority(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	calculator := NewUpdatePriorityCalculator(nil, nil, nil, &test.FakeRecommendationProcessor{})
	pod1 := test.Pod().WithName("POD1").AddContainer(test.BuildTestContainer(containerName, "2", "")).Get()
	pod2 := test.Pod().WithName("POD2").AddContainer(test.BuildTestContainer(containerName, "4", "")).Get()
	pod3 := test.Pod().WithName("POD3").AddContainer(test.BuildTestContainer(containerName, "1", "")).Get()
	pod4 := test.Pod().WithName("POD4").AddContainer(test.BuildTestContainer(containerName, "3", "")).Get()
	recommendation := test.Recommendation().WithContainer(containerName).WithTarget("10", "").Get()
	timestampNow := pod1.Status.StartTime.Time.Add(time.Hour * 24)
	calculator.AddPod(pod1, recommendation, timestampNow)
	calculator.AddPod(pod2, recommendation, timestampNow)
	calculator.AddPod(pod3, recommendation, timestampNow)
	calculator.AddPod(pod4, recommendation, timestampNow)
	result := calculator.GetSortedPods(NewDefaultPodEvictionAdmission())
	assert.Exactly(t, []*apiv1.Pod{pod3, pod1, pod4, pod2}, result, "Wrong priority order")
}
func TestSortPriorityMultiResource(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	calculator := NewUpdatePriorityCalculator(nil, nil, nil, &test.FakeRecommendationProcessor{})
	pod1 := test.Pod().WithName("POD1").AddContainer(test.BuildTestContainer(containerName, "4", "60M")).Get()
	pod2 := test.Pod().WithName("POD2").AddContainer(test.BuildTestContainer(containerName, "3", "90M")).Get()
	recommendation := test.Recommendation().WithContainer(containerName).WithTarget("6", "100M").Get()
	timestampNow := pod1.Status.StartTime.Time.Add(time.Hour * 24)
	calculator.AddPod(pod1, recommendation, timestampNow)
	calculator.AddPod(pod2, recommendation, timestampNow)
	result := calculator.GetSortedPods(NewDefaultPodEvictionAdmission())
	assert.Exactly(t, []*apiv1.Pod{pod1, pod2}, result, "Wrong priority order")
}
func TestSortPriorityMultiContainers(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	containerName2 := "container2"
	pod1 := test.Pod().WithName("POD1").AddContainer(test.BuildTestContainer(containerName, "3", "10M")).Get()
	pod2 := test.Pod().WithName("POD2").AddContainer(test.BuildTestContainer(containerName, "4", "10M")).Get()
	container2 := test.BuildTestContainer(containerName2, "2", "20M")
	pod2.Spec.Containers = append(pod2.Spec.Containers, container2)
	recommendation := test.Recommendation().WithContainer(containerName).WithTarget("6", "20M").Get()
	cpuRec, _ := resource.ParseQuantity("4")
	memRec, _ := resource.ParseQuantity("20M")
	container2rec := vpa_types.RecommendedContainerResources{ContainerName: containerName2, Target: map[apiv1.ResourceName]resource.Quantity{apiv1.ResourceCPU: cpuRec, apiv1.ResourceMemory: memRec}}
	recommendation.ContainerRecommendations = append(recommendation.ContainerRecommendations, container2rec)
	timestampNow := pod1.Status.StartTime.Time.Add(time.Hour * 24)
	calculator := NewUpdatePriorityCalculator(nil, nil, nil, &test.FakeRecommendationProcessor{})
	calculator.AddPod(pod1, recommendation, timestampNow)
	calculator.AddPod(pod2, recommendation, timestampNow)
	podPriority1 := calculator.getUpdatePriority(pod1, recommendation)
	assert.Equal(t, 2.0, podPriority1.resourceDiff)
	podPriority2 := calculator.getUpdatePriority(pod2, recommendation)
	assert.Equal(t, 1.0, podPriority2.resourceDiff)
	result := calculator.GetSortedPods(NewDefaultPodEvictionAdmission())
	assert.Exactly(t, []*apiv1.Pod{pod1, pod2}, result, "Wrong priority order")
}
func TestSortPriorityResourcesDecrease(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	calculator := NewUpdatePriorityCalculator(nil, nil, nil, &test.FakeRecommendationProcessor{})
	pod1 := test.Pod().WithName("POD1").AddContainer(test.BuildTestContainer(containerName, "4", "")).Get()
	pod2 := test.Pod().WithName("POD2").AddContainer(test.BuildTestContainer(containerName, "7", "")).Get()
	pod3 := test.Pod().WithName("POD3").AddContainer(test.BuildTestContainer(containerName, "10", "")).Get()
	recommendation := test.Recommendation().WithContainer(containerName).WithTarget("5", "").Get()
	timestampNow := pod1.Status.StartTime.Time.Add(time.Hour * 24)
	calculator.AddPod(pod1, recommendation, timestampNow)
	calculator.AddPod(pod2, recommendation, timestampNow)
	calculator.AddPod(pod3, recommendation, timestampNow)
	result := calculator.GetSortedPods(NewDefaultPodEvictionAdmission())
	assert.Exactly(t, []*apiv1.Pod{pod1, pod3, pod2}, result, "Wrong priority order")
}
func TestUpdateNotRequired(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	calculator := NewUpdatePriorityCalculator(nil, nil, nil, &test.FakeRecommendationProcessor{})
	pod1 := test.Pod().WithName("POD1").AddContainer(test.BuildTestContainer(containerName, "4", "")).Get()
	recommendation := test.Recommendation().WithContainer(containerName).WithTarget("4", "").Get()
	timestampNow := pod1.Status.StartTime.Time.Add(time.Hour * 24)
	calculator.AddPod(pod1, recommendation, timestampNow)
	result := calculator.GetSortedPods(NewDefaultPodEvictionAdmission())
	assert.Exactly(t, []*apiv1.Pod{}, result, "Pod should not be updated")
}
func TestUpdateRequiredOnMilliQuantities(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	calculator := NewUpdatePriorityCalculator(nil, nil, nil, &test.FakeRecommendationProcessor{})
	pod1 := test.Pod().WithName("POD1").AddContainer(test.BuildTestContainer(containerName, "10m", "")).Get()
	recommendation := test.Recommendation().WithContainer(containerName).WithTarget("900m", "").Get()
	timestampNow := pod1.Status.StartTime.Time.Add(time.Hour * 24)
	calculator.AddPod(pod1, recommendation, timestampNow)
	result := calculator.GetSortedPods(NewDefaultPodEvictionAdmission())
	assert.Exactly(t, []*apiv1.Pod{pod1}, result, "Pod should be updated")
}
func TestUseProcessor(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	processedRecommendation := test.Recommendation().WithContainer(containerName).WithTarget("4", "10M").Get()
	recommendationProcessor := &test.RecommendationProcessorMock{}
	recommendationProcessor.On("Apply").Return(processedRecommendation, nil)
	calculator := NewUpdatePriorityCalculator(nil, nil, nil, recommendationProcessor)
	pod1 := test.Pod().WithName("POD1").AddContainer(test.BuildTestContainer(containerName, "4", "10M")).Get()
	recommendation := test.Recommendation().WithContainer(containerName).WithTarget("5", "5M").Get()
	timestampNow := pod1.Status.StartTime.Time.Add(time.Hour * 24)
	calculator.AddPod(pod1, recommendation, timestampNow)
	result := calculator.GetSortedPods(NewDefaultPodEvictionAdmission())
	assert.Exactly(t, []*apiv1.Pod{}, result, "Pod should not be updated")
}
func TestUpdateLonglivedPods(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	calculator := NewUpdatePriorityCalculator(nil, nil, &UpdateConfig{MinChangePriority: 0.5}, &test.FakeRecommendationProcessor{})
	pods := []*apiv1.Pod{test.Pod().WithName("POD1").AddContainer(test.BuildTestContainer(containerName, "4", "")).Get(), test.Pod().WithName("POD2").AddContainer(test.BuildTestContainer(containerName, "1", "")).Get(), test.Pod().WithName("POD3").AddContainer(test.BuildTestContainer(containerName, "7", "")).Get()}
	recommendation := test.Recommendation().WithContainer(containerName).WithTarget("5", "").WithLowerBound("1", "").WithUpperBound("6", "").Get()
	timestampNow := pods[0].Status.StartTime.Time.Add(time.Hour * 13)
	for i := 0; i < 3; i++ {
		calculator.AddPod(pods[i], recommendation, timestampNow)
	}
	result := calculator.GetSortedPods(NewDefaultPodEvictionAdmission())
	assert.Exactly(t, []*apiv1.Pod{pods[1], pods[2]}, result, "Exactly POD2 and POD3 should be updated")
}
func TestUpdateShortlivedPods(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	calculator := NewUpdatePriorityCalculator(nil, nil, &UpdateConfig{MinChangePriority: 0.5}, &test.FakeRecommendationProcessor{})
	pods := []*apiv1.Pod{test.Pod().WithName("POD1").AddContainer(test.BuildTestContainer(containerName, "4", "")).Get(), test.Pod().WithName("POD2").AddContainer(test.BuildTestContainer(containerName, "1", "")).Get(), test.Pod().WithName("POD3").AddContainer(test.BuildTestContainer(containerName, "7", "")).Get()}
	recommendation := test.Recommendation().WithContainer(containerName).WithTarget("5", "").WithLowerBound("1", "").WithUpperBound("6", "").Get()
	timestampNow := pods[0].Status.StartTime.Time.Add(time.Hour * 11)
	for i := 0; i < 3; i++ {
		calculator.AddPod(pods[i], recommendation, timestampNow)
	}
	result := calculator.GetSortedPods(NewDefaultPodEvictionAdmission())
	assert.Exactly(t, []*apiv1.Pod{pods[2]}, result, "Only POD3 should be updated")
}
func TestUpdatePodWithQuickOOM(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	calculator := NewUpdatePriorityCalculator(nil, nil, &UpdateConfig{MinChangePriority: 0.5}, &test.FakeRecommendationProcessor{})
	pod := test.Pod().WithName("POD1").AddContainer(test.BuildTestContainer(containerName, "4", "")).Get()
	timestampNow := pod.Status.StartTime.Time.Add(time.Hour * 11)
	pod.Status.ContainerStatuses = []apiv1.ContainerStatus{{LastTerminationState: apiv1.ContainerState{Terminated: &apiv1.ContainerStateTerminated{Reason: "OOMKilled", FinishedAt: metav1.NewTime(timestampNow.Add(-1 * 3 * time.Minute)), StartedAt: metav1.NewTime(timestampNow.Add(-1 * 5 * time.Minute))}}}}
	recommendation := test.Recommendation().WithContainer(containerName).WithTarget("5", "").WithLowerBound("1", "").WithUpperBound("6", "").Get()
	calculator.AddPod(pod, recommendation, timestampNow)
	result := calculator.GetSortedPods(NewDefaultPodEvictionAdmission())
	assert.Exactly(t, []*apiv1.Pod{pod}, result, "Pod should be updated")
}
func TestDontUpdatePodWithOOMAfterLongRun(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	calculator := NewUpdatePriorityCalculator(nil, nil, &UpdateConfig{MinChangePriority: 0.5}, &test.FakeRecommendationProcessor{})
	pod := test.Pod().WithName("POD1").AddContainer(test.BuildTestContainer(containerName, "4", "")).Get()
	timestampNow := pod.Status.StartTime.Time.Add(time.Hour * 11)
	pod.Status.ContainerStatuses = []apiv1.ContainerStatus{{LastTerminationState: apiv1.ContainerState{Terminated: &apiv1.ContainerStateTerminated{Reason: "OOMKilled", FinishedAt: metav1.NewTime(timestampNow.Add(-1 * 3 * time.Minute)), StartedAt: metav1.NewTime(timestampNow.Add(-1 * 60 * time.Minute))}}}}
	recommendation := test.Recommendation().WithContainer(containerName).WithTarget("5", "").WithLowerBound("1", "").WithUpperBound("6", "").Get()
	calculator.AddPod(pod, recommendation, timestampNow)
	result := calculator.GetSortedPods(NewDefaultPodEvictionAdmission())
	assert.Exactly(t, []*apiv1.Pod{}, result, "Pod shouldn't be updated")
}
func TestDontUpdatePodWithOOMOnlyOnOneContainer(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	calculator := NewUpdatePriorityCalculator(nil, nil, &UpdateConfig{MinChangePriority: 0.5}, &test.FakeRecommendationProcessor{})
	pod := test.Pod().WithName("POD1").AddContainer(test.BuildTestContainer(containerName, "4", "")).Get()
	timestampNow := pod.Status.StartTime.Time.Add(time.Hour * 11)
	pod.Status.ContainerStatuses = []apiv1.ContainerStatus{{LastTerminationState: apiv1.ContainerState{Terminated: &apiv1.ContainerStateTerminated{Reason: "OOMKilled", FinishedAt: metav1.NewTime(timestampNow.Add(-1 * 3 * time.Minute)), StartedAt: metav1.NewTime(timestampNow.Add(-1 * 5 * time.Minute))}}}, {}}
	recommendation := test.Recommendation().WithContainer(containerName).WithTarget("5", "").WithLowerBound("1", "").WithUpperBound("6", "").Get()
	calculator.AddPod(pod, recommendation, timestampNow)
	result := calculator.GetSortedPods(NewDefaultPodEvictionAdmission())
	assert.Exactly(t, []*apiv1.Pod{}, result, "Pod shouldn't be updated")
}
func TestNoPods(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	calculator := NewUpdatePriorityCalculator(nil, nil, nil, &test.FakeRecommendationProcessor{})
	result := calculator.GetSortedPods(NewDefaultPodEvictionAdmission())
	assert.Exactly(t, []*apiv1.Pod{}, result)
}

type pod1Admission struct{}

func (p *pod1Admission) LoopInit([]*apiv1.Pod, map[*vpa_types.VerticalPodAutoscaler][]*apiv1.Pod) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (p *pod1Admission) Admit(pod *apiv1.Pod, recommendation *vpa_types.RecommendedPodResources) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return pod.Name == "POD1"
}
func (p *pod1Admission) CleanUp() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func TestAdmission(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	calculator := NewUpdatePriorityCalculator(nil, nil, nil, &test.FakeRecommendationProcessor{})
	pod1 := test.Pod().WithName("POD1").AddContainer(test.BuildTestContainer(containerName, "2", "")).Get()
	pod2 := test.Pod().WithName("POD2").AddContainer(test.BuildTestContainer(containerName, "4", "")).Get()
	pod3 := test.Pod().WithName("POD3").AddContainer(test.BuildTestContainer(containerName, "1", "")).Get()
	pod4 := test.Pod().WithName("POD4").AddContainer(test.BuildTestContainer(containerName, "3", "")).Get()
	recommendation := test.Recommendation().WithContainer(containerName).WithTarget("10", "").Get()
	timestampNow := pod1.Status.StartTime.Time.Add(time.Hour * 24)
	calculator.AddPod(pod1, recommendation, timestampNow)
	calculator.AddPod(pod2, recommendation, timestampNow)
	calculator.AddPod(pod3, recommendation, timestampNow)
	calculator.AddPod(pod4, recommendation, timestampNow)
	result := calculator.GetSortedPods(&pod1Admission{})
	assert.Exactly(t, []*apiv1.Pod{pod1}, result, "Wrong priority order")
}
func TestNoRecommendationForContainer(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	calculator := NewUpdatePriorityCalculator(nil, nil, nil, &test.FakeRecommendationProcessor{})
	pod := test.Pod().WithName("POD1").AddContainer(test.BuildTestContainer(containerName, "5", "10")).Get()
	result := calculator.getUpdatePriority(pod, nil)
	assert.NotNil(t, result)
}
