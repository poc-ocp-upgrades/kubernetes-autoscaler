package builder

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/azure"
	"k8s.io/autoscaler/cluster-autoscaler/config"
)

var AvailableCloudProviders = []string{azure.ProviderName}

const DefaultCloudProvider = azure.ProviderName

func buildCloudProvider(opts config.AutoscalingOptions, do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch opts.CloudProviderName {
	case azure.ProviderName:
		return azure.BuildAzure(opts, do, rl)
	}
	return nil
}
