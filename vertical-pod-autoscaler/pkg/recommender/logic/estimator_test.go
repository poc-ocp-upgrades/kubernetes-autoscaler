package logic

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/recommender/model"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/recommender/util"
)

var (
	anyTime		= time.Unix(0, 0)
	testRequest	= model.Resources{model.ResourceCPU: model.CPUAmountFromCores(3.14), model.ResourceMemory: model.MemoryAmountFromBytes(3.14e9)}
)

func TestPercentileEstimator(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cpuHistogram := util.NewHistogram(model.CPUHistogramOptions)
	cpuHistogram.AddSample(1.0, 1.0, anyTime)
	cpuHistogram.AddSample(2.0, 1.0, anyTime)
	cpuHistogram.AddSample(3.0, 1.0, anyTime)
	memoryPeaksHistogram := util.NewHistogram(model.MemoryHistogramOptions)
	memoryPeaksHistogram.AddSample(1e9, 1.0, anyTime)
	memoryPeaksHistogram.AddSample(2e9, 1.0, anyTime)
	memoryPeaksHistogram.AddSample(3e9, 1.0, anyTime)
	CPUPercentile := 0.2
	MemoryPercentile := 0.5
	estimator := NewPercentileEstimator(CPUPercentile, MemoryPercentile)
	resourceEstimation := estimator.GetResourceEstimation(&model.AggregateContainerState{AggregateCPUUsage: cpuHistogram, AggregateMemoryPeaks: memoryPeaksHistogram})
	maxRelativeError := 0.05
	assert.InEpsilon(t, 1.0, model.CoresFromCPUAmount(resourceEstimation[model.ResourceCPU]), maxRelativeError)
	assert.InEpsilon(t, 2e9, model.BytesFromMemoryAmount(resourceEstimation[model.ResourceMemory]), maxRelativeError)
}
func TestConfidenceMultiplier(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	baseEstimator := NewConstEstimator(model.Resources{model.ResourceCPU: model.CPUAmountFromCores(3.14), model.ResourceMemory: model.MemoryAmountFromBytes(3.14e9)})
	testedEstimator := &confidenceMultiplier{0.1, 2.0, baseEstimator}
	s := model.NewAggregateContainerState()
	timestamp := anyTime
	for i := 1; i <= 9; i++ {
		s.AddSample(&model.ContainerUsageSample{timestamp, model.CPUAmountFromCores(1.0), testRequest[model.ResourceCPU], model.ResourceCPU})
		timestamp = timestamp.Add(time.Minute * 2)
	}
	assert.Equal(t, 0.00625, getConfidence(s))
	resourceEstimation := testedEstimator.GetResourceEstimation(s)
	assert.Equal(t, 907.46, model.CoresFromCPUAmount(resourceEstimation[model.ResourceCPU]))
}
func TestConfidenceMultiplierNoHistory(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	baseEstimator := NewConstEstimator(model.Resources{model.ResourceCPU: model.CPUAmountFromCores(3.14), model.ResourceMemory: model.MemoryAmountFromBytes(3.14e9)})
	testedEstimator1 := &confidenceMultiplier{1.0, 1.0, baseEstimator}
	testedEstimator2 := &confidenceMultiplier{1.0, -1.0, baseEstimator}
	s := model.NewAggregateContainerState()
	assert.Equal(t, model.ResourceAmount(1e14), testedEstimator1.GetResourceEstimation(s)[model.ResourceCPU])
	assert.Equal(t, model.ResourceAmount(0), testedEstimator2.GetResourceEstimation(s)[model.ResourceCPU])
}
func TestMarginEstimator(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	marginFraction := 0.1
	baseEstimator := NewConstEstimator(model.Resources{model.ResourceCPU: model.CPUAmountFromCores(3.14), model.ResourceMemory: model.MemoryAmountFromBytes(3.14e9)})
	testedEstimator := &marginEstimator{marginFraction: marginFraction, baseEstimator: baseEstimator}
	s := model.NewAggregateContainerState()
	resourceEstimation := testedEstimator.GetResourceEstimation(s)
	assert.Equal(t, 3.14*1.1, model.CoresFromCPUAmount(resourceEstimation[model.ResourceCPU]))
	assert.Equal(t, 3.14e9*1.1, model.BytesFromMemoryAmount(resourceEstimation[model.ResourceMemory]))
}
func TestMinResourcesEstimator(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	minResources := model.Resources{model.ResourceCPU: model.CPUAmountFromCores(0.2), model.ResourceMemory: model.MemoryAmountFromBytes(4e8)}
	baseEstimator := NewConstEstimator(model.Resources{model.ResourceCPU: model.CPUAmountFromCores(3.14), model.ResourceMemory: model.MemoryAmountFromBytes(2e7)})
	testedEstimator := &minResourcesEstimator{minResources: minResources, baseEstimator: baseEstimator}
	s := model.NewAggregateContainerState()
	resourceEstimation := testedEstimator.GetResourceEstimation(s)
	assert.Equal(t, 3.14, model.CoresFromCPUAmount(resourceEstimation[model.ResourceCPU]))
	assert.Equal(t, 4e8, model.BytesFromMemoryAmount(resourceEstimation[model.ResourceMemory]))
}
