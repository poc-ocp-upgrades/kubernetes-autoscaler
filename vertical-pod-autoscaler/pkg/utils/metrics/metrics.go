package metrics

import (
 "net/http"
 "time"
 "github.com/golang/glog"
 "github.com/prometheus/client_golang/prometheus"
 "github.com/prometheus/client_golang/prometheus/promhttp"
 _ "k8s.io/kubernetes/pkg/client/metrics/prometheus"
)

type ExecutionTimer struct {
 histo *prometheus.HistogramVec
 start time.Time
 last  time.Time
}

const (
 TopMetricsNamespace = "vpa_"
)

func Initialize(address string, healthCheck *HealthCheck) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 go func() {
  http.Handle("/metrics", promhttp.Handler())
  if healthCheck != nil {
   http.Handle("/health-check", healthCheck)
  }
  err := http.ListenAndServe(address, nil)
  glog.Fatalf("Failed to start metrics: %v", err)
 }()
}
func NewExecutionTimer(histo *prometheus.HistogramVec) *ExecutionTimer {
 _logClusterCodePath()
 defer _logClusterCodePath()
 now := time.Now()
 return &ExecutionTimer{histo: histo, start: now, last: now}
}
func (t *ExecutionTimer) ObserveStep(step string) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 now := time.Now()
 (*t.histo).WithLabelValues(step).Observe(now.Sub(t.last).Seconds())
 t.last = now
}
func (t *ExecutionTimer) ObserveTotal() {
 _logClusterCodePath()
 defer _logClusterCodePath()
 (*t.histo).WithLabelValues("total").Observe(time.Now().Sub(t.start).Seconds())
}
func CreateExecutionTimeMetric(namespace string, help string) *prometheus.HistogramVec {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return prometheus.NewHistogramVec(prometheus.HistogramOpts{Namespace: namespace, Name: "execution_latency_seconds", Help: help, Buckets: []float64{0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1.0, 2.0, 5.0, 10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0, 80.0, 90.0, 100.0, 120.0, 150.0, 240.0, 300.0}}, []string{"step"})
}
