package gke

import (
	"fmt"
	"io"
	"os"
	"time"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/gce"
	"k8s.io/autoscaler/cluster-autoscaler/config"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	"k8s.io/autoscaler/cluster-autoscaler/utils/gpu"
	"k8s.io/klog"
	kubeletapis "k8s.io/kubernetes/pkg/kubelet/apis"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

const (
	ProviderNameGKE = "gke"
)
const (
	maxAutoprovisionedSize	= 1000
	minAutoprovisionedSize	= 0
)

var autoprovisionedMachineTypes = []string{"n1-standard-1", "n1-standard-2", "n1-standard-4", "n1-standard-8", "n1-standard-16", "n1-highcpu-2", "n1-highcpu-4", "n1-highcpu-8", "n1-highcpu-16", "n1-highmem-2", "n1-highmem-4", "n1-highmem-8", "n1-highmem-16"}

type GkeCloudProvider struct {
	gkeManager			GkeManager
	resourceLimiterFromFlags	*cloudprovider.ResourceLimiter
}

func BuildGkeCloudProvider(gkeManager GkeManager, resourceLimiter *cloudprovider.ResourceLimiter) (*GkeCloudProvider, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &GkeCloudProvider{gkeManager: gkeManager, resourceLimiterFromFlags: resourceLimiter}, nil
}
func (gke *GkeCloudProvider) Cleanup() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	gke.gkeManager.Cleanup()
	return nil
}
func (gke *GkeCloudProvider) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ProviderNameGKE
}
func (gke *GkeCloudProvider) NodeGroups() []cloudprovider.NodeGroup {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	migs := gke.gkeManager.GetMigs()
	result := make([]cloudprovider.NodeGroup, 0, len(migs))
	for _, mig := range migs {
		result = append(result, mig.Config)
	}
	return result
}
func (gke *GkeCloudProvider) NodeGroupForNode(node *apiv1.Node) (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ref, err := gce.GceRefFromProviderId(node.Spec.ProviderID)
	if err != nil {
		return nil, err
	}
	mig, err := gke.gkeManager.GetMigForInstance(ref)
	return mig, err
}
func (gke *GkeCloudProvider) Pricing() (cloudprovider.PricingModel, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &gce.GcePriceModel{}, nil
}
func (gke *GkeCloudProvider) GetAvailableMachineTypes() ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return autoprovisionedMachineTypes, nil
}
func (gke *GkeCloudProvider) NewNodeGroup(machineType string, labels map[string]string, systemLabels map[string]string, taints []apiv1.Taint, extraResources map[string]resource.Quantity) (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	nodePoolName := fmt.Sprintf("%s-%s-%d", nodeAutoprovisioningPrefix, machineType, time.Now().Unix())
	zone, found := systemLabels[kubeletapis.LabelZoneFailureDomain]
	if !found {
		return nil, cloudprovider.ErrIllegalConfiguration
	}
	if gpuRequest, found := extraResources[gpu.ResourceNvidiaGPU]; found {
		gpuType, found := systemLabels[gpu.GPULabel]
		if !found {
			return nil, cloudprovider.ErrIllegalConfiguration
		}
		gpuCount, err := getNormalizedGpuCount(gpuRequest.Value())
		if err != nil {
			return nil, err
		}
		extraResources[gpu.ResourceNvidiaGPU] = *resource.NewQuantity(gpuCount, resource.DecimalSI)
		err = validateGpuConfig(gpuType, gpuCount, zone, machineType)
		if err != nil {
			return nil, err
		}
		nodePoolName = fmt.Sprintf("%s-%s-gpu-%d", nodeAutoprovisioningPrefix, machineType, time.Now().Unix())
		labels[gpu.GPULabel] = gpuType
		taint := apiv1.Taint{Effect: apiv1.TaintEffectNoSchedule, Key: gpu.ResourceNvidiaGPU, Value: "present"}
		taints = append(taints, taint)
	}
	mig := &GkeMig{gceRef: gce.GceRef{Project: gke.gkeManager.GetProjectId(), Zone: zone, Name: nodePoolName + "-temporary-mig"}, gkeManager: gke.gkeManager, autoprovisioned: true, exist: false, nodePoolName: nodePoolName, minSize: minAutoprovisionedSize, maxSize: maxAutoprovisionedSize, spec: &MigSpec{MachineType: machineType, Labels: labels, Taints: taints, ExtraResources: extraResources}}
	if _, err := gke.gkeManager.GetMigTemplateNode(mig); err != nil {
		return nil, fmt.Errorf("Failed to build node from spec: %v", err)
	}
	return mig, nil
}
func (gke *GkeCloudProvider) GetResourceLimiter() (*cloudprovider.ResourceLimiter, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	resourceLimiter, err := gke.gkeManager.GetResourceLimiter()
	if err != nil {
		return nil, err
	}
	if resourceLimiter != nil {
		return resourceLimiter, nil
	}
	return gke.resourceLimiterFromFlags, nil
}
func (gke *GkeCloudProvider) Refresh() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gke.gkeManager.Refresh()
}
func (gke *GkeCloudProvider) GetClusterInfo() (projectId, location, clusterName string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gke.gkeManager.GetProjectId(), gke.gkeManager.GetLocation(), gke.gkeManager.GetClusterName()
}
func (gke *GkeCloudProvider) GetNodeLocations() []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gke.gkeManager.GetNodeLocations()
}
func (gke *GkeCloudProvider) GetInstanceID(node *apiv1.Node) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return node.Spec.ProviderID
}

type MigSpec struct {
	MachineType	string
	Labels		map[string]string
	Taints		[]apiv1.Taint
	ExtraResources	map[string]resource.Quantity
}
type GkeMig struct {
	gceRef		gce.GceRef
	gkeManager	GkeManager
	minSize		int
	maxSize		int
	autoprovisioned	bool
	exist		bool
	nodePoolName	string
	spec		*MigSpec
}

func (mig *GkeMig) GceRef() gce.GceRef {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return mig.gceRef
}
func (mig *GkeMig) NodePoolName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return mig.nodePoolName
}
func (mig *GkeMig) Spec() *MigSpec {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return mig.spec
}
func (mig *GkeMig) MaxSize() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return mig.maxSize
}
func (mig *GkeMig) MinSize() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return mig.minSize
}
func (mig *GkeMig) TargetSize() (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !mig.exist {
		return 0, nil
	}
	size, err := mig.gkeManager.GetMigSize(mig)
	return int(size), err
}
func (mig *GkeMig) IncreaseSize(delta int) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if delta <= 0 {
		return fmt.Errorf("size increase must be positive")
	}
	size, err := mig.gkeManager.GetMigSize(mig)
	if err != nil {
		return err
	}
	if int(size)+delta > mig.MaxSize() {
		return fmt.Errorf("size increase too large - desired:%d max:%d", int(size)+delta, mig.MaxSize())
	}
	return mig.gkeManager.SetMigSize(mig, size+int64(delta))
}
func (mig *GkeMig) DecreaseTargetSize(delta int) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if delta >= 0 {
		return fmt.Errorf("size decrease must be negative")
	}
	size, err := mig.gkeManager.GetMigSize(mig)
	if err != nil {
		return err
	}
	nodes, err := mig.gkeManager.GetMigNodes(mig)
	if err != nil {
		return err
	}
	if int(size)+delta < len(nodes) {
		return fmt.Errorf("attempt to delete existing nodes targetSize:%d delta:%d existingNodes: %d", size, delta, len(nodes))
	}
	return mig.gkeManager.SetMigSize(mig, size+int64(delta))
}
func (mig *GkeMig) Belongs(node *apiv1.Node) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ref, err := gce.GceRefFromProviderId(node.Spec.ProviderID)
	if err != nil {
		return false, err
	}
	targetMig, err := mig.gkeManager.GetMigForInstance(ref)
	if err != nil {
		return false, err
	}
	if targetMig == nil {
		return false, fmt.Errorf("%s doesn't belong to a known mig", node.Name)
	}
	if targetMig.Id() != mig.Id() {
		return false, nil
	}
	return true, nil
}
func (mig *GkeMig) DeleteNodes(nodes []*apiv1.Node) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	size, err := mig.gkeManager.GetMigSize(mig)
	if err != nil {
		return err
	}
	if int(size) <= mig.MinSize() {
		return fmt.Errorf("min size reached, nodes will not be deleted")
	}
	refs := make([]*gce.GceRef, 0, len(nodes))
	for _, node := range nodes {
		belongs, err := mig.Belongs(node)
		if err != nil {
			return err
		}
		if !belongs {
			return fmt.Errorf("%s belong to a different mig than %s", node.Name, mig.Id())
		}
		gceref, err := gce.GceRefFromProviderId(node.Spec.ProviderID)
		if err != nil {
			return err
		}
		refs = append(refs, gceref)
	}
	return mig.gkeManager.DeleteInstances(refs)
}
func (mig *GkeMig) Id() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gce.GenerateMigUrl(mig.gceRef)
}
func (mig *GkeMig) Debug() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%s (%d:%d)", mig.Id(), mig.MinSize(), mig.MaxSize())
}
func (mig *GkeMig) Nodes() ([]cloudprovider.Instance, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	instanceNames, err := mig.gkeManager.GetMigNodes(mig)
	if err != nil {
		return nil, err
	}
	instances := make([]cloudprovider.Instance, 0, len(instanceNames))
	for _, instanceName := range instanceNames {
		instances = append(instances, cloudprovider.Instance{Id: instanceName})
	}
	return instances, nil
}
func (mig *GkeMig) Exist() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return mig.exist
}
func (mig *GkeMig) Create() (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !mig.exist && mig.autoprovisioned {
		return mig.gkeManager.CreateNodePool(mig)
	}
	return nil, fmt.Errorf("Cannot create non-autoprovisioned node group")
}
func (mig *GkeMig) Delete() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if mig.exist && mig.autoprovisioned {
		return mig.gkeManager.DeleteNodePool(mig)
	}
	return fmt.Errorf("Cannot delete non-autoprovisioned node group")
}
func (mig *GkeMig) Autoprovisioned() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return mig.autoprovisioned
}
func (mig *GkeMig) TemplateNodeInfo() (*schedulercache.NodeInfo, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	node, err := mig.gkeManager.GetMigTemplateNode(mig)
	if err != nil {
		return nil, err
	}
	nodeInfo := schedulercache.NewNodeInfo(cloudprovider.BuildKubeProxy(mig.Id()))
	nodeInfo.SetNode(node)
	return nodeInfo, nil
}
func BuildGKE(opts config.AutoscalingOptions, do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if do.DiscoverySpecified() {
		klog.Fatal("GKE gets nodegroup specification via API, command line specs are not allowed")
	}
	var config io.ReadCloser
	if opts.CloudConfig != "" {
		var err error
		config, err = os.Open(opts.CloudConfig)
		if err != nil {
			klog.Fatalf("Couldn't open cloud provider configuration %s: %#v", opts.CloudConfig, err)
		}
		defer config.Close()
	}
	mode := ModeGKE
	if opts.NodeAutoprovisioningEnabled {
		mode = ModeGKENAP
	}
	manager, err := CreateGkeManager(config, mode, opts.ClusterName, opts.Regional)
	if err != nil {
		klog.Fatalf("Failed to create GKE Manager: %v", err)
	}
	provider, err := BuildGkeCloudProvider(manager, rl)
	if err != nil {
		klog.Fatalf("Failed to create GKE cloud provider: %v", err)
	}
	registerMetrics()
	gce.RegisterMetrics()
	return provider
}
