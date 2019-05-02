package kubemark

import (
 apiv1 "k8s.io/api/core/v1"
 "k8s.io/apimachinery/pkg/api/resource"
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
 "k8s.io/autoscaler/cluster-autoscaler/config"
 "k8s.io/autoscaler/cluster-autoscaler/utils/errors"
 "k8s.io/klog"
)

const (
 ProviderName = "kubemark"
)

type KubemarkCloudProvider struct{}

func BuildKubemarkCloudProvider(kubemarkController interface{}, specs []string, resourceLimiter *cloudprovider.ResourceLimiter) (*KubemarkCloudProvider, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return nil, cloudprovider.ErrNotImplemented
}
func (kubemark *KubemarkCloudProvider) Name() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return ""
}
func (kubemark *KubemarkCloudProvider) NodeGroups() []cloudprovider.NodeGroup {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return []cloudprovider.NodeGroup{}
}
func (kubemark *KubemarkCloudProvider) Pricing() (cloudprovider.PricingModel, errors.AutoscalerError) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return nil, cloudprovider.ErrNotImplemented
}
func (kubemark *KubemarkCloudProvider) NodeGroupForNode(node *apiv1.Node) (cloudprovider.NodeGroup, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return nil, cloudprovider.ErrNotImplemented
}
func (kubemark *KubemarkCloudProvider) GetAvailableMachineTypes() ([]string, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return []string{}, cloudprovider.ErrNotImplemented
}
func (kubemark *KubemarkCloudProvider) NewNodeGroup(machineType string, labels map[string]string, systemLabels map[string]string, taints []apiv1.Taint, extraResources map[string]resource.Quantity) (cloudprovider.NodeGroup, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return nil, cloudprovider.ErrNotImplemented
}
func (kubemark *KubemarkCloudProvider) GetResourceLimiter() (*cloudprovider.ResourceLimiter, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return nil, cloudprovider.ErrNotImplemented
}
func (kubemark *KubemarkCloudProvider) Refresh() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return cloudprovider.ErrNotImplemented
}
func (kubemark *KubemarkCloudProvider) GetInstanceID(node *apiv1.Node) string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return ""
}
func (kubemark *KubemarkCloudProvider) Cleanup() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return cloudprovider.ErrNotImplemented
}
func BuildKubemark(opts config.AutoscalingOptions, do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
 _logClusterCodePath()
 defer _logClusterCodePath()
 klog.Fatal("Failed to create Kubemark cloud provider: only supported on Linux")
 return nil
}
