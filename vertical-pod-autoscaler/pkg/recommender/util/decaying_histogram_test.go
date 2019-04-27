package util

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
)

var (
	startTime = time.Unix(1234567890, 0)
)

func TestPercentilesEmptyDecayingHistogram(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h := NewDecayingHistogram(testHistogramOptions, time.Hour)
	for p := -0.5; p <= 1.5; p += 0.5 {
		assert.Equal(t, 0.0, h.Percentile(p))
	}
}
func TestSimpleDecay(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h := NewDecayingHistogram(testHistogramOptions, time.Hour)
	h.AddSample(2, 1000, startTime)
	h.AddSample(1, 1, startTime.Add(time.Hour*20))
	assert.InEpsilon(t, 2, h.Percentile(0.999), valueEpsilon)
	assert.InEpsilon(t, 3, h.Percentile(1.0), valueEpsilon)
}
func TestLongtermDecay(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h := NewDecayingHistogram(testHistogramOptions, time.Hour)
	h.AddSample(2, 1, startTime)
	h.AddSample(1, 1, startTime.Add(time.Hour*101))
	assert.InEpsilon(t, 2, h.Percentile(1.0), valueEpsilon)
}
func TestDecayingHistogramPercentiles(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h := NewDecayingHistogram(testHistogramOptions, time.Hour)
	timestamp := startTime
	for i := 1; i <= 4; i++ {
		h.AddSample(float64(i), float64(i), timestamp)
		timestamp = timestamp.Add(time.Hour)
	}
	assert.InEpsilon(t, 2, h.Percentile(0.00), valueEpsilon)
	assert.InEpsilon(t, 2, h.Percentile(0.02), valueEpsilon)
	assert.InEpsilon(t, 3, h.Percentile(0.03), valueEpsilon)
	assert.InEpsilon(t, 3, h.Percentile(0.10), valueEpsilon)
	assert.InEpsilon(t, 4, h.Percentile(0.11), valueEpsilon)
	assert.InEpsilon(t, 4, h.Percentile(0.34), valueEpsilon)
	assert.InEpsilon(t, 5, h.Percentile(0.35), valueEpsilon)
	assert.InEpsilon(t, 5, h.Percentile(1.00), valueEpsilon)
}
func TestNoDecay(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h := NewDecayingHistogram(testHistogramOptions, time.Hour)
	for i := 1; i <= 4; i++ {
		h.AddSample(float64(i), float64(i), startTime)
	}
	assert.InEpsilon(t, 2, h.Percentile(0.0), valueEpsilon)
	assert.InEpsilon(t, 3, h.Percentile(0.2), valueEpsilon)
	assert.InEpsilon(t, 2, h.Percentile(0.1), valueEpsilon)
	assert.InEpsilon(t, 3, h.Percentile(0.3), valueEpsilon)
	assert.InEpsilon(t, 4, h.Percentile(0.4), valueEpsilon)
	assert.InEpsilon(t, 4, h.Percentile(0.5), valueEpsilon)
	assert.InEpsilon(t, 4, h.Percentile(0.6), valueEpsilon)
	assert.InEpsilon(t, 5, h.Percentile(0.7), valueEpsilon)
	assert.InEpsilon(t, 5, h.Percentile(0.8), valueEpsilon)
	assert.InEpsilon(t, 5, h.Percentile(0.9), valueEpsilon)
	assert.InEpsilon(t, 5, h.Percentile(1.0), valueEpsilon)
}
func TestDecayingHistogramMerge(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h1 := NewDecayingHistogram(testHistogramOptions, time.Hour)
	h1.AddSample(1, 1, startTime)
	h1.AddSample(2, 1, startTime.Add(time.Hour))
	h2 := NewDecayingHistogram(testHistogramOptions, time.Hour)
	h2.AddSample(2, 1, startTime.Add(time.Hour*2))
	h2.AddSample(3, 1, startTime.Add(time.Hour))
	expected := NewDecayingHistogram(testHistogramOptions, time.Hour)
	expected.AddSample(2, 1, startTime.Add(time.Hour*2))
	expected.AddSample(2, 1, startTime.Add(time.Hour))
	expected.AddSample(3, 1, startTime.Add(time.Hour))
	expected.AddSample(1, 1, startTime)
	h1.Merge(h2)
	assert.True(t, h1.Equals(expected))
}
func TestDecayingHistogramSaveToCheckpoint(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	d := &decayingHistogram{histogram: *NewHistogram(testHistogramOptions).(*histogram), halfLife: time.Hour, referenceTimestamp: time.Time{}}
	d.AddSample(2, 1, startTime.Add(time.Hour*100))
	assert.NotEqual(t, d.referenceTimestamp, time.Time{})
	checkpoint, err := d.SaveToChekpoint()
	assert.NoError(t, err)
	assert.Equal(t, checkpoint.ReferenceTimestamp.Time, d.referenceTimestamp)
	assert.NotEmpty(t, checkpoint.BucketWeights)
	assert.NotZero(t, checkpoint.TotalWeight)
}
func TestDecayingHistogramLoadFromCheckpoint(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	location, _ := time.LoadLocation("UTC")
	timestamp := time.Date(2018, time.January, 2, 3, 4, 5, 0, location)
	checkpoint := vpa_types.HistogramCheckpoint{TotalWeight: 6.0, BucketWeights: map[int]uint32{0: 1}, ReferenceTimestamp: metav1.NewTime(timestamp)}
	d := &decayingHistogram{histogram: *NewHistogram(testHistogramOptions).(*histogram), halfLife: time.Hour, referenceTimestamp: time.Time{}}
	d.LoadFromCheckpoint(&checkpoint)
	assert.False(t, d.histogram.IsEmpty())
	assert.Equal(t, timestamp, d.referenceTimestamp)
}
