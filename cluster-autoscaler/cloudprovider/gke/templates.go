package gke

import (
	"fmt"
	"math/rand"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/gce"
	"k8s.io/autoscaler/cluster-autoscaler/utils/gpu"
)

type GkeTemplateBuilder struct {
	gce.GceTemplateBuilder
	projectId	string
}

func (t *GkeTemplateBuilder) BuildNodeFromMigSpec(mig *GkeMig, cpu int64, mem int64) (*apiv1.Node, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if mig.Spec() == nil {
		return nil, fmt.Errorf("no spec in mig %s", mig.GceRef().Name)
	}
	node := apiv1.Node{}
	nodeName := fmt.Sprintf("%s-autoprovisioned-template-%d", mig.GceRef().Name, rand.Int63())
	node.ObjectMeta = metav1.ObjectMeta{Name: nodeName, SelfLink: fmt.Sprintf("/api/v1/nodes/%s", nodeName), Labels: map[string]string{}}
	capacity, err := t.BuildCapacity(cpu, mem, nil)
	if err != nil {
		return nil, err
	}
	if gpuRequest, found := mig.Spec().ExtraResources[gpu.ResourceNvidiaGPU]; found {
		capacity[gpu.ResourceNvidiaGPU] = gpuRequest.DeepCopy()
	}
	kubeReserved := t.BuildKubeReserved(cpu, mem)
	node.Status = apiv1.NodeStatus{Capacity: capacity, Allocatable: t.CalculateAllocatable(capacity, kubeReserved)}
	labels, err := buildLabelsForAutoprovisionedMig(mig, nodeName)
	if err != nil {
		return nil, err
	}
	node.Labels = labels
	node.Spec.Taints = mig.Spec().Taints
	node.Status.Conditions = cloudprovider.BuildReadyConditions()
	return &node, nil
}
func (t *GkeTemplateBuilder) BuildKubeReserved(cpu, physicalMemory int64) apiv1.ResourceList {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cpuReservedMillicores := PredictKubeReservedCpuMillicores(cpu * 1000)
	memoryReserved := PredictKubeReservedMemory(physicalMemory)
	reserved := apiv1.ResourceList{}
	reserved[apiv1.ResourceCPU] = *resource.NewMilliQuantity(cpuReservedMillicores, resource.DecimalSI)
	reserved[apiv1.ResourceMemory] = *resource.NewQuantity(memoryReserved, resource.BinarySI)
	return reserved
}
func buildLabelsForAutoprovisionedMig(mig *GkeMig, nodeName string) (map[string]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	labels, err := gce.BuildGenericLabels(mig.GceRef(), mig.Spec().MachineType, nodeName)
	if err != nil {
		return nil, err
	}
	for k, v := range mig.Spec().Labels {
		if existingValue, found := labels[k]; found {
			if v != existingValue {
				return map[string]string{}, fmt.Errorf("conflict in labels requested: %s=%s  present: %s=%s", k, v, k, existingValue)
			}
		} else {
			labels[k] = v
		}
	}
	return labels, nil
}
