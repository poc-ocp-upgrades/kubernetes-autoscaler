package model

import (
	"testing"
	"time"
	"github.com/golang/glog"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
)

var (
	testPodID	= PodID{"namespace-1", "pod-1"}
	testContainerID	= ContainerID{testPodID, "container-1"}
	testVpaID	= VpaID{"namespace-1", "vpa-1"}
	testLabels	= map[string]string{"label-1": "value-1"}
	emptyLabels	= map[string]string{}
	testSelectorStr	= "label-1 = value-1"
)

func makeTestUsageSample() *ContainerUsageSampleWithKey {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &ContainerUsageSampleWithKey{ContainerUsageSample{MeasureStart: testTimestamp, Usage: 1.0, Request: testRequest[ResourceCPU], Resource: ResourceCPU}, testContainerID}
}
func TestClusterAddSample(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cluster := NewClusterState()
	cluster.AddOrUpdatePod(testPodID, testLabels, apiv1.PodRunning)
	assert.NoError(t, cluster.AddOrUpdateContainer(testContainerID, testRequest))
	assert.NoError(t, cluster.AddSample(makeTestUsageSample()))
	containerStats := cluster.Pods[testPodID].Containers["container-1"]
	assert.Equal(t, testTimestamp, containerStats.LastCPUSampleStart)
}
func TestClusterGCAggregateContainerStateDeletesOld(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cluster := NewClusterState()
	vpa := addTestVpa(cluster)
	addTestPod(cluster)
	assert.NoError(t, cluster.AddOrUpdateContainer(testContainerID, testRequest))
	usageSample := makeTestUsageSample()
	assert.NoError(t, cluster.AddSample(usageSample))
	assert.NotEmpty(t, cluster.aggregateStateMap)
	assert.NotEmpty(t, vpa.aggregateContainerStates)
	cluster.GarbageCollectAggregateCollectionStates(usageSample.MeasureStart.Add(9 * 24 * time.Hour))
	assert.Empty(t, cluster.aggregateStateMap)
	assert.Empty(t, vpa.aggregateContainerStates)
}
func TestClusterGCAggregateContainerStateLeavesValid(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cluster := NewClusterState()
	vpa := addTestVpa(cluster)
	addTestPod(cluster)
	assert.NoError(t, cluster.AddOrUpdateContainer(testContainerID, testRequest))
	usageSample := makeTestUsageSample()
	assert.NoError(t, cluster.AddSample(usageSample))
	assert.NotEmpty(t, cluster.aggregateStateMap)
	assert.NotEmpty(t, vpa.aggregateContainerStates)
	cluster.GarbageCollectAggregateCollectionStates(usageSample.MeasureStart.Add(7 * 24 * time.Hour))
	assert.NotEmpty(t, cluster.aggregateStateMap)
	assert.NotEmpty(t, vpa.aggregateContainerStates)
}
func TestAddSampleAfterAggregateContainerStateGCed(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cluster := NewClusterState()
	vpa := addTestVpa(cluster)
	pod := addTestPod(cluster)
	addTestContainer(cluster)
	assert.NoError(t, cluster.AddOrUpdateContainer(testContainerID, testRequest))
	usageSample := makeTestUsageSample()
	assert.NoError(t, cluster.AddSample(usageSample))
	assert.NotEmpty(t, cluster.aggregateStateMap)
	assert.NotEmpty(t, vpa.aggregateContainerStates)
	aggregateStateKey := cluster.aggregateStateKeyForContainerID(testContainerID)
	assert.Contains(t, vpa.aggregateContainerStates, aggregateStateKey)
	gcTimestamp := usageSample.MeasureStart.Add(10 * 24 * time.Hour)
	cluster.GarbageCollectAggregateCollectionStates(gcTimestamp)
	assert.Empty(t, cluster.aggregateStateMap)
	assert.Empty(t, vpa.aggregateContainerStates)
	assert.Contains(t, pod.Containers, testContainerID.ContainerName)
	newUsageSample := &ContainerUsageSampleWithKey{ContainerUsageSample{MeasureStart: gcTimestamp.Add(1 * time.Hour), Usage: usageSample.Usage, Request: usageSample.Request, Resource: usageSample.Resource}, testContainerID}
	assert.NoError(t, cluster.AddSample(newUsageSample))
	assert.Contains(t, vpa.aggregateContainerStates, aggregateStateKey)
}
func TestClusterRecordOOM(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cluster := NewClusterState()
	cluster.AddOrUpdatePod(testPodID, testLabels, apiv1.PodRunning)
	assert.NoError(t, cluster.AddOrUpdateContainer(testContainerID, testRequest))
	assert.NoError(t, cluster.RecordOOM(testContainerID, time.Unix(0, 0), ResourceAmount(10)))
	aggregation := cluster.findOrCreateAggregateContainerState(testContainerID)
	assert.NotEmpty(t, aggregation.AggregateMemoryPeaks)
}
func TestMissingKeys(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cluster := NewClusterState()
	err := cluster.AddSample(makeTestUsageSample())
	assert.EqualError(t, err, "KeyError: {namespace-1 pod-1}")
	err = cluster.RecordOOM(testContainerID, time.Unix(0, 0), ResourceAmount(10))
	assert.EqualError(t, err, "KeyError: {namespace-1 pod-1}")
	err = cluster.AddOrUpdateContainer(testContainerID, testRequest)
	assert.EqualError(t, err, "KeyError: {namespace-1 pod-1}")
}
func addVpa(cluster *ClusterState, id VpaID, selector string) *Vpa {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var apiObject vpa_types.VerticalPodAutoscaler
	apiObject.Namespace = id.Namespace
	apiObject.Name = id.VpaName
	apiObject.Spec.Selector, _ = metav1.ParseToLabelSelector(selector)
	err := cluster.AddOrUpdateVpa(&apiObject)
	if err != nil {
		glog.Fatalf("AddOrUpdateVpa() failed: %v", err)
	}
	return cluster.Vpas[id]
}
func addTestVpa(cluster *ClusterState) *Vpa {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return addVpa(cluster, testVpaID, testSelectorStr)
}
func addTestPod(cluster *ClusterState) *PodState {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cluster.AddOrUpdatePod(testPodID, testLabels, apiv1.PodRunning)
	return cluster.Pods[testPodID]
}
func addTestContainer(cluster *ClusterState) *ContainerState {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cluster.AddOrUpdateContainer(testContainerID, testRequest)
	return cluster.GetContainer(testContainerID)
}
func TestAddVpaThenAddPod(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cluster := NewClusterState()
	vpa := addTestVpa(cluster)
	assert.Empty(t, vpa.aggregateContainerStates)
	addTestPod(cluster)
	addTestContainer(cluster)
	aggregateStateKey := cluster.aggregateStateKeyForContainerID(testContainerID)
	assert.Contains(t, vpa.aggregateContainerStates, aggregateStateKey)
}
func TestAddPodThenAddVpa(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cluster := NewClusterState()
	addTestPod(cluster)
	addTestContainer(cluster)
	vpa := addTestVpa(cluster)
	aggregateStateKey := cluster.aggregateStateKeyForContainerID(testContainerID)
	assert.Contains(t, vpa.aggregateContainerStates, aggregateStateKey)
}
func TestChangePodLabels(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cluster := NewClusterState()
	vpa := addTestVpa(cluster)
	addTestPod(cluster)
	addTestContainer(cluster)
	aggregateStateKey := cluster.aggregateStateKeyForContainerID(testContainerID)
	assert.Contains(t, vpa.aggregateContainerStates, aggregateStateKey)
	cluster.AddOrUpdatePod(testPodID, emptyLabels, apiv1.PodRunning)
	aggregateStateKey = cluster.aggregateStateKeyForContainerID(testContainerID)
	assert.NotContains(t, vpa.aggregateContainerStates, aggregateStateKey)
}
func TestUpdatePodSelector(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cluster := NewClusterState()
	vpa := addTestVpa(cluster)
	addTestPod(cluster)
	addTestContainer(cluster)
	vpa = addVpa(cluster, testVpaID, "label-1 in (value-1,value-2)")
	assert.Contains(t, vpa.aggregateContainerStates, cluster.aggregateStateKeyForContainerID(testContainerID))
	vpa = addVpa(cluster, testVpaID, "label-1 = value-2")
	assert.NotContains(t, vpa.aggregateContainerStates, cluster.aggregateStateKeyForContainerID(testContainerID))
	vpa = addVpa(cluster, testVpaID, "label-1 = value-1")
	assert.Contains(t, vpa.aggregateContainerStates, cluster.aggregateStateKeyForContainerID(testContainerID))
}
func TestEqualAggregateStateKey(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cluster := NewClusterState()
	pod := addTestPod(cluster)
	key1 := cluster.MakeAggregateStateKey(pod, "container-1")
	key2 := cluster.MakeAggregateStateKey(pod, "container-1")
	assert.True(t, key1 == key2)
}
func TestTwoPodsWithSameLabels(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	podID1 := PodID{"namespace-1", "pod-1"}
	podID2 := PodID{"namespace-1", "pod-2"}
	containerID1 := ContainerID{podID1, "foo-container"}
	containerID2 := ContainerID{podID2, "foo-container"}
	cluster := NewClusterState()
	cluster.AddOrUpdatePod(podID1, testLabels, apiv1.PodRunning)
	cluster.AddOrUpdatePod(podID2, testLabels, apiv1.PodRunning)
	cluster.AddOrUpdateContainer(containerID1, testRequest)
	cluster.AddOrUpdateContainer(containerID2, testRequest)
	assert.Equal(t, 1, len(cluster.aggregateStateMap))
}
func TestTwoPodsWithDifferentNamespaces(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	podID1 := PodID{"namespace-1", "foo-pod"}
	podID2 := PodID{"namespace-2", "foo-pod"}
	containerID1 := ContainerID{podID1, "foo-container"}
	containerID2 := ContainerID{podID2, "foo-container"}
	cluster := NewClusterState()
	cluster.AddOrUpdatePod(podID1, testLabels, apiv1.PodRunning)
	cluster.AddOrUpdatePod(podID2, testLabels, apiv1.PodRunning)
	cluster.AddOrUpdateContainer(containerID1, testRequest)
	cluster.AddOrUpdateContainer(containerID2, testRequest)
	assert.Equal(t, 2, len(cluster.aggregateStateMap))
	assert.Equal(t, 1, len(cluster.labelSetMap))
}
func TestEmptySelector(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cluster := NewClusterState()
	vpa := addVpa(cluster, testVpaID, "")
	cluster.AddOrUpdatePod(testPodID, testLabels, apiv1.PodRunning)
	containerID1 := ContainerID{testPodID, "foo"}
	assert.NoError(t, cluster.AddOrUpdateContainer(containerID1, testRequest))
	anotherPodID := PodID{"namespace-1", "pod-2"}
	cluster.AddOrUpdatePod(anotherPodID, emptyLabels, apiv1.PodRunning)
	containerID2 := ContainerID{anotherPodID, "foo"}
	assert.NoError(t, cluster.AddOrUpdateContainer(containerID2, testRequest))
	assert.Contains(t, vpa.aggregateContainerStates, cluster.aggregateStateKeyForContainerID(containerID1))
	assert.Contains(t, vpa.aggregateContainerStates, cluster.aggregateStateKeyForContainerID(containerID2))
}
