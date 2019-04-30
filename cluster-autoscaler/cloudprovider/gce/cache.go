package gce

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	gce "google.golang.org/api/compute/v1"
	"k8s.io/klog"
)

type MigInformation struct {
	Config		Mig
	Basename	string
}
type MachineTypeKey struct {
	Zone		string
	MachineType	string
}
type GceCache struct {
	migs		[]*MigInformation
	instancesCache	map[GceRef]Mig
	resourceLimiter	*cloudprovider.ResourceLimiter
	machinesCache	map[MachineTypeKey]*gce.MachineType
	cacheMutex	sync.Mutex
	migsMutex	sync.Mutex
	GceService	AutoscalingGceClient
}

func NewGceCache(gceService AutoscalingGceClient) GceCache {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return GceCache{migs: []*MigInformation{}, instancesCache: map[GceRef]Mig{}, machinesCache: map[MachineTypeKey]*gce.MachineType{}, GceService: gceService}
}
func (gc *GceCache) RegisterMig(mig Mig) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gc.migsMutex.Lock()
	defer gc.migsMutex.Unlock()
	for i := range gc.migs {
		if oldMig := gc.migs[i].Config; oldMig.GceRef() == mig.GceRef() {
			if !reflect.DeepEqual(oldMig, mig) {
				gc.migs[i].Config = mig
				klog.V(4).Infof("Updated Mig %s", mig.GceRef().String())
				return true
			}
			return false
		}
	}
	klog.V(1).Infof("Registering %s", mig.GceRef().String())
	gc.migs = append(gc.migs, &MigInformation{Config: mig})
	return true
}
func (gc *GceCache) UnregisterMig(toBeRemoved Mig) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gc.migsMutex.Lock()
	defer gc.migsMutex.Unlock()
	newMigs := make([]*MigInformation, 0, len(gc.migs))
	found := false
	for _, mig := range gc.migs {
		if mig.Config.GceRef() == toBeRemoved.GceRef() {
			klog.V(1).Infof("Unregistered Mig %s", toBeRemoved.GceRef().String())
			found = true
		} else {
			newMigs = append(newMigs, mig)
		}
	}
	gc.migs = newMigs
	return found
}
func (gc *GceCache) GetMigs() []*MigInformation {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gc.migsMutex.Lock()
	defer gc.migsMutex.Unlock()
	migs := make([]*MigInformation, 0, len(gc.migs))
	for _, mig := range gc.migs {
		migs = append(migs, &MigInformation{Basename: mig.Basename, Config: mig.Config})
	}
	return migs
}
func (gc *GceCache) updateMigBasename(ref GceRef, basename string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gc.migsMutex.Lock()
	defer gc.migsMutex.Unlock()
	for _, mig := range gc.migs {
		if mig.Config.GceRef() == ref {
			mig.Basename = basename
		}
	}
}
func (gc *GceCache) GetMigForInstance(instance *GceRef) (Mig, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gc.cacheMutex.Lock()
	defer gc.cacheMutex.Unlock()
	if mig, found := gc.instancesCache[*instance]; found {
		return mig, nil
	}
	for _, mig := range gc.GetMigs() {
		if mig.Config.GceRef().Project == instance.Project && mig.Config.GceRef().Zone == instance.Zone && strings.HasPrefix(instance.Name, mig.Basename) {
			if err := gc.regenerateCache(); err != nil {
				return nil, fmt.Errorf("Error while looking for MIG for instance %+v, error: %v", *instance, err)
			}
			if mig, found := gc.instancesCache[*instance]; found {
				return mig, nil
			}
			return nil, fmt.Errorf("Instance %+v does not belong to any configured MIG", *instance)
		}
	}
	return nil, nil
}
func (gc *GceCache) RegenerateInstancesCache() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gc.cacheMutex.Lock()
	defer gc.cacheMutex.Unlock()
	return gc.regenerateCache()
}
func (gc *GceCache) regenerateCache() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	newInstancesCache := make(map[GceRef]Mig)
	for _, migInfo := range gc.GetMigs() {
		mig := migInfo.Config
		klog.V(4).Infof("Regenerating MIG information for %s", mig.GceRef().String())
		basename, err := gc.GceService.FetchMigBasename(mig.GceRef())
		if err != nil {
			return err
		}
		gc.updateMigBasename(mig.GceRef(), basename)
		instances, err := gc.GceService.FetchMigInstances(mig.GceRef())
		if err != nil {
			klog.V(4).Infof("Failed MIG info request for %s: %v", mig.GceRef().String(), err)
			return err
		}
		for _, ref := range instances {
			newInstancesCache[ref] = mig
		}
	}
	gc.instancesCache = newInstancesCache
	return nil
}
func (gc *GceCache) SetResourceLimiter(resourceLimiter *cloudprovider.ResourceLimiter) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gc.cacheMutex.Lock()
	defer gc.cacheMutex.Unlock()
	gc.resourceLimiter = resourceLimiter
}
func (gc *GceCache) GetResourceLimiter() (*cloudprovider.ResourceLimiter, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gc.cacheMutex.Lock()
	defer gc.cacheMutex.Unlock()
	return gc.resourceLimiter, nil
}
func (gc *GceCache) GetMachineFromCache(machineType string, zone string) *gce.MachineType {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gc.cacheMutex.Lock()
	defer gc.cacheMutex.Unlock()
	return gc.machinesCache[MachineTypeKey{zone, machineType}]
}
func (gc *GceCache) AddMachineToCache(machineType string, zone string, machine *gce.MachineType) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gc.cacheMutex.Lock()
	defer gc.cacheMutex.Unlock()
	gc.machinesCache[MachineTypeKey{zone, machineType}] = machine
}
func (gc *GceCache) SetMachinesCache(machinesCache map[MachineTypeKey]*gce.MachineType) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gc.cacheMutex.Lock()
	defer gc.cacheMutex.Unlock()
	gc.machinesCache = machinesCache
}
