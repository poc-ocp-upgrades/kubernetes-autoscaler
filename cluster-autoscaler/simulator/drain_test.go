package simulator

import (
	"testing"
	apiv1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	. "k8s.io/autoscaler/cluster-autoscaler/utils/test"
	"k8s.io/kubernetes/pkg/kubelet/types"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
	"github.com/stretchr/testify/assert"
)

func TestFastGetPodsToMove(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pod1 := &apiv1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "ns"}}
	_, err := FastGetPodsToMove(schedulercache.NewNodeInfo(pod1), true, true, nil)
	assert.Error(t, err)
	pod2 := &apiv1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod2", Namespace: "ns", OwnerReferences: GenerateOwnerReferences("rs", "ReplicaSet", "extensions/v1beta1", "")}}
	r2, err := FastGetPodsToMove(schedulercache.NewNodeInfo(pod2), true, true, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r2))
	assert.Equal(t, pod2, r2[0])
	pod3 := &apiv1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod3", Namespace: "kube-system", Annotations: map[string]string{types.ConfigMirrorAnnotationKey: "something"}}}
	r3, err := FastGetPodsToMove(schedulercache.NewNodeInfo(pod3), true, true, nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(r3))
	pod4 := &apiv1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod4", Namespace: "ns", OwnerReferences: GenerateOwnerReferences("ds", "DaemonSet", "extensions/v1beta1", "")}}
	r4, err := FastGetPodsToMove(schedulercache.NewNodeInfo(pod2, pod3, pod4), true, true, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r4))
	assert.Equal(t, pod2, r4[0])
	pod5 := &apiv1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod5", Namespace: "kube-system", OwnerReferences: GenerateOwnerReferences("rs", "ReplicaSet", "extensions/v1beta1", "")}}
	_, err = FastGetPodsToMove(schedulercache.NewNodeInfo(pod5), true, true, nil)
	assert.Error(t, err)
	pod6 := &apiv1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod6", Namespace: "ns", OwnerReferences: GenerateOwnerReferences("rs", "ReplicaSet", "extensions/v1beta1", "")}, Spec: apiv1.PodSpec{Volumes: []apiv1.Volume{{VolumeSource: apiv1.VolumeSource{EmptyDir: &apiv1.EmptyDirVolumeSource{}}}}}}
	_, err = FastGetPodsToMove(schedulercache.NewNodeInfo(pod6), true, true, nil)
	assert.Error(t, err)
	pod7 := &apiv1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod7", Namespace: "ns", OwnerReferences: GenerateOwnerReferences("rs", "ReplicaSet", "extensions/v1beta1", "")}, Spec: apiv1.PodSpec{Volumes: []apiv1.Volume{{VolumeSource: apiv1.VolumeSource{GitRepo: &apiv1.GitRepoVolumeSource{Repository: "my-repo"}}}}}}
	r7, err := FastGetPodsToMove(schedulercache.NewNodeInfo(pod7), true, true, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r7))
	pod8 := &apiv1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod8", Namespace: "ns", OwnerReferences: GenerateOwnerReferences("rs", "ReplicaSet", "extensions/v1beta1", ""), Labels: map[string]string{"critical": "true"}}, Spec: apiv1.PodSpec{}}
	one := intstr.FromInt(1)
	pdb8 := &policyv1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{Name: "foobar", Namespace: "ns"}, Spec: policyv1.PodDisruptionBudgetSpec{MinAvailable: &one, Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"critical": "true"}}}, Status: policyv1.PodDisruptionBudgetStatus{PodDisruptionsAllowed: 0}}
	_, err = FastGetPodsToMove(schedulercache.NewNodeInfo(pod8), true, true, []*policyv1.PodDisruptionBudget{pdb8})
	assert.Error(t, err)
	pod9 := &apiv1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod9", Namespace: "ns", OwnerReferences: GenerateOwnerReferences("rs", "ReplicaSet", "extensions/v1beta1", ""), Labels: map[string]string{"critical": "true"}}, Spec: apiv1.PodSpec{}}
	pdb9 := &policyv1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{Name: "foobar", Namespace: "ns"}, Spec: policyv1.PodDisruptionBudgetSpec{MinAvailable: &one, Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"critical": "true"}}}, Status: policyv1.PodDisruptionBudgetStatus{PodDisruptionsAllowed: 1}}
	r9, err := FastGetPodsToMove(schedulercache.NewNodeInfo(pod9), true, true, []*policyv1.PodDisruptionBudget{pdb9})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r9))
}
