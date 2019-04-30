package openshiftmachineapi

import (
	"reflect"
	clusterclientset "github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/config"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

const (
	ProviderName = "openshift-machine-api"
)

var _ cloudprovider.CloudProvider = (*provider)(nil)

type provider struct {
	controller	*machineController
	providerName	string
	resourceLimiter	*cloudprovider.ResourceLimiter
}

func (p *provider) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return p.providerName
}
func (p *provider) GetResourceLimiter() (*cloudprovider.ResourceLimiter, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return p.resourceLimiter, nil
}
func (p *provider) NodeGroups() []cloudprovider.NodeGroup {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var result []cloudprovider.NodeGroup
	nodegroups, err := p.controller.nodeGroups()
	if err != nil {
		klog.Errorf("error getting node groups: %v", err)
		return nil
	}
	for _, ng := range nodegroups {
		klog.V(4).Infof("discovered node group: %s", ng.Debug())
		result = append(result, ng)
	}
	return result
}
func (p *provider) NodeGroupForNode(node *apiv1.Node) (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ng, err := p.controller.nodeGroupForNode(node)
	if err != nil {
		return nil, err
	}
	if ng == nil || reflect.ValueOf(ng).IsNil() {
		return nil, nil
	}
	return ng, nil
}
func (*provider) Pricing() (cloudprovider.PricingModel, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, cloudprovider.ErrNotImplemented
}
func (*provider) GetAvailableMachineTypes() ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return []string{}, nil
}
func (*provider) NewNodeGroup(machineType string, labels map[string]string, systemLabels map[string]string, taints []apiv1.Taint, extraResources map[string]resource.Quantity) (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, cloudprovider.ErrNotImplemented
}
func (*provider) Cleanup() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (p *provider) Refresh() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (p *provider) GetInstanceID(node *apiv1.Node) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return node.Spec.ProviderID
}
func newProvider(name string, rl *cloudprovider.ResourceLimiter, controller *machineController) (cloudprovider.CloudProvider, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &provider{providerName: name, resourceLimiter: rl, controller: controller}, nil
}
func BuildOpenShiftMachineAPI(opts config.AutoscalingOptions, do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var err error
	var externalConfig *rest.Config
	externalConfig, err = rest.InClusterConfig()
	if err != nil && err != rest.ErrNotInCluster {
		klog.Fatal(err)
	}
	if opts.KubeConfigPath != "" {
		externalConfig, err = clientcmd.BuildConfigFromFlags("", opts.KubeConfigPath)
		if err != nil {
			klog.Fatalf("cannot build config: %v", err)
		}
	}
	kubeclient, err := kubernetes.NewForConfig(externalConfig)
	if err != nil {
		klog.Fatalf("create kube clientset failed: %v", err)
	}
	clusterclient, err := clusterclientset.NewForConfig(externalConfig)
	if err != nil {
		klog.Fatalf("create cluster clientset failed: %v", err)
	}
	enableMachineDeployments := false
	controller, err := newMachineController(kubeclient, clusterclient, enableMachineDeployments)
	if err != nil {
		klog.Fatal(err)
	}
	stopCh := make(chan struct{})
	if err := controller.run(stopCh); err != nil {
		klog.Fatal(err)
	}
	provider, err := newProvider(ProviderName, rl, controller)
	if err != nil {
		klog.Fatal(err)
	}
	return provider
}
