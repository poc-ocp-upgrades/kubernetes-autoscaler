package model

import (
	"github.com/golang/glog"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type ResourceName string
type ResourceAmount int64
type Resources map[ResourceName]ResourceAmount

const (
	ResourceCPU		ResourceName	= "cpu"
	ResourceMemory		ResourceName	= "memory"
	MaxResourceAmount			= ResourceAmount(1e14)
)

func CPUAmountFromCores(cores float64) ResourceAmount {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return resourceAmountFromFloat(cores * 1000.0)
}
func CoresFromCPUAmount(cpuAmount ResourceAmount) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return float64(cpuAmount) / 1000.0
}
func QuantityFromCPUAmount(cpuAmount ResourceAmount) resource.Quantity {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return *resource.NewScaledQuantity(int64(cpuAmount), -3)
}
func MemoryAmountFromBytes(bytes float64) ResourceAmount {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return resourceAmountFromFloat(bytes)
}
func BytesFromMemoryAmount(memoryAmount ResourceAmount) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return float64(memoryAmount)
}
func QuantityFromMemoryAmount(memoryAmount ResourceAmount) resource.Quantity {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return *resource.NewScaledQuantity(int64(memoryAmount), 0)
}
func ScaleResource(amount ResourceAmount, factor float64) ResourceAmount {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return resourceAmountFromFloat(float64(amount) * factor)
}
func ResourcesAsResourceList(resources Resources) apiv1.ResourceList {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(apiv1.ResourceList)
	for key, resourceAmount := range resources {
		var newKey apiv1.ResourceName
		var quantity resource.Quantity
		switch key {
		case ResourceCPU:
			newKey = apiv1.ResourceCPU
			quantity = QuantityFromCPUAmount(resourceAmount)
		case ResourceMemory:
			newKey = apiv1.ResourceMemory
			quantity = QuantityFromMemoryAmount(resourceAmount)
		default:
			glog.Errorf("Cannot translate %v resource name", key)
			continue
		}
		result[newKey] = quantity
	}
	return result
}
func RoundResourceAmount(amount, unit ResourceAmount) ResourceAmount {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ResourceAmount(int64(amount) - int64(amount)%int64(unit))
}
func ResourceAmountMax(amount1, amount2 ResourceAmount) ResourceAmount {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if amount1 > amount2 {
		return amount1
	}
	return amount2
}
func resourceAmountFromFloat(amount float64) ResourceAmount {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if amount < 0 {
		return ResourceAmount(0)
	} else if amount > float64(MaxResourceAmount) {
		return MaxResourceAmount
	} else {
		return ResourceAmount(amount)
	}
}

type PodID struct {
	Namespace	string
	PodName		string
}
type ContainerID struct {
	PodID
	ContainerName	string
}
type VpaID struct {
	Namespace	string
	VpaName		string
}
