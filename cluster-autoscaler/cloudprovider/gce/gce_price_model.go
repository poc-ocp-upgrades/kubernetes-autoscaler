package gce

import (
 "math"
 "time"
 apiv1 "k8s.io/api/core/v1"
 "k8s.io/autoscaler/cluster-autoscaler/utils/gpu"
 "k8s.io/autoscaler/cluster-autoscaler/utils/units"
 kubeletapis "k8s.io/kubernetes/pkg/kubelet/apis"
)

type GcePriceModel struct{}

const (
 cpuPricePerHour         = 0.033174
 memoryPricePerHourPerGb = 0.004446
 preemptibleDiscount     = 0.00698 / 0.033174
 gpuPricePerHour         = 0.700
 preemptibleLabel        = "cloud.google.com/gke-preemptible"
)

var (
 instancePrices    = map[string]float64{"n1-standard-1": 0.0475, "n1-standard-2": 0.0950, "n1-standard-4": 0.1900, "n1-standard-8": 0.3800, "n1-standard-16": 0.7600, "n1-standard-32": 1.5200, "n1-standard-64": 3.0400, "f1-micro": 0.0076, "g1-small": 0.0257, "n1-highmem-2": 0.1184, "n1-highmem-4": 0.2368, "n1-highmem-8": 0.4736, "n1-highmem-16": 0.9472, "n1-highmem-32": 1.8944, "n1-highmem-64": 3.7888, "n1-highcpu-2": 0.0709, "n1-highcpu-4": 0.1418, "n1-highcpu-8": 0.2836, "n1-highcpu-16": 0.5672, "n1-highcpu-32": 1.1344, "n1-highcpu-64": 2.2688}
 preemptiblePrices = map[string]float64{"n1-standard-1": 0.0100, "n1-standard-2": 0.0200, "n1-standard-4": 0.0400, "n1-standard-8": 0.0800, "n1-standard-16": 0.1600, "n1-standard-32": 0.3200, "n1-standard-64": 0.6400, "f1-micro": 0.0035, "g1-small": 0.0070, "n1-highmem-2": 0.0250, "n1-highmem-4": 0.0500, "n1-highmem-8": 0.1000, "n1-highmem-16": 0.2000, "n1-highmem-32": 0.4000, "n1-highmem-64": 0.8000, "n1-highcpu-2": 0.0150, "n1-highcpu-4": 0.0300, "n1-highcpu-8": 0.0600, "n1-highcpu-16": 0.1200, "n1-highcpu-32": 0.2400, "n1-highcpu-64": 0.4800}
)

func (model *GcePriceModel) NodePrice(node *apiv1.Node, startTime time.Time, endTime time.Time) (float64, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 price := 0.0
 basePriceFound := false
 if node.Labels != nil {
  if machineType, found := node.Labels[kubeletapis.LabelInstanceType]; found {
   var priceMapToUse map[string]float64
   if node.Labels[preemptibleLabel] == "true" {
    priceMapToUse = preemptiblePrices
   } else {
    priceMapToUse = instancePrices
   }
   if basePricePerHour, found := priceMapToUse[machineType]; found {
    price = basePricePerHour * getHours(startTime, endTime)
    basePriceFound = true
   }
  }
 }
 if !basePriceFound {
  price = getBasePrice(node.Status.Capacity, startTime, endTime)
  if node.Labels != nil && node.Labels[preemptibleLabel] == "true" {
   price = price * preemptibleDiscount
  }
 }
 price += getAdditionalPrice(node.Status.Capacity, startTime, endTime)
 return price, nil
}
func getHours(startTime time.Time, endTime time.Time) float64 {
 _logClusterCodePath()
 defer _logClusterCodePath()
 minutes := math.Ceil(float64(endTime.Sub(startTime)) / float64(time.Minute))
 hours := minutes / 60.0
 return hours
}
func (model *GcePriceModel) PodPrice(pod *apiv1.Pod, startTime time.Time, endTime time.Time) (float64, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 price := 0.0
 for _, container := range pod.Spec.Containers {
  price += getBasePrice(container.Resources.Requests, startTime, endTime)
  price += getAdditionalPrice(container.Resources.Requests, startTime, endTime)
 }
 return price, nil
}
func getBasePrice(resources apiv1.ResourceList, startTime time.Time, endTime time.Time) float64 {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if len(resources) == 0 {
  return 0
 }
 hours := getHours(startTime, endTime)
 price := 0.0
 cpu := resources[apiv1.ResourceCPU]
 mem := resources[apiv1.ResourceMemory]
 price += float64(cpu.MilliValue()) / 1000.0 * cpuPricePerHour * hours
 price += float64(mem.Value()) / float64(units.GiB) * memoryPricePerHourPerGb * hours
 return price
}
func getAdditionalPrice(resources apiv1.ResourceList, startTime time.Time, endTime time.Time) float64 {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if len(resources) == 0 {
  return 0
 }
 hours := getHours(startTime, endTime)
 price := 0.0
 gpu := resources[gpu.ResourceNvidiaGPU]
 price += float64(gpu.MilliValue()) / 1000.0 * gpuPricePerHour * hours
 return price
}
