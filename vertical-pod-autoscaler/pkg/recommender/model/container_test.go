package model

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/recommender/util"
)

var (
	timeLayout		= "2006-01-02 15:04:05"
	testTimestamp, _	= time.Parse(timeLayout, "2017-04-18 17:35:05")
	TestRequest		= Resources{ResourceCPU: CPUAmountFromCores(2.3), ResourceMemory: MemoryAmountFromBytes(5e8)}
)

const (
	kb	= 1024
	mb	= 1024 * kb
)

func newUsageSample(timestamp time.Time, usage int64, resource ResourceName) *ContainerUsageSample {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &ContainerUsageSample{MeasureStart: timestamp, Usage: ResourceAmount(usage), Request: TestRequest[resource], Resource: resource}
}

type ContainerTest struct {
	mockCPUHistogram	*util.MockHistogram
	mockMemoryHistogram	*util.MockHistogram
	aggregateContainerState	*AggregateContainerState
	container		*ContainerState
}

func newContainerTest() ContainerTest {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mockCPUHistogram := new(util.MockHistogram)
	mockMemoryHistogram := new(util.MockHistogram)
	aggregateContainerState := &AggregateContainerState{AggregateCPUUsage: mockCPUHistogram, AggregateMemoryPeaks: mockMemoryHistogram}
	container := &ContainerState{Request: TestRequest, aggregator: aggregateContainerState}
	return ContainerTest{mockCPUHistogram: mockCPUHistogram, mockMemoryHistogram: mockMemoryHistogram, aggregateContainerState: aggregateContainerState, container: container}
}
func TestAggregateContainerUsageSamples(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	test := newContainerTest()
	c := test.container
	timeStep := MemoryAggregationInterval / 2
	test.mockCPUHistogram.On("AddSample", 3.14, 2.3, testTimestamp)
	test.mockCPUHistogram.On("AddSample", 6.28, 2.3, testTimestamp.Add(timeStep))
	test.mockCPUHistogram.On("AddSample", 1.57, 2.3, testTimestamp.Add(2*timeStep))
	memoryAggregationWindowEnd := testTimestamp.Add(MemoryAggregationInterval)
	test.mockMemoryHistogram.On("AddSample", 5.0, 1.0, memoryAggregationWindowEnd)
	test.mockMemoryHistogram.On("SubtractSample", 5.0, 1.0, memoryAggregationWindowEnd)
	test.mockMemoryHistogram.On("AddSample", 10.0, 1.0, memoryAggregationWindowEnd)
	memoryAggregationWindowEnd = memoryAggregationWindowEnd.Add(MemoryAggregationInterval)
	test.mockMemoryHistogram.On("AddSample", 2.0, 1.0, memoryAggregationWindowEnd)
	assert.True(t, c.AddSample(newUsageSample(testTimestamp, 3140, ResourceCPU)))
	assert.True(t, c.AddSample(newUsageSample(testTimestamp, 5, ResourceMemory)))
	assert.True(t, c.AddSample(newUsageSample(testTimestamp.Add(timeStep), 6280, ResourceCPU)))
	assert.True(t, c.AddSample(newUsageSample(testTimestamp.Add(timeStep), 10, ResourceMemory)))
	assert.True(t, c.AddSample(newUsageSample(testTimestamp.Add(2*timeStep), 1570, ResourceCPU)))
	assert.True(t, c.AddSample(newUsageSample(testTimestamp.Add(2*timeStep), 2, ResourceMemory)))
	assert.False(t, c.AddSample(newUsageSample(testTimestamp.Add(2*timeStep), 1000, ResourceCPU)))
	assert.False(t, c.AddSample(newUsageSample(testTimestamp.Add(4*timeStep), -1000, ResourceCPU)))
	assert.False(t, c.AddSample(newUsageSample(testTimestamp.Add(4*timeStep), -1000, ResourceMemory)))
}
func TestRecordOOMIncreasedByBumpUp(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	test := newContainerTest()
	memoryAggregationWindowEnd := testTimestamp.Add(MemoryAggregationInterval)
	test.mockMemoryHistogram.On("AddSample", 1200.0*mb, 1.0, memoryAggregationWindowEnd)
	assert.NoError(t, test.container.RecordOOM(testTimestamp, ResourceAmount(1000*mb)))
}
func TestRecordOOMDontRunAway(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	test := newContainerTest()
	memoryAggregationWindowEnd := testTimestamp.Add(MemoryAggregationInterval)
	test.mockMemoryHistogram.On("AddSample", 1200.0*mb, 1.0, memoryAggregationWindowEnd)
	assert.NoError(t, test.container.RecordOOM(testTimestamp, ResourceAmount(1000*mb)))
	assert.NoError(t, test.container.RecordOOM(testTimestamp, ResourceAmount(999*mb)))
	assert.NoError(t, test.container.RecordOOM(testTimestamp, ResourceAmount(999*mb)))
	test.mockMemoryHistogram.On("SubtractSample", 1200.0*mb, 1.0, memoryAggregationWindowEnd)
	test.mockMemoryHistogram.On("AddSample", 2400.0*mb, 1.0, memoryAggregationWindowEnd)
	assert.NoError(t, test.container.RecordOOM(testTimestamp, ResourceAmount(2000*mb)))
}
func TestRecordOOMIncreasedByMin(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	test := newContainerTest()
	memoryAggregationWindowEnd := testTimestamp.Add(MemoryAggregationInterval)
	test.mockMemoryHistogram.On("AddSample", 101.0*mb, 1.0, memoryAggregationWindowEnd)
	assert.NoError(t, test.container.RecordOOM(testTimestamp, ResourceAmount(1*mb)))
}
func TestRecordOOMMaxedWithKnownSample(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	test := newContainerTest()
	memoryAggregationWindowEnd := testTimestamp.Add(MemoryAggregationInterval)
	test.mockMemoryHistogram.On("AddSample", 3000.0*mb, 1.0, memoryAggregationWindowEnd)
	assert.True(t, test.container.AddSample(newUsageSample(testTimestamp, 3000*mb, ResourceMemory)))
	test.mockMemoryHistogram.On("SubtractSample", 3000.0*mb, 1.0, memoryAggregationWindowEnd)
	test.mockMemoryHistogram.On("AddSample", 3600.0*mb, 1.0, memoryAggregationWindowEnd)
	assert.NoError(t, test.container.RecordOOM(testTimestamp, ResourceAmount(1000*mb)))
}
func TestRecordOOMDiscardsOldSample(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	test := newContainerTest()
	memoryAggregationWindowEnd := testTimestamp.Add(MemoryAggregationInterval)
	test.mockMemoryHistogram.On("AddSample", 1000.0*mb, 1.0, memoryAggregationWindowEnd)
	assert.True(t, test.container.AddSample(newUsageSample(testTimestamp, 1000*mb, ResourceMemory)))
	assert.Error(t, test.container.RecordOOM(testTimestamp.Add(-30*time.Hour), ResourceAmount(1000*mb)))
}
func TestRecordOOMInNewWindow(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	test := newContainerTest()
	memoryAggregationWindowEnd := testTimestamp.Add(MemoryAggregationInterval)
	test.mockMemoryHistogram.On("AddSample", 2000.0*mb, 1.0, memoryAggregationWindowEnd)
	assert.True(t, test.container.AddSample(newUsageSample(testTimestamp, 2000*mb, ResourceMemory)))
	memoryAggregationWindowEnd = memoryAggregationWindowEnd.Add(2 * MemoryAggregationInterval)
	test.mockMemoryHistogram.On("AddSample", 2400.0*mb, 1.0, memoryAggregationWindowEnd)
	assert.NoError(t, test.container.RecordOOM(testTimestamp.Add(2*MemoryAggregationInterval), ResourceAmount(1000*mb)))
}
