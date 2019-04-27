package azure

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/klog"
)

var virtualMachineRE = regexp.MustCompile(`^azure://(?:.*)/providers/microsoft.compute/virtualmachines/(.+)$`)

type asgCache struct {
	registeredAsgs		[]cloudprovider.NodeGroup
	instanceToAsg		map[azureRef]cloudprovider.NodeGroup
	notInRegisteredAsg	map[azureRef]bool
	mutex			sync.Mutex
	interrupt		chan struct{}
}

func newAsgCache() (*asgCache, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cache := &asgCache{registeredAsgs: make([]cloudprovider.NodeGroup, 0), instanceToAsg: make(map[azureRef]cloudprovider.NodeGroup), notInRegisteredAsg: make(map[azureRef]bool), interrupt: make(chan struct{})}
	go wait.Until(func() {
		cache.mutex.Lock()
		defer cache.mutex.Unlock()
		if err := cache.regenerate(); err != nil {
			klog.Errorf("Error while regenerating Asg cache: %v", err)
		}
	}, time.Hour, cache.interrupt)
	return cache, nil
}
func (m *asgCache) Register(asg cloudprovider.NodeGroup) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for i := range m.registeredAsgs {
		if existing := m.registeredAsgs[i]; strings.EqualFold(existing.Id(), asg.Id()) {
			if reflect.DeepEqual(existing, asg) {
				return false
			}
			m.registeredAsgs[i] = asg
			klog.V(4).Infof("ASG %q updated", asg.Id())
			m.invalidateUnownedInstanceCache()
			return true
		}
	}
	klog.V(4).Infof("Registering ASG %q", asg.Id())
	m.registeredAsgs = append(m.registeredAsgs, asg)
	m.invalidateUnownedInstanceCache()
	return true
}
func (m *asgCache) invalidateUnownedInstanceCache() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	klog.V(4).Info("Invalidating unowned instance cache")
	m.notInRegisteredAsg = make(map[azureRef]bool)
}
func (m *asgCache) Unregister(asg cloudprovider.NodeGroup) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	updated := make([]cloudprovider.NodeGroup, 0, len(m.registeredAsgs))
	changed := false
	for _, existing := range m.registeredAsgs {
		if strings.EqualFold(existing.Id(), asg.Id()) {
			klog.V(1).Infof("Unregistered ASG %s", asg.Id())
			changed = true
			continue
		}
		updated = append(updated, existing)
	}
	m.registeredAsgs = updated
	return changed
}
func (m *asgCache) get() []cloudprovider.NodeGroup {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.registeredAsgs
}
func (m *asgCache) FindForInstance(instance *azureRef, vmType string) (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	inst := azureRef{Name: strings.ToLower(instance.Name)}
	if m.notInRegisteredAsg[inst] {
		return nil, nil
	}
	if vmType == vmTypeVMSS {
		if ok := virtualMachineRE.Match([]byte(inst.Name)); ok {
			klog.V(3).Infof("Instance %q is not managed by vmss, omit it in autoscaler", instance.Name)
			m.notInRegisteredAsg[inst] = true
			return nil, nil
		}
	}
	if vmType == vmTypeStandard {
		if ok := virtualMachineRE.Match([]byte(inst.Name)); !ok {
			klog.V(3).Infof("Instance %q is not in Azure resource ID format, omit it in autoscaler", instance.Name)
			m.notInRegisteredAsg[inst] = true
			return nil, nil
		}
	}
	if asg := m.getInstanceFromCache(inst.Name); asg != nil {
		return asg, nil
	}
	if err := m.regenerate(); err != nil {
		return nil, fmt.Errorf("error while looking for ASG for instance %q, error: %v", instance.Name, err)
	}
	if asg := m.getInstanceFromCache(inst.Name); asg != nil {
		return asg, nil
	}
	m.notInRegisteredAsg[inst] = true
	return nil, nil
}
func (m *asgCache) Cleanup() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	close(m.interrupt)
}
func (m *asgCache) regenerate() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	newCache := make(map[azureRef]cloudprovider.NodeGroup)
	for _, nsg := range m.registeredAsgs {
		instances, err := nsg.Nodes()
		if err != nil {
			return err
		}
		for _, instance := range instances {
			ref := azureRef{Name: instance.Id}
			newCache[ref] = nsg
		}
	}
	m.instanceToAsg = newCache
	return nil
}
func (m *asgCache) getInstanceFromCache(providerID string) cloudprovider.NodeGroup {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for instanceID, asg := range m.instanceToAsg {
		if strings.EqualFold(instanceID.GetKey(), providerID) {
			return asg
		}
	}
	return nil
}
