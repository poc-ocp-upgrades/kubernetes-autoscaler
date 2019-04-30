package aws

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/config/dynamic"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"k8s.io/klog"
)

const scaleToZeroSupported = true

type asgCache struct {
	registeredAsgs		[]*asg
	asgToInstances		map[AwsRef][]AwsInstanceRef
	instanceToAsg		map[AwsInstanceRef]*asg
	mutex			sync.Mutex
	service			autoScalingWrapper
	interrupt		chan struct{}
	asgAutoDiscoverySpecs	[]cloudprovider.ASGAutoDiscoveryConfig
	explicitlyConfigured	map[AwsRef]bool
}
type asg struct {
	AwsRef
	minSize			int
	maxSize			int
	curSize			int
	AvailabilityZones	[]string
	LaunchTemplateName	string
	LaunchTemplateVersion	string
	LaunchConfigurationName	string
	Tags			[]*autoscaling.TagDescription
}

func newASGCache(service autoScalingWrapper, explicitSpecs []string, autoDiscoverySpecs []cloudprovider.ASGAutoDiscoveryConfig) (*asgCache, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registry := &asgCache{registeredAsgs: make([]*asg, 0), service: service, asgToInstances: make(map[AwsRef][]AwsInstanceRef), instanceToAsg: make(map[AwsInstanceRef]*asg), interrupt: make(chan struct{}), asgAutoDiscoverySpecs: autoDiscoverySpecs, explicitlyConfigured: make(map[AwsRef]bool)}
	if err := registry.parseExplicitAsgs(explicitSpecs); err != nil {
		return nil, err
	}
	return registry, nil
}
func (m *asgCache) parseExplicitAsgs(specs []string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, spec := range specs {
		asg, err := m.buildAsgFromSpec(spec)
		if err != nil {
			return fmt.Errorf("failed to parse node group spec: %v", err)
		}
		m.explicitlyConfigured[asg.AwsRef] = true
		m.register(asg)
	}
	return nil
}
func (m *asgCache) register(asg *asg) *asg {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i := range m.registeredAsgs {
		if existing := m.registeredAsgs[i]; existing.AwsRef == asg.AwsRef {
			if reflect.DeepEqual(existing, asg) {
				return existing
			}
			klog.V(4).Infof("Updating ASG %s", asg.AwsRef.Name)
			if !m.explicitlyConfigured[asg.AwsRef] {
				existing.minSize = asg.minSize
				existing.maxSize = asg.maxSize
			}
			existing.curSize = asg.curSize
			existing.AvailabilityZones = asg.AvailabilityZones
			existing.LaunchConfigurationName = asg.LaunchConfigurationName
			existing.LaunchTemplateName = asg.LaunchTemplateName
			existing.LaunchTemplateVersion = asg.LaunchTemplateVersion
			existing.Tags = asg.Tags
			return existing
		}
	}
	klog.V(1).Infof("Registering ASG %s", asg.AwsRef.Name)
	m.registeredAsgs = append(m.registeredAsgs, asg)
	return asg
}
func (m *asgCache) unregister(a *asg) *asg {
	_logClusterCodePath()
	defer _logClusterCodePath()
	updated := make([]*asg, 0, len(m.registeredAsgs))
	var changed *asg
	for _, existing := range m.registeredAsgs {
		if existing.AwsRef == a.AwsRef {
			klog.V(1).Infof("Unregistered ASG %s", a.AwsRef.Name)
			changed = a
			continue
		}
		updated = append(updated, existing)
	}
	m.registeredAsgs = updated
	return changed
}
func (m *asgCache) buildAsgFromSpec(spec string) (*asg, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s, err := dynamic.SpecFromString(spec, scaleToZeroSupported)
	if err != nil {
		return nil, fmt.Errorf("failed to parse node group spec: %v", err)
	}
	asg := &asg{AwsRef: AwsRef{Name: s.Name}, minSize: s.MinSize, maxSize: s.MaxSize}
	return asg, nil
}
func (m *asgCache) Get() []*asg {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.registeredAsgs
}
func (m *asgCache) FindForInstance(instance AwsInstanceRef) *asg {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.findForInstance(instance)
}
func (m *asgCache) findForInstance(instance AwsInstanceRef) *asg {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if asg, found := m.instanceToAsg[instance]; found {
		return asg
	}
	return nil
}
func (m *asgCache) InstancesByAsg(ref AwsRef) ([]AwsInstanceRef, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if instances, found := m.asgToInstances[ref]; found {
		return instances, nil
	}
	return nil, fmt.Errorf("Error while looking for instances of ASG: %s", ref)
}
func (m *asgCache) SetAsgSize(asg *asg, size int) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	params := &autoscaling.SetDesiredCapacityInput{AutoScalingGroupName: aws.String(asg.Name), DesiredCapacity: aws.Int64(int64(size)), HonorCooldown: aws.Bool(false)}
	klog.V(0).Infof("Setting asg %s size to %d", asg.Name, size)
	_, err := m.service.SetDesiredCapacity(params)
	if err != nil {
		return err
	}
	asg.curSize = size
	return nil
}
func (m *asgCache) DeleteInstances(instances []*AwsInstanceRef) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if len(instances) == 0 {
		return nil
	}
	commonAsg := m.findForInstance(*instances[0])
	if commonAsg == nil {
		return fmt.Errorf("can't delete instance %s, which is not part of an ASG", instances[0].Name)
	}
	for _, instance := range instances {
		asg := m.findForInstance(*instance)
		if asg != commonAsg {
			instanceIds := make([]string, len(instances))
			for i, instance := range instances {
				instanceIds[i] = instance.Name
			}
			return fmt.Errorf("can't delete instances %s as they belong to at least two different ASGs (%s and %s)", strings.Join(instanceIds, ","), commonAsg.Name, asg.Name)
		}
	}
	for _, instance := range instances {
		params := &autoscaling.TerminateInstanceInAutoScalingGroupInput{InstanceId: aws.String(instance.Name), ShouldDecrementDesiredCapacity: aws.Bool(true)}
		resp, err := m.service.TerminateInstanceInAutoScalingGroup(params)
		if err != nil {
			return err
		}
		commonAsg.curSize--
		klog.V(4).Infof(*resp.Activity.Description)
	}
	return nil
}
func (m *asgCache) fetchAutoAsgNames() ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	groupNames := make([]string, 0)
	for _, spec := range m.asgAutoDiscoverySpecs {
		names, err := m.service.getAutoscalingGroupNamesByTags(spec.Tags)
		if err != nil {
			return nil, fmt.Errorf("cannot autodiscover ASGs: %s", err)
		}
		groupNames = append(groupNames, names...)
	}
	return groupNames, nil
}
func (m *asgCache) buildAsgNames() ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	refreshNames := make([]string, len(m.explicitlyConfigured))
	i := 0
	for k := range m.explicitlyConfigured {
		refreshNames[i] = k.Name
		i++
	}
	autoDiscoveredNames, err := m.fetchAutoAsgNames()
	if err != nil {
		return nil, err
	}
	for _, name := range autoDiscoveredNames {
		autoRef := AwsRef{Name: name}
		if m.explicitlyConfigured[autoRef] {
			continue
		}
		refreshNames = append(refreshNames, name)
	}
	return refreshNames, nil
}
func (m *asgCache) regenerate() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	newInstanceToAsgCache := make(map[AwsInstanceRef]*asg)
	newAsgToInstancesCache := make(map[AwsRef][]AwsInstanceRef)
	refreshNames, err := m.buildAsgNames()
	if err != nil {
		return err
	}
	klog.V(4).Infof("Regenerating instance to ASG map for ASGs: %v", refreshNames)
	groups, err := m.service.getAutoscalingGroupsByNames(refreshNames)
	if err != nil {
		return err
	}
	exists := make(map[AwsRef]bool)
	for _, group := range groups {
		asg, err := m.buildAsgFromAWS(group)
		if err != nil {
			return err
		}
		exists[asg.AwsRef] = true
		asg = m.register(asg)
		newAsgToInstancesCache[asg.AwsRef] = make([]AwsInstanceRef, len(group.Instances))
		for i, instance := range group.Instances {
			ref := m.buildInstanceRefFromAWS(instance)
			newInstanceToAsgCache[ref] = asg
			newAsgToInstancesCache[asg.AwsRef][i] = ref
		}
	}
	for _, asg := range m.registeredAsgs {
		if !exists[asg.AwsRef] && !m.explicitlyConfigured[asg.AwsRef] {
			m.unregister(asg)
		}
	}
	m.asgToInstances = newAsgToInstancesCache
	m.instanceToAsg = newInstanceToAsgCache
	return nil
}
func (m *asgCache) buildAsgFromAWS(g *autoscaling.Group) (*asg, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	spec := dynamic.NodeGroupSpec{Name: aws.StringValue(g.AutoScalingGroupName), MinSize: int(aws.Int64Value(g.MinSize)), MaxSize: int(aws.Int64Value(g.MaxSize)), SupportScaleToZero: scaleToZeroSupported}
	if verr := spec.Validate(); verr != nil {
		return nil, fmt.Errorf("failed to create node group spec: %v", verr)
	}
	launchTemplateName, launchTemplateVersion := m.buildLaunchTemplateParams(g)
	asg := &asg{AwsRef: AwsRef{Name: spec.Name}, minSize: spec.MinSize, maxSize: spec.MaxSize, curSize: int(aws.Int64Value(g.DesiredCapacity)), AvailabilityZones: aws.StringValueSlice(g.AvailabilityZones), LaunchConfigurationName: aws.StringValue(g.LaunchConfigurationName), LaunchTemplateName: launchTemplateName, LaunchTemplateVersion: launchTemplateVersion, Tags: g.Tags}
	return asg, nil
}
func (m *asgCache) buildLaunchTemplateParams(g *autoscaling.Group) (string, string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if g.LaunchTemplate != nil {
		return aws.StringValue(g.LaunchTemplate.LaunchTemplateName), aws.StringValue(g.LaunchTemplate.Version)
	}
	return "", ""
}
func (m *asgCache) buildInstanceRefFromAWS(instance *autoscaling.Instance) AwsInstanceRef {
	_logClusterCodePath()
	defer _logClusterCodePath()
	providerID := fmt.Sprintf("aws:///%s/%s", aws.StringValue(instance.AvailabilityZone), aws.StringValue(instance.InstanceId))
	return AwsInstanceRef{ProviderID: providerID, Name: aws.StringValue(instance.InstanceId)}
}
func (m *asgCache) Cleanup() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	close(m.interrupt)
}
