package gce

import (
 "fmt"
 "io"
 "os"
 "strings"
 apiv1 "k8s.io/api/core/v1"
 "k8s.io/apimachinery/pkg/api/resource"
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
 "k8s.io/autoscaler/cluster-autoscaler/config"
 "k8s.io/autoscaler/cluster-autoscaler/utils/errors"
 "k8s.io/klog"
 schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

const (
 ProviderNameGCE = "gce"
)

type GceCloudProvider struct {
 gceManager               GceManager
 resourceLimiterFromFlags *cloudprovider.ResourceLimiter
}

func BuildGceCloudProvider(gceManager GceManager, resourceLimiter *cloudprovider.ResourceLimiter) (*GceCloudProvider, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &GceCloudProvider{gceManager: gceManager, resourceLimiterFromFlags: resourceLimiter}, nil
}
func (gce *GceCloudProvider) Cleanup() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 gce.gceManager.Cleanup()
 return nil
}
func (gce *GceCloudProvider) Name() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return ProviderNameGCE
}
func (gce *GceCloudProvider) NodeGroups() []cloudprovider.NodeGroup {
 _logClusterCodePath()
 defer _logClusterCodePath()
 migs := gce.gceManager.GetMigs()
 result := make([]cloudprovider.NodeGroup, 0, len(migs))
 for _, mig := range migs {
  result = append(result, mig.Config)
 }
 return result
}
func (gce *GceCloudProvider) NodeGroupForNode(node *apiv1.Node) (cloudprovider.NodeGroup, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 ref, err := GceRefFromProviderId(node.Spec.ProviderID)
 if err != nil {
  return nil, err
 }
 mig, err := gce.gceManager.GetMigForInstance(ref)
 return mig, err
}
func (gce *GceCloudProvider) Pricing() (cloudprovider.PricingModel, errors.AutoscalerError) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &GcePriceModel{}, nil
}
func (gce *GceCloudProvider) GetAvailableMachineTypes() ([]string, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return []string{}, nil
}
func (gce *GceCloudProvider) NewNodeGroup(machineType string, labels map[string]string, systemLabels map[string]string, taints []apiv1.Taint, extraResources map[string]resource.Quantity) (cloudprovider.NodeGroup, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return nil, cloudprovider.ErrNotImplemented
}
func (gce *GceCloudProvider) GetResourceLimiter() (*cloudprovider.ResourceLimiter, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 resourceLimiter, err := gce.gceManager.GetResourceLimiter()
 if err != nil {
  return nil, err
 }
 if resourceLimiter != nil {
  return resourceLimiter, nil
 }
 return gce.resourceLimiterFromFlags, nil
}
func (gce *GceCloudProvider) Refresh() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return gce.gceManager.Refresh()
}
func (gce *GceCloudProvider) GetInstanceID(node *apiv1.Node) string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return node.Spec.ProviderID
}

type GceRef struct {
 Project string
 Zone    string
 Name    string
}

func (ref GceRef) String() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return fmt.Sprintf("%s/%s/%s", ref.Project, ref.Zone, ref.Name)
}
func GceRefFromProviderId(id string) (*GceRef, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 splitted := strings.Split(id[6:], "/")
 if len(splitted) != 3 {
  return nil, fmt.Errorf("Wrong id: expected format gce://<project-id>/<zone>/<name>, got %v", id)
 }
 return &GceRef{Project: splitted[0], Zone: splitted[1], Name: splitted[2]}, nil
}

type Mig interface {
 cloudprovider.NodeGroup
 GceRef() GceRef
}
type gceMig struct {
 gceRef     GceRef
 gceManager GceManager
 minSize    int
 maxSize    int
}

func (mig *gceMig) GceRef() GceRef {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return mig.gceRef
}
func (mig *gceMig) MaxSize() int {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return mig.maxSize
}
func (mig *gceMig) MinSize() int {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return mig.minSize
}
func (mig *gceMig) TargetSize() (int, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 size, err := mig.gceManager.GetMigSize(mig)
 return int(size), err
}
func (mig *gceMig) IncreaseSize(delta int) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if delta <= 0 {
  return fmt.Errorf("size increase must be positive")
 }
 size, err := mig.gceManager.GetMigSize(mig)
 if err != nil {
  return err
 }
 if int(size)+delta > mig.MaxSize() {
  return fmt.Errorf("size increase too large - desired:%d max:%d", int(size)+delta, mig.MaxSize())
 }
 return mig.gceManager.SetMigSize(mig, size+int64(delta))
}
func (mig *gceMig) DecreaseTargetSize(delta int) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if delta >= 0 {
  return fmt.Errorf("size decrease must be negative")
 }
 size, err := mig.gceManager.GetMigSize(mig)
 if err != nil {
  return err
 }
 nodes, err := mig.gceManager.GetMigNodes(mig)
 if err != nil {
  return err
 }
 if int(size)+delta < len(nodes) {
  return fmt.Errorf("attempt to delete existing nodes targetSize:%d delta:%d existingNodes: %d", size, delta, len(nodes))
 }
 return mig.gceManager.SetMigSize(mig, size+int64(delta))
}
func (mig *gceMig) Belongs(node *apiv1.Node) (bool, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 ref, err := GceRefFromProviderId(node.Spec.ProviderID)
 if err != nil {
  return false, err
 }
 targetMig, err := mig.gceManager.GetMigForInstance(ref)
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
func (mig *gceMig) DeleteNodes(nodes []*apiv1.Node) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 size, err := mig.gceManager.GetMigSize(mig)
 if err != nil {
  return err
 }
 if int(size) <= mig.MinSize() {
  return fmt.Errorf("min size reached, nodes will not be deleted")
 }
 refs := make([]*GceRef, 0, len(nodes))
 for _, node := range nodes {
  belongs, err := mig.Belongs(node)
  if err != nil {
   return err
  }
  if !belongs {
   return fmt.Errorf("%s belong to a different mig than %s", node.Name, mig.Id())
  }
  gceref, err := GceRefFromProviderId(node.Spec.ProviderID)
  if err != nil {
   return err
  }
  refs = append(refs, gceref)
 }
 return mig.gceManager.DeleteInstances(refs)
}
func (mig *gceMig) Id() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return GenerateMigUrl(mig.gceRef)
}
func (mig *gceMig) Debug() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return fmt.Sprintf("%s (%d:%d)", mig.Id(), mig.MinSize(), mig.MaxSize())
}
func (mig *gceMig) Nodes() ([]cloudprovider.Instance, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 instanceNames, err := mig.gceManager.GetMigNodes(mig)
 if err != nil {
  return nil, err
 }
 instances := make([]cloudprovider.Instance, 0, len(instanceNames))
 for _, instanceName := range instanceNames {
  instances = append(instances, cloudprovider.Instance{Id: instanceName})
 }
 return instances, nil
}
func (mig *gceMig) Exist() bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return true
}
func (mig *gceMig) Create() (cloudprovider.NodeGroup, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return nil, cloudprovider.ErrNotImplemented
}
func (mig *gceMig) Delete() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return cloudprovider.ErrNotImplemented
}
func (mig *gceMig) Autoprovisioned() bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return false
}
func (mig *gceMig) TemplateNodeInfo() (*schedulercache.NodeInfo, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 node, err := mig.gceManager.GetMigTemplateNode(mig)
 if err != nil {
  return nil, err
 }
 nodeInfo := schedulercache.NewNodeInfo(cloudprovider.BuildKubeProxy(mig.Id()))
 nodeInfo.SetNode(node)
 return nodeInfo, nil
}
func BuildGCE(opts config.AutoscalingOptions, do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
 _logClusterCodePath()
 defer _logClusterCodePath()
 var config io.ReadCloser
 if opts.CloudConfig != "" {
  var err error
  config, err = os.Open(opts.CloudConfig)
  if err != nil {
   klog.Fatalf("Couldn't open cloud provider configuration %s: %#v", opts.CloudConfig, err)
  }
  defer config.Close()
 }
 manager, err := CreateGceManager(config, do, opts.Regional)
 if err != nil {
  klog.Fatalf("Failed to create GCE Manager: %v", err)
 }
 provider, err := BuildGceCloudProvider(manager, rl)
 if err != nil {
  klog.Fatalf("Failed to create GCE cloud provider: %v", err)
 }
 RegisterMetrics()
 return provider
}
