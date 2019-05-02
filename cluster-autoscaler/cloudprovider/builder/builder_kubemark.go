package builder

import (
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/kubemark"
 "k8s.io/autoscaler/cluster-autoscaler/config"
)

var AvailableCloudProviders = []string{kubemark.ProviderName}

const DefaultCloudProvider = kubemark.ProviderName

func buildCloudProvider(opts config.AutoscalingOptions, do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
 _logClusterCodePath()
 defer _logClusterCodePath()
 switch opts.CloudProviderName {
 case kubemark.ProviderName:
  return kubemark.BuildKubemark(opts, do, rl)
 }
 return nil
}
