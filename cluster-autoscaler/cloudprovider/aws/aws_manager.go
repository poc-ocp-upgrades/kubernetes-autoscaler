package aws

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"gopkg.in/gcfg.v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/utils/gpu"
	"k8s.io/klog"
	provider_aws "k8s.io/kubernetes/pkg/cloudprovider/providers/aws"
	kubeletapis "k8s.io/kubernetes/pkg/kubelet/apis"
)

const (
	operationWaitTimeout	= 5 * time.Second
	operationPollInterval	= 100 * time.Millisecond
	maxRecordsReturnedByAPI	= 100
	maxAsgNamesPerDescribe	= 50
	refreshInterval		= 10 * time.Second
)

type AwsManager struct {
	autoScalingService	autoScalingWrapper
	ec2Service		ec2Wrapper
	asgCache		*asgCache
	lastRefresh		time.Time
}
type asgTemplate struct {
	InstanceType	*instanceType
	Region		string
	Zone		string
	Tags		[]*autoscaling.TagDescription
}

func getRegion(cfg ...*aws.Config) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	region, present := os.LookupEnv("AWS_REGION")
	if !present {
		svc := ec2metadata.New(session.New(), cfg...)
		if r, err := svc.Region(); err == nil {
			region = r
		}
	}
	return region
}
func createAWSManagerInternal(configReader io.Reader, discoveryOpts cloudprovider.NodeGroupDiscoveryOptions, autoScalingService *autoScalingWrapper, ec2Service *ec2Wrapper) (*AwsManager, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if configReader != nil {
		var cfg provider_aws.CloudConfig
		if err := gcfg.ReadInto(&cfg, configReader); err != nil {
			klog.Errorf("Couldn't read config: %v", err)
			return nil, err
		}
	}
	if autoScalingService == nil || ec2Service == nil {
		sess := session.New(aws.NewConfig().WithRegion(getRegion()))
		if autoScalingService == nil {
			autoScalingService = &autoScalingWrapper{autoscaling.New(sess)}
		}
		if ec2Service == nil {
			ec2Service = &ec2Wrapper{ec2.New(sess)}
		}
	}
	specs, err := discoveryOpts.ParseASGAutoDiscoverySpecs()
	if err != nil {
		return nil, err
	}
	cache, err := newASGCache(*autoScalingService, discoveryOpts.NodeGroupSpecs, specs)
	if err != nil {
		return nil, err
	}
	manager := &AwsManager{autoScalingService: *autoScalingService, ec2Service: *ec2Service, asgCache: cache}
	if err := manager.forceRefresh(); err != nil {
		return nil, err
	}
	return manager, nil
}
func CreateAwsManager(configReader io.Reader, discoveryOpts cloudprovider.NodeGroupDiscoveryOptions) (*AwsManager, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return createAWSManagerInternal(configReader, discoveryOpts, nil, nil)
}
func (m *AwsManager) Refresh() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m.lastRefresh.Add(refreshInterval).After(time.Now()) {
		return nil
	}
	return m.forceRefresh()
}
func (m *AwsManager) forceRefresh() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := m.asgCache.regenerate(); err != nil {
		klog.Errorf("Failed to regenerate ASG cache: %v", err)
		return err
	}
	m.lastRefresh = time.Now()
	klog.V(2).Infof("Refreshed ASG list, next refresh after %v", m.lastRefresh.Add(refreshInterval))
	return nil
}
func (m *AwsManager) GetAsgForInstance(instance AwsInstanceRef) *asg {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.asgCache.FindForInstance(instance)
}
func (m *AwsManager) Cleanup() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.asgCache.Cleanup()
}
func (m *AwsManager) getAsgs() []*asg {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.asgCache.Get()
}
func (m *AwsManager) SetAsgSize(asg *asg, size int) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.asgCache.SetAsgSize(asg, size)
}
func (m *AwsManager) DeleteInstances(instances []*AwsInstanceRef) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.asgCache.DeleteInstances(instances)
}
func (m *AwsManager) GetAsgNodes(ref AwsRef) ([]AwsInstanceRef, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.asgCache.InstancesByAsg(ref)
}
func (m *AwsManager) getAsgTemplate(asg *asg) (*asgTemplate, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(asg.AvailabilityZones) < 1 {
		return nil, fmt.Errorf("Unable to get first AvailabilityZone for ASG %q", asg.Name)
	}
	az := asg.AvailabilityZones[0]
	region := az[0 : len(az)-1]
	if len(asg.AvailabilityZones) > 1 {
		klog.Warningf("Found multiple availability zones for ASG %q; using %s\n", asg.Name, az)
	}
	instanceTypeName, err := m.buildInstanceType(asg)
	if err != nil {
		return nil, err
	}
	if t, ok := InstanceTypes[instanceTypeName]; ok {
		return &asgTemplate{InstanceType: t, Region: region, Zone: az, Tags: asg.Tags}, nil
	}
	return nil, fmt.Errorf("ASG %q uses the unknown EC2 instance type %q", asg.Name, instanceTypeName)
}
func (m *AwsManager) buildInstanceType(asg *asg) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if asg.LaunchConfigurationName != "" {
		return m.autoScalingService.getInstanceTypeByLCName(asg.LaunchConfigurationName)
	} else if asg.LaunchTemplateName != "" && asg.LaunchTemplateVersion != "" {
		return m.ec2Service.getInstanceTypeByLT(asg.LaunchTemplateName, asg.LaunchTemplateVersion)
	}
	return "", errors.New("Unable to get instance type from launch config or launch template")
}
func (m *AwsManager) buildNodeFromTemplate(asg *asg, template *asgTemplate) (*apiv1.Node, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	node := apiv1.Node{}
	nodeName := fmt.Sprintf("%s-asg-%d", asg.Name, rand.Int63())
	node.ObjectMeta = metav1.ObjectMeta{Name: nodeName, SelfLink: fmt.Sprintf("/api/v1/nodes/%s", nodeName), Labels: map[string]string{}}
	node.Status = apiv1.NodeStatus{Capacity: apiv1.ResourceList{}}
	node.Status.Capacity[apiv1.ResourcePods] = *resource.NewQuantity(110, resource.DecimalSI)
	node.Status.Capacity[apiv1.ResourceCPU] = *resource.NewQuantity(template.InstanceType.VCPU, resource.DecimalSI)
	node.Status.Capacity[gpu.ResourceNvidiaGPU] = *resource.NewQuantity(template.InstanceType.GPU, resource.DecimalSI)
	node.Status.Capacity[apiv1.ResourceMemory] = *resource.NewQuantity(template.InstanceType.MemoryMb*1024*1024, resource.DecimalSI)
	node.Status.Allocatable = node.Status.Capacity
	node.Labels = cloudprovider.JoinStringMaps(node.Labels, extractLabelsFromAsg(template.Tags))
	node.Labels = cloudprovider.JoinStringMaps(node.Labels, buildGenericLabels(template, nodeName))
	node.Spec.Taints = extractTaintsFromAsg(template.Tags)
	node.Status.Conditions = cloudprovider.BuildReadyConditions()
	return &node, nil
}
func buildGenericLabels(template *asgTemplate, nodeName string) map[string]string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(map[string]string)
	result[kubeletapis.LabelArch] = cloudprovider.DefaultArch
	result[kubeletapis.LabelOS] = cloudprovider.DefaultOS
	result[kubeletapis.LabelInstanceType] = template.InstanceType.InstanceType
	result[kubeletapis.LabelZoneRegion] = template.Region
	result[kubeletapis.LabelZoneFailureDomain] = template.Zone
	result[kubeletapis.LabelHostname] = nodeName
	return result
}
func extractLabelsFromAsg(tags []*autoscaling.TagDescription) map[string]string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(map[string]string)
	for _, tag := range tags {
		k := *tag.Key
		v := *tag.Value
		splits := strings.Split(k, "k8s.io/cluster-autoscaler/node-template/label/")
		if len(splits) > 1 {
			label := splits[1]
			if label != "" {
				result[label] = v
			}
		}
	}
	return result
}
func extractTaintsFromAsg(tags []*autoscaling.TagDescription) []apiv1.Taint {
	_logClusterCodePath()
	defer _logClusterCodePath()
	taints := make([]apiv1.Taint, 0)
	for _, tag := range tags {
		k := *tag.Key
		v := *tag.Value
		r, _ := regexp.Compile("(.*):(?:NoSchedule|NoExecute|PreferNoSchedule)")
		if r.MatchString(v) {
			splits := strings.Split(k, "k8s.io/cluster-autoscaler/node-template/taint/")
			if len(splits) > 1 {
				values := strings.SplitN(v, ":", 2)
				if len(values) > 1 {
					taints = append(taints, apiv1.Taint{Key: splits[1], Value: values[0], Effect: apiv1.TaintEffect(values[1])})
				}
			}
		}
	}
	return taints
}
