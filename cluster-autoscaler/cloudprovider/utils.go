package cloudprovider

import (
	"fmt"
	"math/rand"
	"time"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
)

const (
	DefaultArch			= "amd64"
	DefaultOS			= "linux"
	KubeProxyCpuRequestMillis	= 100
)

func BuildReadyConditions() []apiv1.NodeCondition {
	_logClusterCodePath()
	defer _logClusterCodePath()
	lastTransition := time.Now().Add(-time.Minute)
	return []apiv1.NodeCondition{{Type: apiv1.NodeReady, Status: apiv1.ConditionTrue, LastTransitionTime: metav1.Time{Time: lastTransition}}, {Type: apiv1.NodeNetworkUnavailable, Status: apiv1.ConditionFalse, LastTransitionTime: metav1.Time{Time: lastTransition}}, {Type: apiv1.NodeOutOfDisk, Status: apiv1.ConditionFalse, LastTransitionTime: metav1.Time{Time: lastTransition}}, {Type: apiv1.NodeMemoryPressure, Status: apiv1.ConditionFalse, LastTransitionTime: metav1.Time{Time: lastTransition}}}
}
func BuildKubeProxy(name string) *apiv1.Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &apiv1.Pod{ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("kube-proxy-%s-%d", name, rand.Int63()), Namespace: "kube-system", Annotations: map[string]string{kubetypes.ConfigSourceAnnotationKey: kubetypes.FileSource, kubetypes.CriticalPodAnnotationKey: "true", kubetypes.ConfigMirrorAnnotationKey: "1234567890abcdef"}, Labels: map[string]string{"component": "kube-proxy", "tier": "node"}}, Spec: apiv1.PodSpec{Containers: []apiv1.Container{{Image: "kubeproxy", Resources: apiv1.ResourceRequirements{Requests: apiv1.ResourceList{apiv1.ResourceCPU: *resource.NewMilliQuantity(int64(KubeProxyCpuRequestMillis), resource.DecimalSI)}}}}}, Status: apiv1.PodStatus{Phase: apiv1.PodRunning, Conditions: []apiv1.PodCondition{{Type: apiv1.PodReady, Status: apiv1.ConditionTrue}}}}
}
func JoinStringMaps(items ...map[string]string) map[string]string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(map[string]string)
	for _, m := range items {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
