package builder

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/aws"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/azure"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/gce"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/gke"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/openshiftmachineapi"
	"k8s.io/autoscaler/cluster-autoscaler/config"
)

var AvailableCloudProviders = []string{aws.ProviderName, azure.ProviderName, gce.ProviderNameGCE, gke.ProviderNameGKE, alicloud.ProviderName, openshiftmachineapi.ProviderName}

const DefaultCloudProvider = gce.ProviderNameGCE

func buildCloudProvider(opts config.AutoscalingOptions, do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch opts.CloudProviderName {
	case gce.ProviderNameGCE:
		return gce.BuildGCE(opts, do, rl)
	case gke.ProviderNameGKE:
		return gke.BuildGKE(opts, do, rl)
	case aws.ProviderName:
		return aws.BuildAWS(opts, do, rl)
	case azure.ProviderName:
		return azure.BuildAzure(opts, do, rl)
	case alicloud.ProviderName:
		return alicloud.BuildAlicloud(opts, do, rl)
	case openshiftmachineapi.ProviderName:
		return openshiftmachineapi.BuildOpenShiftMachineAPI(opts, do, rl)
	}
	return nil
}
