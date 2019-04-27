package aws

import (
	"fmt"
	"testing"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMoreThen50Groups(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	service := &AutoScalingMock{}
	autoScalingWrapper := &autoScalingWrapper{autoScaling: service}
	names := make([]string, 51)
	for i := 0; i < len(names); i++ {
		names[i] = fmt.Sprintf("asg-%d", i)
	}
	service.On("DescribeAutoScalingGroupsPages", &autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: aws.StringSlice(names[:50]), MaxRecords: aws.Int64(maxRecordsReturnedByAPI)}, mock.AnythingOfType("func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool")).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool)
		fn(testNamedDescribeAutoScalingGroupsOutput("asg-1", 1, "test-instance-id"), false)
	}).Return(nil)
	service.On("DescribeAutoScalingGroupsPages", &autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: aws.StringSlice([]string{"asg-50"}), MaxRecords: aws.Int64(maxRecordsReturnedByAPI)}, mock.AnythingOfType("func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool")).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool)
		fn(testNamedDescribeAutoScalingGroupsOutput("asg-2", 1, "test-instance-id"), false)
	}).Return(nil)
	asgs, err := autoScalingWrapper.getAutoscalingGroupsByNames(names)
	assert.Nil(t, err)
	assert.Equal(t, len(asgs), 2)
	assert.Equal(t, *asgs[0].AutoScalingGroupName, "asg-1")
	assert.Equal(t, *asgs[1].AutoScalingGroupName, "asg-2")
}
