package alicloud

import (
	"fmt"
	"strings"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/config"
	"k8s.io/autoscaler/cluster-autoscaler/config/dynamic"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	"k8s.io/klog"
	"os"
)

const (
	ProviderName = "alicloud"
)

type aliCloudProvider struct {
	manager		*AliCloudManager
	asgs		[]*Asg
	resourceLimiter	*cloudprovider.ResourceLimiter
}

func BuildAliCloudProvider(manager *AliCloudManager, discoveryOpts cloudprovider.NodeGroupDiscoveryOptions, resourceLimiter *cloudprovider.ResourceLimiter) (cloudprovider.CloudProvider, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if discoveryOpts.StaticDiscoverySpecified() {
		return buildStaticallyDiscoveringProvider(manager, discoveryOpts.NodeGroupSpecs, resourceLimiter)
	}
	if discoveryOpts.AutoDiscoverySpecified() {
		return nil, fmt.Errorf("only support static discovery scaling group in alicloud for now")
	}
	return nil, fmt.Errorf("failed to build alicloud provider: node group specs must be specified")
}
func buildStaticallyDiscoveringProvider(manager *AliCloudManager, specs []string, resourceLimiter *cloudprovider.ResourceLimiter) (*aliCloudProvider, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	acp := &aliCloudProvider{manager: manager, asgs: make([]*Asg, 0), resourceLimiter: resourceLimiter}
	for _, spec := range specs {
		if err := acp.addNodeGroup(spec); err != nil {
			klog.Warningf("failed to add node group to alicloud provider with spec: %s", spec)
			return nil, err
		}
	}
	return acp, nil
}
func (ali *aliCloudProvider) addNodeGroup(spec string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	asg, err := buildAsgFromSpec(spec, ali.manager)
	if err != nil {
		klog.Errorf("failed to build ASG from spec,because of %s", err.Error())
		return err
	}
	ali.addAsg(asg)
	return nil
}
func (ali *aliCloudProvider) addAsg(asg *Asg) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ali.asgs = append(ali.asgs, asg)
	ali.manager.RegisterAsg(asg)
}
func (ali *aliCloudProvider) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ProviderName
}
func (ali *aliCloudProvider) NodeGroups() []cloudprovider.NodeGroup {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make([]cloudprovider.NodeGroup, 0, len(ali.asgs))
	for _, asg := range ali.asgs {
		result = append(result, asg)
	}
	return result
}
func (ali *aliCloudProvider) NodeGroupForNode(node *apiv1.Node) (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	instanceId, err := ecsInstanceIdFromProviderId(node.Spec.ProviderID)
	if err != nil {
		klog.Errorf("failed to get instance Id from provider Id:%s,because of %s", node.Spec.ProviderID, err.Error())
		return nil, err
	}
	return ali.manager.GetAsgForInstance(instanceId)
}
func (ali *aliCloudProvider) Pricing() (cloudprovider.PricingModel, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, cloudprovider.ErrNotImplemented
}
func (ali *aliCloudProvider) GetAvailableMachineTypes() ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return []string{}, nil
}
func (ali *aliCloudProvider) NewNodeGroup(machineType string, labels map[string]string, systemLabels map[string]string, taints []apiv1.Taint, extraResources map[string]resource.Quantity) (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, cloudprovider.ErrNotImplemented
}
func (ali *aliCloudProvider) GetResourceLimiter() (*cloudprovider.ResourceLimiter, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ali.resourceLimiter, nil
}
func (ali *aliCloudProvider) Refresh() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (ali *aliCloudProvider) Cleanup() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (ali *aliCloudProvider) GetInstanceID(node *apiv1.Node) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return node.Spec.ProviderID
}

type AliRef struct {
	ID	string
	Region	string
}

func ecsInstanceIdFromProviderId(id string) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	parts := strings.Split(id, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("AliCloud: unexpected ProviderID format, providerID=%s", id)
	}
	return parts[1], nil
}
func buildAsgFromSpec(value string, manager *AliCloudManager) (*Asg, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	spec, err := dynamic.SpecFromString(value, true)
	if err != nil {
		return nil, fmt.Errorf("failed to parse node group spec: %v", err)
	}
	_, err = manager.aService.getScalingGroupByID(spec.Name)
	if err != nil {
		klog.Errorf("your scaling group: %s does not exist", spec.Name)
		return nil, err
	}
	asg := buildAsg(manager, spec.MinSize, spec.MaxSize, spec.Name, manager.cfg.getRegion())
	return asg, nil
}
func buildAsg(manager *AliCloudManager, minSize int, maxSize int, id string, regionId string) *Asg {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &Asg{manager: manager, minSize: minSize, maxSize: maxSize, regionId: regionId, id: id}
}
func BuildAlicloud(opts config.AutoscalingOptions, do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var aliManager *AliCloudManager
	var aliError error
	if opts.CloudConfig != "" {
		config, fileErr := os.Open(opts.CloudConfig)
		if fileErr != nil {
			klog.Fatalf("Couldn't open cloud provider configuration %s: %#v", opts.CloudConfig, fileErr)
		}
		defer config.Close()
		aliManager, aliError = CreateAliCloudManager(config)
	} else {
		aliManager, aliError = CreateAliCloudManager(nil)
	}
	if aliError != nil {
		klog.Fatalf("Failed to create Alicloud Manager: %v", aliError)
	}
	cloudProvider, err := BuildAliCloudProvider(aliManager, do, rl)
	if err != nil {
		klog.Fatalf("Failed to create Alicloud cloud provider: %v", err)
	}
	return cloudProvider
}
