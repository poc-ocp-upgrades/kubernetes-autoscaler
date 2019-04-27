package gke

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	caNamespace = "cluster_autoscaler"
)

var (
	requestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: caNamespace, Name: "gke_request_count", Help: "Counter of GKE API requests for each verb and API resource."}, []string{"resource", "verb"})
)

func registerMetrics() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	prometheus.MustRegister(requestCounter)
}
func registerRequest(resource string, verb string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	requestCounter.WithLabelValues(resource, verb).Add(1.0)
}
