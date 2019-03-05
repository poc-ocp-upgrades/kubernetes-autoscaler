/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package aws

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	provider_aws "k8s.io/kubernetes/pkg/cloudprovider/providers/aws"
	kubeletapis "k8s.io/kubernetes/pkg/kubelet/apis"
)

func TestBuildGenericLabels(t *testing.T) {
	labels := buildGenericLabels(&asgTemplate{
		InstanceType: &instanceType{
			InstanceType: "c4.large",
			VCPU:         2,
			MemoryMb:     3840,
		},
		Region: "us-east-1",
	}, "sillyname")
	assert.Equal(t, "us-east-1", labels[kubeletapis.LabelZoneRegion])
	assert.Equal(t, "sillyname", labels[kubeletapis.LabelHostname])
	assert.Equal(t, "c4.large", labels[kubeletapis.LabelInstanceType])
	assert.Equal(t, cloudprovider.DefaultArch, labels[kubeletapis.LabelArch])
	assert.Equal(t, cloudprovider.DefaultOS, labels[kubeletapis.LabelOS])
}

func TestExtractLabelsFromAsg(t *testing.T) {
	tags := []*autoscaling.TagDescription{
		{
			Key:   aws.String("k8s.io/cluster-autoscaler/node-template/label/foo"),
			Value: aws.String("bar"),
		},
		{
			Key:   aws.String("bar"),
			Value: aws.String("baz"),
		},
	}

	labels := extractLabelsFromAsg(tags)

	assert.Equal(t, 1, len(labels))
	assert.Equal(t, "bar", labels["foo"])
}

func TestExtractTaintsFromAsg(t *testing.T) {
	tags := []*autoscaling.TagDescription{
		{
			Key:   aws.String("k8s.io/cluster-autoscaler/node-template/taint/dedicated"),
			Value: aws.String("foo:NoSchedule"),
		},
		{
			Key:   aws.String("bar"),
			Value: aws.String("baz"),
		},
	}

	expectedTaints := []apiv1.Taint{
		{
			Key:    "dedicated",
			Value:  "foo",
			Effect: apiv1.TaintEffectNoSchedule,
		},
	}

	taints := extractTaintsFromAsg(tags)
	assert.Equal(t, 1, len(taints))
	assert.Equal(t, makeTaintSet(expectedTaints), makeTaintSet(taints))
}

func makeTaintSet(taints []apiv1.Taint) map[apiv1.Taint]bool {
	set := make(map[apiv1.Taint]bool)
	for _, taint := range taints {
		set[taint] = true
	}
	return set
}
func TestBuildAsg(t *testing.T) {
	do := cloudprovider.NodeGroupDiscoveryOptions{}
	m, err := createAWSManagerInternal(nil, do, &testService)
	assert.NoError(t, err)

	asg, err := m.buildAsgFromSpec("1:5:test-asg")
	assert.NoError(t, err)
	assert.Equal(t, asg.MinSize(), 1)
	assert.Equal(t, asg.MaxSize(), 5)
	assert.Equal(t, asg.Id(), "test-asg")
	assert.Equal(t, asg.Name, "test-asg")
	assert.Equal(t, asg.Debug(), "test-asg (1:5)")

	_, err = m.buildAsgFromSpec("a")
	assert.Error(t, err)
	_, err = m.buildAsgFromSpec("a:b:c")
	assert.Error(t, err)
	_, err = m.buildAsgFromSpec("1:")
	assert.Error(t, err)
	_, err = m.buildAsgFromSpec("1:2:")
	assert.Error(t, err)
}

func validateAsg(t *testing.T, asg *Asg, name string, minSize int, maxSize int) {
	assert.Equal(t, name, asg.Name)
	assert.Equal(t, minSize, asg.minSize)
	assert.Equal(t, maxSize, asg.maxSize)
}

func TestFetchExplicitAsgs(t *testing.T) {
	min, max, groupname := 1, 10, "coolasg"

	s := &AutoScalingMock{}
	s.On("DescribeAutoScalingGroups", &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{aws.String(groupname)},
		MaxRecords:            aws.Int64(1),
	}).Return(&autoscaling.DescribeAutoScalingGroupsOutput{
		AutoScalingGroups: []*autoscaling.Group{
			{AutoScalingGroupName: aws.String(groupname)},
		},
	})

	s.On("DescribeAutoScalingGroupsPages",
		&autoscaling.DescribeAutoScalingGroupsInput{
			AutoScalingGroupNames: aws.StringSlice([]string{groupname}),
			MaxRecords:            aws.Int64(maxRecordsReturnedByAPI),
		},
		mock.AnythingOfType("func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool"),
	).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool)
		fn(&autoscaling.DescribeAutoScalingGroupsOutput{
			AutoScalingGroups: []*autoscaling.Group{
				{AutoScalingGroupName: aws.String(groupname)},
			}}, false)
	}).Return(nil)

	do := cloudprovider.NodeGroupDiscoveryOptions{
		// Register the same node group twice with different max nodes.
		// The intention is to test that the asgs.Register method will update
		// the node group instead of registering it twice.
		NodeGroupSpecs: []string{
			fmt.Sprintf("%d:%d:%s", min, max-1, groupname),
			fmt.Sprintf("%d:%d:%s", min, max, groupname),
		},
	}
	// fetchExplicitASGs is called at manager creation time.
	m, err := createAWSManagerInternal(nil, do, &autoScalingWrapper{s})
	assert.NoError(t, err)

	asgs := m.asgCache.get()
	assert.Equal(t, 1, len(asgs))
	validateAsg(t, asgs[0].config, groupname, min, max)
}

func TestFetchAutoAsgs(t *testing.T) {
	min, max := 1, 10
	groupname, tags := "coolasg", []string{"tag", "anothertag"}

	s := &AutoScalingMock{}
	// Lookup groups associated with tags
	s.On("DescribeTagsPages",
		&autoscaling.DescribeTagsInput{
			Filters: []*autoscaling.Filter{
				{Name: aws.String("key"), Values: aws.StringSlice([]string{tags[0]})},
				{Name: aws.String("key"), Values: aws.StringSlice([]string{tags[1]})},
			},
			MaxRecords: aws.Int64(maxRecordsReturnedByAPI),
		},
		mock.AnythingOfType("func(*autoscaling.DescribeTagsOutput, bool) bool"),
	).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(*autoscaling.DescribeTagsOutput, bool) bool)
		fn(&autoscaling.DescribeTagsOutput{
			Tags: []*autoscaling.TagDescription{
				{ResourceId: aws.String(groupname)},
				{ResourceId: aws.String(groupname)},
			}}, false)
	}).Return(nil).Once()

	// Describe the group to register it, then again to generate the instance
	// cache.
	s.On("DescribeAutoScalingGroupsPages",
		&autoscaling.DescribeAutoScalingGroupsInput{
			AutoScalingGroupNames: aws.StringSlice([]string{groupname}),
			MaxRecords:            aws.Int64(maxRecordsReturnedByAPI),
		},
		mock.AnythingOfType("func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool"),
	).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool)
		fn(&autoscaling.DescribeAutoScalingGroupsOutput{
			AutoScalingGroups: []*autoscaling.Group{{
				AutoScalingGroupName: aws.String(groupname),
				MinSize:              aws.Int64(int64(min)),
				MaxSize:              aws.Int64(int64(max)),
			}}}, false)
	}).Return(nil).Twice()

	do := cloudprovider.NodeGroupDiscoveryOptions{
		NodeGroupAutoDiscoverySpecs: []string{fmt.Sprintf("asg:tag=%s", strings.Join(tags, ","))},
	}

	// fetchAutoASGs is called at manager creation time, via forceRefresh
	m, err := createAWSManagerInternal(nil, do, &autoScalingWrapper{s})
	assert.NoError(t, err)

	asgs := m.asgCache.get()
	assert.Equal(t, 1, len(asgs))
	validateAsg(t, asgs[0].config, groupname, min, max)

	// Simulate the previously discovered ASG disappearing
	s.On("DescribeTagsPages",
		&autoscaling.DescribeTagsInput{
			Filters: []*autoscaling.Filter{
				{Name: aws.String("key"), Values: aws.StringSlice([]string{tags[0]})},
				{Name: aws.String("key"), Values: aws.StringSlice([]string{tags[1]})},
			},
			MaxRecords: aws.Int64(maxRecordsReturnedByAPI),
		},
		mock.AnythingOfType("func(*autoscaling.DescribeTagsOutput, bool) bool"),
	).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(*autoscaling.DescribeTagsOutput, bool) bool)
		fn(&autoscaling.DescribeTagsOutput{Tags: []*autoscaling.TagDescription{}}, false)
	}).Return(nil).Once()

	err = m.fetchAutoAsgs()
	assert.NoError(t, err)
	assert.Empty(t, m.asgCache.get())
}

type ServiceDescriptor struct {
	name                         string
	region                       string
	signingRegion, signingMethod string
	signingName                  string
}

func TestOverridesActiveConfig(t *testing.T) {
	tests := []struct {
		name string

		reader io.Reader
		aws    provider_aws.Services

		expectError        bool
		active             bool
		servicesOverridden []ServiceDescriptor
	}{
		{
			"No overrides",
			strings.NewReader(`
				[global]
				`),
			nil,
			false, false,
			[]ServiceDescriptor{},
		},
		{
			"Missing Service Name",
			strings.NewReader(`
				[global]
				[ServiceOverride "1"]
				Region=sregion
				URL=https://s3.foo.bar
				SigningRegion=sregion
				SigningMethod = sign
				`),
			nil,
			true, false,
			[]ServiceDescriptor{},
		},
		{
			"Missing Service Region",
			strings.NewReader(`
				[global]
				[ServiceOverride "1"]
				Service=s3
				URL=https://s3.foo.bar
				SigningRegion=sregion
				SigningMethod = sign
				`),
			nil,
			true, false,
			[]ServiceDescriptor{},
		},
		{
			"Missing URL",
			strings.NewReader(`
				[global]
				[ServiceOverride "1"]
				Service="s3"
				Region=sregion
				SigningRegion=sregion
				SigningMethod = sign
				`),
			nil,
			true, false,
			[]ServiceDescriptor{},
		},
		{
			"Missing Signing Region",
			strings.NewReader(`
				[global]
				[ServiceOverride "1"]
				Service=s3
				Region=sregion
				URL=https://s3.foo.bar
				SigningMethod = sign
				`),
			nil,
			true, false,
			[]ServiceDescriptor{},
		},
		{
			"Active Overrides",
			strings.NewReader(`
				[Global]
				[ServiceOverride "1"]
				Service = "s3      "
				Region = sregion
				URL = https://s3.foo.bar
				SigningRegion = sregion
				SigningMethod = v4
				`),
			nil,
			false, true,
			[]ServiceDescriptor{{name: "s3", region: "sregion", signingRegion: "sregion", signingMethod: "v4"}},
		},
		{
			"Multiple Overridden Services",
			strings.NewReader(`
				[Global]
				vpc = vpc-abc1234567
				[ServiceOverride "1"]
				Service=s3
				Region=sregion1
				URL=https://s3.foo.bar
				SigningRegion=sregion1
				SigningMethod = v4
				[ServiceOverride "2"]
				Service=ec2
				Region=sregion2
				URL=https://ec2.foo.bar
				SigningRegion=sregion2
				SigningMethod = v4
				`),
			nil,
			false, true,
			[]ServiceDescriptor{{name: "s3", region: "sregion1", signingRegion: "sregion1", signingMethod: "v4"},
				{name: "ec2", region: "sregion2", signingRegion: "sregion2", signingMethod: "v4"}},
		},
		{
			"Duplicate Services",
			strings.NewReader(`
				[Global]
				vpc = vpc-abc1234567
				[ServiceOverride "1"]
				Service=s3
				Region=sregion1
				URL=https://s3.foo.bar
				SigningRegion=sregion
				SigningMethod = sign
				[ServiceOverride "2"]
				Service=s3
				Region=sregion1
				URL=https://s3.foo.bar
				SigningRegion=sregion
				SigningMethod = sign
				`),
			nil,
			true, false,
			[]ServiceDescriptor{},
		},
		{
			"Multiple Overridden Services in Multiple regions",
			strings.NewReader(`
				[global]
				[ServiceOverride "1"]
			 	Service=s3
				Region=region1
				URL=https://s3.foo.bar
				SigningRegion=sregion1
				[ServiceOverride "2"]
				Service=ec2
				Region=region2
				URL=https://ec2.foo.bar
				SigningRegion=sregion
				SigningMethod = v4
				`),
			nil,
			false, true,
			[]ServiceDescriptor{{name: "s3", region: "region1", signingRegion: "sregion1", signingMethod: ""},
				{name: "ec2", region: "region2", signingRegion: "sregion", signingMethod: "v4"}},
		},
		{
			"Multiple regions, Same Service",
			strings.NewReader(`
				[global]
				[ServiceOverride "1"]
				Service=s3
				Region=region1
				URL=https://s3.foo.bar
				SigningRegion=sregion1
				SigningMethod = v3
				[ServiceOverride "2"]
				Service=s3
				Region=region2
				URL=https://s3.foo.bar
				SigningRegion=sregion1
				SigningMethod = v4
				SigningName = "name"
				`),
			nil,
			false, true,
			[]ServiceDescriptor{{name: "s3", region: "region1", signingRegion: "sregion1", signingMethod: "v3"},
				{name: "s3", region: "region2", signingRegion: "sregion1", signingMethod: "v4", signingName: "name"}},
		},
	}

	for _, test := range tests {
		t.Logf("Running test case %s", test.name)
		cfg, err := readAWSCloudConfig(test.reader)
		if err == nil {
			err = validateOverrides(cfg)
		}
		if test.expectError {
			assert.Error(t, err)
			continue
		}
		if len(cfg.ServiceOverride) != len(test.servicesOverridden) {
			t.Errorf("Expected %d overridden services, received %d for case %s",
				len(test.servicesOverridden), len(cfg.ServiceOverride), test.name)
			continue
		}
		for _, sd := range test.servicesOverridden {
			var found *struct {
				Service       string
				Region        string
				URL           string
				SigningRegion string
				SigningMethod string
				SigningName   string
			}
			for _, v := range cfg.ServiceOverride {
				if v.Service == sd.name && v.Region == sd.region {
					found = v
					break
				}
			}
			if found == nil {
				t.Errorf("Missing override for service %s in case %s",
					sd.name, test.name)
			} else {
				if found.SigningRegion != sd.signingRegion {
					t.Errorf("Expected signing region '%s', received '%s' for case %s",
						sd.signingRegion, found.SigningRegion, test.name)
				}
				if found.SigningMethod != sd.signingMethod {
					t.Errorf("Expected signing method '%s', received '%s' for case %s",
						sd.signingMethod, found.SigningRegion, test.name)
				}
				targetName := fmt.Sprintf("https://%s.foo.bar", sd.name)
				if found.URL != targetName {
					t.Errorf("Expected Endpoint '%s', received '%s' for case %s",
						targetName, found.URL, test.name)
				}
				if found.SigningName != sd.signingName {
					t.Errorf("Expected signing name '%s', received '%s' for case %s",
						sd.signingName, found.SigningName, test.name)
				}
				fn := getResolver(cfg)
				ep1, e := fn(sd.name, sd.region, nil)
				assert.Nil(t, e)
				assert.NotNil(t, ep1)
				assert.Equal(t, targetName, ep1.URL)
				assert.Equal(t, sd.signingRegion, ep1.SigningRegion)
				assert.Equal(t, sd.signingMethod, ep1.SigningMethod)
			}
		}
	}
}
