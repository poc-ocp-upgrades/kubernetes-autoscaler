package metrics

import (
 "time"
 "k8s.io/autoscaler/cluster-autoscaler/utils/errors"
 "k8s.io/autoscaler/cluster-autoscaler/utils/gpu"
 _ "k8s.io/kubernetes/pkg/client/metrics/prometheus"
 "github.com/prometheus/client_golang/prometheus"
 "k8s.io/klog"
)

type NodeScaleDownReason string
type FailedScaleUpReason string
type FunctionLabel string
type NodeGroupType string

const (
 caNamespace                                  = "cluster_autoscaler"
 readyLabel                                   = "ready"
 unreadyLabel                                 = "unready"
 startingLabel                                = "notStarted"
 unregisteredLabel                            = "unregistered"
 longUnregisteredLabel                        = "longUnregistered"
 Underutilized            NodeScaleDownReason = "underutilized"
 Empty                    NodeScaleDownReason = "empty"
 Unready                  NodeScaleDownReason = "unready"
 APIError                 FailedScaleUpReason = "apiCallError"
 Timeout                  FailedScaleUpReason = "timeout"
 autoscaledGroup          NodeGroupType       = "autoscaled"
 autoprovisionedGroup     NodeGroupType       = "autoprovisioned"
 LogLongDurationThreshold                     = 5 * time.Second
)
const (
 ScaleDown                  FunctionLabel = "scaleDown"
 ScaleDownNodeDeletion      FunctionLabel = "scaleDown:nodeDeletion"
 ScaleDownFindNodesToRemove FunctionLabel = "scaleDown:findNodesToRemove"
 ScaleDownMiscOperations    FunctionLabel = "scaleDown:miscOperations"
 ScaleUp                    FunctionLabel = "scaleUp"
 FindUnneeded               FunctionLabel = "findUnneeded"
 UpdateState                FunctionLabel = "updateClusterState"
 FilterOutSchedulable       FunctionLabel = "filterOutSchedulable"
 Main                       FunctionLabel = "main"
 Poll                       FunctionLabel = "poll"
 Reconfigure                FunctionLabel = "reconfigure"
 Autoscaling                FunctionLabel = "autoscaling"
)

var (
 clusterSafeToAutoscale = prometheus.NewGauge(prometheus.GaugeOpts{Namespace: caNamespace, Name: "cluster_safe_to_autoscale", Help: "Whether or not cluster is healthy enough for autoscaling. 1 if it is, 0 otherwise."})
 nodesCount             = prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: caNamespace, Name: "nodes_count", Help: "Number of nodes in cluster."}, []string{"state"})
 nodeGroupsCount        = prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: caNamespace, Name: "node_groups_count", Help: "Number of node groups managed by CA."}, []string{"node_group_type"})
 unschedulablePodsCount = prometheus.NewGauge(prometheus.GaugeOpts{Namespace: caNamespace, Name: "unschedulable_pods_count", Help: "Number of unschedulable pods in the cluster."})
 lastActivity           = prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: caNamespace, Name: "last_activity", Help: "Last time certain part of CA logic executed."}, []string{"activity"})
 functionDuration       = prometheus.NewHistogramVec(prometheus.HistogramOpts{Namespace: caNamespace, Name: "function_duration_seconds", Help: "Time taken by various parts of CA main loop.", Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1.0, 2.5, 5.0, 7.5, 10.0, 12.5, 15.0, 17.5, 20.0, 22.5, 25.0, 27.5, 30.0, 50.0, 75.0, 100.0, 1000.0}}, []string{"function"})
 errorsCount            = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: caNamespace, Name: "errors_total", Help: "The number of CA loops failed due to an error."}, []string{"type"})
 scaleUpCount           = prometheus.NewCounter(prometheus.CounterOpts{Namespace: caNamespace, Name: "scaled_up_nodes_total", Help: "Number of nodes added by CA."})
 gpuScaleUpCount        = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: caNamespace, Name: "scaled_up_gpu_nodes_total", Help: "Number of GPU nodes added by CA, by GPU name."}, []string{"gpu_name"})
 failedScaleUpCount     = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: caNamespace, Name: "failed_scale_ups_total", Help: "Number of times scale-up operation has failed."}, []string{"reason"})
 scaleDownCount         = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: caNamespace, Name: "scaled_down_nodes_total", Help: "Number of nodes removed by CA."}, []string{"reason"})
 gpuScaleDownCount      = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: caNamespace, Name: "scaled_down_gpu_nodes_total", Help: "Number of GPU nodes removed by CA, by reason and GPU name."}, []string{"reason", "gpu_name"})
 evictionsCount         = prometheus.NewCounter(prometheus.CounterOpts{Namespace: caNamespace, Name: "evicted_pods_total", Help: "Number of pods evicted by CA"})
 unneededNodesCount     = prometheus.NewGauge(prometheus.GaugeOpts{Namespace: caNamespace, Name: "unneeded_nodes_count", Help: "Number of nodes currently considered unneeded by CA."})
 napEnabled             = prometheus.NewGauge(prometheus.GaugeOpts{Namespace: caNamespace, Name: "nap_enabled", Help: "Whether or not Node Autoprovisioning is enabled. 1 if it is, 0 otherwise."})
 nodeGroupCreationCount = prometheus.NewCounter(prometheus.CounterOpts{Namespace: caNamespace, Name: "created_node_groups_total", Help: "Number of node groups created by Node Autoprovisioning."})
 nodeGroupDeletionCount = prometheus.NewCounter(prometheus.CounterOpts{Namespace: caNamespace, Name: "deleted_node_groups_total", Help: "Number of node groups deleted by Node Autoprovisioning."})
)

func RegisterAll() {
 _logClusterCodePath()
 defer _logClusterCodePath()
 prometheus.MustRegister(clusterSafeToAutoscale)
 prometheus.MustRegister(nodesCount)
 prometheus.MustRegister(nodeGroupsCount)
 prometheus.MustRegister(unschedulablePodsCount)
 prometheus.MustRegister(lastActivity)
 prometheus.MustRegister(functionDuration)
 prometheus.MustRegister(errorsCount)
 prometheus.MustRegister(scaleUpCount)
 prometheus.MustRegister(gpuScaleUpCount)
 prometheus.MustRegister(failedScaleUpCount)
 prometheus.MustRegister(scaleDownCount)
 prometheus.MustRegister(gpuScaleDownCount)
 prometheus.MustRegister(evictionsCount)
 prometheus.MustRegister(unneededNodesCount)
 prometheus.MustRegister(napEnabled)
 prometheus.MustRegister(nodeGroupCreationCount)
 prometheus.MustRegister(nodeGroupDeletionCount)
}
func UpdateDurationFromStart(label FunctionLabel, start time.Time) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 duration := time.Now().Sub(start)
 UpdateDuration(label, duration)
}
func UpdateDuration(label FunctionLabel, duration time.Duration) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if duration > LogLongDurationThreshold && label != ScaleDown {
  klog.V(4).Infof("Function %s took %v to complete", label, duration)
 }
 functionDuration.WithLabelValues(string(label)).Observe(duration.Seconds())
}
func UpdateLastTime(label FunctionLabel, now time.Time) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 lastActivity.WithLabelValues(string(label)).Set(float64(now.Unix()))
}
func UpdateClusterSafeToAutoscale(safe bool) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if safe {
  clusterSafeToAutoscale.Set(1)
 } else {
  clusterSafeToAutoscale.Set(0)
 }
}
func UpdateNodesCount(ready, unready, starting, longUnregistered, unregistered int) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 nodesCount.WithLabelValues(readyLabel).Set(float64(ready))
 nodesCount.WithLabelValues(unreadyLabel).Set(float64(unready))
 nodesCount.WithLabelValues(startingLabel).Set(float64(starting))
 nodesCount.WithLabelValues(longUnregisteredLabel).Set(float64(longUnregistered))
 nodesCount.WithLabelValues(unregisteredLabel).Set(float64(unregistered))
}
func UpdateNodeGroupsCount(autoscaled, autoprovisioned int) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 nodeGroupsCount.WithLabelValues(string(autoscaledGroup)).Set(float64(autoscaled))
 nodeGroupsCount.WithLabelValues(string(autoprovisionedGroup)).Set(float64(autoprovisioned))
}
func UpdateUnschedulablePodsCount(podsCount int) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 unschedulablePodsCount.Set(float64(podsCount))
}
func RegisterError(err errors.AutoscalerError) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 errorsCount.WithLabelValues(string(err.Type())).Add(1.0)
}
func RegisterScaleUp(nodesCount int, gpuType string) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 scaleUpCount.Add(float64(nodesCount))
 if gpuType != gpu.MetricsNoGPU {
  gpuScaleUpCount.WithLabelValues(gpuType).Add(float64(nodesCount))
 }
}
func RegisterFailedScaleUp(reason FailedScaleUpReason) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 failedScaleUpCount.WithLabelValues(string(reason)).Inc()
}
func RegisterScaleDown(nodesCount int, gpuType string, reason NodeScaleDownReason) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 scaleDownCount.WithLabelValues(string(reason)).Add(float64(nodesCount))
 if gpuType != gpu.MetricsNoGPU {
  gpuScaleDownCount.WithLabelValues(string(reason), gpuType).Add(float64(nodesCount))
 }
}
func RegisterEvictions(podsCount int) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 evictionsCount.Add(float64(podsCount))
}
func UpdateUnneededNodesCount(nodesCount int) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 unneededNodesCount.Set(float64(nodesCount))
}
func UpdateNapEnabled(enabled bool) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if enabled {
  napEnabled.Set(1)
 } else {
  napEnabled.Set(0)
 }
}
func RegisterNodeGroupCreation() {
 _logClusterCodePath()
 defer _logClusterCodePath()
 nodeGroupCreationCount.Add(1.0)
}
func RegisterNodeGroupDeletion() {
 _logClusterCodePath()
 defer _logClusterCodePath()
 nodeGroupDeletionCount.Add(1.0)
}
