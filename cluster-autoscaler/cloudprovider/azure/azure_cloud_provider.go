package azure

import (
	"io"
	"os"
	"strings"
	"k8s.io/klog"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/config"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
)

const (
	ProviderName = "azure"
)

type AzureCloudProvider struct {
	azureManager	*AzureManager
	resourceLimiter	*cloudprovider.ResourceLimiter
}

func BuildAzureCloudProvider(azureManager *AzureManager, resourceLimiter *cloudprovider.ResourceLimiter) (cloudprovider.CloudProvider, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	azure := &AzureCloudProvider{azureManager: azureManager, resourceLimiter: resourceLimiter}
	return azure, nil
}
func (azure *AzureCloudProvider) Cleanup() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	azure.azureManager.Cleanup()
	return nil
}
func (azure *AzureCloudProvider) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "azure"
}
func (azure *AzureCloudProvider) NodeGroups() []cloudprovider.NodeGroup {
	_logClusterCodePath()
	defer _logClusterCodePath()
	asgs := azure.azureManager.getAsgs()
	ngs := make([]cloudprovider.NodeGroup, len(asgs))
	for i, asg := range asgs {
		ngs[i] = asg
	}
	return ngs
}
func (azure *AzureCloudProvider) NodeGroupForNode(node *apiv1.Node) (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	klog.V(6).Infof("Searching for node group for the node: %s\n", node.Spec.ProviderID)
	ref := &azureRef{Name: node.Spec.ProviderID}
	return azure.azureManager.GetAsgForInstance(ref)
}
func (azure *AzureCloudProvider) Pricing() (cloudprovider.PricingModel, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, cloudprovider.ErrNotImplemented
}
func (azure *AzureCloudProvider) GetAvailableMachineTypes() ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return []string{}, nil
}
func (azure *AzureCloudProvider) NewNodeGroup(machineType string, labels map[string]string, systemLabels map[string]string, taints []apiv1.Taint, extraResources map[string]resource.Quantity) (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, cloudprovider.ErrNotImplemented
}
func (azure *AzureCloudProvider) GetResourceLimiter() (*cloudprovider.ResourceLimiter, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return azure.resourceLimiter, nil
}
func (azure *AzureCloudProvider) Refresh() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return azure.azureManager.Refresh()
}
func (azure *AzureCloudProvider) GetInstanceID(node *apiv1.Node) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return strings.ToLower(node.Spec.ProviderID)
}

type azureRef struct{ Name string }

func (m *azureRef) GetKey() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Name
}
func BuildAzure(opts config.AutoscalingOptions, do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var config io.ReadCloser
	if opts.CloudConfig != "" {
		klog.Infof("Creating Azure Manager using cloud-config file: %v", opts.CloudConfig)
		var err error
		config, err := os.Open(opts.CloudConfig)
		if err != nil {
			klog.Fatalf("Couldn't open cloud provider configuration %s: %#v", opts.CloudConfig, err)
		}
		defer config.Close()
	} else {
		klog.Info("Creating Azure Manager with default configuration.")
	}
	manager, err := CreateAzureManager(config, do)
	if err != nil {
		klog.Fatalf("Failed to create Azure Manager: %v", err)
	}
	provider, err := BuildAzureCloudProvider(manager, rl)
	if err != nil {
		klog.Fatalf("Failed to create Azure cloud provider: %v", err)
	}
	return provider
}
