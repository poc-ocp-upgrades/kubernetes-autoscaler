package openshiftmachineapi

import (
	"strconv"
	"github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	nodeGroupMinSizeAnnotationKey	= "machine.openshift.io/cluster-api-autoscaler-node-group-min-size"
	nodeGroupMaxSizeAnnotationKey	= "machine.openshift.io/cluster-api-autoscaler-node-group-max-size"
)

var (
	errMissingMinAnnotation	= errors.New("missing min annotation")
	errMissingMaxAnnotation	= errors.New("missing max annotation")
	errInvalidMinAnnotation	= errors.New("invalid min annotation")
	errInvalidMaxAnnotation	= errors.New("invalid max annotation")
)

func minSize(annotations map[string]string) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	val, found := annotations[nodeGroupMinSizeAnnotationKey]
	if !found {
		return 0, errMissingMinAnnotation
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0, errors.Wrapf(err, "%s", errInvalidMinAnnotation)
	}
	return i, nil
}
func maxSize(annotations map[string]string) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	val, found := annotations[nodeGroupMaxSizeAnnotationKey]
	if !found {
		return 0, errMissingMaxAnnotation
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0, errors.Wrapf(err, "%s", errInvalidMaxAnnotation)
	}
	return i, nil
}
func parseScalingBounds(annotations map[string]string) (int, int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	minSize, err := minSize(annotations)
	if err != nil && err != errMissingMinAnnotation {
		return 0, 0, err
	}
	if minSize < 0 {
		return 0, 0, errInvalidMinAnnotation
	}
	maxSize, err := maxSize(annotations)
	if err != nil && err != errMissingMaxAnnotation {
		return 0, 0, err
	}
	if maxSize < 0 {
		return 0, 0, errInvalidMaxAnnotation
	}
	if maxSize < minSize {
		return 0, 0, errInvalidMaxAnnotation
	}
	return minSize, maxSize, nil
}
func machineOwnerRef(machine *v1beta1.Machine) *metav1.OwnerReference {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, ref := range machine.OwnerReferences {
		if ref.Kind == "MachineSet" && ref.Name != "" {
			return ref.DeepCopy()
		}
	}
	return nil
}
func machineIsOwnedByMachineSet(machine *v1beta1.Machine, machineSet *v1beta1.MachineSet) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if ref := machineOwnerRef(machine); ref != nil {
		return ref.UID == machineSet.UID
	}
	return false
}
func machineSetMachineDeploymentRef(machineSet *v1beta1.MachineSet) *metav1.OwnerReference {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, ref := range machineSet.OwnerReferences {
		if ref.Kind == "MachineDeployment" {
			return ref.DeepCopy()
		}
	}
	return nil
}
func machineSetHasMachineDeploymentOwnerRef(machineSet *v1beta1.MachineSet) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return machineSetMachineDeploymentRef(machineSet) != nil
}
func machineSetIsOwnedByMachineDeployment(machineSet *v1beta1.MachineSet, machineDeployment *v1beta1.MachineDeployment) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if ref := machineSetMachineDeploymentRef(machineSet); ref != nil {
		return ref.UID == machineDeployment.UID
	}
	return false
}
