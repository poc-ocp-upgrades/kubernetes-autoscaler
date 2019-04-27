package util

import (
	"time"
	"github.com/stretchr/testify/mock"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
)

type MockHistogram struct{ mock.Mock }

func (m *MockHistogram) Percentile(percentile float64) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := m.Called(percentile)
	return args.Get(0).(float64)
}
func (m *MockHistogram) AddSample(value float64, weight float64, time time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.Called(value, weight, time)
}
func (m *MockHistogram) SubtractSample(value float64, weight float64, time time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.Called(value, weight, time)
}
func (m *MockHistogram) IsEmpty() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := m.Called()
	return args.Bool(0)
}
func (m *MockHistogram) Equals(other Histogram) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := m.Called()
	return args.Bool(0)
}
func (m *MockHistogram) Merge(other Histogram) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.Called(other)
}
func (m *MockHistogram) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := m.Called()
	return args.String(0)
}
func (m *MockHistogram) SaveToChekpoint() (*vpa_types.HistogramCheckpoint, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &vpa_types.HistogramCheckpoint{}, nil
}
func (m *MockHistogram) LoadFromCheckpoint(checkpoint *vpa_types.HistogramCheckpoint) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
