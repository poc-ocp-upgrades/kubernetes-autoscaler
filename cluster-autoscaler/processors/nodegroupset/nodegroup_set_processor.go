package nodegroupset

import (
 "fmt"
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
 "k8s.io/autoscaler/cluster-autoscaler/context"
 "k8s.io/autoscaler/cluster-autoscaler/utils/errors"
 schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

type ScaleUpInfo struct {
 Group       cloudprovider.NodeGroup
 CurrentSize int
 NewSize     int
 MaxSize     int
}

func (s ScaleUpInfo) String() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return fmt.Sprintf("{%v %v->%v (max: %v)}", s.Group.Id(), s.CurrentSize, s.NewSize, s.MaxSize)
}

type NodeGroupSetProcessor interface {
 FindSimilarNodeGroups(context *context.AutoscalingContext, nodeGroup cloudprovider.NodeGroup, nodeInfosForGroups map[string]*schedulercache.NodeInfo) ([]cloudprovider.NodeGroup, errors.AutoscalerError)
 BalanceScaleUpBetweenGroups(context *context.AutoscalingContext, groups []cloudprovider.NodeGroup, newNodes int) ([]ScaleUpInfo, errors.AutoscalerError)
 CleanUp()
}
type NoOpNodeGroupSetProcessor struct{}

func (n *NoOpNodeGroupSetProcessor) FindSimilarNodeGroups(context *context.AutoscalingContext, nodeGroup cloudprovider.NodeGroup, nodeInfosForGroups map[string]*schedulercache.NodeInfo) ([]cloudprovider.NodeGroup, errors.AutoscalerError) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return []cloudprovider.NodeGroup{}, nil
}
func (n *NoOpNodeGroupSetProcessor) BalanceScaleUpBetweenGroups(context *context.AutoscalingContext, groups []cloudprovider.NodeGroup, newNodes int) ([]ScaleUpInfo, errors.AutoscalerError) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return []ScaleUpInfo{}, nil
}
func (n *NoOpNodeGroupSetProcessor) CleanUp() {
 _logClusterCodePath()
 defer _logClusterCodePath()
}
func NewDefaultNodeGroupSetProcessor() NodeGroupSetProcessor {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &BalancingNodeGroupSetProcessor{}
}
