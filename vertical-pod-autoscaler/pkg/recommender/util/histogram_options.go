package util

import (
	"errors"
	"fmt"
	"math"
)

type HistogramOptions interface {
	NumBuckets() int
	FindBucket(value float64) int
	GetBucketStart(bucket int) float64
	Epsilon() float64
}

func NewLinearHistogramOptions(maxValue float64, bucketSize float64, epsilon float64) (HistogramOptions, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if maxValue <= 0.0 || bucketSize <= 0.0 || epsilon <= 0.0 {
		return nil, errors.New("maxValue and bucketSize must both be positive")
	}
	numBuckets := int(math.Ceil(maxValue/bucketSize)) + 1
	return &linearHistogramOptions{numBuckets, bucketSize, epsilon}, nil
}
func NewExponentialHistogramOptions(maxValue float64, firstBucketSize float64, ratio float64, epsilon float64) (HistogramOptions, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if maxValue <= 0.0 || firstBucketSize <= 0.0 || ratio <= 1.0 || epsilon <= 0.0 {
		return nil, errors.New("maxValue, firstBucketSize and epsilon must be > 0.0, ratio must be > 1.0")
	}
	numBuckets := int(math.Ceil(log(ratio, maxValue*(ratio-1)/firstBucketSize+1))) + 1
	return &exponentialHistogramOptions{numBuckets, firstBucketSize, ratio, epsilon}, nil
}

type linearHistogramOptions struct {
	numBuckets	int
	bucketSize	float64
	epsilon		float64
}
type exponentialHistogramOptions struct {
	numBuckets		int
	firstBucketSize	float64
	ratio			float64
	epsilon			float64
}

func (o *linearHistogramOptions) NumBuckets() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return o.numBuckets
}
func (o *linearHistogramOptions) FindBucket(value float64) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bucket := int(value / o.bucketSize)
	if bucket < 0 {
		return 0
	}
	if bucket >= o.numBuckets {
		return o.numBuckets - 1
	}
	return bucket
}
func (o *linearHistogramOptions) GetBucketStart(bucket int) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if bucket < 0 || bucket >= o.numBuckets {
		panic(fmt.Sprintf("index %d out of range [0..%d]", bucket, o.numBuckets-1))
	}
	return float64(bucket) * o.bucketSize
}
func (o *linearHistogramOptions) Epsilon() float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return o.epsilon
}
func (o *exponentialHistogramOptions) NumBuckets() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return o.numBuckets
}
func (o *exponentialHistogramOptions) FindBucket(value float64) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if value < o.firstBucketSize {
		return 0
	}
	bucket := int(log(o.ratio, value*(o.ratio-1)/o.firstBucketSize+1))
	if bucket >= o.numBuckets {
		return o.numBuckets - 1
	}
	return bucket
}
func (o *exponentialHistogramOptions) GetBucketStart(bucket int) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if bucket < 0 || bucket >= o.numBuckets {
		panic(fmt.Sprintf("index %d out of range [0..%d]", bucket, o.numBuckets-1))
	}
	if bucket == 0 {
		return 0.0
	}
	return o.firstBucketSize * (math.Pow(o.ratio, float64(bucket)) - 1) / (o.ratio - 1)
}
func (o *exponentialHistogramOptions) Epsilon() float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return o.epsilon
}
func log(base, x float64) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return math.Log(x) / math.Log(base)
}
