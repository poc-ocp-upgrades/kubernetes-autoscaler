package aws

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	kubeletapis "k8s.io/kubernetes/pkg/kubelet/apis"
)

func TestGetRegion(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	key := "AWS_REGION"
	originalRegion, originalPresent := os.LookupEnv(key)
	defer func(region string, present bool) {
		os.Unsetenv(key)
		if present {
			os.Setenv(key, region)
		}
	}(originalRegion, originalPresent)
	expected1 := "the-shire-1"
	os.Setenv(key, expected1)
	assert.Equal(t, expected1, getRegion())
	expected2 := "mordor-2"
	expected2a := expected2 + "a"
	os.Unsetenv(key)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expected2a))
	}))
	cfg := aws.NewConfig().WithEndpoint(server.URL)
	assert.Equal(t, expected2, getRegion(cfg))
}
func TestBuildGenericLabels(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	labels := buildGenericLabels(&asgTemplate{InstanceType: &instanceType{InstanceType: "c4.large", VCPU: 2, MemoryMb: 3840}, Region: "us-east-1"}, "sillyname")
	assert.Equal(t, "us-east-1", labels[kubeletapis.LabelZoneRegion])
	assert.Equal(t, "sillyname", labels[kubeletapis.LabelHostname])
	assert.Equal(t, "c4.large", labels[kubeletapis.LabelInstanceType])
	assert.Equal(t, cloudprovider.DefaultArch, labels[kubeletapis.LabelArch])
	assert.Equal(t, cloudprovider.DefaultOS, labels[kubeletapis.LabelOS])
}
func TestExtractLabelsFromAsg(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	tags := []*autoscaling.TagDescription{{Key: aws.String("k8s.io/cluster-autoscaler/node-template/label/foo"), Value: aws.String("bar")}, {Key: aws.String("bar"), Value: aws.String("baz")}}
	labels := extractLabelsFromAsg(tags)
	assert.Equal(t, 1, len(labels))
	assert.Equal(t, "bar", labels["foo"])
}
func TestExtractTaintsFromAsg(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	tags := []*autoscaling.TagDescription{{Key: aws.String("k8s.io/cluster-autoscaler/node-template/taint/dedicated"), Value: aws.String("foo:NoSchedule")}, {Key: aws.String("k8s.io/cluster-autoscaler/node-template/taint/group"), Value: aws.String("bar:NoExecute")}, {Key: aws.String("k8s.io/cluster-autoscaler/node-template/taint/app"), Value: aws.String("fizz:PreferNoSchedule")}, {Key: aws.String("bar"), Value: aws.String("baz")}, {Key: aws.String("k8s.io/cluster-autoscaler/node-template/taint/blank"), Value: aws.String("")}, {Key: aws.String("k8s.io/cluster-autoscaler/node-template/taint/nosplit"), Value: aws.String("some_value")}}
	expectedTaints := []apiv1.Taint{{Key: "dedicated", Value: "foo", Effect: apiv1.TaintEffectNoSchedule}, {Key: "group", Value: "bar", Effect: apiv1.TaintEffectNoExecute}, {Key: "app", Value: "fizz", Effect: apiv1.TaintEffectPreferNoSchedule}}
	taints := extractTaintsFromAsg(tags)
	assert.Equal(t, 3, len(taints))
	assert.Equal(t, makeTaintSet(expectedTaints), makeTaintSet(taints))
}
func makeTaintSet(taints []apiv1.Taint) map[apiv1.Taint]bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	set := make(map[apiv1.Taint]bool)
	for _, taint := range taints {
		set[taint] = true
	}
	return set
}
func TestFetchExplicitAsgs(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	min, max, groupname := 1, 10, "coolasg"
	s := &AutoScalingMock{}
	s.On("DescribeAutoScalingGroups", &autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: []*string{aws.String(groupname)}, MaxRecords: aws.Int64(1)}).Return(&autoscaling.DescribeAutoScalingGroupsOutput{AutoScalingGroups: []*autoscaling.Group{{AutoScalingGroupName: aws.String(groupname)}}})
	s.On("DescribeAutoScalingGroupsPages", &autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: aws.StringSlice([]string{groupname}), MaxRecords: aws.Int64(maxRecordsReturnedByAPI)}, mock.AnythingOfType("func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool")).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool)
		fn(&autoscaling.DescribeAutoScalingGroupsOutput{AutoScalingGroups: []*autoscaling.Group{{AutoScalingGroupName: aws.String(groupname)}}}, false)
	}).Return(nil)
	do := cloudprovider.NodeGroupDiscoveryOptions{NodeGroupSpecs: []string{fmt.Sprintf("%d:%d:%s", min, max, groupname), fmt.Sprintf("%d:%d:%s", min, max-1, groupname)}}
	m, err := createAWSManagerInternal(nil, do, &autoScalingWrapper{s}, nil)
	assert.NoError(t, err)
	asgs := m.asgCache.Get()
	assert.Equal(t, 1, len(asgs))
	validateAsg(t, asgs[0], groupname, min, max)
}
func TestBuildInstanceType(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ltName, ltVersion, instanceType := "launcher", "1", "t2.large"
	s := &EC2Mock{}
	s.On("DescribeLaunchTemplateVersions", &ec2.DescribeLaunchTemplateVersionsInput{LaunchTemplateName: aws.String(ltName), Versions: []*string{aws.String(ltVersion)}}).Return(&ec2.DescribeLaunchTemplateVersionsOutput{LaunchTemplateVersions: []*ec2.LaunchTemplateVersion{{LaunchTemplateData: &ec2.ResponseLaunchTemplateData{InstanceType: aws.String(instanceType)}}}})
	m, err := createAWSManagerInternal(nil, cloudprovider.NodeGroupDiscoveryOptions{}, nil, &ec2Wrapper{s})
	assert.NoError(t, err)
	asg := asg{LaunchTemplateName: ltName, LaunchTemplateVersion: ltVersion}
	builtInstanceType, err := m.buildInstanceType(&asg)
	assert.NoError(t, err)
	assert.Equal(t, instanceType, builtInstanceType)
}
func TestGetASGTemplate(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	const (
		knownInstanceType	= "t3.micro"
		region			= "us-east-1"
		az			= region + "a"
		ltName			= "launcher"
		ltVersion		= "1"
	)
	tags := []*autoscaling.TagDescription{{Key: aws.String("k8s.io/cluster-autoscaler/node-template/taint/dedicated"), Value: aws.String("foo:NoSchedule")}}
	tests := []struct {
		description		string
		instanceType		string
		availabilityZones	[]string
		error			bool
	}{{"insufficient availability zones", knownInstanceType, []string{}, true}, {"single availability zone", knownInstanceType, []string{az}, false}, {"multiple availability zones", knownInstanceType, []string{az, "us-west-1b"}, false}, {"unknown instance type", "nonexistent.xlarge", []string{az}, true}}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			s := &EC2Mock{}
			s.On("DescribeLaunchTemplateVersions", &ec2.DescribeLaunchTemplateVersionsInput{LaunchTemplateName: aws.String(ltName), Versions: []*string{aws.String(ltVersion)}}).Return(&ec2.DescribeLaunchTemplateVersionsOutput{LaunchTemplateVersions: []*ec2.LaunchTemplateVersion{{LaunchTemplateData: &ec2.ResponseLaunchTemplateData{InstanceType: aws.String(test.instanceType)}}}})
			m, err := createAWSManagerInternal(nil, cloudprovider.NodeGroupDiscoveryOptions{}, nil, &ec2Wrapper{s})
			assert.NoError(t, err)
			asg := &asg{AwsRef: AwsRef{Name: "sample"}, AvailabilityZones: test.availabilityZones, LaunchTemplateName: ltName, LaunchTemplateVersion: ltVersion, Tags: tags}
			template, err := m.getAsgTemplate(asg)
			if test.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if assert.NotNil(t, template) {
					assert.Equal(t, test.instanceType, template.InstanceType.InstanceType)
					assert.Equal(t, region, template.Region)
					assert.Equal(t, test.availabilityZones[0], template.Zone)
					assert.Equal(t, tags, template.Tags)
				}
			}
		})
	}
}
func TestFetchAutoAsgs(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	min, max := 1, 10
	groupname, tags := "coolasg", []string{"tag", "anothertag"}
	s := &AutoScalingMock{}
	expectedTagsInput := &autoscaling.DescribeTagsInput{Filters: []*autoscaling.Filter{{Name: aws.String("key"), Values: aws.StringSlice([]string{tags[0]})}, {Name: aws.String("key"), Values: aws.StringSlice([]string{tags[1]})}}, MaxRecords: aws.Int64(maxRecordsReturnedByAPI)}
	s.On("DescribeTagsPages", mock.MatchedBy(tagsMatcher(expectedTagsInput)), mock.AnythingOfType("func(*autoscaling.DescribeTagsOutput, bool) bool")).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(*autoscaling.DescribeTagsOutput, bool) bool)
		fn(&autoscaling.DescribeTagsOutput{Tags: []*autoscaling.TagDescription{{ResourceId: aws.String(groupname)}, {ResourceId: aws.String(groupname)}}}, false)
	}).Return(nil).Once()
	s.On("DescribeAutoScalingGroupsPages", &autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: aws.StringSlice([]string{groupname}), MaxRecords: aws.Int64(maxRecordsReturnedByAPI)}, mock.AnythingOfType("func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool")).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool)
		fn(&autoscaling.DescribeAutoScalingGroupsOutput{AutoScalingGroups: []*autoscaling.Group{{AutoScalingGroupName: aws.String(groupname), MinSize: aws.Int64(int64(min)), MaxSize: aws.Int64(int64(max))}}}, false)
	}).Return(nil).Twice()
	do := cloudprovider.NodeGroupDiscoveryOptions{NodeGroupAutoDiscoverySpecs: []string{fmt.Sprintf("asg:tag=%s", strings.Join(tags, ","))}}
	m, err := createAWSManagerInternal(nil, do, &autoScalingWrapper{s}, nil)
	assert.NoError(t, err)
	asgs := m.asgCache.Get()
	assert.Equal(t, 1, len(asgs))
	validateAsg(t, asgs[0], groupname, min, max)
	s.On("DescribeTagsPages", mock.MatchedBy(tagsMatcher(expectedTagsInput)), mock.AnythingOfType("func(*autoscaling.DescribeTagsOutput, bool) bool")).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(*autoscaling.DescribeTagsOutput, bool) bool)
		fn(&autoscaling.DescribeTagsOutput{Tags: []*autoscaling.TagDescription{}}, false)
	}).Return(nil).Once()
	err = m.asgCache.regenerate()
	assert.NoError(t, err)
	assert.Empty(t, m.asgCache.Get())
}
func tagsMatcher(expected *autoscaling.DescribeTagsInput) func(*autoscaling.DescribeTagsInput) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(actual *autoscaling.DescribeTagsInput) bool {
		expectedTags := flatTagSlice(expected.Filters)
		actualTags := flatTagSlice(actual.Filters)
		return *expected.MaxRecords == *actual.MaxRecords && reflect.DeepEqual(expectedTags, actualTags)
	}
}
func flatTagSlice(filters []*autoscaling.Filter) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	tags := []string{}
	for _, filter := range filters {
		tags = append(tags, aws.StringValueSlice(filter.Values)...)
	}
	sort.Strings(tags)
	return tags
}
