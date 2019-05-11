package model

import (
	"time"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/recommender/util"
)

var (
	MemoryAggregationWindowLength	= time.Hour * 8 * 24
	MemoryAggregationInterval		= time.Hour * 24
	CPUHistogramOptions				= cpuHistogramOptions()
	MemoryHistogramOptions			= memoryHistogramOptions()
	HistogramBucketSizeGrowth		= 0.05
	MemoryHistogramDecayHalfLife	= time.Hour * 24
	CPUHistogramDecayHalfLife		= time.Hour * 24
)

const (
	minSampleWeight	= 0.1
	epsilon			= 0.001 * minSampleWeight
)

func cpuHistogramOptions() util.HistogramOptions {
	_logClusterCodePath()
	defer _logClusterCodePath()
	options, err := util.NewExponentialHistogramOptions(1000.0, 0.01, 1.+HistogramBucketSizeGrowth, epsilon)
	if err != nil {
		panic("Invalid CPU histogram options")
	}
	return options
}
func memoryHistogramOptions() util.HistogramOptions {
	_logClusterCodePath()
	defer _logClusterCodePath()
	options, err := util.NewExponentialHistogramOptions(1e12, 1e7, 1.+HistogramBucketSizeGrowth, epsilon)
	if err != nil {
		panic("Invalid memory histogram options")
	}
	return options
}
