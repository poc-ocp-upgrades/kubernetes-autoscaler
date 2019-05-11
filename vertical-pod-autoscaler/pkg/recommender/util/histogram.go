package util

import (
	"fmt"
	"strings"
	"time"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
)

const (
	MaxCheckpointWeight uint32 = 10000
)

type Histogram interface {
	Percentile(percentile float64) float64
	AddSample(value float64, weight float64, time time.Time)
	SubtractSample(value float64, weight float64, time time.Time)
	Merge(other Histogram)
	IsEmpty() bool
	Equals(other Histogram) bool
	String() string
	SaveToChekpoint() (*vpa_types.HistogramCheckpoint, error)
	LoadFromCheckpoint(*vpa_types.HistogramCheckpoint) error
}

func NewHistogram(options HistogramOptions) Histogram {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &histogram{options: options, bucketWeight: make([]float64, options.NumBuckets()), totalWeight: 0.0, minBucket: options.NumBuckets() - 1, maxBucket: 0}
}

type histogram struct {
	options			HistogramOptions
	bucketWeight	[]float64
	totalWeight		float64
	minBucket		int
	maxBucket		int
}

func (h *histogram) AddSample(value float64, weight float64, time time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if weight < 0.0 {
		panic("sample weight must be non-negative")
	}
	bucket := h.options.FindBucket(value)
	h.bucketWeight[bucket] += weight
	h.totalWeight += weight
	if bucket < h.minBucket && h.bucketWeight[bucket] >= h.options.Epsilon() {
		h.minBucket = bucket
	}
	if bucket > h.maxBucket && h.bucketWeight[bucket] >= h.options.Epsilon() {
		h.maxBucket = bucket
	}
}
func safeSubtract(value, sub, epsilon float64) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	value -= sub
	if value < epsilon {
		return 0.0
	}
	return value
}
func (h *histogram) SubtractSample(value float64, weight float64, time time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if weight < 0.0 {
		panic("sample weight must be non-negative")
	}
	bucket := h.options.FindBucket(value)
	epsilon := h.options.Epsilon()
	h.totalWeight = safeSubtract(h.totalWeight, weight, epsilon)
	h.bucketWeight[bucket] = safeSubtract(h.bucketWeight[bucket], weight, epsilon)
	h.updateMinAndMaxBucket()
}
func (h *histogram) Merge(other Histogram) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	o := other.(*histogram)
	if h.options != o.options {
		panic("can't merge histograms with different options")
	}
	for bucket := o.minBucket; bucket <= o.maxBucket; bucket++ {
		h.bucketWeight[bucket] += o.bucketWeight[bucket]
	}
	h.totalWeight += o.totalWeight
	if o.minBucket < h.minBucket {
		h.minBucket = o.minBucket
	}
	if o.maxBucket > h.maxBucket {
		h.maxBucket = o.maxBucket
	}
}
func (h *histogram) Percentile(percentile float64) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if h.IsEmpty() {
		return 0.0
	}
	partialSum := 0.0
	threshold := percentile * h.totalWeight
	bucket := h.minBucket
	for ; bucket < h.maxBucket; bucket++ {
		partialSum += h.bucketWeight[bucket]
		if partialSum >= threshold {
			break
		}
	}
	if bucket < h.options.NumBuckets()-1 {
		return h.options.GetBucketStart(bucket + 1)
	}
	return h.options.GetBucketStart(bucket)
}
func (h *histogram) IsEmpty() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return h.bucketWeight[h.minBucket] < h.options.Epsilon()
}
func (h *histogram) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	lines := []string{fmt.Sprintf("minBucket: %d, maxBucket: %d, totalWeight: %.3f", h.minBucket, h.maxBucket, h.totalWeight), "%-tile\tvalue"}
	for i := 0; i <= 100; i += 5 {
		lines = append(lines, fmt.Sprintf("%d\t%.3f", i, h.Percentile(0.01*float64(i))))
	}
	return strings.Join(lines, "\n")
}
func (h *histogram) Equals(other Histogram) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	h2, typesMatch := other.(*histogram)
	if !typesMatch || h.options != h2.options || h.minBucket != h2.minBucket || h.maxBucket != h2.maxBucket {
		return false
	}
	for bucket := h.minBucket; bucket <= h.maxBucket; bucket++ {
		diff := h.bucketWeight[bucket] - h2.bucketWeight[bucket]
		if diff > 1e-15 || diff < -1e-15 {
			return false
		}
	}
	return true
}
func (h *histogram) updateMinAndMaxBucket() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	epsilon := h.options.Epsilon()
	lastBucket := h.options.NumBuckets() - 1
	for h.bucketWeight[h.minBucket] < epsilon && h.minBucket < lastBucket {
		h.minBucket++
	}
	for h.bucketWeight[h.maxBucket] < epsilon && h.maxBucket > 0 {
		h.maxBucket--
	}
}
func (h *histogram) SaveToChekpoint() (*vpa_types.HistogramCheckpoint, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := vpa_types.HistogramCheckpoint{BucketWeights: make(map[int]uint32)}
	result.TotalWeight = h.totalWeight
	max := 0.
	for bucket := h.minBucket; bucket <= h.maxBucket; bucket++ {
		if h.bucketWeight[bucket] > max {
			max = h.bucketWeight[bucket]
		}
	}
	ratio := float64(MaxCheckpointWeight) / max
	for bucket := h.minBucket; bucket <= h.maxBucket; bucket++ {
		newWeight := uint32(round(h.bucketWeight[bucket] * ratio))
		if newWeight > 0 {
			result.BucketWeights[bucket] = newWeight
		}
	}
	return &result, nil
}
func (h *histogram) LoadFromCheckpoint(checkpoint *vpa_types.HistogramCheckpoint) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if checkpoint == nil {
		return fmt.Errorf("Cannot load from empty checkpoint")
	}
	if checkpoint.TotalWeight < 0.0 {
		return fmt.Errorf("Cannot load checkpoint with negative weight %v", checkpoint.TotalWeight)
	}
	sum := int64(0)
	for bucket, weight := range checkpoint.BucketWeights {
		sum += int64(weight)
		if bucket >= h.options.NumBuckets() {
			return fmt.Errorf("Checkpoint has bucket %v that is exceeding histogram buckets %v", bucket, h.options.NumBuckets())
		}
		if bucket < 0 {
			return fmt.Errorf("Checkpoint has a negative bucket %v", bucket)
		}
	}
	if sum == 0 {
		return nil
	}
	ratio := checkpoint.TotalWeight / float64(sum)
	for bucket, weight := range checkpoint.BucketWeights {
		if bucket < h.minBucket {
			h.minBucket = bucket
		}
		if bucket > h.maxBucket {
			h.maxBucket = bucket
		}
		h.bucketWeight[bucket] += float64(weight) * ratio
	}
	h.totalWeight += checkpoint.TotalWeight
	return nil
}
func (h *histogram) scale(factor float64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if factor < 0.0 {
		panic("scale factor must be non-negative")
	}
	for bucket := h.minBucket; bucket <= h.maxBucket; bucket++ {
		h.bucketWeight[bucket] *= factor
	}
	h.totalWeight *= factor
	h.updateMinAndMaxBucket()
}
