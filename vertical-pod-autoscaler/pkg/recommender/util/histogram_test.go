package util

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
)

var (
	valueEpsilon		= 1e-15
	weightEpsilon		= 1e-15
	anyTime			= time.Time{}
	testHistogramOptions, _	= NewLinearHistogramOptions(10.0, 1.0, weightEpsilon)
)

func TestPercentilesEmptyHistogram(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h := NewHistogram(testHistogramOptions)
	for p := -0.5; p <= 1.5; p += 0.5 {
		assert.Equal(t, 0.0, h.Percentile(p))
	}
}
func TestPercentiles(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h := NewHistogram(testHistogramOptions)
	for i := 1; i <= 4; i++ {
		h.AddSample(float64(i), float64(i), anyTime)
	}
	assert.InEpsilon(t, 2, h.Percentile(0.0), valueEpsilon)
	assert.InEpsilon(t, 2, h.Percentile(0.1), valueEpsilon)
	assert.InEpsilon(t, 3, h.Percentile(0.2), valueEpsilon)
	assert.InEpsilon(t, 3, h.Percentile(0.3), valueEpsilon)
	assert.InEpsilon(t, 4, h.Percentile(0.4), valueEpsilon)
	assert.InEpsilon(t, 4, h.Percentile(0.5), valueEpsilon)
	assert.InEpsilon(t, 4, h.Percentile(0.6), valueEpsilon)
	assert.InEpsilon(t, 5, h.Percentile(0.7), valueEpsilon)
	assert.InEpsilon(t, 5, h.Percentile(0.8), valueEpsilon)
	assert.InEpsilon(t, 5, h.Percentile(0.9), valueEpsilon)
	assert.InEpsilon(t, 5, h.Percentile(1.0), valueEpsilon)
}
func TestPercentileOutOfBounds(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	options, err := NewLinearHistogramOptions(1.0, 0.1, weightEpsilon)
	assert.Nil(t, err)
	h := NewHistogram(options)
	h.AddSample(0.1, 0.1, anyTime)
	h.AddSample(0.2, 0.2, anyTime)
	assert.InEpsilon(t, 0.2, h.Percentile(-0.1), valueEpsilon)
	assert.InEpsilon(t, 0.3, h.Percentile(1.1), valueEpsilon)
	h.AddSample(0.0, 0.1, anyTime)
	h.AddSample(1.0, 0.2, anyTime)
	assert.InEpsilon(t, 0.1, h.Percentile(-0.1), valueEpsilon)
	assert.InEpsilon(t, 1.0, h.Percentile(1.1), valueEpsilon)
}
func TestEmptyHistogram(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	options, err := NewLinearHistogramOptions(1.0, 0.1, weightEpsilon)
	assert.Nil(t, err)
	h := NewHistogram(options)
	assert.True(t, h.IsEmpty())
	h.AddSample(0.1, weightEpsilon*2.5, anyTime)
	assert.False(t, h.IsEmpty())
	h.SubtractSample(0.1, weightEpsilon, anyTime)
	assert.False(t, h.IsEmpty())
	h.SubtractSample(0.1, weightEpsilon, anyTime)
	assert.True(t, h.IsEmpty())
}
func TestNonEmptyOnEpsilonAddition(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	options, err := NewLinearHistogramOptions(1.0, 0.1, weightEpsilon)
	assert.Nil(t, err)
	h := NewHistogram(options)
	assert.True(t, h.IsEmpty())
	h.AddSample(9.9, weightEpsilon*3, anyTime)
	assert.False(t, h.IsEmpty())
	h.AddSample(0.1, weightEpsilon*0.3, anyTime)
	assert.False(t, h.IsEmpty())
	h.AddSample(999.9, weightEpsilon*0.3, anyTime)
	assert.False(t, h.IsEmpty())
}
func TestHistogramMerge(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h1 := NewHistogram(testHistogramOptions)
	h1.AddSample(1, 1, anyTime)
	h1.AddSample(2, 1, anyTime)
	h2 := NewHistogram(testHistogramOptions)
	h2.AddSample(2, 1, anyTime)
	h2.AddSample(3, 1, anyTime)
	expected := NewHistogram(testHistogramOptions)
	expected.AddSample(1, 1, anyTime)
	expected.AddSample(2, 1, anyTime)
	expected.AddSample(2, 1, anyTime)
	expected.AddSample(3, 1, anyTime)
	h1.Merge(h2)
	assert.True(t, h1.Equals(expected))
}
func TestHistogramSaveToCheckpointEmpty(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h := NewHistogram(testHistogramOptions)
	s, err := h.SaveToChekpoint()
	assert.NoError(t, err)
	assert.Equal(t, 0., s.TotalWeight)
	assert.Len(t, s.BucketWeights, 0)
}
func TestHistogramSaveToCheckpoint(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h := NewHistogram(testHistogramOptions)
	h.AddSample(1, 1, anyTime)
	s, err := h.SaveToChekpoint()
	assert.NoError(t, err)
	bucket := testHistogramOptions.FindBucket(1)
	assert.Equal(t, 1., s.TotalWeight)
	assert.Len(t, s.BucketWeights, 1)
	assert.Contains(t, s.BucketWeights, bucket)
	assert.Equal(t, MaxCheckpointWeight, s.BucketWeights[bucket])
}
func TestHistogramSaveToCheckpointDropsRelativelySmallValues(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h := NewHistogram(testHistogramOptions)
	v1, w1 := 1., 1.
	v2, w2 := 2., 100000.
	h.AddSample(v1, w1, anyTime)
	h.AddSample(v2, w2, anyTime)
	bucket1 := testHistogramOptions.FindBucket(v1)
	bucket2 := testHistogramOptions.FindBucket(v2)
	assert.NotEqualf(t, bucket1, bucket2, "For this test %v and %v have to be stored in different buckets", v1, v2)
	assert.True(t, w1 < (w2/float64(MaxCheckpointWeight))/2, "w1 to be omitted has to be less than (0.5*w2)/MaxCheckpointWeight")
	s, err := h.SaveToChekpoint()
	assert.NoError(t, err)
	assert.Equal(t, 100001., s.TotalWeight)
	assert.Len(t, s.BucketWeights, 1)
	assert.Contains(t, s.BucketWeights, bucket2)
	assert.Equal(t, MaxCheckpointWeight, s.BucketWeights[bucket2])
}
func TestHistogramSaveToCheckpointForMultipleValues(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h := NewHistogram(testHistogramOptions)
	v1, w1 := 1., 1.
	v2, w2 := 2., 10000.
	v3, w3 := 3., 50.
	h.AddSample(v1, w1, anyTime)
	h.AddSample(v2, w2, anyTime)
	h.AddSample(v3, w3, anyTime)
	bucket1 := testHistogramOptions.FindBucket(v1)
	bucket2 := testHistogramOptions.FindBucket(v2)
	bucket3 := testHistogramOptions.FindBucket(v3)
	assert.Truef(t, areUnique(bucket1, bucket2, bucket3), "For this test values %v have to be stored in different buckets", []float64{v1, v2, v3})
	s, err := h.SaveToChekpoint()
	assert.NoError(t, err)
	assert.Equal(t, 10051., s.TotalWeight)
	assert.Len(t, s.BucketWeights, 3)
	assert.Equal(t, uint32(1), s.BucketWeights[bucket1])
	assert.Equal(t, uint32(10000), s.BucketWeights[bucket2])
	assert.Equal(t, uint32(50), s.BucketWeights[bucket3])
}
func TestHistogramLoadFromCheckpoint(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	checkpoint := vpa_types.HistogramCheckpoint{TotalWeight: 6.0, BucketWeights: map[int]uint32{0: 1, 1: 2}}
	h := histogram{options: testHistogramOptions, bucketWeight: make([]float64, testHistogramOptions.NumBuckets()), totalWeight: 0.0, minBucket: testHistogramOptions.NumBuckets() - 1, maxBucket: 0}
	err := h.LoadFromCheckpoint(&checkpoint)
	assert.NoError(t, err)
	assert.Equal(t, 6.0, h.totalWeight)
	assert.Equal(t, 2.0, h.bucketWeight[0])
	assert.Equal(t, 4.0, h.bucketWeight[1])
}
func TestHistogramLoadFromCheckpointReturnsErrorOnNegativeBucket(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	checkpoint := vpa_types.HistogramCheckpoint{TotalWeight: 1.0, BucketWeights: map[int]uint32{-1: 1}}
	h := NewHistogram(testHistogramOptions)
	err := h.LoadFromCheckpoint(&checkpoint)
	assert.Error(t, err)
}
func TestHistogramLoadFromCheckpointReturnsErrorOnInvalidBucket(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	checkpoint := vpa_types.HistogramCheckpoint{TotalWeight: 1.0, BucketWeights: map[int]uint32{99: 1}}
	h := NewHistogram(testHistogramOptions)
	err := h.LoadFromCheckpoint(&checkpoint)
	assert.Error(t, err)
}
func TestHistogramLoadFromCheckpointReturnsErrorNegativeTotaWeight(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	checkpoint := vpa_types.HistogramCheckpoint{TotalWeight: -1.0, BucketWeights: map[int]uint32{}}
	h := NewHistogram(testHistogramOptions)
	err := h.LoadFromCheckpoint(&checkpoint)
	assert.Error(t, err)
}
func TestHistogramLoadFromCheckpointReturnsErrorOnNilInput(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h := NewHistogram(testHistogramOptions)
	err := h.LoadFromCheckpoint(nil)
	assert.Error(t, err)
}
func areUnique(values ...interface{}) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	dict := make(map[interface{}]bool)
	for i, v := range values {
		dict[v] = true
		if len(dict) != i+1 {
			return false
		}
	}
	return true
}
