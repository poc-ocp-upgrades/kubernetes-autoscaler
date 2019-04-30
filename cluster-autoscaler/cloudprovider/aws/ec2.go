package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type ec2I interface {
	DescribeLaunchTemplateVersions(input *ec2.DescribeLaunchTemplateVersionsInput) (*ec2.DescribeLaunchTemplateVersionsOutput, error)
}
type ec2Wrapper struct{ ec2I }

func (m ec2Wrapper) getInstanceTypeByLT(name string, version string) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	params := &ec2.DescribeLaunchTemplateVersionsInput{LaunchTemplateName: aws.String(name), Versions: []*string{aws.String(version)}}
	describeData, err := m.DescribeLaunchTemplateVersions(params)
	if err != nil {
		return "", err
	}
	if len(describeData.LaunchTemplateVersions) == 0 {
		return "", fmt.Errorf("Unable to find template versions")
	}
	lt := describeData.LaunchTemplateVersions[0]
	instanceType := lt.LaunchTemplateData.InstanceType
	if instanceType == nil {
		return "", fmt.Errorf("Unable to find instance type within launch template")
	}
	return aws.StringValue(instanceType), nil
}
