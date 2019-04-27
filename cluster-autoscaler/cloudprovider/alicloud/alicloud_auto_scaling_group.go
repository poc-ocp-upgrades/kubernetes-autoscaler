package alicloud

import (
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/klog"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

type Asg struct {
	manager		*AliCloudManager
	minSize		int
	maxSize		int
	regionId	string
	id		string
}

func (asg *Asg) MaxSize() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return asg.maxSize
}
func (asg *Asg) MinSize() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return asg.minSize
}
func (asg *Asg) TargetSize() (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	size, err := asg.manager.GetAsgSize(asg)
	return int(size), err
}
func (asg *Asg) IncreaseSize(delta int) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	klog.Infof("increase ASG:%s with %d nodes", asg.Id(), delta)
	if delta <= 0 {
		return fmt.Errorf("size increase must be positive")
	}
	size, err := asg.manager.GetAsgSize(asg)
	if err != nil {
		klog.Errorf("failed to get ASG size because of %s", err.Error())
		return err
	}
	if int(size)+delta > asg.MaxSize() {
		return fmt.Errorf("size increase is too large - desired:%d max:%d", int(size)+delta, asg.MaxSize())
	}
	return asg.manager.SetAsgSize(asg, size+int64(delta))
}
func (asg *Asg) DecreaseTargetSize(delta int) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	klog.V(4).Infof("Aliyun: DecreaseTargetSize() with args: %v", delta)
	if delta >= 0 {
		return fmt.Errorf("size decrease size must be negative")
	}
	size, err := asg.manager.GetAsgSize(asg)
	if err != nil {
		klog.Errorf("failed to get ASG size because of %s", err.Error())
		return err
	}
	nodes, err := asg.manager.GetAsgNodes(asg)
	if err != nil {
		klog.Errorf("failed to get ASG nodes because of %s", err.Error())
		return err
	}
	if int(size)+delta < len(nodes) {
		return fmt.Errorf("attempt to delete existing nodes targetSize:%d delta:%d existingNodes: %d", size, delta, len(nodes))
	}
	return asg.manager.SetAsgSize(asg, size+int64(delta))
}
func (asg *Asg) Belongs(node *apiv1.Node) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	instanceId, err := ecsInstanceIdFromProviderId(node.Spec.ProviderID)
	if err != nil {
		return false, err
	}
	targetAsg, err := asg.manager.GetAsgForInstance(instanceId)
	if err != nil {
		return false, err
	}
	if targetAsg == nil {
		return false, fmt.Errorf("%s doesn't belong to a known Asg", node.Name)
	}
	if targetAsg.Id() != asg.Id() {
		return false, nil
	}
	return true, nil
}
func (asg *Asg) DeleteNodes(nodes []*apiv1.Node) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	size, err := asg.manager.GetAsgSize(asg)
	if err != nil {
		klog.Errorf("failed to get ASG size because of %s", err.Error())
		return err
	}
	if int(size) <= asg.MinSize() {
		return fmt.Errorf("min size reached, nodes will not be deleted")
	}
	nodeIds := make([]string, 0, len(nodes))
	for _, node := range nodes {
		belongs, err := asg.Belongs(node)
		if err != nil {
			klog.Errorf("failed to check whether node:%s is belong to asg:%s", node.GetName(), asg.Id())
			return err
		}
		if belongs != true {
			return fmt.Errorf("%s belongs to a different asg than %s", node.Name, asg.Id())
		}
		instanceId, err := ecsInstanceIdFromProviderId(node.Spec.ProviderID)
		if err != nil {
			klog.Errorf("failed to find instanceId from providerId,because of %s", err.Error())
			return err
		}
		nodeIds = append(nodeIds, instanceId)
	}
	return asg.manager.DeleteInstances(nodeIds)
}
func (asg *Asg) Id() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return asg.id
}
func (asg *Asg) RegionId() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return asg.regionId
}
func (asg *Asg) Debug() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%s (%d:%d)", asg.Id(), asg.MinSize(), asg.MaxSize())
}
func (asg *Asg) Nodes() ([]cloudprovider.Instance, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	instanceNames, err := asg.manager.GetAsgNodes(asg)
	if err != nil {
		return nil, err
	}
	instances := make([]cloudprovider.Instance, 0, len(instanceNames))
	for _, instanceName := range instanceNames {
		instances = append(instances, cloudprovider.Instance{Id: instanceName})
	}
	return instances, nil
}
func (asg *Asg) TemplateNodeInfo() (*schedulercache.NodeInfo, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	template, err := asg.manager.getAsgTemplate(asg.id)
	if err != nil {
		return nil, err
	}
	node, err := asg.manager.buildNodeFromTemplate(asg, template)
	if err != nil {
		klog.Errorf("failed to build instanceType:%v from template in ASG:%s,because of %s", template.InstanceType, asg.Id(), err.Error())
		return nil, err
	}
	nodeInfo := schedulercache.NewNodeInfo(cloudprovider.BuildKubeProxy(asg.id))
	nodeInfo.SetNode(node)
	return nodeInfo, nil
}
func (asg *Asg) Exist() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return true
}
func (asg *Asg) Create() (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, cloudprovider.ErrNotImplemented
}
func (asg *Asg) Autoprovisioned() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return false
}
func (asg *Asg) Delete() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cloudprovider.ErrNotImplemented
}
