package eviction

import (
	"fmt"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/utils/test"
	appsinformer "k8s.io/client-go/informers/apps/v1"
	coreinformer "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/pkg/api/testapi"
)

type podWithExpectations struct {
	pod		*apiv1.Pod
	canEvict	bool
	evictionSuccess	bool
}

func TestEvictReplicatedByController(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	rc := apiv1.ReplicationController{ObjectMeta: metav1.ObjectMeta{Name: "rc", Namespace: "default", SelfLink: testapi.Default.SelfLink("replicationcontrollers", "rc")}, TypeMeta: metav1.TypeMeta{Kind: "ReplicationController"}}
	index := 0
	generatePod := func() test.PodBuilder {
		index++
		return test.Pod().WithName(fmt.Sprintf("test-%v", index)).WithCreator(&rc.ObjectMeta, &rc.TypeMeta)
	}
	testCases := []struct {
		replicas		int32
		evictionTollerance	float64
		pods			[]podWithExpectations
	}{{replicas: 3, evictionTollerance: 0.5, pods: []podWithExpectations{{pod: generatePod().Get(), canEvict: true, evictionSuccess: true}, {pod: generatePod().Get(), canEvict: true, evictionSuccess: false}, {pod: generatePod().Get(), canEvict: true, evictionSuccess: false}}}, {replicas: 4, evictionTollerance: 0.5, pods: []podWithExpectations{{pod: generatePod().Get(), canEvict: true, evictionSuccess: true}, {pod: generatePod().Get(), canEvict: true, evictionSuccess: true}, {pod: generatePod().Get(), canEvict: true, evictionSuccess: false}, {pod: generatePod().Get(), canEvict: true, evictionSuccess: false}}}, {replicas: 4, evictionTollerance: 0.5, pods: []podWithExpectations{{pod: generatePod().Get(), canEvict: true, evictionSuccess: true}, {pod: generatePod().Get(), canEvict: true, evictionSuccess: false}, {pod: generatePod().Get(), canEvict: true, evictionSuccess: false}}}, {replicas: 3, evictionTollerance: 0.1, pods: []podWithExpectations{{pod: generatePod().Get(), canEvict: true, evictionSuccess: true}, {pod: generatePod().Get(), canEvict: true, evictionSuccess: false}, {pod: generatePod().Get(), canEvict: true, evictionSuccess: false}}}, {replicas: 3, evictionTollerance: 0.1, pods: []podWithExpectations{{pod: generatePod().Get(), canEvict: false, evictionSuccess: false}, {pod: generatePod().Get(), canEvict: false, evictionSuccess: false}}}, {replicas: 3, evictionTollerance: 0.5, pods: []podWithExpectations{{pod: generatePod().Get(), canEvict: false, evictionSuccess: false}, {pod: generatePod().WithPhase(apiv1.PodPending).Get(), canEvict: true, evictionSuccess: true}, {pod: generatePod().Get(), canEvict: false, evictionSuccess: false}}}, {replicas: 4, evictionTollerance: 0.5, pods: []podWithExpectations{{pod: generatePod().Get(), canEvict: false, evictionSuccess: false}, {pod: generatePod().WithPhase(apiv1.PodPending).Get(), canEvict: true, evictionSuccess: true}, {pod: generatePod().WithPhase(apiv1.PodPending).Get(), canEvict: true, evictionSuccess: true}, {pod: generatePod().WithPhase(apiv1.PodPending).Get(), canEvict: true, evictionSuccess: true}}}}
	for tcIndex, testCase := range testCases {
		rc.Spec = apiv1.ReplicationControllerSpec{Replicas: &testCase.replicas}
		pods := make([]*apiv1.Pod, 0, len(testCase.pods))
		for _, p := range testCase.pods {
			pods = append(pods, p.pod)
		}
		factory, _ := getEvictionRestrictionFactory(&rc, nil, nil, 2, testCase.evictionTollerance)
		eviction := factory.NewPodsEvictionRestriction(pods)
		for i, p := range testCase.pods {
			assert.Equalf(t, p.canEvict, eviction.CanEvict(p.pod), "TC %v - unexpected CanEvict result for pod-%v %#v", tcIndex, i, p.pod)
		}
		for i, p := range testCase.pods {
			err := eviction.Evict(p.pod, test.FakeEventRecorder())
			if p.evictionSuccess {
				assert.NoErrorf(t, err, "TC %v - unexpected Evict result for pod-%v %#v", tcIndex, i, p.pod)
			} else {
				assert.Errorf(t, err, "TC %v - unexpected Evict result for pod-%v %#v", tcIndex, i, p.pod)
			}
		}
	}
}
func TestEvictReplicatedByReplicaSet(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	replicas := int32(5)
	livePods := 5
	rs := appsv1.ReplicaSet{ObjectMeta: metav1.ObjectMeta{Name: "rs", Namespace: "default", SelfLink: testapi.Default.SelfLink("replicasets", "rs")}, TypeMeta: metav1.TypeMeta{Kind: "ReplicaSet"}, Spec: appsv1.ReplicaSetSpec{Replicas: &replicas}}
	pods := make([]*apiv1.Pod, livePods)
	for i := range pods {
		pods[i] = test.Pod().WithName(getTestPodName(i)).WithCreator(&rs.ObjectMeta, &rs.TypeMeta).Get()
	}
	factory, _ := getEvictionRestrictionFactory(nil, &rs, nil, 2, 0.5)
	eviction := factory.NewPodsEvictionRestriction(pods)
	for _, pod := range pods {
		assert.True(t, eviction.CanEvict(pod))
	}
	for _, pod := range pods[:2] {
		err := eviction.Evict(pod, test.FakeEventRecorder())
		assert.Nil(t, err, "Should evict with no error")
	}
	for _, pod := range pods[2:] {
		err := eviction.Evict(pod, test.FakeEventRecorder())
		assert.Error(t, err, "Error expected")
	}
}
func TestEvictReplicatedByStatefulSet(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	replicas := int32(5)
	livePods := 5
	ss := appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Name: "ss", Namespace: "default", SelfLink: "/apiv1s/apps/v1/namespaces/default/statefulsets/ss"}, TypeMeta: metav1.TypeMeta{Kind: "StatefulSet"}, Spec: appsv1.StatefulSetSpec{Replicas: &replicas}}
	pods := make([]*apiv1.Pod, livePods)
	for i := range pods {
		pods[i] = test.Pod().WithName(getTestPodName(i)).WithCreator(&ss.ObjectMeta, &ss.TypeMeta).Get()
	}
	factory, _ := getEvictionRestrictionFactory(nil, nil, &ss, 2, 0.5)
	eviction := factory.NewPodsEvictionRestriction(pods)
	for _, pod := range pods {
		assert.True(t, eviction.CanEvict(pod))
	}
	for _, pod := range pods[:2] {
		err := eviction.Evict(pod, test.FakeEventRecorder())
		assert.Nil(t, err, "Should evict with no error")
	}
	for _, pod := range pods[2:] {
		err := eviction.Evict(pod, test.FakeEventRecorder())
		assert.Error(t, err, "Error expected")
	}
}
func TestEvictReplicatedByJob(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	job := batchv1.Job{ObjectMeta: metav1.ObjectMeta{Name: "job", Namespace: "default", SelfLink: "/apiv1s/apps/v1/namespaces/default/jobs/job"}, TypeMeta: metav1.TypeMeta{Kind: "Job"}}
	livePods := 5
	pods := make([]*apiv1.Pod, livePods)
	for i := range pods {
		pods[i] = test.Pod().WithName(getTestPodName(i)).WithCreator(&job.ObjectMeta, &job.TypeMeta).Get()
	}
	factory, _ := getEvictionRestrictionFactory(nil, nil, nil, 2, 0.5)
	eviction := factory.NewPodsEvictionRestriction(pods)
	for _, pod := range pods {
		assert.True(t, eviction.CanEvict(pod))
	}
	for _, pod := range pods[:2] {
		err := eviction.Evict(pod, test.FakeEventRecorder())
		assert.Nil(t, err, "Should evict with no error")
	}
	for _, pod := range pods[2:] {
		err := eviction.Evict(pod, test.FakeEventRecorder())
		assert.Error(t, err, "Error expected")
	}
}
func TestEvictTooFewReplicas(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	replicas := int32(5)
	livePods := 5
	rc := apiv1.ReplicationController{ObjectMeta: metav1.ObjectMeta{Name: "rc", Namespace: "default", SelfLink: testapi.Default.SelfLink("replicationcontrollers", "rc")}, TypeMeta: metav1.TypeMeta{Kind: "ReplicationController"}, Spec: apiv1.ReplicationControllerSpec{Replicas: &replicas}}
	pods := make([]*apiv1.Pod, livePods)
	for i := range pods {
		pods[i] = test.Pod().WithName(getTestPodName(i)).WithCreator(&rc.ObjectMeta, &rc.TypeMeta).Get()
	}
	factory, _ := getEvictionRestrictionFactory(&rc, nil, nil, 10, 0.5)
	eviction := factory.NewPodsEvictionRestriction(pods)
	for _, pod := range pods {
		assert.False(t, eviction.CanEvict(pod))
	}
	for _, pod := range pods {
		err := eviction.Evict(pod, test.FakeEventRecorder())
		assert.Error(t, err, "Error expected")
	}
}
func TestEvictionTolerance(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	replicas := int32(5)
	livePods := 5
	tolerance := 0.8
	rc := apiv1.ReplicationController{ObjectMeta: metav1.ObjectMeta{Name: "rc", Namespace: "default", SelfLink: testapi.Default.SelfLink("replicationcontrollers", "rc")}, TypeMeta: metav1.TypeMeta{Kind: "ReplicationController"}, Spec: apiv1.ReplicationControllerSpec{Replicas: &replicas}}
	pods := make([]*apiv1.Pod, livePods)
	for i := range pods {
		pods[i] = test.Pod().WithName(getTestPodName(i)).WithCreator(&rc.ObjectMeta, &rc.TypeMeta).Get()
	}
	factory, _ := getEvictionRestrictionFactory(&rc, nil, nil, 2, tolerance)
	eviction := factory.NewPodsEvictionRestriction(pods)
	for _, pod := range pods {
		assert.True(t, eviction.CanEvict(pod))
	}
	for _, pod := range pods[:4] {
		err := eviction.Evict(pod, test.FakeEventRecorder())
		assert.Nil(t, err, "Should evict with no error")
	}
	for _, pod := range pods[4:] {
		err := eviction.Evict(pod, test.FakeEventRecorder())
		assert.Error(t, err, "Error expected")
	}
}
func TestEvictAtLeastOne(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	replicas := int32(5)
	livePods := 5
	tolerance := 0.1
	rc := apiv1.ReplicationController{ObjectMeta: metav1.ObjectMeta{Name: "rc", Namespace: "default", SelfLink: testapi.Default.SelfLink("replicationcontrollers", "rc")}, TypeMeta: metav1.TypeMeta{Kind: "ReplicationController"}, Spec: apiv1.ReplicationControllerSpec{Replicas: &replicas}}
	pods := make([]*apiv1.Pod, livePods)
	for i := range pods {
		pods[i] = test.Pod().WithName(getTestPodName(i)).WithCreator(&rc.ObjectMeta, &rc.TypeMeta).Get()
	}
	factory, _ := getEvictionRestrictionFactory(&rc, nil, nil, 2, tolerance)
	eviction := factory.NewPodsEvictionRestriction(pods)
	for _, pod := range pods {
		assert.True(t, eviction.CanEvict(pod))
	}
	for _, pod := range pods[:1] {
		err := eviction.Evict(pod, test.FakeEventRecorder())
		assert.Nil(t, err, "Should evict with no error")
	}
	for _, pod := range pods[1:] {
		err := eviction.Evict(pod, test.FakeEventRecorder())
		assert.Error(t, err, "Error expected")
	}
}
func getEvictionRestrictionFactory(rc *apiv1.ReplicationController, rs *appsv1.ReplicaSet, ss *appsv1.StatefulSet, minReplicas int, evictionToleranceFraction float64) (PodsEvictionRestrictionFactory, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	kubeClient := &fake.Clientset{}
	rcInformer := coreinformer.NewReplicationControllerInformer(kubeClient, apiv1.NamespaceAll, 0*time.Second, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	rsInformer := appsinformer.NewReplicaSetInformer(kubeClient, apiv1.NamespaceAll, 0*time.Second, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	ssInformer := appsinformer.NewStatefulSetInformer(kubeClient, apiv1.NamespaceAll, 0*time.Second, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	if rc != nil {
		err := rcInformer.GetIndexer().Add(rc)
		if err != nil {
			return nil, fmt.Errorf("Error adding object to cache: %v", err)
		}
	}
	if rs != nil {
		err := rsInformer.GetIndexer().Add(rs)
		if err != nil {
			return nil, fmt.Errorf("Error adding object to cache: %v", err)
		}
	}
	if ss != nil {
		err := ssInformer.GetIndexer().Add(ss)
		if err != nil {
			return nil, fmt.Errorf("Error adding object to cache: %v", err)
		}
	}
	return &podsEvictionRestrictionFactoryImpl{client: kubeClient, rsInformer: rsInformer, rcInformer: rcInformer, ssInformer: ssInformer, minReplicas: minReplicas, evictionToleranceFraction: evictionToleranceFraction}, nil
}
func getTestPodName(index int) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("test-%v", index)
}
