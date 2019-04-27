package simulator

import (
	"testing"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	core "k8s.io/client-go/testing"
	"k8s.io/kubernetes/pkg/kubelet/types"
	"github.com/stretchr/testify/assert"
)

func TestRequiredPodsForNode(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pod1 := apiv1.Pod{ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "pod1", SelfLink: "pod1"}}
	pod2 := apiv1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod2", Namespace: "kube-system", SelfLink: "pod2", Annotations: map[string]string{types.ConfigMirrorAnnotationKey: "something"}}}
	fakeClient := &fake.Clientset{}
	fakeClient.Fake.AddReactor("list", "pods", func(action core.Action) (bool, runtime.Object, error) {
		return true, &apiv1.PodList{Items: []apiv1.Pod{pod1, pod2}}, nil
	})
	pods, err := GetRequiredPodsForNode("node1", fakeClient)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(pods))
	assert.Equal(t, "pod2", pods[0].Name)
}
