package builder

import (
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/openshiftmachineapi"
 "k8s.io/autoscaler/cluster-autoscaler/config"
)

var AvailableCloudProviders = []string{openshiftmachineapi.ProviderName}

const DefaultCloudProvider = openshiftmachineapi.ProviderName

func buildCloudProvider(opts config.AutoscalingOptions, do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
 _logClusterCodePath()
 defer _logClusterCodePath()
 switch opts.CloudProviderName {
 case openshiftmachineapi.ProviderName:
  return openshiftmachineapi.BuildOpenShiftMachineAPI(opts, do, rl)
 }
 return nil
}
