package logic

import (
	"testing"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/utils/test"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/utils/vpa"
	vpa_api_util "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/utils/vpa"
)

func TestUpdateResourceRequests(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type testCase struct {
		pod		*apiv1.Pod
		vpas		[]*vpa_types.VerticalPodAutoscaler
		expectedAction	bool
		expectedMem	string
		expectedCPU	string
		annotations	vpa_api_util.ContainerToAnnotationsMap
	}
	containerName := "container1"
	vpaName := "vpa1"
	labels := map[string]string{"app": "testingApp"}
	vpaBuilder := test.VerticalPodAutoscaler().WithName(vpaName).WithContainer(containerName).WithTarget("2", "200Mi").WithMinAllowed("1", "100Mi").WithMaxAllowed("3", "1Gi").WithSelector("app = testingApp")
	vpa := vpaBuilder.Get()
	uninitialized := test.Pod().WithName("test_uninitialized").AddContainer(test.BuildTestContainer(containerName, "", "")).Get()
	uninitialized.ObjectMeta.Labels = labels
	initialized := test.Pod().WithName("test_initialized").AddContainer(test.BuildTestContainer(containerName, "1", "100Mi")).Get()
	initialized.ObjectMeta.Labels = labels
	mismatchedVPA := vpaBuilder.WithSelector("app = differentApp").Get()
	offVPA := vpaBuilder.WithUpdateMode(vpa_types.UpdateModeOff).Get()
	targetBelowMinVPA := vpaBuilder.WithTarget("3", "150Mi").WithMinAllowed("4", "300Mi").WithMaxAllowed("5", "1Gi").Get()
	targetAboveMaxVPA := vpaBuilder.WithTarget("7", "2Gi").WithMinAllowed("4", "300Mi").WithMaxAllowed("5", "1Gi").Get()
	vpaWithHighMemory := vpaBuilder.WithTarget("2", "1000Mi").WithMaxAllowed("3", "3Gi").Get()
	vpaWithEmptyRecommendation := vpaBuilder.Get()
	vpaWithEmptyRecommendation.Status.Recommendation = &vpa_types.RecommendedPodResources{}
	vpaWithNilRecommendation := vpaBuilder.Get()
	vpaWithNilRecommendation.Status.Recommendation = nil
	testCases := []testCase{{pod: uninitialized, vpas: []*vpa_types.VerticalPodAutoscaler{vpa}, expectedAction: true, expectedMem: "200Mi", expectedCPU: "2"}, {pod: uninitialized, vpas: []*vpa_types.VerticalPodAutoscaler{vpa}, expectedAction: true, expectedMem: "200Mi", expectedCPU: "2"}, {pod: uninitialized, vpas: []*vpa_types.VerticalPodAutoscaler{targetBelowMinVPA}, expectedAction: true, expectedMem: "300Mi", expectedCPU: "4", annotations: vpa_api_util.ContainerToAnnotationsMap{containerName: []string{"cpu capped to minAllowed", "memory capped to minAllowed"}}}, {pod: uninitialized, vpas: []*vpa_types.VerticalPodAutoscaler{targetAboveMaxVPA}, expectedAction: true, expectedMem: "1Gi", expectedCPU: "5", annotations: vpa_api_util.ContainerToAnnotationsMap{containerName: []string{"cpu capped to maxAllowed", "memory capped to maxAllowed"}}}, {pod: initialized, vpas: []*vpa_types.VerticalPodAutoscaler{vpa}, expectedAction: true, expectedMem: "200Mi", expectedCPU: "2"}, {pod: initialized, vpas: []*vpa_types.VerticalPodAutoscaler{vpaWithHighMemory}, expectedAction: true, expectedMem: "1000Mi", expectedCPU: "2"}, {pod: uninitialized, vpas: []*vpa_types.VerticalPodAutoscaler{mismatchedVPA}, expectedAction: false}, {pod: uninitialized, vpas: []*vpa_types.VerticalPodAutoscaler{offVPA}, expectedAction: false}, {pod: uninitialized, vpas: []*vpa_types.VerticalPodAutoscaler{offVPA, vpa}, expectedAction: true, expectedMem: "200Mi", expectedCPU: "2"}, {pod: initialized, vpas: []*vpa_types.VerticalPodAutoscaler{vpaWithEmptyRecommendation}, expectedAction: true, expectedMem: "0", expectedCPU: "0"}, {pod: initialized, vpas: []*vpa_types.VerticalPodAutoscaler{vpaWithNilRecommendation}, expectedAction: true, expectedMem: "0", expectedCPU: "0"}}
	for _, tc := range testCases {
		vpaNamespaceLister := &test.VerticalPodAutoscalerListerMock{}
		vpaNamespaceLister.On("List").Return(tc.vpas, nil)
		vpaLister := &test.VerticalPodAutoscalerListerMock{}
		vpaLister.On("VerticalPodAutoscalers", "default").Return(vpaNamespaceLister)
		recommendationProvider := &recommendationProvider{vpaLister: vpaLister, recommendationProcessor: api.NewCappingRecommendationProcessor()}
		resources, annotations, name, err := recommendationProvider.GetContainersResourcesForPod(tc.pod)
		if tc.expectedAction {
			assert.Equal(t, vpaName, name)
			assert.Nil(t, err)
			assert.Equal(t, len(resources), 1)
			expectedCPU, err := resource.ParseQuantity(tc.expectedCPU)
			assert.NoError(t, err)
			cpuRequest := resources[0].Requests[apiv1.ResourceCPU]
			assert.Equal(t, expectedCPU.Value(), cpuRequest.Value(), "cpu request doesn't match")
			expectedMemory, err := resource.ParseQuantity(tc.expectedMem)
			assert.NoError(t, err)
			memoryRequest := resources[0].Requests[apiv1.ResourceMemory]
			assert.Equal(t, expectedMemory.Value(), memoryRequest.Value(), "memory request doesn't match")
			assert.Len(t, annotations, len(tc.annotations))
			if len(tc.annotations) > 0 {
				for annotationKey, annotationValues := range tc.annotations {
					assert.Len(t, annotations[annotationKey], len(annotationValues))
					for _, annotation := range annotationValues {
						assert.Contains(t, annotations[annotationKey], annotation)
					}
				}
			}
		} else {
			assert.Equal(t, len(resources), 0)
		}
	}
}
