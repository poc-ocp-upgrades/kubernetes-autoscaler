package openshiftmachineapi

import (
 "fmt"
 "path"
 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
 machinev1beta1 "github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset/typed/machine/v1beta1"
 metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
 "k8s.io/utils/pointer"
)

type machineDeploymentScalableResource struct {
 machineapiClient  machinev1beta1.MachineV1beta1Interface
 controller        *machineController
 machineDeployment *v1beta1.MachineDeployment
 maxSize           int
 minSize           int
}

var _ scalableResource = (*machineDeploymentScalableResource)(nil)

func (r machineDeploymentScalableResource) ID() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return path.Join(r.Namespace(), r.Name())
}
func (r machineDeploymentScalableResource) MaxSize() int {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return r.maxSize
}
func (r machineDeploymentScalableResource) MinSize() int {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return r.minSize
}
func (r machineDeploymentScalableResource) Name() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return r.machineDeployment.Name
}
func (r machineDeploymentScalableResource) Namespace() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return r.machineDeployment.Namespace
}
func (r machineDeploymentScalableResource) Nodes() ([]string, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result := []string{}
 if err := r.controller.filterAllMachineSets(func(machineSet *v1beta1.MachineSet) error {
  if machineSetIsOwnedByMachineDeployment(machineSet, r.machineDeployment) {
   names, err := r.controller.machineSetNodeNames(machineSet)
   if err != nil {
    return err
   }
   result = append(result, names...)
  }
  return nil
 }); err != nil {
  return nil, err
 }
 return result, nil
}
func (r machineDeploymentScalableResource) Replicas() int32 {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return pointer.Int32PtrDerefOr(r.machineDeployment.Spec.Replicas, 0)
}
func (r machineDeploymentScalableResource) SetSize(nreplicas int32) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 machineDeployment, err := r.machineapiClient.MachineDeployments(r.Namespace()).Get(r.Name(), metav1.GetOptions{})
 if err != nil {
  return fmt.Errorf("unable to get MachineDeployment %q: %v", r.ID(), err)
 }
 machineDeployment = machineDeployment.DeepCopy()
 machineDeployment.Spec.Replicas = &nreplicas
 _, err = r.machineapiClient.MachineDeployments(r.Namespace()).Update(machineDeployment)
 if err != nil {
  return fmt.Errorf("unable to update number of replicas of machineDeployment %q: %v", r.ID(), err)
 }
 return nil
}
func newMachineDeploymentScalableResource(controller *machineController, machineDeployment *v1beta1.MachineDeployment) (*machineDeploymentScalableResource, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 minSize, maxSize, err := parseScalingBounds(machineDeployment.Annotations)
 if err != nil {
  return nil, fmt.Errorf("error validating min/max annotations: %v", err)
 }
 return &machineDeploymentScalableResource{machineapiClient: controller.clusterClientset.MachineV1beta1(), controller: controller, machineDeployment: machineDeployment, maxSize: maxSize, minSize: minSize}, nil
}
