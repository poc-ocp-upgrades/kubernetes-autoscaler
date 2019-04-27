package model

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/recommender/util"
)

var (
	testPodID1		= PodID{"namespace-1", "pod-1"}
	testPodID2		= PodID{"namespace-1", "pod-2"}
	testContainerID1	= ContainerID{testPodID1, "container-1"}
	testRequest		= Resources{ResourceCPU: CPUAmountFromCores(3.14), ResourceMemory: MemoryAmountFromBytes(3.14e9)}
)

func addTestCPUSample(cluster *ClusterState, container ContainerID, cpuCores float64) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sample := ContainerUsageSampleWithKey{Container: container, ContainerUsageSample: ContainerUsageSample{MeasureStart: testTimestamp, Usage: CPUAmountFromCores(cpuCores), Request: testRequest[ResourceCPU], Resource: ResourceCPU}}
	return cluster.AddSample(&sample)
}
func addTestMemorySample(cluster *ClusterState, container ContainerID, memoryBytes float64) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sample := ContainerUsageSampleWithKey{Container: container, ContainerUsageSample: ContainerUsageSample{MeasureStart: testTimestamp, Usage: MemoryAmountFromBytes(memoryBytes), Request: testRequest[ResourceMemory], Resource: ResourceMemory}}
	return cluster.AddSample(&sample)
}
func TestAggregateStateByContainerName(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cluster := NewClusterState()
	cluster.AddOrUpdatePod(testPodID1, testLabels, apiv1.PodRunning)
	otherLabels := labels.Set{"label-2": "value-2"}
	cluster.AddOrUpdatePod(testPodID2, otherLabels, apiv1.PodRunning)
	containers := []ContainerID{{testPodID1, "app-A"}, {testPodID1, "app-B"}, {testPodID2, "app-A"}, {testPodID2, "app-C"}}
	for _, c := range containers {
		assert.NoError(t, cluster.AddOrUpdateContainer(c, testRequest))
	}
	assert.NoError(t, addTestCPUSample(cluster, containers[0], 1.0))
	assert.NoError(t, addTestCPUSample(cluster, containers[1], 5.0))
	assert.NoError(t, addTestCPUSample(cluster, containers[2], 3.0))
	assert.NoError(t, addTestCPUSample(cluster, containers[3], 5.0))
	assert.NoError(t, addTestMemorySample(cluster, containers[0], 2e9))
	assert.NoError(t, addTestMemorySample(cluster, containers[1], 10e9))
	assert.NoError(t, addTestMemorySample(cluster, containers[2], 4e9))
	assert.NoError(t, addTestMemorySample(cluster, containers[3], 10e9))
	aggregateResources := AggregateStateByContainerName(cluster.aggregateStateMap)
	assert.Contains(t, aggregateResources, "app-A")
	assert.Contains(t, aggregateResources, "app-B")
	assert.Contains(t, aggregateResources, "app-C")
	assert.Equal(t, 2, aggregateResources["app-A"].TotalSamplesCount)
	assert.Equal(t, 1, aggregateResources["app-B"].TotalSamplesCount)
	assert.Equal(t, 1, aggregateResources["app-C"].TotalSamplesCount)
	expectedCPUHistogram := util.NewDecayingHistogram(CPUHistogramOptions, CPUHistogramDecayHalfLife)
	expectedCPUHistogram.Merge(cluster.findOrCreateAggregateContainerState(containers[0]).AggregateCPUUsage)
	expectedCPUHistogram.Merge(cluster.findOrCreateAggregateContainerState(containers[2]).AggregateCPUUsage)
	actualCPUHistogram := aggregateResources["app-A"].AggregateCPUUsage
	expectedMemoryHistogram := util.NewDecayingHistogram(MemoryHistogramOptions, MemoryHistogramDecayHalfLife)
	expectedMemoryHistogram.AddSample(2e9, 1.0, cluster.GetContainer(containers[0]).WindowEnd)
	expectedMemoryHistogram.AddSample(4e9, 1.0, cluster.GetContainer(containers[2]).WindowEnd)
	actualMemoryHistogram := aggregateResources["app-A"].AggregateMemoryPeaks
	assert.True(t, expectedCPUHistogram.Equals(actualCPUHistogram), "Expected:\n%s\nActual:\n%s", expectedCPUHistogram, actualCPUHistogram)
	assert.True(t, expectedMemoryHistogram.Equals(actualMemoryHistogram), "Expected:\n%s\nActual:\n%s", expectedMemoryHistogram, actualMemoryHistogram)
}
func TestAggregateContainerStateSaveToCheckpoint(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	location, _ := time.LoadLocation("UTC")
	cs := NewAggregateContainerState()
	t1, t2 := time.Date(2018, time.January, 1, 2, 3, 4, 0, location), time.Date(2018, time.February, 1, 2, 3, 4, 0, location)
	cs.FirstSampleStart = t1
	cs.LastSampleStart = t2
	cs.TotalSamplesCount = 10
	cs.AggregateCPUUsage.AddSample(1, 33, t2)
	cs.AggregateMemoryPeaks.AddSample(1, 55, t1)
	cs.AggregateMemoryPeaks.AddSample(10000000, 55, t1)
	checkpoint, err := cs.SaveToCheckpoint()
	assert.NoError(t, err)
	assert.Equal(t, t1, checkpoint.FirstSampleStart.Time)
	assert.Equal(t, t2, checkpoint.LastSampleStart.Time)
	assert.Equal(t, 10, checkpoint.TotalSamplesCount)
	assert.Equal(t, SupportedCheckpointVersion, checkpoint.Version)
	assert.Len(t, checkpoint.CPUHistogram.BucketWeights, 1)
	assert.Len(t, checkpoint.MemoryHistogram.BucketWeights, 2)
}
func TestAggregateContainerStateLoadFromCheckpointFailsForVersionMismatch(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	checkpoint := vpa_types.VerticalPodAutoscalerCheckpointStatus{Version: "foo"}
	cs := NewAggregateContainerState()
	err := cs.LoadFromCheckpoint(&checkpoint)
	assert.Error(t, err)
}
func TestAggregateContainerStateLoadFromCheckpoint(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	location, _ := time.LoadLocation("UTC")
	t1, t2 := time.Date(2018, time.January, 1, 2, 3, 4, 0, location), time.Date(2018, time.February, 1, 2, 3, 4, 0, location)
	checkpoint := vpa_types.VerticalPodAutoscalerCheckpointStatus{Version: SupportedCheckpointVersion, FirstSampleStart: metav1.NewTime(t1), LastSampleStart: metav1.NewTime(t2), TotalSamplesCount: 20, MemoryHistogram: vpa_types.HistogramCheckpoint{BucketWeights: map[int]uint32{0: 10}, TotalWeight: 33.0}, CPUHistogram: vpa_types.HistogramCheckpoint{BucketWeights: map[int]uint32{0: 10}, TotalWeight: 44.0}}
	cs := NewAggregateContainerState()
	err := cs.LoadFromCheckpoint(&checkpoint)
	assert.NoError(t, err)
	assert.Equal(t, t1, cs.FirstSampleStart)
	assert.Equal(t, t2, cs.LastSampleStart)
	assert.Equal(t, 20, cs.TotalSamplesCount)
	assert.False(t, cs.AggregateCPUUsage.IsEmpty())
	assert.False(t, cs.AggregateMemoryPeaks.IsEmpty())
}
func TestAggregateContainerStateIsExpired(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cs := NewAggregateContainerState()
	cs.LastSampleStart = testTimestamp
	assert.False(t, cs.isExpired(testTimestamp.Add(7*24*time.Hour)))
	assert.True(t, cs.isExpired(testTimestamp.Add(8*24*time.Hour)))
}
