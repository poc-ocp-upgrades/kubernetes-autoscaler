package openshiftmachineapi

import (
	"fmt"
	"path"
	"github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	machinev1beta1 "github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset/typed/machine/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

type machineSetScalableResource struct {
	machineapiClient	machinev1beta1.MachineV1beta1Interface
	controller		*machineController
	machineSet		*v1beta1.MachineSet
	maxSize			int
	minSize			int
}

var _ scalableResource = (*machineSetScalableResource)(nil)

func (r machineSetScalableResource) ID() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return path.Join(r.Namespace(), r.Name())
}
func (r machineSetScalableResource) MaxSize() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.maxSize
}
func (r machineSetScalableResource) MinSize() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.minSize
}
func (r machineSetScalableResource) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.machineSet.Name
}
func (r machineSetScalableResource) Namespace() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.machineSet.Namespace
}
func (r machineSetScalableResource) Nodes() ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.controller.machineSetNodeNames(r.machineSet)
}
func (r machineSetScalableResource) Replicas() int32 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return pointer.Int32PtrDerefOr(r.machineSet.Spec.Replicas, 0)
}
func (r machineSetScalableResource) SetSize(nreplicas int32) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	machineSet, err := r.machineapiClient.MachineSets(r.Namespace()).Get(r.Name(), metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("unable to get MachineSet %q: %v", r.ID(), err)
	}
	machineSet = machineSet.DeepCopy()
	machineSet.Spec.Replicas = &nreplicas
	_, err = r.machineapiClient.MachineSets(r.Namespace()).Update(machineSet)
	if err != nil {
		return fmt.Errorf("unable to update number of replicas of machineset %q: %v", r.ID(), err)
	}
	return nil
}
func newMachineSetScalableResource(controller *machineController, machineSet *v1beta1.MachineSet) (*machineSetScalableResource, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	minSize, maxSize, err := parseScalingBounds(machineSet.Annotations)
	if err != nil {
		return nil, fmt.Errorf("error validating min/max annotations: %v", err)
	}
	return &machineSetScalableResource{machineapiClient: controller.clusterClientset.MachineV1beta1(), controller: controller, machineSet: machineSet, maxSize: maxSize, minSize: minSize}, nil
}
