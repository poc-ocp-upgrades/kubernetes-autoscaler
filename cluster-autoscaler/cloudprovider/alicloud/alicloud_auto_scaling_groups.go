package alicloud

import (
 "sync"
 "time"
 "k8s.io/apimachinery/pkg/util/wait"
 "k8s.io/klog"
)

type autoScalingGroups struct {
 registeredAsgs           []*asgInformation
 instanceToAsg            map[string]*Asg
 cacheMutex               sync.Mutex
 instancesNotInManagedAsg map[string]struct{}
 service                  *autoScalingWrapper
}

func newAutoScalingGroups(service *autoScalingWrapper) *autoScalingGroups {
 _logClusterCodePath()
 defer _logClusterCodePath()
 registry := &autoScalingGroups{registeredAsgs: make([]*asgInformation, 0), service: service, instanceToAsg: make(map[string]*Asg), instancesNotInManagedAsg: make(map[string]struct{})}
 go wait.Forever(func() {
  registry.cacheMutex.Lock()
  defer registry.cacheMutex.Unlock()
  if err := registry.regenerateCache(); err != nil {
   klog.Errorf("failed to do regenerating ASG cache,because of %s", err.Error())
  }
 }, time.Hour)
 return registry
}
func (m *autoScalingGroups) Register(asg *Asg) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 m.cacheMutex.Lock()
 defer m.cacheMutex.Unlock()
 m.registeredAsgs = append(m.registeredAsgs, &asgInformation{config: asg})
}
func (m *autoScalingGroups) FindForInstance(instanceId string) (*Asg, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 m.cacheMutex.Lock()
 defer m.cacheMutex.Unlock()
 if config, found := m.instanceToAsg[instanceId]; found {
  return config, nil
 }
 if _, found := m.instancesNotInManagedAsg[instanceId]; found {
  return nil, nil
 }
 if err := m.regenerateCache(); err != nil {
  return nil, err
 }
 if config, found := m.instanceToAsg[instanceId]; found {
  return config, nil
 }
 m.instancesNotInManagedAsg[instanceId] = struct{}{}
 return nil, nil
}
func (m *autoScalingGroups) regenerateCache() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 newCache := make(map[string]*Asg)
 for _, asg := range m.registeredAsgs {
  instances, err := m.service.getScalingInstancesByGroup(asg.config.id)
  if err != nil {
   return err
  }
  for _, instance := range instances {
   newCache[instance.InstanceId] = asg.config
  }
 }
 m.instanceToAsg = newCache
 return nil
}
