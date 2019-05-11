package builder

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/gce"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/gke"
	"k8s.io/autoscaler/cluster-autoscaler/config"
)

var AvailableCloudProviders = []string{gce.ProviderNameGCE, gke.ProviderNameGKE}

const DefaultCloudProvider = gce.ProviderNameGCE

func buildCloudProvider(opts config.AutoscalingOptions, do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch opts.CloudProviderName {
	case gce.ProviderNameGCE:
		return gce.BuildGCE(opts, do, rl)
	case gke.ProviderNameGKE:
		return gke.BuildGKE(opts, do, rl)
	}
	return nil
}
