package aws

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/config"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	"k8s.io/klog"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

const (
	ProviderName = "aws"
)

type awsCloudProvider struct {
	awsManager	*AwsManager
	resourceLimiter	*cloudprovider.ResourceLimiter
}

func BuildAwsCloudProvider(awsManager *AwsManager, resourceLimiter *cloudprovider.ResourceLimiter) (cloudprovider.CloudProvider, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	aws := &awsCloudProvider{awsManager: awsManager, resourceLimiter: resourceLimiter}
	return aws, nil
}
func (aws *awsCloudProvider) Cleanup() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	aws.awsManager.Cleanup()
	return nil
}
func (aws *awsCloudProvider) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ProviderName
}
func (aws *awsCloudProvider) NodeGroups() []cloudprovider.NodeGroup {
	_logClusterCodePath()
	defer _logClusterCodePath()
	asgs := aws.awsManager.getAsgs()
	ngs := make([]cloudprovider.NodeGroup, len(asgs))
	for i, asg := range asgs {
		ngs[i] = &AwsNodeGroup{asg: asg, awsManager: aws.awsManager}
	}
	return ngs
}
func (aws *awsCloudProvider) NodeGroupForNode(node *apiv1.Node) (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(node.Spec.ProviderID) == 0 {
		klog.Warningf("Node %v has no providerId", node.Name)
		return nil, nil
	}
	ref, err := AwsRefFromProviderId(node.Spec.ProviderID)
	if err != nil {
		return nil, err
	}
	asg := aws.awsManager.GetAsgForInstance(*ref)
	if asg == nil {
		return nil, nil
	}
	return &AwsNodeGroup{asg: asg, awsManager: aws.awsManager}, nil
}
func (aws *awsCloudProvider) Pricing() (cloudprovider.PricingModel, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, cloudprovider.ErrNotImplemented
}
func (aws *awsCloudProvider) GetAvailableMachineTypes() ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return []string{}, nil
}
func (aws *awsCloudProvider) NewNodeGroup(machineType string, labels map[string]string, systemLabels map[string]string, taints []apiv1.Taint, extraResources map[string]resource.Quantity) (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, cloudprovider.ErrNotImplemented
}
func (aws *awsCloudProvider) GetResourceLimiter() (*cloudprovider.ResourceLimiter, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return aws.resourceLimiter, nil
}
func (aws *awsCloudProvider) Refresh() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return aws.awsManager.Refresh()
}
func (aws *awsCloudProvider) GetInstanceID(node *apiv1.Node) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return node.Spec.ProviderID
}

type AwsRef struct{ Name string }
type AwsInstanceRef struct {
	ProviderID	string
	Name		string
}

var validAwsRefIdRegex = regexp.MustCompile(`^aws\:\/\/\/[-0-9a-z]*\/[-0-9a-z]*$`)

func AwsRefFromProviderId(id string) (*AwsInstanceRef, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if validAwsRefIdRegex.FindStringSubmatch(id) == nil {
		return nil, fmt.Errorf("Wrong id: expected format aws:///<zone>/<name>, got %v", id)
	}
	splitted := strings.Split(id[7:], "/")
	return &AwsInstanceRef{ProviderID: id, Name: splitted[1]}, nil
}

type AwsNodeGroup struct {
	awsManager	*AwsManager
	asg		*asg
}

func (ng *AwsNodeGroup) MaxSize() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ng.asg.maxSize
}
func (ng *AwsNodeGroup) MinSize() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ng.asg.minSize
}
func (ng *AwsNodeGroup) TargetSize() (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ng.asg.curSize, nil
}
func (ng *AwsNodeGroup) Exist() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return true
}
func (ng *AwsNodeGroup) Create() (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, cloudprovider.ErrAlreadyExist
}
func (ng *AwsNodeGroup) Autoprovisioned() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return false
}
func (ng *AwsNodeGroup) Delete() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cloudprovider.ErrNotImplemented
}
func (ng *AwsNodeGroup) IncreaseSize(delta int) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if delta <= 0 {
		return fmt.Errorf("size increase must be positive")
	}
	size := ng.asg.curSize
	if size+delta > ng.asg.maxSize {
		return fmt.Errorf("size increase too large - desired:%d max:%d", size+delta, ng.asg.maxSize)
	}
	return ng.awsManager.SetAsgSize(ng.asg, size+delta)
}
func (ng *AwsNodeGroup) DecreaseTargetSize(delta int) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if delta >= 0 {
		return fmt.Errorf("size decrease size must be negative")
	}
	size := ng.asg.curSize
	nodes, err := ng.awsManager.GetAsgNodes(ng.asg.AwsRef)
	if err != nil {
		return err
	}
	if int(size)+delta < len(nodes) {
		return fmt.Errorf("attempt to delete existing nodes targetSize:%d delta:%d existingNodes: %d", size, delta, len(nodes))
	}
	return ng.awsManager.SetAsgSize(ng.asg, size+delta)
}
func (ng *AwsNodeGroup) Belongs(node *apiv1.Node) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ref, err := AwsRefFromProviderId(node.Spec.ProviderID)
	if err != nil {
		return false, err
	}
	targetAsg := ng.awsManager.GetAsgForInstance(*ref)
	if targetAsg == nil {
		return false, fmt.Errorf("%s doesn't belong to a known asg", node.Name)
	}
	if targetAsg.AwsRef != ng.asg.AwsRef {
		return false, nil
	}
	return true, nil
}
func (ng *AwsNodeGroup) DeleteNodes(nodes []*apiv1.Node) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	size := ng.asg.curSize
	if int(size) <= ng.MinSize() {
		return fmt.Errorf("min size reached, nodes will not be deleted")
	}
	refs := make([]*AwsInstanceRef, 0, len(nodes))
	for _, node := range nodes {
		belongs, err := ng.Belongs(node)
		if err != nil {
			return err
		}
		if belongs != true {
			return fmt.Errorf("%s belongs to a different asg than %s", node.Name, ng.Id())
		}
		awsref, err := AwsRefFromProviderId(node.Spec.ProviderID)
		if err != nil {
			return err
		}
		refs = append(refs, awsref)
	}
	return ng.awsManager.DeleteInstances(refs)
}
func (ng *AwsNodeGroup) Id() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ng.asg.Name
}
func (ng *AwsNodeGroup) Debug() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%s (%d:%d)", ng.Id(), ng.MinSize(), ng.MaxSize())
}
func (ng *AwsNodeGroup) Nodes() ([]cloudprovider.Instance, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	asgNodes, err := ng.awsManager.GetAsgNodes(ng.asg.AwsRef)
	if err != nil {
		return nil, err
	}
	instances := make([]cloudprovider.Instance, len(asgNodes))
	for i, asgNode := range asgNodes {
		instances[i] = cloudprovider.Instance{Id: asgNode.ProviderID}
	}
	return instances, nil
}
func (ng *AwsNodeGroup) TemplateNodeInfo() (*schedulercache.NodeInfo, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	template, err := ng.awsManager.getAsgTemplate(ng.asg)
	if err != nil {
		return nil, err
	}
	node, err := ng.awsManager.buildNodeFromTemplate(ng.asg, template)
	if err != nil {
		return nil, err
	}
	nodeInfo := schedulercache.NewNodeInfo(cloudprovider.BuildKubeProxy(ng.asg.Name))
	nodeInfo.SetNode(node)
	return nodeInfo, nil
}
func BuildAWS(opts config.AutoscalingOptions, do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var config io.ReadCloser
	if opts.CloudConfig != "" {
		var err error
		config, err = os.Open(opts.CloudConfig)
		if err != nil {
			klog.Fatalf("Couldn't open cloud provider configuration %s: %#v", opts.CloudConfig, err)
		}
		defer config.Close()
	}
	manager, err := CreateAwsManager(config, do)
	if err != nil {
		klog.Fatalf("Failed to create AWS Manager: %v", err)
	}
	provider, err := BuildAwsCloudProvider(manager, rl)
	if err != nil {
		klog.Fatalf("Failed to create AWS cloud provider: %v", err)
	}
	return provider
}
