package nodegroupset

import (
 "math"
 apiv1 "k8s.io/api/core/v1"
 "k8s.io/apimachinery/pkg/api/resource"
 kubeletapis "k8s.io/kubernetes/pkg/kubelet/apis"
 schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

const (
 MaxAllocatableDifferenceRatio = 0.05
 MaxFreeDifferenceRatio        = 0.05
)

type NodeInfoComparator func(n1, n2 *schedulercache.NodeInfo) bool

func compareResourceMapsWithTolerance(resources map[apiv1.ResourceName][]resource.Quantity, maxDifferenceRatio float64) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 for _, qtyList := range resources {
  if len(qtyList) != 2 {
   return false
  }
  larger := math.Max(float64(qtyList[0].MilliValue()), float64(qtyList[1].MilliValue()))
  smaller := math.Min(float64(qtyList[0].MilliValue()), float64(qtyList[1].MilliValue()))
  if larger-smaller > larger*maxDifferenceRatio {
   return false
  }
 }
 return true
}
func IsNodeInfoSimilar(n1, n2 *schedulercache.NodeInfo) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 capacity := make(map[apiv1.ResourceName][]resource.Quantity)
 allocatable := make(map[apiv1.ResourceName][]resource.Quantity)
 free := make(map[apiv1.ResourceName][]resource.Quantity)
 nodes := []*schedulercache.NodeInfo{n1, n2}
 for _, node := range nodes {
  for res, quantity := range node.Node().Status.Capacity {
   capacity[res] = append(capacity[res], quantity)
  }
  for res, quantity := range node.Node().Status.Allocatable {
   allocatable[res] = append(allocatable[res], quantity)
  }
  requested := node.RequestedResource()
  for res, quantity := range (&requested).ResourceList() {
   freeRes := node.Node().Status.Allocatable[res].DeepCopy()
   freeRes.Sub(quantity)
   free[res] = append(free[res], freeRes)
  }
 }
 for _, qtyList := range capacity {
  if len(qtyList) != 2 || qtyList[0].Cmp(qtyList[1]) != 0 {
   return false
  }
 }
 if !compareResourceMapsWithTolerance(allocatable, MaxAllocatableDifferenceRatio) {
  return false
 }
 if !compareResourceMapsWithTolerance(free, MaxFreeDifferenceRatio) {
  return false
 }
 ignoredLabels := map[string]bool{kubeletapis.LabelHostname: true, kubeletapis.LabelZoneFailureDomain: true, kubeletapis.LabelZoneRegion: true, "beta.kubernetes.io/fluentd-ds-ready": true}
 labels := make(map[string][]string)
 for _, node := range nodes {
  for label, value := range node.Node().ObjectMeta.Labels {
   ignore, _ := ignoredLabels[label]
   if !ignore {
    labels[label] = append(labels[label], value)
   }
  }
 }
 for _, labelValues := range labels {
  if len(labelValues) != 2 || labelValues[0] != labelValues[1] {
   return false
  }
 }
 return true
}
