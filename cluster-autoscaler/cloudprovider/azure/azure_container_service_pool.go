package azure

import (
 "fmt"
 "strings"
 "github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2018-03-31/containerservice"
 "k8s.io/klog"
 apiv1 "k8s.io/api/core/v1"
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
 "k8s.io/autoscaler/cluster-autoscaler/config/dynamic"
 schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

type ContainerServiceAgentPool struct {
 azureRef
 manager           *AzureManager
 util              *AzUtil
 minSize           int
 maxSize           int
 serviceType       string
 clusterName       string
 resourceGroup     string
 nodeResourceGroup string
}

func NewContainerServiceAgentPool(spec *dynamic.NodeGroupSpec, am *AzureManager) (*ContainerServiceAgentPool, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 asg := &ContainerServiceAgentPool{azureRef: azureRef{Name: spec.Name}, minSize: spec.MinSize, maxSize: spec.MaxSize, manager: am}
 asg.util = &AzUtil{manager: am}
 asg.serviceType = am.config.VMType
 asg.clusterName = am.config.ClusterName
 asg.resourceGroup = am.config.ResourceGroup
 if am.config.VMType == vmTypeAKS {
  asg.nodeResourceGroup = am.config.NodeResourceGroup
 } else {
  asg.nodeResourceGroup = am.config.ResourceGroup
 }
 return asg, nil
}
func (agentPool *ContainerServiceAgentPool) GetAKSAgentPool(agentProfiles *[]containerservice.ManagedClusterAgentPoolProfile) (ret *containerservice.ManagedClusterAgentPoolProfile) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 for _, value := range *agentProfiles {
  profileName := *value.Name
  klog.V(5).Infof("AKS AgentPool profile name: %s", profileName)
  if strings.EqualFold(profileName, agentPool.azureRef.Name) {
   return &value
  }
 }
 return nil
}
func (agentPool *ContainerServiceAgentPool) GetACSAgentPool(agentProfiles *[]containerservice.AgentPoolProfile) (ret *containerservice.AgentPoolProfile) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 for _, value := range *agentProfiles {
  profileName := *value.Name
  klog.V(5).Infof("ACS AgentPool profile name: %s", profileName)
  if strings.EqualFold(profileName, agentPool.azureRef.Name) {
   return &value
  }
 }
 for _, value := range *agentProfiles {
  profileName := *value.Name
  poolName := agentPool.azureRef.Name + "pool0"
  klog.V(5).Infof("Workaround match check - ACS AgentPool Profile: %s <=> Poolname: %s", profileName, poolName)
  if strings.EqualFold(profileName, poolName) {
   return &value
  }
 }
 return nil
}
func (agentPool *ContainerServiceAgentPool) getAKSNodeCount() (count int, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 ctx, cancel := getContextWithCancel()
 defer cancel()
 managedCluster, err := agentPool.manager.azClient.managedContainerServicesClient.Get(ctx, agentPool.resourceGroup, agentPool.clusterName)
 if err != nil {
  klog.Errorf("Failed to get AKS cluster (name:%q): %v", agentPool.clusterName, err)
  return -1, err
 }
 pool := agentPool.GetAKSAgentPool(managedCluster.AgentPoolProfiles)
 if pool == nil {
  return -1, fmt.Errorf("could not find pool with name: %s", agentPool.azureRef)
 }
 if pool.Count != nil {
  return int(*pool.Count), nil
 }
 return 0, nil
}
func (agentPool *ContainerServiceAgentPool) getACSNodeCount() (count int, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 ctx, cancel := getContextWithCancel()
 defer cancel()
 acsCluster, err := agentPool.manager.azClient.containerServicesClient.Get(ctx, agentPool.resourceGroup, agentPool.clusterName)
 if err != nil {
  klog.Errorf("Failed to get ACS cluster (name:%q): %v", agentPool.clusterName, err)
  return -1, err
 }
 pool := agentPool.GetACSAgentPool(acsCluster.AgentPoolProfiles)
 if pool == nil {
  return -1, fmt.Errorf("could not find pool with name: %s", agentPool.azureRef)
 }
 if pool.Count != nil {
  return int(*pool.Count), nil
 }
 return 0, nil
}
func (agentPool *ContainerServiceAgentPool) setAKSNodeCount(count int) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 ctx, cancel := getContextWithCancel()
 defer cancel()
 managedCluster, err := agentPool.manager.azClient.managedContainerServicesClient.Get(ctx, agentPool.resourceGroup, agentPool.clusterName)
 if err != nil {
  klog.Errorf("Failed to get AKS cluster (name:%q): %v", agentPool.clusterName, err)
  return err
 }
 pool := agentPool.GetAKSAgentPool(managedCluster.AgentPoolProfiles)
 if pool == nil {
  return fmt.Errorf("could not find pool with name: %s", agentPool.azureRef)
 }
 klog.Infof("Current size: %d, Target size requested: %d", *pool.Count, count)
 updateCtx, updateCancel := getContextWithCancel()
 defer updateCancel()
 *pool.Count = int32(count)
 aksClient := agentPool.manager.azClient.managedContainerServicesClient
 future, err := aksClient.CreateOrUpdate(updateCtx, agentPool.resourceGroup, agentPool.clusterName, managedCluster)
 if err != nil {
  klog.Errorf("Failed to update AKS cluster (%q): %v", agentPool.clusterName, err)
  return err
 }
 err = future.WaitForCompletionRef(updateCtx, aksClient.Client)
 isSuccess, realError := isSuccessHTTPResponse(future.Response(), err)
 if isSuccess {
  klog.V(3).Infof("aksClient.CreateOrUpdate for aks cluster %q success", agentPool.clusterName)
  return nil
 }
 klog.Errorf("aksClient.CreateOrUpdate for aks cluster %q failed: %v", agentPool.clusterName, realError)
 return realError
}
func (agentPool *ContainerServiceAgentPool) setACSNodeCount(count int) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 ctx, cancel := getContextWithCancel()
 defer cancel()
 acsCluster, err := agentPool.manager.azClient.containerServicesClient.Get(ctx, agentPool.resourceGroup, agentPool.clusterName)
 if err != nil {
  klog.Errorf("Failed to get ACS cluster (name:%q): %v", agentPool.clusterName, err)
  return err
 }
 pool := agentPool.GetACSAgentPool(acsCluster.AgentPoolProfiles)
 if pool == nil {
  return fmt.Errorf("could not find pool with name: %s", agentPool.azureRef)
 }
 klog.Infof("Current size: %d, Target size requested: %d", *pool.Count, count)
 updateCtx, updateCancel := getContextWithCancel()
 defer updateCancel()
 *pool.Count = int32(count)
 acsClient := agentPool.manager.azClient.containerServicesClient
 future, err := acsClient.CreateOrUpdate(updateCtx, agentPool.resourceGroup, agentPool.clusterName, acsCluster)
 if err != nil {
  klog.Errorf("Failed to update ACS cluster (%q): %v", agentPool.clusterName, err)
  return err
 }
 err = future.WaitForCompletionRef(updateCtx, acsClient.Client)
 isSuccess, realError := isSuccessHTTPResponse(future.Response(), err)
 if isSuccess {
  klog.V(3).Infof("acsClient.CreateOrUpdate for acs cluster %q success", agentPool.clusterName)
  return nil
 }
 klog.Errorf("acsClient.CreateOrUpdate for acs cluster %q failed: %v", agentPool.clusterName, realError)
 return realError
}
func (agentPool *ContainerServiceAgentPool) GetNodeCount() (count int, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if agentPool.serviceType == vmTypeAKS {
  return agentPool.getAKSNodeCount()
 }
 return agentPool.getACSNodeCount()
}
func (agentPool *ContainerServiceAgentPool) SetNodeCount(count int) (err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if agentPool.serviceType == vmTypeAKS {
  return agentPool.setAKSNodeCount(count)
 }
 return agentPool.setACSNodeCount(count)
}
func (agentPool *ContainerServiceAgentPool) GetProviderID(name string) string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return "azure://" + strings.ToLower(name)
}
func (agentPool *ContainerServiceAgentPool) GetName(providerID string) (string, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 providerID = strings.TrimPrefix(providerID, "azure://")
 ctx, cancel := getContextWithCancel()
 defer cancel()
 vms, err := agentPool.manager.azClient.virtualMachinesClient.List(ctx, agentPool.nodeResourceGroup)
 if err != nil {
  return "", err
 }
 for _, vm := range vms {
  if strings.EqualFold(*vm.ID, providerID) {
   return *vm.Name, nil
  }
 }
 return "", fmt.Errorf("VM list empty")
}
func (agentPool *ContainerServiceAgentPool) MaxSize() int {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return agentPool.maxSize
}
func (agentPool *ContainerServiceAgentPool) MinSize() int {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return agentPool.minSize
}
func (agentPool *ContainerServiceAgentPool) TargetSize() (int, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return agentPool.GetNodeCount()
}
func (agentPool *ContainerServiceAgentPool) SetSize(targetSize int) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if targetSize > agentPool.MaxSize() || targetSize < agentPool.MinSize() {
  klog.Errorf("Target size %d requested outside Max: %d, Min: %d", targetSize, agentPool.MaxSize(), agentPool.MaxSize())
  return fmt.Errorf("Target size %d requested outside Max: %d, Min: %d", targetSize, agentPool.MaxSize(), agentPool.MinSize())
 }
 klog.V(2).Infof("Setting size for cluster (%q) with new count (%d)", agentPool.clusterName, targetSize)
 if agentPool.serviceType == vmTypeAKS {
  return agentPool.setAKSNodeCount(targetSize)
 }
 return agentPool.setACSNodeCount(targetSize)
}
func (agentPool *ContainerServiceAgentPool) IncreaseSize(delta int) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if delta <= 0 {
  return fmt.Errorf("Size increase must be +ve")
 }
 currentSize, err := agentPool.TargetSize()
 if err != nil {
  return err
 }
 targetSize := int(currentSize) + delta
 if targetSize > agentPool.MaxSize() {
  return fmt.Errorf("Size increase request of %d more than max size %d set", targetSize, agentPool.MaxSize())
 }
 return agentPool.SetSize(targetSize)
}
func (agentPool *ContainerServiceAgentPool) DeleteNodesInternal(providerIDs []string) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 currentSize, err := agentPool.TargetSize()
 if err != nil {
  return err
 }
 targetSize := currentSize
 for _, providerID := range providerIDs {
  klog.Infof("ProviderID got to delete: %s", providerID)
  nodeName, err := agentPool.GetName(providerID)
  if err != nil {
   return err
  }
  klog.Infof("VM name got to delete: %s", nodeName)
  err = agentPool.util.DeleteVirtualMachine(agentPool.nodeResourceGroup, nodeName)
  if err != nil {
   klog.Error(err)
   return err
  }
  targetSize--
 }
 if currentSize != targetSize {
  agentPool.SetSize(targetSize)
 }
 return nil
}
func (agentPool *ContainerServiceAgentPool) DeleteNodes(nodes []*apiv1.Node) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 var providerIDs []string
 for _, node := range nodes {
  klog.Infof("Node: %s", node.Spec.ProviderID)
  providerIDs = append(providerIDs, node.Spec.ProviderID)
 }
 for _, p := range providerIDs {
  klog.Infof("ProviderID before calling acsmgr: %s", p)
 }
 return agentPool.DeleteNodesInternal(providerIDs)
}
func (agentPool *ContainerServiceAgentPool) IsContainerServiceNode(tags map[string]*string) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 poolName := tags["poolName"]
 if poolName != nil {
  klog.V(5).Infof("Matching agentPool name: %s with tag name: %s", agentPool.azureRef.Name, *poolName)
  if strings.EqualFold(*poolName, agentPool.azureRef.Name) {
   return true
  }
 }
 return false
}
func (agentPool *ContainerServiceAgentPool) GetNodes() ([]string, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 ctx, cancel := getContextWithCancel()
 defer cancel()
 vmList, err := agentPool.manager.azClient.virtualMachinesClient.List(ctx, agentPool.nodeResourceGroup)
 if err != nil {
  klog.Errorf("Azure client list vm error : %v", err)
  return nil, err
 }
 var nodeArray []string
 for _, node := range vmList {
  klog.V(5).Infof("Node Name: %s, ID: %s", *node.Name, *node.ID)
  if agentPool.IsContainerServiceNode(node.Tags) {
   providerID := agentPool.GetProviderID(*node.ID)
   klog.V(5).Infof("Returning back the providerID: %s", providerID)
   nodeArray = append(nodeArray, providerID)
  }
 }
 return nodeArray, nil
}
func (agentPool *ContainerServiceAgentPool) DecreaseTargetSize(delta int) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if delta >= 0 {
  klog.Errorf("Size decrease error: %d", delta)
  return fmt.Errorf("Size decrease must be negative")
 }
 currentSize, err := agentPool.TargetSize()
 if err != nil {
  klog.Error(err)
  return err
 }
 nodes, err := agentPool.GetNodes()
 if err != nil {
  klog.Error(err)
  return err
 }
 targetSize := int(currentSize) + delta
 if targetSize < len(nodes) {
  return fmt.Errorf("attempt to delete existing nodes targetSize:%d delta:%d existingNodes: %d", currentSize, delta, len(nodes))
 }
 return agentPool.SetSize(targetSize)
}
func (agentPool *ContainerServiceAgentPool) Id() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return agentPool.azureRef.Name
}
func (agentPool *ContainerServiceAgentPool) Debug() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return fmt.Sprintf("%s (%d:%d)", agentPool.Id(), agentPool.MinSize(), agentPool.MaxSize())
}
func (agentPool *ContainerServiceAgentPool) Nodes() ([]cloudprovider.Instance, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 instanceNames, err := agentPool.GetNodes()
 if err != nil {
  return nil, err
 }
 instances := make([]cloudprovider.Instance, 0, len(instanceNames))
 for _, instanceName := range instanceNames {
  instances = append(instances, cloudprovider.Instance{Id: instanceName})
 }
 return instances, nil
}
func (agentPool *ContainerServiceAgentPool) TemplateNodeInfo() (*schedulercache.NodeInfo, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return nil, cloudprovider.ErrNotImplemented
}
func (agentPool *ContainerServiceAgentPool) Exist() bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return true
}
func (agentPool *ContainerServiceAgentPool) Create() (cloudprovider.NodeGroup, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return nil, cloudprovider.ErrAlreadyExist
}
func (agentPool *ContainerServiceAgentPool) Delete() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return cloudprovider.ErrNotImplemented
}
func (agentPool *ContainerServiceAgentPool) Autoprovisioned() bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return false
}
