package openshiftmachineapi

import (
 "fmt"
 "time"
 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
 machinev1beta1 "github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset/typed/machine/v1beta1"
 apiv1 "k8s.io/api/core/v1"
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
 schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

const (
 machineDeleteAnnotationKey = "machine.openshift.io/cluster-api-delete-machine"
 machineAnnotationKey       = "machine.openshift.io/machine"
 debugFormat                = "%s (min: %d, max: %d, replicas: %d)"
)

type nodegroup struct {
 machineapiClient  machinev1beta1.MachineV1beta1Interface
 machineController *machineController
 scalableResource  scalableResource
}

var _ cloudprovider.NodeGroup = (*nodegroup)(nil)

func (ng *nodegroup) Name() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return ng.scalableResource.Name()
}
func (ng *nodegroup) Namespace() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return ng.scalableResource.Namespace()
}
func (ng *nodegroup) MinSize() int {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return ng.scalableResource.MinSize()
}
func (ng *nodegroup) MaxSize() int {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return ng.scalableResource.MaxSize()
}
func (ng *nodegroup) TargetSize() (int, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return int(ng.scalableResource.Replicas()), nil
}
func (ng *nodegroup) IncreaseSize(delta int) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if delta <= 0 {
  return fmt.Errorf("size increase must be positive")
 }
 size := int(ng.scalableResource.Replicas())
 if size+delta > ng.MaxSize() {
  return fmt.Errorf("size increase too large - desired:%d max:%d", size+delta, ng.MaxSize())
 }
 return ng.scalableResource.SetSize(int32(size + delta))
}
func (ng *nodegroup) DeleteNodes(nodes []*apiv1.Node) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 for _, node := range nodes {
  actualNodeGroup, err := ng.machineController.nodeGroupForNode(node)
  if err != nil {
   return nil
  }
  if actualNodeGroup.Id() != ng.Id() {
   return fmt.Errorf("node %q doesn't belong to node group %q", node.Spec.ProviderID, ng.Id())
  }
  machine, err := ng.machineController.findMachineByNodeProviderID(node)
  if err != nil {
   return err
  }
  if machine == nil {
   return fmt.Errorf("unknown machine for node %q", node.Spec.ProviderID)
  }
  machine = machine.DeepCopy()
  if machine.Annotations == nil {
   machine.Annotations = map[string]string{}
  }
  machine.Annotations[machineDeleteAnnotationKey] = time.Now().String()
  if _, err := ng.machineapiClient.Machines(machine.Namespace).Update(machine); err != nil {
   return fmt.Errorf("failed to update machine %s/%s: %v", machine.Namespace, machine.Name, err)
  }
 }
 if int(ng.scalableResource.Replicas())-len(nodes) <= 0 {
  return fmt.Errorf("unable to delete %d machines in %q, machine replicas are <= 0", len(nodes), ng.Id())
 }
 return ng.scalableResource.SetSize(ng.scalableResource.Replicas() - int32(len(nodes)))
}
func (ng *nodegroup) DecreaseTargetSize(delta int) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if delta >= 0 {
  return fmt.Errorf("size decrease must be negative")
 }
 size, err := ng.TargetSize()
 if err != nil {
  return err
 }
 nodes, err := ng.Nodes()
 if err != nil {
  return err
 }
 if size+delta < len(nodes) {
  return fmt.Errorf("attempt to delete existing nodes targetSize:%d delta:%d existingNodes: %d", size, delta, len(nodes))
 }
 return ng.scalableResource.SetSize(int32(size + delta))
}
func (ng *nodegroup) Id() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return ng.scalableResource.ID()
}
func (ng *nodegroup) Debug() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return fmt.Sprintf(debugFormat, ng.Id(), ng.MinSize(), ng.MaxSize(), ng.scalableResource.Replicas())
}
func (ng *nodegroup) Nodes() ([]cloudprovider.Instance, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 nodes, err := ng.scalableResource.Nodes()
 if err != nil {
  return nil, err
 }
 instances := make([]cloudprovider.Instance, len(nodes))
 for i := range nodes {
  instances[i] = cloudprovider.Instance{Id: nodes[i]}
 }
 return instances, nil
}
func (ng *nodegroup) TemplateNodeInfo() (*schedulercache.NodeInfo, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return nil, cloudprovider.ErrNotImplemented
}
func (ng *nodegroup) Exist() bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return true
}
func (ng *nodegroup) Create() (cloudprovider.NodeGroup, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return nil, cloudprovider.ErrAlreadyExist
}
func (ng *nodegroup) Delete() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return cloudprovider.ErrNotImplemented
}
func (ng *nodegroup) Autoprovisioned() bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return false
}
func newNodegroupFromMachineSet(controller *machineController, machineSet *v1beta1.MachineSet) (*nodegroup, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 scalableResource, err := newMachineSetScalableResource(controller, machineSet)
 if err != nil {
  return nil, err
 }
 return &nodegroup{machineapiClient: controller.clusterClientset.MachineV1beta1(), machineController: controller, scalableResource: scalableResource}, nil
}
func newNodegroupFromMachineDeployment(controller *machineController, machineDeployment *v1beta1.MachineDeployment) (*nodegroup, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 scalableResource, err := newMachineDeploymentScalableResource(controller, machineDeployment)
 if err != nil {
  return nil, err
 }
 return &nodegroup{machineapiClient: controller.clusterClientset.MachineV1beta1(), machineController: controller, scalableResource: scalableResource}, nil
}
