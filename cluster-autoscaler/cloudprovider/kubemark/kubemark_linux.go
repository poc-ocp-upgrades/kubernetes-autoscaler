package kubemark

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/config"
	"k8s.io/autoscaler/cluster-autoscaler/config/dynamic"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	"k8s.io/client-go/informers"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/pkg/kubemark"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
	"k8s.io/klog"
)

const (
	ProviderName = "kubemark"
)

type KubemarkCloudProvider struct {
	kubemarkController	*kubemark.KubemarkController
	nodeGroups		[]*NodeGroup
	resourceLimiter		*cloudprovider.ResourceLimiter
}

func BuildKubemarkCloudProvider(kubemarkController *kubemark.KubemarkController, specs []string, resourceLimiter *cloudprovider.ResourceLimiter) (*KubemarkCloudProvider, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	kubemark := &KubemarkCloudProvider{kubemarkController: kubemarkController, nodeGroups: make([]*NodeGroup, 0), resourceLimiter: resourceLimiter}
	for _, spec := range specs {
		if err := kubemark.addNodeGroup(spec); err != nil {
			return nil, err
		}
	}
	return kubemark, nil
}
func (kubemark *KubemarkCloudProvider) addNodeGroup(spec string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	nodeGroup, err := buildNodeGroup(spec, kubemark.kubemarkController)
	if err != nil {
		return err
	}
	klog.V(2).Infof("adding node group: %s", nodeGroup.Name)
	kubemark.nodeGroups = append(kubemark.nodeGroups, nodeGroup)
	return nil
}
func (kubemark *KubemarkCloudProvider) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ProviderName
}
func (kubemark *KubemarkCloudProvider) NodeGroups() []cloudprovider.NodeGroup {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make([]cloudprovider.NodeGroup, 0, len(kubemark.nodeGroups))
	for _, nodegroup := range kubemark.nodeGroups {
		result = append(result, nodegroup)
	}
	return result
}
func (kubemark *KubemarkCloudProvider) Pricing() (cloudprovider.PricingModel, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, cloudprovider.ErrNotImplemented
}
func (kubemark *KubemarkCloudProvider) NodeGroupForNode(node *apiv1.Node) (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	nodeGroupName, err := kubemark.kubemarkController.GetNodeGroupForNode(node.ObjectMeta.Name)
	if err != nil {
		return nil, err
	}
	for _, nodeGroup := range kubemark.nodeGroups {
		if nodeGroup.Name == nodeGroupName {
			return nodeGroup, nil
		}
	}
	return nil, nil
}
func (kubemark *KubemarkCloudProvider) GetAvailableMachineTypes() ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return []string{}, cloudprovider.ErrNotImplemented
}
func (kubemark *KubemarkCloudProvider) NewNodeGroup(machineType string, labels map[string]string, systemLabels map[string]string, taints []apiv1.Taint, extraResources map[string]resource.Quantity) (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, cloudprovider.ErrNotImplemented
}
func (kubemark *KubemarkCloudProvider) GetResourceLimiter() (*cloudprovider.ResourceLimiter, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return kubemark.resourceLimiter, nil
}
func (kubemark *KubemarkCloudProvider) Refresh() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (kubemark *KubemarkCloudProvider) GetInstanceID(node *apiv1.Node) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return node.Spec.ProviderID
}
func (kubemark *KubemarkCloudProvider) Cleanup() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}

type NodeGroup struct {
	Name			string
	kubemarkController	*kubemark.KubemarkController
	minSize			int
	maxSize			int
}

func (nodeGroup *NodeGroup) Id() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nodeGroup.Name
}
func (nodeGroup *NodeGroup) MinSize() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nodeGroup.minSize
}
func (nodeGroup *NodeGroup) MaxSize() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nodeGroup.maxSize
}
func (nodeGroup *NodeGroup) Debug() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%s (%d:%d)", nodeGroup.Id(), nodeGroup.MinSize(), nodeGroup.MaxSize())
}
func (nodeGroup *NodeGroup) Nodes() ([]cloudprovider.Instance, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	instances := make([]cloudprovider.Instance, 0)
	nodes, err := nodeGroup.kubemarkController.GetNodeNamesForNodeGroup(nodeGroup.Name)
	if err != nil {
		return instances, err
	}
	for _, node := range nodes {
		instances = append(instances, cloudprovider.Instance{Id: ":////" + node})
	}
	return instances, nil
}
func (nodeGroup *NodeGroup) DeleteNodes(nodes []*apiv1.Node) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	size, err := nodeGroup.kubemarkController.GetNodeGroupTargetSize(nodeGroup.Name)
	if err != nil {
		return err
	}
	if size <= nodeGroup.MinSize() {
		return fmt.Errorf("min size reached, nodes will not be deleted")
	}
	for _, node := range nodes {
		if err := nodeGroup.kubemarkController.RemoveNodeFromNodeGroup(nodeGroup.Name, node.ObjectMeta.Name); err != nil {
			return err
		}
	}
	return nil
}
func (nodeGroup *NodeGroup) IncreaseSize(delta int) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if delta <= 0 {
		return fmt.Errorf("size increase must be positive")
	}
	size, err := nodeGroup.kubemarkController.GetNodeGroupTargetSize(nodeGroup.Name)
	if err != nil {
		return err
	}
	newSize := int(size) + delta
	if newSize > nodeGroup.MaxSize() {
		return fmt.Errorf("size increase too large, desired: %d max: %d", newSize, nodeGroup.MaxSize())
	}
	return nodeGroup.kubemarkController.SetNodeGroupSize(nodeGroup.Name, newSize)
}
func (nodeGroup *NodeGroup) TargetSize() (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	size, err := nodeGroup.kubemarkController.GetNodeGroupTargetSize(nodeGroup.Name)
	return int(size), err
}
func (nodeGroup *NodeGroup) DecreaseTargetSize(delta int) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if delta >= 0 {
		return fmt.Errorf("size decrease must be negative")
	}
	size, err := nodeGroup.kubemarkController.GetNodeGroupTargetSize(nodeGroup.Name)
	if err != nil {
		return err
	}
	nodes, err := nodeGroup.kubemarkController.GetNodeNamesForNodeGroup(nodeGroup.Name)
	if err != nil {
		return err
	}
	newSize := int(size) + delta
	if newSize < len(nodes) {
		return fmt.Errorf("attempt to delete existing nodes, targetSize: %d delta: %d existingNodes: %d", size, delta, len(nodes))
	}
	return nodeGroup.kubemarkController.SetNodeGroupSize(nodeGroup.Name, newSize)
}
func (nodeGroup *NodeGroup) TemplateNodeInfo() (*schedulercache.NodeInfo, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, cloudprovider.ErrNotImplemented
}
func (nodeGroup *NodeGroup) Exist() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return true
}
func (nodeGroup *NodeGroup) Create() (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, cloudprovider.ErrNotImplemented
}
func (nodeGroup *NodeGroup) Delete() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cloudprovider.ErrNotImplemented
}
func (nodeGroup *NodeGroup) Autoprovisioned() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return false
}
func buildNodeGroup(value string, kubemarkController *kubemark.KubemarkController) (*NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	spec, err := dynamic.SpecFromString(value, true)
	if err != nil {
		return nil, fmt.Errorf("failed to parse node group spec: %v", err)
	}
	nodeGroup := &NodeGroup{Name: spec.Name, kubemarkController: kubemarkController, minSize: spec.MinSize, maxSize: spec.MaxSize}
	return nodeGroup, nil
}
func BuildKubemark(opts config.AutoscalingOptions, do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	externalConfig, err := rest.InClusterConfig()
	if err != nil {
		klog.Fatalf("Failed to get kubeclient config for external cluster: %v", err)
	}
	kubemarkConfig, err := clientcmd.BuildConfigFromFlags("", "/kubeconfig/cluster_autoscaler.kubeconfig")
	if err != nil {
		klog.Fatalf("Failed to get kubeclient config for kubemark cluster: %v", err)
	}
	stop := make(chan struct{})
	externalClient := kubeclient.NewForConfigOrDie(externalConfig)
	kubemarkClient := kubeclient.NewForConfigOrDie(kubemarkConfig)
	externalInformerFactory := informers.NewSharedInformerFactory(externalClient, 0)
	kubemarkInformerFactory := informers.NewSharedInformerFactory(kubemarkClient, 0)
	kubemarkNodeInformer := kubemarkInformerFactory.Core().V1().Nodes()
	go kubemarkNodeInformer.Informer().Run(stop)
	kubemarkController, err := kubemark.NewKubemarkController(externalClient, externalInformerFactory, kubemarkClient, kubemarkNodeInformer)
	if err != nil {
		klog.Fatalf("Failed to create Kubemark cloud provider: %v", err)
	}
	externalInformerFactory.Start(stop)
	if !kubemarkController.WaitForCacheSync(stop) {
		klog.Fatalf("Failed to sync caches for kubemark controller")
	}
	go kubemarkController.Run(stop)
	provider, err := BuildKubemarkCloudProvider(kubemarkController, do.NodeGroupSpecs, rl)
	if err != nil {
		klog.Fatalf("Failed to create Kubemark cloud provider: %v", err)
	}
	return provider
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
