package nodegroups

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/context"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
)

type NodeGroupManager interface {
	CreateNodeGroup(context *context.AutoscalingContext, nodeGroup cloudprovider.NodeGroup) (CreateNodeGroupResult, errors.AutoscalerError)
	RemoveUnneededNodeGroups(context *context.AutoscalingContext) error
	CleanUp()
}
type NoOpNodeGroupManager struct{}
type CreateNodeGroupResult struct {
	MainCreatedNodeGroup	cloudprovider.NodeGroup
	ExtraCreatedNodeGroups	[]cloudprovider.NodeGroup
}

func (*NoOpNodeGroupManager) CreateNodeGroup(context *context.AutoscalingContext, nodeGroup cloudprovider.NodeGroup) (CreateNodeGroupResult, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return CreateNodeGroupResult{}, errors.NewAutoscalerError(errors.InternalError, "not implemented")
}
func (*NoOpNodeGroupManager) RemoveUnneededNodeGroups(context *context.AutoscalingContext) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (*NoOpNodeGroupManager) CleanUp() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func NewDefaultNodeGroupManager() NodeGroupManager {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &NoOpNodeGroupManager{}
}
