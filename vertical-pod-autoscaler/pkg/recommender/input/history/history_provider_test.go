package history

import (
	"fmt"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/recommender/model"
)

const (
	cpuQuery	= "container_cpu_usage_seconds_total{job=\"kubernetes-cadvisor\", pod_name=~\".+\"}[8d]"
	memoryQuery	= "container_memory_usage_bytes{job=\"kubernetes-cadvisor\", pod_name=~\".+\"}[8d]"
	labelsQuery	= "up{job=\"kubernetes-pods\"}[8d]"
)

type mockPrometheusClient struct{ mock.Mock }

func (m *mockPrometheusClient) GetTimeseries(query string) ([]Timeseries, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := m.Called(query)
	var returnArg []Timeseries
	if args.Get(0) != nil {
		returnArg = args.Get(0).([]Timeseries)
	}
	return returnArg, args.Error(1)
}
func TestGetEmptyClusterHistory(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	mockClient := mockPrometheusClient{}
	historyProvider := prometheusHistoryProvider{prometheusClient: &mockClient}
	mockClient.On("GetTimeseries", mock.AnythingOfType("string")).Times(3).Return([]Timeseries{}, nil)
	tss, err := historyProvider.GetClusterHistory()
	assert.Nil(t, err)
	assert.NotNil(t, tss)
	assert.Empty(t, tss)
}
func TestPrometheusError(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	mockClient := mockPrometheusClient{}
	historyProvider := prometheusHistoryProvider{prometheusClient: &mockClient}
	mockClient.On("GetTimeseries", mock.AnythingOfType("string")).Times(3).Return(nil, fmt.Errorf("bla"))
	_, err := historyProvider.GetClusterHistory()
	assert.NotNil(t, err)
}
func TestGetCPUSamples(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	mockClient := mockPrometheusClient{}
	historyProvider := prometheusHistoryProvider{prometheusClient: &mockClient}
	mockClient.On("GetTimeseries", cpuQuery).Return([]Timeseries{{Labels: map[string]string{"namespace": "default", "pod_name": "pod", "name": "container"}, Samples: []Sample{{Value: 5.5, Timestamp: time.Unix(1, 0)}}}}, nil)
	mockClient.On("GetTimeseries", memoryQuery).Return([]Timeseries{}, nil)
	mockClient.On("GetTimeseries", labelsQuery).Return([]Timeseries{}, nil)
	podID := model.PodID{Namespace: "default", PodName: "pod"}
	podHistory := &PodHistory{LastLabels: map[string]string{}, Samples: map[string][]model.ContainerUsageSample{"container": {{MeasureStart: time.Unix(1, 0), Usage: model.CPUAmountFromCores(5.5), Resource: model.ResourceCPU}}}}
	histories, err := historyProvider.GetClusterHistory()
	assert.Nil(t, err)
	assert.Equal(t, histories, map[model.PodID]*PodHistory{podID: podHistory})
}
func TestGetMemorySamples(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	mockClient := mockPrometheusClient{}
	historyProvider := prometheusHistoryProvider{prometheusClient: &mockClient}
	mockClient.On("GetTimeseries", cpuQuery).Return([]Timeseries{}, nil)
	mockClient.On("GetTimeseries", memoryQuery).Return([]Timeseries{{Labels: map[string]string{"namespace": "default", "pod_name": "pod", "name": "container"}, Samples: []Sample{{Value: 12345, Timestamp: time.Unix(1, 0)}}}}, nil)
	mockClient.On("GetTimeseries", labelsQuery).Return([]Timeseries{}, nil)
	podID := model.PodID{Namespace: "default", PodName: "pod"}
	podHistory := &PodHistory{LastLabels: map[string]string{}, Samples: map[string][]model.ContainerUsageSample{"container": {{MeasureStart: time.Unix(1, 0), Usage: model.MemoryAmountFromBytes(12345), Resource: model.ResourceMemory}}}}
	histories, err := historyProvider.GetClusterHistory()
	assert.Nil(t, err)
	assert.Equal(t, histories, map[model.PodID]*PodHistory{podID: podHistory})
}
func TestGetLabels(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	mockClient := mockPrometheusClient{}
	historyProvider := prometheusHistoryProvider{prometheusClient: &mockClient}
	mockClient.On("GetTimeseries", cpuQuery).Return([]Timeseries{}, nil)
	mockClient.On("GetTimeseries", memoryQuery).Return([]Timeseries{}, nil)
	mockClient.On("GetTimeseries", labelsQuery).Return([]Timeseries{{Labels: map[string]string{"kubernetes_namespace": "default", "kubernetes_pod_name": "pod", "pod_label_x": "y"}, Samples: []Sample{{Value: 1, Timestamp: time.Unix(10, 0)}}}, {Labels: map[string]string{"kubernetes_namespace": "default", "kubernetes_pod_name": "pod", "pod_label_x": "z"}, Samples: []Sample{{Value: 1, Timestamp: time.Unix(20, 0)}}}}, nil)
	podID := model.PodID{Namespace: "default", PodName: "pod"}
	podHistory := &PodHistory{LastLabels: map[string]string{"x": "z"}, LastSeen: time.Unix(20, 0), Samples: map[string][]model.ContainerUsageSample{}}
	histories, err := historyProvider.GetClusterHistory()
	assert.Nil(t, err)
	assert.Equal(t, histories, map[model.PodID]*PodHistory{podID: podHistory})
}
