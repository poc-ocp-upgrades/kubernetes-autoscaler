package builder

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/aws"
	"k8s.io/autoscaler/cluster-autoscaler/config"
)

var AvailableCloudProviders = []string{aws.ProviderName}

const DefaultCloudProvider = aws.ProviderName

func buildCloudProvider(opts config.AutoscalingOptions, do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch opts.CloudProviderName {
	case aws.ProviderName:
		return aws.BuildAWS(opts, do, rl)
	}
	return nil
}
