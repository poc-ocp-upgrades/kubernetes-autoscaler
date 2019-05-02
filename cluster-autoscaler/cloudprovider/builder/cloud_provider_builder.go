package builder

import (
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
 "k8s.io/autoscaler/cluster-autoscaler/config"
 "k8s.io/autoscaler/cluster-autoscaler/context"
 "k8s.io/klog"
)

func NewCloudProvider(opts config.AutoscalingOptions) cloudprovider.CloudProvider {
 _logClusterCodePath()
 defer _logClusterCodePath()
 klog.V(1).Infof("Building %s cloud provider.", opts.CloudProviderName)
 do := cloudprovider.NodeGroupDiscoveryOptions{NodeGroupSpecs: opts.NodeGroups, NodeGroupAutoDiscoverySpecs: opts.NodeGroupAutoDiscovery}
 rl := context.NewResourceLimiterFromAutoscalingOptions(opts)
 if opts.CloudProviderName == "" {
  klog.Warning("Returning a nil cloud provider")
  return nil
 }
 provider := buildCloudProvider(opts, do, rl)
 if provider != nil {
  return provider
 }
 klog.Fatalf("Unknown cloud provider: %s", opts.CloudProviderName)
 return nil
}
