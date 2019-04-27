package tpu

import (
	"reflect"
	"testing"
	"k8s.io/apimachinery/pkg/api/resource"
	. "k8s.io/autoscaler/cluster-autoscaler/utils/test"
	apiv1 "k8s.io/api/core/v1"
)

var (
	ResourceTPUV2			= ResourceTPUPrefix + "v2"
	ResourceTPUPreemptibleV2	= ResourceTPUPrefix + "preemptible-v2"
)

type requests map[apiv1.ResourceName]int64
type containerSpecs []requests

func testContainer(requests requests) apiv1.Container {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	container := apiv1.Container{Resources: apiv1.ResourceRequirements{Requests: apiv1.ResourceList{}}}
	for res, request := range requests {
		container.Resources.Requests[res] = *resource.NewQuantity(request, resource.DecimalSI)
	}
	return container
}
func testPod(name string, containers ...requests) *apiv1.Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pod := BuildTestPod(name, 0, 0)
	for _, requests := range containers[1:] {
		pod.Spec.Containers = append(pod.Spec.Containers, testContainer(requests))
	}
	return pod
}
func TestClearTPURequests(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cpuPod := testPod("cpuPod", requests{apiv1.ResourceCPU: 10})
	memoryPod := testPod("memoryPod", requests{apiv1.ResourceMemory: 100})
	cpuMemoryPod := testPod("cpuMemoryPod", requests{apiv1.ResourceCPU: 10, apiv1.ResourceMemory: 30}, requests{apiv1.ResourceMemory: 20})
	tpuPod := testPod("tpuPod", requests{apiv1.ResourceName(ResourceTPUV2): 1})
	sanitizedTPUPod := testPod("tpuPod", requests{})
	preemptibleTPUPod := testPod("preemptibleTPUPod", requests{apiv1.ResourceName(ResourceTPUPreemptibleV2): 1})
	sanitizedPreemptibleTPUPod := testPod("preemptibleTPUPod", requests{})
	tpuMemoryPod := testPod("tpuMemoryPod", requests{apiv1.ResourceName(ResourceTPUV2): 1, apiv1.ResourceMemory: 30}, requests{apiv1.ResourceName(ResourceTPUV2): 2, apiv1.ResourceMemory: 13})
	sanitizedTPUMemoryPod := testPod("tpuMemoryPod", requests{apiv1.ResourceMemory: 30}, requests{apiv1.ResourceMemory: 13})
	podsWithoutTPUs := []*apiv1.Pod{cpuPod, memoryPod, cpuMemoryPod}
	mixedPods := []*apiv1.Pod{cpuPod, tpuPod, preemptibleTPUPod, memoryPod}
	sanitizedMixedPods := []*apiv1.Pod{cpuPod, sanitizedTPUPod, sanitizedPreemptibleTPUPod, memoryPod}
	podsWithTPUs := []*apiv1.Pod{tpuPod, preemptibleTPUPod, tpuMemoryPod}
	sanitizedPodsWithTPUs := []*apiv1.Pod{sanitizedTPUPod, sanitizedPreemptibleTPUPod, sanitizedTPUMemoryPod}
	testCases := []struct {
		desc		string
		pods		[]*apiv1.Pod
		expected	[]*apiv1.Pod
	}{{"Empty pod list", []*apiv1.Pod{}, []*apiv1.Pod{}}, {desc: "Pods without TPU request", pods: podsWithoutTPUs, expected: podsWithoutTPUs}, {desc: "Pods with & without TPU request", pods: mixedPods, expected: sanitizedMixedPods}, {desc: "Just TPU pods", pods: podsWithTPUs, expected: sanitizedPodsWithTPUs}}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			copied := make([]*apiv1.Pod, len(tc.pods))
			for i, pod := range tc.pods {
				copied[i] = pod.DeepCopy()
			}
			actual := ClearTPURequests(tc.pods)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Pod list should be as expected, got: %v, want: %v", actual, tc.expected)
			}
			if !reflect.DeepEqual(tc.pods, copied) {
				t.Errorf("Original pod list shouldn't be modified, got: %v, want: %v", tc.pods, copied)
			}
		})
	}
}
