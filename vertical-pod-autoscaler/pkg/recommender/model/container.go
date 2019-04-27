package model

import (
	"fmt"
	"time"
)

const (
	OOMBumpUpRatio	float64	= 1.2
	OOMMinBumpUp	float64	= 100 * 1024 * 1024
)

type ContainerUsageSample struct {
	MeasureStart	time.Time
	Usage		ResourceAmount
	Request		ResourceAmount
	Resource	ResourceName
}
type ContainerState struct {
	Request			Resources
	LastCPUSampleStart	time.Time
	memoryPeak		ResourceAmount
	oomPeak			ResourceAmount
	WindowEnd		time.Time
	lastMemorySampleStart	time.Time
	aggregator		ContainerStateAggregator
}

func NewContainerState(request Resources, aggregator ContainerStateAggregator) *ContainerState {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &ContainerState{Request: request, LastCPUSampleStart: time.Time{}, WindowEnd: time.Time{}, lastMemorySampleStart: time.Time{}, aggregator: aggregator}
}
func (sample *ContainerUsageSample) isValid(expectedResource ResourceName) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return sample.Usage >= 0 && sample.Resource == expectedResource
}
func (container *ContainerState) addCPUSample(sample *ContainerUsageSample) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !sample.isValid(ResourceCPU) || !sample.MeasureStart.After(container.LastCPUSampleStart) {
		return false
	}
	container.aggregator.AddSample(sample)
	container.LastCPUSampleStart = sample.MeasureStart
	return true
}
func (container *ContainerState) GetMaxMemoryPeak() ResourceAmount {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ResourceAmountMax(container.memoryPeak, container.oomPeak)
}
func (container *ContainerState) addMemorySample(sample *ContainerUsageSample, isOOM bool) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ts := sample.MeasureStart
	if !sample.isValid(ResourceMemory) || ts.Before(container.lastMemorySampleStart) {
		return false
	}
	container.lastMemorySampleStart = ts
	if container.WindowEnd.IsZero() {
		container.WindowEnd = ts
	}
	addNewPeak := false
	if ts.Before(container.WindowEnd) {
		oldMaxMem := container.GetMaxMemoryPeak()
		if oldMaxMem != 0 && sample.Usage > oldMaxMem {
			oldPeak := ContainerUsageSample{MeasureStart: container.WindowEnd, Usage: oldMaxMem, Request: sample.Request, Resource: ResourceMemory}
			container.aggregator.SubtractSample(&oldPeak)
			addNewPeak = true
		}
	} else {
		shift := truncate(ts.Sub(container.WindowEnd), MemoryAggregationInterval) + MemoryAggregationInterval
		container.WindowEnd = container.WindowEnd.Add(shift)
		container.memoryPeak = 0
		container.oomPeak = 0
		addNewPeak = true
	}
	if addNewPeak {
		newPeak := ContainerUsageSample{MeasureStart: container.WindowEnd, Usage: sample.Usage, Request: sample.Request, Resource: ResourceMemory}
		container.aggregator.AddSample(&newPeak)
		if isOOM {
			container.oomPeak = sample.Usage
		} else {
			container.memoryPeak = sample.Usage
		}
	}
	return true
}
func (container *ContainerState) RecordOOM(timestamp time.Time, requestedMemory ResourceAmount) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if timestamp.Before(container.WindowEnd.Add(-1 * MemoryAggregationInterval)) {
		return fmt.Errorf("OOM event will be discarded - it is too old (%v)", timestamp)
	}
	memoryUsed := ResourceAmountMax(requestedMemory, container.memoryPeak)
	memoryNeeded := ResourceAmountMax(memoryUsed+MemoryAmountFromBytes(OOMMinBumpUp), ScaleResource(memoryUsed, OOMBumpUpRatio))
	oomMemorySample := ContainerUsageSample{MeasureStart: timestamp, Usage: memoryNeeded, Resource: ResourceMemory}
	if !container.addMemorySample(&oomMemorySample, true) {
		return fmt.Errorf("Adding OOM sample failed")
	}
	return nil
}
func (container *ContainerState) AddSample(sample *ContainerUsageSample) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch sample.Resource {
	case ResourceCPU:
		return container.addCPUSample(sample)
	case ResourceMemory:
		return container.addMemorySample(sample, false)
	default:
		return false
	}
}
func truncate(d, m time.Duration) time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m <= 0 {
		return d
	}
	return d - d%m
}
