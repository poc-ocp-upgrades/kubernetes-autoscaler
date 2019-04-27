package test

import (
	"fmt"
	"time"
	"github.com/stretchr/testify/mock"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	vpa_lister "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/listers/autoscaling.k8s.io/v1beta1"
	v1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/record"
)

var (
	timeLayout		= "2006-01-02 15:04:05"
	testTimestamp, _	= time.Parse(timeLayout, "2017-04-18 17:35:05")
)

func BuildTestContainer(containerName, cpu, mem string) apiv1.Container {
	_logClusterCodePath()
	defer _logClusterCodePath()
	container := apiv1.Container{Name: containerName, Resources: apiv1.ResourceRequirements{Requests: apiv1.ResourceList{}}}
	if len(cpu) > 0 {
		cpuVal, _ := resource.ParseQuantity(cpu)
		container.Resources.Requests[apiv1.ResourceCPU] = cpuVal
	}
	if len(mem) > 0 {
		memVal, _ := resource.ParseQuantity(mem)
		container.Resources.Requests[apiv1.ResourceMemory] = memVal
	}
	return container
}
func BuildTestPolicy(containerName, minCPU, maxCPU, minMemory, maxMemory string) *vpa_types.PodResourcePolicy {
	_logClusterCodePath()
	defer _logClusterCodePath()
	minCPUVal, _ := resource.ParseQuantity(minCPU)
	maxCPUVal, _ := resource.ParseQuantity(maxCPU)
	minMemVal, _ := resource.ParseQuantity(minMemory)
	maxMemVal, _ := resource.ParseQuantity(maxMemory)
	return &vpa_types.PodResourcePolicy{ContainerPolicies: []vpa_types.ContainerResourcePolicy{{ContainerName: containerName, MinAllowed: apiv1.ResourceList{apiv1.ResourceMemory: minMemVal, apiv1.ResourceCPU: minCPUVal}, MaxAllowed: apiv1.ResourceList{apiv1.ResourceMemory: maxMemVal, apiv1.ResourceCPU: maxCPUVal}}}}
}
func Resources(cpu, mem string) apiv1.ResourceList {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(apiv1.ResourceList)
	if len(cpu) > 0 {
		cpuVal, _ := resource.ParseQuantity(cpu)
		result[apiv1.ResourceCPU] = cpuVal
	}
	if len(mem) > 0 {
		memVal, _ := resource.ParseQuantity(mem)
		result[apiv1.ResourceMemory] = memVal
	}
	return result
}

type RecommenderAPIMock struct{ mock.Mock }

func (m *RecommenderAPIMock) GetRecommendation(spec *apiv1.PodSpec) (*vpa_types.RecommendedPodResources, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := m.Called(spec)
	var returnArg *vpa_types.RecommendedPodResources
	if args.Get(0) != nil {
		returnArg = args.Get(0).(*vpa_types.RecommendedPodResources)
	}
	return returnArg, args.Error(1)
}

type RecommenderMock struct{ mock.Mock }

func (m *RecommenderMock) Get(spec *apiv1.PodSpec) (*vpa_types.RecommendedPodResources, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := m.Called(spec)
	var returnArg *vpa_types.RecommendedPodResources
	if args.Get(0) != nil {
		returnArg = args.Get(0).(*vpa_types.RecommendedPodResources)
	}
	return returnArg, args.Error(1)
}

type PodsEvictionRestrictionMock struct{ mock.Mock }

func (m *PodsEvictionRestrictionMock) Evict(pod *apiv1.Pod, eventRecorder record.EventRecorder) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := m.Called(pod, eventRecorder)
	return args.Error(0)
}
func (m *PodsEvictionRestrictionMock) CanEvict(pod *apiv1.Pod) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := m.Called(pod)
	return args.Bool(0)
}

type PodListerMock struct{ mock.Mock }

func (m *PodListerMock) Pods(namespace string) v1.PodNamespaceLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := m.Called(namespace)
	var returnArg v1.PodNamespaceLister
	if args.Get(0) != nil {
		returnArg = args.Get(0).(v1.PodNamespaceLister)
	}
	return returnArg
}
func (m *PodListerMock) List(selector labels.Selector) (ret []*apiv1.Pod, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := m.Called()
	var returnArg []*apiv1.Pod
	if args.Get(0) != nil {
		returnArg = args.Get(0).([]*apiv1.Pod)
	}
	return returnArg, args.Error(1)
}
func (m *PodListerMock) Get(name string) (*apiv1.Pod, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, fmt.Errorf("unimplemented")
}

type VerticalPodAutoscalerListerMock struct{ mock.Mock }

func (m *VerticalPodAutoscalerListerMock) List(selector labels.Selector) (ret []*vpa_types.VerticalPodAutoscaler, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := m.Called()
	var returnArg []*vpa_types.VerticalPodAutoscaler
	if args.Get(0) != nil {
		returnArg = args.Get(0).([]*vpa_types.VerticalPodAutoscaler)
	}
	return returnArg, args.Error(1)
}
func (m *VerticalPodAutoscalerListerMock) VerticalPodAutoscalers(namespace string) vpa_lister.VerticalPodAutoscalerNamespaceLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := m.Called(namespace)
	var returnArg vpa_lister.VerticalPodAutoscalerNamespaceLister
	if args.Get(0) != nil {
		returnArg = args.Get(0).(vpa_lister.VerticalPodAutoscalerNamespaceLister)
	}
	return returnArg
}
func (m *VerticalPodAutoscalerListerMock) Get(name string) (*vpa_types.VerticalPodAutoscaler, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, fmt.Errorf("unimplemented")
}

type RecommendationProcessorMock struct{ mock.Mock }

func (m *RecommendationProcessorMock) Apply(podRecommendation *vpa_types.RecommendedPodResources, policy *vpa_types.PodResourcePolicy, conditions []vpa_types.VerticalPodAutoscalerCondition, pod *apiv1.Pod) (*vpa_types.RecommendedPodResources, map[string][]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := m.Called()
	var returnArg *vpa_types.RecommendedPodResources
	if args.Get(0) != nil {
		returnArg = args.Get(0).(*vpa_types.RecommendedPodResources)
	}
	var annotations map[string][]string
	if args.Get(1) != nil {
		annotations = args.Get(1).(map[string][]string)
	}
	return returnArg, annotations, args.Error(1)
}

type FakeRecommendationProcessor struct{}

func (f *FakeRecommendationProcessor) Apply(podRecommendation *vpa_types.RecommendedPodResources, policy *vpa_types.PodResourcePolicy, conditions []vpa_types.VerticalPodAutoscalerCondition, pod *apiv1.Pod) (*vpa_types.RecommendedPodResources, map[string][]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return podRecommendation, nil, nil
}

type fakeEventRecorder struct{}

func (f *fakeEventRecorder) Event(object runtime.Object, eventtype, reason, message string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (f *fakeEventRecorder) Eventf(object runtime.Object, eventtype, reason, messageFmt string, args ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (f *fakeEventRecorder) PastEventf(object runtime.Object, timestamp metav1.Time, eventtype, reason, messageFmt string, args ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (f *fakeEventRecorder) AnnotatedEventf(object runtime.Object, annotations map[string]string, eventtype, reason, messageFmt string, args ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func FakeEventRecorder() record.EventRecorder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &fakeEventRecorder{}
}
