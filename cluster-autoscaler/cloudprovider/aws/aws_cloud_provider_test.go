package aws

import (
	"testing"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
)

type AutoScalingMock struct{ mock.Mock }

func (a *AutoScalingMock) DescribeAutoScalingGroupsPages(i *autoscaling.DescribeAutoScalingGroupsInput, fn func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := a.Called(i, fn)
	return args.Error(0)
}
func (a *AutoScalingMock) DescribeLaunchConfigurations(i *autoscaling.DescribeLaunchConfigurationsInput) (*autoscaling.DescribeLaunchConfigurationsOutput, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := a.Called(i)
	return args.Get(0).(*autoscaling.DescribeLaunchConfigurationsOutput), nil
}
func (a *AutoScalingMock) DescribeTagsPages(i *autoscaling.DescribeTagsInput, fn func(*autoscaling.DescribeTagsOutput, bool) bool) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := a.Called(i, fn)
	return args.Error(0)
}
func (a *AutoScalingMock) SetDesiredCapacity(input *autoscaling.SetDesiredCapacityInput) (*autoscaling.SetDesiredCapacityOutput, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := a.Called(input)
	return args.Get(0).(*autoscaling.SetDesiredCapacityOutput), nil
}
func (a *AutoScalingMock) TerminateInstanceInAutoScalingGroup(input *autoscaling.TerminateInstanceInAutoScalingGroupInput) (*autoscaling.TerminateInstanceInAutoScalingGroupOutput, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := a.Called(input)
	return args.Get(0).(*autoscaling.TerminateInstanceInAutoScalingGroupOutput), nil
}

type EC2Mock struct{ mock.Mock }

func (e *EC2Mock) DescribeLaunchTemplateVersions(i *ec2.DescribeLaunchTemplateVersionsInput) (*ec2.DescribeLaunchTemplateVersionsOutput, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := e.Called(i)
	return args.Get(0).(*ec2.DescribeLaunchTemplateVersionsOutput), nil
}

var testService = autoScalingWrapper{&AutoScalingMock{}}
var testAwsManager = &AwsManager{asgCache: &asgCache{registeredAsgs: make([]*asg, 0), asgToInstances: make(map[AwsRef][]AwsInstanceRef), instanceToAsg: make(map[AwsInstanceRef]*asg), interrupt: make(chan struct{}), service: testService}, autoScalingService: testService}

func newTestAwsManagerWithService(service autoScaling, autoDiscoverySpecs []cloudprovider.ASGAutoDiscoveryConfig) *AwsManager {
	_logClusterCodePath()
	defer _logClusterCodePath()
	wrapper := autoScalingWrapper{service}
	return &AwsManager{autoScalingService: wrapper, asgCache: &asgCache{registeredAsgs: make([]*asg, 0), asgToInstances: make(map[AwsRef][]AwsInstanceRef), instanceToAsg: make(map[AwsInstanceRef]*asg), explicitlyConfigured: make(map[AwsRef]bool), interrupt: make(chan struct{}), asgAutoDiscoverySpecs: autoDiscoverySpecs, service: wrapper}}
}
func newTestAwsManagerWithAsgs(t *testing.T, service autoScaling, specs []string) *AwsManager {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := newTestAwsManagerWithService(service, nil)
	m.asgCache.parseExplicitAsgs(specs)
	return m
}
func newTestAwsManagerWithAutoAsgs(t *testing.T, service autoScaling, specs []string, autoDiscoverySpecs []cloudprovider.ASGAutoDiscoveryConfig) *AwsManager {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := newTestAwsManagerWithService(service, autoDiscoverySpecs)
	m.asgCache.parseExplicitAsgs(specs)
	return m
}
func testNamedDescribeAutoScalingGroupsOutput(groupName string, desiredCap int64, instanceIds ...string) *autoscaling.DescribeAutoScalingGroupsOutput {
	_logClusterCodePath()
	defer _logClusterCodePath()
	instances := []*autoscaling.Instance{}
	for _, id := range instanceIds {
		instances = append(instances, &autoscaling.Instance{InstanceId: aws.String(id), AvailabilityZone: aws.String("us-east-1a")})
	}
	return &autoscaling.DescribeAutoScalingGroupsOutput{AutoScalingGroups: []*autoscaling.Group{{AutoScalingGroupName: aws.String(groupName), DesiredCapacity: aws.Int64(desiredCap), MinSize: aws.Int64(1), MaxSize: aws.Int64(5), Instances: instances}}}
}
func testProvider(t *testing.T, m *AwsManager) *awsCloudProvider {
	_logClusterCodePath()
	defer _logClusterCodePath()
	resourceLimiter := cloudprovider.NewResourceLimiter(map[string]int64{cloudprovider.ResourceNameCores: 1, cloudprovider.ResourceNameMemory: 10000000}, map[string]int64{cloudprovider.ResourceNameCores: 10, cloudprovider.ResourceNameMemory: 100000000})
	provider, err := BuildAwsCloudProvider(m, resourceLimiter)
	assert.NoError(t, err)
	return provider.(*awsCloudProvider)
}
func TestBuildAwsCloudProvider(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	resourceLimiter := cloudprovider.NewResourceLimiter(map[string]int64{cloudprovider.ResourceNameCores: 1, cloudprovider.ResourceNameMemory: 10000000}, map[string]int64{cloudprovider.ResourceNameCores: 10, cloudprovider.ResourceNameMemory: 100000000})
	_, err := BuildAwsCloudProvider(testAwsManager, resourceLimiter)
	assert.NoError(t, err)
}
func TestName(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	provider := testProvider(t, testAwsManager)
	assert.Equal(t, provider.Name(), ProviderName)
}
func TestNodeGroups(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	provider := testProvider(t, newTestAwsManagerWithAsgs(t, testService, []string{"1:5:test-asg"}))
	nodeGroups := provider.NodeGroups()
	assert.Equal(t, len(nodeGroups), 1)
	assert.Equal(t, nodeGroups[0].Id(), "test-asg")
	assert.Equal(t, nodeGroups[0].MinSize(), 1)
	assert.Equal(t, nodeGroups[0].MaxSize(), 5)
}
func TestAutoDiscoveredNodeGroups(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	service := &AutoScalingMock{}
	provider := testProvider(t, newTestAwsManagerWithAutoAsgs(t, service, []string{}, []cloudprovider.ASGAutoDiscoveryConfig{{Tags: map[string]string{"test": ""}}}))
	service.On("DescribeTagsPages", &autoscaling.DescribeTagsInput{Filters: []*autoscaling.Filter{{Name: aws.String("key"), Values: aws.StringSlice([]string{"test"})}}, MaxRecords: aws.Int64(maxRecordsReturnedByAPI)}, mock.AnythingOfType("func(*autoscaling.DescribeTagsOutput, bool) bool")).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(*autoscaling.DescribeTagsOutput, bool) bool)
		fn(&autoscaling.DescribeTagsOutput{Tags: []*autoscaling.TagDescription{{ResourceId: aws.String("auto-asg")}}}, false)
	}).Return(nil).Once()
	service.On("DescribeAutoScalingGroupsPages", &autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: aws.StringSlice([]string{"auto-asg"}), MaxRecords: aws.Int64(maxRecordsReturnedByAPI)}, mock.AnythingOfType("func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool")).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool)
		fn(testNamedDescribeAutoScalingGroupsOutput("auto-asg", 1, "test-instance-id"), false)
	}).Return(nil)
	provider.Refresh()
	nodeGroups := provider.NodeGroups()
	assert.Equal(t, len(nodeGroups), 1)
	assert.Equal(t, nodeGroups[0].Id(), "auto-asg")
	assert.Equal(t, nodeGroups[0].MinSize(), 1)
	assert.Equal(t, nodeGroups[0].MaxSize(), 5)
}
func TestNodeGroupForNode(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	node := &apiv1.Node{Spec: apiv1.NodeSpec{ProviderID: "aws:///us-east-1a/test-instance-id"}}
	service := &AutoScalingMock{}
	provider := testProvider(t, newTestAwsManagerWithAsgs(t, service, []string{"1:5:test-asg"}))
	service.On("DescribeAutoScalingGroupsPages", &autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: aws.StringSlice([]string{"test-asg"}), MaxRecords: aws.Int64(maxRecordsReturnedByAPI)}, mock.AnythingOfType("func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool")).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool)
		fn(testNamedDescribeAutoScalingGroupsOutput("test-asg", 1, "test-instance-id"), false)
	}).Return(nil)
	provider.Refresh()
	group, err := provider.NodeGroupForNode(node)
	assert.NoError(t, err)
	assert.Equal(t, group.Id(), "test-asg")
	assert.Equal(t, group.MinSize(), 1)
	assert.Equal(t, group.MaxSize(), 5)
	nodes, err := group.Nodes()
	assert.NoError(t, err)
	assert.Equal(t, []cloudprovider.Instance{{Id: "aws:///us-east-1a/test-instance-id"}}, nodes)
	service.AssertNumberOfCalls(t, "DescribeAutoScalingGroupsPages", 1)
	nodeNotInGroup := &apiv1.Node{Spec: apiv1.NodeSpec{ProviderID: "aws:///us-east-1a/test-instance-id-not-in-group"}}
	group, err = provider.NodeGroupForNode(nodeNotInGroup)
	assert.NoError(t, err)
	assert.Nil(t, group)
	service.AssertNumberOfCalls(t, "DescribeAutoScalingGroupsPages", 1)
}
func TestNodeGroupForNodeWithNoProviderId(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	node := &apiv1.Node{Spec: apiv1.NodeSpec{ProviderID: ""}}
	service := &AutoScalingMock{}
	provider := testProvider(t, newTestAwsManagerWithAsgs(t, service, []string{"1:5:test-asg"}))
	group, err := provider.NodeGroupForNode(node)
	assert.NoError(t, err)
	assert.Equal(t, group, nil)
}
func TestAwsRefFromProviderId(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, err := AwsRefFromProviderId("aws123")
	assert.Error(t, err)
	_, err = AwsRefFromProviderId("aws://test-az/test-instance-id")
	assert.Error(t, err)
	awsRef, err := AwsRefFromProviderId("aws:///us-east-1a/i-260942b3")
	assert.NoError(t, err)
	assert.Equal(t, awsRef, &AwsInstanceRef{Name: "i-260942b3", ProviderID: "aws:///us-east-1a/i-260942b3"})
}
func TestTargetSize(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	service := &AutoScalingMock{}
	provider := testProvider(t, newTestAwsManagerWithAsgs(t, service, []string{"1:5:test-asg"}))
	asgs := provider.NodeGroups()
	service.On("DescribeAutoScalingGroupsPages", &autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: aws.StringSlice([]string{"test-asg"}), MaxRecords: aws.Int64(maxRecordsReturnedByAPI)}, mock.AnythingOfType("func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool")).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool)
		fn(testNamedDescribeAutoScalingGroupsOutput("test-asg", 2, "test-instance-id", "second-test-instance-id"), false)
	}).Return(nil)
	provider.Refresh()
	targetSize, err := asgs[0].TargetSize()
	assert.Equal(t, targetSize, 2)
	assert.NoError(t, err)
	service.AssertNumberOfCalls(t, "DescribeAutoScalingGroupsPages", 1)
}
func TestIncreaseSize(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	service := &AutoScalingMock{}
	provider := testProvider(t, newTestAwsManagerWithAsgs(t, service, []string{"1:5:test-asg"}))
	asgs := provider.NodeGroups()
	service.On("SetDesiredCapacity", &autoscaling.SetDesiredCapacityInput{AutoScalingGroupName: aws.String(asgs[0].Id()), DesiredCapacity: aws.Int64(3), HonorCooldown: aws.Bool(false)}).Return(&autoscaling.SetDesiredCapacityOutput{})
	service.On("DescribeAutoScalingGroupsPages", &autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: aws.StringSlice([]string{"test-asg"}), MaxRecords: aws.Int64(maxRecordsReturnedByAPI)}, mock.AnythingOfType("func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool")).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool)
		fn(testNamedDescribeAutoScalingGroupsOutput("test-asg", 2, "test-instance-id", "second-test-instance-id"), false)
	}).Return(nil)
	provider.Refresh()
	initialSize, err := asgs[0].TargetSize()
	assert.NoError(t, err)
	assert.Equal(t, 2, initialSize)
	err = asgs[0].IncreaseSize(1)
	assert.NoError(t, err)
	service.AssertNumberOfCalls(t, "SetDesiredCapacity", 1)
	service.AssertNumberOfCalls(t, "DescribeAutoScalingGroupsPages", 1)
	newSize, err := asgs[0].TargetSize()
	assert.NoError(t, err)
	assert.Equal(t, 3, newSize)
}
func TestBelongs(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	service := &AutoScalingMock{}
	provider := testProvider(t, newTestAwsManagerWithAsgs(t, service, []string{"1:5:test-asg"}))
	asgs := provider.NodeGroups()
	service.On("DescribeAutoScalingGroupsPages", &autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: aws.StringSlice([]string{asgs[0].Id()}), MaxRecords: aws.Int64(maxRecordsReturnedByAPI)}, mock.AnythingOfType("func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool")).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool)
		fn(testNamedDescribeAutoScalingGroupsOutput("test-asg", 1, "test-instance-id"), false)
	}).Return(nil)
	provider.Refresh()
	invalidNode := &apiv1.Node{Spec: apiv1.NodeSpec{ProviderID: "aws:///us-east-1a/invalid-instance-id"}}
	_, err := asgs[0].(*AwsNodeGroup).Belongs(invalidNode)
	assert.Error(t, err)
	service.AssertNumberOfCalls(t, "DescribeAutoScalingGroupsPages", 1)
	validNode := &apiv1.Node{Spec: apiv1.NodeSpec{ProviderID: "aws:///us-east-1a/test-instance-id"}}
	belongs, err := asgs[0].(*AwsNodeGroup).Belongs(validNode)
	assert.Equal(t, belongs, true)
	assert.NoError(t, err)
	service.AssertNumberOfCalls(t, "DescribeAutoScalingGroupsPages", 1)
}
func TestDeleteNodes(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	service := &AutoScalingMock{}
	provider := testProvider(t, newTestAwsManagerWithAsgs(t, service, []string{"1:5:test-asg"}))
	asgs := provider.NodeGroups()
	service.On("TerminateInstanceInAutoScalingGroup", &autoscaling.TerminateInstanceInAutoScalingGroupInput{InstanceId: aws.String("test-instance-id"), ShouldDecrementDesiredCapacity: aws.Bool(true)}).Return(&autoscaling.TerminateInstanceInAutoScalingGroupOutput{Activity: &autoscaling.Activity{Description: aws.String("Deleted instance")}})
	service.On("DescribeAutoScalingGroupsPages", &autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: aws.StringSlice([]string{"test-asg"}), MaxRecords: aws.Int64(maxRecordsReturnedByAPI)}, mock.AnythingOfType("func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool")).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool)
		fn(testNamedDescribeAutoScalingGroupsOutput("test-asg", 2, "test-instance-id", "second-test-instance-id"), false)
	}).Return(nil)
	provider.Refresh()
	initialSize, err := asgs[0].TargetSize()
	assert.NoError(t, err)
	assert.Equal(t, 2, initialSize)
	node := &apiv1.Node{Spec: apiv1.NodeSpec{ProviderID: "aws:///us-east-1a/test-instance-id"}}
	err = asgs[0].DeleteNodes([]*apiv1.Node{node})
	assert.NoError(t, err)
	service.AssertNumberOfCalls(t, "TerminateInstanceInAutoScalingGroup", 1)
	service.AssertNumberOfCalls(t, "DescribeAutoScalingGroupsPages", 1)
	newSize, err := asgs[0].TargetSize()
	assert.NoError(t, err)
	assert.Equal(t, 1, newSize)
}
func TestDeleteNodesAfterMultipleRefreshes(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	service := &AutoScalingMock{}
	manager := newTestAwsManagerWithAsgs(t, service, []string{"1:5:test-asg"})
	provider := testProvider(t, manager)
	asgs := provider.NodeGroups()
	service.On("TerminateInstanceInAutoScalingGroup", &autoscaling.TerminateInstanceInAutoScalingGroupInput{InstanceId: aws.String("test-instance-id"), ShouldDecrementDesiredCapacity: aws.Bool(true)}).Return(&autoscaling.TerminateInstanceInAutoScalingGroupOutput{Activity: &autoscaling.Activity{Description: aws.String("Deleted instance")}})
	service.On("DescribeAutoScalingGroupsPages", &autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: aws.StringSlice([]string{"test-asg"}), MaxRecords: aws.Int64(maxRecordsReturnedByAPI)}, mock.AnythingOfType("func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool")).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool)
		fn(testNamedDescribeAutoScalingGroupsOutput("test-asg", 2, "test-instance-id", "second-test-instance-id"), false)
	}).Return(nil)
	provider.Refresh()
	manager.forceRefresh()
	node := &apiv1.Node{Spec: apiv1.NodeSpec{ProviderID: "aws:///us-east-1a/test-instance-id"}}
	err := asgs[0].DeleteNodes([]*apiv1.Node{node})
	assert.NoError(t, err)
}
func TestGetResourceLimiter(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	service := &AutoScalingMock{}
	m := newTestAwsManagerWithService(service, nil)
	provider := testProvider(t, m)
	_, err := provider.GetResourceLimiter()
	assert.NoError(t, err)
}
func TestCleanup(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	provider := testProvider(t, testAwsManager)
	err := provider.Cleanup()
	assert.NoError(t, err)
}
